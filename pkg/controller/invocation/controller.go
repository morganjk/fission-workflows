package invocation

import (
	"context"
	"errors"
	"time"

	"github.com/fission/fission-workflows/pkg/api"
	"github.com/fission/fission-workflows/pkg/controller"
	"github.com/fission/fission-workflows/pkg/controller/expr"
	"github.com/fission/fission-workflows/pkg/fes"
	"github.com/fission/fission-workflows/pkg/scheduler"
	"github.com/fission/fission-workflows/pkg/types/aggregates"
	"github.com/fission/fission-workflows/pkg/types/events"
	"github.com/fission/fission-workflows/pkg/util/labels"
	"github.com/fission/fission-workflows/pkg/util/pubsub"
	log "github.com/sirupsen/logrus"
)

const (
	NotificationBuffer = 100
	evalQueueSize      = 50
)

var wfiLog = log.WithField("component", "controller-wi")

type Controller struct {
	invokeCache   fes.CacheReader
	wfCache       fes.CacheReader
	taskAPI       *api.Task
	invocationAPI *api.Invocation
	stateStore    *expr.Store
	scheduler     *scheduler.WorkflowScheduler
	sub           *pubsub.Subscription
	cancelFn      context.CancelFunc
	evalPolicy    controller.Rule
	evalCache     *controller.EvalCache

	// evalQueue is a queue of invocation ids
	evalQueue chan string
}

func NewController(invokeCache fes.CacheReader, wfCache fes.CacheReader, workflowScheduler *scheduler.WorkflowScheduler,
	taskAPI *api.Task, invocationAPI *api.Invocation, stateStore *expr.Store) *Controller {
	ctr := &Controller{
		invokeCache:   invokeCache,
		wfCache:       wfCache,
		scheduler:     workflowScheduler,
		taskAPI:       taskAPI,
		invocationAPI: invocationAPI,
		evalQueue:     make(chan string, evalQueueSize),
		evalCache:     controller.NewEvalCache(),
		stateStore:    stateStore,

		// States maintains an active cache of currently running invocations, with execution related data.
		// This state information is considered preemptable and can be removed or lost at any time.
		//states: map[string]*ControlState{},
	}
	ctr.evalPolicy = defaultPolicy(ctr)

	return ctr
}

func (cr *Controller) Init(sctx context.Context) error {
	ctx, cancelFn := context.WithCancel(sctx)
	cr.cancelFn = cancelFn

	// Subscribe to invocation creations and task events.
	selector := labels.In(fes.PubSubLabelAggregateType, "invocation", "function")
	if invokePub, ok := cr.invokeCache.(pubsub.Publisher); ok {
		cr.sub = invokePub.Subscribe(pubsub.SubscriptionOptions{
			Buffer:   NotificationBuffer,
			Selector: selector,
		})

		// Invocation Notification listener
		go func(ctx context.Context) {
			for {
				select {
				case notification := <-cr.sub.Ch:
					cr.handleMsg(notification)
				case <-ctx.Done():
					wfiLog.WithField("ctx.err", ctx.Err()).Debug("Notification listener closed.")
					return
				}
			}
		}(ctx)
	}

	// process evaluation queue
	go func(ctx context.Context) {
		for {
			select {
			case eval := <-cr.evalQueue:
				go cr.Evaluate(eval) // TODO limit number of goroutines
			case <-ctx.Done():
				return
			}
		}
	}(ctx)

	return nil
}

func (cr *Controller) handleMsg(msg pubsub.Msg) error {
	wfiLog.WithField("labels", msg.Labels()).Debug("Handling invocation notification.")
	switch n := msg.(type) {
	case *fes.Notification:
		cr.Notify(n)
	default:
		wfiLog.WithField("notification", n).Warn("Ignoring unknown notification type")
	}
	return nil
}

func (cr *Controller) Notify(msg *fes.Notification) error {
	wfiLog.WithFields(log.Fields{
		"notification": msg.EventType,
		"labels":       msg.Labels(),
	}).Debug("Handling invocation notification!")

	switch msg.EventType {
	case events.Invocation_INVOCATION_COMPLETED.String():
		fallthrough
	case events.Invocation_INVOCATION_CANCELED.String():
		fallthrough
	case events.Invocation_INVOCATION_FAILED.String():
		wfi, ok := msg.Payload.(*aggregates.WorkflowInvocation)
		if !ok {
			log.Warn("Event did not contain invocation payload", msg)
		}
		// TODO mark to clean up later instead
		cr.stateStore.Delete(wfi.ID())
		cr.evalCache.Del(wfi.ID())
		log.Infof("Removed invocation %v from eval state", wfi.ID())
	case events.Task_TASK_FAILED.String():
		fallthrough
	case events.Task_TASK_SUCCEEDED.String():
		fallthrough
	case events.Invocation_INVOCATION_CREATED.String():
		wfi, ok := msg.Payload.(*aggregates.WorkflowInvocation)
		if !ok {
			panic(msg)
		}
		cr.submitEval(wfi.ID())
	default:
		wfiLog.Debugf("Controller ignored event type: %v", msg.EventType)
	}
	return nil
}

func (cr *Controller) Tick(tick uint64) error {
	// Short loop: invocations the controller is actively tracking
	var err error
	if tick%2 != 0 {
		err = cr.checkEvalCaches()
	}

	// Long loop: to check if there are any orphans
	if tick%4 != 0 {
		err = cr.checkModelCaches()
	}

	return err
}

func (cr *Controller) checkEvalCaches() error {
	for id, state := range cr.evalCache.List() {
		last, ok := state.Last()
		if !ok {
			continue
		}

		reevaluateAt := last.Timestamp.Add(time.Duration(100) * time.Millisecond)
		if time.Now().UnixNano() > reevaluateAt.UnixNano() {
			cr.submitEval(id)
		}
	}
	return nil
}

// checkCaches iterates over the current caches submitting evaluation for invocation when needed
func (cr *Controller) checkModelCaches() error {
	// Short control loop
	entities := cr.invokeCache.List()
	for _, entity := range entities {
		if _, ok := cr.evalCache.Get(entity.Id); ok {
			continue
		}

		wi := aggregates.NewWorkflowInvocation(entity.Id)
		err := cr.invokeCache.Get(wi)
		if err != nil {
			log.Errorf("Failed to read '%v' from cache: %v.", wi.Aggregate(), err)
			continue
		}

		if !wi.Status.Finished() {
			cr.submitEval(wi.ID())
		}
	}
	return nil
}

func (cr *Controller) submitEval(ids ...string) bool {
	for _, id := range ids {
		select {
		case cr.evalQueue <- id:
			return true
			// ok
		default:
			wfiLog.Warnf("Eval queue is full; dropping eval task for '%v'", id)
			return false
		}
	}
	return true
}

func (cr *Controller) Evaluate(invocationID string) {
	// Fetch and attempt to claim the evaluation
	evalState := cr.evalCache.GetOrCreate(invocationID)
	select {
	case <-evalState.Lock():
		defer evalState.Free()
	default:
		// TODO provide option to wait for a lock
		wfiLog.Debugf("Failed to obtain access to invocation %s", invocationID)
		return
	}
	log.Debugf("evaluating invocation %s", invocationID)

	// Fetch the workflow invocation for the provided invocation id
	wfi := aggregates.NewWorkflowInvocation(invocationID)
	err := cr.invokeCache.Get(wfi)
	// TODO move to rule
	if err != nil && wfi.WorkflowInvocation == nil {
		log.Errorf("controller failed to get invocation for invocation id '%s': %v", invocationID, err)
		return
	}
	// TODO move to rule
	if wfi.Status.Finished() {
		wfiLog.Debugf("No need to evaluate finished invocation %v", invocationID)
		return
	}

	// Fetch the workflow relevant to the invocation
	wf := aggregates.NewWorkflow(wfi.Spec.WorkflowId)
	err = cr.wfCache.Get(wf)
	// TODO move to rule
	if err != nil && wf.Workflow == nil {
		log.Errorf("controller failed to get workflow '%s' for invocation id '%s': %v", wfi.Spec.WorkflowId,
			invocationID, err)
		return
	}

	// Evaluate invocation
	record := controller.NewEvalRecord() // TODO implement rulepath + cause

	ec := NewEvalContext(evalState, wf.Workflow, wfi.WorkflowInvocation)

	action := cr.evalPolicy.Eval(ec)
	record.Action = action
	if action == nil {
		return
	}

	// Execute action
	err = action.Apply()
	if err != nil {
		log.Errorf("Action '%T' failed: %v", action, err)
		record.Error = err
	}

	// Record this evaluation
	evalState.Record(record)
}

func (cr *Controller) Close() error {
	wfiLog.Info("Closing invocation controller...")
	if invokePub, ok := cr.invokeCache.(pubsub.Publisher); ok {
		err := invokePub.Unsubscribe(cr.sub)
		if err != nil {
			return err
		}
	}

	cr.cancelFn()
	return nil
}

func (cr *Controller) createFailAction(invocationID string, err error) controller.Action {
	return &ActionFail{
		API:          cr.invocationAPI,
		InvocationID: invocationID,
		Err:          err,
	}
}

func defaultPolicy(ctr *Controller) controller.Rule {
	return &controller.RuleEvalUntilAction{
		Rules: []controller.Rule{
			&controller.RuleTimedOut{
				OnTimedOut: &ActionFail{
					API: ctr.invocationAPI,
					Err: errors.New("timed out"),
				},
				Timeout: time.Duration(10) * time.Minute,
			},
			&controller.RuleExceededErrorCount{
				OnExceeded: &ActionFail{
					API: ctr.invocationAPI,
				},
				MaxErrorCount: 0,
			},
			&RuleHasCompleted{},
			&RuleCheckIfCompleted{
				InvocationAPI: ctr.invocationAPI,
			},
			&RuleWorkflowIsReady{},
			&RuleSchedule{
				Scheduler:     ctr.scheduler,
				InvocationAPI: ctr.invocationAPI,
				FunctionAPI:   ctr.taskAPI,
				StateStore:    ctr.stateStore,
			},
		},
	}
}

// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/validate"
)

// WorkflowInvocationStatusTasks workflow invocation status tasks
// swagger:model workflowInvocationStatusTasks
type WorkflowInvocationStatusTasks map[string]TaskInvocation

// Validate validates this workflow invocation status tasks
func (m WorkflowInvocationStatusTasks) Validate(formats strfmt.Registry) error {
	var res []error

	if err := validate.Required("", "body", WorkflowInvocationStatusTasks(m)); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

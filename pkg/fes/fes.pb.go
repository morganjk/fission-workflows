// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pkg/fes/fes.proto

/*
Package fes is a generated protocol buffer package.

It is generated from these files:
	pkg/fes/fes.proto

It has these top-level messages:
	Aggregate
	Event
	EventHints
	DummyEvent
*/
package fes

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/timestamp"
import google_protobuf1 "github.com/golang/protobuf/ptypes/any"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Aggregate struct {
	Id   string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Type string `protobuf:"bytes,2,opt,name=type" json:"type,omitempty"`
}

func (m *Aggregate) Reset()                    { *m = Aggregate{} }
func (m *Aggregate) String() string            { return proto.CompactTextString(m) }
func (*Aggregate) ProtoMessage()               {}
func (*Aggregate) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Aggregate) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Aggregate) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

type Event struct {
	Id        string                     `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Type      string                     `protobuf:"bytes,2,opt,name=type" json:"type,omitempty"`
	Aggregate *Aggregate                 `protobuf:"bytes,3,opt,name=aggregate" json:"aggregate,omitempty"`
	Timestamp *google_protobuf.Timestamp `protobuf:"bytes,4,opt,name=timestamp" json:"timestamp,omitempty"`
	Data      *google_protobuf1.Any      `protobuf:"bytes,5,opt,name=data" json:"data,omitempty"`
	Parent    *Aggregate                 `protobuf:"bytes,6,opt,name=parent" json:"parent,omitempty"`
	Hints     *EventHints                `protobuf:"bytes,7,opt,name=hints" json:"hints,omitempty"`
	Metadata  map[string]string          `protobuf:"bytes,8,rep,name=metadata" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

func (m *Event) Reset()                    { *m = Event{} }
func (m *Event) String() string            { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()               {}
func (*Event) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Event) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Event) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *Event) GetAggregate() *Aggregate {
	if m != nil {
		return m.Aggregate
	}
	return nil
}

func (m *Event) GetTimestamp() *google_protobuf.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *Event) GetData() *google_protobuf1.Any {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *Event) GetParent() *Aggregate {
	if m != nil {
		return m.Parent
	}
	return nil
}

func (m *Event) GetHints() *EventHints {
	if m != nil {
		return m.Hints
	}
	return nil
}

func (m *Event) GetMetadata() map[string]string {
	if m != nil {
		return m.Metadata
	}
	return nil
}

// EventHints is a collection of optional metadata that help components in the event store to improve performance.
type EventHints struct {
	Completed bool `protobuf:"varint,1,opt,name=completed" json:"completed,omitempty"`
}

func (m *EventHints) Reset()                    { *m = EventHints{} }
func (m *EventHints) String() string            { return proto.CompactTextString(m) }
func (*EventHints) ProtoMessage()               {}
func (*EventHints) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *EventHints) GetCompleted() bool {
	if m != nil {
		return m.Completed
	}
	return false
}

type DummyEvent struct {
	Msg string `protobuf:"bytes,1,opt,name=msg" json:"msg,omitempty"`
}

func (m *DummyEvent) Reset()                    { *m = DummyEvent{} }
func (m *DummyEvent) String() string            { return proto.CompactTextString(m) }
func (*DummyEvent) ProtoMessage()               {}
func (*DummyEvent) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *DummyEvent) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

func init() {
	proto.RegisterType((*Aggregate)(nil), "fission.workflows.eventstore.Aggregate")
	proto.RegisterType((*Event)(nil), "fission.workflows.eventstore.Event")
	proto.RegisterType((*EventHints)(nil), "fission.workflows.eventstore.EventHints")
	proto.RegisterType((*DummyEvent)(nil), "fission.workflows.eventstore.DummyEvent")
}

func init() { proto.RegisterFile("pkg/fes/fes.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 373 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x92, 0xc1, 0x4b, 0xf3, 0x30,
	0x18, 0xc6, 0x69, 0xbb, 0xee, 0x5b, 0xdf, 0xf1, 0x7d, 0xec, 0x0b, 0x3b, 0xd4, 0x32, 0x74, 0xf4,
	0x62, 0xf1, 0x90, 0xe2, 0xbc, 0x0c, 0x05, 0x65, 0xe2, 0xc0, 0xcb, 0x2e, 0xc5, 0x93, 0xb7, 0xcc,
	0xa5, 0xb1, 0xac, 0x69, 0x4a, 0x93, 0x6d, 0xf4, 0xef, 0xf5, 0x1f, 0x91, 0xa6, 0xeb, 0x8a, 0x0a,
	0x63, 0x1e, 0x0a, 0x69, 0xf2, 0xfc, 0xde, 0x3c, 0xcf, 0xfb, 0x06, 0xfe, 0xe7, 0x6b, 0x16, 0xc6,
	0x54, 0x56, 0x1f, 0xce, 0x0b, 0xa1, 0x04, 0x1a, 0xc5, 0x89, 0x94, 0x89, 0xc8, 0xf0, 0x4e, 0x14,
	0xeb, 0x38, 0x15, 0x3b, 0x89, 0xe9, 0x96, 0x66, 0x4a, 0x2a, 0x51, 0x50, 0xef, 0x82, 0x09, 0xc1,
	0x52, 0x1a, 0x6a, 0xed, 0x72, 0x13, 0x87, 0x2a, 0xe1, 0x54, 0x2a, 0xc2, 0xf3, 0x1a, 0xf7, 0xce,
	0xbe, 0x0b, 0x48, 0x56, 0xd6, 0x47, 0x7e, 0x08, 0xce, 0x8c, 0xb1, 0x82, 0x32, 0xa2, 0x28, 0xfa,
	0x07, 0x66, 0xb2, 0x72, 0x8d, 0xb1, 0x11, 0x38, 0x91, 0x99, 0xac, 0x10, 0x82, 0x8e, 0x2a, 0x73,
	0xea, 0x9a, 0x7a, 0x47, 0xaf, 0xfd, 0x0f, 0x0b, 0xec, 0x79, 0x75, 0xf7, 0x29, 0x6a, 0x34, 0x07,
	0x87, 0x34, 0xe5, 0x5d, 0x6b, 0x6c, 0x04, 0xfd, 0xc9, 0x25, 0x3e, 0x16, 0x06, 0x1f, 0xdc, 0x44,
	0x2d, 0x89, 0xa6, 0xe0, 0x1c, 0x32, 0xb9, 0x1d, 0x5d, 0xc6, 0xc3, 0x75, 0x28, 0xdc, 0x84, 0xc2,
	0x2f, 0x8d, 0x22, 0x6a, 0xc5, 0x28, 0x80, 0xce, 0x8a, 0x28, 0xe2, 0xda, 0x1a, 0x1a, 0xfe, 0x80,
	0x66, 0x59, 0x19, 0x69, 0x05, 0x7a, 0x80, 0x6e, 0x4e, 0x0a, 0x9a, 0x29, 0xb7, 0xfb, 0x3b, 0x9f,
	0x7b, 0x0c, 0xdd, 0x83, 0xfd, 0x9e, 0x64, 0x4a, 0xba, 0x7f, 0x34, 0x1f, 0x1c, 0xe7, 0x75, 0x0f,
	0x9f, 0x2b, 0x7d, 0x54, 0x63, 0x68, 0x01, 0x3d, 0x4e, 0x15, 0xd1, 0x76, 0x7b, 0x63, 0x2b, 0xe8,
	0x4f, 0xae, 0x4f, 0x28, 0x81, 0x17, 0x7b, 0x66, 0x9e, 0xa9, 0xa2, 0x8c, 0x0e, 0x25, 0xbc, 0x3b,
	0xf8, 0xfb, 0xe5, 0x08, 0x0d, 0xc0, 0x5a, 0xd3, 0x72, 0x3f, 0xb0, 0x6a, 0x89, 0x86, 0x60, 0x6f,
	0x49, 0xba, 0x69, 0x46, 0x56, 0xff, 0xdc, 0x9a, 0x53, 0xc3, 0xbf, 0x02, 0x68, 0x0d, 0xa2, 0x11,
	0x38, 0x6f, 0x82, 0xe7, 0x29, 0x55, 0xb4, 0x1e, 0x78, 0x2f, 0x6a, 0x37, 0xfc, 0x73, 0x80, 0xa7,
	0x0d, 0xe7, 0x65, 0xfd, 0x2a, 0x06, 0x60, 0x71, 0xc9, 0x9a, 0x5b, 0xb8, 0x64, 0x8f, 0xf6, 0xab,
	0x15, 0x53, 0xb9, 0xec, 0xea, 0x9e, 0xdf, 0x7c, 0x06, 0x00, 0x00, 0xff, 0xff, 0xcd, 0xef, 0xc9,
	0x75, 0xdf, 0x02, 0x00, 0x00,
}

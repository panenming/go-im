// Code generated by protoc-gen-go. DO NOT EDIT.
// source: message.proto

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	message.proto

It has these top-level messages:
	Message
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
	context "golang.org/x/net/context"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

type Message struct {
	Id     int32  `protobuf:"varint,1,opt,name=Id" json:"Id,omitempty"`
	Fr     string `protobuf:"bytes,2,opt,name=Fr" json:"Fr,omitempty"`
	To     string `protobuf:"bytes,3,opt,name=To" json:"To,omitempty"`
	Txt    string `protobuf:"bytes,4,opt,name=Txt" json:"Txt,omitempty"`
	Status int32  `protobuf:"varint,5,opt,name=Status" json:"Status,omitempty"`
	St     int64  `protobuf:"varint,6,opt,name=St" json:"St,omitempty"`
	Type   int32  `protobuf:"varint,7,opt,name=Type" json:"Type,omitempty"`
	Extra  string `protobuf:"bytes,8,opt,name=Extra" json:"Extra,omitempty"`
}

func (m *Message) Reset()                    { *m = Message{} }
func (m *Message) String() string            { return proto1.CompactTextString(m) }
func (*Message) ProtoMessage()               {}
func (*Message) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Message) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Message) GetFr() string {
	if m != nil {
		return m.Fr
	}
	return ""
}

func (m *Message) GetTo() string {
	if m != nil {
		return m.To
	}
	return ""
}

func (m *Message) GetTxt() string {
	if m != nil {
		return m.Txt
	}
	return ""
}

func (m *Message) GetStatus() int32 {
	if m != nil {
		return m.Status
	}
	return 0
}

func (m *Message) GetSt() int64 {
	if m != nil {
		return m.St
	}
	return 0
}

func (m *Message) GetType() int32 {
	if m != nil {
		return m.Type
	}
	return 0
}

func (m *Message) GetExtra() string {
	if m != nil {
		return m.Extra
	}
	return ""
}

func init() {
	proto1.RegisterType((*Message)(nil), "proto.Message")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for MessageService service

type MessageServiceClient interface {
	CreateMsg(ctx context.Context, in *Message, opts ...client.CallOption) (*Message, error)
}

type messageServiceClient struct {
	c           client.Client
	serviceName string
}

func NewMessageServiceClient(serviceName string, c client.Client) MessageServiceClient {
	if c == nil {
		c = client.NewClient()
	}
	if len(serviceName) == 0 {
		serviceName = "proto"
	}
	return &messageServiceClient{
		c:           c,
		serviceName: serviceName,
	}
}

func (c *messageServiceClient) CreateMsg(ctx context.Context, in *Message, opts ...client.CallOption) (*Message, error) {
	req := c.c.NewRequest(c.serviceName, "MessageService.CreateMsg", in)
	out := new(Message)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for MessageService service

type MessageServiceHandler interface {
	CreateMsg(context.Context, *Message, *Message) error
}

func RegisterMessageServiceHandler(s server.Server, hdlr MessageServiceHandler, opts ...server.HandlerOption) {
	s.Handle(s.NewHandler(&MessageService{hdlr}, opts...))
}

type MessageService struct {
	MessageServiceHandler
}

func (h *MessageService) CreateMsg(ctx context.Context, in *Message, out *Message) error {
	return h.MessageServiceHandler.CreateMsg(ctx, in, out)
}

func init() { proto1.RegisterFile("message.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 204 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcd, 0x4d, 0x2d, 0x2e,
	0x4e, 0x4c, 0x4f, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x53, 0x4a, 0x73, 0x19,
	0xb9, 0xd8, 0x7d, 0x21, 0x12, 0x42, 0x7c, 0x5c, 0x4c, 0x9e, 0x29, 0x12, 0x8c, 0x0a, 0x8c, 0x1a,
	0xac, 0x41, 0x4c, 0x9e, 0x29, 0x20, 0xbe, 0x5b, 0x91, 0x04, 0x93, 0x02, 0xa3, 0x06, 0x67, 0x10,
	0x93, 0x5b, 0x11, 0x88, 0x1f, 0x92, 0x2f, 0xc1, 0x0c, 0xe1, 0x87, 0xe4, 0x0b, 0x09, 0x70, 0x31,
	0x87, 0x54, 0x94, 0x48, 0xb0, 0x80, 0x05, 0x40, 0x4c, 0x21, 0x31, 0x2e, 0xb6, 0xe0, 0x92, 0xc4,
	0x92, 0xd2, 0x62, 0x09, 0x56, 0xb0, 0x29, 0x50, 0x1e, 0x48, 0x67, 0x70, 0x89, 0x04, 0x9b, 0x02,
	0xa3, 0x06, 0x73, 0x10, 0x53, 0x70, 0x89, 0x90, 0x10, 0x17, 0x4b, 0x48, 0x65, 0x41, 0xaa, 0x04,
	0x3b, 0x58, 0x15, 0x98, 0x2d, 0x24, 0xc2, 0xc5, 0xea, 0x5a, 0x51, 0x52, 0x94, 0x28, 0xc1, 0x01,
	0x36, 0x0f, 0xc2, 0x31, 0xb2, 0xe7, 0xe2, 0x83, 0x3a, 0x2f, 0x38, 0xb5, 0xa8, 0x2c, 0x33, 0x39,
	0x55, 0x48, 0x97, 0x8b, 0xd3, 0xb9, 0x28, 0x35, 0xb1, 0x24, 0xd5, 0xb7, 0x38, 0x5d, 0x88, 0x0f,
	0xe2, 0x1b, 0x3d, 0xa8, 0x1a, 0x29, 0x34, 0xbe, 0x12, 0x43, 0x12, 0x1b, 0x58, 0xc0, 0x18, 0x10,
	0x00, 0x00, 0xff, 0xff, 0xf7, 0x45, 0x70, 0x83, 0xff, 0x00, 0x00, 0x00,
}

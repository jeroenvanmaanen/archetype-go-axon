// Code generated by protoc-gen-go. DO NOT EDIT.
// source: example.proto

package example

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Greeting struct {
	Message              string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Greeting) Reset()         { *m = Greeting{} }
func (m *Greeting) String() string { return proto.CompactTextString(m) }
func (*Greeting) ProtoMessage()    {}
func (*Greeting) Descriptor() ([]byte, []int) {
	return fileDescriptor_15a1dc8d40dadaa6, []int{0}
}

func (m *Greeting) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Greeting.Unmarshal(m, b)
}
func (m *Greeting) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Greeting.Marshal(b, m, deterministic)
}
func (m *Greeting) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Greeting.Merge(m, src)
}
func (m *Greeting) XXX_Size() int {
	return xxx_messageInfo_Greeting.Size(m)
}
func (m *Greeting) XXX_DiscardUnknown() {
	xxx_messageInfo_Greeting.DiscardUnknown(m)
}

var xxx_messageInfo_Greeting proto.InternalMessageInfo

func (m *Greeting) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type Acknowledgement struct {
	Message              string   `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Acknowledgement) Reset()         { *m = Acknowledgement{} }
func (m *Acknowledgement) String() string { return proto.CompactTextString(m) }
func (*Acknowledgement) ProtoMessage()    {}
func (*Acknowledgement) Descriptor() ([]byte, []int) {
	return fileDescriptor_15a1dc8d40dadaa6, []int{1}
}

func (m *Acknowledgement) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Acknowledgement.Unmarshal(m, b)
}
func (m *Acknowledgement) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Acknowledgement.Marshal(b, m, deterministic)
}
func (m *Acknowledgement) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Acknowledgement.Merge(m, src)
}
func (m *Acknowledgement) XXX_Size() int {
	return xxx_messageInfo_Acknowledgement.Size(m)
}
func (m *Acknowledgement) XXX_DiscardUnknown() {
	xxx_messageInfo_Acknowledgement.DiscardUnknown(m)
}

var xxx_messageInfo_Acknowledgement proto.InternalMessageInfo

func (m *Acknowledgement) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_15a1dc8d40dadaa6, []int{2}
}

func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (m *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(m, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

type GreetCommand struct {
	AggregateIdentifier  string    `protobuf:"bytes,1,opt,name=aggregateIdentifier,proto3" json:"aggregateIdentifier,omitempty"`
	Message              *Greeting `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *GreetCommand) Reset()         { *m = GreetCommand{} }
func (m *GreetCommand) String() string { return proto.CompactTextString(m) }
func (*GreetCommand) ProtoMessage()    {}
func (*GreetCommand) Descriptor() ([]byte, []int) {
	return fileDescriptor_15a1dc8d40dadaa6, []int{3}
}

func (m *GreetCommand) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GreetCommand.Unmarshal(m, b)
}
func (m *GreetCommand) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GreetCommand.Marshal(b, m, deterministic)
}
func (m *GreetCommand) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GreetCommand.Merge(m, src)
}
func (m *GreetCommand) XXX_Size() int {
	return xxx_messageInfo_GreetCommand.Size(m)
}
func (m *GreetCommand) XXX_DiscardUnknown() {
	xxx_messageInfo_GreetCommand.DiscardUnknown(m)
}

var xxx_messageInfo_GreetCommand proto.InternalMessageInfo

func (m *GreetCommand) GetAggregateIdentifier() string {
	if m != nil {
		return m.AggregateIdentifier
	}
	return ""
}

func (m *GreetCommand) GetMessage() *Greeting {
	if m != nil {
		return m.Message
	}
	return nil
}

type RecordCommand struct {
	AggregateIdentifier  string   `protobuf:"bytes,1,opt,name=aggregateIdentifier,proto3" json:"aggregateIdentifier,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RecordCommand) Reset()         { *m = RecordCommand{} }
func (m *RecordCommand) String() string { return proto.CompactTextString(m) }
func (*RecordCommand) ProtoMessage()    {}
func (*RecordCommand) Descriptor() ([]byte, []int) {
	return fileDescriptor_15a1dc8d40dadaa6, []int{4}
}

func (m *RecordCommand) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RecordCommand.Unmarshal(m, b)
}
func (m *RecordCommand) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RecordCommand.Marshal(b, m, deterministic)
}
func (m *RecordCommand) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RecordCommand.Merge(m, src)
}
func (m *RecordCommand) XXX_Size() int {
	return xxx_messageInfo_RecordCommand.Size(m)
}
func (m *RecordCommand) XXX_DiscardUnknown() {
	xxx_messageInfo_RecordCommand.DiscardUnknown(m)
}

var xxx_messageInfo_RecordCommand proto.InternalMessageInfo

func (m *RecordCommand) GetAggregateIdentifier() string {
	if m != nil {
		return m.AggregateIdentifier
	}
	return ""
}

type StopCommand struct {
	AggregateIdentifier  string   `protobuf:"bytes,1,opt,name=aggregateIdentifier,proto3" json:"aggregateIdentifier,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StopCommand) Reset()         { *m = StopCommand{} }
func (m *StopCommand) String() string { return proto.CompactTextString(m) }
func (*StopCommand) ProtoMessage()    {}
func (*StopCommand) Descriptor() ([]byte, []int) {
	return fileDescriptor_15a1dc8d40dadaa6, []int{5}
}

func (m *StopCommand) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StopCommand.Unmarshal(m, b)
}
func (m *StopCommand) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StopCommand.Marshal(b, m, deterministic)
}
func (m *StopCommand) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StopCommand.Merge(m, src)
}
func (m *StopCommand) XXX_Size() int {
	return xxx_messageInfo_StopCommand.Size(m)
}
func (m *StopCommand) XXX_DiscardUnknown() {
	xxx_messageInfo_StopCommand.DiscardUnknown(m)
}

var xxx_messageInfo_StopCommand proto.InternalMessageInfo

func (m *StopCommand) GetAggregateIdentifier() string {
	if m != nil {
		return m.AggregateIdentifier
	}
	return ""
}

type GreetedEvent struct {
	Message              *Greeting `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *GreetedEvent) Reset()         { *m = GreetedEvent{} }
func (m *GreetedEvent) String() string { return proto.CompactTextString(m) }
func (*GreetedEvent) ProtoMessage()    {}
func (*GreetedEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_15a1dc8d40dadaa6, []int{6}
}

func (m *GreetedEvent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GreetedEvent.Unmarshal(m, b)
}
func (m *GreetedEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GreetedEvent.Marshal(b, m, deterministic)
}
func (m *GreetedEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GreetedEvent.Merge(m, src)
}
func (m *GreetedEvent) XXX_Size() int {
	return xxx_messageInfo_GreetedEvent.Size(m)
}
func (m *GreetedEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_GreetedEvent.DiscardUnknown(m)
}

var xxx_messageInfo_GreetedEvent proto.InternalMessageInfo

func (m *GreetedEvent) GetMessage() *Greeting {
	if m != nil {
		return m.Message
	}
	return nil
}

type StartedRecordingEvent struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StartedRecordingEvent) Reset()         { *m = StartedRecordingEvent{} }
func (m *StartedRecordingEvent) String() string { return proto.CompactTextString(m) }
func (*StartedRecordingEvent) ProtoMessage()    {}
func (*StartedRecordingEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_15a1dc8d40dadaa6, []int{7}
}

func (m *StartedRecordingEvent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StartedRecordingEvent.Unmarshal(m, b)
}
func (m *StartedRecordingEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StartedRecordingEvent.Marshal(b, m, deterministic)
}
func (m *StartedRecordingEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StartedRecordingEvent.Merge(m, src)
}
func (m *StartedRecordingEvent) XXX_Size() int {
	return xxx_messageInfo_StartedRecordingEvent.Size(m)
}
func (m *StartedRecordingEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_StartedRecordingEvent.DiscardUnknown(m)
}

var xxx_messageInfo_StartedRecordingEvent proto.InternalMessageInfo

type StoppedRecordingEvent struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StoppedRecordingEvent) Reset()         { *m = StoppedRecordingEvent{} }
func (m *StoppedRecordingEvent) String() string { return proto.CompactTextString(m) }
func (*StoppedRecordingEvent) ProtoMessage()    {}
func (*StoppedRecordingEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_15a1dc8d40dadaa6, []int{8}
}

func (m *StoppedRecordingEvent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StoppedRecordingEvent.Unmarshal(m, b)
}
func (m *StoppedRecordingEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StoppedRecordingEvent.Marshal(b, m, deterministic)
}
func (m *StoppedRecordingEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StoppedRecordingEvent.Merge(m, src)
}
func (m *StoppedRecordingEvent) XXX_Size() int {
	return xxx_messageInfo_StoppedRecordingEvent.Size(m)
}
func (m *StoppedRecordingEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_StoppedRecordingEvent.DiscardUnknown(m)
}

var xxx_messageInfo_StoppedRecordingEvent proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Greeting)(nil), "org.leialearns.grpc.example.Greeting")
	proto.RegisterType((*Acknowledgement)(nil), "org.leialearns.grpc.example.Acknowledgement")
	proto.RegisterType((*Empty)(nil), "org.leialearns.grpc.example.Empty")
	proto.RegisterType((*GreetCommand)(nil), "org.leialearns.grpc.example.GreetCommand")
	proto.RegisterType((*RecordCommand)(nil), "org.leialearns.grpc.example.RecordCommand")
	proto.RegisterType((*StopCommand)(nil), "org.leialearns.grpc.example.StopCommand")
	proto.RegisterType((*GreetedEvent)(nil), "org.leialearns.grpc.example.GreetedEvent")
	proto.RegisterType((*StartedRecordingEvent)(nil), "org.leialearns.grpc.example.StartedRecordingEvent")
	proto.RegisterType((*StoppedRecordingEvent)(nil), "org.leialearns.grpc.example.StoppedRecordingEvent")
}

func init() {
	proto.RegisterFile("example.proto", fileDescriptor_15a1dc8d40dadaa6)
}

var fileDescriptor_15a1dc8d40dadaa6 = []byte{
	// 333 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x53, 0x5f, 0x4b, 0x3a, 0x41,
	0x14, 0x75, 0xe5, 0xa7, 0xfe, 0xbc, 0x66, 0xc1, 0xf4, 0x4f, 0xec, 0x45, 0x86, 0x02, 0xa1, 0x58,
	0xc5, 0x3e, 0x80, 0x58, 0x48, 0xf4, 0x54, 0xe8, 0x9b, 0x0f, 0xc1, 0xb4, 0x7b, 0x1b, 0x16, 0x77,
	0x67, 0x86, 0xbb, 0x83, 0xd5, 0x63, 0x5f, 0xb9, 0x4f, 0x10, 0xee, 0x1f, 0x24, 0x89, 0x65, 0x93,
	0x1e, 0x77, 0xcf, 0xb9, 0xe7, 0x9e, 0x73, 0xb8, 0x03, 0x6d, 0x7c, 0x13, 0x91, 0x09, 0xd1, 0x35,
	0xa4, 0xad, 0x66, 0x67, 0x9a, 0xa4, 0x1b, 0x62, 0x20, 0x42, 0x14, 0xa4, 0x62, 0x57, 0x92, 0xf1,
	0xdc, 0x8c, 0xc2, 0xcf, 0xe1, 0xff, 0x1d, 0x21, 0xda, 0x40, 0x49, 0xd6, 0x81, 0x46, 0x84, 0x71,
	0x2c, 0x24, 0x76, 0x9c, 0x9e, 0xd3, 0x6f, 0xce, 0xf2, 0x4f, 0x7e, 0x09, 0x07, 0x13, 0x6f, 0xa9,
	0xf4, 0x6b, 0x88, 0xbe, 0xc4, 0x08, 0x95, 0x2d, 0x20, 0x37, 0xa0, 0x36, 0x8d, 0x8c, 0x7d, 0xe7,
	0x1f, 0x0e, 0xec, 0x25, 0xe2, 0xb7, 0x3a, 0x8a, 0x84, 0xf2, 0xd9, 0x10, 0x0e, 0x85, 0x94, 0x84,
	0x52, 0x58, 0xbc, 0xf7, 0x51, 0xd9, 0xe0, 0x25, 0x40, 0xca, 0xe6, 0x7f, 0x82, 0xd8, 0x78, 0xb3,
	0xa5, 0xda, 0x73, 0xfa, 0xad, 0xd1, 0x85, 0x5b, 0x90, 0xc6, 0xcd, 0xa3, 0x6c, 0xcc, 0x4c, 0xa0,
	0x3d, 0x43, 0x4f, 0x93, 0xbf, 0xb3, 0x07, 0x3e, 0x86, 0xd6, 0xdc, 0x6a, 0xb3, 0xbb, 0xc0, 0x43,
	0x56, 0x03, 0xfa, 0xd3, 0xd5, 0xba, 0xba, 0xf1, 0xf7, 0xea, 0x7e, 0x1f, 0xea, 0x14, 0x8e, 0xe7,
	0x56, 0x90, 0x45, 0x3f, 0xcd, 0x16, 0x28, 0x99, 0x28, 0xa7, 0x80, 0x36, 0x66, 0x1b, 0x18, 0x7d,
	0x56, 0x61, 0x3f, 0xf5, 0x40, 0x73, 0xa4, 0x55, 0xe0, 0x21, 0x7b, 0x82, 0x5a, 0xf2, 0x87, 0x95,
	0xdb, 0xde, 0xbd, 0x2a, 0xa4, 0x6d, 0x9d, 0x07, 0xaf, 0xb0, 0x05, 0x34, 0xf3, 0xd9, 0x98, 0xf1,
	0xc2, 0xe1, 0xe4, 0x5c, 0xba, 0xe5, 0x7c, 0xf0, 0xca, 0xd0, 0x61, 0x33, 0xa8, 0xa7, 0x01, 0x4b,
	0x09, 0x97, 0xe0, 0xf0, 0x0a, 0x7b, 0x84, 0x7f, 0xeb, 0xee, 0xfe, 0x4e, 0xf1, 0xe6, 0x64, 0x71,
	0x14, 0x93, 0x37, 0x30, 0x4b, 0x39, 0x58, 0xe3, 0x83, 0x0c, 0x7f, 0xae, 0x27, 0xef, 0xf2, 0xfa,
	0x2b, 0x00, 0x00, 0xff, 0xff, 0x79, 0x95, 0x86, 0x14, 0xa8, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// GreeterServiceClient is the client API for GreeterService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type GreeterServiceClient interface {
	Greet(ctx context.Context, in *Greeting, opts ...grpc.CallOption) (*Acknowledgement, error)
	Greetings(ctx context.Context, in *Empty, opts ...grpc.CallOption) (GreeterService_GreetingsClient, error)
	Record(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	Stop(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
}

type greeterServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGreeterServiceClient(cc grpc.ClientConnInterface) GreeterServiceClient {
	return &greeterServiceClient{cc}
}

func (c *greeterServiceClient) Greet(ctx context.Context, in *Greeting, opts ...grpc.CallOption) (*Acknowledgement, error) {
	out := new(Acknowledgement)
	err := c.cc.Invoke(ctx, "/org.leialearns.grpc.example.GreeterService/Greet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *greeterServiceClient) Greetings(ctx context.Context, in *Empty, opts ...grpc.CallOption) (GreeterService_GreetingsClient, error) {
	stream, err := c.cc.NewStream(ctx, &_GreeterService_serviceDesc.Streams[0], "/org.leialearns.grpc.example.GreeterService/Greetings", opts...)
	if err != nil {
		return nil, err
	}
	x := &greeterServiceGreetingsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type GreeterService_GreetingsClient interface {
	Recv() (*Greeting, error)
	grpc.ClientStream
}

type greeterServiceGreetingsClient struct {
	grpc.ClientStream
}

func (x *greeterServiceGreetingsClient) Recv() (*Greeting, error) {
	m := new(Greeting)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *greeterServiceClient) Record(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/org.leialearns.grpc.example.GreeterService/Record", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *greeterServiceClient) Stop(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/org.leialearns.grpc.example.GreeterService/Stop", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GreeterServiceServer is the server API for GreeterService service.
type GreeterServiceServer interface {
	Greet(context.Context, *Greeting) (*Acknowledgement, error)
	Greetings(*Empty, GreeterService_GreetingsServer) error
	Record(context.Context, *Empty) (*Empty, error)
	Stop(context.Context, *Empty) (*Empty, error)
}

// UnimplementedGreeterServiceServer can be embedded to have forward compatible implementations.
type UnimplementedGreeterServiceServer struct {
}

func (*UnimplementedGreeterServiceServer) Greet(ctx context.Context, req *Greeting) (*Acknowledgement, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Greet not implemented")
}
func (*UnimplementedGreeterServiceServer) Greetings(req *Empty, srv GreeterService_GreetingsServer) error {
	return status.Errorf(codes.Unimplemented, "method Greetings not implemented")
}
func (*UnimplementedGreeterServiceServer) Record(ctx context.Context, req *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Record not implemented")
}
func (*UnimplementedGreeterServiceServer) Stop(ctx context.Context, req *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Stop not implemented")
}

func RegisterGreeterServiceServer(s *grpc.Server, srv GreeterServiceServer) {
	s.RegisterService(&_GreeterService_serviceDesc, srv)
}

func _GreeterService_Greet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Greeting)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServiceServer).Greet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/org.leialearns.grpc.example.GreeterService/Greet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServiceServer).Greet(ctx, req.(*Greeting))
	}
	return interceptor(ctx, in, info, handler)
}

func _GreeterService_Greetings_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GreeterServiceServer).Greetings(m, &greeterServiceGreetingsServer{stream})
}

type GreeterService_GreetingsServer interface {
	Send(*Greeting) error
	grpc.ServerStream
}

type greeterServiceGreetingsServer struct {
	grpc.ServerStream
}

func (x *greeterServiceGreetingsServer) Send(m *Greeting) error {
	return x.ServerStream.SendMsg(m)
}

func _GreeterService_Record_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServiceServer).Record(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/org.leialearns.grpc.example.GreeterService/Record",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServiceServer).Record(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _GreeterService_Stop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreeterServiceServer).Stop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/org.leialearns.grpc.example.GreeterService/Stop",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreeterServiceServer).Stop(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _GreeterService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "org.leialearns.grpc.example.GreeterService",
	HandlerType: (*GreeterServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Greet",
			Handler:    _GreeterService_Greet_Handler,
		},
		{
			MethodName: "Record",
			Handler:    _GreeterService_Record_Handler,
		},
		{
			MethodName: "Stop",
			Handler:    _GreeterService_Stop_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Greetings",
			Handler:       _GreeterService_Greetings_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "example.proto",
}

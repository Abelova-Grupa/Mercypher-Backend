// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v6.30.2
// source: relay-service.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	RelayService_SendMessage_FullMethodName = "/relay.RelayService/SendMessage"
	RelayService_GetMessages_FullMethodName = "/relay.RelayService/GetMessages"
)

// RelayServiceClient is the client API for RelayService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RelayServiceClient interface {
	SendMessage(ctx context.Context, in *ChatMessage, opts ...grpc.CallOption) (*Status, error)
	GetMessages(ctx context.Context, in *UserId, opts ...grpc.CallOption) (grpc.ServerStreamingClient[ChatMessage], error)
}

type relayServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRelayServiceClient(cc grpc.ClientConnInterface) RelayServiceClient {
	return &relayServiceClient{cc}
}

func (c *relayServiceClient) SendMessage(ctx context.Context, in *ChatMessage, opts ...grpc.CallOption) (*Status, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Status)
	err := c.cc.Invoke(ctx, RelayService_SendMessage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relayServiceClient) GetMessages(ctx context.Context, in *UserId, opts ...grpc.CallOption) (grpc.ServerStreamingClient[ChatMessage], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &RelayService_ServiceDesc.Streams[0], RelayService_GetMessages_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[UserId, ChatMessage]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type RelayService_GetMessagesClient = grpc.ServerStreamingClient[ChatMessage]

// RelayServiceServer is the server API for RelayService service.
// All implementations must embed UnimplementedRelayServiceServer
// for forward compatibility.
type RelayServiceServer interface {
	SendMessage(context.Context, *ChatMessage) (*Status, error)
	GetMessages(*UserId, grpc.ServerStreamingServer[ChatMessage]) error
	mustEmbedUnimplementedRelayServiceServer()
}

// UnimplementedRelayServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedRelayServiceServer struct{}

func (UnimplementedRelayServiceServer) SendMessage(context.Context, *ChatMessage) (*Status, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
}
func (UnimplementedRelayServiceServer) GetMessages(*UserId, grpc.ServerStreamingServer[ChatMessage]) error {
	return status.Errorf(codes.Unimplemented, "method GetMessages not implemented")
}
func (UnimplementedRelayServiceServer) mustEmbedUnimplementedRelayServiceServer() {}
func (UnimplementedRelayServiceServer) testEmbeddedByValue()                      {}

// UnsafeRelayServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RelayServiceServer will
// result in compilation errors.
type UnsafeRelayServiceServer interface {
	mustEmbedUnimplementedRelayServiceServer()
}

func RegisterRelayServiceServer(s grpc.ServiceRegistrar, srv RelayServiceServer) {
	// If the following call pancis, it indicates UnimplementedRelayServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&RelayService_ServiceDesc, srv)
}

func _RelayService_SendMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChatMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelayServiceServer).SendMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RelayService_SendMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelayServiceServer).SendMessage(ctx, req.(*ChatMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _RelayService_GetMessages_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(UserId)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(RelayServiceServer).GetMessages(m, &grpc.GenericServerStream[UserId, ChatMessage]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type RelayService_GetMessagesServer = grpc.ServerStreamingServer[ChatMessage]

// RelayService_ServiceDesc is the grpc.ServiceDesc for RelayService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RelayService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "relay.RelayService",
	HandlerType: (*RelayServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendMessage",
			Handler:    _RelayService_SendMessage_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetMessages",
			Handler:       _RelayService_GetMessages_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "relay-service.proto",
}

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ConfigCenterClient is the client API for ConfigCenter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ConfigCenterClient interface {
	Set(ctx context.Context, in *ServiceConfig, opts ...grpc.CallOption) (*Result, error)
	Watch(ctx context.Context, in *ServiceConfig, opts ...grpc.CallOption) (ConfigCenter_WatchClient, error)
}

type configCenterClient struct {
	cc grpc.ClientConnInterface
}

func NewConfigCenterClient(cc grpc.ClientConnInterface) ConfigCenterClient {
	return &configCenterClient{cc}
}

func (c *configCenterClient) Set(ctx context.Context, in *ServiceConfig, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := c.cc.Invoke(ctx, "/config.api.ConfigCenter/Set", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configCenterClient) Watch(ctx context.Context, in *ServiceConfig, opts ...grpc.CallOption) (ConfigCenter_WatchClient, error) {
	stream, err := c.cc.NewStream(ctx, &ConfigCenter_ServiceDesc.Streams[0], "/config.api.ConfigCenter/Watch", opts...)
	if err != nil {
		return nil, err
	}
	x := &configCenterWatchClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ConfigCenter_WatchClient interface {
	Recv() (*ServiceConfig, error)
	grpc.ClientStream
}

type configCenterWatchClient struct {
	grpc.ClientStream
}

func (x *configCenterWatchClient) Recv() (*ServiceConfig, error) {
	m := new(ServiceConfig)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ConfigCenterServer is the server API for ConfigCenter service.
// All implementations must embed UnimplementedConfigCenterServer
// for forward compatibility
type ConfigCenterServer interface {
	Set(context.Context, *ServiceConfig) (*Result, error)
	Watch(*ServiceConfig, ConfigCenter_WatchServer) error
	mustEmbedUnimplementedConfigCenterServer()
}

// UnimplementedConfigCenterServer must be embedded to have forward compatible implementations.
type UnimplementedConfigCenterServer struct {
}

func (UnimplementedConfigCenterServer) Set(context.Context, *ServiceConfig) (*Result, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Set not implemented")
}
func (UnimplementedConfigCenterServer) Watch(*ServiceConfig, ConfigCenter_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "method Watch not implemented")
}
func (UnimplementedConfigCenterServer) mustEmbedUnimplementedConfigCenterServer() {}

// UnsafeConfigCenterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ConfigCenterServer will
// result in compilation errors.
type UnsafeConfigCenterServer interface {
	mustEmbedUnimplementedConfigCenterServer()
}

func RegisterConfigCenterServer(s grpc.ServiceRegistrar, srv ConfigCenterServer) {
	s.RegisterService(&ConfigCenter_ServiceDesc, srv)
}

func _ConfigCenter_Set_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ServiceConfig)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigCenterServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/config.api.ConfigCenter/Set",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigCenterServer).Set(ctx, req.(*ServiceConfig))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConfigCenter_Watch_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ServiceConfig)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ConfigCenterServer).Watch(m, &configCenterWatchServer{stream})
}

type ConfigCenter_WatchServer interface {
	Send(*ServiceConfig) error
	grpc.ServerStream
}

type configCenterWatchServer struct {
	grpc.ServerStream
}

func (x *configCenterWatchServer) Send(m *ServiceConfig) error {
	return x.ServerStream.SendMsg(m)
}

// ConfigCenter_ServiceDesc is the grpc.ServiceDesc for ConfigCenter service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ConfigCenter_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "config.api.ConfigCenter",
	HandlerType: (*ConfigCenterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Set",
			Handler:    _ConfigCenter_Set_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Watch",
			Handler:       _ConfigCenter_Watch_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "config.proto",
}
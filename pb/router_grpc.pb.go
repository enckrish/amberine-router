// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: pb/router.proto

package pb

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

// RouterClient is the client API for Router service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RouterClient interface {
	Init_Type0(ctx context.Context, in *InitRequest_Type0, opts ...grpc.CallOption) (*InitResponse_Type0, error)
	RouteLog_Type0(ctx context.Context, opts ...grpc.CallOption) (Router_RouteLog_Type0Client, error)
}

type routerClient struct {
	cc grpc.ClientConnInterface
}

func NewRouterClient(cc grpc.ClientConnInterface) RouterClient {
	return &routerClient{cc}
}

func (c *routerClient) Init_Type0(ctx context.Context, in *InitRequest_Type0, opts ...grpc.CallOption) (*InitResponse_Type0, error) {
	out := new(InitResponse_Type0)
	err := c.cc.Invoke(ctx, "/Router/init_Type0", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *routerClient) RouteLog_Type0(ctx context.Context, opts ...grpc.CallOption) (Router_RouteLog_Type0Client, error) {
	stream, err := c.cc.NewStream(ctx, &Router_ServiceDesc.Streams[0], "/Router/routeLog_Type0", opts...)
	if err != nil {
		return nil, err
	}
	x := &routerRouteLog_Type0Client{stream}
	return x, nil
}

type Router_RouteLog_Type0Client interface {
	Send(*AnalyzerRequest_Type0) error
	Recv() (*AnalyzerResponse, error)
	grpc.ClientStream
}

type routerRouteLog_Type0Client struct {
	grpc.ClientStream
}

func (x *routerRouteLog_Type0Client) Send(m *AnalyzerRequest_Type0) error {
	return x.ClientStream.SendMsg(m)
}

func (x *routerRouteLog_Type0Client) Recv() (*AnalyzerResponse, error) {
	m := new(AnalyzerResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// RouterServer is the server API for Router service.
// All implementations must embed UnimplementedRouterServer
// for forward compatibility
type RouterServer interface {
	Init_Type0(context.Context, *InitRequest_Type0) (*InitResponse_Type0, error)
	RouteLog_Type0(Router_RouteLog_Type0Server) error
	mustEmbedUnimplementedRouterServer()
}

// UnimplementedRouterServer must be embedded to have forward compatible implementations.
type UnimplementedRouterServer struct {
}

func (UnimplementedRouterServer) Init_Type0(context.Context, *InitRequest_Type0) (*InitResponse_Type0, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Init_Type0 not implemented")
}
func (UnimplementedRouterServer) RouteLog_Type0(Router_RouteLog_Type0Server) error {
	return status.Errorf(codes.Unimplemented, "method RouteLog_Type0 not implemented")
}
func (UnimplementedRouterServer) mustEmbedUnimplementedRouterServer() {}

// UnsafeRouterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RouterServer will
// result in compilation errors.
type UnsafeRouterServer interface {
	mustEmbedUnimplementedRouterServer()
}

func RegisterRouterServer(s grpc.ServiceRegistrar, srv RouterServer) {
	s.RegisterService(&Router_ServiceDesc, srv)
}

func _Router_Init_Type0_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InitRequest_Type0)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RouterServer).Init_Type0(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Router/init_Type0",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RouterServer).Init_Type0(ctx, req.(*InitRequest_Type0))
	}
	return interceptor(ctx, in, info, handler)
}

func _Router_RouteLog_Type0_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(RouterServer).RouteLog_Type0(&routerRouteLog_Type0Server{stream})
}

type Router_RouteLog_Type0Server interface {
	Send(*AnalyzerResponse) error
	Recv() (*AnalyzerRequest_Type0, error)
	grpc.ServerStream
}

type routerRouteLog_Type0Server struct {
	grpc.ServerStream
}

func (x *routerRouteLog_Type0Server) Send(m *AnalyzerResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *routerRouteLog_Type0Server) Recv() (*AnalyzerRequest_Type0, error) {
	m := new(AnalyzerRequest_Type0)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Router_ServiceDesc is the grpc.ServiceDesc for Router service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Router_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Router",
	HandlerType: (*RouterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "init_Type0",
			Handler:    _Router_Init_Type0_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "routeLog_Type0",
			Handler:       _Router_RouteLog_Type0_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "pb/router.proto",
}

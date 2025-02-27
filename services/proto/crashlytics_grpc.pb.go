// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v5.27.3
// source: crashlytics.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	Crashlytics_RecordError_FullMethodName = "/crashlytics.v1.Crashlytics/RecordError"
)

// CrashlyticsClient is the client API for Crashlytics service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CrashlyticsClient interface {
	RecordError(ctx context.Context, in *Error, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type crashlyticsClient struct {
	cc grpc.ClientConnInterface
}

func NewCrashlyticsClient(cc grpc.ClientConnInterface) CrashlyticsClient {
	return &crashlyticsClient{cc}
}

func (c *crashlyticsClient) RecordError(ctx context.Context, in *Error, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Crashlytics_RecordError_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CrashlyticsServer is the server API for Crashlytics service.
// All implementations must embed UnimplementedCrashlyticsServer
// for forward compatibility
type CrashlyticsServer interface {
	RecordError(context.Context, *Error) (*emptypb.Empty, error)
	mustEmbedUnimplementedCrashlyticsServer()
}

// UnimplementedCrashlyticsServer must be embedded to have forward compatible implementations.
type UnimplementedCrashlyticsServer struct {
}

func (UnimplementedCrashlyticsServer) RecordError(context.Context, *Error) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecordError not implemented")
}
func (UnimplementedCrashlyticsServer) mustEmbedUnimplementedCrashlyticsServer() {}

// UnsafeCrashlyticsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CrashlyticsServer will
// result in compilation errors.
type UnsafeCrashlyticsServer interface {
	mustEmbedUnimplementedCrashlyticsServer()
}

func RegisterCrashlyticsServer(s grpc.ServiceRegistrar, srv CrashlyticsServer) {
	s.RegisterService(&Crashlytics_ServiceDesc, srv)
}

func _Crashlytics_RecordError_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Error)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CrashlyticsServer).RecordError(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Crashlytics_RecordError_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CrashlyticsServer).RecordError(ctx, req.(*Error))
	}
	return interceptor(ctx, in, info, handler)
}

// Crashlytics_ServiceDesc is the grpc.ServiceDesc for Crashlytics service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Crashlytics_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "crashlytics.v1.Crashlytics",
	HandlerType: (*CrashlyticsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RecordError",
			Handler:    _Crashlytics_RecordError_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "crashlytics.proto",
}

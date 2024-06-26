// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: account_service.proto

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
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	AccountService_GetCosts_FullMethodName        = "/chatbot.account.v1.AccountService/GetCosts"
	AccountService_GetPayments_FullMethodName     = "/chatbot.account.v1.AccountService/GetPayments"
	AccountService_GetBalanceSheet_FullMethodName = "/chatbot.account.v1.AccountService/GetBalanceSheet"
)

// AccountServiceClient is the client API for AccountService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AccountServiceClient interface {
	GetCosts(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Costs, error)
	GetPayments(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Payments, error)
	GetBalanceSheet(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*BalanceSheet, error)
}

type accountServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAccountServiceClient(cc grpc.ClientConnInterface) AccountServiceClient {
	return &accountServiceClient{cc}
}

func (c *accountServiceClient) GetCosts(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Costs, error) {
	out := new(Costs)
	err := c.cc.Invoke(ctx, AccountService_GetCosts_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) GetPayments(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Payments, error) {
	out := new(Payments)
	err := c.cc.Invoke(ctx, AccountService_GetPayments_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *accountServiceClient) GetBalanceSheet(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*BalanceSheet, error) {
	out := new(BalanceSheet)
	err := c.cc.Invoke(ctx, AccountService_GetBalanceSheet_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AccountServiceServer is the server API for AccountService service.
// All implementations must embed UnimplementedAccountServiceServer
// for forward compatibility
type AccountServiceServer interface {
	GetCosts(context.Context, *emptypb.Empty) (*Costs, error)
	GetPayments(context.Context, *emptypb.Empty) (*Payments, error)
	GetBalanceSheet(context.Context, *emptypb.Empty) (*BalanceSheet, error)
	mustEmbedUnimplementedAccountServiceServer()
}

// UnimplementedAccountServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAccountServiceServer struct {
}

func (UnimplementedAccountServiceServer) GetCosts(context.Context, *emptypb.Empty) (*Costs, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCosts not implemented")
}
func (UnimplementedAccountServiceServer) GetPayments(context.Context, *emptypb.Empty) (*Payments, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPayments not implemented")
}
func (UnimplementedAccountServiceServer) GetBalanceSheet(context.Context, *emptypb.Empty) (*BalanceSheet, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBalanceSheet not implemented")
}
func (UnimplementedAccountServiceServer) mustEmbedUnimplementedAccountServiceServer() {}

// UnsafeAccountServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AccountServiceServer will
// result in compilation errors.
type UnsafeAccountServiceServer interface {
	mustEmbedUnimplementedAccountServiceServer()
}

func RegisterAccountServiceServer(s grpc.ServiceRegistrar, srv AccountServiceServer) {
	s.RegisterService(&AccountService_ServiceDesc, srv)
}

func _AccountService_GetCosts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).GetCosts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_GetCosts_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).GetCosts(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_GetPayments_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).GetPayments(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_GetPayments_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).GetPayments(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _AccountService_GetBalanceSheet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AccountServiceServer).GetBalanceSheet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AccountService_GetBalanceSheet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AccountServiceServer).GetBalanceSheet(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// AccountService_ServiceDesc is the grpc.ServiceDesc for AccountService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AccountService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chatbot.account.v1.AccountService",
	HandlerType: (*AccountServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetCosts",
			Handler:    _AccountService_GetCosts_Handler,
		},
		{
			MethodName: "GetPayments",
			Handler:    _AccountService_GetPayments_Handler,
		},
		{
			MethodName: "GetBalanceSheet",
			Handler:    _AccountService_GetBalanceSheet_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "account_service.proto",
}

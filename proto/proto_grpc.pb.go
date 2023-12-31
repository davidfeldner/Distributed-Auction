// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.1
// source: proto/proto.proto

package proto

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

// AuctionServiceClient is the client API for AuctionService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuctionServiceClient interface {
	Bid(ctx context.Context, in *BidRequest, opts ...grpc.CallOption) (*Ack, error)
	Result(ctx context.Context, in *EmptyMessage, opts ...grpc.CallOption) (*ResultResponse, error)
	ServerIsDesynchronized(ctx context.Context, in *ResultResponse, opts ...grpc.CallOption) (*Ack, error)
}

type auctionServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAuctionServiceClient(cc grpc.ClientConnInterface) AuctionServiceClient {
	return &auctionServiceClient{cc}
}

func (c *auctionServiceClient) Bid(ctx context.Context, in *BidRequest, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := c.cc.Invoke(ctx, "/main.AuctionService/Bid", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *auctionServiceClient) Result(ctx context.Context, in *EmptyMessage, opts ...grpc.CallOption) (*ResultResponse, error) {
	out := new(ResultResponse)
	err := c.cc.Invoke(ctx, "/main.AuctionService/Result", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *auctionServiceClient) ServerIsDesynchronized(ctx context.Context, in *ResultResponse, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := c.cc.Invoke(ctx, "/main.AuctionService/ServerIsDesynchronized", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuctionServiceServer is the server API for AuctionService service.
// All implementations must embed UnimplementedAuctionServiceServer
// for forward compatibility
type AuctionServiceServer interface {
	Bid(context.Context, *BidRequest) (*Ack, error)
	Result(context.Context, *EmptyMessage) (*ResultResponse, error)
	ServerIsDesynchronized(context.Context, *ResultResponse) (*Ack, error)
	mustEmbedUnimplementedAuctionServiceServer()
}

// UnimplementedAuctionServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAuctionServiceServer struct {
}

func (UnimplementedAuctionServiceServer) Bid(context.Context, *BidRequest) (*Ack, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Bid not implemented")
}
func (UnimplementedAuctionServiceServer) Result(context.Context, *EmptyMessage) (*ResultResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Result not implemented")
}
func (UnimplementedAuctionServiceServer) ServerIsDesynchronized(context.Context, *ResultResponse) (*Ack, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ServerIsDesynchronized not implemented")
}
func (UnimplementedAuctionServiceServer) mustEmbedUnimplementedAuctionServiceServer() {}

// UnsafeAuctionServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuctionServiceServer will
// result in compilation errors.
type UnsafeAuctionServiceServer interface {
	mustEmbedUnimplementedAuctionServiceServer()
}

func RegisterAuctionServiceServer(s grpc.ServiceRegistrar, srv AuctionServiceServer) {
	s.RegisterService(&AuctionService_ServiceDesc, srv)
}

func _AuctionService_Bid_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BidRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuctionServiceServer).Bid(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.AuctionService/Bid",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuctionServiceServer).Bid(ctx, req.(*BidRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuctionService_Result_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuctionServiceServer).Result(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.AuctionService/Result",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuctionServiceServer).Result(ctx, req.(*EmptyMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _AuctionService_ServerIsDesynchronized_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResultResponse)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuctionServiceServer).ServerIsDesynchronized(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.AuctionService/ServerIsDesynchronized",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuctionServiceServer).ServerIsDesynchronized(ctx, req.(*ResultResponse))
	}
	return interceptor(ctx, in, info, handler)
}

// AuctionService_ServiceDesc is the grpc.ServiceDesc for AuctionService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AuctionService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "main.AuctionService",
	HandlerType: (*AuctionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Bid",
			Handler:    _AuctionService_Bid_Handler,
		},
		{
			MethodName: "Result",
			Handler:    _AuctionService_Result_Handler,
		},
		{
			MethodName: "ServerIsDesynchronized",
			Handler:    _AuctionService_ServerIsDesynchronized_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/proto.proto",
}

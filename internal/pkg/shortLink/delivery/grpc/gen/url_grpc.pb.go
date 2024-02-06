// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: proto/url.proto

package gen

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

// ShortLinkClient is the client API for ShortLink service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ShortLinkClient interface {
	CreateLink(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
	GetLink(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
}

type shortLinkClient struct {
	cc grpc.ClientConnInterface
}

func NewShortLinkClient(cc grpc.ClientConnInterface) ShortLinkClient {
	return &shortLinkClient{cc}
}

func (c *shortLinkClient) CreateLink(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/shortener.ShortLink/CreateLink", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortLinkClient) GetLink(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/shortener.ShortLink/GetLink", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ShortLinkServer is the server API for ShortLink service.
// All implementations must embed UnimplementedShortLinkServer
// for forward compatibility
type ShortLinkServer interface {
	CreateLink(context.Context, *Request) (*Response, error)
	GetLink(context.Context, *Request) (*Response, error)
	mustEmbedUnimplementedShortLinkServer()
}

// UnimplementedShortLinkServer must be embedded to have forward compatible implementations.
type UnimplementedShortLinkServer struct {
}

func (UnimplementedShortLinkServer) CreateLink(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateLink not implemented")
}
func (UnimplementedShortLinkServer) GetLink(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLink not implemented")
}
func (UnimplementedShortLinkServer) mustEmbedUnimplementedShortLinkServer() {}

// UnsafeShortLinkServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ShortLinkServer will
// result in compilation errors.
type UnsafeShortLinkServer interface {
	mustEmbedUnimplementedShortLinkServer()
}

func RegisterShortLinkServer(s grpc.ServiceRegistrar, srv ShortLinkServer) {
	s.RegisterService(&ShortLink_ServiceDesc, srv)
}

func _ShortLink_CreateLink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortLinkServer).CreateLink(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.ShortLink/CreateLink",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortLinkServer).CreateLink(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShortLink_GetLink_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortLinkServer).GetLink(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shortener.ShortLink/GetLink",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortLinkServer).GetLink(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

// ShortLink_ServiceDesc is the grpc.ServiceDesc for ShortLink service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ShortLink_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "shortener.ShortLink",
	HandlerType: (*ShortLinkServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateLink",
			Handler:    _ShortLink_CreateLink_Handler,
		},
		{
			MethodName: "GetLink",
			Handler:    _ShortLink_GetLink_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/url.proto",
}

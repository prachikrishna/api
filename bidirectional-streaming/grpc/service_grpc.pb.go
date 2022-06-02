// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package grpc

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

// SearchBookServiceClient is the client API for SearchBookService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SearchBookServiceClient interface {
	FindMatch(ctx context.Context, opts ...grpc.CallOption) (SearchBookService_FindMatchClient, error)
}

type searchBookServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSearchBookServiceClient(cc grpc.ClientConnInterface) SearchBookServiceClient {
	return &searchBookServiceClient{cc}
}

func (c *searchBookServiceClient) FindMatch(ctx context.Context, opts ...grpc.CallOption) (SearchBookService_FindMatchClient, error) {
	stream, err := c.cc.NewStream(ctx, &SearchBookService_ServiceDesc.Streams[0], "/v1.SearchBookService/FindMatch", opts...)
	if err != nil {
		return nil, err
	}
	x := &searchBookServiceFindMatchClient{stream}
	return x, nil
}

type SearchBookService_FindMatchClient interface {
	Send(*Request) error
	Recv() (*Response, error)
	grpc.ClientStream
}

type searchBookServiceFindMatchClient struct {
	grpc.ClientStream
}

func (x *searchBookServiceFindMatchClient) Send(m *Request) error {
	return x.ClientStream.SendMsg(m)
}

func (x *searchBookServiceFindMatchClient) Recv() (*Response, error) {
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SearchBookServiceServer is the server API for SearchBookService service.
// All implementations must embed UnimplementedSearchBookServiceServer
// for forward compatibility
type SearchBookServiceServer interface {
	FindMatch(SearchBookService_FindMatchServer) error
	mustEmbedUnimplementedSearchBookServiceServer()
}

// UnimplementedSearchBookServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSearchBookServiceServer struct {
}

func (UnimplementedSearchBookServiceServer) FindMatch(SearchBookService_FindMatchServer) error {
	return status.Errorf(codes.Unimplemented, "method FindMatch not implemented")
}
func (UnimplementedSearchBookServiceServer) mustEmbedUnimplementedSearchBookServiceServer() {}

// UnsafeSearchBookServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SearchBookServiceServer will
// result in compilation errors.
type UnsafeSearchBookServiceServer interface {
	mustEmbedUnimplementedSearchBookServiceServer()
}

func RegisterSearchBookServiceServer(s grpc.ServiceRegistrar, srv SearchBookServiceServer) {
	s.RegisterService(&SearchBookService_ServiceDesc, srv)
}

func _SearchBookService_FindMatch_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SearchBookServiceServer).FindMatch(&searchBookServiceFindMatchServer{stream})
}

type SearchBookService_FindMatchServer interface {
	Send(*Response) error
	Recv() (*Request, error)
	grpc.ServerStream
}

type searchBookServiceFindMatchServer struct {
	grpc.ServerStream
}

func (x *searchBookServiceFindMatchServer) Send(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *searchBookServiceFindMatchServer) Recv() (*Request, error) {
	m := new(Request)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SearchBookService_ServiceDesc is the grpc.ServiceDesc for SearchBookService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SearchBookService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "v1.SearchBookService",
	HandlerType: (*SearchBookServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "FindMatch",
			Handler:       _SearchBookService_FindMatch_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "service.proto",
}

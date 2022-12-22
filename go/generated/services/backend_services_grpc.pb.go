// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.11
// source: backend_services.proto

package services

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

// CollectionServiceClient is the client API for CollectionService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CollectionServiceClient interface {
	StoreCollection(ctx context.Context, in *StoreCollectionRequest, opts ...grpc.CallOption) (*StoreCollectionResponse, error)
	FindCollection(ctx context.Context, in *FindCollectionRequest, opts ...grpc.CallOption) (*FindCollectionResponse, error)
	CreateSchema(ctx context.Context, in *CreateSchemaRequest, opts ...grpc.CallOption) (*CreateSchemaResponse, error)
	ValidateCollection(ctx context.Context, in *ValidateCollectionRequest, opts ...grpc.CallOption) (*ValidateCollectionResponse, error)
	InvalidateSchema(ctx context.Context, in *InvalidateSchemaRequest, opts ...grpc.CallOption) (*InvalidateSchemaResponse, error)
}

type collectionServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCollectionServiceClient(cc grpc.ClientConnInterface) CollectionServiceClient {
	return &collectionServiceClient{cc}
}

func (c *collectionServiceClient) StoreCollection(ctx context.Context, in *StoreCollectionRequest, opts ...grpc.CallOption) (*StoreCollectionResponse, error) {
	out := new(StoreCollectionResponse)
	err := c.cc.Invoke(ctx, "/services.CollectionService/StoreCollection", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *collectionServiceClient) FindCollection(ctx context.Context, in *FindCollectionRequest, opts ...grpc.CallOption) (*FindCollectionResponse, error) {
	out := new(FindCollectionResponse)
	err := c.cc.Invoke(ctx, "/services.CollectionService/FindCollection", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *collectionServiceClient) CreateSchema(ctx context.Context, in *CreateSchemaRequest, opts ...grpc.CallOption) (*CreateSchemaResponse, error) {
	out := new(CreateSchemaResponse)
	err := c.cc.Invoke(ctx, "/services.CollectionService/CreateSchema", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *collectionServiceClient) ValidateCollection(ctx context.Context, in *ValidateCollectionRequest, opts ...grpc.CallOption) (*ValidateCollectionResponse, error) {
	out := new(ValidateCollectionResponse)
	err := c.cc.Invoke(ctx, "/services.CollectionService/ValidateCollection", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *collectionServiceClient) InvalidateSchema(ctx context.Context, in *InvalidateSchemaRequest, opts ...grpc.CallOption) (*InvalidateSchemaResponse, error) {
	out := new(InvalidateSchemaResponse)
	err := c.cc.Invoke(ctx, "/services.CollectionService/InvalidateSchema", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CollectionServiceServer is the server API for CollectionService service.
// All implementations must embed UnimplementedCollectionServiceServer
// for forward compatibility
type CollectionServiceServer interface {
	StoreCollection(context.Context, *StoreCollectionRequest) (*StoreCollectionResponse, error)
	FindCollection(context.Context, *FindCollectionRequest) (*FindCollectionResponse, error)
	CreateSchema(context.Context, *CreateSchemaRequest) (*CreateSchemaResponse, error)
	ValidateCollection(context.Context, *ValidateCollectionRequest) (*ValidateCollectionResponse, error)
	InvalidateSchema(context.Context, *InvalidateSchemaRequest) (*InvalidateSchemaResponse, error)
	mustEmbedUnimplementedCollectionServiceServer()
}

// UnimplementedCollectionServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCollectionServiceServer struct {
}

func (UnimplementedCollectionServiceServer) StoreCollection(context.Context, *StoreCollectionRequest) (*StoreCollectionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StoreCollection not implemented")
}
func (UnimplementedCollectionServiceServer) FindCollection(context.Context, *FindCollectionRequest) (*FindCollectionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindCollection not implemented")
}
func (UnimplementedCollectionServiceServer) CreateSchema(context.Context, *CreateSchemaRequest) (*CreateSchemaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSchema not implemented")
}
func (UnimplementedCollectionServiceServer) ValidateCollection(context.Context, *ValidateCollectionRequest) (*ValidateCollectionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidateCollection not implemented")
}
func (UnimplementedCollectionServiceServer) InvalidateSchema(context.Context, *InvalidateSchemaRequest) (*InvalidateSchemaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InvalidateSchema not implemented")
}
func (UnimplementedCollectionServiceServer) mustEmbedUnimplementedCollectionServiceServer() {}

// UnsafeCollectionServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CollectionServiceServer will
// result in compilation errors.
type UnsafeCollectionServiceServer interface {
	mustEmbedUnimplementedCollectionServiceServer()
}

func RegisterCollectionServiceServer(s grpc.ServiceRegistrar, srv CollectionServiceServer) {
	s.RegisterService(&CollectionService_ServiceDesc, srv)
}

func _CollectionService_StoreCollection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StoreCollectionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CollectionServiceServer).StoreCollection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.CollectionService/StoreCollection",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CollectionServiceServer).StoreCollection(ctx, req.(*StoreCollectionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CollectionService_FindCollection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindCollectionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CollectionServiceServer).FindCollection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.CollectionService/FindCollection",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CollectionServiceServer).FindCollection(ctx, req.(*FindCollectionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CollectionService_CreateSchema_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSchemaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CollectionServiceServer).CreateSchema(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.CollectionService/CreateSchema",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CollectionServiceServer).CreateSchema(ctx, req.(*CreateSchemaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CollectionService_ValidateCollection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValidateCollectionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CollectionServiceServer).ValidateCollection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.CollectionService/ValidateCollection",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CollectionServiceServer).ValidateCollection(ctx, req.(*ValidateCollectionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CollectionService_InvalidateSchema_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InvalidateSchemaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CollectionServiceServer).InvalidateSchema(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.CollectionService/InvalidateSchema",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CollectionServiceServer).InvalidateSchema(ctx, req.(*InvalidateSchemaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CollectionService_ServiceDesc is the grpc.ServiceDesc for CollectionService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CollectionService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "services.CollectionService",
	HandlerType: (*CollectionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StoreCollection",
			Handler:    _CollectionService_StoreCollection_Handler,
		},
		{
			MethodName: "FindCollection",
			Handler:    _CollectionService_FindCollection_Handler,
		},
		{
			MethodName: "CreateSchema",
			Handler:    _CollectionService_CreateSchema_Handler,
		},
		{
			MethodName: "ValidateCollection",
			Handler:    _CollectionService_ValidateCollection_Handler,
		},
		{
			MethodName: "InvalidateSchema",
			Handler:    _CollectionService_InvalidateSchema_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "backend_services.proto",
}

// KeyServiceClient is the client API for KeyService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KeyServiceClient interface {
	VerifyMessageAttestation(ctx context.Context, in *VerifyMessageAttestationRequest, opts ...grpc.CallOption) (*VerifyMessageAttestationResponse, error)
}

type keyServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewKeyServiceClient(cc grpc.ClientConnInterface) KeyServiceClient {
	return &keyServiceClient{cc}
}

func (c *keyServiceClient) VerifyMessageAttestation(ctx context.Context, in *VerifyMessageAttestationRequest, opts ...grpc.CallOption) (*VerifyMessageAttestationResponse, error) {
	out := new(VerifyMessageAttestationResponse)
	err := c.cc.Invoke(ctx, "/services.KeyService/VerifyMessageAttestation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KeyServiceServer is the server API for KeyService service.
// All implementations must embed UnimplementedKeyServiceServer
// for forward compatibility
type KeyServiceServer interface {
	VerifyMessageAttestation(context.Context, *VerifyMessageAttestationRequest) (*VerifyMessageAttestationResponse, error)
	mustEmbedUnimplementedKeyServiceServer()
}

// UnimplementedKeyServiceServer must be embedded to have forward compatible implementations.
type UnimplementedKeyServiceServer struct {
}

func (UnimplementedKeyServiceServer) VerifyMessageAttestation(context.Context, *VerifyMessageAttestationRequest) (*VerifyMessageAttestationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyMessageAttestation not implemented")
}
func (UnimplementedKeyServiceServer) mustEmbedUnimplementedKeyServiceServer() {}

// UnsafeKeyServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KeyServiceServer will
// result in compilation errors.
type UnsafeKeyServiceServer interface {
	mustEmbedUnimplementedKeyServiceServer()
}

func RegisterKeyServiceServer(s grpc.ServiceRegistrar, srv KeyServiceServer) {
	s.RegisterService(&KeyService_ServiceDesc, srv)
}

func _KeyService_VerifyMessageAttestation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerifyMessageAttestationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeyServiceServer).VerifyMessageAttestation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/services.KeyService/VerifyMessageAttestation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeyServiceServer).VerifyMessageAttestation(ctx, req.(*VerifyMessageAttestationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// KeyService_ServiceDesc is the grpc.ServiceDesc for KeyService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var KeyService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "services.KeyService",
	HandlerType: (*KeyServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "VerifyMessageAttestation",
			Handler:    _KeyService_VerifyMessageAttestation_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "backend_services.proto",
}

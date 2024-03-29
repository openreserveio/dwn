// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
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

const (
	RecordService_FindRecord_FullMethodName       = "/services.RecordService/FindRecord"
	RecordService_CreateSchema_FullMethodName     = "/services.RecordService/CreateSchema"
	RecordService_ValidateRecord_FullMethodName   = "/services.RecordService/ValidateRecord"
	RecordService_InvalidateSchema_FullMethodName = "/services.RecordService/InvalidateSchema"
	RecordService_Query_FullMethodName            = "/services.RecordService/Query"
	RecordService_Write_FullMethodName            = "/services.RecordService/Write"
	RecordService_Commit_FullMethodName           = "/services.RecordService/Commit"
	RecordService_Delete_FullMethodName           = "/services.RecordService/Delete"
)

// RecordServiceClient is the client API for RecordService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RecordServiceClient interface {
	FindRecord(ctx context.Context, in *FindRecordRequest, opts ...grpc.CallOption) (*FindRecordResponse, error)
	CreateSchema(ctx context.Context, in *CreateSchemaRequest, opts ...grpc.CallOption) (*CreateSchemaResponse, error)
	ValidateRecord(ctx context.Context, in *ValidateRecordRequest, opts ...grpc.CallOption) (*ValidateRecordResponse, error)
	InvalidateSchema(ctx context.Context, in *InvalidateSchemaRequest, opts ...grpc.CallOption) (*InvalidateSchemaResponse, error)
	// Net new go simplify the service
	Query(ctx context.Context, in *QueryRecordRequest, opts ...grpc.CallOption) (*QueryRecordResponse, error)
	Write(ctx context.Context, in *WriteRecordRequest, opts ...grpc.CallOption) (*WriteRecordResponse, error)
	Commit(ctx context.Context, in *CommitRecordRequest, opts ...grpc.CallOption) (*CommitRecordResponse, error)
	Delete(ctx context.Context, in *DeleteRecordRequest, opts ...grpc.CallOption) (*DeleteRecordResponse, error)
}

type recordServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRecordServiceClient(cc grpc.ClientConnInterface) RecordServiceClient {
	return &recordServiceClient{cc}
}

func (c *recordServiceClient) FindRecord(ctx context.Context, in *FindRecordRequest, opts ...grpc.CallOption) (*FindRecordResponse, error) {
	out := new(FindRecordResponse)
	err := c.cc.Invoke(ctx, RecordService_FindRecord_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *recordServiceClient) CreateSchema(ctx context.Context, in *CreateSchemaRequest, opts ...grpc.CallOption) (*CreateSchemaResponse, error) {
	out := new(CreateSchemaResponse)
	err := c.cc.Invoke(ctx, RecordService_CreateSchema_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *recordServiceClient) ValidateRecord(ctx context.Context, in *ValidateRecordRequest, opts ...grpc.CallOption) (*ValidateRecordResponse, error) {
	out := new(ValidateRecordResponse)
	err := c.cc.Invoke(ctx, RecordService_ValidateRecord_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *recordServiceClient) InvalidateSchema(ctx context.Context, in *InvalidateSchemaRequest, opts ...grpc.CallOption) (*InvalidateSchemaResponse, error) {
	out := new(InvalidateSchemaResponse)
	err := c.cc.Invoke(ctx, RecordService_InvalidateSchema_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *recordServiceClient) Query(ctx context.Context, in *QueryRecordRequest, opts ...grpc.CallOption) (*QueryRecordResponse, error) {
	out := new(QueryRecordResponse)
	err := c.cc.Invoke(ctx, RecordService_Query_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *recordServiceClient) Write(ctx context.Context, in *WriteRecordRequest, opts ...grpc.CallOption) (*WriteRecordResponse, error) {
	out := new(WriteRecordResponse)
	err := c.cc.Invoke(ctx, RecordService_Write_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *recordServiceClient) Commit(ctx context.Context, in *CommitRecordRequest, opts ...grpc.CallOption) (*CommitRecordResponse, error) {
	out := new(CommitRecordResponse)
	err := c.cc.Invoke(ctx, RecordService_Commit_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *recordServiceClient) Delete(ctx context.Context, in *DeleteRecordRequest, opts ...grpc.CallOption) (*DeleteRecordResponse, error) {
	out := new(DeleteRecordResponse)
	err := c.cc.Invoke(ctx, RecordService_Delete_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RecordServiceServer is the server API for RecordService service.
// All implementations must embed UnimplementedRecordServiceServer
// for forward compatibility
type RecordServiceServer interface {
	FindRecord(context.Context, *FindRecordRequest) (*FindRecordResponse, error)
	CreateSchema(context.Context, *CreateSchemaRequest) (*CreateSchemaResponse, error)
	ValidateRecord(context.Context, *ValidateRecordRequest) (*ValidateRecordResponse, error)
	InvalidateSchema(context.Context, *InvalidateSchemaRequest) (*InvalidateSchemaResponse, error)
	// Net new go simplify the service
	Query(context.Context, *QueryRecordRequest) (*QueryRecordResponse, error)
	Write(context.Context, *WriteRecordRequest) (*WriteRecordResponse, error)
	Commit(context.Context, *CommitRecordRequest) (*CommitRecordResponse, error)
	Delete(context.Context, *DeleteRecordRequest) (*DeleteRecordResponse, error)
	mustEmbedUnimplementedRecordServiceServer()
}

// UnimplementedRecordServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRecordServiceServer struct {
}

func (UnimplementedRecordServiceServer) FindRecord(context.Context, *FindRecordRequest) (*FindRecordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindRecord not implemented")
}
func (UnimplementedRecordServiceServer) CreateSchema(context.Context, *CreateSchemaRequest) (*CreateSchemaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSchema not implemented")
}
func (UnimplementedRecordServiceServer) ValidateRecord(context.Context, *ValidateRecordRequest) (*ValidateRecordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidateRecord not implemented")
}
func (UnimplementedRecordServiceServer) InvalidateSchema(context.Context, *InvalidateSchemaRequest) (*InvalidateSchemaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InvalidateSchema not implemented")
}
func (UnimplementedRecordServiceServer) Query(context.Context, *QueryRecordRequest) (*QueryRecordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Query not implemented")
}
func (UnimplementedRecordServiceServer) Write(context.Context, *WriteRecordRequest) (*WriteRecordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Write not implemented")
}
func (UnimplementedRecordServiceServer) Commit(context.Context, *CommitRecordRequest) (*CommitRecordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Commit not implemented")
}
func (UnimplementedRecordServiceServer) Delete(context.Context, *DeleteRecordRequest) (*DeleteRecordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedRecordServiceServer) mustEmbedUnimplementedRecordServiceServer() {}

// UnsafeRecordServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RecordServiceServer will
// result in compilation errors.
type UnsafeRecordServiceServer interface {
	mustEmbedUnimplementedRecordServiceServer()
}

func RegisterRecordServiceServer(s grpc.ServiceRegistrar, srv RecordServiceServer) {
	s.RegisterService(&RecordService_ServiceDesc, srv)
}

func _RecordService_FindRecord_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindRecordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecordServiceServer).FindRecord(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RecordService_FindRecord_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecordServiceServer).FindRecord(ctx, req.(*FindRecordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RecordService_CreateSchema_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSchemaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecordServiceServer).CreateSchema(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RecordService_CreateSchema_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecordServiceServer).CreateSchema(ctx, req.(*CreateSchemaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RecordService_ValidateRecord_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValidateRecordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecordServiceServer).ValidateRecord(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RecordService_ValidateRecord_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecordServiceServer).ValidateRecord(ctx, req.(*ValidateRecordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RecordService_InvalidateSchema_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InvalidateSchemaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecordServiceServer).InvalidateSchema(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RecordService_InvalidateSchema_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecordServiceServer).InvalidateSchema(ctx, req.(*InvalidateSchemaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RecordService_Query_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryRecordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecordServiceServer).Query(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RecordService_Query_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecordServiceServer).Query(ctx, req.(*QueryRecordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RecordService_Write_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WriteRecordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecordServiceServer).Write(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RecordService_Write_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecordServiceServer).Write(ctx, req.(*WriteRecordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RecordService_Commit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommitRecordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecordServiceServer).Commit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RecordService_Commit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecordServiceServer).Commit(ctx, req.(*CommitRecordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RecordService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRecordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecordServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RecordService_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecordServiceServer).Delete(ctx, req.(*DeleteRecordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RecordService_ServiceDesc is the grpc.ServiceDesc for RecordService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RecordService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "services.RecordService",
	HandlerType: (*RecordServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FindRecord",
			Handler:    _RecordService_FindRecord_Handler,
		},
		{
			MethodName: "CreateSchema",
			Handler:    _RecordService_CreateSchema_Handler,
		},
		{
			MethodName: "ValidateRecord",
			Handler:    _RecordService_ValidateRecord_Handler,
		},
		{
			MethodName: "InvalidateSchema",
			Handler:    _RecordService_InvalidateSchema_Handler,
		},
		{
			MethodName: "Query",
			Handler:    _RecordService_Query_Handler,
		},
		{
			MethodName: "Write",
			Handler:    _RecordService_Write_Handler,
		},
		{
			MethodName: "Commit",
			Handler:    _RecordService_Commit_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _RecordService_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "backend_services.proto",
}

const (
	HookService_RegisterHook_FullMethodName             = "/services.HookService/RegisterHook"
	HookService_UpdateHook_FullMethodName               = "/services.HookService/UpdateHook"
	HookService_GetHookByRecordId_FullMethodName        = "/services.HookService/GetHookByRecordId"
	HookService_GetHooksForRecord_FullMethodName        = "/services.HookService/GetHooksForRecord"
	HookService_NotifyHooksOfRecordEvent_FullMethodName = "/services.HookService/NotifyHooksOfRecordEvent"
)

// HookServiceClient is the client API for HookService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HookServiceClient interface {
	RegisterHook(ctx context.Context, in *RegisterHookRequest, opts ...grpc.CallOption) (*RegisterHookResponse, error)
	UpdateHook(ctx context.Context, in *UpdateHookRequest, opts ...grpc.CallOption) (*UpdateHookResponse, error)
	GetHookByRecordId(ctx context.Context, in *GetHookByRecordIdRequest, opts ...grpc.CallOption) (*GetHookByRecordIdResponse, error)
	GetHooksForRecord(ctx context.Context, in *GetHooksForRecordRequest, opts ...grpc.CallOption) (*GetHooksForRecordResponse, error)
	NotifyHooksOfRecordEvent(ctx context.Context, in *NotifyHooksOfRecordEventRequest, opts ...grpc.CallOption) (*NotifyHooksOfRecordEventResponse, error)
}

type hookServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewHookServiceClient(cc grpc.ClientConnInterface) HookServiceClient {
	return &hookServiceClient{cc}
}

func (c *hookServiceClient) RegisterHook(ctx context.Context, in *RegisterHookRequest, opts ...grpc.CallOption) (*RegisterHookResponse, error) {
	out := new(RegisterHookResponse)
	err := c.cc.Invoke(ctx, HookService_RegisterHook_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hookServiceClient) UpdateHook(ctx context.Context, in *UpdateHookRequest, opts ...grpc.CallOption) (*UpdateHookResponse, error) {
	out := new(UpdateHookResponse)
	err := c.cc.Invoke(ctx, HookService_UpdateHook_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hookServiceClient) GetHookByRecordId(ctx context.Context, in *GetHookByRecordIdRequest, opts ...grpc.CallOption) (*GetHookByRecordIdResponse, error) {
	out := new(GetHookByRecordIdResponse)
	err := c.cc.Invoke(ctx, HookService_GetHookByRecordId_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hookServiceClient) GetHooksForRecord(ctx context.Context, in *GetHooksForRecordRequest, opts ...grpc.CallOption) (*GetHooksForRecordResponse, error) {
	out := new(GetHooksForRecordResponse)
	err := c.cc.Invoke(ctx, HookService_GetHooksForRecord_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *hookServiceClient) NotifyHooksOfRecordEvent(ctx context.Context, in *NotifyHooksOfRecordEventRequest, opts ...grpc.CallOption) (*NotifyHooksOfRecordEventResponse, error) {
	out := new(NotifyHooksOfRecordEventResponse)
	err := c.cc.Invoke(ctx, HookService_NotifyHooksOfRecordEvent_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HookServiceServer is the server API for HookService service.
// All implementations must embed UnimplementedHookServiceServer
// for forward compatibility
type HookServiceServer interface {
	RegisterHook(context.Context, *RegisterHookRequest) (*RegisterHookResponse, error)
	UpdateHook(context.Context, *UpdateHookRequest) (*UpdateHookResponse, error)
	GetHookByRecordId(context.Context, *GetHookByRecordIdRequest) (*GetHookByRecordIdResponse, error)
	GetHooksForRecord(context.Context, *GetHooksForRecordRequest) (*GetHooksForRecordResponse, error)
	NotifyHooksOfRecordEvent(context.Context, *NotifyHooksOfRecordEventRequest) (*NotifyHooksOfRecordEventResponse, error)
	mustEmbedUnimplementedHookServiceServer()
}

// UnimplementedHookServiceServer must be embedded to have forward compatible implementations.
type UnimplementedHookServiceServer struct {
}

func (UnimplementedHookServiceServer) RegisterHook(context.Context, *RegisterHookRequest) (*RegisterHookResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterHook not implemented")
}
func (UnimplementedHookServiceServer) UpdateHook(context.Context, *UpdateHookRequest) (*UpdateHookResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateHook not implemented")
}
func (UnimplementedHookServiceServer) GetHookByRecordId(context.Context, *GetHookByRecordIdRequest) (*GetHookByRecordIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetHookByRecordId not implemented")
}
func (UnimplementedHookServiceServer) GetHooksForRecord(context.Context, *GetHooksForRecordRequest) (*GetHooksForRecordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetHooksForRecord not implemented")
}
func (UnimplementedHookServiceServer) NotifyHooksOfRecordEvent(context.Context, *NotifyHooksOfRecordEventRequest) (*NotifyHooksOfRecordEventResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NotifyHooksOfRecordEvent not implemented")
}
func (UnimplementedHookServiceServer) mustEmbedUnimplementedHookServiceServer() {}

// UnsafeHookServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HookServiceServer will
// result in compilation errors.
type UnsafeHookServiceServer interface {
	mustEmbedUnimplementedHookServiceServer()
}

func RegisterHookServiceServer(s grpc.ServiceRegistrar, srv HookServiceServer) {
	s.RegisterService(&HookService_ServiceDesc, srv)
}

func _HookService_RegisterHook_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterHookRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HookServiceServer).RegisterHook(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: HookService_RegisterHook_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HookServiceServer).RegisterHook(ctx, req.(*RegisterHookRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HookService_UpdateHook_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateHookRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HookServiceServer).UpdateHook(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: HookService_UpdateHook_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HookServiceServer).UpdateHook(ctx, req.(*UpdateHookRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HookService_GetHookByRecordId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetHookByRecordIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HookServiceServer).GetHookByRecordId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: HookService_GetHookByRecordId_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HookServiceServer).GetHookByRecordId(ctx, req.(*GetHookByRecordIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HookService_GetHooksForRecord_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetHooksForRecordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HookServiceServer).GetHooksForRecord(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: HookService_GetHooksForRecord_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HookServiceServer).GetHooksForRecord(ctx, req.(*GetHooksForRecordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _HookService_NotifyHooksOfRecordEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NotifyHooksOfRecordEventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HookServiceServer).NotifyHooksOfRecordEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: HookService_NotifyHooksOfRecordEvent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HookServiceServer).NotifyHooksOfRecordEvent(ctx, req.(*NotifyHooksOfRecordEventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// HookService_ServiceDesc is the grpc.ServiceDesc for HookService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var HookService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "services.HookService",
	HandlerType: (*HookServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterHook",
			Handler:    _HookService_RegisterHook_Handler,
		},
		{
			MethodName: "UpdateHook",
			Handler:    _HookService_UpdateHook_Handler,
		},
		{
			MethodName: "GetHookByRecordId",
			Handler:    _HookService_GetHookByRecordId_Handler,
		},
		{
			MethodName: "GetHooksForRecord",
			Handler:    _HookService_GetHooksForRecord_Handler,
		},
		{
			MethodName: "NotifyHooksOfRecordEvent",
			Handler:    _HookService_NotifyHooksOfRecordEvent_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "backend_services.proto",
}

const (
	KeyService_VerifyMessageAttestation_FullMethodName = "/services.KeyService/VerifyMessageAttestation"
)

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
	err := c.cc.Invoke(ctx, KeyService_VerifyMessageAttestation_FullMethodName, in, out, opts...)
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
		FullMethod: KeyService_VerifyMessageAttestation_FullMethodName,
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

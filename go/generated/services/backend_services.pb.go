// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.9
// source: backend_services.proto

package services

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Status int32

const (
	Status_OK        Status = 0
	Status_NOT_FOUND Status = 1
	Status_ERROR     Status = 2
	Status_CONFLICT  Status = 3
)

// Enum value maps for Status.
var (
	Status_name = map[int32]string{
		0: "OK",
		1: "NOT_FOUND",
		2: "ERROR",
		3: "CONFLICT",
	}
	Status_value = map[string]int32{
		"OK":        0,
		"NOT_FOUND": 1,
		"ERROR":     2,
		"CONFLICT":  3,
	}
)

func (x Status) Enum() *Status {
	p := new(Status)
	*p = x
	return p
}

func (x Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Status) Descriptor() protoreflect.EnumDescriptor {
	return file_backend_services_proto_enumTypes[0].Descriptor()
}

func (Status) Type() protoreflect.EnumType {
	return &file_backend_services_proto_enumTypes[0]
}

func (x Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Status.Descriptor instead.
func (Status) EnumDescriptor() ([]byte, []int) {
	return file_backend_services_proto_rawDescGZIP(), []int{0}
}

type QueryType int32

const (
	QueryType_SINGLE_COLLECTION_BY_ID_SCHEMA_URI QueryType = 0
)

// Enum value maps for QueryType.
var (
	QueryType_name = map[int32]string{
		0: "SINGLE_COLLECTION_BY_ID_SCHEMA_URI",
	}
	QueryType_value = map[string]int32{
		"SINGLE_COLLECTION_BY_ID_SCHEMA_URI": 0,
	}
)

func (x QueryType) Enum() *QueryType {
	p := new(QueryType)
	*p = x
	return p
}

func (x QueryType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (QueryType) Descriptor() protoreflect.EnumDescriptor {
	return file_backend_services_proto_enumTypes[1].Descriptor()
}

func (QueryType) Type() protoreflect.EnumType {
	return &file_backend_services_proto_enumTypes[1]
}

func (x QueryType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use QueryType.Descriptor instead.
func (QueryType) EnumDescriptor() ([]byte, []int) {
	return file_backend_services_proto_rawDescGZIP(), []int{1}
}

type CommonStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status  Status `protobuf:"varint,1,opt,name=status,proto3,enum=services.Status" json:"status,omitempty"`
	Details string `protobuf:"bytes,2,opt,name=details,proto3" json:"details,omitempty"`
}

func (x *CommonStatus) Reset() {
	*x = CommonStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_backend_services_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommonStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommonStatus) ProtoMessage() {}

func (x *CommonStatus) ProtoReflect() protoreflect.Message {
	mi := &file_backend_services_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommonStatus.ProtoReflect.Descriptor instead.
func (*CommonStatus) Descriptor() ([]byte, []int) {
	return file_backend_services_proto_rawDescGZIP(), []int{0}
}

func (x *CommonStatus) GetStatus() Status {
	if x != nil {
		return x.Status
	}
	return Status_OK
}

func (x *CommonStatus) GetDetails() string {
	if x != nil {
		return x.Details
	}
	return ""
}

type StoreCollectionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *StoreCollectionRequest) Reset() {
	*x = StoreCollectionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_backend_services_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StoreCollectionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoreCollectionRequest) ProtoMessage() {}

func (x *StoreCollectionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_backend_services_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StoreCollectionRequest.ProtoReflect.Descriptor instead.
func (*StoreCollectionRequest) Descriptor() ([]byte, []int) {
	return file_backend_services_proto_rawDescGZIP(), []int{1}
}

type StoreCollectionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status *CommonStatus `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *StoreCollectionResponse) Reset() {
	*x = StoreCollectionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_backend_services_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StoreCollectionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoreCollectionResponse) ProtoMessage() {}

func (x *StoreCollectionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_backend_services_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StoreCollectionResponse.ProtoReflect.Descriptor instead.
func (*StoreCollectionResponse) Descriptor() ([]byte, []int) {
	return file_backend_services_proto_rawDescGZIP(), []int{2}
}

func (x *StoreCollectionResponse) GetStatus() *CommonStatus {
	if x != nil {
		return x.Status
	}
	return nil
}

type FindCollectionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	QueryType        QueryType `protobuf:"varint,1,opt,name=queryType,proto3,enum=services.QueryType" json:"queryType,omitempty"`
	CollectionItemId string    `protobuf:"bytes,2,opt,name=collectionItemId,proto3" json:"collectionItemId,omitempty"`
	SchemaURI        string    `protobuf:"bytes,3,opt,name=schemaURI,proto3" json:"schemaURI,omitempty"`
}

func (x *FindCollectionRequest) Reset() {
	*x = FindCollectionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_backend_services_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindCollectionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindCollectionRequest) ProtoMessage() {}

func (x *FindCollectionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_backend_services_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindCollectionRequest.ProtoReflect.Descriptor instead.
func (*FindCollectionRequest) Descriptor() ([]byte, []int) {
	return file_backend_services_proto_rawDescGZIP(), []int{3}
}

func (x *FindCollectionRequest) GetQueryType() QueryType {
	if x != nil {
		return x.QueryType
	}
	return QueryType_SINGLE_COLLECTION_BY_ID_SCHEMA_URI
}

func (x *FindCollectionRequest) GetCollectionItemId() string {
	if x != nil {
		return x.CollectionItemId
	}
	return ""
}

func (x *FindCollectionRequest) GetSchemaURI() string {
	if x != nil {
		return x.SchemaURI
	}
	return ""
}

type FindCollectionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status *CommonStatus `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *FindCollectionResponse) Reset() {
	*x = FindCollectionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_backend_services_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindCollectionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindCollectionResponse) ProtoMessage() {}

func (x *FindCollectionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_backend_services_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindCollectionResponse.ProtoReflect.Descriptor instead.
func (*FindCollectionResponse) Descriptor() ([]byte, []int) {
	return file_backend_services_proto_rawDescGZIP(), []int{4}
}

func (x *FindCollectionResponse) GetStatus() *CommonStatus {
	if x != nil {
		return x.Status
	}
	return nil
}

type CreateSchemaRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CreateSchemaRequest) Reset() {
	*x = CreateSchemaRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_backend_services_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateSchemaRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateSchemaRequest) ProtoMessage() {}

func (x *CreateSchemaRequest) ProtoReflect() protoreflect.Message {
	mi := &file_backend_services_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateSchemaRequest.ProtoReflect.Descriptor instead.
func (*CreateSchemaRequest) Descriptor() ([]byte, []int) {
	return file_backend_services_proto_rawDescGZIP(), []int{5}
}

type CreateSchemaResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status *CommonStatus `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *CreateSchemaResponse) Reset() {
	*x = CreateSchemaResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_backend_services_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateSchemaResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateSchemaResponse) ProtoMessage() {}

func (x *CreateSchemaResponse) ProtoReflect() protoreflect.Message {
	mi := &file_backend_services_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateSchemaResponse.ProtoReflect.Descriptor instead.
func (*CreateSchemaResponse) Descriptor() ([]byte, []int) {
	return file_backend_services_proto_rawDescGZIP(), []int{6}
}

func (x *CreateSchemaResponse) GetStatus() *CommonStatus {
	if x != nil {
		return x.Status
	}
	return nil
}

type ValidateCollectionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SchemaURI string `protobuf:"bytes,1,opt,name=schemaURI,proto3" json:"schemaURI,omitempty"`
	Document  []byte `protobuf:"bytes,2,opt,name=document,proto3" json:"document,omitempty"`
}

func (x *ValidateCollectionRequest) Reset() {
	*x = ValidateCollectionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_backend_services_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ValidateCollectionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValidateCollectionRequest) ProtoMessage() {}

func (x *ValidateCollectionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_backend_services_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValidateCollectionRequest.ProtoReflect.Descriptor instead.
func (*ValidateCollectionRequest) Descriptor() ([]byte, []int) {
	return file_backend_services_proto_rawDescGZIP(), []int{7}
}

func (x *ValidateCollectionRequest) GetSchemaURI() string {
	if x != nil {
		return x.SchemaURI
	}
	return ""
}

func (x *ValidateCollectionRequest) GetDocument() []byte {
	if x != nil {
		return x.Document
	}
	return nil
}

type ValidateCollectionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status *CommonStatus `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *ValidateCollectionResponse) Reset() {
	*x = ValidateCollectionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_backend_services_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ValidateCollectionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValidateCollectionResponse) ProtoMessage() {}

func (x *ValidateCollectionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_backend_services_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValidateCollectionResponse.ProtoReflect.Descriptor instead.
func (*ValidateCollectionResponse) Descriptor() ([]byte, []int) {
	return file_backend_services_proto_rawDescGZIP(), []int{8}
}

func (x *ValidateCollectionResponse) GetStatus() *CommonStatus {
	if x != nil {
		return x.Status
	}
	return nil
}

type InvalidateSchemaRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SchemaURI string `protobuf:"bytes,1,opt,name=schemaURI,proto3" json:"schemaURI,omitempty"`
}

func (x *InvalidateSchemaRequest) Reset() {
	*x = InvalidateSchemaRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_backend_services_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InvalidateSchemaRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InvalidateSchemaRequest) ProtoMessage() {}

func (x *InvalidateSchemaRequest) ProtoReflect() protoreflect.Message {
	mi := &file_backend_services_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InvalidateSchemaRequest.ProtoReflect.Descriptor instead.
func (*InvalidateSchemaRequest) Descriptor() ([]byte, []int) {
	return file_backend_services_proto_rawDescGZIP(), []int{9}
}

func (x *InvalidateSchemaRequest) GetSchemaURI() string {
	if x != nil {
		return x.SchemaURI
	}
	return ""
}

type InvalidateSchemaResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status *CommonStatus `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *InvalidateSchemaResponse) Reset() {
	*x = InvalidateSchemaResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_backend_services_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InvalidateSchemaResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InvalidateSchemaResponse) ProtoMessage() {}

func (x *InvalidateSchemaResponse) ProtoReflect() protoreflect.Message {
	mi := &file_backend_services_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InvalidateSchemaResponse.ProtoReflect.Descriptor instead.
func (*InvalidateSchemaResponse) Descriptor() ([]byte, []int) {
	return file_backend_services_proto_rawDescGZIP(), []int{10}
}

func (x *InvalidateSchemaResponse) GetStatus() *CommonStatus {
	if x != nil {
		return x.Status
	}
	return nil
}

var File_backend_services_proto protoreflect.FileDescriptor

var file_backend_services_proto_rawDesc = []byte{
	0x0a, 0x16, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x73, 0x22, 0x52, 0x0a, 0x0c, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x28, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x10, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x18, 0x0a, 0x07,
	0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x64,
	0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x22, 0x18, 0x0a, 0x16, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x43,
	0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x22, 0x49, 0x0a, 0x17, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x94, 0x01, 0x0a, 0x15,
	0x46, 0x69, 0x6e, 0x64, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x31, 0x0a, 0x09, 0x71, 0x75, 0x65, 0x72, 0x79, 0x54, 0x79,
	0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x73, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x54, 0x79, 0x70, 0x65, 0x52, 0x09, 0x71,
	0x75, 0x65, 0x72, 0x79, 0x54, 0x79, 0x70, 0x65, 0x12, 0x2a, 0x0a, 0x10, 0x63, 0x6f, 0x6c, 0x6c,
	0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x74, 0x65, 0x6d, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x10, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x74,
	0x65, 0x6d, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x55, 0x52,
	0x49, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x55,
	0x52, 0x49, 0x22, 0x48, 0x0a, 0x16, 0x46, 0x69, 0x6e, 0x64, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x15, 0x0a, 0x13,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x22, 0x46, 0x0a, 0x14, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x63, 0x68,
	0x65, 0x6d, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x55, 0x0a, 0x19, 0x56,
	0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x63, 0x68, 0x65,
	0x6d, 0x61, 0x55, 0x52, 0x49, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x63, 0x68,
	0x65, 0x6d, 0x61, 0x55, 0x52, 0x49, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65,
	0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65,
	0x6e, 0x74, 0x22, 0x4c, 0x0a, 0x1a, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6f,
	0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x2e, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x16, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x43, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x22, 0x37, 0x0a, 0x17, 0x49, 0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x53, 0x63,
	0x68, 0x65, 0x6d, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x73,
	0x63, 0x68, 0x65, 0x6d, 0x61, 0x55, 0x52, 0x49, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x55, 0x52, 0x49, 0x22, 0x4a, 0x0a, 0x18, 0x49, 0x6e, 0x76,
	0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73,
	0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x2a, 0x38, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12,
	0x06, 0x0a, 0x02, 0x4f, 0x4b, 0x10, 0x00, 0x12, 0x0d, 0x0a, 0x09, 0x4e, 0x4f, 0x54, 0x5f, 0x46,
	0x4f, 0x55, 0x4e, 0x44, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10,
	0x02, 0x12, 0x0c, 0x0a, 0x08, 0x43, 0x4f, 0x4e, 0x46, 0x4c, 0x49, 0x43, 0x54, 0x10, 0x03, 0x2a,
	0x33, 0x0a, 0x09, 0x51, 0x75, 0x65, 0x72, 0x79, 0x54, 0x79, 0x70, 0x65, 0x12, 0x26, 0x0a, 0x22,
	0x53, 0x49, 0x4e, 0x47, 0x4c, 0x45, 0x5f, 0x43, 0x4f, 0x4c, 0x4c, 0x45, 0x43, 0x54, 0x49, 0x4f,
	0x4e, 0x5f, 0x42, 0x59, 0x5f, 0x49, 0x44, 0x5f, 0x53, 0x43, 0x48, 0x45, 0x4d, 0x41, 0x5f, 0x55,
	0x52, 0x49, 0x10, 0x00, 0x32, 0xd5, 0x03, 0x0a, 0x11, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x58, 0x0a, 0x0f, 0x53, 0x74,
	0x6f, 0x72, 0x65, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x20, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x43, 0x6f,
	0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x21, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x53, 0x74, 0x6f, 0x72, 0x65,
	0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x12, 0x55, 0x0a, 0x0e, 0x46, 0x69, 0x6e, 0x64, 0x43, 0x6f, 0x6c, 0x6c,
	0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1f, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x73, 0x2e, 0x46, 0x69, 0x6e, 0x64, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x73, 0x2e, 0x46, 0x69, 0x6e, 0x64, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4f, 0x0a, 0x0c, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x12, 0x1d, 0x2e, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x63, 0x68,
	0x65, 0x6d, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x53, 0x63, 0x68, 0x65,
	0x6d, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x61, 0x0a, 0x12,
	0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x23, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x56, 0x61,
	0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x73, 0x2e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6c, 0x6c, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12,
	0x5b, 0x0a, 0x10, 0x49, 0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x53, 0x63, 0x68,
	0x65, 0x6d, 0x61, 0x12, 0x21, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x49,
	0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x73, 0x2e, 0x49, 0x6e, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x53, 0x63, 0x68, 0x65,
	0x6d, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x34, 0x5a, 0x32,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x72,
	0x65, 0x73, 0x65, 0x72, 0x76, 0x65, 0x69, 0x6f, 0x2f, 0x64, 0x77, 0x6e, 0x2f, 0x67, 0x6f, 0x2f,
	0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_backend_services_proto_rawDescOnce sync.Once
	file_backend_services_proto_rawDescData = file_backend_services_proto_rawDesc
)

func file_backend_services_proto_rawDescGZIP() []byte {
	file_backend_services_proto_rawDescOnce.Do(func() {
		file_backend_services_proto_rawDescData = protoimpl.X.CompressGZIP(file_backend_services_proto_rawDescData)
	})
	return file_backend_services_proto_rawDescData
}

var file_backend_services_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_backend_services_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_backend_services_proto_goTypes = []interface{}{
	(Status)(0),                        // 0: services.Status
	(QueryType)(0),                     // 1: services.QueryType
	(*CommonStatus)(nil),               // 2: services.CommonStatus
	(*StoreCollectionRequest)(nil),     // 3: services.StoreCollectionRequest
	(*StoreCollectionResponse)(nil),    // 4: services.StoreCollectionResponse
	(*FindCollectionRequest)(nil),      // 5: services.FindCollectionRequest
	(*FindCollectionResponse)(nil),     // 6: services.FindCollectionResponse
	(*CreateSchemaRequest)(nil),        // 7: services.CreateSchemaRequest
	(*CreateSchemaResponse)(nil),       // 8: services.CreateSchemaResponse
	(*ValidateCollectionRequest)(nil),  // 9: services.ValidateCollectionRequest
	(*ValidateCollectionResponse)(nil), // 10: services.ValidateCollectionResponse
	(*InvalidateSchemaRequest)(nil),    // 11: services.InvalidateSchemaRequest
	(*InvalidateSchemaResponse)(nil),   // 12: services.InvalidateSchemaResponse
}
var file_backend_services_proto_depIdxs = []int32{
	0,  // 0: services.CommonStatus.status:type_name -> services.Status
	2,  // 1: services.StoreCollectionResponse.status:type_name -> services.CommonStatus
	1,  // 2: services.FindCollectionRequest.queryType:type_name -> services.QueryType
	2,  // 3: services.FindCollectionResponse.status:type_name -> services.CommonStatus
	2,  // 4: services.CreateSchemaResponse.status:type_name -> services.CommonStatus
	2,  // 5: services.ValidateCollectionResponse.status:type_name -> services.CommonStatus
	2,  // 6: services.InvalidateSchemaResponse.status:type_name -> services.CommonStatus
	3,  // 7: services.CollectionService.StoreCollection:input_type -> services.StoreCollectionRequest
	5,  // 8: services.CollectionService.FindCollection:input_type -> services.FindCollectionRequest
	7,  // 9: services.CollectionService.CreateSchema:input_type -> services.CreateSchemaRequest
	9,  // 10: services.CollectionService.ValidateCollection:input_type -> services.ValidateCollectionRequest
	11, // 11: services.CollectionService.InvalidateSchema:input_type -> services.InvalidateSchemaRequest
	4,  // 12: services.CollectionService.StoreCollection:output_type -> services.StoreCollectionResponse
	6,  // 13: services.CollectionService.FindCollection:output_type -> services.FindCollectionResponse
	8,  // 14: services.CollectionService.CreateSchema:output_type -> services.CreateSchemaResponse
	10, // 15: services.CollectionService.ValidateCollection:output_type -> services.ValidateCollectionResponse
	12, // 16: services.CollectionService.InvalidateSchema:output_type -> services.InvalidateSchemaResponse
	12, // [12:17] is the sub-list for method output_type
	7,  // [7:12] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_backend_services_proto_init() }
func file_backend_services_proto_init() {
	if File_backend_services_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_backend_services_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CommonStatus); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_backend_services_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StoreCollectionRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_backend_services_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StoreCollectionResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_backend_services_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FindCollectionRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_backend_services_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FindCollectionResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_backend_services_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateSchemaRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_backend_services_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateSchemaResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_backend_services_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ValidateCollectionRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_backend_services_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ValidateCollectionResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_backend_services_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InvalidateSchemaRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_backend_services_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InvalidateSchemaResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_backend_services_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_backend_services_proto_goTypes,
		DependencyIndexes: file_backend_services_proto_depIdxs,
		EnumInfos:         file_backend_services_proto_enumTypes,
		MessageInfos:      file_backend_services_proto_msgTypes,
	}.Build()
	File_backend_services_proto = out.File
	file_backend_services_proto_rawDesc = nil
	file_backend_services_proto_goTypes = nil
	file_backend_services_proto_depIdxs = nil
}

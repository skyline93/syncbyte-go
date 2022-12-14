// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.6
// source: internal/proto/syncbyte.proto

package proto

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

type BackupRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Datasetname string          `protobuf:"bytes,1,opt,name=datasetname,proto3" json:"datasetname,omitempty"`
	Iscompress  bool            `protobuf:"varint,2,opt,name=iscompress,proto3" json:"iscompress,omitempty"`
	SourceOpts  *SourceOptions  `protobuf:"bytes,3,opt,name=sourceOpts,proto3" json:"sourceOpts,omitempty"`
	BackendOpts *BackendOptions `protobuf:"bytes,4,opt,name=backendOpts,proto3" json:"backendOpts,omitempty"`
}

func (x *BackupRequest) Reset() {
	*x = BackupRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_syncbyte_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BackupRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BackupRequest) ProtoMessage() {}

func (x *BackupRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_syncbyte_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BackupRequest.ProtoReflect.Descriptor instead.
func (*BackupRequest) Descriptor() ([]byte, []int) {
	return file_internal_proto_syncbyte_proto_rawDescGZIP(), []int{0}
}

func (x *BackupRequest) GetDatasetname() string {
	if x != nil {
		return x.Datasetname
	}
	return ""
}

func (x *BackupRequest) GetIscompress() bool {
	if x != nil {
		return x.Iscompress
	}
	return false
}

func (x *BackupRequest) GetSourceOpts() *SourceOptions {
	if x != nil {
		return x.SourceOpts
	}
	return nil
}

func (x *BackupRequest) GetBackendOpts() *BackendOptions {
	if x != nil {
		return x.BackendOpts
	}
	return nil
}

type BackupResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Jobid string `protobuf:"bytes,1,opt,name=jobid,proto3" json:"jobid,omitempty"`
}

func (x *BackupResponse) Reset() {
	*x = BackupResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_syncbyte_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BackupResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BackupResponse) ProtoMessage() {}

func (x *BackupResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_syncbyte_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BackupResponse.ProtoReflect.Descriptor instead.
func (*BackupResponse) Descriptor() ([]byte, []int) {
	return file_internal_proto_syncbyte_proto_rawDescGZIP(), []int{1}
}

func (x *BackupResponse) GetJobid() string {
	if x != nil {
		return x.Jobid
	}
	return ""
}

type RestoreRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Datasetname  string          `protobuf:"bytes,1,opt,name=datasetname,proto3" json:"datasetname,omitempty"`
	Isuncompress bool            `protobuf:"varint,2,opt,name=isuncompress,proto3" json:"isuncompress,omitempty"`
	DestOpts     *DestOptions    `protobuf:"bytes,3,opt,name=destOpts,proto3" json:"destOpts,omitempty"`
	BackendOpts  *BackendOptions `protobuf:"bytes,4,opt,name=backendOpts,proto3" json:"backendOpts,omitempty"`
}

func (x *RestoreRequest) Reset() {
	*x = RestoreRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_syncbyte_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RestoreRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RestoreRequest) ProtoMessage() {}

func (x *RestoreRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_syncbyte_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RestoreRequest.ProtoReflect.Descriptor instead.
func (*RestoreRequest) Descriptor() ([]byte, []int) {
	return file_internal_proto_syncbyte_proto_rawDescGZIP(), []int{2}
}

func (x *RestoreRequest) GetDatasetname() string {
	if x != nil {
		return x.Datasetname
	}
	return ""
}

func (x *RestoreRequest) GetIsuncompress() bool {
	if x != nil {
		return x.Isuncompress
	}
	return false
}

func (x *RestoreRequest) GetDestOpts() *DestOptions {
	if x != nil {
		return x.DestOpts
	}
	return nil
}

func (x *RestoreRequest) GetBackendOpts() *BackendOptions {
	if x != nil {
		return x.BackendOpts
	}
	return nil
}

type RestoreResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Jobid string `protobuf:"bytes,1,opt,name=jobid,proto3" json:"jobid,omitempty"`
}

func (x *RestoreResponse) Reset() {
	*x = RestoreResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_syncbyte_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RestoreResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RestoreResponse) ProtoMessage() {}

func (x *RestoreResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_syncbyte_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RestoreResponse.ProtoReflect.Descriptor instead.
func (*RestoreResponse) Descriptor() ([]byte, []int) {
	return file_internal_proto_syncbyte_proto_rawDescGZIP(), []int{3}
}

func (x *RestoreResponse) GetJobid() string {
	if x != nil {
		return x.Jobid
	}
	return ""
}

type SourceOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name     string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Server   string `protobuf:"bytes,2,opt,name=server,proto3" json:"server,omitempty"`
	User     string `protobuf:"bytes,3,opt,name=user,proto3" json:"user,omitempty"`
	Password string `protobuf:"bytes,4,opt,name=password,proto3" json:"password,omitempty"`
	Dbname   string `protobuf:"bytes,5,opt,name=dbname,proto3" json:"dbname,omitempty"`
	Version  string `protobuf:"bytes,6,opt,name=version,proto3" json:"version,omitempty"`
	Dbtype   string `protobuf:"bytes,7,opt,name=dbtype,proto3" json:"dbtype,omitempty"`
	Port     int32  `protobuf:"varint,8,opt,name=port,proto3" json:"port,omitempty"`
}

func (x *SourceOptions) Reset() {
	*x = SourceOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_syncbyte_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SourceOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SourceOptions) ProtoMessage() {}

func (x *SourceOptions) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_syncbyte_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SourceOptions.ProtoReflect.Descriptor instead.
func (*SourceOptions) Descriptor() ([]byte, []int) {
	return file_internal_proto_syncbyte_proto_rawDescGZIP(), []int{4}
}

func (x *SourceOptions) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *SourceOptions) GetServer() string {
	if x != nil {
		return x.Server
	}
	return ""
}

func (x *SourceOptions) GetUser() string {
	if x != nil {
		return x.User
	}
	return ""
}

func (x *SourceOptions) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *SourceOptions) GetDbname() string {
	if x != nil {
		return x.Dbname
	}
	return ""
}

func (x *SourceOptions) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *SourceOptions) GetDbtype() string {
	if x != nil {
		return x.Dbtype
	}
	return ""
}

func (x *SourceOptions) GetPort() int32 {
	if x != nil {
		return x.Port
	}
	return 0
}

type DestOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name     string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Server   string `protobuf:"bytes,2,opt,name=server,proto3" json:"server,omitempty"`
	User     string `protobuf:"bytes,3,opt,name=user,proto3" json:"user,omitempty"`
	Password string `protobuf:"bytes,4,opt,name=password,proto3" json:"password,omitempty"`
	Dbname   string `protobuf:"bytes,5,opt,name=dbname,proto3" json:"dbname,omitempty"`
	Version  string `protobuf:"bytes,6,opt,name=version,proto3" json:"version,omitempty"`
	Dbtype   string `protobuf:"bytes,7,opt,name=dbtype,proto3" json:"dbtype,omitempty"`
	Port     int32  `protobuf:"varint,8,opt,name=port,proto3" json:"port,omitempty"`
}

func (x *DestOptions) Reset() {
	*x = DestOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_syncbyte_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DestOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DestOptions) ProtoMessage() {}

func (x *DestOptions) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_syncbyte_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DestOptions.ProtoReflect.Descriptor instead.
func (*DestOptions) Descriptor() ([]byte, []int) {
	return file_internal_proto_syncbyte_proto_rawDescGZIP(), []int{5}
}

func (x *DestOptions) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *DestOptions) GetServer() string {
	if x != nil {
		return x.Server
	}
	return ""
}

func (x *DestOptions) GetUser() string {
	if x != nil {
		return x.User
	}
	return ""
}

func (x *DestOptions) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *DestOptions) GetDbname() string {
	if x != nil {
		return x.Dbname
	}
	return ""
}

func (x *DestOptions) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *DestOptions) GetDbtype() string {
	if x != nil {
		return x.Dbtype
	}
	return ""
}

func (x *DestOptions) GetPort() int32 {
	if x != nil {
		return x.Port
	}
	return 0
}

type BackendOptions struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Endpoint  string `protobuf:"bytes,1,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	Accesskey string `protobuf:"bytes,2,opt,name=accesskey,proto3" json:"accesskey,omitempty"`
	Secretkey string `protobuf:"bytes,3,opt,name=secretkey,proto3" json:"secretkey,omitempty"`
	Bucket    string `protobuf:"bytes,4,opt,name=bucket,proto3" json:"bucket,omitempty"`
}

func (x *BackendOptions) Reset() {
	*x = BackendOptions{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_syncbyte_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BackendOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BackendOptions) ProtoMessage() {}

func (x *BackendOptions) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_syncbyte_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BackendOptions.ProtoReflect.Descriptor instead.
func (*BackendOptions) Descriptor() ([]byte, []int) {
	return file_internal_proto_syncbyte_proto_rawDescGZIP(), []int{6}
}

func (x *BackendOptions) GetEndpoint() string {
	if x != nil {
		return x.Endpoint
	}
	return ""
}

func (x *BackendOptions) GetAccesskey() string {
	if x != nil {
		return x.Accesskey
	}
	return ""
}

func (x *BackendOptions) GetSecretkey() string {
	if x != nil {
		return x.Secretkey
	}
	return ""
}

func (x *BackendOptions) GetBucket() string {
	if x != nil {
		return x.Bucket
	}
	return ""
}

type GetJobRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Jobid string `protobuf:"bytes,1,opt,name=jobid,proto3" json:"jobid,omitempty"`
}

func (x *GetJobRequest) Reset() {
	*x = GetJobRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_syncbyte_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetJobRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetJobRequest) ProtoMessage() {}

func (x *GetJobRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_syncbyte_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetJobRequest.ProtoReflect.Descriptor instead.
func (*GetJobRequest) Descriptor() ([]byte, []int) {
	return file_internal_proto_syncbyte_proto_rawDescGZIP(), []int{7}
}

func (x *GetJobRequest) GetJobid() string {
	if x != nil {
		return x.Jobid
	}
	return ""
}

type GetJobResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status string `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *GetJobResponse) Reset() {
	*x = GetJobResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_syncbyte_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetJobResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetJobResponse) ProtoMessage() {}

func (x *GetJobResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_syncbyte_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetJobResponse.ProtoReflect.Descriptor instead.
func (*GetJobResponse) Descriptor() ([]byte, []int) {
	return file_internal_proto_syncbyte_proto_rawDescGZIP(), []int{8}
}

func (x *GetJobResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

var File_internal_proto_syncbyte_proto protoreflect.FileDescriptor

var file_internal_proto_syncbyte_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x73, 0x79, 0x6e, 0x63, 0x62, 0x79, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc0, 0x01, 0x0a, 0x0d, 0x42, 0x61, 0x63, 0x6b, 0x75,
	0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x61, 0x74, 0x61,
	0x73, 0x65, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64,
	0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x69, 0x73,
	0x63, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0a,
	0x69, 0x73, 0x63, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x12, 0x34, 0x0a, 0x0a, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x4f, 0x70, 0x74, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x4f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x52, 0x0a, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x4f, 0x70, 0x74, 0x73,
	0x12, 0x37, 0x0a, 0x0b, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x4f, 0x70, 0x74, 0x73, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x42, 0x61,
	0x63, 0x6b, 0x65, 0x6e, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x0b, 0x62, 0x61,
	0x63, 0x6b, 0x65, 0x6e, 0x64, 0x4f, 0x70, 0x74, 0x73, 0x22, 0x26, 0x0a, 0x0e, 0x42, 0x61, 0x63,
	0x6b, 0x75, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6a,
	0x6f, 0x62, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6a, 0x6f, 0x62, 0x69,
	0x64, 0x22, 0xbf, 0x01, 0x0a, 0x0e, 0x52, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x61, 0x74, 0x61, 0x73, 0x65, 0x74, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x61, 0x74, 0x61, 0x73,
	0x65, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x22, 0x0a, 0x0c, 0x69, 0x73, 0x75, 0x6e, 0x63, 0x6f,
	0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x69, 0x73,
	0x75, 0x6e, 0x63, 0x6f, 0x6d, 0x70, 0x72, 0x65, 0x73, 0x73, 0x12, 0x2e, 0x0a, 0x08, 0x64, 0x65,
	0x73, 0x74, 0x4f, 0x70, 0x74, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x65, 0x73, 0x74, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x52, 0x08, 0x64, 0x65, 0x73, 0x74, 0x4f, 0x70, 0x74, 0x73, 0x12, 0x37, 0x0a, 0x0b, 0x62, 0x61,
	0x63, 0x6b, 0x65, 0x6e, 0x64, 0x4f, 0x70, 0x74, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x42, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x4f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x0b, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x4f,
	0x70, 0x74, 0x73, 0x22, 0x27, 0x0a, 0x0f, 0x52, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6a, 0x6f, 0x62, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6a, 0x6f, 0x62, 0x69, 0x64, 0x22, 0xc9, 0x01, 0x0a,
	0x0d, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x73,
	0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x12, 0x1a,
	0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x62,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x62, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06,
	0x64, 0x62, 0x74, 0x79, 0x70, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x62,
	0x74, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x22, 0xc7, 0x01, 0x0a, 0x0b, 0x44, 0x65, 0x73,
	0x74, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x65,
	0x72, 0x76, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x73, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73,
	0x77, 0x6f, 0x72, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73,
	0x77, 0x6f, 0x72, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x62, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x62, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x62, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x62, 0x74, 0x79, 0x70, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x6f,
	0x72, 0x74, 0x22, 0x80, 0x01, 0x0a, 0x0e, 0x42, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x4f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e,
	0x74, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x6b, 0x65, 0x79, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x6b, 0x65, 0x79, 0x12,
	0x1c, 0x0a, 0x09, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x6b, 0x65, 0x79, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x6b, 0x65, 0x79, 0x12, 0x16, 0x0a,
	0x06, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x62,
	0x75, 0x63, 0x6b, 0x65, 0x74, 0x22, 0x25, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x4a, 0x6f, 0x62, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6a, 0x6f, 0x62, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6a, 0x6f, 0x62, 0x69, 0x64, 0x22, 0x28, 0x0a, 0x0e,
	0x47, 0x65, 0x74, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x32, 0xc5, 0x01, 0x0a, 0x05, 0x41, 0x67, 0x65, 0x6e, 0x74,
	0x12, 0x3c, 0x0a, 0x0b, 0x53, 0x74, 0x61, 0x72, 0x74, 0x42, 0x61, 0x63, 0x6b, 0x75, 0x70, 0x12,
	0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x42, 0x61, 0x63, 0x6b, 0x75, 0x70, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x42, 0x61,
	0x63, 0x6b, 0x75, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x3f,
	0x0a, 0x0c, 0x53, 0x74, 0x61, 0x72, 0x74, 0x52, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x12, 0x15,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65,
	0x73, 0x74, 0x6f, 0x72, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12,
	0x3d, 0x0a, 0x0c, 0x47, 0x65, 0x74, 0x4a, 0x6f, 0x62, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12,
	0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x4a, 0x6f, 0x62, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65,
	0x74, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x10,
	0x5a, 0x0e, 0x73, 0x79, 0x6e, 0x63, 0x62, 0x79, 0x74, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_proto_syncbyte_proto_rawDescOnce sync.Once
	file_internal_proto_syncbyte_proto_rawDescData = file_internal_proto_syncbyte_proto_rawDesc
)

func file_internal_proto_syncbyte_proto_rawDescGZIP() []byte {
	file_internal_proto_syncbyte_proto_rawDescOnce.Do(func() {
		file_internal_proto_syncbyte_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_proto_syncbyte_proto_rawDescData)
	})
	return file_internal_proto_syncbyte_proto_rawDescData
}

var file_internal_proto_syncbyte_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_internal_proto_syncbyte_proto_goTypes = []interface{}{
	(*BackupRequest)(nil),   // 0: proto.BackupRequest
	(*BackupResponse)(nil),  // 1: proto.BackupResponse
	(*RestoreRequest)(nil),  // 2: proto.RestoreRequest
	(*RestoreResponse)(nil), // 3: proto.RestoreResponse
	(*SourceOptions)(nil),   // 4: proto.SourceOptions
	(*DestOptions)(nil),     // 5: proto.DestOptions
	(*BackendOptions)(nil),  // 6: proto.BackendOptions
	(*GetJobRequest)(nil),   // 7: proto.GetJobRequest
	(*GetJobResponse)(nil),  // 8: proto.GetJobResponse
}
var file_internal_proto_syncbyte_proto_depIdxs = []int32{
	4, // 0: proto.BackupRequest.sourceOpts:type_name -> proto.SourceOptions
	6, // 1: proto.BackupRequest.backendOpts:type_name -> proto.BackendOptions
	5, // 2: proto.RestoreRequest.destOpts:type_name -> proto.DestOptions
	6, // 3: proto.RestoreRequest.backendOpts:type_name -> proto.BackendOptions
	0, // 4: proto.Agent.StartBackup:input_type -> proto.BackupRequest
	2, // 5: proto.Agent.StartRestore:input_type -> proto.RestoreRequest
	7, // 6: proto.Agent.GetJobStatus:input_type -> proto.GetJobRequest
	1, // 7: proto.Agent.StartBackup:output_type -> proto.BackupResponse
	3, // 8: proto.Agent.StartRestore:output_type -> proto.RestoreResponse
	8, // 9: proto.Agent.GetJobStatus:output_type -> proto.GetJobResponse
	7, // [7:10] is the sub-list for method output_type
	4, // [4:7] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_internal_proto_syncbyte_proto_init() }
func file_internal_proto_syncbyte_proto_init() {
	if File_internal_proto_syncbyte_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_proto_syncbyte_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BackupRequest); i {
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
		file_internal_proto_syncbyte_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BackupResponse); i {
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
		file_internal_proto_syncbyte_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RestoreRequest); i {
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
		file_internal_proto_syncbyte_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RestoreResponse); i {
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
		file_internal_proto_syncbyte_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SourceOptions); i {
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
		file_internal_proto_syncbyte_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DestOptions); i {
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
		file_internal_proto_syncbyte_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BackendOptions); i {
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
		file_internal_proto_syncbyte_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetJobRequest); i {
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
		file_internal_proto_syncbyte_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetJobResponse); i {
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
			RawDescriptor: file_internal_proto_syncbyte_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_proto_syncbyte_proto_goTypes,
		DependencyIndexes: file_internal_proto_syncbyte_proto_depIdxs,
		MessageInfos:      file_internal_proto_syncbyte_proto_msgTypes,
	}.Build()
	File_internal_proto_syncbyte_proto = out.File
	file_internal_proto_syncbyte_proto_rawDesc = nil
	file_internal_proto_syncbyte_proto_goTypes = nil
	file_internal_proto_syncbyte_proto_depIdxs = nil
}

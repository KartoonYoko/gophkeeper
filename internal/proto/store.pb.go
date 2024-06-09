// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v5.26.1
// source: proto/store.proto

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

type DataTypeEnum int32

const (
	DataTypeEnum_DATA_TYPE_CREDENTIALS DataTypeEnum = 0
	DataTypeEnum_DATA_TYPE_TEXT        DataTypeEnum = 1
	DataTypeEnum_DATA_TYPE_BINARY      DataTypeEnum = 2
	DataTypeEnum_DATA_TYPE_BANK_CARD   DataTypeEnum = 3
)

// Enum value maps for DataTypeEnum.
var (
	DataTypeEnum_name = map[int32]string{
		0: "DATA_TYPE_CREDENTIALS",
		1: "DATA_TYPE_TEXT",
		2: "DATA_TYPE_BINARY",
		3: "DATA_TYPE_BANK_CARD",
	}
	DataTypeEnum_value = map[string]int32{
		"DATA_TYPE_CREDENTIALS": 0,
		"DATA_TYPE_TEXT":        1,
		"DATA_TYPE_BINARY":      2,
		"DATA_TYPE_BANK_CARD":   3,
	}
)

func (x DataTypeEnum) Enum() *DataTypeEnum {
	p := new(DataTypeEnum)
	*p = x
	return p
}

func (x DataTypeEnum) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (DataTypeEnum) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_store_proto_enumTypes[0].Descriptor()
}

func (DataTypeEnum) Type() protoreflect.EnumType {
	return &file_proto_store_proto_enumTypes[0]
}

func (x DataTypeEnum) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use DataTypeEnum.Descriptor instead.
func (DataTypeEnum) EnumDescriptor() ([]byte, []int) {
	return file_proto_store_proto_rawDescGZIP(), []int{0}
}

type GetMetaDataListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetMetaDataListRequest) Reset() {
	*x = GetMetaDataListRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_store_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMetaDataListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMetaDataListRequest) ProtoMessage() {}

func (x *GetMetaDataListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_store_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMetaDataListRequest.ProtoReflect.Descriptor instead.
func (*GetMetaDataListRequest) Descriptor() ([]byte, []int) {
	return file_proto_store_proto_rawDescGZIP(), []int{0}
}

type GetMetaDataListResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []*GetMetaDataListItemResponse `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *GetMetaDataListResponse) Reset() {
	*x = GetMetaDataListResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_store_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMetaDataListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMetaDataListResponse) ProtoMessage() {}

func (x *GetMetaDataListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_store_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMetaDataListResponse.ProtoReflect.Descriptor instead.
func (*GetMetaDataListResponse) Descriptor() ([]byte, []int) {
	return file_proto_store_proto_rawDescGZIP(), []int{1}
}

func (x *GetMetaDataListResponse) GetItems() []*GetMetaDataListItemResponse {
	if x != nil {
		return x.Items
	}
	return nil
}

type GetMetaDataListItemResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          string       `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Description string       `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Type        DataTypeEnum `protobuf:"varint,3,opt,name=type,proto3,enum=proto.DataTypeEnum" json:"type,omitempty"`
	Hash        string       `protobuf:"bytes,4,opt,name=hash,proto3" json:"hash,omitempty"`
}

func (x *GetMetaDataListItemResponse) Reset() {
	*x = GetMetaDataListItemResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_store_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMetaDataListItemResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMetaDataListItemResponse) ProtoMessage() {}

func (x *GetMetaDataListItemResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_store_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMetaDataListItemResponse.ProtoReflect.Descriptor instead.
func (*GetMetaDataListItemResponse) Descriptor() ([]byte, []int) {
	return file_proto_store_proto_rawDescGZIP(), []int{2}
}

func (x *GetMetaDataListItemResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *GetMetaDataListItemResponse) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *GetMetaDataListItemResponse) GetType() DataTypeEnum {
	if x != nil {
		return x.Type
	}
	return DataTypeEnum_DATA_TYPE_CREDENTIALS
}

func (x *GetMetaDataListItemResponse) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

type SaveDataRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id                    string       `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Description           string       `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Type                  DataTypeEnum `protobuf:"varint,3,opt,name=type,proto3,enum=proto.DataTypeEnum" json:"type,omitempty"`
	Hash                  string       `protobuf:"bytes,4,opt,name=hash,proto3" json:"hash,omitempty"`
	ModificationTimestamp int64        `protobuf:"varint,5,opt,name=modification_timestamp,json=modificationTimestamp,proto3" json:"modification_timestamp,omitempty"`
	Data                  []byte       `protobuf:"bytes,6,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *SaveDataRequest) Reset() {
	*x = SaveDataRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_store_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveDataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveDataRequest) ProtoMessage() {}

func (x *SaveDataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_store_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveDataRequest.ProtoReflect.Descriptor instead.
func (*SaveDataRequest) Descriptor() ([]byte, []int) {
	return file_proto_store_proto_rawDescGZIP(), []int{3}
}

func (x *SaveDataRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *SaveDataRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *SaveDataRequest) GetType() DataTypeEnum {
	if x != nil {
		return x.Type
	}
	return DataTypeEnum_DATA_TYPE_CREDENTIALS
}

func (x *SaveDataRequest) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *SaveDataRequest) GetModificationTimestamp() int64 {
	if x != nil {
		return x.ModificationTimestamp
	}
	return 0
}

func (x *SaveDataRequest) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type SaveDataResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DataId string `protobuf:"bytes,1,opt,name=data_id,json=dataId,proto3" json:"data_id,omitempty"`
}

func (x *SaveDataResponse) Reset() {
	*x = SaveDataResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_store_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveDataResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveDataResponse) ProtoMessage() {}

func (x *SaveDataResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_store_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveDataResponse.ProtoReflect.Descriptor instead.
func (*SaveDataResponse) Descriptor() ([]byte, []int) {
	return file_proto_store_proto_rawDescGZIP(), []int{4}
}

func (x *SaveDataResponse) GetDataId() string {
	if x != nil {
		return x.DataId
	}
	return ""
}

type GetDataByIDRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetDataByIDRequest) Reset() {
	*x = GetDataByIDRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_store_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetDataByIDRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDataByIDRequest) ProtoMessage() {}

func (x *GetDataByIDRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_store_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDataByIDRequest.ProtoReflect.Descriptor instead.
func (*GetDataByIDRequest) Descriptor() ([]byte, []int) {
	return file_proto_store_proto_rawDescGZIP(), []int{5}
}

func (x *GetDataByIDRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type GetDataByIDResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type        DataTypeEnum `protobuf:"varint,1,opt,name=type,proto3,enum=proto.DataTypeEnum" json:"type,omitempty"`
	Description string       `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Hash        string       `protobuf:"bytes,3,opt,name=hash,proto3" json:"hash,omitempty"`
	Data        []byte       `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *GetDataByIDResponse) Reset() {
	*x = GetDataByIDResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_store_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetDataByIDResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetDataByIDResponse) ProtoMessage() {}

func (x *GetDataByIDResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_store_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetDataByIDResponse.ProtoReflect.Descriptor instead.
func (*GetDataByIDResponse) Descriptor() ([]byte, []int) {
	return file_proto_store_proto_rawDescGZIP(), []int{6}
}

func (x *GetDataByIDResponse) GetType() DataTypeEnum {
	if x != nil {
		return x.Type
	}
	return DataTypeEnum_DATA_TYPE_CREDENTIALS
}

func (x *GetDataByIDResponse) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *GetDataByIDResponse) GetHash() string {
	if x != nil {
		return x.Hash
	}
	return ""
}

func (x *GetDataByIDResponse) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_proto_store_proto protoreflect.FileDescriptor

var file_proto_store_proto_rawDesc = []byte{
	0x0a, 0x11, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x18, 0x0a, 0x16, 0x47, 0x65,
	0x74, 0x4d, 0x65, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x22, 0x53, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x44,
	0x61, 0x74, 0x61, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x38, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x22,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x44, 0x61,
	0x74, 0x61, 0x4c, 0x69, 0x73, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0x8c, 0x01, 0x0a, 0x1b, 0x47, 0x65,
	0x74, 0x4d, 0x65, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x4c, 0x69, 0x73, 0x74, 0x49, 0x74, 0x65,
	0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x27, 0x0a, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x54, 0x79, 0x70, 0x65, 0x45, 0x6e, 0x75, 0x6d, 0x52, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x22, 0xcb, 0x01, 0x0a, 0x0f, 0x53, 0x61, 0x76,
	0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x20, 0x0a, 0x0b,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x27,
	0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x54, 0x79, 0x70, 0x65, 0x45, 0x6e, 0x75,
	0x6d, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12, 0x35, 0x0a, 0x16, 0x6d,
	0x6f, 0x64, 0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52, 0x15, 0x6d, 0x6f, 0x64,
	0x69, 0x66, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x2b, 0x0a, 0x10, 0x53, 0x61, 0x76, 0x65, 0x44, 0x61,
	0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x64, 0x61,
	0x74, 0x61, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x61, 0x74,
	0x61, 0x49, 0x64, 0x22, 0x24, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x44, 0x61, 0x74, 0x61, 0x42, 0x79,
	0x49, 0x44, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x88, 0x01, 0x0a, 0x13, 0x47, 0x65,
	0x74, 0x44, 0x61, 0x74, 0x61, 0x42, 0x79, 0x49, 0x44, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x27, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x54, 0x79, 0x70, 0x65,
	0x45, 0x6e, 0x75, 0x6d, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04,
	0x68, 0x61, 0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68,
	0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x2a, 0x6c, 0x0a, 0x0c, 0x44, 0x61, 0x74, 0x61, 0x54, 0x79, 0x70, 0x65,
	0x45, 0x6e, 0x75, 0x6d, 0x12, 0x19, 0x0a, 0x15, 0x44, 0x41, 0x54, 0x41, 0x5f, 0x54, 0x59, 0x50,
	0x45, 0x5f, 0x43, 0x52, 0x45, 0x44, 0x45, 0x4e, 0x54, 0x49, 0x41, 0x4c, 0x53, 0x10, 0x00, 0x12,
	0x12, 0x0a, 0x0e, 0x44, 0x41, 0x54, 0x41, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x54, 0x45, 0x58,
	0x54, 0x10, 0x01, 0x12, 0x14, 0x0a, 0x10, 0x44, 0x41, 0x54, 0x41, 0x5f, 0x54, 0x59, 0x50, 0x45,
	0x5f, 0x42, 0x49, 0x4e, 0x41, 0x52, 0x59, 0x10, 0x02, 0x12, 0x17, 0x0a, 0x13, 0x44, 0x41, 0x54,
	0x41, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x42, 0x41, 0x4e, 0x4b, 0x5f, 0x43, 0x41, 0x52, 0x44,
	0x10, 0x03, 0x32, 0xe3, 0x01, 0x0a, 0x0c, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x3b, 0x0a, 0x08, 0x53, 0x61, 0x76, 0x65, 0x44, 0x61, 0x74, 0x61, 0x12,
	0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x61, 0x76, 0x65, 0x44, 0x61, 0x74, 0x61,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x53, 0x61, 0x76, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x44, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x44, 0x61, 0x74, 0x61, 0x42, 0x79, 0x49, 0x44, 0x12,
	0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x44, 0x61, 0x74, 0x61, 0x42,
	0x79, 0x49, 0x44, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x44, 0x61, 0x74, 0x61, 0x42, 0x79, 0x49, 0x44, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x50, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74,
	0x61, 0x44, 0x61, 0x74, 0x61, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x4c, 0x69, 0x73,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x47, 0x65, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x44, 0x61, 0x74, 0x61, 0x4c, 0x69, 0x73, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x32, 0x5a, 0x30, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x4b, 0x61, 0x72, 0x74, 0x6f, 0x6f, 0x6e, 0x59, 0x6f,
	0x6b, 0x6f, 0x2f, 0x67, 0x6f, 0x70, 0x68, 0x6b, 0x65, 0x65, 0x70, 0x65, 0x72, 0x2f, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_store_proto_rawDescOnce sync.Once
	file_proto_store_proto_rawDescData = file_proto_store_proto_rawDesc
)

func file_proto_store_proto_rawDescGZIP() []byte {
	file_proto_store_proto_rawDescOnce.Do(func() {
		file_proto_store_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_store_proto_rawDescData)
	})
	return file_proto_store_proto_rawDescData
}

var file_proto_store_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_store_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_proto_store_proto_goTypes = []interface{}{
	(DataTypeEnum)(0),                   // 0: proto.DataTypeEnum
	(*GetMetaDataListRequest)(nil),      // 1: proto.GetMetaDataListRequest
	(*GetMetaDataListResponse)(nil),     // 2: proto.GetMetaDataListResponse
	(*GetMetaDataListItemResponse)(nil), // 3: proto.GetMetaDataListItemResponse
	(*SaveDataRequest)(nil),             // 4: proto.SaveDataRequest
	(*SaveDataResponse)(nil),            // 5: proto.SaveDataResponse
	(*GetDataByIDRequest)(nil),          // 6: proto.GetDataByIDRequest
	(*GetDataByIDResponse)(nil),         // 7: proto.GetDataByIDResponse
}
var file_proto_store_proto_depIdxs = []int32{
	3, // 0: proto.GetMetaDataListResponse.items:type_name -> proto.GetMetaDataListItemResponse
	0, // 1: proto.GetMetaDataListItemResponse.type:type_name -> proto.DataTypeEnum
	0, // 2: proto.SaveDataRequest.type:type_name -> proto.DataTypeEnum
	0, // 3: proto.GetDataByIDResponse.type:type_name -> proto.DataTypeEnum
	4, // 4: proto.StoreService.SaveData:input_type -> proto.SaveDataRequest
	6, // 5: proto.StoreService.GetDataByID:input_type -> proto.GetDataByIDRequest
	1, // 6: proto.StoreService.GetMetaDataList:input_type -> proto.GetMetaDataListRequest
	5, // 7: proto.StoreService.SaveData:output_type -> proto.SaveDataResponse
	7, // 8: proto.StoreService.GetDataByID:output_type -> proto.GetDataByIDResponse
	2, // 9: proto.StoreService.GetMetaDataList:output_type -> proto.GetMetaDataListResponse
	7, // [7:10] is the sub-list for method output_type
	4, // [4:7] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_proto_store_proto_init() }
func file_proto_store_proto_init() {
	if File_proto_store_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_store_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetMetaDataListRequest); i {
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
		file_proto_store_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetMetaDataListResponse); i {
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
		file_proto_store_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetMetaDataListItemResponse); i {
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
		file_proto_store_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SaveDataRequest); i {
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
		file_proto_store_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SaveDataResponse); i {
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
		file_proto_store_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetDataByIDRequest); i {
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
		file_proto_store_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetDataByIDResponse); i {
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
			RawDescriptor: file_proto_store_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_store_proto_goTypes,
		DependencyIndexes: file_proto_store_proto_depIdxs,
		EnumInfos:         file_proto_store_proto_enumTypes,
		MessageInfos:      file_proto_store_proto_msgTypes,
	}.Build()
	File_proto_store_proto = out.File
	file_proto_store_proto_rawDesc = nil
	file_proto_store_proto_goTypes = nil
	file_proto_store_proto_depIdxs = nil
}

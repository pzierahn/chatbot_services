// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.3
// source: account_service.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Overview struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Payments []*Payment    `protobuf:"bytes,1,rep,name=payments,proto3" json:"payments,omitempty"`
	Usage    []*ModelUsage `protobuf:"bytes,2,rep,name=usage,proto3" json:"usage,omitempty"`
	Balance  int32         `protobuf:"varint,3,opt,name=balance,proto3" json:"balance,omitempty"`
}

func (x *Overview) Reset() {
	*x = Overview{}
	if protoimpl.UnsafeEnabled {
		mi := &file_account_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Overview) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Overview) ProtoMessage() {}

func (x *Overview) ProtoReflect() protoreflect.Message {
	mi := &file_account_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Overview.ProtoReflect.Descriptor instead.
func (*Overview) Descriptor() ([]byte, []int) {
	return file_account_service_proto_rawDescGZIP(), []int{0}
}

func (x *Overview) GetPayments() []*Payment {
	if x != nil {
		return x.Payments
	}
	return nil
}

func (x *Overview) GetUsage() []*ModelUsage {
	if x != nil {
		return x.Usage
	}
	return nil
}

func (x *Overview) GetBalance() int32 {
	if x != nil {
		return x.Balance
	}
	return 0
}

type ModelUsage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Model    string `protobuf:"bytes,1,opt,name=model,proto3" json:"model,omitempty"`
	Input    uint32 `protobuf:"varint,2,opt,name=input,proto3" json:"input,omitempty"`
	Output   uint32 `protobuf:"varint,3,opt,name=output,proto3" json:"output,omitempty"`
	Costs    uint32 `protobuf:"varint,4,opt,name=costs,proto3" json:"costs,omitempty"`
	Requests uint32 `protobuf:"varint,5,opt,name=requests,proto3" json:"requests,omitempty"`
}

func (x *ModelUsage) Reset() {
	*x = ModelUsage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_account_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ModelUsage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ModelUsage) ProtoMessage() {}

func (x *ModelUsage) ProtoReflect() protoreflect.Message {
	mi := &file_account_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ModelUsage.ProtoReflect.Descriptor instead.
func (*ModelUsage) Descriptor() ([]byte, []int) {
	return file_account_service_proto_rawDescGZIP(), []int{1}
}

func (x *ModelUsage) GetModel() string {
	if x != nil {
		return x.Model
	}
	return ""
}

func (x *ModelUsage) GetInput() uint32 {
	if x != nil {
		return x.Input
	}
	return 0
}

func (x *ModelUsage) GetOutput() uint32 {
	if x != nil {
		return x.Output
	}
	return 0
}

func (x *ModelUsage) GetCosts() uint32 {
	if x != nil {
		return x.Costs
	}
	return 0
}

func (x *ModelUsage) GetRequests() uint32 {
	if x != nil {
		return x.Requests
	}
	return 0
}

type Usage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Models []*ModelUsage `protobuf:"bytes,1,rep,name=models,proto3" json:"models,omitempty"`
}

func (x *Usage) Reset() {
	*x = Usage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_account_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Usage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Usage) ProtoMessage() {}

func (x *Usage) ProtoReflect() protoreflect.Message {
	mi := &file_account_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Usage.ProtoReflect.Descriptor instead.
func (*Usage) Descriptor() ([]byte, []int) {
	return file_account_service_proto_rawDescGZIP(), []int{2}
}

func (x *Usage) GetModels() []*ModelUsage {
	if x != nil {
		return x.Models
	}
	return nil
}

type Payment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Date   *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=date,proto3" json:"date,omitempty"`
	Amount uint32                 `protobuf:"varint,3,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (x *Payment) Reset() {
	*x = Payment{}
	if protoimpl.UnsafeEnabled {
		mi := &file_account_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Payment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Payment) ProtoMessage() {}

func (x *Payment) ProtoReflect() protoreflect.Message {
	mi := &file_account_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Payment.ProtoReflect.Descriptor instead.
func (*Payment) Descriptor() ([]byte, []int) {
	return file_account_service_proto_rawDescGZIP(), []int{3}
}

func (x *Payment) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Payment) GetDate() *timestamppb.Timestamp {
	if x != nil {
		return x.Date
	}
	return nil
}

func (x *Payment) GetAmount() uint32 {
	if x != nil {
		return x.Amount
	}
	return 0
}

type Payments struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []*Payment `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *Payments) Reset() {
	*x = Payments{}
	if protoimpl.UnsafeEnabled {
		mi := &file_account_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Payments) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Payments) ProtoMessage() {}

func (x *Payments) ProtoReflect() protoreflect.Message {
	mi := &file_account_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Payments.ProtoReflect.Descriptor instead.
func (*Payments) Descriptor() ([]byte, []int) {
	return file_account_service_proto_rawDescGZIP(), []int{4}
}

func (x *Payments) GetItems() []*Payment {
	if x != nil {
		return x.Items
	}
	return nil
}

var File_account_service_proto protoreflect.FileDescriptor

var file_account_service_proto_rawDesc = []byte{
	0x0a, 0x15, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x63, 0x68, 0x61, 0x74, 0x62, 0x6f, 0x74,
	0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70,
	0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x93, 0x01, 0x0a, 0x08, 0x4f, 0x76,
	0x65, 0x72, 0x76, 0x69, 0x65, 0x77, 0x12, 0x37, 0x0a, 0x08, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e,
	0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x62,
	0x6f, 0x74, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61,
	0x79, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x08, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12,
	0x34, 0x0a, 0x05, 0x75, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e,
	0x2e, 0x63, 0x68, 0x61, 0x74, 0x62, 0x6f, 0x74, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x55, 0x73, 0x61, 0x67, 0x65, 0x52, 0x05,
	0x75, 0x73, 0x61, 0x67, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x22,
	0x82, 0x01, 0x0a, 0x0a, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x55, 0x73, 0x61, 0x67, 0x65, 0x12, 0x14,
	0x0a, 0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6d,
	0x6f, 0x64, 0x65, 0x6c, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x75,
	0x74, 0x70, 0x75, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x6f, 0x75, 0x74, 0x70,
	0x75, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x73, 0x74, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x05, 0x63, 0x6f, 0x73, 0x74, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x72, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x73, 0x22, 0x3f, 0x0a, 0x05, 0x55, 0x73, 0x61, 0x67, 0x65, 0x12, 0x36, 0x0a,
	0x06, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e,
	0x63, 0x68, 0x61, 0x74, 0x62, 0x6f, 0x74, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e,
	0x76, 0x31, 0x2e, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x55, 0x73, 0x61, 0x67, 0x65, 0x52, 0x06, 0x6d,
	0x6f, 0x64, 0x65, 0x6c, 0x73, 0x22, 0x61, 0x0a, 0x07, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x2e, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x04, 0x64, 0x61, 0x74, 0x65,
	0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x3d, 0x0a, 0x08, 0x50, 0x61, 0x79, 0x6d,
	0x65, 0x6e, 0x74, 0x73, 0x12, 0x31, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x62, 0x6f, 0x74, 0x2e, 0x61, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74,
	0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x32, 0xd2, 0x01, 0x0a, 0x07, 0x41, 0x63, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x12, 0x3d, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x55, 0x73, 0x61, 0x67, 0x65, 0x12,
	0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x19, 0x2e, 0x63, 0x68, 0x61, 0x74, 0x62, 0x6f,
	0x74, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x43, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74,
	0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x1c, 0x2e, 0x63, 0x68, 0x61, 0x74,
	0x62, 0x6f, 0x74, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x50,
	0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x43, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x4f, 0x76,
	0x65, 0x72, 0x76, 0x69, 0x65, 0x77, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x1c,
	0x2e, 0x63, 0x68, 0x61, 0x74, 0x62, 0x6f, 0x74, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x2e, 0x76, 0x31, 0x2e, 0x4f, 0x76, 0x65, 0x72, 0x76, 0x69, 0x65, 0x77, 0x42, 0x09, 0x5a, 0x07,
	0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_account_service_proto_rawDescOnce sync.Once
	file_account_service_proto_rawDescData = file_account_service_proto_rawDesc
)

func file_account_service_proto_rawDescGZIP() []byte {
	file_account_service_proto_rawDescOnce.Do(func() {
		file_account_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_account_service_proto_rawDescData)
	})
	return file_account_service_proto_rawDescData
}

var file_account_service_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_account_service_proto_goTypes = []any{
	(*Overview)(nil),              // 0: chatbot.account.v1.Overview
	(*ModelUsage)(nil),            // 1: chatbot.account.v1.ModelUsage
	(*Usage)(nil),                 // 2: chatbot.account.v1.Usage
	(*Payment)(nil),               // 3: chatbot.account.v1.Payment
	(*Payments)(nil),              // 4: chatbot.account.v1.Payments
	(*timestamppb.Timestamp)(nil), // 5: google.protobuf.Timestamp
	(*emptypb.Empty)(nil),         // 6: google.protobuf.Empty
}
var file_account_service_proto_depIdxs = []int32{
	3, // 0: chatbot.account.v1.Overview.payments:type_name -> chatbot.account.v1.Payment
	1, // 1: chatbot.account.v1.Overview.usage:type_name -> chatbot.account.v1.ModelUsage
	1, // 2: chatbot.account.v1.Usage.models:type_name -> chatbot.account.v1.ModelUsage
	5, // 3: chatbot.account.v1.Payment.date:type_name -> google.protobuf.Timestamp
	3, // 4: chatbot.account.v1.Payments.items:type_name -> chatbot.account.v1.Payment
	6, // 5: chatbot.account.v1.Account.GetUsage:input_type -> google.protobuf.Empty
	6, // 6: chatbot.account.v1.Account.GetPayments:input_type -> google.protobuf.Empty
	6, // 7: chatbot.account.v1.Account.GetOverview:input_type -> google.protobuf.Empty
	2, // 8: chatbot.account.v1.Account.GetUsage:output_type -> chatbot.account.v1.Usage
	4, // 9: chatbot.account.v1.Account.GetPayments:output_type -> chatbot.account.v1.Payments
	0, // 10: chatbot.account.v1.Account.GetOverview:output_type -> chatbot.account.v1.Overview
	8, // [8:11] is the sub-list for method output_type
	5, // [5:8] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_account_service_proto_init() }
func file_account_service_proto_init() {
	if File_account_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_account_service_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Overview); i {
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
		file_account_service_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*ModelUsage); i {
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
		file_account_service_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*Usage); i {
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
		file_account_service_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*Payment); i {
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
		file_account_service_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*Payments); i {
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
			RawDescriptor: file_account_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_account_service_proto_goTypes,
		DependencyIndexes: file_account_service_proto_depIdxs,
		MessageInfos:      file_account_service_proto_msgTypes,
	}.Build()
	File_account_service_proto = out.File
	file_account_service_proto_rawDesc = nil
	file_account_service_proto_goTypes = nil
	file_account_service_proto_depIdxs = nil
}

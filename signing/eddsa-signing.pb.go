// Copyright © 2019 Binance
//
// This file is part of Binance. The full Binance copyright notice, including
// terms governing use, modification, and redistribution, is contained in the
// file LICENSE at the root of the source code distribution tree.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        (unknown)
// source: eddsa-signing.proto

package signing

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type SignMessagePoint struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	X     []byte `protobuf:"bytes,1,opt,name=x,proto3" json:"x,omitempty"`
	Y     []byte `protobuf:"bytes,2,opt,name=y,proto3" json:"y,omitempty"`
	Index int32  `protobuf:"varint,3,opt,name=index,proto3" json:"index,omitempty"`
}

func (x *SignMessagePoint) Reset() {
	*x = SignMessagePoint{}
	if protoimpl.UnsafeEnabled {
		mi := &file_eddsa_signing_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SignMessagePoint) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignMessagePoint) ProtoMessage() {}

func (x *SignMessagePoint) ProtoReflect() protoreflect.Message {
	mi := &file_eddsa_signing_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignMessagePoint.ProtoReflect.Descriptor instead.
func (*SignMessagePoint) Descriptor() ([]byte, []int) {
	return file_eddsa_signing_proto_rawDescGZIP(), []int{0}
}

func (x *SignMessagePoint) GetX() []byte {
	if x != nil {
		return x.X
	}
	return nil
}

func (x *SignMessagePoint) GetY() []byte {
	if x != nil {
		return x.Y
	}
	return nil
}

func (x *SignMessagePoint) GetIndex() int32 {
	if x != nil {
		return x.Index
	}
	return 0
}

type SignMessageS struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	S []byte `protobuf:"bytes,1,opt,name=s,proto3" json:"s,omitempty"`
}

func (x *SignMessageS) Reset() {
	*x = SignMessageS{}
	if protoimpl.UnsafeEnabled {
		mi := &file_eddsa_signing_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SignMessageS) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignMessageS) ProtoMessage() {}

func (x *SignMessageS) ProtoReflect() protoreflect.Message {
	mi := &file_eddsa_signing_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignMessageS.ProtoReflect.Descriptor instead.
func (*SignMessageS) Descriptor() ([]byte, []int) {
	return file_eddsa_signing_proto_rawDescGZIP(), []int{1}
}

func (x *SignMessageS) GetS() []byte {
	if x != nil {
		return x.S
	}
	return nil
}

type SignMessagePointList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Li      [][]byte `protobuf:"bytes,1,rep,name=li,proto3" json:"li,omitempty"`
	Indexes []int32  `protobuf:"varint,2,rep,packed,name=indexes,proto3" json:"indexes,omitempty"`
}

func (x *SignMessagePointList) Reset() {
	*x = SignMessagePointList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_eddsa_signing_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SignMessagePointList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SignMessagePointList) ProtoMessage() {}

func (x *SignMessagePointList) ProtoReflect() protoreflect.Message {
	mi := &file_eddsa_signing_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SignMessagePointList.ProtoReflect.Descriptor instead.
func (*SignMessagePointList) Descriptor() ([]byte, []int) {
	return file_eddsa_signing_proto_rawDescGZIP(), []int{2}
}

func (x *SignMessagePointList) GetLi() [][]byte {
	if x != nil {
		return x.Li
	}
	return nil
}

func (x *SignMessagePointList) GetIndexes() []int32 {
	if x != nil {
		return x.Indexes
	}
	return nil
}

type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Content:
	//	*Message_Smp
	//	*Message_Sms
	//	*Message_Smpl
	Content isMessage_Content `protobuf_oneof:"content"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_eddsa_signing_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_eddsa_signing_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_eddsa_signing_proto_rawDescGZIP(), []int{3}
}

func (m *Message) GetContent() isMessage_Content {
	if m != nil {
		return m.Content
	}
	return nil
}

func (x *Message) GetSmp() *SignMessagePoint {
	if x, ok := x.GetContent().(*Message_Smp); ok {
		return x.Smp
	}
	return nil
}

func (x *Message) GetSms() *SignMessageS {
	if x, ok := x.GetContent().(*Message_Sms); ok {
		return x.Sms
	}
	return nil
}

func (x *Message) GetSmpl() *SignMessagePointList {
	if x, ok := x.GetContent().(*Message_Smpl); ok {
		return x.Smpl
	}
	return nil
}

type isMessage_Content interface {
	isMessage_Content()
}

type Message_Smp struct {
	Smp *SignMessagePoint `protobuf:"bytes,1,opt,name=smp,proto3,oneof"`
}

type Message_Sms struct {
	Sms *SignMessageS `protobuf:"bytes,2,opt,name=sms,proto3,oneof"`
}

type Message_Smpl struct {
	Smpl *SignMessagePointList `protobuf:"bytes,3,opt,name=smpl,proto3,oneof"`
}

func (*Message_Smp) isMessage_Content() {}

func (*Message_Sms) isMessage_Content() {}

func (*Message_Smpl) isMessage_Content() {}

var File_eddsa_signing_proto protoreflect.FileDescriptor

var file_eddsa_signing_proto_rawDesc = []byte{
	0x0a, 0x13, 0x65, 0x64, 0x64, 0x73, 0x61, 0x2d, 0x73, 0x69, 0x67, 0x6e, 0x69, 0x6e, 0x67, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x44, 0x0a, 0x10, 0x53, 0x69, 0x67, 0x6e, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x0c, 0x0a, 0x01, 0x78, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x01, 0x78, 0x12, 0x0c, 0x0a, 0x01, 0x79, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x01, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x22, 0x1c, 0x0a, 0x0c, 0x53,
	0x69, 0x67, 0x6e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x53, 0x12, 0x0c, 0x0a, 0x01, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x01, 0x73, 0x22, 0x40, 0x0a, 0x14, 0x53, 0x69, 0x67,
	0x6e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x4c, 0x69, 0x73,
	0x74, 0x12, 0x0e, 0x0a, 0x02, 0x6c, 0x69, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x02, 0x6c,
	0x69, 0x12, 0x18, 0x0a, 0x07, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x05, 0x52, 0x07, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x73, 0x22, 0x8b, 0x01, 0x0a, 0x07,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x25, 0x0a, 0x03, 0x73, 0x6d, 0x70, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x48, 0x00, 0x52, 0x03, 0x73, 0x6d, 0x70, 0x12, 0x21,
	0x0a, 0x03, 0x73, 0x6d, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x53, 0x69,
	0x67, 0x6e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x53, 0x48, 0x00, 0x52, 0x03, 0x73, 0x6d,
	0x73, 0x12, 0x2b, 0x0a, 0x04, 0x73, 0x6d, 0x70, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x15, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x50, 0x6f, 0x69,
	0x6e, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x48, 0x00, 0x52, 0x04, 0x73, 0x6d, 0x70, 0x6c, 0x42, 0x09,
	0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x42, 0x09, 0x5a, 0x07, 0x73, 0x69, 0x67,
	0x6e, 0x69, 0x6e, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_eddsa_signing_proto_rawDescOnce sync.Once
	file_eddsa_signing_proto_rawDescData = file_eddsa_signing_proto_rawDesc
)

func file_eddsa_signing_proto_rawDescGZIP() []byte {
	file_eddsa_signing_proto_rawDescOnce.Do(func() {
		file_eddsa_signing_proto_rawDescData = protoimpl.X.CompressGZIP(file_eddsa_signing_proto_rawDescData)
	})
	return file_eddsa_signing_proto_rawDescData
}

var file_eddsa_signing_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_eddsa_signing_proto_goTypes = []interface{}{
	(*SignMessagePoint)(nil),     // 0: SignMessagePoint
	(*SignMessageS)(nil),         // 1: SignMessageS
	(*SignMessagePointList)(nil), // 2: SignMessagePointList
	(*Message)(nil),              // 3: Message
}
var file_eddsa_signing_proto_depIdxs = []int32{
	0, // 0: Message.smp:type_name -> SignMessagePoint
	1, // 1: Message.sms:type_name -> SignMessageS
	2, // 2: Message.smpl:type_name -> SignMessagePointList
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_eddsa_signing_proto_init() }
func file_eddsa_signing_proto_init() {
	if File_eddsa_signing_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_eddsa_signing_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SignMessagePoint); i {
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
		file_eddsa_signing_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SignMessageS); i {
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
		file_eddsa_signing_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SignMessagePointList); i {
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
		file_eddsa_signing_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
	file_eddsa_signing_proto_msgTypes[3].OneofWrappers = []interface{}{
		(*Message_Smp)(nil),
		(*Message_Sms)(nil),
		(*Message_Smpl)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_eddsa_signing_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_eddsa_signing_proto_goTypes,
		DependencyIndexes: file_eddsa_signing_proto_depIdxs,
		MessageInfos:      file_eddsa_signing_proto_msgTypes,
	}.Build()
	File_eddsa_signing_proto = out.File
	file_eddsa_signing_proto_rawDesc = nil
	file_eddsa_signing_proto_goTypes = nil
	file_eddsa_signing_proto_depIdxs = nil
}

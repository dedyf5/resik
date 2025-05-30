// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v4.25.1
// source: app/grpc/handler/general/home.proto

package general

import (
	common "github.com/dedyf5/resik/app/grpc/proto/common"
	status "github.com/dedyf5/resik/app/grpc/proto/status"
	response "github.com/dedyf5/resik/core/app/response"
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

type HomeReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Common *common.Request `protobuf:"bytes,1,opt,name=Common,proto3" json:"common"`  
}

func (x *HomeReq) Reset() {
	*x = HomeReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_grpc_handler_general_home_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HomeReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HomeReq) ProtoMessage() {}

func (x *HomeReq) ProtoReflect() protoreflect.Message {
	mi := &file_app_grpc_handler_general_home_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HomeReq.ProtoReflect.Descriptor instead.
func (*HomeReq) Descriptor() ([]byte, []int) {
	return file_app_grpc_handler_general_home_proto_rawDescGZIP(), []int{0}
}

func (x *HomeReq) GetCommon() *common.Request {
	if x != nil {
		return x.Common
	}
	return nil
}

type HomeRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status *status.Status `protobuf:"bytes,1,opt,name=Status,proto3" json:"status"`  
	Data   *response.App  `protobuf:"bytes,2,opt,name=Data,proto3" json:"data"`      
}

func (x *HomeRes) Reset() {
	*x = HomeRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_app_grpc_handler_general_home_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HomeRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HomeRes) ProtoMessage() {}

func (x *HomeRes) ProtoReflect() protoreflect.Message {
	mi := &file_app_grpc_handler_general_home_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HomeRes.ProtoReflect.Descriptor instead.
func (*HomeRes) Descriptor() ([]byte, []int) {
	return file_app_grpc_handler_general_home_proto_rawDescGZIP(), []int{1}
}

func (x *HomeRes) GetStatus() *status.Status {
	if x != nil {
		return x.Status
	}
	return nil
}

func (x *HomeRes) GetData() *response.App {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_app_grpc_handler_general_home_proto protoreflect.FileDescriptor

var file_app_grpc_handler_general_home_proto_rawDesc = []byte{
	0x0a, 0x23, 0x61, 0x70, 0x70, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x68, 0x61, 0x6e, 0x64, 0x6c,
	0x65, 0x72, 0x2f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x2f, 0x68, 0x6f, 0x6d, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x6c, 0x1a, 0x23,
	0x61, 0x70, 0x70, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x22, 0x61, 0x70, 0x70, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x61, 0x70,
	0x70, 0x2f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2f, 0x61, 0x70, 0x70, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x32, 0x0a, 0x07, 0x48, 0x6f, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x12,
	0x27, 0x0a, 0x06, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0f, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x52, 0x06, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x22, 0x54, 0x0a, 0x07, 0x48, 0x6f, 0x6d, 0x65,
	0x52, 0x65, 0x73, 0x12, 0x26, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x52, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x21, 0x0a, 0x04, 0x44,
	0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x72, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x41, 0x70, 0x70, 0x52, 0x04, 0x44, 0x61, 0x74, 0x61, 0x42, 0x32,
	0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x64, 0x65, 0x64,
	0x79, 0x66, 0x35, 0x2f, 0x72, 0x65, 0x73, 0x69, 0x6b, 0x2f, 0x61, 0x70, 0x70, 0x2f, 0x67, 0x72,
	0x70, 0x63, 0x2f, 0x68, 0x61, 0x6e, 0x64, 0x6c, 0x65, 0x72, 0x2f, 0x67, 0x65, 0x6e, 0x65, 0x72,
	0x61, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_app_grpc_handler_general_home_proto_rawDescOnce sync.Once
	file_app_grpc_handler_general_home_proto_rawDescData = file_app_grpc_handler_general_home_proto_rawDesc
)

func file_app_grpc_handler_general_home_proto_rawDescGZIP() []byte {
	file_app_grpc_handler_general_home_proto_rawDescOnce.Do(func() {
		file_app_grpc_handler_general_home_proto_rawDescData = protoimpl.X.CompressGZIP(file_app_grpc_handler_general_home_proto_rawDescData)
	})
	return file_app_grpc_handler_general_home_proto_rawDescData
}

var file_app_grpc_handler_general_home_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_app_grpc_handler_general_home_proto_goTypes = []interface{}{
	(*HomeReq)(nil),        // 0: general.HomeReq
	(*HomeRes)(nil),        // 1: general.HomeRes
	(*common.Request)(nil), // 2: common.Request
	(*status.Status)(nil),  // 3: status.Status
	(*response.App)(nil),   // 4: response.App
}
var file_app_grpc_handler_general_home_proto_depIdxs = []int32{
	2, // 0: general.HomeReq.Common:type_name -> common.Request
	3, // 1: general.HomeRes.Status:type_name -> status.Status
	4, // 2: general.HomeRes.Data:type_name -> response.App
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_app_grpc_handler_general_home_proto_init() }
func file_app_grpc_handler_general_home_proto_init() {
	if File_app_grpc_handler_general_home_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_app_grpc_handler_general_home_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HomeReq); i {
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
		file_app_grpc_handler_general_home_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HomeRes); i {
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
			RawDescriptor: file_app_grpc_handler_general_home_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_app_grpc_handler_general_home_proto_goTypes,
		DependencyIndexes: file_app_grpc_handler_general_home_proto_depIdxs,
		MessageInfos:      file_app_grpc_handler_general_home_proto_msgTypes,
	}.Build()
	File_app_grpc_handler_general_home_proto = out.File
	file_app_grpc_handler_general_home_proto_rawDesc = nil
	file_app_grpc_handler_general_home_proto_goTypes = nil
	file_app_grpc_handler_general_home_proto_depIdxs = nil
}

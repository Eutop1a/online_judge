// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v5.26.1
// source: submission_service.proto

package pb

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

type SubmitRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// @inject_tag: json:"user_id" form:"user_id"
	UserId int64 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	// @inject_tag: json:"language" form:"language"
	Language int32 `protobuf:"varint,2,opt,name=language,proto3" json:"language,omitempty"`
	// @inject_tag: json:"code" form:"code"
	Code string `protobuf:"bytes,3,opt,name=code,proto3" json:"code,omitempty"`
	// @inject_tag: json:"input" form:"input"
	Input []string `protobuf:"bytes,4,rep,name=input,proto3" json:"input,omitempty"`
	// @inject_tag: json:"expected" form:"expected"
	Expected []string `protobuf:"bytes,5,rep,name=expected,proto3" json:"expected,omitempty"`
	// @inject_tag: json:"time_limit" form:"time_limit"
	TimeLimit int32 `protobuf:"varint,6,opt,name=time_limit,json=timeLimit,proto3" json:"time_limit,omitempty"`
	// @inject_tag: json:"memory_limit" form:"memory_limit"
	MemoryLimit int32 `protobuf:"varint,7,opt,name=memory_limit,json=memoryLimit,proto3" json:"memory_limit,omitempty"`
	// @inject_tag: json:"total_num" form:"total_num"
	TotalNum int32 `protobuf:"varint,8,opt,name=total_num,json=totalNum,proto3" json:"total_num,omitempty"`
}

func (x *SubmitRequest) Reset() {
	*x = SubmitRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_submission_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubmitRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitRequest) ProtoMessage() {}

func (x *SubmitRequest) ProtoReflect() protoreflect.Message {
	mi := &file_submission_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubmitRequest.ProtoReflect.Descriptor instead.
func (*SubmitRequest) Descriptor() ([]byte, []int) {
	return file_submission_service_proto_rawDescGZIP(), []int{0}
}

func (x *SubmitRequest) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *SubmitRequest) GetLanguage() int32 {
	if x != nil {
		return x.Language
	}
	return 0
}

func (x *SubmitRequest) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *SubmitRequest) GetInput() []string {
	if x != nil {
		return x.Input
	}
	return nil
}

func (x *SubmitRequest) GetExpected() []string {
	if x != nil {
		return x.Expected
	}
	return nil
}

func (x *SubmitRequest) GetTimeLimit() int32 {
	if x != nil {
		return x.TimeLimit
	}
	return 0
}

func (x *SubmitRequest) GetMemoryLimit() int32 {
	if x != nil {
		return x.MemoryLimit
	}
	return 0
}

func (x *SubmitRequest) GetTotalNum() int32 {
	if x != nil {
		return x.TotalNum
	}
	return 0
}

type SubmitResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// @inject_tag: json:"user_id" form:"user_id"
	UserId int64 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	// @inject_tag: json:"status" form:"status"
	Status int32 `protobuf:"varint,2,opt,name=status,proto3" json:"status,omitempty"`
	// @inject_tag: json:"pass_num" form:"pass_num"
	PassNum int32 `protobuf:"varint,3,opt,name=pass_num,json=passNum,proto3" json:"pass_num,omitempty"`
	// @inject_tag: json:"total_num" form:"total_num"
	TotalNum int32 `protobuf:"varint,4,opt,name=total_num,json=totalNum,proto3" json:"total_num,omitempty"`
	// @inject_tag: json:"memory_usage" form:"memory_usage"
	MemoryUsage int32 `protobuf:"varint,5,opt,name=memory_usage,json=memoryUsage,proto3" json:"memory_usage,omitempty"`
	// @inject_tag: json:"runtime" form:"runtime"
	Runtime int32 `protobuf:"varint,6,opt,name=runtime,proto3" json:"runtime,omitempty"`
}

func (x *SubmitResponse) Reset() {
	*x = SubmitResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_submission_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubmitResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitResponse) ProtoMessage() {}

func (x *SubmitResponse) ProtoReflect() protoreflect.Message {
	mi := &file_submission_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubmitResponse.ProtoReflect.Descriptor instead.
func (*SubmitResponse) Descriptor() ([]byte, []int) {
	return file_submission_service_proto_rawDescGZIP(), []int{1}
}

func (x *SubmitResponse) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *SubmitResponse) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (x *SubmitResponse) GetPassNum() int32 {
	if x != nil {
		return x.PassNum
	}
	return 0
}

func (x *SubmitResponse) GetTotalNum() int32 {
	if x != nil {
		return x.TotalNum
	}
	return 0
}

func (x *SubmitResponse) GetMemoryUsage() int32 {
	if x != nil {
		return x.MemoryUsage
	}
	return 0
}

func (x *SubmitResponse) GetRuntime() int32 {
	if x != nil {
		return x.Runtime
	}
	return 0
}

var File_submission_service_proto protoreflect.FileDescriptor

var file_submission_service_proto_rawDesc = []byte{
	0x0a, 0x18, 0x73, 0x75, 0x62, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x22, 0xe9,
	0x01, 0x0a, 0x0d, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x61, 0x6e,
	0x67, 0x75, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x6c, 0x61, 0x6e,
	0x67, 0x75, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x70,
	0x75, 0x74, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x12,
	0x1a, 0x0a, 0x08, 0x65, 0x78, 0x70, 0x65, 0x63, 0x74, 0x65, 0x64, 0x18, 0x05, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x08, 0x65, 0x78, 0x70, 0x65, 0x63, 0x74, 0x65, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x74,
	0x69, 0x6d, 0x65, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x09, 0x74, 0x69, 0x6d, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x6d, 0x65,
	0x6d, 0x6f, 0x72, 0x79, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0b, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x1b, 0x0a,
	0x09, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x6e, 0x75, 0x6d, 0x18, 0x08, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x08, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x4e, 0x75, 0x6d, 0x22, 0xb6, 0x01, 0x0a, 0x0e, 0x53,
	0x75, 0x62, 0x6d, 0x69, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x17, 0x0a,
	0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06,
	0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x19,
	0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x5f, 0x6e, 0x75, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x07, 0x70, 0x61, 0x73, 0x73, 0x4e, 0x75, 0x6d, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x6f, 0x74,
	0x61, 0x6c, 0x5f, 0x6e, 0x75, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x74, 0x6f,
	0x74, 0x61, 0x6c, 0x4e, 0x75, 0x6d, 0x12, 0x21, 0x0a, 0x0c, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79,
	0x5f, 0x75, 0x73, 0x61, 0x67, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x6d, 0x65,
	0x6d, 0x6f, 0x72, 0x79, 0x55, 0x73, 0x61, 0x67, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x72, 0x75, 0x6e,
	0x74, 0x69, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x72, 0x75, 0x6e, 0x74,
	0x69, 0x6d, 0x65, 0x32, 0x41, 0x0a, 0x0a, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x12, 0x33, 0x0a, 0x0a, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x43, 0x6f, 0x64, 0x65, 0x12,
	0x11, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x12, 0x2e, 0x70, 0x62, 0x2e, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x07, 0x5a, 0x05, 0x2e, 0x2f, 0x3b, 0x70, 0x62, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_submission_service_proto_rawDescOnce sync.Once
	file_submission_service_proto_rawDescData = file_submission_service_proto_rawDesc
)

func file_submission_service_proto_rawDescGZIP() []byte {
	file_submission_service_proto_rawDescOnce.Do(func() {
		file_submission_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_submission_service_proto_rawDescData)
	})
	return file_submission_service_proto_rawDescData
}

var file_submission_service_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_submission_service_proto_goTypes = []interface{}{
	(*SubmitRequest)(nil),  // 0: pb.SubmitRequest
	(*SubmitResponse)(nil), // 1: pb.SubmitResponse
}
var file_submission_service_proto_depIdxs = []int32{
	0, // 0: pb.Submission.SubmitCode:input_type -> pb.SubmitRequest
	1, // 1: pb.Submission.SubmitCode:output_type -> pb.SubmitResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_submission_service_proto_init() }
func file_submission_service_proto_init() {
	if File_submission_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_submission_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubmitRequest); i {
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
		file_submission_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubmitResponse); i {
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
			RawDescriptor: file_submission_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_submission_service_proto_goTypes,
		DependencyIndexes: file_submission_service_proto_depIdxs,
		MessageInfos:      file_submission_service_proto_msgTypes,
	}.Build()
	File_submission_service_proto = out.File
	file_submission_service_proto_rawDesc = nil
	file_submission_service_proto_goTypes = nil
	file_submission_service_proto_depIdxs = nil
}

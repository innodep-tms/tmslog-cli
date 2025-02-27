// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v5.26.1
// source: process/grpc/logsrv.proto

package grpc

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

type LogDataSearchCondition struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FromStr       string   `protobuf:"bytes,5,opt,name=FromStr,proto3" json:"FromStr,omitempty"`
	ToStr         string   `protobuf:"bytes,4,opt,name=ToStr,proto3" json:"ToStr,omitempty"`
	ServiceId     string   `protobuf:"bytes,3,opt,name=ServiceId,proto3" json:"ServiceId,omitempty"`
	ServiceName   string   `protobuf:"bytes,1,opt,name=ServiceName,proto3" json:"ServiceName,omitempty"`
	LogLevel      []string `protobuf:"bytes,2,rep,name=LogLevel,proto3" json:"LogLevel,omitempty"`
	Message       string   `protobuf:"bytes,6,opt,name=Message,proto3" json:"Message,omitempty"`
	Format        string   `protobuf:"bytes,7,opt,name=Format,proto3" json:"Format,omitempty"`
	Columns       string   `protobuf:"bytes,8,opt,name=Columns,proto3" json:"Columns,omitempty"`
	IgnoreNewline bool     `protobuf:"varint,9,opt,name=IgnoreNewline,proto3" json:"IgnoreNewline,omitempty"`
	TimeFormat    string   `protobuf:"bytes,10,opt,name=TimeFormat,proto3" json:"TimeFormat,omitempty"`
	Tail          string   `protobuf:"bytes,11,opt,name=Tail,proto3" json:"Tail,omitempty"`
	TimeLocale    string   `protobuf:"bytes,12,opt,name=TimeLocale,proto3" json:"TimeLocale,omitempty"`
}

func (x *LogDataSearchCondition) Reset() {
	*x = LogDataSearchCondition{}
	if protoimpl.UnsafeEnabled {
		mi := &file_process_grpc_logsrv_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogDataSearchCondition) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogDataSearchCondition) ProtoMessage() {}

func (x *LogDataSearchCondition) ProtoReflect() protoreflect.Message {
	mi := &file_process_grpc_logsrv_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogDataSearchCondition.ProtoReflect.Descriptor instead.
func (*LogDataSearchCondition) Descriptor() ([]byte, []int) {
	return file_process_grpc_logsrv_proto_rawDescGZIP(), []int{0}
}

func (x *LogDataSearchCondition) GetFromStr() string {
	if x != nil {
		return x.FromStr
	}
	return ""
}

func (x *LogDataSearchCondition) GetToStr() string {
	if x != nil {
		return x.ToStr
	}
	return ""
}

func (x *LogDataSearchCondition) GetServiceId() string {
	if x != nil {
		return x.ServiceId
	}
	return ""
}

func (x *LogDataSearchCondition) GetServiceName() string {
	if x != nil {
		return x.ServiceName
	}
	return ""
}

func (x *LogDataSearchCondition) GetLogLevel() []string {
	if x != nil {
		return x.LogLevel
	}
	return nil
}

func (x *LogDataSearchCondition) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *LogDataSearchCondition) GetFormat() string {
	if x != nil {
		return x.Format
	}
	return ""
}

func (x *LogDataSearchCondition) GetColumns() string {
	if x != nil {
		return x.Columns
	}
	return ""
}

func (x *LogDataSearchCondition) GetIgnoreNewline() bool {
	if x != nil {
		return x.IgnoreNewline
	}
	return false
}

func (x *LogDataSearchCondition) GetTimeFormat() string {
	if x != nil {
		return x.TimeFormat
	}
	return ""
}

func (x *LogDataSearchCondition) GetTail() string {
	if x != nil {
		return x.Tail
	}
	return ""
}

func (x *LogDataSearchCondition) GetTimeLocale() string {
	if x != nil {
		return x.TimeLocale
	}
	return ""
}

type LogData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             int32  `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	RegDate        string `protobuf:"bytes,2,opt,name=RegDate,proto3" json:"RegDate,omitempty"`
	ServiceId      string `protobuf:"bytes,3,opt,name=ServiceId,proto3" json:"ServiceId,omitempty"`
	ServiceName    string `protobuf:"bytes,4,opt,name=ServiceName,proto3" json:"ServiceName,omitempty"`
	ServiceVersion string `protobuf:"bytes,5,opt,name=ServiceVersion,proto3" json:"ServiceVersion,omitempty"`
	LogLevel       string `protobuf:"bytes,6,opt,name=LogLevel,proto3" json:"LogLevel,omitempty"`
	Message        string `protobuf:"bytes,7,opt,name=Message,proto3" json:"Message,omitempty"`
	Caller         string `protobuf:"bytes,8,opt,name=Caller,proto3" json:"Caller,omitempty"`
	StackTrace     string `protobuf:"bytes,9,opt,name=StackTrace,proto3" json:"StackTrace,omitempty"`
}

func (x *LogData) Reset() {
	*x = LogData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_process_grpc_logsrv_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogData) ProtoMessage() {}

func (x *LogData) ProtoReflect() protoreflect.Message {
	mi := &file_process_grpc_logsrv_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogData.ProtoReflect.Descriptor instead.
func (*LogData) Descriptor() ([]byte, []int) {
	return file_process_grpc_logsrv_proto_rawDescGZIP(), []int{1}
}

func (x *LogData) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *LogData) GetRegDate() string {
	if x != nil {
		return x.RegDate
	}
	return ""
}

func (x *LogData) GetServiceId() string {
	if x != nil {
		return x.ServiceId
	}
	return ""
}

func (x *LogData) GetServiceName() string {
	if x != nil {
		return x.ServiceName
	}
	return ""
}

func (x *LogData) GetServiceVersion() string {
	if x != nil {
		return x.ServiceVersion
	}
	return ""
}

func (x *LogData) GetLogLevel() string {
	if x != nil {
		return x.LogLevel
	}
	return ""
}

func (x *LogData) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *LogData) GetCaller() string {
	if x != nil {
		return x.Caller
	}
	return ""
}

func (x *LogData) GetStackTrace() string {
	if x != nil {
		return x.StackTrace
	}
	return ""
}

var File_process_grpc_logsrv_proto protoreflect.FileDescriptor

var file_process_grpc_logsrv_proto_rawDesc = []byte{
	0x0a, 0x19, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x6c,
	0x6f, 0x67, 0x73, 0x72, 0x76, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x67, 0x72, 0x70,
	0x63, 0x22, 0xea, 0x02, 0x0a, 0x16, 0x4c, 0x6f, 0x67, 0x44, 0x61, 0x74, 0x61, 0x53, 0x65, 0x61,
	0x72, 0x63, 0x68, 0x43, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x18, 0x0a, 0x07,
	0x46, 0x72, 0x6f, 0x6d, 0x53, 0x74, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x46,
	0x72, 0x6f, 0x6d, 0x53, 0x74, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x54, 0x6f, 0x53, 0x74, 0x72, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x54, 0x6f, 0x53, 0x74, 0x72, 0x12, 0x1c, 0x0a, 0x09,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08,
	0x4c, 0x6f, 0x67, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08,
	0x4c, 0x6f, 0x67, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x12, 0x18, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x43, 0x6f,
	0x6c, 0x75, 0x6d, 0x6e, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x43, 0x6f, 0x6c,
	0x75, 0x6d, 0x6e, 0x73, 0x12, 0x24, 0x0a, 0x0d, 0x49, 0x67, 0x6e, 0x6f, 0x72, 0x65, 0x4e, 0x65,
	0x77, 0x6c, 0x69, 0x6e, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0d, 0x49, 0x67, 0x6e,
	0x6f, 0x72, 0x65, 0x4e, 0x65, 0x77, 0x6c, 0x69, 0x6e, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x54, 0x69,
	0x6d, 0x65, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a,
	0x54, 0x69, 0x6d, 0x65, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x61,
	0x69, 0x6c, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x54, 0x61, 0x69, 0x6c, 0x12, 0x1e,
	0x0a, 0x0a, 0x54, 0x69, 0x6d, 0x65, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x65, 0x18, 0x0c, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x54, 0x69, 0x6d, 0x65, 0x4c, 0x6f, 0x63, 0x61, 0x6c, 0x65, 0x22, 0x89,
	0x02, 0x0a, 0x07, 0x4c, 0x6f, 0x67, 0x44, 0x61, 0x74, 0x61, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x52, 0x65,
	0x67, 0x44, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x52, 0x65, 0x67,
	0x44, 0x61, 0x74, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x49,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x49, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4e, 0x61, 0x6d,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x4e, 0x61, 0x6d, 0x65, 0x12, 0x26, 0x0a, 0x0e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08,
	0x4c, 0x6f, 0x67, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x4c, 0x6f, 0x67, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x12, 0x18, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x43, 0x61, 0x6c, 0x6c, 0x65, 0x72, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x43, 0x61, 0x6c, 0x6c, 0x65, 0x72, 0x12, 0x1e, 0x0a, 0x0a, 0x53, 0x74,
	0x61, 0x63, 0x6b, 0x54, 0x72, 0x61, 0x63, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a,
	0x53, 0x74, 0x61, 0x63, 0x6b, 0x54, 0x72, 0x61, 0x63, 0x65, 0x32, 0x41, 0x0a, 0x03, 0x6c, 0x6f,
	0x67, 0x12, 0x3a, 0x0a, 0x07, 0x72, 0x65, 0x61, 0x64, 0x4c, 0x6f, 0x67, 0x12, 0x1c, 0x2e, 0x67,
	0x72, 0x70, 0x63, 0x2e, 0x4c, 0x6f, 0x67, 0x44, 0x61, 0x74, 0x61, 0x53, 0x65, 0x61, 0x72, 0x63,
	0x68, 0x43, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x0d, 0x2e, 0x67, 0x72, 0x70,
	0x63, 0x2e, 0x4c, 0x6f, 0x67, 0x44, 0x61, 0x74, 0x61, 0x22, 0x00, 0x30, 0x01, 0x42, 0x35, 0x5a,
	0x33, 0x62, 0x69, 0x74, 0x62, 0x75, 0x63, 0x6b, 0x65, 0x74, 0x2e, 0x6f, 0x72, 0x67, 0x2f, 0x69,
	0x6e, 0x6e, 0x6f, 0x64, 0x65, 0x70, 0x2f, 0x6e, 0x74, 0x6d, 0x73, 0x2d, 0x6c, 0x6f, 0x67, 0x2d,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x2f,
	0x67, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_process_grpc_logsrv_proto_rawDescOnce sync.Once
	file_process_grpc_logsrv_proto_rawDescData = file_process_grpc_logsrv_proto_rawDesc
)

func file_process_grpc_logsrv_proto_rawDescGZIP() []byte {
	file_process_grpc_logsrv_proto_rawDescOnce.Do(func() {
		file_process_grpc_logsrv_proto_rawDescData = protoimpl.X.CompressGZIP(file_process_grpc_logsrv_proto_rawDescData)
	})
	return file_process_grpc_logsrv_proto_rawDescData
}

var file_process_grpc_logsrv_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_process_grpc_logsrv_proto_goTypes = []interface{}{
	(*LogDataSearchCondition)(nil), // 0: grpc.LogDataSearchCondition
	(*LogData)(nil),                // 1: grpc.LogData
}
var file_process_grpc_logsrv_proto_depIdxs = []int32{
	0, // 0: grpc.log.readLog:input_type -> grpc.LogDataSearchCondition
	1, // 1: grpc.log.readLog:output_type -> grpc.LogData
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_process_grpc_logsrv_proto_init() }
func file_process_grpc_logsrv_proto_init() {
	if File_process_grpc_logsrv_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_process_grpc_logsrv_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogDataSearchCondition); i {
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
		file_process_grpc_logsrv_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogData); i {
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
			RawDescriptor: file_process_grpc_logsrv_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_process_grpc_logsrv_proto_goTypes,
		DependencyIndexes: file_process_grpc_logsrv_proto_depIdxs,
		MessageInfos:      file_process_grpc_logsrv_proto_msgTypes,
	}.Build()
	File_process_grpc_logsrv_proto = out.File
	file_process_grpc_logsrv_proto_rawDesc = nil
	file_process_grpc_logsrv_proto_goTypes = nil
	file_process_grpc_logsrv_proto_depIdxs = nil
}

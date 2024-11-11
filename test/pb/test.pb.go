// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v3.12.4
// source: test.proto

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

type Mode int32

const (
	Mode_MODE_A Mode = 0
	Mode_MODE_B Mode = 1
	Mode_MODE_C Mode = 2
)

// Enum value maps for Mode.
var (
	Mode_name = map[int32]string{
		0: "MODE_A",
		1: "MODE_B",
		2: "MODE_C",
	}
	Mode_value = map[string]int32{
		"MODE_A": 0,
		"MODE_B": 1,
		"MODE_C": 2,
	}
)

func (x Mode) Enum() *Mode {
	p := new(Mode)
	*p = x
	return p
}

func (x Mode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Mode) Descriptor() protoreflect.EnumDescriptor {
	return file_test_proto_enumTypes[0].Descriptor()
}

func (Mode) Type() protoreflect.EnumType {
	return &file_test_proto_enumTypes[0]
}

func (x Mode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Mode.Descriptor instead.
func (Mode) EnumDescriptor() ([]byte, []int) {
	return file_test_proto_rawDescGZIP(), []int{0}
}

type Small struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	I32  int32  `protobuf:"varint,1,opt,name=i32,proto3" json:"i32,omitempty"`
	Flag bool   `protobuf:"varint,2,opt,name=flag,proto3" json:"flag,omitempty"`
	Str  string `protobuf:"bytes,3,opt,name=str,proto3" json:"str,omitempty"`
}

func (x *Small) Reset() {
	*x = Small{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Small) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Small) ProtoMessage() {}

func (x *Small) ProtoReflect() protoreflect.Message {
	mi := &file_test_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Small.ProtoReflect.Descriptor instead.
func (*Small) Descriptor() ([]byte, []int) {
	return file_test_proto_rawDescGZIP(), []int{0}
}

func (x *Small) GetI32() int32 {
	if x != nil {
		return x.I32
	}
	return 0
}

func (x *Small) GetFlag() bool {
	if x != nil {
		return x.Flag
	}
	return false
}

func (x *Small) GetStr() string {
	if x != nil {
		return x.Str
	}
	return ""
}

type Vec2D struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	X []*Vec2D_Vec1D `protobuf:"bytes,1,rep,name=_,proto3" json:"_,omitempty"`
}

func (x *Vec2D) Reset() {
	*x = Vec2D{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Vec2D) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Vec2D) ProtoMessage() {}

func (x *Vec2D) ProtoReflect() protoreflect.Message {
	mi := &file_test_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Vec2D.ProtoReflect.Descriptor instead.
func (*Vec2D) Descriptor() ([]byte, []int) {
	return file_test_proto_rawDescGZIP(), []int{1}
}

func (x *Vec2D) GetX() []*Vec2D_Vec1D {
	if x != nil {
		return x.X
	}
	return nil
}

type ArrMap struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	X map[string]*ArrMap_Array `protobuf:"bytes,1,rep,name=_,proto3" json:"_,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *ArrMap) Reset() {
	*x = ArrMap{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ArrMap) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ArrMap) ProtoMessage() {}

func (x *ArrMap) ProtoReflect() protoreflect.Message {
	mi := &file_test_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ArrMap.ProtoReflect.Descriptor instead.
func (*ArrMap) Descriptor() ([]byte, []int) {
	return file_test_proto_rawDescGZIP(), []int{2}
}

func (x *ArrMap) GetX() map[string]*ArrMap_Array {
	if x != nil {
		return x.X
	}
	return nil
}

type Main struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	I32     int32            `protobuf:"varint,1,opt,name=i32,proto3" json:"i32,omitempty"`
	U32     uint32           `protobuf:"varint,2,opt,name=u32,proto3" json:"u32,omitempty"`
	I64     int64            `protobuf:"varint,3,opt,name=i64,proto3" json:"i64,omitempty"`
	U64     uint64           `protobuf:"varint,4,opt,name=u64,proto3" json:"u64,omitempty"`
	Flag    bool             `protobuf:"varint,5,opt,name=flag,proto3" json:"flag,omitempty"`
	Mode    Mode             `protobuf:"varint,6,opt,name=mode,proto3,enum=test.Mode" json:"mode,omitempty"`
	Str     string           `protobuf:"bytes,7,opt,name=str,proto3" json:"str,omitempty"`
	Data    []byte           `protobuf:"bytes,8,opt,name=data,proto3" json:"data,omitempty"`
	F32     float32          `protobuf:"fixed32,9,opt,name=f32,proto3" json:"f32,omitempty"`
	F64     float64          `protobuf:"fixed64,10,opt,name=f64,proto3" json:"f64,omitempty"`
	Object  *Small           `protobuf:"bytes,11,opt,name=object,proto3" json:"object,omitempty"`
	I32V    []int32          `protobuf:"varint,12,rep,packed,name=i32v,proto3" json:"i32v,omitempty"`
	U64V    []uint64         `protobuf:"varint,13,rep,packed,name=u64v,proto3" json:"u64v,omitempty"`
	Strv    []string         `protobuf:"bytes,14,rep,name=strv,proto3" json:"strv,omitempty"`
	Datav   [][]byte         `protobuf:"bytes,15,rep,name=datav,proto3" json:"datav,omitempty"`
	F32V    []float32        `protobuf:"fixed32,16,rep,packed,name=f32v,proto3" json:"f32v,omitempty"`
	F64V    []float64        `protobuf:"fixed64,17,rep,packed,name=f64v,proto3" json:"f64v,omitempty"`
	Flags   []bool           `protobuf:"varint,18,rep,packed,name=flags,proto3" json:"flags,omitempty"`
	Objectv []*Small         `protobuf:"bytes,19,rep,name=objectv,proto3" json:"objectv,omitempty"`
	TU32    uint32           `protobuf:"fixed32,20,opt,name=t_u32,json=tU32,proto3" json:"t_u32,omitempty"`
	TI32    int32            `protobuf:"fixed32,21,opt,name=t_i32,json=tI32,proto3" json:"t_i32,omitempty"`
	TS32    int32            `protobuf:"zigzag32,22,opt,name=t_s32,json=tS32,proto3" json:"t_s32,omitempty"`
	TU64    uint64           `protobuf:"fixed64,23,opt,name=t_u64,json=tU64,proto3" json:"t_u64,omitempty"`
	TI64    int64            `protobuf:"fixed64,24,opt,name=t_i64,json=tI64,proto3" json:"t_i64,omitempty"`
	TS64    int64            `protobuf:"zigzag64,25,opt,name=t_s64,json=tS64,proto3" json:"t_s64,omitempty"`
	Index   map[string]int32 `protobuf:"bytes,26,rep,name=index,proto3" json:"index,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	Objects map[int32]*Small `protobuf:"bytes,27,rep,name=objects,proto3" json:"objects,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Matrix  *Vec2D           `protobuf:"bytes,28,opt,name=matrix,proto3" json:"matrix,omitempty"`
	Vector  []*ArrMap        `protobuf:"bytes,29,rep,name=vector,proto3" json:"vector,omitempty"`
	Arrays  *ArrMap          `protobuf:"bytes,30,opt,name=arrays,proto3" json:"arrays,omitempty"`
}

func (x *Main) Reset() {
	*x = Main{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Main) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Main) ProtoMessage() {}

func (x *Main) ProtoReflect() protoreflect.Message {
	mi := &file_test_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Main.ProtoReflect.Descriptor instead.
func (*Main) Descriptor() ([]byte, []int) {
	return file_test_proto_rawDescGZIP(), []int{3}
}

func (x *Main) GetI32() int32 {
	if x != nil {
		return x.I32
	}
	return 0
}

func (x *Main) GetU32() uint32 {
	if x != nil {
		return x.U32
	}
	return 0
}

func (x *Main) GetI64() int64 {
	if x != nil {
		return x.I64
	}
	return 0
}

func (x *Main) GetU64() uint64 {
	if x != nil {
		return x.U64
	}
	return 0
}

func (x *Main) GetFlag() bool {
	if x != nil {
		return x.Flag
	}
	return false
}

func (x *Main) GetMode() Mode {
	if x != nil {
		return x.Mode
	}
	return Mode_MODE_A
}

func (x *Main) GetStr() string {
	if x != nil {
		return x.Str
	}
	return ""
}

func (x *Main) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *Main) GetF32() float32 {
	if x != nil {
		return x.F32
	}
	return 0
}

func (x *Main) GetF64() float64 {
	if x != nil {
		return x.F64
	}
	return 0
}

func (x *Main) GetObject() *Small {
	if x != nil {
		return x.Object
	}
	return nil
}

func (x *Main) GetI32V() []int32 {
	if x != nil {
		return x.I32V
	}
	return nil
}

func (x *Main) GetU64V() []uint64 {
	if x != nil {
		return x.U64V
	}
	return nil
}

func (x *Main) GetStrv() []string {
	if x != nil {
		return x.Strv
	}
	return nil
}

func (x *Main) GetDatav() [][]byte {
	if x != nil {
		return x.Datav
	}
	return nil
}

func (x *Main) GetF32V() []float32 {
	if x != nil {
		return x.F32V
	}
	return nil
}

func (x *Main) GetF64V() []float64 {
	if x != nil {
		return x.F64V
	}
	return nil
}

func (x *Main) GetFlags() []bool {
	if x != nil {
		return x.Flags
	}
	return nil
}

func (x *Main) GetObjectv() []*Small {
	if x != nil {
		return x.Objectv
	}
	return nil
}

func (x *Main) GetTU32() uint32 {
	if x != nil {
		return x.TU32
	}
	return 0
}

func (x *Main) GetTI32() int32 {
	if x != nil {
		return x.TI32
	}
	return 0
}

func (x *Main) GetTS32() int32 {
	if x != nil {
		return x.TS32
	}
	return 0
}

func (x *Main) GetTU64() uint64 {
	if x != nil {
		return x.TU64
	}
	return 0
}

func (x *Main) GetTI64() int64 {
	if x != nil {
		return x.TI64
	}
	return 0
}

func (x *Main) GetTS64() int64 {
	if x != nil {
		return x.TS64
	}
	return 0
}

func (x *Main) GetIndex() map[string]int32 {
	if x != nil {
		return x.Index
	}
	return nil
}

func (x *Main) GetObjects() map[int32]*Small {
	if x != nil {
		return x.Objects
	}
	return nil
}

func (x *Main) GetMatrix() *Vec2D {
	if x != nil {
		return x.Matrix
	}
	return nil
}

func (x *Main) GetVector() []*ArrMap {
	if x != nil {
		return x.Vector
	}
	return nil
}

func (x *Main) GetArrays() *ArrMap {
	if x != nil {
		return x.Arrays
	}
	return nil
}

type Vec2D_Vec1D struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	X []float32 `protobuf:"fixed32,1,rep,packed,name=_,proto3" json:"_,omitempty"`
}

func (x *Vec2D_Vec1D) Reset() {
	*x = Vec2D_Vec1D{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Vec2D_Vec1D) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Vec2D_Vec1D) ProtoMessage() {}

func (x *Vec2D_Vec1D) ProtoReflect() protoreflect.Message {
	mi := &file_test_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Vec2D_Vec1D.ProtoReflect.Descriptor instead.
func (*Vec2D_Vec1D) Descriptor() ([]byte, []int) {
	return file_test_proto_rawDescGZIP(), []int{1, 0}
}

func (x *Vec2D_Vec1D) GetX() []float32 {
	if x != nil {
		return x.X
	}
	return nil
}

type ArrMap_Array struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	X []float32 `protobuf:"fixed32,1,rep,packed,name=_,proto3" json:"_,omitempty"`
}

func (x *ArrMap_Array) Reset() {
	*x = ArrMap_Array{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ArrMap_Array) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ArrMap_Array) ProtoMessage() {}

func (x *ArrMap_Array) ProtoReflect() protoreflect.Message {
	mi := &file_test_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ArrMap_Array.ProtoReflect.Descriptor instead.
func (*ArrMap_Array) Descriptor() ([]byte, []int) {
	return file_test_proto_rawDescGZIP(), []int{2, 0}
}

func (x *ArrMap_Array) GetX() []float32 {
	if x != nil {
		return x.X
	}
	return nil
}

var File_test_proto protoreflect.FileDescriptor

var file_test_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x74, 0x65,
	0x73, 0x74, 0x22, 0x3f, 0x0a, 0x05, 0x53, 0x6d, 0x61, 0x6c, 0x6c, 0x12, 0x10, 0x0a, 0x03, 0x69,
	0x33, 0x32, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x69, 0x33, 0x32, 0x12, 0x12, 0x0a,
	0x04, 0x66, 0x6c, 0x61, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x66, 0x6c, 0x61,
	0x67, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x74, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x73, 0x74, 0x72, 0x22, 0x3d, 0x0a, 0x05, 0x56, 0x65, 0x63, 0x32, 0x44, 0x12, 0x1e, 0x0a, 0x01,
	0x5f, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x56,
	0x65, 0x63, 0x32, 0x44, 0x2e, 0x56, 0x65, 0x63, 0x31, 0x44, 0x52, 0x00, 0x1a, 0x14, 0x0a, 0x05,
	0x56, 0x65, 0x63, 0x31, 0x44, 0x12, 0x0b, 0x0a, 0x01, 0x5f, 0x18, 0x01, 0x20, 0x03, 0x28, 0x02,
	0x52, 0x00, 0x22, 0x88, 0x01, 0x0a, 0x06, 0x41, 0x72, 0x72, 0x4d, 0x61, 0x70, 0x12, 0x1f, 0x0a,
	0x01, 0x5f, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e,
	0x41, 0x72, 0x72, 0x4d, 0x61, 0x70, 0x2e, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x00, 0x1a, 0x14,
	0x0a, 0x05, 0x41, 0x72, 0x72, 0x61, 0x79, 0x12, 0x0b, 0x0a, 0x01, 0x5f, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x02, 0x52, 0x00, 0x1a, 0x47, 0x0a, 0x05, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x28, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12,
	0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x41, 0x72, 0x72, 0x4d, 0x61, 0x70, 0x2e, 0x41, 0x72, 0x72,
	0x61, 0x79, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xfa, 0x06,
	0x0a, 0x04, 0x4d, 0x61, 0x69, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x33, 0x32, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x03, 0x69, 0x33, 0x32, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x33, 0x32, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x03, 0x75, 0x33, 0x32, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x36,
	0x34, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x69, 0x36, 0x34, 0x12, 0x10, 0x0a, 0x03,
	0x75, 0x36, 0x34, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x03, 0x75, 0x36, 0x34, 0x12, 0x12,
	0x0a, 0x04, 0x66, 0x6c, 0x61, 0x67, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x66, 0x6c,
	0x61, 0x67, 0x12, 0x1e, 0x0a, 0x04, 0x6d, 0x6f, 0x64, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x0a, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x4d, 0x6f, 0x64, 0x65, 0x52, 0x04, 0x6d, 0x6f,
	0x64, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x74, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x73, 0x74, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x08, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x10, 0x0a, 0x03, 0x66, 0x33, 0x32, 0x18,
	0x09, 0x20, 0x01, 0x28, 0x02, 0x52, 0x03, 0x66, 0x33, 0x32, 0x12, 0x10, 0x0a, 0x03, 0x66, 0x36,
	0x34, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x66, 0x36, 0x34, 0x12, 0x23, 0x0a, 0x06,
	0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x74,
	0x65, 0x73, 0x74, 0x2e, 0x53, 0x6d, 0x61, 0x6c, 0x6c, 0x52, 0x06, 0x6f, 0x62, 0x6a, 0x65, 0x63,
	0x74, 0x12, 0x12, 0x0a, 0x04, 0x69, 0x33, 0x32, 0x76, 0x18, 0x0c, 0x20, 0x03, 0x28, 0x05, 0x52,
	0x04, 0x69, 0x33, 0x32, 0x76, 0x12, 0x12, 0x0a, 0x04, 0x75, 0x36, 0x34, 0x76, 0x18, 0x0d, 0x20,
	0x03, 0x28, 0x04, 0x52, 0x04, 0x75, 0x36, 0x34, 0x76, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x74, 0x72,
	0x76, 0x18, 0x0e, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x73, 0x74, 0x72, 0x76, 0x12, 0x14, 0x0a,
	0x05, 0x64, 0x61, 0x74, 0x61, 0x76, 0x18, 0x0f, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x05, 0x64, 0x61,
	0x74, 0x61, 0x76, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x33, 0x32, 0x76, 0x18, 0x10, 0x20, 0x03, 0x28,
	0x02, 0x52, 0x04, 0x66, 0x33, 0x32, 0x76, 0x12, 0x12, 0x0a, 0x04, 0x66, 0x36, 0x34, 0x76, 0x18,
	0x11, 0x20, 0x03, 0x28, 0x01, 0x52, 0x04, 0x66, 0x36, 0x34, 0x76, 0x12, 0x14, 0x0a, 0x05, 0x66,
	0x6c, 0x61, 0x67, 0x73, 0x18, 0x12, 0x20, 0x03, 0x28, 0x08, 0x52, 0x05, 0x66, 0x6c, 0x61, 0x67,
	0x73, 0x12, 0x25, 0x0a, 0x07, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x76, 0x18, 0x13, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x53, 0x6d, 0x61, 0x6c, 0x6c, 0x52,
	0x07, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x76, 0x12, 0x13, 0x0a, 0x05, 0x74, 0x5f, 0x75, 0x33,
	0x32, 0x18, 0x14, 0x20, 0x01, 0x28, 0x07, 0x52, 0x04, 0x74, 0x55, 0x33, 0x32, 0x12, 0x13, 0x0a,
	0x05, 0x74, 0x5f, 0x69, 0x33, 0x32, 0x18, 0x15, 0x20, 0x01, 0x28, 0x0f, 0x52, 0x04, 0x74, 0x49,
	0x33, 0x32, 0x12, 0x13, 0x0a, 0x05, 0x74, 0x5f, 0x73, 0x33, 0x32, 0x18, 0x16, 0x20, 0x01, 0x28,
	0x11, 0x52, 0x04, 0x74, 0x53, 0x33, 0x32, 0x12, 0x13, 0x0a, 0x05, 0x74, 0x5f, 0x75, 0x36, 0x34,
	0x18, 0x17, 0x20, 0x01, 0x28, 0x06, 0x52, 0x04, 0x74, 0x55, 0x36, 0x34, 0x12, 0x13, 0x0a, 0x05,
	0x74, 0x5f, 0x69, 0x36, 0x34, 0x18, 0x18, 0x20, 0x01, 0x28, 0x10, 0x52, 0x04, 0x74, 0x49, 0x36,
	0x34, 0x12, 0x13, 0x0a, 0x05, 0x74, 0x5f, 0x73, 0x36, 0x34, 0x18, 0x19, 0x20, 0x01, 0x28, 0x12,
	0x52, 0x04, 0x74, 0x53, 0x36, 0x34, 0x12, 0x2b, 0x0a, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18,
	0x1a, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x4d, 0x61, 0x69,
	0x6e, 0x2e, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x05, 0x69, 0x6e,
	0x64, 0x65, 0x78, 0x12, 0x31, 0x0a, 0x07, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x18, 0x1b,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x4d, 0x61, 0x69, 0x6e,
	0x2e, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x6f,
	0x62, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x12, 0x23, 0x0a, 0x06, 0x6d, 0x61, 0x74, 0x72, 0x69, 0x78,
	0x18, 0x1c, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x56, 0x65,
	0x63, 0x32, 0x44, 0x52, 0x06, 0x6d, 0x61, 0x74, 0x72, 0x69, 0x78, 0x12, 0x24, 0x0a, 0x06, 0x76,
	0x65, 0x63, 0x74, 0x6f, 0x72, 0x18, 0x1d, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x74, 0x65,
	0x73, 0x74, 0x2e, 0x41, 0x72, 0x72, 0x4d, 0x61, 0x70, 0x52, 0x06, 0x76, 0x65, 0x63, 0x74, 0x6f,
	0x72, 0x12, 0x24, 0x0a, 0x06, 0x61, 0x72, 0x72, 0x61, 0x79, 0x73, 0x18, 0x1e, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0c, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x41, 0x72, 0x72, 0x4d, 0x61, 0x70, 0x52,
	0x06, 0x61, 0x72, 0x72, 0x61, 0x79, 0x73, 0x1a, 0x38, 0x0a, 0x0a, 0x49, 0x6e, 0x64, 0x65, 0x78,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38,
	0x01, 0x1a, 0x47, 0x0a, 0x0c, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x21, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x53, 0x6d, 0x61, 0x6c, 0x6c, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x2a, 0x2a, 0x0a, 0x04, 0x4d, 0x6f,
	0x64, 0x65, 0x12, 0x0a, 0x0a, 0x06, 0x4d, 0x4f, 0x44, 0x45, 0x5f, 0x41, 0x10, 0x00, 0x12, 0x0a,
	0x0a, 0x06, 0x4d, 0x4f, 0x44, 0x45, 0x5f, 0x42, 0x10, 0x01, 0x12, 0x0a, 0x0a, 0x06, 0x4d, 0x4f,
	0x44, 0x45, 0x5f, 0x43, 0x10, 0x02, 0x42, 0x2d, 0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x65, 0x74, 0x65, 0x72, 0x72, 0x6b, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2d, 0x67, 0x6f, 0x2f, 0x74, 0x65, 0x73, 0x74, 0x2f,
	0x70, 0x62, 0xf8, 0x01, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_test_proto_rawDescOnce sync.Once
	file_test_proto_rawDescData = file_test_proto_rawDesc
)

func file_test_proto_rawDescGZIP() []byte {
	file_test_proto_rawDescOnce.Do(func() {
		file_test_proto_rawDescData = protoimpl.X.CompressGZIP(file_test_proto_rawDescData)
	})
	return file_test_proto_rawDescData
}

var file_test_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_test_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_test_proto_goTypes = []interface{}{
	(Mode)(0),            // 0: test.Mode
	(*Small)(nil),        // 1: test.Small
	(*Vec2D)(nil),        // 2: test.Vec2D
	(*ArrMap)(nil),       // 3: test.ArrMap
	(*Main)(nil),         // 4: test.Main
	(*Vec2D_Vec1D)(nil),  // 5: test.Vec2D.Vec1D
	(*ArrMap_Array)(nil), // 6: test.ArrMap.Array
	nil,                  // 7: test.ArrMap.Entry
	nil,                  // 8: test.Main.IndexEntry
	nil,                  // 9: test.Main.ObjectsEntry
}
var file_test_proto_depIdxs = []int32{
	5,  // 0: test.Vec2D._:type_name -> test.Vec2D.Vec1D
	7,  // 1: test.ArrMap._:type_name -> test.ArrMap.Entry
	0,  // 2: test.Main.mode:type_name -> test.Mode
	1,  // 3: test.Main.object:type_name -> test.Small
	1,  // 4: test.Main.objectv:type_name -> test.Small
	8,  // 5: test.Main.index:type_name -> test.Main.IndexEntry
	9,  // 6: test.Main.objects:type_name -> test.Main.ObjectsEntry
	2,  // 7: test.Main.matrix:type_name -> test.Vec2D
	3,  // 8: test.Main.vector:type_name -> test.ArrMap
	3,  // 9: test.Main.arrays:type_name -> test.ArrMap
	6,  // 10: test.ArrMap.Entry.value:type_name -> test.ArrMap.Array
	1,  // 11: test.Main.ObjectsEntry.value:type_name -> test.Small
	12, // [12:12] is the sub-list for method output_type
	12, // [12:12] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_test_proto_init() }
func file_test_proto_init() {
	if File_test_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_test_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Small); i {
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
		file_test_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Vec2D); i {
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
		file_test_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ArrMap); i {
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
		file_test_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Main); i {
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
		file_test_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Vec2D_Vec1D); i {
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
		file_test_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ArrMap_Array); i {
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
			RawDescriptor: file_test_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_test_proto_goTypes,
		DependencyIndexes: file_test_proto_depIdxs,
		EnumInfos:         file_test_proto_enumTypes,
		MessageInfos:      file_test_proto_msgTypes,
	}.Build()
	File_test_proto = out.File
	file_test_proto_rawDesc = nil
	file_test_proto_goTypes = nil
	file_test_proto_depIdxs = nil
}

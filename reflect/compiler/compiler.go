package compiler

import (
	"errors"
	"unsafe"

	"google.golang.org/protobuf/proto"
	pb "google.golang.org/protobuf/types/descriptorpb"
)

/*
#cgo LDFLAGS: -pthread -lstdc++ -static-libstdc++ -static-libgcc -lprotoc -lprotobuf

#include <stdbool.h>
#include <stdint.h>

struct Buffer {
    uintptr_t data;
    size_t size;
};

extern void FreeBuffer(struct Buffer* buf);
extern bool ParseProto(const struct Buffer* input, struct Buffer* output);
*/
import "C"

func ParseProto(data []byte) (*pb.FileDescriptorProto, error) {
	input := C.struct_Buffer{
		data: C.ulong(uintptr(unsafe.Pointer(unsafe.SliceData(data)))),
		size: C.size_t(len(data)),
	}
	output := C.struct_Buffer{}
	done := C.ParseProto(&input, &output)
	defer C.FreeBuffer(&output)
	if !bool(done) {
		msg := C.GoStringN((*C.char)(unsafe.Pointer(uintptr(output.data))), C.int(output.size))
		return nil, errors.New(msg)
	}
	raw := C.GoBytes(unsafe.Pointer(uintptr(output.data)), C.int(output.size))
	file := &pb.FileDescriptorProto{}
	err := proto.Unmarshal(raw, file)
	if err != nil {
		return nil, err
	}
	return file, nil
}

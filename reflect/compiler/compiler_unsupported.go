//go:build !linux || !cgo

package compiler

import (
	pb "google.golang.org/protobuf/types/descriptorpb"
)

func ParseProto(data []byte) (*pb.FileDescriptorProto, error) {
	return nil, ErrUnsupported
}

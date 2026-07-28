//go:build !linux || !cgo

package compiler

import (
	pb "google.golang.org/protobuf/types/descriptorpb"
)

// ParseProto returns ErrUnsupported when the native parser is unavailable.
func ParseProto(data []byte) (*pb.FileDescriptorProto, error) {
	return nil, ErrUnsupported
}

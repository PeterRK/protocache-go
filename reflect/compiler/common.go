// Package compiler parses .proto source text into protobuf file descriptors.
//
// Parsing requires Linux, cgo, libprotoc, and libprotobuf. On unsupported
// platforms ParseProto returns ErrUnsupported.
package compiler

import (
	"errors"
)

// ErrUnsupported indicates that the native .proto parser is unavailable on this platform.
var ErrUnsupported = errors.New("reflect/compiler is only supported on linux with cgo enabled")

package compiler

import (
	"errors"
)

var ErrUnsupported = errors.New("reflect/compiler is only supported on linux with cgo enabled")

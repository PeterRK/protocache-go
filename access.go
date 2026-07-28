// Package protocache provides zero-copy access to ProtoCache data, conversion
// from protobuf messages, and optional lightweight compression.
//
// ProtoCache data uses a little-endian binary format. Read-only views returned
// by this package borrow their input buffer: the buffer must remain alive and
// must not be modified while a view or a value derived from it is in use.
package protocache

import (
	"math"
	"math/bits"
	"unsafe"
)

// EnumValue is the underlying value used by generated ProtoCache enum types.
type EnumValue int32

// Enum is implemented by generated ProtoCache enum types.
type Enum interface {
	~int32
}

// CastEnumArray converts enum values to T without allocating; the result shares storage with vec.
func CastEnumArray[T Enum](vec []EnumValue) []T {
	return *(*[]T)(unsafe.Pointer(&vec))
}

func extractBytes(data []byte) []byte {
	if len(data) == 0 {
		return nil
	}
	b0 := data[0]
	if (b0 & 0x80) == 0 {
		if b0&3 != 0 {
			return nil
		}
		size := int(b0 >> 2)
		if 1+size > len(data) {
			return nil
		}
		return data[1 : 1+size]
	}

	off := 1
	mark := uint32(b0 & 0x7f)
	for sft := 7; sft < 32; sft += 7 {
		if off >= len(data) {
			return nil
		}
		b := uint32(data[off])
		off++
		if (b & 0x80) != 0 {
			mark |= (b & 0x7f) << sft
		} else {
			mark |= b << sft
			if mark&3 != 0 {
				return nil
			}
			size := int(mark >> 2)
			if off+size > len(data) {
				return nil
			}
			return data[off : off+size]
		}
	}
	return nil
}

func extractString(data []byte) string {
	raw := extractBytes(data)
	return *(*string)(unsafe.Pointer(&raw))
}

func extractBoolArray(data []byte) []bool {
	raw := extractBytes(data)
	return *(*[]bool)(unsafe.Pointer(&raw))
}

// Field is a zero-copy view of one encoded message, array, or map field; its zero value is invalid.
type Field struct {
	data []byte
}

// IsValid reports whether f refers to an encoded field.
func (f *Field) IsValid() bool {
	return len(f.data) != 0
}

// RawWords returns the encoded field words without copying; the result must be treated as read-only.
func (f *Field) RawWords() []uint32 {
	if len(f.data) == 0 {
		return nil
	}
	return castBytesToWords(f.data)
}

// GetBool decodes f as a bool, returning false when f has the wrong width.
func (f *Field) GetBool() bool {
	if len(f.data) != 4 {
		return false
	}
	return f.data[0] != 0
}

// GetEnumValue decodes f as an EnumValue, returning zero when f has the wrong width.
func (f *Field) GetEnumValue() EnumValue {
	return EnumValue(f.GetUint32())
}

// GetUint32 decodes f as a uint32, returning zero when f has the wrong width.
func (f *Field) GetUint32() uint32 {
	if len(f.data) != 4 {
		return 0
	}
	return getUint32(f.data)
}

// GetInt32 decodes f as an int32, returning zero when f has the wrong width.
func (f *Field) GetInt32() int32 {
	return int32(f.GetUint32())
}

// GetUint64 decodes f as a uint64, returning zero when f has the wrong width.
func (f *Field) GetUint64() uint64 {
	if len(f.data) != 8 {
		return 0
	}
	return getUint64(f.data)
}

// GetInt64 decodes f as an int64, returning zero when f has the wrong width.
func (f *Field) GetInt64() int64 {
	return int64(f.GetUint64())
}

// GetFloat32 decodes f as a float32, returning zero when f has the wrong width.
func (f *Field) GetFloat32() float32 {
	if len(f.data) != 4 {
		return 0
	}
	return math.Float32frombits(getUint32(f.data))
}

// GetFloat64 decodes f as a float64, returning zero when f has the wrong width.
func (f *Field) GetFloat64() float64 {
	if len(f.data) != 8 {
		return 0
	}
	return math.Float64frombits(getUint64(f.data))
}

// GetObject returns the encoded object referenced by f without copying, or nil if f is not a valid object field.
func (f *Field) GetObject() []byte {
	if len(f.data) < 4 {
		return nil
	}
	mark := getUint32(f.data)
	if (mark & 3) != 3 {
		return f.data
	}
	off := mark & 0xfffffffc
	if off >= uint32(cap(f.data)) {
		return nil
	}
	base := unsafe.Slice(unsafe.SliceData(f.data), cap(f.data))
	return base[off:]
}

// GetBytes decodes f as bytes without copying; the result borrows the original encoded buffer.
func (f *Field) GetBytes() []byte {
	return extractBytes(f.GetObject())
}

// GetString decodes f as a string without copying; the result borrows the original encoded buffer.
func (f *Field) GetString() string {
	return extractString(f.GetObject())
}

// GetBoolArray decodes f as a bool slice without copying.
func (f *Field) GetBoolArray() []bool {
	return extractBoolArray(f.GetObject())
}

// GetEnumValueArray decodes f as an enum slice without copying.
func (f *Field) GetEnumValueArray() []EnumValue {
	arr := AsArray(f.GetObject())
	return arr.EnumValue()
}

// GetInt32Array decodes f as an int32 slice without copying.
func (f *Field) GetInt32Array() []int32 {
	arr := AsArray(f.GetObject())
	return arr.Int32()
}

// GetUint32Array decodes f as a uint32 slice without copying.
func (f *Field) GetUint32Array() []uint32 {
	arr := AsArray(f.GetObject())
	return arr.Uint32()
}

// GetInt64Array decodes f as an int64 slice without copying.
func (f *Field) GetInt64Array() []int64 {
	arr := AsArray(f.GetObject())
	return arr.Int64()
}

// GetUint64Array decodes f as a uint64 slice without copying.
func (f *Field) GetUint64Array() []uint64 {
	arr := AsArray(f.GetObject())
	return arr.Uint64()
}

// GetFloat32Array decodes f as a float32 slice without copying.
func (f *Field) GetFloat32Array() []float32 {
	arr := AsArray(f.GetObject())
	return arr.Float32()
}

// GetFloat64Array decodes f as a float64 slice without copying.
func (f *Field) GetFloat64Array() []float64 {
	arr := AsArray(f.GetObject())
	return arr.Float64()
}

// GetMessage decodes f as a Message view.
func (f *Field) GetMessage() Message {
	return AsMessage(f.GetObject())
}

// GetArray decodes f as an Array view.
func (f *Field) GetArray() Array {
	return AsArray(f.GetObject())
}

// GetMap decodes f as a Map view.
func (f *Field) GetMap() Map {
	return AsMap(f.GetObject())
}

// Message is a zero-copy read-only view of an encoded ProtoCache message.
type Message struct {
	data []byte
}

func count32(v uint32) uint32 {
	return uint32(bits.OnesCount32(v) + bits.OnesCount32(v&0xaaaaaaaa))
}

func count64(v uint64) uint32 {
	return uint32(bits.OnesCount64(v) + bits.OnesCount64(v&0xaaaaaaaaaaaaaaaa))
}

// AsMessage validates the outer message header and returns a view, or an invalid Message for malformed data.
func AsMessage(data []byte) Message {
	if len(data) < 4 {
		return Message{}
	}
	section := uint32(data[0])
	if uint32(len(data)) < 4+8*section {
		return Message{}
	}
	return Message{data: data}
}

// IsValid reports whether m contains a valid outer message header.
func (m *Message) IsValid() bool {
	return len(m.data) != 0
}

// HasField reports whether the zero-based ProtoCache field id is present.
func (m *Message) HasField(id uint16) bool {
	if len(m.data) == 0 {
		return false
	}
	section := uint32(m.data[0])
	if id < 12 {
		v := getUint32(m.data) >> 8
		width := (v >> id * 2) & 3
		return width != 0
	}
	a, b := uint32((id-12)/25), uint32((id-12)%25)
	if a >= section {
		return false
	}
	v := getUint64(m.data[4+a*8:])
	width := (v >> b * 2) & 3
	return width != 0
}

func (m *Message) locateField(id uint16) (off uint32, width uint32, ok bool) {
	if len(m.data) == 0 {
		return 0, 0, false
	}
	section := uint32(m.data[0])
	off = 1 + section*2
	width = 0
	if id < 12 {
		v := getUint32(m.data) >> 8
		width = (v >> (uint32(id) << 1)) & 3
		if width == 0 {
			return 0, 0, false
		}
		off += count32(v & ^(uint32(0xffffffff) << (uint32(id) << 1)))
	} else {
		a, b := uint32((id-12)/25), uint32((id-12)%25)
		if a >= section {
			return 0, 0, false
		}
		v := getUint64(m.data[4+a*8:])
		width = uint32(v>>(b<<1)) & 3
		if width == 0 {
			return 0, 0, false
		}
		off += uint32(v >> 50)
		off += count64(v & ^(uint64(0xffffffffffffffff) << (b << 1)))
	}
	off *= 4
	width *= 4
	if off+width > uint32(len(m.data)) {
		return 0, 0, false
	}
	return off, width, true
}

// DetectInlined returns the inline portion of m without copying; it is intended for generated code.
func (m *Message) DetectInlined() []byte {
	if len(m.data) == 0 {
		return nil
	}
	section := uint16(m.data[0])
	last := uint16(11)
	if section != 0 {
		last = 12 + section*25 - 1
	}
	off := uint32((1 + uint32(section)*2) * 4)
	for {
		if pos, width, ok := m.locateField(last); ok {
			off = pos + width
			break
		}
		if last == 0 {
			break
		}
		last--
	}
	if off > uint32(len(m.data)) {
		return nil
	}
	return m.data[:off]
}

// GetField returns the field with the zero-based ProtoCache field id, or an invalid Field when absent or malformed.
func (m *Message) GetField(id uint16) Field {
	off, width, ok := m.locateField(id)
	if !ok {
		return Field{}
	}
	return Field{data: m.data[off : off+width]}
}

// Array is a zero-copy schema-independent view of an encoded ProtoCache array.
type Array struct {
	data  []byte
	size  uint32
	width uint32
}

// AsArray validates the array header and returns a view, or an invalid Array for malformed data.
func AsArray(data []byte) Array {
	if len(data) < 4 {
		return Array{}
	}
	mark := getUint32(data)
	arr := Array{data: data[4:]}
	arr.size = mark >> 2
	arr.width = (mark & 3) * 4
	if arr.width == 0 || arr.width*arr.size > uint32(len(arr.data)) {
		return Array{}
	}
	return arr
}

// IsValid reports whether a contains a valid array header and inline body.
func (a *Array) IsValid() bool {
	return a.data != nil
}

// Size returns the number of elements in a.
func (a *Array) Size() uint32 {
	return a.size
}

// Get returns element i as a Field, or an invalid Field when i is out of range.
func (a *Array) Get(i uint32) Field {
	if i >= a.size {
		return Field{}
	}
	off := i * a.width
	return Field{data: a.data[off : off+a.width]}
}

// EnumValue returns the elements as enum values without copying.
func (a *Array) EnumValue() []EnumValue {
	p := unsafe.Pointer(unsafe.SliceData(a.data))
	return unsafe.Slice((*EnumValue)(p), a.size)
}

// Int32 returns the elements as int32 values without copying; the encoded width must match int32.
func (a *Array) Int32() []int32 {
	p := unsafe.Pointer(unsafe.SliceData(a.data))
	return unsafe.Slice((*int32)(p), a.size)
}

// Uint32 returns the elements as uint32 values without copying; the encoded width must match uint32.
func (a *Array) Uint32() []uint32 {
	p := unsafe.Pointer(unsafe.SliceData(a.data))
	return unsafe.Slice((*uint32)(p), a.size)
}

// Int64 returns the elements as int64 values without copying; the encoded width must match int64.
func (a *Array) Int64() []int64 {
	p := unsafe.Pointer(unsafe.SliceData(a.data))
	return unsafe.Slice((*int64)(p), a.size)
}

// Uint64 returns the elements as uint64 values without copying; the encoded width must match uint64.
func (a *Array) Uint64() []uint64 {
	p := unsafe.Pointer(unsafe.SliceData(a.data))
	return unsafe.Slice((*uint64)(p), a.size)
}

// Float32 returns the elements as float32 values without copying; the encoded width must match float32.
func (a *Array) Float32() []float32 {
	p := unsafe.Pointer(unsafe.SliceData(a.data))
	return unsafe.Slice((*float32)(p), a.size)
}

// Float64 returns the elements as float64 values without copying; the encoded width must match float64.
func (a *Array) Float64() []float64 {
	p := unsafe.Pointer(unsafe.SliceData(a.data))
	return unsafe.Slice((*float64)(p), a.size)
}

// Map is a zero-copy view of an encoded ProtoCache map.
type Map struct {
	core     PerfectHashTable
	body     uint32
	keyWidth uint16
	valWidth uint16
}

// AsMap validates the map header and returns a view, or an invalid Map for malformed data.
func AsMap(data []byte) Map {
	m := Map{}
	if !m.core.InitFromEncoded(data) {
		return Map{}
	}
	// Map bodies are written on a 32-bit word boundary, so the payload starts
	// after the padded perfect-hash header rather than the raw header byte size.
	m.body = calcWordSize(uint32(len(m.core.data))) * 4
	m.core.data = data
	m.keyWidth = uint16(getUint32(data)>>28) & 0xc
	m.valWidth = uint16(getUint32(data)>>26) & 0xc
	if m.body+uint32(m.keyWidth+m.valWidth)*m.core.size > uint32(len(data)) {
		return Map{}
	}
	return m
}

// IsValid reports whether m contains a valid map header and inline body.
func (m *Map) IsValid() bool {
	return m.core.IsValid()
}

// Size returns the number of entries in m.
func (m *Map) Size() uint32 {
	return m.core.size
}

// Key returns entry i key, or an invalid Field when i is out of range.
func (m *Map) Key(i uint32) Field {
	if i >= m.core.size {
		return Field{}
	}
	off := m.body + i*uint32(m.keyWidth+m.valWidth)
	return Field{data: m.core.data[off : off+uint32(m.keyWidth)]}
}

// Value returns entry i value, or an invalid Field when i is out of range.
func (m *Map) Value(i uint32) Field {
	if i >= m.core.size {
		return Field{}
	}
	off := m.body + i*uint32(m.keyWidth+m.valWidth) + uint32(m.keyWidth)
	return Field{data: m.core.data[off : off+uint32(m.valWidth)]}
}

// FindByString returns the value for key, or an invalid Field when key is absent.
func (m *Map) FindByString(key string) Field {
	idx := m.core.Lookup(castStrToBytes(key))
	field := m.Key(idx)
	if field.GetString() != key {
		return Field{}
	}
	return m.Value(idx)
}

// FindByUint32 returns the value for key, or an invalid Field when key is absent.
func (m *Map) FindByUint32(key uint32) Field {
	var raw [4]byte
	putUint32(raw[:], key)
	idx := m.core.Lookup(raw[:])
	field := m.Key(idx)
	if field.GetUint32() != key {
		return Field{}
	}
	return m.Value(idx)
}

// FindByInt32 returns the value for key, or an invalid Field when key is absent.
func (m *Map) FindByInt32(key int32) Field {
	return m.FindByUint32(uint32(key))
}

// FindByUint64 returns the value for key, or an invalid Field when key is absent.
func (m *Map) FindByUint64(key uint64) Field {
	var raw [8]byte
	putUint64(raw[:], key)
	idx := m.core.Lookup(raw[:])
	field := m.Key(idx)
	if field.GetUint64() != key {
		return Field{}
	}
	return m.Value(idx)
}

// FindByInt64 returns the value for key, or an invalid Field when key is absent.
func (m *Map) FindByInt64(key int64) Field {
	return m.FindByUint64(uint64(key))
}

// BoolArray is a zero-copy view of an encoded bool array.
type BoolArray struct {
	core []bool
}

// AsBoolArray returns a bool-array view of data.
func AsBoolArray(data []byte) BoolArray {
	return BoolArray{core: extractBoolArray(data)}
}

// IsValid reports whether a was decoded successfully.
func (a *BoolArray) IsValid() bool {
	return a.core != nil
}

// Size returns the number of elements in a.
func (a *BoolArray) Size() uint32 {
	return uint32(len(a.core))
}

// Get returns element i and panics when i is out of range.
func (a *BoolArray) Get(i uint32) bool {
	return a.core[i]
}

// Raw returns all elements without copying; the result borrows the encoded buffer.
func (a *BoolArray) Raw() []bool {
	return a.core
}

// EnumArray is a zero-copy view of an encoded array of T values.
type EnumArray[T Enum] struct {
	core []T
}

// AsEnumArray returns an enum-array view of data.
func AsEnumArray[T Enum](data []byte) EnumArray[T] {
	arr := AsArray(data)
	core := arr.EnumValue()
	return EnumArray[T]{core: *(*[]T)(unsafe.Pointer(&core))}
}

// IsValid reports whether a was decoded successfully.
func (a *EnumArray[T]) IsValid() bool {
	return a.core != nil
}

// Size returns the number of elements in a.
func (a *EnumArray[T]) Size() uint32 {
	return uint32(len(a.core))
}

// Get returns element i and panics when i is out of range.
func (a *EnumArray[T]) Get(i uint32) T {
	return a.core[i]
}

// Raw returns all elements without copying; the result borrows the encoded buffer.
func (a *EnumArray[T]) Raw() []T {
	return a.core
}

type scalarArrayValue interface {
	int32 | uint32 | int64 | uint64 | float32 | float64
}

type scalarArray[T scalarArrayValue] struct {
	core []T
}

// IsValid reports whether a was decoded successfully.
func (a *scalarArray[T]) IsValid() bool {
	return a.core != nil
}

// Size returns the number of elements in a.
func (a *scalarArray[T]) Size() uint32 {
	return uint32(len(a.core))
}

// Get returns element i and panics when i is out of range.
func (a *scalarArray[T]) Get(i uint32) T {
	return a.core[i]
}

// Raw returns all elements without copying; the result borrows the encoded buffer.
func (a *scalarArray[T]) Raw() []T {
	return a.core
}

// Int32Array is a zero-copy view of an encoded int32 array.
type Int32Array = scalarArray[int32]

// AsInt32Array returns an int32-array view of data.
func AsInt32Array(data []byte) Int32Array {
	arr := AsArray(data)
	return Int32Array{core: arr.Int32()}
}

// Uint32Array is a zero-copy view of an encoded uint32 array.
type Uint32Array = scalarArray[uint32]

// AsUint32Array returns a uint32-array view of data.
func AsUint32Array(data []byte) Uint32Array {
	arr := AsArray(data)
	return Uint32Array{core: arr.Uint32()}
}

// Int64Array is a zero-copy view of an encoded int64 array.
type Int64Array = scalarArray[int64]

// AsInt64Array returns an int64-array view of data.
func AsInt64Array(data []byte) Int64Array {
	arr := AsArray(data)
	return Int64Array{core: arr.Int64()}
}

// Uint64Array is a zero-copy view of an encoded uint64 array.
type Uint64Array = scalarArray[uint64]

// AsUint64Array returns a uint64-array view of data.
func AsUint64Array(data []byte) Uint64Array {
	arr := AsArray(data)
	return Uint64Array{core: arr.Uint64()}
}

// Float32Array is a zero-copy view of an encoded float32 array.
type Float32Array = scalarArray[float32]

// AsFloat32Array returns a float32-array view of data.
func AsFloat32Array(data []byte) Float32Array {
	arr := AsArray(data)
	return Float32Array{core: arr.Float32()}
}

// Float64Array is a zero-copy view of an encoded float64 array.
type Float64Array = scalarArray[float64]

// AsFloat64Array returns a float64-array view of data.
func AsFloat64Array(data []byte) Float64Array {
	arr := AsArray(data)
	return Float64Array{core: arr.Float64()}
}

// StringArray is a zero-copy view of an encoded string array.
type StringArray struct {
	core Array
}

// AsStringArray returns a string-array view of data.
func AsStringArray(data []byte) StringArray {
	return StringArray{core: AsArray(data)}
}

// IsValid reports whether a was decoded successfully.
func (a *StringArray) IsValid() bool {
	return a.core.IsValid()
}

// Size returns the number of elements in a.
func (a *StringArray) Size() uint32 {
	return a.core.Size()
}

// Get returns element i, or an empty string when i is out of range.
func (a *StringArray) Get(i uint32) string {
	field := a.core.Get(i)
	return field.GetString()
}

// BytesArray is a zero-copy view of an encoded bytes array.
type BytesArray struct {
	core Array
}

// AsBytesArray returns a bytes-array view of data.
func AsBytesArray(data []byte) BytesArray {
	return BytesArray{core: AsArray(data)}
}

// IsValid reports whether a was decoded successfully.
func (a *BytesArray) IsValid() bool {
	return a.core.IsValid()
}

// Size returns the number of elements in a.
func (a *BytesArray) Size() uint32 {
	return a.core.Size()
}

// Get returns element i without copying, or nil when i is out of range.
func (a *BytesArray) Get(i uint32) []byte {
	field := a.core.Get(i)
	return field.GetBytes()
}

// DetectBytes returns the complete encoded bytes object, or nil for malformed data; it is intended for generated code.
func DetectBytes(data []byte) []byte {
	raw := extractBytes(data)
	if raw == nil {
		return nil
	}
	head := int(uintptr(unsafe.Pointer(unsafe.SliceData(raw))) - uintptr(unsafe.Pointer(unsafe.SliceData(data))))
	size := (head + len(raw) + 3) &^ 3
	if size > len(data) {
		return nil
	}
	return data[:size]
}

// DetectObject returns the referenced non-inline object bytes, or nil for invalid, scalar, or inline fields.
func (f *Field) DetectObject() []byte {
	obj := f.GetObject()
	if obj == nil || unsafe.SliceData(obj) == unsafe.SliceData(f.data) {
		return nil
	}
	return obj
}

// DetectShrink returns data truncated after part, or nil when the ranges are inconsistent; it is intended for generated code.
func DetectShrink(data, obj, part []byte) []byte {
	if len(part) == 0 {
		return nil
	}
	tail := int(uintptr(unsafe.Pointer(unsafe.SliceData(obj)))-uintptr(unsafe.Pointer(unsafe.SliceData(data)))) + len(part)
	if tail > len(data) {
		return nil
	}
	return data[:tail]
}

// DetectArray returns the complete encoded array, or nil for malformed data; it is intended for generated code.
func DetectArray(data []byte, detect func([]byte) []byte) []byte {
	a := AsArray(data)
	if !a.IsValid() {
		return nil
	}
	compactEnd := 4 + int(a.size*a.width)
	if detect == nil {
		return data[:compactEnd]
	}
	for i := a.size; i > 0; i-- {
		field := a.Get(i - 1)
		obj := field.DetectObject()
		if obj == nil {
			continue
		}
		base := 4 + int((i-1)*a.width)
		off := int(uintptr(unsafe.Pointer(unsafe.SliceData(obj))) - uintptr(unsafe.Pointer(unsafe.SliceData(field.data))))
		part := detect(obj)
		if len(part) == 0 {
			return nil
		}
		tail := base + off + len(part)
		if tail > len(data) {
			return nil
		}
		return data[:tail]
	}
	return data[:compactEnd]
}

// DetectMap returns the complete encoded map, or nil for malformed data; it is intended for generated code.
func DetectMap(data []byte, detectKey func([]byte) []byte, detectValue func([]byte) []byte) []byte {
	m := AsMap(data)
	if !m.IsValid() {
		return nil
	}
	compactEnd := int(m.body + uint32(m.keyWidth+m.valWidth)*m.core.size)
	if detectKey == nil && detectValue == nil {
		return m.core.data[:compactEnd]
	}
	for i := m.core.size; i > 0; i-- {
		idx := i - 1
		if detectValue != nil {
			field := m.Value(idx)
			obj := field.DetectObject()
			if obj != nil {
				base := int(m.body + idx*uint32(m.keyWidth+m.valWidth) + uint32(m.keyWidth))
				off := int(uintptr(unsafe.Pointer(unsafe.SliceData(obj))) - uintptr(unsafe.Pointer(unsafe.SliceData(field.data))))
				part := detectValue(obj)
				if len(part) == 0 {
					return nil
				}
				tail := base + off + len(part)
				if tail > len(m.core.data) {
					return nil
				}
				return m.core.data[:tail]
			}
		}
		if detectKey != nil {
			field := m.Key(idx)
			obj := field.DetectObject()
			if obj != nil {
				base := int(m.body + idx*uint32(m.keyWidth+m.valWidth))
				off := int(uintptr(unsafe.Pointer(unsafe.SliceData(obj))) - uintptr(unsafe.Pointer(unsafe.SliceData(field.data))))
				part := detectKey(obj)
				if len(part) == 0 {
					return nil
				}
				tail := base + off + len(part)
				if tail > len(m.core.data) {
					return nil
				}
				return m.core.data[:tail]
			}
		}
	}
	return m.core.data[:compactEnd]
}

// CheckVisited reports whether id is set in bitmap; it is intended for generated code and requires a sufficiently large bitmap.
func CheckVisited(bitmap []byte, id uint16) bool {
	return (bitmap[id>>3] & byte(1<<(id&7))) != 0
}

// Visit sets id in bitmap; it is intended for generated code and requires a sufficiently large bitmap.
func Visit(bitmap []byte, id uint16) {
	bitmap[id>>3] |= byte(1 << (id & 7))
}

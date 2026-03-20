package protocache

import (
	"math"
	"unsafe"
)

type EnumValue int32

type Enum interface {
	~int32
}

func CastEnumArray[T Enum](vec []EnumValue) []T {
	return *(*[]T)(unsafe.Pointer(&vec))
}

func fullSizeSlice[T any](data []T) []T {
	return unsafe.Slice(unsafe.SliceData(data), cap(data))
}

func extractBytes(data []byte) []byte {
	off := 0
	mark := uint32(0)
	for sft := 0; sft < 32; sft += 7 {
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

type Field struct {
	data []byte
}

func (f Field) IsValid() bool {
	return len(f.data) != 0
}

func (f Field) RawWords() []uint32 {
	if len(f.data) == 0 {
		return nil
	}
	return castBytesToWords(f.data)
}

func (f Field) GetBool() bool {
	if len(f.data) != 4 {
		return false
	}
	return f.data[0] != 0
}

func (f Field) GetEnumValue() EnumValue {
	return EnumValue(f.GetUint32())
}

func (f Field) GetUint32() uint32 {
	if len(f.data) != 4 {
		return 0
	}
	return getUint32(f.data)
}

func (f Field) GetInt32() int32 {
	return int32(f.GetUint32())
}

func (f Field) GetUint64() uint64 {
	if len(f.data) != 8 {
		return 0
	}
	return getUint64(f.data)
}

func (f Field) GetInt64() int64 {
	return int64(f.GetUint64())
}

func (f Field) GetFloat32() float32 {
	if len(f.data) != 4 {
		return 0
	}
	return math.Float32frombits(getUint32(f.data))
}

func (f Field) GetFloat64() float64 {
	if len(f.data) != 8 {
		return 0
	}
	return math.Float64frombits(getUint64(f.data))
}

func (f Field) GetObject() []byte {
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
	return fullSizeSlice(f.data)[off:]
}

func (f Field) GetBytes() []byte {
	return extractBytes(f.GetObject())
}

func (f Field) GetString() string {
	return extractString(f.GetObject())
}

func (f Field) GetBoolArray() []bool {
	return extractBoolArray(f.GetObject())
}

func (f Field) GetEnumValueArray() []EnumValue {
	arr := AsArray(f.GetObject())
	return arr.EnumValue()
}

func (f Field) GetInt32Array() []int32 {
	arr := AsArray(f.GetObject())
	return arr.Int32()
}

func (f Field) GetUint32Array() []uint32 {
	arr := AsArray(f.GetObject())
	return arr.Uint32()
}

func (f Field) GetInt64Array() []int64 {
	arr := AsArray(f.GetObject())
	return arr.Int64()
}

func (f Field) GetUint64Array() []uint64 {
	arr := AsArray(f.GetObject())
	return arr.Uint64()
}

func (f Field) GetFloat32Array() []float32 {
	arr := AsArray(f.GetObject())
	return arr.Float32()
}

func (f Field) GetFloat64Array() []float64 {
	arr := AsArray(f.GetObject())
	return arr.Float64()
}

func (f Field) GetMessage() Message {
	return AsMessage(f.GetObject())
}

func (f Field) GetArray() Array {
	return AsArray(f.GetObject())
}

func (f Field) GetMap() Map {
	return AsMap(f.GetObject())
}

type Message struct {
	data []byte
}

func count32(v uint32) uint32 {
	v = (v & 0x33333333) + ((v >> 2) & 0x33333333)
	v = v + (v >> 4)
	v = (v & 0x0f0f0f0f) + ((v >> 8) & 0x0f0f0f0f)
	v = v + (v >> 16)
	return v & 0xff
}

func count64(v uint64) uint32 {
	v = (v & 0x3333333333333333) + ((v >> 2) & 0x3333333333333333)
	v = v + (v >> 4)
	v = (v & 0x0f0f0f0f0f0f0f0f) + ((v >> 8) & 0x0f0f0f0f0f0f0f0f)
	v = v + (v >> 16)
	v = v + (v >> 32)
	return uint32(v) & 0xff
}

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

func (m *Message) IsValid() bool {
	return len(m.data) != 0
}

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

func (m *Message) GetField(id uint16) Field {
	off, width, ok := m.locateField(id)
	if !ok {
		return Field{}
	}
	return Field{data: m.data[off : off+width]}
}

type Array struct {
	data  []byte
	size  uint32
	width uint32
}

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

func (a *Array) IsValid() bool {
	return a.data != nil
}

func (a *Array) Size() uint32 {
	return a.size
}

func (a *Array) Get(i uint32) Field {
	if i >= a.size {
		return Field{}
	}
	off := i * a.width
	return Field{data: a.data[off : off+a.width]}
}

func (a *Array) EnumValue() []EnumValue {
	p := unsafe.Pointer(unsafe.SliceData(a.data))
	return unsafe.Slice((*EnumValue)(p), a.size)
}

func (a *Array) Int32() []int32 {
	p := unsafe.Pointer(unsafe.SliceData(a.data))
	return unsafe.Slice((*int32)(p), a.size)
}

func (a *Array) Uint32() []uint32 {
	p := unsafe.Pointer(unsafe.SliceData(a.data))
	return unsafe.Slice((*uint32)(p), a.size)
}

func (a *Array) Int64() []int64 {
	p := unsafe.Pointer(unsafe.SliceData(a.data))
	return unsafe.Slice((*int64)(p), a.size)
}

func (a *Array) Uint64() []uint64 {
	p := unsafe.Pointer(unsafe.SliceData(a.data))
	return unsafe.Slice((*uint64)(p), a.size)
}

func (a *Array) Float32() []float32 {
	p := unsafe.Pointer(unsafe.SliceData(a.data))
	return unsafe.Slice((*float32)(p), a.size)
}

func (a *Array) Float64() []float64 {
	p := unsafe.Pointer(unsafe.SliceData(a.data))
	return unsafe.Slice((*float64)(p), a.size)
}

type Map struct {
	core     PerfectHash
	body     uint32
	keyWidth uint16
	valWidth uint16
}

func AsMap(data []byte) Map {
	m := Map{}
	if !m.core.Init(data) {
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

func (m *Map) IsValid() bool {
	return m.core.IsValid()
}

func (m *Map) Size() uint32 {
	return m.core.size
}

func (m *Map) Key(i uint32) Field {
	if i >= m.core.size {
		return Field{}
	}
	off := m.body + i*uint32(m.keyWidth+m.valWidth)
	return Field{data: m.core.data[off : off+uint32(m.keyWidth)]}
}

func (m *Map) Value(i uint32) Field {
	if i >= m.core.size {
		return Field{}
	}
	off := m.body + i*uint32(m.keyWidth+m.valWidth) + uint32(m.keyWidth)
	return Field{data: m.core.data[off : off+uint32(m.valWidth)]}
}

func (m *Map) FindByString(key string) Field {
	idx := m.core.Locate(castStrToBytes(key))
	field := m.Key(idx)
	if field.GetString() != key {
		return Field{}
	}
	return m.Value(idx)
}

func (m *Map) FindByUint32(key uint32) Field {
	var raw [4]byte
	putUint32(raw[:], key)
	idx := m.core.Locate(raw[:])
	field := m.Key(idx)
	if field.GetUint32() != key {
		return Field{}
	}
	return m.Value(idx)
}

func (m *Map) FindByInt32(key int32) Field {
	return m.FindByUint32(uint32(key))
}

func (m *Map) FindByUint64(key uint64) Field {
	var raw [8]byte
	putUint64(raw[:], key)
	idx := m.core.Locate(raw[:])
	field := m.Key(idx)
	if field.GetUint64() != key {
		return Field{}
	}
	return m.Value(idx)
}

func (m *Map) FindByInt64(key int64) Field {
	return m.FindByUint64(uint64(key))
}

type BoolArray struct {
	core []bool
}

func AsBoolArray(data []byte) BoolArray {
	return BoolArray{core: extractBoolArray(data)}
}

func (a *BoolArray) IsValid() bool {
	return a.core != nil
}

func (a *BoolArray) Size() uint32 {
	return uint32(len(a.core))
}

func (a *BoolArray) Get(i uint32) bool {
	return a.core[i]
}

func (a *BoolArray) Raw() []bool {
	return a.core
}

type EnumArray[T Enum] struct {
	core []T
}

func AsEnumArray[T Enum](data []byte) EnumArray[T] {
	arr := AsArray(data)
	core := arr.EnumValue()
	return EnumArray[T]{core: *(*[]T)(unsafe.Pointer(&core))}
}

func (a *EnumArray[T]) IsValid() bool {
	return a.core != nil
}

func (a *EnumArray[T]) Size() uint32 {
	return uint32(len(a.core))
}

func (a *EnumArray[T]) Get(i uint32) T {
	return a.core[i]
}

func (a *EnumArray[T]) Raw() []T {
	return a.core
}

type Int32Array struct {
	core []int32
}

func AsInt32Array(data []byte) Int32Array {
	arr := AsArray(data)
	return Int32Array{core: arr.Int32()}
}

func (a *Int32Array) IsValid() bool {
	return a.core != nil
}

func (a *Int32Array) Size() uint32 {
	return uint32(len(a.core))
}

func (a *Int32Array) Get(i uint32) int32 {
	return a.core[i]
}

func (a *Int32Array) Raw() []int32 {
	return a.core
}

type Uint32Array struct {
	core []uint32
}

func AsUint32Array(data []byte) Uint32Array {
	arr := AsArray(data)
	return Uint32Array{core: arr.Uint32()}
}

func (a *Uint32Array) IsValid() bool {
	return a.core != nil
}

func (a *Uint32Array) Size() uint32 {
	return uint32(len(a.core))
}

func (a *Uint32Array) Get(i uint32) uint32 {
	return a.core[i]
}

func (a *Uint32Array) Raw() []uint32 {
	return a.core
}

type Int64Array struct {
	core []int64
}

func AsInt64Array(data []byte) Int64Array {
	arr := AsArray(data)
	return Int64Array{core: arr.Int64()}
}

func (a *Int64Array) IsValid() bool {
	return a.core != nil
}

func (a *Int64Array) Size() uint32 {
	return uint32(len(a.core))
}

func (a *Int64Array) Get(i uint32) int64 {
	return a.core[i]
}

func (a *Int64Array) Raw() []int64 {
	return a.core
}

type Uint64Array struct {
	core []uint64
}

func AsUint64Array(data []byte) Uint64Array {
	arr := AsArray(data)
	return Uint64Array{core: arr.Uint64()}
}

func (a *Uint64Array) IsValid() bool {
	return a.core != nil
}

func (a *Uint64Array) Size() uint32 {
	return uint32(len(a.core))
}

func (a *Uint64Array) Get(i uint32) uint64 {
	return a.core[i]
}

func (a *Uint64Array) Raw() []uint64 {
	return a.core
}

type Float32Array struct {
	core []float32
}

func AsFloat32Array(data []byte) Float32Array {
	arr := AsArray(data)
	return Float32Array{core: arr.Float32()}
}

func (a *Float32Array) IsValid() bool {
	return a.core != nil
}

func (a *Float32Array) Size() uint32 {
	return uint32(len(a.core))
}

func (a *Float32Array) Get(i uint32) float32 {
	return a.core[i]
}

func (a *Float32Array) Raw() []float32 {
	return a.core
}

type Float64Array struct {
	core []float64
}

func AsFloat64Array(data []byte) Float64Array {
	arr := AsArray(data)
	return Float64Array{core: arr.Float64()}
}

func (a *Float64Array) IsValid() bool {
	return a.core != nil
}

func (a *Float64Array) Size() uint32 {
	return uint32(len(a.core))
}

func (a *Float64Array) Get(i uint32) float64 {
	return a.core[i]
}

func (a *Float64Array) Raw() []float64 {
	return a.core
}

type StringArray struct {
	core Array
}

func AsStringArray(data []byte) StringArray {
	return StringArray{core: AsArray(data)}
}

func (a *StringArray) IsValid() bool {
	return a.core.IsValid()
}

func (a *StringArray) Size() uint32 {
	return a.core.Size()
}

func (a *StringArray) Get(i uint32) string {
	field := a.core.Get(i)
	return field.GetString()
}

type BytesArray struct {
	core Array
}

func AsBytesArray(data []byte) BytesArray {
	return BytesArray{core: AsArray(data)}
}

func (a *BytesArray) IsValid() bool {
	return a.core.IsValid()
}

func (a *BytesArray) Size() uint32 {
	return a.core.Size()
}

func (a *BytesArray) Get(i uint32) []byte {
	field := a.core.Get(i)
	return field.GetBytes()
}

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

// DetectObject returns the referenced object bytes for a non-inline object field.
// It returns nil for invalid fields, scalar fields, and inline object fields.
func (f Field) DetectObject() []byte {
	obj := f.GetObject()
	if obj == nil || unsafe.SliceData(obj) == unsafe.SliceData(f.data) {
		return nil
	}
	return obj
}

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

package protocache

import (
	"errors"
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func Serialize(obj proto.Message) ([]byte, error) {
	data, err := encodeMessage(obj.ProtoReflect())
	if err != nil {
		return nil, err
	}
	return castToBytes(data), nil
}

func SerializeWords(data []uint32, err error) ([]byte, error) {
	return WordsToBytes(data), err
}

func EncodeScalar[T scalar](v T) ([]uint32, error) {
	return encodeScalar(v), nil
}

func EncodeBool(v bool) ([]uint32, error) {
	return encodeBool(v), nil
}

func EncodeBytes(data []byte) ([]uint32, error) {
	return encodeBytes(data)
}

func EncodeString(str string) ([]uint32, error) {
	return encodeString(str)
}

func EncodeMessageParts(parts [][]uint32) ([]uint32, error) {
	return encodeMessageParts(parts)
}

func EncodeScalarVector[T scalar](vec []T) ([]uint32, error) {
	return encodeScalarVector(vec)
}

func EncodeBoolArray(vec []bool) ([]uint32, error) {
	tmp := make([]byte, len(vec))
	for i, one := range vec {
		if one {
			tmp[i] = 1
		}
	}
	return encodeBytes(tmp)
}

func EncodeEnumArray[T Enum](vec []T) ([]uint32, error) {
	return EncodeScalarVector(upCast[T, int32](vec))
}

func EncodeStringArray(vec []string) ([]uint32, error) {
	return encodeArray(len(vec), func(i int) ([]uint32, error) {
		return encodeString(vec[i])
	})
}

func EncodeBytesArray(vec [][]byte) ([]uint32, error) {
	return encodeArray(len(vec), func(i int) ([]uint32, error) {
		return encodeBytes(vec[i])
	})
}

func EncodeObjectArray[T any](vec []T, encoder func(T) ([]uint32, error)) ([]uint32, error) {
	return encodeArray(len(vec), func(i int) ([]uint32, error) {
		return encoder(vec[i])
	})
}

type mapParts struct {
	key []uint32
	val []uint32
}

func collectMapParts[K comparable, V any](x map[K]V,
	keyEnc func(K) ([]uint32, error),
	valEnc func(V) ([]uint32, error)) ([]mapParts, error) {
	parts := make([]mapParts, 0, len(x))
	for k, v := range x {
		keyPart, err := keyEnc(k)
		if err != nil {
			return nil, err
		}
		valPart, err := valEnc(v)
		if err != nil {
			return nil, err
		}
		parts = append(parts, mapParts{key: keyPart, val: valPart})
	}
	return parts, nil
}

func EncodeScalarMap[K scalar, V any](x map[K]V,
	keyEnc func(K) ([]uint32, error), valEnc func(V) ([]uint32, error)) ([]uint32, error) {
	if len(x) == 0 {
		return []uint32{5 << 28}, nil
	}
	parts, err := collectMapParts(x, keyEnc, valEnc)
	if err != nil {
		return nil, err
	}
	return encodeMapParts(parts, false)
}

func EncodeStringMap[V any](x map[string]V,
	encoder func(V) ([]uint32, error)) ([]uint32, error) {
	if len(x) == 0 {
		return []uint32{5 << 28}, nil
	}
	parts, err := collectMapParts(x, EncodeString, encoder)
	if err != nil {
		return nil, err
	}
	return encodeMapParts(parts, true)
}

func BytesToWords(data []byte) []uint32 {
	if len(data) == 0 {
		return nil
	}
	if (len(data) & 3) != 0 {
		panic("unaligned field data")
	}
	return castBytesToWords(data)
}

func WordsToBytes(data []uint32) []byte {
	if len(data) == 0 {
		return nil
	}
	return castToBytes(data)
}

func calcWordSize(size uint32) uint32 {
	return (size + 3) / 4
}

func calcOffset(off uint32) uint32 {
	return (off << 2) | 3
}

func encodeMessage(message protoreflect.Message) ([]uint32, error) {
	descriptor := message.Descriptor()
	originFields := descriptor.Fields()
	if originFields.Len() <= 0 {
		return nil, fmt.Errorf("no fields in %s", descriptor.FullName())
	}
	maxId := 1
	for i := 0; i < originFields.Len(); i++ {
		field := originFields.Get(i)
		if field == nil || field.Number() <= 0 {
			return nil, fmt.Errorf("illegal field in %s", descriptor.FullName())
		}
		if maxId < int(field.Number()) {
			maxId = int(field.Number())
		}
	}
	if maxId > (12 + 25*255) {
		return nil, fmt.Errorf("too many fields in %s", descriptor.FullName())
	} else if maxId-originFields.Len() > 6 && maxId > originFields.Len()*2 {
		return nil, fmt.Errorf("message %s is too sparse", descriptor.FullName())
	}

	fields := make([]protoreflect.FieldDescriptor, maxId)
	for i := 0; i < originFields.Len(); i++ {
		field := originFields.Get(i)
		j := field.Number() - 1
		if fields[j] != nil {
			return nil, fmt.Errorf("duplicate field id %d in %s", field.Number(), descriptor.FullName())
		}
		fields[j] = field
	}
	parts := make([][]uint32, len(fields))
	for i, field := range fields {
		if field == nil || !message.Has(field) {
			continue
		}
		var err error
		if field.IsMap() {
			parts[i], err = encodeMap(field, message.Get(field).Map())
		} else if field.IsList() {
			parts[i], err = encodeList(field, message.Get(field).List())
		} else {
			parts[i], err = encodeField(field, message.Get(field))
			if len(parts[i]) == 1 && field.Kind() == protoreflect.MessageKind {
				parts[i] = nil
			}
		}
		if err != nil {
			return nil, err
		}
	}

	if len(fields) == 1 && fields[0].Name() == "_" {
		if len(parts[0]) == 0 {
			if fields[0].IsMap() {
				parts[0] = []uint32{5 << 28}
			} else {
				parts[0] = []uint32{1}
			}
		}
		return parts[0], nil
	}
	return encodeMessageParts(parts)
}

func encodeMessageParts(parts [][]uint32) ([]uint32, error) {
	for len(parts) != 0 && parts[len(parts)-1] == nil {
		parts = parts[:len(parts)-1]
	}
	if len(parts) == 0 {
		return []uint32{0}, nil
	}

	section := (uint32(len(parts)) + 12) / 25
	size := 1 + section*2

	cnt := uint32(0)
	head := section
	n := uint32(12)
	if uint32(len(parts)) < n {
		n = uint32(len(parts))
	}
	for i := uint32(0); i < n; i++ {
		one := parts[i]
		if len(one) < 4 {
			head |= uint32(len(one)) << (8 + i*2)
			size += uint32(len(one))
			cnt += uint32(len(one))
		} else {
			head |= 1 << (8 + i*2)
			size += 1 + uint32(len(one))
			cnt += 1
		}
	}
	for i := 12; i < len(parts); i++ {
		one := parts[i]
		if len(one) < 4 {
			size += uint32(len(one))
		} else {
			size += 1 + uint32(len(one))
		}
		if size >= (1 << 30) {
			return nil, errors.New("message size overflow")
		}
	}

	out := make([]uint32, 1+section*2, size)
	out[0] = head
	table := upCast[uint32, uint64](out[1:])
	for i, k := 12, 0; i < len(parts); k++ {
		next := i + 25
		if next > len(parts) {
			next = len(parts)
		}
		if cnt >= (1 << 14) {
			return nil, errors.New("message table overflow")
		}
		mark := uint64(cnt) << 50
		for j := 0; i < next; j += 2 {
			one := parts[i]
			i++
			if len(one) < 4 {
				mark |= uint64(len(one)) << j
				cnt += uint32(len(one))
			} else {
				mark |= 1 << j
				cnt += 1
			}
		}
		table[k] = mark
	}

	off := uint32(len(out))
	for _, one := range parts {
		if len(one) == 0 {
			continue
		}
		if len(one) < 4 {
			out = append(out, one...)
		} else {
			out = append(out, 0)
		}
	}
	for _, one := range parts {
		if len(one) < 4 {
			off += uint32(len(one))
		} else {
			out[off] = calcOffset(uint32(len(out)) - off)
			out = append(out, one...)
			off++
		}
	}
	if size != uint32(len(out)) {
		panic("size mismatch")
	}
	return out, nil
}

type scalar interface {
	int32 | int64 | uint32 | uint64 | float32 | float64
}

func encodeScalar[T scalar](v T) []uint32 {
	out := []T{v}
	return downCast[T, uint32](out)
}

func encodeBool(v bool) []uint32 {
	out := make([]uint32, 1)
	if v {
		out[0] = 1
	}
	return out
}

func encodeBytes(data []byte) ([]uint32, error) {
	if len(data) >= (1 << 30) {
		return nil, errors.New("too long string")
	}
	var tmp [5]byte
	header := tmp[:0]
	mark := uint32(len(data) << 2)
	for ; (mark & 0xffffff80) != 0; mark >>= 7 {
		header = append(header, byte(0x80|(mark&0x7f)))
	}
	header = append(header, byte(mark))
	out := make([]uint32, calcWordSize(uint32(len(header)+len(data))))
	buf := castToBytes(out)
	copy(buf[:len(header)], header)
	copy(buf[len(header):], data)
	return out, nil
}

func encodeString(str string) ([]uint32, error) {
	return encodeBytes(castStrToBytes(str))
}

func encodeField(field protoreflect.FieldDescriptor, value protoreflect.Value) ([]uint32, error) {
	switch field.Kind() {
	case protoreflect.MessageKind:
		return encodeMessage(value.Message())
	case protoreflect.BytesKind:
		return encodeBytes(value.Bytes())
	case protoreflect.StringKind:
		return encodeString(value.String())
	case protoreflect.DoubleKind:
		return encodeScalar(value.Float()), nil
	case protoreflect.FloatKind:
		return encodeScalar(float32(value.Float())), nil
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return encodeScalar(value.Uint()), nil
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return encodeScalar(value.Int()), nil
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return encodeScalar(uint32(value.Uint())), nil
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return encodeScalar(int32(value.Int())), nil
	case protoreflect.BoolKind:
		return encodeBool(value.Bool()), nil
	case protoreflect.EnumKind:
		return encodeScalar(int32(value.Enum())), nil
	default:
		return nil, fmt.Errorf("unsupported field: %s", field.FullName())
	}
}

func encodeScalarArray[T scalar](size int, get func(i int) T) ([]uint32, error) {
	m := sizeof[T]() / 4
	out := make([]uint32, 1+size*m)
	if len(out) >= (1 << 30) {
		return nil, errors.New("array size overflow")
	}
	out[0] = uint32((size << 2) | m)
	vec := upCast[uint32, T](out[1:])
	for i := 0; i < size; i++ {
		vec[i] = get(i)
	}
	return out, nil
}

func encodeScalarVector[T scalar](src []T) ([]uint32, error) {
	m := sizeof[T]() / 4
	out := make([]uint32, 1+len(src)*m)
	if len(out) >= (1 << 30) {
		return nil, errors.New("array size overflow")
	}
	out[0] = uint32((len(src) << 2) | m)
	vec := upCast[uint32, T](out[1:])
	copy(vec, src)
	return out, nil
}

func bestArraySize(parts [][]uint32) (size, width int) {
	sizes := [3]int{0, 0, 0}
	for _, one := range parts {
		sizes[0] += 1
		sizes[1] += 2
		sizes[2] += 3
		if len(one) <= 1 {
			continue
		}
		sizes[0] += len(one)
		if len(one) <= 2 {
			continue
		}
		sizes[1] += len(one)
		if len(one) <= 3 {
			continue
		}
		sizes[2] += len(one)
	}
	mode := 0
	for i := 1; i < 3; i++ {
		if sizes[i] < sizes[mode] {
			mode = i
		}
	}
	return sizes[mode], mode + 1
}

func bestKvArraySize(parts []mapParts, order []int, key bool) (size, width int) {
	sizes := [3]int{0, 0, 0}
	for _, idx := range order {
		one := parts[idx].val
		if key {
			one = parts[idx].key
		}
		sizes[0] += 1
		sizes[1] += 2
		sizes[2] += 3
		if len(one) <= 1 {
			continue
		}
		sizes[0] += len(one)
		if len(one) <= 2 {
			continue
		}
		sizes[1] += len(one)
		if len(one) <= 3 {
			continue
		}
		sizes[2] += len(one)
	}
	mode := 0
	for i := 1; i < 3; i++ {
		if sizes[i] < sizes[mode] {
			mode = i
		}
	}
	return sizes[mode], mode + 1
}

func encodeArray(size int, get func(i int) ([]uint32, error)) ([]uint32, error) {
	parts := make([][]uint32, size)
	var err error
	for i := 0; i < size; i++ {
		parts[i], err = get(i)
		if err != nil {
			return nil, err
		}
	}
	n, m := bestArraySize(parts)
	n += 1
	if n >= (1 << 30) {
		return nil, errors.New("array size overflow")
	}
	out := make([]uint32, 1+size*m, n)
	out[0] = uint32((size << 2) | m)

	off := 1
	for _, one := range parts {
		if len(one) <= m {
			copy(out[off:], one)
		}
		off += m
	}
	off = 1
	for _, one := range parts {
		if len(one) > m {
			out[off] = calcOffset(uint32(len(out) - off))
			out = append(out, one...)
		}
		off += m
	}
	if n != len(out) {
		panic("size mismatch")
	}
	return out, nil
}

func encodeList(field protoreflect.FieldDescriptor, list protoreflect.List) ([]uint32, error) {
	switch field.Kind() {
	case protoreflect.MessageKind:
		return encodeArray(list.Len(), func(i int) ([]uint32, error) {
			return encodeMessage(list.Get(i).Message())
		})
	case protoreflect.BytesKind:
		return encodeArray(list.Len(), func(i int) ([]uint32, error) {
			return encodeBytes(list.Get(i).Bytes())
		})
	case protoreflect.StringKind:
		return encodeArray(list.Len(), func(i int) ([]uint32, error) {
			return encodeString(list.Get(i).String())
		})
	case protoreflect.DoubleKind:
		return encodeScalarArray(list.Len(), func(i int) float64 {
			return list.Get(i).Float()
		})
	case protoreflect.FloatKind:
		return encodeScalarArray(list.Len(), func(i int) float32 {
			return float32(list.Get(i).Float())
		})
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return encodeScalarArray(list.Len(), func(i int) uint64 {
			return list.Get(i).Uint()
		})
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return encodeScalarArray(list.Len(), func(i int) int64 {
			return list.Get(i).Int()
		})
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return encodeScalarArray(list.Len(), func(i int) uint32 {
			return uint32(list.Get(i).Uint())
		})
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return encodeScalarArray(list.Len(), func(i int) int32 {
			return int32(list.Get(i).Int())
		})
	case protoreflect.BoolKind:
		tmp := make([]byte, list.Len())
		for i := 0; i < len(tmp); i++ {
			if list.Get(i).Bool() {
				tmp[i] = 1
			}
		}
		return encodeBytes(tmp)
	case protoreflect.EnumKind:
		return encodeScalarArray(list.Len(), func(i int) int32 {
			return int32(list.Get(i).Enum())
		})
	default:
		return nil, fmt.Errorf("unsupported field: %s", field.FullName())
	}
	return nil, nil
}

type arrayReader struct {
	parts []mapParts
	curr  int
}

func (r *arrayReader) Reset() {
	r.curr = 0
}

func (r *arrayReader) Total() int {
	return len(r.parts)
}

type scalarReader struct{ arrayReader }

func (r *scalarReader) Next() []byte {
	key := castToBytes(r.parts[r.curr].key)
	r.curr++
	return key
}

type stringReader struct{ arrayReader }

func (r *stringReader) Next() []byte {
	key := castToBytes(r.parts[r.curr].key)
	r.curr++
	return extractBytes(key)
}

func encodeMap(field protoreflect.FieldDescriptor, pack protoreflect.Map) ([]uint32, error) {
	kField := field.MapKey()
	vField := field.MapValue()
	size := pack.Len()
	parts := make([]mapParts, 0, size)

	var err error
	pack.Range(func(key protoreflect.MapKey, val protoreflect.Value) bool {
		entry := mapParts{}
		entry.key, err = encodeField(kField, key.Value())
		if err != nil {
			return false
		}
		entry.val, err = encodeField(vField, val)
		if err != nil {
			return false
		}
		parts = append(parts, entry)
		return true
	})
	if err != nil {
		return nil, err
	}

	return encodeMapParts(parts, kField.Kind() == protoreflect.StringKind)
}

func encodeMapParts(parts []mapParts, stringKey bool) ([]uint32, error) {
	size := len(parts)
	var index perfectHashTable
	order := make([]int, size)
	build := func(src hashKeySource) {
		index = buildPerfectHashTable(src)
		if !index.isValid() {
			return
		}
		src.Reset()
		for i := 0; i < size; i++ {
			order[index.lookup(src.Next())] = i
		}
	}

	if stringKey {
		build(&stringReader{arrayReader: arrayReader{parts: parts}})
	} else {
		build(&scalarReader{arrayReader: arrayReader{parts: parts}})
	}
	if !index.isValid() {
		return nil, errors.New("fail to build map")
	}

	n0 := int(calcWordSize(uint32(len(index.encodedBytes()))))
	n1, m1 := bestKvArraySize(parts, order, true)
	n2, m2 := bestKvArraySize(parts, order, false)

	n := n0 + n1 + n2
	if n >= (1 << 30) {
		return nil, errors.New("map size overflow")
	}

	out := make([]uint32, n0+size*(m1+m2), n)
	copy(castToBytes(out), index.encodedBytes())
	out[0] |= uint32((m1 << 30) | (m2 << 28))

	off := n0
	for _, idx := range order {
		if key := parts[idx].key; len(key) <= m1 {
			copy(out[off:], key)
		}
		off += m1
		if val := parts[idx].val; len(val) <= m2 {
			copy(out[off:], val)
		}
		off += m2
	}

	off = n0
	for _, idx := range order {
		if key := parts[idx].key; len(key) > m1 {
			out[off] = calcOffset(uint32(len(out) - off))
			out = append(out, key...)
		}
		off += m1
		if val := parts[idx].val; len(val) > m2 {
			out[off] = calcOffset(uint32(len(out) - off))
			out = append(out, val...)
		}
		off += m2
	}
	if n != len(out) {
		panic("size mismatch")
	}
	return out, nil
}

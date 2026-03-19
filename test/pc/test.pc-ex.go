package pc

import (
	"github.com/peterrk/protocache-go"
	"unsafe"
)

func DETECT_Small(data []byte) []byte {
	msg := protocache.AsMessage(data)
	if !msg.IsValid() {
		return nil
	}
	inlined := msg.DetectInlined()
	if inlined == nil {
		return nil
	}
	compactEnd := len(inlined)
	if field := msg.GetField(_FIELD_Small_str); field.IsValid() {
		if obj := field.DetectObject(); obj != nil {
			part := func() []byte {
				if obj := field.DetectObject(); obj != nil {
					return protocache.DetectBytes(obj)
				}
				obj := field.GetObject()
				return protocache.DetectBytes(obj)
			}()
			if len(part) == 0 {
				return nil
			}
			tail := int(uintptr(unsafe.Pointer(unsafe.SliceData(obj)))-uintptr(unsafe.Pointer(unsafe.SliceData(data)))) + len(part)
			if tail > len(data) {
				return nil
			}
			return data[:tail]
		}
	}
	if compactEnd > len(data) {
		return nil
	}
	return data[:compactEnd]
}

type SmallEX struct {
	__    protocache.MessageEX
	fI32  int32
	fFlag bool
	fStr  string
	fJunk int64
}

func TO_SmallEX(data []byte) *SmallEX {
	out := &SmallEX{}
	out.__.Init(data)
	return out
}

func (m *SmallEX) HasBase() bool { return m.__.HasBase() }

func (m *SmallEX) Serialize() ([]byte, error) {
	words, err := m.serializeWords()
	if err != nil {
		return nil, err
	}
	return protocache.WordsToBytes(words), nil
}

func (m *SmallEX) serializeWords() ([]uint32, error) {
	if m == nil {
		return []uint32{0}, nil
	}
	parts := make([][]uint32, 5)
	if m.__.IsVisited(_FIELD_Small_i32, _FIELD_TOTAL_Small) {
		part, err := func() ([]uint32, error) { return protocache.EncodeInt32(m.fI32), nil }()
		if err != nil {
			return nil, err
		}
		parts[0] = part
	} else {
		field := m.__.RawField(_FIELD_Small_i32)
		if raw := func() []byte {
			if !field.IsValid() {
				return nil
			}
			return protocache.WordsToBytes(protocache.EncodeInt32(field.GetInt32()))
		}(); len(raw) != 0 {
			parts[0] = protocache.BytesToWords(raw)
		}
	}
	if m.__.IsVisited(_FIELD_Small_flag, _FIELD_TOTAL_Small) {
		part, err := func() ([]uint32, error) { return protocache.EncodeBool(m.fFlag), nil }()
		if err != nil {
			return nil, err
		}
		parts[1] = part
	} else {
		field := m.__.RawField(_FIELD_Small_flag)
		if raw := func() []byte {
			if !field.IsValid() {
				return nil
			}
			return protocache.WordsToBytes(protocache.EncodeBool(field.GetBool()))
		}(); len(raw) != 0 {
			parts[1] = protocache.BytesToWords(raw)
		}
	}
	if m.__.IsVisited(_FIELD_Small_str, _FIELD_TOTAL_Small) {
		part, err := protocache.EncodeString(m.fStr)
		if err != nil {
			return nil, err
		}
		parts[3] = part
	} else {
		field := m.__.RawField(_FIELD_Small_str)
		if raw := func() []byte {
			if obj := field.DetectObject(); obj != nil {
				return protocache.DetectBytes(obj)
			}
			obj := field.GetObject()
			return protocache.DetectBytes(obj)
		}(); len(raw) != 0 {
			parts[3] = protocache.BytesToWords(raw)
		}
	}
	if m.__.IsVisited(_FIELD_Small_junk, _FIELD_TOTAL_Small) {
		part, err := func() ([]uint32, error) { return protocache.EncodeInt64(m.fJunk), nil }()
		if err != nil {
			return nil, err
		}
		parts[4] = part
	} else {
		field := m.__.RawField(_FIELD_Small_junk)
		if raw := func() []byte {
			if !field.IsValid() {
				return nil
			}
			return protocache.WordsToBytes(protocache.EncodeInt64(field.GetInt64()))
		}(); len(raw) != 0 {
			parts[4] = protocache.BytesToWords(raw)
		}
	}
	return protocache.EncodeMessageParts(parts)
}

func (m *SmallEX) GetI32() int32 {
	if m.__.IsVisited(_FIELD_Small_i32, _FIELD_TOTAL_Small) {
		return m.fI32
	}
	field := m.__.RawField(_FIELD_Small_i32)
	m.fI32 = field.GetInt32()
	m.__.Visit(_FIELD_Small_i32, _FIELD_TOTAL_Small)
	return m.fI32
}

func (m *SmallEX) SetI32(v int32) {
	m.fI32 = v
	m.__.Visit(_FIELD_Small_i32, _FIELD_TOTAL_Small)
}

func (m *SmallEX) GetFlag() bool {
	if m.__.IsVisited(_FIELD_Small_flag, _FIELD_TOTAL_Small) {
		return m.fFlag
	}
	field := m.__.RawField(_FIELD_Small_flag)
	m.fFlag = field.GetBool()
	m.__.Visit(_FIELD_Small_flag, _FIELD_TOTAL_Small)
	return m.fFlag
}

func (m *SmallEX) SetFlag(v bool) {
	m.fFlag = v
	m.__.Visit(_FIELD_Small_flag, _FIELD_TOTAL_Small)
}

func (m *SmallEX) GetStr() string {
	if m.__.IsVisited(_FIELD_Small_str, _FIELD_TOTAL_Small) {
		return m.fStr
	}
	field := m.__.RawField(_FIELD_Small_str)
	m.fStr = field.GetString()
	m.__.Visit(_FIELD_Small_str, _FIELD_TOTAL_Small)
	return m.fStr
}

func (m *SmallEX) SetStr(v string) {
	m.fStr = v
	m.__.Visit(_FIELD_Small_str, _FIELD_TOTAL_Small)
}

func (m *SmallEX) GetJunk() int64 {
	if m.__.IsVisited(_FIELD_Small_junk, _FIELD_TOTAL_Small) {
		return m.fJunk
	}
	field := m.__.RawField(_FIELD_Small_junk)
	m.fJunk = field.GetInt64()
	m.__.Visit(_FIELD_Small_junk, _FIELD_TOTAL_Small)
	return m.fJunk
}

func (m *SmallEX) SetJunk(v int64) {
	m.fJunk = v
	m.__.Visit(_FIELD_Small_junk, _FIELD_TOTAL_Small)
}

func DETECT_Vec2D_Vec1D(data []byte) []byte {
	return protocache.DetectArray(data, nil)
}

type Vec2D_Vec1DEX []float32

func TO_Vec2D_Vec1DEX(data []byte) Vec2D_Vec1DEX {
	var out Vec2D_Vec1DEX
	arr := protocache.AsFloat32Array(data)
	out = append([]float32(nil), arr.Raw()...)
	return out
}

func (x Vec2D_Vec1DEX) Serialize() ([]byte, error) {
	words, err := serializeVec2D_Vec1DEX(x)
	if err != nil {
		return nil, err
	}
	return protocache.WordsToBytes(words), nil
}

func serializeVec2D_Vec1DEX(x Vec2D_Vec1DEX) ([]uint32, error) {
	if len(x) == 0 {
		return []uint32{1}, nil
	}
	return protocache.EncodeFloat32Array([]float32(x))
}

func DETECT_Vec2D(data []byte) []byte {
	return protocache.DetectArray(data, DETECT_Vec2D_Vec1D)
}

type Vec2DEX []Vec2D_Vec1DEX

func TO_Vec2DEX(data []byte) Vec2DEX {
	var out Vec2DEX
	arr := protocache.AsArray(data)
	out = make([]Vec2D_Vec1DEX, int(arr.Size()))
	for i := uint32(0); i < arr.Size(); i++ {
		elem := arr.Get(i)
		out[i] = TO_Vec2D_Vec1DEX(elem.GetObject())
	}
	return out
}

func (x Vec2DEX) Serialize() ([]byte, error) {
	words, err := serializeVec2DEX(x)
	if err != nil {
		return nil, err
	}
	return protocache.WordsToBytes(words), nil
}

func serializeVec2DEX(x Vec2DEX) ([]uint32, error) {
	if len(x) == 0 {
		return []uint32{1}, nil
	}
	return protocache.EncodeObjectArray(len(x), func(i int) ([]uint32, error) {
		return func() ([]uint32, error) {
			data, err := x[i].Serialize()
			if err != nil {
				return nil, err
			}
			return protocache.BytesToWords(data), nil
		}()
	})
}

func DETECT_ArrMap_Array(data []byte) []byte {
	return protocache.DetectArray(data, nil)
}

type ArrMap_ArrayEX []float32

func TO_ArrMap_ArrayEX(data []byte) ArrMap_ArrayEX {
	var out ArrMap_ArrayEX
	arr := protocache.AsFloat32Array(data)
	out = append([]float32(nil), arr.Raw()...)
	return out
}

func (x ArrMap_ArrayEX) Serialize() ([]byte, error) {
	words, err := serializeArrMap_ArrayEX(x)
	if err != nil {
		return nil, err
	}
	return protocache.WordsToBytes(words), nil
}

func serializeArrMap_ArrayEX(x ArrMap_ArrayEX) ([]uint32, error) {
	if len(x) == 0 {
		return []uint32{1}, nil
	}
	return protocache.EncodeFloat32Array([]float32(x))
}

func DETECT_ArrMap(data []byte) []byte {
	return protocache.DetectMap(data, protocache.DetectBytes, DETECT_ArrMap_Array)
}

type ArrMapEX map[string]ArrMap_ArrayEX

func TO_ArrMapEX(data []byte) ArrMapEX {
	var out ArrMapEX
	pack := protocache.AsMap(data)
	out = make(map[string]ArrMap_ArrayEX, int(pack.Size()))
	for i := uint32(0); i < pack.Size(); i++ {
		keyField := pack.Key(i)
		key := keyField.GetString()
		valField := pack.Value(i)
		out[key] = TO_ArrMap_ArrayEX(valField.GetObject())
	}
	return out
}

func (x ArrMapEX) Serialize() ([]byte, error) {
	words, err := serializeArrMapEX(x)
	if err != nil {
		return nil, err
	}
	return protocache.WordsToBytes(words), nil
}

func serializeArrMapEX(x ArrMapEX) ([]uint32, error) {
	if len(x) == 0 {
		return []uint32{5 << 28}, nil
	}
	keys := make([][]uint32, 0, len(x))
	vals := make([][]uint32, 0, len(x))
	for k, v := range x {
		keyPart, err := protocache.EncodeString(k)
		if err != nil {
			return nil, err
		}
		valPart, err := func() ([]uint32, error) {
			data, err := v.Serialize()
			if err != nil {
				return nil, err
			}
			return protocache.BytesToWords(data), nil
		}()
		if err != nil {
			return nil, err
		}
		keys = append(keys, keyPart)
		vals = append(vals, valPart)
	}
	return protocache.EncodeMapParts(keys, vals, true)
}

func DETECT_Main(data []byte) []byte {
	msg := protocache.AsMessage(data)
	if !msg.IsValid() {
		return nil
	}
	inlined := msg.DetectInlined()
	if inlined == nil {
		return nil
	}
	compactEnd := len(inlined)
	if field := msg.GetField(_FIELD_Main_modev); field.IsValid() {
		if obj := field.DetectObject(); obj != nil {
			part := func() []byte {
				if obj := field.DetectObject(); obj != nil {
					return protocache.DetectArray(obj, nil)
				}
				obj := field.GetObject()
				return protocache.DetectArray(obj, nil)
			}()
			if len(part) == 0 {
				return nil
			}
			tail := int(uintptr(unsafe.Pointer(unsafe.SliceData(obj)))-uintptr(unsafe.Pointer(unsafe.SliceData(data)))) + len(part)
			if tail > len(data) {
				return nil
			}
			return data[:tail]
		}
	}
	if field := msg.GetField(_FIELD_Main_arrays); field.IsValid() {
		if obj := field.DetectObject(); obj != nil {
			part := func() []byte {
				if obj := field.DetectObject(); obj != nil {
					return DETECT_ArrMap(obj)
				}
				obj := field.GetObject()
				return DETECT_ArrMap(obj)
			}()
			if len(part) == 0 {
				return nil
			}
			tail := int(uintptr(unsafe.Pointer(unsafe.SliceData(obj)))-uintptr(unsafe.Pointer(unsafe.SliceData(data)))) + len(part)
			if tail > len(data) {
				return nil
			}
			return data[:tail]
		}
	}
	if field := msg.GetField(_FIELD_Main_vector); field.IsValid() {
		if obj := field.DetectObject(); obj != nil {
			part := func() []byte {
				if obj := field.DetectObject(); obj != nil {
					return protocache.DetectArray(obj, DETECT_ArrMap)
				}
				obj := field.GetObject()
				return protocache.DetectArray(obj, DETECT_ArrMap)
			}()
			if len(part) == 0 {
				return nil
			}
			tail := int(uintptr(unsafe.Pointer(unsafe.SliceData(obj)))-uintptr(unsafe.Pointer(unsafe.SliceData(data)))) + len(part)
			if tail > len(data) {
				return nil
			}
			return data[:tail]
		}
	}
	if field := msg.GetField(_FIELD_Main_matrix); field.IsValid() {
		if obj := field.DetectObject(); obj != nil {
			part := func() []byte {
				if obj := field.DetectObject(); obj != nil {
					return DETECT_Vec2D(obj)
				}
				obj := field.GetObject()
				return DETECT_Vec2D(obj)
			}()
			if len(part) == 0 {
				return nil
			}
			tail := int(uintptr(unsafe.Pointer(unsafe.SliceData(obj)))-uintptr(unsafe.Pointer(unsafe.SliceData(data)))) + len(part)
			if tail > len(data) {
				return nil
			}
			return data[:tail]
		}
	}
	if field := msg.GetField(_FIELD_Main_objects); field.IsValid() {
		if obj := field.DetectObject(); obj != nil {
			part := func() []byte {
				if obj := field.DetectObject(); obj != nil {
					return protocache.DetectMap(obj, nil, DETECT_Small)
				}
				obj := field.GetObject()
				return protocache.DetectMap(obj, nil, DETECT_Small)
			}()
			if len(part) == 0 {
				return nil
			}
			tail := int(uintptr(unsafe.Pointer(unsafe.SliceData(obj)))-uintptr(unsafe.Pointer(unsafe.SliceData(data)))) + len(part)
			if tail > len(data) {
				return nil
			}
			return data[:tail]
		}
	}
	if field := msg.GetField(_FIELD_Main_index); field.IsValid() {
		if obj := field.DetectObject(); obj != nil {
			part := func() []byte {
				if obj := field.DetectObject(); obj != nil {
					return protocache.DetectMap(obj, protocache.DetectBytes, nil)
				}
				obj := field.GetObject()
				return protocache.DetectMap(obj, protocache.DetectBytes, nil)
			}()
			if len(part) == 0 {
				return nil
			}
			tail := int(uintptr(unsafe.Pointer(unsafe.SliceData(obj)))-uintptr(unsafe.Pointer(unsafe.SliceData(data)))) + len(part)
			if tail > len(data) {
				return nil
			}
			return data[:tail]
		}
	}
	if field := msg.GetField(_FIELD_Main_objectv); field.IsValid() {
		if obj := field.DetectObject(); obj != nil {
			part := func() []byte {
				if obj := field.DetectObject(); obj != nil {
					return protocache.DetectArray(obj, DETECT_Small)
				}
				obj := field.GetObject()
				return protocache.DetectArray(obj, DETECT_Small)
			}()
			if len(part) == 0 {
				return nil
			}
			tail := int(uintptr(unsafe.Pointer(unsafe.SliceData(obj)))-uintptr(unsafe.Pointer(unsafe.SliceData(data)))) + len(part)
			if tail > len(data) {
				return nil
			}
			return data[:tail]
		}
	}
	if field := msg.GetField(_FIELD_Main_flags); field.IsValid() {
		if obj := field.DetectObject(); obj != nil {
			part := func() []byte {
				if obj := field.DetectObject(); obj != nil {
					return protocache.DetectBytes(obj)
				}
				obj := field.GetObject()
				return protocache.DetectBytes(obj)
			}()
			if len(part) == 0 {
				return nil
			}
			tail := int(uintptr(unsafe.Pointer(unsafe.SliceData(obj)))-uintptr(unsafe.Pointer(unsafe.SliceData(data)))) + len(part)
			if tail > len(data) {
				return nil
			}
			return data[:tail]
		}
	}
	if field := msg.GetField(_FIELD_Main_f64v); field.IsValid() {
		if obj := field.DetectObject(); obj != nil {
			part := func() []byte {
				if obj := field.DetectObject(); obj != nil {
					return protocache.DetectArray(obj, nil)
				}
				obj := field.GetObject()
				return protocache.DetectArray(obj, nil)
			}()
			if len(part) == 0 {
				return nil
			}
			tail := int(uintptr(unsafe.Pointer(unsafe.SliceData(obj)))-uintptr(unsafe.Pointer(unsafe.SliceData(data)))) + len(part)
			if tail > len(data) {
				return nil
			}
			return data[:tail]
		}
	}
	if field := msg.GetField(_FIELD_Main_f32v); field.IsValid() {
		if obj := field.DetectObject(); obj != nil {
			part := func() []byte {
				if obj := field.DetectObject(); obj != nil {
					return protocache.DetectArray(obj, nil)
				}
				obj := field.GetObject()
				return protocache.DetectArray(obj, nil)
			}()
			if len(part) == 0 {
				return nil
			}
			tail := int(uintptr(unsafe.Pointer(unsafe.SliceData(obj)))-uintptr(unsafe.Pointer(unsafe.SliceData(data)))) + len(part)
			if tail > len(data) {
				return nil
			}
			return data[:tail]
		}
	}
	if field := msg.GetField(_FIELD_Main_datav); field.IsValid() {
		if obj := field.DetectObject(); obj != nil {
			part := func() []byte {
				if obj := field.DetectObject(); obj != nil {
					return protocache.DetectArray(obj, protocache.DetectBytes)
				}
				obj := field.GetObject()
				return protocache.DetectArray(obj, protocache.DetectBytes)
			}()
			if len(part) == 0 {
				return nil
			}
			tail := int(uintptr(unsafe.Pointer(unsafe.SliceData(obj)))-uintptr(unsafe.Pointer(unsafe.SliceData(data)))) + len(part)
			if tail > len(data) {
				return nil
			}
			return data[:tail]
		}
	}
	if field := msg.GetField(_FIELD_Main_strv); field.IsValid() {
		if obj := field.DetectObject(); obj != nil {
			part := func() []byte {
				if obj := field.DetectObject(); obj != nil {
					return protocache.DetectArray(obj, protocache.DetectBytes)
				}
				obj := field.GetObject()
				return protocache.DetectArray(obj, protocache.DetectBytes)
			}()
			if len(part) == 0 {
				return nil
			}
			tail := int(uintptr(unsafe.Pointer(unsafe.SliceData(obj)))-uintptr(unsafe.Pointer(unsafe.SliceData(data)))) + len(part)
			if tail > len(data) {
				return nil
			}
			return data[:tail]
		}
	}
	if field := msg.GetField(_FIELD_Main_u64v); field.IsValid() {
		if obj := field.DetectObject(); obj != nil {
			part := func() []byte {
				if obj := field.DetectObject(); obj != nil {
					return protocache.DetectArray(obj, nil)
				}
				obj := field.GetObject()
				return protocache.DetectArray(obj, nil)
			}()
			if len(part) == 0 {
				return nil
			}
			tail := int(uintptr(unsafe.Pointer(unsafe.SliceData(obj)))-uintptr(unsafe.Pointer(unsafe.SliceData(data)))) + len(part)
			if tail > len(data) {
				return nil
			}
			return data[:tail]
		}
	}
	if field := msg.GetField(_FIELD_Main_i32v); field.IsValid() {
		if obj := field.DetectObject(); obj != nil {
			part := func() []byte {
				if obj := field.DetectObject(); obj != nil {
					return protocache.DetectArray(obj, nil)
				}
				obj := field.GetObject()
				return protocache.DetectArray(obj, nil)
			}()
			if len(part) == 0 {
				return nil
			}
			tail := int(uintptr(unsafe.Pointer(unsafe.SliceData(obj)))-uintptr(unsafe.Pointer(unsafe.SliceData(data)))) + len(part)
			if tail > len(data) {
				return nil
			}
			return data[:tail]
		}
	}
	if field := msg.GetField(_FIELD_Main_object); field.IsValid() {
		if obj := field.DetectObject(); obj != nil {
			part := func() []byte {
				if obj := field.DetectObject(); obj != nil {
					return DETECT_Small(obj)
				}
				obj := field.GetObject()
				return DETECT_Small(obj)
			}()
			if len(part) == 0 {
				return nil
			}
			tail := int(uintptr(unsafe.Pointer(unsafe.SliceData(obj)))-uintptr(unsafe.Pointer(unsafe.SliceData(data)))) + len(part)
			if tail > len(data) {
				return nil
			}
			return data[:tail]
		}
	}
	if field := msg.GetField(_FIELD_Main_data); field.IsValid() {
		if obj := field.DetectObject(); obj != nil {
			part := func() []byte {
				if obj := field.DetectObject(); obj != nil {
					return protocache.DetectBytes(obj)
				}
				obj := field.GetObject()
				return protocache.DetectBytes(obj)
			}()
			if len(part) == 0 {
				return nil
			}
			tail := int(uintptr(unsafe.Pointer(unsafe.SliceData(obj)))-uintptr(unsafe.Pointer(unsafe.SliceData(data)))) + len(part)
			if tail > len(data) {
				return nil
			}
			return data[:tail]
		}
	}
	if field := msg.GetField(_FIELD_Main_str); field.IsValid() {
		if obj := field.DetectObject(); obj != nil {
			part := func() []byte {
				if obj := field.DetectObject(); obj != nil {
					return protocache.DetectBytes(obj)
				}
				obj := field.GetObject()
				return protocache.DetectBytes(obj)
			}()
			if len(part) == 0 {
				return nil
			}
			tail := int(uintptr(unsafe.Pointer(unsafe.SliceData(obj)))-uintptr(unsafe.Pointer(unsafe.SliceData(data)))) + len(part)
			if tail > len(data) {
				return nil
			}
			return data[:tail]
		}
	}
	if compactEnd > len(data) {
		return nil
	}
	return data[:compactEnd]
}

type MainEX struct {
	__       protocache.MessageEX
	fI32     int32
	fU32     uint32
	fI64     int64
	fU64     uint64
	fFlag    bool
	fMode    Mode
	fStr     string
	fData    []byte
	fF32     float32
	fF64     float64
	fObject  *SmallEX
	fI32V    []int32
	fU64V    []uint64
	fStrv    []string
	fDatav   [][]byte
	fF32V    []float32
	fF64V    []float64
	fFlags   []bool
	fObjectv []*SmallEX
	fTU32    uint32
	fTI32    int32
	fTS32    int32
	fTU64    uint64
	fTI64    int64
	fTS64    int64
	fIndex   map[string]int32
	fObjects map[int32]*SmallEX
	fMatrix  Vec2DEX
	fVector  []ArrMapEX
	fArrays  ArrMapEX
	fModev   []Mode
}

func TO_MainEX(data []byte) *MainEX {
	out := &MainEX{}
	out.__.Init(data)
	return out
}

func (m *MainEX) HasBase() bool { return m.__.HasBase() }

func (m *MainEX) Serialize() ([]byte, error) {
	words, err := m.serializeWords()
	if err != nil {
		return nil, err
	}
	return protocache.WordsToBytes(words), nil
}

func (m *MainEX) serializeWords() ([]uint32, error) {
	if m == nil {
		return []uint32{0}, nil
	}
	parts := make([][]uint32, 32)
	if m.__.IsVisited(_FIELD_Main_i32, _FIELD_TOTAL_Main) {
		part, err := func() ([]uint32, error) { return protocache.EncodeInt32(m.fI32), nil }()
		if err != nil {
			return nil, err
		}
		parts[0] = part
	} else {
		field := m.__.RawField(_FIELD_Main_i32)
		if part := func() []uint32 {
			if !field.IsValid() {
				return nil
			}
			return field.RawWords()
		}(); len(part) != 0 {
			parts[0] = part
		}
	}
	if m.__.IsVisited(_FIELD_Main_u32, _FIELD_TOTAL_Main) {
		part, err := func() ([]uint32, error) { return protocache.EncodeUint32(m.fU32), nil }()
		if err != nil {
			return nil, err
		}
		parts[1] = part
	} else {
		field := m.__.RawField(_FIELD_Main_u32)
		if part := func() []uint32 {
			if !field.IsValid() {
				return nil
			}
			return field.RawWords()
		}(); len(part) != 0 {
			parts[1] = part
		}
	}
	if m.__.IsVisited(_FIELD_Main_i64, _FIELD_TOTAL_Main) {
		part, err := func() ([]uint32, error) { return protocache.EncodeInt64(m.fI64), nil }()
		if err != nil {
			return nil, err
		}
		parts[2] = part
	} else {
		field := m.__.RawField(_FIELD_Main_i64)
		if part := func() []uint32 {
			if !field.IsValid() {
				return nil
			}
			return field.RawWords()
		}(); len(part) != 0 {
			parts[2] = part
		}
	}
	if m.__.IsVisited(_FIELD_Main_u64, _FIELD_TOTAL_Main) {
		part, err := func() ([]uint32, error) { return protocache.EncodeUint64(m.fU64), nil }()
		if err != nil {
			return nil, err
		}
		parts[3] = part
	} else {
		field := m.__.RawField(_FIELD_Main_u64)
		if part := func() []uint32 {
			if !field.IsValid() {
				return nil
			}
			return field.RawWords()
		}(); len(part) != 0 {
			parts[3] = part
		}
	}
	if m.__.IsVisited(_FIELD_Main_flag, _FIELD_TOTAL_Main) {
		part, err := func() ([]uint32, error) { return protocache.EncodeBool(m.fFlag), nil }()
		if err != nil {
			return nil, err
		}
		parts[4] = part
	} else {
		field := m.__.RawField(_FIELD_Main_flag)
		if part := func() []uint32 {
			if !field.IsValid() {
				return nil
			}
			return field.RawWords()
		}(); len(part) != 0 {
			parts[4] = part
		}
	}
	if m.__.IsVisited(_FIELD_Main_mode, _FIELD_TOTAL_Main) {
		part, err := func() ([]uint32, error) { return protocache.EncodeInt32(int32(m.fMode)), nil }()
		if err != nil {
			return nil, err
		}
		parts[5] = part
	} else {
		field := m.__.RawField(_FIELD_Main_mode)
		if part := func() []uint32 {
			if !field.IsValid() {
				return nil
			}
			return field.RawWords()
		}(); len(part) != 0 {
			parts[5] = part
		}
	}
	if m.__.IsVisited(_FIELD_Main_str, _FIELD_TOTAL_Main) {
		part, err := protocache.EncodeString(m.fStr)
		if err != nil {
			return nil, err
		}
		parts[6] = part
	} else {
		field := m.__.RawField(_FIELD_Main_str)
		if raw := func() []byte {
			if obj := field.DetectObject(); obj != nil {
				return protocache.DetectBytes(obj)
			}
			obj := field.GetObject()
			return protocache.DetectBytes(obj)
		}(); len(raw) != 0 {
			parts[6] = protocache.BytesToWords(raw)
		}
	}
	if m.__.IsVisited(_FIELD_Main_data, _FIELD_TOTAL_Main) {
		part, err := protocache.EncodeBytes(m.fData)
		if err != nil {
			return nil, err
		}
		parts[7] = part
	} else {
		field := m.__.RawField(_FIELD_Main_data)
		if raw := func() []byte {
			if obj := field.DetectObject(); obj != nil {
				return protocache.DetectBytes(obj)
			}
			obj := field.GetObject()
			return protocache.DetectBytes(obj)
		}(); len(raw) != 0 {
			parts[7] = protocache.BytesToWords(raw)
		}
	}
	if m.__.IsVisited(_FIELD_Main_f32, _FIELD_TOTAL_Main) {
		part, err := func() ([]uint32, error) { return protocache.EncodeFloat32(m.fF32), nil }()
		if err != nil {
			return nil, err
		}
		parts[8] = part
	} else {
		field := m.__.RawField(_FIELD_Main_f32)
		if part := func() []uint32 {
			if !field.IsValid() {
				return nil
			}
			return field.RawWords()
		}(); len(part) != 0 {
			parts[8] = part
		}
	}
	if m.__.IsVisited(_FIELD_Main_f64, _FIELD_TOTAL_Main) {
		part, err := func() ([]uint32, error) { return protocache.EncodeFloat64(m.fF64), nil }()
		if err != nil {
			return nil, err
		}
		parts[9] = part
	} else {
		field := m.__.RawField(_FIELD_Main_f64)
		if part := func() []uint32 {
			if !field.IsValid() {
				return nil
			}
			return field.RawWords()
		}(); len(part) != 0 {
			parts[9] = part
		}
	}
	if m.__.IsVisited(_FIELD_Main_object, _FIELD_TOTAL_Main) {
		if m.fObject != nil {
			part, err := m.fObject.serializeWords()
			if err != nil {
				return nil, err
			}
			if len(part) > 1 {
				parts[10] = part
			}
		}
	} else {
		field := m.__.RawField(_FIELD_Main_object)
		if raw := func() []byte {
			if obj := field.DetectObject(); obj != nil {
				return DETECT_Small(obj)
			}
			obj := field.GetObject()
			return DETECT_Small(obj)
		}(); len(raw) != 0 {
			parts[10] = protocache.BytesToWords(raw)
		}
	}
	if m.__.IsVisited(_FIELD_Main_i32v, _FIELD_TOTAL_Main) {
		if len(m.fI32V) != 0 {
			part, err := protocache.EncodeInt32Array(m.fI32V)
			if err != nil {
				return nil, err
			}
			parts[11] = part
		}
	} else {
		field := m.__.RawField(_FIELD_Main_i32v)
		if raw := func() []byte {
			if obj := field.DetectObject(); obj != nil {
				return protocache.DetectArray(obj, nil)
			}
			obj := field.GetObject()
			return protocache.DetectArray(obj, nil)
		}(); len(raw) != 0 {
			parts[11] = protocache.BytesToWords(raw)
		}
	}
	if m.__.IsVisited(_FIELD_Main_u64v, _FIELD_TOTAL_Main) {
		if len(m.fU64V) != 0 {
			part, err := protocache.EncodeUint64Array(m.fU64V)
			if err != nil {
				return nil, err
			}
			parts[12] = part
		}
	} else {
		field := m.__.RawField(_FIELD_Main_u64v)
		if raw := func() []byte {
			if obj := field.DetectObject(); obj != nil {
				return protocache.DetectArray(obj, nil)
			}
			obj := field.GetObject()
			return protocache.DetectArray(obj, nil)
		}(); len(raw) != 0 {
			parts[12] = protocache.BytesToWords(raw)
		}
	}
	if m.__.IsVisited(_FIELD_Main_strv, _FIELD_TOTAL_Main) {
		if len(m.fStrv) != 0 {
			part, err := protocache.EncodeStringArray(m.fStrv)
			if err != nil {
				return nil, err
			}
			parts[13] = part
		}
	} else {
		field := m.__.RawField(_FIELD_Main_strv)
		if raw := func() []byte {
			if obj := field.DetectObject(); obj != nil {
				return protocache.DetectArray(obj, protocache.DetectBytes)
			}
			obj := field.GetObject()
			return protocache.DetectArray(obj, protocache.DetectBytes)
		}(); len(raw) != 0 {
			parts[13] = protocache.BytesToWords(raw)
		}
	}
	if m.__.IsVisited(_FIELD_Main_datav, _FIELD_TOTAL_Main) {
		if len(m.fDatav) != 0 {
			part, err := protocache.EncodeBytesArray(m.fDatav)
			if err != nil {
				return nil, err
			}
			parts[14] = part
		}
	} else {
		field := m.__.RawField(_FIELD_Main_datav)
		if raw := func() []byte {
			if obj := field.DetectObject(); obj != nil {
				return protocache.DetectArray(obj, protocache.DetectBytes)
			}
			obj := field.GetObject()
			return protocache.DetectArray(obj, protocache.DetectBytes)
		}(); len(raw) != 0 {
			parts[14] = protocache.BytesToWords(raw)
		}
	}
	if m.__.IsVisited(_FIELD_Main_f32v, _FIELD_TOTAL_Main) {
		if len(m.fF32V) != 0 {
			part, err := protocache.EncodeFloat32Array(m.fF32V)
			if err != nil {
				return nil, err
			}
			parts[15] = part
		}
	} else {
		field := m.__.RawField(_FIELD_Main_f32v)
		if raw := func() []byte {
			if obj := field.DetectObject(); obj != nil {
				return protocache.DetectArray(obj, nil)
			}
			obj := field.GetObject()
			return protocache.DetectArray(obj, nil)
		}(); len(raw) != 0 {
			parts[15] = protocache.BytesToWords(raw)
		}
	}
	if m.__.IsVisited(_FIELD_Main_f64v, _FIELD_TOTAL_Main) {
		if len(m.fF64V) != 0 {
			part, err := protocache.EncodeFloat64Array(m.fF64V)
			if err != nil {
				return nil, err
			}
			parts[16] = part
		}
	} else {
		field := m.__.RawField(_FIELD_Main_f64v)
		if raw := func() []byte {
			if obj := field.DetectObject(); obj != nil {
				return protocache.DetectArray(obj, nil)
			}
			obj := field.GetObject()
			return protocache.DetectArray(obj, nil)
		}(); len(raw) != 0 {
			parts[16] = protocache.BytesToWords(raw)
		}
	}
	if m.__.IsVisited(_FIELD_Main_flags, _FIELD_TOTAL_Main) {
		if len(m.fFlags) != 0 {
			part, err := protocache.EncodeBoolArray(m.fFlags)
			if err != nil {
				return nil, err
			}
			parts[17] = part
		}
	} else {
		field := m.__.RawField(_FIELD_Main_flags)
		if raw := func() []byte {
			if obj := field.DetectObject(); obj != nil {
				return protocache.DetectBytes(obj)
			}
			obj := field.GetObject()
			return protocache.DetectBytes(obj)
		}(); len(raw) != 0 {
			parts[17] = protocache.BytesToWords(raw)
		}
	}
	if m.__.IsVisited(_FIELD_Main_objectv, _FIELD_TOTAL_Main) {
		if len(m.fObjectv) != 0 {
			part, err := protocache.EncodeObjectArray(len(m.fObjectv), func(i int) ([]uint32, error) {
				if m.fObjectv[i] == nil {
					return []uint32{0}, nil
				}
				return m.fObjectv[i].serializeWords()
			})
			if err != nil {
				return nil, err
			}
			parts[18] = part
		}
	} else {
		field := m.__.RawField(_FIELD_Main_objectv)
		if raw := func() []byte {
			if obj := field.DetectObject(); obj != nil {
				return protocache.DetectArray(obj, DETECT_Small)
			}
			obj := field.GetObject()
			return protocache.DetectArray(obj, DETECT_Small)
		}(); len(raw) != 0 {
			parts[18] = protocache.BytesToWords(raw)
		}
	}
	if m.__.IsVisited(_FIELD_Main_t_u32, _FIELD_TOTAL_Main) {
		part, err := func() ([]uint32, error) { return protocache.EncodeUint32(m.fTU32), nil }()
		if err != nil {
			return nil, err
		}
		parts[19] = part
	} else {
		field := m.__.RawField(_FIELD_Main_t_u32)
		if part := func() []uint32 {
			if !field.IsValid() {
				return nil
			}
			return field.RawWords()
		}(); len(part) != 0 {
			parts[19] = part
		}
	}
	if m.__.IsVisited(_FIELD_Main_t_i32, _FIELD_TOTAL_Main) {
		part, err := func() ([]uint32, error) { return protocache.EncodeInt32(m.fTI32), nil }()
		if err != nil {
			return nil, err
		}
		parts[20] = part
	} else {
		field := m.__.RawField(_FIELD_Main_t_i32)
		if part := func() []uint32 {
			if !field.IsValid() {
				return nil
			}
			return field.RawWords()
		}(); len(part) != 0 {
			parts[20] = part
		}
	}
	if m.__.IsVisited(_FIELD_Main_t_s32, _FIELD_TOTAL_Main) {
		part, err := func() ([]uint32, error) { return protocache.EncodeInt32(m.fTS32), nil }()
		if err != nil {
			return nil, err
		}
		parts[21] = part
	} else {
		field := m.__.RawField(_FIELD_Main_t_s32)
		if part := func() []uint32 {
			if !field.IsValid() {
				return nil
			}
			return field.RawWords()
		}(); len(part) != 0 {
			parts[21] = part
		}
	}
	if m.__.IsVisited(_FIELD_Main_t_u64, _FIELD_TOTAL_Main) {
		part, err := func() ([]uint32, error) { return protocache.EncodeUint64(m.fTU64), nil }()
		if err != nil {
			return nil, err
		}
		parts[22] = part
	} else {
		field := m.__.RawField(_FIELD_Main_t_u64)
		if part := func() []uint32 {
			if !field.IsValid() {
				return nil
			}
			return field.RawWords()
		}(); len(part) != 0 {
			parts[22] = part
		}
	}
	if m.__.IsVisited(_FIELD_Main_t_i64, _FIELD_TOTAL_Main) {
		part, err := func() ([]uint32, error) { return protocache.EncodeInt64(m.fTI64), nil }()
		if err != nil {
			return nil, err
		}
		parts[23] = part
	} else {
		field := m.__.RawField(_FIELD_Main_t_i64)
		if part := func() []uint32 {
			if !field.IsValid() {
				return nil
			}
			return field.RawWords()
		}(); len(part) != 0 {
			parts[23] = part
		}
	}
	if m.__.IsVisited(_FIELD_Main_t_s64, _FIELD_TOTAL_Main) {
		part, err := func() ([]uint32, error) { return protocache.EncodeInt64(m.fTS64), nil }()
		if err != nil {
			return nil, err
		}
		parts[24] = part
	} else {
		field := m.__.RawField(_FIELD_Main_t_s64)
		if part := func() []uint32 {
			if !field.IsValid() {
				return nil
			}
			return field.RawWords()
		}(); len(part) != 0 {
			parts[24] = part
		}
	}
	if m.__.IsVisited(_FIELD_Main_index, _FIELD_TOTAL_Main) {
		if len(m.fIndex) != 0 {
			keys := make([][]uint32, 0, len(m.fIndex))
			vals := make([][]uint32, 0, len(m.fIndex))
			for k, v := range m.fIndex {
				keyPart, err := protocache.EncodeString(k)
				if err != nil {
					return nil, err
				}
				valPart, err := func() ([]uint32, error) { return protocache.EncodeInt32(v), nil }()
				if err != nil {
					return nil, err
				}
				keys = append(keys, keyPart)
				vals = append(vals, valPart)
			}
			part, err := protocache.EncodeMapParts(keys, vals, true)
			if err != nil {
				return nil, err
			}
			parts[25] = part
		}
	} else {
		field := m.__.RawField(_FIELD_Main_index)
		if raw := func() []byte {
			if obj := field.DetectObject(); obj != nil {
				return protocache.DetectMap(obj, protocache.DetectBytes, nil)
			}
			obj := field.GetObject()
			return protocache.DetectMap(obj, protocache.DetectBytes, nil)
		}(); len(raw) != 0 {
			parts[25] = protocache.BytesToWords(raw)
		}
	}
	if m.__.IsVisited(_FIELD_Main_objects, _FIELD_TOTAL_Main) {
		if len(m.fObjects) != 0 {
			keys := make([][]uint32, 0, len(m.fObjects))
			vals := make([][]uint32, 0, len(m.fObjects))
			for k, v := range m.fObjects {
				keyPart, err := func() ([]uint32, error) { return protocache.EncodeInt32(k), nil }()
				if err != nil {
					return nil, err
				}
				valPart, err := v.serializeWords()
				if err != nil {
					return nil, err
				}
				if len(valPart) <= 1 {
					valPart = nil
				}
				keys = append(keys, keyPart)
				vals = append(vals, valPart)
			}
			part, err := protocache.EncodeMapParts(keys, vals, false)
			if err != nil {
				return nil, err
			}
			parts[26] = part
		}
	} else {
		field := m.__.RawField(_FIELD_Main_objects)
		if raw := func() []byte {
			if obj := field.DetectObject(); obj != nil {
				return protocache.DetectMap(obj, nil, DETECT_Small)
			}
			obj := field.GetObject()
			return protocache.DetectMap(obj, nil, DETECT_Small)
		}(); len(raw) != 0 {
			parts[26] = protocache.BytesToWords(raw)
		}
	}
	if m.__.IsVisited(_FIELD_Main_matrix, _FIELD_TOTAL_Main) {
		if m.fMatrix != nil {
			part, err := func() ([]uint32, error) {
				data, err := m.fMatrix.Serialize()
				if err != nil {
					return nil, err
				}
				return protocache.BytesToWords(data), nil
			}()
			if err != nil {
				return nil, err
			}
			if len(part) > 1 {
				parts[27] = part
			}
		}
	} else {
		field := m.__.RawField(_FIELD_Main_matrix)
		if raw := func() []byte {
			if obj := field.DetectObject(); obj != nil {
				return DETECT_Vec2D(obj)
			}
			obj := field.GetObject()
			return DETECT_Vec2D(obj)
		}(); len(raw) != 0 {
			parts[27] = protocache.BytesToWords(raw)
		}
	}
	if m.__.IsVisited(_FIELD_Main_vector, _FIELD_TOTAL_Main) {
		if len(m.fVector) != 0 {
			part, err := protocache.EncodeObjectArray(len(m.fVector), func(i int) ([]uint32, error) {
				return func() ([]uint32, error) {
					data, err := m.fVector[i].Serialize()
					if err != nil {
						return nil, err
					}
					return protocache.BytesToWords(data), nil
				}()
			})
			if err != nil {
				return nil, err
			}
			parts[28] = part
		}
	} else {
		field := m.__.RawField(_FIELD_Main_vector)
		if raw := func() []byte {
			if obj := field.DetectObject(); obj != nil {
				return protocache.DetectArray(obj, DETECT_ArrMap)
			}
			obj := field.GetObject()
			return protocache.DetectArray(obj, DETECT_ArrMap)
		}(); len(raw) != 0 {
			parts[28] = protocache.BytesToWords(raw)
		}
	}
	if m.__.IsVisited(_FIELD_Main_arrays, _FIELD_TOTAL_Main) {
		if m.fArrays != nil {
			part, err := func() ([]uint32, error) {
				data, err := m.fArrays.Serialize()
				if err != nil {
					return nil, err
				}
				return protocache.BytesToWords(data), nil
			}()
			if err != nil {
				return nil, err
			}
			if len(part) > 1 {
				parts[29] = part
			}
		}
	} else {
		field := m.__.RawField(_FIELD_Main_arrays)
		if raw := func() []byte {
			if obj := field.DetectObject(); obj != nil {
				return DETECT_ArrMap(obj)
			}
			obj := field.GetObject()
			return DETECT_ArrMap(obj)
		}(); len(raw) != 0 {
			parts[29] = protocache.BytesToWords(raw)
		}
	}
	if m.__.IsVisited(_FIELD_Main_modev, _FIELD_TOTAL_Main) {
		if len(m.fModev) != 0 {
			part, err := protocache.EncodeEnumArray(m.fModev)
			if err != nil {
				return nil, err
			}
			parts[31] = part
		}
	} else {
		field := m.__.RawField(_FIELD_Main_modev)
		if raw := func() []byte {
			if obj := field.DetectObject(); obj != nil {
				return protocache.DetectArray(obj, nil)
			}
			obj := field.GetObject()
			return protocache.DetectArray(obj, nil)
		}(); len(raw) != 0 {
			parts[31] = protocache.BytesToWords(raw)
		}
	}
	return protocache.EncodeMessageParts(parts)
}

func (m *MainEX) GetI32() int32 {
	if m.__.IsVisited(_FIELD_Main_i32, _FIELD_TOTAL_Main) {
		return m.fI32
	}
	field := m.__.RawField(_FIELD_Main_i32)
	m.fI32 = field.GetInt32()
	m.__.Visit(_FIELD_Main_i32, _FIELD_TOTAL_Main)
	return m.fI32
}

func (m *MainEX) SetI32(v int32) {
	m.fI32 = v
	m.__.Visit(_FIELD_Main_i32, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetU32() uint32 {
	if m.__.IsVisited(_FIELD_Main_u32, _FIELD_TOTAL_Main) {
		return m.fU32
	}
	field := m.__.RawField(_FIELD_Main_u32)
	m.fU32 = field.GetUint32()
	m.__.Visit(_FIELD_Main_u32, _FIELD_TOTAL_Main)
	return m.fU32
}

func (m *MainEX) SetU32(v uint32) {
	m.fU32 = v
	m.__.Visit(_FIELD_Main_u32, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetI64() int64 {
	if m.__.IsVisited(_FIELD_Main_i64, _FIELD_TOTAL_Main) {
		return m.fI64
	}
	field := m.__.RawField(_FIELD_Main_i64)
	m.fI64 = field.GetInt64()
	m.__.Visit(_FIELD_Main_i64, _FIELD_TOTAL_Main)
	return m.fI64
}

func (m *MainEX) SetI64(v int64) {
	m.fI64 = v
	m.__.Visit(_FIELD_Main_i64, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetU64() uint64 {
	if m.__.IsVisited(_FIELD_Main_u64, _FIELD_TOTAL_Main) {
		return m.fU64
	}
	field := m.__.RawField(_FIELD_Main_u64)
	m.fU64 = field.GetUint64()
	m.__.Visit(_FIELD_Main_u64, _FIELD_TOTAL_Main)
	return m.fU64
}

func (m *MainEX) SetU64(v uint64) {
	m.fU64 = v
	m.__.Visit(_FIELD_Main_u64, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetFlag() bool {
	if m.__.IsVisited(_FIELD_Main_flag, _FIELD_TOTAL_Main) {
		return m.fFlag
	}
	field := m.__.RawField(_FIELD_Main_flag)
	m.fFlag = field.GetBool()
	m.__.Visit(_FIELD_Main_flag, _FIELD_TOTAL_Main)
	return m.fFlag
}

func (m *MainEX) SetFlag(v bool) {
	m.fFlag = v
	m.__.Visit(_FIELD_Main_flag, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetMode() Mode {
	if m.__.IsVisited(_FIELD_Main_mode, _FIELD_TOTAL_Main) {
		return m.fMode
	}
	field := m.__.RawField(_FIELD_Main_mode)
	m.fMode = Mode(field.GetEnumValue())
	m.__.Visit(_FIELD_Main_mode, _FIELD_TOTAL_Main)
	return m.fMode
}

func (m *MainEX) SetMode(v Mode) {
	m.fMode = v
	m.__.Visit(_FIELD_Main_mode, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetStr() string {
	if m.__.IsVisited(_FIELD_Main_str, _FIELD_TOTAL_Main) {
		return m.fStr
	}
	field := m.__.RawField(_FIELD_Main_str)
	m.fStr = field.GetString()
	m.__.Visit(_FIELD_Main_str, _FIELD_TOTAL_Main)
	return m.fStr
}

func (m *MainEX) SetStr(v string) {
	m.fStr = v
	m.__.Visit(_FIELD_Main_str, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetData() []byte {
	if m.__.IsVisited(_FIELD_Main_data, _FIELD_TOTAL_Main) {
		return m.fData
	}
	field := m.__.RawField(_FIELD_Main_data)
	if data := field.GetBytes(); data != nil {
		m.fData = append([]byte(nil), data...)
	}
	m.__.Visit(_FIELD_Main_data, _FIELD_TOTAL_Main)
	return m.fData
}

func (m *MainEX) SetData(v []byte) {
	if v == nil {
		m.fData = nil
	} else {
		m.fData = append([]byte(nil), v...)
	}
	m.__.Visit(_FIELD_Main_data, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetF32() float32 {
	if m.__.IsVisited(_FIELD_Main_f32, _FIELD_TOTAL_Main) {
		return m.fF32
	}
	field := m.__.RawField(_FIELD_Main_f32)
	m.fF32 = field.GetFloat32()
	m.__.Visit(_FIELD_Main_f32, _FIELD_TOTAL_Main)
	return m.fF32
}

func (m *MainEX) SetF32(v float32) {
	m.fF32 = v
	m.__.Visit(_FIELD_Main_f32, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetF64() float64 {
	if m.__.IsVisited(_FIELD_Main_f64, _FIELD_TOTAL_Main) {
		return m.fF64
	}
	field := m.__.RawField(_FIELD_Main_f64)
	m.fF64 = field.GetFloat64()
	m.__.Visit(_FIELD_Main_f64, _FIELD_TOTAL_Main)
	return m.fF64
}

func (m *MainEX) SetF64(v float64) {
	m.fF64 = v
	m.__.Visit(_FIELD_Main_f64, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetObject() *SmallEX {
	if m.__.IsVisited(_FIELD_Main_object, _FIELD_TOTAL_Main) {
		return m.fObject
	}
	field := m.__.RawField(_FIELD_Main_object)
	m.fObject = TO_SmallEX(field.GetObject())
	m.__.Visit(_FIELD_Main_object, _FIELD_TOTAL_Main)
	return m.fObject
}

func (m *MainEX) SetObject(v *SmallEX) {
	m.fObject = v
	m.__.Visit(_FIELD_Main_object, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetI32V() []int32 {
	if m.__.IsVisited(_FIELD_Main_i32v, _FIELD_TOTAL_Main) {
		return m.fI32V
	}
	field := m.__.RawField(_FIELD_Main_i32v)
	m.fI32V = append([]int32(nil), field.GetInt32Array()...)
	m.__.Visit(_FIELD_Main_i32v, _FIELD_TOTAL_Main)
	return m.fI32V
}

func (m *MainEX) SetI32V(v []int32) {
	if v == nil {
		m.fI32V = nil
	} else {
		m.fI32V = append(m.fI32V[:0], v...)
	}
	m.__.Visit(_FIELD_Main_i32v, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetU64V() []uint64 {
	if m.__.IsVisited(_FIELD_Main_u64v, _FIELD_TOTAL_Main) {
		return m.fU64V
	}
	field := m.__.RawField(_FIELD_Main_u64v)
	m.fU64V = append([]uint64(nil), field.GetUint64Array()...)
	m.__.Visit(_FIELD_Main_u64v, _FIELD_TOTAL_Main)
	return m.fU64V
}

func (m *MainEX) SetU64V(v []uint64) {
	if v == nil {
		m.fU64V = nil
	} else {
		m.fU64V = append(m.fU64V[:0], v...)
	}
	m.__.Visit(_FIELD_Main_u64v, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetStrv() []string {
	if m.__.IsVisited(_FIELD_Main_strv, _FIELD_TOTAL_Main) {
		return m.fStrv
	}
	field := m.__.RawField(_FIELD_Main_strv)
	arr := protocache.AsStringArray(field.GetObject())
	m.fStrv = make([]string, int(arr.Size()))
	for i := uint32(0); i < arr.Size(); i++ {
		m.fStrv[i] = arr.Get(i)
	}
	m.__.Visit(_FIELD_Main_strv, _FIELD_TOTAL_Main)
	return m.fStrv
}

func (m *MainEX) SetStrv(v []string) {
	if v == nil {
		m.fStrv = nil
	} else {
		m.fStrv = append(m.fStrv[:0], v...)
	}
	m.__.Visit(_FIELD_Main_strv, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetDatav() [][]byte {
	if m.__.IsVisited(_FIELD_Main_datav, _FIELD_TOTAL_Main) {
		return m.fDatav
	}
	field := m.__.RawField(_FIELD_Main_datav)
	arr := protocache.AsBytesArray(field.GetObject())
	m.fDatav = make([][]byte, int(arr.Size()))
	for i := uint32(0); i < arr.Size(); i++ {
		if data := arr.Get(i); data != nil {
			m.fDatav[i] = append([]byte(nil), data...)
		}
	}
	m.__.Visit(_FIELD_Main_datav, _FIELD_TOTAL_Main)
	return m.fDatav
}

func (m *MainEX) SetDatav(v [][]byte) {
	if v == nil {
		m.fDatav = nil
	} else {
		m.fDatav = make([][]byte, len(v))
		for i := range v {
			if v[i] != nil {
				m.fDatav[i] = append([]byte(nil), v[i]...)
			}
		}
	}
	m.__.Visit(_FIELD_Main_datav, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetF32V() []float32 {
	if m.__.IsVisited(_FIELD_Main_f32v, _FIELD_TOTAL_Main) {
		return m.fF32V
	}
	field := m.__.RawField(_FIELD_Main_f32v)
	m.fF32V = append([]float32(nil), field.GetFloat32Array()...)
	m.__.Visit(_FIELD_Main_f32v, _FIELD_TOTAL_Main)
	return m.fF32V
}

func (m *MainEX) SetF32V(v []float32) {
	if v == nil {
		m.fF32V = nil
	} else {
		m.fF32V = append(m.fF32V[:0], v...)
	}
	m.__.Visit(_FIELD_Main_f32v, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetF64V() []float64 {
	if m.__.IsVisited(_FIELD_Main_f64v, _FIELD_TOTAL_Main) {
		return m.fF64V
	}
	field := m.__.RawField(_FIELD_Main_f64v)
	m.fF64V = append([]float64(nil), field.GetFloat64Array()...)
	m.__.Visit(_FIELD_Main_f64v, _FIELD_TOTAL_Main)
	return m.fF64V
}

func (m *MainEX) SetF64V(v []float64) {
	if v == nil {
		m.fF64V = nil
	} else {
		m.fF64V = append(m.fF64V[:0], v...)
	}
	m.__.Visit(_FIELD_Main_f64v, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetFlags() []bool {
	if m.__.IsVisited(_FIELD_Main_flags, _FIELD_TOTAL_Main) {
		return m.fFlags
	}
	field := m.__.RawField(_FIELD_Main_flags)
	m.fFlags = append([]bool(nil), field.GetBoolArray()...)
	m.__.Visit(_FIELD_Main_flags, _FIELD_TOTAL_Main)
	return m.fFlags
}

func (m *MainEX) SetFlags(v []bool) {
	if v == nil {
		m.fFlags = nil
	} else {
		m.fFlags = append(m.fFlags[:0], v...)
	}
	m.__.Visit(_FIELD_Main_flags, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetObjectv() []*SmallEX {
	if m.__.IsVisited(_FIELD_Main_objectv, _FIELD_TOTAL_Main) {
		return m.fObjectv
	}
	field := m.__.RawField(_FIELD_Main_objectv)
	arr := field.GetArray()
	m.fObjectv = make([]*SmallEX, int(arr.Size()))
	for i := uint32(0); i < arr.Size(); i++ {
		elem := arr.Get(i)
		m.fObjectv[i] = TO_SmallEX(elem.GetObject())
	}
	m.__.Visit(_FIELD_Main_objectv, _FIELD_TOTAL_Main)
	return m.fObjectv
}

func (m *MainEX) SetObjectv(v []*SmallEX) {
	if v == nil {
		m.fObjectv = nil
	} else {
		m.fObjectv = append(m.fObjectv[:0], v...)
	}
	m.__.Visit(_FIELD_Main_objectv, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetTU32() uint32 {
	if m.__.IsVisited(_FIELD_Main_t_u32, _FIELD_TOTAL_Main) {
		return m.fTU32
	}
	field := m.__.RawField(_FIELD_Main_t_u32)
	m.fTU32 = field.GetUint32()
	m.__.Visit(_FIELD_Main_t_u32, _FIELD_TOTAL_Main)
	return m.fTU32
}

func (m *MainEX) SetTU32(v uint32) {
	m.fTU32 = v
	m.__.Visit(_FIELD_Main_t_u32, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetTI32() int32 {
	if m.__.IsVisited(_FIELD_Main_t_i32, _FIELD_TOTAL_Main) {
		return m.fTI32
	}
	field := m.__.RawField(_FIELD_Main_t_i32)
	m.fTI32 = field.GetInt32()
	m.__.Visit(_FIELD_Main_t_i32, _FIELD_TOTAL_Main)
	return m.fTI32
}

func (m *MainEX) SetTI32(v int32) {
	m.fTI32 = v
	m.__.Visit(_FIELD_Main_t_i32, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetTS32() int32 {
	if m.__.IsVisited(_FIELD_Main_t_s32, _FIELD_TOTAL_Main) {
		return m.fTS32
	}
	field := m.__.RawField(_FIELD_Main_t_s32)
	m.fTS32 = field.GetInt32()
	m.__.Visit(_FIELD_Main_t_s32, _FIELD_TOTAL_Main)
	return m.fTS32
}

func (m *MainEX) SetTS32(v int32) {
	m.fTS32 = v
	m.__.Visit(_FIELD_Main_t_s32, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetTU64() uint64 {
	if m.__.IsVisited(_FIELD_Main_t_u64, _FIELD_TOTAL_Main) {
		return m.fTU64
	}
	field := m.__.RawField(_FIELD_Main_t_u64)
	m.fTU64 = field.GetUint64()
	m.__.Visit(_FIELD_Main_t_u64, _FIELD_TOTAL_Main)
	return m.fTU64
}

func (m *MainEX) SetTU64(v uint64) {
	m.fTU64 = v
	m.__.Visit(_FIELD_Main_t_u64, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetTI64() int64 {
	if m.__.IsVisited(_FIELD_Main_t_i64, _FIELD_TOTAL_Main) {
		return m.fTI64
	}
	field := m.__.RawField(_FIELD_Main_t_i64)
	m.fTI64 = field.GetInt64()
	m.__.Visit(_FIELD_Main_t_i64, _FIELD_TOTAL_Main)
	return m.fTI64
}

func (m *MainEX) SetTI64(v int64) {
	m.fTI64 = v
	m.__.Visit(_FIELD_Main_t_i64, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetTS64() int64 {
	if m.__.IsVisited(_FIELD_Main_t_s64, _FIELD_TOTAL_Main) {
		return m.fTS64
	}
	field := m.__.RawField(_FIELD_Main_t_s64)
	m.fTS64 = field.GetInt64()
	m.__.Visit(_FIELD_Main_t_s64, _FIELD_TOTAL_Main)
	return m.fTS64
}

func (m *MainEX) SetTS64(v int64) {
	m.fTS64 = v
	m.__.Visit(_FIELD_Main_t_s64, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetIndex() map[string]int32 {
	if m.__.IsVisited(_FIELD_Main_index, _FIELD_TOTAL_Main) {
		return m.fIndex
	}
	field := m.__.RawField(_FIELD_Main_index)
	pack := field.GetMap()
	m.fIndex = make(map[string]int32, int(pack.Size()))
	for i := uint32(0); i < pack.Size(); i++ {
		keyField := pack.Key(i)
		key := keyField.GetString()
		valField := pack.Value(i)
		m.fIndex[key] = valField.GetInt32()
	}
	m.__.Visit(_FIELD_Main_index, _FIELD_TOTAL_Main)
	return m.fIndex
}

func (m *MainEX) SetIndex(v map[string]int32) {
	if v == nil {
		m.fIndex = nil
	} else {
		m.fIndex = make(map[string]int32, len(v))
		for k, one := range v {
			m.fIndex[k] = one
		}
	}
	m.__.Visit(_FIELD_Main_index, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetObjects() map[int32]*SmallEX {
	if m.__.IsVisited(_FIELD_Main_objects, _FIELD_TOTAL_Main) {
		return m.fObjects
	}
	field := m.__.RawField(_FIELD_Main_objects)
	pack := field.GetMap()
	m.fObjects = make(map[int32]*SmallEX, int(pack.Size()))
	for i := uint32(0); i < pack.Size(); i++ {
		keyField := pack.Key(i)
		key := keyField.GetInt32()
		valField := pack.Value(i)
		m.fObjects[key] = TO_SmallEX(valField.GetObject())
	}
	m.__.Visit(_FIELD_Main_objects, _FIELD_TOTAL_Main)
	return m.fObjects
}

func (m *MainEX) SetObjects(v map[int32]*SmallEX) {
	if v == nil {
		m.fObjects = nil
	} else {
		m.fObjects = make(map[int32]*SmallEX, len(v))
		for k, one := range v {
			m.fObjects[k] = one
		}
	}
	m.__.Visit(_FIELD_Main_objects, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetMatrix() Vec2DEX {
	if m.__.IsVisited(_FIELD_Main_matrix, _FIELD_TOTAL_Main) {
		return m.fMatrix
	}
	field := m.__.RawField(_FIELD_Main_matrix)
	m.fMatrix = TO_Vec2DEX(field.GetObject())
	m.__.Visit(_FIELD_Main_matrix, _FIELD_TOTAL_Main)
	return m.fMatrix
}

func (m *MainEX) SetMatrix(v Vec2DEX) {
	m.fMatrix = v
	m.__.Visit(_FIELD_Main_matrix, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetVector() []ArrMapEX {
	if m.__.IsVisited(_FIELD_Main_vector, _FIELD_TOTAL_Main) {
		return m.fVector
	}
	field := m.__.RawField(_FIELD_Main_vector)
	arr := field.GetArray()
	m.fVector = make([]ArrMapEX, int(arr.Size()))
	for i := uint32(0); i < arr.Size(); i++ {
		elem := arr.Get(i)
		m.fVector[i] = TO_ArrMapEX(elem.GetObject())
	}
	m.__.Visit(_FIELD_Main_vector, _FIELD_TOTAL_Main)
	return m.fVector
}

func (m *MainEX) SetVector(v []ArrMapEX) {
	if v == nil {
		m.fVector = nil
	} else {
		m.fVector = append(m.fVector[:0], v...)
	}
	m.__.Visit(_FIELD_Main_vector, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetArrays() ArrMapEX {
	if m.__.IsVisited(_FIELD_Main_arrays, _FIELD_TOTAL_Main) {
		return m.fArrays
	}
	field := m.__.RawField(_FIELD_Main_arrays)
	m.fArrays = TO_ArrMapEX(field.GetObject())
	m.__.Visit(_FIELD_Main_arrays, _FIELD_TOTAL_Main)
	return m.fArrays
}

func (m *MainEX) SetArrays(v ArrMapEX) {
	m.fArrays = v
	m.__.Visit(_FIELD_Main_arrays, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetModev() []Mode {
	if m.__.IsVisited(_FIELD_Main_modev, _FIELD_TOTAL_Main) {
		return m.fModev
	}
	field := m.__.RawField(_FIELD_Main_modev)
	m.fModev = append([]Mode(nil), protocache.CastEnumArray[Mode](field.GetEnumValueArray())...)
	m.__.Visit(_FIELD_Main_modev, _FIELD_TOTAL_Main)
	return m.fModev
}

func (m *MainEX) SetModev(v []Mode) {
	if v == nil {
		m.fModev = nil
	} else {
		m.fModev = append(m.fModev[:0], v...)
	}
	m.__.Visit(_FIELD_Main_modev, _FIELD_TOTAL_Main)
}

func DETECT_CyclicA(data []byte) []byte {
	msg := protocache.AsMessage(data)
	if !msg.IsValid() {
		return nil
	}
	inlined := msg.DetectInlined()
	if inlined == nil {
		return nil
	}
	compactEnd := len(inlined)
	if field := msg.GetField(_FIELD_CyclicA_cyclic); field.IsValid() {
		if obj := field.DetectObject(); obj != nil {
			part := func() []byte {
				if obj := field.DetectObject(); obj != nil {
					return DETECT_CyclicB(obj)
				}
				obj := field.GetObject()
				return DETECT_CyclicB(obj)
			}()
			if len(part) == 0 {
				return nil
			}
			tail := int(uintptr(unsafe.Pointer(unsafe.SliceData(obj)))-uintptr(unsafe.Pointer(unsafe.SliceData(data)))) + len(part)
			if tail > len(data) {
				return nil
			}
			return data[:tail]
		}
	}
	if compactEnd > len(data) {
		return nil
	}
	return data[:compactEnd]
}

type CyclicAEX struct {
	__      protocache.MessageEX
	fValue  int32
	fCyclic *CyclicBEX
}

func TO_CyclicAEX(data []byte) *CyclicAEX {
	out := &CyclicAEX{}
	out.__.Init(data)
	return out
}

func (m *CyclicAEX) HasBase() bool { return m.__.HasBase() }

func (m *CyclicAEX) Serialize() ([]byte, error) {
	words, err := m.serializeWords()
	if err != nil {
		return nil, err
	}
	return protocache.WordsToBytes(words), nil
}

func (m *CyclicAEX) serializeWords() ([]uint32, error) {
	if m == nil {
		return []uint32{0}, nil
	}
	parts := make([][]uint32, 2)
	if m.__.IsVisited(_FIELD_CyclicA_value, _FIELD_TOTAL_CyclicA) {
		part, err := func() ([]uint32, error) { return protocache.EncodeInt32(m.fValue), nil }()
		if err != nil {
			return nil, err
		}
		parts[0] = part
	} else {
		field := m.__.RawField(_FIELD_CyclicA_value)
		if raw := func() []byte {
			if !field.IsValid() {
				return nil
			}
			return protocache.WordsToBytes(protocache.EncodeInt32(field.GetInt32()))
		}(); len(raw) != 0 {
			parts[0] = protocache.BytesToWords(raw)
		}
	}
	if m.__.IsVisited(_FIELD_CyclicA_cyclic, _FIELD_TOTAL_CyclicA) {
		if m.fCyclic != nil {
			part, err := m.fCyclic.serializeWords()
			if err != nil {
				return nil, err
			}
			if len(part) > 1 {
				parts[1] = part
			}
		}
	} else {
		field := m.__.RawField(_FIELD_CyclicA_cyclic)
		if raw := func() []byte {
			if obj := field.DetectObject(); obj != nil {
				return DETECT_CyclicB(obj)
			}
			obj := field.GetObject()
			return DETECT_CyclicB(obj)
		}(); len(raw) != 0 {
			parts[1] = protocache.BytesToWords(raw)
		}
	}
	return protocache.EncodeMessageParts(parts)
}

func (m *CyclicAEX) GetValue() int32 {
	if m.__.IsVisited(_FIELD_CyclicA_value, _FIELD_TOTAL_CyclicA) {
		return m.fValue
	}
	field := m.__.RawField(_FIELD_CyclicA_value)
	m.fValue = field.GetInt32()
	m.__.Visit(_FIELD_CyclicA_value, _FIELD_TOTAL_CyclicA)
	return m.fValue
}

func (m *CyclicAEX) SetValue(v int32) {
	m.fValue = v
	m.__.Visit(_FIELD_CyclicA_value, _FIELD_TOTAL_CyclicA)
}

func (m *CyclicAEX) GetCyclic() *CyclicBEX {
	if m.__.IsVisited(_FIELD_CyclicA_cyclic, _FIELD_TOTAL_CyclicA) {
		return m.fCyclic
	}
	field := m.__.RawField(_FIELD_CyclicA_cyclic)
	m.fCyclic = TO_CyclicBEX(field.GetObject())
	m.__.Visit(_FIELD_CyclicA_cyclic, _FIELD_TOTAL_CyclicA)
	return m.fCyclic
}

func (m *CyclicAEX) SetCyclic(v *CyclicBEX) {
	m.fCyclic = v
	m.__.Visit(_FIELD_CyclicA_cyclic, _FIELD_TOTAL_CyclicA)
}

func DETECT_CyclicB(data []byte) []byte {
	msg := protocache.AsMessage(data)
	if !msg.IsValid() {
		return nil
	}
	inlined := msg.DetectInlined()
	if inlined == nil {
		return nil
	}
	compactEnd := len(inlined)
	if field := msg.GetField(_FIELD_CyclicB_cyclic); field.IsValid() {
		if obj := field.DetectObject(); obj != nil {
			part := func() []byte {
				if obj := field.DetectObject(); obj != nil {
					return DETECT_CyclicA(obj)
				}
				obj := field.GetObject()
				return DETECT_CyclicA(obj)
			}()
			if len(part) == 0 {
				return nil
			}
			tail := int(uintptr(unsafe.Pointer(unsafe.SliceData(obj)))-uintptr(unsafe.Pointer(unsafe.SliceData(data)))) + len(part)
			if tail > len(data) {
				return nil
			}
			return data[:tail]
		}
	}
	if compactEnd > len(data) {
		return nil
	}
	return data[:compactEnd]
}

type CyclicBEX struct {
	__      protocache.MessageEX
	fValue  int32
	fCyclic *CyclicAEX
}

func TO_CyclicBEX(data []byte) *CyclicBEX {
	out := &CyclicBEX{}
	out.__.Init(data)
	return out
}

func (m *CyclicBEX) HasBase() bool { return m.__.HasBase() }

func (m *CyclicBEX) Serialize() ([]byte, error) {
	words, err := m.serializeWords()
	if err != nil {
		return nil, err
	}
	return protocache.WordsToBytes(words), nil
}

func (m *CyclicBEX) serializeWords() ([]uint32, error) {
	if m == nil {
		return []uint32{0}, nil
	}
	parts := make([][]uint32, 2)
	if m.__.IsVisited(_FIELD_CyclicB_value, _FIELD_TOTAL_CyclicB) {
		part, err := func() ([]uint32, error) { return protocache.EncodeInt32(m.fValue), nil }()
		if err != nil {
			return nil, err
		}
		parts[0] = part
	} else {
		field := m.__.RawField(_FIELD_CyclicB_value)
		if raw := func() []byte {
			if !field.IsValid() {
				return nil
			}
			return protocache.WordsToBytes(protocache.EncodeInt32(field.GetInt32()))
		}(); len(raw) != 0 {
			parts[0] = protocache.BytesToWords(raw)
		}
	}
	if m.__.IsVisited(_FIELD_CyclicB_cyclic, _FIELD_TOTAL_CyclicB) {
		if m.fCyclic != nil {
			part, err := m.fCyclic.serializeWords()
			if err != nil {
				return nil, err
			}
			if len(part) > 1 {
				parts[1] = part
			}
		}
	} else {
		field := m.__.RawField(_FIELD_CyclicB_cyclic)
		if raw := func() []byte {
			if obj := field.DetectObject(); obj != nil {
				return DETECT_CyclicA(obj)
			}
			obj := field.GetObject()
			return DETECT_CyclicA(obj)
		}(); len(raw) != 0 {
			parts[1] = protocache.BytesToWords(raw)
		}
	}
	return protocache.EncodeMessageParts(parts)
}

func (m *CyclicBEX) GetValue() int32 {
	if m.__.IsVisited(_FIELD_CyclicB_value, _FIELD_TOTAL_CyclicB) {
		return m.fValue
	}
	field := m.__.RawField(_FIELD_CyclicB_value)
	m.fValue = field.GetInt32()
	m.__.Visit(_FIELD_CyclicB_value, _FIELD_TOTAL_CyclicB)
	return m.fValue
}

func (m *CyclicBEX) SetValue(v int32) {
	m.fValue = v
	m.__.Visit(_FIELD_CyclicB_value, _FIELD_TOTAL_CyclicB)
}

func (m *CyclicBEX) GetCyclic() *CyclicAEX {
	if m.__.IsVisited(_FIELD_CyclicB_cyclic, _FIELD_TOTAL_CyclicB) {
		return m.fCyclic
	}
	field := m.__.RawField(_FIELD_CyclicB_cyclic)
	m.fCyclic = TO_CyclicAEX(field.GetObject())
	m.__.Visit(_FIELD_CyclicB_cyclic, _FIELD_TOTAL_CyclicB)
	return m.fCyclic
}

func (m *CyclicBEX) SetCyclic(v *CyclicAEX) {
	m.fCyclic = v
	m.__.Visit(_FIELD_CyclicB_cyclic, _FIELD_TOTAL_CyclicB)
}

func DETECT_Deprecated_Valid(data []byte) []byte {
	msg := protocache.AsMessage(data)
	if !msg.IsValid() {
		return nil
	}
	inlined := msg.DetectInlined()
	if inlined == nil {
		return nil
	}
	compactEnd := len(inlined)
	if compactEnd > len(data) {
		return nil
	}
	return data[:compactEnd]
}

type Deprecated_ValidEX struct {
	__   protocache.MessageEX
	fVal int32
}

func TO_Deprecated_ValidEX(data []byte) *Deprecated_ValidEX {
	out := &Deprecated_ValidEX{}
	out.__.Init(data)
	return out
}

func (m *Deprecated_ValidEX) HasBase() bool { return m.__.HasBase() }

func (m *Deprecated_ValidEX) Serialize() ([]byte, error) {
	words, err := m.serializeWords()
	if err != nil {
		return nil, err
	}
	return protocache.WordsToBytes(words), nil
}

func (m *Deprecated_ValidEX) serializeWords() ([]uint32, error) {
	if m == nil {
		return []uint32{0}, nil
	}
	parts := make([][]uint32, 1)
	if m.__.IsVisited(_FIELD_Deprecated_Valid_val, _FIELD_TOTAL_Deprecated_Valid) {
		part, err := func() ([]uint32, error) { return protocache.EncodeInt32(m.fVal), nil }()
		if err != nil {
			return nil, err
		}
		parts[0] = part
	} else {
		field := m.__.RawField(_FIELD_Deprecated_Valid_val)
		if raw := func() []byte {
			if !field.IsValid() {
				return nil
			}
			return protocache.WordsToBytes(protocache.EncodeInt32(field.GetInt32()))
		}(); len(raw) != 0 {
			parts[0] = protocache.BytesToWords(raw)
		}
	}
	return protocache.EncodeMessageParts(parts)
}

func (m *Deprecated_ValidEX) GetVal() int32 {
	if m.__.IsVisited(_FIELD_Deprecated_Valid_val, _FIELD_TOTAL_Deprecated_Valid) {
		return m.fVal
	}
	field := m.__.RawField(_FIELD_Deprecated_Valid_val)
	m.fVal = field.GetInt32()
	m.__.Visit(_FIELD_Deprecated_Valid_val, _FIELD_TOTAL_Deprecated_Valid)
	return m.fVal
}

func (m *Deprecated_ValidEX) SetVal(v int32) {
	m.fVal = v
	m.__.Visit(_FIELD_Deprecated_Valid_val, _FIELD_TOTAL_Deprecated_Valid)
}

func DETECT_Deprecated(data []byte) []byte {
	msg := protocache.AsMessage(data)
	if !msg.IsValid() {
		return nil
	}
	inlined := msg.DetectInlined()
	if inlined == nil {
		return nil
	}
	compactEnd := len(inlined)
	if compactEnd > len(data) {
		return nil
	}
	return data[:compactEnd]
}

type DeprecatedEX struct {
	__    protocache.MessageEX
	fJunk int32
}

func TO_DeprecatedEX(data []byte) *DeprecatedEX {
	out := &DeprecatedEX{}
	out.__.Init(data)
	return out
}

func (m *DeprecatedEX) HasBase() bool { return m.__.HasBase() }

func (m *DeprecatedEX) Serialize() ([]byte, error) {
	words, err := m.serializeWords()
	if err != nil {
		return nil, err
	}
	return protocache.WordsToBytes(words), nil
}

func (m *DeprecatedEX) serializeWords() ([]uint32, error) {
	if m == nil {
		return []uint32{0}, nil
	}
	parts := make([][]uint32, 1)
	if m.__.IsVisited(_FIELD_Deprecated_junk, _FIELD_TOTAL_Deprecated) {
		part, err := func() ([]uint32, error) { return protocache.EncodeInt32(m.fJunk), nil }()
		if err != nil {
			return nil, err
		}
		parts[0] = part
	} else {
		field := m.__.RawField(_FIELD_Deprecated_junk)
		if raw := func() []byte {
			if !field.IsValid() {
				return nil
			}
			return protocache.WordsToBytes(protocache.EncodeInt32(field.GetInt32()))
		}(); len(raw) != 0 {
			parts[0] = protocache.BytesToWords(raw)
		}
	}
	return protocache.EncodeMessageParts(parts)
}

func (m *DeprecatedEX) GetJunk() int32 {
	if m.__.IsVisited(_FIELD_Deprecated_junk, _FIELD_TOTAL_Deprecated) {
		return m.fJunk
	}
	field := m.__.RawField(_FIELD_Deprecated_junk)
	m.fJunk = field.GetInt32()
	m.__.Visit(_FIELD_Deprecated_junk, _FIELD_TOTAL_Deprecated)
	return m.fJunk
}

func (m *DeprecatedEX) SetJunk(v int32) {
	m.fJunk = v
	m.__.Visit(_FIELD_Deprecated_junk, _FIELD_TOTAL_Deprecated)
}

func DETECT_ModeDict_Value(data []byte) []byte {
	return protocache.DetectArray(data, nil)
}

type ModeDict_ValueEX []Mode

func TO_ModeDict_ValueEX(data []byte) ModeDict_ValueEX {
	var out ModeDict_ValueEX
	arr := protocache.AsEnumArray[Mode](data)
	out = append([]Mode(nil), arr.Raw()...)
	return out
}

func (x ModeDict_ValueEX) Serialize() ([]byte, error) {
	words, err := serializeModeDict_ValueEX(x)
	if err != nil {
		return nil, err
	}
	return protocache.WordsToBytes(words), nil
}

func serializeModeDict_ValueEX(x ModeDict_ValueEX) ([]uint32, error) {
	if len(x) == 0 {
		return []uint32{1}, nil
	}
	return protocache.EncodeEnumArray([]Mode(x))
}

func DETECT_ModeDict(data []byte) []byte {
	return protocache.DetectMap(data, nil, DETECT_ModeDict_Value)
}

type ModeDictEX map[int32]ModeDict_ValueEX

func TO_ModeDictEX(data []byte) ModeDictEX {
	var out ModeDictEX
	pack := protocache.AsMap(data)
	out = make(map[int32]ModeDict_ValueEX, int(pack.Size()))
	for i := uint32(0); i < pack.Size(); i++ {
		keyField := pack.Key(i)
		key := keyField.GetInt32()
		valField := pack.Value(i)
		out[key] = TO_ModeDict_ValueEX(valField.GetObject())
	}
	return out
}

func (x ModeDictEX) Serialize() ([]byte, error) {
	words, err := serializeModeDictEX(x)
	if err != nil {
		return nil, err
	}
	return protocache.WordsToBytes(words), nil
}

func serializeModeDictEX(x ModeDictEX) ([]uint32, error) {
	if len(x) == 0 {
		return []uint32{5 << 28}, nil
	}
	keys := make([][]uint32, 0, len(x))
	vals := make([][]uint32, 0, len(x))
	for k, v := range x {
		keyPart, err := func() ([]uint32, error) { return protocache.EncodeInt32(k), nil }()
		if err != nil {
			return nil, err
		}
		valPart, err := func() ([]uint32, error) {
			data, err := v.Serialize()
			if err != nil {
				return nil, err
			}
			return protocache.BytesToWords(data), nil
		}()
		if err != nil {
			return nil, err
		}
		keys = append(keys, keyPart)
		vals = append(vals, valPart)
	}
	return protocache.EncodeMapParts(keys, vals, false)
}

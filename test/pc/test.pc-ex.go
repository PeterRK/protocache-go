package pc

import "github.com/peterrk/protocache-go"

func DETECT_Small(data []byte) []byte {
	msg := protocache.AsMessage(data)
	if !msg.IsValid() {
		return nil
	}
	inlined := msg.DetectInlined()
	if inlined == nil {
		return nil
	}
	field := msg.GetField(_FIELD_Small_str)
	if obj := field.DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, protocache.DetectBytes(obj))
	}
	return inlined
}

type SmallEX struct {
	source  protocache.Message
	visited [(_FIELD_TOTAL_Small + 7) / 8]byte
	fI32    int32
	fFlag   bool
	fStr    string
	fJunk   int64
}

func TO_SmallEX(data []byte) *SmallEX {
	out := &SmallEX{}
	out.source = protocache.AsMessage(data)
	return out
}

func (m *SmallEX) HasSource() bool { return m.source.IsValid() }

func ENCODE_Small(m *SmallEX) ([]uint32, error) {
	if m == nil {
		return []uint32{0}, nil
	}
	parts := make([][]uint32, _FIELD_TOTAL_Small)
	if !protocache.CheckVisited(m.visited[:], _FIELD_Small_i32) {
		field := m.source.GetField(_FIELD_Small_i32)
		parts[_FIELD_Small_i32] = field.RawWords()
	} else if m.fI32 != 0 {
		parts[_FIELD_Small_i32], _ = protocache.EncodeScalar[int32](m.fI32)
	}
	if !protocache.CheckVisited(m.visited[:], _FIELD_Small_flag) {
		field := m.source.GetField(_FIELD_Small_flag)
		parts[_FIELD_Small_flag] = field.RawWords()
	} else if m.fFlag {
		parts[_FIELD_Small_flag], _ = protocache.EncodeBool(m.fFlag)
	}
	if !protocache.CheckVisited(m.visited[:], _FIELD_Small_str) {
		field := m.source.GetField(_FIELD_Small_str)
		if obj := field.DetectObject(); obj != nil {
			parts[_FIELD_Small_str] = protocache.BytesToWords(protocache.DetectBytes(obj))
		} else {
			parts[_FIELD_Small_str] = field.RawWords()
		}
	} else if len(m.fStr) != 0 {
		part, err := protocache.EncodeString(m.fStr)
		if err != nil {
			return nil, err
		}
		parts[_FIELD_Small_str] = part
	}
	if !protocache.CheckVisited(m.visited[:], _FIELD_Small_junk) {
		field := m.source.GetField(_FIELD_Small_junk)
		parts[_FIELD_Small_junk] = field.RawWords()
	} else if m.fJunk != 0 {
		parts[_FIELD_Small_junk], _ = protocache.EncodeScalar[int64](m.fJunk)
	}
	return protocache.EncodeMessageParts(parts)
}

func (m *SmallEX) Serialize() ([]byte, error) { return protocache.SerializeWords(ENCODE_Small(m)) }

func (m *SmallEX) GetI32() int32 {
	if protocache.CheckVisited(m.visited[:], _FIELD_Small_i32) {
		return m.fI32
	}
	field := m.source.GetField(_FIELD_Small_i32)
	m.fI32 = field.GetInt32()
	protocache.Visit(m.visited[:], _FIELD_Small_i32)
	return m.fI32
}

func (m *SmallEX) SetI32(v int32) {
	m.fI32 = v
	protocache.Visit(m.visited[:], _FIELD_Small_i32)
}

func (m *SmallEX) GetFlag() bool {
	if protocache.CheckVisited(m.visited[:], _FIELD_Small_flag) {
		return m.fFlag
	}
	field := m.source.GetField(_FIELD_Small_flag)
	m.fFlag = field.GetBool()
	protocache.Visit(m.visited[:], _FIELD_Small_flag)
	return m.fFlag
}

func (m *SmallEX) SetFlag(v bool) {
	m.fFlag = v
	protocache.Visit(m.visited[:], _FIELD_Small_flag)
}

func (m *SmallEX) GetStr() string {
	if protocache.CheckVisited(m.visited[:], _FIELD_Small_str) {
		return m.fStr
	}
	field := m.source.GetField(_FIELD_Small_str)
	m.fStr = field.GetString()
	protocache.Visit(m.visited[:], _FIELD_Small_str)
	return m.fStr
}

func (m *SmallEX) SetStr(v string) {
	m.fStr = v
	protocache.Visit(m.visited[:], _FIELD_Small_str)
}

func (m *SmallEX) GetJunk() int64 {
	if protocache.CheckVisited(m.visited[:], _FIELD_Small_junk) {
		return m.fJunk
	}
	field := m.source.GetField(_FIELD_Small_junk)
	m.fJunk = field.GetInt64()
	protocache.Visit(m.visited[:], _FIELD_Small_junk)
	return m.fJunk
}

func (m *SmallEX) SetJunk(v int64) {
	m.fJunk = v
	protocache.Visit(m.visited[:], _FIELD_Small_junk)
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

func ENCODE_Vec2D_Vec1D(x Vec2D_Vec1DEX) ([]uint32, error) {
	if len(x) == 0 {
		return []uint32{1}, nil
	}
	return protocache.EncodeScalarVector[float32]([]float32(x))
}

func (x Vec2D_Vec1DEX) Serialize() ([]byte, error) {
	return protocache.SerializeWords(ENCODE_Vec2D_Vec1D(x))
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

func ENCODE_Vec2D(x Vec2DEX) ([]uint32, error) {
	if len(x) == 0 {
		return []uint32{1}, nil
	}
	return protocache.EncodeObjectArray(x, ENCODE_Vec2D_Vec1D)
}

func (x Vec2DEX) Serialize() ([]byte, error) { return protocache.SerializeWords(ENCODE_Vec2D(x)) }

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

func ENCODE_ArrMap_Array(x ArrMap_ArrayEX) ([]uint32, error) {
	if len(x) == 0 {
		return []uint32{1}, nil
	}
	return protocache.EncodeScalarVector[float32]([]float32(x))
}

func (x ArrMap_ArrayEX) Serialize() ([]byte, error) {
	return protocache.SerializeWords(ENCODE_ArrMap_Array(x))
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

func ENCODE_ArrMap(x ArrMapEX) ([]uint32, error) {
	return protocache.EncodeStringMap(x, ENCODE_ArrMap_Array)
}

func (x ArrMapEX) Serialize() ([]byte, error) { return protocache.SerializeWords(ENCODE_ArrMap(x)) }

func DETECT_Main(data []byte) []byte {
	msg := protocache.AsMessage(data)
	if !msg.IsValid() {
		return nil
	}
	inlined := msg.DetectInlined()
	if inlined == nil {
		return nil
	}
	field := msg.GetField(_FIELD_Main_modev)
	if obj := field.DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, protocache.DetectArray(obj, nil))
	}
	field = msg.GetField(_FIELD_Main_arrays)
	if obj := field.DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, DETECT_ArrMap(obj))
	}
	field = msg.GetField(_FIELD_Main_vector)
	if obj := field.DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, protocache.DetectArray(obj, DETECT_ArrMap))
	}
	field = msg.GetField(_FIELD_Main_matrix)
	if obj := field.DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, DETECT_Vec2D(obj))
	}
	field = msg.GetField(_FIELD_Main_objects)
	if obj := field.DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, protocache.DetectMap(obj, nil, DETECT_Small))
	}
	field = msg.GetField(_FIELD_Main_index)
	if obj := field.DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, protocache.DetectMap(obj, protocache.DetectBytes, nil))
	}
	field = msg.GetField(_FIELD_Main_objectv)
	if obj := field.DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, protocache.DetectArray(obj, DETECT_Small))
	}
	field = msg.GetField(_FIELD_Main_flags)
	if obj := field.DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, protocache.DetectBytes(obj))
	}
	field = msg.GetField(_FIELD_Main_f64v)
	if obj := field.DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, protocache.DetectArray(obj, nil))
	}
	field = msg.GetField(_FIELD_Main_f32v)
	if obj := field.DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, protocache.DetectArray(obj, nil))
	}
	field = msg.GetField(_FIELD_Main_datav)
	if obj := field.DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, protocache.DetectArray(obj, protocache.DetectBytes))
	}
	field = msg.GetField(_FIELD_Main_strv)
	if obj := field.DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, protocache.DetectArray(obj, protocache.DetectBytes))
	}
	field = msg.GetField(_FIELD_Main_u64v)
	if obj := field.DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, protocache.DetectArray(obj, nil))
	}
	field = msg.GetField(_FIELD_Main_i32v)
	if obj := field.DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, protocache.DetectArray(obj, nil))
	}
	field = msg.GetField(_FIELD_Main_object)
	if obj := field.DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, DETECT_Small(obj))
	}
	field = msg.GetField(_FIELD_Main_data)
	if obj := field.DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, protocache.DetectBytes(obj))
	}
	field = msg.GetField(_FIELD_Main_str)
	if obj := field.DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, protocache.DetectBytes(obj))
	}
	return inlined
}

type MainEX struct {
	source   protocache.Message
	visited  [(_FIELD_TOTAL_Main + 7) / 8]byte
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
	out.source = protocache.AsMessage(data)
	return out
}

func (m *MainEX) HasSource() bool { return m.source.IsValid() }

func ENCODE_Main(m *MainEX) ([]uint32, error) {
	if m == nil {
		return []uint32{0}, nil
	}
	parts := make([][]uint32, _FIELD_TOTAL_Main)
	if !protocache.CheckVisited(m.visited[:], _FIELD_Main_i32) {
		field := m.source.GetField(_FIELD_Main_i32)
		parts[_FIELD_Main_i32] = field.RawWords()
	} else if m.fI32 != 0 {
		parts[_FIELD_Main_i32], _ = protocache.EncodeScalar[int32](m.fI32)
	}
	if !protocache.CheckVisited(m.visited[:], _FIELD_Main_u32) {
		field := m.source.GetField(_FIELD_Main_u32)
		parts[_FIELD_Main_u32] = field.RawWords()
	} else if m.fU32 != 0 {
		parts[_FIELD_Main_u32], _ = protocache.EncodeScalar[uint32](m.fU32)
	}
	if !protocache.CheckVisited(m.visited[:], _FIELD_Main_i64) {
		field := m.source.GetField(_FIELD_Main_i64)
		parts[_FIELD_Main_i64] = field.RawWords()
	} else if m.fI64 != 0 {
		parts[_FIELD_Main_i64], _ = protocache.EncodeScalar[int64](m.fI64)
	}
	if !protocache.CheckVisited(m.visited[:], _FIELD_Main_u64) {
		field := m.source.GetField(_FIELD_Main_u64)
		parts[_FIELD_Main_u64] = field.RawWords()
	} else if m.fU64 != 0 {
		parts[_FIELD_Main_u64], _ = protocache.EncodeScalar[uint64](m.fU64)
	}
	if !protocache.CheckVisited(m.visited[:], _FIELD_Main_flag) {
		field := m.source.GetField(_FIELD_Main_flag)
		parts[_FIELD_Main_flag] = field.RawWords()
	} else if m.fFlag {
		parts[_FIELD_Main_flag], _ = protocache.EncodeBool(m.fFlag)
	}
	if !protocache.CheckVisited(m.visited[:], _FIELD_Main_mode) {
		field := m.source.GetField(_FIELD_Main_mode)
		parts[_FIELD_Main_mode] = field.RawWords()
	} else if m.fMode != 0 {
		parts[_FIELD_Main_mode], _ = protocache.EncodeScalar[int32](int32(m.fMode))
	}
	if !protocache.CheckVisited(m.visited[:], _FIELD_Main_str) {
		field := m.source.GetField(_FIELD_Main_str)
		if obj := field.DetectObject(); obj != nil {
			parts[_FIELD_Main_str] = protocache.BytesToWords(protocache.DetectBytes(obj))
		} else {
			parts[_FIELD_Main_str] = field.RawWords()
		}
	} else if len(m.fStr) != 0 {
		part, err := protocache.EncodeString(m.fStr)
		if err != nil {
			return nil, err
		}
		parts[_FIELD_Main_str] = part
	}
	if !protocache.CheckVisited(m.visited[:], _FIELD_Main_data) {
		field := m.source.GetField(_FIELD_Main_data)
		if obj := field.DetectObject(); obj != nil {
			parts[_FIELD_Main_data] = protocache.BytesToWords(protocache.DetectBytes(obj))
		} else {
			parts[_FIELD_Main_data] = field.RawWords()
		}
	} else if len(m.fData) != 0 {
		part, err := protocache.EncodeBytes(m.fData)
		if err != nil {
			return nil, err
		}
		parts[_FIELD_Main_data] = part
	}
	if !protocache.CheckVisited(m.visited[:], _FIELD_Main_f32) {
		field := m.source.GetField(_FIELD_Main_f32)
		parts[_FIELD_Main_f32] = field.RawWords()
	} else if m.fF32 != 0 {
		parts[_FIELD_Main_f32], _ = protocache.EncodeScalar[float32](m.fF32)
	}
	if !protocache.CheckVisited(m.visited[:], _FIELD_Main_f64) {
		field := m.source.GetField(_FIELD_Main_f64)
		parts[_FIELD_Main_f64] = field.RawWords()
	} else if m.fF64 != 0 {
		parts[_FIELD_Main_f64], _ = protocache.EncodeScalar[float64](m.fF64)
	}
	if !protocache.CheckVisited(m.visited[:], _FIELD_Main_object) {
		field := m.source.GetField(_FIELD_Main_object)
		if obj := field.DetectObject(); obj != nil {
			parts[_FIELD_Main_object] = protocache.BytesToWords(DETECT_Small(obj))
		} else {
			parts[_FIELD_Main_object] = field.RawWords()
		}
	} else if m.fObject != nil {
		part, err := ENCODE_Small(m.fObject)
		if err != nil {
			return nil, err
		}
		if len(part) > 1 {
			parts[_FIELD_Main_object] = part
		}
	}
	if !protocache.CheckVisited(m.visited[:], _FIELD_Main_i32v) {
		field := m.source.GetField(_FIELD_Main_i32v)
		if obj := field.DetectObject(); obj != nil {
			parts[_FIELD_Main_i32v] = protocache.BytesToWords(protocache.DetectArray(obj, nil))
		} else {
			parts[_FIELD_Main_i32v] = field.RawWords()
		}
	} else if len(m.fI32V) != 0 {
		part, err := protocache.EncodeScalarVector[int32](m.fI32V)
		if err != nil {
			return nil, err
		}
		parts[_FIELD_Main_i32v] = part
	}
	if !protocache.CheckVisited(m.visited[:], _FIELD_Main_u64v) {
		field := m.source.GetField(_FIELD_Main_u64v)
		if obj := field.DetectObject(); obj != nil {
			parts[_FIELD_Main_u64v] = protocache.BytesToWords(protocache.DetectArray(obj, nil))
		} else {
			parts[_FIELD_Main_u64v] = field.RawWords()
		}
	} else if len(m.fU64V) != 0 {
		part, err := protocache.EncodeScalarVector[uint64](m.fU64V)
		if err != nil {
			return nil, err
		}
		parts[_FIELD_Main_u64v] = part
	}
	if !protocache.CheckVisited(m.visited[:], _FIELD_Main_strv) {
		field := m.source.GetField(_FIELD_Main_strv)
		if obj := field.DetectObject(); obj != nil {
			parts[_FIELD_Main_strv] = protocache.BytesToWords(protocache.DetectArray(obj, protocache.DetectBytes))
		} else {
			parts[_FIELD_Main_strv] = field.RawWords()
		}
	} else if len(m.fStrv) != 0 {
		part, err := protocache.EncodeStringArray(m.fStrv)
		if err != nil {
			return nil, err
		}
		parts[_FIELD_Main_strv] = part
	}
	if !protocache.CheckVisited(m.visited[:], _FIELD_Main_datav) {
		field := m.source.GetField(_FIELD_Main_datav)
		if obj := field.DetectObject(); obj != nil {
			parts[_FIELD_Main_datav] = protocache.BytesToWords(protocache.DetectArray(obj, protocache.DetectBytes))
		} else {
			parts[_FIELD_Main_datav] = field.RawWords()
		}
	} else if len(m.fDatav) != 0 {
		part, err := protocache.EncodeBytesArray(m.fDatav)
		if err != nil {
			return nil, err
		}
		parts[_FIELD_Main_datav] = part
	}
	if !protocache.CheckVisited(m.visited[:], _FIELD_Main_f32v) {
		field := m.source.GetField(_FIELD_Main_f32v)
		if obj := field.DetectObject(); obj != nil {
			parts[_FIELD_Main_f32v] = protocache.BytesToWords(protocache.DetectArray(obj, nil))
		} else {
			parts[_FIELD_Main_f32v] = field.RawWords()
		}
	} else if len(m.fF32V) != 0 {
		part, err := protocache.EncodeScalarVector[float32](m.fF32V)
		if err != nil {
			return nil, err
		}
		parts[_FIELD_Main_f32v] = part
	}
	if !protocache.CheckVisited(m.visited[:], _FIELD_Main_f64v) {
		field := m.source.GetField(_FIELD_Main_f64v)
		if obj := field.DetectObject(); obj != nil {
			parts[_FIELD_Main_f64v] = protocache.BytesToWords(protocache.DetectArray(obj, nil))
		} else {
			parts[_FIELD_Main_f64v] = field.RawWords()
		}
	} else if len(m.fF64V) != 0 {
		part, err := protocache.EncodeScalarVector[float64](m.fF64V)
		if err != nil {
			return nil, err
		}
		parts[_FIELD_Main_f64v] = part
	}
	if !protocache.CheckVisited(m.visited[:], _FIELD_Main_flags) {
		field := m.source.GetField(_FIELD_Main_flags)
		if obj := field.DetectObject(); obj != nil {
			parts[_FIELD_Main_flags] = protocache.BytesToWords(protocache.DetectBytes(obj))
		} else {
			parts[_FIELD_Main_flags] = field.RawWords()
		}
	} else if len(m.fFlags) != 0 {
		part, err := protocache.EncodeBoolArray(m.fFlags)
		if err != nil {
			return nil, err
		}
		parts[_FIELD_Main_flags] = part
	}
	if !protocache.CheckVisited(m.visited[:], _FIELD_Main_objectv) {
		field := m.source.GetField(_FIELD_Main_objectv)
		if obj := field.DetectObject(); obj != nil {
			parts[_FIELD_Main_objectv] = protocache.BytesToWords(protocache.DetectArray(obj, DETECT_Small))
		} else {
			parts[_FIELD_Main_objectv] = field.RawWords()
		}
	} else if len(m.fObjectv) != 0 {
		part, err := protocache.EncodeObjectArray(m.fObjectv, ENCODE_Small)
		if err != nil {
			return nil, err
		}
		parts[_FIELD_Main_objectv] = part
	}
	if !protocache.CheckVisited(m.visited[:], _FIELD_Main_t_u32) {
		field := m.source.GetField(_FIELD_Main_t_u32)
		parts[_FIELD_Main_t_u32] = field.RawWords()
	} else if m.fTU32 != 0 {
		parts[_FIELD_Main_t_u32], _ = protocache.EncodeScalar[uint32](m.fTU32)
	}
	if !protocache.CheckVisited(m.visited[:], _FIELD_Main_t_i32) {
		field := m.source.GetField(_FIELD_Main_t_i32)
		parts[_FIELD_Main_t_i32] = field.RawWords()
	} else if m.fTI32 != 0 {
		parts[_FIELD_Main_t_i32], _ = protocache.EncodeScalar[int32](m.fTI32)
	}
	if !protocache.CheckVisited(m.visited[:], _FIELD_Main_t_s32) {
		field := m.source.GetField(_FIELD_Main_t_s32)
		parts[_FIELD_Main_t_s32] = field.RawWords()
	} else if m.fTS32 != 0 {
		parts[_FIELD_Main_t_s32], _ = protocache.EncodeScalar[int32](m.fTS32)
	}
	if !protocache.CheckVisited(m.visited[:], _FIELD_Main_t_u64) {
		field := m.source.GetField(_FIELD_Main_t_u64)
		parts[_FIELD_Main_t_u64] = field.RawWords()
	} else if m.fTU64 != 0 {
		parts[_FIELD_Main_t_u64], _ = protocache.EncodeScalar[uint64](m.fTU64)
	}
	if !protocache.CheckVisited(m.visited[:], _FIELD_Main_t_i64) {
		field := m.source.GetField(_FIELD_Main_t_i64)
		parts[_FIELD_Main_t_i64] = field.RawWords()
	} else if m.fTI64 != 0 {
		parts[_FIELD_Main_t_i64], _ = protocache.EncodeScalar[int64](m.fTI64)
	}
	if !protocache.CheckVisited(m.visited[:], _FIELD_Main_t_s64) {
		field := m.source.GetField(_FIELD_Main_t_s64)
		parts[_FIELD_Main_t_s64] = field.RawWords()
	} else if m.fTS64 != 0 {
		parts[_FIELD_Main_t_s64], _ = protocache.EncodeScalar[int64](m.fTS64)
	}
	if !protocache.CheckVisited(m.visited[:], _FIELD_Main_index) {
		field := m.source.GetField(_FIELD_Main_index)
		if obj := field.DetectObject(); obj != nil {
			parts[_FIELD_Main_index] = protocache.BytesToWords(protocache.DetectMap(obj, protocache.DetectBytes, nil))
		} else {
			parts[_FIELD_Main_index] = field.RawWords()
		}
	} else if len(m.fIndex) != 0 {
		part, err := protocache.EncodeStringMap(m.fIndex, protocache.EncodeScalar[int32])
		if err != nil {
			return nil, err
		}
		parts[_FIELD_Main_index] = part
	}
	if !protocache.CheckVisited(m.visited[:], _FIELD_Main_objects) {
		field := m.source.GetField(_FIELD_Main_objects)
		if obj := field.DetectObject(); obj != nil {
			parts[_FIELD_Main_objects] = protocache.BytesToWords(protocache.DetectMap(obj, nil, DETECT_Small))
		} else {
			parts[_FIELD_Main_objects] = field.RawWords()
		}
	} else if len(m.fObjects) != 0 {
		part, err := protocache.EncodeScalarMap(m.fObjects,
			protocache.EncodeScalar[int32],
			ENCODE_Small)
		if err != nil {
			return nil, err
		}
		parts[_FIELD_Main_objects] = part
	}
	if !protocache.CheckVisited(m.visited[:], _FIELD_Main_matrix) {
		field := m.source.GetField(_FIELD_Main_matrix)
		if obj := field.DetectObject(); obj != nil {
			parts[_FIELD_Main_matrix] = protocache.BytesToWords(DETECT_Vec2D(obj))
		} else {
			parts[_FIELD_Main_matrix] = field.RawWords()
		}
	} else if m.fMatrix != nil {
		part, err := ENCODE_Vec2D(m.fMatrix)
		if err != nil {
			return nil, err
		}
		if len(part) > 1 {
			parts[_FIELD_Main_matrix] = part
		}
	}
	if !protocache.CheckVisited(m.visited[:], _FIELD_Main_vector) {
		field := m.source.GetField(_FIELD_Main_vector)
		if obj := field.DetectObject(); obj != nil {
			parts[_FIELD_Main_vector] = protocache.BytesToWords(protocache.DetectArray(obj, DETECT_ArrMap))
		} else {
			parts[_FIELD_Main_vector] = field.RawWords()
		}
	} else if len(m.fVector) != 0 {
		part, err := protocache.EncodeObjectArray(m.fVector, ENCODE_ArrMap)
		if err != nil {
			return nil, err
		}
		parts[_FIELD_Main_vector] = part
	}
	if !protocache.CheckVisited(m.visited[:], _FIELD_Main_arrays) {
		field := m.source.GetField(_FIELD_Main_arrays)
		if obj := field.DetectObject(); obj != nil {
			parts[_FIELD_Main_arrays] = protocache.BytesToWords(DETECT_ArrMap(obj))
		} else {
			parts[_FIELD_Main_arrays] = field.RawWords()
		}
	} else if m.fArrays != nil {
		part, err := ENCODE_ArrMap(m.fArrays)
		if err != nil {
			return nil, err
		}
		if len(part) > 1 {
			parts[_FIELD_Main_arrays] = part
		}
	}
	if !protocache.CheckVisited(m.visited[:], _FIELD_Main_modev) {
		field := m.source.GetField(_FIELD_Main_modev)
		if obj := field.DetectObject(); obj != nil {
			parts[_FIELD_Main_modev] = protocache.BytesToWords(protocache.DetectArray(obj, nil))
		} else {
			parts[_FIELD_Main_modev] = field.RawWords()
		}
	} else if len(m.fModev) != 0 {
		part, err := protocache.EncodeEnumArray(m.fModev)
		if err != nil {
			return nil, err
		}
		parts[_FIELD_Main_modev] = part
	}
	return protocache.EncodeMessageParts(parts)
}

func (m *MainEX) Serialize() ([]byte, error) { return protocache.SerializeWords(ENCODE_Main(m)) }

func (m *MainEX) GetI32() int32 {
	if protocache.CheckVisited(m.visited[:], _FIELD_Main_i32) {
		return m.fI32
	}
	field := m.source.GetField(_FIELD_Main_i32)
	m.fI32 = field.GetInt32()
	protocache.Visit(m.visited[:], _FIELD_Main_i32)
	return m.fI32
}

func (m *MainEX) SetI32(v int32) {
	m.fI32 = v
	protocache.Visit(m.visited[:], _FIELD_Main_i32)
}

func (m *MainEX) GetU32() uint32 {
	if protocache.CheckVisited(m.visited[:], _FIELD_Main_u32) {
		return m.fU32
	}
	field := m.source.GetField(_FIELD_Main_u32)
	m.fU32 = field.GetUint32()
	protocache.Visit(m.visited[:], _FIELD_Main_u32)
	return m.fU32
}

func (m *MainEX) SetU32(v uint32) {
	m.fU32 = v
	protocache.Visit(m.visited[:], _FIELD_Main_u32)
}

func (m *MainEX) GetI64() int64 {
	if protocache.CheckVisited(m.visited[:], _FIELD_Main_i64) {
		return m.fI64
	}
	field := m.source.GetField(_FIELD_Main_i64)
	m.fI64 = field.GetInt64()
	protocache.Visit(m.visited[:], _FIELD_Main_i64)
	return m.fI64
}

func (m *MainEX) SetI64(v int64) {
	m.fI64 = v
	protocache.Visit(m.visited[:], _FIELD_Main_i64)
}

func (m *MainEX) GetU64() uint64 {
	if protocache.CheckVisited(m.visited[:], _FIELD_Main_u64) {
		return m.fU64
	}
	field := m.source.GetField(_FIELD_Main_u64)
	m.fU64 = field.GetUint64()
	protocache.Visit(m.visited[:], _FIELD_Main_u64)
	return m.fU64
}

func (m *MainEX) SetU64(v uint64) {
	m.fU64 = v
	protocache.Visit(m.visited[:], _FIELD_Main_u64)
}

func (m *MainEX) GetFlag() bool {
	if protocache.CheckVisited(m.visited[:], _FIELD_Main_flag) {
		return m.fFlag
	}
	field := m.source.GetField(_FIELD_Main_flag)
	m.fFlag = field.GetBool()
	protocache.Visit(m.visited[:], _FIELD_Main_flag)
	return m.fFlag
}

func (m *MainEX) SetFlag(v bool) {
	m.fFlag = v
	protocache.Visit(m.visited[:], _FIELD_Main_flag)
}

func (m *MainEX) GetMode() Mode {
	if protocache.CheckVisited(m.visited[:], _FIELD_Main_mode) {
		return m.fMode
	}
	field := m.source.GetField(_FIELD_Main_mode)
	m.fMode = Mode(field.GetEnumValue())
	protocache.Visit(m.visited[:], _FIELD_Main_mode)
	return m.fMode
}

func (m *MainEX) SetMode(v Mode) {
	m.fMode = v
	protocache.Visit(m.visited[:], _FIELD_Main_mode)
}

func (m *MainEX) GetStr() string {
	if protocache.CheckVisited(m.visited[:], _FIELD_Main_str) {
		return m.fStr
	}
	field := m.source.GetField(_FIELD_Main_str)
	m.fStr = field.GetString()
	protocache.Visit(m.visited[:], _FIELD_Main_str)
	return m.fStr
}

func (m *MainEX) SetStr(v string) {
	m.fStr = v
	protocache.Visit(m.visited[:], _FIELD_Main_str)
}

func (m *MainEX) GetData() []byte {
	if protocache.CheckVisited(m.visited[:], _FIELD_Main_data) {
		return m.fData
	}
	field := m.source.GetField(_FIELD_Main_data)
	if data := field.GetBytes(); data != nil {
		m.fData = append([]byte(nil), data...)
	}
	protocache.Visit(m.visited[:], _FIELD_Main_data)
	return m.fData
}

func (m *MainEX) SetData(v []byte) {
	if v == nil {
		m.fData = nil
	} else {
		m.fData = append([]byte(nil), v...)
	}
	protocache.Visit(m.visited[:], _FIELD_Main_data)
}

func (m *MainEX) GetF32() float32 {
	if protocache.CheckVisited(m.visited[:], _FIELD_Main_f32) {
		return m.fF32
	}
	field := m.source.GetField(_FIELD_Main_f32)
	m.fF32 = field.GetFloat32()
	protocache.Visit(m.visited[:], _FIELD_Main_f32)
	return m.fF32
}

func (m *MainEX) SetF32(v float32) {
	m.fF32 = v
	protocache.Visit(m.visited[:], _FIELD_Main_f32)
}

func (m *MainEX) GetF64() float64 {
	if protocache.CheckVisited(m.visited[:], _FIELD_Main_f64) {
		return m.fF64
	}
	field := m.source.GetField(_FIELD_Main_f64)
	m.fF64 = field.GetFloat64()
	protocache.Visit(m.visited[:], _FIELD_Main_f64)
	return m.fF64
}

func (m *MainEX) SetF64(v float64) {
	m.fF64 = v
	protocache.Visit(m.visited[:], _FIELD_Main_f64)
}

func (m *MainEX) GetObject() *SmallEX {
	if protocache.CheckVisited(m.visited[:], _FIELD_Main_object) {
		return m.fObject
	}
	field := m.source.GetField(_FIELD_Main_object)
	m.fObject = TO_SmallEX(field.GetObject())
	protocache.Visit(m.visited[:], _FIELD_Main_object)
	return m.fObject
}

func (m *MainEX) SetObject(v *SmallEX) {
	m.fObject = v
	protocache.Visit(m.visited[:], _FIELD_Main_object)
}

func (m *MainEX) GetI32V() []int32 {
	if protocache.CheckVisited(m.visited[:], _FIELD_Main_i32v) {
		return m.fI32V
	}
	field := m.source.GetField(_FIELD_Main_i32v)
	m.fI32V = append([]int32(nil), field.GetInt32Array()...)
	protocache.Visit(m.visited[:], _FIELD_Main_i32v)
	return m.fI32V
}

func (m *MainEX) SetI32V(v []int32) {
	if v == nil {
		m.fI32V = nil
	} else {
		m.fI32V = append(m.fI32V[:0], v...)
	}
	protocache.Visit(m.visited[:], _FIELD_Main_i32v)
}

func (m *MainEX) GetU64V() []uint64 {
	if protocache.CheckVisited(m.visited[:], _FIELD_Main_u64v) {
		return m.fU64V
	}
	field := m.source.GetField(_FIELD_Main_u64v)
	m.fU64V = append([]uint64(nil), field.GetUint64Array()...)
	protocache.Visit(m.visited[:], _FIELD_Main_u64v)
	return m.fU64V
}

func (m *MainEX) SetU64V(v []uint64) {
	if v == nil {
		m.fU64V = nil
	} else {
		m.fU64V = append(m.fU64V[:0], v...)
	}
	protocache.Visit(m.visited[:], _FIELD_Main_u64v)
}

func (m *MainEX) GetStrv() []string {
	if protocache.CheckVisited(m.visited[:], _FIELD_Main_strv) {
		return m.fStrv
	}
	field := m.source.GetField(_FIELD_Main_strv)
	arr := protocache.AsStringArray(field.GetObject())
	m.fStrv = make([]string, int(arr.Size()))
	for i := uint32(0); i < arr.Size(); i++ {
		m.fStrv[i] = arr.Get(i)
	}
	protocache.Visit(m.visited[:], _FIELD_Main_strv)
	return m.fStrv
}

func (m *MainEX) SetStrv(v []string) {
	if v == nil {
		m.fStrv = nil
	} else {
		m.fStrv = append(m.fStrv[:0], v...)
	}
	protocache.Visit(m.visited[:], _FIELD_Main_strv)
}

func (m *MainEX) GetDatav() [][]byte {
	if protocache.CheckVisited(m.visited[:], _FIELD_Main_datav) {
		return m.fDatav
	}
	field := m.source.GetField(_FIELD_Main_datav)
	arr := protocache.AsBytesArray(field.GetObject())
	m.fDatav = make([][]byte, int(arr.Size()))
	for i := uint32(0); i < arr.Size(); i++ {
		if data := arr.Get(i); data != nil {
			m.fDatav[i] = append([]byte(nil), data...)
		}
	}
	protocache.Visit(m.visited[:], _FIELD_Main_datav)
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
	protocache.Visit(m.visited[:], _FIELD_Main_datav)
}

func (m *MainEX) GetF32V() []float32 {
	if protocache.CheckVisited(m.visited[:], _FIELD_Main_f32v) {
		return m.fF32V
	}
	field := m.source.GetField(_FIELD_Main_f32v)
	m.fF32V = append([]float32(nil), field.GetFloat32Array()...)
	protocache.Visit(m.visited[:], _FIELD_Main_f32v)
	return m.fF32V
}

func (m *MainEX) SetF32V(v []float32) {
	if v == nil {
		m.fF32V = nil
	} else {
		m.fF32V = append(m.fF32V[:0], v...)
	}
	protocache.Visit(m.visited[:], _FIELD_Main_f32v)
}

func (m *MainEX) GetF64V() []float64 {
	if protocache.CheckVisited(m.visited[:], _FIELD_Main_f64v) {
		return m.fF64V
	}
	field := m.source.GetField(_FIELD_Main_f64v)
	m.fF64V = append([]float64(nil), field.GetFloat64Array()...)
	protocache.Visit(m.visited[:], _FIELD_Main_f64v)
	return m.fF64V
}

func (m *MainEX) SetF64V(v []float64) {
	if v == nil {
		m.fF64V = nil
	} else {
		m.fF64V = append(m.fF64V[:0], v...)
	}
	protocache.Visit(m.visited[:], _FIELD_Main_f64v)
}

func (m *MainEX) GetFlags() []bool {
	if protocache.CheckVisited(m.visited[:], _FIELD_Main_flags) {
		return m.fFlags
	}
	field := m.source.GetField(_FIELD_Main_flags)
	m.fFlags = append([]bool(nil), field.GetBoolArray()...)
	protocache.Visit(m.visited[:], _FIELD_Main_flags)
	return m.fFlags
}

func (m *MainEX) SetFlags(v []bool) {
	if v == nil {
		m.fFlags = nil
	} else {
		m.fFlags = append(m.fFlags[:0], v...)
	}
	protocache.Visit(m.visited[:], _FIELD_Main_flags)
}

func (m *MainEX) GetObjectv() []*SmallEX {
	if protocache.CheckVisited(m.visited[:], _FIELD_Main_objectv) {
		return m.fObjectv
	}
	field := m.source.GetField(_FIELD_Main_objectv)
	arr := field.GetArray()
	m.fObjectv = make([]*SmallEX, int(arr.Size()))
	for i := uint32(0); i < arr.Size(); i++ {
		elem := arr.Get(i)
		m.fObjectv[i] = TO_SmallEX(elem.GetObject())
	}
	protocache.Visit(m.visited[:], _FIELD_Main_objectv)
	return m.fObjectv
}

func (m *MainEX) SetObjectv(v []*SmallEX) {
	if v == nil {
		m.fObjectv = nil
	} else {
		m.fObjectv = append(m.fObjectv[:0], v...)
	}
	protocache.Visit(m.visited[:], _FIELD_Main_objectv)
}

func (m *MainEX) GetTU32() uint32 {
	if protocache.CheckVisited(m.visited[:], _FIELD_Main_t_u32) {
		return m.fTU32
	}
	field := m.source.GetField(_FIELD_Main_t_u32)
	m.fTU32 = field.GetUint32()
	protocache.Visit(m.visited[:], _FIELD_Main_t_u32)
	return m.fTU32
}

func (m *MainEX) SetTU32(v uint32) {
	m.fTU32 = v
	protocache.Visit(m.visited[:], _FIELD_Main_t_u32)
}

func (m *MainEX) GetTI32() int32 {
	if protocache.CheckVisited(m.visited[:], _FIELD_Main_t_i32) {
		return m.fTI32
	}
	field := m.source.GetField(_FIELD_Main_t_i32)
	m.fTI32 = field.GetInt32()
	protocache.Visit(m.visited[:], _FIELD_Main_t_i32)
	return m.fTI32
}

func (m *MainEX) SetTI32(v int32) {
	m.fTI32 = v
	protocache.Visit(m.visited[:], _FIELD_Main_t_i32)
}

func (m *MainEX) GetTS32() int32 {
	if protocache.CheckVisited(m.visited[:], _FIELD_Main_t_s32) {
		return m.fTS32
	}
	field := m.source.GetField(_FIELD_Main_t_s32)
	m.fTS32 = field.GetInt32()
	protocache.Visit(m.visited[:], _FIELD_Main_t_s32)
	return m.fTS32
}

func (m *MainEX) SetTS32(v int32) {
	m.fTS32 = v
	protocache.Visit(m.visited[:], _FIELD_Main_t_s32)
}

func (m *MainEX) GetTU64() uint64 {
	if protocache.CheckVisited(m.visited[:], _FIELD_Main_t_u64) {
		return m.fTU64
	}
	field := m.source.GetField(_FIELD_Main_t_u64)
	m.fTU64 = field.GetUint64()
	protocache.Visit(m.visited[:], _FIELD_Main_t_u64)
	return m.fTU64
}

func (m *MainEX) SetTU64(v uint64) {
	m.fTU64 = v
	protocache.Visit(m.visited[:], _FIELD_Main_t_u64)
}

func (m *MainEX) GetTI64() int64 {
	if protocache.CheckVisited(m.visited[:], _FIELD_Main_t_i64) {
		return m.fTI64
	}
	field := m.source.GetField(_FIELD_Main_t_i64)
	m.fTI64 = field.GetInt64()
	protocache.Visit(m.visited[:], _FIELD_Main_t_i64)
	return m.fTI64
}

func (m *MainEX) SetTI64(v int64) {
	m.fTI64 = v
	protocache.Visit(m.visited[:], _FIELD_Main_t_i64)
}

func (m *MainEX) GetTS64() int64 {
	if protocache.CheckVisited(m.visited[:], _FIELD_Main_t_s64) {
		return m.fTS64
	}
	field := m.source.GetField(_FIELD_Main_t_s64)
	m.fTS64 = field.GetInt64()
	protocache.Visit(m.visited[:], _FIELD_Main_t_s64)
	return m.fTS64
}

func (m *MainEX) SetTS64(v int64) {
	m.fTS64 = v
	protocache.Visit(m.visited[:], _FIELD_Main_t_s64)
}

func (m *MainEX) GetIndex() map[string]int32 {
	if protocache.CheckVisited(m.visited[:], _FIELD_Main_index) {
		return m.fIndex
	}
	field := m.source.GetField(_FIELD_Main_index)
	pack := field.GetMap()
	m.fIndex = make(map[string]int32, int(pack.Size()))
	for i := uint32(0); i < pack.Size(); i++ {
		keyField := pack.Key(i)
		key := keyField.GetString()
		valField := pack.Value(i)
		m.fIndex[key] = valField.GetInt32()
	}
	protocache.Visit(m.visited[:], _FIELD_Main_index)
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
	protocache.Visit(m.visited[:], _FIELD_Main_index)
}

func (m *MainEX) GetObjects() map[int32]*SmallEX {
	if protocache.CheckVisited(m.visited[:], _FIELD_Main_objects) {
		return m.fObjects
	}
	field := m.source.GetField(_FIELD_Main_objects)
	pack := field.GetMap()
	m.fObjects = make(map[int32]*SmallEX, int(pack.Size()))
	for i := uint32(0); i < pack.Size(); i++ {
		keyField := pack.Key(i)
		key := keyField.GetInt32()
		valField := pack.Value(i)
		m.fObjects[key] = TO_SmallEX(valField.GetObject())
	}
	protocache.Visit(m.visited[:], _FIELD_Main_objects)
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
	protocache.Visit(m.visited[:], _FIELD_Main_objects)
}

func (m *MainEX) GetMatrix() Vec2DEX {
	if protocache.CheckVisited(m.visited[:], _FIELD_Main_matrix) {
		return m.fMatrix
	}
	field := m.source.GetField(_FIELD_Main_matrix)
	m.fMatrix = TO_Vec2DEX(field.GetObject())
	protocache.Visit(m.visited[:], _FIELD_Main_matrix)
	return m.fMatrix
}

func (m *MainEX) SetMatrix(v Vec2DEX) {
	m.fMatrix = v
	protocache.Visit(m.visited[:], _FIELD_Main_matrix)
}

func (m *MainEX) GetVector() []ArrMapEX {
	if protocache.CheckVisited(m.visited[:], _FIELD_Main_vector) {
		return m.fVector
	}
	field := m.source.GetField(_FIELD_Main_vector)
	arr := field.GetArray()
	m.fVector = make([]ArrMapEX, int(arr.Size()))
	for i := uint32(0); i < arr.Size(); i++ {
		elem := arr.Get(i)
		m.fVector[i] = TO_ArrMapEX(elem.GetObject())
	}
	protocache.Visit(m.visited[:], _FIELD_Main_vector)
	return m.fVector
}

func (m *MainEX) SetVector(v []ArrMapEX) {
	if v == nil {
		m.fVector = nil
	} else {
		m.fVector = append(m.fVector[:0], v...)
	}
	protocache.Visit(m.visited[:], _FIELD_Main_vector)
}

func (m *MainEX) GetArrays() ArrMapEX {
	if protocache.CheckVisited(m.visited[:], _FIELD_Main_arrays) {
		return m.fArrays
	}
	field := m.source.GetField(_FIELD_Main_arrays)
	m.fArrays = TO_ArrMapEX(field.GetObject())
	protocache.Visit(m.visited[:], _FIELD_Main_arrays)
	return m.fArrays
}

func (m *MainEX) SetArrays(v ArrMapEX) {
	m.fArrays = v
	protocache.Visit(m.visited[:], _FIELD_Main_arrays)
}

func (m *MainEX) GetModev() []Mode {
	if protocache.CheckVisited(m.visited[:], _FIELD_Main_modev) {
		return m.fModev
	}
	field := m.source.GetField(_FIELD_Main_modev)
	m.fModev = append([]Mode(nil), protocache.CastEnumArray[Mode](field.GetEnumValueArray())...)
	protocache.Visit(m.visited[:], _FIELD_Main_modev)
	return m.fModev
}

func (m *MainEX) SetModev(v []Mode) {
	if v == nil {
		m.fModev = nil
	} else {
		m.fModev = append(m.fModev[:0], v...)
	}
	protocache.Visit(m.visited[:], _FIELD_Main_modev)
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
	field := msg.GetField(_FIELD_CyclicA_cyclic)
	if obj := field.DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, DETECT_CyclicB(obj))
	}
	return inlined
}

type CyclicAEX struct {
	source  protocache.Message
	visited [(_FIELD_TOTAL_CyclicA + 7) / 8]byte
	fValue  int32
	fCyclic *CyclicBEX
}

func TO_CyclicAEX(data []byte) *CyclicAEX {
	out := &CyclicAEX{}
	out.source = protocache.AsMessage(data)
	return out
}

func (m *CyclicAEX) HasSource() bool { return m.source.IsValid() }

func ENCODE_CyclicA(m *CyclicAEX) ([]uint32, error) {
	if m == nil {
		return []uint32{0}, nil
	}
	parts := make([][]uint32, _FIELD_TOTAL_CyclicA)
	if !protocache.CheckVisited(m.visited[:], _FIELD_CyclicA_value) {
		field := m.source.GetField(_FIELD_CyclicA_value)
		parts[_FIELD_CyclicA_value] = field.RawWords()
	} else if m.fValue != 0 {
		parts[_FIELD_CyclicA_value], _ = protocache.EncodeScalar[int32](m.fValue)
	}
	if !protocache.CheckVisited(m.visited[:], _FIELD_CyclicA_cyclic) {
		field := m.source.GetField(_FIELD_CyclicA_cyclic)
		if obj := field.DetectObject(); obj != nil {
			parts[_FIELD_CyclicA_cyclic] = protocache.BytesToWords(DETECT_CyclicB(obj))
		} else {
			parts[_FIELD_CyclicA_cyclic] = field.RawWords()
		}
	} else if m.fCyclic != nil {
		part, err := ENCODE_CyclicB(m.fCyclic)
		if err != nil {
			return nil, err
		}
		if len(part) > 1 {
			parts[_FIELD_CyclicA_cyclic] = part
		}
	}
	return protocache.EncodeMessageParts(parts)
}

func (m *CyclicAEX) Serialize() ([]byte, error) { return protocache.SerializeWords(ENCODE_CyclicA(m)) }

func (m *CyclicAEX) GetValue() int32 {
	if protocache.CheckVisited(m.visited[:], _FIELD_CyclicA_value) {
		return m.fValue
	}
	field := m.source.GetField(_FIELD_CyclicA_value)
	m.fValue = field.GetInt32()
	protocache.Visit(m.visited[:], _FIELD_CyclicA_value)
	return m.fValue
}

func (m *CyclicAEX) SetValue(v int32) {
	m.fValue = v
	protocache.Visit(m.visited[:], _FIELD_CyclicA_value)
}

func (m *CyclicAEX) GetCyclic() *CyclicBEX {
	if protocache.CheckVisited(m.visited[:], _FIELD_CyclicA_cyclic) {
		return m.fCyclic
	}
	field := m.source.GetField(_FIELD_CyclicA_cyclic)
	m.fCyclic = TO_CyclicBEX(field.GetObject())
	protocache.Visit(m.visited[:], _FIELD_CyclicA_cyclic)
	return m.fCyclic
}

func (m *CyclicAEX) SetCyclic(v *CyclicBEX) {
	m.fCyclic = v
	protocache.Visit(m.visited[:], _FIELD_CyclicA_cyclic)
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
	field := msg.GetField(_FIELD_CyclicB_cyclic)
	if obj := field.DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, DETECT_CyclicA(obj))
	}
	return inlined
}

type CyclicBEX struct {
	source  protocache.Message
	visited [(_FIELD_TOTAL_CyclicB + 7) / 8]byte
	fValue  int32
	fCyclic *CyclicAEX
}

func TO_CyclicBEX(data []byte) *CyclicBEX {
	out := &CyclicBEX{}
	out.source = protocache.AsMessage(data)
	return out
}

func (m *CyclicBEX) HasSource() bool { return m.source.IsValid() }

func ENCODE_CyclicB(m *CyclicBEX) ([]uint32, error) {
	if m == nil {
		return []uint32{0}, nil
	}
	parts := make([][]uint32, _FIELD_TOTAL_CyclicB)
	if !protocache.CheckVisited(m.visited[:], _FIELD_CyclicB_value) {
		field := m.source.GetField(_FIELD_CyclicB_value)
		parts[_FIELD_CyclicB_value] = field.RawWords()
	} else if m.fValue != 0 {
		parts[_FIELD_CyclicB_value], _ = protocache.EncodeScalar[int32](m.fValue)
	}
	if !protocache.CheckVisited(m.visited[:], _FIELD_CyclicB_cyclic) {
		field := m.source.GetField(_FIELD_CyclicB_cyclic)
		if obj := field.DetectObject(); obj != nil {
			parts[_FIELD_CyclicB_cyclic] = protocache.BytesToWords(DETECT_CyclicA(obj))
		} else {
			parts[_FIELD_CyclicB_cyclic] = field.RawWords()
		}
	} else if m.fCyclic != nil {
		part, err := ENCODE_CyclicA(m.fCyclic)
		if err != nil {
			return nil, err
		}
		if len(part) > 1 {
			parts[_FIELD_CyclicB_cyclic] = part
		}
	}
	return protocache.EncodeMessageParts(parts)
}

func (m *CyclicBEX) Serialize() ([]byte, error) { return protocache.SerializeWords(ENCODE_CyclicB(m)) }

func (m *CyclicBEX) GetValue() int32 {
	if protocache.CheckVisited(m.visited[:], _FIELD_CyclicB_value) {
		return m.fValue
	}
	field := m.source.GetField(_FIELD_CyclicB_value)
	m.fValue = field.GetInt32()
	protocache.Visit(m.visited[:], _FIELD_CyclicB_value)
	return m.fValue
}

func (m *CyclicBEX) SetValue(v int32) {
	m.fValue = v
	protocache.Visit(m.visited[:], _FIELD_CyclicB_value)
}

func (m *CyclicBEX) GetCyclic() *CyclicAEX {
	if protocache.CheckVisited(m.visited[:], _FIELD_CyclicB_cyclic) {
		return m.fCyclic
	}
	field := m.source.GetField(_FIELD_CyclicB_cyclic)
	m.fCyclic = TO_CyclicAEX(field.GetObject())
	protocache.Visit(m.visited[:], _FIELD_CyclicB_cyclic)
	return m.fCyclic
}

func (m *CyclicBEX) SetCyclic(v *CyclicAEX) {
	m.fCyclic = v
	protocache.Visit(m.visited[:], _FIELD_CyclicB_cyclic)
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
	return inlined
}

type Deprecated_ValidEX struct {
	source  protocache.Message
	visited [(_FIELD_TOTAL_Deprecated_Valid + 7) / 8]byte
	fVal    int32
}

func TO_Deprecated_ValidEX(data []byte) *Deprecated_ValidEX {
	out := &Deprecated_ValidEX{}
	out.source = protocache.AsMessage(data)
	return out
}

func (m *Deprecated_ValidEX) HasSource() bool { return m.source.IsValid() }

func ENCODE_Deprecated_Valid(m *Deprecated_ValidEX) ([]uint32, error) {
	if m == nil {
		return []uint32{0}, nil
	}
	parts := make([][]uint32, _FIELD_TOTAL_Deprecated_Valid)
	if !protocache.CheckVisited(m.visited[:], _FIELD_Deprecated_Valid_val) {
		field := m.source.GetField(_FIELD_Deprecated_Valid_val)
		parts[_FIELD_Deprecated_Valid_val] = field.RawWords()
	} else if m.fVal != 0 {
		parts[_FIELD_Deprecated_Valid_val], _ = protocache.EncodeScalar[int32](m.fVal)
	}
	return protocache.EncodeMessageParts(parts)
}

func (m *Deprecated_ValidEX) Serialize() ([]byte, error) {
	return protocache.SerializeWords(ENCODE_Deprecated_Valid(m))
}

func (m *Deprecated_ValidEX) GetVal() int32 {
	if protocache.CheckVisited(m.visited[:], _FIELD_Deprecated_Valid_val) {
		return m.fVal
	}
	field := m.source.GetField(_FIELD_Deprecated_Valid_val)
	m.fVal = field.GetInt32()
	protocache.Visit(m.visited[:], _FIELD_Deprecated_Valid_val)
	return m.fVal
}

func (m *Deprecated_ValidEX) SetVal(v int32) {
	m.fVal = v
	protocache.Visit(m.visited[:], _FIELD_Deprecated_Valid_val)
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
	return inlined
}

type DeprecatedEX struct {
	source  protocache.Message
	visited [(_FIELD_TOTAL_Deprecated + 7) / 8]byte
	fJunk   int32
}

func TO_DeprecatedEX(data []byte) *DeprecatedEX {
	out := &DeprecatedEX{}
	out.source = protocache.AsMessage(data)
	return out
}

func (m *DeprecatedEX) HasSource() bool { return m.source.IsValid() }

func ENCODE_Deprecated(m *DeprecatedEX) ([]uint32, error) {
	if m == nil {
		return []uint32{0}, nil
	}
	parts := make([][]uint32, _FIELD_TOTAL_Deprecated)
	if !protocache.CheckVisited(m.visited[:], _FIELD_Deprecated_junk) {
		field := m.source.GetField(_FIELD_Deprecated_junk)
		parts[_FIELD_Deprecated_junk] = field.RawWords()
	} else if m.fJunk != 0 {
		parts[_FIELD_Deprecated_junk], _ = protocache.EncodeScalar[int32](m.fJunk)
	}
	return protocache.EncodeMessageParts(parts)
}

func (m *DeprecatedEX) Serialize() ([]byte, error) {
	return protocache.SerializeWords(ENCODE_Deprecated(m))
}

func (m *DeprecatedEX) GetJunk() int32 {
	if protocache.CheckVisited(m.visited[:], _FIELD_Deprecated_junk) {
		return m.fJunk
	}
	field := m.source.GetField(_FIELD_Deprecated_junk)
	m.fJunk = field.GetInt32()
	protocache.Visit(m.visited[:], _FIELD_Deprecated_junk)
	return m.fJunk
}

func (m *DeprecatedEX) SetJunk(v int32) {
	m.fJunk = v
	protocache.Visit(m.visited[:], _FIELD_Deprecated_junk)
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

func ENCODE_ModeDict_Value(x ModeDict_ValueEX) ([]uint32, error) {
	if len(x) == 0 {
		return []uint32{1}, nil
	}
	return protocache.EncodeEnumArray([]Mode(x))
}

func (x ModeDict_ValueEX) Serialize() ([]byte, error) {
	return protocache.SerializeWords(ENCODE_ModeDict_Value(x))
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

func ENCODE_ModeDict(x ModeDictEX) ([]uint32, error) {
	return protocache.EncodeScalarMap(x,
		protocache.EncodeScalar[int32],
		ENCODE_ModeDict_Value)
}

func (x ModeDictEX) Serialize() ([]byte, error) { return protocache.SerializeWords(ENCODE_ModeDict(x)) }

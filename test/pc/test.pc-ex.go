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
	if obj := msg.GetField(_FIELD_Small_str).DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, protocache.DetectBytes(obj))
	}
	return inlined
}

type SmallEX struct {
	meta  protocache.MessageEX
	fI32  int32
	fFlag bool
	fStr  string
	fJunk int64
}

func TO_SmallEX(data []byte) *SmallEX {
	out := &SmallEX{}
	out.meta.Init(data)
	return out
}

func (m *SmallEX) HasBase() bool { return m.meta.HasBase() }

func (m *SmallEX) Encode() ([]uint32, error) {
	if m == nil {
		return []uint32{0}, nil
	}
	parts := make([][]uint32, 5)
	if !m.meta.IsVisited(_FIELD_Small_i32, _FIELD_TOTAL_Small) {
		field := m.meta.RawField(_FIELD_Small_i32)
		parts[0] = field.RawWords()
	} else if m.fI32 != 0 {
		parts[0] = protocache.EncodeInt32(m.fI32)
	}
	if !m.meta.IsVisited(_FIELD_Small_flag, _FIELD_TOTAL_Small) {
		field := m.meta.RawField(_FIELD_Small_flag)
		parts[1] = field.RawWords()
	} else if m.fFlag {
		parts[1] = protocache.EncodeBool(m.fFlag)
	}
	if !m.meta.IsVisited(_FIELD_Small_str, _FIELD_TOTAL_Small) {
		field := m.meta.RawField(_FIELD_Small_str)
		if obj := field.DetectObject(); obj != nil {
			parts[3] = protocache.BytesToWords(protocache.DetectBytes(obj))
		} else {
			parts[3] = field.RawWords()
		}
	} else if len(m.fStr) != 0 {
		part, err := protocache.EncodeString(m.fStr)
		if err != nil {
			return nil, err
		}
		parts[3] = part
	}
	if !m.meta.IsVisited(_FIELD_Small_junk, _FIELD_TOTAL_Small) {
		field := m.meta.RawField(_FIELD_Small_junk)
		parts[4] = field.RawWords()
	} else if m.fJunk != 0 {
		parts[4] = protocache.EncodeInt64(m.fJunk)
	}
	return protocache.EncodeMessageParts(parts)
}

func (m *SmallEX) Serialize() ([]byte, error) { return protocache.SerializeEncoded(m.Encode()) }

func (m *SmallEX) GetI32() int32 {
	if m.meta.IsVisited(_FIELD_Small_i32, _FIELD_TOTAL_Small) {
		return m.fI32
	}
	field := m.meta.RawField(_FIELD_Small_i32)
	m.fI32 = field.GetInt32()
	m.meta.Visit(_FIELD_Small_i32, _FIELD_TOTAL_Small)
	return m.fI32
}

func (m *SmallEX) SetI32(v int32) {
	m.fI32 = v
	m.meta.Visit(_FIELD_Small_i32, _FIELD_TOTAL_Small)
}

func (m *SmallEX) GetFlag() bool {
	if m.meta.IsVisited(_FIELD_Small_flag, _FIELD_TOTAL_Small) {
		return m.fFlag
	}
	field := m.meta.RawField(_FIELD_Small_flag)
	m.fFlag = field.GetBool()
	m.meta.Visit(_FIELD_Small_flag, _FIELD_TOTAL_Small)
	return m.fFlag
}

func (m *SmallEX) SetFlag(v bool) {
	m.fFlag = v
	m.meta.Visit(_FIELD_Small_flag, _FIELD_TOTAL_Small)
}

func (m *SmallEX) GetStr() string {
	if m.meta.IsVisited(_FIELD_Small_str, _FIELD_TOTAL_Small) {
		return m.fStr
	}
	field := m.meta.RawField(_FIELD_Small_str)
	m.fStr = field.GetString()
	m.meta.Visit(_FIELD_Small_str, _FIELD_TOTAL_Small)
	return m.fStr
}

func (m *SmallEX) SetStr(v string) {
	m.fStr = v
	m.meta.Visit(_FIELD_Small_str, _FIELD_TOTAL_Small)
}

func (m *SmallEX) GetJunk() int64 {
	if m.meta.IsVisited(_FIELD_Small_junk, _FIELD_TOTAL_Small) {
		return m.fJunk
	}
	field := m.meta.RawField(_FIELD_Small_junk)
	m.fJunk = field.GetInt64()
	m.meta.Visit(_FIELD_Small_junk, _FIELD_TOTAL_Small)
	return m.fJunk
}

func (m *SmallEX) SetJunk(v int64) {
	m.fJunk = v
	m.meta.Visit(_FIELD_Small_junk, _FIELD_TOTAL_Small)
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

func (x Vec2D_Vec1DEX) Encode() ([]uint32, error) {
	if len(x) == 0 {
		return []uint32{1}, nil
	}
	return protocache.EncodeFloat32Array([]float32(x))
}

func (x Vec2D_Vec1DEX) Serialize() ([]byte, error) { return protocache.SerializeEncoded(x.Encode()) }

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

func (x Vec2DEX) Encode() ([]uint32, error) {
	if len(x) == 0 {
		return []uint32{1}, nil
	}
	return protocache.EncodeObjectArray(len(x), func(i int) ([]uint32, error) { return x[i].Encode() })
}

func (x Vec2DEX) Serialize() ([]byte, error) { return protocache.SerializeEncoded(x.Encode()) }

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

func (x ArrMap_ArrayEX) Encode() ([]uint32, error) {
	if len(x) == 0 {
		return []uint32{1}, nil
	}
	return protocache.EncodeFloat32Array([]float32(x))
}

func (x ArrMap_ArrayEX) Serialize() ([]byte, error) { return protocache.SerializeEncoded(x.Encode()) }

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

func (x ArrMapEX) Encode() ([]uint32, error) {
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
		valPart, err := v.Encode()
		if err != nil {
			return nil, err
		}
		keys = append(keys, keyPart)
		vals = append(vals, valPart)
	}
	return protocache.EncodeMapParts(keys, vals, true)
}

func (x ArrMapEX) Serialize() ([]byte, error) { return protocache.SerializeEncoded(x.Encode()) }

func DETECT_Main(data []byte) []byte {
	msg := protocache.AsMessage(data)
	if !msg.IsValid() {
		return nil
	}
	inlined := msg.DetectInlined()
	if inlined == nil {
		return nil
	}
	if obj := msg.GetField(_FIELD_Main_modev).DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, protocache.DetectArray(obj, nil))
	}
	if obj := msg.GetField(_FIELD_Main_arrays).DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, DETECT_ArrMap(obj))
	}
	if obj := msg.GetField(_FIELD_Main_vector).DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, protocache.DetectArray(obj, DETECT_ArrMap))
	}
	if obj := msg.GetField(_FIELD_Main_matrix).DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, DETECT_Vec2D(obj))
	}
	if obj := msg.GetField(_FIELD_Main_objects).DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, protocache.DetectMap(obj, nil, DETECT_Small))
	}
	if obj := msg.GetField(_FIELD_Main_index).DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, protocache.DetectMap(obj, protocache.DetectBytes, nil))
	}
	if obj := msg.GetField(_FIELD_Main_objectv).DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, protocache.DetectArray(obj, DETECT_Small))
	}
	if obj := msg.GetField(_FIELD_Main_flags).DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, protocache.DetectBytes(obj))
	}
	if obj := msg.GetField(_FIELD_Main_f64v).DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, protocache.DetectArray(obj, nil))
	}
	if obj := msg.GetField(_FIELD_Main_f32v).DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, protocache.DetectArray(obj, nil))
	}
	if obj := msg.GetField(_FIELD_Main_datav).DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, protocache.DetectArray(obj, protocache.DetectBytes))
	}
	if obj := msg.GetField(_FIELD_Main_strv).DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, protocache.DetectArray(obj, protocache.DetectBytes))
	}
	if obj := msg.GetField(_FIELD_Main_u64v).DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, protocache.DetectArray(obj, nil))
	}
	if obj := msg.GetField(_FIELD_Main_i32v).DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, protocache.DetectArray(obj, nil))
	}
	if obj := msg.GetField(_FIELD_Main_object).DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, DETECT_Small(obj))
	}
	if obj := msg.GetField(_FIELD_Main_data).DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, protocache.DetectBytes(obj))
	}
	if obj := msg.GetField(_FIELD_Main_str).DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, protocache.DetectBytes(obj))
	}
	return inlined
}

type MainEX struct {
	meta     protocache.MessageEX
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
	out.meta.Init(data)
	return out
}

func (m *MainEX) HasBase() bool { return m.meta.HasBase() }

func (m *MainEX) Encode() ([]uint32, error) {
	if m == nil {
		return []uint32{0}, nil
	}
	parts := make([][]uint32, 32)
	if !m.meta.IsVisited(_FIELD_Main_i32, _FIELD_TOTAL_Main) {
		field := m.meta.RawField(_FIELD_Main_i32)
		parts[0] = field.RawWords()
	} else if m.fI32 != 0 {
		parts[0] = protocache.EncodeInt32(m.fI32)
	}
	if !m.meta.IsVisited(_FIELD_Main_u32, _FIELD_TOTAL_Main) {
		field := m.meta.RawField(_FIELD_Main_u32)
		parts[1] = field.RawWords()
	} else if m.fU32 != 0 {
		parts[1] = protocache.EncodeUint32(m.fU32)
	}
	if !m.meta.IsVisited(_FIELD_Main_i64, _FIELD_TOTAL_Main) {
		field := m.meta.RawField(_FIELD_Main_i64)
		parts[2] = field.RawWords()
	} else if m.fI64 != 0 {
		parts[2] = protocache.EncodeInt64(m.fI64)
	}
	if !m.meta.IsVisited(_FIELD_Main_u64, _FIELD_TOTAL_Main) {
		field := m.meta.RawField(_FIELD_Main_u64)
		parts[3] = field.RawWords()
	} else if m.fU64 != 0 {
		parts[3] = protocache.EncodeUint64(m.fU64)
	}
	if !m.meta.IsVisited(_FIELD_Main_flag, _FIELD_TOTAL_Main) {
		field := m.meta.RawField(_FIELD_Main_flag)
		parts[4] = field.RawWords()
	} else if m.fFlag {
		parts[4] = protocache.EncodeBool(m.fFlag)
	}
	if !m.meta.IsVisited(_FIELD_Main_mode, _FIELD_TOTAL_Main) {
		field := m.meta.RawField(_FIELD_Main_mode)
		parts[5] = field.RawWords()
	} else if m.fMode != 0 {
		parts[5] = protocache.EncodeInt32(int32(m.fMode))
	}
	if !m.meta.IsVisited(_FIELD_Main_str, _FIELD_TOTAL_Main) {
		field := m.meta.RawField(_FIELD_Main_str)
		if obj := field.DetectObject(); obj != nil {
			parts[6] = protocache.BytesToWords(protocache.DetectBytes(obj))
		} else {
			parts[6] = field.RawWords()
		}
	} else if len(m.fStr) != 0 {
		part, err := protocache.EncodeString(m.fStr)
		if err != nil {
			return nil, err
		}
		parts[6] = part
	}
	if !m.meta.IsVisited(_FIELD_Main_data, _FIELD_TOTAL_Main) {
		field := m.meta.RawField(_FIELD_Main_data)
		if obj := field.DetectObject(); obj != nil {
			parts[7] = protocache.BytesToWords(protocache.DetectBytes(obj))
		} else {
			parts[7] = field.RawWords()
		}
	} else if len(m.fData) != 0 {
		part, err := protocache.EncodeBytes(m.fData)
		if err != nil {
			return nil, err
		}
		parts[7] = part
	}
	if !m.meta.IsVisited(_FIELD_Main_f32, _FIELD_TOTAL_Main) {
		field := m.meta.RawField(_FIELD_Main_f32)
		parts[8] = field.RawWords()
	} else if m.fF32 != 0 {
		parts[8] = protocache.EncodeFloat32(m.fF32)
	}
	if !m.meta.IsVisited(_FIELD_Main_f64, _FIELD_TOTAL_Main) {
		field := m.meta.RawField(_FIELD_Main_f64)
		parts[9] = field.RawWords()
	} else if m.fF64 != 0 {
		parts[9] = protocache.EncodeFloat64(m.fF64)
	}
	if !m.meta.IsVisited(_FIELD_Main_object, _FIELD_TOTAL_Main) {
		field := m.meta.RawField(_FIELD_Main_object)
		if obj := field.DetectObject(); obj != nil {
			parts[10] = protocache.BytesToWords(DETECT_Small(obj))
		} else {
			parts[10] = field.RawWords()
		}
	} else if m.fObject != nil {
		part, err := m.fObject.Encode()
		if err != nil {
			return nil, err
		}
		if len(part) > 1 {
			parts[10] = part
		}
	}
	if !m.meta.IsVisited(_FIELD_Main_i32v, _FIELD_TOTAL_Main) {
		field := m.meta.RawField(_FIELD_Main_i32v)
		if obj := field.DetectObject(); obj != nil {
			parts[11] = protocache.BytesToWords(protocache.DetectArray(obj, nil))
		} else {
			parts[11] = field.RawWords()
		}
	} else if len(m.fI32V) != 0 {
		part, err := protocache.EncodeInt32Array(m.fI32V)
		if err != nil {
			return nil, err
		}
		parts[11] = part
	}
	if !m.meta.IsVisited(_FIELD_Main_u64v, _FIELD_TOTAL_Main) {
		field := m.meta.RawField(_FIELD_Main_u64v)
		if obj := field.DetectObject(); obj != nil {
			parts[12] = protocache.BytesToWords(protocache.DetectArray(obj, nil))
		} else {
			parts[12] = field.RawWords()
		}
	} else if len(m.fU64V) != 0 {
		part, err := protocache.EncodeUint64Array(m.fU64V)
		if err != nil {
			return nil, err
		}
		parts[12] = part
	}
	if !m.meta.IsVisited(_FIELD_Main_strv, _FIELD_TOTAL_Main) {
		field := m.meta.RawField(_FIELD_Main_strv)
		if obj := field.DetectObject(); obj != nil {
			parts[13] = protocache.BytesToWords(protocache.DetectArray(obj, protocache.DetectBytes))
		} else {
			parts[13] = field.RawWords()
		}
	} else if len(m.fStrv) != 0 {
		part, err := protocache.EncodeStringArray(m.fStrv)
		if err != nil {
			return nil, err
		}
		parts[13] = part
	}
	if !m.meta.IsVisited(_FIELD_Main_datav, _FIELD_TOTAL_Main) {
		field := m.meta.RawField(_FIELD_Main_datav)
		if obj := field.DetectObject(); obj != nil {
			parts[14] = protocache.BytesToWords(protocache.DetectArray(obj, protocache.DetectBytes))
		} else {
			parts[14] = field.RawWords()
		}
	} else if len(m.fDatav) != 0 {
		part, err := protocache.EncodeBytesArray(m.fDatav)
		if err != nil {
			return nil, err
		}
		parts[14] = part
	}
	if !m.meta.IsVisited(_FIELD_Main_f32v, _FIELD_TOTAL_Main) {
		field := m.meta.RawField(_FIELD_Main_f32v)
		if obj := field.DetectObject(); obj != nil {
			parts[15] = protocache.BytesToWords(protocache.DetectArray(obj, nil))
		} else {
			parts[15] = field.RawWords()
		}
	} else if len(m.fF32V) != 0 {
		part, err := protocache.EncodeFloat32Array(m.fF32V)
		if err != nil {
			return nil, err
		}
		parts[15] = part
	}
	if !m.meta.IsVisited(_FIELD_Main_f64v, _FIELD_TOTAL_Main) {
		field := m.meta.RawField(_FIELD_Main_f64v)
		if obj := field.DetectObject(); obj != nil {
			parts[16] = protocache.BytesToWords(protocache.DetectArray(obj, nil))
		} else {
			parts[16] = field.RawWords()
		}
	} else if len(m.fF64V) != 0 {
		part, err := protocache.EncodeFloat64Array(m.fF64V)
		if err != nil {
			return nil, err
		}
		parts[16] = part
	}
	if !m.meta.IsVisited(_FIELD_Main_flags, _FIELD_TOTAL_Main) {
		field := m.meta.RawField(_FIELD_Main_flags)
		if obj := field.DetectObject(); obj != nil {
			parts[17] = protocache.BytesToWords(protocache.DetectBytes(obj))
		} else {
			parts[17] = field.RawWords()
		}
	} else if len(m.fFlags) != 0 {
		part, err := protocache.EncodeBoolArray(m.fFlags)
		if err != nil {
			return nil, err
		}
		parts[17] = part
	}
	if !m.meta.IsVisited(_FIELD_Main_objectv, _FIELD_TOTAL_Main) {
		field := m.meta.RawField(_FIELD_Main_objectv)
		if obj := field.DetectObject(); obj != nil {
			parts[18] = protocache.BytesToWords(protocache.DetectArray(obj, DETECT_Small))
		} else {
			parts[18] = field.RawWords()
		}
	} else if len(m.fObjectv) != 0 {
		part, err := protocache.EncodeObjectArray(len(m.fObjectv), func(i int) ([]uint32, error) { return m.fObjectv[i].Encode() })
		if err != nil {
			return nil, err
		}
		parts[18] = part
	}
	if !m.meta.IsVisited(_FIELD_Main_t_u32, _FIELD_TOTAL_Main) {
		field := m.meta.RawField(_FIELD_Main_t_u32)
		parts[19] = field.RawWords()
	} else if m.fTU32 != 0 {
		parts[19] = protocache.EncodeUint32(m.fTU32)
	}
	if !m.meta.IsVisited(_FIELD_Main_t_i32, _FIELD_TOTAL_Main) {
		field := m.meta.RawField(_FIELD_Main_t_i32)
		parts[20] = field.RawWords()
	} else if m.fTI32 != 0 {
		parts[20] = protocache.EncodeInt32(m.fTI32)
	}
	if !m.meta.IsVisited(_FIELD_Main_t_s32, _FIELD_TOTAL_Main) {
		field := m.meta.RawField(_FIELD_Main_t_s32)
		parts[21] = field.RawWords()
	} else if m.fTS32 != 0 {
		parts[21] = protocache.EncodeInt32(m.fTS32)
	}
	if !m.meta.IsVisited(_FIELD_Main_t_u64, _FIELD_TOTAL_Main) {
		field := m.meta.RawField(_FIELD_Main_t_u64)
		parts[22] = field.RawWords()
	} else if m.fTU64 != 0 {
		parts[22] = protocache.EncodeUint64(m.fTU64)
	}
	if !m.meta.IsVisited(_FIELD_Main_t_i64, _FIELD_TOTAL_Main) {
		field := m.meta.RawField(_FIELD_Main_t_i64)
		parts[23] = field.RawWords()
	} else if m.fTI64 != 0 {
		parts[23] = protocache.EncodeInt64(m.fTI64)
	}
	if !m.meta.IsVisited(_FIELD_Main_t_s64, _FIELD_TOTAL_Main) {
		field := m.meta.RawField(_FIELD_Main_t_s64)
		parts[24] = field.RawWords()
	} else if m.fTS64 != 0 {
		parts[24] = protocache.EncodeInt64(m.fTS64)
	}
	if !m.meta.IsVisited(_FIELD_Main_index, _FIELD_TOTAL_Main) {
		field := m.meta.RawField(_FIELD_Main_index)
		if obj := field.DetectObject(); obj != nil {
			parts[25] = protocache.BytesToWords(protocache.DetectMap(obj, protocache.DetectBytes, nil))
		} else {
			parts[25] = field.RawWords()
		}
	} else if len(m.fIndex) != 0 {
		keys := make([][]uint32, 0, len(m.fIndex))
		vals := make([][]uint32, 0, len(m.fIndex))
		for k, v := range m.fIndex {
			keyPart, err := protocache.EncodeString(k)
			if err != nil {
				return nil, err
			}
			valPart := protocache.EncodeInt32(v)
			keys = append(keys, keyPart)
			vals = append(vals, valPart)
		}
		part, err := protocache.EncodeMapParts(keys, vals, true)
		if err != nil {
			return nil, err
		}
		parts[25] = part
	}
	if !m.meta.IsVisited(_FIELD_Main_objects, _FIELD_TOTAL_Main) {
		field := m.meta.RawField(_FIELD_Main_objects)
		if obj := field.DetectObject(); obj != nil {
			parts[26] = protocache.BytesToWords(protocache.DetectMap(obj, nil, DETECT_Small))
		} else {
			parts[26] = field.RawWords()
		}
	} else if len(m.fObjects) != 0 {
		keys := make([][]uint32, 0, len(m.fObjects))
		vals := make([][]uint32, 0, len(m.fObjects))
		for k, v := range m.fObjects {
			keyPart := protocache.EncodeInt32(k)
			valPart, err := v.Encode()
			if err != nil {
				return nil, err
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
	if !m.meta.IsVisited(_FIELD_Main_matrix, _FIELD_TOTAL_Main) {
		field := m.meta.RawField(_FIELD_Main_matrix)
		if obj := field.DetectObject(); obj != nil {
			parts[27] = protocache.BytesToWords(DETECT_Vec2D(obj))
		} else {
			parts[27] = field.RawWords()
		}
	} else if m.fMatrix != nil {
		part, err := m.fMatrix.Encode()
		if err != nil {
			return nil, err
		}
		if len(part) > 1 {
			parts[27] = part
		}
	}
	if !m.meta.IsVisited(_FIELD_Main_vector, _FIELD_TOTAL_Main) {
		field := m.meta.RawField(_FIELD_Main_vector)
		if obj := field.DetectObject(); obj != nil {
			parts[28] = protocache.BytesToWords(protocache.DetectArray(obj, DETECT_ArrMap))
		} else {
			parts[28] = field.RawWords()
		}
	} else if len(m.fVector) != 0 {
		part, err := protocache.EncodeObjectArray(len(m.fVector), func(i int) ([]uint32, error) { return m.fVector[i].Encode() })
		if err != nil {
			return nil, err
		}
		parts[28] = part
	}
	if !m.meta.IsVisited(_FIELD_Main_arrays, _FIELD_TOTAL_Main) {
		field := m.meta.RawField(_FIELD_Main_arrays)
		if obj := field.DetectObject(); obj != nil {
			parts[29] = protocache.BytesToWords(DETECT_ArrMap(obj))
		} else {
			parts[29] = field.RawWords()
		}
	} else if m.fArrays != nil {
		part, err := m.fArrays.Encode()
		if err != nil {
			return nil, err
		}
		if len(part) > 1 {
			parts[29] = part
		}
	}
	if !m.meta.IsVisited(_FIELD_Main_modev, _FIELD_TOTAL_Main) {
		field := m.meta.RawField(_FIELD_Main_modev)
		if obj := field.DetectObject(); obj != nil {
			parts[31] = protocache.BytesToWords(protocache.DetectArray(obj, nil))
		} else {
			parts[31] = field.RawWords()
		}
	} else if len(m.fModev) != 0 {
		part, err := protocache.EncodeEnumArray(m.fModev)
		if err != nil {
			return nil, err
		}
		parts[31] = part
	}
	return protocache.EncodeMessageParts(parts)
}

func (m *MainEX) Serialize() ([]byte, error) { return protocache.SerializeEncoded(m.Encode()) }

func (m *MainEX) GetI32() int32 {
	if m.meta.IsVisited(_FIELD_Main_i32, _FIELD_TOTAL_Main) {
		return m.fI32
	}
	field := m.meta.RawField(_FIELD_Main_i32)
	m.fI32 = field.GetInt32()
	m.meta.Visit(_FIELD_Main_i32, _FIELD_TOTAL_Main)
	return m.fI32
}

func (m *MainEX) SetI32(v int32) {
	m.fI32 = v
	m.meta.Visit(_FIELD_Main_i32, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetU32() uint32 {
	if m.meta.IsVisited(_FIELD_Main_u32, _FIELD_TOTAL_Main) {
		return m.fU32
	}
	field := m.meta.RawField(_FIELD_Main_u32)
	m.fU32 = field.GetUint32()
	m.meta.Visit(_FIELD_Main_u32, _FIELD_TOTAL_Main)
	return m.fU32
}

func (m *MainEX) SetU32(v uint32) {
	m.fU32 = v
	m.meta.Visit(_FIELD_Main_u32, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetI64() int64 {
	if m.meta.IsVisited(_FIELD_Main_i64, _FIELD_TOTAL_Main) {
		return m.fI64
	}
	field := m.meta.RawField(_FIELD_Main_i64)
	m.fI64 = field.GetInt64()
	m.meta.Visit(_FIELD_Main_i64, _FIELD_TOTAL_Main)
	return m.fI64
}

func (m *MainEX) SetI64(v int64) {
	m.fI64 = v
	m.meta.Visit(_FIELD_Main_i64, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetU64() uint64 {
	if m.meta.IsVisited(_FIELD_Main_u64, _FIELD_TOTAL_Main) {
		return m.fU64
	}
	field := m.meta.RawField(_FIELD_Main_u64)
	m.fU64 = field.GetUint64()
	m.meta.Visit(_FIELD_Main_u64, _FIELD_TOTAL_Main)
	return m.fU64
}

func (m *MainEX) SetU64(v uint64) {
	m.fU64 = v
	m.meta.Visit(_FIELD_Main_u64, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetFlag() bool {
	if m.meta.IsVisited(_FIELD_Main_flag, _FIELD_TOTAL_Main) {
		return m.fFlag
	}
	field := m.meta.RawField(_FIELD_Main_flag)
	m.fFlag = field.GetBool()
	m.meta.Visit(_FIELD_Main_flag, _FIELD_TOTAL_Main)
	return m.fFlag
}

func (m *MainEX) SetFlag(v bool) {
	m.fFlag = v
	m.meta.Visit(_FIELD_Main_flag, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetMode() Mode {
	if m.meta.IsVisited(_FIELD_Main_mode, _FIELD_TOTAL_Main) {
		return m.fMode
	}
	field := m.meta.RawField(_FIELD_Main_mode)
	m.fMode = Mode(field.GetEnumValue())
	m.meta.Visit(_FIELD_Main_mode, _FIELD_TOTAL_Main)
	return m.fMode
}

func (m *MainEX) SetMode(v Mode) {
	m.fMode = v
	m.meta.Visit(_FIELD_Main_mode, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetStr() string {
	if m.meta.IsVisited(_FIELD_Main_str, _FIELD_TOTAL_Main) {
		return m.fStr
	}
	field := m.meta.RawField(_FIELD_Main_str)
	m.fStr = field.GetString()
	m.meta.Visit(_FIELD_Main_str, _FIELD_TOTAL_Main)
	return m.fStr
}

func (m *MainEX) SetStr(v string) {
	m.fStr = v
	m.meta.Visit(_FIELD_Main_str, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetData() []byte {
	if m.meta.IsVisited(_FIELD_Main_data, _FIELD_TOTAL_Main) {
		return m.fData
	}
	field := m.meta.RawField(_FIELD_Main_data)
	if data := field.GetBytes(); data != nil {
		m.fData = append([]byte(nil), data...)
	}
	m.meta.Visit(_FIELD_Main_data, _FIELD_TOTAL_Main)
	return m.fData
}

func (m *MainEX) SetData(v []byte) {
	if v == nil {
		m.fData = nil
	} else {
		m.fData = append([]byte(nil), v...)
	}
	m.meta.Visit(_FIELD_Main_data, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetF32() float32 {
	if m.meta.IsVisited(_FIELD_Main_f32, _FIELD_TOTAL_Main) {
		return m.fF32
	}
	field := m.meta.RawField(_FIELD_Main_f32)
	m.fF32 = field.GetFloat32()
	m.meta.Visit(_FIELD_Main_f32, _FIELD_TOTAL_Main)
	return m.fF32
}

func (m *MainEX) SetF32(v float32) {
	m.fF32 = v
	m.meta.Visit(_FIELD_Main_f32, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetF64() float64 {
	if m.meta.IsVisited(_FIELD_Main_f64, _FIELD_TOTAL_Main) {
		return m.fF64
	}
	field := m.meta.RawField(_FIELD_Main_f64)
	m.fF64 = field.GetFloat64()
	m.meta.Visit(_FIELD_Main_f64, _FIELD_TOTAL_Main)
	return m.fF64
}

func (m *MainEX) SetF64(v float64) {
	m.fF64 = v
	m.meta.Visit(_FIELD_Main_f64, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetObject() *SmallEX {
	if m.meta.IsVisited(_FIELD_Main_object, _FIELD_TOTAL_Main) {
		return m.fObject
	}
	field := m.meta.RawField(_FIELD_Main_object)
	m.fObject = TO_SmallEX(field.GetObject())
	m.meta.Visit(_FIELD_Main_object, _FIELD_TOTAL_Main)
	return m.fObject
}

func (m *MainEX) SetObject(v *SmallEX) {
	m.fObject = v
	m.meta.Visit(_FIELD_Main_object, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetI32V() []int32 {
	if m.meta.IsVisited(_FIELD_Main_i32v, _FIELD_TOTAL_Main) {
		return m.fI32V
	}
	field := m.meta.RawField(_FIELD_Main_i32v)
	m.fI32V = append([]int32(nil), field.GetInt32Array()...)
	m.meta.Visit(_FIELD_Main_i32v, _FIELD_TOTAL_Main)
	return m.fI32V
}

func (m *MainEX) SetI32V(v []int32) {
	if v == nil {
		m.fI32V = nil
	} else {
		m.fI32V = append(m.fI32V[:0], v...)
	}
	m.meta.Visit(_FIELD_Main_i32v, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetU64V() []uint64 {
	if m.meta.IsVisited(_FIELD_Main_u64v, _FIELD_TOTAL_Main) {
		return m.fU64V
	}
	field := m.meta.RawField(_FIELD_Main_u64v)
	m.fU64V = append([]uint64(nil), field.GetUint64Array()...)
	m.meta.Visit(_FIELD_Main_u64v, _FIELD_TOTAL_Main)
	return m.fU64V
}

func (m *MainEX) SetU64V(v []uint64) {
	if v == nil {
		m.fU64V = nil
	} else {
		m.fU64V = append(m.fU64V[:0], v...)
	}
	m.meta.Visit(_FIELD_Main_u64v, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetStrv() []string {
	if m.meta.IsVisited(_FIELD_Main_strv, _FIELD_TOTAL_Main) {
		return m.fStrv
	}
	field := m.meta.RawField(_FIELD_Main_strv)
	arr := protocache.AsStringArray(field.GetObject())
	m.fStrv = make([]string, int(arr.Size()))
	for i := uint32(0); i < arr.Size(); i++ {
		m.fStrv[i] = arr.Get(i)
	}
	m.meta.Visit(_FIELD_Main_strv, _FIELD_TOTAL_Main)
	return m.fStrv
}

func (m *MainEX) SetStrv(v []string) {
	if v == nil {
		m.fStrv = nil
	} else {
		m.fStrv = append(m.fStrv[:0], v...)
	}
	m.meta.Visit(_FIELD_Main_strv, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetDatav() [][]byte {
	if m.meta.IsVisited(_FIELD_Main_datav, _FIELD_TOTAL_Main) {
		return m.fDatav
	}
	field := m.meta.RawField(_FIELD_Main_datav)
	arr := protocache.AsBytesArray(field.GetObject())
	m.fDatav = make([][]byte, int(arr.Size()))
	for i := uint32(0); i < arr.Size(); i++ {
		if data := arr.Get(i); data != nil {
			m.fDatav[i] = append([]byte(nil), data...)
		}
	}
	m.meta.Visit(_FIELD_Main_datav, _FIELD_TOTAL_Main)
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
	m.meta.Visit(_FIELD_Main_datav, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetF32V() []float32 {
	if m.meta.IsVisited(_FIELD_Main_f32v, _FIELD_TOTAL_Main) {
		return m.fF32V
	}
	field := m.meta.RawField(_FIELD_Main_f32v)
	m.fF32V = append([]float32(nil), field.GetFloat32Array()...)
	m.meta.Visit(_FIELD_Main_f32v, _FIELD_TOTAL_Main)
	return m.fF32V
}

func (m *MainEX) SetF32V(v []float32) {
	if v == nil {
		m.fF32V = nil
	} else {
		m.fF32V = append(m.fF32V[:0], v...)
	}
	m.meta.Visit(_FIELD_Main_f32v, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetF64V() []float64 {
	if m.meta.IsVisited(_FIELD_Main_f64v, _FIELD_TOTAL_Main) {
		return m.fF64V
	}
	field := m.meta.RawField(_FIELD_Main_f64v)
	m.fF64V = append([]float64(nil), field.GetFloat64Array()...)
	m.meta.Visit(_FIELD_Main_f64v, _FIELD_TOTAL_Main)
	return m.fF64V
}

func (m *MainEX) SetF64V(v []float64) {
	if v == nil {
		m.fF64V = nil
	} else {
		m.fF64V = append(m.fF64V[:0], v...)
	}
	m.meta.Visit(_FIELD_Main_f64v, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetFlags() []bool {
	if m.meta.IsVisited(_FIELD_Main_flags, _FIELD_TOTAL_Main) {
		return m.fFlags
	}
	field := m.meta.RawField(_FIELD_Main_flags)
	m.fFlags = append([]bool(nil), field.GetBoolArray()...)
	m.meta.Visit(_FIELD_Main_flags, _FIELD_TOTAL_Main)
	return m.fFlags
}

func (m *MainEX) SetFlags(v []bool) {
	if v == nil {
		m.fFlags = nil
	} else {
		m.fFlags = append(m.fFlags[:0], v...)
	}
	m.meta.Visit(_FIELD_Main_flags, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetObjectv() []*SmallEX {
	if m.meta.IsVisited(_FIELD_Main_objectv, _FIELD_TOTAL_Main) {
		return m.fObjectv
	}
	field := m.meta.RawField(_FIELD_Main_objectv)
	arr := field.GetArray()
	m.fObjectv = make([]*SmallEX, int(arr.Size()))
	for i := uint32(0); i < arr.Size(); i++ {
		elem := arr.Get(i)
		m.fObjectv[i] = TO_SmallEX(elem.GetObject())
	}
	m.meta.Visit(_FIELD_Main_objectv, _FIELD_TOTAL_Main)
	return m.fObjectv
}

func (m *MainEX) SetObjectv(v []*SmallEX) {
	if v == nil {
		m.fObjectv = nil
	} else {
		m.fObjectv = append(m.fObjectv[:0], v...)
	}
	m.meta.Visit(_FIELD_Main_objectv, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetTU32() uint32 {
	if m.meta.IsVisited(_FIELD_Main_t_u32, _FIELD_TOTAL_Main) {
		return m.fTU32
	}
	field := m.meta.RawField(_FIELD_Main_t_u32)
	m.fTU32 = field.GetUint32()
	m.meta.Visit(_FIELD_Main_t_u32, _FIELD_TOTAL_Main)
	return m.fTU32
}

func (m *MainEX) SetTU32(v uint32) {
	m.fTU32 = v
	m.meta.Visit(_FIELD_Main_t_u32, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetTI32() int32 {
	if m.meta.IsVisited(_FIELD_Main_t_i32, _FIELD_TOTAL_Main) {
		return m.fTI32
	}
	field := m.meta.RawField(_FIELD_Main_t_i32)
	m.fTI32 = field.GetInt32()
	m.meta.Visit(_FIELD_Main_t_i32, _FIELD_TOTAL_Main)
	return m.fTI32
}

func (m *MainEX) SetTI32(v int32) {
	m.fTI32 = v
	m.meta.Visit(_FIELD_Main_t_i32, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetTS32() int32 {
	if m.meta.IsVisited(_FIELD_Main_t_s32, _FIELD_TOTAL_Main) {
		return m.fTS32
	}
	field := m.meta.RawField(_FIELD_Main_t_s32)
	m.fTS32 = field.GetInt32()
	m.meta.Visit(_FIELD_Main_t_s32, _FIELD_TOTAL_Main)
	return m.fTS32
}

func (m *MainEX) SetTS32(v int32) {
	m.fTS32 = v
	m.meta.Visit(_FIELD_Main_t_s32, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetTU64() uint64 {
	if m.meta.IsVisited(_FIELD_Main_t_u64, _FIELD_TOTAL_Main) {
		return m.fTU64
	}
	field := m.meta.RawField(_FIELD_Main_t_u64)
	m.fTU64 = field.GetUint64()
	m.meta.Visit(_FIELD_Main_t_u64, _FIELD_TOTAL_Main)
	return m.fTU64
}

func (m *MainEX) SetTU64(v uint64) {
	m.fTU64 = v
	m.meta.Visit(_FIELD_Main_t_u64, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetTI64() int64 {
	if m.meta.IsVisited(_FIELD_Main_t_i64, _FIELD_TOTAL_Main) {
		return m.fTI64
	}
	field := m.meta.RawField(_FIELD_Main_t_i64)
	m.fTI64 = field.GetInt64()
	m.meta.Visit(_FIELD_Main_t_i64, _FIELD_TOTAL_Main)
	return m.fTI64
}

func (m *MainEX) SetTI64(v int64) {
	m.fTI64 = v
	m.meta.Visit(_FIELD_Main_t_i64, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetTS64() int64 {
	if m.meta.IsVisited(_FIELD_Main_t_s64, _FIELD_TOTAL_Main) {
		return m.fTS64
	}
	field := m.meta.RawField(_FIELD_Main_t_s64)
	m.fTS64 = field.GetInt64()
	m.meta.Visit(_FIELD_Main_t_s64, _FIELD_TOTAL_Main)
	return m.fTS64
}

func (m *MainEX) SetTS64(v int64) {
	m.fTS64 = v
	m.meta.Visit(_FIELD_Main_t_s64, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetIndex() map[string]int32 {
	if m.meta.IsVisited(_FIELD_Main_index, _FIELD_TOTAL_Main) {
		return m.fIndex
	}
	field := m.meta.RawField(_FIELD_Main_index)
	pack := field.GetMap()
	m.fIndex = make(map[string]int32, int(pack.Size()))
	for i := uint32(0); i < pack.Size(); i++ {
		keyField := pack.Key(i)
		key := keyField.GetString()
		valField := pack.Value(i)
		m.fIndex[key] = valField.GetInt32()
	}
	m.meta.Visit(_FIELD_Main_index, _FIELD_TOTAL_Main)
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
	m.meta.Visit(_FIELD_Main_index, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetObjects() map[int32]*SmallEX {
	if m.meta.IsVisited(_FIELD_Main_objects, _FIELD_TOTAL_Main) {
		return m.fObjects
	}
	field := m.meta.RawField(_FIELD_Main_objects)
	pack := field.GetMap()
	m.fObjects = make(map[int32]*SmallEX, int(pack.Size()))
	for i := uint32(0); i < pack.Size(); i++ {
		keyField := pack.Key(i)
		key := keyField.GetInt32()
		valField := pack.Value(i)
		m.fObjects[key] = TO_SmallEX(valField.GetObject())
	}
	m.meta.Visit(_FIELD_Main_objects, _FIELD_TOTAL_Main)
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
	m.meta.Visit(_FIELD_Main_objects, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetMatrix() Vec2DEX {
	if m.meta.IsVisited(_FIELD_Main_matrix, _FIELD_TOTAL_Main) {
		return m.fMatrix
	}
	field := m.meta.RawField(_FIELD_Main_matrix)
	m.fMatrix = TO_Vec2DEX(field.GetObject())
	m.meta.Visit(_FIELD_Main_matrix, _FIELD_TOTAL_Main)
	return m.fMatrix
}

func (m *MainEX) SetMatrix(v Vec2DEX) {
	m.fMatrix = v
	m.meta.Visit(_FIELD_Main_matrix, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetVector() []ArrMapEX {
	if m.meta.IsVisited(_FIELD_Main_vector, _FIELD_TOTAL_Main) {
		return m.fVector
	}
	field := m.meta.RawField(_FIELD_Main_vector)
	arr := field.GetArray()
	m.fVector = make([]ArrMapEX, int(arr.Size()))
	for i := uint32(0); i < arr.Size(); i++ {
		elem := arr.Get(i)
		m.fVector[i] = TO_ArrMapEX(elem.GetObject())
	}
	m.meta.Visit(_FIELD_Main_vector, _FIELD_TOTAL_Main)
	return m.fVector
}

func (m *MainEX) SetVector(v []ArrMapEX) {
	if v == nil {
		m.fVector = nil
	} else {
		m.fVector = append(m.fVector[:0], v...)
	}
	m.meta.Visit(_FIELD_Main_vector, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetArrays() ArrMapEX {
	if m.meta.IsVisited(_FIELD_Main_arrays, _FIELD_TOTAL_Main) {
		return m.fArrays
	}
	field := m.meta.RawField(_FIELD_Main_arrays)
	m.fArrays = TO_ArrMapEX(field.GetObject())
	m.meta.Visit(_FIELD_Main_arrays, _FIELD_TOTAL_Main)
	return m.fArrays
}

func (m *MainEX) SetArrays(v ArrMapEX) {
	m.fArrays = v
	m.meta.Visit(_FIELD_Main_arrays, _FIELD_TOTAL_Main)
}

func (m *MainEX) GetModev() []Mode {
	if m.meta.IsVisited(_FIELD_Main_modev, _FIELD_TOTAL_Main) {
		return m.fModev
	}
	field := m.meta.RawField(_FIELD_Main_modev)
	m.fModev = append([]Mode(nil), protocache.CastEnumArray[Mode](field.GetEnumValueArray())...)
	m.meta.Visit(_FIELD_Main_modev, _FIELD_TOTAL_Main)
	return m.fModev
}

func (m *MainEX) SetModev(v []Mode) {
	if v == nil {
		m.fModev = nil
	} else {
		m.fModev = append(m.fModev[:0], v...)
	}
	m.meta.Visit(_FIELD_Main_modev, _FIELD_TOTAL_Main)
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
	if obj := msg.GetField(_FIELD_CyclicA_cyclic).DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, DETECT_CyclicB(obj))
	}
	return inlined
}

type CyclicAEX struct {
	meta    protocache.MessageEX
	fValue  int32
	fCyclic *CyclicBEX
}

func TO_CyclicAEX(data []byte) *CyclicAEX {
	out := &CyclicAEX{}
	out.meta.Init(data)
	return out
}

func (m *CyclicAEX) HasBase() bool { return m.meta.HasBase() }

func (m *CyclicAEX) Encode() ([]uint32, error) {
	if m == nil {
		return []uint32{0}, nil
	}
	parts := make([][]uint32, 2)
	if !m.meta.IsVisited(_FIELD_CyclicA_value, _FIELD_TOTAL_CyclicA) {
		field := m.meta.RawField(_FIELD_CyclicA_value)
		parts[0] = field.RawWords()
	} else if m.fValue != 0 {
		parts[0] = protocache.EncodeInt32(m.fValue)
	}
	if !m.meta.IsVisited(_FIELD_CyclicA_cyclic, _FIELD_TOTAL_CyclicA) {
		field := m.meta.RawField(_FIELD_CyclicA_cyclic)
		if obj := field.DetectObject(); obj != nil {
			parts[1] = protocache.BytesToWords(DETECT_CyclicB(obj))
		} else {
			parts[1] = field.RawWords()
		}
	} else if m.fCyclic != nil {
		part, err := m.fCyclic.Encode()
		if err != nil {
			return nil, err
		}
		if len(part) > 1 {
			parts[1] = part
		}
	}
	return protocache.EncodeMessageParts(parts)
}

func (m *CyclicAEX) Serialize() ([]byte, error) { return protocache.SerializeEncoded(m.Encode()) }

func (m *CyclicAEX) GetValue() int32 {
	if m.meta.IsVisited(_FIELD_CyclicA_value, _FIELD_TOTAL_CyclicA) {
		return m.fValue
	}
	field := m.meta.RawField(_FIELD_CyclicA_value)
	m.fValue = field.GetInt32()
	m.meta.Visit(_FIELD_CyclicA_value, _FIELD_TOTAL_CyclicA)
	return m.fValue
}

func (m *CyclicAEX) SetValue(v int32) {
	m.fValue = v
	m.meta.Visit(_FIELD_CyclicA_value, _FIELD_TOTAL_CyclicA)
}

func (m *CyclicAEX) GetCyclic() *CyclicBEX {
	if m.meta.IsVisited(_FIELD_CyclicA_cyclic, _FIELD_TOTAL_CyclicA) {
		return m.fCyclic
	}
	field := m.meta.RawField(_FIELD_CyclicA_cyclic)
	m.fCyclic = TO_CyclicBEX(field.GetObject())
	m.meta.Visit(_FIELD_CyclicA_cyclic, _FIELD_TOTAL_CyclicA)
	return m.fCyclic
}

func (m *CyclicAEX) SetCyclic(v *CyclicBEX) {
	m.fCyclic = v
	m.meta.Visit(_FIELD_CyclicA_cyclic, _FIELD_TOTAL_CyclicA)
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
	if obj := msg.GetField(_FIELD_CyclicB_cyclic).DetectObject(); obj != nil {
		return protocache.DetectShrink(data, obj, DETECT_CyclicA(obj))
	}
	return inlined
}

type CyclicBEX struct {
	meta    protocache.MessageEX
	fValue  int32
	fCyclic *CyclicAEX
}

func TO_CyclicBEX(data []byte) *CyclicBEX {
	out := &CyclicBEX{}
	out.meta.Init(data)
	return out
}

func (m *CyclicBEX) HasBase() bool { return m.meta.HasBase() }

func (m *CyclicBEX) Encode() ([]uint32, error) {
	if m == nil {
		return []uint32{0}, nil
	}
	parts := make([][]uint32, 2)
	if !m.meta.IsVisited(_FIELD_CyclicB_value, _FIELD_TOTAL_CyclicB) {
		field := m.meta.RawField(_FIELD_CyclicB_value)
		parts[0] = field.RawWords()
	} else if m.fValue != 0 {
		parts[0] = protocache.EncodeInt32(m.fValue)
	}
	if !m.meta.IsVisited(_FIELD_CyclicB_cyclic, _FIELD_TOTAL_CyclicB) {
		field := m.meta.RawField(_FIELD_CyclicB_cyclic)
		if obj := field.DetectObject(); obj != nil {
			parts[1] = protocache.BytesToWords(DETECT_CyclicA(obj))
		} else {
			parts[1] = field.RawWords()
		}
	} else if m.fCyclic != nil {
		part, err := m.fCyclic.Encode()
		if err != nil {
			return nil, err
		}
		if len(part) > 1 {
			parts[1] = part
		}
	}
	return protocache.EncodeMessageParts(parts)
}

func (m *CyclicBEX) Serialize() ([]byte, error) { return protocache.SerializeEncoded(m.Encode()) }

func (m *CyclicBEX) GetValue() int32 {
	if m.meta.IsVisited(_FIELD_CyclicB_value, _FIELD_TOTAL_CyclicB) {
		return m.fValue
	}
	field := m.meta.RawField(_FIELD_CyclicB_value)
	m.fValue = field.GetInt32()
	m.meta.Visit(_FIELD_CyclicB_value, _FIELD_TOTAL_CyclicB)
	return m.fValue
}

func (m *CyclicBEX) SetValue(v int32) {
	m.fValue = v
	m.meta.Visit(_FIELD_CyclicB_value, _FIELD_TOTAL_CyclicB)
}

func (m *CyclicBEX) GetCyclic() *CyclicAEX {
	if m.meta.IsVisited(_FIELD_CyclicB_cyclic, _FIELD_TOTAL_CyclicB) {
		return m.fCyclic
	}
	field := m.meta.RawField(_FIELD_CyclicB_cyclic)
	m.fCyclic = TO_CyclicAEX(field.GetObject())
	m.meta.Visit(_FIELD_CyclicB_cyclic, _FIELD_TOTAL_CyclicB)
	return m.fCyclic
}

func (m *CyclicBEX) SetCyclic(v *CyclicAEX) {
	m.fCyclic = v
	m.meta.Visit(_FIELD_CyclicB_cyclic, _FIELD_TOTAL_CyclicB)
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
	meta protocache.MessageEX
	fVal int32
}

func TO_Deprecated_ValidEX(data []byte) *Deprecated_ValidEX {
	out := &Deprecated_ValidEX{}
	out.meta.Init(data)
	return out
}

func (m *Deprecated_ValidEX) HasBase() bool { return m.meta.HasBase() }

func (m *Deprecated_ValidEX) Encode() ([]uint32, error) {
	if m == nil {
		return []uint32{0}, nil
	}
	parts := make([][]uint32, 1)
	if !m.meta.IsVisited(_FIELD_Deprecated_Valid_val, _FIELD_TOTAL_Deprecated_Valid) {
		field := m.meta.RawField(_FIELD_Deprecated_Valid_val)
		parts[0] = field.RawWords()
	} else if m.fVal != 0 {
		parts[0] = protocache.EncodeInt32(m.fVal)
	}
	return protocache.EncodeMessageParts(parts)
}

func (m *Deprecated_ValidEX) Serialize() ([]byte, error) {
	return protocache.SerializeEncoded(m.Encode())
}

func (m *Deprecated_ValidEX) GetVal() int32 {
	if m.meta.IsVisited(_FIELD_Deprecated_Valid_val, _FIELD_TOTAL_Deprecated_Valid) {
		return m.fVal
	}
	field := m.meta.RawField(_FIELD_Deprecated_Valid_val)
	m.fVal = field.GetInt32()
	m.meta.Visit(_FIELD_Deprecated_Valid_val, _FIELD_TOTAL_Deprecated_Valid)
	return m.fVal
}

func (m *Deprecated_ValidEX) SetVal(v int32) {
	m.fVal = v
	m.meta.Visit(_FIELD_Deprecated_Valid_val, _FIELD_TOTAL_Deprecated_Valid)
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
	meta  protocache.MessageEX
	fJunk int32
}

func TO_DeprecatedEX(data []byte) *DeprecatedEX {
	out := &DeprecatedEX{}
	out.meta.Init(data)
	return out
}

func (m *DeprecatedEX) HasBase() bool { return m.meta.HasBase() }

func (m *DeprecatedEX) Encode() ([]uint32, error) {
	if m == nil {
		return []uint32{0}, nil
	}
	parts := make([][]uint32, 1)
	if !m.meta.IsVisited(_FIELD_Deprecated_junk, _FIELD_TOTAL_Deprecated) {
		field := m.meta.RawField(_FIELD_Deprecated_junk)
		parts[0] = field.RawWords()
	} else if m.fJunk != 0 {
		parts[0] = protocache.EncodeInt32(m.fJunk)
	}
	return protocache.EncodeMessageParts(parts)
}

func (m *DeprecatedEX) Serialize() ([]byte, error) { return protocache.SerializeEncoded(m.Encode()) }

func (m *DeprecatedEX) GetJunk() int32 {
	if m.meta.IsVisited(_FIELD_Deprecated_junk, _FIELD_TOTAL_Deprecated) {
		return m.fJunk
	}
	field := m.meta.RawField(_FIELD_Deprecated_junk)
	m.fJunk = field.GetInt32()
	m.meta.Visit(_FIELD_Deprecated_junk, _FIELD_TOTAL_Deprecated)
	return m.fJunk
}

func (m *DeprecatedEX) SetJunk(v int32) {
	m.fJunk = v
	m.meta.Visit(_FIELD_Deprecated_junk, _FIELD_TOTAL_Deprecated)
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

func (x ModeDict_ValueEX) Encode() ([]uint32, error) {
	if len(x) == 0 {
		return []uint32{1}, nil
	}
	return protocache.EncodeEnumArray([]Mode(x))
}

func (x ModeDict_ValueEX) Serialize() ([]byte, error) { return protocache.SerializeEncoded(x.Encode()) }

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

func (x ModeDictEX) Encode() ([]uint32, error) {
	if len(x) == 0 {
		return []uint32{5 << 28}, nil
	}
	keys := make([][]uint32, 0, len(x))
	vals := make([][]uint32, 0, len(x))
	for k, v := range x {
		keyPart := protocache.EncodeInt32(k)
		valPart, err := v.Encode()
		if err != nil {
			return nil, err
		}
		keys = append(keys, keyPart)
		vals = append(vals, valPart)
	}
	return protocache.EncodeMapParts(keys, vals, false)
}

func (x ModeDictEX) Serialize() ([]byte, error) { return protocache.SerializeEncoded(x.Encode()) }

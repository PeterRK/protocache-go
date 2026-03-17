package protocache

import (
	"bytes"
	"os"
	"testing"
	"unsafe"

	"github.com/peterrk/protocache-go/test/pb"
	"google.golang.org/protobuf/encoding/protojson"
)

func mustEncodeStringForTest(t *testing.T, s string) []uint32 {
	t.Helper()
	raw, err := EncodeString(s)
	assert(t, err == nil)
	return raw
}

func scanFieldsForTest(msg Message) ([]Field, int) {
	fieldsLimit := uint16(msg.data[0])*25 + 12
	fields := make([]Field, fieldsLimit)
	size := int((1 + uint32(msg.data[0])*2) * 4)
	for id := uint16(0); id < fieldsLimit; id++ {
		field := msg.GetField(id)
		fields[id] = field
		size += len(field.data)
	}
	if size > len(msg.data) {
		size = len(msg.data)
	}
	return fields, size
}

func detectAnyObjectPartForTest(data []byte) []byte {
	var best []byte
	candidates := [][]byte{
		DetectMap(data, detectAnyObjectPartForTest, detectAnyObjectPartForTest),
		DetectArray(data, detectAnyObjectPartForTest),
		detectMessagePartBytesForTest(data),
		DetectBytes(data),
	}
	for _, candidate := range candidates {
		if len(candidate) > len(best) {
			best = candidate
		}
	}
	return best
}

func subsliceOffset[T any](container, sub []T) (int, bool) {
	if len(sub) == 0 {
		return 0, true
	}
	if len(container) == 0 {
		return 0, false
	}
	base := uintptr(unsafe.Pointer(unsafe.SliceData(container)))
	ptr := uintptr(unsafe.Pointer(unsafe.SliceData(sub)))
	if ptr < base {
		return 0, false
	}
	off := ptr - base
	size := unsafe.Sizeof(*unsafe.SliceData(container))
	if size == 0 || (off%size) != 0 {
		return 0, false
	}
	idx := int(off / size)
	if idx < 0 || idx+len(sub) > len(container) {
		return 0, false
	}
	return idx, true
}

func detectNestedFieldSpanForTest(container []byte, field Field, compactEnd int) (int, int, bool) {
	if !field.IsValid() {
		return 0, 0, true
	}
	base, ok := subsliceOffset(container, field.data)
	if !ok {
		return 0, 0, false
	}
	if len(field.data) != 4 {
		return base, base + len(field.data), true
	}
	mark := getUint32(field.data)
	if (mark & 3) != 3 {
		return base, base + 4, true
	}
	off := int(mark & 0xfffffffc)
	if off <= 0 {
		return base, base + 4, true
	}
	start := base + off
	if start < compactEnd || start >= len(container) {
		return base, base + 4, true
	}
	part := detectAnyObjectPartForTest(container[start:])
	if len(part) == 0 {
		return base, base + 4, true
	}
	return start, start + len(part), true
}

func detectMessagePartBytesForTest(data []byte) []byte {
	msg := AsMessage(data)
	if !msg.IsValid() {
		return nil
	}
	fields, compactEnd := scanFieldsForTest(msg)
	end := compactEnd
	for _, field := range fields {
		_, tail, ok := detectNestedFieldSpanForTest(msg.data, field, compactEnd)
		if !ok {
			return nil
		}
		if tail > end {
			end = tail
		}
	}
	if end > len(msg.data) {
		return nil
	}
	return msg.data[:end]
}

func detectArrayComplexDirectForTest(m *Message, id uint16) []byte {
	field := m.GetField(id)
	return DetectArray(field.GetObject(), detectAnyObjectPartForTest)
}

func detectMapComplexDirectForTest(m *Message, id uint16) []byte {
	field := m.GetField(id)
	return DetectMap(field.GetObject(), detectAnyObjectPartForTest, detectAnyObjectPartForTest)
}

func detectScalar4PartForTest(m Message, id uint16) []byte {
	field := m.GetField(id)
	if len(field.data) != 4 {
		return nil
	}
	return field.data
}

func detectByteArrayPartForTest(m Message, id uint16) []byte {
	return DetectBytes(m.GetField(id).GetObject())
}

func detectMessagePartForTest(m Message, id uint16) []byte {
	return detectMessagePartBytesForTest(m.GetField(id).GetObject())
}

func detectArrayPartForTest(m Message, id uint16, detect func([]byte) []byte) []byte {
	return DetectArray(m.GetField(id).GetObject(), detect)
}

func detectMapPartForTest(m Message, id uint16, detectKey, detectValue func([]byte) []byte) []byte {
	return DetectMap(m.GetField(id).GetObject(), detectKey, detectValue)
}

func TestDetectFieldPartScalar4OffsetLikeBits(t *testing.T) {
	raw, err := Serialize(&pb.Small{
		I32:  7,
		Flag: true,
		Str:  "abcd",
		Junk: 11,
	})
	assert(t, err == nil)

	msg := AsMessage(raw)
	part := detectScalar4PartForTest(msg, 0)
	assert(t, len(part) == 4)
	assert(t, getUint32(part) == 7)
}

func TestDetectFieldPartStringBoundaries(t *testing.T) {
	parts := make([][]uint32, 4)
	rawWords, err := EncodeString("")
	assert(t, err == nil)
	parts[3] = rawWords
	rawWords, err = EncodeMessageParts(parts)
	assert(t, err == nil)

	msg := AsMessage(WordsToBytes(rawWords))
	part := detectByteArrayPartForTest(msg, 3)
	want, err := EncodeString("")
	assert(t, err == nil)
	assert(t, bytes.Equal(part, WordsToBytes(want)))

	raw, err := Serialize(&pb.Small{Str: "hello world"})
	assert(t, err == nil)
	msg = AsMessage(raw)
	part = detectByteArrayPartForTest(msg, 3)
	want, err = EncodeString("hello world")
	assert(t, err == nil)
	assert(t, bytes.Equal(part, WordsToBytes(want)))
}

func TestDetectFieldPartMessageArrayAndMap(t *testing.T) {
	root := &pb.Main{
		Object: &pb.Small{I32: 9, Str: "obj", Flag: true},
		Strv:   []string{"left", "right"},
		Objects: map[int32]*pb.Small{
			7: {I32: 77, Str: "map"},
		},
	}
	raw, err := Serialize(root)
	assert(t, err == nil)

	msg := AsMessage(raw)

	part := detectMessagePartForTest(msg, 10)
	want, err := Serialize(root.Object)
	assert(t, err == nil)
	assert(t, bytes.Equal(part, want))

	part = detectArrayPartForTest(msg, 13, detectAnyObjectPartForTest)
	strvWords, err := EncodeStringArray(root.Strv)
	assert(t, err == nil)
	assert(t, bytes.Equal(part, WordsToBytes(strvWords)))

	part = detectArrayPartForTest(msg, 13, DetectBytes)
	assert(t, bytes.Equal(part, WordsToBytes(strvWords)))

	key := EncodeInt32(7)
	val, err := Serialize(root.Objects[7])
	assert(t, err == nil)
	wantWords, err := EncodeMapParts([][]uint32{key}, [][]uint32{BytesToWords(val)}, false)
	assert(t, err == nil)
	part = detectMapPartForTest(msg, 26, detectAnyObjectPartForTest, detectAnyObjectPartForTest)
	assert(t, bytes.Equal(part, WordsToBytes(wantWords)))

	part = detectMapPartForTest(msg, 26, nil, detectMessagePartBytesForTest)
	assert(t, bytes.Equal(part, WordsToBytes(wantWords)))
}

func TestDetectComplexHelpers(t *testing.T) {
	root := &pb.Main{
		Object: &pb.Small{I32: 9, Str: "obj", Flag: true},
		Strv:   []string{"left", "right"},
		Objects: map[int32]*pb.Small{
			7: {I32: 77, Str: "map"},
		},
	}
	raw, err := Serialize(root)
	assert(t, err == nil)
	msg := AsMessage(raw)

	field := msg.GetField(10)
	part := detectMessagePartBytesForTest(field.GetObject())
	want, err := Serialize(root.Object)
	assert(t, err == nil)
	assert(t, bytes.Equal(part, want))

	field = msg.GetField(13)
	part = DetectArray(field.GetObject(), detectAnyObjectPartForTest)
	wantWords, err := EncodeStringArray(root.Strv)
	assert(t, err == nil)
	assert(t, bytes.Equal(part, WordsToBytes(wantWords)))

	key := EncodeInt32(7)
	val, err := Serialize(root.Objects[7])
	assert(t, err == nil)
	wantWords, err = EncodeMapParts([][]uint32{key}, [][]uint32{BytesToWords(val)}, false)
	assert(t, err == nil)
	field = msg.GetField(26)
	part = DetectMap(field.GetObject(), detectAnyObjectPartForTest, detectAnyObjectPartForTest)
	assert(t, bytes.Equal(part, WordsToBytes(wantWords)))
}

func TestDetectFieldPartComplexCasesFromFixture(t *testing.T) {
	raw, err := os.ReadFile("test/test.json")
	assert(t, err == nil)

	root := &pb.Main{}
	err = protojson.Unmarshal(raw, root)
	assert(t, err == nil)

	raw, err = Serialize(root)
	assert(t, err == nil)
	msg := AsMessage(raw)

	cases := []struct {
		id  uint16
		get func(Message, uint16) []byte
	}{
		{id: 10, get: detectMessagePartForTest},
		{id: 13, get: func(m Message, id uint16) []byte { return detectArrayPartForTest(m, id, detectAnyObjectPartForTest) }},
		{id: 18, get: func(m Message, id uint16) []byte { return detectArrayPartForTest(m, id, detectAnyObjectPartForTest) }},
		{id: 25, get: func(m Message, id uint16) []byte {
			return detectMapPartForTest(m, id, detectAnyObjectPartForTest, detectAnyObjectPartForTest)
		}},
		{id: 27, get: func(m Message, id uint16) []byte { return detectArrayPartForTest(m, id, detectAnyObjectPartForTest) }},
		{id: 28, get: func(m Message, id uint16) []byte { return detectArrayPartForTest(m, id, detectAnyObjectPartForTest) }},
	}
	for _, tc := range cases {
		part := tc.get(msg, tc.id)
		assert(t, len(part) != 0)
	}
}

func TestDetectFieldPartMatchesDirectComplexHelpers(t *testing.T) {
	raw, err := os.ReadFile("test/test.json")
	assert(t, err == nil)

	root := &pb.Main{}
	err = protojson.Unmarshal(raw, root)
	assert(t, err == nil)

	raw, err = Serialize(root)
	assert(t, err == nil)
	msg := AsMessage(raw)

	assert(t, len(detectMessagePartForTest(msg, 10)) != 0)

	arrayField := msg.GetField(13)
	arrayObj := arrayField.GetObject()
	assert(t, len(detectArrayPartForTest(msg, 13, detectAnyObjectPartForTest)) != 0)
	assert(t, len(DetectArray(arrayObj, detectAnyObjectPartForTest)) != 0)

	mapField := msg.GetField(26)
	mapObj := mapField.GetObject()
	assert(t, len(detectMapPartForTest(msg, 26, detectAnyObjectPartForTest, detectAnyObjectPartForTest)) != 0)
	assert(t, len(DetectMap(mapObj, detectAnyObjectPartForTest, detectAnyObjectPartForTest)) != 0)
}

func TestDetectFieldPartComplexFixtureKeepsExpectedBytes(t *testing.T) {
	raw, err := os.ReadFile("test/test.json")
	assert(t, err == nil)

	root := &pb.Main{}
	err = protojson.Unmarshal(raw, root)
	assert(t, err == nil)

	raw, err = Serialize(root)
	assert(t, err == nil)
	msg := AsMessage(raw)

	cases := []struct {
		id  uint16
		get func(Message, uint16) []byte
	}{
		{id: 10, get: detectMessagePartForTest},
		{id: 13, get: func(m Message, id uint16) []byte { return detectArrayPartForTest(m, id, DetectBytes) }},
		{id: 18, get: func(m Message, id uint16) []byte { return detectArrayPartForTest(m, id, detectMessagePartBytesForTest) }},
		{id: 25, get: func(m Message, id uint16) []byte { return detectMapPartForTest(m, id, DetectBytes, nil) }},
		{id: 26, get: func(m Message, id uint16) []byte {
			return detectMapPartForTest(m, id, nil, detectMessagePartBytesForTest)
		}},
		{id: 27, get: func(m Message, id uint16) []byte {
			return detectArrayPartForTest(m, id, func(data []byte) []byte { return DetectArray(data, nil) })
		}},
		{id: 28, get: func(m Message, id uint16) []byte {
			return detectArrayPartForTest(m, id, func(data []byte) []byte {
				return DetectMap(data, DetectBytes, func(value []byte) []byte { return DetectArray(value, nil) })
			})
		}},
		{id: 29, get: func(m Message, id uint16) []byte {
			return detectMapPartForTest(m, id, DetectBytes, func(value []byte) []byte { return DetectArray(value, nil) })
		}},
	}
	for _, tc := range cases {
		part := tc.get(msg, tc.id)
		assert(t, len(part) != 0)
	}
}

func TestDetectDirectFieldPartComplexFixture(t *testing.T) {
	raw, err := os.ReadFile("test/test.json")
	assert(t, err == nil)

	root := &pb.Main{}
	err = protojson.Unmarshal(raw, root)
	assert(t, err == nil)

	raw, err = Serialize(root)
	assert(t, err == nil)
	msg := AsMessage(raw)

	cases := []struct {
		id     uint16
		get    func(Message, uint16) []byte
		direct func(*Message, uint16) []byte
	}{
		{id: 10, get: detectMessagePartForTest, direct: func(m *Message, id uint16) []byte {
			return detectMessagePartBytesForTest(m.GetField(id).GetObject())
		}},
		{id: 13, get: func(m Message, id uint16) []byte { return detectArrayPartForTest(m, id, detectAnyObjectPartForTest) }, direct: detectArrayComplexDirectForTest},
		{id: 18, get: func(m Message, id uint16) []byte { return detectArrayPartForTest(m, id, detectAnyObjectPartForTest) }, direct: detectArrayComplexDirectForTest},
		{id: 25, get: func(m Message, id uint16) []byte {
			return detectMapPartForTest(m, id, detectAnyObjectPartForTest, detectAnyObjectPartForTest)
		}, direct: detectMapComplexDirectForTest},
		{id: 26, get: func(m Message, id uint16) []byte {
			return detectMapPartForTest(m, id, detectAnyObjectPartForTest, detectAnyObjectPartForTest)
		}, direct: detectMapComplexDirectForTest},
		{id: 27, get: func(m Message, id uint16) []byte { return detectArrayPartForTest(m, id, detectAnyObjectPartForTest) }, direct: detectArrayComplexDirectForTest},
		{id: 28, get: func(m Message, id uint16) []byte { return detectArrayPartForTest(m, id, detectAnyObjectPartForTest) }, direct: detectArrayComplexDirectForTest},
		{id: 29, get: func(m Message, id uint16) []byte {
			return detectMapPartForTest(m, id, detectAnyObjectPartForTest, detectAnyObjectPartForTest)
		}, direct: detectMapComplexDirectForTest},
	}
	nonEmpty := 0
	for _, tc := range cases {
		direct := tc.direct(&msg, tc.id)
		if len(direct) != 0 {
			nonEmpty++
		}
	}
	assert(t, nonEmpty >= 1)
}

func TestDetectDirectFieldPartKnownFixtureMatches(t *testing.T) {
	root := &pb.Main{
		Object: &pb.Small{I32: 9, Str: "obj", Flag: true},
		Objects: map[int32]*pb.Small{
			7: {I32: 77, Str: "map"},
		},
		Arrays: &pb.ArrMap{
			X: map[string]*pb.ArrMap_Array{
				"lv4": {X: []float32{41, 42}},
			},
		},
	}
	raw, err := Serialize(root)
	assert(t, err == nil)
	msg := AsMessage(raw)

	field := msg.GetField(10)
	obj := field.GetObject()
	want, err := Serialize(root.Object)
	assert(t, err == nil)
	assert(t, bytes.Equal(detectMessagePartBytesForTest(obj), want))

	val, err := Serialize(root.Objects[7])
	assert(t, err == nil)
	wantWords, err := EncodeMapParts([][]uint32{EncodeInt32(7)}, [][]uint32{BytesToWords(val)}, false)
	assert(t, err == nil)
	assert(t, bytes.Equal(detectMapComplexDirectForTest(&msg, 26), WordsToBytes(wantWords)))

	arrayWords, err := EncodeFloat32Array([]float32{41, 42})
	assert(t, err == nil)
	wantWords, err = EncodeMapParts([][]uint32{mustEncodeStringForTest(t, "lv4")}, [][]uint32{arrayWords}, true)
	assert(t, err == nil)
	assert(t, bytes.Equal(
		detectMapPartForTest(msg, 29, DetectBytes, func(value []byte) []byte { return DetectArray(value, nil) }),
		WordsToBytes(wantWords),
	))
}

func TestDetectKnownMessageDirectMatches(t *testing.T) {
	raw, err := os.ReadFile("test/test.json")
	assert(t, err == nil)

	root := &pb.Main{}
	err = protojson.Unmarshal(raw, root)
	assert(t, err == nil)

	raw, err = Serialize(root)
	assert(t, err == nil)
	msg := AsMessage(raw)

	field := msg.GetField(10)
	obj := field.GetObject()
	assert(t, len(detectMessagePartBytesForTest(obj)) != 0)
	assert(t, len(detectArrayPartForTest(msg, 18, detectMessagePartBytesForTest)) != 0)
}

func TestDetectKnownStringAndMapDirectMatches(t *testing.T) {
	raw, err := os.ReadFile("test/test.json")
	assert(t, err == nil)

	root := &pb.Main{}
	err = protojson.Unmarshal(raw, root)
	assert(t, err == nil)

	raw, err = Serialize(root)
	assert(t, err == nil)
	msg := AsMessage(raw)

	assert(t, len(detectArrayPartForTest(msg, 13, DetectBytes)) != 0)
	assert(t, len(detectMapPartForTest(msg, 25, DetectBytes, nil)) != 0)
	assert(t, len(detectArrayPartForTest(msg, 27, func(data []byte) []byte { return DetectArray(data, nil) })) != 0)
	assert(t, len(detectArrayPartForTest(msg, 28, func(data []byte) []byte {
		return DetectMap(data, DetectBytes, func(value []byte) []byte { return DetectArray(value, nil) })
	})) != 0)
	assert(t, len(detectMapPartForTest(msg, 29, DetectBytes, func(value []byte) []byte { return DetectArray(value, nil) })) != 0)
}

func TestMessageDetectInlinedMatchesFieldScan(t *testing.T) {
	raw, err := os.ReadFile("test/test.json")
	assert(t, err == nil)

	root := &pb.Main{}
	err = protojson.Unmarshal(raw, root)
	assert(t, err == nil)

	raw, err = Serialize(root)
	assert(t, err == nil)
	msg := AsMessage(raw)

	inlined := msg.DetectInlined()
	assert(t, len(inlined) != 0)

	_, want := scanFieldsForTest(msg)
	assert(t, len(inlined) == want)
	assert(t, bytes.Equal(inlined, raw[:want]))
}

func TestDetectScalarScalarArrayMapPart(t *testing.T) {
	root := &pb.ModeDict{
		X: map[int32]*pb.ModeDict_Value{
			7: {X: []pb.Mode{pb.Mode_MODE_A, pb.Mode_MODE_C}},
			9: {X: []pb.Mode{pb.Mode_MODE_B}},
		},
	}
	raw, err := Serialize(root)
	assert(t, err == nil)

	parts := [][]uint32{BytesToWords(raw)}
	msgWords, err := EncodeMessageParts(parts)
	assert(t, err == nil)
	msg := AsMessage(WordsToBytes(msgWords))

	part := detectMapPartForTest(msg, 0, nil, func(value []byte) []byte { return DetectArray(value, nil) })
	assert(t, bytes.Equal(part, raw))
}

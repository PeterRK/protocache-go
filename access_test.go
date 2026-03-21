package protocache

import (
	"os"
	"math/rand"
	"testing"
	"unsafe"

	"github.com/peterrk/protocache-go/test/pb"
	"google.golang.org/protobuf/encoding/protojson"
)

var benchCountSink uint32
var benchLocateSink uint32
var benchBytesSink []byte

func sumPairs32(v uint32) uint32 {
	sum := uint32(0)
	for i := 0; i < 16; i++ {
		sum += (v >> (i * 2)) & 3
	}
	return sum
}

func sumPairs64(v uint64) uint32 {
	sum := uint32(0)
	for i := 0; i < 32; i++ {
		sum += uint32((v >> (i * 2)) & 3)
	}
	return sum
}

func TestCount32(t *testing.T) {
	cases := []uint32{
		0,
		0xffffffff,
		0xaaaaaaaa,
		0x55555555,
		0x33333333,
		0xcccccccc,
		0x12345678,
		0x87654321,
	}
	for _, v := range cases {
		got := count32(v)
		want := sumPairs32(v)
		if got != want {
			t.Fatalf("count32(%#08x)=%d want %d", v, got, want)
		}
	}

	rng := rand.New(rand.NewSource(1))
	for i := 0; i < 10000; i++ {
		v := rng.Uint32()
		got := count32(v)
		want := sumPairs32(v)
		if got != want {
			t.Fatalf("count32(%#08x)=%d want %d", v, got, want)
		}
	}
}

func TestCount64(t *testing.T) {
	cases := []uint64{
		0,
		0xffffffffffffffff,
		0xaaaaaaaaaaaaaaaa,
		0x5555555555555555,
		0x3333333333333333,
		0xcccccccccccccccc,
		0x0123456789abcdef,
		0xfedcba9876543210,
	}
	for _, v := range cases {
		got := count64(v)
		want := sumPairs64(v)
		if got != want {
			t.Fatalf("count64(%#016x)=%d want %d", v, got, want)
		}
	}

	rng := rand.New(rand.NewSource(1))
	for i := 0; i < 10000; i++ {
		v := rng.Uint64()
		got := count64(v)
		want := sumPairs64(v)
		if got != want {
			t.Fatalf("count64(%#016x)=%d want %d", v, got, want)
		}
	}
}

func TestExtractBytes(t *testing.T) {
	longPayload := []byte("abcdefghijklmnopqrstuvwxyz1234567890ABCD")
	cases := [][]byte{
		nil,
		{},
		{0},
		{4, 't'},
		{8, 't', 'e'},
		{12, 't', 'e', 's'},
		{16, 't', 'e', 's', 't'},
		append([]byte{0xa0, 0x01}, longPayload...),
	}
	expected := [][]byte{
		nil,
		nil,
		{},
		[]byte("t"),
		[]byte("te"),
		[]byte("tes"),
		[]byte("test"),
		longPayload,
	}
	for i, raw := range cases {
		got := extractBytes(raw)
		want := expected[i]
		if string(got) != string(want) || (got == nil) != (want == nil) {
			t.Fatalf("extractBytes(%v)=%v want %v", raw, got, want)
		}
	}
}

func BenchmarkExtractBytesShortLegacy(b *testing.B) {
	data := []byte{16, 't', 'e', 's', 't'}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out := extractBytesLegacy(data)
		benchCountSink += uint32(len(out))
	}
}

func BenchmarkExtractBytesShortCurrent(b *testing.B) {
	data := []byte{16, 't', 'e', 's', 't'}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out := extractBytes(data)
		benchCountSink += uint32(len(out))
	}
}

func BenchmarkExtractBytesLongLegacy(b *testing.B) {
	data := append([]byte{0xa0, 0x01}, []byte("abcdefghijklmnopqrstuvwxyz1234567890ABCD")...)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out := extractBytesLegacy(data)
		benchCountSink += uint32(len(out))
	}
}

func BenchmarkExtractBytesLongCurrent(b *testing.B) {
	data := append([]byte{0xa0, 0x01}, []byte("abcdefghijklmnopqrstuvwxyz1234567890ABCD")...)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out := extractBytes(data)
		benchCountSink += uint32(len(out))
	}
}

func count32Legacy(v uint32) uint32 {
	v = (v & 0x33333333) + ((v >> 2) & 0x33333333)
	v = v + (v >> 4)
	v = (v & 0x0f0f0f0f) + ((v >> 8) & 0x0f0f0f0f)
	v = v + (v >> 16)
	return v & 0xff
}

func count64Legacy(v uint64) uint32 {
	v = (v & 0x3333333333333333) + ((v >> 2) & 0x3333333333333333)
	v = v + (v >> 4)
	v = (v & 0x0f0f0f0f0f0f0f0f) + ((v >> 8) & 0x0f0f0f0f0f0f0f0f)
	v = v + (v >> 16)
	v = v + (v >> 32)
	return uint32(v) & 0xff
}

func extractBytesLegacy(data []byte) []byte {
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

func getObjectLegacy(f *Field) []byte {
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
	return unsafe.Slice(unsafe.SliceData(f.data), cap(f.data))[off:]
}

func locateFieldLegacy(m *Message, id uint16) (off uint32, width uint32, ok bool) {
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
		off += count32Legacy(v & ^(uint32(0xffffffff) << (uint32(id) << 1)))
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
		off += count64Legacy(v & ^(uint64(0xffffffffffffffff) << (b << 1)))
	}
	off *= 4
	width *= 4
	if off+width > uint32(len(m.data)) {
		return 0, 0, false
	}
	return off, width, true
}

func loadBenchMessage(b *testing.B) Message {
	b.Helper()

	raw, err := os.ReadFile("test/test.json")
	if err != nil {
		b.Fatal(err)
	}
	message := &pb.Main{}
	if err := protojson.Unmarshal(raw, message); err != nil {
		b.Fatal(err)
	}
	encoded, err := Serialize(message)
	if err != nil {
		b.Fatal(err)
	}
	msg := AsMessage(encoded)
	if !msg.IsValid() {
		b.Fatal("invalid bench message")
	}
	return msg
}

func BenchmarkCount32Legacy(b *testing.B) {
	rng := rand.New(rand.NewSource(1))
	data := make([]uint32, 1024)
	for i := range data {
		data[i] = rng.Uint32()
	}

	b.ReportAllocs()
	var sink uint32
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += count32Legacy(data[i&1023])
	}
	benchCountSink = sink
}

func BenchmarkCount32Current(b *testing.B) {
	rng := rand.New(rand.NewSource(1))
	data := make([]uint32, 1024)
	for i := range data {
		data[i] = rng.Uint32()
	}

	b.ReportAllocs()
	var sink uint32
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += count32(data[i&1023])
	}
	benchCountSink = sink
}

func BenchmarkCount64Legacy(b *testing.B) {
	rng := rand.New(rand.NewSource(1))
	data := make([]uint64, 1024)
	for i := range data {
		data[i] = rng.Uint64()
	}

	b.ReportAllocs()
	var sink uint32
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += count64Legacy(data[i&1023])
	}
	benchCountSink = sink
}

func BenchmarkCount64Current(b *testing.B) {
	rng := rand.New(rand.NewSource(1))
	data := make([]uint64, 1024)
	for i := range data {
		data[i] = rng.Uint64()
	}

	b.ReportAllocs()
	var sink uint32
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sink += count64(data[i&1023])
	}
	benchCountSink = sink
}

func BenchmarkLocateFieldLegacy(b *testing.B) {
	msg := loadBenchMessage(b)
	ids := []uint16{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 16, 17, 21, 22, 26, 27, 28, 29, 31}

	b.ReportAllocs()
	var sum uint32
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := ids[i%len(ids)]
		off, width, ok := locateFieldLegacy(&msg, id)
		if ok {
			sum += off + width
		}
	}
	benchLocateSink = sum
}

func BenchmarkLocateFieldCurrent(b *testing.B) {
	msg := loadBenchMessage(b)
	ids := []uint16{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 16, 17, 21, 22, 26, 27, 28, 29, 31}

	b.ReportAllocs()
	var sum uint32
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		id := ids[i%len(ids)]
		off, width, ok := msg.locateField(id)
		if ok {
			sum += off + width
		}
	}
	benchLocateSink = sum
}

func BenchmarkGetObjectInlineLegacy(b *testing.B) {
	field := Field{data: []byte{1, 2, 3, 4}}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchBytesSink = getObjectLegacy(&field)
	}
}

func BenchmarkGetObjectInlineCurrent(b *testing.B) {
	field := Field{data: []byte{1, 2, 3, 4}}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchBytesSink = field.GetObject()
	}
}

func BenchmarkGetObjectOffsetLegacy(b *testing.B) {
	backing := make([]byte, 32)
	putUint32(backing[:4], 8|3)
	copy(backing[8:], []byte("payload"))
	field := Field{data: backing[:4:len(backing)]}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchBytesSink = getObjectLegacy(&field)
	}
}

func BenchmarkGetObjectOffsetCurrent(b *testing.B) {
	backing := make([]byte, 32)
	putUint32(backing[:4], 8|3)
	copy(backing[8:], []byte("payload"))
	field := Field{data: backing[:4:len(backing)]}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchBytesSink = field.GetObject()
	}
}

//go:build ignore

package test

// This file preserves the optional FlatBuffers/Fory comparison benchmarks.
// It is excluded from normal builds and from go mod tidy so the main module
// does not acquire dependencies used only by historical comparisons.

import (
	"fmt"
	"os"
	"testing"

	"github.com/peterrk/protocache-go/test/fb"
	"github.com/peterrk/protocache-go/test/fr"
	"github.com/peterrk/protocache-go/test/pb"
	"github.com/peterrk/protocache-go/test/pc"
	"google.golang.org/protobuf/proto"
)

func TestProtoCacheBenchmark(t *testing.T) {
	var junk Junk
	raw, err := os.ReadFile("test.pc")
	if err != nil {
		t.Fatal(err)
	}
	root := pc.AS_Main(raw)
	junk.traversePcMain(root)
	fmt.Println(junk.fuse())
}

func TestProtobufBenchmark(t *testing.T) {
	var junk Junk
	raw, err := os.ReadFile("test.pb")
	if err != nil {
		t.Fatal(err)
	}
	root := &pb.Main{}
	if err := proto.Unmarshal(raw, root); err != nil {
		t.Fatal(err)
	}
	junk.traversePbMain(root)
	fmt.Println(junk.fuse())

	junk = Junk{}
	junk.traversePbMessage(root.ProtoReflect())
	fmt.Println(junk.fuse())
}

func TestFlatbuffersBenchmark(t *testing.T) {
	var junk Junk
	raw, err := os.ReadFile("test.fb")
	if err != nil {
		t.Fatal(err)
	}
	root := fb.GetRootAsMain(raw, 0)
	junk.traverseFbMain(root)
	fmt.Println(junk.fuse())
}

func TestForyBenchmark(t *testing.T) {
	var junk Junk
	raw, err := os.ReadFile("test.fr")
	if err != nil {
		t.Fatal(err)
	}
	f, err := fr.New()
	if err != nil {
		t.Fatal(err)
	}
	var root fr.Main
	if err := f.Deserialize(raw, &root); err != nil {
		t.Fatal(err)
	}
	junk.traverseForyMain(&root)
	fmt.Println(junk.fuse())
}

func BenchmarkFory(b *testing.B) {
	raw, err := os.ReadFile("test.fr")
	if err != nil {
		b.Fatal(err)
	}
	f, err := fr.New()
	if err != nil {
		b.Fatal(err)
	}
	var junk Junk
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var root fr.Main
		if err := f.Deserialize(raw, &root); err != nil {
			b.Fatal(err)
		}
		junk.traverseForyMain(&root)
	}
	benchmarkFuse = junk.fuse()
}

func BenchmarkForySerialize(b *testing.B) {
	f, err := fr.New()
	if err != nil {
		b.Fatal(err)
	}
	root := fr.CreateObject()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out, err := f.Serialize(root)
		if err != nil {
			b.Fatal(err)
		}
		benchmarkBytes = out
	}
}

func BenchmarkFlatbuffers(b *testing.B) {
	raw, err := os.ReadFile("test.fb")
	if err != nil {
		b.Fatal(err)
	}
	var junk Junk
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		root := fb.GetRootAsMain(raw, 0)
		junk.traverseFbMain(root)
	}
	benchmarkFuse = junk.fuse()
}

func (p *Junk) traverseForySmall(root *fr.Small) {
	if root == nil {
		return
	}
	p.u32 += uint32(root.I32)
	p.consumeBool(root.Flag)
	p.consumeString(root.Str)
}

func (p *Junk) traverseForyVec2D(root *fr.Vec2D) {
	if root == nil {
		return
	}
	for _, one := range root.X {
		for _, v := range one.X {
			p.f32 += v
		}
	}
}

func (p *Junk) traverseForyArrMap(root *fr.ArrMap) {
	if root == nil {
		return
	}
	for key, val := range root.X {
		p.consumeString(key)
		for _, v := range val.X {
			p.f32 += v
		}
	}
}

func (p *Junk) traverseForyMain(root *fr.Main) {
	p.u32 += uint32(root.I32) + root.U32 + uint32(root.Mode)
	p.consumeBool(root.Flag)
	p.u32 += uint32(root.TI32) + uint32(root.TS32) + root.TU32
	for _, v := range root.I32v {
		p.u32 += uint32(v)
	}
	p.u64 += uint64(root.I64) + root.U64 +
		uint64(root.TI64) + uint64(root.TS64) + root.TU64
	for _, v := range root.U64v {
		p.u64 += v
	}
	for _, v := range root.Flags {
		p.consumeBool(v)
	}
	p.consumeString(root.Str)
	p.consumeBytes(root.Data)
	for _, v := range root.Strv {
		p.consumeString(v)
	}
	for _, v := range root.Datav {
		p.consumeBytes(v)
	}

	p.f32 += root.F32
	for _, v := range root.F32v {
		p.f32 += v
	}
	p.f64 += root.F64
	for _, v := range root.F64v {
		p.f64 += v
	}

	p.traverseForySmall(root.Object)
	for i := range root.Objectv {
		p.traverseForySmall(&root.Objectv[i])
	}

	for key, val := range root.Index {
		p.consumeString(key)
		p.u32 += uint32(val)
	}
	for key, val := range root.Objects {
		p.u32 += uint32(key)
		p.traverseForySmall(&val)
	}

	p.traverseForyVec2D(root.Matrix)
	for i := range root.Vector {
		p.traverseForyArrMap(&root.Vector[i])
	}
	p.traverseForyArrMap(root.Arrays)
}

func (p *Junk) traverseFbSmall(root *fb.Small) {
	p.u32 += uint32(root.I32())
	p.consumeBool(root.Flag())
	p.consumeBytes(root.Str())
}

func (p *Junk) traverseFbVec2D(root *fb.Vec2D) {
	var unit fb.Vec1D
	for i := 0; i < root.AliasLength(); i++ {
		root.Alias(&unit, i)
		for j := 0; j < unit.AliasLength(); j++ {
			p.f32 += unit.Alias(j)
		}
	}
}

func (p *Junk) traverseFbArrMap(root *fb.ArrMap) {
	var unit fb.Array
	var pair fb.ArrMapEntry
	for i := 0; i < root.AliasLength(); i++ {
		root.Alias(&pair, i)
		p.consumeBytes(pair.Key())
		pair.Value(&unit)
		for j := 0; j < unit.AliasLength(); j++ {
			p.f32 += unit.Alias(j)
		}
	}
}

func (p *Junk) traverseFbMain(root *fb.Main) {
	p.u32 += uint32(root.I32()) + root.U32() + uint32(root.Mode())
	p.consumeBool(root.Flag())
	p.u32 += uint32(root.TI32()) + uint32(root.TS32()) + root.TU32()
	for i := 0; i < root.I32vLength(); i++ {
		p.u32 += uint32(root.I32v(i))
	}
	p.u64 += uint64(root.I64()) + root.U64() +
		uint64(root.TI64()) + uint64(root.TS64()) + root.TU64()
	for i := 0; i < root.U64vLength(); i++ {
		p.u64 += root.U64v(i)
	}
	for i := 0; i < root.FlagsLength(); i++ {
		p.consumeBool(root.Flags(i))
	}
	p.consumeBytes(root.Str())

	data := make([]byte, root.DataLength())
	for i := range data {
		data[i] = byte(root.Data(i))
	}
	p.consumeBytes(data)
	for i := 0; i < root.StrvLength(); i++ {
		p.consumeBytes(root.Strv(i))
	}

	var bytesValue fb.Bytes
	for i := 0; i < root.DatavLength(); i++ {
		root.Datav(&bytesValue, i)
		data := make([]byte, bytesValue.AliasLength())
		for j := range data {
			data[j] = byte(bytesValue.Alias(j))
		}
		p.consumeBytes(data)
	}

	p.f32 += root.F32()
	for i := 0; i < root.F32vLength(); i++ {
		p.f32 += root.F32v(i)
	}
	p.f64 += root.F64()
	for i := 0; i < root.F64vLength(); i++ {
		p.f64 += root.F64v(i)
	}

	var small fb.Small
	root.Object(&small)
	p.traverseFbSmall(&small)
	for i := 0; i < root.ObjectvLength(); i++ {
		root.Objectv(&small, i)
		p.traverseFbSmall(&small)
	}

	var pair1 fb.Map1Entry
	for i := 0; i < root.IndexLength(); i++ {
		root.Index(&pair1, i)
		p.consumeBytes(pair1.Key())
		p.u32 += uint32(pair1.Value())
	}

	var pair2 fb.Map2Entry
	for i := 0; i < root.ObjectsLength(); i++ {
		root.Objects(&pair2, i)
		p.u32 += uint32(pair2.Key())
		pair2.Value(&small)
		p.traverseFbSmall(&small)
	}

	var vec2d fb.Vec2D
	root.Matrix(&vec2d)
	p.traverseFbVec2D(&vec2d)

	var arrMap fb.ArrMap
	for i := 0; i < root.VectorLength(); i++ {
		root.Vector(&arrMap, i)
		p.traverseFbArrMap(&arrMap)
	}
	root.Arrays(&arrMap)
	p.traverseFbArrMap(&arrMap)
}

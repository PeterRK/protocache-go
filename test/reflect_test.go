package test

import (
	"errors"
	"os"
	"testing"

	pc "github.com/peterrk/protocache-go"
	"github.com/peterrk/protocache-go/reflect"
	"github.com/peterrk/protocache-go/reflect/compiler"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
)

func parseProtoForTest(t testing.TB, path string) *descriptorpb.FileDescriptorProto {
	t.Helper()

	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	proto, err := compiler.ParseProto(raw)
	if errors.Is(err, compiler.ErrUnsupported) {
		t.Skip(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	return proto
}

func TestReflection(t *testing.T) {
	proto := parseProtoForTest(t, "reflect-test.proto")

	var taggedPool reflect.DescriptorPool
	assert(t, taggedPool.Register(proto) == nil)
	assert(t, taggedPool.Find("test.Main") != nil)

	proto = parseProtoForTest(t, "test.proto")

	var pool reflect.DescriptorPool
	assert(t, pool.Register(proto) == nil)

	root := pool.Find("test.Main")
	assert(t, root != nil)

	field := root.Lookup("f64")
	assert(t, field != nil)
	assert(t, !field.IsRepeated())
	assert(t, field.Value() == reflect.TypeFloat64)

	field = root.Lookup("strv")
	assert(t, field != nil)
	assert(t, field.IsRepeated())
	assert(t, field.Value() == reflect.TypeString)

	field = root.Lookup("mode")
	assert(t, field != nil)
	assert(t, !field.IsRepeated())
	assert(t, field.Value() == reflect.TypeEnum)

	field = root.Lookup("object")
	assert(t, field != nil)
	assert(t, !field.IsRepeated())
	assert(t, field.Value() == reflect.TypeMessage)
	object := pool.Find(field.ValueType())
	assert(t, object != nil)
	assert(t, object.Alias() == nil)
	field = object.Lookup("flag")
	assert(t, field != nil)
	assert(t, !field.IsRepeated())
	assert(t, field.Value() == reflect.TypeBool)

	field = root.Lookup("index")
	assert(t, field != nil)
	assert(t, field.IsRepeated())
	assert(t, field.IsMap())

	field = root.Lookup("matrix")
	assert(t, field != nil)
	assert(t, !field.IsRepeated())
	assert(t, field.Value() == reflect.TypeMessage)
	object = field.ValueDescriptor()
	assert(t, object != nil)
	field = object.Alias()
	assert(t, field != nil)
	assert(t, field.IsRepeated())
	assert(t, !field.IsMap())
	assert(t, field.Value() == reflect.TypeMessage)
	object = field.ValueDescriptor()
	assert(t, object != nil)
	field = object.Alias()
	assert(t, field != nil)
	assert(t, field.IsRepeated())
	assert(t, !field.IsMap())
	assert(t, field.Value() == reflect.TypeFloat32)

	field = root.Lookup("arrays")
	assert(t, field != nil)
	assert(t, !field.IsRepeated())
	assert(t, field.Value() == reflect.TypeMessage)
	object = field.ValueDescriptor()
	assert(t, object != nil)
	field = object.Alias()
	assert(t, field != nil)
	assert(t, field.IsRepeated())
	assert(t, field.IsMap())
	assert(t, field.Key() == reflect.TypeString)
	assert(t, field.Value() == reflect.TypeMessage)
	object = field.ValueDescriptor()
	assert(t, object != nil)
	field = object.Alias()
	assert(t, field != nil)
	assert(t, field.IsRepeated())
	assert(t, !field.IsMap())
	assert(t, field.Value() == reflect.TypeFloat32)

	raw := loadMain(t)
	msg := pc.AsMessage(raw)
	field = root.Lookup("f32")
	assert(t, field != nil)
	unit := msg.GetField(field.Id())
	assert(t, unit.GetFloat32() == -2.1)
}

/*
	func TestBenchmarkReflect(t *testing.T) {
		raw, err := os.ReadFile("test.proto")
		assert(t, err == nil)

		proto, err := compiler.ParseProto(raw)
		assert(t, err == nil)

		var pool reflect.DescriptorPool
		assert(t, pool.Register(proto) == nil)

		descriptor := pool.Find("test.Main")
		assert(t, descriptor != nil)

		var junk Junk
		raw, err = os.ReadFile("test.pc")
		if err != nil {
			fmt.Println(err)
			return
		}
		root := pc.AsMessage(raw)
		junk.traversePcMessage(descriptor, root)
		fmt.Println(junk.fuse())
	}
*/

func BenchmarkProtoCacheReflect(b *testing.B) {
	proto := parseProtoForTest(b, "test.proto")
	var pool reflect.DescriptorPool
	if err := pool.Register(proto); err != nil {
		b.Fatal(err)
	}
	descriptor := pool.Find("test.Main")
	if descriptor == nil {
		b.Fatal("fail to get root")
	}
	raw, err := os.ReadFile("test.pc")
	if err != nil {
		b.Fatal(err)
	}
	var junk Junk
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		root := pc.AsMessage(raw)
		junk.traversePcMessage(descriptor, root)
	}
	benchmarkFuse = junk.fuse()
}

func (p *Junk) traversePcMessage(descriptor *reflect.Descriptor, root pc.Message) {
	descriptor.Traverse(func(name string, field *reflect.Field) bool {
		unit := root.GetField(field.Id())
		if unit.IsValid() {
			p.traversePcField(field, unit)
		}
		return true
	})
}

func (p *Junk) traversePcField(descriptor *reflect.Field, field pc.Field) {
	if descriptor.IsMap() {
		pack := field.GetMap()
		for i := uint32(0); i < pack.Size(); i++ {
			key := pack.Key(i)
			switch descriptor.Key() {
			case reflect.TypeString:
				p.consumeString(key.GetString())
			case reflect.TypeUint64:
				p.u64 += key.GetUint64()
			case reflect.TypeUint32:
				p.u32 += key.GetUint32()
			case reflect.TypeInt64:
				p.u64 += uint64(key.GetInt64())
			case reflect.TypeInt32:
				p.u32 += uint32(key.GetInt32())
			}
			p.accessPcField(descriptor, pack.Value(i))
		}
	} else if descriptor.IsRepeated() {
		switch descriptor.Value() {
		case reflect.TypeMessage:
			array := field.GetArray()
			subtype := descriptor.ValueDescriptor()
			if alias := subtype.Alias(); alias != nil {
				for i := uint32(0); i < array.Size(); i++ {
					p.traversePcField(alias, array.Get(i))
				}
			} else {
				for i := uint32(0); i < array.Size(); i++ {
					unit := array.Get(i)
					p.traversePcMessage(subtype, unit.GetMessage())
				}
			}
		case reflect.TypeBytes:
			array := field.GetArray()
			for i := uint32(0); i < array.Size(); i++ {
				unit := array.Get(i)
				p.consumeBytes(unit.GetBytes())
			}
		case reflect.TypeString:
			array := field.GetArray()
			for i := uint32(0); i < array.Size(); i++ {
				unit := array.Get(i)
				p.consumeString(unit.GetString())
			}
		case reflect.TypeFloat64:
			for _, v := range field.GetFloat64Array() {
				p.f64 += v
			}
		case reflect.TypeFloat32:
			for _, v := range field.GetFloat32Array() {
				p.f32 += v
			}
		case reflect.TypeUint64:
			for _, v := range field.GetUint64Array() {
				p.u64 += v
			}
		case reflect.TypeUint32:
			for _, v := range field.GetUint32Array() {
				p.u32 += v
			}
		case reflect.TypeInt64:
			for _, v := range field.GetInt64Array() {
				p.u64 += uint64(v)
			}
		case reflect.TypeInt32:
			for _, v := range field.GetInt32Array() {
				p.u32 += uint32(v)
			}
		case reflect.TypeBool:
			for _, v := range field.GetBoolArray() {
				p.consumeBool(v)
			}
		case reflect.TypeEnum:
			for _, v := range field.GetEnumValueArray() {
				p.u32 += uint32(v)
			}
		}

	} else {
		p.accessPcField(descriptor, field)
	}
}

func (p *Junk) accessPcField(descriptor *reflect.Field, field pc.Field) {
	switch descriptor.Value() {
	case reflect.TypeMessage:
		subtype := descriptor.ValueDescriptor()
		if alias := subtype.Alias(); alias != nil {
			p.traversePcField(alias, field)
		} else {
			p.traversePcMessage(subtype, field.GetMessage())
		}
	case reflect.TypeBytes:
		p.consumeBytes(field.GetBytes())
	case reflect.TypeString:
		p.consumeString(field.GetString())
	case reflect.TypeFloat64:
		p.f64 += field.GetFloat64()
	case reflect.TypeFloat32:
		p.f32 += field.GetFloat32()
	case reflect.TypeUint64:
		p.u64 += field.GetUint64()
	case reflect.TypeUint32:
		p.u32 += field.GetUint32()
	case reflect.TypeInt64:
		p.u64 += uint64(field.GetInt64())
	case reflect.TypeInt32:
		p.u32 += uint32(field.GetInt32())
	case reflect.TypeBool:
		p.consumeBool(field.GetBool())
	case reflect.TypeEnum:
		p.u32 += uint32(field.GetEnumValue())
	}
}

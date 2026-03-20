package test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"testing"

	"github.com/peterrk/protocache-go"
	pcreflect "github.com/peterrk/protocache-go/reflect"
	"github.com/peterrk/protocache-go/test/pb"
	"github.com/peterrk/protocache-go/test/pc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

func assert(t *testing.T, state bool) {
	if !state {
		t.FailNow()
	}
}

func buildDynamicMessage(t *testing.T, file *descriptorpb.FileDescriptorProto, name protoreflect.Name) protoreflect.Message {
	t.Helper()

	desc, err := protodesc.NewFile(file, nil)
	if err != nil {
		t.Fatal(err)
	}
	message := desc.Messages().ByName(name)
	if message == nil {
		t.Fatalf("missing message %s", name)
	}
	return dynamicpb.NewMessage(message)
}

func TestEmptyViews(t *testing.T) {
	var field protocache.Field
	assert(t, !field.IsValid())
	assert(t, field.GetInt32() == 0)
	assert(t, !field.GetBool())
	assert(t, field.GetString() == "")
	assert(t, field.GetBytes() == nil)
	msgField := field.GetMessage()
	arrField := field.GetArray()
	mapField := field.GetMap()
	assert(t, !msgField.IsValid())
	assert(t, !arrField.IsValid())
	assert(t, !mapField.IsValid())

	var msg protocache.Message
	assert(t, !msg.IsValid())
	assert(t, !msg.HasField(0))
	msgField0 := msg.GetField(0)
	assert(t, !msgField0.IsValid())

	var arr protocache.Array
	assert(t, !arr.IsValid())
	assert(t, arr.Size() == 0)
	arrField0 := arr.Get(0)
	assert(t, !arrField0.IsValid())

	var pack protocache.Map
	assert(t, !pack.IsValid())
	assert(t, pack.Size() == 0)
	key0 := pack.Key(0)
	val0 := pack.Value(0)
	found0 := pack.FindByInt32(0)
	assert(t, !key0.IsValid())
	assert(t, !val0.IsValid())
	assert(t, !found0.IsValid())
}

func TestTinyFixtures(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		raw, err := os.ReadFile("test-empty.json")
		assert(t, err == nil)

		message := &pb.Main{}
		err = protojson.Unmarshal(raw, message)
		assert(t, err == nil)

		cooked, err := protocache.Serialize(message)
		assert(t, err == nil)
		assert(t, len(cooked) == 4)
		msg := protocache.AsMessage(cooked)
		assert(t, msg.IsValid())
	})

	t.Run("tiny", func(t *testing.T) {
		raw, err := os.ReadFile("test-tiny.json")
		assert(t, err == nil)

		message := &pb.Main{}
		err = protojson.Unmarshal(raw, message)
		assert(t, err == nil)

		cooked, err := protocache.Serialize(message)
		assert(t, err == nil)
		assert(t, len(cooked) == 8)

		msg := protocache.AsMessage(cooked)
		assert(t, msg.IsValid())
		assert(t, msg.HasField(0))
		assert(t, !msg.HasField(1))
		field0 := msg.GetField(0)
		assert(t, field0.GetInt32() == 123)
	})
}

func TestProtoCache(t *testing.T) {
	raw, err := os.ReadFile("test.json")
	assert(t, err == nil)

	message := &pb.Main{}
	err = protojson.Unmarshal(raw, message)
	assert(t, err == nil)

	raw, err = protocache.Serialize(message)
	assert(t, err == nil)

	//	raw, err := os.ReadFile("test.pc")
	//	assert(t, err == nil)

	root := pc.AS_Main(raw)
	assert(t, root.IsValid())

	assert(t, root.GetI32() == -999)
	assert(t, root.GetU32() == 1234)
	assert(t, root.GetI64() == -9876543210)
	assert(t, root.GetU64() == 98765432123456789)
	assert(t, root.GetFlag())
	assert(t, root.GetMode() == pc.Mode_MODE_C)
	assert(t, root.GetStr() == "Hello World!")
	assert(t, bytes.Equal(root.GetData(), []byte("abc123!?$*&()'-=@~")))
	assert(t, root.GetF32() == -2.1)
	assert(t, root.GetF64() == 1.0)

	leaf := root.GetObject()
	assert(t, leaf.IsValid())
	assert(t, leaf.GetI32() == 88)
	assert(t, !leaf.GetFlag())
	assert(t, leaf.GetStr() == "tmp")

	i32v := root.GetI32V()
	assert(t, len(i32v) == 2)
	assert(t, i32v[0] == 1)
	assert(t, i32v[1] == 2)

	u64v := root.GetU64V()
	assert(t, len(u64v) == 1)
	assert(t, u64v[0] == 12345678987654321)

	expectedStrv := []string{
		"abc", "apple", "banana", "orange", "pear", "grape",
		"strawberry", "cherry", "mango", "watermelon"}
	strv := root.GetStrv()
	assert(t, strv.Size() == uint32(len(expectedStrv)))
	for i := 0; i < len(expectedStrv); i++ {
		assert(t, strv.Get(uint32(i)) == expectedStrv[i])
	}

	f32v := root.GetF32V()
	assert(t, len(f32v) == 2)
	assert(t, f32v[0] == 1.1)
	assert(t, f32v[1] == 2.2)

	f64v := root.GetF64V()
	assert(t, len(f64v) == 5)
	assert(t, f64v[0] == 9.9)
	assert(t, f64v[1] == 8.8)
	assert(t, f64v[2] == 7.7)
	assert(t, f64v[3] == 6.6)
	assert(t, f64v[4] == 5.5)

	expectedFlags := []bool{true, true, false, true, false, false, false}
	flags := root.GetFlags()
	assert(t, len(flags) == len(expectedFlags))
	for i := 0; i < len(expectedFlags); i++ {
		assert(t, flags[i] == expectedFlags[i])
	}

	objs := root.GetObjectv()
	assert(t, objs.Size() == 3)
	obj0 := objs.Get(0)
	assert(t, obj0.GetI32() == 1)
	obj1 := objs.Get(1)
	assert(t, obj1.GetFlag())
	obj2 := objs.Get(2)
	assert(t, obj2.GetStr() == "good luck!")

	map1 := root.GetIndex()
	assert(t, map1.Size() == 6)
	val1, found := map1.Find("abc-1")
	assert(t, found && val1 == 1)
	val1, found = map1.Find("abc-2")
	assert(t, found && val1 == 2)
	_, found = map1.Find("abc-3")
	assert(t, !found)
	_, found = map1.Find("abc-4")
	assert(t, !found)

	map2 := root.GetObjects()
	assert(t, map2.Size() == 4)
	for i := uint32(0); i < map2.Size(); i++ {
		key, val := map2.Key(i), map2.Value(i)
		assert(t, key == val.GetI32())
	}

	matrix := root.GetMatrix()
	assert(t, matrix.Size() == 3)
	line := matrix.Get(2)
	assert(t, line.Size() == 3)
	assert(t, line.Get(2) == 9)

	vector := root.GetVector()
	assert(t, vector.Size() == 2)
	map3 := vector.Get(0)
	assert(t, map3.Size() == 2)
	val3, found := map3.Find("lv2")
	assert(t, found)
	assert(t, val3.Size() == 2)
	assert(t, val3.Get(0) == 21)
	assert(t, val3.Get(1) == 22)

	map4 := root.GetArrays()
	val4, found := map4.Find("lv5")
	assert(t, found)
	assert(t, val4.Size() == 2)
	assert(t, val4.Get(0) == 51)
	assert(t, val4.Get(1) == 52)
}

func TestAlias(t *testing.T) {
	raw, err := os.ReadFile("test-alias.json")
	assert(t, err == nil)

	message := &pb.Main{}
	err = protojson.Unmarshal(raw, message)
	assert(t, err == nil)

	raw, err = protocache.Serialize(message)
	assert(t, err == nil)

	assert(t, len(raw) == 68)
	assert(t, binary.LittleEndian.Uint32(raw[20:]) == 0x0d)
	assert(t, binary.LittleEndian.Uint32(raw[24:]) == 1)
	assert(t, binary.LittleEndian.Uint32(raw[28:]) == 1)
}

func TestCompress(t *testing.T) {
	raw, err := os.ReadFile("test.json")
	assert(t, err == nil)

	message := &pb.Main{}
	err = protojson.Unmarshal(raw, message)
	assert(t, err == nil)

	raw, err = protocache.Serialize(message)
	assert(t, err == nil)

	cooked := protocache.Compress(raw)
	assert(t, len(cooked) != 0 && len(cooked) < len(raw))
	back, err := protocache.Decompress(cooked)
	assert(t, err == nil && bytes.Equal(raw, back))
}

func TestBigObject(t *testing.T) {
	const fieldsCount = 1000

	fields := make([]*descriptorpb.FieldDescriptorProto, 0, fieldsCount)
	for i := 1; i <= fieldsCount; i++ {
		fields = append(fields, &descriptorpb.FieldDescriptorProto{
			Name:   proto.String(fmt.Sprintf("i%d", i)),
			Number: proto.Int32(int32(i)),
			Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
			Type:   descriptorpb.FieldDescriptorProto_TYPE_INT32.Enum(),
		})
	}

	file := &descriptorpb.FileDescriptorProto{
		Syntax:  proto.String("proto3"),
		Name:    proto.String("big-object.proto"),
		Package: proto.String("cov"),
		MessageType: []*descriptorpb.DescriptorProto{
			{
				Name:  proto.String("Object"),
				Field: fields,
			},
		},
	}

	message := buildDynamicMessage(t, file, "Object")
	descriptor := message.Descriptor()
	for i := 1; i <= fieldsCount; i++ {
		field := descriptor.Fields().ByName(protoreflect.Name(fmt.Sprintf("i%d", i)))
		if field == nil {
			t.Fatalf("missing dynamic field i%d", i)
		}
		message.Set(field, protoreflect.ValueOfInt32(int32(i)))
	}

	raw, err := protocache.Serialize(message.Interface())
	assert(t, err == nil)

	msg := protocache.AsMessage(raw)
	assert(t, msg.IsValid())

	var pool pcreflect.DescriptorPool
	assert(t, pool.Register(file))
	root := pool.Find("cov.Object")
	assert(t, root != nil)

	for i := 1; i <= fieldsCount; i++ {
		field := root.Lookup(fmt.Sprintf("i%d", i))
		if field == nil {
			t.Fatalf("missing reflected field i%d", i)
		}
		assert(t, field.Id() == uint16(i-1))
		unit := msg.GetField(uint16(i - 1))
		assert(t, unit.GetInt32() == int32(i))
	}
}

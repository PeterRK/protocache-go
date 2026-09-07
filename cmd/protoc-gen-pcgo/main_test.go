package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func TestDeprecatedDeclarationsAreExcluded(t *testing.T) {
	file := &descriptorpb.FileDescriptorProto{
		Name:    proto.String("deprecated.proto"),
		Package: proto.String("fixture"),
		Syntax:  proto.String("proto3"),
		Options: &descriptorpb.FileOptions{
			GoPackage: proto.String("example.com/fixture;fixture"),
		},
		EnumType: []*descriptorpb.EnumDescriptorProto{
			{
				Name: proto.String("ActiveEnum"),
				Value: []*descriptorpb.EnumValueDescriptorProto{
					{Name: proto.String("ACTIVE"), Number: proto.Int32(0)},
					{
						Name:    proto.String("OLD"),
						Number:  proto.Int32(1),
						Options: &descriptorpb.EnumValueOptions{Deprecated: proto.Bool(true)},
					},
				},
			},
			{
				Name:    proto.String("OldEnum"),
				Options: &descriptorpb.EnumOptions{Deprecated: proto.Bool(true)},
				Value: []*descriptorpb.EnumValueDescriptorProto{
					{Name: proto.String("OLD_ENUM_VALUE"), Number: proto.Int32(0)},
				},
			},
		},
		MessageType: []*descriptorpb.DescriptorProto{
			{
				Name: proto.String("Active"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{
						Name:   proto.String("active"),
						Number: proto.Int32(1),
						Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:   descriptorpb.FieldDescriptorProto_TYPE_INT32.Enum(),
					},
					{
						Name:    proto.String("old"),
						Number:  proto.Int32(2),
						Label:   descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:    descriptorpb.FieldDescriptorProto_TYPE_INT32.Enum(),
						Options: &descriptorpb.FieldOptions{Deprecated: proto.Bool(true)},
					},
				},
				NestedType: []*descriptorpb.DescriptorProto{
					{
						Name:    proto.String("OldNestedMessage"),
						Options: &descriptorpb.MessageOptions{Deprecated: proto.Bool(true)},
						Field: []*descriptorpb.FieldDescriptorProto{
							{
								Name:   proto.String("value"),
								Number: proto.Int32(1),
								Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
								Type:   descriptorpb.FieldDescriptorProto_TYPE_INT32.Enum(),
							},
						},
					},
				},
				EnumType: []*descriptorpb.EnumDescriptorProto{
					{
						Name:    proto.String("OldNestedEnum"),
						Options: &descriptorpb.EnumOptions{Deprecated: proto.Bool(true)},
						Value: []*descriptorpb.EnumValueDescriptorProto{
							{Name: proto.String("OLD_NESTED_VALUE"), Number: proto.Int32(0)},
						},
					},
				},
			},
			{
				Name:    proto.String("OldMessage"),
				Options: &descriptorpb.MessageOptions{Deprecated: proto.Bool(true)},
				Field: []*descriptorpb.FieldDescriptorProto{
					{
						Name:   proto.String("value"),
						Number: proto.Int32(1),
						Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						Type:   descriptorpb.FieldDescriptorProto_TYPE_INT32.Enum(),
					},
				},
			},
		},
	}
	request := &pluginpb.CodeGeneratorRequest{
		FileToGenerate: []string{file.GetName()},
		ProtoFile:      []*descriptorpb.FileDescriptorProto{file},
	}
	gen, err := (protogen.Options{}).New(request)
	if err != nil {
		t.Fatal(err)
	}

	oldTypeBook, oldAliasBook := typeBook, aliasBook
	typeBook = make(map[string]Type)
	aliasBook = make(map[string]Alias)
	defer func() {
		typeBook = oldTypeBook
		aliasBook = oldAliasBook
	}()

	for _, one := range gen.Files {
		CollectEnums(string(one.GoImportPath), one.Enums)
		CollectMessages(string(one.GoImportPath), one.Messages)
	}
	generated := gen.FilesByPath[file.GetName()]
	if err := GenFile(gen, generated, Options{Relative: true}); err != nil {
		t.Fatal(err)
	}
	if err := GenEXFile(gen, generated, Options{Relative: true}); err != nil {
		t.Fatal(err)
	}
	response := gen.Response()
	if response.GetError() != "" {
		t.Fatal(response.GetError())
	}

	var output strings.Builder
	for _, one := range response.File {
		output.WriteString(one.GetContent())
	}
	source := output.String()
	for _, want := range []string{
		"type ActiveEnum protocache.EnumValue",
		"ActiveEnum_ACTIVE",
		"type Active struct",
		"func (m *Active) GetActive() int32",
		"type ActiveEX struct",
	} {
		if !strings.Contains(source, want) {
			t.Errorf("generated output does not contain %q", want)
		}
	}
	for _, unwanted := range []string{
		"ActiveEnum_OLD",
		"OldEnum",
		"GetOld",
		"SetOld",
		"_FIELD_Active_old",
		"OldNestedMessage",
		"OldNestedEnum",
		"OldMessage",
	} {
		if strings.Contains(source, unwanted) {
			t.Errorf("generated output contains deprecated declaration %q", unwanted)
		}
	}
}

func TestGeneratedShortAliases(t *testing.T) {
	file := &descriptorpb.FileDescriptorProto{
		Name: proto.String("short_alias.proto"), Syntax: proto.String("proto3"),
		Package: proto.String("fixture"),
		Options: &descriptorpb.FileOptions{GoPackage: proto.String("example.com/fixture")},
	}
	holder := &descriptorpb.DescriptorProto{Name: proto.String("Holder")}
	for i, row := range []struct {
		name, field string
		kind        descriptorpb.FieldDescriptorProto_Type
	}{
		{"BoolRow", "bools", descriptorpb.FieldDescriptorProto_TYPE_BOOL},
		{"IntRow", "ints", descriptorpb.FieldDescriptorProto_TYPE_INT32},
		{"LongRow", "longs", descriptorpb.FieldDescriptorProto_TYPE_INT64},
		{"StringRow", "strings", descriptorpb.FieldDescriptorProto_TYPE_STRING},
	} {
		file.MessageType = append(file.MessageType, &descriptorpb.DescriptorProto{
			Name: proto.String(row.name), Field: []*descriptorpb.FieldDescriptorProto{{
				Name: proto.String("_"), Number: proto.Int32(1), Type: row.kind.Enum(),
				Label: descriptorpb.FieldDescriptorProto_LABEL_REPEATED.Enum(),
			}},
		})
		holder.Field = append(holder.Field, &descriptorpb.FieldDescriptorProto{
			Name: proto.String(row.field), Number: proto.Int32(int32(i + 1)),
			Type: descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(), TypeName: proto.String(".fixture." + row.name),
			Label: descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
		})
	}
	file.MessageType = append(file.MessageType, holder)
	gen, err := (protogen.Options{}).New(&pluginpb.CodeGeneratorRequest{
		FileToGenerate: []string{file.GetName()}, ProtoFile: []*descriptorpb.FileDescriptorProto{file},
	})
	if err != nil {
		t.Fatal(err)
	}
	oldTypes, oldAliases := typeBook, aliasBook
	typeBook, aliasBook = make(map[string]Type), make(map[string]Alias)
	defer func() { typeBook, aliasBook = oldTypes, oldAliases }()
	generated := gen.Files[0]
	CollectMessages(string(generated.GoImportPath), generated.Messages)
	if err := GenFile(gen, generated, Options{Relative: true}); err != nil {
		t.Fatal(err)
	}
	if err := GenEXFile(gen, generated, Options{Relative: true}); err != nil {
		t.Fatal(err)
	}
	response := gen.Response()
	if response.GetError() != "" {
		t.Fatal(response.GetError())
	}

	// Compile and execute the generated readers and EX writers against this runtime.
	dir := t.TempDir()
	args := []string{"test", "-count=1", "-gcflags=github.com/peterrk/protocache-go=-d=checkptr=2"}
	for _, out := range response.File {
		name := filepath.Join(dir, filepath.Base(out.GetName()))
		if err := os.WriteFile(name, []byte(out.GetContent()), 0600); err != nil {
			t.Fatal(err)
		}
		args = append(args, name)
	}
	const fixture = `package fixture

import (
	"encoding/binary"
	"fmt"
	"slices"
	"testing"
)

func TestGeneratedShortAliases(t *testing.T) {
	for size := 0; size <= 4; size++ {
		t.Run(fmt.Sprint(size), func(t *testing.T) {
			bools := BoolRowEX{false, true, false, true}[:size]
			ints := IntRowEX{-1, 0, 1, 2}[:size]
			longs := LongRowEX{-1, 1<<33 - 1, 2<<33 - 1, 3<<33 - 1}[:size]
			strings := StringRowEX{"x", "x", "x", "x"}[:size]
			for _, row := range []struct {
				name   string
				value  interface{ Serialize() ([]byte, error) }
				header uint32
				bytes  int
			}{
				{"bool", bools, uint32(size << 2), (1 + size + 3) / 4 * 4},
				{"int", ints, uint32(size<<2 | 1), 4 + size*4},
				{"long", longs, uint32(size<<2 | 2), 4 + size*8},
				{"string", strings, uint32(size<<2 | 1), 4 + size*4},
			} {
				t.Run(row.name, func(t *testing.T) {
					raw, err := row.value.Serialize()
					if err != nil {
						t.Fatal(err)
					}
					if len(raw) != row.bytes || uint32(raw[0]) != row.header {
						t.Fatalf("root encoding = %x, want %d bytes, header %x", raw, row.bytes, row.header)
					}
					if size == 0 && binary.LittleEndian.Uint32(raw) != row.header {
						t.Fatalf("wrong empty header: %x", raw)
					}
					switch row.name {
					case "bool":
						if !slices.Equal(TO_BoolRowEX(raw), bools) {
							t.Fatal("bool values lost")
						}
					case "int":
						if !slices.Equal(TO_IntRowEX(raw), ints) {
							t.Fatal("int values lost")
						}
					case "long":
						if !slices.Equal(TO_LongRowEX(raw), longs) {
							t.Fatal("long values lost")
						}
					case "string":
						if !slices.Equal(TO_StringRowEX(raw), strings) {
							t.Fatal("string values lost")
						}
					}
				})
			}

			holder := TO_HolderEX(nil)
			holder.SetBools(bools)
			holder.SetInts(ints)
			holder.SetLongs(longs)
			holder.SetStrings(strings)
			// Exercise new encoding, source replay, and encoding after getters.
			for pass := 0; pass < 3; pass++ {
				raw, err := holder.Serialize()
				if err != nil {
					t.Fatal(err)
				}
				view := TO_HolderEX(raw)
				if !slices.Equal(view.GetBools(), bools) || !slices.Equal(view.GetInts(), ints) ||
					!slices.Equal(view.GetLongs(), longs) || !slices.Equal(view.GetStrings(), strings) {
					t.Fatalf("nested values lost on pass %d: %x", pass, raw)
				}
				if pass == 0 {
					holder = TO_HolderEX(raw)
				} else {
					holder = view
				}
			}
		})
	}
}
`
	name := filepath.Join(dir, "short_alias_test.go")
	if err := os.WriteFile(name, []byte(fixture), 0600); err != nil {
		t.Fatal(err)
	}
	args = append(args, name)
	if output, err := exec.Command("go", args...).CombinedOutput(); err != nil {
		t.Fatalf("generated code tests failed: %v\n%s", err, output)
	}
}

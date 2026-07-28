package main

import (
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

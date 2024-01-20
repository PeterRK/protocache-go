package reflect

import (
	"math"
	"strings"

	pb "google.golang.org/protobuf/types/descriptorpb"
)

type FieldType uint8

const (
	TYPE_NONE    FieldType = 0
	TYPE_MESSAGE FieldType = 1
	TYPE_BYTES   FieldType = 2
	TYPE_STRING  FieldType = 3
	TYPE_FLOAT64 FieldType = 4
	TYPE_FLOAT32 FieldType = 5
	TYPE_UINT64  FieldType = 6
	TYPE_UINT32  FieldType = 7
	TYPE_INT64   FieldType = 8
	TYPE_INT32   FieldType = 9
	TYPE_BOOL    FieldType = 10
	TYPE_ENUM    FieldType = 11
	TYPE_UNKNOWN FieldType = 255
)

type Field struct {
	id              uint16
	repeated        bool
	key             FieldType
	value           FieldType
	valueType       string
	valueDescriptor *Descriptor
}

func (f *Field) Id() uint16 {
	return f.id
}

func (f *Field) IsRepeated() bool {
	return f.repeated
}

func (f *Field) Key() FieldType {
	return f.key
}

func (f *Field) Value() FieldType {
	return f.value
}

func (f *Field) ValueType() string {
	return f.valueType
}

func (f *Field) ValueDescriptor() *Descriptor {
	return f.valueDescriptor
}

func (f *Field) IsValid() bool {
	return f.value != TYPE_NONE
}

func (f *Field) IsMap() bool {
	return f.key != TYPE_NONE
}

type Descriptor struct {
	alias  Field
	fields map[string]*Field
}

func (d *Descriptor) Alias() *Field {
	if d.alias.IsValid() {
		return &d.alias
	}
	return nil
}

func (d *Descriptor) Lookup(name string) *Field {
	return d.fields[name]
}

func (d *Descriptor) Traverse(doit func(name string, field *Field) bool) {
	for k, v := range d.fields {
		if !doit(k, v) {
			break
		}
	}
}

type DescriptorPool struct {
	enum map[string]struct{}
	pool map[string]*Descriptor
}

func (p *DescriptorPool) Register(proto *pb.FileDescriptorProto) bool {
	if p.enum == nil {
		p.enum = make(map[string]struct{})
	}
	if p.pool == nil {
		p.pool = make(map[string]*Descriptor)
	}
	for _, one := range proto.GetEnumType() {
		if one.GetOptions() != nil && one.GetOptions().GetDeprecated() {
			continue
		}
		p.enum[clacFullname(proto.GetPackage(), one.GetName())] = struct{}{}
	}
	for _, one := range proto.GetMessageType() {
		if one.GetOptions() != nil && one.GetOptions().GetDeprecated() {
			continue
		}
		if !p.register(proto.GetPackage(), one) {
			return false
		}
	}
	for name, descriptor := range p.pool {
		if descriptor.alias.id != 0 && p.fixUnknownType(name, descriptor) {
			descriptor.alias.id = 0
		}
	}
	return true
}

func convertType(field *pb.FieldDescriptorProto) FieldType {
	if field.Type == nil {
		return TYPE_UNKNOWN
	}
	switch *field.Type {
	case pb.FieldDescriptorProto_TYPE_MESSAGE:
		return TYPE_MESSAGE
	case pb.FieldDescriptorProto_TYPE_BYTES:
		return TYPE_BYTES
	case pb.FieldDescriptorProto_TYPE_STRING:
		return TYPE_STRING
	case pb.FieldDescriptorProto_TYPE_DOUBLE:
		return TYPE_FLOAT64
	case pb.FieldDescriptorProto_TYPE_FLOAT:
		return TYPE_FLOAT32
	case pb.FieldDescriptorProto_TYPE_FIXED64,
		pb.FieldDescriptorProto_TYPE_UINT64:
		return TYPE_UINT64
	case pb.FieldDescriptorProto_TYPE_FIXED32,
		pb.FieldDescriptorProto_TYPE_UINT32:
		return TYPE_UINT32
	case pb.FieldDescriptorProto_TYPE_SFIXED64,
		pb.FieldDescriptorProto_TYPE_SINT64,
		pb.FieldDescriptorProto_TYPE_INT64:
		return TYPE_INT64
	case pb.FieldDescriptorProto_TYPE_SFIXED32,
		pb.FieldDescriptorProto_TYPE_SINT32,
		pb.FieldDescriptorProto_TYPE_INT32:
		return TYPE_INT32
	case pb.FieldDescriptorProto_TYPE_BOOL:
		return TYPE_BOOL
	case pb.FieldDescriptorProto_TYPE_ENUM:
		return TYPE_ENUM
	default:
		return TYPE_NONE
	}
}

func canBeKey(t FieldType) bool {
	switch t {
	case TYPE_STRING,
		TYPE_UINT64,
		TYPE_UINT32,
		TYPE_INT64,
		TYPE_INT32:
		return true
	default:
		return false
	}
}

func (p *DescriptorPool) register(ns string, proto *pb.DescriptorProto) bool {
	fullname := clacFullname(ns, proto.GetName())
	if p.pool[fullname] != nil {
		return false
	}

	for _, one := range proto.GetEnumType() {
		if one.GetOptions() != nil && one.GetOptions().GetDeprecated() {
			continue
		}
		p.enum[clacFullname(fullname, one.GetName())] = struct{}{}
	}

	mapEntries := make(map[string]*pb.DescriptorProto)
	for _, one := range proto.GetNestedType() {
		if options := one.GetOptions(); options != nil {
			if options.GetDeprecated() {
				continue
			}
			if options.GetMapEntry() {
				mapEntries[one.GetName()] = one
				continue
			}
		}
		if !p.register(fullname, one) {
			return false
		}
	}

	convertField := func(src *pb.FieldDescriptorProto, out *Field) bool {
		if src == nil {
			return false
		}
		out.repeated = src.GetLabel() == pb.FieldDescriptorProto_LABEL_REPEATED
		out.value = convertType(src)
		if out.value == TYPE_NONE {
			return false
		}
		if out.value == TYPE_MESSAGE || out.value == TYPE_UNKNOWN {
			entry := mapEntries[src.GetTypeName()]
			if entry != nil {
				out.key = convertType(entry.Field[0])
				out.value = convertType(entry.Field[1])
				if !canBeKey(out.key) || out.value == TYPE_NONE {
					return false
				}
				out.valueType = entry.Field[1].GetTypeName()
			} else {
				out.valueType = src.GetTypeName()
			}
		}
		return true
	}

	descriptor := &Descriptor{}
	descriptor.alias.id = math.MaxUint16
	if len(proto.Field) == 1 && proto.Field[0].GetName() == "_" {
		if !convertField(proto.Field[0], &descriptor.alias) {
			return false
		}
	} else {
		descriptor.fields = make(map[string]*Field, len(proto.Field))
		for _, one := range proto.Field {
			if one.GetOptions() != nil && one.GetOptions().GetDeprecated() {
				continue
			}
			if one.GetNumber() <= 0 {
				return false
			}
			field := &Field{
				id: uint16(one.GetNumber() - 1),
			}
			if !convertField(one, field) {
				return false
			}
			descriptor.fields[one.GetName()] = field
		}
	}
	p.pool[fullname] = descriptor
	return true
}

func (p *DescriptorPool) Find(fullname string) *Descriptor {
	descriptor := p.pool[fullname]
	if descriptor == nil {
		return nil
	}
	if descriptor.alias.id == 0 {
		return descriptor
	}
	if descriptor.alias.id != math.MaxUint16 {
		return nil
	}
	descriptor.alias.id--
	if !p.fixUnknownType(fullname, descriptor) {
		return nil
	}
	descriptor.alias.id = 0
	return descriptor
}

func (p *DescriptorPool) fixUnknownType(fullname string, descriptor *Descriptor) bool {
	bindType := func(name string, field *Field) bool {
		if _, found := p.enum[name]; found {
			field.value = TYPE_ENUM
			field.valueType = ""
			return true
		}
		if descriptor := p.pool[name]; descriptor != nil {
			field.value = TYPE_MESSAGE
			field.valueType = name
			field.valueDescriptor = descriptor
			return true
		}
		return false
	}

	checkType := func(field *Field) bool {
		if field.value != TYPE_UNKNOWN {
			return true
		}
		if len(field.valueType) == 0 {
			return false
		}
		if bindType(field.valueType, field) {
			return true
		}
		if bindType(fullname+"."+field.valueType, field) {
			return true
		}
		name := fullname
		for {
			pos := strings.LastIndexByte(name, '.')
			if pos < 0 {
				break
			}
			if bindType(name[:pos+1]+field.valueType, field) {
				return true
			}
			name = name[:pos]
		}
		return false
	}

	if field := descriptor.Alias(); field != nil {
		if !checkType(field) {
			return false
		}
	} else {
		for _, f := range descriptor.fields {
			if !checkType(f) {
				return false
			}
		}
	}
	return true
}

func clacFullname(ns, name string) string {
	if len(ns) == 0 {
		return name
	}
	return ns + "." + name
}

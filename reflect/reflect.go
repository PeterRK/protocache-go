package reflect

import (
	"errors"
	"fmt"
	"math"
	"strings"

	pb "google.golang.org/protobuf/types/descriptorpb"
)

type FieldType uint8

const (
	TypeNone    FieldType = 0
	TypeMessage FieldType = 1
	TypeBytes   FieldType = 2
	TypeString  FieldType = 3
	TypeFloat64 FieldType = 4
	TypeFloat32 FieldType = 5
	TypeUint64  FieldType = 6
	TypeUint32  FieldType = 7
	TypeInt64   FieldType = 8
	TypeInt32   FieldType = 9
	TypeBool    FieldType = 10
	TypeEnum    FieldType = 11
	TypeUnknown FieldType = 255
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
	return f.value != TypeNone
}

func (f *Field) IsMap() bool {
	return f.key != TypeNone
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

var (
	ErrDuplicateDescriptor = errors.New("duplicate descriptor")
	ErrInvalidField        = errors.New("invalid field")
	ErrInvalidFieldNumber  = errors.New("invalid field number")
	ErrInvalidMapKey       = errors.New("invalid map key type")
	ErrUnknownType         = errors.New("unknown type")
)

func (p *DescriptorPool) Register(proto *pb.FileDescriptorProto) error {
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
		p.enum[calcFullname(proto.GetPackage(), one.GetName())] = struct{}{}
	}
	for _, one := range proto.GetMessageType() {
		if one.GetOptions() != nil && one.GetOptions().GetDeprecated() {
			continue
		}
		if err := p.register(proto.GetPackage(), one); err != nil {
			return err
		}
	}
	for name, descriptor := range p.pool {
		if descriptor.alias.id == 0 {
			continue
		}
		if !p.fixUnknownType(name, descriptor) {
			return fmt.Errorf("%w: %s", ErrUnknownType, name)
		}
		descriptor.alias.id = 0
	}
	return nil
}

func convertType(field *pb.FieldDescriptorProto) FieldType {
	if field.Type == nil {
		return TypeUnknown
	}
	switch *field.Type {
	case pb.FieldDescriptorProto_TYPE_MESSAGE:
		return TypeMessage
	case pb.FieldDescriptorProto_TYPE_BYTES:
		return TypeBytes
	case pb.FieldDescriptorProto_TYPE_STRING:
		return TypeString
	case pb.FieldDescriptorProto_TYPE_DOUBLE:
		return TypeFloat64
	case pb.FieldDescriptorProto_TYPE_FLOAT:
		return TypeFloat32
	case pb.FieldDescriptorProto_TYPE_FIXED64,
		pb.FieldDescriptorProto_TYPE_UINT64:
		return TypeUint64
	case pb.FieldDescriptorProto_TYPE_FIXED32,
		pb.FieldDescriptorProto_TYPE_UINT32:
		return TypeUint32
	case pb.FieldDescriptorProto_TYPE_SFIXED64,
		pb.FieldDescriptorProto_TYPE_SINT64,
		pb.FieldDescriptorProto_TYPE_INT64:
		return TypeInt64
	case pb.FieldDescriptorProto_TYPE_SFIXED32,
		pb.FieldDescriptorProto_TYPE_SINT32,
		pb.FieldDescriptorProto_TYPE_INT32:
		return TypeInt32
	case pb.FieldDescriptorProto_TYPE_BOOL:
		return TypeBool
	case pb.FieldDescriptorProto_TYPE_ENUM:
		return TypeEnum
	default:
		return TypeNone
	}
}

func canBeKey(t FieldType) bool {
	switch t {
	case TypeString,
		TypeUint64,
		TypeUint32,
		TypeInt64,
		TypeInt32:
		return true
	default:
		return false
	}
}

func (p *DescriptorPool) register(ns string, proto *pb.DescriptorProto) error {
	fullname := calcFullname(ns, proto.GetName())
	if p.pool[fullname] != nil {
		return fmt.Errorf("%w: %s", ErrDuplicateDescriptor, fullname)
	}

	for _, one := range proto.GetEnumType() {
		if one.GetOptions() != nil && one.GetOptions().GetDeprecated() {
			continue
		}
		p.enum[calcFullname(fullname, one.GetName())] = struct{}{}
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
		if err := p.register(fullname, one); err != nil {
			return err
		}
	}

	convertField := func(src *pb.FieldDescriptorProto, out *Field) error {
		if src == nil {
			return ErrInvalidField
		}
		out.repeated = src.GetLabel() == pb.FieldDescriptorProto_LABEL_REPEATED
		out.value = convertType(src)
		if out.value == TypeNone {
			return fmt.Errorf("%w: %s", ErrInvalidField, src.GetName())
		}
		if out.value == TypeMessage || out.value == TypeUnknown {
			entry := mapEntries[src.GetTypeName()]
			if entry != nil {
				out.key = convertType(entry.Field[0])
				out.value = convertType(entry.Field[1])
				if !canBeKey(out.key) || out.value == TypeNone {
					return fmt.Errorf("%w: %s", ErrInvalidMapKey, src.GetName())
				}
				out.valueType = entry.Field[1].GetTypeName()
			} else {
				out.valueType = src.GetTypeName()
			}
		}
		return nil
	}

	descriptor := &Descriptor{}
	descriptor.alias.id = math.MaxUint16
	if len(proto.Field) == 1 && proto.Field[0].GetName() == "_" {
		if err := convertField(proto.Field[0], &descriptor.alias); err != nil {
			return fmt.Errorf("%w in %s._: %v", ErrInvalidField, fullname, err)
		}
	} else {
		descriptor.fields = make(map[string]*Field, len(proto.Field))
		for _, one := range proto.Field {
			if one.GetOptions() != nil && one.GetOptions().GetDeprecated() {
				continue
			}
			if one.GetNumber() <= 0 {
				return fmt.Errorf("%w: %s.%s", ErrInvalidFieldNumber, fullname, one.GetName())
			}
			field := &Field{
				id: uint16(one.GetNumber() - 1),
			}
			if err := convertField(one, field); err != nil {
				return fmt.Errorf("%w in %s.%s: %v", ErrInvalidField, fullname, one.GetName(), err)
			}
			descriptor.fields[one.GetName()] = field
		}
	}
	p.pool[fullname] = descriptor
	return nil
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
			field.value = TypeEnum
			field.valueType = ""
			return true
		}
		if descriptor := p.pool[name]; descriptor != nil {
			field.value = TypeMessage
			field.valueType = name
			field.valueDescriptor = descriptor
			return true
		}
		return false
	}

	checkType := func(field *Field) bool {
		if field.value != TypeUnknown {
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

func calcFullname(ns, name string) string {
	if len(ns) == 0 {
		return name
	}
	return ns + "." + name
}

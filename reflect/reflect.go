// Package reflect provides schema descriptors used to inspect ProtoCache
// messages dynamically.
//
// A DescriptorPool is populated from protobuf file descriptors. It is not safe
// for concurrent registration or lookup.
package reflect

import (
	"errors"
	"fmt"
	"math"
	"strings"

	pb "google.golang.org/protobuf/types/descriptorpb"
)

// FieldType identifies the ProtoCache representation of a protobuf field.
type FieldType uint8

const (
	// TypeNone denotes the absence of a field type.
	TypeNone FieldType = 0
	// TypeMessage denotes an embedded message.
	TypeMessage FieldType = 1
	// TypeBytes denotes a byte string.
	TypeBytes FieldType = 2
	// TypeString denotes a UTF-8 string.
	TypeString FieldType = 3
	// TypeFloat64 denotes a 64-bit floating-point value.
	TypeFloat64 FieldType = 4
	// TypeFloat32 denotes a 32-bit floating-point value.
	TypeFloat32 FieldType = 5
	// TypeUint64 denotes an unsigned 64-bit integer.
	TypeUint64 FieldType = 6
	// TypeUint32 denotes an unsigned 32-bit integer.
	TypeUint32 FieldType = 7
	// TypeInt64 denotes a signed 64-bit integer.
	TypeInt64 FieldType = 8
	// TypeInt32 denotes a signed 32-bit integer.
	TypeInt32 FieldType = 9
	// TypeBool denotes a boolean value.
	TypeBool FieldType = 10
	// TypeEnum denotes a protobuf enum value.
	TypeEnum FieldType = 11
	// TypeUnknown denotes a referenced type that has not been resolved.
	TypeUnknown FieldType = 255
)

// Field describes one protobuf field in ProtoCache terms.
type Field struct {
	id              uint16
	repeated        bool
	key             FieldType
	value           FieldType
	valueType       string
	valueDescriptor *Descriptor
}

// Id returns the zero-based ProtoCache field id, equal to the protobuf field number minus one.
func (f *Field) Id() uint16 {
	return f.id
}

// IsRepeated reports whether the field is repeated.
func (f *Field) IsRepeated() bool {
	return f.repeated
}

// Key returns the map key type, or TypeNone for a non-map field.
func (f *Field) Key() FieldType {
	return f.key
}

// Value returns the field value or map value type.
func (f *Field) Value() FieldType {
	return f.value
}

// ValueType returns the fully qualified protobuf type name for an unresolved or message value.
func (f *Field) ValueType() string {
	return f.valueType
}

// ValueDescriptor returns the descriptor of a message value, or nil for non-message values.
func (f *Field) ValueDescriptor() *Descriptor {
	return f.valueDescriptor
}

// IsValid reports whether f describes a supported field.
func (f *Field) IsValid() bool {
	return f.value != TypeNone
}

// IsMap reports whether f describes a map field.
func (f *Field) IsMap() bool {
	return f.key != TypeNone
}

// Descriptor describes a protobuf message or a ProtoCache alias.
type Descriptor struct {
	alias  Field
	fields map[string]*Field
}

// Alias returns the aliased field for an alias descriptor, or nil for a regular message.
func (d *Descriptor) Alias() *Field {
	if d.alias.IsValid() {
		return &d.alias
	}
	return nil
}

// Lookup returns the field named name, or nil when no such field exists.
func (d *Descriptor) Lookup(name string) *Field {
	return d.fields[name]
}

// Traverse visits fields in unspecified order until doit returns false.
func (d *Descriptor) Traverse(doit func(name string, field *Field) bool) {
	for k, v := range d.fields {
		if !doit(k, v) {
			break
		}
	}
}

// DescriptorPool stores descriptors registered from protobuf files and is not safe for concurrent use.
type DescriptorPool struct {
	enum map[string]struct{}
	pool map[string]*Descriptor
}

var (
	// ErrDuplicateDescriptor reports a duplicate fully qualified descriptor name.
	ErrDuplicateDescriptor = errors.New("duplicate descriptor")
	// ErrInvalidField reports an unsupported or malformed field declaration.
	ErrInvalidField = errors.New("invalid field")
	// ErrInvalidFieldNumber reports a non-positive protobuf field number.
	ErrInvalidFieldNumber = errors.New("invalid field number")
	// ErrInvalidMapKey reports a map key type unsupported by ProtoCache.
	ErrInvalidMapKey = errors.New("invalid map key type")
	// ErrUnknownType reports a referenced protobuf type that could not be resolved.
	ErrUnknownType = errors.New("unknown type")
)

// Register adds the non-deprecated declarations in proto and resolves their referenced types.
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

// Find returns the descriptor with the fully qualified protobuf name, or nil when absent or unresolved.
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

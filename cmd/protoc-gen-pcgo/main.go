package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/peterrk/slices"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func debugPrintln(args ...any) {
	fmt.Fprintln(os.Stderr, args...)
}

func debugPrintf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
}

const NoneKind = protoreflect.Kind(0)

type Alias struct {
	Key       protoreflect.Kind
	Value     protoreflect.Kind
	ValueType string
}

type Type struct {
	Package string
	GoName  string
}

var (
	typeBook  = make(map[string]Type)
	aliasBook = make(map[string]Alias)

	builinTypes = map[string]string{
		"Bytes":        "[]byte",
		"String":       "string",
		"Float64":      "float64",
		"Float32":      "float32",
		"Uint64":       "uint64",
		"Int64":        "int64",
		"Uint32":       "uint32",
		"Int32":        "int32",
		"Bool":         "bool",
		"Float64Array": "[]float64",
		"Float32Array": "[]float32",
		"Uint64Array":  "[]uint64",
		"Int64Array":   "[]int64",
		"Uint32Array":  "[]uint32",
		"Int32Array":   "[]int32",
		"BoolArray":    "[]bool",
	}
)

func (a *Alias) CalcName(imports map[string]string) (string, error) {
	if a.Key == NoneKind {
		switch a.Value {
		case protoreflect.MessageKind:
			if t, got := typeBook[a.ValueType]; !got {
				return "", fmt.Errorf("unknown type: %s", a.ValueType)
			} else if pkg, got := imports[t.Package]; got {
				return "ARRAY_" + pkg + t.GoName, nil
			} else {
				return "ARRAY_" + t.GoName, nil
			}
		case protoreflect.BytesKind:
			return "protocache.BytesArray", nil
		case protoreflect.StringKind:
			return "protocache.StringArray", nil
		case protoreflect.DoubleKind:
			return "protocache.Float64Array", nil
		case protoreflect.FloatKind:
			return "protocache.Float32Array", nil
		case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
			return "protocache.Uint64Array", nil
		case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
			return "protocache.Int64Array", nil
		case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
			return "protocache.Uint32Array", nil
		case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
			return "protocache.Int32Array", nil
		case protoreflect.BoolKind:
			return "protocache.BoolArray", nil
		case protoreflect.EnumKind:
			if t, got := typeBook[a.ValueType]; !got {
				return "", fmt.Errorf("unknown type: %s", a.ValueType)
			} else if pkg, got := imports[t.Package]; got {
				return fmt.Sprintf("protocache.EnumArray[%s.%s]", pkg, t.GoName), nil
			} else {
				return fmt.Sprintf("protocache.EnumArray[%s]", t.GoName), nil
			}
		default:
			return "", fmt.Errorf("unsupported type: %v", a.Value)
		}
	}

	prefix := "MAP_"
	switch a.Key {
	case protoreflect.StringKind:
		prefix += "string_"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		prefix += "uint64_"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		prefix += "int64_"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		prefix += "uint32_"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		prefix += "int32_"
	default:
		return "", fmt.Errorf("unsupported key type: %v", a.Key)
	}

	switch a.Value {
	case protoreflect.MessageKind,
		protoreflect.EnumKind:
		if t, got := typeBook[a.ValueType]; !got {
			return "", fmt.Errorf("unknown type: %s", a.ValueType)
		} else if pkg, got := imports[t.Package]; got {
			return prefix + pkg + t.GoName, nil
		} else {
			return prefix + t.GoName, nil
		}
	case protoreflect.BytesKind:
		return prefix + "bytes", nil
	case protoreflect.StringKind:
		return prefix + "string", nil
	case protoreflect.DoubleKind:
		return prefix + "float64", nil
	case protoreflect.FloatKind:
		return prefix + "float32", nil
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return prefix + "uint64", nil
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return prefix + "int64", nil
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return prefix + "uint32", nil
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return prefix + "int32", nil
	case protoreflect.BoolKind:
		return prefix + "bool", nil
	default:
		return "", fmt.Errorf("unsupported value type: %v", a.Value)
	}
}

func CollectEnums(pkg string, list []*protogen.Enum) {
	for _, one := range list {
		fullname := string(one.Desc.FullName())
		typeBook[fullname] = Type{
			Package: pkg,
			GoName:  one.GoIdent.GoName,
		}
	}
}

func CollectMessages(pkg string, list []*protogen.Message) {
	for _, one := range list {
		if one.Desc.IsMapEntry() {
			continue
		}
		fullname := string(one.Desc.FullName())
		CollectEnums(pkg, one.Enums)
		CollectMessages(pkg, one.Messages)
		typeBook[fullname] = Type{
			Package: pkg,
			GoName:  one.GoIdent.GoName,
		}
		if len(one.Fields) != 1 || one.Fields[0].Desc.Name() != "_" {
			continue
		}

		var alias Alias
		field := one.Fields[0].Desc
		if field.IsMap() {
			alias.Key = field.MapKey().Kind()
			field = field.MapValue()
		} else if field.IsList() {
			alias.Key = NoneKind
		} else {
			continue
		}
		alias.Value = field.Kind()
		if alias.Value == protoreflect.EnumKind {
			alias.ValueType = string(field.Enum().FullName())
		} else if alias.Value == protoreflect.MessageKind {
			alias.ValueType = string(field.Message().FullName())
		}
		aliasBook[fullname] = alias
	}
}

func main() {
	protogen.Options{}.Run(func(gen *protogen.Plugin) error {
		for _, file := range gen.Files {
			CollectEnums(string(file.GoImportPath), file.Enums)
			CollectMessages(string(file.GoImportPath), file.Messages)
		}

		//debugPrintln(typeBook)
		//debugPrintln(aliasBook)

		for _, file := range gen.Files {
			if !file.Generate {
				continue
			}
			err := GenFile(gen, file)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func CollectImports(pkg string, list []*protogen.Message, book map[string]string) {
	for _, one := range list {
		if one.Desc.IsMapEntry() {
			continue
		}
		CollectImports(pkg, one.Messages, book)
		for _, field := range one.Fields {
			desc := field.Desc
			if desc.IsMap() {
				desc = desc.MapValue()
			}
			var valueType string
			if desc.Kind() == protoreflect.EnumKind {
				valueType = string(desc.Enum().FullName())
			} else if desc.Kind() == protoreflect.MessageKind {
				valueType = string(desc.Message().FullName())
			} else {
				continue
			}
			t, got := typeBook[valueType]
			if got && t.Package != pkg {
				book[t.Package] = ""
			}
		}
	}
}

func GenEnums(g *protogen.GeneratedFile, list []*protogen.Enum) error {
	for _, one := range list {
		g.P()
		g.P("type ", one.GoIdent.GoName, " protocache.EnumValue")
		g.P()
		g.P("const (")
		for _, v := range one.Values {
			g.P("	", v.GoIdent.GoName, " ", one.GoIdent.GoName, " = ", v.Desc.Number())
		}
		g.P(")")
	}
	return nil
}

func GenMessages(g *protogen.GeneratedFile, imports map[string]string,
	derived map[string]Alias, list []*protogen.Message) error {
	getPkgPrefix := func(fullname string) string {
		if t, got := typeBook[fullname]; got {
			if pkg, got := imports[t.Package]; got {
				return pkg + "."
			}
		}
		return ""
	}

	for _, one := range list {
		if one.Desc.IsMapEntry() {
			continue
		}
		GenEnums(g, one.Enums)
		GenMessages(g, imports, derived, one.Messages)

		alias, ok := aliasBook[string(one.Desc.FullName())]
		if ok {
			name, err := alias.CalcName(imports)
			if err != nil {
				return err
			}
			derived[name] = alias
			g.P()
			g.P("type ", one.GoIdent.GoName, " = ", name)
			g.P()
			g.P("func AS_", one.GoIdent.GoName, "(data []byte)", one.GoIdent.GoName, " {")
			if strings.HasPrefix(name, "protocache.") {
				g.P("	return protocache.As", name[11:], "(data)")
			} else {
				g.P("	return AS_", name, "(data)")
			}
			g.P("}")
			continue
		}

		g.P()
		g.P("const (")
		for _, field := range one.Fields {
			g.P("	_FIELD_", one.GoIdent.GoName, "_", field.Desc.Name(), " uint16 = ", field.Desc.Number()-1)
		}
		g.P(")")
		g.P()
		g.P("type ", one.GoIdent.GoName, " struct { core protocache.Message }")
		g.P()
		g.P("func AS_", one.GoIdent.GoName, "(data []byte) ", one.GoIdent.GoName,
			"{ return ", one.GoIdent.GoName, "{core: protocache.AsMessage(data)}}")
		g.P()
		g.P("func (m *", one.GoIdent.GoName, ") IsValid() bool {return m.core.IsValid()}")

		handleSimpleField := func(field *protogen.Field, t string) {
			g.P("func (m *", one.GoIdent.GoName, ") Get", field.GoName, "() ", builinTypes[t], " {")
			g.P("	field := m.core.GetField(_FIELD_", one.GoIdent.GoName, "_", field.Desc.Name(), ")")
			g.P("	return field.Get", t, "()")
			g.P("}")
		}

		handleStingrArray := func(field *protogen.Field, t string) {
			g.P("func (m *", one.GoIdent.GoName, ") Get", field.GoName, "() protocache.", t, "Array {")
			g.P("	field := m.core.GetField(_FIELD_", one.GoIdent.GoName, "_", field.Desc.Name(), ")")
			g.P("	return protocache.As", t, "Array(field.GetObject())")
			g.P("}")
		}

		handleComlexField := func(field *protogen.Field, t string) {
			g.P("func (m *", one.GoIdent.GoName, ") Get", field.GoName, "() ", t, " {")
			g.P("	field := m.core.GetField(_FIELD_", one.GoIdent.GoName, "_", field.Desc.Name(), ")")
			g.P("	return AS_", t, "(field.GetObject())")
			g.P("}")
		}

		for _, field := range one.Fields {
			desc := field.Desc
			g.P()
			if desc.IsMap() {
				vDesc := desc.MapValue()
				alias := Alias{
					Key:   desc.MapKey().Kind(),
					Value: vDesc.Kind(),
				}
				if alias.Value == protoreflect.EnumKind {
					alias.ValueType = string(vDesc.Enum().FullName())
				} else if alias.Value == protoreflect.MessageKind {
					alias.ValueType = string(vDesc.Message().FullName())
				}
				name, err := alias.CalcName(imports)
				if err != nil {
					return err
				}
				derived[name] = alias
				handleComlexField(field, name)
			} else if desc.IsList() {
				switch desc.Kind() {
				case protoreflect.MessageKind:
					alias := Alias{
						Key:       NoneKind,
						Value:     desc.Kind(),
						ValueType: string(desc.Message().FullName()),
					}
					name, err := alias.CalcName(imports)
					if err != nil {
						return err
					}
					derived[name] = alias
					handleComlexField(field, name)
				case protoreflect.BytesKind:
					handleStingrArray(field, "Bytes")
				case protoreflect.StringKind:
					handleStingrArray(field, "String")
				case protoreflect.DoubleKind:
					handleSimpleField(field, "Float64Array")
				case protoreflect.FloatKind:
					handleSimpleField(field, "Float32Array")
				case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
					handleSimpleField(field, "Uint64Array")
				case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
					handleSimpleField(field, "Int64Array")
				case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
					handleSimpleField(field, "Uint32Array")
				case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
					handleSimpleField(field, "Int32Array")
				case protoreflect.BoolKind:
					handleSimpleField(field, "BoolArray")
				case protoreflect.EnumKind:
					pkgPrefix := getPkgPrefix(string(field.Enum.Desc.FullName()))
					typeName := pkgPrefix + field.Enum.GoIdent.GoName
					g.P("func (m *", one.GoIdent.GoName, ") Get", field.GoName, "() []", typeName, " {")
					g.P("	field := m.core.GetField(_FIELD_", one.GoIdent.GoName, "_", field.Desc.Name(), ")")
					g.P("	return field.GetEnumArray[", typeName, "]()")
					g.P("}")
				default:
					return fmt.Errorf("unsupported type: %v", desc.Kind())
				}
			} else {
				switch desc.Kind() {
				case protoreflect.MessageKind:
					pkgPrefix := getPkgPrefix(string(field.Message.Desc.FullName()))
					typeName := field.Message.GoIdent.GoName
					g.P("func (m *", one.GoIdent.GoName, ") Get", field.GoName, "() ", pkgPrefix, typeName, " {")
					g.P("	field := m.core.GetField(_FIELD_", one.GoIdent.GoName, "_", field.Desc.Name(), ")")
					g.P("	return ", pkgPrefix, "AS_", typeName, "(field.GetObject())")
					g.P("}")
				case protoreflect.BytesKind:
					handleSimpleField(field, "Bytes")
				case protoreflect.StringKind:
					handleSimpleField(field, "String")
				case protoreflect.DoubleKind:
					handleSimpleField(field, "Float64")
				case protoreflect.FloatKind:
					handleSimpleField(field, "Float32")
				case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
					handleSimpleField(field, "Uint64")
				case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
					handleSimpleField(field, "Int64")
				case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
					handleSimpleField(field, "Uint32")
				case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
					handleSimpleField(field, "Int32")
				case protoreflect.BoolKind:
					handleSimpleField(field, "Bool")
				case protoreflect.EnumKind:
					pkgPrefix := getPkgPrefix(string(field.Enum.Desc.FullName()))
					typeName := pkgPrefix + field.Enum.GoIdent.GoName
					g.P("func (m *", one.GoIdent.GoName, ") Get", field.GoName, "() ", typeName, " {")
					g.P("	field := m.core.GetField(_FIELD_", one.GoIdent.GoName, "_", field.Desc.Name(), ")")
					g.P("	return ", typeName, "(field.GetEnumValue())")
					g.P("}")
				default:
					return fmt.Errorf("unsupported type: %v", desc.Kind())
				}
			}
		}
	}
	return nil
}

func GenFile(gen *protogen.Plugin, file *protogen.File) error {
	imports := make(map[string]string)
	CollectImports(string(file.GoImportPath), file.Messages, imports)

	pkgs := make([]string, 0, len(imports))
	for one, _ := range imports {
		pkgs = append(pkgs, one)
	}
	slices.Sort(pkgs)
	for i, one := range pkgs {
		imports[one] = fmt.Sprintf("p%d", i+1)
	}

	filename := file.GeneratedFilenamePrefix + ".pc.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P("package ", file.GoPackageName)
	g.P()
	g.P("import (")
	g.P(`	"github.com/peterrk/protocache-go"`)
	for name, mark := range imports {
		g.P("	", mark, ` "`, name, `"`)
	}
	g.P(")")

	if err := GenEnums(g, file.Enums); err != nil {
		return err
	}

	derived := make(map[string]Alias)
	if err := GenMessages(g, imports, derived, file.Messages); err != nil {
		return err
	}

	names := make([]string, 0, len(derived))
	for one, _ := range derived {
		if !strings.HasPrefix(one, "protocache.") {
			names = append(names, one)
		}
	}
	slices.Sort(names)
	//debugPrintln(names)

	handleSimpleKV := func(name, kt, vt string) {
		g.P("func (x *", name, ") Value(i uint32) ", builinTypes[vt], " {")
		g.P("	field := x.core.Value(i)")
		g.P("	return field.Get", vt, "()")
		g.P("}")
		g.P()
		g.P("func (x *", name, ") Find(key ", builinTypes[kt], ") (", builinTypes[vt], ", bool) {")
		g.P("	field := x.core.FindBy", kt, "(key)")
		g.P("	return field.Get", vt, "(), field.IsValid()")
		g.P("}")
	}

	getValueType := func(fullname string) (prefix string, name string) {
		t := typeBook[fullname]
		if pkg, got := imports[t.Package]; got {
			return pkg + ".", t.GoName
		} else {
			return "", t.GoName
		}
	}
	for _, name := range names {
		g.P()
		alias := derived[name]
		if alias.Key == NoneKind {
			if alias.Value != protoreflect.MessageKind {
				panic("message array only")
			}
			g.P("type ", name, " struct { core protocache.Array }")
			g.P()
			g.P("func AS_", name, "(data []byte) ", name, " {")
			g.P("	return ", name, "{core: protocache.AsArray(data)}")
			g.P("}")
			g.P()
			prefix, value := getValueType(alias.ValueType)
			g.P("func (x *", name, ") Get(i uint32) ", prefix, value, " {")
			g.P("	field := x.core.Get(i)")
			g.P("	return ", prefix, "AS_", value, "(field.GetObject())")
			g.P("}")
		} else {
			g.P("type ", name, " struct { core protocache.Map }")
			g.P()
			g.P("func AS_", name, "(data []byte) ", name, " {")
			g.P("	return ", name, "{core: protocache.AsMap(data)}")
			g.P("}")

			var key string
			switch alias.Key {
			case protoreflect.StringKind:
				key = "String"
			case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
				key = "Uint64"
			case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
				key = "Int64"
			case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
				key = "Uint32"
			case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
				key = "Int32"
			default:
				panic("unsupported key type")
			}

			g.P()
			g.P("func (x *", name, ") Key(i uint32) ", builinTypes[key], " {")
			g.P("	field := x.core.Key(i)")
			g.P("	return field.Get", key, "()")
			g.P("}")
			g.P()

			switch alias.Value {
			case protoreflect.MessageKind:
				prefix, value := getValueType(alias.ValueType)
				g.P("func (x *", name, ") Value(i uint32) ", prefix, value, " {")
				g.P("	field := x.core.Value(i)")
				g.P("	return ", prefix, "AS_", value, "(field.GetObject())")
				g.P("}")
				g.P()
				g.P("func (x *", name, ") Find(key ", builinTypes[key], ") (", prefix, value, ", bool) {")
				g.P("	field := x.core.FindBy", key, "(key)")
				g.P("	return ", prefix, "AS_", value, "(field.GetObject()), field.IsValid()")
				g.P("}")
			case protoreflect.BytesKind:
				handleSimpleKV(name, key, "Bytes")
			case protoreflect.StringKind:
				handleSimpleKV(name, key, "String")
			case protoreflect.DoubleKind:
				handleSimpleKV(name, key, "Float64")
			case protoreflect.FloatKind:
				handleSimpleKV(name, key, "Float32")
			case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
				handleSimpleKV(name, key, "Uint64")
			case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
				handleSimpleKV(name, key, "Int64")
			case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
				handleSimpleKV(name, key, "Uint32")
			case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
				handleSimpleKV(name, key, "Int32")
			case protoreflect.BoolKind:
				handleSimpleKV(name, key, "Bool")
			case protoreflect.EnumKind:
				prefix, value := getValueType(alias.ValueType)
				g.P("func (x *", name, ") Value(i uint32) ", prefix, value, " {")
				g.P("	field := x.core.Value(i)")
				g.P("	return ", prefix, value, "(field.GetEnumValue())")
				g.P("}")
				g.P()
				g.P("func (x *", name, ") Find(key ", builinTypes[key], ") (", prefix, value, ", bool) {")
				g.P("	field := x.core.FindBy", key, "(key)")
				g.P("	return ", prefix, value, "(field.GetEnumValue()), field.IsValid()")
				g.P("}")
			default:
				panic("unsupported value type")
			}
		}
		g.P()
		g.P("func (x *", name, ") IsValid() bool {return x.core.IsValid()}")
		g.P()
		g.P("func (x *", name, ") Size() uint32 {return x.core.Size()}")
	}
	return nil
}

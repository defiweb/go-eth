package abi

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/defiweb/go-anymapper"
	"github.com/defiweb/go-sigparser"
)

// Type represents a type that can be marshaled to and from ABI.
//
// https://docs.soliditylang.org/en/develop/abi-spec.html#strict-encoding-mode
type Type interface {
	// DynamicType indicates whether the type is dynamic.
	DynamicType() bool
	// EncodeABI returns the ABI encoding of the value.
	EncodeABI() (Words, error)
	// DecodeABI sets the value from the ABI encoding.
	DecodeABI(Words) (int, error)
}

func ParseType(signature string) (TypeDefinition, error) {
	return DefaultConfig.ParseType(signature)
}

func NewType(param TypeDefinition) (Type, error) {
	return DefaultConfig.NewType(param)
}

type TypeDefinition struct {
	Name          string
	Type          string
	CanonicalType string
	Arrays        []int
	Tuple         []TypeDefinition
}

var DefaultConfig *Config

type Config struct {
	mapper *anymapper.Mapper
	types  map[string]*typ
}

type typ struct {
	factory   func() Type
	canonical string
}

func (c *Config) SetMapper(mapper *anymapper.Mapper) {
	c.mapper = mapper
}

func (c *Config) RegisterType(name string, canonical string, factory func() Type) {
	c.types[name] = &typ{
		canonical: canonical,
		factory:   factory,
	}
}

func (c *Config) ParseType(signature string) (TypeDefinition, error) {
	param, err := sigparser.ParseParameter(signature)
	if err != nil {
		return TypeDefinition{}, err
	}
	return c.convertParam(param), nil
}

func (c *Config) ParseMethod(signature string) (*Method, error) {
	sig, err := sigparser.ParseSignature(signature)
	if err != nil {
		return nil, err
	}
	return &Method{
		Name:    sig.Name,
		Inputs:  c.convertParamList(sig.Inputs),
		Outputs: c.convertParamList(sig.Outputs),
		Config:  c,
	}, nil
}

func (c *Config) NewType(def TypeDefinition) (Type, error) {
	var typ Type
	switch {
	case len(def.Arrays) > 0:
		cpy := def.Copy()
		cpy.Arrays = cpy.Arrays[1:]
		if def.Arrays[0] < 0 {
			typ = NewArray(cpy)
		} else {
			typ = NewFixedArray(cpy, def.Arrays[0])
		}
	case len(def.Tuple) > 0:
		tuple := NewTupleOfSize(len(def.Tuple))
		for idx, elem := range def.Tuple {
			typ, err := c.NewType(elem)
			if err != nil {
				return nil, err
			}
			name := elem.Name
			if name == "" {
				name = strconv.Itoa(idx)
			}
			if err := tuple.Set(idx, name, typ); err != nil {
				return nil, err
			}
		}
		typ = tuple
	default:
		if t, ok := c.types[def.Type]; ok {
			typ = t.factory()
		}
	}
	if typ == nil {
		return nil, fmt.Errorf("unknown type %q", def.Type)
	}
	if typ, ok := typ.(configAware); ok {
		typ.SetConfig(c)
	}
	return typ, nil
}

func (c *Config) NewTypeList(defs []TypeDefinition) (*TupleType, error) {
	tuple := NewTupleOfSize(len(defs))
	for idx, def := range defs {
		typ, err := c.NewType(def)
		if err != nil {
			return nil, err
		}
		name := def.Name
		if name == "" {
			name = strconv.Itoa(idx)
		}
		if err := tuple.Set(idx, name, typ); err != nil {
			return nil, err
		}
	}
	return tuple, nil
}

func (c *Config) convertParam(param sigparser.Parameter) TypeDefinition {
	def := TypeDefinition{
		Name:          param.Name,
		Type:          param.Type,
		CanonicalType: c.canonicalizeParam(param),
	}
	def.Tuple = make([]TypeDefinition, len(param.Tuple))
	def.Arrays = make([]int, len(param.Arrays))
	for i, t := range param.Tuple {
		def.Tuple[i] = c.convertParam(t)
	}
	copy(def.Arrays, param.Arrays)
	return def
}

func (c *Config) convertParamList(params []sigparser.Parameter) []TypeDefinition {
	defs := make([]TypeDefinition, len(params))
	for i, param := range params {
		defs[i] = c.convertParam(param)
	}
	return defs
}

func (c *Config) canonicalizeParam(param sigparser.Parameter) string {
	var buf strings.Builder
	if len(param.Type) > 0 {
		if t, ok := c.types[param.Type]; ok {
			buf.WriteString(t.canonical)
		} else {
			buf.WriteString(param.Type)
		}
	} else {
		buf.WriteByte('(')
		for i, elem := range param.Tuple {
			if i > 0 {
				buf.WriteString(",")
			}
			buf.WriteString(c.canonicalizeParam(elem))
		}
		buf.WriteByte(')')
	}
	for _, arr := range param.Arrays {
		if arr < 0 {
			buf.WriteString("[]")
		} else {
			buf.WriteByte('[')
			buf.WriteString(strconv.Itoa(arr))
			buf.WriteByte(']')
		}
	}
	return buf.String()
}

func (c *Config) Copy() *Config {
	cpy := &Config{
		types:  make(map[string]*typ),
		mapper: c.mapper,
	}
	for k, v := range c.types {
		cpy.types[k] = v
	}
	return cpy
}

func (t TypeDefinition) String() string {
	var buf strings.Builder
	if len(t.Type) > 0 {
		buf.WriteString(t.Type)
	} else {
		buf.WriteByte('(')
		for i, elem := range t.Tuple {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(elem.String())
		}
		buf.WriteByte(')')
	}
	for _, arr := range t.Arrays {
		if arr < 0 {
			buf.WriteString("[]")
		} else {
			buf.WriteByte('[')
			buf.WriteString(strconv.Itoa(arr))
			buf.WriteByte(']')
		}
	}
	return buf.String()
}

func (t TypeDefinition) Copy() TypeDefinition {
	cpy := TypeDefinition{
		Name:          t.Name,
		Type:          t.Type,
		CanonicalType: t.CanonicalType,
	}
	cpy.Tuple = make([]TypeDefinition, len(t.Tuple))
	cpy.Arrays = make([]int, len(t.Arrays))
	for i, t := range t.Tuple {
		cpy.Tuple[i] = t.Copy()
	}
	copy(cpy.Arrays, t.Arrays)
	return cpy
}

type configAware interface {
	SetConfig(*Config)
}

func init() {
	m := anymapper.DefaultMapper.Copy()
	m.Tag = "abi"
	m.FieldMapper = func(name string) string {
		if len(name) > 0 {
			name = strings.ToLower(name[:1]) + name[1:]
		}
		return name
	}
	DefaultConfig = &Config{
		types:  map[string]*typ{},
		mapper: m,
	}
	DefaultConfig.RegisterType("uint", "uint256", func() Type { return NewUint(32) })
	DefaultConfig.RegisterType("int", "int256", func() Type { return NewInt(32) })
	DefaultConfig.RegisterType("bool", "bool", func() Type { return NewBool() })
	DefaultConfig.RegisterType("bytes", "bytes", func() Type { return NewBytes() })
	DefaultConfig.RegisterType("string", "string", func() Type { return NewString() })
	DefaultConfig.RegisterType("address", "address", func() Type { return NewAddress() })
	for i := 0; i <= 32; i++ {
		i := i
		uintType := fmt.Sprintf("uint%d", i*8)
		intType := fmt.Sprintf("int%d", i*8)
		bytesType := fmt.Sprintf("bytes%d", i)
		DefaultConfig.RegisterType(uintType, uintType, func() Type { return NewUint(i) })
		DefaultConfig.RegisterType(intType, intType, func() Type { return NewInt(i) })
		DefaultConfig.RegisterType(bytesType, bytesType, func() Type { return NewFixedBytes(i) })
	}
}

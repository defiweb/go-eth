package abi

import (
	"fmt"
	"strings"
)

// Type is a representation of a type like uint256 or address. The type can be
// used to create a new value of that type, but it cannot store a value.
type Type interface {
	// IsDynamic indicates whether the type is dynamic.
	IsDynamic() bool

	// CanonicalType returns the canonical name of the type. In case of a
	// tuple, the canonical name is the canonical name of the tuple's
	// elements, separated by commas and enclosed in parentheses. Arrays
	// are represented by the canonical name of the element type followed
	// by square brackets with the array size.
	CanonicalType() string

	// String returns the user-friendly name of the type.
	String() string

	// Value creates a new zero value for the type.
	Value() Value
}

// ParseType parses a type signature and returns a new Type.
//
// A type can be either an elementary type like uint256 or a tuple type. Tuple
// types are denoted by parentheses, with the optional keyword "tuple" before
// the parentheses. Parameter names are optional.
//
// The generated types can be used to create new values, which can then be used
// to encode or decode ABI data.
//
// Custom types may be added to the ABI.Types, this will allow the parser to
// handle them.
//
// The following examples are valid type signatures:
//
//	uint256
//	(uint256 a,bytes32 b)
//	tuple(uint256 a, bytes32 b)[]
//
// This function is equivalent to calling Parser.ParseType with the default
// configuration.
func ParseType(signature string) (Type, error) {
	return Default.ParseType(signature)
}

// ParseStruct parses a struct definition and returns a new Type.
//
// It is similar to ParseType, but accepts a struct definition instead of a
// type signature.
//
// For example, the following two calls are equivalent:
//
//	ParseType("(uint256 a, bytes32 b)")
//	ParseStruct("struct { uint256 a; bytes32 b; }")
func ParseStruct(definition string) (Type, error) {
	return Default.ParseStruct(definition)
}

// MustParseType is like ParseType but panics on error.
func MustParseType(signature string) Type {
	t, err := ParseType(signature)
	if err != nil {
		panic(err)
	}
	return t
}

// MustParseStruct is like ParseStruct but panics on error.
func MustParseStruct(definition string) Type {
	t, err := ParseStruct(definition)
	if err != nil {
		panic(err)
	}
	return t
}

// ParseType parses a type signature and returns a new Type.
//
// See ParseType for more information.
func (a *ABI) ParseType(signature string) (Type, error) {
	return parseType(a, nil, signature)
}

// ParseStruct parses a struct definition and returns a new Type.
//
// See ParseStruct for more information.
func (a *ABI) ParseStruct(definition string) (Type, error) {
	return parseStruct(a, nil, definition)
}

// AliasType wraps another type and gives it a different type name. The canonical
// type name is the same as the wrapped type.
type AliasType struct {
	alias string
	typ   Type
}

// NewAliasType creates a new alias type.
func NewAliasType(alias string, typ Type) *AliasType {
	return &AliasType{alias: alias, typ: typ}
}

// Type returns the aliased type.
func (a *AliasType) Type() Type {
	return a.typ
}

// IsDynamic implements the Type interface.
func (a *AliasType) IsDynamic() bool {
	return a.typ.IsDynamic()
}

// CanonicalType implements the Type interface.
func (a *AliasType) CanonicalType() string {
	return a.typ.CanonicalType()
}

// String implements the Type interface.
func (a *AliasType) String() string {
	return a.alias
}

// Value implements the Type interface.
func (a *AliasType) Value() Value {
	return a.typ.Value()
}

// TupleType represents a tuple type.
type TupleType struct {
	elems []TupleTypeElem
}

// TupleTypeElem is an element of a tuple.
type TupleTypeElem struct {
	// Name of the tuple element. It is used when mapping values from and to
	// maps and structures. If the name is empty, when creating a new value
	// the name will be set to argN, where N is the index of the element.
	Name string

	// Type is the type of the element.
	Type Type
}

// NewTupleType creates a new tuple type with the given elements.
func NewTupleType(elems ...TupleTypeElem) *TupleType {
	return &TupleType{elems: elems}
}

// Size returns the number of elements in the tuple.
func (t *TupleType) Size() int {
	return len(t.elems)
}

// Elements returns the tuple elements.
func (t *TupleType) Elements() []TupleTypeElem {
	cpy := make([]TupleTypeElem, len(t.elems))
	copy(cpy, t.elems)
	return cpy
}

// IsDynamic implements the Type interface.
func (t *TupleType) IsDynamic() bool {
	for _, elem := range t.elems {
		if elem.Type.IsDynamic() {
			return true
		}
	}
	return false
}

// CanonicalType implements the Type interface.
func (t *TupleType) CanonicalType() string {
	var buf strings.Builder
	buf.WriteString("(")
	for i, elem := range t.elems {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(elem.Type.CanonicalType())
	}
	buf.WriteString(")")
	return buf.String()
}

// String implements the Type interface.
func (t *TupleType) String() string {
	var buf strings.Builder
	buf.WriteString("(")
	for i, elem := range t.elems {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(elem.Type.String())
		if len(elem.Name) > 0 {
			buf.WriteString(" ")
			buf.WriteString(elem.Name)
		}
	}
	buf.WriteString(")")
	return buf.String()
}

// Value implements the Type interface.
func (t *TupleType) Value() Value {
	v := make(TupleValue, len(t.elems))
	for i, elem := range t.elems {
		v[i] = TupleValueElem{
			Name:  elem.Name,
			Value: elem.Type.Value(),
		}
		if len(elem.Name) == 0 {
			v[i].Name = fmt.Sprintf("arg%d", i)
		}
	}
	return &v
}

// EventTupleType represents a tuple type for event inputs. It works just like
// TupleType, but elements can be marked as indexed. When creating a new value,
// the indexed elements will be created first, followed by the non-indexed
// elements.
type EventTupleType struct {
	elems   []EventTupleElem
	indexed int
}

// EventTupleElem is an element of an event tuple.
type EventTupleElem struct {
	// Name of the tuple element. It is used when mapping values from and to
	// maps and structures. If the name is empty, when creating a new value,
	// the name will be set to topicN or dataN, where N is the index of the
	// topic or data element. Topics are counted from 1 because the first topic
	// is the event signature.
	Name string

	// Indexed indicates whether the element is indexed.
	Indexed bool

	// Type is the type of the element.
	Type Type
}

// NewEventTupleType creates a new tuple type with the given elements.
func NewEventTupleType(elems ...EventTupleElem) *EventTupleType {
	indexed := 0
	for _, elem := range elems {
		if elem.Indexed {
			indexed++
		}
	}
	return &EventTupleType{elems: elems, indexed: indexed}
}

// Size returns the number of elements in the tuple.
func (t *EventTupleType) Size() int {
	return len(t.elems)
}

// IndexedSize returns the number of indexed elements in the tuple.
func (t *EventTupleType) IndexedSize() int {
	return t.indexed
}

// DataSize returns the number of non-indexed elements in the tuple.
func (t *EventTupleType) DataSize() int {
	return len(t.elems) - t.indexed
}

// Elements returns the tuple elements.
func (t *EventTupleType) Elements() []EventTupleElem {
	cpy := make([]EventTupleElem, len(t.elems))
	copy(cpy, t.elems)
	return cpy
}

// TopicsTuple returns the tuple of indexed arguments.
//
// If the type is indexed and dynamic, the type will be converted to bytes32.
func (t *EventTupleType) TopicsTuple() *TupleType {
	topics := make([]TupleTypeElem, 0, t.indexed)
	for _, elem := range t.elems {
		if !elem.Indexed {
			continue
		}
		name := elem.Name
		if len(name) == 0 {
			name = fmt.Sprintf("topic%d", len(topics))
		}
		typ := elem.Type
		if typ.IsDynamic() {
			typ = &FixedBytesType{size: 32}
		}
		topics = append(topics, TupleTypeElem{
			Name: name,
			Type: typ,
		})
	}
	return &TupleType{elems: topics}
}

// DataTuple returns the tuple of non-indexed arguments.
func (t *EventTupleType) DataTuple() *TupleType {
	data := make([]TupleTypeElem, 0, len(t.elems)-t.indexed)
	for _, elem := range t.elems {
		if elem.Indexed {
			continue
		}
		name := elem.Name
		if len(name) == 0 {
			name = fmt.Sprintf("data%d", len(data))
		}
		data = append(data, TupleTypeElem{
			Name: name,
			Type: elem.Type,
		})
	}
	return &TupleType{elems: data}
}

// IsDynamic implements the Type interface.
func (t *EventTupleType) IsDynamic() bool {
	for _, elem := range t.elems {
		if elem.Type.IsDynamic() {
			return true
		}
	}
	return false
}

// CanonicalType implements the Type interface.
func (t *EventTupleType) CanonicalType() string {
	var buf strings.Builder
	buf.WriteString("(")
	for i, elem := range t.elems {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(elem.Type.CanonicalType())
	}
	buf.WriteString(")")
	return buf.String()
}

// String implements the Type interface.
func (t *EventTupleType) String() string {
	var buf strings.Builder
	buf.WriteString("(")
	for i, elem := range t.elems {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(elem.Type.String())
		if elem.Indexed {
			buf.WriteString(" indexed")
		}
		if len(elem.Name) > 0 {
			buf.WriteString(" ")
			buf.WriteString(elem.Name)
		}
	}
	buf.WriteString(")")
	return buf.String()
}

// Value implements the Type interface.
func (t *EventTupleType) Value() Value {
	v := make(TupleValue, len(t.elems))
	// Fills tuple in such a way that indexed fields are first.
	dataIdx, topicIdx := 0, 0
	for _, elem := range t.elems {
		idx := 0
		if elem.Indexed {
			idx = topicIdx
			topicIdx++
		} else {
			idx = dataIdx + t.indexed
			dataIdx++
		}
		v[idx] = TupleValueElem{
			Name:  elem.Name,
			Value: elem.Type.Value(),
		}
		if len(elem.Name) == 0 {
			if elem.Indexed {
				v[idx].Name = fmt.Sprintf("topic%d", topicIdx)
			} else {
				v[idx].Name = fmt.Sprintf("data%d", dataIdx-1)
			}
		}
	}
	return &v
}

// ArrayType represents an unbounded array type.
type ArrayType struct {
	typ Type
}

// NewArrayType creates a dynamic array type with the given element type.
func NewArrayType(typ Type) *ArrayType {
	return &ArrayType{typ: typ}
}

// ElementType returns the type of the array elements.
func (a *ArrayType) ElementType() Type {
	return a.typ
}

// IsDynamic implements the Type interface.
func (a *ArrayType) IsDynamic() bool {
	return true
}

// CanonicalType implements the Type interface.
func (a *ArrayType) CanonicalType() string {
	return a.typ.CanonicalType() + "[]"
}

// String implements the Type interface.
func (a *ArrayType) String() string {
	return a.typ.String() + "[]"
}

// Value implements the Type interface.
func (a *ArrayType) Value() Value {
	return &ArrayValue{Type: a.typ}
}

// FixedArrayType represents a fixed-size array type.
type FixedArrayType struct {
	typ  Type
	size int
}

// NewFixedArrayType creates a new fixed array type with the given element type
// and size.
func NewFixedArrayType(typ Type, size int) *FixedArrayType {
	if size <= 0 {
		panic(fmt.Errorf("abi: invalid array size %d", size))
	}
	return &FixedArrayType{typ: typ, size: size}
}

// Size returns the size of the array.
func (f *FixedArrayType) Size() int {
	return f.size
}

// ElementType returns the type of the array elements.
func (f *FixedArrayType) ElementType() Type {
	return f.typ
}

// IsDynamic implements the Type interface.
func (f *FixedArrayType) IsDynamic() bool {
	return false
}

// CanonicalType implements the Type interface.
func (f *FixedArrayType) CanonicalType() string {
	return f.typ.CanonicalType() + fmt.Sprintf("[%d]", f.size)
}

// String implements the Type interface.
func (f *FixedArrayType) String() string {
	return f.typ.String() + fmt.Sprintf("[%d]", f.size)
}

// Value implements the Type interface.
func (f *FixedArrayType) Value() Value {
	elems := make([]Value, f.size)
	for i := range elems {
		elems[i] = f.typ.Value()
	}
	return (*FixedArrayValue)(&elems)
}

// BytesType represents a bytes type.
type BytesType struct{}

// NewBytesType creates a new "bytes" type.
func NewBytesType() *BytesType {
	return &BytesType{}
}

// IsDynamic implements the Type interface.
func (b *BytesType) IsDynamic() bool {
	return true
}

// CanonicalType implements the Type interface.
func (b *BytesType) CanonicalType() string {
	return "bytes"
}

// String implements the Type interface.
func (b *BytesType) String() string {
	return "bytes"
}

// Value implements the Type interface.
func (b *BytesType) Value() Value {
	return &BytesValue{}
}

// StringType represents a string type.
type StringType struct{}

// NewStringType creates a new "string" type.
func NewStringType() *StringType {
	return &StringType{}
}

// Type implements the Type interface.
func (s *StringType) String() string {
	return "string"
}

// IsDynamic implements the Type interface.
func (s *StringType) IsDynamic() bool {
	return true
}

// CanonicalType implements the Type interface.
func (s *StringType) CanonicalType() string {
	return "string"
}

// Value implements the Type interface.
func (s *StringType) Value() Value {
	return new(StringValue)
}

// FixedBytesType represents a fixed-size bytes type.
type FixedBytesType struct{ size int }

// NewFixedBytesType creates a new fixed-size bytes type with the given size.
// The size must be between 1 and 32.
func NewFixedBytesType(size int) *FixedBytesType {
	if size < 0 || size > 32 {
		panic(fmt.Sprintf("abi: invalid fixed bytes size %d", size))
	}
	return &FixedBytesType{size: size}
}

// Size returns the size of the bytes type.
func (f *FixedBytesType) Size() int {
	return f.size
}

// IsDynamic implements the Type interface.
func (f *FixedBytesType) IsDynamic() bool {
	return false
}

// CanonicalType implements the Type interface.
func (f *FixedBytesType) CanonicalType() string {
	return fmt.Sprintf("bytes%d", f.size)
}

// String implements the Type interface.
func (f *FixedBytesType) String() string {
	return fmt.Sprintf("bytes%d", f.size)
}

// Value implements the Type interface.
func (f *FixedBytesType) Value() Value {
	b := make(FixedBytesValue, f.size)
	return &b
}

// UintType represents an unsigned integer type.
type UintType struct{ size int }

// NewUintType creates a new "uint" type with the given size. The size must be
// between 8 and 256 and a multiple of 8.
func NewUintType(size int) *UintType {
	if size < 0 || size > 256 || size%8 != 0 {
		panic(fmt.Errorf("abi: invalid uint size %d", size))
	}
	return &UintType{size: size}
}

// Size returns the size of the uint type.
func (u *UintType) Size() int {
	return u.size
}

// IsDynamic implements the Type interface.
func (u *UintType) IsDynamic() bool {
	return false
}

// CanonicalType implements the Type interface.
func (u *UintType) CanonicalType() string {
	return fmt.Sprintf("uint%d", u.size)
}

// String implements the Type interface.
func (u *UintType) String() string {
	return fmt.Sprintf("uint%d", u.size)
}

// Value implements the Type interface.
func (u *UintType) Value() Value {
	return &UintValue{Size: u.size}
}

// IntType represents a signed integer type.
type IntType struct{ size int }

// NewIntType creates a new "int" type with the given size. The size must be
// between 8 and 256 and a multiple of 8.
func NewIntType(size int) *IntType {
	if size < 0 || size > 256 || size%8 != 0 {
		panic(fmt.Errorf("abi: invalid int size %d", size))
	}
	return &IntType{size: size}
}

// Size returns the size of the int type.
func (i *IntType) Size() int {
	return i.size
}

// Type implements the Type interface.
func (i *IntType) String() string {
	return fmt.Sprintf("int%d", i.size)
}

// IsDynamic implements the Type interface.
func (i *IntType) IsDynamic() bool {
	return false
}

// CanonicalType implements the Type interface.
func (i *IntType) CanonicalType() string {
	return fmt.Sprintf("int%d", i.size)
}

// Value implements the Type interface.
func (i *IntType) Value() Value {
	return &IntValue{Size: i.size}
}

// BoolType represents a boolean type.
type BoolType struct{}

// NewBoolType creates a new "bool" type.
func NewBoolType() *BoolType {
	return &BoolType{}
}

// IsDynamic implements the Type interface.
func (b *BoolType) IsDynamic() bool {
	return false
}

// CanonicalType implements the Type interface.
func (b *BoolType) CanonicalType() string {
	return "bool"
}

// String implements the Type interface.
func (b *BoolType) String() string {
	return "bool"
}

// Value implements the Type interface.
func (b *BoolType) Value() Value {
	return new(BoolValue)
}

// AddressType represents an address type.
type AddressType struct{}

// NewAddressType creates a new "address" type.
func NewAddressType() *AddressType {
	return &AddressType{}
}

// IsDynamic implements the Type interface.
func (a *AddressType) IsDynamic() bool {
	return false
}

// CanonicalType implements the Type interface.
func (a *AddressType) CanonicalType() string {
	return "address"
}

// String implements the Type interface.
func (a *AddressType) String() string {
	return "address"
}

// Value implements the Type interface.
func (a *AddressType) Value() Value {
	return new(AddressValue)
}

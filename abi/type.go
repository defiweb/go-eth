package abi

import (
	"fmt"
	"strings"
)

// Type represents an ABI type. A type cannot have a value, but can be used to
// create values.
type Type interface {
	// New returns a new zero value of this type.
	New() Value

	// Type returns the user-friendly name of the type.
	Type() string

	// CanonicalType returns the canonical name of the type. In case of a
	// tuple, the canonical name is the canonical name of the tuple's
	// elements, separated by commas and enclosed in parentheses. Arrays
	// are represented by the canonical name of the element type followed
	// by square brackets with the array size.
	CanonicalType() string
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

// New implements the Type interface.
func (a *AliasType) New() Value {
	return a.typ.New()
}

// Type implements the Type interface.
func (a *AliasType) Type() string {
	return a.alias
}

// CanonicalType implements the Type interface.
func (a *AliasType) CanonicalType() string {
	return a.typ.CanonicalType()
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

// New implements the Type interface.
func (t *TupleType) New() Value {
	v := make(TupleValue, len(t.elems))
	for i, elem := range t.elems {
		v[i] = TupleValueElem{
			Name:  elem.Name,
			Value: elem.Type.New(),
		}
		if len(elem.Name) == 0 {
			v[i].Name = fmt.Sprintf("arg%d", i)
		}
	}
	return &v
}

// Type implements the Type interface.
func (t *TupleType) Type() string {
	var buf strings.Builder
	buf.WriteString("(")
	for i, elem := range t.elems {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(elem.Type.Type())
		if len(elem.Name) > 0 {
			buf.WriteString(" ")
			buf.WriteString(elem.Name)
		}
	}
	buf.WriteString(")")
	return buf.String()
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

// New implements the Type interface.
func (t *EventTupleType) New() Value {
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
			Value: elem.Type.New(),
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

// Type implements the Type interface.
func (t *EventTupleType) Type() string {
	var buf strings.Builder
	buf.WriteString("(")
	for i, elem := range t.elems {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(elem.Type.Type())
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

// New implements the Type interface.
func (a *ArrayType) New() Value {
	return &ArrayValue{Type: a.typ}
}

// Type implements the Type interface.
func (a *ArrayType) Type() string {
	return a.typ.Type() + "[]"
}

// CanonicalType implements the Type interface.
func (a *ArrayType) CanonicalType() string {
	return a.typ.CanonicalType() + "[]"
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

// New implements the Type interface.
func (f *FixedArrayType) New() Value {
	elems := make([]Value, f.size)
	for i := range elems {
		elems[i] = f.typ.New()
	}
	return (*FixedArrayValue)(&elems)
}

// Type implements the Type interface.
func (f *FixedArrayType) Type() string {
	return f.typ.Type() + fmt.Sprintf("[%d]", f.size)
}

// CanonicalType implements the Type interface.
func (f *FixedArrayType) CanonicalType() string {
	return f.typ.CanonicalType() + fmt.Sprintf("[%d]", f.size)
}

// BytesType represents a bytes type.
type BytesType struct{}

// NewBytesType creates a new "bytes" type.
func NewBytesType() *BytesType {
	return &BytesType{}
}

// New implements the Type interface.
func (b *BytesType) New() Value {
	return &BytesValue{}
}

// Type implements the Type interface.
func (b *BytesType) Type() string {
	return "bytes"
}

// CanonicalType implements the Type interface.
func (b *BytesType) CanonicalType() string {
	return "bytes"
}

// StringType represents a string type.
type StringType struct{}

// NewStringType creates a new "string" type.
func NewStringType() *StringType {
	return &StringType{}
}

// New implements the Type interface.
func (s *StringType) New() Value {
	return new(StringValue)
}

// Type implements the Type interface.
func (s *StringType) Type() string {
	return "string"
}

// CanonicalType implements the Type interface.
func (s *StringType) CanonicalType() string {
	return "string"
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

// New implements the Type interface.
func (f *FixedBytesType) New() Value {
	b := make(FixedBytesValue, f.size)
	return &b
}

// Type implements the Type interface.
func (f *FixedBytesType) Type() string {
	return fmt.Sprintf("bytes%d", f.size)
}

// CanonicalType implements the Type interface.
func (f *FixedBytesType) CanonicalType() string {
	return fmt.Sprintf("bytes%d", f.size)
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

// New implements the Type interface.
func (u *UintType) New() Value {
	return &UintValue{Size: u.size}
}

// Type implements the Type interface.
func (u *UintType) Type() string {
	return fmt.Sprintf("uint%d", u.size)
}

// CanonicalType implements the Type interface.
func (u *UintType) CanonicalType() string {
	return fmt.Sprintf("uint%d", u.size)
}

// IntType represents an signed integer type.
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

// New implements the Type interface.
func (i *IntType) New() Value {
	return &IntValue{Size: i.size}
}

// Type implements the Type interface.
func (i *IntType) Type() string {
	return fmt.Sprintf("int%d", i.size)
}

// CanonicalType implements the Type interface.
func (i *IntType) CanonicalType() string {
	return fmt.Sprintf("int%d", i.size)
}

// BoolType represents a boolean type.
type BoolType struct{}

// NewBoolType creates a new "bool" type.
func NewBoolType() *BoolType {
	return &BoolType{}
}

// New implements the Type interface.
func (b *BoolType) New() Value {
	return new(BoolValue)
}

// Type implements the Type interface.
func (b *BoolType) Type() string {
	return "bool"
}

// CanonicalType implements the Type interface.
func (b *BoolType) CanonicalType() string {
	return "bool"
}

// AddressType represents an address type.
type AddressType struct{}

// NewAddressType creates a new "address" type.
func NewAddressType() *AddressType {
	return &AddressType{}
}

// New implements the Type interface.
func (a *AddressType) New() Value {
	return new(AddressValue)
}

// Type implements the Type interface.
func (a *AddressType) Type() string {
	return "address"
}

// CanonicalType implements the Type interface.
func (a *AddressType) CanonicalType() string {
	return "address"
}

package abi

import (
	"fmt"
	"strings"
)

// Type represents an ABI type.
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

type TupleType struct {
	elems []TupleTypeElem
}

type TupleTypeElem struct {
	Name string
	Type Type
}

// NewTupleType creates a new tuple type with the given elements.
func NewTupleType(elems ...TupleTypeElem) *TupleType {
	return &TupleType{elems: elems}
}

func (t *TupleType) Size() int {
	return len(t.elems)
}

func (t *TupleType) Elements() []TupleTypeElem {
	cpy := make([]TupleTypeElem, len(t.elems))
	copy(cpy, t.elems)
	return cpy
}

func (t *TupleType) New() Value {
	v := NewTupleOfSize(len(t.elems))
	for i, elem := range t.elems {
		v.elems[i] = elem.Type.New()
		v.names[i] = elem.Name
		if len(elem.Name) == 0 {
			v.names[i] = fmt.Sprintf("arg%d", i)
		}
	}
	return v
}

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

type EventTupleType struct {
	elems   []EventTupleTypeElem
	indexed int
}

type EventTupleTypeElem struct {
	Name    string
	Indexed bool
	Type    Type
}

// NewEventTupleType creates a new tuple type with the given elements.
func NewEventTupleType(elems ...EventTupleTypeElem) *EventTupleType {
	indexed := 0
	for _, elem := range elems {
		if elem.Indexed {
			indexed++
		}
	}
	return &EventTupleType{elems: elems, indexed: indexed}
}

func (t *EventTupleType) Size() int {
	return len(t.elems)
}

func (t *EventTupleType) IndexedSize() int {
	return t.indexed
}

func (t *EventTupleType) DataSize() int {
	return len(t.elems) - t.indexed
}

func (t *EventTupleType) Elements() []EventTupleTypeElem {
	cpy := make([]EventTupleTypeElem, len(t.elems))
	copy(cpy, t.elems)
	return cpy
}

func (t *EventTupleType) New() Value {
	tuple := NewTupleOfSize(len(t.elems))
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
		tuple.elems[idx] = elem.Type.New()
		tuple.names[idx] = elem.Name
		if len(elem.Name) == 0 {
			if elem.Indexed {
				tuple.names[idx] = fmt.Sprintf("topic%d", topicIdx)
			} else {
				tuple.names[idx] = fmt.Sprintf("data%d", dataIdx-1)
			}
		}
	}
	return tuple
}

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

type ArrayType struct {
	typ Type
}

// NewArrayType creates a dynamic array type with the given element type.
func NewArrayType(typ Type) *ArrayType {
	return &ArrayType{typ: typ}
}

func (a *ArrayType) Element() Type {
	return a.typ
}

func (a *ArrayType) New() Value {
	return NewArray(a.typ)
}

func (a *ArrayType) Type() string {
	return a.typ.Type() + "[]"
}

func (a *ArrayType) CanonicalType() string {
	return a.typ.CanonicalType() + "[]"
}

type FixedArrayType struct {
	typ  Type
	size int
}

// NewFixedArrayType creates a new fixed array type with the given element type
// and size. The size must be greater than zero.
func NewFixedArrayType(typ Type, size int) *FixedArrayType {
	if size <= 0 {
		panic(fmt.Errorf("abi: invalid array size %d", size))
	}
	return &FixedArrayType{typ: typ, size: size}
}

func (f *FixedArrayType) Size() int {
	return f.size
}

func (f *FixedArrayType) Element() Type {
	return f.typ
}

func (f *FixedArrayType) New() Value {
	return NewFixedArray(f.typ, f.size)
}

func (f *FixedArrayType) Type() string {
	return f.typ.Type() + fmt.Sprintf("[%d]", f.size)
}

func (f *FixedArrayType) CanonicalType() string {
	return f.typ.CanonicalType() + fmt.Sprintf("[%d]", f.size)
}

type BytesType struct{}

// NewBytesType creates a new "bytes" type.
func NewBytesType() *BytesType {
	return &BytesType{}
}

func (b *BytesType) New() Value {
	return NewBytes()
}

func (b *BytesType) Type() string {
	return "bytes"
}

func (b *BytesType) CanonicalType() string {
	return "bytes"
}

type StringType struct{}

// NewStringType creates a new "string" type.
func NewStringType() *StringType {
	return &StringType{}
}

func (s *StringType) New() Value {
	return NewString()
}

func (s *StringType) Type() string {
	return "string"
}

func (s *StringType) CanonicalType() string {
	return "string"
}

type FixedBytesType struct{ size int }

func NewFixedBytesType(size int) *FixedBytesType {
	if size < 0 || size > 32 {
		panic(fmt.Sprintf("abi: invalid fixed bytes size %d", size))
	}
	return &FixedBytesType{size: size}
}

func (f *FixedBytesType) Size() int {
	return f.size
}

func (f *FixedBytesType) New() Value {
	return NewFixedBytes(f.size)
}

func (f *FixedBytesType) Type() string {
	return fmt.Sprintf("bytes%d", f.size)
}

func (f *FixedBytesType) CanonicalType() string {
	return fmt.Sprintf("bytes%d", f.size)
}

type UintType struct{ size int }

// NewUintType creates a new "uint" type with the given size. The size is in
// bytes and must be between 1 and 32.
func NewUintType(size int) *UintType {
	if size < 0 || size > 32 {
		panic(fmt.Errorf("abi: invalid uint size %d", size))
	}
	return &UintType{size: size}
}

func (u *UintType) New() Value {
	return NewUint(u.size)
}

func (u *UintType) Size() int {
	return u.size
}

func (u *UintType) Type() string {
	return fmt.Sprintf("uint%d", u.size*8)
}

func (u *UintType) CanonicalType() string {
	return fmt.Sprintf("uint%d", u.size*8)
}

type IntType struct{ size int }

// NewIntType creates a new "int" type with the given size. The size is in
// bytes and must be between 1 and 32.
func NewIntType(size int) *IntType {
	if size < 0 || size > 32 {
		panic(fmt.Errorf("abi: invalid int size %d", size))
	}
	return &IntType{size: size}
}

func (i *IntType) New() Value {
	return NewInt(i.size)
}

func (i *IntType) Size() int {
	return i.size
}

func (i *IntType) Type() string {
	return fmt.Sprintf("int%d", i.size*8)
}

func (i *IntType) CanonicalType() string {
	return fmt.Sprintf("int%d", i.size*8)
}

type BoolType struct{}

// NewBoolType creates a new "bool" type.
func NewBoolType() *BoolType {
	return &BoolType{}
}

func (b *BoolType) New() Value {
	return NewBool()
}

func (b *BoolType) Type() string {
	return "bool"
}

func (b *BoolType) CanonicalType() string {
	return "bool"
}

type AddressType struct{}

// NewAddressType creates a new "address" type.
func NewAddressType() *AddressType {
	return &AddressType{}
}

func (a *AddressType) New() Value {
	return NewAddress()
}

func (a *AddressType) Type() string {
	return "address"
}

func (a *AddressType) CanonicalType() string {
	return "address"
}

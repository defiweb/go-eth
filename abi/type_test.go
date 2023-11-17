package abi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type nullType struct{}
type nullValue struct{}
type dynamicNullType struct{ nullType }

func (n nullType) Value() Value                    { return new(nullValue) }
func (n nullType) String() string                  { return "null" }
func (n nullType) IsDynamic() bool                 { return false }
func (n nullType) CanonicalType() string           { return "null" }
func (n nullValue) IsDynamic() bool                { return false }
func (n nullValue) EncodeABI() (Words, error)      { return nil, nil }
func (n nullValue) DecodeABI(_ Words) (int, error) { return 0, nil }
func (n dynamicNullType) IsDynamic() bool          { return true }

func TestAliasType(t *testing.T) {
	v := NewAliasType("alias", nullType{})
	assert.Equal(t, &nullValue{}, v.Value())
	assert.Equal(t, "alias", v.String())
	assert.Equal(t, "null", v.CanonicalType())
}

func TestTupleType(t *testing.T) {
	v := NewTupleType(
		TupleTypeElem{Name: "foo", Type: nullType{}},
		TupleTypeElem{Name: "bar", Type: nullType{}},
		TupleTypeElem{Type: nullType{}},
	)
	assert.Equal(t, &TupleValue{
		{Name: "foo", Value: &nullValue{}},
		{Name: "bar", Value: &nullValue{}},
		{Name: "arg2", Value: &nullValue{}},
	}, v.Value())
	assert.Equal(t, "(null foo, null bar, null)", v.String())
	assert.Equal(t, "(null,null,null)", v.CanonicalType())
}

func TestEventTupleType(t *testing.T) {
	v := NewEventTupleType(
		EventTupleElem{Name: "foo", Type: nullType{}},
		EventTupleElem{Name: "bar", Type: nullType{}, Indexed: true},
		EventTupleElem{Name: "qux", Type: dynamicNullType{}, Indexed: true},
		EventTupleElem{Type: nullType{}, Indexed: true},
		EventTupleElem{Type: nullType{}},
	)
	assert.Equal(t, &TupleValue{
		{Name: "bar", Value: &nullValue{}},
		{Name: "qux", Value: &nullValue{}},
		{Name: "topic3", Value: &nullValue{}},
		{Name: "foo", Value: &nullValue{}},
		{Name: "data1", Value: &nullValue{}},
	}, v.Value())
	assert.Equal(t, "(null foo, null indexed bar, null indexed qux, null indexed, null)", v.String())
	assert.Equal(t, "(null,null,null,null,null)", v.CanonicalType())
	assert.Equal(t, "(null bar, bytes32 qux, null topic3)", v.TopicsTuple().String())
	assert.Equal(t, "(null foo, null data1)", v.DataTuple().String())
}

func TestArrayType(t *testing.T) {
	v := NewArrayType(&nullType{})
	assert.Equal(t, "null[]", v.String())
	assert.Equal(t, "null[]", v.CanonicalType())
	assert.Equal(t, &ArrayValue{
		Type:  &nullType{},
		Elems: nil,
	}, v.Value())
}

func TestFixedArrayType(t *testing.T) {
	v := NewFixedArrayType(&nullType{}, 2)
	assert.Equal(t, "null[2]", v.String())
	assert.Equal(t, "null[2]", v.CanonicalType())
	assert.Equal(t, &FixedArrayValue{
		&nullValue{},
		&nullValue{},
	}, v.Value())
}

func TestBytesType(t *testing.T) {
	v := NewBytesType()
	assert.Equal(t, "bytes", v.String())
	assert.Equal(t, "bytes", v.CanonicalType())
	assert.Equal(t, &BytesValue{}, v.Value())
}

func TestStringType(t *testing.T) {
	v := NewStringType()
	assert.Equal(t, "string", v.String())
	assert.Equal(t, "string", v.CanonicalType())
	assert.Equal(t, new(StringValue), v.Value())
}

func TestFixedBytesType(t *testing.T) {
	v := NewFixedBytesType(2)
	assert.Equal(t, "bytes2", v.String())
	assert.Equal(t, "bytes2", v.CanonicalType())
	assert.Equal(t, &FixedBytesValue{0, 0}, v.Value())
}

func TestUintType(t *testing.T) {
	v := NewUintType(256)
	assert.Equal(t, "uint256", v.String())
	assert.Equal(t, "uint256", v.CanonicalType())
	assert.Equal(t, &UintValue{Size: 256}, v.Value())
}

func TestIntType(t *testing.T) {
	v := NewIntType(256)
	assert.Equal(t, "int256", v.String())
	assert.Equal(t, "int256", v.CanonicalType())
	assert.Equal(t, &IntValue{Size: 256}, v.Value())
}

func TestBoolType(t *testing.T) {
	v := NewBoolType()
	assert.Equal(t, "bool", v.String())
	assert.Equal(t, "bool", v.CanonicalType())
	assert.Equal(t, new(BoolValue), v.Value())
}

func TestAddressType(t *testing.T) {
	v := NewAddressType()
	assert.Equal(t, "address", v.String())
	assert.Equal(t, "address", v.CanonicalType())
	assert.Equal(t, new(AddressValue), v.Value())
}

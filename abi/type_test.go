package abi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type nullType struct{}
type nullValue struct{}

func (n nullType) Value() Value                    { return new(nullValue) }
func (n nullType) String() string                  { return "null" }
func (n nullType) CanonicalType() string           { return "null" }
func (n nullValue) IsDynamic() bool                { return false }
func (n nullValue) EncodeABI() (Words, error)      { return nil, nil }
func (n nullValue) DecodeABI(_ Words) (int, error) { return 0, nil }

func TestAliasType(t *testing.T) {
	v := NewAliasType("alias", nullType{})
	assert.Equal(t, &nullValue{}, v.Value())
	assert.Equal(t, v.String(), "alias")
	assert.Equal(t, v.CanonicalType(), "null")
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
	assert.Equal(t, v.String(), "(null foo, null bar, null)")
	assert.Equal(t, v.CanonicalType(), "(null,null,null)")
}

func TestEventTupleType(t *testing.T) {
	v := NewEventTupleType(
		EventTupleElem{Name: "foo", Type: nullType{}},
		EventTupleElem{Name: "bar", Type: nullType{}, Indexed: true},
		EventTupleElem{Type: nullType{}, Indexed: true},
		EventTupleElem{Type: nullType{}},
	)
	assert.Equal(t, &TupleValue{
		{Name: "bar", Value: &nullValue{}},
		{Name: "topic2", Value: &nullValue{}},
		{Name: "foo", Value: &nullValue{}},
		{Name: "data1", Value: &nullValue{}},
	}, v.Value())
	assert.Equal(t, v.String(), "(null foo, null indexed bar, null indexed, null)")
	assert.Equal(t, v.CanonicalType(), "(null,null,null,null)")
}

func TestArrayType(t *testing.T) {
	v := NewArrayType(&nullType{})
	assert.Equal(t, &ArrayValue{
		Type:  &nullType{},
		Elems: nil,
	}, v.Value())
	assert.Equal(t, v.String(), "null[]")
	assert.Equal(t, v.CanonicalType(), "null[]")
}

func TestFixedArrayType(t *testing.T) {
	v := NewFixedArrayType(&nullType{}, 2)
	assert.Equal(t, &FixedArrayValue{
		&nullValue{},
		&nullValue{},
	}, v.Value())
	assert.Equal(t, v.String(), "null[2]")
	assert.Equal(t, v.CanonicalType(), "null[2]")
}

func TestBytesType(t *testing.T) {
	v := NewBytesType()
	assert.Equal(t, &BytesValue{}, v.Value())
	assert.Equal(t, v.String(), "bytes")
	assert.Equal(t, v.CanonicalType(), "bytes")
}

func TestStringType(t *testing.T) {
	v := NewStringType()
	assert.Equal(t, new(StringValue), v.Value())
	assert.Equal(t, v.String(), "string")
	assert.Equal(t, v.CanonicalType(), "string")
}

func TestFixedBytesType(t *testing.T) {
	v := NewFixedBytesType(2)
	assert.Equal(t, &FixedBytesValue{0, 0}, v.Value())
	assert.Equal(t, v.String(), "bytes2")
	assert.Equal(t, v.CanonicalType(), "bytes2")
}

func TestUintType(t *testing.T) {
	v := NewUintType(256)
	assert.Equal(t, &UintValue{Size: 256}, v.Value())
	assert.Equal(t, v.String(), "uint256")
	assert.Equal(t, v.CanonicalType(), "uint256")
}

func TestIntType(t *testing.T) {
	v := NewIntType(256)
	assert.Equal(t, &IntValue{Size: 256}, v.Value())
	assert.Equal(t, v.String(), "int256")
	assert.Equal(t, v.CanonicalType(), "int256")
}

func TestBoolType(t *testing.T) {
	v := NewBoolType()
	assert.Equal(t, new(BoolValue), v.Value())
	assert.Equal(t, v.String(), "bool")
	assert.Equal(t, v.CanonicalType(), "bool")
}

func TestAddressType(t *testing.T) {
	v := NewAddressType()
	assert.Equal(t, new(AddressValue), v.Value())
	assert.Equal(t, v.String(), "address")
	assert.Equal(t, v.CanonicalType(), "address")
}

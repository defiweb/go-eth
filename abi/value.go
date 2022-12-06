package abi

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/defiweb/go-eth/types"

	"github.com/defiweb/go-anymapper"
)

// Value represents a value that can be encoded to and from ABI.
//
// Values are used as an intermediate representation during encoding and
// decoding ABI data. Usually, they are not used outside the abi package.
//
// When data is encoded using Encoder, the values provided to Encoder are
// mapped to Value instances using the anymapper package, and then they are used
// to encode the ABI data. When the data is decoded using Decoder, the Value
// instances are used to decode the ABI data, and then the values are mapped to
// the target types.
type Value interface {
	// IsDynamic indicates whether the type is dynamic.
	IsDynamic() bool

	// EncodeABI returns the ABI encoding of the value.
	EncodeABI() (Words, error)

	// DecodeABI sets the value from the ABI encoded data.
	DecodeABI(Words) (int, error)
}

// TupleValue is a value of tuple type.
//
// During encoding, the TupleValue can be mapped from a struct or a map where
// keys or struct fields are used as tuple element names.
//
// During decoding, the TupleValue can be mapped to a struct or a map where
// tuple element names are used as keys or struct fields.
type TupleValue []TupleValueElem

// TupleValueElem is an element of tuple value.
type TupleValueElem struct {
	// Name of the tuple element. It is used when mapping values from and to
	// maps and structures.
	Name string

	// Value is the value of the tuple element. It is used to encode and decode
	// the ABI data.
	Value Value
}

// IsDynamic implements the Value interface.
func (t *TupleValue) IsDynamic() bool {
	for _, elem := range *t {
		if elem.Value.IsDynamic() {
			return true
		}
	}
	return false
}

// EncodeABI implements the Value interface.
func (t *TupleValue) EncodeABI() (Words, error) {
	elems := make([]Value, len(*t))
	for i, elem := range *t {
		elems[i] = elem.Value
	}
	return encodeTuple(elems)
}

// DecodeABI implements the Value interface.
func (t *TupleValue) DecodeABI(words Words) (int, error) {
	elems := make([]Value, len(*t))
	for i, elem := range *t {
		elems[i] = elem.Value
	}
	return decodeTuple(&elems, words)
}

// MapFrom implements the anymapper.MapFrom interface.
func (t *TupleValue) MapFrom(m *anymapper.Mapper, src reflect.Value) error {
	vals := make(map[string]Value, len(*t))
	for _, elem := range *t {
		vals[elem.Name] = elem.Value
	}
	return m.MapRefl(src, reflect.ValueOf(vals))
}

// MapTo implements the anymapper.MapTo interface.
func (t *TupleValue) MapTo(m *anymapper.Mapper, dest reflect.Value) error {
	vals := make(map[string]Value, len(*t))
	for _, elem := range *t {
		vals[elem.Name] = elem.Value
	}
	return m.MapRefl(reflect.ValueOf(vals), dest)
}

// ArrayValue is a value of array type.
//
// During encoding, the ArrayValue can be mapped from a slice or an array.
//
// During decoding the ArrayValue is mapped to a slice or an array of the
// same size.
type ArrayValue struct {
	Elems []Value
	Type  Type
}

// IsDynamic implements the Value interface.
func (a *ArrayValue) IsDynamic() bool {
	return true
}

// EncodeABI implements the Value interface.
func (a *ArrayValue) EncodeABI() (Words, error) {
	return encodeArray(a.Elems)
}

// DecodeABI implements the Value interface.
func (a *ArrayValue) DecodeABI(data Words) (int, error) {
	return decodeArray(&a.Elems, data, a.Type)
}

// MapFrom implements the anymapper.MapFrom interface.
func (a *ArrayValue) MapFrom(m *anymapper.Mapper, src reflect.Value) error {
	if src.Kind() != reflect.Slice && src.Kind() != reflect.Array {
		return fmt.Errorf("abi: cannot map array from %s", src.Kind())
	}
	a.Elems = make([]Value, src.Len())
	for i := 0; i < src.Len(); i++ {
		a.Elems[i] = a.Type.Value()
	}
	return m.MapRefl(src, reflect.ValueOf(&a.Elems))
}

// MapTo implements the anymapper.MapTo interface.
func (a *ArrayValue) MapTo(m *anymapper.Mapper, dest reflect.Value) error {
	return m.MapRefl(reflect.ValueOf(&a.Elems), dest)
}

// FixedArrayValue is a value of fixed array type. The size of a slice is
// assumed to be equal to the size of the type.
//
// During encoding, the FixedArrayValue can be mapped from a slice or an array.
//
// During decoding the FixedArrayValue is mapped to a slice or an array of the
// same size.
type FixedArrayValue []Value

// IsDynamic implements the Value interface.
func (a FixedArrayValue) IsDynamic() bool {
	return false
}

// EncodeABI implements the Value interface.
func (a FixedArrayValue) EncodeABI() (Words, error) {
	return encodeFixedArray(a)
}

// DecodeABI implements the Value interface.
func (a FixedArrayValue) DecodeABI(data Words) (int, error) {
	return decodeFixedArray((*[]Value)(&a), data)
}

// MapFrom implements the anymapper.MapFrom interface.
func (a FixedArrayValue) MapFrom(m *anymapper.Mapper, src reflect.Value) error {
	if src.Kind() != reflect.Slice && src.Kind() != reflect.Array {
		return fmt.Errorf("abi: cannot map %s to array[%d]", src.Type(), len(a))
	}
	if src.Len() != len(a) {
		return fmt.Errorf("abi: cannot map %d elements to array[%d]", src.Len(), len(a))
	}
	return m.MapRefl(src, reflect.ValueOf((*[]Value)(&a)))
}

// MapTo implements the anymapper.MapTo interface.
func (a FixedArrayValue) MapTo(m *anymapper.Mapper, dest reflect.Value) error {
	return m.MapRefl(reflect.ValueOf(([]Value)(a)), dest)
}

// BytesValue is a value of bytes type.
//
// During encoding ad decoding, the BytesValue can be mapped using the slice
// rules described in the documentation of anymapper package.
type BytesValue []byte

// IsDynamic implements the Value interface.
func (b *BytesValue) IsDynamic() bool {
	return true
}

// EncodeABI implements the Value interface.
func (b *BytesValue) EncodeABI() (Words, error) {
	return encodeBytes(*b)
}

// DecodeABI implements the Value interface.
func (b *BytesValue) DecodeABI(data Words) (int, error) {
	return decodeBytes((*[]byte)(b), data)
}

// MapFrom implements the anymapper.MapFrom interface.
func (b *BytesValue) MapFrom(m *anymapper.Mapper, src reflect.Value) error {
	return m.MapRefl(src, reflect.ValueOf((*[]byte)(b)))
}

// MapTo implements the anymapper.MapTo interface.
func (b *BytesValue) MapTo(m *anymapper.Mapper, dest reflect.Value) error {
	return m.MapRefl(reflect.ValueOf((*[]byte)(b)), dest)
}

// StringValue is a value of bytes type.
//
// During encoding ad decoding, the StringValue is mapped using the string
// rules described in the documentation of anymapper package.
type StringValue string

// IsDynamic implements the Value interface.
func (s *StringValue) IsDynamic() bool {
	return true
}

// EncodeABI implements the Value interface.
func (s *StringValue) EncodeABI() (Words, error) {
	return encodeBytes([]byte(*s))
}

// DecodeABI implements the Value interface.
func (s *StringValue) DecodeABI(data Words) (int, error) {
	var b []byte
	if _, err := decodeBytes(&b, data); err != nil {
		return 0, err
	}
	*s = StringValue(b)
	return 1, nil
}

// MapFrom implements the anymapper.MapFrom interface.
func (s *StringValue) MapFrom(m *anymapper.Mapper, src reflect.Value) error {
	return m.MapRefl(src, reflect.ValueOf((*string)(s)))
}

// MapTo implements the anymapper.MapTo interface.
func (s *StringValue) MapTo(m *anymapper.Mapper, dest reflect.Value) error {
	return m.MapRefl(reflect.ValueOf((*string)(s)), dest)
}

// FixedBytesValue is a value of fixed bytes type. The size of a slice is
// assumed to be equal to the size of the bytesN type.
//
// During encoding and decoding, the FixedBytesValue is mapped using the slice
// rules described in the documentation of anymapper package. Both values must
// have the same size.
type FixedBytesValue []byte

// IsDynamic implements the Value interface.
func (b FixedBytesValue) IsDynamic() bool {
	return false
}

// EncodeABI implements the Value interface.
func (b FixedBytesValue) EncodeABI() (Words, error) {
	return encodeFixedBytes(b, len(b))
}

// DecodeABI implements the Value interface.
func (b FixedBytesValue) DecodeABI(data Words) (int, error) {
	return decodeFixedBytes((*[]byte)(&b), data, len(b))
}

// MapFrom implements the anymapper.MapFrom interface.
func (b FixedBytesValue) MapFrom(m *anymapper.Mapper, src reflect.Value) error {
	var dst []byte
	if err := m.MapRefl(src, reflect.ValueOf(&dst)); err != nil {
		return err
	}
	if len(dst) != len(b) {
		return fmt.Errorf("abi: cannot map %d bytes to bytes%d", len(dst), len(b))
	}
	copy(b, dst)
	return nil
}

// MapTo implements the anymapper.MapTo interface.
func (b FixedBytesValue) MapTo(m *anymapper.Mapper, dest reflect.Value) error {
	return m.MapRefl(reflect.ValueOf((*[]byte)(&b)), dest)
}

// UintValue is a value of uintN types.
//
// During encoding, the UintValue is mapped to the *big.Int type using the
// rules described in the documentation of anymapper package.
//
// During decoding, the UintValue is mapped from the *big.Int type using the
// rules described in the documentation of anymapper package.
type UintValue struct {
	big.Int
	Size int
}

// IsDynamic implements the Value interface.
func (u *UintValue) IsDynamic() bool {
	return false
}

// EncodeABI implements the Value interface.
func (u *UintValue) EncodeABI() (Words, error) {
	if u.Size < 8 || u.Size > 256 || u.Size%8 != 0 {
		return nil, fmt.Errorf("abi: invalid uint size: %d", u.Size)
	}
	return encodeUint(&u.Int, u.Size)
}

// DecodeABI implements the Value interface.
func (u *UintValue) DecodeABI(words Words) (int, error) {
	if u.Size < 8 || u.Size > 256 || u.Size%8 != 0 {
		return 0, fmt.Errorf("abi: invalid uint size: %d", u.Size)
	}
	return decodeUint(&u.Int, words, u.Size)
}

// MapFrom implements the anymapper.MapFrom interface.
func (u *UintValue) MapFrom(m *anymapper.Mapper, src reflect.Value) error {
	return m.MapRefl(src, reflect.ValueOf(&u.Int))
}

// MapTo implements the anymapper.MapTo interface.
func (u *UintValue) MapTo(m *anymapper.Mapper, dest reflect.Value) error {
	return m.MapRefl(reflect.ValueOf(&u.Int), dest)
}

// IntValue is a value of intN types.
//
// During encoding, the IntValue is mapped to the *big.Int type using the
// rules described in the documentation of anymapper package.
//
// During decoding, the IntValue is mapped from the *big.Int type using the
// rules described in the documentation of anymapper package.
type IntValue struct {
	big.Int
	Size int
}

// IsDynamic implements the Value interface.
func (i *IntValue) IsDynamic() bool {
	return false
}

// EncodeABI implements the Value interface.
func (i *IntValue) EncodeABI() (Words, error) {
	if i.Size < 8 || i.Size > 256 || i.Size%8 != 0 {
		return nil, fmt.Errorf("abi: invalid int size: %d", i.Size)
	}
	return encodeInt(&i.Int, i.Size)
}

// DecodeABI implements the Value interface.
func (i *IntValue) DecodeABI(words Words) (int, error) {
	if i.Size < 8 || i.Size > 256 || i.Size%8 != 0 {
		return 0, fmt.Errorf("abi: invalid int size: %d", i.Size)
	}
	return decodeInt(&i.Int, words, i.Size)
}

// MapFrom implements the anymapper.MapFrom interface.
func (i *IntValue) MapFrom(m *anymapper.Mapper, src reflect.Value) error {
	return m.MapRefl(src, reflect.ValueOf(&i.Int))
}

// MapTo implements the anymapper.MapTo interface.
func (i *IntValue) MapTo(m *anymapper.Mapper, dest reflect.Value) error {
	return m.MapRefl(reflect.ValueOf(&i.Int), dest)
}

// BoolValue is a value of bool type.
//
// During encoding and decoding, the BoolValue is mapped using the bool rules
// described in the documentation of anymapper package.
type BoolValue bool

// IsDynamic implements the Value interface.
func (b *BoolValue) IsDynamic() bool {
	return false
}

// EncodeABI implements the Value interface.
func (b *BoolValue) EncodeABI() (Words, error) {
	return encodeBool(bool(*b)), nil
}

// DecodeABI implements the Value interface.
func (b *BoolValue) DecodeABI(words Words) (int, error) {
	return decodeBool((*bool)(b), words)
}

// MapFrom implements the anymapper.MapFrom interface.
func (b *BoolValue) MapFrom(m *anymapper.Mapper, src reflect.Value) error {
	return m.MapRefl(src, reflect.ValueOf((*bool)(b)))
}

// MapTo implements the anymapper.MapTo interface.
func (b *BoolValue) MapTo(m *anymapper.Mapper, dest reflect.Value) error {
	return m.MapRefl(reflect.ValueOf((*bool)(b)), dest)
}

// AddressValue is a value of address type.
//
// During encoding, the AddressValue can be mapped to the types.Address type,
// string as a hex-encoded address. For other types, the rules for []byte slice
// described in the documentation of anymapper package are used.
type AddressValue types.Address

// IsDynamic implements the Value interface.
func (a *AddressValue) IsDynamic() bool {
	return false
}

// EncodeABI implements the Value interface.
func (a *AddressValue) EncodeABI() (Words, error) {
	return encodeAddress(types.Address(*a))
}

// DecodeABI implements the Value interface.
func (a *AddressValue) DecodeABI(words Words) (int, error) {
	return decodeAddress((*types.Address)(a), words)
}

// MapFrom implements the anymapper.MapFrom interface.
func (a *AddressValue) MapFrom(m *anymapper.Mapper, src reflect.Value) error {
	var err error
	var addr types.Address
	if !m.StrictTypes && src.Type().Kind() == reflect.String {
		addr, err = types.HexToAddress(src.String())
		if err != nil {
			return fmt.Errorf("abi: cannot convert string to address: %v", err)
		}
	} else {
		if err := m.MapRefl(src, reflect.ValueOf(&addr)); err != nil {
			return err
		}
	}
	*a = AddressValue(addr)
	return nil
}

// MapTo implements the anymapper.MapTo interface.
func (a *AddressValue) MapTo(m *anymapper.Mapper, dest reflect.Value) error {
	if !m.StrictTypes && dest.Type().Kind() == reflect.String {
		dest.SetString(types.Address(*a).String())
		return nil
	}
	addr := types.Address(*a)
	return m.MapRefl(reflect.ValueOf(addr), dest)
}

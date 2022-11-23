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

type TupleValue []TupleValueElem

type TupleValueElem struct {
	Value Value
	Name  string
}

func (t *TupleValue) IsDynamic() bool {
	for _, elem := range *t {
		if elem.Value.IsDynamic() {
			return true
		}
	}
	return false
}

func (t *TupleValue) EncodeABI() (Words, error) {
	elems := make([]Value, len(*t))
	for i, elem := range *t {
		elems[i] = elem.Value
	}
	return encodeTuple(elems)
}

func (t *TupleValue) DecodeABI(words Words) (int, error) {
	elems := make([]Value, len(*t))
	for i, elem := range *t {
		elems[i] = elem.Value
	}
	return decodeTuple(&elems, words)
}

func (t *TupleValue) MapFrom(m *anymapper.Mapper, src reflect.Value) error {
	vals := make(map[string]Value)
	for _, elem := range *t {
		vals[elem.Name] = elem.Value
	}
	return m.MapRefl(src, reflect.ValueOf(vals))
}

func (t *TupleValue) MapInto(m *anymapper.Mapper, dest reflect.Value) error {
	vals := make(map[string]Value)
	for _, elem := range *t {
		vals[elem.Name] = elem.Value
	}
	return m.MapRefl(reflect.ValueOf(vals), dest)
}

type ArrayValue struct {
	Elems []Value
	Type  Type
}

func (a *ArrayValue) IsDynamic() bool {
	return true
}

func (a *ArrayValue) EncodeABI() (Words, error) {
	return encodeArray(a.Elems)
}

func (a *ArrayValue) DecodeABI(data Words) (int, error) {
	return decodeArray(&a.Elems, data, a.Type)
}

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

func (a *ArrayValue) MapInto(m *anymapper.Mapper, dest reflect.Value) error {
	return m.MapRefl(reflect.ValueOf(&a.Elems), dest)
}

type FixedArrayValue []Value

func (a FixedArrayValue) IsDynamic() bool {
	return false
}

func (a FixedArrayValue) EncodeABI() (Words, error) {
	return encodeFixedArray(a)
}

func (a FixedArrayValue) DecodeABI(data Words) (int, error) {
	return decodeFixedArray((*[]Value)(&a), data)
}

func (a FixedArrayValue) MapFrom(m *anymapper.Mapper, src reflect.Value) error {
	if src.Kind() != reflect.Slice && src.Kind() != reflect.Array {
		return fmt.Errorf("abi: cannot map %s to array[%d]", src.Type(), len(a))
	}
	if src.Len() != len(a) {
		return fmt.Errorf("abi: cannot map %d elements to array[%d]", src.Len(), len(a))
	}
	return m.MapRefl(src, reflect.ValueOf((*[]Value)(&a)))
}

func (a FixedArrayValue) MapInto(m *anymapper.Mapper, dest reflect.Value) error {
	return m.MapRefl(reflect.ValueOf(([]Value)(a)), dest)
}

type BytesValue []byte

func (b *BytesValue) IsDynamic() bool {
	return true
}

func (b *BytesValue) EncodeABI() (Words, error) {
	return encodeBytes(*b)
}

func (b *BytesValue) DecodeABI(data Words) (int, error) {
	return decodeBytes((*[]byte)(b), data)
}

func (b *BytesValue) MapFrom(m *anymapper.Mapper, src reflect.Value) error {
	return m.MapRefl(src, reflect.ValueOf((*[]byte)(b)))
}

func (b *BytesValue) MapInto(m *anymapper.Mapper, dest reflect.Value) error {
	return m.MapRefl(reflect.ValueOf((*[]byte)(b)), dest)
}

type StringValue string

func (s *StringValue) IsDynamic() bool {
	return true
}

func (s *StringValue) EncodeABI() (Words, error) {
	return encodeBytes([]byte(*s))
}

func (s *StringValue) DecodeABI(data Words) (int, error) {
	var b []byte
	if _, err := decodeBytes(&b, data); err != nil {
		return 0, err
	}
	*s = StringValue(b)
	return 1, nil
}

func (s *StringValue) MapFrom(m *anymapper.Mapper, src reflect.Value) error {
	return m.MapRefl(src, reflect.ValueOf((*string)(s)))
}

func (s *StringValue) MapInto(m *anymapper.Mapper, dest reflect.Value) error {
	return m.MapRefl(reflect.ValueOf((*string)(s)), dest)
}

type FixedBytesValue []byte

func (b FixedBytesValue) IsDynamic() bool {
	return false
}

func (b FixedBytesValue) EncodeABI() (Words, error) {
	return encodeFixedBytes(b, len(b))
}

func (b FixedBytesValue) DecodeABI(data Words) (int, error) {
	return decodeFixedBytes((*[]byte)(&b), data, len(b))
}

func (b FixedBytesValue) MapFrom(m *anymapper.Mapper, src reflect.Value) error {
	return m.MapRefl(src, reflect.ValueOf((*[]byte)(&b)))
}

func (b FixedBytesValue) MapInto(m *anymapper.Mapper, dest reflect.Value) error {
	return m.MapRefl(reflect.ValueOf((*[]byte)(&b)), dest)
}

type UintValue struct {
	big.Int
	Size int
}

func (u *UintValue) IsDynamic() bool {
	return false
}

func (u *UintValue) EncodeABI() (Words, error) {
	if u.Size < 8 || u.Size > 256 || u.Size%8 != 0 {
		return nil, fmt.Errorf("abi: invalid uint size: %d", u.Size)
	}
	return encodeUint(&u.Int, u.Size)
}

func (u *UintValue) DecodeABI(words Words) (int, error) {
	if u.Size < 8 || u.Size > 256 || u.Size%8 != 0 {
		return 0, fmt.Errorf("abi: invalid uint size: %d", u.Size)
	}
	return decodeUint(&u.Int, words, u.Size)
}

func (u *UintValue) MapFrom(m *anymapper.Mapper, src reflect.Value) error {
	var val *big.Int
	if err := m.MapRefl(src, reflect.ValueOf(&val)); err != nil {
		return err
	}
	u.Set(val)
	return nil
}

func (u *UintValue) MapInto(m *anymapper.Mapper, dest reflect.Value) error {
	return m.MapRefl(reflect.ValueOf(&u.Int), dest)
}

type IntValue struct {
	big.Int
	Size int
}

func (i *IntValue) IsDynamic() bool {
	return false
}

func (i *IntValue) EncodeABI() (Words, error) {
	if i.Size < 8 || i.Size > 256 || i.Size%8 != 0 {
		return nil, fmt.Errorf("abi: invalid int size: %d", i.Size)
	}
	return encodeInt(&i.Int, i.Size)
}

func (i *IntValue) DecodeABI(words Words) (int, error) {
	if i.Size < 8 || i.Size > 256 || i.Size%8 != 0 {
		return 0, fmt.Errorf("abi: invalid int size: %d", i.Size)
	}
	return decodeInt(&i.Int, words, i.Size)
}

func (i *IntValue) MapFrom(m *anymapper.Mapper, src reflect.Value) error {
	var val *big.Int
	if err := m.MapRefl(src, reflect.ValueOf(&val)); err != nil {
		return err
	}
	i.Set(val)
	return nil
}

func (i *IntValue) MapInto(m *anymapper.Mapper, dest reflect.Value) error {
	return m.MapRefl(reflect.ValueOf(&i.Int), dest)
}

type BoolValue bool

func (b *BoolValue) IsDynamic() bool {
	return false
}

func (b *BoolValue) EncodeABI() (Words, error) {
	return encodeBool(bool(*b)), nil
}

func (b *BoolValue) DecodeABI(words Words) (int, error) {
	return decodeBool((*bool)(b), words)
}

func (b *BoolValue) MapFrom(m *anymapper.Mapper, src reflect.Value) error {
	var val bool
	if err := m.MapRefl(src, reflect.ValueOf(&val)); err != nil {
		return err
	}
	*b = BoolValue(val)
	return nil
}

func (b *BoolValue) MapInto(m *anymapper.Mapper, dest reflect.Value) error {
	val := bool(*b)
	return m.MapRefl(reflect.ValueOf(val), dest)
}

type AddressValue types.Address

func (a *AddressValue) IsDynamic() bool {
	return false
}

func (a *AddressValue) EncodeABI() (Words, error) {
	return encodeAddress(types.Address(*a))
}

func (a *AddressValue) DecodeABI(words Words) (int, error) {
	return decodeAddress((*types.Address)(a), words)
}

func (a *AddressValue) MapFrom(m *anymapper.Mapper, src reflect.Value) error {
	if !m.StrictTypes && src.Type().Kind() == reflect.String {
		addr, err := types.HexToAddress(src.String())
		if err != nil {
			return fmt.Errorf("abi: cannot convert string to address: %v", err)
		}
		*a = AddressValue(addr)
		return nil
	}
	var addr types.Address
	if err := m.MapRefl(src, reflect.ValueOf(&addr)); err != nil {
		return err
	}
	*a = AddressValue(addr)
	return nil
}

func (a *AddressValue) MapInto(m *anymapper.Mapper, dest reflect.Value) error {
	if !m.StrictTypes && dest.Type().Kind() == reflect.String {
		dest.SetString(types.Address(*a).String())
		return nil
	}
	addr := types.Address(*a)
	return m.MapRefl(reflect.ValueOf(addr), dest)
}

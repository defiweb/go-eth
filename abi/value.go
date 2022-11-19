package abi

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/defiweb/go-eth/hexutil"
	"github.com/defiweb/go-eth/types"

	"github.com/defiweb/go-anymapper"
)

// Value represents a value that can be marshaled to and from ABI.
//
// https://docs.soliditylang.org/en/develop/abi-spec.html#strict-encoding-mode
type Value interface {
	// DynamicType indicates whether the type is dynamic.
	DynamicType() bool

	// EncodeABI returns the ABI encoding of the value.
	EncodeABI() (Words, error)

	// DecodeABI sets the value from the ABI encoding.
	DecodeABI(Words) (int, error)
}

type TupleValue struct {
	elems []Value
	names []string
}

func (t *TupleValue) Size() int {
	return len(t.elems)
}

func (t *TupleValue) Elements() []Value {
	return t.elems
}

func (t *TupleValue) Names() []string {
	return t.names
}

func (t *TupleValue) Map() map[string]any {
	m := make(map[string]any)
	for i, name := range t.names {
		m[name] = t.elems[i]
	}
	return m
}

func (t *TupleValue) Add(name string, v Value) {
	t.elems = append(t.elems, v)
	t.names = append(t.names, name)
}

func (t *TupleValue) Set(idx int, name string, typ Value) error {
	if idx < 0 || idx >= len(t.elems) {
		return fmt.Errorf("abi: index out of range: %d", idx)
	}
	t.names[idx] = name
	t.elems[idx] = typ
	return nil
}

func (t *TupleValue) Get(idx int) Value {
	if idx < 0 || idx >= len(t.elems) {
		return nil
	}
	return t.elems[idx]
}

func (t *TupleValue) GetName(idx int) string {
	if idx < 0 || idx >= len(t.elems) {
		return ""
	}
	return t.names[idx]
}

func (t *TupleValue) DynamicType() bool {
	for _, p := range t.elems {
		if p.DynamicType() {
			return true
		}
	}
	return false
}

func (t *TupleValue) EncodeABI() (Words, error) {
	return encodeTuple(t.elems)
}

func (t *TupleValue) DecodeABI(words Words) (int, error) {
	return decodeTuple(&t.elems, words)
}

func (t *TupleValue) MapFrom(m *anymapper.Mapper, src reflect.Value) error {
	return m.MapRefl(src, reflect.ValueOf(t.Map()))
}

func (t *TupleValue) MapInto(m *anymapper.Mapper, dest reflect.Value) error {
	return m.MapRefl(reflect.ValueOf(t.Map()), dest)
}

type ArrayValue struct {
	elems []Value
	typ   Type
}

func (a *ArrayValue) Length() int {
	return len(a.elems)
}

func (a *ArrayValue) Elements() []Value {
	return a.elems
}

func (a *ArrayValue) Type() Type {
	return a.typ
}

func (a *ArrayValue) Add(v Value) {
	a.elems = append(a.elems, v)
}

func (a *ArrayValue) Set(idx int, v Value) error {
	if idx < 0 || idx >= len(a.elems) {
		return fmt.Errorf("abi: array index out of bounds")
	}
	a.elems[idx] = v
	return nil
}

func (a *ArrayValue) Get(idx int) (Value, error) {
	if idx < 0 || idx >= len(a.elems) {
		return nil, fmt.Errorf("abi: array index out of bounds")
	}
	if len(a.elems) == 0 {
		return nil, nil
	}
	return a.elems[idx], nil
}

func (a *ArrayValue) DynamicType() bool {
	return true
}

func (a *ArrayValue) EncodeABI() (Words, error) {
	return encodeArray(a.elems)
}

func (a *ArrayValue) DecodeABI(data Words) (int, error) {
	return decodeArray(&a.elems, data, a.typ)
}

func (a *ArrayValue) MapFrom(m *anymapper.Mapper, src reflect.Value) error {
	if src.Kind() != reflect.Slice && src.Kind() != reflect.Array {
		return fmt.Errorf("abi: cannot map array from %s", src.Kind())
	}
	a.elems = make([]Value, src.Len())
	for i := 0; i < src.Len(); i++ {
		a.elems[i] = a.typ.New()
	}
	return m.MapRefl(src, reflect.ValueOf(&a.elems))
}

func (a *ArrayValue) MapInto(m *anymapper.Mapper, dest reflect.Value) error {
	return m.MapRefl(reflect.ValueOf(&a.elems), dest)
}

type FixedArrayValue struct {
	elems []Value
	typ   Type
}

func (a *FixedArrayValue) Size() int {
	return len(a.elems)
}

func (a *FixedArrayValue) Elements() []Value {
	return a.elems
}

func (a *FixedArrayValue) Type() Type {
	return a.typ
}

func (a *FixedArrayValue) Set(idx int, v Value) error {
	if idx < 0 || idx >= len(a.elems) {
		return fmt.Errorf("abi: array index out of bounds")
	}
	a.elems[idx] = v
	return nil
}

func (a *FixedArrayValue) Get(idx int) (Value, error) {
	if idx < 0 || idx >= len(a.elems) {
		return nil, fmt.Errorf("abi: array index out of bounds")
	}
	if len(a.elems) == 0 {
		return nil, nil
	}
	return a.elems[idx], nil
}

func (a *FixedArrayValue) DynamicType() bool {
	return false
}

func (a *FixedArrayValue) EncodeABI() (Words, error) {
	return encodeFixedArray(a.elems)
}

func (a *FixedArrayValue) DecodeABI(data Words) (int, error) {
	return decodeFixedArray(&a.elems, data, a.typ, len(a.elems))
}

func (a *FixedArrayValue) MapFrom(m *anymapper.Mapper, src reflect.Value) error {
	if src.Kind() != reflect.Slice && src.Kind() != reflect.Array {
		return fmt.Errorf("abi: cannot map %s to array[%d]", src.Type(), len(a.elems))
	}
	if src.Len() != len(a.elems) {
		return fmt.Errorf("abi: cannot map %d elements to array[%d]", src.Len(), len(a.elems))
	}
	for i := 0; i < len(a.elems); i++ {
		a.elems[i] = a.typ.New()
	}
	return m.MapRefl(src, reflect.ValueOf(&a.elems))
}

func (a *FixedArrayValue) MapInto(m *anymapper.Mapper, dest reflect.Value) error {
	return m.MapRefl(reflect.ValueOf(&a.elems), dest)
}

type BytesValue struct {
	data []byte
}

func (b *BytesValue) Length() int {
	return len(b.data)
}

func (b *BytesValue) Bytes() []byte {
	return b.data
}

func (b *BytesValue) String() string {
	return string(b.data)
}

func (b *BytesValue) Hex() string {
	return hexutil.BytesToHex(b.data)
}

func (b *BytesValue) SetBytes(d []byte) {
	b.data = d
}

func (b *BytesValue) SetString(s string) {
	b.data = []byte(s)
}

func (b *BytesValue) SetHex(s string) error {
	data, err := hexutil.HexToBytes(s)
	if err != nil {
		return err
	}
	b.data = data
	return nil
}

func (b *BytesValue) DynamicType() bool {
	return true
}

func (b *BytesValue) EncodeABI() (Words, error) {
	return encodeBytes(b.data)
}

func (b *BytesValue) DecodeABI(data Words) (int, error) {
	return decodeBytes(&b.data, data)
}

func (b *BytesValue) MapFrom(m *anymapper.Mapper, src reflect.Value) error {
	return m.MapRefl(src, reflect.ValueOf(&b.data))
}

func (b *BytesValue) MapInto(m *anymapper.Mapper, dest reflect.Value) error {
	return m.MapRefl(reflect.ValueOf(&b.data), dest)
}

type StringValue struct {
	data []byte
}

func (s *StringValue) Length() int {
	return len(s.data)
}

func (s *StringValue) Bytes() []byte {
	return s.data
}

func (s *StringValue) String() string {
	return string(s.data)
}

func (s *StringValue) SetBytes(v []byte) {
	s.data = v
}

func (s *StringValue) SetString(v string) {
	s.data = []byte(v)
}

func (s *StringValue) DynamicType() bool {
	return true
}

func (s *StringValue) EncodeABI() (Words, error) {
	return encodeBytes(s.data)
}

func (s *StringValue) DecodeABI(data Words) (int, error) {
	return decodeBytes(&s.data, data)
}

func (s *StringValue) MapFrom(m *anymapper.Mapper, src reflect.Value) error {
	return m.MapRefl(src, reflect.ValueOf(&s.data))
}

func (s *StringValue) MapInto(m *anymapper.Mapper, dest reflect.Value) error {
	return m.MapRefl(reflect.ValueOf(&s.data), dest)
}

type FixedBytesValue struct {
	data []byte
}

func (b *FixedBytesValue) Size() int {
	return len(b.data)
}

func (b *FixedBytesValue) Bytes() []byte {
	return b.data
}

func (b *FixedBytesValue) String() string {
	return string(b.data)
}

func (b *FixedBytesValue) Hex() string {
	return hexutil.BytesToHex(b.data)
}

func (b *FixedBytesValue) SetBytesPadRight(d []byte) error {
	if len(d) > len(b.data) {
		return fmt.Errorf("abi: cannot set %d bytes into bytes%d", len(d), len(b.data))
	}
	copy(b.data, d)
	for i := len(d); i < len(b.data); i++ {
		b.data[i] = 0
	}
	return nil
}

func (b *FixedBytesValue) SetBytesPadLeft(d []byte) error {
	if len(d) > len(b.data) {
		return fmt.Errorf("abi: cannot set %d bytes into bytes%d", len(d), len(b.data))
	}
	copy(b.data[len(b.data)-len(d):], d)
	for i := 0; i < len(b.data)-len(d); i++ {
		b.data[i] = 0
	}
	return nil
}

func (b *FixedBytesValue) SetString(s string) error {
	return b.SetBytesPadRight([]byte(s))
}

func (b *FixedBytesValue) SetHex(s string) error {
	data, err := hexutil.HexToBytes(s)
	if err != nil {
		return err
	}
	return b.SetBytesPadLeft(data)
}

func (b *FixedBytesValue) DynamicType() bool {
	return false
}

func (b *FixedBytesValue) EncodeABI() (Words, error) {
	return encodeFixedBytes(b.data)
}

func (b *FixedBytesValue) DecodeABI(data Words) (int, error) {
	return decodeFixedBytes(&b.data, data, len(b.data))
}

func (b *FixedBytesValue) MapFrom(m *anymapper.Mapper, src reflect.Value) error {
	var data []byte
	if err := m.MapRefl(src, reflect.ValueOf(&b)); err != nil {
		return err
	}
	return b.SetBytesPadLeft(data)
}

func (b *FixedBytesValue) MapInto(m *anymapper.Mapper, dest reflect.Value) error {
	return m.MapRefl(reflect.ValueOf(&b.data), dest)
}

type UintValue struct {
	val  *big.Int
	size int
}

func (u *UintValue) Bytes() []byte {
	return u.val.Bytes()
}

func (u *UintValue) String() string {
	return u.val.String()
}

func (u *UintValue) Hex() string {
	return hexutil.BigIntToHex(u.val)
}

func (u *UintValue) Uint64() (uint64, error) {
	if u.size > 8 {
		return 0, fmt.Errorf("abi: cannot convert uint%d to uint64", u.size*8)
	}
	return u.val.Uint64(), nil
}

func (u *UintValue) SetBytes(d []byte) error {
	if len(d) > u.size {
		return fmt.Errorf("abi: cannot set %d bytes into uint%d", len(d), u.size*8)
	}
	u.val.SetBytes(d)
	return nil
}

func (u *UintValue) SetHex(s string) error {
	data, err := hexutil.HexToBytes(s)
	if err != nil {
		return err
	}
	return u.SetBytes(data)
}

func (u *UintValue) SetBigInt(i *big.Int) error {
	if i.BitLen() > u.size*8 {
		return fmt.Errorf("abi: cannot set %d-bit integer into uint%d", i.BitLen(), u.size*8)
	}
	u.val.Set(i)
	return nil
}

func (u *UintValue) SetUint64(i uint64) error {
	return u.SetBigInt(new(big.Int).SetUint64(i))
}

func (u *UintValue) DynamicType() bool {
	return false
}

func (u *UintValue) EncodeABI() (Words, error) {
	return encodeUint(u.val, u.size)
}

func (u *UintValue) DecodeABI(words Words) (int, error) {
	return decodeUint(u.val, words)
}

func (u *UintValue) MapFrom(m *anymapper.Mapper, src reflect.Value) error {
	return m.MapRefl(src, reflect.ValueOf(&u.val))
}

func (u *UintValue) MapInto(m *anymapper.Mapper, dest reflect.Value) error {
	return m.MapRefl(reflect.ValueOf(&u.val), dest)
}

type IntValue struct {
	val  *big.Int
	size int
}

func (u *IntValue) Bytes() []byte {
	return u.val.Bytes()
}

func (u *IntValue) String() string {
	return u.val.String()
}

func (u *IntValue) Hex() string {
	return hexutil.BigIntToHex(u.val)
}

func (u *IntValue) Int64() (int64, error) {
	if u.size > 8 {
		return 0, fmt.Errorf("abi: cannot convert int%d to int64", u.size*8)
	}
	return u.val.Int64(), nil
}

func (u *IntValue) SetBytes(d []byte) error {
	if len(d) > u.size {
		return fmt.Errorf("abi: cannot set %d bytes into int%d", len(d), u.size*8)
	}
	u.val.SetBytes(d)
	return nil
}

func (u *IntValue) SetHex(s string) error {
	data, err := hexutil.HexToBytes(s)
	if err != nil {
		return err
	}
	return u.SetBytes(data)
}

func (u *IntValue) SetBigInt(i *big.Int) error {
	if signedBitLen(i) > u.size*8 {
		return fmt.Errorf("abi: cannot set %d-bit integer into int%d", i.BitLen(), u.size*8)
	}
	u.val.Set(i)
	return nil
}

func (u *IntValue) SetInt64(i int64) error {
	return u.SetBigInt(new(big.Int).SetInt64(i))
}

func (u *IntValue) DynamicType() bool {
	return false
}

func (u *IntValue) EncodeABI() (Words, error) {
	return encodeInt(u.val, u.size)
}

func (u *IntValue) DecodeABI(words Words) (int, error) {
	return decodeInt(u.val, words)
}

func (u *IntValue) MapFrom(m *anymapper.Mapper, src reflect.Value) error {
	return m.MapRefl(src, reflect.ValueOf(&u.val))
}

func (u *IntValue) MapInto(m *anymapper.Mapper, dest reflect.Value) error {
	return m.MapRefl(reflect.ValueOf(&u.val), dest)
}

type BoolValue bool

func (b *BoolValue) Bool() bool {
	return bool(*b)
}

func (b *BoolValue) SetBool(v bool) {
	*b = BoolValue(v)
}

func (b *BoolValue) DynamicType() bool {
	return false
}

func (b *BoolValue) EncodeABI() (Words, error) {
	return encodeBool(bool(*b)), nil
}

func (b *BoolValue) DecodeABI(words Words) (int, error) {
	return decodeBool((*bool)(b), words)
}

func (b *BoolValue) MapFrom(m *anymapper.Mapper, src reflect.Value) error {
	return m.MapRefl(src, reflect.ValueOf(b))
}

func (b *BoolValue) MapInto(m *anymapper.Mapper, dest reflect.Value) error {
	return m.MapRefl(reflect.ValueOf(b), dest)
}

type AddressValue types.Address

func (a *AddressValue) Address() types.Address {
	return types.Address(*a)
}

func (a *AddressValue) SetAddress(addr types.Address) {
	*a = AddressValue(addr)
}

func (a *AddressValue) DynamicType() bool {
	return false
}

func (a *AddressValue) EncodeABI() (Words, error) {
	var w Word
	copy(w[WordLength-types.AddressLength:], a[:])
	return Words{w}, nil
}

func (a *AddressValue) DecodeABI(words Words) (int, error) {
	if len(words) == 0 {
		return 0, fmt.Errorf("abi: cannot unmarshal address from empty value")
	}
	copy(a[:], words[0][WordLength-types.AddressLength:])
	return 1, nil
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
	return m.MapRefl(src, reflect.ValueOf(a))
}

func (a *AddressValue) MapInto(m *anymapper.Mapper, dest reflect.Value) error {
	if !m.StrictTypes && dest.Type().Kind() == reflect.String {
		dest.SetString(a.Address().String())
		return nil
	}
	return m.MapRefl(reflect.ValueOf(a), dest)
}

package abi

import (
	"fmt"
	"math/big"
	"reflect"
	"strconv"

	"github.com/defiweb/go-eth/hexutil"
	"github.com/defiweb/go-eth/types"

	"github.com/defiweb/go-anymapper"
)

type TupleType struct {
	elems []Type
	names []string
}

func NewTuple() *TupleType {
	return &TupleType{}
}

func NewTupleOfSize(size int) *TupleType {
	return &TupleType{
		elems: make([]Type, size),
		names: make([]string, size),
	}
}

func NewTupleOfElements(elems ...Type) *TupleType {
	names := make([]string, len(elems))
	for i := range elems {
		names[i] = strconv.Itoa(i)
	}
	return &TupleType{
		elems: elems,
		names: names,
	}
}

func (t *TupleType) Elements() []Type {
	return t.elems
}

func (t *TupleType) ElementsMap() map[string]any {
	m := make(map[string]any)
	for i, name := range t.names {
		m[name] = t.elems[i]
	}
	return m
}

func (t *TupleType) Length() int {
	return len(t.elems)
}

func (t *TupleType) Add(name string, v Type) {
	t.elems = append(t.elems, v)
	t.names = append(t.names, name)
}

func (t *TupleType) Set(idx int, name string, typ Type) error {
	if idx < 0 || idx >= len(t.elems) {
		return fmt.Errorf("abi: index out of range: %d", idx)
	}
	t.names[idx] = name
	t.elems[idx] = typ
	return nil
}

func (t *TupleType) Get(idx int) Type {
	if idx < 0 || idx >= len(t.elems) {
		return nil
	}
	return t.elems[idx]
}

func (t *TupleType) Name(idx int) string {
	if idx < 0 || idx >= len(t.elems) {
		return ""
	}
	return t.names[idx]
}

func (t *TupleType) DynamicType() bool {
	for _, p := range t.elems {
		if p.DynamicType() {
			return true
		}
	}
	return false
}

func (t *TupleType) EncodeABI() (Words, error) {
	return encodeTuple(t.elems)
}

func (t *TupleType) DecodeABI(words Words) (int, error) {
	return decodeTuple(&t.elems, words)
}

func (t *TupleType) MapFrom(m *anymapper.Mapper, src reflect.Value) error {
	return m.MapRefl(src, reflect.ValueOf(t.ElementsMap()))
}

func (t *TupleType) MapInto(m *anymapper.Mapper, dest reflect.Value) error {
	return m.MapRefl(reflect.ValueOf(t.ElementsMap()), dest)
}

type ArrayType struct {
	elems []Type
	def   TypeDefinition
	cfg   *Config
}

func NewArray(param TypeDefinition) *ArrayType {
	return &ArrayType{
		def: param,
		cfg: DefaultConfig,
	}
}

func NewArrayOfElements(elems ...Type) *ArrayType {
	return &ArrayType{
		elems: elems,
		cfg:   DefaultConfig,
	}
}

func NewArrayOfSize(size int) *ArrayType {
	return &ArrayType{
		elems: make([]Type, size),
		cfg:   DefaultConfig,
	}
}

func (a *ArrayType) Elements() []Type {
	return a.elems
}

func (a *ArrayType) Type() TypeDefinition {
	return a.def
}

func (a *ArrayType) Size() int {
	return len(a.elems)
}

func (a *ArrayType) Add(v Type) {
	a.elems = append(a.elems, v)
}

func (a *ArrayType) Set(idx int, v Type) error {
	if idx < 0 || idx >= len(a.elems) {
		return fmt.Errorf("abi: array index out of bounds")
	}
	a.elems[idx] = v
	return nil
}

func (a *ArrayType) Get(idx int) (Type, error) {
	if idx < 0 || idx >= len(a.elems) {
		return nil, fmt.Errorf("abi: array index out of bounds")
	}
	if len(a.elems) == 0 {
		return nil, nil
	}
	return a.elems[idx], nil
}

func (a *ArrayType) SetConfig(c *Config) {
	a.cfg = c
}

func (a *ArrayType) DynamicType() bool {
	return true
}

func (a *ArrayType) EncodeABI() (Words, error) {
	return encodeArray(a.elems)
}

func (a *ArrayType) DecodeABI(data Words) (int, error) {
	return decodeArray(&a.elems, data, a.cfg, a.def)
}

func (a *ArrayType) MapFrom(m *anymapper.Mapper, src reflect.Value) error {
	return m.MapRefl(src, reflect.ValueOf(&a.elems))
}

func (a *ArrayType) MapInto(m *anymapper.Mapper, dest reflect.Value) error {
	return m.MapRefl(reflect.ValueOf(&a.elems), dest)
}

type FixedArrayType struct {
	elems []Type
	def   TypeDefinition
	cfg   *Config
}

func NewFixedArray(def TypeDefinition, size int) *FixedArrayType {
	if size < 0 {
		panic(fmt.Errorf("abi: negative array size"))
	}
	return &FixedArrayType{
		elems: make([]Type, size),
		def:   def,
		cfg:   DefaultConfig,
	}
}

func (a *FixedArrayType) Elements() []Type {
	return a.elems
}

func (a *FixedArrayType) Type() TypeDefinition {
	return a.def
}

func (a *FixedArrayType) Size() int {
	return len(a.elems)
}

func (a *FixedArrayType) Set(idx int, v Type) error {
	if idx < 0 || idx >= len(a.elems) {
		return fmt.Errorf("abi: array index out of bounds")
	}
	a.elems[idx] = v
	return nil
}

func (a *FixedArrayType) Get(idx int) (Type, error) {
	if idx < 0 || idx >= len(a.elems) {
		return nil, fmt.Errorf("abi: array index out of bounds")
	}
	if len(a.elems) == 0 {
		return nil, nil
	}
	return a.elems[idx], nil
}

func (a *FixedArrayType) SetConfig(c *Config) {
	a.cfg = c
}

func (a *FixedArrayType) DynamicType() bool {
	return false
}

func (a *FixedArrayType) EncodeABI() (Words, error) {
	return encodeFixedArray(a.elems)
}

func (a *FixedArrayType) DecodeABI(data Words) (int, error) {
	return decodeFixedArray(&a.elems, data, a.cfg, a.def, len(a.elems))
}

func (a *FixedArrayType) MapFrom(m *anymapper.Mapper, src reflect.Value) error {
	return m.MapRefl(src, reflect.ValueOf(&a.elems))
}

func (a *FixedArrayType) MapInto(m *anymapper.Mapper, dest reflect.Value) error {
	return m.MapRefl(reflect.ValueOf(&a.elems), dest)
}

type BytesType struct {
	data []byte
}

func NewBytes() *BytesType {
	return &BytesType{}
}

func (b *BytesType) Bytes() []byte {
	return b.data
}

func (b *BytesType) String() string {
	return string(b.data)
}

func (b *BytesType) Hex() string {
	return hexutil.BytesToHex(b.data)
}

func (b *BytesType) BigInt() *big.Int {
	return new(big.Int).SetBytes(b.data)
}

func (b *BytesType) Uint64() (uint64, error) {
	return decodeUint64(b.data)
}

func (b *BytesType) SetBytes(d []byte) {
	b.data = d
}

func (b *BytesType) SetString(s string) {
	b.data = []byte(s)
}

func (b *BytesType) SetHex(s string) error {
	data, err := hexutil.HexToBytes(s)
	if err != nil {
		return err
	}
	b.data = data
	return nil
}

func (b *BytesType) SetBigInt(i *big.Int) {
	if i == nil || i.Sign() == 0 {
		*b = BytesType{}
		return
	}
	b.data = i.Bytes()
}

func (b *BytesType) SetUint64(i uint64) {
	b.data = encodeUint64(i)
}

func (b *BytesType) DynamicType() bool {
	return true
}

func (b *BytesType) EncodeABI() (Words, error) {
	return encodeBytes(b.data)
}

func (b *BytesType) DecodeABI(data Words) (int, error) {
	return decodeBytes(&b.data, data)
}

func (b *BytesType) MapFrom(m *anymapper.Mapper, src reflect.Value) error {
	return m.MapRefl(src, reflect.ValueOf(&b.data))
}

func (b *BytesType) MapInto(m *anymapper.Mapper, dest reflect.Value) error {
	return m.MapRefl(reflect.ValueOf(&b.data), dest)
}

type StringType struct {
	data []byte
}

func NewString() *StringType {
	return &StringType{}
}

func (s *StringType) Bytes() []byte {
	return s.data
}

func (s *StringType) String() string {
	return string(s.data)
}

func (s *StringType) SetBytes(v []byte) {
	s.data = v
}

func (s *StringType) SetString(v string) {
	s.data = []byte(v)
}

func (s *StringType) DynamicType() bool {
	return true
}

func (s *StringType) EncodeABI() (Words, error) {
	return encodeBytes(s.data)
}

func (s *StringType) DecodeABI(data Words) (int, error) {
	return decodeBytes(&s.data, data)
}

func (s *StringType) MapFrom(m *anymapper.Mapper, src reflect.Value) error {
	return m.MapRefl(src, reflect.ValueOf(&s.data))
}

func (s *StringType) MapInto(m *anymapper.Mapper, dest reflect.Value) error {
	return m.MapRefl(reflect.ValueOf(&s.data), dest)
}

func NewFixedBytes(size int) *FixedBytesType {
	if size < 0 || size > 32 {
		panic(fmt.Sprintf("abi: invalid fixed bytes size %d", size))
	}
	return &FixedBytesType{data: make([]byte, size)}
}

type FixedBytesType struct {
	data []byte
}

func (b *FixedBytesType) Bytes() []byte {
	return b.data
}

func (b *FixedBytesType) String() string {
	return string(b.data)
}

func (b *FixedBytesType) Hex() string {
	return hexutil.BytesToHex(b.data)
}

func (b *FixedBytesType) BigInt() *big.Int {
	return new(big.Int).SetBytes(b.data)
}

func (b *FixedBytesType) Uint64() (uint64, error) {
	return decodeUint64(b.data)
}

func (b *FixedBytesType) SetBytesPadRight(d []byte) error {
	if len(d) > len(b.data) {
		return fmt.Errorf("abi: cannot set %d bytes into bytes%d", len(d), len(b.data))
	}
	copy(b.data, d)
	for i := len(d); i < len(b.data); i++ {
		b.data[i] = 0
	}
	return nil
}

func (b *FixedBytesType) SetBytesPadLeft(d []byte) error {
	if len(d) > len(b.data) {
		return fmt.Errorf("abi: cannot set %d bytes into bytes%d", len(d), len(b.data))
	}
	copy(b.data[len(b.data)-len(d):], d)
	for i := 0; i < len(b.data)-len(d); i++ {
		b.data[i] = 0
	}
	return nil
}

func (b *FixedBytesType) SetString(s string) error {
	return b.SetBytesPadRight([]byte(s))
}

func (b *FixedBytesType) SetHex(s string) error {
	data, err := hexutil.HexToBytes(s)
	if err != nil {
		return err
	}
	return b.SetBytesPadLeft(data)
}

func (b *FixedBytesType) SetBigInt(i *big.Int) error {
	return b.SetBytesPadLeft(i.Bytes())
}

func (b *FixedBytesType) SetUint64(i uint64) error {
	return b.SetBytesPadLeft(encodeUint64(i))
}

func (b *FixedBytesType) DynamicType() bool {
	return false
}

func (b *FixedBytesType) EncodeABI() (Words, error) {
	return encodeFixedBytes(b.data)
}

func (b *FixedBytesType) DecodeABI(data Words) (int, error) {
	return decodeFixedBytes(&b.data, data, len(b.data))
}

func (b *FixedBytesType) MapFrom(m *anymapper.Mapper, src reflect.Value) error {
	return m.MapRefl(src, reflect.ValueOf(&b.data))
}

func (b *FixedBytesType) MapInto(m *anymapper.Mapper, dest reflect.Value) error {
	return m.MapRefl(reflect.ValueOf(&b.data), dest)
}

type UintType struct {
	val  *big.Int
	size int
}

func NewUint(size int) *UintType {
	if size < 0 || size > 32 {
		panic(fmt.Errorf("abi: invalid uint size %d", size))
	}
	return &UintType{val: new(big.Int), size: size}
}

func (u *UintType) Bytes() []byte {
	return u.val.Bytes()
}

func (u *UintType) String() string {
	return u.val.String()
}

func (u *UintType) Hex() string {
	return hexutil.BigIntToHex(u.val)
}

func (u *UintType) Uint64() (uint64, error) {
	if u.size > 8 {
		return 0, fmt.Errorf("abi: cannot convert uint%d to uint64", u.size*8)
	}
	return u.val.Uint64(), nil
}

func (u *UintType) SetBytes(d []byte) error {
	if len(d) > u.size {
		return fmt.Errorf("abi: cannot set %d bytes into uint%d", len(d), u.size*8)
	}
	u.val.SetBytes(d)
	return nil
}

func (u *UintType) SetHex(s string) error {
	data, err := hexutil.HexToBytes(s)
	if err != nil {
		return err
	}
	return u.SetBytes(data)
}

func (u *UintType) SetBigInt(i *big.Int) error {
	if i.BitLen() > u.size*8 {
		return fmt.Errorf("abi: cannot set %d-bit integer into uint%d", i.BitLen(), u.size*8)
	}
	u.val.Set(i)
	return nil
}

func (u *UintType) SetUint64(i uint64) error {
	return u.SetBigInt(new(big.Int).SetUint64(i))
}

func (u *UintType) DynamicType() bool {
	return false
}

func (u *UintType) EncodeABI() (Words, error) {
	return encodeUint(u.val, u.size)
}

func (u *UintType) DecodeABI(words Words) (int, error) {
	return decodeUint(u.val, words)
}

func (u *UintType) MapFrom(m *anymapper.Mapper, src reflect.Value) error {
	return m.MapRefl(src, reflect.ValueOf(&u.val))
}

func (u *UintType) MapInto(m *anymapper.Mapper, dest reflect.Value) error {
	return m.MapRefl(reflect.ValueOf(&u.val), dest)
}

type IntType struct {
	val  *big.Int
	size int
}

func NewInt(size int) *IntType {
	if size < 0 || size > 32 {
		panic(fmt.Errorf("abi: invalid uint size %d", size))
	}
	return &IntType{val: new(big.Int), size: size}
}

func (u *IntType) Bytes() []byte {
	return u.val.Bytes()
}

func (u *IntType) String() string {
	return u.val.String()
}

func (u *IntType) Hex() string {
	return hexutil.BigIntToHex(u.val)
}

func (u *IntType) Int64() (int64, error) {
	if u.size > 8 {
		return 0, fmt.Errorf("abi: cannot convert int%d to int64", u.size*8)
	}
	return u.val.Int64(), nil
}

func (u *IntType) SetBytes(d []byte) error {
	if len(d) > u.size {
		return fmt.Errorf("abi: cannot set %d bytes into int%d", len(d), u.size*8)
	}
	u.val.SetBytes(d)
	return nil
}

func (u *IntType) SetHex(s string) error {
	data, err := hexutil.HexToBytes(s)
	if err != nil {
		return err
	}
	return u.SetBytes(data)
}

func (u *IntType) SetBigInt(i *big.Int) error {
	if i.BitLen() > u.size*8 {
		return fmt.Errorf("abi: cannot set %d-bit integer into int%d", i.BitLen(), u.size*8)
	}
	u.val.Set(i)
	return nil
}

func (u *IntType) SetInt64(i int64) error {
	return u.SetBigInt(new(big.Int).SetInt64(i))
}

func (u *IntType) DynamicType() bool {
	return false
}

func (u *IntType) EncodeABI() (Words, error) {
	return encodeInt(u.val, u.size)
}

func (u *IntType) DecodeABI(words Words) (int, error) {
	return decodeInt(u.val, words)
}

func (u *IntType) MapFrom(m *anymapper.Mapper, src reflect.Value) error {
	return m.MapRefl(src, reflect.ValueOf(&u.val))
}

func (u *IntType) MapInto(m *anymapper.Mapper, dest reflect.Value) error {
	return m.MapRefl(reflect.ValueOf(&u.val), dest)
}

type BoolType bool

func NewBool() *BoolType {
	return new(BoolType)
}

func (b *BoolType) Bool() bool {
	return bool(*b)
}

func (b *BoolType) SetBool(v bool) {
	*b = BoolType(v)
}

func (b *BoolType) DynamicType() bool {
	return false
}

func (b *BoolType) EncodeABI() (Words, error) {
	return encodeBool(bool(*b)), nil
}

func (b *BoolType) DecodeABI(words Words) (int, error) {
	return decodeBool((*bool)(b), words)
}

func (b *BoolType) MapFrom(m *anymapper.Mapper, src reflect.Value) error {
	return m.MapRefl(src, reflect.ValueOf(b))
}

func (b *BoolType) MapInto(m *anymapper.Mapper, dest reflect.Value) error {
	return m.MapRefl(reflect.ValueOf(b), dest)
}

type AddressType types.Address

func NewAddress() *AddressType {
	return (*AddressType)(&types.Address{})
}

func (a *AddressType) Address() types.Address {
	return types.Address(*a)
}

func (a *AddressType) SetAddress(addr types.Address) {
	*a = AddressType(addr)
}

func (a *AddressType) DynamicType() bool {
	return false
}

func (a *AddressType) EncodeABI() (Words, error) {
	var w Word
	copy(w[WordLength-types.AddressLength:], a[:])
	return Words{w}, nil
}

func (a *AddressType) DecodeABI(words Words) (int, error) {
	if len(words) == 0 {
		return 0, fmt.Errorf("abi: cannot unmarshal address from empty value")
	}
	copy(a[:], words[0][WordLength-types.AddressLength:])
	return 1, nil
}

func (a *AddressType) MapFrom(m *anymapper.Mapper, src reflect.Value) error {
	if !m.StrictTypes {
		switch src.Kind() {
		case reflect.String:
			addr, err := types.HexToAddress(src.String())
			if err != nil {
				return anymapper.NewInvalidMappingError(src.Type(), reflect.TypeOf(a), err.Error())
			}
			*a = AddressType(addr)
			return nil
		case reflect.Slice:
			if src.Type().Elem().Kind() == reflect.Uint8 {
				addr, err := types.BytesToAddress(src.Bytes())
				if err != nil {
					return anymapper.NewInvalidMappingError(src.Type(), reflect.TypeOf(a), err.Error())
				}
				*a = AddressType(addr)
				return nil
			}
		case reflect.Array:
			if src.Type().Elem().Kind() == reflect.Uint8 {
				if src.Len() != types.AddressLength {
					return anymapper.NewInvalidMappingError(src.Type(), reflect.TypeOf(a), "array length must be 20")
				}
				var addr types.Address
				copy(addr[:], src.Bytes())
				*a = AddressType(addr)
				return nil
			}
		}
	}
	return m.MapRefl(src, reflect.ValueOf(a))
}

func (a *AddressType) MapInto(m *anymapper.Mapper, dest reflect.Value) error {
	if !m.StrictTypes {
		switch dest.Kind() {
		case reflect.String:
			dest.SetString(a.Address().String())
			return nil
		case reflect.Slice:
			if dest.Type().Elem().Kind() == reflect.Uint8 {
				dest.SetBytes(a.Address().Bytes())
				return nil
			}
		case reflect.Array:
			if dest.Type().Elem().Kind() == reflect.Uint8 {
				if dest.Len() != types.AddressLength {
					return anymapper.NewInvalidMappingError(reflect.TypeOf(a), dest.Type(), "array length must be 20")
				}
				for i := 0; i < types.AddressLength; i++ {
					dest.Index(i).SetUint(uint64(a[i]))
				}
				return nil
			}
		}
	}
	return m.MapRefl(reflect.ValueOf(a), dest)
}

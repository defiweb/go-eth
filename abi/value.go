package abi

import (
	"fmt"
	"math"
	"math/big"
	"reflect"

	"github.com/defiweb/go-eth/hexutil"
	"github.com/defiweb/go-eth/types"
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
func (t *TupleValue) MapFrom(m Mapper, src any) error {
	vals := make(map[string]Value, len(*t))
	for _, elem := range *t {
		vals[elem.Name] = elem.Value
	}
	if err := m.Map(src, vals); err != nil {
		return fmt.Errorf("abi: cannot map tuple from %s: %w", reflect.TypeOf(src), err)
	}
	return nil
}

// MapTo implements the anymapper.MapTo interface.
func (t *TupleValue) MapTo(m Mapper, dst any) error {
	vals := make(map[string]Value, len(*t))
	for _, elem := range *t {
		vals[elem.Name] = elem.Value
	}
	if err := m.Map(vals, dst); err != nil {
		return fmt.Errorf("abi: cannot map tuple to %s: %w", reflect.TypeOf(dst), err)
	}
	return nil
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
func (a *ArrayValue) MapFrom(m Mapper, src any) error {
	srcRef := reflect.ValueOf(src)
	if srcRef.Kind() != reflect.Slice && srcRef.Kind() != reflect.Array {
		return fmt.Errorf("abi: cannot map array from %s", srcRef.Kind())
	}
	a.Elems = make([]Value, srcRef.Len())
	for i := 0; i < srcRef.Len(); i++ {
		a.Elems[i] = a.Type.Value()
	}
	if err := m.Map(src, &a.Elems); err != nil {
		return fmt.Errorf("abi: cannot map array from %s: %w", reflect.TypeOf(src), err)
	}
	return nil
}

// MapTo implements the anymapper.MapTo interface.
func (a *ArrayValue) MapTo(m Mapper, dst any) error {
	if err := m.Map(&a.Elems, dst); err != nil {
		return fmt.Errorf("abi: cannot map array to %s: %w", reflect.TypeOf(dst), err)
	}
	return nil
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
func (a FixedArrayValue) MapFrom(m Mapper, src any) error {
	srcRef := reflect.ValueOf(src)
	if srcRef.Kind() != reflect.Slice && srcRef.Kind() != reflect.Array {
		return fmt.Errorf("abi: cannot map %s to array[%d]", srcRef.Type(), len(a))
	}
	if srcRef.Len() != len(a) {
		return fmt.Errorf("abi: cannot map %d elements to array[%d]", srcRef.Len(), len(a))
	}
	if err := m.Map(src, (*[]Value)(&a)); err != nil {
		return fmt.Errorf("abi: cannot map array from %s: %w", reflect.TypeOf(src), err)
	}
	return nil
}

// MapTo implements the anymapper.MapTo interface.
func (a FixedArrayValue) MapTo(m Mapper, dst any) error {
	if err := m.Map(([]Value)(a), dst); err != nil {
		return fmt.Errorf("abi: cannot map array to %s: %w", reflect.TypeOf(dst), err)
	}
	return nil
}

// BytesValue is a value of bytes type.
//
// During encoding ad decoding, the BytesValue can be mapped using the slice
// rules described in the documentation of anymapper package.
type BytesValue []byte

// Bytes returns the value of the BytesValue.
func (b *BytesValue) Bytes() []byte {
	return *b
}

// SetBytes sets the value of the BytesValue.
func (b *BytesValue) SetBytes(data []byte) {
	*b = data
}

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
func (b *BytesValue) MapFrom(m Mapper, src any) error {
	srcRef := reflect.ValueOf(src)
	switch srcRef.Type().Kind() {
	case reflect.Slice, reflect.Array:
		if srcRef.Type().Elem().Kind() != reflect.Uint8 {
			return fmt.Errorf("abi: cannot map %s to bytes", srcRef.Type())
		}
		if err := m.Map(src, (*[]byte)(b)); err != nil {
			return fmt.Errorf("abi: cannot map %s to bytes: %v", srcRef.Type(), err)
		}
	case reflect.String:
		bin, err := hexutil.HexToBytes(srcRef.String())
		if err != nil {
			return fmt.Errorf("abi: cannot map %s to bytes: %v", srcRef.Type(), err)
		}
		*b = bin
	default:
		return fmt.Errorf("abi: cannot map %s to bytes", srcRef.Type())
	}
	return nil
}

// MapTo implements the anymapper.MapTo interface.
func (b *BytesValue) MapTo(m Mapper, dst any) error {
	dstRef := reflect.ValueOf(dst).Elem()
	switch dstRef.Type().Kind() {
	case reflect.Slice, reflect.Array:
		if dstRef.Type().Elem().Kind() != reflect.Uint8 {
			return fmt.Errorf("abi: cannot map bytes to %s", dstRef.Type())
		}
		if err := m.Map((*[]byte)(b), &dst); err != nil {
			return fmt.Errorf("abi: cannot map bytes to %s: %v", dstRef.Type(), err)
		}
	case reflect.String:
		dstRef.SetString(hexutil.BytesToHex(*b))
	default:
		return fmt.Errorf("abi: cannot map bytes to %s", dstRef.Type())
	}
	return nil
}

// StringValue is a value of bytes type.
//
// During encoding ad decoding, the StringValue is mapped using the string
// rules described in the documentation of anymapper package.
type StringValue string

// String returns the value of the StringValue.
func (s *StringValue) String() string {
	return string(*s)
}

// SetString sets the value of the StringValue.
func (s *StringValue) SetString(str string) {
	*s = StringValue(str)
}

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
func (s *StringValue) MapFrom(m Mapper, src any) error {
	srcRef := reflect.ValueOf(src)
	switch srcRef.Type().Kind() {
	case reflect.Slice:
		if srcRef.Type().Elem().Kind() != reflect.Uint8 {
			return fmt.Errorf("abi: cannot map %s to string", srcRef.Type())
		}
		if err := m.Map(src, (*string)(s)); err != nil {
			return fmt.Errorf("abi: cannot map %s to string: %v", srcRef.Type(), err)
		}
	case reflect.String:
		*s = StringValue(srcRef.String())
	default:
		return fmt.Errorf("abi: cannot map %s to string", srcRef.Type())
	}
	return nil
}

// MapTo implements the anymapper.MapTo interface.
func (s *StringValue) MapTo(m Mapper, dst any) error {
	dstRef := reflect.ValueOf(dst).Elem()
	switch dstRef.Type().Kind() {
	case reflect.Slice:
		if dstRef.Type().Elem().Kind() != reflect.Uint8 {
			return fmt.Errorf("abi: cannot map string to %s", dstRef.Type())
		}
		if err := m.Map((*string)(s), &dst); err != nil {
			return fmt.Errorf("abi: cannot map string to %s: %v", dstRef.Type(), err)
		}
	case reflect.String:
		dstRef.SetString(string(*s))
	default:
		return fmt.Errorf("abi: cannot map string to %s", dstRef.Type())
	}
	return nil
}

// FixedBytesValue is a value of fixed bytes type. The size of a slice is
// assumed to be equal to the size of the bytesN type.
//
// During encoding and decoding, the FixedBytesValue is mapped using the slice
// rules described in the documentation of anymapper package. Both values must
// have the same size.
type FixedBytesValue []byte

// Bytes returns the value of the FixedBytesValue.
func (b *FixedBytesValue) Bytes() []byte {
	return *b
}

// SetBytes sets the value of the FixedBytesValue.
func (b *FixedBytesValue) SetBytes(data []byte) error {
	if len(data) != len(*b) {
		return fmt.Errorf("abi: cannot set bytes of length %d to bytes%d", len(data), len(*b))
	}
	*b = data
	return nil
}

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
func (b FixedBytesValue) MapFrom(m Mapper, src any) error {
	srcRef := reflect.ValueOf(src)
	switch srcRef.Type().Kind() {
	case reflect.Slice:
		if srcRef.Type().Elem().Kind() != reflect.Uint8 {
			return fmt.Errorf("abi: cannot map %s to bytes", srcRef.Type())
		}
		bin := srcRef.Bytes()
		if len(bin) != len(b) {
			return fmt.Errorf("abi: cannot map %d bytes to bytes%d", len(bin), len(b))
		}
		copy(b, bin)
	case reflect.Array:
		if srcRef.Type().Elem().Kind() != reflect.Uint8 {
			return fmt.Errorf("abi: cannot map %s to bytes", srcRef.Type())
		}
		var bin []byte
		if err := m.Map(src, &bin); err != nil {
			return fmt.Errorf("abi: cannot map %s to bytes%d: %v", srcRef.Type(), len(b), err)
		}
		if len(bin) != len(b) {
			return fmt.Errorf("abi: cannot map %d bytes to bytes%d", len(bin), len(b))
		}
		copy(b, bin)
	case reflect.String:
		bin, err := hexutil.HexToBytes(srcRef.String())
		if err != nil {
			return fmt.Errorf("abi: cannot map %s to bytes%d: %v", srcRef.Type(), len(b), err)
		}
		if len(bin) != len(b) {
			return fmt.Errorf("abi: cannot map %d bytes to bytes%d", len(bin), len(b))
		}
		copy(b, bin)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if len(b) != 32 {
			return fmt.Errorf("abi: cannot map %s to bytes%d: only bytes32 is supported", srcRef.Type(), len(b))
		}
		x := newUintX(256)
		_ = x.SetUint64(srcRef.Uint())
		bin := x.Bytes()
		for i := 0; i < len(b)-len(bin); i++ {
			b[i] = 0
		}
		copy(b[len(b)-len(bin):], bin)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if len(b) != 32 {
			return fmt.Errorf("abi: cannot map %s to bytes%d: only bytes32 is supported", srcRef.Type(), len(b))
		}
		x := newIntX(256)
		_ = x.SetInt64(srcRef.Int())
		bin := x.Bytes()
		for i := 0; i < len(b)-len(bin); i++ {
			b[i] = 0
		}
		copy(b[len(b)-len(bin):], bin)
	default:
		switch srcTyp := srcRef.Interface().(type) {
		case big.Int:
			if len(b) != 32 {
				return fmt.Errorf("abi: cannot map %s to bytes%d: only bytes32 is supported", srcRef.Type(), len(b))
			}
			x := newIntX(len(b) * 8)
			if err := x.SetBigInt(&srcTyp); err != nil {
				return fmt.Errorf("abi: cannot map %s to bytes%d: %v", srcRef.Type(), len(b), err)
			}
			bin := x.Bytes()
			for i := 0; i < len(b)-len(bin); i++ {
				b[i] = 0
			}
			copy(b[len(b)-len(bin):], bin)
		case types.Number:
			if len(b) != 32 {
				return fmt.Errorf("abi: cannot map %s to bytes%d: only bytes32 is supported", srcRef.Type(), len(b))
			}
			x := newIntX(len(b) * 8)
			if err := x.SetBigInt(srcTyp.Big()); err != nil {
				return fmt.Errorf("abi: cannot map %s to bytes%d: %v", srcRef.Type(), len(b), err)
			}
			bin := x.Bytes()
			for i := 0; i < len(b)-len(bin); i++ {
				b[i] = 0
			}
			copy(b[len(b)-len(bin):], bin)
		case types.BlockNumber:
			if srcTyp.Big().Sign() < 0 {
				return fmt.Errorf("abi: cannot map %s to bytes%d: latest, earliest and pending are not supported", srcRef.Type(), len(b))
			}
			if len(b) != 32 {
				return fmt.Errorf("abi: cannot map %s to bytes%d: only bytes32 is supported", srcRef.Type(), len(b))
			}
			x := newIntX(len(b) * 8)
			if err := x.SetBigInt(srcTyp.Big()); err != nil {
				return fmt.Errorf("abi: cannot map %s to bytes%d: %v", srcRef.Type(), len(b), err)
			}
			bin := x.Bytes()
			for i := 0; i < len(b)-len(bin); i++ {
				b[i] = 0
			}
			copy(b[len(b)-len(bin):], bin)
		default:
			return fmt.Errorf("abi: cannot map %s to bytes", srcRef.Type())
		}
	}
	return nil
}

// MapTo implements the anymapper.MapTo interface.
func (b FixedBytesValue) MapTo(m Mapper, dst any) error {
	dstRef := reflect.ValueOf(dst).Elem()
	switch dstRef.Type().Kind() {
	case reflect.Slice:
		if dstRef.Type().Elem().Kind() != reflect.Uint8 {
			return fmt.Errorf("abi: cannot map bytes to %s", dstRef.Type())
		}
		dstRef.Set(reflect.ValueOf([]byte(b)))
	case reflect.Array:
		if dstRef.Type().Elem().Kind() != reflect.Uint8 {
			return fmt.Errorf("abi: cannot map bytes to %s", dstRef.Type())
		}
		if dstRef.Len() != len(b) {
			return fmt.Errorf("abi: cannot map bytes%d to %s: length mismatch", len(b), dstRef.Type())
		}
		for i := 0; i < dstRef.Len(); i++ {
			dstRef.Index(i).SetUint(uint64(b[i]))
		}
	case reflect.String:
		dstRef.SetString(hexutil.BytesToHex(b))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if len(b) != 32 {
			return fmt.Errorf("abi: cannot map bytes%d to %s: only bytes32 is supported", len(b), dstRef.Type())
		}
		x := newUintX(256)
		if err := x.SetBytes(b); err != nil {
			return fmt.Errorf("abi: cannot map bytes%d to %s: %v", len(b), dstRef.Type(), err)
		}
		u64, err := x.Uint64()
		if err != nil {
			return fmt.Errorf("abi: cannot map bytes%d to %s: %v", len(b), dstRef.Type(), err)
		}
		dstRef.SetUint(u64)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if len(b) != 32 {
			return fmt.Errorf("abi: cannot map bytes%d to %s: only bytes32 is supported", len(b), dstRef.Type())
		}
		x := newIntX(256)
		if err := x.SetBytes(b); err != nil {
			return fmt.Errorf("abi: cannot map bytes%d to %s: %v", len(b), dstRef.Type(), err)
		}
		i64, err := x.Int64()
		if err != nil {
			return fmt.Errorf("abi: cannot map bytes%d to %s: %v", len(b), dstRef.Type(), err)
		}
		if dstRef.OverflowInt(i64) {
			return fmt.Errorf("abi: cannot map bytes%d to %s: %v", len(b), dstRef.Type(), err)
		}
		dstRef.SetInt(i64)
	default:
		switch dstRef.Interface().(type) {
		case big.Int:
			if len(b) != 32 {
				return fmt.Errorf("abi: cannot map bytes%d to %s: only bytes32 is supported", len(b), dstRef.Type())
			}
			x := newIntX(256)
			if err := x.SetBytes(b); err != nil {
				return fmt.Errorf("abi: cannot map bytes%d to %s: %v", len(b), dstRef.Type(), err)
			}
			dstRef.Set(reflect.ValueOf(x.BigInt()).Elem())
		case types.Number:
			if len(b) != 32 {
				return fmt.Errorf("abi: cannot map bytes%d to %s: only bytes32 is supported", len(b), dstRef.Type())
			}
			x := newIntX(256)
			if err := x.SetBytes(b); err != nil {
				return fmt.Errorf("abi: cannot map bytes%d to %s: %v", len(b), dstRef.Type(), err)
			}
			dstRef.Set(reflect.ValueOf(types.BigIntToNumber(x.BigInt())))
		case types.BlockNumber:
			if len(b) != 32 {
				return fmt.Errorf("abi: cannot map bytes%d to %s: only bytes32 is supported", len(b), dstRef.Type())
			}
			x := newIntX(256)
			if err := x.SetBytes(b); err != nil {
				return fmt.Errorf("abi: cannot map bytes%d to %s: %v", len(b), dstRef.Type(), err)
			}
			dstRef.Set(reflect.ValueOf(types.BigIntToBlockNumber(x.BigInt())))
		default:
			return fmt.Errorf("abi: cannot map bytes%d to %s", len(b), dstRef.Type())
		}
	}
	return nil
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
func (u *UintValue) MapFrom(m Mapper, src any) error {
	srcRef := reflect.ValueOf(src)
	switch srcRef.Type().Kind() {
	case reflect.String:
		bn, err := hexutil.HexToBigInt(srcRef.String())
		if err != nil {
			return fmt.Errorf("abi: cannot map %s to uint%d: %v", srcRef.Type(), u.Size, err)
		}
		if bn.Sign() < 0 {
			return fmt.Errorf("abi: cannot map %s to uint%d: negative value", srcRef.Type(), u.Size)
		}
		if bn.BitLen() > u.Size {
			return fmt.Errorf("abi: cannot map %s to uint%d: value too large", srcRef.Type(), u.Size)
		}
		u.Int = *bn
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i64 := srcRef.Int()
		if i64 < 0 {
			return fmt.Errorf("abi: cannot map %s to uint%d: negative value", srcRef.Type(), u.Size)
		}
		if !canSetUint(uint64(i64), u.Size) {
			return fmt.Errorf("abi: cannot map value to uint%d: value too large", u.Size)
		}
		u.Int.SetInt64(i64)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if !canSetUint(srcRef.Uint(), u.Size) {
			return fmt.Errorf("abi: cannot map value to uint%d: value too large", u.Size)
		}
		u.Int.SetUint64(srcRef.Uint())
	default:
		switch srcTyp := srcRef.Interface().(type) {
		case big.Int:
			if srcTyp.Sign() < 0 {
				return fmt.Errorf("abi: cannot map %s to uint%d: negative value", srcRef.Type(), u.Size)
			}
			if srcTyp.BitLen() > u.Size {
				return fmt.Errorf("abi: cannot map %s to uint%d: value too large", srcRef.Type(), u.Size)
			}
			u.Int = srcTyp
		case types.Number:
			bn := srcTyp.Big()
			if bn.Sign() < 0 {
				return fmt.Errorf("abi: cannot map %s to uint%d: negative value", srcRef.Type(), u.Size)
			}
			if bn.BitLen() > u.Size {
				return fmt.Errorf("abi: cannot map %s to uint%d: value too large", srcRef.Type(), u.Size)
			}
			u.Int = *bn
		case types.BlockNumber:
			bn := srcTyp.Big()
			if bn.Sign() < 0 {
				return fmt.Errorf("abi: cannot map %s to uint%d: negative value", srcRef.Type(), u.Size)
			}
			if bn.BitLen() > u.Size {
				return fmt.Errorf("abi: cannot map %s to uint%d: value too large", srcRef.Type(), u.Size)
			}
			u.Int = *bn
		default:
			return fmt.Errorf("abi: cannot map %s to uint%d", srcRef.Type(), u.Size)
		}
	}
	return nil
}

// MapTo implements the anymapper.MapTo interface.
func (u *UintValue) MapTo(m Mapper, dst any) error {
	dstRef := reflect.ValueOf(dst).Elem()
	switch dstRef.Type().Kind() {
	case reflect.String:
		dstRef.SetString(hexutil.BigIntToHex(&u.Int))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if signedBitLen(&u.Int) > dstRef.Type().Bits() {
			return fmt.Errorf("abi: cannot map uint%d to %s: value too large", u.Size, dstRef.Type())
		}
		dstRef.SetInt(u.Int64())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if u.Int.BitLen() > dstRef.Type().Bits() {
			return fmt.Errorf("abi: cannot map uint%d to %s: value too large", u.Size, dstRef.Type())
		}
		dstRef.SetUint(u.Uint64())
	default:
		switch dstRef.Interface().(type) {
		case big.Int:
			dstRef.Set(reflect.ValueOf(u.Int))
		case types.Number:
			dstRef.Set(reflect.ValueOf(types.BigIntToNumber(&u.Int)))
		case types.BlockNumber:
			dstRef.Set(reflect.ValueOf(types.BigIntToBlockNumber(&u.Int)))
		default:
			return fmt.Errorf("abi: cannot map uint%d to %s", u.Size, dstRef.Type())
		}
	}
	return nil
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
func (i *IntValue) MapFrom(m Mapper, src any) error {
	srcRef := reflect.ValueOf(src)
	switch srcRef.Type().Kind() {
	case reflect.String:
		bn, err := hexutil.HexToBigInt(srcRef.String())
		if err != nil {
			return fmt.Errorf("abi: cannot map %s to int%d: %v", srcRef.Type(), i.Size, err)
		}
		if signedBitLen(bn) > i.Size {
			return fmt.Errorf("abi: cannot map value to int%d: value too large", i.Size)
		}
		i.Int = *bn
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if !canSetInt(srcRef.Int(), i.Size) {
			return fmt.Errorf("abi: cannot map value to int%d: value too large", i.Size)
		}
		i.Int.SetInt64(srcRef.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u64 := srcRef.Uint()
		if u64 > math.MaxInt64 {
			return fmt.Errorf("abi: cannot map value to int%d: value too large", i.Size)
		}
		if !canSetInt(int64(u64), i.Size) {
			return fmt.Errorf("abi: cannot map value to int%d: value too large", i.Size)
		}
		i.Int.SetUint64(u64)
	default:
		switch srcTyp := srcRef.Interface().(type) {
		case big.Int:
			if signedBitLen(&srcTyp) > i.Size {
				return fmt.Errorf("abi: cannot map %s to uint%d: value too large", srcRef.Type(), i.Size)
			}
			i.Int = srcTyp
		case types.Number:
			bn := srcTyp.Big()
			if signedBitLen(bn) > i.Size {
				return fmt.Errorf("abi: cannot map %s to uint%d: value too large", srcRef.Type(), i.Size)
			}
			i.Int = *bn
		case types.BlockNumber:
			bn := srcTyp.Big()
			if bn.Sign() < 0 {
				return fmt.Errorf("abi: cannot map %s to uint%d: latest, earliest and pending are not supported", srcRef.Type(), i.Size)
			}
			if signedBitLen(bn) > i.Size {
				return fmt.Errorf("abi: cannot map %s to uint%d: value too large", srcRef.Type(), i.Size)
			}
			i.Int = *bn
		default:
			return fmt.Errorf("abi: cannot map %s to uint%d", srcRef.Type(), i.Size)
		}
	}
	return nil
}

// MapTo implements the anymapper.MapTo interface.
func (i *IntValue) MapTo(m Mapper, dst any) error {
	dstRef := reflect.ValueOf(dst).Elem()
	switch dstRef.Type().Kind() {
	case reflect.String:
		dstRef.SetString(hexutil.BigIntToHex(&i.Int))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if signedBitLen(&i.Int) > dstRef.Type().Bits() {
			return fmt.Errorf("abi: cannot map int%d to %s: value too large", i.Size, dstRef.Type())
		}
		dstRef.SetInt(i.Int64())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if i.Int.BitLen() > dstRef.Type().Bits() {
			return fmt.Errorf("abi: cannot map int%d to %s: value too large", i.Size, dstRef.Type())
		}
		if i.Sign() < 0 {
			return fmt.Errorf("abi: cannot map int%d to %s: value is negative", i.Size, dstRef.Type())
		}
		dstRef.SetUint(i.Uint64())
	default:
		switch dstRef.Interface().(type) {
		case big.Int:
			dstRef.Set(reflect.ValueOf(i.Int))
		case types.Number:
			dstRef.Set(reflect.ValueOf(types.BigIntToNumber(&i.Int)))
		case types.BlockNumber:
			if i.Sign() < 0 {
				return fmt.Errorf("abi: cannot map int%d to %s: value is negative", i.Size, dstRef.Type())
			}
			dstRef.Set(reflect.ValueOf(types.BigIntToBlockNumber(&i.Int)))
		default:
			return fmt.Errorf("abi: cannot map int%d to %s", i.Size, dstRef.Type())
		}
	}
	return nil
}

// BoolValue is a value of bool type.
//
// During encoding and decoding, the BoolValue is mapped using the bool rules
// described in the documentation of anymapper package.
type BoolValue bool

// SetBool sets the value of the BoolValue.
func (b *BoolValue) SetBool(v bool) {
	*b = BoolValue(v)
}

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
func (b *BoolValue) MapFrom(m Mapper, src any) error {
	srcRef := reflect.ValueOf(src)
	switch srcRef.Type().Kind() {
	case reflect.Bool:
		*b = BoolValue(srcRef.Bool())
	default:
		return fmt.Errorf("abi: cannot map %s to bool", srcRef.Type())
	}
	return nil
}

// MapTo implements the anymapper.MapTo interface.
func (b *BoolValue) MapTo(m Mapper, dst any) error {
	dstRef := reflect.ValueOf(dst).Elem()
	switch dstRef.Type().Kind() {
	case reflect.Bool:
		dstRef.SetBool(bool(*b))
	default:
		return fmt.Errorf("abi: cannot map bool to %s", dstRef.Type())
	}
	return nil
}

// AddressValue is a value of address type.
//
// During encoding, the AddressValue can be mapped to the types.Address type,
// string as a hex-encoded address. For other types, the rules for []byte slice
// described in the documentation of anymapper package are used.
type AddressValue types.Address

// Address returns the address value.
func (a *AddressValue) Address() types.Address {
	return types.Address(*a)
}

// SetAddress sets the value of the AddressValue.
func (a *AddressValue) SetAddress(v types.Address) {
	*a = AddressValue(v)
}

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
func (a *AddressValue) MapFrom(m Mapper, src any) error {
	srcRef := reflect.ValueOf(src)
	switch srcRef.Type().Kind() {
	case reflect.String:
		addr, err := types.HexToAddress(srcRef.String())
		if err != nil {
			return fmt.Errorf("abi: cannot map %s to address: %v", srcRef.Type(), err)
		}
		*a = AddressValue(addr)
	case reflect.Slice, reflect.Array:
		if srcRef.Type().Elem().Kind() != reflect.Uint8 {
			return fmt.Errorf("abi: cannot map %s to address", srcRef.Type())
		}
		if srcRef.Len() != types.AddressLength {
			return fmt.Errorf("abi: cannot map %s to address: length mismatch", srcRef.Type())
		}
		var bin []byte
		if err := m.Map(src, &bin); err != nil {
			return err
		}
		*a = AddressValue(types.MustBytesToAddress(bin))
	default:
		return fmt.Errorf("abi: cannot map %s to address", srcRef.Type())
	}
	return nil
}

// MapTo implements the anymapper.MapTo interface.
func (a *AddressValue) MapTo(m Mapper, dst any) error {
	dstRef := reflect.ValueOf(dst).Elem()
	switch dstRef.Type().Kind() {
	case reflect.String:
		dstRef.SetString(types.Address(*a).String())
	case reflect.Slice:
		if dstRef.Type().Elem().Kind() != reflect.Uint8 {
			return fmt.Errorf("abi: cannot map address to %s", dstRef.Type())
		}
		dstRef.SetBytes((*a)[:])
	case reflect.Array:
		if dstRef.Type().Elem().Kind() != reflect.Uint8 {
			return fmt.Errorf("abi: cannot map address to %s", dstRef.Type())
		}
		if dstRef.Len() != types.AddressLength {
			return fmt.Errorf("abi: cannot map address to %s: length mismatch", dstRef.Type())
		}
		for i := 0; i < dstRef.Len()-types.AddressLength; i++ {
			dstRef.Index(i).SetUint(0)
		}
		for i := 0; i < types.AddressLength; i++ {
			dstRef.Index(dstRef.Len() - types.AddressLength + i).SetUint(uint64((*a)[i]))
		}
	default:
		return fmt.Errorf("abi: cannot map address to %s", dstRef.Type())
	}
	return nil
}

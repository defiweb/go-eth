package abi

import (
	"bytes"
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/defiweb/go-eth/hexutil"
	"github.com/defiweb/go-eth/types"
)

func TestEncodeABI(t *testing.T) {
	tests := []struct {
		name    string
		val     Value
		abi     Words
		wantErr bool
	}{
		// TupleValue:
		{
			name: "tuple/empty",
			val:  NewTupleValue(),
			abi:  nil,
		},
		{
			name: "tuple/static",
			val: NewTupleValue(
				TupleValueElem{Value: NewBoolValue().Set(true)},
				TupleValueElem{Value: NewBoolValue().Set(true)},
			),
			abi: Words{padL("1"), padL("1")},
		},
		{
			name: "tuple/dynamic",
			val: NewTupleValue(
				TupleValueElem{Value: NewBytesValue().SetBytes([]byte{1, 2, 3})},
			),
			abi: Words{
				padL("20"),     // offset
				padL("3"),      // length
				padR("010203"), // data
			},
		},
		{
			name: "tuple/static-and-dynamic",
			val: NewTupleValue(
				TupleValueElem{Value: NewBoolValue().Set(true)},
				TupleValueElem{Value: NewBytesValue().SetBytes([]byte{1, 2, 3})},
			),
			abi: Words{
				padL("1"),      // arg0
				padL("40"),     // offset to arg1
				padL("3"),      // length of arg1
				padR("010203"), // arg1
			},
		},
		{
			name: "tuple/nested",
			val: NewTupleValue(
				TupleValueElem{
					Value: NewTupleValue(
						TupleValueElem{Value: NewBytesValue().SetBytes([]byte{1, 2, 3}), Name: "x"},
					),
					Name: "a",
				},
				TupleValueElem{
					Value: NewTupleValue(
						TupleValueElem{Value: NewBytesValue().SetBytes([]byte{4, 5, 6}), Name: "x"},
					),
					Name: "b",
				},
			),
			abi: Words{
				padL("40"),     // offset to a
				padL("a0"),     // offset to b
				padL("20"),     // offset to a.x
				padL("3"),      // length of a.x
				padR("010203"), // a.x
				padL("20"),     // offset to b.x
				padL("3"),      // length of b.x
				padR("040506"), // b.x
			},
		},
		// ArrayValue:
		{
			name: "array/empty",
			val:  NewArrayValue(NewBoolType()),
			abi:  Words{padL("0")},
		},
		{
			name: "array/two-static-elements",
			val:  NewArrayValue(NewBoolType(), NewBoolValue().Set(true), NewBoolValue().Set(true)),
			abi: Words{
				padL("2"), // array length
				padL("1"), // first element
				padL("1"), // second element
			},
		},
		{
			name: "array/two-dynamic-elements",
			val: NewArrayValue(
				NewBytesType(),
				NewBytesValue().SetBytes([]byte{1, 2, 3}),
				NewBytesValue().SetBytes([]byte{4, 5, 6}),
			),
			abi: Words{
				padL("2"),      // array length
				padL("40"),     // offset to first element
				padL("80"),     // offset to second element
				padL("3"),      // length of first element
				padR("010203"), // first element
				padL("3"),      // length of second element
				padR("040506"), // second element
			},
		},
		// FixedArrayValue:
		{
			name: "fixed-array/empty",
			val:  NewFixedArrayValue(NewBoolType(), 0),
			abi:  nil,
		},
		{
			name: "fixed-array/two-static-elements",
			val: NewFixedArrayValue(NewBoolType(), 2).
				SetElem(0, NewBoolValue().Set(true)).
				SetElem(1, NewBoolValue().Set(true)),
			abi: Words{
				padL("1"), // first element
				padL("1"), // second element
			},
		},
		{
			name: "fixed-array/two-dynamic-elements",
			val: NewFixedArrayValue(NewBytesType(), 2).
				SetElem(0, NewBytesValue().SetBytes([]byte{1, 2, 3})).
				SetElem(1, NewBytesValue().SetBytes([]byte{4, 5, 6})),
			abi: Words{
				padL("40"),     // offset to first element
				padL("80"),     // offset to second element
				padL("3"),      // length of first element
				padR("010203"), // first element
				padL("3"),      // length of second element
				padR("040506"), // second element
			},
		},
		// BytesValue:
		{
			name: "bytes/empty",
			val:  NewBytesValue(),
			abi:  Words{padL("0")},
		},
		{
			name: "bytes/one-word",
			val:  NewBytesValue().SetBytes([]byte{1, 2, 3}),
			abi: Words{
				padL("3"),      // length
				padR("010203"), // data
			},
		},
		{
			name: "bytes/two-words",
			val:  NewBytesValue().SetBytes(bytes.Repeat([]byte{0x01}, 33)),
			abi: Words{
				padL("21"), // length
				padR("0101010101010101010101010101010101010101010101010101010101010101"), // data
				padR("01"), // data
			},
		},
		// StringValue:
		{
			name: "string/empty",
			val:  NewStringValue(),
			abi:  Words{padL("0")},
		},
		{
			name: "string/one-word",
			val:  NewStringValue().SetString("abc"),
			abi: Words{
				padL("3"),      // length
				padR("616263"), // data
			},
		},
		{
			name: "string/two-words",
			val:  NewStringValue().SetString(strings.Repeat("a", 33)),
			abi: Words{
				padL("21"), // length
				padR("6161616161616161616161616161616161616161616161616161616161616161"), // data
				padR("61"), // data
			},
		},
		// FixedBytesValue:
		{
			name: "fixed/bytes-empty",
			val:  NewFixedBytesValue(1),
			abi:  Words{padL("0")},
		},
		{
			name: "fixed/bytes-non-empty",
			val:  NewFixedBytesValue(32).SetBytesPadRight([]byte{1, 2, 3}),
			abi: Words{
				padR("010203"), // data
			},
		},
		// UintValue:
		{
			name: "uint256/0",
			val:  NewUintValue(256),
			abi:  Words{padL("0")},
		},
		{
			name: "uint256/MaxUint256",
			val:  NewUintValue(256).SetBigInt(MaxUint256),
			abi:  Words{padR("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")},
		},
		// IntValue:
		{
			name: "int256/0",
			val:  NewIntValue(256).SetBigInt(big.NewInt(0)),
			abi:  Words{padL("0")},
		},
		{
			name: "int256/1",
			val:  NewIntValue(256).SetBigInt(big.NewInt(1)),
			abi:  Words{padL("1")},
		},
		{
			name: "int256/-1",
			val:  NewIntValue(256).SetBigInt(big.NewInt(-1)),
			abi:  Words{padR("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")},
		},
		{
			name: "int256/MaxInt256",
			val:  NewIntValue(256).SetBigInt(MaxInt256),
			abi:  Words{padR("7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")},
		},
		{
			name: "int256/MinInt256",
			val:  NewIntValue(256).SetBigInt(MinInt256),
			abi:  Words{padR("8000000000000000000000000000000000000000000000000000000000000000")},
		},
		{
			name: "int8/127",
			val:  NewIntValue(8).SetBigInt(big.NewInt(127)),
			abi:  Words{padL("7f")},
		},
		{
			name: "int8/-128",
			val:  NewIntValue(8).SetBigInt(big.NewInt(-128)),
			abi:  Words{padR("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80")},
		},
		// BoolValue:
		{
			name: "bool/true",
			val:  NewBoolValue().Set(true),
			abi:  Words{padL("1")},
		},
		{
			name: "bool/false",
			val:  NewBoolValue().Set(false),
			abi:  Words{padL("0")},
		},
		// AddressValue:
		{
			name: "address/empty",
			val:  NewAddressValue(),
			abi:  Words{padL("0")},
		},
		{
			name: "address/non-empty",
			val:  NewAddressValue().SetAddress(types.MustHexToAddress("0x0102030405060708090a0b0c0d0e0f1011121314")),
			abi:  Words{padL("0102030405060708090a0b0c0d0e0f1011121314")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.val.EncodeABI()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.abi, got)
			}
		})
	}
}

func TestDecodeABI(t *testing.T) {
	tests := []struct {
		name     string
		val      Value
		abi      Words
		assertFn func(t *testing.T, val Value)
		wantErr  bool
	}{
		{
			name: "tuple/empty",
			val:  NewTupleValue(),
			abi:  Words{},
			assertFn: func(t *testing.T, val Value) {
				assert.Equal(t, 0, val.(*TupleValue).Size())
			},
		},
		{
			name: "tuple/static",
			val: NewTupleValue(
				TupleValueElem{Value: NewBoolValue()},
				TupleValueElem{Value: NewBoolValue()},
			),
			abi: Words{
				padL("1"),
				padL("1"),
			},
			assertFn: func(t *testing.T, val Value) {
				assert.Equal(t, 2, val.(*TupleValue).Size())
				assert.True(t, val.(*TupleValue).Elem(0).Value.(*BoolValue).Bool())
				assert.True(t, val.(*TupleValue).Elem(1).Value.(*BoolValue).Bool())
			},
		},
		{
			name: "tuple/dynamic",
			val: NewTupleValue(
				TupleValueElem{Value: NewBytesValue()},
			),
			abi: Words{
				padL("20"),     // offset
				padL("3"),      // length
				padR("010203"), // data
			},
			assertFn: func(t *testing.T, val Value) {
				assert.Equal(t, 1, val.(*TupleValue).Size())
				assert.Equal(t, []byte{1, 2, 3}, val.(*TupleValue).Elem(0).Value.(*BytesValue).Bytes())
			},
		},
		{
			name: "tuple/static-and-dynamic",
			val: NewTupleValue(
				TupleValueElem{Value: NewBoolValue()},
				TupleValueElem{Value: NewBytesValue()},
			),
			abi: Words{
				padL("1"),      // arg0
				padL("40"),     // offset to arg1
				padL("3"),      // length of arg1
				padR("010203"), // arg1
			},
			assertFn: func(t *testing.T, val Value) {
				assert.Equal(t, 2, val.(*TupleValue).Size())
				assert.True(t, val.(*TupleValue).Elem(0).Value.(*BoolValue).Bool())
				assert.Equal(t, []byte{1, 2, 3}, val.(*TupleValue).Elem(1).Value.(*BytesValue).Bytes())
			},
		},
		{
			name: "tuple/nested",
			val: NewTupleValue(
				TupleValueElem{
					Value: NewTupleValue(
						TupleValueElem{Value: NewBytesValue(), Name: "x"},
					),
					Name: "a",
				},
				TupleValueElem{
					Value: NewTupleValue(
						TupleValueElem{Value: NewBytesValue(), Name: "x"},
					),
					Name: "b",
				},
			),
			abi: Words{
				padL("40"),     // offset to a
				padL("a0"),     // offset to b
				padL("20"),     // offset to a.x
				padL("3"),      // length of a.x
				padR("010203"), // a.x
				padL("20"),     // offset to b.x
				padL("3"),      // length of b.x
				padR("040506"), // b.x
			},
			assertFn: func(t *testing.T, val Value) {
				assert.Equal(t, 2, val.(*TupleValue).Size())
				assert.Equal(t, []byte{1, 2, 3}, val.(*TupleValue).Elem(0).Value.(*TupleValue).Elem(0).Value.(*BytesValue).Bytes())
				assert.Equal(t, []byte{4, 5, 6}, val.(*TupleValue).Elem(1).Value.(*TupleValue).Elem(0).Value.(*BytesValue).Bytes())
			},
		},
		// ArrayValue:
		{
			name:    "array/empty",
			val:     NewArrayValue(NewBoolType()),
			abi:     Words{},
			wantErr: true,
		},
		{
			name: "array/empty-2",
			val:  NewArrayValue(NewBoolType()),
			abi:  Words{padL("0")},
			assertFn: func(t *testing.T, val Value) {
				assert.Equal(t, 0, val.(*ArrayValue).Length())
			},
		},
		{
			name: "array/two-static-elements",
			val:  NewArrayValue(NewBoolType(), NewBoolValue(), NewBoolValue()),
			abi: Words{
				padL("2"), // array length
				padL("1"), // first element
				padL("1"), // second element
			},
			assertFn: func(t *testing.T, val Value) {
				assert.Equal(t, 2, val.(*ArrayValue).Length())
				assert.True(t, val.(*ArrayValue).Elem(0).(*BoolValue).Bool())
				assert.True(t, val.(*ArrayValue).Elem(1).(*BoolValue).Bool())
			},
		},
		{
			name: "array/two-dynamic-elements",
			val: NewArrayValue(
				NewBytesType(),
				NewBytesValue(),
				NewBytesValue(),
			),
			abi: Words{
				padL("2"),      // array length
				padL("40"),     // offset to first element
				padL("80"),     // offset to second element
				padL("3"),      // length of first element
				padR("010203"), // first element
				padL("3"),      // length of second element
				padR("040506"), // second element
			},
			assertFn: func(t *testing.T, val Value) {
				assert.Equal(t, 2, val.(*ArrayValue).Length())
				assert.Equal(t, []byte{1, 2, 3}, val.(*ArrayValue).Elem(0).(*BytesValue).Bytes())
				assert.Equal(t, []byte{4, 5, 6}, val.(*ArrayValue).Elem(1).(*BytesValue).Bytes())
			},
		},
		// FixedArrayValue:
		{
			name:    "fixed-array/empty",
			val:     NewFixedArrayValue(NewBoolType(), 0),
			abi:     Words{},
			wantErr: true,
		},
		{
			name: "fixed-array/empty-2",
			val:  NewFixedArrayValue(NewBoolType(), 0),
			abi:  Words{padL("0")},
			assertFn: func(t *testing.T, val Value) {
				assert.Equal(t, 0, val.(*FixedArrayValue).Size())
			},
		},
		{
			name: "fixed-array/two-static-elements",
			val: NewFixedArrayValue(NewBoolType(), 2).
				SetElem(0, NewBoolValue()).
				SetElem(1, NewBoolValue()),
			abi: Words{
				padL("1"), // first element
				padL("1"), // second element
			},
			assertFn: func(t *testing.T, val Value) {
				assert.Equal(t, 2, val.(*FixedArrayValue).Size())
				assert.True(t, val.(*FixedArrayValue).Elem(0).(*BoolValue).Bool())
				assert.True(t, val.(*FixedArrayValue).Elem(1).(*BoolValue).Bool())
			},
		},
		{
			name: "fixed-array/two-dynamic-elements",
			val: NewFixedArrayValue(NewBytesType(), 2).
				SetElem(0, NewBytesValue()).
				SetElem(1, NewBytesValue()),
			abi: Words{
				padL("40"),     // offset to first element
				padL("80"),     // offset to second element
				padL("3"),      // length of first element
				padR("010203"), // first element
				padL("3"),      // length of second element
				padR("040506"), // second element
			},
			assertFn: func(t *testing.T, val Value) {
				assert.Equal(t, 2, val.(*FixedArrayValue).Size())
				assert.Equal(t, []byte{1, 2, 3}, val.(*FixedArrayValue).Elem(0).(*BytesValue).Bytes())
				assert.Equal(t, []byte{4, 5, 6}, val.(*FixedArrayValue).Elem(1).(*BytesValue).Bytes())
			},
		},
		// BytesValue:
		{
			name:    "bytes/empty-1",
			val:     NewBytesValue(),
			abi:     Words{},
			wantErr: true,
		},
		{
			name: "bytes/empty-2",
			val:  NewBytesValue(),
			abi:  Words{padL("0")},
			assertFn: func(t *testing.T, val Value) {
				assert.Equal(t, []byte{}, val.(*BytesValue).Bytes())
			},
		},
		{
			name: "bytes/one-word",
			val:  NewBytesValue(),
			abi: Words{
				padL("3"),      // length
				padR("010203"), // data
			},
			assertFn: func(t *testing.T, val Value) {
				assert.Equal(t, []byte{1, 2, 3}, val.(*BytesValue).Bytes())
			},
		},
		{
			name: "bytes/two-words",
			val:  NewBytesValue(),
			abi: Words{
				padL("21"), // length
				padR("0101010101010101010101010101010101010101010101010101010101010101"), // data
				padR("01"), // data
			},
			assertFn: func(t *testing.T, val Value) {
				assert.Equal(t, bytes.Repeat([]byte{0x01}, 33), val.(*BytesValue).Bytes())
			},
		},
		// StringValue:
		{
			name:    "string/empty-1",
			val:     NewStringValue(),
			abi:     Words{},
			wantErr: true,
		},
		{
			name: "string/empty-2",
			val:  NewStringValue(),
			abi:  Words{padL("0")},
			assertFn: func(t *testing.T, val Value) {
				assert.Equal(t, []byte{}, val.(*StringValue).Bytes())
			},
		},
		{
			name: "string/one-word",
			val:  NewStringValue(),
			abi: Words{
				padL("3"),      // length
				padR("616263"), // data
			},
			assertFn: func(t *testing.T, val Value) {
				assert.Equal(t, "abc", val.(*StringValue).String())
			},
		},
		{
			name: "string/two-words",
			val:  NewStringValue(),
			abi: Words{
				padL("21"), // length
				padR("6161616161616161616161616161616161616161616161616161616161616161"), // data
				padR("61"), // data
			},
			assertFn: func(t *testing.T, val Value) {
				assert.Equal(t, strings.Repeat("a", 33), val.(*StringValue).String())
			},
		},
		// FixedBytesValue:
		{
			name:    "fixed/bytes-empty-1",
			val:     NewFixedBytesValue(1),
			abi:     Words{},
			wantErr: true,
		},
		{
			name: "fixed/bytes-empty-2",
			val:  NewFixedBytesValue(1),
			abi:  Words{padL("0")},
			assertFn: func(t *testing.T, val Value) {
				assert.Equal(t, []byte{0}, val.(*FixedBytesValue).Bytes())
			},
		},
		{
			name: "fixed/bytes-non-empty",
			val:  NewFixedBytesValue(4),
			abi: Words{
				padR("010203"), // data
			},
			assertFn: func(t *testing.T, val Value) {
				assert.Equal(t, []byte{1, 2, 3, 0}, val.(*FixedBytesValue).Bytes())
			},
		},
		// UintValue:
		{
			name: "uint256/0",
			val:  NewUintValue(64),
			abi:  Words{padL("0")},
			assertFn: func(t *testing.T, val Value) {
				assert.Equal(t, uint64(0), val.(*UintValue).Uint64())
			},
		},
		{
			name: "uint256/MaxUint256",
			val:  NewUintValue(256),
			abi:  Words{padR("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")},
			assertFn: func(t *testing.T, val Value) {
				assert.Equal(t, MaxUint256, val.(*UintValue).BigInt())
			},
		},
		// IntValue:
		{
			name: "int256/0",
			val:  NewIntValue(256),
			abi:  Words{padL("0")},
			assertFn: func(t *testing.T, val Value) {
				assert.Equal(t, big.NewInt(0), val.(*IntValue).BigInt())
			},
		},
		{
			name: "int256/1",
			val:  NewIntValue(256),
			abi:  Words{padL("1")},
			assertFn: func(t *testing.T, val Value) {
				assert.Equal(t, big.NewInt(1), val.(*IntValue).BigInt())
			},
		},
		{
			name: "int256/-1",
			val:  NewIntValue(256),
			abi:  Words{padR("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")},
			assertFn: func(t *testing.T, val Value) {
				assert.Equal(t, big.NewInt(-1), val.(*IntValue).BigInt())
			},
		},
		{
			name: "int256/MaxInt256",
			val:  NewIntValue(256),
			abi:  Words{padR("7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")},
			assertFn: func(t *testing.T, val Value) {
				assert.Equal(t, MaxInt256, val.(*IntValue).BigInt())
			},
		},
		{
			name: "int256/MinInt256",
			val:  NewIntValue(256),
			abi:  Words{padR("8000000000000000000000000000000000000000000000000000000000000000")},
			assertFn: func(t *testing.T, val Value) {
				assert.Equal(t, MinInt256, val.(*IntValue).BigInt())
			},
		},
		{
			name: "int8/127",
			val:  NewIntValue(8),
			abi:  Words{padL("7f")},
			assertFn: func(t *testing.T, val Value) {
				assert.Equal(t, big.NewInt(127), val.(*IntValue).BigInt())
			},
		},
		{
			name: "int8/-128",
			val:  NewIntValue(8),
			abi:  Words{padR("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80")},
			assertFn: func(t *testing.T, val Value) {
				assert.Equal(t, big.NewInt(-128), val.(*IntValue).BigInt())
			},
		},
		// BoolValue:
		{
			name: "bool/true",
			val:  NewBoolValue(),
			abi:  Words{padL("1")},
			assertFn: func(t *testing.T, val Value) {
				assert.True(t, val.(*BoolValue).Bool())
			},
		},
		{
			name: "bool/false",
			val:  NewBoolValue(),
			abi:  Words{padL("0")},
			assertFn: func(t *testing.T, val Value) {
				assert.False(t, val.(*BoolValue).Bool())
			},
		},
		// AddressValue:
		{
			name: "address/empty",
			val:  NewAddressValue(),
			abi:  Words{padL("0")},
			assertFn: func(t *testing.T, val Value) {
				assert.Equal(t, types.Address{}, val.(*AddressValue).Address())
			},
		},
		{
			name: "address/non-empty",
			val:  NewAddressValue(),
			abi:  Words{padL("0102030405060708090a0b0c0d0e0f1011121314")},
			assertFn: func(t *testing.T, val Value) {
				assert.Equal(t, types.MustHexToAddress("0x0102030405060708090a0b0c0d0e0f1011121314"), val.(*AddressValue).Address())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.val.DecodeABI(tt.abi)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				tt.assertFn(t, tt.val)
			}
		})
	}
}

func padL(h string) (w Word) {
	_ = (&w).SetBytesPadLeft(hexutil.MustHexToBytes(h))
	return w
}

func padR(h string) (w Word) {
	_ = (&w).SetBytesPadRight(hexutil.MustHexToBytes(h))
	return w
}

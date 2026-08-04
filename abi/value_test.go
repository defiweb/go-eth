package abi

import (
	"bytes"
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/defiweb/go-eth/hexutil"
	"github.com/defiweb/go-eth/types"
)

func TestEncodeABI(t *testing.T) {
	tests := []struct {
		name    string
		val     Value
		arg     any
		want    Words
		wantErr bool
	}{
		// TupleValue:
		{
			name: "tuple#empty",
			val:  &TupleValue{},
			arg:  struct{}{},
			want: Words(nil),
		},
		{
			name: "tuple#static",
			val: &TupleValue{
				TupleValueElem{Value: new(BoolValue), Name: "a"},
				TupleValueElem{Value: new(BoolValue), Name: "b"},
			},
			arg:  map[string]bool{"a": true, "b": true},
			want: Words{padL("01"), padL("01")},
		},
		{
			name: "tuple#dynamic",
			val: &TupleValue{
				TupleValueElem{Value: new(BytesValue), Name: "a"},
			},
			arg: map[string][]byte{"a": {1, 2, 3}},
			want: Words{
				padL("20"),     // offset
				padL("03"),     // length
				padR("010203"), // data
			},
		},
		{
			name: "tuple#static-and-dynamic",
			val: &TupleValue{
				TupleValueElem{Value: new(BoolValue), Name: "a"},
				TupleValueElem{Value: new(BytesValue), Name: "b"},
			},
			arg: map[string]interface{}{"a": true, "b": []byte{1, 2, 3}},
			want: Words{
				padL("01"),     // a
				padL("40"),     // offset to b
				padL("03"),     // length of b
				padR("010203"), // b
			},
		},
		{
			name: "tuple#nested",
			val: &TupleValue{
				TupleValueElem{
					Value: &TupleValue{
						TupleValueElem{Value: new(BytesValue), Name: "x"},
					},
					Name: "a",
				},
				TupleValueElem{
					Value: &TupleValue{
						TupleValueElem{Value: new(BytesValue), Name: "x"},
					},
					Name: "b",
				},
			},
			arg: map[string]interface{}{
				"a": map[string][]byte{"x": {1, 2, 3}},
				"b": map[string][]byte{"x": {4, 5, 6}},
			},
			want: Words{
				padL("40"),     // offset to a
				padL("a0"),     // offset to b
				padL("20"),     // offset to a.x
				padL("03"),     // length of a.x
				padR("010203"), // a.x
				padL("20"),     // offset to b.x
				padL("03"),     // length of b.x
				padR("040506"), // b.x
			},
		},
		// ArrayValue:
		{
			name: "array#empty",
			val:  &ArrayValue{Type: NewBoolType()},
			arg:  []bool{},
			want: Words{padL("0")},
		},
		{
			name: "array#two-static-elements",
			val:  &ArrayValue{Type: NewBoolType()},
			arg:  []bool{true, true},
			want: Words{
				padL("02"), // array length
				padL("01"), // first element
				padL("01"), // second element
			},
		},
		{
			name: "array#two-dynamic-elements",
			val:  &ArrayValue{Type: NewBytesType()},
			arg:  [][]byte{{1, 2, 3}, {4, 5, 6}},
			want: Words{
				padL("02"),     // array length
				padL("40"),     // offset to first element
				padL("80"),     // offset to second element
				padL("03"),     // length of first element
				padR("010203"), // first element
				padL("03"),     // length of second element
				padR("040506"), // second element
			},
		},
		// FixedArrayValue:
		{
			name: "fixed-array#empty",
			val:  &FixedArrayValue{},
			arg:  [0]bool{},
			want: nil,
		},
		{
			name: "fixed-array#two-static-elements",
			val:  &FixedArrayValue{new(BoolValue), new(BoolValue)},
			arg:  []bool{true, true},
			want: Words{
				padL("01"), // first element
				padL("01"), // second element
			},
		},
		{
			name: "fixed-array#two-dynamic-elements",
			val:  &FixedArrayValue{new(BytesValue), new(BytesValue)},
			arg:  [][]byte{{1, 2, 3}, {4, 5, 6}},
			want: Words{
				padL("40"),     // offset to first element
				padL("80"),     // offset to second element
				padL("03"),     // length of first element
				padR("010203"), // first element
				padL("03"),     // length of second element
				padR("040506"), // second element
			},
		},
		// BytesValue:
		{
			name: "bytes#empty",
			val:  new(BytesValue),
			arg:  []byte{},
			want: Words{padL("00")},
		},
		{
			name: "bytes#one-word",
			val:  new(BytesValue),
			arg:  []byte{1, 2, 3},
			want: Words{
				padL("03"),     // length
				padR("010203"), // data
			},
		},
		{
			name: "bytes#two-words",
			val:  new(BytesValue),
			arg:  bytes.Repeat([]byte{1}, 33),
			want: Words{
				padL("21"), // length
				padR("0101010101010101010101010101010101010101010101010101010101010101"), // data
				padR("01"), // data
			},
		},

		// StringValue:
		{
			name: "string#empty",
			val:  new(StringValue),
			arg:  "",
			want: Words{padL("0")},
		},
		{
			name: "string#one-word",
			val:  new(StringValue),
			arg:  "abc",
			want: Words{
				padL("03"),     // length
				padR("616263"), // data
			},
		},
		{
			name: "string#two-words",
			val:  new(StringValue),
			arg:  strings.Repeat("a", 33),
			want: Words{
				padL("21"), // length
				padR("6161616161616161616161616161616161616161616161616161616161616161"), // data
				padR("61"), // data
			},
		},
		// FixedBytesValue:
		{
			name: "fixed#bytes-empty",
			val:  make(FixedBytesValue, 0),
			arg:  []byte{},
			want: Words{padL("00")},
		},
		{
			name: "fixed#bytes-non-empty",
			val:  make(FixedBytesValue, 3),
			arg:  []byte{1, 2, 3},
			want: Words{padR("010203")},
		},
		// UintValue:
		{
			name: "uint8#0",
			val:  &UintValue{Size: 8},
			arg:  big.NewInt(0),
			want: Words{padL("00")},
		},
		{
			name: "uint256#0",
			val:  &UintValue{Size: 256},
			arg:  big.NewInt(0),
			want: Words{padL("00")},
		},
		{
			name: "uint256#MaxUint256",
			val:  &UintValue{Size: 256},
			arg:  new(big.Int).Set(MaxUint[256]),
			want: Words{padR("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")},
		},
		// IntValue:
		{
			name: "int8#0",
			val:  &IntValue{Size: 8},
			arg:  big.NewInt(0),
			want: Words{padL("00")},
		},
		{
			name: "int256#0",
			val:  &IntValue{Size: 256},
			arg:  big.NewInt(0),
			want: Words{padL("00")},
		},
		{
			name: "int256#1",
			val:  &IntValue{Size: 256},
			arg:  big.NewInt(1),
			want: Words{padL("01")},
		},
		{
			name: "int256#-1",
			val:  &IntValue{Size: 256},
			arg:  big.NewInt(-1),
			want: Words{padR("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")},
		},
		{
			name: "int256#MaxInt256",
			val:  &IntValue{Size: 256},
			arg:  new(big.Int).Set(MaxInt[256]),
			want: Words{padR("7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")},
		},
		{
			name: "int256#MinInt256",
			val:  &IntValue{Size: 256},
			arg:  new(big.Int).Set(MinInt[256]),
			want: Words{padR("8000000000000000000000000000000000000000000000000000000000000000")},
		},
		{
			name: "int8#127",
			val:  &IntValue{Size: 8},
			arg:  big.NewInt(127),
			want: Words{padL("7f")},
		},
		{
			name: "int8#-128",
			val:  &IntValue{Size: 8},
			arg:  big.NewInt(-128),
			want: Words{padR("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80")},
		},
		// BoolValue:
		{
			name: "bool#true",
			val:  new(BoolValue),
			arg:  true,
			want: Words{padL("01")},
		},
		{
			name: "bool#false",
			val:  new(BoolValue),
			arg:  false,
			want: Words{padL("00")},
		},
		// AddressValue:
		{
			name: "address#empty",
			val:  new(AddressValue),
			arg:  types.Address{},
			want: Words{padL("00")},
		},
		{
			name: "address#non-empty",
			val:  new(AddressValue),
			arg:  types.MustAddressFromHex("0x0102030405060708090a0b0c0d0e0f1011121314"),
			want: Words{padL("0102030405060708090a0b0c0d0e0f1011121314")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, Default.Mapper.Map(tt.arg, tt.val))
			enc, err := tt.val.EncodeABI()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, enc)
			}
		})
	}
}

func TestDecodeABI(t *testing.T) {
	tests := []struct {
		name    string
		abi     Words
		val     Value
		want    Value
		wantErr bool
	}{
		// TupleValue:
		{
			name: "tuple#empty",
			abi:  Words{},
			val:  new(TupleValue),
			want: new(TupleValue),
		},
		{
			name: "tuple#static",
			abi: Words{
				padL("01"),
				padL("01"),
			},
			val: &TupleValue{
				TupleValueElem{Value: new(BoolValue), Name: "a"},
				TupleValueElem{Value: new(BoolValue), Name: "b"},
			},
			want: &TupleValue{
				TupleValueElem{Value: func() *BoolValue { b := BoolValue(true); return &b }(), Name: "a"},
				TupleValueElem{Value: func() *BoolValue { b := BoolValue(true); return &b }(), Name: "b"},
			},
		},
		{
			name: "tuple#dynamic",
			abi: Words{
				padL("20"),     // offset
				padL("03"),     // length
				padR("010203"), // data
			},
			val: &TupleValue{
				TupleValueElem{Value: new(BytesValue), Name: "a"},
			},
			want: &TupleValue{
				TupleValueElem{Value: func() *BytesValue { b := BytesValue([]byte{1, 2, 3}); return &b }(), Name: "a"},
			},
		},
		{
			name: "tuple#static-and-dynamic",
			abi: Words{
				padL("01"),     // a
				padL("40"),     // offset to b
				padL("03"),     // length of b
				padR("010203"), // b
			},
			val: &TupleValue{
				TupleValueElem{Value: new(BoolValue), Name: "a"},
				TupleValueElem{Value: new(BytesValue), Name: "b"},
			},
			want: &TupleValue{
				TupleValueElem{Value: func() *BoolValue { b := BoolValue(true); return &b }(), Name: "a"},
				TupleValueElem{Value: func() *BytesValue { b := BytesValue([]byte{1, 2, 3}); return &b }(), Name: "b"},
			},
		},
		{
			name: "tuple#nested",
			abi: Words{
				padL("40"),     // offset to a
				padL("a0"),     // offset to b
				padL("20"),     // offset to a.x
				padL("03"),     // length of a.x
				padR("010203"), // a.x
				padL("20"),     // offset to b.x
				padL("03"),     // length of b.x
				padR("040506"), // b.x
			},
			val: &TupleValue{
				TupleValueElem{
					Value: &TupleValue{
						TupleValueElem{Value: new(BytesValue), Name: "x"},
					},
					Name: "a",
				},
				TupleValueElem{
					Value: &TupleValue{
						TupleValueElem{Value: new(BytesValue), Name: "x"},
					},
					Name: "b",
				},
			},
			want: &TupleValue{
				TupleValueElem{
					Value: &TupleValue{
						TupleValueElem{Value: func() *BytesValue { b := BytesValue([]byte{1, 2, 3}); return &b }(), Name: "x"},
					},
					Name: "a",
				},
				TupleValueElem{
					Value: &TupleValue{
						TupleValueElem{Value: func() *BytesValue { b := BytesValue([]byte{4, 5, 6}); return &b }(), Name: "x"},
					},
					Name: "b",
				},
			},
		},
		// ArrayValue:
		{
			name: "array#empty",
			abi:  Words{padL("00")},
			val:  &ArrayValue{Type: NewBoolType()},
			want: &ArrayValue{Type: NewBoolType(), Elems: []Value{}},
		},
		{
			name: "array#two-static-elements",
			abi: Words{
				padL("02"), // array length
				padL("01"), // first element
				padL("01"), // second element
			},
			val: &ArrayValue{Type: NewBoolType()},
			want: &ArrayValue{
				Type: NewBoolType(),
				Elems: []Value{
					func() *BoolValue { b := BoolValue(true); return &b }(),
					func() *BoolValue { b := BoolValue(true); return &b }(),
				},
			},
		},
		{
			name: "array#two-dynamic-elements",
			abi: Words{
				padL("02"),     // array length
				padL("40"),     // offset to first element
				padL("80"),     // offset to second element
				padL("03"),     // length of first element
				padR("010203"), // first element
				padL("03"),     // length of second element
				padR("040506"), // second element
			},
			val: &ArrayValue{Type: NewBytesType()},
			want: &ArrayValue{
				Type: NewBytesType(),
				Elems: []Value{
					func() *BytesValue { b := BytesValue([]byte{1, 2, 3}); return &b }(),
					func() *BytesValue { b := BytesValue([]byte{4, 5, 6}); return &b }(),
				},
			},
		},
		// FixedArrayValue:
		{
			name: "fixed-array#empty",
			val:  &FixedArrayValue{},
			abi:  Words{padL("0")},
			want: &FixedArrayValue{},
		},
		{
			name: "fixed-array#two-static-elements",
			abi: Words{
				padL("01"), // first element
				padL("01"), // second element
			},
			val: &FixedArrayValue{new(BoolValue), new(BoolValue)},
			want: &FixedArrayValue{
				func() *BoolValue { b := BoolValue(true); return &b }(),
				func() *BoolValue { b := BoolValue(true); return &b }(),
			},
		},
		{
			name: "fixed-array#two-dynamic-elements",
			abi: Words{
				padL("40"),     // offset to first element
				padL("80"),     // offset to second element
				padL("03"),     // length of first element
				padR("010203"), // first element
				padL("03"),     // length of second element
				padR("040506"), // second element
			},
			val: &FixedArrayValue{new(BytesValue), new(BytesValue)},
			want: &FixedArrayValue{
				func() *BytesValue { b := BytesValue([]byte{1, 2, 3}); return &b }(),
				func() *BytesValue { b := BytesValue([]byte{4, 5, 6}); return &b }(),
			},
		},
		// BytesValue:
		{
			name: "bytes#empty",
			abi:  Words{padL("00")},
			val:  new(BytesValue),
			want: func() *BytesValue { b := BytesValue([]byte{}); return &b }(),
		},
		{
			name: "bytes#one-word",
			abi: Words{
				padL("03"),     // length
				padR("010203"), // data
			},
			val:  new(BytesValue),
			want: func() *BytesValue { b := BytesValue([]byte{1, 2, 3}); return &b }(),
		},
		{
			name: "bytes#two-words",
			abi: Words{
				padL("21"), // length
				padR("0101010101010101010101010101010101010101010101010101010101010101"), // data
				padR("01"), // data
			},
			val:  new(BytesValue),
			want: func() *BytesValue { b := BytesValue(bytes.Repeat([]byte{1}, 33)); return &b }(),
		},
		// StringValue:
		{
			name: "string#empty",
			abi:  Words{padL("0")},
			val:  new(StringValue),
			want: func() *StringValue { s := StringValue(""); return &s }(),
		},
		{
			name: "string#one-word",
			abi: Words{
				padL("03"),     // length
				padR("616263"), // data
			},
			val:  new(StringValue),
			want: func() *StringValue { s := StringValue("abc"); return &s }(),
		},
		{
			name: "string#two-words",
			abi: Words{
				padL("21"), // length
				padR("6161616161616161616161616161616161616161616161616161616161616161"), // data
				padR("61"), // data
			},
			val:  new(StringValue),
			want: func() *StringValue { s := StringValue(strings.Repeat("a", 33)); return &s }(),
		},
		// FixedBytesValue:
		{
			name: "fixed#bytes-empty",
			abi:  Words{padL("0")},
			val:  make(FixedBytesValue, 0),
			want: make(FixedBytesValue, 0),
		},
		{
			name: "fixed#bytes-non-empty",
			abi: Words{
				padR("010203"),
			},
			val:  make(FixedBytesValue, 3),
			want: func() FixedBytesValue { b := FixedBytesValue([]byte{1, 2, 3}); return b }(),
		},
		// UintValue:
		{
			name: "uint8#0",
			abi:  Words{padL("00")},
			val:  &UintValue{Size: 8},
			want: func() *UintValue { u := UintValue{Size: 8, Int: *big.NewInt(0)}; return &u }(),
		},
		{
			name: "uint8#1",
			abi:  Words{padL("01")},
			val:  &UintValue{Size: 8},
			want: func() *UintValue { u := UintValue{Size: 8, Int: *big.NewInt(1)}; return &u }(),
		},
		{
			name: "uint256#0",
			abi:  Words{padL("00")},
			val:  &UintValue{Size: 256},
			want: func() *UintValue { u := UintValue{Size: 256, Int: *big.NewInt(0)}; return &u }(),
		},
		{
			name: "uint256#MaxUint256",
			abi:  Words{padR("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")},
			val:  &UintValue{Size: 256},
			want: func() *UintValue { u := UintValue{Size: 256, Int: *MaxUint[256]}; return &u }(),
		},
		// IntValue:
		{
			name: "int8#0",
			abi:  Words{padL("00")},
			val:  &IntValue{Size: 8},
			want: func() *IntValue { i := IntValue{Size: 8, Int: *big.NewInt(0)}; return &i }(),
		},
		{
			name: "int8#1",
			abi:  Words{padL("01")},
			val:  &IntValue{Size: 8},
			want: func() *IntValue { i := IntValue{Size: 8, Int: *big.NewInt(1)}; return &i }(),
		},
		{
			name: "int8#-1",
			abi:  Words{padL("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")},
			val:  &IntValue{Size: 8},
			want: func() *IntValue { i := IntValue{Size: 8, Int: *big.NewInt(-1)}; return &i }(),
		},
		{
			name: "int256#0",
			abi:  Words{padL("00")},
			val:  &IntValue{Size: 256},
			want: func() *IntValue { i := IntValue{Size: 256, Int: *big.NewInt(0)}; return &i }(),
		},
		{
			name: "int256#1",
			abi:  Words{padL("01")},
			val:  &IntValue{Size: 256},
			want: func() *IntValue { i := IntValue{Size: 256, Int: *big.NewInt(1)}; return &i }(),
		},
		{
			name: "int256#-1",
			abi:  Words{padR("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")},
			val:  &IntValue{Size: 256},
			want: func() *IntValue { i := IntValue{Size: 256, Int: *big.NewInt(-1)}; return &i }(),
		},
		{
			name: "int256#MaxInt256",
			abi:  Words{padR("7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")},
			val:  &IntValue{Size: 256},
			want: func() *IntValue { i := IntValue{Size: 256, Int: *MaxInt[256]}; return &i }(),
		},
		{
			name: "int256#MinInt256",
			abi:  Words{padR("8000000000000000000000000000000000000000000000000000000000000000")},
			val:  &IntValue{Size: 256},
			want: func() *IntValue { i := IntValue{Size: 256, Int: *MinInt[256]}; return &i }(),
		},
		{
			name: "int8#127",
			abi:  Words{padL("7f")},
			val:  &IntValue{Size: 256},
			want: func() *IntValue { i := IntValue{Size: 256, Int: *big.NewInt(127)}; return &i }(),
		},
		{
			name: "int8#-128",
			abi:  Words{padR("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80")},
			val:  &IntValue{Size: 256},
			want: func() *IntValue { i := IntValue{Size: 256, Int: *big.NewInt(-128)}; return &i }(),
		},
		// BoolValue:
		{
			name: "bool#true",
			abi:  Words{padL("01")},
			val:  new(BoolValue),
			want: func() *BoolValue { b := BoolValue(true); return &b }(),
		},
		{
			name: "bool#false",
			abi:  Words{padL("00")},
			val:  new(BoolValue),
			want: func() *BoolValue { b := BoolValue(false); return &b }(),
		},
		// AddressValue:
		{
			name: "address#empty",
			abi:  Words{padL("00")},
			val:  new(AddressValue),
			want: func() *AddressValue { a := AddressValue{}; return &a }(),
		},
		{
			name: "address#non-empty",
			abi:  Words{padL("0102030405060708090a0b0c0d0e0f1011121314")},
			val:  new(AddressValue),
			want: func() *AddressValue {
				a := types.MustAddressFromHex("0102030405060708090a0b0c0d0e0f1011121314")
				return (*AddressValue)(&a)
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.val.DecodeABI(tt.abi)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, tt.val)
			}
		})
	}
}

func TestMapFrom(t *testing.T) {
	tests := []struct {
		name    string
		val     Value
		data    any
		want    Words
		wantErr bool
	}{
		// BytesValue:
		{
			name: "[]byte->bytes",
			val:  new(BytesValue),
			data: []byte{1, 2, 3},
			want: Words{
				padL("03"),     // length
				padR("010203"), // data
			},
		},
		{
			name: "array->bytes",
			val:  new(BytesValue),
			data: [3]byte{1, 2, 3},
			want: Words{
				padL("03"),     // length
				padR("010203"), // data
			},
		},
		{
			name: "string->bytes",
			val:  new(BytesValue),
			data: "0x2a",
			want: Words{
				padL("01"), // length
				padR("2a"), // data
			},
		},
		{
			name:    "int64->bytes",
			val:     new(BytesValue),
			data:    int64(42),
			wantErr: true,
		},
		{
			name:    "big.Int->bytes",
			val:     new(BytesValue),
			data:    big.NewInt(42),
			wantErr: true,
		},
		// StringValue:
		{
			name: "string->string",
			val:  new(StringValue),
			data: "foo",
			want: Words{
				padL("03"),     // length
				padR("666f6f"), // data
			},
		},
		{
			name: "[]byte->string",
			val:  new(StringValue),
			data: []byte{1, 2, 3},
			want: Words{
				padL("03"),     // length
				padR("010203"), // data
			},
		},
		{
			name:    "array->string",
			val:     new(StringValue),
			data:    [3]byte{1, 2, 3},
			wantErr: true,
		},
		{
			name:    "int64->string",
			val:     new(StringValue),
			data:    int64(42),
			wantErr: true,
		},
		{
			name:    "big.Int->string",
			val:     new(StringValue),
			data:    big.NewInt(42),
			wantErr: true,
		},
		{
			name:    "types.Address->string",
			val:     new(StringValue),
			data:    types.MustAddressFromHex("0102030405060708090a0b0c0d0e0f1011121314"),
			wantErr: true,
		},
		{
			name:    "types.Hash->string",
			val:     new(StringValue),
			data:    types.MustHashFromHex("0102030405060708090a0b0c0d0e0f10111213145566778899aabbccddeeff00", types.PadNone),
			wantErr: true,
		},
		{
			name:    "types.Number->string",
			val:     new(StringValue),
			data:    types.MustNumberFromHex("2a"),
			wantErr: true,
		},
		// FixedBytesValue
		{
			name: "[]byte->bytes4",
			val:  make(FixedBytesValue, 4),
			data: []byte{1, 2, 3, 4},
			want: Words{
				padR("01020304"),
			},
		},
		{
			name:    "[]byte->bytes32#InvalidLength",
			val:     make(FixedBytesValue, 32),
			data:    []byte{1, 2, 3, 4},
			wantErr: true,
		},
		{
			name: "array->bytes4",
			val:  make(FixedBytesValue, 4),
			data: [4]byte{1, 2, 3, 4},
			want: Words{
				padR("01020304"),
			},
		},
		{
			name:    "array->bytes32#InvalidLength",
			val:     make(FixedBytesValue, 32),
			data:    [4]byte{1, 2, 3, 4},
			wantErr: true,
		},
		{
			name: "string->bytes4",
			val:  make(FixedBytesValue, 4),
			data: "0x01020304",
			want: Words{
				padR("01020304"),
			},
		},
		{
			name:    "string->bytes32#InvalidLength",
			val:     make(FixedBytesValue, 32),
			data:    "0x01020304",
			wantErr: true,
		},
		{
			name: "uint64->bytes32",
			val:  make(FixedBytesValue, 32),
			data: uint64(42),
			want: Words{
				padL("2a"),
			},
		},
		{
			name:    "uint64->bytes16",
			val:     make(FixedBytesValue, 16),
			data:    uint64(42),
			wantErr: true,
		},
		{
			name: "int64->bytes32",
			val:  make(FixedBytesValue, 32),
			data: int64(42),
			want: Words{
				padL("2a"),
			},
		},
		{
			name:    "int64->bytes32/negative",
			val:     make(FixedBytesValue, 32),
			data:    int64(-42),
			wantErr: true,
		},
		{
			name:    "int64->bytes16",
			val:     make(FixedBytesValue, 16),
			data:    uint64(42),
			wantErr: true,
		},
		{
			name: "big.Int->bytes32",
			val:  make(FixedBytesValue, 32),
			data: big.NewInt(42),
			want: Words{
				padL("2a"),
			},
		},
		{
			name:    "big.Int->bytes16",
			val:     make(FixedBytesValue, 16),
			data:    big.NewInt(42),
			wantErr: true,
		},
		{
			name:    "big.Int->bytes32/negative",
			val:     make(FixedBytesValue, 32),
			data:    big.NewInt(-42),
			wantErr: true,
		},
		{
			name:    "big.Int->bytes16",
			val:     make(FixedBytesValue, 16),
			data:    big.NewInt(42),
			wantErr: true,
		},
		{
			name:    "types.Address->bytes32",
			val:     make(FixedBytesValue, 32),
			data:    types.MustAddressFromHex("0102030405060708090a0b0c0d0e0f1011121314"),
			wantErr: true,
		},
		{
			name: "types.Address->bytes20",
			val:  make(FixedBytesValue, 20),
			data: types.MustAddressFromHex("0102030405060708090a0b0c0d0e0f1011121314"),
			want: Words{
				padR("0102030405060708090a0b0c0d0e0f1011121314"),
			},
		},
		{
			name: "types.Hash->bytes32",
			val:  make(FixedBytesValue, 32),
			data: types.MustHashFromHex("0102030405060708090a0b0c0d0e0f10111213145566778899aabbccddeeff00", types.PadNone),
			want: Words{
				padR("0102030405060708090a0b0c0d0e0f10111213145566778899aabbccddeeff00"),
			},
		},
		{
			name: "types.Number->bytes32",
			val:  make(FixedBytesValue, 32),
			data: types.MustNumberFromHex("2a"),
			want: Words{
				padL("2a"),
			},
		},
		{
			name:    "types.Number->bytes16",
			val:     make(FixedBytesValue, 16),
			data:    types.MustNumberFromHex("2a"),
			wantErr: true,
		},
		// UintValue:
		{
			name: "string->uint256",
			val:  &UintValue{Size: 256},
			data: "0x2a",
			want: Words{
				padL("2a"),
			},
		},
		{
			name: "string->uint256#odd",
			val:  &UintValue{Size: 256},
			data: "0x2",
			want: Words{
				padL("02"),
			},
		},
		{
			name:    "string->uint256#negative",
			val:     &UintValue{Size: 256},
			data:    "-0x2a",
			wantErr: true,
		},
		{
			name:    "[]byte->uint256",
			val:     &UintValue{Size: 256},
			data:    []byte{1},
			wantErr: true,
		},
		{
			name:    "array->uint256",
			val:     &UintValue{Size: 256},
			data:    [1]byte{1},
			wantErr: true,
		},
		{
			name: "int64->uint256",
			val:  &UintValue{Size: 256},
			data: int64(42),
			want: Words{
				padL("2a"),
			},
		},
		{
			name:    "int64->uint256#negative",
			val:     &UintValue{Size: 256},
			data:    int64(-42),
			wantErr: true,
		},
		{
			name: "uint16->uint8#255",
			val:  &UintValue{Size: 8},
			data: uint16(255),
			want: Words{
				padL("ff"),
			},
		},
		{
			name:    "uint16->uint8#256",
			val:     &UintValue{Size: 8},
			data:    uint16(256),
			wantErr: true,
		},
		{
			name: "uint64->uint256",
			val:  &UintValue{Size: 256},
			data: uint64(42),
			want: Words{
				padL("2a"),
			},
		},
		{
			name: "uint16->uint8#255",
			val:  &UintValue{Size: 8},
			data: uint16(255),
			want: Words{
				padL("ff"),
			},
		},
		{
			name:    "uint16->uint8#256",
			val:     &UintValue{Size: 8},
			data:    uint16(256),
			wantErr: true,
		},
		{
			name: "big.Int->uint256",
			val:  &UintValue{Size: 256},
			data: big.NewInt(42),
			want: Words{
				padL("2a"),
			},
		},
		{
			name:    "big.Int->uint256#negative",
			val:     &UintValue{Size: 256},
			data:    big.NewInt(-42),
			wantErr: true,
		},
		{
			name:    "big.Int->uint8#overflow",
			val:     &UintValue{Size: 8},
			data:    big.NewInt(256),
			wantErr: true,
		},
		{
			name:    "types.Address->uint256",
			val:     &UintValue{Size: 256},
			data:    types.MustAddressFromHex("0102030405060708090a0b0c0d0e0f1011121314"),
			wantErr: true,
		},
		{
			name:    "types.Hash->uint256",
			val:     &UintValue{Size: 256},
			data:    types.MustHashFromHex("0102030405060708090a0b0c0d0e0f10111213145566778899aabbccddeeff00", types.PadNone),
			wantErr: true,
		},
		{
			name: "types.Number->uint256",
			val:  &UintValue{Size: 256},
			data: types.MustNumberFromHex("2a"),
			want: Words{
				padL("2a"),
			},
		},
		// IntValue:
		{
			name: "string->int",
			val:  &IntValue{Size: 256},
			data: "0x2a",
			want: Words{
				padL("2a"),
			},
		},
		{
			name: "string->int#odd",
			val:  &IntValue{Size: 256},
			data: "0x2",
			want: Words{
				padL("02"),
			},
		},
		{
			name: "string->int#negative",
			val:  &IntValue{Size: 256},
			data: "-0x2a",
			want: Words{
				padL("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd6"),
			},
		},
		{
			name:    "[]byte->int",
			val:     &IntValue{Size: 256},
			data:    []byte{1},
			wantErr: true,
		},
		{
			name:    "array->int",
			val:     &IntValue{Size: 256},
			data:    [1]byte{1},
			wantErr: true,
		},
		{
			name: "int64->int256",
			val:  &IntValue{Size: 256},
			data: int64(42),
			want: Words{
				padL("2a"),
			},
		},
		{
			name: "int64->int256#negative",
			val:  &IntValue{Size: 256},
			data: int64(-42),
			want: Words{
				padL("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd6"),
			},
		},
		{
			name: "int16->int8#127",
			val:  &IntValue{Size: 8},
			data: int16(127),
			want: Words{
				padL("7f"),
			},
		},
		{
			name:    "int16->int8#128",
			val:     &IntValue{Size: 8},
			data:    int16(128),
			wantErr: true,
		},
		{
			name: "uint64->int256",
			val:  &IntValue{Size: 256},
			data: uint64(42),
			want: Words{
				padL("2a"),
			},
		},
		{
			name: "uint16->int8#128",
			val:  &IntValue{Size: 8},
			data: uint16(127),
			want: Words{
				padL("7f"),
			},
		},
		{
			name:    "uint16->int8#128",
			val:     &IntValue{Size: 8},
			data:    uint16(128),
			wantErr: true,
		},
		{
			name: "big.Int->int256",
			val:  &IntValue{Size: 256},
			data: big.NewInt(42),
			want: Words{
				padL("2a"),
			},
		},
		{
			name: "big.Int->int256#negative",
			val:  &IntValue{Size: 256},
			data: big.NewInt(-42),
			want: Words{
				padL("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd6"),
			},
		},
		{
			name:    "big.Int->int8#overflow",
			val:     &IntValue{Size: 8},
			data:    big.NewInt(256),
			wantErr: true,
		},
		{
			name:    "types.Address->int256",
			val:     &IntValue{Size: 256},
			data:    types.MustAddressFromHex("0102030405060708090a0b0c0d0e0f1011121314"),
			wantErr: true,
		},
		{
			name:    "types.Hash->int256",
			val:     &IntValue{Size: 256},
			data:    types.MustHashFromHex("0102030405060708090a0b0c0d0e0f10111213145566778899aabbccddeeff00", types.PadNone),
			wantErr: true,
		},
		{
			name: "types.Number->int256",
			val:  &IntValue{Size: 256},
			data: types.MustNumberFromHex("2a"),
			want: Words{
				padL("2a"),
			},
		},
		// BoolValue:
		{
			name: "bool->bool",
			val:  new(BoolValue),
			data: true,
			want: Words{
				padL("01"),
			},
		},
		{
			name:    "string->bool",
			val:     new(BoolValue),
			data:    "true",
			wantErr: true,
		},
		{
			name:    "int->bool",
			val:     new(BoolValue),
			data:    1,
			wantErr: true,
		},
		// AddressValue:
		{
			name: "string->address",
			val:  new(AddressValue),
			data: "0x1234567890123456789012345678901234567890",
			want: Words{
				padL("1234567890123456789012345678901234567890"), // data
			},
		},
		{
			name:    "string->address/short",
			val:     new(AddressValue),
			data:    "0x12",
			wantErr: true,
		},
		{
			name:    "string->address/long",
			val:     new(AddressValue),
			data:    "0x123456789012345678901234567890123456789012",
			wantErr: true,
		},
		{
			name: "[]byte->address",
			val:  new(AddressValue),
			data: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			want: Words{
				padL("0102030405060708090a0b0c0d0e0f1011121314"), // data
			},
		},
		{
			name:    "[]byte->address/short",
			val:     new(AddressValue),
			data:    []byte{1, 2},
			wantErr: true,
		},
		{
			name:    "[]byte->address/long",
			val:     new(AddressValue),
			data:    []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21},
			wantErr: true,
		},
		{
			name: "array->address",
			val:  new(AddressValue),
			data: [20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			want: Words{
				padL("0102030405060708090a0b0c0d0e0f1011121314"), // data
			},
		},
		{
			name:    "array->address/short",
			val:     new(AddressValue),
			data:    [2]byte{1, 2},
			wantErr: true,
		},
		{
			name:    "array->address/long",
			val:     new(AddressValue),
			data:    [21]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21},
			wantErr: true,
		},
		{
			name:    "int->address",
			val:     new(AddressValue),
			data:    42,
			wantErr: true,
		},
		{
			name: "types.Address->address",
			val:  new(AddressValue),
			data: types.MustAddressFromHex("0102030405060708090a0b0c0d0e0f1011121314"),
			want: Words{
				padL("0102030405060708090a0b0c0d0e0f1011121314"),
			},
		},
		{
			name:    "types.Hash->address",
			val:     new(AddressValue),
			data:    types.MustHashFromHex("0102030405060708090a0b0c0d0e0f10111213145566778899aabbccddeeff00", types.PadNone),
			wantErr: true,
		},
		{
			name:    "types.Number->address",
			val:     new(AddressValue),
			data:    types.MustNumberFromHex("2a"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Default.Mapper.Map(tt.data, tt.val)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				abi, err := tt.val.EncodeABI()
				require.NoError(t, err)
				assert.Equal(t, tt.want, abi)
			}
		})
	}
}

func TestMapTo(t *testing.T) {
	tests := []struct {
		name    string
		arg     any
		val     Value
		want    any
		wantErr bool
	}{
		// BytesValue:
		{
			name: "bytes->byte[]",
			arg:  new([]byte),
			val:  func() Value { b := BytesValue([]byte{1, 2, 3}); return &b }(),
			want: func() *[]byte { b := []byte{1, 2, 3}; return &b }(),
		},
		{
			name: "bytes->array",
			arg:  new([3]byte),
			val:  func() Value { b := BytesValue([]byte{1, 2, 3}); return &b }(),
			want: func() *[3]byte { b := [3]byte{1, 2, 3}; return &b }(),
		},
		{
			name: "bytes->string",
			arg:  new(string),
			val:  func() Value { v := BytesValue("foo"); return &v }(),
			want: func() *string { s := "0x666f6f"; return &s }(),
		},
		{
			name:    "bytes->int64",
			arg:     new(int64),
			val:     func() Value { b := BytesValue([]byte{1, 2, 3}); return &b }(),
			wantErr: true,
		},
		{
			name:    "bytes->big.Int",
			arg:     new(big.Int),
			val:     func() Value { b := BytesValue([]byte{1, 2, 3}); return &b }(),
			wantErr: true,
		},
		// StringValue:
		{
			name: "string->string",
			arg:  new(string),
			val:  func() Value { v := StringValue("foo"); return &v }(),
			want: func() *string { s := "foo"; return &s }(),
		},
		{
			name: "string->[]byte",
			arg:  new([]byte),
			val:  func() Value { s := StringValue("foo"); return &s }(),
			want: func() *[]byte { b := []byte("foo"); return &b }(),
		},
		{
			name:    "string->array",
			arg:     new([3]byte),
			val:     func() Value { s := StringValue("foo"); return &s }(),
			wantErr: true,
		},
		{
			name:    "string->int64",
			arg:     new(int64),
			val:     func() Value { s := StringValue("foo"); return &s }(),
			wantErr: true,
		},
		{
			name:    "string->big.Int",
			arg:     new(big.Int),
			val:     func() Value { s := StringValue("foo"); return &s }(),
			wantErr: true,
		},
		{
			name:    "string->types.Address",
			arg:     new(types.Address),
			val:     func() Value { s := StringValue("0x0102030405060708090a0b0c0d0e0f1011121314"); return &s }(),
			wantErr: true,
		},
		{
			name: "string->types.Hash",
			arg:  new(types.Hash),
			val: func() Value {
				s := StringValue("0x0102030405060708090a0b0c0d0e0f10111213145566778899aabbccddeeff00")
				return &s
			}(),
			wantErr: true,
		},
		{
			name: "string->types.Number",
			arg:  new(types.Number),
			val: func() Value {
				s := StringValue("0x0102030405060708090a0b0c0d0e0f10111213145566778899aabbccddeeff00")
				return &s
			}(),
			wantErr: true,
		},
		// FixedBytesValue:
		{
			name: "bytes4->[]byte",
			arg:  new([]byte),
			val:  func() Value { v := FixedBytesValue([]byte{1, 2, 3, 4}); return &v }(),
			want: func() *[]byte { b := []byte{1, 2, 3, 4}; return &b }(),
		},
		{
			name: "bytes4->array",
			arg:  new([4]byte),
			val:  func() Value { v := FixedBytesValue([]byte{1, 2, 3, 4}); return &v }(),
			want: func() *[4]byte { b := [4]byte{1, 2, 3, 4}; return &b }(),
		},
		{
			name: "bytes32->string",
			arg:  new(string),
			val: func() Value {
				v := FixedBytesValue(hexutil.MustHexToBytes("0x0102030405060708090a0b0c0d0e0f10111213145566778899aabbccddeeff00"))
				return &v
			}(),
			want: func() *string { s := "0x0102030405060708090a0b0c0d0e0f10111213145566778899aabbccddeeff00"; return &s }(),
		},
		{
			name: "bytes32->uint64",
			arg:  new(uint64),
			val: func() Value {
				v := FixedBytesValue(hexutil.MustHexToBytes("0x000000000000000000000000000000000000000000000000000000000000002a"))
				return &v
			}(),
			want: func() *uint64 { i := uint64(42); return &i }(),
		},
		{
			name: "bytes16->uint64",
			arg:  new(uint64),
			val: func() Value {
				v := FixedBytesValue(hexutil.MustHexToBytes("0x0000000000000000000000000000002a"))
				return &v
			}(),
			wantErr: true,
		},
		{
			name: "bytes32->int64",
			arg:  new(int64),
			val: func() Value {
				v := FixedBytesValue(hexutil.MustHexToBytes("0x000000000000000000000000000000000000000000000000000000000000002a"))
				return &v
			}(),
			want: func() *int64 { i := int64(42); return &i }(),
		},
		{
			name: "bytes32->int64#negative",
			arg:  new(int64),
			val: func() Value {
				v := FixedBytesValue(hexutil.MustHexToBytes("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd6"))
				return &v
			}(),
			wantErr: true,
		},
		{
			name: "bytes16->int64",
			arg:  new(int64),
			val: func() Value {
				v := FixedBytesValue(hexutil.MustHexToBytes("0x0000000000000000000000000000002a"))
				return &v
			}(),
			wantErr: true,
		},
		{
			name: "bytes32->big.Int",
			arg:  new(big.Int),
			val: func() Value {
				v := FixedBytesValue(hexutil.MustHexToBytes("0x000000000000000000000000000000000000000000000000000000000000002a"))
				return &v
			}(),
			want: func() *big.Int { i := big.NewInt(42); return i }(),
		},
		{
			name: "bytes16->big.Int",
			arg:  new(big.Int),
			val: func() Value {
				v := FixedBytesValue(hexutil.MustHexToBytes("0x0000000000000000000000000000002a"))
				return &v
			}(),
			wantErr: true,
		},
		{
			name: "bytes32->types.Address",
			arg:  new(types.Address),
			val: func() Value {
				v := FixedBytesValue(hexutil.MustHexToBytes("0x0102030405060708090a0b0c0d0e0f1011121314000000000000000000000000"))
				return &v
			}(),
			wantErr: true,
		},
		{
			name: "bytes20->types.Address",
			arg:  new(types.Address),
			val: func() Value {
				v := FixedBytesValue(hexutil.MustHexToBytes("0x0102030405060708090a0b0c0d0e0f1011121314"))
				return &v
			}(),
			want: types.MustAddressFromHexPtr("0x0102030405060708090a0b0c0d0e0f1011121314"),
		},
		{
			name: "bytes32->types.Hash",
			arg:  new(types.Hash),
			val: func() Value {
				v := FixedBytesValue(hexutil.MustHexToBytes("0x0102030405060708090a0b0c0d0e0f1011121314112233445566778899aabbcc"))
				return &v
			}(),
			want: types.MustHashFromHexPtr("0x0102030405060708090a0b0c0d0e0f1011121314112233445566778899aabbcc", types.PadNone),
		},
		{
			name: "bytes32->types.Number",
			arg:  new(types.Number),
			val: func() Value {
				v := FixedBytesValue(hexutil.MustHexToBytes("0x000000000000000000000000000000000000000000000000000000000000002a"))
				return &v
			}(),
			want: types.MustNumberFromHexPtr("0x2a"),
		},
		{
			name: "bytes16->types.Number",
			arg:  new(types.Number),
			val: func() Value {
				v := FixedBytesValue(hexutil.MustHexToBytes("0x0000000000000000000000000000002a"))
				return &v
			}(),
			wantErr: true,
		},
		// UintValue:
		{
			name: "uint256->string",
			arg:  new(string),
			val:  func() Value { i := &UintValue{Size: 256}; i.SetUint64(42); return i }(),
			want: func() *string { s := "0x2a"; return &s }(),
		},
		{
			name:    "uint256->[]byte",
			arg:     new([]byte),
			val:     func() Value { i := &UintValue{Size: 256}; i.SetUint64(42); return i }(),
			wantErr: true,
		},
		{
			name:    "uint256->array",
			arg:     new([1]byte),
			val:     func() Value { i := &UintValue{Size: 256}; i.SetUint64(42); return i }(),
			wantErr: true,
		},
		{
			name: "uint256->int64",
			arg:  new(int64),
			val:  func() Value { i := &UintValue{Size: 256}; i.SetUint64(42); return i }(),
			want: func() *int64 { i := int64(42); return &i }(),
		},
		{
			name: "uint16->uint8#fit",
			arg:  new(uint8),
			val:  func() Value { i := &UintValue{Size: 16}; i.SetUint64(42); return i }(),
			want: func() *uint8 { i := uint8(42); return &i }(),
		},
		{
			name:    "uint16->uint8#overflow",
			arg:     new(uint8),
			val:     func() Value { i := &UintValue{Size: 16}; i.SetUint64(256); return i }(),
			wantErr: true,
		},
		{
			name: "uint16->int8#fit",
			arg:  new(int8),
			val:  func() Value { i := &UintValue{Size: 16}; i.SetUint64(42); return i }(),
			want: func() *int8 { i := int8(42); return &i }(),
		},
		{
			name: "uint16->int8#overflow",

			arg:     new(int8),
			val:     func() Value { i := &UintValue{Size: 16}; i.SetUint64(128); return i }(),
			wantErr: true,
		},
		{
			name: "uint256->big.Int",
			arg:  new(big.Int),
			val:  func() Value { i := &UintValue{Size: 256}; i.SetUint64(42); return i }(),
			want: func() *big.Int { i := big.NewInt(42); return i }(),
		},
		{
			name:    "uint256->types.Address",
			arg:     new(types.Address),
			val:     func() Value { i := &UintValue{Size: 256}; i.SetUint64(42); return i }(),
			wantErr: true,
		},
		{
			name:    "uint256->types.Hash",
			arg:     new(types.Hash),
			val:     func() Value { i := &UintValue{Size: 256}; i.SetUint64(42); return i }(),
			wantErr: true,
		},
		{
			name: "uint256->types.Number",
			arg:  new(types.Number),
			val:  func() Value { i := &UintValue{Size: 256}; i.SetUint64(42); return i }(),
			want: types.MustNumberFromHexPtr("0x2a"),
		},
		// IntValue:
		{
			name: "int256->string",
			arg:  new(string),
			val:  func() Value { v := &IntValue{Size: 256}; v.SetUint64(42); return v }(),
			want: func() *string { s := "0x2a"; return &s }(),
		},
		{
			name: "int256->string#negative",
			arg:  new(string),
			val:  func() Value { v := &IntValue{Size: 256}; v.SetInt64(-42); return v }(),
			want: func() *string { s := "-0x2a"; return &s }(),
		},
		{
			name:    "int256->[]byte",
			arg:     new([]byte),
			val:     func() Value { i := &IntValue{Size: 256}; i.SetUint64(42); return i }(),
			wantErr: true,
		},
		{
			name:    "int256->array",
			arg:     new([1]byte),
			val:     func() Value { i := &IntValue{Size: 256}; i.SetUint64(42); return i }(),
			wantErr: true,
		},
		{
			name: "int256->int64",
			arg:  new(int64),
			val:  func() Value { i := &IntValue{Size: 256}; i.SetUint64(42); return i }(),
			want: func() *int64 { i := int64(42); return &i }(),
		},
		{
			name: "int256->int64#negative",
			arg:  new(int64),
			val:  func() Value { i := &IntValue{Size: 256}; i.SetInt64(-42); return i }(),
			want: func() *int64 { i := int64(-42); return &i }(),
		},
		{
			name: "int16->uint8#fit",
			arg:  new(uint8),
			val:  func() Value { i := &IntValue{Size: 16}; i.SetUint64(42); return i }(),
			want: func() *uint8 { i := uint8(42); return &i }(),
		},
		{
			name:    "int16->uint8#overflow",
			arg:     new(uint8),
			val:     func() Value { i := &IntValue{Size: 16}; i.SetUint64(256); return i }(),
			wantErr: true,
		},
		{
			name: "int16->int8#fit",
			arg:  new(int8),
			val:  func() Value { i := &IntValue{Size: 16}; i.SetUint64(42); return i }(),
			want: func() *int8 { i := int8(42); return &i }(),
		},
		{
			name:    "int16->int8#overflow",
			arg:     new(int8),
			val:     func() Value { i := &IntValue{Size: 16}; i.SetUint64(128); return i }(),
			wantErr: true,
		},
		{
			name: "int256->big.Int",
			arg:  new(big.Int),
			val:  func() Value { i := &IntValue{Size: 256}; i.SetUint64(42); return i }(),
			want: func() *big.Int { i := big.NewInt(42); return i }(),
		},
		{
			name:    "int256->types.Address",
			arg:     new(types.Address),
			val:     func() Value { i := &IntValue{Size: 256}; i.SetUint64(42); return i }(),
			wantErr: true,
		},
		{
			name:    "int256->types.Hash",
			arg:     new(types.Hash),
			val:     func() Value { i := &IntValue{Size: 256}; i.SetUint64(42); return i }(),
			wantErr: true,
		},
		{
			name: "int256->types.Number",
			arg:  new(types.Number),
			val:  func() Value { i := &IntValue{Size: 256}; i.SetUint64(42); return i }(),
			want: types.MustNumberFromHexPtr("0x2a"),
		},
		// BoolValue:
		{
			name: "bool->bool",
			arg:  new(bool),
			val:  func() Value { b := BoolValue(true); return &b }(),
			want: func() *bool { b := true; return &b }(),
		},
		{
			name:    "bool->string",
			arg:     new(string),
			val:     func() Value { v := BoolValue(true); return &v }(),
			wantErr: true,
		},
		{
			name:    "bool->int64",
			arg:     new(int64),
			val:     func() *BoolValue { b := BoolValue(true); return &b }(),
			wantErr: true,
		},
		// AddressValue:
		{
			name: "address->string",
			arg:  new(string),
			val:  (*AddressValue)(types.MustAddressFromHexPtr("0x1234567890123456789012345678901234567890")),
			want: func() *string { s := "0x1234567890123456789012345678901234567890"; return &s }(),
		},
		{
			name: "address->[]byte",
			arg:  new([]byte),
			val:  (*AddressValue)(types.MustAddressFromHexPtr("0x1234567890123456789012345678901234567890")),
			want: func() *[]byte {
				b := []byte{0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90}
				return &b
			}(),
		},
		{
			name: "address->array",
			arg:  new([20]byte),
			val:  (*AddressValue)(types.MustAddressFromHexPtr("0x1234567890123456789012345678901234567890")),
			want: func() *[20]byte {
				b := [20]byte{0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12, 0x34, 0x56, 0x78, 0x90}
				return &b
			}(),
		},
		{
			name:    "address->int64",
			arg:     new(int64),
			val:     (*AddressValue)(types.MustAddressFromHexPtr("0x1234567890123456789012345678901234567890")),
			wantErr: true,
		},
		{
			name:    "address->big.Int",
			arg:     new(big.Int),
			val:     (*AddressValue)(types.MustAddressFromHexPtr("0x1234567890123456789012345678901234567890")),
			wantErr: true,
		},
		{
			name: "address->types.Address",
			arg:  new(types.Address),
			val:  (*AddressValue)(types.MustAddressFromHexPtr("0x1234567890123456789012345678901234567890")),
			want: types.MustAddressFromHexPtr("0x1234567890123456789012345678901234567890"),
		},
		{
			name:    "address->types.Hash",
			arg:     new(types.Hash),
			val:     (*AddressValue)(types.MustAddressFromHexPtr("0x1234567890123456789012345678901234567890")),
			wantErr: true,
		},
		{
			name:    "address->types.Number",
			arg:     new(types.Number),
			val:     (*AddressValue)(types.MustAddressFromHexPtr("0x1234567890123456789012345678901234567890")),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Default.Mapper.Map(tt.val, tt.arg)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tt.want, tt.arg)
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

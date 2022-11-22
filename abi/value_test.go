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
		wantABI Words
		wantErr bool
	}{
		// TupleValue:
		{
			name:    "tuple/empty",
			val:     &TupleValue{},
			arg:     struct{}{},
			wantABI: Words(nil),
		},
		{
			name: "tuple/static",
			val: &TupleValue{
				TupleValueElem{Value: new(BoolValue), Name: "a"},
				TupleValueElem{Value: new(BoolValue), Name: "b"},
			},
			arg:     map[string]bool{"a": true, "b": true},
			wantABI: Words{padL("1"), padL("1")},
		},
		{
			name: "tuple/dynamic",
			val: &TupleValue{
				TupleValueElem{Value: new(BytesValue), Name: "a"},
			},
			arg: map[string][]byte{"a": {1, 2, 3}},
			wantABI: Words{
				padL("20"),     // offset
				padL("3"),      // length
				padR("010203"), // data
			},
		},
		{
			name: "tuple/static-and-dynamic",
			val: &TupleValue{
				TupleValueElem{Value: new(BoolValue), Name: "a"},
				TupleValueElem{Value: new(BytesValue), Name: "b"},
			},
			arg: map[string]interface{}{"a": true, "b": []byte{1, 2, 3}},
			wantABI: Words{
				padL("1"),      // a
				padL("40"),     // offset to b
				padL("3"),      // length of b
				padR("010203"), // b
			},
		},
		{
			name: "tuple/nested",
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
			wantABI: Words{
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
			name:    "array/empty",
			val:     &ArrayValue{Type: NewBoolType()},
			arg:     []bool{},
			wantABI: Words{padL("0")},
		},
		{
			name: "array/two-static-elements",
			val:  &ArrayValue{Type: NewBoolType()},
			arg:  []bool{true, true},
			wantABI: Words{
				padL("2"), // array length
				padL("1"), // first element
				padL("1"), // second element
			},
		},
		{
			name: "array/two-dynamic-elements",
			val:  &ArrayValue{Type: NewBytesType()},
			arg:  [][]byte{{1, 2, 3}, {4, 5, 6}},
			wantABI: Words{
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
			name:    "fixed-array/empty",
			val:     &FixedArrayValue{},
			arg:     [0]bool{},
			wantABI: nil,
		},
		{
			name: "fixed-array/two-static-elements",
			val:  &FixedArrayValue{new(BoolValue), new(BoolValue)},
			arg:  []bool{true, true},
			wantABI: Words{
				padL("1"), // first element
				padL("1"), // second element
			},
		},
		{
			name: "fixed-array/two-dynamic-elements",
			val:  &FixedArrayValue{new(BytesValue), new(BytesValue)},
			arg:  [][]byte{{1, 2, 3}, {4, 5, 6}},
			wantABI: Words{
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
			name:    "bytes/empty",
			val:     new(BytesValue),
			arg:     []byte{},
			wantABI: Words{padL("0")},
		},
		{
			name: "bytes/one-word",
			val:  new(BytesValue),
			arg:  []byte{1, 2, 3},
			wantABI: Words{
				padL("3"),      // length
				padR("010203"), // data
			},
		},
		{
			name: "bytes/two-words",
			val:  new(BytesValue),
			arg:  bytes.Repeat([]byte{1}, 33),
			wantABI: Words{
				padL("21"), // length
				padR("0101010101010101010101010101010101010101010101010101010101010101"), // data
				padR("01"), // data
			},
		},

		// StringValue:
		{
			name:    "string/empty",
			val:     new(StringValue),
			arg:     "",
			wantABI: Words{padL("0")},
		},
		{
			name: "string/one-word",
			val:  new(StringValue),
			arg:  "abc",
			wantABI: Words{
				padL("3"),      // length
				padR("616263"), // data
			},
		},
		{
			name: "string/two-words",
			val:  new(StringValue),
			arg:  strings.Repeat("a", 33),
			wantABI: Words{
				padL("21"), // length
				padR("6161616161616161616161616161616161616161616161616161616161616161"), // data
				padR("61"), // data
			},
		},
		// FixedBytesValue:
		{
			name:    "fixed/bytes-empty",
			val:     make(FixedBytesValue, 0),
			arg:     []byte{},
			wantABI: Words{padL("0")},
		},
		{
			name: "fixed/bytes-non-empty",
			val:  make(FixedBytesValue, 3),
			arg:  []byte{1, 2, 3},
			wantABI: Words{
				padR("010203"), // data
			},
		},
		// UintValue:
		{
			name:    "uint256/0",
			val:     &UintValue{Size: 256},
			arg:     big.NewInt(0),
			wantABI: Words{padL("0")},
		},
		{
			name:    "uint256/MaxUint256",
			val:     &UintValue{Size: 256},
			arg:     new(big.Int).Set(MaxUint[256]),
			wantABI: Words{padR("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")},
		},
		// IntValue:
		{
			name:    "int256/0",
			val:     &IntValue{Size: 256},
			arg:     big.NewInt(0),
			wantABI: Words{padL("0")},
		},
		{
			name:    "int256/1",
			val:     &IntValue{Size: 256},
			arg:     big.NewInt(1),
			wantABI: Words{padL("1")},
		},
		{
			name:    "int256/-1",
			val:     &IntValue{Size: 256},
			arg:     big.NewInt(-1),
			wantABI: Words{padR("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")},
		},
		{
			name:    "int256/MaxInt256",
			val:     &IntValue{Size: 256},
			arg:     new(big.Int).Set(MaxInt[256]),
			wantABI: Words{padR("7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")},
		},
		{
			name:    "int256/MinInt256",
			val:     &IntValue{Size: 256},
			arg:     new(big.Int).Set(MinInt[256]),
			wantABI: Words{padR("8000000000000000000000000000000000000000000000000000000000000000")},
		},
		{
			name:    "int8/127",
			val:     &IntValue{Size: 256},
			arg:     big.NewInt(127),
			wantABI: Words{padL("7f")},
		},
		{
			name:    "int8/-128",
			val:     &IntValue{Size: 256},
			arg:     big.NewInt(-128),
			wantABI: Words{padR("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80")},
		},
		// BoolValue:
		{
			name:    "bool/true",
			val:     new(BoolValue),
			arg:     true,
			wantABI: Words{padL("1")},
		},
		{
			name:    "bool/false",
			val:     new(BoolValue),
			arg:     false,
			wantABI: Words{padL("0")},
		},
		// AddressValue:
		{
			name:    "address/empty",
			val:     new(AddressValue),
			arg:     types.Address{},
			wantABI: Words{padL("0")},
		},
		{
			name:    "address/non-empty",
			val:     new(AddressValue),
			arg:     types.MustHexToAddress("0x0102030405060708090a0b0c0d0e0f1011121314"),
			wantABI: Words{padL("0102030405060708090a0b0c0d0e0f1011121314")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, DefaultConfig.Mapper.Map(tt.arg, tt.val))
			enc, err := tt.val.EncodeABI()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantABI, enc)
			}
		})
	}
}

func TestDecodeABI(t *testing.T) {
	tests := []struct {
		name    string
		abi     Words
		val     Value
		wantVal Value
		wantErr bool
	}{
		// TupleValue:
		{
			name:    "tuple/empty",
			abi:     Words{},
			val:     new(TupleValue),
			wantVal: new(TupleValue),
		},
		{
			name: "tuple/static",
			abi: Words{
				padL("1"),
				padL("1"),
			},
			val: &TupleValue{
				TupleValueElem{Value: new(BoolValue), Name: "a"},
				TupleValueElem{Value: new(BoolValue), Name: "b"},
			},
			wantVal: &TupleValue{
				TupleValueElem{Value: func() *BoolValue { b := BoolValue(true); return &b }(), Name: "a"},
				TupleValueElem{Value: func() *BoolValue { b := BoolValue(true); return &b }(), Name: "b"},
			},
		},
		{
			name: "tuple/dynamic",
			abi: Words{
				padL("20"),     // offset
				padL("3"),      // length
				padR("010203"), // data
			},
			val: &TupleValue{
				TupleValueElem{Value: new(BytesValue), Name: "a"},
			},
			wantVal: &TupleValue{
				TupleValueElem{Value: func() *BytesValue { b := BytesValue([]byte{1, 2, 3}); return &b }(), Name: "a"},
			},
		},
		{
			name: "tuple/static-and-dynamic",
			abi: Words{
				padL("1"),      // a
				padL("40"),     // offset to b
				padL("3"),      // length of b
				padR("010203"), // b
			},
			val: &TupleValue{
				TupleValueElem{Value: new(BoolValue), Name: "a"},
				TupleValueElem{Value: new(BytesValue), Name: "b"},
			},
			wantVal: &TupleValue{
				TupleValueElem{Value: func() *BoolValue { b := BoolValue(true); return &b }(), Name: "a"},
				TupleValueElem{Value: func() *BytesValue { b := BytesValue([]byte{1, 2, 3}); return &b }(), Name: "b"},
			},
		},
		{
			name: "tuple/nested",
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
			wantVal: &TupleValue{
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
			name:    "array/empty",
			abi:     Words{padL("0")},
			val:     &ArrayValue{Type: NewBoolType()},
			wantVal: &ArrayValue{Type: NewBoolType(), Elems: []Value{}},
		},
		{
			name: "array/two-static-elements",
			abi: Words{
				padL("2"), // array length
				padL("1"), // first element
				padL("1"), // second element
			},
			val: &ArrayValue{Type: NewBoolType()},
			wantVal: &ArrayValue{
				Type: NewBoolType(),
				Elems: []Value{
					func() *BoolValue { b := BoolValue(true); return &b }(),
					func() *BoolValue { b := BoolValue(true); return &b }(),
				},
			},
		},
		{
			name: "array/two-dynamic-elements",
			abi: Words{
				padL("2"),      // array length
				padL("40"),     // offset to first element
				padL("80"),     // offset to second element
				padL("3"),      // length of first element
				padR("010203"), // first element
				padL("3"),      // length of second element
				padR("040506"), // second element
			},
			val: &ArrayValue{Type: NewBytesType()},
			wantVal: &ArrayValue{
				Type: NewBytesType(),
				Elems: []Value{
					func() *BytesValue { b := BytesValue([]byte{1, 2, 3}); return &b }(),
					func() *BytesValue { b := BytesValue([]byte{4, 5, 6}); return &b }(),
				},
			},
		},
		// FixedArrayValue:
		{
			name:    "fixed-array/empty",
			val:     &FixedArrayValue{},
			abi:     Words{padL("0")},
			wantVal: &FixedArrayValue{},
		},
		{
			name: "fixed-array/two-static-elements",
			abi: Words{
				padL("1"), // first element
				padL("1"), // second element
			},
			val: &FixedArrayValue{new(BoolValue), new(BoolValue)},
			wantVal: &FixedArrayValue{
				func() *BoolValue { b := BoolValue(true); return &b }(),
				func() *BoolValue { b := BoolValue(true); return &b }(),
			},
		},
		{
			name: "fixed-array/two-dynamic-elements",
			abi: Words{
				padL("40"),     // offset to first element
				padL("80"),     // offset to second element
				padL("3"),      // length of first element
				padR("010203"), // first element
				padL("3"),      // length of second element
				padR("040506"), // second element
			},
			val: &FixedArrayValue{new(BytesValue), new(BytesValue)},
			wantVal: &FixedArrayValue{
				func() *BytesValue { b := BytesValue([]byte{1, 2, 3}); return &b }(),
				func() *BytesValue { b := BytesValue([]byte{4, 5, 6}); return &b }(),
			},
		},
		// BytesValue:
		{
			name:    "bytes/empty",
			abi:     Words{padL("0")},
			val:     new(BytesValue),
			wantVal: func() *BytesValue { b := BytesValue([]byte{}); return &b }(),
		},
		{
			name: "bytes/one-word",
			abi: Words{
				padL("3"),      // length
				padR("010203"), // data
			},
			val:     new(BytesValue),
			wantVal: func() *BytesValue { b := BytesValue([]byte{1, 2, 3}); return &b }(),
		},
		{
			name: "bytes/two-words",
			abi: Words{
				padL("21"), // length
				padR("0101010101010101010101010101010101010101010101010101010101010101"), // data
				padR("01"), // data
			},
			val:     new(BytesValue),
			wantVal: func() *BytesValue { b := BytesValue(bytes.Repeat([]byte{1}, 33)); return &b }(),
		},
		// StringValue:
		{
			name:    "string/empty",
			abi:     Words{padL("0")},
			val:     new(StringValue),
			wantVal: func() *StringValue { s := StringValue(""); return &s }(),
		},
		{
			name: "string/one-word",
			abi: Words{
				padL("3"),      // length
				padR("616263"), // data
			},
			val:     new(StringValue),
			wantVal: func() *StringValue { s := StringValue("abc"); return &s }(),
		},
		{
			name: "string/two-words",
			abi: Words{
				padL("21"), // length
				padR("6161616161616161616161616161616161616161616161616161616161616161"), // data
				padR("61"), // data
			},
			val:     new(StringValue),
			wantVal: func() *StringValue { s := StringValue(strings.Repeat("a", 33)); return &s }(),
		},
		// FixedBytesValue:
		{
			name:    "fixed/bytes-empty",
			abi:     Words{padL("0")},
			val:     make(FixedBytesValue, 0),
			wantVal: make(FixedBytesValue, 0),
		},
		{
			name: "fixed/bytes-non-empty",
			abi: Words{
				padR("010203"), // data
			},
			val:     make(FixedBytesValue, 3),
			wantVal: func() FixedBytesValue { b := FixedBytesValue([]byte{1, 2, 3}); return b }(),
		},
		// UintValue:
		{
			name:    "uint256/0",
			abi:     Words{padL("0")},
			val:     &UintValue{Size: 256},
			wantVal: func() *UintValue { u := UintValue{Size: 256, Int: *big.NewInt(0)}; return &u }(),
		},
		{
			name:    "uint256/MaxUint256",
			abi:     Words{padR("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")},
			val:     &UintValue{Size: 256},
			wantVal: func() *UintValue { u := UintValue{Size: 256, Int: *MaxUint[256]}; return &u }(),
		},
		// IntValue:
		{
			name:    "int256/0",
			abi:     Words{padL("0")},
			val:     &IntValue{Size: 256},
			wantVal: func() *IntValue { i := IntValue{Size: 256, Int: *big.NewInt(0)}; return &i }(),
		},
		{
			name:    "int256/1",
			abi:     Words{padL("1")},
			val:     &IntValue{Size: 256},
			wantVal: func() *IntValue { i := IntValue{Size: 256, Int: *big.NewInt(1)}; return &i }(),
		},
		{
			name:    "int256/-1",
			abi:     Words{padR("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")},
			val:     &IntValue{Size: 256},
			wantVal: func() *IntValue { i := IntValue{Size: 256, Int: *big.NewInt(-1)}; return &i }(),
		},
		{
			name:    "int256/MaxInt256",
			abi:     Words{padR("7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")},
			val:     &IntValue{Size: 256},
			wantVal: func() *IntValue { i := IntValue{Size: 256, Int: *MaxInt[256]}; return &i }(),
		},
		{
			name:    "int256/MinInt256",
			abi:     Words{padR("8000000000000000000000000000000000000000000000000000000000000000")},
			val:     &IntValue{Size: 256},
			wantVal: func() *IntValue { i := IntValue{Size: 256, Int: *MinInt[256]}; return &i }(),
		},
		{
			name:    "int8/127",
			abi:     Words{padL("7f")},
			val:     &IntValue{Size: 256},
			wantVal: func() *IntValue { i := IntValue{Size: 256, Int: *big.NewInt(127)}; return &i }(),
		},
		{
			name:    "int8/-128",
			abi:     Words{padR("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80")},
			val:     &IntValue{Size: 256},
			wantVal: func() *IntValue { i := IntValue{Size: 256, Int: *big.NewInt(-128)}; return &i }(),
		},
		// BoolValue:
		{
			name:    "bool/true",
			abi:     Words{padL("1")},
			val:     new(BoolValue),
			wantVal: func() *BoolValue { b := BoolValue(true); return &b }(),
		},
		{
			name:    "bool/false",
			abi:     Words{padL("0")},
			val:     new(BoolValue),
			wantVal: func() *BoolValue { b := BoolValue(false); return &b }(),
		},
		// AddressValue:
		{
			name:    "address/empty",
			abi:     Words{padL("0")},
			val:     new(AddressValue),
			wantVal: func() *AddressValue { a := AddressValue{}; return &a }(),
		},
		{
			name: "address/non-empty",
			abi:  Words{padL("0102030405060708090a0b0c0d0e0f1011121314")},
			val:  new(AddressValue),
			wantVal: func() *AddressValue {
				a := types.MustHexToAddress("0102030405060708090a0b0c0d0e0f1011121314")
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
				assert.Equal(t, tt.wantVal, tt.val)
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

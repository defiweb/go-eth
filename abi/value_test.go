package abi

import (
	"bytes"
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
		enc     any
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
				TupleValueElem{Value: NewBoolValue(), Name: "arg0"},
				TupleValueElem{Value: NewBoolValue(), Name: "arg1"},
			),
			enc: map[string]any{"arg0": true, "arg1": true},
			abi: Words{padL("1"), padL("1")},
		},
		{
			name: "tuple/dynamic",
			val: NewTupleValue(
				TupleValueElem{Value: NewBytesValue(), Name: "arg0"},
			),
			enc: map[string]any{"arg0": []byte{0x01, 0x02, 0x03}},
			abi: Words{
				padL("20"),     // offset
				padL("3"),      // length
				padR("010203"), // data
			},
		},
		{
			name: "tuple/static-and-dynamic",
			val: NewTupleValue(
				TupleValueElem{Value: NewBoolValue(), Name: "arg0"},
				TupleValueElem{Value: NewBytesValue(), Name: "arg1"},
			),
			enc: map[string]any{"arg0": true, "arg1": []byte{0x01, 0x02, 0x03}},
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
			enc: map[string]any{
				"a": map[string]any{"x": []byte{0x01, 0x02, 0x03}},
				"b": map[string]any{"x": []byte{0x04, 0x05, 0x06}},
			},
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
			val:  NewArrayValue(NewBoolType(), NewBoolValue(), NewBoolValue()),
			enc:  []any{true, true},
			abi: Words{
				padL("2"), // array length
				padL("1"), // first element
				padL("1"), // second element
			},
		},
		{
			name: "array/two-dynamic-elements",
			val:  NewArrayValue(NewBytesType(), NewBytesValue(), NewBytesValue()),
			enc:  []any{[]byte{0x01, 0x02, 0x03}, []byte{0x04, 0x05, 0x06}},
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
			name: "fixed/array-empty",
			val:  NewFixedArrayValue(NewBoolType(), 0),
			abi:  nil,
		},
		{
			name: "fixed/array-two-static-elements",
			val:  NewFixedArrayValue(NewBoolType(), 2),
			enc:  []any{true, true},
			abi: Words{
				padL("1"), // first element
				padL("1"), // second element
			},
		},
		{
			name: "fixed/array-two-dynamic-elements",
			val:  NewFixedArrayValue(NewBytesType(), 2),
			enc:  []any{[]byte{0x01, 0x02, 0x03}, []byte{0x04, 0x05, 0x06}},
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
			val:  NewBytesValue(),
			enc:  []byte{0x01, 0x02, 0x03},
			abi: Words{
				padL("3"),      // length
				padR("010203"), // data
			},
		},
		{
			name: "bytes/two-words",
			val:  NewBytesValue(),
			enc:  bytes.Repeat([]byte{0x01}, 33),
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
			val:  NewStringValue(),
			enc:  "abc",
			abi: Words{
				padL("3"),      // length
				padR("616263"), // data
			},
		},
		{
			name: "string/two-words",
			val:  NewStringValue(),
			enc:  strings.Repeat("a", 33),
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
			val:  NewFixedBytesValue(32),
			enc:  []byte{0x01, 0x02, 0x03},
			abi: Words{
				padR("010203"), // data
			},
		},
		// UintValue:
		{
			name: "uint256/0",
			val:  NewUintValue(256),
			enc:  uint64(0x00),
			abi:  Words{padL("0")},
		},
		{
			name: "uint256/MaxUint256",
			val:  NewUintValue(256),
			enc:  MaxUint256,
			abi:  Words{padR("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")},
		},
		// IntValue:
		{
			name: "int256/0",
			val:  NewIntValue(256),
			enc:  int64(0x00),
			abi:  Words{padL("0")},
		},
		{
			name: "int256/1",
			val:  NewIntValue(256),
			enc:  int64(0x01),
			abi:  Words{padL("1")},
		},
		{
			name: "int256/-1",
			val:  NewIntValue(256),
			enc:  int64(-0x01),
			abi:  Words{padR("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")},
		},
		{
			name: "int256/MaxInt256",
			val:  NewIntValue(256),
			enc:  MaxInt256,
			abi:  Words{padR("7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")},
		},
		{
			name: "int256/MinInt256",
			val:  NewIntValue(256),
			enc:  MinInt256,
			abi:  Words{padR("8000000000000000000000000000000000000000000000000000000000000000")},
		},
		{
			name: "int8/127",
			val:  NewIntValue(8),
			enc:  int64(127),
			abi:  Words{padL("7f")},
		},
		{
			name: "int8/-128",
			val:  NewIntValue(8),
			enc:  int64(-128),
			abi:  Words{padR("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80")},
		},
		// BoolValue:
		{
			name: "bool/true",
			val:  NewBoolValue(),
			enc:  true,
			abi:  Words{padL("1")},
		},
		{
			name: "bool/false",
			val:  NewBoolValue(),
			enc:  false,
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
			val:  NewAddressValue(),
			enc:  types.MustHexToAddress("0x0102030405060708090a0b0c0d0e0f1011121314"),
			abi:  Words{padL("0102030405060708090a0b0c0d0e0f1011121314")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.enc != nil {
				_, err := EncodeValue(tt.val, tt.enc)
				require.NoError(t, err)
			}
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

func padL(h string) (w Word) {
	_ = (&w).SetBytesPadLeft(hexutil.MustHexToBytes(h))
	return w
}

func padR(h string) (w Word) {
	_ = (&w).SetBytesPadRight(hexutil.MustHexToBytes(h))
	return w
}

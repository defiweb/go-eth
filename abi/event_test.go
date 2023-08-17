package abi

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/defiweb/go-eth/hexutil"
	"github.com/defiweb/go-eth/types"
)

func TestParseEvent(t *testing.T) {
	tests := []struct {
		signature string
		expected  string
		wantErr   bool
	}{
		{signature: "foo((uint256,bytes32)[])", expected: "event foo((uint256, bytes32)[])"},
		{signature: "foo((uint256 a, bytes32 b)[] c)", expected: "event foo((uint256 a, bytes32 b)[] c)"},
		{signature: "event foo(tuple(uint256 a, bytes32 b)[] c)", expected: "event foo((uint256 a, bytes32 b)[] c)"},
		{signature: "foo(uint256)(uint256)", wantErr: true},
		{signature: "constructor(uint256)", wantErr: true},
		{signature: "error foo(uint256)", wantErr: true},
		{signature: "function foo(uint256)", wantErr: true},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			e, err := ParseEvent(tt.signature)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, e.String())
			}
		})
	}
}

func TestEvent_DecodeValue(t *testing.T) {
	tests := []struct {
		signature string
		arg       any
		topics    []string
		data      string
		expected  any
		wantErr   bool
	}{
		{
			signature: "foo(uint256)",
			arg:       map[string]any{},
			topics:    []string{"0x2fbebd3821c4e005fbe0a9002cc1bd25dc266d788dba1dbcb39cc66a07e7b38b"},
			data:      "0000000000000000000000000000000000000000000000000000000000000001",
			expected:  map[string]any{"data0": big.NewInt(1)},
		},
		{
			signature: "foo((bytes a))",
			arg:       map[string]any{"data0": map[string]any{"a": []byte{}}},
			topics:    []string{"0x7a699d0514ec3b3aad6ef3992fd4993cacccf6906a6b41200cbd7c24d4dde537"},
			data:      "000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000004deadbeef00000000000000000000000000000000000000000000000000000000",
			expected:  map[string]any{"data0": map[string]any{"a": []byte{0xde, 0xad, 0xbe, 0xef}}},
		},
		{
			signature: "foo(address indexed a, address indexed b, (bytes32 ca, address cb, bytes cc) c, (uint128 da, uint32 db) d)",
			arg: map[string]any{
				"a": &types.Address{},
				"b": &types.Address{},
				"c": map[string]any{
					"ca": &types.Hash{},
					"cb": &types.Address{},
					"cc": []byte{},
				},
				"d": map[string]any{
					"da": big.NewInt(0),
					"db": big.NewInt(0),
				},
			},
			topics: []string{
				"0xe2f95411c22cc63e49510640786b88092a37c3926a08ba565b77c8d0fa08bb27",
				"0x0000000000000000000000001F7acDa376eF37EC371235a094113dF9Cb4EfEe1",
				"0x00000000000000000000000068E527780872cda0216Ba0d8fBD58b67a5D5e351",
			},
			data: "0000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000005b0000000000000000000000000000000000000000000000000000000064dcea4063cff3f05ab6a55e7e49095371098bf6455d27c6ed5b6e3c3178c661a821f7290000000000000000000000005b1742020856ad976469f498a296465d7803a7d100000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000014140310010a12090d020c11130e060b040507080f000000000000000000000000",
			expected: map[string]any{
				"a": types.MustAddressFromHexPtr("0x1F7acDa376eF37EC371235a094113dF9Cb4EfEe1"),
				"b": types.MustAddressFromHexPtr("0x68E527780872cda0216Ba0d8fBD58b67a5D5e351"),
				"c": map[string]any{
					"ca": types.MustHashFromHexPtr("0x63cff3f05ab6a55e7e49095371098bf6455d27c6ed5b6e3c3178c661a821f729", types.PadNone),
					"cb": types.MustAddressFromHexPtr("0x5b1742020856ad976469f498a296465d7803a7d1"),
					"cc": hexutil.MustHexToBytes("0x140310010a12090d020c11130e060b040507080f"),
				},
				"d": map[string]any{
					"da": big.NewInt(91),
					"db": big.NewInt(1692199488),
				},
			},
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			c, err := ParseEvent(tt.signature)
			require.NoError(t, err)
			var topics []types.Hash
			for _, topic := range tt.topics {
				topics = append(topics, types.MustHashFromHex(topic, types.PadNone))
			}
			err = c.DecodeValue(topics, hexutil.MustHexToBytes(tt.data), &tt.arg)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, tt.arg)
			}
		})
	}
}

func TestEvent_DecodeValues(t *testing.T) {
	tests := []struct {
		signature string
		args      []any
		topics    []string
		data      string
		expected  []any
		wantErr   bool
	}{
		{
			signature: "foo(uint256)",
			args:      []any{&big.Int{}},
			topics:    []string{"0x2fbebd3821c4e005fbe0a9002cc1bd25dc266d788dba1dbcb39cc66a07e7b38b"},
			data:      "0000000000000000000000000000000000000000000000000000000000000001",
			expected:  []any{big.NewInt(1)},
		},

		{
			signature: "foo((bytes a))",
			args:      []any{map[string]any{}},
			topics:    []string{"0x7a699d0514ec3b3aad6ef3992fd4993cacccf6906a6b41200cbd7c24d4dde537"},
			data:      "000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000004deadbeef00000000000000000000000000000000000000000000000000000000",
			expected:  []any{map[string]any{"a": []byte{0xde, 0xad, 0xbe, 0xef}}},
		},
		{
			signature: "foo(address indexed a, address indexed b, (bytes32 ca, address cb, bytes cc) c, (uint128 da, uint32 db) d)",
			args: []any{
				&types.Address{},
				&types.Address{},
				map[string]any{
					"ca": &types.Hash{},
					"cb": &types.Address{},
					"cc": []byte{},
				},
				map[string]any{
					"da": &big.Int{},
					"db": &big.Int{},
				},
			},
			topics: []string{
				"0xe2f95411c22cc63e49510640786b88092a37c3926a08ba565b77c8d0fa08bb27",
				"0x0000000000000000000000001F7acDa376eF37EC371235a094113dF9Cb4EfEe1",
				"0x00000000000000000000000068E527780872cda0216Ba0d8fBD58b67a5D5e351",
			},
			data: "0000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000005b0000000000000000000000000000000000000000000000000000000064dcea4063cff3f05ab6a55e7e49095371098bf6455d27c6ed5b6e3c3178c661a821f7290000000000000000000000005b1742020856ad976469f498a296465d7803a7d100000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000014140310010a12090d020c11130e060b040507080f000000000000000000000000",
			expected: []any{
				types.MustAddressFromHexPtr("0x1F7acDa376eF37EC371235a094113dF9Cb4EfEe1"),
				types.MustAddressFromHexPtr("0x68E527780872cda0216Ba0d8fBD58b67a5D5e351"),
				map[string]any{
					"ca": types.MustHashFromHexPtr("0x63cff3f05ab6a55e7e49095371098bf6455d27c6ed5b6e3c3178c661a821f729", types.PadNone),
					"cb": types.MustAddressFromHexPtr("0x5b1742020856ad976469f498a296465d7803a7d1"),
					"cc": hexutil.MustHexToBytes("0x140310010a12090d020c11130e060b040507080f"),
				},
				map[string]any{
					"da": big.NewInt(91),
					"db": big.NewInt(1692199488),
				},
			},
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			c, err := ParseEvent(tt.signature)
			require.NoError(t, err)
			var topics []types.Hash
			for _, topic := range tt.topics {
				topics = append(topics, types.MustHashFromHex(topic, types.PadNone))
			}
			err = c.DecodeValues(topics, hexutil.MustHexToBytes(tt.data), tt.args...)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				for i, arg := range tt.args {
					assert.Equal(t, tt.expected[i], arg)
				}
			}
		})
	}
}

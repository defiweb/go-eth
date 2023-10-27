package abi

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/defiweb/go-eth/hexutil"
)

func TestParseMethod(t *testing.T) {
	tests := []struct {
		signature string
		expected  string
		wantErr   bool
	}{
		{signature: "foo((uint256,bytes32)[])(uint256)", expected: "function foo((uint256, bytes32)[]) returns (uint256)"},
		{signature: "foo((uint256 a, bytes32 b)[] c)(uint256 d)", expected: "function foo((uint256 a, bytes32 b)[] c) returns (uint256 d)"},
		{signature: "function foo(tuple(uint256 a, bytes32 b)[] memory c) pure returns (uint256 d)", expected: "function foo((uint256 a, bytes32 b)[] c) pure returns (uint256 d)"},
		{signature: "event foo(uint256)", wantErr: true},
		{signature: "error foo(uint256)", wantErr: true},
		{signature: "constructor(uint256)", wantErr: true},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			m, err := ParseMethod(tt.signature)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, m.String())
			}
		})
	}
}

func TestMethod_EncodeArgs(t *testing.T) {
	tests := []struct {
		signature string
		arg       []any
		expected  string
	}{
		{signature: "foo()", arg: nil, expected: "c2985578"},
		{signature: "foo(uint256)", arg: []any{1}, expected: "2fbebd380000000000000000000000000000000000000000000000000000000000000001"},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			c, err := ParseMethod(tt.signature)
			require.NoError(t, err)
			enc, err := c.EncodeArgs(tt.arg...)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, hex.EncodeToString(enc))
		})
	}
}

func TestMethod_DecodeArg(t *testing.T) {
	tests := []struct {
		signature string
		arg       any
		data      string
		expected  any
		wantErr   bool
	}{
		{signature: "foo(uint256)", arg: map[string]any{}, data: "2fbebd380000000000000000000000000000000000000000000000000000000000000001", expected: map[string]any{"arg0": big.NewInt(1)}},
		{signature: "foo(uint256)", arg: map[string]any{}, data: "aabbccdd0000000000000000000000000000000000000000000000000000000000000001", wantErr: true},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			c, err := ParseMethod(tt.signature)
			require.NoError(t, err)
			err = c.DecodeArg(hexutil.MustHexToBytes(tt.data), &tt.arg)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, tt.arg)
			}
		})
	}
}

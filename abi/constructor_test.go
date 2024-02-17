package abi

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseConstructor(t *testing.T) {
	tests := []struct {
		signature string
		expected  string
		wantErr   bool
	}{
		{signature: "constructor()", expected: "constructor()"},
		{signature: "constructor(uint256)", expected: "constructor(uint256)"},
		{signature: "((uint256, bytes32)[])", expected: "constructor((uint256, bytes32)[])"},
		{signature: "((uint256 a,bytes32 b)[] a)", expected: "constructor((uint256 a, bytes32 b)[] a)"},
		{signature: "constructor(tuple(uint256 a, bytes32 b)[] memory c)", expected: "constructor((uint256 a, bytes32 b)[] c)"},
		{signature: "foo(uint256)(uint256)", wantErr: true},
		{signature: "event foo(uint256)", wantErr: true},
		{signature: "error foo(uint256)", wantErr: true},
		{signature: "function foo(uint256)", wantErr: true},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			c, err := ParseConstructor(tt.signature)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, c.String())
			}
		})
	}
}

func TestConstructor_EncodeArgs(t *testing.T) {
	tests := []struct {
		signature string
		arg       []any
		expected  string
	}{
		{signature: "constructor()", arg: nil, expected: "aabb"},
		{signature: "constructor(uint256)", arg: []any{1}, expected: "aabb0000000000000000000000000000000000000000000000000000000000000001"},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			c, err := ParseConstructor(tt.signature)
			require.NoError(t, err)
			enc, err := c.EncodeArgs([]byte{0xAA, 0xBB}, tt.arg...)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, hex.EncodeToString(enc))
		})
	}
}

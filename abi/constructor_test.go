package abi

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/defiweb/go-eth/hexutil"
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

func TestConstructor_Text(t *testing.T) {
	tests := []struct {
		name      string
		signature string
		data      string
		expected  string
	}{
		{
			name:      "named arg",
			signature: "constructor(uint256 amount)",
			data:      "000000000000000000000000000000000000000000000000000000000000012c",
			expected:  "constructor(amount=300)",
		},
		{
			name:      "unnamed arg",
			signature: "constructor(uint256)",
			data:      "000000000000000000000000000000000000000000000000000000000000012c",
			expected:  "constructor(arg0=300)",
		},
		{
			name:      "multiple named args",
			signature: "constructor(uint256 code, bool flag)",
			data: "0000000000000000000000000000000000000000000000000000000000000042" +
				"0000000000000000000000000000000000000000000000000000000000000001",
			expected: "constructor(code=66, flag=true)",
		},
		{
			name:      "no args",
			signature: "constructor()",
			data:      "",
			expected:  "constructor()",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := MustParseConstructor(tt.signature)
			assert.Equal(t, tt.expected, c.Text(hexutil.MustHexToBytes(tt.data)))
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

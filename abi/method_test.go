package abi

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseMethod(t *testing.T) {
	tests := []struct {
		signature string
		expected  string
		wantErr   bool
	}{
		{signature: "foo((uint256,bytes32)[])(uint256)", expected: "function foo((uint256, bytes32)[]) returns (uint256)"},
		{signature: "foo((uint256 a, bytes32 b)[] c)(uint256 d)", expected: "function foo((uint256 a, bytes32 b)[] c) returns (uint256 d)"},
		{signature: "function foo(tuple(uint256 a, bytes32 b)[] memory c) pure returns (uint256 d)", expected: "function foo((uint256 a, bytes32 b)[] c) returns (uint256 d)"},
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

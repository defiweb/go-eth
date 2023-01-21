package abi

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

package abi

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/defiweb/go-eth/hexutil"
)

type mockError struct {
	data any
}

func (m *mockError) Error() string {
	return "mock error"
}

func (m *mockError) RPCErrorData() any {
	return m.data
}

func TestParseError(t *testing.T) {
	tests := []struct {
		signature string
		expected  string
		wantErr   bool
	}{
		{signature: "foo((uint256,bytes32)[])", expected: "error foo((uint256, bytes32)[])"},
		{signature: "foo((uint256 a, bytes32 b)[] c)", expected: "error foo((uint256 a, bytes32 b)[] c)"},
		{signature: "error foo(tuple(uint256 a, bytes32 b)[] c)", expected: "error foo((uint256 a, bytes32 b)[] c)"},
		{signature: "foo(uint256)(uint256)", wantErr: true},
		{signature: "event foo(uint256)", wantErr: true},
		{signature: "function foo(uint256)", wantErr: true},
		{signature: "constructor(uint256)", wantErr: true},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			e, err := ParseError(tt.signature)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, e.String())
			}
		})
	}
}

func TestError_Is(t *testing.T) {
	e, err := ParseError("error foo(uint256)")
	require.NoError(t, err)

	assert.True(t, e.Is(hexutil.MustHexToBytes("0x2fbebd38000000000000000000000000000000000000000000000000000000000000012c")))
	assert.False(t, e.Is(hexutil.MustHexToBytes("0xaabbccdd000000000000000000000000000000000000000000000000000000000000012c")))
}

func TestError_ToError(t *testing.T) {
	e, err := ParseError("error foo(uint256)")
	require.NoError(t, err)

	// Custom error
	t.Run("custom error", func(t *testing.T) {
		customErr := e.ToError(hexutil.MustHexToBytes("0x2fbebd38000000000000000000000000000000000000000000000000000000000000012c"))
		require.NotNil(t, customErr)
		assert.Equal(t, "error: foo", customErr.Error())
	})

	// Unknown error
	t.Run("unknown error", func(t *testing.T) {
		unkErr := e.ToError(hexutil.MustHexToBytes("0x112233440000000000000000000000000000000000000000000000000000000000000000"))
		require.Nil(t, unkErr)
	})
}

func TestError_HandleError(t *testing.T) {
	e, err := ParseError("error foo(uint256)")
	require.NoError(t, err)

	// Custom error
	t.Run("custom error", func(t *testing.T) {
		callErr := &mockError{data: hexutil.MustHexToBytes("0x2fbebd38000000000000000000000000000000000000000000000000000000000000012c")}
		customErr := e.HandleError(callErr)
		require.NotNil(t, customErr)
		assert.Equal(t, "error: foo", customErr.Error())
	})

	// Unknown error
	t.Run("unknown error", func(t *testing.T) {
		callErr := &mockError{data: hexutil.MustHexToBytes("0x112233440000000000000000000000000000000000000000000000000000000000000000")}
		unkErr := e.HandleError(callErr)
		require.NotNil(t, unkErr)
		assert.Equal(t, callErr, unkErr)
	})

	// Nil
	t.Run("nil", func(t *testing.T) {
		require.Nil(t, e.HandleError(nil))
	})

	// Not a byte slice
	t.Run("not a byte slice", func(t *testing.T) {
		callErr := &mockError{data: "not a byte slice"}
		require.Equal(t, callErr, e.HandleError(callErr))
	})

	// Not a RPC call error
	t.Run("not a RPC call error", func(t *testing.T) {
		require.Equal(t, errors.New("not a RPC call error"), e.HandleError(errors.New("not a RPC call error")))
	})
}

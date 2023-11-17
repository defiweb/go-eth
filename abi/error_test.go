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

	customErr := e.ToError(hexutil.MustHexToBytes("0x2fbebd38000000000000000000000000000000000000000000000000000000000000012c"))
	require.NotNil(t, customErr)
	assert.Equal(t, "error: foo", customErr.Error())
}

func TestHandleError(t *testing.T) {
	t.Run("ReturnsNil_WhenErrorIsNil", func(t *testing.T) {
		e := NewError("foo", nil)
		result := e.HandleError(nil)
		assert.Nil(t, result)
	})
	t.Run("ReturnsOriginalError_WhenErrorDoesNotImplementRPCErrorData", func(t *testing.T) {
		e := NewError("foo", nil)
		originalErr := errors.New("original error")
		result := e.HandleError(originalErr)
		assert.Equal(t, originalErr, result)
	})
	t.Run("ReturnsOriginalError_WhenRPCErrorDataIsNotByteSlice", func(t *testing.T) {
		e := NewError("foo", nil)
		originalErr := &mockError{data: "not a byte slice"}
		result := e.HandleError(originalErr)
		assert.Equal(t, originalErr, result)
	})
	t.Run("ReturnsCustomError_WhenRPCErrorDataIsByteSlice", func(t *testing.T) {
		e := NewError("foo", nil)
		println(e.fourBytes.Hex())
		originalErr := &mockError{data: hexutil.MustHexToBytes("0xc2985578")}
		result := e.HandleError(originalErr)
		assert.IsType(t, CustomError{}, result)
	})
}

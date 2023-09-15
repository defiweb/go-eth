package transport

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/defiweb/go-eth/hexutil"
)

func TestNewRPCError(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		message  string
		data     any
		expected *RPCError
	}{
		{
			name:    "error with non-hex data",
			code:    ErrCodeGeneral,
			message: "Unauthorized access",
			data:    "some data",
			expected: &RPCError{
				Code:    ErrCodeGeneral,
				Message: "Unauthorized access",
				Data:    "some data",
			},
		},
		{
			name:    "error with hex data",
			code:    ErrCodeGeneral,
			message: "Invalid request",
			data:    "0x68656c6c6f",
			expected: &RPCError{
				Code:    ErrCodeGeneral,
				Message: "Invalid request",
				Data:    hexutil.MustHexToBytes("0x68656c6c6f"),
			},
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := NewRPCError(tt.code, tt.message, tt.data)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

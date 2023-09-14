package hexutil

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBigIntToHex(t *testing.T) {
	tests := []struct {
		name     string
		input    *big.Int
		expected string
	}{
		{"nil input", nil, "0x0"},
		{"zero value", big.NewInt(0), "0x0"},
		{"positive value", big.NewInt(26), "0x1a"},
		{"negative value", big.NewInt(-26), "-0x1a"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, BigIntToHex(tt.input))
		})
	}
}

func TestHexToBigInt(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *big.Int
		err      error
	}{
		{"zero", "0x0", big.NewInt(0), nil},
		{"zero without prefix", "0", big.NewInt(0), nil},
		{"valid positive hex", "0x1a", big.NewInt(26), nil},
		{"valid positive hex without prefix", "1a", big.NewInt(26), nil},
		{"valid negative hex", "-0x1a", big.NewInt(-26), nil},
		{"valid negative hex without prefix", "-1a", big.NewInt(-26), nil},
		{"valid positive single char hex", "0xa", big.NewInt(10), nil},
		{"valid negative single char hex", "-0xa", big.NewInt(-10), nil},
		{"empty string", "", nil, fmt.Errorf("invalid hex string")},
		{"invalid hex", "0x1g", nil, fmt.Errorf("invalid hex string")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := HexToBigInt(tt.input)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBytesToHex(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{"empty bytes", []byte{}, "0x"},
		{"non-empty bytes", []byte("abc"), "0x616263"},
		{"bytes with zeros", []byte{0, 1, 2}, "0x000102"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, BytesToHex(tt.input))
		})
	}
}

func TestHexToBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []byte
		err      error
	}{
		{"empty string", "", []byte{}, nil},
		{"empty data", "0x", []byte{}, nil},
		{"valid hex", "0x616263", []byte("abc"), nil},
		{"valid hex without prefix", "616263", []byte("abc"), nil},
		{"single zero", "0", []byte{0}, nil},
		{"invalid hex", "0x1", nil, fmt.Errorf("invalid hex string, length must be even")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := HexToBytes(tt.input)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

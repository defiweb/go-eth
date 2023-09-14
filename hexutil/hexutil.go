package hexutil

import (
	"encoding/hex"
	"fmt"
	"math/big"
)

// BigIntToHex returns the hex representation of the given big integer.
// The hex string is prefixed with "0x". Negative numbers are prefixed with
// "-0x".
func BigIntToHex(x *big.Int) string {
	if x == nil {
		return "0x0"
	}
	sign := x.Sign()
	switch {
	case sign == 0:
		return "0x0"
	case sign > 0:
		return "0x" + x.Text(16)
	default:
		return "-0x" + x.Text(16)[1:]
	}
}

// HexToBigInt returns the big integer representation of the given hex string.
// The hex string may be prefixed with "0x".
func HexToBigInt(h string) (*big.Int, error) {
	isNeg := len(h) > 1 && h[0] == '-'
	if isNeg {
		h = h[1:]
	}
	if has0xPrefix(h) {
		h = h[2:]
	}
	x, ok := new(big.Int).SetString(h, 16)
	if !ok {
		return nil, fmt.Errorf("invalid hex string")
	}
	if isNeg {
		x.Neg(x)
	}
	return x, nil
}

func MustHexToBigInt(h string) *big.Int {
	x, err := HexToBigInt(h)
	if err != nil {
		panic(err)
	}
	return x
}

// BytesToHex returns the hex representation of the given bytes. The hex string
// is always even-length and prefixed with "0x".
func BytesToHex(b []byte) string {
	r := make([]byte, len(b)*2+2)
	copy(r, `0x`)
	hex.Encode(r[2:], b)
	return string(r)
}

// HexToBytes returns the bytes representation of the given hex string.
// The number of hex digits must be even. The hex string may be prefixed with
// "0x".
func HexToBytes(h string) ([]byte, error) {
	if len(h) == 0 {
		return []byte{}, nil
	}
	if has0xPrefix(h) {
		h = h[2:]
	}
	if len(h) == 1 && h[0] == '0' {
		return []byte{0}, nil
	}
	if len(h) == 0 {
		return []byte{}, nil
	}
	if len(h)%2 != 0 {
		return nil, fmt.Errorf("invalid hex string, length must be even")
	}
	return hex.DecodeString(h)
}

func MustHexToBytes(h string) []byte {
	b, err := HexToBytes(h)
	if err != nil {
		panic(err)
	}
	return b
}

// has0xPrefix returns true if the given byte slice starts with "0x".
func has0xPrefix(h string) bool {
	return len(h) >= 2 && h[0] == '0' && (h[1] == 'x' || h[1] == 'X')
}

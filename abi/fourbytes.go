package abi

import "github.com/defiweb/go-eth/hexutil"

// FourBytes is a 4-byte method selector.
type FourBytes [4]byte

// Bytes returns the four bytes as a byte slice.
func (f FourBytes) Bytes() []byte {
	return f[:]
}

// Hex returns the four bytes as a hex string.
func (f FourBytes) Hex() string {
	return hexutil.BytesToHex(f[:])
}

// String returns the four bytes as a hex string.
func (f FourBytes) String() string {
	return f.Hex()
}

// Match returns true if the given jsonABI data matches the four byte selector.
func (f FourBytes) Match(data []byte) bool {
	if len(data) < 4 {
		return false
	}
	return f == FourBytes{data[0], data[1], data[2], data[3]}
}

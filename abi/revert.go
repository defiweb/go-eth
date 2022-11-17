package abi

import (
	"bytes"
	"fmt"
)

// revertPrefix is the prefix of revert messages. It is the first 4 bytes of the
// keccak256 hash of the string "Error(string)".
var revertPrefix = []byte{0x08, 0xc3, 0x79, 0xa0}

// IsRevert returns true if the data is a revert message. It does not check
// whether the data is a valid revert message, hence DecodeRevert may still
// return an error.
func IsRevert(data []byte) bool {
	return len(data) >= 4 && bytes.Equal(data[:4], revertPrefix)
}

// DecodeRevert decodes the revert data returned by contract calls.
func DecodeRevert(data []byte) (string, error) {
	if len(data) < 4 {
		return "", fmt.Errorf("abi: invalid data length %d", len(data))
	}
	if !bytes.Equal(data[:4], revertPrefix) {
		return "", fmt.Errorf("abi: invalid revert prefix %x", data[:4])
	}
	t := NewTupleOfElements(NewString()) // Equivalent to NewType("(string)").
	if _, err := t.DecodeABI(BytesToWords(data[4:])); err != nil {
		return "", err
	}
	return t.Get(0).(*StringValue).String(), nil
}

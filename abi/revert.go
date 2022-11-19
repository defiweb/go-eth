package abi

import (
	"fmt"
)

// Revert is the Error instance for revert responses.
var Revert = NewError("Error", NewTupleType(TupleTypeElem{Name: "error", Type: NewStringType()}))

// revertPrefix is the prefix of revert messages. It is the first 4 bytes of the
// keccak256 hash of the string "Error(string)".
var revertPrefix = FourBytes{0x08, 0xc3, 0x79, 0xa0}

// IsRevert returns true if the data is a revert message. It does not check
// whether the data is a valid revert message, hence DecodeRevert may still
// return an error.
func IsRevert(data []byte) bool {
	return revertPrefix.Match(data)
}

// DecodeRevert decodes the revert data returned by contract calls.
func DecodeRevert(data []byte) (string, error) {
	// The code below is a slightly optimized version of
	// Revert.DecodeValues(data).
	if !revertPrefix.Match(data) {
		return "", fmt.Errorf("abi: invalid revert prefix")
	}
	s := NewStringValue()
	t := NewTupleValue(TupleValueElem{Value: s})
	if _, err := t.DecodeABI(BytesToWords(data[4:])); err != nil {
		return "", err
	}
	return s.String(), nil
}

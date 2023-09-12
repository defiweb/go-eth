package abi

import "fmt"

// Revert is the Error instance for revert responses.
var Revert = NewError("Error", NewTupleType(TupleTypeElem{Name: "error", Type: NewStringType()}))

// revertPrefix is the prefix of revert messages. It is the first 4 bytes of the
// keccak256 hash of the string "Error(string)".
var revertPrefix = FourBytes{0x08, 0xc3, 0x79, 0xa0}

// RevertError represents an error returned by contract calls when the call
// reverts.
type RevertError struct {
	Reason string
}

// Error implements the error interface.
func (e RevertError) Error() string {
	return fmt.Sprintf("revert: %s", e.Reason)
}

// IsRevert returns true if the data has the revert prefix.
func IsRevert(data []byte) bool {
	return revertPrefix.Match(data) && (len(data)-4)%WordLength == 0
}

// DecodeRevert decodes the revert data returned by contract calls.
// If the data is not a valid revert message, it returns an empty string.
func DecodeRevert(data []byte) string {
	// The code below is a slightly optimized version of
	// Revert.DecodeValues(data).
	if !IsRevert(data) {
		return ""
	}
	s := new(StringValue)
	t := TupleValue{TupleValueElem{Value: s}}
	if _, err := t.DecodeABI(BytesToWords(data[4:])); err != nil {
		return ""
	}
	return string(*s)
}

// ToRevertError converts the revert data returned by contract calls into a RevertError.
// If the data does not contain a valid revert message, it returns nil.
func ToRevertError(data []byte) error {
	if !IsRevert(data) {
		return nil
	}
	return RevertError{Reason: DecodeRevert(data)}
}

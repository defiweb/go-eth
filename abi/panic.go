package abi

import (
	"fmt"
	"math/big"
)

// Panic is the Error instance for panic responses.
var Panic = NewError("Panic", NewTupleType(TupleTypeElem{Name: "error", Type: NewUintType(256)}))

// panicPrefix is the prefix of panic messages. It is the first 4 bytes of the
// keccak256 hash of the string "Panic(uint256)".
var panicPrefix = FourBytes{0x4e, 0x48, 0x7b, 0x71}

// PanicError represents an error returned by contract calls when the call
// panics.
type PanicError struct {
	Code *big.Int
}

// Error implements the error interface.
func (e PanicError) Error() string {
	return fmt.Sprintf("panic: %s", e.Code.String())
}

// IsPanic returns true if the data has the panic prefix.
func IsPanic(data []byte) bool {
	return panicPrefix.Match(data) && len(data) == 36
}

// DecodePanic decodes the panic data returned by contract calls.
// If the data is not a valid panic message, it returns nil.
func DecodePanic(data []byte) *big.Int {
	// The code below is a slightly optimized version of
	// Panic.DecodeValues(data).
	if !IsPanic(data) {
		return nil
	}
	s := &UintValue{Size: 256}
	t := TupleValue{TupleValueElem{Value: s}}
	if _, err := t.DecodeABI(BytesToWords(data[4:])); err != nil {
		return nil
	}
	return &s.Int
}

// ToPanicError converts the panic data returned by contract calls into a PanicError.
// If the data does not contain a valid panic message, it returns nil.
func ToPanicError(data []byte) error {
	if !IsPanic(data) {
		return nil
	}
	return PanicError{Code: DecodePanic(data)}
}

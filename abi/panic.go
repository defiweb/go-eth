package abi

import (
	"math/big"
)

// Panic is the Error instance for panic responses.
var Panic = NewError("Panic", NewTupleType(TupleTypeElem{Name: "error", Type: NewUintType(256)}))

// panicPrefix is the prefix of panic messages. It is the first 4 bytes of the
// keccak256 hash of the string "Panic(uint256)".
var panicPrefix = FourBytes{0x4e, 0x48, 0x7b, 0x71}

// IsPanic returns true if the data has the panic prefix. It does not check
// whether the data is a valid panic message.
func IsPanic(data []byte) bool {
	return panicPrefix.Match(data)
}

// DecodePanic decodes the panic data returned by contract calls.
// If the data is not a valid panic message, it returns nil.
func DecodePanic(data []byte) *big.Int {
	// The code below is a slightly optimized version of
	// Panic.DecodeValues(data).
	if !panicPrefix.Match(data) {
		return nil
	}
	s := &UintValue{Size: 256}
	t := TupleValue{TupleValueElem{Value: s}}
	if _, err := t.DecodeABI(BytesToWords(data[4:])); err != nil {
		return nil
	}
	return &s.Int
}

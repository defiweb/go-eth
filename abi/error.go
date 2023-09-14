package abi

import (
	"fmt"

	"github.com/defiweb/go-eth/crypto"
)

// CustomError represents an custom error returned by a contract call.
type CustomError struct {
	Type *Error // The error type.
	Data []byte // The error data returned by the contract call.
}

// Error implements the error interface.
func (e CustomError) Error() string {
	return fmt.Sprintf("error: %s", e.Type.Name())
}

// Error represents an error in an ABI. The error can be used to decode errors
// returned by a contract call.
type Error struct {
	name   string
	inputs *TupleType
	abi    *ABI

	fourBytes FourBytes
	signature string
}

// NewError creates a new Error instance.
func NewError(name string, inputs *TupleType) *Error {
	return Default.NewError(name, inputs)
}

// ParseError parses an error signature and returns a new Error.
//
// An error signature is similar to a method signature, but returns no values.
// It can be optionally prefixed with the "error" keyword.
//
// The following examples are valid signatures:
//
//	foo((uint256,bytes32)[])
//	foo((uint256 a, bytes32 b)[] c)
//	error foo(tuple(uint256 a, bytes32 b)[] c)
//
// This function is equivalent to calling Parser.ParseError with the default
// configuration.
func ParseError(signature string) (*Error, error) {
	return Default.ParseError(signature)
}

// MustParseError is like ParseError but panics on error.
func MustParseError(signature string) *Error {
	e, err := ParseError(signature)
	if err != nil {
		panic(err)
	}
	return e
}

// NewError creates a new Error instance.
func (a *ABI) NewError(name string, inputs *TupleType) *Error {
	m := &Error{
		name:   name,
		inputs: inputs,
		abi:    a,
	}
	m.generateSignature()
	m.calculateFourBytes()
	return m
}

// ParseError parses an error signature and returns a new Error.
//
// See ParseError for more information.
func (a *ABI) ParseError(signature string) (*Error, error) {
	return parseError(a, signature)
}

// Name returns the name of the error.
func (m *Error) Name() string {
	return m.name
}

// Inputs returns the input arguments of the error as a tuple type.
func (m *Error) Inputs() *TupleType {
	return m.inputs
}

// FourBytes is the first four bytes of the Keccak256 hash of the error
// signature.
func (m *Error) FourBytes() FourBytes {
	return m.fourBytes
}

// Signature returns the error signature, that is, the error name and the
// canonical type of error arguments.
func (m *Error) Signature() string {
	return m.signature
}

// Is returns true if the ABI encoded data is an error of this type.
func (m *Error) Is(data []byte) bool {
	return m.fourBytes.Match(data) && (len(data)-4)%WordLength == 0
}

// DecodeValue decodes the error into a map or structure. If a structure is
// given, it must have fields with the same names as error arguments.
func (m *Error) DecodeValue(data []byte, val any) error {
	if m.fourBytes.Match(data) {
		return fmt.Errorf("abi: selector mismatch for error %s", m.name)
	}
	return m.abi.DecodeValue(m.inputs, data[4:], val)
}

// DecodeValues decodes the error into a map or structure. If a structure is
// given, it must have fields with the same names as error arguments.
func (m *Error) DecodeValues(data []byte, vals ...any) error {
	if m.fourBytes.Match(data) {
		return fmt.Errorf("abi: selector mismatch for error %s", m.name)
	}
	return m.abi.DecodeValues(m.inputs, data[4:], vals...)
}

// ToError converts the error data returned by contract calls into a CustomError.
// If the data does not contain a valid error message, it returns nil.
func (m *Error) ToError(data []byte) error {
	if !m.fourBytes.Match(data) {
		return nil
	}
	return CustomError{
		Type: m,
		Data: data[4:],
	}
}

// String returns the human-readable signature of the error.
func (m *Error) String() string {
	return "error " + m.name + m.inputs.String()
}

func (m *Error) generateSignature() {
	m.signature = m.name + m.inputs.CanonicalType()
}

func (m *Error) calculateFourBytes() {
	id := crypto.Keccak256([]byte(m.Signature()))
	copy(m.fourBytes[:], id[:4])
}

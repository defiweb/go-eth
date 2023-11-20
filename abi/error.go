package abi

import (
	"errors"
	"fmt"

	"github.com/defiweb/go-eth/crypto"
)

// CustomError represents a custom error returned by a contract call.
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
	if inputs == nil {
		inputs = NewTupleType()
	}
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
	return Default.MustParseError(signature)
}

// NewError creates a new Error instance.
//
// This method is rarely used, see ParseError for a more convenient way to
// create a new Error.
func (a *ABI) NewError(name string, inputs *TupleType) *Error {
	e := &Error{
		name:   name,
		inputs: inputs,
		abi:    a,
	}
	e.generateSignature()
	e.calculateFourBytes()
	return e
}

// ParseError parses an error signature and returns a new Error.
//
// See ParseError for more information.
func (a *ABI) ParseError(signature string) (*Error, error) {
	return parseError(a, nil, signature)
}

// MustParseError is like ParseError but panics on error.
func (a *ABI) MustParseError(signature string) *Error {
	m, err := a.ParseError(signature)
	if err != nil {
		panic(err)
	}
	return m
}

// Name returns the name of the error.
func (e *Error) Name() string {
	return e.name
}

// Inputs returns the input arguments of the error as a tuple type.
func (e *Error) Inputs() *TupleType {
	return e.inputs
}

// FourBytes is the first four bytes of the Keccak256 hash of the error
// signature.
func (e *Error) FourBytes() FourBytes {
	return e.fourBytes
}

// Signature returns the error signature, that is, the error name and the
// canonical type of error arguments.
func (e *Error) Signature() string {
	return e.signature
}

// Is returns true if the ABI encoded data is an error of this type.
func (e *Error) Is(data []byte) bool {
	return e.fourBytes.Match(data) && (len(data)-4)%WordLength == 0
}

// DecodeValue decodes the error into a map or structure. If a structure is
// given, it must have fields with the same names as error arguments.
func (e *Error) DecodeValue(data []byte, val any) error {
	if e.fourBytes.Match(data) {
		return fmt.Errorf("abi: selector mismatch for error %s", e.name)
	}
	return e.abi.DecodeValue(e.inputs, data[4:], val)
}

// MustDecodeValue is like DecodeValue but panics on error.
func (e *Error) MustDecodeValue(data []byte, val any) {
	err := e.DecodeValue(data, val)
	if err != nil {
		panic(err)
	}
}

// DecodeValues decodes the error into a map or structure. If a structure is
// given, it must have fields with the same names as error arguments.
func (e *Error) DecodeValues(data []byte, vals ...any) error {
	if e.fourBytes.Match(data) {
		return fmt.Errorf("abi: selector mismatch for error %s", e.name)
	}
	return e.abi.DecodeValues(e.inputs, data[4:], vals...)
}

// MustDecodeValues is like DecodeValues but panics on error.
func (e *Error) MustDecodeValues(data []byte, vals ...any) {
	err := e.DecodeValues(data, vals...)
	if err != nil {
		panic(err)
	}
}

// ToError converts the error data returned by contract calls into a CustomError.
// If the data does not contain a valid error message, it returns nil.
func (e *Error) ToError(data []byte) error {
	if !e.fourBytes.Match(data) {
		return nil
	}
	return CustomError{
		Type: e,
		Data: data[4:],
	}
}

// HandleError converts an error returned by a contract call to a custom error
// if possible. If provider error is nil, it returns nil.
func (e *Error) HandleError(err error) error {
	if err == nil {
		return nil
	}
	var dataErr interface{ RPCErrorData() any }
	if !errors.As(err, &dataErr) {
		return err
	}
	data, ok := dataErr.RPCErrorData().([]byte)
	if !ok {
		return err
	}
	if err := e.ToError(data); err != nil {
		return err
	}
	return err
}

// String returns the human-readable signature of the error.
func (e *Error) String() string {
	return "error " + e.name + e.inputs.String()
}

func (e *Error) generateSignature() {
	e.signature = e.name + e.inputs.CanonicalType()
}

func (e *Error) calculateFourBytes() {
	id := crypto.Keccak256([]byte(e.Signature()))
	copy(e.fourBytes[:], id[:4])
}

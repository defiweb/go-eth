package abi

import (
	"errors"
	"fmt"
	"strings"

	"github.com/defiweb/go-eth/crypto"
)

// CustomError represents a custom error returned by a contract call.
type CustomError struct {
	Type *Error // The error type.
	Data []byte // The error data returned by the contract call (including the 4-byte selector).
}

// Values decodes the error data into a map of argument names to values.
// If decoding fails, it returns nil.
func (e CustomError) Values() map[string]any {
	if e.Type == nil || len(e.Data) == 0 {
		return nil
	}
	res := make(map[string]any)
	if err := e.Type.DecodeValue(e.Data, res); err != nil {
		return nil
	}
	return res
}

// Error implements the error interface.
func (e CustomError) Error() string {
	if e.Type == nil {
		return "unknown error"
	}
	return e.Type.Format(e.Data)
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
	if !e.fourBytes.Match(data) {
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
	if !e.fourBytes.Match(data) {
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
		Data: data,
	}
}

// HandleError converts an error returned by a contract call to a custom error
// if possible. If the provided error is nil, it returns nil.
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

// Format returns a human-readable representation of the error, including the
// values of the error arguments.
//
// The data must be the ABI-encoded error data returned by a contract call;
// the 4-byte selector is optional.
func (e *Error) Format(data []byte) string {
	res := make(map[string]any)
	msg := strings.Builder{}
	msg.WriteString("error ")
	msg.WriteString(e.Name())
	if len(data)%32 == 4 {
		if !e.fourBytes.Match(data) {
			msg.WriteString("(selector mismatch)")
			return msg.String()
		}
		data = data[4:]
	}
	if decErr := DecodeValue(e.Inputs(), data, res); decErr != nil {
		msg.WriteString("(")
		msg.WriteString(decErr.Error())
		msg.WriteString(")")
		return msg.String()
	}
	msg.WriteString("(")
	for i, input := range e.Inputs().Elements() {
		if i > 0 {
			msg.WriteString(", ")
		}
		name := input.Name
		if name == "" {
			name = fmt.Sprintf("arg%d", i)
		}
		msg.WriteString(name)
		msg.WriteString("=")
		_, _ = fmt.Fprintf(&msg, "%v", res[name])
	}
	msg.WriteString(")")
	return msg.String()
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

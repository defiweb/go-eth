package abi

import (
	"fmt"

	"github.com/defiweb/go-eth/crypto"
)

// Error represents an error in an jsonABI. The error can be used to decode errors
// returned by a contract call.
type Error struct {
	name   string
	inputs *TupleType
	config *Config

	fourBytes FourBytes
	signature string
}

// NewError creates a new Error instance.
func NewError(name string, inputs *TupleType) *Error {
	return NewErrorWithConfig(name, inputs, DefaultConfig)
}

// NewErrorWithConfig creates a new Error instance with a custom config.
func NewErrorWithConfig(name string, inputs *TupleType, config *Config) *Error {
	m := &Error{
		name:   name,
		inputs: inputs,
		config: config,
	}
	m.generateSignature()
	m.calculateFourBytes()
	return m
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
// canonical type of the error arguments.
func (m *Error) Signature() string {
	return m.signature
}

// Is returns true if the jsonABI encoded data is an error of this type.
func (m *Error) Is(data []byte) bool {
	return m.fourBytes.Match(data)
}

// DecodeValue decodes the error into a map or structure. If a structure is
// given, it must have fields with the same names as the error arguments.
func (m *Error) DecodeValue(data []byte, val any) error {
	if m.fourBytes.Match(data) {
		return fmt.Errorf("abi: selector mismatch for error %s", m.name)
	}
	return NewDecoder(m.config).DecodeValue(m.inputs.Value(), data[4:], val)
}

// DecodeValues decodes the error into a map or structure. If a structure is
// given, it must have fields with the same names as the error arguments.
func (m *Error) DecodeValues(data []byte, vals ...any) error {
	if m.fourBytes.Match(data) {
		return fmt.Errorf("abi: selector mismatch for error %s", m.name)
	}
	return NewDecoder(m.config).DecodeValues(m.inputs.Value().(*TupleValue), data[4:], vals...)
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

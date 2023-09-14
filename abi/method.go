package abi

import (
	"fmt"
	"strings"

	"github.com/defiweb/go-eth/crypto"
)

// Method represents a method in an ABI. The method can be used to encode
// arguments for a method call and decode return values from a method call.
type Method struct {
	name    string
	inputs  *TupleType
	outputs *TupleType
	abi     *ABI

	fourBytes FourBytes
	signature string
}

// NewMethod creates a new Method instance.
func NewMethod(name string, inputs, outputs *TupleType) *Method {
	return Default.NewMethod(name, inputs, outputs)
}

// ParseMethod parses a method signature and returns a new Method.
//
// The method accepts Solidity method signatures, but allows to omit the
// "function" keyword, argument names and the "returns" keyword. Method
// modifiers and argument data location specifiers are allowed, but ignored.
// Tuple types are indicated by parentheses, with the optional keyword "tuple"
// before the parentheses.
//
// The following examples are valid signatures:
//
//	foo((uint256,bytes32)[])(uint256)
//	foo((uint256 a, bytes32 b)[] c)(uint256 d)
//	function foo(tuple(uint256 a, bytes32 b)[] memory c) pure returns (uint256 d)
//
// This function is equivalent to calling Parser.ParseMethod with the default
// configuration.
func ParseMethod(signature string) (*Method, error) {
	return Default.ParseMethod(signature)
}

// MustParseMethod is like ParseMethod but panics on error.
func MustParseMethod(signature string) *Method {
	m, err := ParseMethod(signature)
	if err != nil {
		panic(err)
	}
	return m
}

// NewMethod creates a new Method instance.
func (a *ABI) NewMethod(name string, inputs, outputs *TupleType) *Method {
	m := &Method{
		name:    name,
		inputs:  inputs,
		outputs: outputs,
		abi:     a,
	}
	m.generateSignature()
	m.calculateFourBytes()
	return m
}

// ParseMethod parses a method signature and returns a new Method.
//
// See ParseMethod for more information.
func (a *ABI) ParseMethod(signature string) (*Method, error) {
	return parseMethod(a, signature)
}

// Name returns the name of the method.
func (m *Method) Name() string {
	return m.name
}

// Inputs returns the input arguments of the method as a tuple type.
func (m *Method) Inputs() *TupleType {
	return m.inputs
}

// Outputs returns the output values of the method as a tuple type.
func (m *Method) Outputs() *TupleType {
	return m.outputs
}

// FourBytes is the first four bytes of the Keccak256 hash of the method
// signature. It is also known as a "function selector."
func (m *Method) FourBytes() FourBytes {
	return m.fourBytes
}

// Signature returns the method signature, that is, the method name and the
// canonical types of the input arguments.
func (m *Method) Signature() string {
	return m.signature
}

// EncodeArg encodes arguments for a method call using a provided map or
// structure. The map or structure must have fields with the same names as
// the method arguments.
func (m *Method) EncodeArg(arg any) ([]byte, error) {
	encoded, err := m.abi.EncodeValue(m.inputs, arg)
	if err != nil {
		return nil, err
	}
	return append(m.fourBytes.Bytes(), encoded...), nil
}

// EncodeArgs encodes arguments for a method call.
func (m *Method) EncodeArgs(args ...any) ([]byte, error) {
	encoded, err := m.abi.EncodeValues(m.inputs, args...)
	if err != nil {
		return nil, err
	}
	return append(m.fourBytes.Bytes(), encoded...), nil
}

// DecodeArg decodes ABI-encoded arguments a method call.
func (m *Method) DecodeArg(data []byte, arg any) error {
	if !m.fourBytes.Match(data[:4]) {
		return fmt.Errorf(
			"abi: calldata signature 0x%x do not match method signature %s",
			data[:4],
			m.fourBytes,
		)
	}
	return m.abi.DecodeValue(m.inputs, data[4:], arg)
}

// DecodeArgs decodes ABI-encoded arguments a method call.
func (m *Method) DecodeArgs(data []byte, args ...any) error {
	if !m.fourBytes.Match(data[:4]) {
		return fmt.Errorf(
			"abi: calldata signature 0x%x do not match method signature %s",
			data[:4],
			m.fourBytes,
		)
	}
	return m.abi.DecodeValues(m.inputs, data[4:], args...)
}

// DecodeValue decodes the values returned by a method call into a map or
// structure. If a structure is given, it must have fields with the same names
// as the values returned by the method.
func (m *Method) DecodeValue(data []byte, val any) error {
	return m.abi.DecodeValue(m.outputs, data, val)
}

// DecodeValues decodes return values from a method call to a given values.
func (m *Method) DecodeValues(data []byte, vals ...any) error {
	return m.abi.DecodeValues(m.outputs, data, vals...)
}

// String returns the human-readable signature of the method.
func (m *Method) String() string {
	var buf strings.Builder
	buf.WriteString("function ")
	buf.WriteString(m.name)
	buf.WriteString(m.inputs.String())
	if m.outputs.Size() > 0 {
		buf.WriteString(" returns ")
		buf.WriteString(m.outputs.String())
	}
	return buf.String()
}

func (m *Method) generateSignature() {
	m.signature = m.name + m.inputs.CanonicalType()
}

func (m *Method) calculateFourBytes() {
	id := crypto.Keccak256([]byte(m.Signature()))
	copy(m.fourBytes[:], id[:4])
}

package abi

import (
	"fmt"
	"strings"

	"github.com/defiweb/go-eth/crypto"
)

type StateMutability int

const (
	StateMutabilityUnknown StateMutability = iota
	StateMutabilityPure
	StateMutabilityView
	StateMutabilityNonPayable
	StateMutabilityPayable
)

func StateMutabilityFromString(s string) StateMutability {
	switch strings.ToLower(s) {
	case "pure":
		return StateMutabilityPure
	case "view":
		return StateMutabilityView
	case "nonpayable":
		return StateMutabilityNonPayable
	case "payable":
		return StateMutabilityPayable
	default:
		return StateMutabilityUnknown
	}
}

func (m StateMutability) String() string {
	switch m {
	case StateMutabilityPure:
		return "pure"
	case StateMutabilityView:
		return "view"
	case StateMutabilityNonPayable:
		return "nonpayable"
	case StateMutabilityPayable:
		return "payable"
	default:
		return "unknown"
	}
}

// Method represents a method in an ABI. The method can be used to encode
// arguments for a method call and decode return values from a method call.
type Method struct {
	name            string
	inputs          *TupleType
	outputs         *TupleType
	stateMutability StateMutability
	abi             *ABI

	fourBytes FourBytes
	signature string
}

// NewMethod creates a new Method instance.
//
// This method is rarely used, see ParseMethod for a more convenient way to
// create a new Method.
func NewMethod(name string, inputs, outputs *TupleType, mutability StateMutability) *Method {
	if inputs == nil {
		inputs = NewTupleType()
	}
	if outputs == nil {
		outputs = NewTupleType()
	}
	return Default.NewMethod(name, inputs, outputs, mutability)
}

// ParseMethod parses a method signature and returns a new Method.
//
// The method accepts Solidity method signatures, but allows to omit the
// "function" keyword, argument names and the "returns" keyword. Method
// modifiers and argument data location specifiers are allowed. Tuple types are
// indicated by parentheses, with the optional keyword "tuple" before the
// parentheses.
//
// If argument names are omitted, the default "argN" is used, where N is the
// argument index. Similarly, if return value names are omitted, "argN" is also
// used.
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
	return Default.MustParseMethod(signature)
}

// NewMethod creates a new Method instance.
func (a *ABI) NewMethod(name string, inputs, outputs *TupleType, mutability StateMutability) *Method {
	m := &Method{
		name:            name,
		inputs:          inputs,
		outputs:         outputs,
		stateMutability: mutability,
		abi:             a,
	}
	m.generateSignature()
	m.calculateFourBytes()
	return m
}

// ParseMethod parses a method signature and returns a new Method.
//
// See ParseMethod for more information.
func (a *ABI) ParseMethod(signature string) (*Method, error) {
	return parseMethod(a, nil, signature)
}

// MustParseMethod is like ParseMethod but panics on error.
func (a *ABI) MustParseMethod(signature string) *Method {
	m, err := a.ParseMethod(signature)
	if err != nil {
		panic(err)
	}
	return m
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

// StateMutability returns the state mutability of the method.
func (m *Method) StateMutability() StateMutability {
	return m.stateMutability
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
// structure.
//
// Provided struct or map must have fields that match the names of the method's
// arguments.
//
// The return value is a ABI-encoded data prefixed with the method selector.
func (m *Method) EncodeArg(arg any) ([]byte, error) {
	encoded, err := m.abi.EncodeValue(m.inputs, arg)
	if err != nil {
		return nil, err
	}
	return append(m.fourBytes.Bytes(), encoded...), nil
}

// MustEncodeArg is like EncodeArg but panics on error.
func (m *Method) MustEncodeArg(arg any) []byte {
	encoded, err := m.EncodeArg(arg)
	if err != nil {
		panic(err)
	}
	return encoded
}

// EncodeArgs encodes arguments for a method call using a provided list of
// arguments.
//
// The return value is a ABI-encoded data prefixed with the method selector.
func (m *Method) EncodeArgs(args ...any) ([]byte, error) {
	encoded, err := m.abi.EncodeValues(m.inputs, args...)
	if err != nil {
		return nil, err
	}
	return append(m.fourBytes.Bytes(), encoded...), nil
}

// MustEncodeArgs is like EncodeArgs but panics on error.
func (m *Method) MustEncodeArgs(args ...any) []byte {
	encoded, err := m.EncodeArgs(args...)
	if err != nil {
		panic(err)
	}
	return encoded
}

// DecodeArg decodes an ABI-encoded data into a provided map or struct.
//
// Provided struct or map must have fields that match the names of the method's
// arguments.
//
// Provided data must be prefixed with the method selector.
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

// MustDecodeArg is like DecodeArg but panics on error.
func (m *Method) MustDecodeArg(data []byte, arg any) {
	if err := m.DecodeArg(data, arg); err != nil {
		panic(err)
	}
}

// DecodeArgs decodes an ABI-encoded data into a provided list of arguments.
//
// Provided data must be prefixed with the method selector.
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

// MustDecodeArgs is like DecodeArgs but panics on error.
func (m *Method) MustDecodeArgs(data []byte, args ...any) {
	if err := m.DecodeArgs(data, args...); err != nil {
		panic(err)
	}
}

// DecodeValue decodes an ABI-encoded data into a provided map or struct.
//
// Provided struct or map must have fields that match the names of the method's
// return values.
func (m *Method) DecodeValue(data []byte, val any) error {
	return m.abi.DecodeValue(m.outputs, data, val)
}

// MustDecodeValue is like DecodeValue but panics on error.
func (m *Method) MustDecodeValue(data []byte, val any) {
	if err := m.DecodeValue(data, val); err != nil {
		panic(err)
	}
}

// DecodeValues decodes an ABI-encoded data into a provided list of return
// variables.
func (m *Method) DecodeValues(data []byte, vals ...any) error {
	return m.abi.DecodeValues(m.outputs, data, vals...)
}

// MustDecodeValues is like DecodeValues but panics on error.
func (m *Method) MustDecodeValues(data []byte, vals ...any) {
	if err := m.DecodeValues(data, vals...); err != nil {
		panic(err)
	}
}

// String returns the human-readable signature of the method.
func (m *Method) String() string {
	var buf strings.Builder
	buf.WriteString("function ")
	buf.WriteString(m.name)
	buf.WriteString(m.inputs.String())
	if m.stateMutability != StateMutabilityUnknown {
		buf.WriteString(" ")
		buf.WriteString(m.stateMutability.String())
	}
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

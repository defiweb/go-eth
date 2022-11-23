package abi

import (
	"strings"

	"github.com/defiweb/go-eth/crypto"
)

// Method represents a method in an ABI. The method can be used to encode
// arguments for a method call and decode return values from a method call.
type Method struct {
	name    string
	inputs  *TupleType
	outputs *TupleType
	config  *Config

	fourBytes FourBytes
	signature string
}

// NewMethod creates a new Method instance.
func NewMethod(name string, inputs, outputs *TupleType) *Method {
	return NewMethodWithConfig(name, inputs, outputs, DefaultConfig)
}

// NewMethodWithConfig creates a new Method instance with a custom config.
func NewMethodWithConfig(name string, inputs, outputs *TupleType, config *Config) *Method {
	m := &Method{
		name:    name,
		inputs:  inputs,
		outputs: outputs,
		config:  config,
	}
	m.generateSignature()
	m.calculateFourBytes()
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
func (m *Method) EncodeArg(val any) ([]byte, error) {
	encoded, err := NewEncoder(m.config).EncodeValue(m.inputs.Value(), val)
	if err != nil {
		return nil, err
	}
	return append(m.fourBytes.Bytes(), encoded...), nil
}

// EncodeArgs encodes arguments for a method call.
func (m *Method) EncodeArgs(args ...any) ([]byte, error) {
	encoded, err := NewEncoder(m.config).EncodeValues(m.inputs.Value().(*TupleValue), args...)
	if err != nil {
		return nil, err
	}
	return append(m.fourBytes.Bytes(), encoded...), nil
}

// DecodeValue decodes the values returned by a method call into a map or
// structure. If a structure is given, it must have fields with the same names
// as the values returned by the method.
func (m *Method) DecodeValue(data []byte, val any) error {
	return NewDecoder(m.config).DecodeValue(m.outputs.Value(), data, val)
}

// DecodeValues decodes return values from a method call to a given values.
func (m *Method) DecodeValues(data []byte, vals ...any) error {
	return NewDecoder(m.config).DecodeValues(m.outputs.Value().(*TupleValue), data, vals...)
}

// String returns the human-readable signature of the method.
func (m *Method) String() string {
	var buf strings.Builder
	buf.WriteString("function ")
	buf.WriteString(m.name)
	buf.WriteString(m.inputs.Type())
	if m.outputs.Size() > 0 {
		buf.WriteString(" returns ")
		buf.WriteString(m.outputs.Type())
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

package abi

import (
	"strings"

	"github.com/defiweb/go-eth/crypto"
	"github.com/defiweb/go-eth/hexutil"
)

type FourBytes [4]byte

type Method struct {
	name    string
	inputs  *TupleType
	outputs *TupleType
	config  *Config

	fourBytes FourBytes
	signature string
}

func NewMethod(name string, inputs, outputs *TupleType) *Method {
	return NewMethodWithConfig(name, inputs, outputs, DefaultConfig)
}

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

func (m *Method) Name() string {
	return m.name
}

func (m *Method) Inputs() *TupleType {
	return m.inputs
}

func (m *Method) Outputs() *TupleType {
	return m.outputs
}

func (m *Method) FourBytes() FourBytes {
	return m.fourBytes
}

func (m *Method) Signature() string {
	return m.signature
}

func (m *Method) EncodeArg(val any) ([]byte, error) {
	encoded, err := NewEncoder(m.config).EncodeValue(m.inputs.New(), val)
	if err != nil {
		return nil, err
	}
	return append(m.fourBytes.Bytes(), encoded...), nil
}

func (m *Method) EncodeArgs(args ...any) ([]byte, error) {
	encoded, err := NewEncoder(m.config).EncodeValues(m.inputs.New().(*TupleValue), args...)
	if err != nil {
		return nil, err
	}
	return append(m.fourBytes.Bytes(), encoded...), nil
}

func (m *Method) DecodeValue(data []byte, val any) error {
	return NewDecoder(m.config).DecodeValue(m.outputs.New(), data, val)
}

func (m *Method) DecodeValues(data []byte, vals ...any) error {
	return NewDecoder(m.config).DecodeValues(m.outputs.New().(*TupleValue), data, vals...)
}

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

func (f FourBytes) Bytes() []byte {
	return f[:]
}

func (f FourBytes) Hex() string {
	return hexutil.BytesToHex(f[:])
}

func (f FourBytes) String() string {
	return f.Hex()
}

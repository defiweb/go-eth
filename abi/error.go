package abi

import (
	"bytes"
	"fmt"

	"github.com/defiweb/go-eth/crypto"
)

type Error struct {
	name   string
	inputs *TupleType
	config *Config

	fourBytes FourBytes
	signature string
}

func NewError(name string, inputs *TupleType) *Error {
	return NewErrorWithConfig(name, inputs, DefaultConfig)
}

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

func (m *Error) Name() string {
	return m.name
}

func (m *Error) Inputs() *TupleType {
	return m.inputs
}

func (m *Error) FourBytes() FourBytes {
	return m.fourBytes
}

func (m *Error) Signature() string {
	return m.signature
}

func (m *Error) Is(data []byte) bool {
	return len(data) >= 4 && bytes.Equal(data[:4], m.fourBytes[:])
}

func (m *Error) DecodeValue(data []byte, val any) error {
	if len(data) < 4 {
		return fmt.Errorf("abi: error data too short")
	}
	if !bytes.Equal(data[:4], m.fourBytes[:]) {
		return fmt.Errorf("abi: selector mismatch for error %s", m.name)
	}
	return NewDecoder(m.config).DecodeValue(m.inputs.New(), data, val)
}

func (m *Error) DecodeValues(data []byte, vals ...any) error {
	if len(data) < 4 {
		return fmt.Errorf("abi: error data too short")
	}
	if !bytes.Equal(data[:4], m.fourBytes[:]) {
		return fmt.Errorf("abi: selector mismatch for error %s", m.name)
	}
	return NewDecoder(m.config).DecodeValues(m.inputs.New().(*TupleValue), data, vals...)
}

func (m *Error) String() string {
	return "error " + m.name + m.inputs.Type()
}

func (m *Error) generateSignature() {
	m.signature = m.name + m.inputs.CanonicalType()
}

func (m *Error) calculateFourBytes() {
	id := crypto.Keccak256([]byte(m.Signature()))
	copy(m.fourBytes[:], id[:4])
}
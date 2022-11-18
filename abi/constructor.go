package abi

import (
	"strings"
)

type Constructor struct {
	inputs *TupleType
	config *Config
}

func NewConstructor(inputs *TupleType) *Constructor {
	return NewConstructorWithConfig(inputs, DefaultConfig)
}

func NewConstructorWithConfig(inputs *TupleType, config *Config) *Constructor {
	return &Constructor{
		inputs: inputs,
		config: config,
	}
}

func (m *Constructor) Inputs() *TupleType {
	return m.inputs
}

func (m *Constructor) EncodeValue(val any) ([]byte, error) {
	encoded, err := NewEncoder(m.config).EncodeValue(m.inputs.New(), val)
	if err != nil {
		return nil, err
	}
	return encoded, nil
}

func (m *Constructor) EncodeValues(args ...any) ([]byte, error) {
	encoded, err := NewEncoder(m.config).EncodeValues(m.inputs.New().(*TupleValue), args...)
	if err != nil {
		return nil, err
	}
	return encoded, nil
}

func (m *Constructor) String() string {
	var buf strings.Builder
	buf.WriteString("constructor(")
	for i, typ := range m.inputs.Elements() {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(typ.Type.Type())
	}
	buf.WriteByte(')')
	return buf.String()
}

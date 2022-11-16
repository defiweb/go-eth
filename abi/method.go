package abi

import (
	"strings"

	"github.com/defiweb/go-eth/crypto"
	"github.com/defiweb/go-eth/hexutil"
)

type FourBytes [4]byte

type Method struct {
	Name    string
	Inputs  *TupleType
	Outputs *TupleType
	Config  *Config
}

func (m *Method) Encode(val any) ([]byte, error) {
	encoded, err := NewEncoder(m.Config).EncodeValue(m.Inputs.New(), val)
	if err != nil {
		return nil, err
	}
	return append(m.FourBytes().Bytes(), encoded...), nil
}

func (m *Method) EncodeArgs(args ...any) ([]byte, error) {
	encoded, err := NewEncoder(m.Config).EncodeValues(m.Inputs.New().(*TupleValue), args...)
	if err != nil {
		return nil, err
	}
	return append(m.FourBytes().Bytes(), encoded...), nil
}

func (m *Method) Decode(data []byte, val any) error {
	return NewDecoder(m.Config).DecodeValue(m.Outputs.New(), data, val)
}

func (m *Method) DecodeValues(data []byte, vals ...any) error {
	return NewDecoder(m.Config).DecodeValues(m.Outputs.New().(*TupleValue), data, vals...)
}

func (m *Method) FourBytes() FourBytes {
	id := crypto.Keccak256([]byte(m.Signature()))
	var f FourBytes
	copy(f[:], id[:4])
	return f
}

func (m *Method) String() string {
	var buf strings.Builder
	buf.WriteString(m.Name)
	buf.WriteByte('(')
	for i, typ := range m.Inputs.Elements() {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(typ.Type.Type())
	}
	buf.WriteByte(')')
	if m.Outputs.Size() > 0 {
		buf.WriteString(" returns (")
		for i, typ := range m.Outputs.Elements() {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(typ.Type.Type())
		}
		buf.WriteByte(')')
	}
	return buf.String()
}

func (m *Method) Signature() string {
	var buf strings.Builder
	buf.WriteString(m.Name)
	buf.WriteByte('(')
	for i, param := range m.Inputs.Elements() {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(param.Type.CanonicalType())
	}
	buf.WriteByte(')')
	return buf.String()
}

func (f FourBytes) Bytes() []byte {
	return f[:]
}

func (f FourBytes) Hex() string {
	return hexutil.BytesToHex(f[:])
}

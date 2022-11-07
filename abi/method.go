package abi

import (
	"strconv"
	"strings"

	"web3rpc/crypto"
	"web3rpc/hexutil"
)

type FourBytes [4]byte

type Method struct {
	Name    string
	Inputs  []TypeDefinition
	Outputs []TypeDefinition
	Config  *Config
}

func ParseMethod(signature string) (*Method, error) {
	return DefaultConfig.ParseMethod(signature)
}

func MustParseMethod(signature string) *Method {
	m, err := ParseMethod(signature)
	if err != nil {
		panic(err)
	}
	return m
}

func (m *Method) Encode(val any) ([]byte, error) {
	tuple, err := m.Config.NewTypeList(m.Inputs)
	if err != nil {
		return nil, err
	}
	encoded, err := EncodeValue(tuple, val)
	if err != nil {
		return nil, err
	}
	return append(m.FourBytes().Bytes(), encoded...), nil
}

func (m *Method) EncodeArgs(args ...any) ([]byte, error) {
	tuple, err := m.Config.NewTypeList(m.Inputs)
	if err != nil {
		return nil, err
	}
	encoded, err := EncodeValues(tuple, args...)
	if err != nil {
		return nil, err
	}
	return append(m.FourBytes().Bytes(), encoded...), nil
}

func (m *Method) Decode(data []byte, val any) error {
	tuple, err := m.Config.NewTypeList(m.Outputs)
	if err != nil {
		return err
	}
	return DecodeValue(tuple, data, val)
}

func (m *Method) DecodeValues(data []byte, vals ...any) error {
	tuple, err := m.Config.NewTypeList(m.Outputs)
	if err != nil {
		return err
	}
	return DecodeValues(tuple, data, vals...)
}

func (m *Method) InputTuple() (*TupleType, error) {
	return m.Config.NewTypeList(m.Inputs)
}

func (m *Method) OutputTuple() (*TupleType, error) {
	return m.Config.NewTypeList(m.Outputs)
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
	for i, param := range m.Inputs {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(param.String())
	}
	buf.WriteByte(')')
	if len(m.Outputs) > 0 {
		buf.WriteString(" returns (")
		for i, param := range m.Outputs {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(param.String())
		}
		buf.WriteByte(')')
	}
	return buf.String()
}

func (m *Method) Signature() string {
	var buf strings.Builder
	buf.WriteString(m.Name)
	buf.WriteByte('(')
	for i, param := range m.Inputs {
		if i > 0 {
			buf.WriteString(",")
		}
		typeSignature(&buf, param)
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

func typeSignature(buf *strings.Builder, typ TypeDefinition) {
	if len(typ.CanonicalType) > 0 {
		buf.WriteString(typ.CanonicalType)
		return
	}
	switch {
	case len(typ.Type) > 0:
		buf.WriteString(typ.Type)
	case len(typ.Tuple) > 0:
		buf.WriteByte('(')
		for i, c := range typ.Tuple {
			if i > 0 {
				buf.WriteByte(',')
			}
			typeSignature(buf, c)
		}
		buf.WriteByte(')')
	}
	if len(typ.Arrays) > 0 {
		for _, a := range typ.Arrays {
			if a < 0 {
				buf.WriteString("[]")
			} else {
				buf.WriteByte('[')
				buf.WriteString(strconv.Itoa(a))
				buf.WriteByte(']')
			}
		}
	}
}

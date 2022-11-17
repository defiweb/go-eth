package abi

import (
	"fmt"
	"strings"

	"github.com/defiweb/go-eth/crypto"
	"github.com/defiweb/go-eth/types"
)

type Event struct {
	Name   string
	Inputs *EventTupleType
	Config *Config
}

func (e *Event) Decode(topics []types.Hash, data []byte, val any) error {
	if len(topics) != e.Inputs.IndexedSize()+1 || len(topics) == 0 {
		return fmt.Errorf("abi: wrong number of topics for event %s", e.Name)
	}
	if topics[0] != e.Topic0() {
		return fmt.Errorf("abi: topic0 mismatch for event %s", e.Name)
	}
	return NewDecoder(e.Config).DecodeValue(
		e.Inputs.New(),
		e.mergeData(topics[1:], data),
		val,
	)
}

func (e *Event) DecodeValues(topics []types.Hash, data []byte, vals ...any) error {
	if len(topics) != e.Inputs.IndexedSize()+1 || len(topics) == 0 {
		return fmt.Errorf("abi: wrong number of topics for event %s", e.Name)
	}
	if topics[0] != e.Topic0() {
		return fmt.Errorf("abi: topic0 mismatch for event %s", e.Name)
	}
	return NewDecoder(e.Config).DecodeValues(
		e.Inputs.New().(*TupleValue),
		e.mergeData(topics[1:], data),
		vals...,
	)
}

func (e *Event) Topic0() types.Hash {
	return crypto.Keccak256([]byte(e.Signature()))
}

func (e *Event) String() string {
	var buf strings.Builder
	buf.WriteString(e.Name)
	buf.WriteByte('(')
	for i, typ := range e.Inputs.Elements() {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(typ.Type.Type())
	}
	return buf.String()
}

func (e *Event) Signature() string {
	var buf strings.Builder
	buf.WriteString(e.Name)
	buf.WriteByte('(')
	for i, param := range e.Inputs.Elements() {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(param.Type.CanonicalType())
	}
	buf.WriteByte(')')
	return buf.String()
}

func (e *Event) mergeData(topics []types.Hash, data []byte) []byte {
	merged := make([]byte, len(topics)*32+len(data))
	for i, topic := range topics {
		copy(merged[i*32:], topic[:])
	}
	copy(merged[len(topics)*32:], data)
	return merged
}

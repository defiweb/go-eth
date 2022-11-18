package abi

import (
	"fmt"

	"github.com/defiweb/go-eth/crypto"
	"github.com/defiweb/go-eth/types"
)

type Event struct {
	name   string
	inputs *EventTupleType
	config *Config

	topic0    types.Hash
	signature string
}

func NewEvent(name string, inputs *EventTupleType) *Event {
	return NewEventWithConfig(name, inputs, DefaultConfig)
}

func NewEventWithConfig(name string, inputs *EventTupleType, config *Config) *Event {
	e := &Event{
		name:   name,
		inputs: inputs,
		config: config,
	}
	e.generateSignature()
	e.calculateTopic0()
	return e
}

func (e *Event) Name() string {
	return e.name
}

func (e *Event) Inputs() *EventTupleType {
	return e.inputs
}

func (e *Event) Topic0() types.Hash {
	return e.topic0
}

func (e *Event) Signature() string {
	return e.signature
}

func (e *Event) DecodeValue(topics []types.Hash, data []byte, val any) error {
	if len(topics) != e.inputs.IndexedSize()+1 {
		return fmt.Errorf("abi: wrong number of topics for event %s", e.name)
	}
	if topics[0] != e.topic0 {
		return fmt.Errorf("abi: topic0 mismatch for event %s", e.name)
	}
	return NewDecoder(e.config).DecodeValue(
		e.inputs.New(),
		e.mergeData(topics[1:], data),
		val,
	)
}

func (e *Event) DecodeValues(topics []types.Hash, data []byte, vals ...any) error {
	if len(topics) != e.inputs.IndexedSize()+1 {
		return fmt.Errorf("abi: wrong number of topics for event %s", e.name)
	}
	if topics[0] != e.topic0 {
		return fmt.Errorf("abi: topic0 mismatch for event %s", e.name)
	}
	return NewDecoder(e.config).DecodeValues(
		e.inputs.New().(*TupleValue),
		e.mergeData(topics[1:], data),
		vals...,
	)
}

func (e *Event) String() string {
	return fmt.Sprintf("event %s%s", e.name, e.inputs.Type())
}

func (e *Event) calculateTopic0() {
	e.topic0 = crypto.Keccak256([]byte(e.signature))
}

func (e *Event) generateSignature() {
	e.signature = fmt.Sprintf("%s%s", e.name, e.inputs.Type())
}

func (e *Event) mergeData(topics []types.Hash, data []byte) []byte {
	if len(topics) == 0 {
		return data
	}
	merged := make([]byte, len(topics)*types.HashLength+len(data))
	for i, topic := range topics {
		copy(merged[i*types.HashLength:], topic[:])
	}
	copy(merged[len(topics)*types.HashLength:], data)
	return merged
}

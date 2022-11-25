package abi

import (
	"fmt"
	"strings"

	"github.com/defiweb/go-eth/crypto"
	"github.com/defiweb/go-eth/types"
)

// Event represents an event in an jsonABI. The event can be used to decode events
// emitted by a contract.
type Event struct {
	name      string
	inputs    *EventTupleType
	anonymous bool
	config    *Config

	topic0    types.Hash
	signature string
}

// NewEvent creates a new Event instance.
func NewEvent(name string, inputs *EventTupleType, anonymous bool) *Event {
	return NewEventWithConfig(name, inputs, anonymous, DefaultConfig)
}

// NewEventWithConfig creates a new Event instance with a custom config.
func NewEventWithConfig(name string, inputs *EventTupleType, anonymous bool, config *Config) *Event {
	e := &Event{
		name:      name,
		inputs:    inputs,
		anonymous: anonymous,
		config:    config,
	}
	e.generateSignature()
	e.calculateTopic0()
	return e
}

// Name returns the name of the event.
func (e *Event) Name() string {
	return e.name
}

// Inputs returns the input arguments of the event as a tuple type.
func (e *Event) Inputs() *EventTupleType {
	return e.inputs
}

// Topic0 returns the first topic of the event, that is, the Keccak256 hash of
// the event signature.
func (e *Event) Topic0() types.Hash {
	return e.topic0
}

// Signature returns the event signature, that is, the event name and the
// canonical type of the input arguments.
func (e *Event) Signature() string {
	return e.signature
}

// DecodeValue decodes the event into a map or structure. If a structure is
// given, it must have fields with the same names as the event arguments.
func (e *Event) DecodeValue(topics []types.Hash, data []byte, val any) error {
	if e.anonymous {
		return NewDecoder(e.config).DecodeValue(
			e.inputs.Value().(*TupleValue),
			data,
			val,
		)
	}
	if len(topics) != e.inputs.IndexedSize()+1 {
		return fmt.Errorf("abi: wrong number of topics for event %s", e.name)
	}
	if topics[0] != e.topic0 {
		return fmt.Errorf("abi: topic0 mismatch for event %s", e.name)
	}
	return NewDecoder(e.config).DecodeValue(
		e.inputs.Value(),
		e.mergeData(topics[1:], data),
		val,
	)
}

// DecodeValues decodes the event into a map or structure. If a structure is
// given, it must have fields with the same names as the event arguments.
func (e *Event) DecodeValues(topics []types.Hash, data []byte, vals ...any) error {
	if e.anonymous {
		return NewDecoder(e.config).DecodeValues(
			e.inputs.Value().(*TupleValue),
			data,
			vals...,
		)
	}
	if len(topics) != e.inputs.IndexedSize()+1 {
		return fmt.Errorf("abi: wrong number of topics for event %s", e.name)
	}
	if topics[0] != e.topic0 {
		return fmt.Errorf("abi: topic0 mismatch for event %s", e.name)
	}
	return NewDecoder(e.config).DecodeValues(
		e.inputs.Value().(*TupleValue),
		e.mergeData(topics[1:], data),
		vals...,
	)
}

// String returns the human-readable signature of the event.
func (e *Event) String() string {
	var buf strings.Builder
	buf.WriteString("event ")
	buf.WriteString(e.name)
	buf.WriteString(e.inputs.String())
	if e.anonymous {
		buf.WriteString(" anonymous")
	}
	return buf.String()
}

func (e *Event) calculateTopic0() {
	e.topic0 = crypto.Keccak256([]byte(e.signature))
}

func (e *Event) generateSignature() {
	e.signature = fmt.Sprintf("%s%s", e.name, e.inputs.CanonicalType())
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

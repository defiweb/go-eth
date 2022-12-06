package abi

import (
	"fmt"
	"strings"

	"github.com/defiweb/go-eth/crypto"
	"github.com/defiweb/go-eth/types"
)

// Event represents an event in an ABI. The event can be used to decode events
// emitted by a contract.
type Event struct {
	name      string
	inputs    *EventTupleType
	anonymous bool
	config    *ABI

	topic0    types.Hash
	signature string
}

// NewEvent creates a new Event instance.
func NewEvent(name string, inputs *EventTupleType, anonymous bool) *Event {
	return Default.NewEvent(name, inputs, anonymous)
}

// ParseEvent parses an event signature and returns a new Event.
//
// An event signature is similar to a method signature, but returns no values.
// It can be optionally prefixed with the "event" keyword.
//
// The following examples are valid signatures:
//
//   foo(int indexed,(uint256,bytes32)[])
//   foo(int indexed a, (uint256 b, bytes32 c)[] d)
//   event foo(int indexed a tuple(uint256 b, bytes32 c)[] d)
//
// This function is equivalent to calling Parser.ParseEvent with the default
// configuration.
func ParseEvent(signature string) (*Event, error) {
	return Default.ParseEvent(signature)
}

// MustParseEvent is like ParseEvent but panics on error.
func MustParseEvent(signature string) *Event {
	e, err := ParseEvent(signature)
	if err != nil {
		panic(err)
	}
	return e
}

// NewEvent creates a new Event instance.
func (a *ABI) NewEvent(name string, inputs *EventTupleType, anonymous bool) *Event {
	e := &Event{
		name:      name,
		inputs:    inputs,
		anonymous: anonymous,
		config:    a,
	}
	e.generateSignature()
	e.calculateTopic0()
	return e
}

// ParseEvent parses an event signature and returns a new Event.
//
// See ParseEvent for more information.
func (a *ABI) ParseEvent(signature string) (*Event, error) {
	return parseEvent(a, signature)
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
		return e.config.DecodeValue(e.inputs, data, val)
	}
	if len(topics) != e.inputs.IndexedSize()+1 {
		return fmt.Errorf("abi: wrong number of topics for event %s", e.name)
	}
	if topics[0] != e.topic0 {
		return fmt.Errorf("abi: topic0 mismatch for event %s", e.name)
	}
	return e.config.DecodeValue(e.inputs, e.mergeData(topics[1:], data), val)
}

// DecodeValues decodes the event into a map or structure. If a structure is
// given, it must have fields with the same names as the event arguments.
func (e *Event) DecodeValues(topics []types.Hash, data []byte, vals ...any) error {
	if e.anonymous {
		return e.config.DecodeValues(e.inputs, data, vals...)
	}
	if len(topics) != e.inputs.IndexedSize()+1 {
		return fmt.Errorf("abi: wrong number of topics for event %s", e.name)
	}
	if topics[0] != e.topic0 {
		return fmt.Errorf("abi: topic0 mismatch for event %s", e.name)
	}
	return e.config.DecodeValues(e.inputs, e.mergeData(topics[1:], data), vals...)
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

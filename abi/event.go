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
	abi       *ABI

	topic0    types.Hash
	signature string
}

// NewEvent creates a new Event instance.
//
// This method is rarely used, see ParseEvent for a more convenient way to
// create a new Event.
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
//	foo(int indexed,(uint256,bytes32)[])
//	foo(int indexed a, (uint256 b, bytes32 c)[] d)
//	event foo(int indexed a tuple(uint256 b, bytes32 c)[] d)
//
// This function is equivalent to calling Parser.ParseEvent with the default
// configuration.
func ParseEvent(signature string) (*Event, error) {
	return Default.ParseEvent(signature)
}

// MustParseEvent is like ParseEvent but panics on error.
func MustParseEvent(signature string) *Event {
	return Default.MustParseEvent(signature)
}

// NewEvent creates a new Event instance.
func (a *ABI) NewEvent(name string, inputs *EventTupleType, anonymous bool) *Event {
	if inputs == nil {
		inputs = NewEventTupleType()
	}
	e := &Event{
		name:      name,
		inputs:    inputs,
		anonymous: anonymous,
		abi:       a,
	}
	e.generateSignature()
	e.calculateTopic0()
	return e
}

// ParseEvent parses an event signature and returns a new Event.
//
// See ParseEvent for more information.
func (a *ABI) ParseEvent(signature string) (*Event, error) {
	return parseEvent(a, nil, signature)
}

// MustParseEvent is like ParseEvent but panics on error.
func (a *ABI) MustParseEvent(signature string) *Event {
	e, err := a.ParseEvent(signature)
	if err != nil {
		panic(err)
	}
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
		return e.abi.DecodeValue(e.inputs, data, val)
	}
	if len(topics) != e.inputs.IndexedSize()+1 {
		return fmt.Errorf("abi: wrong number of topics for event %s", e.name)
	}
	if topics[0] != e.topic0 {
		return fmt.Errorf("abi: topic0 mismatch for event %s", e.name)
	}
	// The anymapper package does not zero out values before decoding into
	// it, therefore we can decode topics and data into the same value.
	if len(topics) > 1 {
		if err := e.abi.DecodeValue(e.inputs.TopicsTuple(), hashSliceToBytes(topics[1:]), val); err != nil {
			return err
		}
	}
	if len(data) > 0 {
		if err := e.abi.DecodeValue(e.inputs.DataTuple(), data, val); err != nil {
			return err
		}
	}
	return nil
}

// MustDecodeValue is like DecodeValue but panics on error.
func (e *Event) MustDecodeValue(topics []types.Hash, data []byte, val any) {
	err := e.DecodeValue(topics, data, val)
	if err != nil {
		panic(err)
	}
}

// DecodeValues decodes the event into a map or structure. If a structure is
// given, it must have fields with the same names as the event arguments.
func (e *Event) DecodeValues(topics []types.Hash, data []byte, vals ...any) error {
	if e.anonymous {
		return e.abi.DecodeValues(e.inputs, data, vals...)
	}
	if len(topics) != e.inputs.IndexedSize()+1 {
		return fmt.Errorf("abi: wrong number of topics for event %s", e.name)
	}
	if topics[0] != e.topic0 {
		return fmt.Errorf("abi: topic0 mismatch for event %s", e.name)
	}
	indexedVals := make([]any, 0, e.inputs.IndexedSize())
	dataVals := make([]any, 0, e.inputs.DataSize())
	for i := range e.inputs.Elements() {
		if i >= len(vals) {
			break
		}
		if e.inputs.Elements()[i].Indexed {
			indexedVals = append(indexedVals, vals[i])
		} else {
			dataVals = append(dataVals, vals[i])
		}
	}
	// The anymapper package does not zero out values before decoding into
	// it, therefore we can decode topics and data into the same value.
	if len(topics) > 1 {
		if err := e.abi.DecodeValues(e.inputs.TopicsTuple(), hashSliceToBytes(topics[1:]), indexedVals...); err != nil {
			return err
		}
	}
	if len(data) > 0 {
		if err := e.abi.DecodeValues(e.inputs.DataTuple(), data, dataVals...); err != nil {
			return err
		}
	}
	return nil
}

// MustDecodeValues is like DecodeValues but panics on error.
func (e *Event) MustDecodeValues(topics []types.Hash, data []byte, vals ...any) {
	err := e.DecodeValues(topics, data, vals...)
	if err != nil {
		panic(err)
	}
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
	e.topic0 = types.Hash(crypto.Keccak256([]byte(e.signature)))
}

func (e *Event) generateSignature() {
	e.signature = fmt.Sprintf("%s%s", e.name, e.inputs.CanonicalType())
}

func hashSliceToBytes(hashes []types.Hash) []byte {
	buf := make([]byte, len(hashes)*types.HashLength)
	for i, hash := range hashes {
		copy(buf[i*types.HashLength:], hash[:])
	}
	return buf
}

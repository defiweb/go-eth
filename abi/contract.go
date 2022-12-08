package abi

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/defiweb/go-sigparser"
)

// Contract provides a high-level API for interacting with a contract. It can
// be created from a JSON ABI definition using the ParseJSON function or from
// a list of signatures using the ParseSignatures function.
type Contract struct {
	constructor        *Constructor
	methods            map[string]*Method
	methodsBySignature map[string]*Method
	events             map[string]*Event
	errors             map[string]*Error
}

// ParseJSON parses the given ABI JSON and returns a Contract instance.
func ParseJSON(data []byte) (*Contract, error) {
	return Default.ParseJSON(data)
}

// ParseSignatures parses list of signatures and returns a Contract instance.
// Signatures must be prefixed with the kind, e.g. "function" or "event".
//
// It accepts signatures in the same format as ParseConstructor, ParseMethod,
// ParseEvent, and ParseError functions.
func ParseSignatures(signatures ...string) (*Contract, error) {
	return Default.ParseSignatures(signatures...)
}

// MustParseJSON is like ParseJSON but panics on error.
func MustParseJSON(data []byte) *Contract {
	abi, err := ParseJSON(data)
	if err != nil {
		panic(err)
	}
	return abi
}

// MustParseSignatures is like ParseSignatures but panics on error.
func MustParseSignatures(signatures ...string) *Contract {
	abi, err := ParseSignatures(signatures...)
	if err != nil {
		panic(err)
	}
	return abi
}

// Constructor returns the contract constructor.
func (a *Contract) Constructor() *Constructor {
	return a.constructor
}

// Method returns the method with the given name.
func (a *Contract) Method(name string) *Method {
	return a.methods[name]
}

// MethodBySignature returns the method with the given signature.
func (a *Contract) MethodBySignature(signature string) *Method {
	return a.methodsBySignature[signature]
}

// Event returns the event with the given name.
func (a *Contract) Event(name string) *Event {
	return a.events[name]
}

// Error returns the error with the given name.
func (a *Contract) Error(name string) *Error {
	return a.errors[name]
}

// ParseJSON parses the given ABI JSON and returns a Contract instance.
func (a *ABI) ParseJSON(data []byte) (*Contract, error) {
	var fields []jsonField
	if err := json.Unmarshal(data, &fields); err != nil {
		return nil, err
	}
	abi := &Contract{
		methods:            make(map[string]*Method),
		methodsBySignature: make(map[string]*Method),
		events:             make(map[string]*Event),
		errors:             make(map[string]*Error),
	}
	for _, f := range fields {
		switch f.Type {
		case "constructor":
			inputs, err := f.Inputs.toTupleType(a)
			if err != nil {
				return nil, err
			}
			abi.constructor = a.NewConstructor(inputs)
		case "function", "":
			inputs, err := f.Inputs.toTupleType(a)
			if err != nil {
				return nil, err
			}
			outputs, err := f.Outputs.toTupleType(a)
			if err != nil {
				return nil, err
			}
			method := a.NewMethod(f.Name, inputs, outputs)
			abi.methods[f.Name] = method
			abi.methodsBySignature[method.Signature()] = method
		case "event":
			inputs, err := f.Inputs.toEventTupleType(a)
			if err != nil {
				return nil, err
			}
			abi.events[f.Name] = a.NewEvent(f.Name, inputs, f.Anonymous)
		case "error":
			inputs, err := f.Inputs.toTupleType(a)
			if err != nil {
				return nil, err
			}
			abi.errors[f.Name] = a.NewError(f.Name, inputs)
		case "fallback":
		case "receive":
		default:
			return nil, fmt.Errorf("unknown type: %s", f.Type)
		}
	}
	return abi, nil
}

// ParseSignatures parses list of signatures and returns a Contract instance.
// Signatures must be prefixed with the kind, e.g. "constructor" or "event".
// For functions, the "function" prefix can be omitted.
func (a *ABI) ParseSignatures(signatures ...string) (*Contract, error) {
	abi := &Contract{
		methods:            make(map[string]*Method),
		methodsBySignature: make(map[string]*Method),
		events:             make(map[string]*Event),
		errors:             make(map[string]*Error),
	}
	for _, s := range signatures {
		sig, err := sigparser.ParseSignature(s)
		if err != nil {
			return nil, err
		}
		switch sig.Kind {
		case sigparser.ConstructorKind:
			constructor, err := newConstructorFromSig(a, sig)
			if err != nil {
				return nil, err
			}
			abi.constructor = constructor
		case sigparser.FunctionKind, sigparser.UnknownKind:
			method, err := newMethodFromSig(a, sig)
			if err != nil {
				return nil, err
			}
			abi.methods[method.Name()] = method
			abi.methodsBySignature[method.Signature()] = method
		case sigparser.EventKind:
			event, err := newEventFromSig(a, sig)
			if err != nil {
				return nil, err
			}
			abi.events[event.Name()] = event
		case sigparser.ErrorKind:
			errsig, err := newErrorFromSig(a, sig)
			if err != nil {
				return nil, err
			}
			abi.errors[errsig.Name()] = errsig
		default:
			return nil, fmt.Errorf("unknown kind: %s", sig.Kind)
		}
	}
	return abi, nil
}

type jsonField struct {
	Type            string         `json:"type"`
	Name            string         `json:"name"`
	Constant        bool           `json:"constant"`
	Anonymous       bool           `json:"anonymous"`
	StateMutability string         `json:"stateMutability"`
	Inputs          jsonParameters `json:"inputs"`
	Outputs         jsonParameters `json:"outputs"`
}

type jsonParameters []jsonParameter

// toTupleType converts parameters to a TupleType type.
func (a jsonParameters) toTupleType(abi *ABI) (*TupleType, error) {
	var elems []TupleTypeElem
	for _, param := range a {
		typ, err := param.toType(abi)
		if err != nil {
			return nil, err
		}
		elems = append(elems, TupleTypeElem{
			Name: param.Name,
			Type: typ,
		})
	}
	return NewTupleType(elems...), nil
}

// toEventTupleType converts parameters to a EventTupleType type.
func (a jsonParameters) toEventTupleType(abi *ABI) (*EventTupleType, error) {
	var elems []EventTupleElem
	for _, param := range a {
		typ, err := param.toType(abi)
		if err != nil {
			return nil, err
		}
		elems = append(elems, EventTupleElem{
			Name:    param.Name,
			Indexed: param.Indexed,
			Type:    typ,
		})
	}
	return NewEventTupleType(elems...), nil
}

type jsonParameter struct {
	Name       string         `json:"name"`
	Type       string         `json:"type"`
	Indexed    bool           `json:"indexed"`
	Components jsonParameters `json:"components"`
}

// toType converts a jsonParameter to a Type.
func (a jsonParameter) toType(abi *ABI) (typ Type, err error) {
	name, arrays, err := parseArrays(a.Type)
	if err != nil {
		return nil, err
	}
	switch {
	case len(arrays) > 0:
		a.Type = name
		if typ, err = a.toType(abi); err != nil {
			return nil, err
		}
		for i := len(arrays) - 1; i >= 0; i-- {
			if arrays[i] == -1 {
				typ = NewArrayType(typ)
			} else {
				typ = NewFixedArrayType(typ, arrays[i])
			}
		}
		return typ, nil
	case len(a.Components) > 0:
		tuple := make([]TupleTypeElem, len(a.Components))
		for i, comp := range a.Components {
			tuple[i].Name = comp.Name
			tuple[i].Type, err = comp.toType(abi)
			if err != nil {
				return nil, err
			}
		}
		return NewTupleType(tuple...), nil
	default:
		if typ = abi.Types[name]; typ != nil {
			return typ, nil
		}
		return nil, fmt.Errorf("abi: unknown type %q", a.Type)
	}
}

// parseArray parses type name and returns the name and array dimensions.
// For example, "uint256[][3]" will return "uint256" and [-1, 3].
// For unbounded arrays, the dimension is -1.
func parseArrays(typ string) (name string, arrays []int, err error) {
	openBracket := strings.Index(typ, "[")
	if openBracket == -1 {
		name = typ
		return
	}
	name = typ[:openBracket]
	for {
		closeBracket := openBracket
		for closeBracket < len(typ) && typ[closeBracket] != ']' {
			closeBracket++
		}
		if openBracket >= closeBracket {
			return "", nil, fmt.Errorf("abi: invalid type %q", typ)
		}
		n := typ[openBracket+1 : closeBracket]
		if len(n) == 0 {
			arrays = append(arrays, -1)
		} else {
			i, err := strconv.Atoi(n)
			if err != nil {
				return "", nil, err
			}
			if i <= 0 {
				return "", nil, fmt.Errorf("abi: invalid array size %d", i)
			}
			arrays = append(arrays, i)
		}
		if closeBracket+1 == len(typ) {
			break
		}
		openBracket = closeBracket + 1
	}
	return
}

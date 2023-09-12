package abi

import (
	"fmt"

	"github.com/defiweb/go-sigparser"
)

// parseType parses a type signature and returns a Type.
func parseType(abi *ABI, signature string) (Type, error) {
	p, err := sigparser.ParseParameter(signature)
	if err != nil {
		return nil, err
	}
	return newTypeFromSig(abi, p)
}

// parseStruct parses a structure definition and returns a Type.
func parseStruct(abi *ABI, signature string) (Type, error) {
	p, err := sigparser.ParseStruct(signature)
	if err != nil {
		return nil, err
	}
	return newTypeFromSig(abi, p)
}

// parseConstructor parses a constructor signature and returns a Constructor.
func parseConstructor(abi *ABI, signature string) (*Constructor, error) {
	s, err := sigparser.ParseSignatureAs(sigparser.ConstructorKind, signature)
	if err != nil {
		return nil, err
	}
	return newConstructorFromSig(abi, s)
}

// parseError parses an error signature and returns an Error.
func parseError(abi *ABI, signature string) (*Error, error) {
	s, err := sigparser.ParseSignatureAs(sigparser.ErrorKind, signature)
	if err != nil {
		return nil, err
	}
	return newErrorFromSig(abi, s)
}

// parseEvent parses an event signature and returns an Event.
func parseEvent(abi *ABI, signature string) (*Event, error) {
	s, err := sigparser.ParseSignatureAs(sigparser.EventKind, signature)
	if err != nil {
		return nil, err
	}
	return newEventFromSig(abi, s)
}

// parseMethod parses a method signature and returns a Method.
func parseMethod(abi *ABI, signature string) (*Method, error) {
	s, err := sigparser.ParseSignatureAs(sigparser.FunctionKind, signature)
	if err != nil {
		return nil, err
	}
	return newMethodFromSig(abi, s)
}

// newConstructorFromSig creates a new constructor from a sigparser.Signature.
func newConstructorFromSig(abi *ABI, s sigparser.Signature) (*Constructor, error) {
	var in []TupleTypeElem
	for _, param := range s.Inputs {
		typ, err := newTypeFromSig(abi, param)
		if err != nil {
			return nil, err
		}
		in = append(in, TupleTypeElem{
			Name: param.Name,
			Type: typ,
		})
	}
	return abi.NewConstructor(NewTupleType(in...)), nil
}

// newErrorFromSig creates a new error from a sigparser.Signature.
func newErrorFromSig(abi *ABI, s sigparser.Signature) (*Error, error) {
	var in []TupleTypeElem
	for _, param := range s.Inputs {
		typ, err := newTypeFromSig(abi, param)
		if err != nil {
			return nil, err
		}
		in = append(in, TupleTypeElem{
			Name: param.Name,
			Type: typ,
		})
	}
	return abi.NewError(s.Name, NewTupleType(in...)), nil
}

// newEventFromSig creates a new event from a sigparser.Signature.
func newEventFromSig(abi *ABI, s sigparser.Signature) (*Event, error) {
	var in []EventTupleElem
	if len(s.Inputs) == 0 {
		return nil, fmt.Errorf("abi: event %q has no inputs", s.Name)
	}
	for _, param := range s.Inputs {
		typ, err := newTypeFromSig(abi, param)
		if err != nil {
			return nil, err
		}
		in = append(in, EventTupleElem{
			Name:    param.Name,
			Indexed: param.Indexed,
			Type:    typ,
		})
	}
	anonymous := false
	for _, param := range s.Modifiers {
		if param == "anonymous" {
			anonymous = true
			break
		}
	}
	return abi.NewEvent(s.Name, NewEventTupleType(in...), anonymous), nil
}

// newMethodFromSig creates a new method from a sigparser.Signature.
func newMethodFromSig(abi *ABI, s sigparser.Signature) (*Method, error) {
	var (
		in  []TupleTypeElem
		out []TupleTypeElem
	)
	for _, param := range s.Inputs {
		typ, err := newTypeFromSig(abi, param)
		if err != nil {
			return nil, err
		}
		in = append(in, TupleTypeElem{
			Name: param.Name,
			Type: typ,
		})
	}
	for _, param := range s.Outputs {
		typ, err := newTypeFromSig(abi, param)
		if err != nil {
			return nil, err
		}
		out = append(out, TupleTypeElem{
			Name: param.Name,
			Type: typ,
		})
	}
	return abi.NewMethod(s.Name, NewTupleType(in...), NewTupleType(out...)), nil
}

// newTypeFromSig creates a new type from a sigparser.Parameter.
func newTypeFromSig(abi *ABI, s sigparser.Parameter) (typ Type, err error) {
	switch {
	case len(s.Arrays) > 0:
		arrays := s.Arrays
		s.Arrays = nil
		if typ, err = newTypeFromSig(abi, s); err != nil {
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
	case len(s.Tuple) > 0:
		tuple := make([]TupleTypeElem, len(s.Tuple))
		for i, param := range s.Tuple {
			tuple[i].Name = param.Name
			tuple[i].Type, err = newTypeFromSig(abi, param)
			if err != nil {
				return nil, err
			}
		}
		return NewTupleType(tuple...), nil
	default:
		if typ = abi.Types[s.Type]; typ != nil {
			return typ, nil
		}
		return nil, fmt.Errorf("abi: unknown type %q", s.Type)
	}
}

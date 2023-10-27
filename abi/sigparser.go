package abi

import (
	"fmt"

	"github.com/defiweb/go-sigparser"
)

// parseType parses a type signature and returns a Type.
//
// The extraTypes map is used to resolve types that are not part of the ABI.
func parseType(abi *ABI, extraTypes map[string]Type, signature string) (Type, error) {
	p, err := sigparser.ParseParameter(signature)
	if err != nil {
		return nil, err
	}
	return newTypeFromSig(abi, extraTypes, p)
}

// parseStruct parses a structure definition and returns a Type.
//
// The extraTypes map is used to resolve types that are not part of the ABI.
func parseStruct(abi *ABI, extraTypes map[string]Type, signature string) (Type, error) {
	p, err := sigparser.ParseStruct(signature)
	if err != nil {
		return nil, err
	}
	return newTypeFromSig(abi, extraTypes, p)
}

// parseConstructor parses a constructor signature and returns a Constructor.
//
// The extraTypes map is used to resolve types that are not part of the ABI.
func parseConstructor(abi *ABI, extraTypes map[string]Type, signature string) (*Constructor, error) {
	s, err := sigparser.ParseSignatureAs(sigparser.ConstructorKind, signature)
	if err != nil {
		return nil, err
	}
	return newConstructorFromSig(abi, extraTypes, s)
}

// parseError parses an error signature and returns an Error.
//
// The extraTypes map is used to resolve types that are not part of the ABI.
func parseError(abi *ABI, extraTypes map[string]Type, signature string) (*Error, error) {
	s, err := sigparser.ParseSignatureAs(sigparser.ErrorKind, signature)
	if err != nil {
		return nil, err
	}
	return newErrorFromSig(abi, extraTypes, s)
}

// parseEvent parses an event signature and returns an Event.
//
// The extraTypes map is used to resolve types that are not part of the ABI.
func parseEvent(abi *ABI, extraTypes map[string]Type, signature string) (*Event, error) {
	s, err := sigparser.ParseSignatureAs(sigparser.EventKind, signature)
	if err != nil {
		return nil, err
	}
	return newEventFromSig(abi, extraTypes, s)
}

// parseMethod parses a method signature and returns a Method.
//
// The extraTypes map is used to resolve types that are not part of the ABI.
func parseMethod(abi *ABI, extraTypes map[string]Type, signature string) (*Method, error) {
	s, err := sigparser.ParseSignatureAs(sigparser.FunctionKind, signature)
	if err != nil {
		return nil, err
	}
	return newMethodFromSig(abi, extraTypes, s)
}

// newConstructorFromSig creates a new constructor from a sigparser.Signature.
//
// The extraTypes map is used to resolve types that are not part of the ABI.
func newConstructorFromSig(abi *ABI, extraTypes map[string]Type, s sigparser.Signature) (*Constructor, error) {
	var in []TupleTypeElem
	for _, param := range s.Inputs {
		typ, err := newTypeFromSig(abi, extraTypes, param)
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
//
// The extraTypes map is used to resolve types that are not part of the ABI.
func newErrorFromSig(abi *ABI, extraTypes map[string]Type, s sigparser.Signature) (*Error, error) {
	var in []TupleTypeElem
	for _, param := range s.Inputs {
		typ, err := newTypeFromSig(abi, extraTypes, param)
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
//
// The extraTypes map is used to resolve types that are not part of the ABI.
func newEventFromSig(abi *ABI, extraTypes map[string]Type, s sigparser.Signature) (*Event, error) {
	var in []EventTupleElem
	if len(s.Inputs) == 0 {
		return nil, fmt.Errorf("abi: event %q has no inputs", s.Name)
	}
	for _, param := range s.Inputs {
		typ, err := newTypeFromSig(abi, extraTypes, param)
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
//
// The extraTypes map is used to resolve types that are not part of the ABI.
func newMethodFromSig(abi *ABI, extraTypes map[string]Type, s sigparser.Signature) (*Method, error) {
	var (
		in  []TupleTypeElem
		out []TupleTypeElem
	)
	for _, param := range s.Inputs {
		typ, err := newTypeFromSig(abi, extraTypes, param)
		if err != nil {
			return nil, err
		}
		in = append(in, TupleTypeElem{
			Name: param.Name,
			Type: typ,
		})
	}
	for _, param := range s.Outputs {
		typ, err := newTypeFromSig(abi, extraTypes, param)
		if err != nil {
			return nil, err
		}
		out = append(out, TupleTypeElem{
			Name: param.Name,
			Type: typ,
		})
	}
	mutability := StateMutabilityUnknown
	for _, modifier := range s.Modifiers {
		switch modifier {
		case "pure":
			mutability = StateMutabilityPure
		case "view":
			mutability = StateMutabilityView
		case "payable":
			mutability = StateMutabilityPayable
		case "nonpayable":
			mutability = StateMutabilityNonPayable
		}
	}
	return abi.NewMethod(s.Name, NewTupleType(in...), NewTupleType(out...), mutability), nil
}

// newTypeFromSig creates a new type from a sigparser.Parameter.
//
// The extraTypes map is used to resolve types that are not part of the ABI.
func newTypeFromSig(abi *ABI, extraTypes map[string]Type, s sigparser.Parameter) (typ Type, err error) {
	switch {
	case len(s.Arrays) > 0:
		arrays := s.Arrays
		s.Arrays = nil
		if typ, err = newTypeFromSig(abi, extraTypes, s); err != nil {
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
			tuple[i].Type, err = newTypeFromSig(abi, extraTypes, param)
			if err != nil {
				return nil, err
			}
		}
		return NewTupleType(tuple...), nil
	default:
		if typ = extraTypes[s.Type]; typ != nil {
			return typ, nil
		}
		if typ = abi.Types[s.Type]; typ != nil {
			return typ, nil
		}
		return nil, fmt.Errorf("abi: unknown type %q", s.Type)
	}
}

package abi

import (
	"fmt"

	"github.com/defiweb/go-sigparser"
)

func ParseType(signature string) (Type, error) {
	return NewParser(DefaultConfig).ParseType(signature)
}

func ParseConstructor(signature string) (*Constructor, error) {
	return NewParser(DefaultConfig).ParseConstructor(signature)
}

func ParseMethod(signature string) (*Method, error) {
	return NewParser(DefaultConfig).ParseMethod(signature)
}

func ParseEvent(signature string) (*Event, error) {
	return NewParser(DefaultConfig).ParseEvent(signature)
}

func ParseError(signature string) (*Error, error) {
	return NewParser(DefaultConfig).ParseError(signature)
}

func MustParseType(signature string) Type {
	t, err := ParseType(signature)
	if err != nil {
		panic(err)
	}
	return t
}

func MustParseConstructor(signature string) *Constructor {
	c, err := ParseConstructor(signature)
	if err != nil {
		panic(err)
	}
	return c
}

func MustParseMethod(signature string) *Method {
	m, err := ParseMethod(signature)
	if err != nil {
		panic(err)
	}
	return m
}

func MustParseEvent(signature string) *Event {
	e, err := ParseEvent(signature)
	if err != nil {
		panic(err)
	}
	return e
}

func MustParseError(signature string) *Error {
	e, err := ParseError(signature)
	if err != nil {
		panic(err)
	}
	return e
}

type Parser struct {
	Config *Config
}

func NewParser(c *Config) *Parser {
	return &Parser{Config: c}
}

func (p *Parser) ParseType(signature string) (Type, error) {
	typ, err := sigparser.ParseParameter(signature)
	if err != nil {
		return nil, err
	}
	return p.toType(typ)
}

func (p *Parser) ParseConstructor(signature string) (*Constructor, error) {
	constructor, err := sigparser.ParseSignature(signature)
	if err != nil {
		return nil, err
	}
	if !isKind(constructor.Kind, sigparser.ConstructorKind, sigparser.UnknownKind) {
		return nil, fmt.Errorf("abi: expected signature type for constructor, got %s", constructor.Kind)
	}
	return p.toConstructor(constructor)
}

func (p *Parser) ParseMethod(signature string) (*Method, error) {
	method, err := sigparser.ParseSignature(signature)
	if err != nil {
		return nil, err
	}
	if !isKind(method.Kind, sigparser.FunctionKind, sigparser.UnknownKind) {
		return nil, fmt.Errorf("abi: expected signature type for method, got %s", method.Kind)
	}
	return p.toMethod(method)
}

func (p *Parser) ParseEvent(signature string) (*Event, error) {
	event, err := sigparser.ParseSignature(signature)
	if err != nil {
		return nil, err
	}
	if !isKind(event.Kind, sigparser.EventKind, sigparser.UnknownKind) {
		return nil, fmt.Errorf("abi: expected signature type for event, got %s", event.Kind)
	}
	return p.toEvent(event)
}

func (p *Parser) ParseError(signature string) (*Error, error) {
	errsig, err := sigparser.ParseSignature(signature)
	if err != nil {
		return nil, err
	}
	if !isKind(errsig.Kind, sigparser.ErrorKind, sigparser.UnknownKind) {
		return nil, fmt.Errorf("abi: unexpected signature type for error, got %s", errsig.Kind)
	}
	return p.toError(errsig)
}

// toType converts a sigparser.Parameter to a Type.
func (p *Parser) toType(param sigparser.Parameter) (typ Type, err error) {
	if len(param.Type) > 0 && len(param.Tuple) > 0 {
		return nil, fmt.Errorf("abi: parameter cannot be both elementary type and tuple: %s", param)
	}
	switch {
	case len(param.Arrays) > 0:
		// The sigparser package return array size in the Arrays field. If the
		// array has multiple dimensions, the size of consecutive dimensions is
		// stored in the Arrays, e.g. for a [2][3] array, the Arrays field
		// contains [2, 3]. Unbounded arrays have a size of -1. We need to
		// convert this to a nested structure of ArrayType and FixedArrayType
		// types.
		cpy := copyParam(param)
		cpy.Arrays = nil
		typ, err = p.toType(cpy)
		if err != nil {
			return nil, err
		}
		for i := len(param.Arrays) - 1; i >= 0; i-- {
			if param.Arrays[i] == -1 {
				typ = NewArrayType(typ)
			} else {
				typ = NewFixedArrayType(typ, param.Arrays[i])
			}
		}
	case len(param.Tuple) > 0:
		// If a parameter is a tuple, we need to convert all its elements
		// recursively.
		tuple := make([]TupleTypeElem, len(param.Tuple))
		for i, param := range param.Tuple {
			elemTyp, err := p.toType(param)
			if err != nil {
				return nil, err
			}
			tuple[i] = TupleTypeElem{
				Name: param.Name,
				Type: elemTyp,
			}
			if err != nil {
				return nil, err
			}
		}
		typ = NewTupleType(tuple...)
	default:
		// If the parameter is not a tuple or array, we look up the type in the
		// Config struct.
		typ = p.Config.Types[param.Type]
	}
	if typ == nil {
		return nil, fmt.Errorf("abi: unknown type %q", param.Type)
	}
	return typ, nil
}

// toConstructor converts a sigparser.Signature to a Constructor.
func (p *Parser) toConstructor(s sigparser.Signature) (*Constructor, error) {
	var inputs []TupleTypeElem
	for _, param := range s.Inputs {
		typ, err := p.toType(param)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, TupleTypeElem{
			Name: param.Name,
			Type: typ,
		})
	}
	return NewConstructorWithConfig(NewTupleType(inputs...), p.Config), nil
}

// toMethod converts a sigparser.Signature to a Method.
func (p *Parser) toMethod(s sigparser.Signature) (*Method, error) {
	var (
		inputs  []TupleTypeElem
		outputs []TupleTypeElem
	)
	for _, param := range s.Inputs {
		typ, err := p.toType(param)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, TupleTypeElem{
			Name: param.Name,
			Type: typ,
		})
	}
	for _, param := range s.Outputs {
		typ, err := p.toType(param)
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, TupleTypeElem{
			Name: param.Name,
			Type: typ,
		})
	}
	return NewMethodWithConfig(s.Name, NewTupleType(inputs...), NewTupleType(outputs...), p.Config), nil
}

func (p *Parser) toEvent(s sigparser.Signature) (*Event, error) {
	var inputs []EventTupleTypeElem
	if len(s.Inputs) == 0 {
		return nil, fmt.Errorf("abi: event %q has no inputs", s.Name)
	}
	for _, param := range s.Inputs {
		typ, err := p.toType(param)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, EventTupleTypeElem{
			Name:    param.Name,
			Indexed: param.Indexed,
			Type:    typ,
		})
	}
	return NewEventWithConfig(s.Name, NewEventTupleType(inputs...), p.Config), nil
}

func (p *Parser) toError(s sigparser.Signature) (*Error, error) {
	var inputs []TupleTypeElem
	if len(s.Inputs) == 0 {
		return nil, fmt.Errorf("abi: event %q has no inputs", s.Name)
	}
	for _, param := range s.Inputs {
		typ, err := p.toType(param)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, TupleTypeElem{
			Name: param.Name,
			Type: typ,
		})
	}
	return NewErrorWithConfig(s.Name, NewTupleType(inputs...), p.Config), nil
}

func copyParam(param sigparser.Parameter) sigparser.Parameter {
	cpy := sigparser.Parameter{
		Type:         param.Type,
		Name:         param.Name,
		Arrays:       make([]int, len(param.Arrays)),
		Tuple:        make([]sigparser.Parameter, len(param.Tuple)),
		Indexed:      param.Indexed,
		DataLocation: param.DataLocation,
	}
	copy(cpy.Arrays, param.Arrays)
	for i, p := range param.Tuple {
		cpy.Tuple[i] = copyParam(p)
	}
	return cpy
}

func isKind(kind sigparser.SignatureKind, kinds ...sigparser.SignatureKind) bool {
	for _, k := range kinds {
		if kind == k {
			return true
		}
	}
	return false
}

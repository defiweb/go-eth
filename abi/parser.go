package abi

import (
	"fmt"

	"github.com/defiweb/go-sigparser"
)

func ParseType(signature string) (Type, error) {
	return NewParser(DefaultConfig).ParseType(signature)
}

func ParseMethod(signature string) (*Method, error) {
	return NewParser(DefaultConfig).ParseMethod(signature)
}

func ParseEvent(signature string) (*Event, error) {
	return NewParser(DefaultConfig).ParseEvent(signature)
}

func MustParseType(signature string) Type {
	typ, err := ParseType(signature)
	if err != nil {
		panic(err)
	}
	return typ
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

type Parser struct {
	Config *Config
}

func NewParser(c *Config) *Parser {
	return &Parser{Config: c}
}

func (c *Parser) ParseType(signature string) (Type, error) {
	p, err := sigparser.ParseParameter(signature)
	if err != nil {
		return nil, err
	}
	return c.toType(p)
}

func (c *Parser) ParseMethod(signature string) (*Method, error) {
	s, err := sigparser.ParseSignature(signature)
	if err != nil {
		return nil, err
	}
	return c.toMethod(s)
}

func (c *Parser) ParseEvent(signature string) (*Event, error) {
	s, err := sigparser.ParseSignature(signature)
	if err != nil {
		return nil, err
	}
	return c.toEvent(s)
}

// toMethod converts a sigparser.Signature to a Method.
func (c *Parser) toMethod(s sigparser.Signature) (*Method, error) {
	var (
		inputs  []TupleTypeElem
		outputs []TupleTypeElem
	)
	for _, p := range s.Inputs {
		typ, err := c.toType(p)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, TupleTypeElem{
			Name: p.Name,
			Type: typ,
		})
	}
	for _, p := range s.Outputs {
		typ, err := c.toType(p)
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, TupleTypeElem{
			Name: p.Name,
			Type: typ,
		})
	}
	return &Method{
		Name:    s.Name,
		Inputs:  NewTupleType(inputs...),
		Outputs: NewTupleType(outputs...),
		Config:  c.Config,
	}, nil
}

func (c *Parser) toEvent(s sigparser.Signature) (*Event, error) {
	var inputs []EventTupleTypeElem
	if len(s.Inputs) == 0 {
		return nil, fmt.Errorf("abi: event %q has no inputs", s.Name)
	}
	for _, p := range s.Inputs {
		typ, err := c.toType(p)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, EventTupleTypeElem{
			Name:    p.Name,
			Indexed: p.Indexed,
			Type:    typ,
		})
	}
	return &Event{
		Name:   s.Name,
		Inputs: NewEventTupleType(inputs...),
		Config: c.Config,
	}, nil
}

// toType converts a sigparser.Parameter to a Type.
func (c *Parser) toType(p sigparser.Parameter) (typ Type, err error) {
	if len(p.Type) > 0 && len(p.Tuple) > 0 {
		return nil, fmt.Errorf("abi: parameter cannot be both elementary type and tuple: %s", p)
	}
	switch {
	case len(p.Arrays) > 0:
		// The sigparser package return array size in the Arrays field. If the
		// array has multiple dimensions, the size of consecutive dimensions is
		// stored in the Arrays, e.g. for a [2][3] array, the Arrays field
		// contains [2, 3]. Unbounded arrays have a size of -1. We need to
		// convert this to a nested structure of ArrayType and FixedArrayType
		// types.
		cpy := copyParam(p)
		cpy.Arrays = nil
		typ, err = c.toType(cpy)
		if err != nil {
			return nil, err
		}
		for i := len(p.Arrays) - 1; i >= 0; i-- {
			if p.Arrays[i] == -1 {
				typ = NewArrayType(typ)
			} else {
				typ = NewFixedArrayType(typ, p.Arrays[i])
			}
		}
	case len(p.Tuple) > 0:
		// If a parameter is a tuple, we need to convert all its elements
		// recursively.
		tuple := make([]TupleTypeElem, len(p.Tuple))
		for i, p := range p.Tuple {
			elemTyp, err := c.toType(p)
			if err != nil {
				return nil, err
			}
			tuple[i] = TupleTypeElem{
				Name: p.Name,
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
		typ = c.Config.Types[p.Type]
	}
	if typ == nil {
		return nil, fmt.Errorf("abi: unknown type %q", p.Type)
	}
	return typ, nil
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

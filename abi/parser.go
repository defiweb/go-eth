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

func (c *Parser) toMethod(s sigparser.Signature) (*Method, error) {
	var (
		err     error
		inputs  Type
		outputs Type
	)
	if len(s.Inputs) > 0 {
		inputs, err = c.toType(sigparser.Parameter{Tuple: s.Inputs})
		if err != nil {
			return nil, err
		}
	}
	if len(s.Outputs) > 0 {
		outputs, err = c.toType(sigparser.Parameter{Tuple: s.Outputs})
		if err != nil {
			return nil, err
		}
	}
	if inputs == nil {
		inputs = NewTupleType()
	}
	if outputs == nil {
		outputs = NewTupleType()
	}
	return &Method{
		Name:    s.Name,
		Inputs:  inputs.(*TupleType),
		Outputs: outputs.(*TupleType),
		Config:  c.Config,
	}, nil
}

func (c *Parser) toType(p sigparser.Parameter) (typ Type, err error) {
	if len(p.Type) > 0 && len(p.Tuple) > 0 {
		return nil, fmt.Errorf("abi: parameter cannot be both elementary type and tuple: %s", p)
	}
	switch {
	case len(p.Arrays) > 0:
		cpy := copyParam(p)
		cpy.Arrays = cpy.Arrays[1:]
		cpyTyp, err := c.toType(cpy)
		if err != nil {
			return nil, err
		}
		if p.Arrays[0] < 0 {
			typ = NewArrayType(cpyTyp)
		} else {
			typ = NewFixedArrayType(cpyTyp, p.Arrays[0])
		}
	case len(p.Tuple) > 0:
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

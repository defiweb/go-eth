package abi

import (
	"fmt"

	"github.com/defiweb/go-sigparser"
)

// ParseType parses a type signature and returns a new Type.
//
// A type can be either an elementary type like uint256 or a tuple type. Tuple
// types are denoted by parentheses, with the optional keyword "tuple" before
// the parentheses. Parameter names are optional.
//
// The generated types can be used to create new values, which can then be used
// to encode or decode ABI data.
//
// Custom types may be added to the Config.Types. This will allow the parser to
// handle custom types.
//
// The following examples are valid type signatures:
//
//   uint256
//   (uint256 a,bytes32 b)
//   tuple(uint256 a, bytes32 b)[]
//
// This function is equivalent to calling Parser.ParseType with the default
// configuration.
func ParseType(signature string) (Type, error) {
	return NewParser(DefaultConfig).ParseType(signature)
}

// ParseConstructor parses a constructor signature and returns a new Constructor.
//
// A constructor signature is similar to a method signature, but it does not
// have a name and returns no values. It can be optionally prefixed with the
// "constructor" keyword.
//
// The following examples are valid signatures:
//
//   ((uint256,bytes32)[])
//   ((uint256 a, bytes32 b)[] c)
//   constructor(tuple(uint256 a, bytes32 b)[] memory c)
//
// This function is equivalent to calling Parser.ParseConstructor with the
// default configuration.
func ParseConstructor(signature string) (*Constructor, error) {
	return NewParser(DefaultConfig).ParseConstructor(signature)
}

// ParseMethod parses a method signature and returns a new Method.
//
// The method accepts Solidity method signatures, but allows to omit the
// "function" keyword, argument names and the "returns" keyword. Method
// modifiers and argument data location specifiers are allowed, but ignored.
// Tuple types are indicated by parentheses, with the optional keyword "tuple"
// before the parentheses.
//
// The following examples are valid signatures:
//
//   foo((uint256,bytes32)[])(uint256)
//   foo((uint256 a, bytes32 b)[] c)(uint256 d)
//   function foo(tuple(uint256 a, bytes32 b)[] memory c) pure returns (uint256 d)
//
// This function is equivalent to calling Parser.ParseMethod with the default
// configuration.
func ParseMethod(signature string) (*Method, error) {
	return NewParser(DefaultConfig).ParseMethod(signature)
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
	return NewParser(DefaultConfig).ParseEvent(signature)
}

// ParseError parses an error signature and returns a new Error.
//
// An error signature is similar to a method signature, but returns no values.
// It can be optionally prefixed with the "error" keyword.
//
// The following examples are valid signatures:
//
//   foo((uint256,bytes32)[])(uint256)
//   foo((uint256 a, bytes32 b)[] c)(uint256 d)
//   error foo(tuple(uint256 a, bytes32 b)[] memory c)
//
// This function is equivalent to calling Parser.ParseError with the default
// configuration.
func ParseError(signature string) (*Error, error) {
	return NewParser(DefaultConfig).ParseError(signature)
}

// MustParseType is like ParseType but panics on error.
func MustParseType(signature string) Type {
	t, err := ParseType(signature)
	if err != nil {
		panic(err)
	}
	return t
}

// MustParseConstructor is like ParseConstructor but panics on error.
func MustParseConstructor(signature string) *Constructor {
	c, err := ParseConstructor(signature)
	if err != nil {
		panic(err)
	}
	return c
}

// MustParseMethod is like ParseMethod but panics on error.
func MustParseMethod(signature string) *Method {
	m, err := ParseMethod(signature)
	if err != nil {
		panic(err)
	}
	return m
}

// MustParseEvent is like ParseEvent but panics on error.
func MustParseEvent(signature string) *Event {
	e, err := ParseEvent(signature)
	if err != nil {
		panic(err)
	}
	return e
}

// MustParseError is like ParseError but panics on error.
func MustParseError(signature string) *Error {
	e, err := ParseError(signature)
	if err != nil {
		panic(err)
	}
	return e
}

// Parser parses method, constructor, event, error and type signatures.
type Parser struct {
	Config *Config
}

// NewParser returns a new Parser with the given configuration.
func NewParser(c *Config) *Parser {
	return &Parser{Config: c}
}

// ParseType parses a type signature and returns a new Type.
//
// See ParseType for more information.
func (p *Parser) ParseType(signature string) (Type, error) {
	typ, err := sigparser.ParseParameter(signature)
	if err != nil {
		return nil, err
	}
	return p.toType(typ)
}

// ParseConstructor parses a constructor signature and returns a new Constructor.
//
// See ParseConstructor for more information.
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

// ParseMethod parses a method signature and returns a new Method.
//
// See ParseMethod for more information.
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

// ParseEvent parses an event signature and returns a new Event.
//
// See ParseEvent for more information.
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

// ParseError parses an error signature and returns a new Error.
//
// See ParseError for more information.
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

func (p *Parser) toType(param sigparser.Parameter) (typ Type, err error) {
	if len(param.Type) > 0 && len(param.Tuple) > 0 {
		return nil, fmt.Errorf("abi: parameter cannot be both elementary type and tuple: %s", param)
	}
	switch {
	case len(param.Arrays) > 0:
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
		typ = p.Config.Types[param.Type]
	}
	if typ == nil {
		return nil, fmt.Errorf("abi: unknown type %q", param.Type)
	}
	return typ, nil
}

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
	var inputs []EventTupleElem
	if len(s.Inputs) == 0 {
		return nil, fmt.Errorf("abi: event %q has no inputs", s.Name)
	}
	for _, param := range s.Inputs {
		typ, err := p.toType(param)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, EventTupleElem{
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

package abi

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/defiweb/go-sigparser"
)

// Contract provides a high-level API for interacting with a contract. It can
// be created from a JSON ABI definition using the ParseJSON function or from
// a list of signatures using the ParseSignatures function.
type Contract struct {
	Constructor        *Constructor
	Methods            map[string]*Method
	MethodsBySignature map[string]*Method
	Events             map[string]*Event
	Errors             map[string]*Error
	Types              map[string]Type // Types defined in the ABI (structs, enums and user-defined Value Types)
}

// IsError returns true if the given error data, returned by a contract call,
// corresponds to a revert, panic, or custom error.
func (c *Contract) IsError(data []byte) bool {
	if IsRevert(data) || IsPanic(data) {
		return true
	}
	for _, err := range c.Errors {
		if err.Is(data) {
			return true
		}
	}
	return false
}

// ToError returns an error if the given error data, returned by a contract
// call, corresponds to a revert, panic, or a custom error. It returns nil if
// the data cannot be recognized as an error.
func (c *Contract) ToError(data []byte) error {
	if IsRevert(data) {
		return RevertError{Reason: DecodeRevert(data)}
	}
	if IsPanic(data) {
		return PanicError{Code: DecodePanic(data)}
	}
	for _, err := range c.Errors {
		if err.Is(data) {
			return CustomError{Type: err, Data: data}
		}
	}
	return nil
}

// HandleError converts an error returned by a contract call to a RevertError,
// PanicError, or CustomError if applicable. If not, it returns the original
// error.
func (c *Contract) HandleError(err error) error {
	if err == nil {
		return nil
	}
	var dataErr interface{ RPCErrorData() any }
	if !errors.As(err, &dataErr) {
		return err
	}
	data, ok := dataErr.RPCErrorData().([]byte)
	if !ok {
		return err
	}
	if err := c.ToError(data); err != nil {
		return err
	}
	return err
}

// RegisterTypes registers types defined in the contract to the given ABI
// instance. This enables the use of types defined in the contract in all
// Parse* methods.
//
// If the type name already exists, it will be overwritten.
func (c *Contract) RegisterTypes(a *ABI) {
	for n, t := range c.Types {
		a.Types[n] = t
	}
}

// LoadJSON loads the ABI from the given JSON file and returns a Contract
// instance.
func LoadJSON(path string) (*Contract, error) {
	return Default.LoadJSON(path)
}

// MustLoadJSON is like LoadJSON but panics on error.
func MustLoadJSON(path string) *Contract {
	return Default.MustLoadJSON(path)
}

// ParseJSON parses the given ABI JSON and returns a Contract instance.
func ParseJSON(data []byte) (*Contract, error) {
	return Default.ParseJSON(data)
}

// MustParseJSON is like ParseJSON but panics on error.
func MustParseJSON(data []byte) *Contract {
	return Default.MustParseJSON(data)
}

// ParseSignatures parses list of signatures and returns a Contract instance.
// Signatures must be prefixed with the kind, e.g. "function" or "event".
//
// It accepts signatures in the same format as ParseConstructor, ParseMethod,
// ParseEvent, and ParseError functions.
func ParseSignatures(signatures ...string) (*Contract, error) {
	return Default.ParseSignatures(signatures...)
}

// MustParseSignatures is like ParseSignatures but panics on error.
func MustParseSignatures(signatures ...string) *Contract {
	return Default.MustParseSignatures(signatures...)
}

// LoadJSON loads the ABI from the given JSON file and returns a Contract
// instance.
func (a *ABI) LoadJSON(path string) (*Contract, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return a.ParseJSON(data)
}

// MustLoadJSON is like LoadJSON but panics on error.
func (a *ABI) MustLoadJSON(path string) *Contract {
	c, err := a.LoadJSON(path)
	if err != nil {
		panic(err)
	}
	return c
}

// ParseJSON parses the given ABI JSON and returns a Contract instance.
func (a *ABI) ParseJSON(data []byte) (*Contract, error) {
	var fields []jsonField
	if err := json.Unmarshal(data, &fields); err != nil {
		return nil, err
	}
	c := &Contract{
		Methods:            make(map[string]*Method),
		MethodsBySignature: make(map[string]*Method),
		Events:             make(map[string]*Event),
		Errors:             make(map[string]*Error),
		Types:              make(map[string]Type),
	}
	for _, f := range fields {
		inputs, err := f.Inputs.toTypes(a)
		if err != nil {
			return nil, err
		}
		outputs, err := f.Outputs.toTypes(a)
		if err != nil {
			return nil, err
		}
		for _, k := range []int{kindType, kindStruct, kindEnum} {
			for n, t := range inputs.internalTypes(k) {
				c.Types[n] = t
			}
			for n, t := range outputs.internalTypes(k) {
				c.Types[n] = t
			}
		}
		switch f.Type {
		case "constructor":
			c.Constructor = a.NewConstructor(inputs.toTupleType())
		case "function", "":
			method := a.NewMethod(
				f.Name,
				inputs.toTupleType(),
				outputs.toTupleType(),
				StateMutabilityFromString(f.StateMutability),
			)
			c.Methods[f.Name] = method
			c.MethodsBySignature[method.Signature()] = method
		case "event":
			c.Events[f.Name] = a.NewEvent(f.Name, inputs.toEventTupleType(), f.Anonymous)
		case "error":
			c.Errors[f.Name] = a.NewError(f.Name, inputs.toTupleType())
		case "fallback":
		case "receive":
		default:
			return nil, fmt.Errorf("unknown type: %s", f.Type)
		}
	}
	return c, nil
}

// MustParseJSON is like ParseJSON but panics on error.
func (a *ABI) MustParseJSON(data []byte) *Contract {
	c, err := a.ParseJSON(data)
	if err != nil {
		panic(err)
	}
	return c
}

// ParseSignatures parses list of signatures and returns a Contract instance.
// Signatures must be prefixed with the kind, e.g. "constructor" or "event".
// For functions, the "function" prefix can be omitted.
func (a *ABI) ParseSignatures(signatures ...string) (*Contract, error) {
	c := &Contract{
		Methods:            make(map[string]*Method),
		MethodsBySignature: make(map[string]*Method),
		Events:             make(map[string]*Event),
		Errors:             make(map[string]*Error),
		Types:              make(map[string]Type),
	}
	extraTypes := map[string]Type{}
	for _, s := range signatures {
		switch sigparser.Kind(s) {
		case sigparser.StructDefinitionInput:
			typ, err := sigparser.ParseStruct(s)
			if err != nil {
				return nil, err
			}
			if typ.Name == "" {
				return nil, errors.New("struct must have a name")
			}
			alias, err := newTypeFromSig(a, extraTypes, typ)
			if err != nil {
				return nil, err
			}
			alias = NewAliasType(typ.Name, alias)
			c.Types[typ.Name] = alias
			extraTypes[typ.Name] = alias
		case sigparser.TupleInput, sigparser.TypeInput, sigparser.ArrayInput:
			typ, err := sigparser.ParseParameter(s)
			if err != nil {
				return nil, err
			}
			if typ.Name == "" {
				return nil, errors.New("type must have a name")
			}
			alias, err := newTypeFromSig(a, extraTypes, typ)
			if err != nil {
				return nil, err
			}
			alias = NewAliasType(typ.Name, alias)
			c.Types[typ.Name] = alias
			extraTypes[typ.Name] = alias
		case sigparser.ConstructorSignatureInput:
			sig, err := sigparser.ParseSignatureAs(sigparser.ConstructorKind, s)
			if err != nil {
				return nil, err
			}
			constructor, err := newConstructorFromSig(a, extraTypes, sig)
			if err != nil {
				return nil, err
			}
			c.Constructor = constructor
		case sigparser.FunctionSignatureInput:
			sig, err := sigparser.ParseSignatureAs(sigparser.FunctionKind, s)
			if err != nil {
				return nil, err
			}
			method, err := newMethodFromSig(a, extraTypes, sig)
			if err != nil {
				return nil, err
			}
			c.Methods[method.Name()] = method
			c.MethodsBySignature[method.Signature()] = method
		case sigparser.EventSignatureInput:
			sig, err := sigparser.ParseSignatureAs(sigparser.EventKind, s)
			if err != nil {
				return nil, err
			}
			event, err := newEventFromSig(a, extraTypes, sig)
			if err != nil {
				return nil, err
			}
			c.Events[event.Name()] = event
		case sigparser.ErrorSignatureInput:
			sig, err := sigparser.ParseSignatureAs(sigparser.ErrorKind, s)
			if err != nil {
				return nil, err
			}
			errsig, err := newErrorFromSig(a, extraTypes, sig)
			if err != nil {
				return nil, err
			}
			c.Errors[errsig.Name()] = errsig
		default:
			return nil, fmt.Errorf("invalid signature: %s", s)
		}
	}
	return c, nil
}

// MustParseSignatures is like ParseSignatures but panics on error.
func (a *ABI) MustParseSignatures(signatures ...string) *Contract {
	c, err := a.ParseSignatures(signatures...)
	if err != nil {
		panic(err)
	}
	return c
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

func (a jsonParameters) toTypes(abi *ABI) (types jsonABITypes, err error) {
	types = make(jsonABITypes, len(a))
	for i, p := range a {
		var typ jsonABIType
		if typ, err = p.toType(abi); err != nil {
			return
		}
		types[i] = typ
	}
	return
}

type jsonParameter struct {
	Name         string         `json:"name"`
	Type         string         `json:"type"`
	InternalType string         `json:"internalType"`
	Indexed      bool           `json:"indexed"`
	Components   jsonParameters `json:"components"`
}

// toType converts a JSON parameter to a parsed type.
//
// It also returns a map of internal types defined in the parameter.
func (a jsonParameter) toType(abi *ABI) (typ jsonABIType, err error) {
	baseTyp, arrays, err := parseArrays(a.Type)
	if err != nil {
		return
	}
	var (
		kind         int
		intName      string
		intNamespace string
	)

	// If the internal type is different from type, it means that the type
	// is a custom type defined in the contract code. We can extract the name
	// and use it as an alias. As a result, when printing signatures, the
	// name from the contract code will be used instead of the type.
	if len(a.InternalType) > 0 && a.Type != a.InternalType {
		kind, intName, intNamespace = parseInternalType(a.InternalType)
	}

	typ.name = a.Name
	typ.indexed = a.Indexed
	typ.internalTypes = make(map[string]jsonABIType)

	switch {
	case len(arrays) > 0:
		elemType := jsonParameter{
			Name:         a.Name,
			Type:         baseTyp,
			InternalType: baseTyp,
			Components:   a.Components,
		}
		if typ, err = elemType.toType(abi); err != nil {
			return
		}
		if len(intName) > 0 {
			typ.typ = NewAliasType(intName, typ.typ)
		}
		typ.elemTyp = typ.typ
		for i := len(arrays) - 1; i >= 0; i-- {
			if arrays[i] == -1 {
				typ.typ = NewArrayType(typ.typ)
			} else {
				typ.typ = NewFixedArrayType(typ.typ, arrays[i])
			}
		}
	case len(a.Components) > 0:
		tuple := make([]TupleTypeElem, len(a.Components))
		for i, jsonComp := range a.Components {
			var compTyp jsonABIType
			compTyp, err = jsonComp.toType(abi)
			if err != nil {
				return
			}
			tuple[i].Type = compTyp.typ
			tuple[i].Name = compTyp.name
			for intName, intTyp := range typ.internalTypes {
				typ.internalTypes[intName] = intTyp
			}
		}
		typ.typ = NewTupleType(tuple...)
		if len(intName) > 0 {
			typ.typ = NewAliasType(intName, typ.typ)
		}
	default:
		typ.typ = abi.Types[baseTyp]
		if typ.typ == nil {
			return jsonABIType{}, fmt.Errorf("abi: unknown type %q", a.Type)
		}
		if len(intName) > 0 {
			typ.typ = NewAliasType(intName, typ.typ)
		}
	}
	if typ.elemTyp == nil {
		typ.elemTyp = typ.typ
	}
	if kind != kindNone {
		typ.internalTypes[a.InternalType] = jsonABIType{
			typ:       typ.elemTyp,
			elemTyp:   typ.elemTyp,
			name:      intName,
			namespace: intNamespace,
			kind:      kind,
		}
	}
	return
}

// jsonABIType represents a ABI type extracted from a JSON ABI.
type jsonABIType struct {
	typ           Type                   // Type of the parameter.
	elemTyp       Type                   // In case of arrays, the type of the elements. Otherwise, the same as typ.
	name          string                 // Name of the parameter.
	namespace     string                 // Namespace of the parameter.
	indexed       bool                   // Whether the parameter is indexed.
	kind          int                    // Kind of the internal parameter, e.g. (enum, struct, type).
	internalTypes map[string]jsonABIType // All internal types defined in the parameter.
}

type jsonABITypes []jsonABIType

func (a jsonABITypes) internalTypes(kind int) map[string]Type {
	internalTypes := make(map[string]Type)
	for _, param := range a {
		for _, v := range param.internalTypes {
			if v.kind != kind {
				continue
			}
			internalTypes[v.name] = v.elemTyp
		}
	}
	return internalTypes
}

// toTupleType converts parameters to a TupleType type.
func (a jsonABITypes) toTupleType() *TupleType {
	tupleElements := make([]TupleTypeElem, len(a))
	for i, param := range a {
		tupleElements[i] = TupleTypeElem{
			Name: param.name,
			Type: param.typ,
		}
	}
	return NewTupleType(tupleElements...)
}

// toEventTupleType converts parameters to a EventTupleType type.
func (a jsonABITypes) toEventTupleType() *EventTupleType {
	tupleElements := make([]EventTupleElem, len(a))
	for i, param := range a {
		tupleElements[i] = EventTupleElem{
			Name:    param.name,
			Indexed: param.indexed,
			Type:    param.typ,
		}
	}
	return NewEventTupleType(tupleElements...)
}

// parseArray parses type name and returns the name and array dimensions.
// For example, "uint256[][3]" will return "uint256" and [-1, 3].
// For unbounded arrays, the dimension is -1.
func parseArrays(typ string) (baseTyp string, arrays []int, err error) {
	openBracket := strings.Index(typ, "[")
	if openBracket == -1 {
		baseTyp = typ
		return
	}
	baseTyp = typ[:openBracket]
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

const (
	kindNone = iota
	kindEnum
	kindStruct
	kindType
)

func parseInternalType(typ string) (int, string, string) {
	if len(typ) == 0 {
		return kindNone, "", ""
	}
	var (
		kind      int
		prefixLen int
	)
	switch {
	case strings.Index(typ, "struct ") == 0:
		kind, prefixLen = kindStruct, 7
	case strings.Index(typ, "enum ") == 0:
		kind, prefixLen = kindEnum, 5
	case !strings.Contains(typ, " "):
		kind, prefixLen = kindType, 0
	default:
		return kindNone, "", ""
	}
	intName := typ[prefixLen:]
	intNamespace := ""
	if bracket := strings.Index(intName, "["); bracket != -1 {
		intName = intName[:bracket]
	}
	if parts := strings.SplitN(intName, ".", 2); len(parts) == 2 {
		intName, intNamespace = parts[1], parts[0]
	}
	return kind, intName, intNamespace
}

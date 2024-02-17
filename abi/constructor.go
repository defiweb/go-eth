package abi

// Constructor represents a constructor in an Contract. The constructor can be used to
// encode arguments for a constructor call.
type Constructor struct {
	inputs *TupleType
	abi    *ABI
}

// NewConstructor creates a new Constructor instance.
//
// This method is rarely used, see ParseConstructor for a more convenient way
// to create a new Constructor.
func NewConstructor(inputs *TupleType) *Constructor {
	return Default.NewConstructor(inputs)
}

// ParseConstructor parses a constructor signature and returns a new Constructor.
//
// A constructor signature is similar to a method signature, but it does not
// have a name and returns no values. It can be optionally prefixed with the
// "constructor" keyword.
//
// The following examples are valid signatures:
//
//	((uint256,bytes32)[])
//	((uint256 a, bytes32 b)[] c)
//	constructor(tuple(uint256 a, bytes32 b)[] memory c)
//
// This function is equivalent to calling Parser.ParseConstructor with the
// default configuration.
func ParseConstructor(signature string) (*Constructor, error) {
	return Default.ParseConstructor(signature)
}

// MustParseConstructor is like ParseConstructor but panics on error.
func MustParseConstructor(signature string) *Constructor {
	return Default.MustParseConstructor(signature)
}

// NewConstructor creates a new Constructor instance.
func (a *ABI) NewConstructor(inputs *TupleType) *Constructor {
	if inputs == nil {
		inputs = NewTupleType()
	}
	return &Constructor{
		inputs: inputs,
		abi:    a,
	}
}

// ParseConstructor parses a constructor signature and returns a new Constructor.
//
// See ParseConstructor for more information.
func (a *ABI) ParseConstructor(signature string) (*Constructor, error) {
	return parseConstructor(a, nil, signature)
}

// MustParseConstructor is like ParseConstructor but panics on error.
func (a *ABI) MustParseConstructor(signature string) *Constructor {
	c, err := a.ParseConstructor(signature)
	if err != nil {
		panic(err)
	}
	return c
}

// Inputs returns the input arguments of the constructor as a tuple type.
func (m *Constructor) Inputs() *TupleType {
	return m.inputs
}

// EncodeArg encodes an argument for a contract deployment.
// The map or structure must have fields with the same names as the
// constructor arguments.
func (m *Constructor) EncodeArg(code []byte, arg any) ([]byte, error) {
	encoded, err := m.abi.EncodeValue(m.inputs, arg)
	if err != nil {
		return nil, err
	}
	input := make([]byte, len(code)+len(encoded))
	copy(input, code)
	copy(input[len(code):], encoded)
	return input, nil
}

// MustEncodeArg is like EncodeArg but panics on error.
func (m *Constructor) MustEncodeArg(code []byte, arg any) []byte {
	encoded, err := m.EncodeArg(code, arg)
	if err != nil {
		panic(err)
	}
	return encoded
}

// EncodeArgs encodes arguments for a contract deployment.
// The map or structure must have fields with the same names as the
// constructor arguments.
func (m *Constructor) EncodeArgs(code []byte, args ...any) ([]byte, error) {
	encoded, err := m.abi.EncodeValues(m.inputs, args...)
	if err != nil {
		return nil, err
	}
	input := make([]byte, len(code)+len(encoded))
	copy(input, code)
	copy(input[len(code):], encoded)
	return input, nil
}

// MustEncodeArgs is like EncodeArgs but panics on error.
func (m *Constructor) MustEncodeArgs(code []byte, args ...any) []byte {
	encoded, err := m.EncodeArgs(code, args...)
	if err != nil {
		panic(err)
	}
	return encoded
}

// String returns the human-readable signature of the constructor.
func (m *Constructor) String() string {
	return "constructor" + m.inputs.String()
}

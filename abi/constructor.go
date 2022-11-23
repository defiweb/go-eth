package abi

// Constructor represents a constructor in an ABI. The constructor can be used to
// encode arguments for a constructor call.
type Constructor struct {
	inputs *TupleType
	config *Config
}

// NewConstructor creates a new Constructor instance.
func NewConstructor(inputs *TupleType) *Constructor {
	return NewConstructorWithConfig(inputs, DefaultConfig)
}

// NewConstructorWithConfig creates a new Constructor instance with a custom
// config.
func NewConstructorWithConfig(inputs *TupleType, config *Config) *Constructor {
	return &Constructor{
		inputs: inputs,
		config: config,
	}
}

// Inputs returns the input arguments of the constructor as a tuple type.
func (m *Constructor) Inputs() *TupleType {
	return m.inputs
}

// EncodeArg encodes arguments for a constructor call using a provided map or
// structure. The map or structure must have fields with the same names as
// the constructor arguments.
func (m *Constructor) EncodeArg(val any) ([]byte, error) {
	encoded, err := NewEncoder(m.config).EncodeValue(m.inputs.Value(), val)
	if err != nil {
		return nil, err
	}
	return encoded, nil
}

// EncodeArgs encodes arguments for a constructor call.
func (m *Constructor) EncodeArgs(args ...any) ([]byte, error) {
	encoded, err := NewEncoder(m.config).EncodeValues(m.inputs.Value().(*TupleValue), args...)
	if err != nil {
		return nil, err
	}
	return encoded, nil
}

// String returns the human-readable signature of the constructor.
func (m *Constructor) String() string {
	return "constructor" + m.inputs.Type()
}

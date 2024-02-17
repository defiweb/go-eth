package main

import (
	"fmt"

	"github.com/defiweb/go-eth/abi"
	"github.com/defiweb/go-eth/hexutil"
)

// BoolFlagsType is a custom type that represents a 256-bit bitfield.
//
// It must implement the abi.Type interface.
type BoolFlagsType struct{}

// IsDynamic returns true if the type is dynamic-length, like string or bytes.
func (b BoolFlagsType) IsDynamic() bool {
	return false
}

// CanonicalType is the type as it would appear in the ABI.
// It must only use the types defined in the ABI specification:
// https://docs.soliditylang.org/en/latest/abi-spec.html
func (b BoolFlagsType) CanonicalType() string {
	return "bytes32"
}

// String returns the custom type name.
func (b BoolFlagsType) String() string {
	return "BoolFlags"
}

// Value returns the zero value for this type.
func (b BoolFlagsType) Value() abi.Value {
	return &BoolFlagsValue{}
}

// BoolFlagsValue is the value of the custom type.
//
// It must implement the abi.Value interface.
type BoolFlagsValue [256]bool

// IsDynamic returns true if the type is dynamic-length, like string or bytes.
func (b BoolFlagsValue) IsDynamic() bool {
	return false
}

// EncodeABI encodes the value to the ABI format.
func (b BoolFlagsValue) EncodeABI() (abi.Words, error) {
	var w abi.Word
	for i, v := range b {
		if v {
			w[i/8] |= 1 << uint(i%8)
		}
	}
	return abi.Words{w}, nil
}

// DecodeABI decodes the value from the ABI format.
func (b *BoolFlagsValue) DecodeABI(words abi.Words) (int, error) {
	if len(words) == 0 {
		return 0, fmt.Errorf("abi: cannot decode BytesFlags from empty data")
	}
	for i, v := range words[0] {
		for j := 0; j < 8; j++ {
			b[i*8+j] = v&(1<<uint(j)) != 0
		}
	}
	return 1, nil
}

// MapFrom and MapTo are optional methods that allow mapping between different
// types.
//
// The abi.Mapper is the instance of the internal mapper that is used to
// perform the mapping. It can be used to map nested types.
//
// Note, that you want to use reflection to implement following methods because
// it would allow to write more generic code mapping functions.

// MapFrom maps value from a different type.
func (b *BoolFlagsValue) MapFrom(_ abi.Mapper, src any) error {
	switch src := src.(type) {
	case [256]bool:
		*b = src
	case []bool:
		if len(src) > 256 {
			return fmt.Errorf("abi: cannot map []bool of length %d to BytesFlags", len(src))
		}
		for i, v := range src {
			b[i] = v
		}
	}
	return nil
}

// MapTo maps value to a different type.
func (b *BoolFlagsValue) MapTo(_ abi.Mapper, dst any) error {
	switch dst := dst.(type) {
	case *[256]bool:
		*dst = *b
	case *[]bool:
		*dst = make([]bool, 256)
		for i, v := range b {
			(*dst)[i] = v
		}
	}
	return nil
}

func main() {
	// Add custom type.
	abi.Default.Types["BoolFlags"] = &BoolFlagsType{}

	// Generate calldata.
	setFlags := abi.MustParseMethod("setFlags(BoolFlags flags)")
	calldata, _ := setFlags.EncodeArgs(
		[]bool{true, false, true, true, false, true, false, true},
	)

	// Print the calldata.
	fmt.Printf("Calldata: %s\n", hexutil.BytesToHex(calldata))
}

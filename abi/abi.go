package abi

import (
	"fmt"
	"unicode"

	"github.com/defiweb/go-anymapper"
)

// Default is the default ABI instance that is used by the package-level
// functions.
var Default *ABI

// ABI structure implements the Ethereum ABI (Application Binary Interface).
//
// It provides methods for working with ABI, such as parsing, encoding and
// decoding data.
//
// The package provides default ABI instance that is used by the package-level
// functions. It is possible to create custom ABI instances and use them
// instead of the default one. To do this, use the Copy method to create a copy
// of the default ABI instance and modify it as needed.
type ABI struct {
	// Types is a map of types that can be used in the ABU.
	// The key is the name of the type, and the value is the type.
	Types map[string]Type

	// Mapper is the instance of the mapper that will be used to map
	// values to and from Contract types.
	Mapper *anymapper.Mapper
}

// Copy returns a copy of the ABI instance.
func (a *ABI) Copy() *ABI {
	cpy := &ABI{
		Types:  make(map[string]Type, len(a.Types)),
		Mapper: a.Mapper.Copy(),
	}
	for k, v := range a.Types {
		cpy.Types[k] = v
	}
	return cpy
}

// fieldMapper lowercase the first letter of the field name. If the field name
// starts with an acronym, it will lowercase the whole acronym. For example:
//  - "User" will be mapped to "user"
// 	- "ID" will be mapped to "id"
// 	- "DAPPName" will be mapped to "dappName"
// Unfortunately, it does not work with field names that contain two acronyms
// next to each other. For example, "DAPPID" will be mapped to "dappid".
var fieldMapper = func(field string) string {
	if len(field) == 0 {
		return field
	}
	runes := []rune(field)
	for i, c := range runes {
		if unicode.IsUpper(c) && (i == 0 || i == len(runes)-1 || !(unicode.IsLower(runes[i+1]))) {
			runes[i] = unicode.ToLower(c)
		}
		if unicode.IsLower(c) {
			break
		}
	}
	return string(runes)
}

func init() {
	mapper := anymapper.DefaultMapper.Copy()
	mapper.Tag = "abi"
	mapper.FieldMapper = fieldMapper

	types := map[string]Type{}
	types["bool"] = NewBoolType()
	types["bytes"] = NewBytesType()
	types["string"] = NewStringType()
	types["address"] = NewAddressType()
	types["int"] = NewAliasType("int", NewIntType(256))
	types["uint"] = NewAliasType("uint", NewUintType(256))
	for i := 1; i <= 32; i++ {
		types[fmt.Sprintf("int%d", i*8)] = NewIntType(i * 8)
		types[fmt.Sprintf("uint%d", i*8)] = NewUintType(i * 8)
		types[fmt.Sprintf("bytes%d", i)] = NewFixedBytesType(i)
	}

	Default = &ABI{
		Types:  types,
		Mapper: mapper,
	}
}

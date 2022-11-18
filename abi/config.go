package abi

import (
	"fmt"
	"unicode"

	"github.com/defiweb/go-anymapper"
)

var DefaultConfig *Config

// Config holds the configuration for the ABI parser, encoder and decoder.
type Config struct {
	// Types is a map of types that can be used in the ABI.
	// The key is the name of the type, and the value is the type.
	Types map[string]Type

	// Mapper is the instance of the mapper that will be used to map
	// values to and from ABI types.
	Mapper *anymapper.Mapper
}

func (c *Config) Copy() *Config {
	cpy := &Config{
		Types:  make(map[string]Type, len(c.Types)),
		Mapper: c.Mapper.Copy(),
	}
	for k, v := range c.Types {
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
	types["int"] = NewAliasType("int", NewIntType(32))
	types["uint"] = NewAliasType("uint", NewUintType(32))
	for i := 1; i <= 32; i++ {
		types[fmt.Sprintf("int%d", i*8)] = NewIntType(i)
		types[fmt.Sprintf("uint%d", i*8)] = NewUintType(i)
		types[fmt.Sprintf("bytes%d", i)] = NewFixedBytesType(i)
	}

	DefaultConfig = &Config{
		Types:  types,
		Mapper: mapper,
	}
}

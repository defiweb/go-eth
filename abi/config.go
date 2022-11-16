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

func init() {
	mapper := anymapper.DefaultMapper.Copy()
	mapper.Tag = "abi"
	mapper.FieldMapper = func(field string) string {
		if len(field) > 0 {
			return string(unicode.ToLower(rune(field[0]))) + field[1:]
		}
		return field
	}

	types := map[string]Type{}
	types["bool"] = NewBoolType()
	types["int"] = NewIntType(32)
	types["uint"] = NewUintType(32)
	types["bytes"] = NewBytesType()
	types["string"] = NewStringType()
	types["address"] = NewAddressType()
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

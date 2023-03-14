package abi

import (
	"fmt"
	"reflect"
	"unicode"

	"github.com/defiweb/go-anymapper"
)

// Default is the default ABI instance that is used by the package-level
// functions.
//
// It is recommended to create a new ABI instance using NewABI rather than
// modifying the default instance, as this can potentially interfere with
// other packages that use the default ABI instance.
var Default = NewABI()

// ABI structure implements the Ethereum ABI (Application Binary Interface).
//
// It provides methods for working with ABI, such as parsing, encoding and
// decoding data.
//
// The package provides default ABI instance that is used by the package-level
// functions. It is possible to create custom ABI instances and use them
// instead of the default one.
type ABI struct {
	// Types is a map of known ABI types.
	// The key is the name of the type, and the value is the type.
	Types map[string]Type

	// Mapper is used to map values to and from ABI types.
	Mapper Mapper
}

// Mapper used to map values to and from ABI types.
type Mapper interface {
	Map(src any, dst any) error
}

// MapFrom maps the value from the ABI Value.
type MapFrom interface {
	MapFrom(m Mapper, src any) error
}

// MapTo maps the value to the ABI Value.
type MapTo interface {
	MapTo(m Mapper, dst any) error
}

// NewABI creates a new ABI instance.
//
// For most use cases, the default ABI instance should be used instead of
// creating a new one.
func NewABI() *ABI {
	mapper := anymapper.New()
	mapper.Tag = "abi"
	mapper.FieldMapper = fieldMapper

	// Those hooks add support for MapTo and MapFrom interfaces.
	//
	// Interfaces are used only if a source or destination type implements
	// the Value interface.
	//
	// If both types implement MapTo/MapFrom, then the method from the value
	// that does NOT implement Value interface is used. This is to ensure
	// that mapping functions defined by the user have higher priority.
	mapper.Hooks = anymapper.Hooks{
		MapFuncHook: func(m *anymapper.Mapper, src, dst reflect.Type) anymapper.MapFunc {
			srcImplMapTo := src.Implements(mapToTy)
			dstImplMapFrom := dst.Implements(mapFromTy)
			switch {
			case srcImplMapTo && dstImplMapFrom:
				if src.Implements(valueTy) {
					return func(m *anymapper.Mapper, src, dst reflect.Value) error {
						return dst.Interface().(MapFrom).MapFrom(m, src.Interface())
					}
				}
				if dst.Implements(valueTy) {
					return func(m *anymapper.Mapper, src, dst reflect.Value) error {
						return src.Interface().(MapTo).MapTo(m, addr(dst).Interface())
					}
				}
			case srcImplMapTo:
				return func(m *anymapper.Mapper, src, dst reflect.Value) error {
					return src.Interface().(MapTo).MapTo(m, addr(dst).Interface())
				}
			case dstImplMapFrom:
				return func(m *anymapper.Mapper, src, dst reflect.Value) error {
					return dst.Interface().(MapFrom).MapFrom(m, src.Interface())
				}
			}
			return nil
		},
		SourceValueHook: func(v reflect.Value) reflect.Value {
			for {
				if _, ok := v.Interface().(MapTo); ok {
					for v.Kind() == reflect.Interface {
						v = v.Elem()
					}
					return v
				}
				if v.Kind() != reflect.Interface && v.Kind() != reflect.Ptr {
					break
				}
				v = v.Elem()
			}
			return reflect.Value{}
		},
		DestinationValueHook: func(v reflect.Value) reflect.Value {
			for {
				// If the destination is a nil interface, then return it.
				if v.Kind() == reflect.Interface && v.IsNil() {
					return v
				}
				// If the destination is a nil pointer, then initialize it.
				if v.Kind() == reflect.Ptr && v.IsNil() {
					if !v.CanSet() {
						return reflect.Value{}
					}
					v.Set(reflect.New(v.Type().Elem()))
				}
				// If the destination implements MapFrom, then return it but
				// first dereference it if it is an interface.
				if _, ok := v.Interface().(MapFrom); ok {
					for v.Kind() == reflect.Interface {
						v = v.Elem()
					}
					return v
				}
				// If the destination is not a pointer or interface, then
				// break the loop and return an empty value. Returning an
				// empty value will cause the anymapper package to ignore
				// this hook.
				if v.Kind() != reflect.Interface && v.Kind() != reflect.Ptr {
					break
				}
				v = v.Elem()
			}
			return reflect.Value{}
		},
	}

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

	return &ABI{
		Types:  types,
		Mapper: mapper,
	}
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

func addr(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		return v
	}
	if v.CanAddr() {
		return v.Addr()
	}
	return v
}

var (
	valueTy   = reflect.TypeOf((*Value)(nil)).Elem()
	mapFromTy = reflect.TypeOf((*MapFrom)(nil)).Elem()
	mapToTy   = reflect.TypeOf((*MapTo)(nil)).Elem()
)

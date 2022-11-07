package abi

import (
	"fmt"
	"math"
	"math/big"
)

func DecodeValue(t Type, abi []byte, val any) error {
	return NewDecoder(DefaultConfig).DecodeValue(t, abi, val)
}

func DecodeValues(t *TupleType, abi []byte, vals ...any) error {
	return NewDecoder(DefaultConfig).DecodeValues(t, abi, vals...)
}

type Decoder struct {
	Config *Config
}

func NewDecoder(c *Config) *Decoder {
	return &Decoder{Config: c}
}

func (d *Decoder) DecodeValue(t Type, abi []byte, val any) error {
	if _, err := t.DecodeABI(BytesToWords(abi)); err != nil {
		return err
	}
	return d.Config.mapper.Map(t, val)
}

func (d *Decoder) DecodeValues(t *TupleType, abi []byte, vals ...any) error {
	if t.Length() != len(vals) {
		return fmt.Errorf("abi: cannot decode tuple, expected %d values, got %d", t.Length(), len(vals))
	}
	if _, err := t.DecodeABI(BytesToWords(abi)); err != nil {
		return err
	}
	for i, elem := range t.Elements() {
		if err := d.Config.mapper.Map(elem, vals[i]); err != nil {
			return err
		}
	}
	return nil
}

// decodeTuple decodes a tuple from the given words and stores the result in the
// given tuple. The tuple must contain the correct number of elements.
func decodeTuple(t *[]Type, w Words) (int, error) {
	if len(w) == 0 {
		return 0, fmt.Errorf("abi: cannot decode tuple from empty data")
	}
	var (
		wordIdx   int
		wordsRead int
	)
	for _, e := range *t {
		if e.DynamicType() {
			offset, err := w[wordIdx].Int()
			if err != nil {
				return 0, err
			}
			if offset%WordLength != 0 {
				return 0, fmt.Errorf("abi: cannot decode tuple, offset not a multiple of word length")
			}
			wordOffset := offset / WordLength
			if wordOffset >= len(w) {
				return 0, fmt.Errorf("abi: cannot decode tuple, offset exceeds data length")
			}
			n, err := e.DecodeABI(w[wordOffset:])
			if err != nil {
				return 0, err
			}
			wordIdx++
			if wordOffset+n > wordsRead {
				wordsRead = wordOffset + n
			}
		} else {
			n, err := e.DecodeABI(w[wordIdx:])
			if err != nil {
				return 0, err
			}
			wordIdx += n
			if wordIdx > wordsRead {
				wordsRead = wordIdx
			}
		}
	}
	return wordsRead, nil
}

// decodeArray decodes a dynamic array from the given words and stores the result
// in the given array. The elements of the array are decoded using the given
// type definition.
func decodeArray(a *[]Type, w Words, c *Config, t TypeDefinition) (int, error) {
	if len(w) == 0 {
		return 0, fmt.Errorf("abi: cannot decode array from empty data")
	}
	size, err := w[0].Int()
	if err != nil {
		return 0, err
	}
	if size >= math.MaxInt {
		return 0, fmt.Errorf("abi: cannot decode array, size too large")
	}
	if size+1 > len(w) {
		return 0, fmt.Errorf("abi: cannot decode array, size exceeds data length")
	}
	*a = make([]Type, size)
	for i := 0; i < size; i++ {
		if (*a)[i], err = c.NewType(t); err != nil {
			return 0, err
		}
	}
	if _, err := decodeTuple(a, w[1:]); err != nil {
		return 0, err
	}
	return size + 1, nil
}

func decodeFixedArray(a *[]Type, w Words, c *Config, t TypeDefinition, size int) (int, error) {
	var err error
	if len(w) == 0 {
		return 0, fmt.Errorf("abi: cannot decode array[%d] from empty value", size)
	}
	*a = make([]Type, size)
	for i := 0; i < size; i++ {
		if (*a)[i], err = c.NewType(t); err != nil {
			return 0, err
		}
	}
	if _, err := decodeTuple(a, w); err != nil {
		return 0, err
	}
	return size, nil
}

// decodeBytes decodes a dynamic byte array from the given words and stores the
// result in the given byte array.
func decodeBytes(b *[]byte, w Words) (int, error) {
	if len(w) == 0 {
		return 0, fmt.Errorf("abi: cannot decode bytes from empty value")
	}
	size, err := w[0].Int()
	if err != nil {
		return 0, err
	}
	if size >= math.MaxInt {
		return 0, fmt.Errorf("abi: cannot decode bytes, size too large")
	}
	if requiredWords(size)+1 > len(w) {
		return 0, fmt.Errorf("abi: cannot decode bytes, size exceeds data length")
	}
	*b = w[1:].Bytes()[0:size]
	return size + 1, nil
}

func decodeFixedBytes(b *[]byte, w Words, size int) (int, error) {
	if len(w) == 0 {
		return 0, fmt.Errorf("abi: cannot unmarshal bytes%d from empty value", size)
	}
	*b = w.Bytes()[0:size]
	return 1, nil
}

func decodeInt(i *big.Int, w Words) (int, error) {
	if len(w) == 0 {
		return 0, fmt.Errorf("abi: cannot decode int from empty value")
	}
	i.Set(w[0].SignedBigInt())
	return 1, nil
}

func decodeUint(i *big.Int, w Words) (int, error) {
	if len(w) == 0 {
		return 0, fmt.Errorf("abi: cannot decode int from empty value")
	}
	i.Set(w[0].UnsignedBigInt())
	return 1, nil
}

func decodeBool(a *bool, w Words) (int, error) {
	if len(w) == 0 {
		return 0, fmt.Errorf("abi: cannot decode bool from empty value")
	}
	*a = w[0].IsZero() == false
	return 1, nil
}

func decodeUint32(b []byte) (uint32, error) {
	if len(b) > 4 {
		return 0, fmt.Errorf("abi: cannot decode uint32, value too large")
	}
	var i uint32
	for _, v := range b {
		i <<= 8
		i |= uint32(v)
	}
	return i, nil
}

func decodeUint64(b []byte) (uint64, error) {
	if len(b) > 8 {
		return 0, fmt.Errorf("abi: cannot decode uint64, value too large")
	}
	var i uint64
	for _, v := range b {
		i <<= 8
		i |= uint64(v)
	}
	return i, nil
}

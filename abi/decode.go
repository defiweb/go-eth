package abi

import (
	"fmt"
	"math/big"

	"github.com/defiweb/go-eth/types"
)

// DecodeValue decodes the given ABI-encoded data into the given value.
// Value must be a pointer to a struct or a map.
func DecodeValue(t Type, abi []byte, val any) error {
	return Default.DecodeValue(t, abi, val)
}

// DecodeValues decodes the given ABI-encoded data into the given values.
// The t type must be a tuple type.
func DecodeValues(t Type, abi []byte, vals ...any) error {
	return Default.DecodeValues(t, abi, vals...)
}

// DecodeValue decodes the given ABI-encoded data into the given value.
// Value must be a pointer to a struct or a map.
func (a *ABI) DecodeValue(t Type, abi []byte, val any) error {
	v := t.Value()
	if _, err := v.DecodeABI(BytesToWords(abi)); err != nil {
		return err
	}
	return a.Mapper.Map(v, val)
}

// DecodeValues decodes the given ABI-encoded data into the given values.
// The t type must be a tuple type.
func (a *ABI) DecodeValues(t Type, abi []byte, vals ...any) error {
	v, ok := t.Value().(*TupleValue)
	if !ok {
		return fmt.Errorf("abi: cannot decode values, expected tuple type")
	}
	if len(*v) != len(vals) {
		return fmt.Errorf("abi: cannot decode tuple, expected %d values, got %d", len(*v), len(vals))
	}
	if _, err := v.DecodeABI(BytesToWords(abi)); err != nil {
		return err
	}
	for i, elem := range *v {
		if err := a.Mapper.Map(elem.Value, vals[i]); err != nil {
			return err
		}
	}
	return nil
}

// decodeTuple decodes a tuple from the given words and stores the result in the
// given tuple. The tuple must contain the correct number of elements.
func decodeTuple(t *[]Value, w Words) (int, error) {
	var (
		wordIdx   int
		wordsRead int
	)
	for _, e := range *t {
		if wordIdx >= len(w) {
			return 0, fmt.Errorf("abi: cannot decode tuple, unexpected end of data")
		}
		if e.IsDynamic() {
			offset, err := readInt(&w[wordIdx])
			if err != nil {
				return 0, fmt.Errorf("abi: cannot decode tuple, invalid offset: %v", err)
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
// in the given array. The elements of the array are decoded to t type.
func decodeArray(a *[]Value, w Words, t Type) (int, error) {
	if len(w) == 0 {
		return 0, fmt.Errorf("abi: cannot decode array from empty data")
	}
	size, err := readInt(&w[0])
	if err != nil {
		return 0, err
	}
	if size+1 > len(w) {
		return 0, fmt.Errorf("abi: cannot decode array, size exceeds data length")
	}
	*a = make([]Value, size)
	for i := 0; i < size; i++ {
		(*a)[i] = t.Value()
	}
	if _, err := decodeTuple(a, w[1:]); err != nil {
		return 0, err
	}
	return size + 1, nil
}

// decodeFixedArray decodes a fixed array from the given words into the values
// in the given array.
func decodeFixedArray(a *[]Value, w Words) (int, error) {
	if len(w) == 0 {
		return 0, fmt.Errorf("abi: cannot decode array[%d] from empty data", len(*a))
	}
	if _, err := decodeTuple(a, w); err != nil {
		return 0, err
	}
	return len(*a), nil
}

// decodeBytes decodes a dynamic byte array from the given words and stores the
// result in the given byte array.
func decodeBytes(b *[]byte, w Words) (int, error) {
	if len(w) == 0 {
		return 0, fmt.Errorf("abi: cannot decode bytes from empty data")
	}
	size, err := readInt(&w[0])
	if err != nil {
		return 0, err
	}
	l := requiredWords(size)
	if l+1 > len(w) {
		return 0, fmt.Errorf("abi: cannot decode bytes, size exceeds data length")
	}
	*b = w[1 : l+1].Bytes()[0:size]
	return size + 1, nil
}

// decodeFixedBytes decodes a fixed byte of the given size from the given words
// and stores the result in the given byte array.
func decodeFixedBytes(b *[]byte, w Words, size int) (int, error) {
	if len(w) == 0 {
		return 0, fmt.Errorf("abi: cannot decode bytes%d from empty data", size)
	}
	if len(*b) != size {
		return 0, fmt.Errorf("abi: cannot decode bytes%d, expected %d bytes, got %d", size, size, len(*b))
	}
	copy(*b, w[0].Bytes()[0:size])
	return 1, nil
}

// decodeInt decodes an integer of the given size from the given words and
// stores the result in the given integer. If the integer is larger than the
// maximum value for the given size, an error is returned.
func decodeInt(v *big.Int, w Words, size int) (int, error) {
	if len(w) == 0 {
		return 0, fmt.Errorf("abi: cannot decode int from empty data")
	}
	x := newIntX(size)
	if err := x.SetBytes(w[0].Bytes()); err != nil {
		return 0, err
	}
	v.Set(x.BigInt())
	return 1, nil
}

// decodeUint decodes an unsigned integer of the given size from the given
// words and stores the result in the given integer. If the integer is larger
// than the maximum value for the given size, an error is returned.
func decodeUint(v *big.Int, w Words, size int) (int, error) {
	if len(w) == 0 {
		return 0, fmt.Errorf("abi: cannot decode int from empty data")
	}
	x := newUintX(size)
	if err := x.SetBytes(w[0].Bytes()); err != nil {
		return 0, err
	}
	v.Set(x.BigInt())
	return 1, nil
}

// decodeBool decodes a boolean from the given words and stores the result in
// the given boolean.
func decodeBool(a *bool, w Words) (int, error) {
	if len(w) == 0 {
		return 0, fmt.Errorf("abi: cannot decode bool from empty data")
	}
	*a = w[0].IsZero() == false
	return 1, nil
}

// decodeAddress decodes an address from the given words and stores the result
// in the given address.
func decodeAddress(v *types.Address, w Words) (int, error) {
	if len(w) == 0 {
		return 0, fmt.Errorf("abi: cannot decode address from empty data")
	}
	*v = types.MustAddressFromBytes(w[0].Bytes()[WordLength-types.AddressLength:])
	return 1, nil
}

// readInt reads an integer from the given word.
func readInt(w *Word) (int, error) {
	i32 := newIntX(32)
	if err := i32.SetBytes(w.Bytes()); err != nil {
		return 0, err
	}
	return i32.Int()
}

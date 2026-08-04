package abi

import (
	"fmt"
	"math/big"

	"github.com/defiweb/go-eth/types"
)

// EncodeValue encodes a value to ABI encoding.
func EncodeValue(t Type, val any) ([]byte, error) {
	return Default.EncodeValue(t, val)
}

// MustEncodeValue is like EncodeValue but panics on error.
func MustEncodeValue(t Type, val any) []byte {
	return Default.MustEncodeValue(t, val)
}

// EncodeValues encodes a list of values to ABI encoding.
// The t type must be a tuple type.
func EncodeValues(t Type, vals ...any) ([]byte, error) {
	return Default.EncodeValues(t, vals...)
}

// MustEncodeValues is like EncodeValues but panics on error.
func MustEncodeValues(t Type, vals ...any) []byte {
	return Default.MustEncodeValues(t, vals...)
}

// EncodeValue encodes a value to ABI encoding.
func (a *ABI) EncodeValue(t Type, val any) ([]byte, error) {
	v := t.Value()
	if err := a.Mapper.Map(val, v); err != nil {
		return nil, err
	}
	words, err := v.EncodeABI()
	if err != nil {
		return nil, err
	}
	return words.Bytes(), nil
}

// MustEncodeValue is like EncodeValue but panics on error.
func (a *ABI) MustEncodeValue(t Type, val any) []byte {
	encoded, err := a.EncodeValue(t, val)
	if err != nil {
		panic(err)
	}
	return encoded
}

// EncodeValues encodes a list of values to ABI encoding.
// The t type must be a tuple type.
func (a *ABI) EncodeValues(t Type, vals ...any) ([]byte, error) {
	v, ok := t.Value().(*TupleValue)
	if !ok {
		return nil, fmt.Errorf("abi: cannot encode values, expected tuple type")
	}
	if len(*v) != len(vals) {
		return nil, fmt.Errorf("abi: expected %d values, got %d", len(*v), len(vals))
	}
	for i, elem := range *v {
		if err := a.Mapper.Map(vals[i], elem.Value); err != nil {
			return nil, err
		}
	}
	words, err := v.EncodeABI()
	if err != nil {
		return nil, err
	}
	return words.Bytes(), nil
}

// MustEncodeValues is like EncodeValues but panics on error.
func (a *ABI) MustEncodeValues(t Type, vals ...any) []byte {
	encoded, err := a.EncodeValues(t, vals...)
	if err != nil {
		panic(err)
	}
	return encoded
}

// encodeTuple encodes a tuple of types.
//
// A tuple consists of two sections: head and tail. The tail section is placed
// after the head section. During encoding, if the element is static, it is
// encoded directly in the head section. If the element is dynamic, it is
// encoded in the tail section, and the offset to the element is placed in the
// head section. The offset is a 256-bit integer (single word) that points to
// the start of the element in the tail section. The offset is relative to the
// beginning of the tuple.
func encodeTuple(t []Value) (Words, error) {
	var (
		head      Words
		tail      Words
		headLen   int
		tailLen   int
		offsetIdx []int // indices of head elements that are offsets
		offsetVal []int // offset values for head elements minus headLen
	)
	for _, p := range t {
		words, err := p.EncodeABI()
		if err != nil {
			return nil, err
		}
		if p.IsDynamic() {
			// At this point, we do not know what the number of words in the
			// head will be, so we cannot calculate the offset. Instead, we
			// store the index of the offset element and the number of words
			// in the tail section. We will calculate the offset later.
			head = append(head, Word{})
			tail = append(tail, words...)
			offsetIdx = append(offsetIdx, len(head)-1) // index of offset element
			offsetVal = append(offsetVal, tailLen)     // number of words in tail section
			headLen += WordLength
			tailLen += len(words) * WordLength
		} else {
			// If a type is not dynamic, it is encoded directly in the head
			// section.
			head = append(head, words...)
			headLen += len(words) * WordLength
		}
		continue
	}
	// Fast path if there are no dynamic elements.
	if len(tail) == 0 {
		return head, nil
	}
	// Calculate the offsets for the dynamic elements as described above.
	for n, i := range offsetIdx {
		if err := writeInt(&head[i], headLen+offsetVal[n]); err != nil {
			return nil, err
		}
	}
	// Append the tail section to the head section.
	words := make(Words, len(head)+len(tail))
	copy(words, head)
	copy(words[len(head):], tail)
	return words, nil
}

// encodeArray encodes a dynamic array.
//
// The array is encoded just like a tuple, except that the first word is the
// number of elements in the array. All array elements must be of the same type.
func encodeArray(a []Value) (Words, error) {
	tuple, err := encodeTuple(a)
	if err != nil {
		return nil, err
	}
	words := make(Words, len(tuple)+1)
	if err := writeInt(&words[0], len(a)); err != nil {
		return nil, err
	}
	copy(words[1:], tuple)
	return words, nil
}

// encodeFixedArray encodes a fixed-size array.
//
// The fixed-size array is encoded just like a tuple. All array elements must be
// of the same type.
func encodeFixedArray(a []Value) (Words, error) {
	return encodeTuple(a)
}

// encodeBytes encodes a dynamic byte sequence.
//
// The byte sequence is encoded as multiple words, padded on the right if
// needed. The length of the byte sequence is encoded as a 256-bit integer
// (single word) before the byte sequence.
func encodeBytes(b []byte) (Words, error) {
	words := make(Words, requiredWords(len(b))+1)
	if err := writeInt(&words[0], len(b)); err != nil {
		return nil, err
	}
	for i, w := range BytesToWords(b) {
		words[i+1] = w
	}
	return words, nil
}

// encodeFixedBytes encodes a fixed-size byte sequence.
//
// The fixed-size byte sequence is encoded in a single word, padded on the
// left if needed.
func encodeFixedBytes(b []byte, size int) (Words, error) {
	word := Word{}
	if len(b) > size {
		return Words{}, fmt.Errorf("abi: cannot encode %d bytes to bytes%d", len(b), size)
	}
	if err := word.SetBytesPadRight(b); err != nil {
		return nil, err
	}
	return Words{word}, nil
}

// encodeInt encodes an integer.
//
// The integer is encoded as two's complement integer. If the integer cannot
// be represented in number of bits specified by the size argument, an error
// is returned.
func encodeInt(v *big.Int, size int) (Words, error) {
	w := Word{}
	x := newIntX(size)
	if err := x.SetBigInt(v); err != nil {
		return nil, err
	}
	if err := x.FillBytes(w[:]); err != nil {
		return nil, err
	}
	return Words{w}, nil
}

// encodeUint encodes an unsigned integer.
//
// The integer is encoded as unsigned integer. If the integer cannot be
// represented in number of bits specified by the size argument, an error
// is returned.
func encodeUint(v *big.Int, size int) (Words, error) {
	w := Word{}
	x := newUintX(size)
	if err := x.SetBigInt(v); err != nil {
		return nil, err
	}
	if err := x.FillBytes(w[:]); err != nil {
		return nil, err
	}
	return Words{w}, nil
}

// encodeBool encodes a boolean.
//
// The boolean is encoded as a single word where the least significant bit
// is the value of the boolean.
func encodeBool(b bool) Words {
	w := Word{}
	if b {
		w[WordLength-1] = 1
	}
	return Words{w}
}

// encodeAddress encodes an address.
//
// An address is encoded as a 160-bit byte sequence, padded on the left.
func encodeAddress(val types.Address) (Words, error) {
	w := Word{}
	if err := w.SetBytesPadLeft(val.Bytes()); err != nil {
		return nil, err
	}
	return Words{w}, nil
}

// writeInt writes an integer to a word.
func writeInt(w *Word, x int) error {
	i32 := newIntX(32)
	if err := i32.SetInt(x); err != nil {
		return err
	}
	return i32.FillBytes(w[:])
}

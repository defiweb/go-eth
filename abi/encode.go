package abi

import (
	"fmt"
	"math/big"
)

func EncodeValue(t Value, val any) ([]byte, error) {
	return NewEncoder(DefaultConfig).EncodeValue(t, val)
}

func EncodeValues(t *TupleValue, vals ...any) ([]byte, error) {
	return NewEncoder(DefaultConfig).EncodeValues(t, vals...)
}

func MustEncodeValue(t Value, val any) []byte {
	b, err := EncodeValue(t, val)
	if err != nil {
		panic(err)
	}
	return b
}

func MustEncodeValues(t *TupleValue, vals ...any) []byte {
	b, err := EncodeValues(t, vals...)
	if err != nil {
		panic(err)
	}
	return b
}

type Encoder struct {
	Config *Config
}

func NewEncoder(c *Config) *Encoder {
	return &Encoder{Config: c}
}

func (e *Encoder) EncodeValue(t Value, val any) ([]byte, error) {
	if err := e.Config.Mapper.Map(val, t); err != nil {
		return nil, err
	}
	words, err := t.EncodeABI()
	if err != nil {
		return nil, err
	}
	return words.Bytes(), nil
}

func (e *Encoder) EncodeValues(t *TupleValue, vals ...any) ([]byte, error) {
	if t.Size() != len(vals) {
		return nil, fmt.Errorf("abi: expected %d values, got %d", t.Size(), len(vals))
	}
	for i, elem := range t.Elements() {
		if err := e.Config.Mapper.Map(vals[i], elem); err != nil {
			return nil, err
		}
	}
	words, err := t.EncodeABI()
	if err != nil {
		return nil, err
	}
	return words.Bytes(), nil
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
		if p.DynamicType() {
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
		if headLen+offsetVal[n] < 0 {
			return nil, fmt.Errorf("abi: element offset overflow")
		}
		head[i].SetInt(headLen + offsetVal[n])
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
// number of elements in the elems.
func encodeArray(a []Value) (Words, error) {
	tuple, err := encodeTuple(a)
	if err != nil {
		return nil, err
	}
	words := make(Words, len(tuple)+1)
	words[0].SetInt(len(a))
	copy(words[1:], tuple)
	return words, nil
}

// encodeFixedArray encodes a fixed-size array.
//
// The fixed-size array is encoded just like a tuple.
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
	words[0].SetInt(len(b))
	words.SetBytesAt(1, b)
	return words, nil
}

// encodeFixedBytes encodes a fixed-size byte sequence.
//
// The fixed-size byte sequence is encoded in a single word, padded on the
// left if needed.
func encodeFixedBytes(b []byte) (Words, error) {
	word := Word{}
	if err := word.SetBytesPadLeft(b); err != nil {
		return nil, err
	}
	return Words{word}, nil
}

// encodeInt encodes an integer.
//
// The integer is encoded as two's complement 256-bit integer.
func encodeInt(i *big.Int, size int) (Words, error) {
	w := Word{}
	bitLen := signedBitLen(i)
	if bitLen > size*8 {
		return Words{w}, fmt.Errorf("abi: cannot encode %d-bit integer into int%d", bitLen, size*8)
	}
	if err := w.SetBigInt(i); err != nil {
		return Words{w}, err
	}
	return Words{w}, nil
}

// encodeUint encodes an unsigned integer.
//
// The integer is encoded as 256-bit unsigned integer.
func encodeUint(i *big.Int, size int) (Words, error) {
	w := Word{}
	if i.Sign() < 0 {
		return Words{w}, fmt.Errorf("abi: cannot encode negative integer to uint%d", size*8)
	}
	bitLen := i.BitLen()
	if bitLen > size*8 {
		return Words{w}, fmt.Errorf("abi: cannot encode %d-bit integer into uint%d", bitLen, size*8)
	}
	if err := w.SetBigInt(i); err != nil {
		return Words{w}, err
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

func encodeUint32(i uint32) []byte {
	b := make([]byte, 4)
	b[0] = byte(i >> 24)
	b[1] = byte(i >> 16)
	b[2] = byte(i >> 8)
	b[3] = byte(i)
	return b
}

// encodeUint64 in the big endian order.
func encodeUint64(i uint64) []byte {
	b := make([]byte, 8)
	b[0] = byte(i >> 56)
	b[1] = byte(i >> 48)
	b[2] = byte(i >> 40)
	b[3] = byte(i >> 32)
	b[4] = byte(i >> 24)
	b[5] = byte(i >> 16)
	b[6] = byte(i >> 8)
	b[7] = byte(i)
	return b
}

// signedBitLen returns the number of bits in the signed integer i.
func signedBitLen(x *big.Int) int {
	if x == nil || x.Sign() == 0 {
		return 0
	}
	if x.Sign() < 0 {
		return new(big.Int).Not(x).BitLen() + 1
	}
	return x.BitLen() + 1
}

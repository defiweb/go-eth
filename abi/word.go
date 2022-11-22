package abi

import (
	"fmt"
	"math/bits"
)

// WordLength is the number of bytes in an EVM word.
const WordLength = 32

// Word represents a 32-bytes EVM word.
type Word [WordLength]byte

func BytesToWords(b []byte) []Word {
	var words Words
	words.SetBytes(b)
	return words
}

// SetBytesPadRight sets the word to the given bytes, padded on the right.
func (w *Word) SetBytesPadRight(b []byte) error {
	if len(b) > WordLength {
		return fmt.Errorf("abi: cannot set %d-byte data to a %d-byte word", len(b), WordLength)
	}
	for i := len(b); i < WordLength; i++ {
		w[i] = 0
	}
	copy((*w)[:], b)
	return nil
}

// SetBytesPadLeft sets the word to the given bytes, padded on the left.
func (w *Word) SetBytesPadLeft(b []byte) error {
	if len(b) > WordLength {
		return fmt.Errorf("abi: cannot set %d-byte data to a %d-byte word", len(b), WordLength)
	}
	for i := 0; i < WordLength-len(b); i++ {
		w[i] = 0
	}
	copy((*w)[WordLength-len(b):], b)
	return nil
}

// Bytes returns the word as a byte slice.
func (w Word) Bytes() []byte {
	return w[:]
}

// IsZero returns true if all bytes in then word are zeros.
func (w Word) IsZero() bool {
	for _, b := range w {
		if b != 0 {
			return false
		}
	}
	return true
}

// LeadingZeros returns the number of leading zero bits.
func (w Word) LeadingZeros() int {
	for i, b := range w {
		if b != 0 {
			return i*8 + bits.LeadingZeros8(b)
		}
	}
	return WordLength * 8
}

// TrailingZeros returns the number of trailing zero bits.
func (w Word) TrailingZeros() int {
	for i := len(w) - 1; i >= 0; i-- {
		if w[i] != 0 {
			return (len(w)-i-1)*8 + bits.TrailingZeros8(w[i])
		}
	}
	return WordLength * 8
}

// Words is a slice of words.
type Words []Word

// SetBytes sets the words to the given bytes.
func (w *Words) SetBytes(b []byte) {
	*w = make([]Word, requiredWords(len(b)))
	for i := 0; i < len(b); i += WordLength {
		if len(b)-i < WordLength {
			copy((*w)[i/WordLength][i%WordLength:], b[i:])
		} else {
			copy((*w)[i/WordLength][:], b[i:i+WordLength])
		}
	}
}

// AppendBytes appends the given bytes to the words.
func (w *Words) AppendBytes(b []byte) {
	for len(b) == 0 {
		return
	}
	c := requiredWords(len(b))
	l := len(*w)
	w.grow(c)
	for i := 0; i < len(b); i += WordLength {
		if len(b)-i < WordLength {
			copy((*w)[l+i/WordLength][i%WordLength:], b[i:])
		} else {
			copy((*w)[l+i/WordLength][:], b[i:i+WordLength])
		}
	}
}

// Bytes returns the words as a byte slice.
func (w Words) Bytes() []byte {
	b := make([]byte, len(w)*WordLength)
	for i, word := range w {
		copy(b[i*WordLength:], word[:])
	}
	return b
}

func (w *Words) grow(n int) {
	w.resize(len(*w) + n)
}

func (w *Words) resize(n int) {
	if cap(*w) < n {
		cpy := make([]Word, len(*w), n)
		copy(cpy, *w)
		*w = cpy
	}
	*w = (*w)[:n]
}

// requiredWords returns the number of words required to store the given number
// of bytes.
func requiredWords(n int) int {
	if n <= 0 {
		return 0
	}
	return 1 + (n-1)/WordLength
}

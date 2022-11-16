package abi

import (
	"fmt"
	"math"
	"math/big"

	"github.com/defiweb/go-eth/hexutil"
)

// WordLength is the number of bytes in an EVM word.
const WordLength = 32

// Word represents a 32-bytes EVM word.
type Word [WordLength]byte

// Words is a slice of words.
type Words []Word

var (
	pow255     = new(big.Int).Lsh(big.NewInt(1), 255)
	pow256     = new(big.Int).Lsh(big.NewInt(1), 256)
	maxInt256  = new(big.Int).Sub(pow255, big.NewInt(1)) // 2^255 - 1
	maxUint256 = new(big.Int).Sub(pow256, big.NewInt(1)) // 2^256 - 1
)

func BytesToWord(b []byte) (Word, error) {
	var word Word
	if err := word.SetBytesPadRight(b); err != nil {
		return Word{}, err
	}
	return word, nil
}

func BytesToWords(b []byte) []Word {
	var words Words
	words.SetBytes(b)
	return words
}

// SetBytesPadRight sets the word to the given bytes, padded on the right.
func (w *Word) SetBytesPadRight(b []byte) error {
	if len(b) > WordLength {
		return fmt.Errorf("abi: cannot set %d bytes to a word of length %d", len(b), WordLength)
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
		return fmt.Errorf("abi: cannot set %d bytes to a word of length %d", len(b), WordLength)
	}
	for i := 0; i < WordLength-len(b); i++ {
		w[i] = 0
	}
	copy((*w)[WordLength-len(b):], b)
	return nil
}

// SetInt sets the word to the given integer.
func (w *Word) SetInt(i int) {
	switch math.MaxInt {
	case math.MaxInt32:
		w.SetUint32(uint32(i))
	default:
		w.SetUint64(uint64(i))
	}
}

// SetUint32 sets the word to the given uint32.
func (w *Word) SetUint32(i uint32) {
	(w)[WordLength-1] = byte(i)
	(w)[WordLength-2] = byte(i >> 8)
	(w)[WordLength-3] = byte(i >> 16)
	(w)[WordLength-4] = byte(i >> 24)
}

// SetUint64 sets the word to the given uint64.
func (w *Word) SetUint64(i uint64) {
	(w)[WordLength-1] = byte(i)
	(w)[WordLength-2] = byte(i >> 8)
	(w)[WordLength-3] = byte(i >> 16)
	(w)[WordLength-4] = byte(i >> 24)
	(w)[WordLength-5] = byte(i >> 32)
	(w)[WordLength-6] = byte(i >> 40)
	(w)[WordLength-7] = byte(i >> 48)
	(w)[WordLength-8] = byte(i >> 56)
}

// SetBigInt sets the word to the given big integer.
func (w *Word) SetBigInt(i *big.Int) error {
	if i == nil || i.Sign() == 0 {
		*w = Word{}
	}
	if i.Sign() < 0 {
		if i.BitLen() > 255 {
			return fmt.Errorf("abi: cannot set negative integer of a size larger than 255 bits")
		}
		n := new(big.Int).Set(i)
		n.And(n, maxUint256)
	} else {
		if i.BitLen() > 256 {
			return fmt.Errorf("abi: cannot set integer of a size larger than 256 bits")
		}
	}
	return w.SetBytesPadLeft(i.Bytes())
}

// Bytes returns the word as a byte slice.
func (w Word) Bytes() []byte {
	return w[:]
}

func (w Word) Hex() string {
	return hexutil.BytesToHex(w[:])
}

// Int returns the word as an int. If the word is larger than the maximum
// integer size, an error is returned.
func (w Word) Int() (int, error) {
	switch math.MaxInt {
	case math.MaxInt32:
		i, err := w.Uint32()
		return int(i), err
	default:
		i, err := w.Uint64()
		return int(i), err
	}
}

// Uint32 returns the word as an uint32. If the word is larger than the maximum
// uint32 size, an error is returned.
func (w Word) Uint32() (uint32, error) {
	for i := 0; i < WordLength-4; i++ {
		if w[i] != 0 {
			return 0, fmt.Errorf("abi: uint32 overflow")
		}
	}
	return uint32(w[WordLength-1]) |
		uint32(w[WordLength-2])<<8 |
		uint32(w[WordLength-3])<<16 |
		uint32(w[WordLength-4])<<24, nil
}

// Uint64 returns the word as an uint64. If the word is larger than the maximum
// uint64 size, an error is returned.
func (w Word) Uint64() (uint64, error) {
	for i := 0; i < WordLength-8; i++ {
		if w[i] != 0 {
			return 0, fmt.Errorf("abi: uint64 overflow")
		}
	}
	return uint64(w[WordLength-1]) |
		uint64(w[WordLength-2])<<8 |
		uint64(w[WordLength-3])<<16 |
		uint64(w[WordLength-4])<<24 |
		uint64(w[WordLength-5])<<32 |
		uint64(w[WordLength-6])<<40 |
		uint64(w[WordLength-7])<<48 |
		uint64(w[WordLength-8])<<56, nil
}

// SignedBigInt returns the words as a signed big integer.
func (w *Word) SignedBigInt() *big.Int {
	i := new(big.Int).SetBytes(w.Bytes())
	if i.Cmp(maxInt256) > 0 {
		i.Add(maxUint256, big.NewInt(0).Neg(i))
		i.Add(i, big.NewInt(1))
		i.Neg(i)
	}
	return new(big.Int).SetBytes(w.Bytes())
}

// UnsignedBigInt returns the words as an unsigned big integer.
func (w *Word) UnsignedBigInt() *big.Int {
	return new(big.Int).SetBytes(w.Bytes())
}

// IsZero returns true if all bytes in then word are zeros.
func (w *Word) IsZero() bool {
	for _, b := range w {
		if b != 0 {
			return false
		}
	}
	return true
}

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

// SetBytesAt sets the words to the given bytes starting at the given index.
func (w *Words) SetBytesAt(index int, b []byte) {
	c := requiredWords(len(b))
	if index+c > len(*w) {
		w.resize(index + c)
	}
	for i := 0; i < len(b); i += WordLength {
		if len(b)-i < WordLength {
			copy((*w)[index+i/WordLength][i%WordLength:], b[i:])
		} else {
			copy((*w)[index+i/WordLength][:], b[i:i+WordLength])
		}
	}
}

// AppendBytes appends the given bytes to the words.
func (w *Words) AppendBytes(b []byte) {
	for len(b) == 0 {
		return
	}
	c := requiredWords(len(b))
	w.grow(c)
	*w = (*w)[:len(*w)+c]
	for i := 0; i < len(b); i += WordLength {
		if len(b)-i < WordLength {
			copy((*w)[len(*w)+i/WordLength][i%WordLength:], b[i:])
		} else {
			copy((*w)[len(*w)+i/WordLength][:], b[i:i+WordLength])
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

func (w Words) Hex() string {
	return hexutil.BytesToHex(w.Bytes())
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

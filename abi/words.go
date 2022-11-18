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
	MaxInt256  = new(big.Int).Sub(pow255, big.NewInt(1)) // 2^255 - 1
	MinInt256  = new(big.Int).Lsh(big.NewInt(-1), 255)   // 1^255
	MaxUint256 = new(big.Int).Sub(pow256, big.NewInt(1)) // 2^256 - 1
	pow255     = new(big.Int).Lsh(big.NewInt(1), 255)
	pow256     = new(big.Int).Lsh(big.NewInt(1), 256)
)

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
	w.SetInt64(int64(i))
}

// SetUint sets the word to the given unsigned integer.
func (w *Word) SetUint(i uint) {
	w.SetUint64(uint64(i))
}

// SetInt64 sets the word to the given int64.
func (w *Word) SetInt64(i int64) {
	_ = w.SetBigInt(new(big.Int).SetInt64(i))
}

// SetUint64 sets the word to the given uint64.
func (w *Word) SetUint64(i uint64) {
	_ = w.SetBytesPadLeft(new(big.Int).SetUint64(i).Bytes())
}

// SetBigInt sets the word to the given big integer.
func (w *Word) SetBigInt(i *big.Int) error {
	if i == nil || i.Sign() == 0 {
		*w = Word{}
	}
	if i.Sign() < 0 {
		bitLen := signedBitLen(i)
		i = new(big.Int).Set(i).And(i, MaxUint256)
		if bitLen > 256 {
			return fmt.Errorf("abi: cannot set %d-bit integer to a word of length 256", bitLen)
		}
	} else {
		bitLen := i.BitLen()
		if bitLen > 256 {
			return fmt.Errorf("abi: cannot set %d-bit integer to a word of length 256", bitLen)
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
	i, err := w.Int64()
	if err != nil {
		return 0, err
	}
	if i > math.MaxInt || i < math.MinInt {
		return 0, fmt.Errorf("abi: int overflow")
	}
	return int(i), nil
}

// Uint returns the word as an uint. If the word is larger than the maximum
// integer size, an error is returned.
func (w Word) Uint() (uint, error) {
	i, err := w.Uint64()
	if err != nil {
		return 0, err
	}
	if i > math.MaxUint {
		return 0, fmt.Errorf("abi: uint overflow")
	}
	return uint(i), nil
}

// Int64 returns the word as an int. If the word is larger than the maximum
// integer size, an error is returned.
func (w Word) Int64() (int64, error) {
	bn := w.SignedBigInt()
	if !bn.IsInt64() {
		return 0, fmt.Errorf("abi: int64 overflow")
	}
	return bn.Int64(), nil
}

// Uint64 returns the word as an uint64. If the word is larger than the maximum
// uint64 size, an error is returned.
func (w Word) Uint64() (uint64, error) {
	bn := w.UnsignedBigInt()
	if !bn.IsUint64() {
		return 0, fmt.Errorf("abi: uint64 overflow")
	}
	return bn.Uint64(), nil
}

// SignedBigInt returns the words as a signed big integer.
func (w *Word) SignedBigInt() *big.Int {
	i := new(big.Int).SetBytes(w.Bytes())
	if i.Cmp(MaxInt256) > 0 {
		i.Add(MaxUint256, big.NewInt(0).Neg(i))
		i.Add(i, big.NewInt(1))
		i.Neg(i)
	}
	return i
}

// UnsignedBigInt returns the words as an unsigned big integer.
func (w *Word) UnsignedBigInt() *big.Int {
	return new(big.Int).SetBytes(w.Bytes())
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

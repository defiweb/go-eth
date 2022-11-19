package abi

import (
	"bytes"
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/defiweb/go-eth/hexutil"
)

func TestWord_SetBytesPadRight(t *testing.T) {
	tests := []struct {
		args    []byte
		want    Word
		wantErr bool
	}{
		{
			args: []byte{0x01},
			want: hexToWord("0x0100000000000000000000000000000000000000000000000000000000000000"),
		},
		{
			args: []byte{0x01, 0x02},
			want: hexToWord("0x0102000000000000000000000000000000000000000000000000000000000000"),
		},
		{
			args:    bytes.Repeat([]byte{0x00}, 33),
			wantErr: true,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			var w Word
			err := w.SetBytesPadRight(tt.args)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, w)
			}
		})
	}
}

func TestWord_SetBytesPadLeft(t *testing.T) {
	tests := []struct {
		args    []byte
		want    Word
		wantErr bool
	}{
		{
			args: []byte{0x01},
			want: hexToWord("0x0000000000000000000000000000000000000000000000000000000000000001"),
		},
		{
			args: []byte{0x01, 0x02},
			want: hexToWord("0x0000000000000000000000000000000000000000000000000000000000000102"),
		},
		{
			args:    bytes.Repeat([]byte{0x00}, 33),
			wantErr: true,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			var w Word
			err := w.SetBytesPadLeft(tt.args)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, w)
			}
		})
	}
}

func TestWord_SetInt64(t *testing.T) {
	tests := []struct {
		arg  int64
		want Word
	}{
		{
			arg:  -1,
			want: hexToWord("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
		},
		{
			arg:  1,
			want: hexToWord("0x0000000000000000000000000000000000000000000000000000000000000001"),
		},
		{
			arg:  math.MinInt64,
			want: hexToWord("0xffffffffffffffffffffffffffffffffffffffffffffffff8000000000000000"),
		},
		{
			arg:  math.MaxInt64,
			want: hexToWord("0x0000000000000000000000000000000000000000000000007fffffffffffffff"),
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			var w Word
			w.SetInt64(tt.arg)
			assert.Equal(t, tt.want, w)
		})
	}
}

func TestWord_SetUint64(t *testing.T) {
	tests := []struct {
		arg  uint64
		want Word
	}{
		{
			arg:  1,
			want: hexToWord("0x0000000000000000000000000000000000000000000000000000000000000001"),
		},
		{
			arg:  math.MaxUint64,
			want: hexToWord("0x000000000000000000000000000000000000000000000000ffffffffffffffff"),
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			var w Word
			w.SetUint64(tt.arg)
			assert.Equal(t, tt.want, w)
		})
	}
}

func TestWord_SetBigInt(t *testing.T) {
	tests := []struct {
		arg     *big.Int
		want    Word
		wantErr bool
	}{
		{
			arg:  big.NewInt(-1),
			want: hexToWord("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
		},
		{
			arg:  big.NewInt(1),
			want: hexToWord("0x0000000000000000000000000000000000000000000000000000000000000001"),
		},
		{
			arg:  MaxInt256,
			want: hexToWord("0x7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
		},
		{
			arg:  MaxUint256,
			want: hexToWord("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
		},
		{
			arg:  new(big.Int).Add(MinInt256, big.NewInt(1)),
			want: hexToWord("0x8000000000000000000000000000000000000000000000000000000000000001"),
		},
		{
			arg:  MinInt256,
			want: hexToWord("0x8000000000000000000000000000000000000000000000000000000000000000"),
		},
		{
			arg:     new(big.Int).Sub(MinInt256, big.NewInt(1)),
			wantErr: true,
		},
		{
			arg:     new(big.Int).Add(MaxUint256, big.NewInt(1)),
			wantErr: true,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			var w Word
			err := w.SetBigInt(tt.arg)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, w)
			}
		})
	}
}

func TestWord_Hex(t *testing.T) {
	assert.Equal(t,
		"0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		hexToWord("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef").Hex(),
	)
}

func TestWord_Int64(t *testing.T) {
	tests := []struct {
		arg     Word
		want    int64
		wantErr bool
	}{
		{
			arg:  hexToWord("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
			want: -1,
		},
		{
			arg:  hexToWord("0x0000000000000000000000000000000000000000000000000000000000000001"),
			want: 1,
		},
		{
			arg:  hexToWord("0xffffffffffffffffffffffffffffffffffffffffffffffff8000000000000000"),
			want: math.MinInt64,
		},
		{
			arg:  hexToWord("0x0000000000000000000000000000000000000000000000007fffffffffffffff"),
			want: math.MaxInt64,
		},
		{
			arg:     hexToWord("0x0000000000000000000000000000000000000000000000010000000000000000"),
			wantErr: true,
		},
		{
			arg:     hexToWord("0xfffffffffffffffffffffffffffffffffffffffffffffffe0000000000000000"),
			wantErr: true,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			got, err := tt.arg.Int64()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestWord_Uint64(t *testing.T) {
	tests := []struct {
		arg     Word
		want    uint64
		wantErr bool
	}{
		{
			arg:  hexToWord("0x0000000000000000000000000000000000000000000000000000000000000001"),
			want: 1,
		},
		{
			arg:  hexToWord("0x000000000000000000000000000000000000000000000000ffffffffffffffff"),
			want: math.MaxUint64,
		},
		{
			arg:     hexToWord("0x0000000000000000000000000000000000000000000000010000000000000000"),
			wantErr: true,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			got, err := tt.arg.Uint64()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestWord_BigInt(t *testing.T) {
	tests := []struct {
		arg  Word
		want *big.Int
	}{
		{
			arg:  hexToWord("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
			want: big.NewInt(-1),
		},
		{
			arg:  hexToWord("0x0000000000000000000000000000000000000000000000000000000000000001"),
			want: big.NewInt(1),
		},
		{
			arg:  hexToWord("0x8000000000000000000000000000000000000000000000000000000000000000"),
			want: MinInt256,
		},
		{
			arg:  hexToWord("0x7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
			want: MaxInt256,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			got := tt.arg.BigInt()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWord_UBigInt(t *testing.T) {
	tests := []struct {
		arg  Word
		want *big.Int
	}{
		{
			arg:  hexToWord("0x0000000000000000000000000000000000000000000000000000000000000001"),
			want: big.NewInt(1),
		},
		{
			arg:  hexToWord("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"),
			want: MaxUint256,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			assert.Equal(t, tt.want, tt.arg.UBigInt())
		})
	}
}

func TestWord_IsZero(t *testing.T) {
	assert.True(t, hexToWord("0x0000000000000000000000000000000000000000000000000000000000000000").IsZero())
	assert.False(t, hexToWord("0x0000000000000000000000000000000000000000000000000000000000000001").IsZero())
	assert.False(t, hexToWord("0x1000000000000000000000000000000000000000000000000000000000000000").IsZero())
}

func TestWords_SetBytes(t *testing.T) {
	tests := []struct {
		arg  []byte
		want Words
	}{
		{
			arg:  []byte{0x01},
			want: Words{hexToWord("0x0100000000000000000000000000000000000000000000000000000000000000")},
		},
		{
			arg:  hexutil.MustHexToBytes("0x0000000000000000000000000000000000000000000000000000000000000001"),
			want: hexToWords("0x0000000000000000000000000000000000000000000000000000000000000001"),
		},
		{
			arg:  hexutil.MustHexToBytes("0x000000000000000000000000000000000000000000000000000000000000000001"),
			want: hexToWords("0x00000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000"),
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			var got Words
			got.SetBytes(tt.arg)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWord_LeadingZeros(t *testing.T) {
	tests := []struct {
		arg  Word
		want int
	}{
		{
			arg:  hexToWord("0x0000000000000000000000000000000000000000000000000000000000000000"),
			want: 256,
		},
		{
			arg:  hexToWord("0x0000000000000000000000000000000000000000000000000000000000000001"),
			want: 255,
		},
		{
			arg:  hexToWord("0x8000000000000000000000000000000000000000000000000000000000000000"),
			want: 0,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			assert.Equal(t, tt.want, tt.arg.LeadingZeros())
		})
	}
}

func TestWord_TrailingZeros(t *testing.T) {
	tests := []struct {
		arg  Word
		want int
	}{
		{
			arg:  hexToWord("0x0000000000000000000000000000000000000000000000000000000000000000"),
			want: 256,
		},
		{
			arg:  hexToWord("0x0000000000000000000000000000000000000000000000000000000000000001"),
			want: 0,
		},
		{
			arg:  hexToWord("0x8000000000000000000000000000000000000000000000000000000000000000"),
			want: 255,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			assert.Equal(t, tt.want, tt.arg.TrailingZeros())
		})
	}
}

func TestWords_AppendBytes(t *testing.T) {
	tests := []struct {
		words Words
		arg   []byte
		want  Words
	}{
		{
			words: Words{},
			arg:   []byte{0x01},
			want:  hexToWords("0x0100000000000000000000000000000000000000000000000000000000000000"),
		},
		{
			words: hexToWords("0x0000000000000000000000000000000000000000000000000000000000000001"),
			arg:   []byte{0x01},
			want:  hexToWords("0x00000000000000000000000000000000000000000000000000000000000000010100000000000000000000000000000000000000000000000000000000000000"),
		},
		{
			words: hexToWords("0x0000000000000000000000000000000000000000000000000000000000000001"),
			arg:   hexutil.MustHexToBytes("0x0000000000000000000000000000000000000000000000000000000000000001"),
			want:  hexToWords("0x00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001"),
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			tt.words.AppendBytes(tt.arg)
			assert.Equal(t, tt.want, tt.words)
		})
	}
}

func hexToWord(h string) Word {
	return BytesToWords(hexutil.MustHexToBytes(h))[0]
}

func hexToWords(h string) Words {
	return BytesToWords(hexutil.MustHexToBytes(h))
}

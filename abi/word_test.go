package abi

import (
	"bytes"
	"fmt"
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

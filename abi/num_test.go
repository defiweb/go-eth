package abi

import (
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntX_Bytes(t *testing.T) {
	tests := []struct {
		val  *intX
		set  *big.Int
		want []byte
	}{
		{
			val:  newIntX(8),
			set:  big.NewInt(0),
			want: []byte{0x00},
		},
		{
			val:  newIntX(8),
			set:  big.NewInt(1),
			want: []byte{0x01},
		},
		{
			val:  newIntX(8),
			set:  big.NewInt(-1),
			want: []byte{0xff},
		},
		{
			val:  newIntX(8),
			set:  big.NewInt(127),
			want: []byte{0x7f},
		},
		{
			val:  newIntX(8),
			set:  big.NewInt(-128),
			want: []byte{0x80},
		},
		{
			val:  newIntX(32),
			set:  big.NewInt(0),
			want: []byte{0x00, 0x00, 0x00, 0x00},
		},
		{
			val:  newIntX(32),
			set:  big.NewInt(1),
			want: []byte{0x00, 0x00, 0x00, 0x01},
		},
		{
			val:  newIntX(32),
			set:  big.NewInt(-1),
			want: []byte{0xff, 0xff, 0xff, 0xff},
		},
		{
			val:  newIntX(32),
			set:  MaxInt[32],
			want: []byte{0x7f, 0xff, 0xff, 0xff},
		},
		{
			val:  newIntX(32),
			set:  MinInt[32],
			want: []byte{0x80, 0x00, 0x00, 0x00},
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			require.NoError(t, tt.val.SetBigInt(tt.set))
			assert.Equal(t, tt.want, tt.val.Bytes())
		})
	}
}

func TestIntX_SetIntUint(t *testing.T) {
	tests := []struct {
		x      uint64
		bitLen int
		want   bool
	}{
		{x: 0, bitLen: 8, want: true},
		{x: 1, bitLen: 8, want: true},
		{x: 255, bitLen: 8, want: true},
		{x: 256, bitLen: 8, want: false},
		{x: 0, bitLen: 32, want: true},
		{x: 1, bitLen: 32, want: true},
		{x: math.MaxUint32, bitLen: 32, want: true},
		{x: math.MaxUint32 + 1, bitLen: 32, want: false},
		{x: math.MaxUint64, bitLen: 64, want: true},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			assert.Equal(t, tt.want, canSetUint(tt.x, tt.bitLen))
		})
	}
}

func TestIntX_SetBytes(t *testing.T) {
	tests := []struct {
		val     *intX
		bytes   []byte
		want    *big.Int
		wantErr bool
	}{
		{
			val:   newIntX(8),
			bytes: []byte{0x00},
			want:  big.NewInt(0),
		},
		{
			val:   newIntX(8),
			bytes: []byte{0x01},
			want:  big.NewInt(1),
		},
		{
			val:   newIntX(8),
			bytes: []byte{0xff},
			want:  big.NewInt(-1),
		},
		{
			val:   newIntX(8),
			bytes: []byte{0x7f},
			want:  big.NewInt(127),
		},
		{
			val:   newIntX(8),
			bytes: []byte{0x80},
			want:  big.NewInt(-128),
		},
		{
			val:   newIntX(32),
			bytes: []byte{0x00, 0x00, 0x00, 0x00},
			want:  big.NewInt(0),
		},
		{
			val:   newIntX(32),
			bytes: []byte{0x00, 0x00, 0x00, 0x01},
			want:  big.NewInt(1),
		},
		{
			val:   newIntX(32),
			bytes: []byte{0xff, 0xff, 0xff, 0xff},
			want:  big.NewInt(-1),
		},
		{
			val:   newIntX(32),
			bytes: []byte{0x7f, 0xff, 0xff, 0xff},
			want:  MaxInt[32],
		},
		{
			val:   newIntX(32),
			bytes: []byte{0x80, 0x00, 0x00, 0x00},
			want:  MinInt[32],
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			err := tt.val.SetBytes(tt.bytes)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, tt.val.val)
			}
		})
	}
}

func TestIntX_FillBytes(t *testing.T) {
	tests := []struct {
		val     *intX
		set     *big.Int
		want    []byte
		wantErr bool
	}{
		{
			val:  newIntX(8),
			set:  big.NewInt(0),
			want: []byte{0x00},
		},
		{
			val:  newIntX(8),
			set:  big.NewInt(1),
			want: []byte{0x01},
		},
		{
			val:  newIntX(8),
			set:  big.NewInt(-1),
			want: []byte{0xff},
		},
		{
			val:  newIntX(8),
			set:  big.NewInt(127),
			want: []byte{0x7f},
		},
		{
			val:  newIntX(8),
			set:  big.NewInt(-128),
			want: []byte{0x80},
		},
		{
			val:  newIntX(8),
			set:  big.NewInt(1),
			want: []byte{0x00, 0x01},
		},
		{
			val:  newIntX(8),
			set:  big.NewInt(-1),
			want: []byte{0xff, 0xff},
		},
		// Too small slice
		{
			val:     newIntX(16),
			set:     big.NewInt(0),
			want:    []byte{0x00},
			wantErr: true,
		},
		// Too large slice
		{
			val:     newIntX(8),
			set:     big.NewInt(0),
			want:    make([]byte, 33),
			wantErr: true,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			require.NoError(t, tt.val.SetBigInt(tt.set))
			res := make([]byte, len(tt.want))
			err := tt.val.FillBytes(res)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, res)
			}
		})
	}
}

func Test_signedBitLen(t *testing.T) {
	tests := []struct {
		arg  *big.Int
		want int
	}{
		{arg: big.NewInt(0), want: 0},
		{arg: MaxInt[256], want: 256},
		{arg: MinInt[256], want: 256},
		{arg: MaxUint[256], want: 257},
		{arg: bigIntStr("-0x010000000000000000"), want: 65},
		{arg: bigIntStr("-0x020000000000000000"), want: 66},
		{arg: bigIntStr("-0x030000000000000000"), want: 67},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			assert.Equal(t, tt.want, signedBitLen(tt.arg))
		})
	}
}

func Test_canSetInt(t *testing.T) {
	tests := []struct {
		x      int64
		bitLen int
		want   bool
	}{
		{x: 0, bitLen: 8, want: true},
		{x: 1, bitLen: 8, want: true},
		{x: -1, bitLen: 8, want: true},
		{x: 127, bitLen: 8, want: true},
		{x: -128, bitLen: 8, want: true},
		{x: 128, bitLen: 8, want: false},
		{x: -129, bitLen: 8, want: false},
		{x: 0, bitLen: 32, want: true},
		{x: 1, bitLen: 32, want: true},
		{x: -1, bitLen: 32, want: true},
		{x: math.MaxInt32, bitLen: 32, want: true},
		{x: math.MinInt32, bitLen: 32, want: true},
		{x: math.MaxInt32 + 1, bitLen: 32, want: false},
		{x: math.MinInt32 - 1, bitLen: 32, want: false},
		{x: math.MaxInt64, bitLen: 64, want: true},
		{x: math.MinInt64, bitLen: 64, want: true},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			assert.Equal(t, tt.want, canSetInt(tt.x, tt.bitLen))
		})
	}
}

func bigIntStr(s string) *big.Int {
	i, ok := new(big.Int).SetString(s, 0)
	if !ok {
		panic("invalid big.Int string")
	}
	return i
}

package abi

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_signedBitLen(t *testing.T) {
	tests := []struct {
		arg  *big.Int
		want int
	}{
		{arg: big.NewInt(0), want: 0},
		{arg: MaxInt256, want: 256},
		{arg: MinInt256, want: 256},
		{arg: MaxUint256, want: 257},
		{arg: bigIntMustSetString("-0x010000000000000000"), want: 65},
		{arg: bigIntMustSetString("-0x020000000000000000"), want: 66},
		{arg: bigIntMustSetString("-0x030000000000000000"), want: 67},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			assert.Equal(t, tt.want, signedBitLen(tt.arg))
		})
	}
}

func BenchmarkSignedBitLen(b *testing.B) {
	v1 := bigIntMustSetString("-0x010000000000000000")
	v2 := bigIntMustSetString("-0x020000000000000000")
	v3 := bigIntMustSetString("-0x030000000000000000")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		signedBitLen(v1)
		signedBitLen(v2)
		signedBitLen(v3)
	}
}

func bigIntMustSetString(s string) *big.Int {
	i, ok := new(big.Int).SetString(s, 0)
	if !ok {
		panic("invalid big.Int string")
	}
	return i
}

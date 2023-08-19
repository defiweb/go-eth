package abi

import (
	"crypto/rand"
	"testing"
)

func BenchmarkFourBytes_Match(b *testing.B) {
	data := make([]byte, 32*6+4)
	_, _ = rand.Read(data)
	f := FourBytes{0x01, 0x02, 0x03, 0x04}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f.Match(data)
	}
}

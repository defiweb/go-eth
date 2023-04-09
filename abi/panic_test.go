package abi

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/defiweb/go-eth/hexutil"
)

func TestPanicPrefix(t *testing.T) {
	assert.Equal(t, panicPrefix, Panic.FourBytes())
}

func TestDecodePanic(t *testing.T) {
	tests := []struct {
		data    []byte
		want    uint64
		wantNil bool
	}{
		{
			data: hexutil.MustHexToBytes("0x4e487b710000000000000000000000000000000000000000000000000000000000000000"),
			want: 0,
		},
		{
			data: hexutil.MustHexToBytes("0x4e487b71000000000000000000000000000000000000000000000000000000000000002a"),
			want: 42,
		},
		{
			// Invalid panic prefix.
			data:    hexutil.MustHexToBytes("0xaaaaaaaa00000000000000000000000000000000000000000000000000000000000000"),
			wantNil: true,
		},
		{
			// Empty panic data.
			data:    hexutil.MustHexToBytes("0x4e487b71"),
			wantNil: true,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			got := DecodePanic(tt.data)
			if tt.wantNil {
				assert.Nil(t, got)
			} else {
				assert.Equal(t, tt.want, got.Uint64())
			}
		})
	}
}

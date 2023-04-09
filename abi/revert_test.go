package abi

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/defiweb/go-eth/hexutil"
)

func TestRevertPrefix(t *testing.T) {
	assert.Equal(t, revertPrefix, Revert.FourBytes())
}

func TestDecodeRevert(t *testing.T) {
	tests := []struct {
		data []byte
		want string
	}{
		{
			data: hexutil.MustHexToBytes("0x08c379a00000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000e726576657274206d657373616765000000000000000000000000000000000000"),
			want: "revert message",
		},
		{
			data: hexutil.MustHexToBytes("0x08c379a00000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000004061616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161"),
			want: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		},
		{
			// Invalid revert prefix.
			data: hexutil.MustHexToBytes("0xaaaaaaaa0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000e726576657274206d657373616765000000000000000000000000000000000000"),
			want: "",
		},
		{
			// Empty revert data.
			data: hexutil.MustHexToBytes("0x08c379a0"),
			want: "",
		},
		{
			// Invalid revert data.
			data: hexutil.MustHexToBytes("0x08c379a0726576657274206d657373616765000000000000000000000000000000000000"),
			want: "",
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			assert.Equal(t, tt.want, DecodeRevert(tt.data))
		})
	}
}

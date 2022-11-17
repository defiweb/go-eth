package abi

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/defiweb/go-eth/hexutil"
)

func TestDecodeRevert(t *testing.T) {
	tests := []struct {
		data    []byte
		want    string
		wantErr bool
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
			data:    hexutil.MustHexToBytes("0xaaaaaaaa0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000e726576657274206d657373616765000000000000000000000000000000000000"),
			wantErr: true,
		},
		{
			// Empty revert data.
			data:    hexutil.MustHexToBytes("0x08c379a0"),
			wantErr: true,
		},
		{
			// Invalid revert data.
			data:    hexutil.MustHexToBytes("0x08c379a0726576657274206d657373616765000000000000000000000000000000000000"),
			wantErr: true,
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			got, err := DecodeRevert(tt.data)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

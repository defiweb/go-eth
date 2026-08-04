package keccak

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeccak256(t *testing.T) {
	tc := []struct {
		data [][]byte
		want string
	}{
		{
			data: [][]byte{[]byte("")},
			want: "c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470",
		},
		{
			data: [][]byte{[]byte("ab")},
			want: "67fad3bfa1e0321bd021ca805ce14876e50acac8ca8532eda8cbf924da565160",
		},
		{
			data: [][]byte{[]byte("a"), []byte("b")},
			want: "67fad3bfa1e0321bd021ca805ce14876e50acac8ca8532eda8cbf924da565160",
		},
	}
	for n, tt := range tc {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			h := Keccak256(tt.data...)
			assert.Equal(t, tt.want, hex.EncodeToString(h[:]))
		})
	}
}

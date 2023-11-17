package crypto

import (
	"golang.org/x/crypto/sha3"

	"github.com/defiweb/go-eth/types"
)

// Keccak256 calculates the Keccak256 hash of the given data.
func Keccak256(data ...[]byte) types.Hash {
	h := sha3.NewLegacyKeccak256()
	for _, i := range data {
		h.Write(i)
	}
	return types.MustHashFromBytes(h.Sum(nil), types.PadNone)
}

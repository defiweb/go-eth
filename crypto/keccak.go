package crypto

import (
	"golang.org/x/crypto/sha3"

	"web3rpc/types"
)

// Keccak256 calculates the Keccak256
func Keccak256(v ...[]byte) types.Hash {
	h := sha3.NewLegacyKeccak256()
	for _, i := range v {
		h.Write(i)
	}
	return types.BytesToHash(h.Sum(nil))
}

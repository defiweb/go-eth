// Package keccak provides Keccak-256 hashing functionality.
package keccak

import (
	"hash"
	"sync"

	"golang.org/x/crypto/sha3"
)

// Hash is a 32-byte Keccak256 hash.
//
// For most use cases, the [types.Hash] type should be used instead.
type Hash [32]byte

// Keccak256 calculates the Keccak256 hash of the given data.
func Keccak256(data ...[]byte) (h Hash) {
	k := keccakPool.Get().(hash.Hash)
	k.Reset()
	for _, i := range data {
		k.Write(i)
	}
	k.Sum(h[:0])
	keccakPool.Put(k)
	return
}

var keccakPool = sync.Pool{
	New: func() any { return sha3.NewLegacyKeccak256() },
}

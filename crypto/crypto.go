// Package crypto provides default implementations of cryptographic functions.
package crypto

import (
	"github.com/defiweb/go-eth/crypto/ecdsa"
	"github.com/defiweb/go-eth/crypto/keccak"
	"github.com/defiweb/go-eth/crypto/kzg4844"
	"github.com/defiweb/go-eth/crypto/primitives"
)

const (
	HashSize          = primitives.HashSize
	AddressSize       = primitives.AddressSize
	PrivateKeySize    = primitives.PrivateKeySize
	PublicKeySize     = primitives.PublicKeySize
	KZGHashSize       = primitives.KZGHashSize
	KZGScalarsPerBlob = primitives.KZGScalarsPerBlob
	KZGScalarSize     = primitives.KZGScalarSize
	KZGBlobSize       = primitives.KZGBlobSize
	KZGCommitmentSize = primitives.KZGCommitmentSize
	KZGProofSize      = primitives.KZGProofSize
	KZGPointSize      = primitives.KZGPointSize
)

type (
	Hash          = primitives.Hash
	Address       = primitives.Address
	PublicKey     = primitives.PublicKey
	PrivateKey    = primitives.PrivateKey
	Signature     = primitives.Signature
	KZGHash       = primitives.KZGHash
	KZGBlob       = primitives.KZGBlob
	KZGCommitment = primitives.KZGCommitment
	KZGProof      = primitives.KZGProof
	KZGPoint      = primitives.KZGPoint
)

// Default implementations of the crypto functions. Can be overridden to use
// alternative implementations.
var (
	Keccak256               = keccak.Keccak256
	ECGenerateKey           = ecdsa.GenerateKey
	ECPublicKeyToAddress    = ecdsa.PublicKeyToAddress
	ECPrivateKeyToPublicKey = ecdsa.PrivateKeyToPublicKey
	ECSignHash              = ecdsa.SignHash
	ECRecoverHash           = ecdsa.RecoverHash
	ECSignMessage           = ecdsa.SignMessage
	ECRecoverMessage        = ecdsa.RecoverMessage
	KZGBlobToCommitment     = kzg4844.BlobToCommitment
	KZGComputeProof         = kzg4844.ComputeProof
	KZGVerifyProof          = kzg4844.VerifyProof
	KZGComputeBlobProof     = kzg4844.ComputeBlobProof
	KZGVerifyBlobProof      = kzg4844.VerifyBlobProof
	KZGComputeBlobHashV1    = kzg4844.ComputeBlobHashV1
)

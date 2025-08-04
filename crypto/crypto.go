// Package crypto provides default implementations of cryptographic functions.
package crypto

import (
	"github.com/defiweb/go-eth/crypto/ecdsa"
	"github.com/defiweb/go-eth/crypto/keccak"
	"github.com/defiweb/go-eth/crypto/kzg4844"
)

// Default implementations of the crypto functions. Can be overridden to use
// alternative implementations.
var (
	Keccak256               = keccak.Keccak256
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

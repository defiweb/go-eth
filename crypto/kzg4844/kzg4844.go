package kzg4844

import (
	"crypto/sha256"
	"math/big"
	"sync"

	kzg4844 "github.com/crate-crypto/go-kzg-4844"

	"github.com/defiweb/go-eth/crypto/primitives"
)

// BLSModulus is the BLS12-381 scalar field modulus.
var BLSModulus = new(big.Int).SetBytes(kzg4844.BlsModulus[:])

// BlobToCommitment computes the KZG commitment for the given blob.
func BlobToCommitment(blob *primitives.KZGBlob) (primitives.KZGCommitment, error) {
	initContext()
	commitment, err := context.BlobToKZGCommitment(
		(*kzg4844.Blob)(blob),
		0,
	)
	if err != nil {
		return primitives.KZGCommitment{}, err
	}
	return (primitives.KZGCommitment)(commitment), nil
}

// ComputeProof computes the KZG proof and claim for the given blob and point.
func ComputeProof(blob *primitives.KZGBlob, point primitives.KZGPoint) (primitives.KZGProof, primitives.KZGPoint, error) {
	initContext()
	proof, claim, err := context.ComputeKZGProof(
		(*kzg4844.Blob)(blob),
		(kzg4844.Scalar)(point),
		0,
	)
	if err != nil {
		return primitives.KZGProof{}, primitives.KZGPoint{}, err
	}
	return (primitives.KZGProof)(proof), (primitives.KZGPoint)(claim), nil
}

// VerifyProof verifies the KZG proof for the given commitment, point, claim,
// and proof.
func VerifyProof(commitment primitives.KZGCommitment, point primitives.KZGPoint, claim primitives.KZGPoint, proof primitives.KZGProof) error {
	initContext()
	return context.VerifyKZGProof(
		(kzg4844.KZGCommitment)(commitment),
		(kzg4844.Scalar)(point),
		(kzg4844.Scalar)(claim),
		(kzg4844.KZGProof)(proof),
	)
}

// ComputeBlobProof computes the KZG proof for the given blob and commitment.
func ComputeBlobProof(blob *primitives.KZGBlob, commitment primitives.KZGCommitment) (primitives.KZGProof, error) {
	initContext()
	proof, err := context.ComputeBlobKZGProof(
		(*kzg4844.Blob)(blob),
		(kzg4844.KZGCommitment)(commitment),
		0,
	)
	if err != nil {
		return primitives.KZGProof{}, err
	}
	return (primitives.KZGProof)(proof), nil
}

// VerifyBlobProof verifies the KZG proof for the given blob, commitment, and proof.
func VerifyBlobProof(blob *primitives.KZGBlob, commitment primitives.KZGCommitment, proof primitives.KZGProof) error {
	initContext()
	return context.VerifyBlobKZGProof(
		(*kzg4844.Blob)(blob),
		(kzg4844.KZGCommitment)(commitment),
		(kzg4844.KZGProof)(proof),
	)
}

// ComputeBlobHashV1 calculates the 'versioned blob hash' of a commitment.
func ComputeBlobHashV1(commit primitives.KZGCommitment) (h primitives.KZGHash) {
	k := sha256.New()
	k.Write(commit[:])
	k.Sum(h[:0])
	h[0] = 0x01
	return
}

// context holds the necessary configuration needed to create and verify proofs.
var context *kzg4844.Context

var once sync.Once

func initContext() {
	once.Do(func() {
		var err error
		context, err = kzg4844.NewContext4096Secure()
		if err != nil {
			panic(err)
		}
	})
}

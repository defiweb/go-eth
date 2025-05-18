package kzg4844

import (
	"crypto/sha256"
	"sync"

	kzg4844 "github.com/crate-crypto/go-kzg-4844"
)

const (
	BlobLength       = 131072 // 128 KiB
	CommitmentLength = 48
	ProofLength      = 48
	PointLength      = 32
	ClaimLength      = 32
)

// Blob represents a 4844 data blob.
type Blob [BlobLength]byte

// Commitment is a serialized commitment to a polynomial.
type Commitment [CommitmentLength]byte

// Proof is a serialized commitment to the quotient polynomial.
type Proof [ProofLength]byte

// Point is a BLS field element.
type Point [PointLength]byte

// Claim is a claimed evaluation value in a specific point.
type Claim [ClaimLength]byte

func BlobToCommitment(blob *Blob) (Commitment, error) {
	initContext()
	commitment, err := context.BlobToKZGCommitment(
		(*kzg4844.Blob)(blob),
		0,
	)
	if err != nil {
		return Commitment{}, err
	}
	return (Commitment)(commitment), nil
}

func ComputeProof(blob *Blob, point Point) (Proof, Claim, error) {
	initContext()
	proof, claim, err := context.ComputeKZGProof(
		(*kzg4844.Blob)(blob),
		(kzg4844.Scalar)(point),
		0,
	)
	if err != nil {
		return Proof{}, Claim{}, err
	}
	return (Proof)(proof), (Claim)(claim), nil
}

func VerifyProof(commitment Commitment, point Point, claim Claim, proof Proof) error {
	initContext()
	return context.VerifyKZGProof(
		(kzg4844.KZGCommitment)(commitment),
		(kzg4844.Scalar)(point),
		(kzg4844.Scalar)(claim),
		(kzg4844.KZGProof)(proof),
	)
}

func ComputeBlobProof(blob *Blob, commitment Commitment) (Proof, error) {
	initContext()
	proof, err := context.ComputeBlobKZGProof(
		(*kzg4844.Blob)(blob),
		(kzg4844.KZGCommitment)(commitment),
		0,
	)
	if err != nil {
		return Proof{}, err
	}
	return (Proof)(proof), nil
}

func VerifyBlobProof(blob *Blob, commitment Commitment, proof Proof) error {
	initContext()
	return context.VerifyBlobKZGProof(
		(*kzg4844.Blob)(blob),
		(kzg4844.KZGCommitment)(commitment),
		(kzg4844.KZGProof)(proof),
	)
}

// ComputeBlobHashV1 calculates the 'versioned blob hash' of a commitment.
func ComputeBlobHashV1(commit Commitment) (h [32]byte) {
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

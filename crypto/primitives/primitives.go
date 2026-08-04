// Package primitives provides the basic data types used by the cryptographic
// functions in the crypto package and its implementation packages.
//
// DO NOT USE THIS PACKAGE DIRECTLY. Use the crypto package instead.
//
// This package is intended to be used by implementations of the functions in
// the crypto package.
package primitives

import (
	"errors"
	"math/big"
)

const (
	// HashSize is the size of a hash in bytes.
	HashSize = 32

	// AddressSize is the size of an Ethereum address in bytes.
	AddressSize = 20

	// PrivateKeySize is the size of an ECDSA private key in bytes.
	PrivateKeySize = 32

	// PublicKeySize is the size of an ECDSA public key in bytes.
	PublicKeySize = 64

	// KZGHashSize is the size of a blob hash in bytes.
	KZGHashSize = 32

	// KZGScalarsPerBlob is the number of field elements in a blob.
	KZGScalarsPerBlob = 4096

	// KZGScalarSize is the size of a single field element in bytes.
	KZGScalarSize = 32

	// KZGBlobSize is the size of a blob in bytes.
	KZGBlobSize = KZGScalarsPerBlob * KZGScalarSize // 128 KiB

	// KZGCommitmentSize is the size of a KZG commitment in bytes.
	KZGCommitmentSize = 48

	// KZGProofSize is the size of a KZG proof in bytes.
	KZGProofSize = 48

	// KZGPointSize is the size of a BLS field element in bytes.
	KZGPointSize = 32
)

// Hash is a 32-byte hash.
type Hash [HashSize]byte

// Address is a 20-byte Ethereum address.
type Address [AddressSize]byte

// PrivateKey is an ECDSA private key, stored as a big-endian 32-byte scalar.
type PrivateKey [PrivateKeySize]byte

// PublicKey is an ECDSA public key, stored as the 64-byte uncompressed
// representation: the X and Y coordinates, each big-endian and left-padded to
// 32 bytes, without the 0x04 prefix.
type PublicKey [PublicKeySize]byte

// Signature is an ECDSA signature.
type Signature struct {
	V *big.Int
	R *big.Int
	S *big.Int
}

// KZGHash is a blob hash.
type KZGHash [KZGHashSize]byte

// KZGBlob represents a 4844 data blob.
type KZGBlob [KZGBlobSize]byte

// KZGCommitment is a serialized commitment to a polynomial.
type KZGCommitment [KZGCommitmentSize]byte

// KZGProof is a serialized commitment to the quotient polynomial.
type KZGProof [KZGProofSize]byte

// KZGPoint is a BLS field element.
type KZGPoint [KZGPointSize]byte

// errNoDisclose is returned by the marshalers of PrivateKey to prevent the
// key material from being serialized by accident.
var errNoDisclose = errors.New("crypto: refusing to serialize a private key")

// IsZero reports whether the key is the zero value, which is not a valid key.
func (k PrivateKey) IsZero() bool {
	return k == PrivateKey{}
}

// Bytes returns a copy of the key as a byte slice.
func (k PrivateKey) Bytes() []byte {
	b := make([]byte, PrivateKeySize)
	copy(b, k[:])
	return b
}

// Scalar returns the key as a big integer.
func (k PrivateKey) Scalar() *big.Int {
	return new(big.Int).SetBytes(k[:])
}

// Zero wipes the key material.
func (k *PrivateKey) Zero() {
	clear(k[:])
}

// String implements the [fmt.Stringer] interface.
//
// It does not disclose the key material.
func (k PrivateKey) String() string {
	return "PrivateKey(redacted)"
}

// GoString implements the [fmt.GoStringer] interface.
//
// It does not disclose the key material.
func (k PrivateKey) GoString() string {
	return "PrivateKey(redacted)"
}

// MarshalText implements the [encoding.TextMarshaler] interface.
//
// It always returns an error, to prevent the key material from being
// serialized by accident.
func (k PrivateKey) MarshalText() ([]byte, error) {
	return nil, errNoDisclose
}

// MarshalJSON implements the [json.Marshaler] interface.
//
// It always returns an error, to prevent the key material from being
// serialized by accident.
func (k PrivateKey) MarshalJSON() ([]byte, error) {
	return nil, errNoDisclose
}

// IsZero reports whether the key is the zero value, which is not a valid key.
func (k PublicKey) IsZero() bool {
	return k == PublicKey{}
}

// Bytes returns a copy of the key as a byte slice, in the 64-byte
// uncompressed representation without the 0x04 prefix.
func (k PublicKey) Bytes() []byte {
	b := make([]byte, PublicKeySize)
	copy(b, k[:])
	return b
}

// X returns the X coordinate of the key.
func (k PublicKey) X() *big.Int {
	return new(big.Int).SetBytes(k[:32])
}

// Y returns the Y coordinate of the key.
func (k PublicKey) Y() *big.Int {
	return new(big.Int).SetBytes(k[32:])
}

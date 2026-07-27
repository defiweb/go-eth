// Package ecdsa provides ECDSA cryptographic functionality for Ethereum.
package ecdsa

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	secpecdsa "github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"

	"github.com/defiweb/go-eth/crypto/keccak"
	"github.com/defiweb/go-eth/crypto/primitives"
)

// AddMessagePrefix adds the Ethereum message prefix to the given data as
// defined in EIP-191.
func AddMessagePrefix(data []byte) []byte {
	return []byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data))
}

// GenerateKey generates a new ECDSA private key.
func GenerateKey() (primitives.PrivateKey, error) {
	pk, err := secp256k1.GeneratePrivateKey()
	if err != nil {
		return primitives.PrivateKey{}, fmt.Errorf("failed to generate ECDSA key: %w", err)
	}
	defer pk.Zero()
	return primitives.PrivateKey(pk.Serialize()), nil
}

// PublicKeyToAddress returns the Ethereum address for the given ECDSA
// public key.
func PublicKeyToAddress(publicKey primitives.PublicKey) (addr primitives.Address) {
	// The address is the last 20 bytes of the Keccak-256 hash of the 64-byte
	// public key: X and Y, each left-padded to 32 bytes.
	h := keccak.Keccak256(publicKey[:])
	copy(addr[:], h[12:])
	return
}

// PrivateKeyToPublicKey converts a private key to a public key.
// If the private key is the zero value, it returns the zero public key.
func PrivateKeyToPublicKey(privateKey primitives.PrivateKey) primitives.PublicKey {
	if privateKey.IsZero() {
		return primitives.PublicKey{}
	}
	priv := secp256k1.PrivKeyFromBytes(privateKey[:])
	defer priv.Zero()
	return newPublicKey(priv.PubKey())
}

// SignHash signs the given hash with the given private key.
func SignHash(privateKey primitives.PrivateKey, hash primitives.Hash) (*primitives.Signature, error) {
	if privateKey.IsZero() {
		return nil, fmt.Errorf("invalid private key")
	}
	priv := secp256k1.PrivKeyFromBytes(privateKey[:])
	defer priv.Zero()
	sig := secpecdsa.SignCompact(priv, hash[:], false)
	v := sig[0] - 27
	copy(sig, sig[1:])
	sig[64] = v
	return &primitives.Signature{
		V: new(big.Int).SetBytes(sig[64:]),
		R: new(big.Int).SetBytes(sig[:32]),
		S: new(big.Int).SetBytes(sig[32:64]),
	}, nil
}

// RecoverHash recovers the Ethereum address from the given hash and
// signature.
func RecoverHash(hash primitives.Hash, signature primitives.Signature) (*primitives.Address, error) {
	if signature.V.BitLen() > 8 {
		return nil, errors.New("invalid signature: V has more than 8 bits")
	}
	if signature.R.BitLen() > 256 {
		return nil, errors.New("invalid signature: R has more than 256 bits")
	}
	if signature.S.BitLen() > 256 {
		return nil, errors.New("invalid signature: S has more than 256 bits")
	}
	v, err := recoveryByte(byte(signature.V.Uint64()))
	if err != nil {
		return nil, err
	}
	rb := signature.R.Bytes()
	sb := signature.S.Bytes()
	bin := make([]byte, 65)
	bin[0] = v
	copy(bin[1+(32-len(rb)):], rb)
	copy(bin[33+(32-len(sb)):], sb)
	pub, _, err := secpecdsa.RecoverCompact(bin, hash[:])
	if err != nil {
		return nil, err
	}
	addr := PublicKeyToAddress(newPublicKey(pub))
	return &addr, nil
}

// SignMessage signs the given message with the given private key.
func SignMessage(key primitives.PrivateKey, data []byte) (*primitives.Signature, error) {
	if key.IsZero() {
		return nil, fmt.Errorf("invalid private key")
	}
	sig, err := SignHash(key, keccak.Keccak256(AddMessagePrefix(data)))
	if err != nil {
		return nil, err
	}
	sig.V = new(big.Int).Add(sig.V, big.NewInt(27))
	return sig, nil
}

// RecoverMessage recovers the Ethereum address from the given message and
// signature.
func RecoverMessage(data []byte, sig primitives.Signature) (*primitives.Address, error) {
	sig.V = new(big.Int).Sub(sig.V, big.NewInt(27))
	return RecoverHash(keccak.Keccak256(AddMessagePrefix(data)), sig)
}

func newPublicKey(pub *secp256k1.PublicKey) (key primitives.PublicKey) {
	// SerializeUncompressed returns 0x04 || X (32 bytes) || Y (32 bytes).
	b := pub.SerializeUncompressed()
	copy(key[:], b[1:])
	return
}

func recoveryByte(v byte) (byte, error) {
	switch v {
	case 0, 27:
		return 27, nil
	case 1, 28:
		return 28, nil
	default:
		return 0, fmt.Errorf("invalid signature: V must be 0, 1, 27, or 28, got %d", v)
	}
}

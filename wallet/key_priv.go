package wallet

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"

	"github.com/defiweb/go-eth/crypto"
	"github.com/defiweb/go-eth/crypto/txsign"
	"github.com/defiweb/go-eth/types"
)

type PrivateKey struct {
	private crypto.PrivateKey
	public  crypto.PublicKey
	address types.Address
}

// NewKeyFromECDSA creates a new private key from a [crypto.PrivateKey].
//
// It key must be a valid secp256k1 scalar.
func NewKeyFromECDSA(prv crypto.PrivateKey) (*PrivateKey, error) {
	if err := validateScalar(prv); err != nil {
		return nil, err
	}
	pub := crypto.ECPrivateKeyToPublicKey(prv)
	return &PrivateKey{
		private: prv,
		public:  pub,
		address: types.Address(crypto.ECPublicKeyToAddress(pub)),
	}, nil
}

// MustNewKeyFromECDSA works like [NewKeyFromECDSA] but panics on error.
func MustNewKeyFromECDSA(prv crypto.PrivateKey) *PrivateKey {
	key, err := NewKeyFromECDSA(prv)
	if err != nil {
		panic(err)
	}
	return key
}

// NewKeyFromBytes creates a new private key from private key bytes.
//
// The input must be exactly [crypto.PrivateKeySize] bytes, big-endian, and a
// valid secp256k1 scalar.
func NewKeyFromBytes(prv []byte) (*PrivateKey, error) {
	if len(prv) != crypto.PrivateKeySize {
		return nil, fmt.Errorf("invalid private key length: got %d, want %d", len(prv), crypto.PrivateKeySize)
	}
	return NewKeyFromECDSA(crypto.PrivateKey(prv))
}

// MustNewKeyFromBytes works like [NewKeyFromBytes] but panics on error.
func MustNewKeyFromBytes(prv []byte) *PrivateKey {
	key, err := NewKeyFromBytes(prv)
	if err != nil {
		panic(err)
	}
	return key
}

// NewRandomKey creates a random private key.
func NewRandomKey() *PrivateKey {
	key, err := crypto.ECGenerateKey()
	if err != nil {
		panic(err)
	}
	return MustNewKeyFromECDSA(key)
}

// PublicKey returns the ECDSA public key.
func (k *PrivateKey) PublicKey() crypto.PublicKey {
	return k.public
}

// PrivateKey returns the ECDSA private key.
func (k *PrivateKey) PrivateKey() crypto.PrivateKey {
	return k.private
}

// JSON returns the JSON representation of the private key.
func (k *PrivateKey) JSON(passphrase string, scryptN, scryptP int) ([]byte, error) {
	key, err := encryptV3Key(k.private, passphrase, scryptN, scryptP)
	if err != nil {
		return nil, err
	}
	return json.Marshal(key)
}

// Address implements the Key interface.
func (k *PrivateKey) Address() types.Address {
	return k.address
}

// SignHash implements the KeyWithHashSigner interface.
func (k *PrivateKey) SignHash(_ context.Context, hash types.Hash) (*types.Signature, error) {
	s, err := crypto.ECSignHash(k.private, crypto.Hash(hash))
	if err != nil {
		return nil, err
	}
	return (*types.Signature)(s), nil
}

// SignMessage implements the Key interface.
func (k *PrivateKey) SignMessage(_ context.Context, data []byte) (*types.Signature, error) {
	s, err := crypto.ECSignMessage(k.private, data)
	if err != nil {
		return nil, err
	}
	return (*types.Signature)(s), nil
}

// SignTransaction implements the Key interface.
func (k *PrivateKey) SignTransaction(_ context.Context, tx types.SignableTransaction) error {
	return txsign.Sign(k.private, tx)
}

// VerifyHash implements the KeyWithHashSigner interface.
func (k *PrivateKey) VerifyHash(_ context.Context, hash types.Hash, sig types.Signature) bool {
	addr, err := crypto.ECRecoverHash(crypto.Hash(hash), crypto.Signature(sig))
	if err != nil {
		return false
	}
	return types.Address(*addr) == k.address
}

// VerifyMessage implements the Key interface.
func (k *PrivateKey) VerifyMessage(_ context.Context, data []byte, sig types.Signature) bool {
	addr, err := crypto.ECRecoverMessage(data, crypto.Signature(sig))
	if err != nil {
		return false
	}
	return types.Address(*addr) == k.address
}

func validateScalar(prv crypto.PrivateKey) error {
	var s secp256k1.ModNScalar
	overflow := s.SetBytes((*[crypto.PrivateKeySize]byte)(&prv))
	zeroBit := s.IsZeroBit()
	s.Zero()
	if overflow|zeroBit != 0 {
		return errors.New("invalid private key: scalar must be in the range [1, N-1]")
	}
	return nil
}

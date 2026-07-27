package wallet

import (
	"context"
	"encoding/json"

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
func NewKeyFromECDSA(prv crypto.PrivateKey) *PrivateKey {
	pub := crypto.ECPrivateKeyToPublicKey(prv)
	return &PrivateKey{
		private: prv,
		public:  pub,
		address: types.Address(crypto.ECPublicKeyToAddress(pub)),
	}
}

// NewKeyFromBytes creates a new private key from private key bytes.
func NewKeyFromBytes(prv []byte) *PrivateKey {
	key := secp256k1.PrivKeyFromBytes(prv)
	defer key.Zero()
	return NewKeyFromECDSA(crypto.PrivateKey(key.Serialize()))
}

// NewRandomKey creates a random private key.
func NewRandomKey() *PrivateKey {
	key, err := crypto.ECGenerateKey()
	if err != nil {
		panic(err)
	}
	return NewKeyFromECDSA(key)
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

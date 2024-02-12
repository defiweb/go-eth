package wallet

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/json"

	"github.com/btcsuite/btcd/btcec/v2"

	"github.com/defiweb/go-eth/crypto"
	"github.com/defiweb/go-eth/types"
)

var s256 = btcec.S256()

type PrivateKey struct {
	private *ecdsa.PrivateKey
	public  *ecdsa.PublicKey
	address types.Address
	sign    crypto.Signer
	recover crypto.Recoverer
}

// NewKeyFromECDSA creates a new private key from an ecdsa.PrivateKey.
func NewKeyFromECDSA(prv *ecdsa.PrivateKey) *PrivateKey {
	return &PrivateKey{
		private: prv,
		public:  &prv.PublicKey,
		address: crypto.ECPublicKeyToAddress(&prv.PublicKey),
		sign:    crypto.ECSigner(prv),
		recover: crypto.ECRecoverer,
	}
}

// NewKeyFromBytes creates a new private key from private key bytes.
func NewKeyFromBytes(prv []byte) *PrivateKey {
	key, _ := btcec.PrivKeyFromBytes(prv)
	return NewKeyFromECDSA(key.ToECDSA())
}

// NewRandomKey creates a random private key.
func NewRandomKey() *PrivateKey {
	key, err := ecdsa.GenerateKey(s256, rand.Reader)
	if err != nil {
		panic(err)
	}
	return NewKeyFromECDSA(key)
}

// PublicKey returns the ECDSA public key.
func (k *PrivateKey) PublicKey() *ecdsa.PublicKey {
	return k.public
}

// PrivateKey returns the ECDSA private key.
func (k *PrivateKey) PrivateKey() *ecdsa.PrivateKey {
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
	return k.sign.SignHash(hash)
}

// SignMessage implements the Key interface.
func (k *PrivateKey) SignMessage(_ context.Context, data []byte) (*types.Signature, error) {
	return k.sign.SignMessage(data)
}

// SignTransaction implements the Key interface.
func (k *PrivateKey) SignTransaction(_ context.Context, tx *types.Transaction) error {
	return k.sign.SignTransaction(tx)
}

// VerifyHash implements the KeyWithHashSigner interface.
func (k *PrivateKey) VerifyHash(_ context.Context, hash types.Hash, sig types.Signature) bool {
	addr, err := k.recover.RecoverHash(hash, sig)
	if err != nil {
		return false
	}
	return *addr == k.address
}

// VerifyMessage implements the Key interface.
func (k *PrivateKey) VerifyMessage(_ context.Context, data []byte, sig types.Signature) bool {
	addr, err := k.recover.RecoverMessage(data, sig)
	if err != nil {
		return false
	}
	return *addr == k.address
}

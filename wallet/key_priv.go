package wallet

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/json"

	"github.com/btcsuite/btcd/btcec"

	"github.com/defiweb/go-eth/crypto"
	"github.com/defiweb/go-eth/types"
)

var s256 = btcec.S256()

type PrivateKey struct {
	private *ecdsa.PrivateKey
	public  *ecdsa.PublicKey
	address types.Address
}

// NewKeyFromECDSA creates a new private key from an ecdsa.PrivateKey.
func NewKeyFromECDSA(prv *ecdsa.PrivateKey) *PrivateKey {
	return &PrivateKey{
		private: prv,
		public:  &prv.PublicKey,
		address: crypto.ECPublicKeyToAddress(&prv.PublicKey),
	}
}

// NewKeyFromBytes creates a new private key from private key bytes.
func NewKeyFromBytes(prv []byte) *PrivateKey {
	key, _ := btcec.PrivKeyFromBytes(s256, prv)
	return NewKeyFromECDSA((*ecdsa.PrivateKey)(key))
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

// SignHash implements the Key interface.
func (k *PrivateKey) SignHash(hash types.Hash) (*types.Signature, error) {
	return crypto.ECSignHash(k.private, hash)
}

// SignMessage implements the Key interface.
func (k *PrivateKey) SignMessage(data []byte) (*types.Signature, error) {
	return crypto.ECSignMessage(k.private, data)
}

// SignTransaction implements the Key interface.
func (k *PrivateKey) SignTransaction(tx *types.Transaction) error {
	return crypto.ECSignTransaction(k.private, tx)
}

// VerifyHash implements the Key interface.
func (k *PrivateKey) VerifyHash(hash types.Hash, sig types.Signature) bool {
	addr, err := crypto.ECRecoverHash(hash, sig)
	if err != nil {
		return false
	}
	return *addr == k.address
}

// VerifyMessage implements the Key interface.
func (k *PrivateKey) VerifyMessage(data []byte, sig types.Signature) bool {
	addr, err := crypto.ECRecoverMessage(data, sig)
	if err != nil {
		return false
	}
	return *addr == k.address
}

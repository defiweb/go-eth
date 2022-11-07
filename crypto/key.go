package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"

	"github.com/btcsuite/btcd/btcec"

	"web3rpc/types"
)

var s256 = btcec.S256()

type Key struct {
	private *ecdsa.PrivateKey
	public  *ecdsa.PublicKey
	address types.Address
}

// NewKeyFromECDSA creates a new key with a private key
func NewKeyFromECDSA(priv *ecdsa.PrivateKey) *Key {
	return &Key{
		private: priv,
		public:  &priv.PublicKey,
		address: pubKeyToAddress(&priv.PublicKey),
	}
}

func NewKeyFromBytes(b []byte) *Key {
	priv, _ := btcec.PrivKeyFromBytes(s256, b)
	return NewKeyFromECDSA((*ecdsa.PrivateKey)(priv))
}

func NewRandomKey() *Key {
	priv, err := ecdsa.GenerateKey(s256, rand.Reader)
	if err != nil {
		panic(err)
	}
	return NewKeyFromECDSA(priv)
}

func (k *Key) Address() types.Address {
	return k.address
}

func (k *Key) Sign(hash types.Hash) (types.Signature, error) {
	sig, err := btcec.SignCompact(s256, (*btcec.PrivateKey)(k.private), hash.Bytes(), false)
	if err != nil {
		return types.Signature{}, err
	}
	v := sig[0] - 27
	copy(sig, sig[1:])
	sig[64] = v
	return types.BytesToSignature(sig), nil
}

func (k *Key) SignMessage(data []byte) (types.Signature, error) {
	return k.Sign(Keccak256(formatMessage(data)))
}

func (k *Key) SignTransaction(tx *types.Transaction) (*types.Transaction, error) {
	r, err := tx.SigningData()
	if err != nil {
		return nil, err
	}
	s, err := k.Sign(Keccak256(r))
	if err != nil {
		return nil, err
	}
	v := uint64(s[types.SignatureLength-1])
	if tx.Type == 0 {
		v = v + 35 + tx.ChainID.Uint64()*2
	}
	tx.Signature = s
	return tx, nil
}

func (k *Key) Verify(hash types.Hash, sig types.Signature) bool {
	addr, err := Ecrecover(hash, sig)
	if err != nil {
		return false
	}
	return addr == k.address
}

func (k *Key) VerifyMessage(data []byte, sig types.Signature) bool {
	return k.Verify(Keccak256(formatMessage(data)), sig)
}

func pubKeyToAddress(pub *ecdsa.PublicKey) (addr types.Address) {
	b := Keccak256(elliptic.Marshal(s256, pub.X, pub.Y)[1:])
	copy(addr[:], b[12:])
	return
}

package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"

	"github.com/btcsuite/btcd/btcec"

	"github.com/defiweb/go-eth/types"
)

var s256 = btcec.S256()

// PublicKeyToAddress returns the Ethereum address for the given ECDSA public key.
func PublicKeyToAddress(pub *ecdsa.PublicKey) (addr types.Address) {
	b := Keccak256(elliptic.Marshal(s256, pub.X, pub.Y)[1:])
	copy(addr[:], b[12:])
	return
}

// Ecrecover recovers the Ethereum address from a signed hash.
func Ecrecover(hash types.Hash, sig types.Signature) (types.Address, error) {
	v := sig[types.SignatureLength-1]
	copy(sig[1:], sig[:types.SignatureLength-1])
	sig[0] = v
	if v < 27 {
		sig[0] += 27
	}
	pub, _, err := btcec.RecoverCompact(s256, sig.Bytes(), hash.Bytes())
	if err != nil {
		return types.ZeroAddress, err
	}
	return PublicKeyToAddress(pub.ToECDSA()), nil
}

// EcrecoverMessage recovers the Ethereum address from a signed message.
func EcrecoverMessage(data []byte, sig types.Signature) (types.Address, error) {
	return Ecrecover(Keccak256(FormatMessage(data)), sig)
}

// EcrecoverTransaction recovers the Ethereum address from a signed transaction.
func EcrecoverTransaction(tx *types.Transaction) (types.Address, error) {
	d, err := tx.SigningHash(Keccak256)
	if err != nil {
		return types.ZeroAddress, err
	}
	addr, err := Ecrecover(d, tx.Signature)
	if err != nil {
		return types.ZeroAddress, err
	}
	return addr, nil
}

// Sign signs the given hash with the given key.
func Sign(key *ecdsa.PrivateKey, hash types.Hash) (types.Signature, error) {
	sig, err := btcec.SignCompact(s256, (*btcec.PrivateKey)(key), hash.Bytes(), false)
	if err != nil {
		return types.Signature{}, err
	}
	v := sig[0] - 27
	copy(sig, sig[1:])
	sig[64] = v
	return types.BytesToSignature(sig), nil
}

// SignMessage signs the given message with the given key.
func SignMessage(key *ecdsa.PrivateKey, data []byte) (types.Signature, error) {
	return Sign(key, Keccak256(FormatMessage(data)))
}

// FormatMessage add the Ethereum message prefix to the given data.
func FormatMessage(data []byte) []byte {
	return []byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data))
}

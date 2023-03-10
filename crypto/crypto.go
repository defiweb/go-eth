package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"math/big"

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
	if sig.V.BitLen() > 8 {
		return types.ZeroAddress, fmt.Errorf("invalid signature V: %d", sig.V)
	}
	v := byte(sig.V.Uint64())
	if v < 27 {
		v += 27
	}
	b := append([]byte{v}, append(sig.R.Bytes(), sig.S.Bytes()...)...)
	pub, _, err := btcec.RecoverCompact(s256, b, hash.Bytes())
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
	if tx.Signature == nil {
		return types.ZeroAddress, fmt.Errorf("signature is missing")
	}
	sig := *tx.Signature
	if tx.Type == types.LegacyTxType && tx.Signature.V.Cmp(big.NewInt(35)) >= 0 {
		x := new(big.Int).Sub(sig.V, big.NewInt(35))

		// Derive the chain ID from the signature.
		chainID := new(big.Int).Div(x, big.NewInt(2))
		if tx.ChainID != nil && *tx.ChainID != chainID.Uint64() {
			return types.ZeroAddress, fmt.Errorf("invalid chain ID: %d", chainID)
		}

		// For legacy transactions with specified CHAIN_ID, the signature V is recalculated as follows:
		// V = CHAIN_ID * 2 + 35 + (V - 27)
		//
		// Because V is always 27 or 28, we can use following formula to derive the original V:
		// V = (V - 35) % 2 + 27
		sig.V = new(big.Int).Add(new(big.Int).Mod(x, big.NewInt(2)), big.NewInt(27))
	}
	hash, err := tx.SigningHash(Keccak256)
	if err != nil {
		return types.ZeroAddress, err
	}
	addr, err := Ecrecover(hash, sig)
	if err != nil {
		return types.ZeroAddress, err
	}
	return addr, nil
}

// SignHash signs the given hash with the given key.
func SignHash(key *ecdsa.PrivateKey, hash types.Hash) (types.Signature, error) {
	sig, err := btcec.SignCompact(s256, (*btcec.PrivateKey)(key), hash.Bytes(), false)
	if err != nil {
		return types.Signature{}, err
	}
	v := sig[0]
	copy(sig, sig[1:])
	sig[64] = v
	return types.MustSignatureFromBytes(sig), nil
}

// SignMessage signs the given message with the given key.
func SignMessage(key *ecdsa.PrivateKey, data []byte) (types.Signature, error) {
	return SignHash(key, Keccak256(FormatMessage(data)))
}

// SignTransaction signs the given transaction with the given key.
func SignTransaction(key *ecdsa.PrivateKey, tx *types.Transaction) error {
	from := PublicKeyToAddress(&key.PublicKey)
	if tx.From != nil && *tx.From != from {
		return fmt.Errorf("invalid signer address: %s", tx.From)
	}
	r, err := tx.SigningHash(Keccak256)
	if err != nil {
		return err
	}
	s, err := SignHash(key, r)
	if err != nil {
		return err
	}
	sv, sr, ss := s.V, s.R, s.S
	if tx.Type == types.LegacyTxType {
		if tx.ChainID != nil {
			sv = new(big.Int).Sub(sv, big.NewInt(27))
			sv = new(big.Int).Add(sv, big.NewInt(35))
			sv = new(big.Int).Add(sv, new(big.Int).SetUint64(*tx.ChainID*2))
		}
	} else {
		sv = new(big.Int).Sub(sv, big.NewInt(27))
	}
	tx.From = &from
	tx.Signature = types.SignatureFromVRSPtr(sv, sr, ss)
	return nil
}

// FormatMessage adds the Ethereum message prefix to the given data.
func FormatMessage(data []byte) []byte {
	return []byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data))
}

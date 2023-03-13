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

// ECPublicKeyToAddress returns the Ethereum address for the given ECDSA public key.
func ECPublicKeyToAddress(pub *ecdsa.PublicKey) (addr types.Address) {
	b := Keccak256(elliptic.Marshal(s256, pub.X, pub.Y)[1:])
	copy(addr[:], b[12:])
	return
}

// ECSignHash signs the given hash with the given private key.
func ECSignHash(key *ecdsa.PrivateKey, hash types.Hash) (*types.Signature, error) {
	if key == nil {
		return nil, fmt.Errorf("missing private key")
	}
	sig, err := btcec.SignCompact(s256, (*btcec.PrivateKey)(key), hash.Bytes(), false)
	if err != nil {
		return nil, err
	}
	v := sig[0]
	copy(sig, sig[1:])
	sig[64] = v
	return types.SignatureFromBytesPtr(sig), nil
}

// ECSignMessage signs the given message with the given private key.
func ECSignMessage(key *ecdsa.PrivateKey, data []byte) (*types.Signature, error) {
	if key == nil {
		return nil, fmt.Errorf("missing private key")
	}
	return ECSignHash(key, Keccak256(AddMessagePrefix(data)))
}

// ECSignTransaction signs the given transaction with the given private key.
func ECSignTransaction(key *ecdsa.PrivateKey, tx *types.Transaction) error {
	if key == nil {
		return fmt.Errorf("missing private key")
	}
	from := ECPublicKeyToAddress(&key.PublicKey)
	if tx.From != nil && *tx.From != from {
		return fmt.Errorf("invalid signer address: %s", tx.From)
	}
	hash, err := signingHash(tx)
	if err != nil {
		return err
	}
	sig, err := ECSignHash(key, hash)
	if err != nil {
		return err
	}
	sv, sr, ss := sig.V, sig.R, sig.S
	switch tx.Type {
	case types.LegacyTxType:
		if tx.ChainID != nil {
			sv = new(big.Int).Sub(sv, big.NewInt(27))
			sv = new(big.Int).Add(sv, new(big.Int).SetUint64(*tx.ChainID*2))
			sv = new(big.Int).Add(sv, big.NewInt(35))
		}
	case types.AccessListTxType:
		sv = new(big.Int).Sub(sv, big.NewInt(27))
	case types.DynamicFeeTxType:
		sv = new(big.Int).Sub(sv, big.NewInt(27))
	default:
		return fmt.Errorf("unsupported transaction type: %d", tx.Type)
	}
	tx.From = &from
	tx.Signature = types.SignatureFromVRSPtr(sv, sr, ss)
	return nil
}

// ECRecoverHash recovers the Ethereum address from the given hash and signature.
func ECRecoverHash(hash types.Hash, sig types.Signature) (*types.Address, error) {
	if sig.V.BitLen() > 8 {
		return nil, fmt.Errorf("invalid signature V: %d", sig.V)
	}
	v := byte(sig.V.Uint64())
	bin := make([]byte, 65)
	bin[0] = v
	copy(bin[1:], sig.R.Bytes())
	copy(bin[33:], sig.S.Bytes())
	pub, _, err := btcec.RecoverCompact(s256, bin, hash.Bytes())
	if err != nil {
		return nil, err
	}
	addr := ECPublicKeyToAddress(pub.ToECDSA())
	return &addr, nil
}

// ECRecoverMessage recovers the Ethereum address from the given message and signature.
func ECRecoverMessage(data []byte, sig types.Signature) (*types.Address, error) {
	return ECRecoverHash(Keccak256(AddMessagePrefix(data)), sig)
}

// ECRecoverTransaction recovers the Ethereum address from the given transaction.
func ECRecoverTransaction(tx *types.Transaction) (*types.Address, error) {
	if tx.Signature == nil {
		return nil, fmt.Errorf("signature is missing")
	}
	sig := *tx.Signature
	switch tx.Type {
	case types.LegacyTxType:
		if tx.Signature.V.Cmp(big.NewInt(35)) >= 0 {
			x := new(big.Int).Sub(sig.V, big.NewInt(35))

			// Derive the chain ID from the signature.
			chainID := new(big.Int).Div(x, big.NewInt(2))
			if tx.ChainID != nil && *tx.ChainID != chainID.Uint64() {
				return nil, fmt.Errorf("invalid chain ID: %d", chainID)
			}

			// Derive the recovery ID from the signature.
			sig.V = new(big.Int).Add(new(big.Int).Mod(x, big.NewInt(2)), big.NewInt(27))
		}
	case types.AccessListTxType:
		sig.V = new(big.Int).Add(sig.V, big.NewInt(27))
	case types.DynamicFeeTxType:
		sig.V = new(big.Int).Add(sig.V, big.NewInt(27))
	default:
		return nil, fmt.Errorf("unsupported transaction type: %d", tx.Type)
	}
	hash, err := signingHash(tx)
	if err != nil {
		return nil, err
	}
	return ECRecoverHash(hash, sig)
}

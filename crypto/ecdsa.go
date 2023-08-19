package crypto

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
	ecdsa2 "github.com/btcsuite/btcd/btcec/v2/ecdsa"

	"github.com/defiweb/go-eth/types"
)

// ECPublicKeyToAddress returns the Ethereum address for the given ECDSA public key.
func ECPublicKeyToAddress(pub *ecdsa.PublicKey) (addr types.Address) {
	publicKey, err := pub.ECDH()
	if err != nil {
		panic(err)
	}
	b := Keccak256(publicKey.Bytes())
	copy(addr[:], b[12:])
	return
}

// ecSignHash signs the given hash with the given private key.
func ecSignHash(key *ecdsa.PrivateKey, hash types.Hash) (*types.Signature, error) {
	if key == nil {
		return nil, fmt.Errorf("missing private key")
	}
	privKey, _ := btcec.PrivKeyFromBytes(key.D.Bytes())
	sig, err := ecdsa2.SignCompact(privKey, hash.Bytes(), false)
	if err != nil {
		return nil, err
	}
	v := sig[0]
	switch v {
	case 27, 28:
		v -= 27
	}
	copy(sig, sig[1:])
	sig[64] = v
	return types.SignatureFromBytesPtr(sig), nil
}

// ecSignMessage signs the given message with the given private key.
func ecSignMessage(key *ecdsa.PrivateKey, data []byte) (*types.Signature, error) {
	if key == nil {
		return nil, fmt.Errorf("missing private key")
	}
	sig, err := ecSignHash(key, Keccak256(AddMessagePrefix(data)))
	if err != nil {
		return nil, err
	}
	sig.V = new(big.Int).Add(sig.V, big.NewInt(27))
	return sig, nil
}

// ecSignTransaction signs the given transaction with the given private key.
func ecSignTransaction(key *ecdsa.PrivateKey, tx *types.Transaction) error {
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
	sig, err := ecSignHash(key, hash)
	if err != nil {
		return err
	}
	sv, sr, ss := sig.V, sig.R, sig.S
	switch tx.Type {
	case types.LegacyTxType:
		if tx.ChainID != nil {
			sv = new(big.Int).Add(sv, new(big.Int).SetUint64(*tx.ChainID*2))
			sv = new(big.Int).Add(sv, big.NewInt(35))
		} else {
			sv = new(big.Int).Add(sv, big.NewInt(27))
		}
	case types.AccessListTxType:
	case types.DynamicFeeTxType:
	default:
		return fmt.Errorf("unsupported transaction type: %d", tx.Type)
	}
	tx.From = &from
	tx.Signature = types.SignatureFromVRSPtr(sv, sr, ss)
	return nil
}

// ecRecoverHash recovers the Ethereum address from the given hash and signature.
func ecRecoverHash(hash types.Hash, sig types.Signature) (*types.Address, error) {
	if sig.V.BitLen() > 8 {
		return nil, errors.New("invalid signature: V has more than 8 bits")
	}
	if sig.R.BitLen() > 256 {
		return nil, errors.New("invalid signature: R has more than 256 bits")
	}
	if sig.S.BitLen() > 256 {
		return nil, errors.New("invalid signature: S has more than 256 bits")
	}
	v := byte(sig.V.Uint64())
	switch v {
	case 0, 1:
		v += 27
	}
	rb := sig.R.Bytes()
	sb := sig.S.Bytes()
	bin := make([]byte, 65)
	bin[0] = v
	copy(bin[1+(32-len(rb)):], rb)
	copy(bin[33+(32-len(sb)):], sb)
	pub, _, err := ecdsa2.RecoverCompact(bin, hash.Bytes())
	if err != nil {
		return nil, err
	}
	addr := ECPublicKeyToAddress(pub.ToECDSA())
	return &addr, nil
}

// ecRecoverMessage recovers the Ethereum address from the given message and signature.
func ecRecoverMessage(data []byte, sig types.Signature) (*types.Address, error) {
	sig.V = new(big.Int).Sub(sig.V, big.NewInt(27))
	return ecRecoverHash(Keccak256(AddMessagePrefix(data)), sig)
}

// ecRecoverTransaction recovers the Ethereum address from the given transaction.
func ecRecoverTransaction(tx *types.Transaction) (*types.Address, error) {
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

			// Derive the recovery byte from the signature.
			sig.V = new(big.Int).Add(new(big.Int).Mod(x, big.NewInt(2)), big.NewInt(27))
		} else {
			sig.V = new(big.Int).Sub(sig.V, big.NewInt(27))
		}
	case types.AccessListTxType:
	case types.DynamicFeeTxType:
	default:
		return nil, fmt.Errorf("unsupported transaction type: %d", tx.Type)
	}
	hash, err := signingHash(tx)
	if err != nil {
		return nil, err
	}
	return ecRecoverHash(hash, sig)
}

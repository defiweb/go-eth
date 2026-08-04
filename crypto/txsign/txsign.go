// Package txsign provides transaction signing and address recovery
// functionality for Ethereum transactions.
package txsign

import (
	"fmt"
	"math/big"

	"github.com/defiweb/go-eth/crypto"
	"github.com/defiweb/go-eth/types"
)

// Sign signs the given transaction with the given private key.
func Sign(key crypto.PrivateKey, tx types.SignableTransaction) error {
	if key.IsZero() {
		return fmt.Errorf("invalid private key")
	}
	if tx == nil {
		return fmt.Errorf("missing transaction")
	}
	sd := tx.GetSigningData()
	ed := types.GetExecutionData(sd)
	hash, err := tx.SigningHash()
	if err != nil {
		return err
	}
	sig, err := crypto.ECSignHash(key, crypto.Hash(hash))
	if err != nil {
		return err
	}
	sv, sr, ss := sig.V, sig.R, sig.S
	if tx.Type() == types.LegacyTxType {
		if sd.ChainID != nil {
			sv = new(big.Int).Add(sv, new(big.Int).SetUint64(*sd.ChainID*2))
			sv = new(big.Int).Add(sv, big.NewInt(35))
		} else {
			sv = new(big.Int).Add(sv, big.NewInt(27))
		}
	}
	sd.SetSignature(types.SignatureFromVRS(sv, sr, ss))
	if ed != nil {
		ed.SetFrom(types.Address(crypto.ECPublicKeyToAddress(crypto.ECPrivateKeyToPublicKey(key))))
	}
	return nil
}

// Recover recovers the Ethereum address from the given transaction's
// signature.
func Recover(tx types.SignableTransaction) (*types.Address, error) {
	std := tx.GetSigningData()
	if std.Signature == nil {
		return nil, fmt.Errorf("signature is missing")
	}
	sig := *std.Signature
	if tx.Type() == types.LegacyTxType {
		if sig.V.Cmp(big.NewInt(35)) >= 0 {
			x := new(big.Int).Sub(sig.V, big.NewInt(35))

			// Derive the chain ID from the signature.
			chainID := new(big.Int).Div(x, big.NewInt(2))
			if std.ChainID != nil && *std.ChainID != chainID.Uint64() {
				return nil, fmt.Errorf("invalid chain ID: %d", chainID)
			}

			// Derive the recovery byte from the signature.
			sig.V = new(big.Int).Add(new(big.Int).Mod(x, big.NewInt(2)), big.NewInt(27))
		} else {
			sig.V = new(big.Int).Sub(sig.V, big.NewInt(27))
		}
	}
	hash, err := tx.SigningHash()
	if err != nil {
		return nil, err
	}
	addr, err := crypto.ECRecoverHash(crypto.Hash(hash), crypto.Signature(sig))
	if err != nil {
		return nil, err
	}
	return (*types.Address)(addr), nil
}

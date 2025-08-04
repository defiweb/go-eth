// Package txsign provides transaction signing and address recovery
// functionality for Ethereum transactions.
package txsign

import (
	"fmt"
	"math/big"

	"github.com/defiweb/go-eth/crypto"
	"github.com/defiweb/go-eth/crypto/ecdsa"
	"github.com/defiweb/go-eth/types"
)

// Sign signs the given transaction with the given private key.
func Sign(key *ecdsa.PrivateKey, tx types.Transaction) error {
	if key == nil {
		return fmt.Errorf("missing private key")
	}
	txd := tx.GetTransactionData()
	hash, err := tx.CalculateSigningHash()
	if err != nil {
		return err
	}
	sig, err := crypto.ECSignHash(key, ecdsa.Hash(hash))
	if err != nil {
		return err
	}
	sv, sr, ss := sig.V, sig.R, sig.S
	if tx.Type() == types.LegacyTxType {
		if txd.ChainID != nil {
			sv = new(big.Int).Add(sv, new(big.Int).SetUint64(*txd.ChainID*2))
			sv = new(big.Int).Add(sv, big.NewInt(35))
		} else {
			sv = new(big.Int).Add(sv, big.NewInt(27))
		}
	}
	txd.SetSignature(types.SignatureFromVRS(sv, sr, ss))
	if cd, ok := tx.(types.HasCallData); ok {
		cd.GetCallData().SetFrom(types.Address(crypto.ECPublicKeyToAddress(crypto.ECPrivateKeyToPublicKey(key))))
	}
	return nil
}

// Recover recovers the Ethereum address from the given transaction's
// signature.
func Recover(tx types.Transaction) (*types.Address, error) {
	txd := tx.GetTransactionData()
	if txd.Signature == nil {
		return nil, fmt.Errorf("signature is missing")
	}
	sig := *txd.Signature
	if tx.Type() == types.LegacyTxType {
		if sig.V.Cmp(big.NewInt(35)) >= 0 {
			x := new(big.Int).Sub(sig.V, big.NewInt(35))

			// Derive the chain ID from the signature.
			chainID := new(big.Int).Div(x, big.NewInt(2))
			if txd.ChainID != nil && *txd.ChainID != chainID.Uint64() {
				return nil, fmt.Errorf("invalid chain ID: %d", chainID)
			}

			// Derive the recovery byte from the signature.
			sig.V = new(big.Int).Add(new(big.Int).Mod(x, big.NewInt(2)), big.NewInt(27))
		} else {
			sig.V = new(big.Int).Sub(sig.V, big.NewInt(27))
		}
	}
	hash, err := tx.CalculateSigningHash()
	if err != nil {
		return nil, err
	}
	addr, err := crypto.ECRecoverHash(ecdsa.Hash(hash), ecdsa.Signature(sig))
	if err != nil {
		return nil, err
	}
	return (*types.Address)(addr), nil
}

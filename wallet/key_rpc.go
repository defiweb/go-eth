package wallet

import (
	"context"
	"fmt"

	"github.com/defiweb/go-eth/crypto"
	"github.com/defiweb/go-eth/crypto/txsign"
	"github.com/defiweb/go-eth/types"
)

// RPCSigningClient is the interface for an Ethereum RPC client that can
// sign messages and transactions.
type RPCSigningClient interface {
	Sign(ctx context.Context, account types.Address, data []byte) (*types.Signature, error)
	SignTransaction(ctx context.Context, tx types.Transaction) ([]byte, error)
}

// KeyRPC is an Ethereum key that uses an RPC client to sign messages and transactions.
type KeyRPC struct {
	client  RPCSigningClient
	address types.Address
	decoder types.RLPTransactionDecoder
}

// NewKeyRPC returns a new KeyRPC.
func NewKeyRPC(client RPCSigningClient, address types.Address) *KeyRPC {
	return &KeyRPC{
		client:  client,
		address: address,
		decoder: types.DefaultTransactionDecoder,
	}
}

// Address implements the Key interface.
func (k *KeyRPC) Address() types.Address {
	return k.address
}

// SignMessage implements the Key interface.
func (k *KeyRPC) SignMessage(ctx context.Context, data []byte) (*types.Signature, error) {
	return k.client.Sign(ctx, k.address, data)
}

// SignTransaction implements the Key interface.
func (k *KeyRPC) SignTransaction(ctx context.Context, tx types.SignableTransaction) error {
	raw, err := k.client.SignTransaction(ctx, tx)
	if err != nil {
		return err
	}
	stx, err := k.decoder.DecodeRLP(raw)
	if err != nil {
		return fmt.Errorf("failed to decode signed transaction: %w", err)
	}
	sd := types.GetSigningData(stx)
	if sd == nil {
		return fmt.Errorf("failed to get signed transaction data")
	}
	tx.SetSigningData(*sd)
	addr, err := txsign.Recover(tx)
	if err != nil {
		return fmt.Errorf("failed to verify signed transaction: %w", err)
	}
	if *addr != k.address {
		return fmt.Errorf("failed to verify signed transaction: recovered address does not match key address")
	}
	return nil
}

// VerifyMessage implements the Key interface.
func (k *KeyRPC) VerifyMessage(_ context.Context, data []byte, sig types.Signature) bool {
	addr, err := crypto.ECRecoverMessage(data, crypto.Signature(sig))
	if err != nil {
		return false
	}
	return types.Address(*addr) == k.address
}

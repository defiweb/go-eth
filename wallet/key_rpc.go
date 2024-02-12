package wallet

import (
	"context"

	"github.com/defiweb/go-eth/crypto"
	"github.com/defiweb/go-eth/types"
)

// RPCSigningClient is the interface for an Ethereum RPC client that can
// sign messages and transactions.
type RPCSigningClient interface {
	Sign(ctx context.Context, account types.Address, data []byte) (*types.Signature, error)
	SignTransaction(ctx context.Context, tx *types.Transaction) ([]byte, *types.Transaction, error)
}

// KeyRPC is an Ethereum key that uses an RPC client to sign messages and transactions.
type KeyRPC struct {
	client  RPCSigningClient
	address types.Address
	recover crypto.Recoverer
}

// NewKeyRPC returns a new KeyRPC.
func NewKeyRPC(client RPCSigningClient, address types.Address) *KeyRPC {
	return &KeyRPC{
		client:  client,
		address: address,
		recover: crypto.ECRecoverer,
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
func (k *KeyRPC) SignTransaction(ctx context.Context, tx *types.Transaction) error {
	_, signedTX, err := k.client.SignTransaction(ctx, tx)
	if err != nil {
		return err
	}
	*tx = *signedTX
	return err
}

// VerifyMessage implements the Key interface.
func (k *KeyRPC) VerifyMessage(_ context.Context, data []byte, sig types.Signature) bool {
	addr, err := k.recover.RecoverMessage(data, sig)
	if err != nil {
		return false
	}
	return *addr == k.address
}

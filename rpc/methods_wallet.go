package rpc

import (
	"context"
	"errors"

	"github.com/defiweb/go-eth/types"
)

// MethodsWallet is a collection of RPC methods that require private keys to
// perform operations.
//
// Note: Public JSON-RPC APIs do not support these methods.
type MethodsWallet struct {
	Context *ClientContext
}

// Accounts performs eth_accounts RPC call.
//
// It returns the list of addresses owned by the client.
func (c *MethodsWallet) Accounts(ctx context.Context) ([]types.Address, error) {
	var res []types.Address
	if err := c.Context.Transport.Call(ctx, &res, "eth_accounts"); err != nil {
		return nil, err
	}
	return res, nil
}

// Sign performs eth_sign RPC call.
//
// It signs the given data with the given address.
func (c *MethodsWallet) Sign(ctx context.Context, account types.Address, data []byte) (*types.Signature, error) {
	var res types.Signature
	if err := c.Context.Transport.Call(ctx, &res, "eth_sign", account, types.Bytes(data)); err != nil {
		return nil, err
	}
	return &res, nil
}

// SignTransaction performs eth_signTransaction RPC call.
//
// It signs the given transaction and returns the raw transaction data.
func (c *MethodsWallet) SignTransaction(ctx context.Context, tx types.Transaction) ([]byte, error) {
	if tx == nil {
		return nil, errors.New("rpc client: transaction is nil")
	}
	var res signTransactionResult
	if err := c.Context.Transport.Call(ctx, &res, "eth_signTransaction", tx); err != nil {
		return nil, err
	}
	return res.Raw, nil
}

// SendTransaction performs eth_sendTransaction RPC call.
//
// It sends a transaction to the network.
func (c *MethodsWallet) SendTransaction(ctx context.Context, tx types.Transaction) (*types.Hash, error) {
	if tx == nil {
		return nil, errors.New("rpc client: transaction is nil")
	}
	var res types.Hash
	if err := c.Context.Transport.Call(ctx, &res, "eth_sendTransaction", tx); err != nil {
		return nil, err
	}
	return &res, nil
}

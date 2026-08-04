package rpc

import (
	"context"

	"github.com/defiweb/go-eth/types"
)

// MethodsPrivateTransaction provides unofficial RPC methods for private
// transactions.
type MethodsPrivateTransaction struct {
	Context *ClientContext
}

func (c *MethodsPrivateTransaction) CancelPrivateTransaction(ctx context.Context, hash types.Hash) (bool, error) {
	var res bool
	if err := c.Context.Transport.Call(ctx, &res, "eth_cancelPrivateTransaction", hash); err != nil {
		return false, err
	}
	return res, nil
}

func (c *MethodsPrivateTransaction) SendPrivateTransaction(ctx context.Context, data []byte) (*types.Hash, error) {
	var res types.Hash
	if err := c.Context.Transport.Call(ctx, &res, "eth_sendPrivateTransaction", types.Bytes(data)); err != nil {
		return nil, err
	}
	return &res, nil
}

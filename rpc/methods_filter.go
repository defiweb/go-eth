package rpc

import (
	"context"
	"math/big"

	"github.com/defiweb/go-eth/types"
)

// MethodsFilter is a collection of RPC methods to interact with filters.
//
// Note: Some JSON-RPC APIs do not support these methods.
type MethodsFilter struct {
	Context *ClientContext
}

// NewFilter implements the RPC interface.
func (c *MethodsFilter) NewFilter(ctx context.Context, query *types.FilterLogsQuery) (*big.Int, error) {
	var res *types.Number
	if err := c.Context.Transport.Call(ctx, &res, "eth_newFilter", query); err != nil {
		return nil, err
	}
	return res.Big(), nil
}

// NewBlockFilter implements the RPC interface.
func (c *MethodsFilter) NewBlockFilter(ctx context.Context) (*big.Int, error) {
	var res *types.Number
	if err := c.Context.Transport.Call(ctx, &res, "eth_newBlockFilter"); err != nil {
		return nil, err
	}
	return res.Big(), nil

}

// NewPendingTransactionFilter implements the RPC interface.
func (c *MethodsFilter) NewPendingTransactionFilter(ctx context.Context) (*big.Int, error) {
	var res *types.Number
	if err := c.Context.Transport.Call(ctx, &res, "eth_newPendingTransactionFilter"); err != nil {
		return nil, err
	}
	return res.Big(), nil
}

// UninstallFilter implements the RPC interface.
func (c *MethodsFilter) UninstallFilter(ctx context.Context, id *big.Int) (bool, error) {
	var res bool
	if err := c.Context.Transport.Call(ctx, &res, "eth_uninstallFilter", types.NumberFromBigInt(id)); err != nil {
		return false, err
	}
	return res, nil
}

// GetFilterChanges implements the RPC interface.
func (c *MethodsFilter) GetFilterChanges(ctx context.Context, id *big.Int) ([]types.Log, error) {
	var res []types.Log
	if err := c.Context.Transport.Call(ctx, &res, "eth_getFilterChanges", types.NumberFromBigInt(id)); err != nil {
		return nil, err
	}
	return res, nil
}

// GetFilterLogs implements the RPC interface.
func (c *MethodsFilter) GetFilterLogs(ctx context.Context, id *big.Int) ([]types.Log, error) {
	var res []types.Log
	if err := c.Context.Transport.Call(ctx, &res, "eth_getFilterLogs", types.NumberFromBigInt(id)); err != nil {
		return nil, err
	}
	return res, nil
}

// GetBlockFilterChanges implements the RPC interface.
func (c *MethodsFilter) GetBlockFilterChanges(ctx context.Context, id *big.Int) ([]types.Hash, error) {
	var res []types.Hash
	if err := c.Context.Transport.Call(ctx, &res, "eth_getFilterChanges", types.NumberFromBigInt(id)); err != nil {
		return nil, err
	}
	return res, nil
}

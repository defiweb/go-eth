package rpc

import (
	"context"
	"fmt"
	"math/big"

	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
)

type baseClient struct {
	transport transport.Transport
}

// ChainID implements the RPC interface.
func (c *baseClient) ChainID(ctx context.Context) (uint64, error) {
	var res types.Number
	if err := c.transport.Call(ctx, &res, "eth_chainId"); err != nil {
		return 0, err
	}
	if !res.Big().IsUint64() {
		return 0, fmt.Errorf("chain id is too big")
	}
	return res.Big().Uint64(), nil
}

// GasPrice implements the RPC interface.
func (c *baseClient) GasPrice(ctx context.Context) (*big.Int, error) {
	var res types.Number
	if err := c.transport.Call(ctx, &res, "eth_gasPrice"); err != nil {
		return nil, err
	}
	return res.Big(), nil
}

// Accounts implements the RPC interface.
func (c *baseClient) Accounts(ctx context.Context) ([]types.Address, error) {
	var res []types.Address
	if err := c.transport.Call(ctx, &res, "eth_accounts"); err != nil {
		return nil, err
	}
	return res, nil
}

// BlockNumber implements the RPC interface.
func (c *baseClient) BlockNumber(ctx context.Context) (*big.Int, error) {
	var res types.Number
	if err := c.transport.Call(ctx, &res, "eth_blockNumber"); err != nil {
		return nil, err
	}
	return res.Big(), nil
}

// GetBalance implements the RPC interface.
func (c *baseClient) GetBalance(ctx context.Context, address types.Address, block types.BlockNumber) (*big.Int, error) {
	var res types.Number
	if err := c.transport.Call(ctx, &res, "eth_getBalance", address, block); err != nil {
		return nil, err
	}
	return res.Big(), nil
}

// GetStorageAt implements the RPC interface.
func (c *baseClient) GetStorageAt(ctx context.Context, account types.Address, key types.Hash, block types.BlockNumber) (*types.Hash, error) {
	var res types.Hash
	if err := c.transport.Call(ctx, &res, "eth_getStorageAt", account, key, block); err != nil {
		return nil, err
	}
	return &res, nil
}

// GetTransactionCount implements the RPC interface.
func (c *baseClient) GetTransactionCount(ctx context.Context, account types.Address, block types.BlockNumber) (uint64, error) {
	var res types.Number
	if err := c.transport.Call(ctx, &res, "eth_getTransactionCount", account, block); err != nil {
		return 0, err
	}
	return res.Big().Uint64(), nil
}

// GetBlockTransactionCountByHash implements the RPC interface.
func (c *baseClient) GetBlockTransactionCountByHash(ctx context.Context, hash types.Hash) (uint64, error) {
	var res types.Number
	if err := c.transport.Call(ctx, &res, "eth_getBlockTransactionCountByHash", hash); err != nil {
		return 0, err
	}
	return res.Big().Uint64(), nil
}

// GetBlockTransactionCountByNumber implements the RPC interface.
func (c *baseClient) GetBlockTransactionCountByNumber(ctx context.Context, number types.BlockNumber) (uint64, error) {
	var res types.Number
	if err := c.transport.Call(ctx, &res, "eth_getBlockTransactionCountByNumber", number); err != nil {
		return 0, err
	}
	return res.Big().Uint64(), nil
}

// GetUncleCountByBlockHash implements the RPC interface.
func (c *baseClient) GetUncleCountByBlockHash(ctx context.Context, hash types.Hash) (uint64, error) {
	var res types.Number
	if err := c.transport.Call(ctx, &res, "eth_getUncleCountByBlockHash", hash); err != nil {
		return 0, err
	}
	return res.Big().Uint64(), nil
}

// GetUncleCountByBlockNumber implements the RPC interface.
func (c *baseClient) GetUncleCountByBlockNumber(ctx context.Context, number types.BlockNumber) (uint64, error) {
	var res types.Number
	if err := c.transport.Call(ctx, &res, "eth_getUncleCountByBlockNumber", number); err != nil {
		return 0, err
	}
	return res.Big().Uint64(), nil
}

// GetCode implements the RPC interface.
func (c *baseClient) GetCode(ctx context.Context, account types.Address, block types.BlockNumber) ([]byte, error) {
	var res types.Bytes
	if err := c.transport.Call(ctx, &res, "eth_getCode", account, block); err != nil {
		return nil, err
	}
	return res.Bytes(), nil
}

// Sign implements the RPC interface.
func (c *baseClient) Sign(ctx context.Context, account types.Address, data []byte) (*types.Signature, error) {
	var res types.Signature
	if err := c.transport.Call(ctx, &res, "eth_sign", account, types.Bytes(data)); err != nil {
		return nil, err
	}
	return &res, nil
}

// SignTransaction implements the RPC interface.
func (c *baseClient) SignTransaction(ctx context.Context, tx types.Transaction) ([]byte, *types.Transaction, error) {
	var res signTransactionResult
	if err := c.transport.Call(ctx, &res, "eth_signTransaction", tx); err != nil {
		return nil, nil, err
	}
	return res.Raw, res.Tx, nil
}

// SendTransaction implements the RPC interface.
func (c *baseClient) SendTransaction(ctx context.Context, tx types.Transaction) (*types.Hash, error) {
	var res types.Hash
	if err := c.transport.Call(ctx, &res, "eth_sendTransaction", tx); err != nil {
		return nil, err
	}
	return &res, nil
}

// SendRawTransaction implements the RPC interface.
func (c *baseClient) SendRawTransaction(ctx context.Context, data []byte) (*types.Hash, error) {
	var res types.Hash
	if err := c.transport.Call(ctx, &res, "eth_sendRawTransaction", types.Bytes(data)); err != nil {
		return nil, err
	}
	return &res, nil
}

// Call implements the RPC interface.
func (c *baseClient) Call(ctx context.Context, call types.Call, block types.BlockNumber) ([]byte, error) {
	var res types.Bytes
	if err := c.transport.Call(ctx, &res, "eth_call", call, block); err != nil {
		return nil, err
	}
	return res, nil
}

// EstimateGas implements the RPC interface.
func (c *baseClient) EstimateGas(ctx context.Context, call types.Call, block types.BlockNumber) (uint64, error) {
	var res types.Number
	if err := c.transport.Call(ctx, &res, "eth_estimateGas", call, block); err != nil {
		return 0, err
	}
	return res.Big().Uint64(), nil
}

// BlockByHash implements the RPC interface.
func (c *baseClient) BlockByHash(ctx context.Context, hash types.Hash, full bool) (*types.Block, error) {
	var res types.Block
	if err := c.transport.Call(ctx, &res, "eth_getBlockByHash", hash, full); err != nil {
		return nil, err
	}
	return &res, nil
}

// BlockByNumber implements the RPC interface.
func (c *baseClient) BlockByNumber(ctx context.Context, number types.BlockNumber, full bool) (*types.Block, error) {
	var res types.Block
	if err := c.transport.Call(ctx, &res, "eth_getBlockByNumber", number, full); err != nil {
		return nil, err
	}
	return &res, nil
}

// GetTransactionByHash implements the RPC interface.
func (c *baseClient) GetTransactionByHash(ctx context.Context, hash types.Hash) (*types.OnChainTransaction, error) {
	var res types.OnChainTransaction
	if err := c.transport.Call(ctx, &res, "eth_getTransactionByHash", hash); err != nil {
		return nil, err
	}
	return &res, nil
}

// GetTransactionByBlockHashAndIndex implements the RPC interface.
func (c *baseClient) GetTransactionByBlockHashAndIndex(ctx context.Context, hash types.Hash, index uint64) (*types.OnChainTransaction, error) {
	var res types.OnChainTransaction
	if err := c.transport.Call(ctx, &res, "eth_getTransactionByBlockHashAndIndex", hash, types.NumberFromUint64(index)); err != nil {
		return nil, err
	}
	return &res, nil
}

// GetTransactionByBlockNumberAndIndex implements the RPC interface.
func (c *baseClient) GetTransactionByBlockNumberAndIndex(ctx context.Context, number types.BlockNumber, index uint64) (*types.OnChainTransaction, error) {
	var res types.OnChainTransaction
	if err := c.transport.Call(ctx, &res, "eth_getTransactionByBlockNumberAndIndex", number, types.NumberFromUint64(index)); err != nil {
		return nil, err
	}
	return &res, nil
}

// GetTransactionReceipt implements the RPC interface.
func (c *baseClient) GetTransactionReceipt(ctx context.Context, hash types.Hash) (*types.TransactionReceipt, error) {
	var res types.TransactionReceipt
	if err := c.transport.Call(ctx, &res, "eth_getTransactionReceipt", hash); err != nil {
		return nil, err
	}
	return &res, nil
}

// GetLogs implements the RPC interface.
func (c *baseClient) GetLogs(ctx context.Context, query types.FilterLogsQuery) ([]types.Log, error) {
	var res []types.Log
	if err := c.transport.Call(ctx, &res, "eth_getLogs", query); err != nil {
		return nil, err
	}
	return res, nil
}

// MaxPriorityFeePerGas implements the RPC interface.
func (c *baseClient) MaxPriorityFeePerGas(ctx context.Context) (*big.Int, error) {
	var res types.Number
	if err := c.transport.Call(ctx, &res, "eth_maxPriorityFeePerGas"); err != nil {
		return nil, err
	}
	return res.Big(), nil
}

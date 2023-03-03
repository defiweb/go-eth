package rpc

import (
	"context"

	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
)

type Client struct {
	transport transport.Transport
}

func NewClient(transport transport.Transport) *Client {
	return &Client{transport: transport}
}

func (c *Client) GasPrice(ctx context.Context) (uint64, error) {
	var res types.Number
	if err := c.transport.Call(ctx, &res, "eth_gasPrice"); err != nil {
		return 0, err
	}
	return res.Big().Uint64(), nil
}

func (c *Client) BlockNumber(ctx context.Context) (uint64, error) {
	var res types.Number
	if err := c.transport.Call(ctx, &res, "eth_blockNumber"); err != nil {
		return 0, err
	}
	return res.Big().Uint64(), nil
}

func (c *Client) GetBalance(ctx context.Context, address types.Address, block types.BlockNumber) (uint64, error) {
	var res types.Number
	if err := c.transport.Call(ctx, &res, "eth_getBalance", address, block); err != nil {
		return 0, err
	}
	return res.Big().Uint64(), nil
}

func (c *Client) GetStorageAt(ctx context.Context, account types.Address, key types.Hash, block types.BlockNumber) (*types.Hash, error) {
	var res types.Hash
	if err := c.transport.Call(ctx, &res, "eth_getStorageAt", account, key, block); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetTransactionCount(ctx context.Context, account types.Address, block types.BlockNumber) (uint64, error) {
	var res types.Number
	if err := c.transport.Call(ctx, &res, "eth_getTransactionCount", account, block); err != nil {
		return 0, err
	}
	return res.Big().Uint64(), nil
}

func (c *Client) GetBlockTransactionCountByHash(ctx context.Context, hash types.Hash) (uint64, error) {
	var res types.Number
	if err := c.transport.Call(ctx, &res, "eth_getBlockTransactionCountByHash", hash); err != nil {
		return 0, err
	}
	return res.Big().Uint64(), nil
}

func (c *Client) GetBlockTransactionCountByNumber(ctx context.Context, number types.BlockNumber) (uint64, error) {
	var res types.Number
	if err := c.transport.Call(ctx, &res, "eth_getBlockTransactionCountByNumber", number); err != nil {
		return 0, err
	}
	return res.Big().Uint64(), nil
}

func (c *Client) GetUncleCountByBlockHash(ctx context.Context, hash types.Hash) (uint64, error) {
	var res types.Number
	if err := c.transport.Call(ctx, &res, "eth_getUncleCountByBlockHash", hash); err != nil {
		return 0, err
	}
	return res.Big().Uint64(), nil
}

func (c *Client) GetUncleCountByBlockNumber(ctx context.Context, number types.BlockNumber) (uint64, error) {
	var res types.Number
	if err := c.transport.Call(ctx, &res, "eth_getUncleCountByBlockNumber", number); err != nil {
		return 0, err
	}
	return res.Big().Uint64(), nil
}

func (c *Client) GetCode(ctx context.Context, account types.Address, block types.BlockNumber) ([]byte, error) {
	var res types.Bytes
	if err := c.transport.Call(ctx, &res, "eth_getCode", account, block); err != nil {
		return nil, err
	}
	return res.Bytes(), nil
}

func (c *Client) Sign(ctx context.Context, account types.Address, data []byte) (types.Signature, error) {
	var res types.Signature
	if err := c.transport.Call(ctx, &res, "eth_sign", account, types.Bytes(data)); err != nil {
		return types.Signature{}, err
	}
	return res, nil
}

func (c *Client) SignTransaction(ctx context.Context, tx types.Transaction) ([]byte, *types.Transaction, error) {
	var res signTransactionResult
	if err := c.transport.Call(ctx, &res, "eth_signTransaction", tx); err != nil {
		return nil, nil, err
	}
	return res.Raw, res.Tx, nil
}

func (c *Client) SendTransaction(ctx context.Context, tx types.Transaction) (*types.Hash, error) {
	var res types.Hash
	if err := c.transport.Call(ctx, &res, "eth_sendTransaction", tx); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) SendRawTransaction(ctx context.Context, data []byte) (*types.Hash, error) {
	var res types.Hash
	if err := c.transport.Call(ctx, &res, "eth_sendRawTransaction", types.Bytes(data)); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) Call(ctx context.Context, call types.Call, block types.BlockNumber) ([]byte, error) {
	var res types.Bytes
	if err := c.transport.Call(ctx, &res, "eth_call", call, block); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) EstimateGas(ctx context.Context, call types.Call, block types.BlockNumber) (uint64, error) {
	var res types.Number
	if err := c.transport.Call(ctx, &res, "eth_estimateGas", call, block); err != nil {
		return 0, err
	}
	return res.Big().Uint64(), nil
}

func (c *Client) BlockByHash(ctx context.Context, hash types.Hash, full bool) (*types.Block, error) {
	var res types.Block
	if err := c.transport.Call(ctx, &res, "eth_getBlockByHash", hash, full); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) BlockByNumber(ctx context.Context, number types.BlockNumber, full bool) (*types.Block, error) {
	var res types.Block
	if err := c.transport.Call(ctx, &res, "eth_getBlockByNumber", number, full); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetTransactionByHash(ctx context.Context, hash types.Hash) (*types.Transaction, error) {
	var res types.Transaction
	if err := c.transport.Call(ctx, &res, "eth_getTransactionByHash", hash); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetTransactionByBlockHashAndIndex(ctx context.Context, hash types.Hash, index uint64) (*types.Transaction, error) {
	var res types.Transaction
	if err := c.transport.Call(ctx, &res, "eth_getTransactionByBlockHashAndIndex", hash, types.NumberFromUint64(index)); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetTransactionByBlockNumberAndIndex(ctx context.Context, number types.BlockNumber, index uint64) (*types.Transaction, error) {
	var res types.Transaction
	if err := c.transport.Call(ctx, &res, "eth_getTransactionByBlockNumberAndIndex", number, types.NumberFromUint64(index)); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetTransactionReceipt(ctx context.Context, hash types.Hash) (*types.TransactionReceipt, error) {
	var res types.TransactionReceipt
	if err := c.transport.Call(ctx, &res, "eth_getTransactionReceipt", hash); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetLogs(ctx context.Context, query types.FilterLogsQuery) ([]types.Log, error) {
	var res []types.Log
	if err := c.transport.Call(ctx, &res, "eth_getLogs", query); err != nil {
		return nil, err
	}
	return res, nil
}

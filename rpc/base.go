package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
)

// baseClient is a base implementation of the RPC interface. It implements
// RPC methods supported by Ethereum nodes.
type baseClient struct {
	transport transport.Transport
}

// ClientVersion implements the RPC interface.
func (c *baseClient) ClientVersion(ctx context.Context) (string, error) {
	var res string
	if err := c.transport.Call(ctx, &res, "web3_clientVersion"); err != nil {
		return "", err
	}
	return res, nil
}

// Listening implements the RPC interface.
func (c *baseClient) Listening(ctx context.Context) (bool, error) {
	var res bool
	if err := c.transport.Call(ctx, &res, "net_listening"); err != nil {
		return false, err
	}
	return res, nil
}

// PeerCount implements the RPC interface.
func (c *baseClient) PeerCount(ctx context.Context) (uint64, error) {
	var res types.Number
	if err := c.transport.Call(ctx, &res, "net_peerCount"); err != nil {
		return 0, err
	}
	return res.Big().Uint64(), nil
}

// ProtocolVersion implements the RPC interface.
func (c *baseClient) ProtocolVersion(ctx context.Context) (uint64, error) {
	var res types.Number
	if err := c.transport.Call(ctx, &res, "eth_protocolVersion"); err != nil {
		return 0, err
	}
	return res.Big().Uint64(), nil
}

// Syncing implements the RPC interface.
func (c *baseClient) Syncing(ctx context.Context) (*types.SyncStatus, error) {
	var res types.SyncStatus
	if err := c.transport.Call(ctx, &res, "eth_syncing"); err != nil {
		return nil, err
	}
	return &res, nil
}

// NetworkID implements the RPC interface.
func (c *baseClient) NetworkID(ctx context.Context) (uint64, error) {
	var res types.Number
	if err := c.transport.Call(ctx, &res, "net_version"); err != nil {
		return 0, err
	}
	return res.Big().Uint64(), nil
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
	if !res.Big().IsUint64() {
		return 0, errors.New("transaction count is too big")
	}
	return res.Big().Uint64(), nil
}

// GetBlockTransactionCountByHash implements the RPC interface.
func (c *baseClient) GetBlockTransactionCountByHash(ctx context.Context, hash types.Hash) (uint64, error) {
	var res types.Number
	if err := c.transport.Call(ctx, &res, "eth_getBlockTransactionCountByHash", hash); err != nil {
		return 0, err
	}
	if !res.Big().IsUint64() {
		return 0, errors.New("transaction count is too big")
	}
	return res.Big().Uint64(), nil
}

// GetBlockTransactionCountByNumber implements the RPC interface.
func (c *baseClient) GetBlockTransactionCountByNumber(ctx context.Context, number types.BlockNumber) (uint64, error) {
	var res types.Number
	if err := c.transport.Call(ctx, &res, "eth_getBlockTransactionCountByNumber", number); err != nil {
		return 0, err
	}
	if !res.Big().IsUint64() {
		return 0, errors.New("transaction count is too big")
	}
	return res.Big().Uint64(), nil
}

// GetUncleCountByBlockHash implements the RPC interface.
func (c *baseClient) GetUncleCountByBlockHash(ctx context.Context, hash types.Hash) (uint64, error) {
	var res types.Number
	if err := c.transport.Call(ctx, &res, "eth_getUncleCountByBlockHash", hash); err != nil {
		return 0, err
	}
	if !res.Big().IsUint64() {
		return 0, errors.New("uncle count is too big")
	}
	return res.Big().Uint64(), nil
}

// GetUncleCountByBlockNumber implements the RPC interface.
func (c *baseClient) GetUncleCountByBlockNumber(ctx context.Context, number types.BlockNumber) (uint64, error) {
	var res types.Number
	if err := c.transport.Call(ctx, &res, "eth_getUncleCountByBlockNumber", number); err != nil {
		return 0, err
	}
	if !res.Big().IsUint64() {
		return 0, errors.New("uncle count is too big")
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
func (c *baseClient) SignTransaction(ctx context.Context, tx *types.Transaction) ([]byte, *types.Transaction, error) {
	if tx == nil {
		return nil, nil, errors.New("rpc client: transaction is nil")
	}
	var res signTransactionResult
	if err := c.transport.Call(ctx, &res, "eth_signTransaction", tx); err != nil {
		return nil, nil, err
	}
	return res.Raw, res.Tx, nil
}

// SendTransaction implements the RPC interface.
func (c *baseClient) SendTransaction(ctx context.Context, tx *types.Transaction) (*types.Hash, *types.Transaction, error) {
	if tx == nil {
		return nil, nil, errors.New("rpc client: transaction is nil")
	}
	var res types.Hash
	if err := c.transport.Call(ctx, &res, "eth_sendTransaction", tx); err != nil {
		return nil, nil, err
	}
	return &res, tx, nil
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
func (c *baseClient) Call(ctx context.Context, call *types.Call, block types.BlockNumber) ([]byte, *types.Call, error) {
	if call == nil {
		return nil, nil, errors.New("rpc client: call is nil")
	}
	var res types.Bytes
	if err := c.transport.Call(ctx, &res, "eth_call", call, block); err != nil {
		return nil, nil, err
	}
	return res, call, nil
}

// EstimateGas implements the RPC interface.
func (c *baseClient) EstimateGas(ctx context.Context, call *types.Call, block types.BlockNumber) (uint64, *types.Call, error) {
	if call == nil {
		return 0, nil, errors.New("rpc client: call is nil")
	}
	var res types.Number
	if err := c.transport.Call(ctx, &res, "eth_estimateGas", call, block); err != nil {
		return 0, nil, err
	}
	if !res.Big().IsUint64() {
		return 0, nil, errors.New("gas estimate is too big")
	}
	return res.Big().Uint64(), call, nil
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

// GetBlockReceipts implements the RPC interface.
func (c *baseClient) GetBlockReceipts(ctx context.Context, block types.BlockNumber) ([]*types.TransactionReceipt, error) {
	var res []*types.TransactionReceipt
	if err := c.transport.Call(ctx, &res, "eth_getBlockReceipts", block); err != nil {
		return nil, err
	}
	return res, nil
}

// GetUncleByBlockHashAndIndex implements the RPC interface.
func (c *baseClient) GetUncleByBlockHashAndIndex(ctx context.Context, hash types.Hash, index uint64) (*types.Block, error) {
	var res types.Block
	if err := c.transport.Call(ctx, &res, "eth_getUncleByBlockHashAndIndex", hash, types.NumberFromUint64(index)); err != nil {
		return nil, err
	}
	return &res, nil
}

// GetUncleByBlockNumberAndIndex implements the RPC interface.
func (c *baseClient) GetUncleByBlockNumberAndIndex(ctx context.Context, number types.BlockNumber, index uint64) (*types.Block, error) {
	var res types.Block
	if err := c.transport.Call(ctx, &res, "eth_getUncleByBlockNumberAndIndex", number, types.NumberFromUint64(index)); err != nil {
		return nil, err
	}
	return &res, nil
}

// GetLogs implements the RPC interface.
func (c *baseClient) GetLogs(ctx context.Context, query *types.FilterLogsQuery) ([]types.Log, error) {
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

// SubscribeLogs implements the RPC interface.
func (c *baseClient) SubscribeLogs(ctx context.Context, query *types.FilterLogsQuery) (chan types.Log, error) {
	return subscribe[types.Log](ctx, c.transport, "logs", query)
}

// SubscribeNewHeads implements the RPC interface.
func (c *baseClient) SubscribeNewHeads(ctx context.Context) (chan types.Block, error) {
	return subscribe[types.Block](ctx, c.transport, "newHeads")
}

// SubscribeNewPendingTransactions implements the RPC interface.
func (c *baseClient) SubscribeNewPendingTransactions(ctx context.Context) (chan types.Hash, error) {
	return subscribe[types.Hash](ctx, c.transport, "newPendingTransactions")
}

// subscribe creates a subscription to the given method and returns a channel
// that will receive the subscription messages. The messages are unmarshalled
// to the T type. The subscription is unsubscribed and channel closed when the
// context is cancelled.
func subscribe[T any](ctx context.Context, t transport.Transport, method string, params ...any) (chan T, error) {
	st, ok := t.(transport.SubscriptionTransport)
	if !ok {
		return nil, errors.New("transport does not support subscriptions")
	}
	rawCh, subID, err := st.Subscribe(ctx, method, params...)
	if err != nil {
		return nil, err
	}
	msgCh := make(chan T)
	go subscriptionRoutine(ctx, st, subID, rawCh, msgCh)
	return msgCh, nil
}

//nolint:errcheck
func subscriptionRoutine[T any](ctx context.Context, t transport.SubscriptionTransport, subID string, rawCh chan json.RawMessage, msgCh chan T) {
	defer close(msgCh)
	defer t.Unsubscribe(ctx, subID)
	for {
		select {
		case <-ctx.Done():
			return
		case raw, ok := <-rawCh:
			if !ok {
				return
			}
			var msg T
			if err := json.Unmarshal(raw, &msg); err != nil {
				continue
			}
			msgCh <- msg
		}
	}
}

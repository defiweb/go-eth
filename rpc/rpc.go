package rpc

import (
	"context"
	"math/big"

	"github.com/defiweb/go-eth/types"
)

// RPC is an RPC client for the Ethereum-compatible nodes.
type RPC interface {
	// ClientVersion performs web3_clientVersion RPC call.
	//
	// It returns the current client version.
	ClientVersion(ctx context.Context) (string, error)

	// Listening performs net_listening RPC call.
	//
	// It returns true if the client is actively listening for network.
	Listening(ctx context.Context) (bool, error)

	// PeerCount performs net_peerCount RPC call.
	//
	// It returns the number of connected peers.
	PeerCount(ctx context.Context) (uint64, error)

	// ProtocolVersion performs eth_protocolVersion RPC call.
	//
	// It returns the current Ethereum protocol version.
	ProtocolVersion(ctx context.Context) (uint64, error)

	// Syncing performs eth_syncing RPC call.
	//
	// It returns an object with data about the sync status or false.
	Syncing(ctx context.Context) (*types.SyncStatus, error)

	// NetworkID performs net_version RPC call.
	//
	// It returns the current network ID.
	NetworkID(ctx context.Context) (uint64, error)

	// ChainID performs eth_chainId RPC call.
	//
	// It returns the current chain ID.
	ChainID(ctx context.Context) (uint64, error)

	// GasPrice performs eth_gasPrice RPC call.
	//
	// It returns the current price per gas in wei.
	GasPrice(ctx context.Context) (*big.Int, error)

	// Accounts performs eth_accounts RPC call.
	//
	// It returns the list of addresses owned by the client.
	Accounts(ctx context.Context) ([]types.Address, error)

	// BlockNumber performs eth_blockNumber RPC call.
	//
	// It returns the current block number.
	BlockNumber(ctx context.Context) (*big.Int, error)

	// GetBalance performs eth_getBalance RPC call.
	//
	// It returns the balance of the account of given address in wei.
	GetBalance(ctx context.Context, address types.Address, block types.BlockNumber) (*big.Int, error)

	// GetStorageAt performs eth_getStorageAt RPC call.
	//
	// It returns the value of key in the contract storage at the given
	// address.
	GetStorageAt(ctx context.Context, account types.Address, key types.Hash, block types.BlockNumber) (*types.Hash, error)

	// GetTransactionCount performs eth_getTransactionCount RPC call.
	//
	// It returns the number of transactions sent from the given address.
	GetTransactionCount(ctx context.Context, account types.Address, block types.BlockNumber) (uint64, error)

	// GetBlockTransactionCountByHash performs eth_getBlockTransactionCountByHash RPC call.
	//
	// It returns the number of transactions in the block with the given hash.
	GetBlockTransactionCountByHash(ctx context.Context, hash types.Hash) (uint64, error)

	// GetBlockTransactionCountByNumber performs eth_getBlockTransactionCountByNumber RPC call.
	//
	// It returns the number of transactions in the block with the given block
	GetBlockTransactionCountByNumber(ctx context.Context, number types.BlockNumber) (uint64, error)

	// GetUncleCountByBlockHash performs eth_getUncleCountByBlockHash RPC call.
	//
	// It returns the number of uncles in the block with the given hash.
	GetUncleCountByBlockHash(ctx context.Context, hash types.Hash) (uint64, error)

	// GetUncleCountByBlockNumber performs eth_getUncleCountByBlockNumber RPC call.
	//
	// It returns the number of uncles in the block with the given block number.
	GetUncleCountByBlockNumber(ctx context.Context, number types.BlockNumber) (uint64, error)

	// GetCode performs eth_getCode RPC call.
	//
	// It returns the contract code at the given address.
	GetCode(ctx context.Context, account types.Address, block types.BlockNumber) ([]byte, error)

	// Sign performs eth_sign RPC call.
	//
	// It signs the given data with the given address.
	Sign(ctx context.Context, account types.Address, data []byte) (*types.Signature, error)

	// SignTransaction performs eth_signTransaction RPC call.
	//
	// It signs the given transaction.
	//
	// If transaction was internally mutated, the mutated call is returned.
	SignTransaction(ctx context.Context, tx *types.Transaction) ([]byte, *types.Transaction, error)

	// SendTransaction performs eth_sendTransaction RPC call.
	//
	// It sends a transaction to the network.
	//
	// If transaction was internally mutated, the mutated call is returned.
	SendTransaction(ctx context.Context, tx *types.Transaction) (*types.Hash, *types.Transaction, error)

	// SendRawTransaction performs eth_sendRawTransaction RPC call.
	//
	// It sends an encoded transaction to the network.
	SendRawTransaction(ctx context.Context, data []byte) (*types.Hash, error)

	// Call performs eth_call RPC call.
	//
	// It executes a new message call immediately without creating a
	// transaction on the blockchain.
	//
	// If call was internally mutated, the mutated call is returned.
	Call(ctx context.Context, call *types.Call, block types.BlockNumber) ([]byte, *types.Call, error)

	// EstimateGas performs eth_estimateGas RPC call.
	//
	// It estimates the gas necessary to execute a specific transaction.
	//
	// If call was internally mutated, the mutated call is returned.
	EstimateGas(ctx context.Context, call *types.Call, block types.BlockNumber) (uint64, *types.Call, error)

	// BlockByHash performs eth_getBlockByHash RPC call.
	//
	// It returns information about a block by hash.
	BlockByHash(ctx context.Context, hash types.Hash, full bool) (*types.Block, error)

	// BlockByNumber performs eth_getBlockByNumber RPC call.
	//
	// It returns the block with the given number.
	BlockByNumber(ctx context.Context, number types.BlockNumber, full bool) (*types.Block, error)

	// GetTransactionByHash performs eth_getTransactionByHash RPC call.
	//
	// It returns the information about a transaction requested by transaction.
	GetTransactionByHash(ctx context.Context, hash types.Hash) (*types.OnChainTransaction, error)

	// GetTransactionByBlockHashAndIndex performs eth_getTransactionByBlockHashAndIndex RPC call.
	//
	// It returns the information about a transaction requested by transaction.
	GetTransactionByBlockHashAndIndex(ctx context.Context, hash types.Hash, index uint64) (*types.OnChainTransaction, error)

	// GetTransactionByBlockNumberAndIndex performs eth_getTransactionByBlockNumberAndIndex RPC call.
	//
	// It returns the information about a transaction requested by transaction.
	GetTransactionByBlockNumberAndIndex(ctx context.Context, number types.BlockNumber, index uint64) (*types.OnChainTransaction, error)

	// GetTransactionReceipt performs eth_getTransactionReceipt RPC call.
	//
	// It returns the receipt of a transaction by transaction hash.
	GetTransactionReceipt(ctx context.Context, hash types.Hash) (*types.TransactionReceipt, error)

	// GetBlockReceipts performs eth_getBlockReceipts RPC call.
	//
	// It returns all transaction receipts for a given block hash or number.
	GetBlockReceipts(ctx context.Context, block types.BlockNumber) ([]*types.TransactionReceipt, error)

	// GetUncleByBlockHashAndIndex performs eth_getUncleByBlockNumberAndIndex RPC call.
	//
	// It returns information about an uncle of a block by number and uncle index position.
	GetUncleByBlockHashAndIndex(ctx context.Context, hash types.Hash, index uint64) (*types.Block, error)

	// GetUncleByBlockNumberAndIndex performs eth_getUncleByBlockNumberAndIndex RPC call.
	//
	// It returns information about an uncle of a block by hash and uncle index position.
	GetUncleByBlockNumberAndIndex(ctx context.Context, number types.BlockNumber, index uint64) (*types.Block, error)

	// NewFilter performs eth_newFilter RPC call.
	//
	// It creates a filter object based on the given filter options. To check
	// if the state has changed, use GetFilterChanges.
	NewFilter(ctx context.Context, query *types.FilterLogsQuery) (*big.Int, error)

	// NewBlockFilter performs eth_newBlockFilter RPC call.
	//
	// It creates a filter in the node, to notify when a new block arrives. To
	// check if the state has changed, use GetBlockFilterChanges.
	NewBlockFilter(ctx context.Context) (*big.Int, error)

	// NewPendingTransactionFilter performs eth_newPendingTransactionFilter RPC call.
	//
	// It creates a filter in the node, to notify when new pending transactions
	// arrive. To check if the state has changed, use GetFilterChanges.
	NewPendingTransactionFilter(ctx context.Context) (*big.Int, error)

	// UninstallFilter performs eth_uninstallFilter RPC call.
	//
	// It uninstalls a filter with given ID. Should always be called when watch
	// is no longer needed.
	UninstallFilter(ctx context.Context, id *big.Int) (bool, error)

	// GetFilterChanges performs eth_getFilterChanges RPC call.
	//
	// It returns an array of logs that occurred since the given filter ID.
	GetFilterChanges(ctx context.Context, id *big.Int) ([]types.Log, error)

	// GetBlockFilterChanges performs eth_getFilterChanges RPC call.
	//
	// It returns an array of block hashes that occurred since the given filter ID.
	GetBlockFilterChanges(ctx context.Context, id *big.Int) ([]types.Hash, error)

	// GetFilterLogs performs eth_getFilterLogs RPC call.
	//
	// It returns an array of all logs matching filter with given ID.
	GetFilterLogs(ctx context.Context, id *big.Int) ([]types.Log, error)

	// GetLogs performs eth_getLogs RPC call.
	//
	// It returns logs that match the given query.
	GetLogs(ctx context.Context, query *types.FilterLogsQuery) ([]types.Log, error)

	// MaxPriorityFeePerGas performs eth_maxPriorityFeePerGas RPC call.
	//
	// It returns the estimated maximum priority fee per gas.
	MaxPriorityFeePerGas(ctx context.Context) (*big.Int, error)

	// SubscribeLogs performs eth_subscribe RPC call with "logs" subscription
	// type.
	//
	// It creates a subscription that will send logs that match the given query.
	//
	// Subscription channel will be closed when the context is canceled.
	SubscribeLogs(ctx context.Context, query *types.FilterLogsQuery) (chan types.Log, error)

	// SubscribeNewHeads performs eth_subscribe RPC call with "newHeads"
	// subscription type.
	//
	// It creates a subscription that will send new block headers.
	//
	// Subscription channel will be closed when the context is canceled.
	SubscribeNewHeads(ctx context.Context) (chan types.Block, error)

	// SubscribeNewPendingTransactions performs eth_subscribe RPC call with
	// "newPendingTransactions" subscription type.
	//
	// It creates a subscription that will send new pending transactions.
	//
	// Subscription channel will be closed when the context is canceled.
	SubscribeNewPendingTransactions(ctx context.Context) (chan types.Hash, error)
}

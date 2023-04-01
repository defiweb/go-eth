package rpc

import (
	"context"
	"math/big"

	"github.com/defiweb/go-eth/types"
)

// RPC is an RPC client for the Ethereum-compatible nodes.
type RPC interface {
	// TODO: web3_clientVersion
	// TODO: web3_sha3
	// TODO: net_version
	// TODO: net_listening
	// TODO: net_peerCount
	// TODO: eth_protocolVersion
	// TODO: eth_syncing
	// TODO: eth_coinbase
	// TODO: eth_mining
	// TODO: eth_hashrate

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
	SignTransaction(ctx context.Context, tx types.Transaction) ([]byte, *types.Transaction, error)

	// SendTransaction performs eth_sendTransaction RPC call.
	//
	// It creates new message call transaction or a contract creation for
	// signed transactions.
	SendTransaction(ctx context.Context, tx types.Transaction) (*types.Hash, error)

	// SendRawTransaction performs eth_sendRawTransaction RPC call.
	//
	// It sends an encoded transaction to the network.
	SendRawTransaction(ctx context.Context, data []byte) (*types.Hash, error)

	// Call performs eth_call RPC call.
	//
	// Creates new message call transaction or a contract creation, if the data
	// field contains code.
	Call(ctx context.Context, call types.Call, block types.BlockNumber) ([]byte, error)

	// EstimateGas performs eth_estimateGas RPC call.
	//
	// It estimates the gas necessary to execute a specific transaction.
	EstimateGas(ctx context.Context, call types.Call, block types.BlockNumber) (uint64, error)

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

	// TODO: eth_getUncleByBlockHashAndIndex
	// TODO: eth_getUncleByBlockNumberAndIndex
	// TODO: eth_getCompilers
	// TODO: eth_compileSolidity
	// TODO: eth_compileLLL
	// TODO: eth_compileSerpent
	// TODO: eth_newFilter
	// TODO: eth_newBlockFilter
	// TODO: eth_newPendingTransactionFilter
	// TODO: eth_uninstallFilter
	// TODO: eth_getFilterChanges
	// TODO: eth_getFilterLogs

	// GetLogs performs eth_getLogs RPC call.
	//
	// It returns logs that match the given query.
	GetLogs(ctx context.Context, query types.FilterLogsQuery) ([]types.Log, error)

	// TODO: eth_getWork
	// TODO: eth_submitWork
	// TODO: eth_submitHashrate

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
	SubscribeLogs(ctx context.Context, query types.FilterLogsQuery) (chan types.Log, error)

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

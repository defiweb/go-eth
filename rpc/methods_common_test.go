package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"math/big"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/defiweb/go-eth/hexutil"
	"github.com/defiweb/go-eth/types"
)

const mockBlockResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": {
	    "number": "0x11",
	    "hash": "0x2222222222222222222222222222222222222222222222222222222222222222",
	    "parentHash": "0x3333333333333333333333333333333333333333333333333333333333333333",
	    "nonce": "0x4444444444444444",
	    "sha3Uncles": "0x5555555555555555555555555555555555555555555555555555555555555555",
	    "logsBloom": "0x66666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666",
	    "transactionsRoot": "0x7777777777777777777777777777777777777777777777777777777777777777",
	    "stateRoot": "0x8888888888888888888888888888888888888888888888888888888888888888",
	    "receiptsRoot": "0x9999999999999999999999999999999999999999999999999999999999999999",
	    "miner": "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	    "difficulty": "0xbbbbbb",
	    "totalDifficulty": "0xcccccc",
	    "extraData": "0x0000000000000000000000000000000000000000000000000000000000000000",
	    "size": "0xdddddd",
	    "gasLimit": "0xeeeeee",
	    "gasUsed": "0xffffff",
	    "timestamp": "0x54e34e8e",
	    "transactions": [
		  {
	        "hash": "0x1111111111111111111111111111111111111111111111111111111111111111",
	        "nonce": "0x22",
	        "blockHash": "0x3333333333333333333333333333333333333333333333333333333333333333",
	        "blockNumber": "0x4444",
	        "transactionIndex": "0x01",
	        "from": "0x5555555555555555555555555555555555555555",
	        "to": "0x6666666666666666666666666666666666666666",
	        "value": "0x2540be400",
	        "gas": "0x76c0",
	        "gasPrice": "0x9184e72a000",
	        "input": "0x777777777777"
		  }
		],
	    "uncles": [
	        "0x8888888888888888888888888888888888888888888888888888888888888888"
		]
	  }
	}
`

const mockOnChainTransactionResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": {
	    "blockHash": "0x1111111111111111111111111111111111111111111111111111111111111111",
	    "blockNumber": "0x22",
	    "from": "0x3333333333333333333333333333333333333333",
	    "gas": "0x76c0",
	    "gasPrice": "0x9184e72a000",
	    "hash": "0x4444444444444444444444444444444444444444444444444444444444444444",
	    "input": "0x555555555555",
	    "nonce": "0x66",
	    "to": "0x7777777777777777777777777777777777777777",
	    "transactionIndex": "0x0",
	    "value": "0x2540be400",
	    "v": "0x88",
	    "r": "0x9999999999999999999999999999999999999999999999999999999999999999",
	    "s": "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	  }
	}
`

const mockGetBalanceRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_getBalance",
	  "params": [
	    "0x1111111111111111111111111111111111111111",
	    "latest"
	  ]
	}
`

const mockGetBalanceResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": "0x0234c8a3397aab58"
	}
`

func TestBaseClient_GetBalance(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsCommon{Transport: httpMock}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockGetBalanceRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockGetBalanceResponse)),
		}, nil
	}

	balance, err := client.GetBalance(
		t.Context(),
		types.MustAddressFromHex("0x1111111111111111111111111111111111111111"),
		types.LatestBlockNumber,
	)

	require.NoError(t, err)
	assert.Equal(t, big.NewInt(158972490234375000), balance)
}

const mockGetCodeRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_getCode",
	  "params": [
	    "0x1111111111111111111111111111111111111111",
	    "0x2"
	  ]
	}
`

const mockGetCodeResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": "0x3333333333333333333333333333333333333333333333333333333333333333"
	}
`

func TestBaseClient_GetCode(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsCommon{Transport: httpMock}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockGetCodeRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockGetCodeResponse)),
		}, nil
	}

	code, err := client.GetCode(
		t.Context(),
		types.MustAddressFromHex("0x1111111111111111111111111111111111111111"),
		types.MustBlockNumberFromHex("0x2"),
	)

	require.NoError(t, err)
	assert.Equal(t, "0x3333333333333333333333333333333333333333333333333333333333333333", hexutil.BytesToHex(code))
}

const mockGetStorageAtRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_getStorageAt",
	  "params": [
	    "0x1111111111111111111111111111111111111111",
	    "0x2222222222222222222222222222222222222222222222222222222222222222",
	    "0x1"
	  ]
	}
`

const mockGetStorageAtResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": "0x3333333333333333333333333333333333333333333333333333333333333333"
	}
`

func TestBaseClient_GetStorageAt(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsCommon{Transport: httpMock}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockGetStorageAtRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockGetStorageAtResponse)),
		}, nil
	}

	storage, err := client.GetStorageAt(
		t.Context(),
		types.MustAddressFromHex("0x1111111111111111111111111111111111111111"),
		types.MustHashFromHex("0x2222222222222222222222222222222222222222222222222222222222222222", types.PadNone),
		types.MustBlockNumberFromHex("0x1"),
	)

	require.NoError(t, err)
	assert.Equal(t, types.MustHashFromHex("0x3333333333333333333333333333333333333333333333333333333333333333", types.PadNone), *storage)
}

const mockGetTransactionCountRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_getTransactionCount",
	  "params": [
	    "0x1111111111111111111111111111111111111111",
	    "0x1"
	  ]
	}
`

const mockGetTransactionCountResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": "0x1"
	}
`

func TestBaseClient_GetTransactionCount(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsCommon{Transport: httpMock}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockGetTransactionCountRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockGetTransactionCountResponse)),
		}, nil
	}

	transactionCount, err := client.GetTransactionCount(
		t.Context(),
		types.MustAddressFromHex("0x1111111111111111111111111111111111111111"),
		types.MustBlockNumberFromHex("0x1"),
	)

	require.NoError(t, err)
	assert.Equal(t, uint64(1), transactionCount)
}

const mockBlockByHashRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_getBlockByHash",
	  "params": [
	    "0x1111111111111111111111111111111111111111111111111111111111111111",
		true
	  ]
	}
`

func TestBaseClient_BlockByHash(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsCommon{Transport: httpMock}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockBlockByHashRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockBlockResponse)),
		}, nil
	}

	block, err := client.BlockByHash(
		t.Context(),
		types.MustHashFromHex("0x1111111111111111111111111111111111111111111111111111111111111111", types.PadNone),
		true,
	)

	require.NoError(t, err)
	assert.Equal(t, types.MustHashFromHex("0x2222222222222222222222222222222222222222222222222222222222222222", types.PadNone), block.Hash)
}

const mockBlockByNumberRequest = `
	{
	  "method": "eth_getBlockByNumber",
	  "params": [
	    "0x1",
		true
	  ],
	  "id": 1,
	  "jsonrpc": "2.0"
	}
`

func TestBaseClient_BlockByNumber(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsCommon{Transport: httpMock}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockBlockByNumberRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockBlockResponse)),
		}, nil
	}

	block, err := client.BlockByNumber(
		t.Context(),
		types.MustBlockNumberFromHex("0x1"),
		true,
	)

	require.NoError(t, err)
	assert.Equal(t, big.NewInt(0x11), block.Number)
	assert.Equal(t, types.MustHashFromHex("0x2222222222222222222222222222222222222222222222222222222222222222", types.PadNone), block.Hash)
	assert.Equal(t, types.MustHashFromHex("0x3333333333333333333333333333333333333333333333333333333333333333", types.PadNone), block.ParentHash)
	assert.Equal(t, hexutil.MustHexToBigInt("0x4444444444444444"), block.Nonce)
	assert.Equal(t, types.MustHashFromHex("0x5555555555555555555555555555555555555555555555555555555555555555", types.PadNone), block.Sha3Uncles)
	assert.Equal(t, hexutil.MustHexToBytes("0x66666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666"), block.LogsBloom)
	assert.Equal(t, types.MustHashFromHex("0x7777777777777777777777777777777777777777777777777777777777777777", types.PadNone), block.TransactionsRoot)
	assert.Equal(t, types.MustHashFromHex("0x8888888888888888888888888888888888888888888888888888888888888888", types.PadNone), block.StateRoot)
	assert.Equal(t, types.MustHashFromHex("0x9999999999999999999999999999999999999999999999999999999999999999", types.PadNone), block.ReceiptsRoot)
	assert.Equal(t, types.MustAddressFromHex("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"), block.Miner)
	assert.Equal(t, hexutil.MustHexToBigInt("0xbbbbbb"), block.Difficulty)
	assert.Equal(t, hexutil.MustHexToBigInt("0xcccccc"), block.TotalDifficulty)
	assert.Equal(t, hexutil.MustHexToBytes("0x0000000000000000000000000000000000000000000000000000000000000000"), block.ExtraData)
	assert.Equal(t, hexutil.MustHexToBigInt("0xdddddd").Uint64(), block.Size)
	assert.Equal(t, hexutil.MustHexToBigInt("0xeeeeee").Uint64(), block.GasLimit)
	assert.Equal(t, hexutil.MustHexToBigInt("0xffffff").Uint64(), block.GasUsed)
	assert.Equal(t, int64(1424182926), block.Timestamp.Unix())
	require.Len(t, block.Transactions, 1)
	require.Len(t, block.Uncles, 1)

	tx := block.Transactions[0].Transaction.(*types.TransactionLegacy)
	assert.Equal(t, types.MustHashFromHexPtr("0x1111111111111111111111111111111111111111111111111111111111111111", types.PadNone), block.Transactions[0].Hash)
	assert.Equal(t, uint64(0x22), *tx.Nonce)
	assert.Equal(t, types.MustAddressFromHexPtr("0x5555555555555555555555555555555555555555"), tx.From)
	assert.Equal(t, types.MustAddressFromHexPtr("0x6666666666666666666666666666666666666666"), tx.To)
	assert.Equal(t, big.NewInt(10000000000), tx.Value)
	assert.Equal(t, uint64(30400), *tx.GasLimit)
	assert.Equal(t, big.NewInt(10000000000000), tx.GasPrice)
	assert.Equal(t, hexutil.MustHexToBytes("0x777777777777"), tx.Input)
	assert.Equal(t, types.MustHashFromHex("0x8888888888888888888888888888888888888888888888888888888888888888", types.PadNone), block.Uncles[0])
}

const mockGetBlockTransactionCountByHashRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_getBlockTransactionCountByHash",
	  "params": [
	    "0x1111111111111111111111111111111111111111111111111111111111111111"
	  ]
	}
`

const mockGetBlockTransactionCountByHashResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": "0x1"
	}
`

func TestBaseClient_GetBlockTransactionCountByHash(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsCommon{Transport: httpMock}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockGetBlockTransactionCountByHashRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockGetBlockTransactionCountByHashResponse)),
		}, nil
	}

	transactionCount, err := client.GetBlockTransactionCountByHash(
		t.Context(),
		types.MustHashFromHex("0x1111111111111111111111111111111111111111111111111111111111111111", types.PadNone),
	)

	require.NoError(t, err)
	assert.Equal(t, uint64(1), transactionCount)
}

const mockGetBlockTransactionCountByNumberRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_getBlockTransactionCountByNumber",
	  "params": [
	    "0x1"
	  ]
	}
`

const mockGetBlockTransactionCountByNumberResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": "0x2"
	}
`

func TestBaseClient_GetBlockTransactionCountByNumber(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsCommon{Transport: httpMock}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockGetBlockTransactionCountByNumberRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockGetBlockTransactionCountByNumberResponse)),
		}, nil
	}

	transactionCount, err := client.GetBlockTransactionCountByNumber(
		t.Context(),
		types.MustBlockNumberFromHex("0x1"),
	)

	require.NoError(t, err)
	assert.Equal(t, uint64(2), transactionCount)
}

const mockGetUncleByBlockHashAndIndexRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_getUncleByBlockHashAndIndex",
	  "params": [
	    "0x1111111111111111111111111111111111111111111111111111111111111111",
	    "0x0"
	  ]
	}
`

func TestBaseClient_GetUncleByBlockHashAndIndex(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsCommon{Transport: httpMock}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockGetUncleByBlockHashAndIndexRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockBlockResponse)),
		}, nil
	}

	block, err := client.GetUncleByBlockHashAndIndex(
		t.Context(),
		types.MustHashFromHex("0x1111111111111111111111111111111111111111111111111111111111111111", types.PadNone),
		0,
	)

	require.NoError(t, err)
	assert.Equal(t, types.MustHashFromHex("0x2222222222222222222222222222222222222222222222222222222222222222", types.PadNone), block.Hash)
}

const mockGetUncleByBlockNumberAndIndexRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_getUncleByBlockNumberAndIndex",
	  "params": [
	    "0x1",
	    "0x2"
	  ]
	}
`

func TestBaseClient_GetUncleByBlockNumberAndIndex(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsCommon{Transport: httpMock}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockGetUncleByBlockNumberAndIndexRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockBlockResponse)),
		}, nil
	}

	block, err := client.GetUncleByBlockNumberAndIndex(
		t.Context(),
		types.MustBlockNumberFromHex("0x1"),
		2,
	)

	require.NoError(t, err)
	assert.Equal(t, types.MustHashFromHex("0x2222222222222222222222222222222222222222222222222222222222222222", types.PadNone), block.Hash)
}

const mockGetUncleCountByBlockHashRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_getUncleCountByBlockHash",
	  "params": [
	    "0x1111111111111111111111111111111111111111111111111111111111111111"
	  ]
	}
`

const mockGetUncleCountByBlockHashResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": "0x1"
	}
`

func TestBaseClient_GetUncleCountByBlockHash(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsCommon{Transport: httpMock}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockGetUncleCountByBlockHashRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockGetUncleCountByBlockHashResponse)),
		}, nil
	}

	uncleCount, err := client.GetUncleCountByBlockHash(
		t.Context(),
		types.MustHashFromHex("0x1111111111111111111111111111111111111111111111111111111111111111", types.PadNone),
	)

	require.NoError(t, err)
	assert.Equal(t, uint64(1), uncleCount)
}

const mockGetUncleCountByBlockNumberRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_getUncleCountByBlockNumber",
	  "params": [
	    "0x1"
	  ]
	}
`

const mockGetUncleCountByBlockNumberResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": "0x2"
	}
`

func TestBaseClient_GetUncleCountByBlockNumber(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsCommon{Transport: httpMock}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockGetUncleCountByBlockNumberRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockGetUncleCountByBlockNumberResponse)),
		}, nil
	}

	uncleCount, err := client.GetUncleCountByBlockNumber(
		t.Context(),
		types.MustBlockNumberFromHex("0x1"),
	)

	require.NoError(t, err)
	assert.Equal(t, uint64(2), uncleCount)
}

// TODO: SubscribeNewHeads

const mockCallRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_call",
	  "params": [
		{
		  "from": "0x1111111111111111111111111111111111111111",
		  "to": "0x2222222222222222222222222222222222222222",
		  "gas": "0x76c0",
		  "gasPrice": "0x9184e72a000",
		  "value": "0x2540be400",
		  "input": "0x3333333333333333333333333333333333333333333333333333333333333333333333333333333333"
		},
	    "0x1"
	  ]
	}
`

const mockCallResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": "0x00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000004000000000000000000000000d9c9cd5f6779558b6e0ed4e6acf6b1947e7fa1f300000000000000000000000078d1ad571a1a09d60d9bbf25894b44e4c8859595000000000000000000000000286834935f4a8cfb4ff4c77d5770c2775ae2b0e7000000000000000000000000b86e2b0ab5a4b1373e40c51a7c712c70ba2f9f8e"
	}
`

func TestBaseClient_Call(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsCommon{Transport: httpMock}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockCallRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockCallResponse)),
		}, nil
	}

	from := types.MustAddressFromHexPtr("0x1111111111111111111111111111111111111111")
	to := types.MustAddressFromHexPtr("0x2222222222222222222222222222222222222222")
	gasLimit := uint64(30400)
	gasPrice := big.NewInt(10000000000000)
	value := big.NewInt(10000000000)
	input := hexutil.MustHexToBytes("0x3333333333333333333333333333333333333333333333333333333333333333333333333333333333")
	response, err := client.Call(
		t.Context(),
		&types.CallLegacy{
			CallData: types.CallData{
				From:     from,
				To:       to,
				GasLimit: &gasLimit,
				Value:    value,
				Input:    input,
			},
			LegacyPriceData: types.LegacyPriceData{
				GasPrice: gasPrice,
			},
		},
		types.MustBlockNumberFromHex("0x1"),
	)

	require.NoError(t, err)
	assert.Equal(t, hexutil.MustHexToBytes("0x00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000004000000000000000000000000d9c9cd5f6779558b6e0ed4e6acf6b1947e7fa1f300000000000000000000000078d1ad571a1a09d60d9bbf25894b44e4c8859595000000000000000000000000286834935f4a8cfb4ff4c77d5770c2775ae2b0e7000000000000000000000000b86e2b0ab5a4b1373e40c51a7c712c70ba2f9f8e"), response)
}

const mockEstimateGasRequest = `
	{
	  "id": 1,
	  "jsonrpc": "2.0",
	  "method": "eth_estimateGas",
	  "params": [
	    {
	      "from": "0x1111111111111111111111111111111111111111",
	      "to": "0x2222222222222222222222222222222222222222",
	      "gas": "0x76c0",
	      "gasPrice": "0x9184e72a000",
	      "value": "0x2540be400",
	      "input": "0x3333333333333333333333333333333333333333333333333333333333333333333333333333333333"
	    },
	    "latest"
	  ]
	}
`

const mockEstimateGasResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": "0x5208"
	}
`

func TestBaseClient_EstimateGas(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsCommon{Transport: httpMock}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockEstimateGasRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockEstimateGasResponse)),
		}, nil
	}

	gasLimit := uint64(30400)
	gas, err := client.EstimateGas(
		t.Context(),
		&types.CallLegacy{
			CallData: types.CallData{
				From:     types.MustAddressFromHexPtr("0x1111111111111111111111111111111111111111"),
				To:       types.MustAddressFromHexPtr("0x2222222222222222222222222222222222222222"),
				GasLimit: &gasLimit,
				Value:    big.NewInt(10000000000),
				Input:    hexutil.MustHexToBytes("0x3333333333333333333333333333333333333333333333333333333333333333333333333333333333"),
			},
			LegacyPriceData: types.LegacyPriceData{
				GasPrice: big.NewInt(10000000000000),
			},
		},
		types.LatestBlockNumber,
	)

	require.NoError(t, err)
	assert.Equal(t, uint64(21000), gas)
}

const mockSendRawTransactionRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_sendRawTransaction",
	  "params": [
	    "0xf893808609184e72a0008276c094d46e8dd67c5d32be8058bb8eb970870f072445678502540be400a9d46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f07244567511a02222222222222222222222222222222222222222222222222222222222222222a03333333333333333333333333333333333333333333333333333333333333333"
	  ]
	}
`

const mockSendRawTransactionResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": "0x1111111111111111111111111111111111111111111111111111111111111111"
	}
`

func TestBaseClient_SendRawTransaction(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsCommon{Transport: httpMock}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockSendRawTransactionRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockSendRawTransactionResponse)),
		}, nil
	}

	txHash, err := client.SendRawTransaction(
		t.Context(),
		hexutil.MustHexToBytes("0xf893808609184e72a0008276c094d46e8dd67c5d32be8058bb8eb970870f072445678502540be400a9d46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f07244567511a02222222222222222222222222222222222222222222222222222222222222222a03333333333333333333333333333333333333333333333333333333333333333"),
	)

	require.NoError(t, err)
	assert.Equal(t, types.MustHashFromHex("0x1111111111111111111111111111111111111111111111111111111111111111", types.PadNone), *txHash)
}

const mockGetTransactionByHashRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_getTransactionByHash",
	  "params": [
	    "0x1111111111111111111111111111111111111111111111111111111111111111"
	  ]
	}
`

func TestBaseClient_GetTransactionByHash(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsCommon{Transport: httpMock}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockGetTransactionByHashRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockOnChainTransactionResponse)),
		}, nil
	}

	onChainTX, err := client.GetTransactionByHash(
		t.Context(),
		types.MustHashFromHex("0x1111111111111111111111111111111111111111111111111111111111111111", types.PadNone),
	)

	require.NoError(t, err)
	assert.Equal(t, types.MustHashFromHexPtr("0x1111111111111111111111111111111111111111111111111111111111111111", types.PadNone), onChainTX.BlockHash)
	assert.Equal(t, big.NewInt(0x22), onChainTX.BlockNumber)
	assert.Equal(t, uint64(0x0), *onChainTX.TransactionIndex)
	assert.Equal(t, types.MustHashFromHexPtr("0x4444444444444444444444444444444444444444444444444444444444444444", types.PadNone), onChainTX.Hash)

	tx := onChainTX.Transaction.(*types.TransactionLegacy)
	assert.Equal(t, types.MustAddressFromHexPtr("0x3333333333333333333333333333333333333333"), tx.From)
	assert.Equal(t, types.MustAddressFromHexPtr("0x7777777777777777777777777777777777777777"), tx.To)
	assert.Equal(t, big.NewInt(10000000000), tx.Value)
	assert.Equal(t, uint64(30400), *tx.GasLimit)
	assert.Equal(t, big.NewInt(10000000000000), tx.GasPrice)
	assert.Equal(t, hexutil.MustHexToBytes("0x555555555555"), tx.Input)
	assert.Equal(t, uint64(0x66), *tx.Nonce)
	assert.Equal(t, uint64(0x88), tx.Signature.V.Uint64())
	assert.Equal(t, hexutil.MustHexToBytes("0x9999999999999999999999999999999999999999999999999999999999999999"), tx.Signature.R.Bytes())
	assert.Equal(t, hexutil.MustHexToBytes("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"), tx.Signature.S.Bytes())
}

const mockGetTransactionByBlockHashAndIndexRequest = `
	{
	  "id": 1,
	  "jsonrpc": "2.0",
	  "method": "eth_getTransactionByBlockHashAndIndex",
	  "params": [
	    "0x1111111111111111111111111111111111111111111111111111111111111111",
	    "0x0"
	  ]
	}
`

func TestBaseClient_GetTransactionByBlockHashAndIndex(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsCommon{Transport: httpMock}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockGetTransactionByBlockHashAndIndexRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockOnChainTransactionResponse)),
		}, nil
	}

	onChainTX, err := client.GetTransactionByBlockHashAndIndex(
		t.Context(),
		types.MustHashFromHex("0x1111111111111111111111111111111111111111111111111111111111111111", types.PadNone),
		0,
	)

	require.NoError(t, err)
	assert.Equal(t, types.MustHashFromHexPtr("0x1111111111111111111111111111111111111111111111111111111111111111", types.PadNone), onChainTX.BlockHash)
	assert.Equal(t, big.NewInt(0x22), onChainTX.BlockNumber)
	assert.Equal(t, uint64(0x0), *onChainTX.TransactionIndex)
	assert.Equal(t, types.MustHashFromHexPtr("0x4444444444444444444444444444444444444444444444444444444444444444", types.PadNone), onChainTX.Hash)
}

const mockGetTransactionByBlockNumberAndIndexRequest = `
	{
	  "id": 1,
	  "jsonrpc": "2.0",
	  "method": "eth_getTransactionByBlockNumberAndIndex",
	  "params": [
	    "0x1",
	    "0x2"
	  ]
	}
`

func TestBaseClient_GetTransactionByBlockNumberAndIndex(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsCommon{Transport: httpMock}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockGetTransactionByBlockNumberAndIndexRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockOnChainTransactionResponse)),
		}, nil
	}

	onChainTX, err := client.GetTransactionByBlockNumberAndIndex(
		t.Context(),
		types.MustBlockNumberFromHex("0x1"),
		2,
	)

	require.NoError(t, err)
	assert.Equal(t, types.MustHashFromHexPtr("0x1111111111111111111111111111111111111111111111111111111111111111", types.PadNone), onChainTX.BlockHash)
	assert.Equal(t, big.NewInt(0x22), onChainTX.BlockNumber)
	assert.Equal(t, uint64(0x0), *onChainTX.TransactionIndex)
	assert.Equal(t, types.MustHashFromHexPtr("0x4444444444444444444444444444444444444444444444444444444444444444", types.PadNone), onChainTX.Hash)
}

const mockGetTransactionReceiptRequest = `
	{
	  "id": 1,
	  "jsonrpc": "2.0",
	  "method": "eth_getTransactionReceipt",
	  "params": [
	    "0x1111111111111111111111111111111111111111111111111111111111111111"
	  ]
	}
`

const mockGetTransactionReceiptResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": {
	    "blockHash": "0x1111111111111111111111111111111111111111111111111111111111111111",
	    "blockNumber": "0x2222",
	    "contractAddress": null,
	    "cumulativeGasUsed": "0x33333",
	    "effectiveGasPrice":"0x4444444444",
	    "from": "0x5555555555555555555555555555555555555555",
	    "gasUsed": "0x66666",
	    "logs": [
		  {
	        "address": "0x7777777777777777777777777777777777777777",
	        "blockHash": "0x1111111111111111111111111111111111111111111111111111111111111111",
	        "blockNumber": "0x2222",
	        "data": "0x000000000000000000000000398137383b3d25c92898c656696e41950e47316b00000000000000000000000000000000000000000000000000000000000cee6100000000000000000000000000000000000000000000000000000000000ac3e100000000000000000000000000000000000000000000000000000000005baf35",
	        "logIndex": "0x8",
	        "removed": false,
	        "topics": [
			  "0x9999999999999999999999999999999999999999999999999999999999999999"
			],
	        "transactionHash": "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	        "transactionIndex": "0x11"
		  }
		],
	    "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000800000200000000000000000000000000000",
	    "status": "0x1",
	    "to": "0x7777777777777777777777777777777777777777",
	    "transactionHash": "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	    "transactionIndex": "0x11",
	    "type": "0x0"
	  }
	}
`

func TestBaseClient_GetTransactionReceipt(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsCommon{Transport: httpMock}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockGetTransactionReceiptRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockGetTransactionReceiptResponse)),
		}, nil
	}

	receipt, err := client.GetTransactionReceipt(
		t.Context(),
		types.MustHashFromHex("0x1111111111111111111111111111111111111111111111111111111111111111", types.PadNone),
	)

	status := uint64(1)
	require.NoError(t, err)
	assert.Equal(t, types.MustHashFromHex("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", types.PadNone), receipt.TransactionHash)
	assert.Equal(t, uint64(17), receipt.TransactionIndex)
	assert.Equal(t, types.MustHashFromHex("0x1111111111111111111111111111111111111111111111111111111111111111", types.PadNone), receipt.BlockHash)
	assert.Equal(t, big.NewInt(0x2222), receipt.BlockNumber)
	assert.Equal(t, (*types.Address)(nil), receipt.ContractAddress)
	assert.Equal(t, hexutil.MustHexToBigInt("0x33333").Uint64(), receipt.CumulativeGasUsed)
	assert.Equal(t, hexutil.MustHexToBigInt("0x4444444444"), receipt.EffectiveGasPrice)
	assert.Equal(t, hexutil.MustHexToBigInt("0x66666").Uint64(), receipt.GasUsed)
	assert.Equal(t, types.MustAddressFromHex("0x5555555555555555555555555555555555555555"), receipt.From)
	assert.Equal(t, types.MustAddressFromHex("0x7777777777777777777777777777777777777777"), receipt.To)
	assert.Equal(t, hexutil.MustHexToBytes("0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000800000200000000000000000000000000000"), receipt.LogsBloom)
	assert.Equal(t, &status, receipt.Status)
	require.Len(t, receipt.Logs, 1)
	assert.Equal(t, types.MustHashFromHexPtr("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", types.PadNone), receipt.Logs[0].TransactionHash)
	assert.Equal(t, uint64(17), *receipt.Logs[0].TransactionIndex)
	assert.Equal(t, types.MustHashFromHexPtr("0x1111111111111111111111111111111111111111111111111111111111111111", types.PadNone), receipt.Logs[0].BlockHash)
	assert.Equal(t, big.NewInt(0x2222), receipt.Logs[0].BlockNumber)
	assert.Equal(t, uint64(8), *receipt.Logs[0].LogIndex)
	assert.Equal(t, hexutil.MustHexToBytes("0x000000000000000000000000398137383b3d25c92898c656696e41950e47316b00000000000000000000000000000000000000000000000000000000000cee6100000000000000000000000000000000000000000000000000000000000ac3e100000000000000000000000000000000000000000000000000000000005baf35"), receipt.Logs[0].Data)
	assert.Equal(t, types.MustAddressFromHex("0x7777777777777777777777777777777777777777"), receipt.Logs[0].Address)
	assert.Equal(t, []types.Hash{types.MustHashFromHex("0x9999999999999999999999999999999999999999999999999999999999999999", types.PadNone)}, receipt.Logs[0].Topics)
	assert.Equal(t, false, receipt.Logs[0].Removed)
}

const mockGetBlockReceiptsRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_getBlockReceipts",
	  "params": [
	    "0x1"
	  ]
	}
`

const mockGetBlockReceiptsResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": [
		{
		  "blockHash": "0x1111111111111111111111111111111111111111111111111111111111111111",
		  "blockNumber": "0x2222",
		  "contractAddress": null,
		  "cumulativeGasUsed": "0x33333",
		  "effectiveGasPrice": "0x4444444444",
		  "from": "0x5555555555555555555555555555555555555555",
		  "gasUsed": "0x66666",
		  "logs": [
			{
			  "address": "0x7777777777777777777777777777777777777777",
			  "blockHash": "0x1111111111111111111111111111111111111111111111111111111111111111",
			  "blockNumber": "0x2222",
			  "data": "0x000000000000000000000000398137383b3d25c92898c656696e41950e47316b00000000000000000000000000000000000000000000000000000000000cee6100000000000000000000000000000000000000000000000000000000000ac3e100000000000000000000000000000000000000000000000000000000005baf35",
			  "logIndex": "0x8",
			  "removed": false,
			  "topics": [
	            "0x9999999999999999999999999999999999999999999999999999999999999999"
			  ],
			  "transactionHash": "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			  "transactionIndex": "0x11"
			}
		  ],
		  "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000800000200000000000000000000000000000",
		  "status": "0x1",
		  "to": "0x7777777777777777777777777777777777777777",
		  "transactionHash": "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		  "transactionIndex": "0x11",
		  "type": "0x0"
		}
	  ]
	}
`

func TestBaseClient_GetBlockReceipts(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsCommon{Transport: httpMock}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockGetBlockReceiptsRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockGetBlockReceiptsResponse)),
		}, nil
	}

	receipts, err := client.GetBlockReceipts(
		t.Context(),
		types.MustBlockNumberFromHex("0x1"),
	)

	require.NoError(t, err)
	require.Len(t, receipts, 1)
	assert.Equal(t, types.MustHashFromHex("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", types.PadNone), receipts[0].TransactionHash)
	assert.Equal(t, uint64(17), receipts[0].TransactionIndex)
	assert.Equal(t, types.MustHashFromHex("0x1111111111111111111111111111111111111111111111111111111111111111", types.PadNone), receipts[0].BlockHash)
	assert.Equal(t, big.NewInt(0x2222), receipts[0].BlockNumber)
	assert.Equal(t, (*types.Address)(nil), receipts[0].ContractAddress)
	assert.Equal(t, hexutil.MustHexToBigInt("0x33333").Uint64(), receipts[0].CumulativeGasUsed)
	assert.Equal(t, hexutil.MustHexToBigInt("0x4444444444"), receipts[0].EffectiveGasPrice)
	assert.Equal(t, hexutil.MustHexToBigInt("0x66666").Uint64(), receipts[0].GasUsed)
	assert.Equal(t, types.MustAddressFromHex("0x5555555555555555555555555555555555555555"), receipts[0].From)
	assert.Equal(t, types.MustAddressFromHex("0x7777777777777777777777777777777777777777"), receipts[0].To)
	assert.Equal(t, hexutil.MustHexToBytes("0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000800000200000000000000000000000000000"), receipts[0].LogsBloom)
	status := uint64(1)
	assert.Equal(t, &status, receipts[0].Status)
	require.Len(t, receipts[0].Logs, 1)
	assert.Equal(t, types.MustHashFromHexPtr("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", types.PadNone), receipts[0].Logs[0].TransactionHash)
	assert.Equal(t, uint64(17), *receipts[0].Logs[0].TransactionIndex)
	assert.Equal(t, types.MustHashFromHexPtr("0x1111111111111111111111111111111111111111111111111111111111111111", types.PadNone), receipts[0].Logs[0].BlockHash)
	assert.Equal(t, big.NewInt(0x2222), receipts[0].Logs[0].BlockNumber)
	assert.Equal(t, uint64(8), *receipts[0].Logs[0].LogIndex)
	assert.Equal(t, hexutil.MustHexToBytes("0x000000000000000000000000398137383b3d25c92898c656696e41950e47316b00000000000000000000000000000000000000000000000000000000000cee6100000000000000000000000000000000000000000000000000000000000ac3e100000000000000000000000000000000000000000000000000000000005baf35"), receipts[0].Logs[0].Data)
	assert.Equal(t, types.MustAddressFromHex("0x7777777777777777777777777777777777777777"), receipts[0].Logs[0].Address)
	assert.Equal(t, []types.Hash{types.MustHashFromHex("0x9999999999999999999999999999999999999999999999999999999999999999", types.PadNone)}, receipts[0].Logs[0].Topics)
	assert.Equal(t, false, receipts[0].Logs[0].Removed)
}

// TODO: SubscribeNewPendingTransactions

const mockGetLogsRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_getLogs",
	  "params": [
		{
		  "fromBlock": "0x1",
		  "toBlock": "0x2",
		  "address": "0x3333333333333333333333333333333333333333",
		  "topics": [
	        "0x4444444444444444444444444444444444444444444444444444444444444444"
		  ]
		}
	  ]
	}
`

const mockGetLogsResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": [
		{
		  "address": "0x3333333333333333333333333333333333333333",
		  "topics": [
	        "0x4444444444444444444444444444444444444444444444444444444444444444"
		  ],
		  "data": "0x68656c6c6f21",
		  "blockNumber": "0x1",
		  "transactionHash": "0x4444444444444444444444444444444444444444444444444444444444444444",
		  "transactionIndex": "0x0",
		  "blockHash": "0x4444444444444444444444444444444444444444444444444444444444444444",
		  "logIndex": "0x0",
		  "removed": false
		}
	  ]
	}
`

func TestBaseClient_GetLogs(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsCommon{Transport: httpMock}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockGetLogsRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockGetLogsResponse)),
		}, nil
	}

	from := types.MustBlockNumberFromHex("0x1")
	to := types.MustBlockNumberFromHex("0x2")
	logs, err := client.GetLogs(t.Context(), &types.FilterLogsQuery{
		FromBlock: &from,
		ToBlock:   &to,
		Address:   []types.Address{types.MustAddressFromHex("0x3333333333333333333333333333333333333333")},
		Topics: [][]types.Hash{
			{types.MustHashFromHex("0x4444444444444444444444444444444444444444444444444444444444444444", types.PadNone)},
		},
	})

	require.NoError(t, err)
	require.Len(t, logs, 1)
	assert.Equal(t, types.MustAddressFromHex("0x3333333333333333333333333333333333333333"), logs[0].Address)
	assert.Equal(t, types.MustHashFromHex("0x4444444444444444444444444444444444444444444444444444444444444444", types.PadNone), logs[0].Topics[0])
	assert.Equal(t, hexutil.MustHexToBytes("0x68656c6c6f21"), logs[0].Data)
	assert.Equal(t, big.NewInt(1), logs[0].BlockNumber)
	assert.Equal(t, types.MustHashFromHexPtr("0x4444444444444444444444444444444444444444444444444444444444444444", types.PadNone), logs[0].TransactionHash)
	assert.Equal(t, uint64(0), *logs[0].TransactionIndex)
	assert.Equal(t, types.MustHashFromHexPtr("0x4444444444444444444444444444444444444444444444444444444444444444", types.PadNone), logs[0].BlockHash)
	assert.Equal(t, uint64(0), *logs[0].LogIndex)
	assert.Equal(t, false, logs[0].Removed)
}

// TODO: SubscribeLogs

const mockChanIDRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_chainId",
	  "params": []
	}
`

const mockChanIDResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,	
	  "result": "0x1"
	}
`

func TestBaseClient_ChainID(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsCommon{Transport: httpMock}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockChanIDRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockChanIDResponse)),
		}, nil
	}

	chainID, err := client.ChainID(t.Context())
	require.NoError(t, err)
	assert.Equal(t, uint64(1), chainID)
}

const mockBlockNumberRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_blockNumber",
	  "params": []
	}
`

const mockBlockNumberResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": "0x1"
	}
`

func TestBaseClient_BlockNumber(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsCommon{Transport: httpMock}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockBlockNumberRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockBlockNumberResponse)),
		}, nil
	}

	blockNumber, err := client.BlockNumber(t.Context())

	require.NoError(t, err)
	assert.Equal(t, big.NewInt(1), blockNumber)
}

const mockGasPriceRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_gasPrice",
	  "params": []
	}
`

const mockGasPriceResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": "0x09184e72a000"
	}
`

func TestBaseClient_GasPrice(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsCommon{Transport: httpMock}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockGasPriceRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockGasPriceResponse)),
		}, nil
	}

	gasPrice, err := client.GasPrice(t.Context())

	require.NoError(t, err)
	assert.Equal(t, big.NewInt(10000000000000), gasPrice)
}

const mockMaxPriorityFeePerGasRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_maxPriorityFeePerGas",
	  "params": []
	}
`

const mockMaxPriorityFeePerGasResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": "0x1"
	}
`

func TestBaseClient_MaxPriorityFeePerGas(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsCommon{Transport: httpMock}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockMaxPriorityFeePerGasRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockMaxPriorityFeePerGasResponse)),
		}, nil
	}

	gasPrice, err := client.MaxPriorityFeePerGas(t.Context())

	require.NoError(t, err)
	assert.Equal(t, hexutil.MustHexToBigInt("0x1"), gasPrice)
}

const mockBlobBaseFeeRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_blobBaseFee",
	  "params": []
	}
`

const mockBlobBaseFeeResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": "0x1"
	}
`

func TestBaseClient_BlobBaseFee(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsCommon{Transport: httpMock}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockBlobBaseFeeRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockBlobBaseFeeResponse)),
		}, nil
	}

	gasPrice, err := client.BlobBaseFee(t.Context())

	require.NoError(t, err)
	assert.Equal(t, hexutil.MustHexToBigInt("0x1"), gasPrice)
}

func TestBaseClient_SubscribeNewHeads(t *testing.T) {
	// TODO: Veirify

	streamMock := newStreamMock(t)
	client := &MethodsCommon{Transport: streamMock}

	ch := make(chan json.RawMessage)
	streamMock.SubscribeMocks = []subscribeMock{
		{
			ArgMethod: "newHeads",
			ArgParams: []any{},
			RetCh:     ch,
			RetID:     "0x1",
			RetErr:    nil,
		},
	}
	streamMock.UnsubscribeMocks = []unsubscribeMock{
		{
			ArgID:     "0x1",
			ResultErr: nil,
		},
	}

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	resultCh, err := client.SubscribeNewHeads(ctx)
	require.NoError(t, err)
	assert.NotNil(t, resultCh)
}

func TestBaseClient_SubscribeNewPendingTransactions(t *testing.T) {
	// TODO: Veirify

	streamMock := newStreamMock(t)
	client := &MethodsCommon{Transport: streamMock}

	ch := make(chan json.RawMessage)
	streamMock.SubscribeMocks = []subscribeMock{
		{
			ArgMethod: "newPendingTransactions",
			ArgParams: []any{},
			RetCh:     ch,
			RetID:     "0x2",
			RetErr:    nil,
		},
	}
	streamMock.UnsubscribeMocks = []unsubscribeMock{
		{
			ArgID:     "0x2",
			ResultErr: nil,
		},
	}

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	resultCh, err := client.SubscribeNewPendingTransactions(ctx)
	require.NoError(t, err)
	assert.NotNil(t, resultCh)
}

func TestBaseClient_SubscribeLogs(t *testing.T) {
	// TODO: Veirify

	streamMock := newStreamMock(t)
	client := &MethodsCommon{Transport: streamMock}

	ch := make(chan json.RawMessage)
	from := types.MustBlockNumberFromHex("0x1")
	to := types.MustBlockNumberFromHex("0x2")
	query := &types.FilterLogsQuery{
		FromBlock: &from,
		ToBlock:   &to,
		Address:   []types.Address{types.MustAddressFromHex("0x3333333333333333333333333333333333333333")},
		Topics: [][]types.Hash{
			{types.MustHashFromHex("0x4444444444444444444444444444444444444444444444444444444444444444", types.PadNone)},
		},
	}
	streamMock.SubscribeMocks = []subscribeMock{
		{
			ArgMethod: "logs",
			ArgParams: []any{query},
			RetCh:     ch,
			RetID:     "0x3",
			RetErr:    nil,
		},
	}
	streamMock.UnsubscribeMocks = []unsubscribeMock{
		{
			ArgID:     "0x3",
			ResultErr: nil,
		},
	}

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	resultCh, err := client.SubscribeLogs(ctx, query)
	require.NoError(t, err)
	assert.NotNil(t, resultCh)
}

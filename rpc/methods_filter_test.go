package rpc

import (
	"bytes"
	"context"
	"io"
	"math/big"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/defiweb/go-eth/types"
)

const mockNewFilterRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_newFilter",
	  "params": [
		{
		  "fromBlock": "0x1",
		  "toBlock": "0x2",
		  "address": "0x3333333333333333333333333333333333333333",
		  "topics": ["0x4444444444444444444444444444444444444444444444444444444444444444"]
		}
	  ]
	}
`

const mockNewFilterResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": "0x1"
	}
`

func TestBaseClient_NewFilter(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsFilter{&ClientContext{Transport: httpMock}}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockNewFilterRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockNewFilterResponse)),
		}, nil
	}

	from := types.MustBlockNumberFromHex("0x1")
	to := types.MustBlockNumberFromHex("0x2")
	id, err := client.NewFilter(context.Background(), &types.FilterLogsQuery{
		FromBlock: &from,
		ToBlock:   &to,
		Address:   []types.Address{types.MustAddressFromHex("0x3333333333333333333333333333333333333333")},
		Topics: [][]types.Hash{
			{types.MustHashFromHex("0x4444444444444444444444444444444444444444444444444444444444444444", types.PadNone)},
		},
	})

	require.NoError(t, err)
	assert.Equal(t, "1", id.String())
}

const mockNewBlockFilterRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_newBlockFilter",
	  "params": []
	}
`

const mockNewBlockFilterResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": "0x2"
	}
`

func TestBaseClient_NewBlockFilter(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsFilter{&ClientContext{Transport: httpMock}}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockNewBlockFilterRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockNewBlockFilterResponse)),
		}, nil
	}

	id, err := client.NewBlockFilter(context.Background())

	require.NoError(t, err)
	assert.Equal(t, "2", id.String())
}

const mockNewPendingTransactionFilterRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_newPendingTransactionFilter",
	  "params": []
	}
`

const mockNewPendingTransactionFilterResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": "0x3"
	}
`

func TestBaseClient_NewPendingTransactionFilter(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsFilter{&ClientContext{Transport: httpMock}}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockNewPendingTransactionFilterRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockNewPendingTransactionFilterResponse)),
		}, nil
	}

	id, err := client.NewPendingTransactionFilter(context.Background())

	require.NoError(t, err)
	assert.Equal(t, "3", id.String())
}

const mockUninstallFilterRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_uninstallFilter",
	  "params": ["0x1"]
	}
`

const mockUninstallFilterResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": true
	}
`

func TestBaseClient_UninstallFilter(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsFilter{&ClientContext{Transport: httpMock}}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockUninstallFilterRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockUninstallFilterResponse)),
		}, nil
	}

	result, err := client.UninstallFilter(context.Background(), big.NewInt(1))

	require.NoError(t, err)
	assert.True(t, result)
}

const mockGetFilterChangesRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_getFilterChanges",
	  "params": ["0x1"]
	}
`

const mockGetFilterChangesResponse = `
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

func TestBaseClient_GetFilterChanges(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsFilter{&ClientContext{Transport: httpMock}}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockGetFilterChangesRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockGetFilterChangesResponse)),
		}, nil
	}

	logs, err := client.GetFilterChanges(context.Background(), big.NewInt(1))

	require.NoError(t, err)
	require.Len(t, logs, 1)
	assert.Equal(t, types.MustAddressFromHex("0x3333333333333333333333333333333333333333"), logs[0].Address)
	assert.Equal(t, types.MustHashFromHex("0x4444444444444444444444444444444444444444444444444444444444444444", types.PadNone), logs[0].Topics[0])
}

const mockGetFilterLogsRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_getFilterLogs",
	  "params": ["0x1"]
	}
`

const mockGetFilterLogsResponse = `
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
		  "blockNumber": "0x2",
		  "transactionHash": "0x5555555555555555555555555555555555555555555555555555555555555555",
		  "transactionIndex": "0x1",
		  "blockHash": "0x5555555555555555555555555555555555555555555555555555555555555555",
		  "logIndex": "0x1",
		  "removed": false
		}
	  ]
	}
`

func TestBaseClient_GetFilterLogs(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsFilter{&ClientContext{Transport: httpMock}}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockGetFilterLogsRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockGetFilterLogsResponse)),
		}, nil
	}

	logs, err := client.GetFilterLogs(context.Background(), big.NewInt(1))

	require.NoError(t, err)
	require.Len(t, logs, 1)
	assert.Equal(t, types.MustAddressFromHex("0x3333333333333333333333333333333333333333"), logs[0].Address)
	assert.Equal(t, types.MustHashFromHex("0x4444444444444444444444444444444444444444444444444444444444444444", types.PadNone), logs[0].Topics[0])
	assert.Equal(t, big.NewInt(2), logs[0].BlockNumber)
}

const mockGetBlockFilterChangesRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_getFilterChanges",
	  "params": ["0x2"]
	}
`

const mockGetBlockFilterChangesResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": [
	    "0x1111111111111111111111111111111111111111111111111111111111111111",
	    "0x2222222222222222222222222222222222222222222222222222222222222222"
	  ]
	}
`

func TestBaseClient_GetBlockFilterChanges(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsFilter{&ClientContext{Transport: httpMock}}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockGetBlockFilterChangesRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockGetBlockFilterChangesResponse)),
		}, nil
	}

	hashes, err := client.GetBlockFilterChanges(context.Background(), big.NewInt(2))

	require.NoError(t, err)
	require.Len(t, hashes, 2)
	assert.Equal(t, types.MustHashFromHex("0x1111111111111111111111111111111111111111111111111111111111111111", types.PadNone), hashes[0])
	assert.Equal(t, types.MustHashFromHex("0x2222222222222222222222222222222222222222222222222222222222222222", types.PadNone), hashes[1])
}

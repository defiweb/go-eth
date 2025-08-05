package rpc

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/defiweb/go-eth/hexutil"
	"github.com/defiweb/go-eth/types"
)

const mockCancelPrivateTransactionRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_cancelPrivateTransaction",
	  "params": [
	    "0x1111111111111111111111111111111111111111111111111111111111111111"
	  ]
	}
`

const mockCancelPrivateTransactionResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": true
	}
`

func TestBaseClient_CancelPrivateTransaction(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsPrivateTransaction{Transport: httpMock}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockCancelPrivateTransactionRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockCancelPrivateTransactionResponse)),
		}, nil
	}

	result, err := client.CancelPrivateTransaction(
		context.Background(),
		types.MustHashFromHex("0x1111111111111111111111111111111111111111111111111111111111111111", types.PadNone),
	)

	require.NoError(t, err)
	assert.True(t, result)
}

const mockSendPrivateTransactionRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_sendPrivateTransaction",
	  "params": [
	    "0xf893808609184e72a0008276c094d46e8dd67c5d32be8058bb8eb970870f072445678502540be400a9d46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f07244567511a02222222222222222222222222222222222222222222222222222222222222222a03333333333333333333333333333333333333333333333333333333333333333"
	  ]
	}
`

const mockSendPrivateTransactionResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": "0x1111111111111111111111111111111111111111111111111111111111111111"
	}
`

func TestBaseClient_SendPrivateTransaction(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsPrivateTransaction{Transport: httpMock}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockSendPrivateTransactionRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockSendPrivateTransactionResponse)),
		}, nil
	}

	txHash, err := client.SendPrivateTransaction(
		context.Background(),
		hexutil.MustHexToBytes("0xf893808609184e72a0008276c094d46e8dd67c5d32be8058bb8eb970870f072445678502540be400a9d46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f07244567511a02222222222222222222222222222222222222222222222222222222222222222a03333333333333333333333333333333333333333333333333333333333333333"),
	)

	require.NoError(t, err)
	assert.Equal(t, types.MustHashFromHex("0x1111111111111111111111111111111111111111111111111111111111111111", types.PadNone), *txHash)
}

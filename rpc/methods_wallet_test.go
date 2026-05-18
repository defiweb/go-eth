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

const mockAccountsRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_accounts",
	  "params": []
	}
`

const mockAccountsResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": [
	    "0x1111111111111111111111111111111111111111",
	    "0x2222222222222222222222222222222222222222"
	  ]
	}
`

func TestBaseClient_Accounts(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsWallet{Transport: httpMock}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockAccountsRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockAccountsResponse)),
		}, nil
	}

	accounts, err := client.Accounts(context.Background())

	require.NoError(t, err)
	require.Len(t, accounts, 2)
	assert.Equal(t, types.MustAddressFromHex("0x1111111111111111111111111111111111111111"), accounts[0])
	assert.Equal(t, types.MustAddressFromHex("0x2222222222222222222222222222222222222222"), accounts[1])
}

const mockSignRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_sign",
	  "params": [
	    "0x1111111111111111111111111111111111111111",
	    "0x48656c6c6f20576f726c64"
	  ]
	}
`

const mockSignResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": "0x3333333333333333333333333333333333333333333333333333333333333333444444444444444444444444444444444444444444444444444444444444444455"
	}
`

func TestBaseClient_Sign(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsWallet{Transport: httpMock}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockSignRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockSignResponse)),
		}, nil
	}

	signature, err := client.Sign(
		context.Background(),
		types.MustAddressFromHex("0x1111111111111111111111111111111111111111"),
		[]byte("Hello World"),
	)

	require.NoError(t, err)
	assert.Equal(t, hexutil.MustHexToBytes("0x3333333333333333333333333333333333333333333333333333333333333333"), signature.R.Bytes())
	assert.Equal(t, hexutil.MustHexToBytes("0x4444444444444444444444444444444444444444444444444444444444444444"), signature.S.Bytes())
	assert.Equal(t, uint64(0x55), signature.V.Uint64())
}

const mockSignTransactionRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_signTransaction",
	  "params": [
		{
		  "from": "0x1111111111111111111111111111111111111111",
		  "to": "0x2222222222222222222222222222222222222222",
		  "gas": "0x5208",
		  "gasPrice": "0x9184e72a000",
		  "value": "0x2540be400",
		  "input": "0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675"
		}
	  ]
	}
`

const mockSignTransactionResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": {
	    "raw": "0xf893808609184e72a0008276c094d46e8dd67c5d32be8058bb8eb970870f072445678502540be400a9d46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f07244567511a02222222222222222222222222222222222222222222222222222222222222222a03333333333333333333333333333333333333333333333333333333333333333"
	  }
	}
`

func TestBaseClient_SignTransaction(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsWallet{Transport: httpMock}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockSignTransactionRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockSignTransactionResponse)),
		}, nil
	}

	tx := types.NewTransactionLegacy()
	tx.SetFrom(*types.MustAddressFromHexPtr("0x1111111111111111111111111111111111111111"))
	tx.SetTo(*types.MustAddressFromHexPtr("0x2222222222222222222222222222222222222222"))
	tx.SetValue(hexutil.MustHexToBigInt("0x2540be400"))
	tx.SetGasLimit(0x5208)
	tx.SetGasPrice(hexutil.MustHexToBigInt("0x9184e72a000"))
	tx.SetInput(hexutil.MustHexToBytes("0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675"))
	rawTx, err := client.SignTransaction(context.Background(), tx)

	require.NoError(t, err)
	assert.Equal(t, "0xf893808609184e72a0008276c094d46e8dd67c5d32be8058bb8eb970870f072445678502540be400a9d46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f07244567511a02222222222222222222222222222222222222222222222222222222222222222a03333333333333333333333333333333333333333333333333333333333333333", hexutil.BytesToHex(rawTx))
}

const mockSendTransactionRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_sendTransaction",
	  "params": [
		{
		  "from": "0x1111111111111111111111111111111111111111",
		  "to": "0x2222222222222222222222222222222222222222",
		  "gas": "0x5208",
		  "gasPrice": "0x9184e72a000",
		  "value": "0x2540be400",
		  "input": "0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675"
		}
	  ]
	}
`

const mockSendTransactionResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": "0x1111111111111111111111111111111111111111111111111111111111111111"
	}
`

func TestBaseClient_SendTransaction(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsWallet{Transport: httpMock}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockSendTransactionRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockSendTransactionResponse)),
		}, nil
	}

	tx := types.NewTransactionLegacy()
	tx.SetFrom(*types.MustAddressFromHexPtr("0x1111111111111111111111111111111111111111"))
	tx.SetTo(*types.MustAddressFromHexPtr("0x2222222222222222222222222222222222222222"))
	tx.SetValue(hexutil.MustHexToBigInt("0x2540be400"))
	tx.SetGasLimit(0x5208)
	tx.SetGasPrice(hexutil.MustHexToBigInt("0x9184e72a000"))
	tx.SetInput(hexutil.MustHexToBytes("0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675"))
	txHash, err := client.SendTransaction(context.Background(), tx)

	require.NoError(t, err)
	assert.Equal(t, types.MustHashFromHex("0x1111111111111111111111111111111111111111111111111111111111111111", types.PadNone), *txHash)
}

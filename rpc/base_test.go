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

	"github.com/defiweb/go-eth/hexutil"
	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
)

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type httpMock struct {
	*transport.HTTP
	Request  *http.Request
	Response *http.Response
}

func newHTTPMock() *httpMock {
	h := &httpMock{}
	h.HTTP, _ = transport.NewHTTP(transport.HTTPOptions{
		URL: "http://localhost",
		HTTPClient: &http.Client{
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				h.Request = req
				return h.Response, nil
			}),
		},
	})
	return h
}

type keyMock struct {
	addressCallback         func() types.Address
	signHashCallback        func(hash types.Hash) (*types.Signature, error)
	signMessageCallback     func(data []byte) (*types.Signature, error)
	signTransactionCallback func(tx *types.Transaction) error
}

func (k *keyMock) Address() types.Address {
	return k.addressCallback()
}

func (k *keyMock) SignHash(hash types.Hash) (*types.Signature, error) {
	return k.signHashCallback(hash)
}

func (k *keyMock) SignMessage(data []byte) (*types.Signature, error) {
	return k.signMessageCallback(data)
}

func (k *keyMock) SignTransaction(tx *types.Transaction) error {
	return k.signTransactionCallback(tx)
}

func (k *keyMock) VerifyHash(hash types.Hash, sig types.Signature) bool {
	return false
}

func (k keyMock) VerifyMessage(data []byte, sig types.Signature) bool {
	return false
}

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
	client := &baseClient{transport: httpMock}

	httpMock.Response = &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(mockChanIDResponse)),
	}

	chainID, err := client.ChainID(context.Background())
	require.NoError(t, err)
	assert.JSONEq(t, mockChanIDRequest, readBody(httpMock.Request))
	assert.Equal(t, uint64(1), chainID)
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
	client := &baseClient{transport: httpMock}

	httpMock.Response = &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(mockGasPriceResponse)),
	}

	gasPrice, err := client.GasPrice(context.Background())
	require.NoError(t, err)
	assert.JSONEq(t, mockGasPriceRequest, readBody(httpMock.Request))
	assert.Equal(t, big.NewInt(10000000000000), gasPrice)
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
	client := &baseClient{transport: httpMock}

	httpMock.Response = &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(mockBlockNumberResponse)),
	}

	blockNumber, err := client.BlockNumber(context.Background())

	require.NoError(t, err)
	assert.JSONEq(t, mockBlockNumberRequest, readBody(httpMock.Request))
	assert.Equal(t, big.NewInt(1), blockNumber)
}

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
	client := &baseClient{transport: httpMock}

	httpMock.Response = &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(mockGetBalanceResponse)),
	}

	balance, err := client.GetBalance(
		context.Background(),
		types.MustAddressFromHex("0x1111111111111111111111111111111111111111"),
		types.LatestBlockNumber,
	)

	require.NoError(t, err)
	assert.JSONEq(t, mockGetBalanceRequest, readBody(httpMock.Request))
	assert.Equal(t, big.NewInt(158972490234375000), balance)
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
	client := &baseClient{transport: httpMock}

	httpMock.Response = &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(mockGetStorageAtResponse)),
	}

	storage, err := client.GetStorageAt(
		context.Background(),
		types.MustAddressFromHex("0x1111111111111111111111111111111111111111"),
		types.MustHashFromHex("0x2222222222222222222222222222222222222222222222222222222222222222", types.PadNone),
		types.MustBlockNumberFromHex("0x1"),
	)

	require.NoError(t, err)
	assert.JSONEq(t, mockGetStorageAtRequest, readBody(httpMock.Request))
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
	client := &baseClient{transport: httpMock}

	httpMock.Response = &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(mockGetTransactionCountResponse)),
	}

	transactionCount, err := client.GetTransactionCount(
		context.Background(),
		types.MustAddressFromHex("0x1111111111111111111111111111111111111111"),
		types.MustBlockNumberFromHex("0x1"),
	)

	require.NoError(t, err)
	assert.JSONEq(t, mockGetTransactionCountRequest, readBody(httpMock.Request))
	assert.Equal(t, uint64(1), transactionCount)
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
	client := &baseClient{transport: httpMock}

	httpMock.Response = &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(mockGetBlockTransactionCountByHashResponse)),
	}

	transactionCount, err := client.GetBlockTransactionCountByHash(
		context.Background(),
		types.MustHashFromHex("0x1111111111111111111111111111111111111111111111111111111111111111", types.PadNone),
	)
	require.NoError(t, err)
	assert.JSONEq(t, mockGetBlockTransactionCountByHashRequest, readBody(httpMock.Request))
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
	client := &baseClient{transport: httpMock}

	httpMock.Response = &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(mockGetBlockTransactionCountByNumberResponse)),
	}

	transactionCount, err := client.GetBlockTransactionCountByNumber(
		context.Background(),
		types.MustBlockNumberFromHex("0x1"),
	)
	require.NoError(t, err)
	assert.JSONEq(t, mockGetBlockTransactionCountByNumberRequest, readBody(httpMock.Request))
	assert.Equal(t, uint64(2), transactionCount)
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
	client := &baseClient{transport: httpMock}

	httpMock.Response = &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(mockGetUncleCountByBlockHashResponse)),
	}

	uncleCount, err := client.GetUncleCountByBlockHash(
		context.Background(),
		types.MustHashFromHex("0x1111111111111111111111111111111111111111111111111111111111111111", types.PadNone),
	)
	require.NoError(t, err)
	assert.JSONEq(t, mockGetUncleCountByBlockHashRequest, readBody(httpMock.Request))
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
	client := &baseClient{transport: httpMock}

	httpMock.Response = &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(mockGetUncleCountByBlockNumberResponse)),
	}

	uncleCount, err := client.GetUncleCountByBlockNumber(
		context.Background(),
		types.MustBlockNumberFromHex("0x1"),
	)
	require.NoError(t, err)
	assert.JSONEq(t, mockGetUncleCountByBlockNumberRequest, readBody(httpMock.Request))
	assert.Equal(t, uint64(2), uncleCount)
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
	client := &baseClient{transport: httpMock}

	httpMock.Response = &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(mockGetCodeResponse)),
	}

	code, err := client.GetCode(
		context.Background(),
		types.MustAddressFromHex("0x1111111111111111111111111111111111111111"),
		types.MustBlockNumberFromHex("0x2"),
	)
	require.NoError(t, err)
	assert.JSONEq(t, mockGetCodeRequest, readBody(httpMock.Request))
	assert.Equal(t, hexToBytes("0x3333333333333333333333333333333333333333333333333333333333333333"), code)
}

const mockSignRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_sign",
	  "params": [
		"0x1111111111111111111111111111111111111111",
		"0x416c6c20796f75722062617365206172652062656c6f6e6720746f207573"
	  ]
	}
`

const mockSignResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": "0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"
	}
`

func TestBaseClient_Sign(t *testing.T) {
	httpMock := newHTTPMock()
	client := &baseClient{transport: httpMock}

	httpMock.Response = &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(mockSignResponse)),
	}

	signature, err := client.Sign(
		context.Background(),
		types.MustAddressFromHex("0x1111111111111111111111111111111111111111"),
		[]byte("All your base are belong to us"),
	)
	require.NoError(t, err)
	assert.JSONEq(t, mockSignRequest, readBody(httpMock.Request))
	assert.Equal(t, types.MustSignatureFromHex("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"), *signature)
}

const mockSignTransactionRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_signTransaction",
	  "params": [
		{
		  "from": "0xb60e8dd61c5d32be8058bb8eb970870f07233155",
		  "to": "0xd46e8dd67c5d32be8058bb8eb970870f07244567",
		  "gas": "0x76c0",
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
		"raw": "0xf893808609184e72a0008276c094d46e8dd67c5d32be8058bb8eb970870f072445678502540be400a9d46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f07244567511a02222222222222222222222222222222222222222222222222222222222222222a03333333333333333333333333333333333333333333333333333333333333333",
		"tx": {
		  "nonce": "0x0",
		  "gasPrice": "0x09184e72a000",
		  "gas": "0x76c0",
		  "to": "0xd46e8dd67c5d32be8058bb8eb970870f07244567",
		  "value": "0x2540be400",
		  "input": "0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675",
		  "v": "0x11",
		  "r": "0x2222222222222222222222222222222222222222222222222222222222222222",
		  "s": "0x3333333333333333333333333333333333333333333333333333333333333333"
		}
	  }
	}
`

func TestBaseClient_SignTransaction(t *testing.T) {
	httpMock := newHTTPMock()
	client := &baseClient{transport: httpMock}

	httpMock.Response = &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(mockSignTransactionResponse)),
	}

	from := types.MustAddressFromHex("0xb60e8dd61c5d32be8058bb8eb970870f07233155")
	to := types.MustAddressFromHex("0xd46e8dd67c5d32be8058bb8eb970870f07244567")
	gasLimit := uint64(30400)
	chainID := uint64(1)
	raw, tx, err := client.SignTransaction(
		context.Background(),
		types.Transaction{
			ChainID: &chainID,
			Call: types.Call{
				From:     &from,
				To:       &to,
				GasLimit: &gasLimit,
				GasPrice: big.NewInt(10000000000000),
				Value:    big.NewInt(10000000000),
				Input:    hexToBytes("0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675"),
			},
		},
	)
	require.NoError(t, err)
	assert.JSONEq(t, mockSignTransactionRequest, readBody(httpMock.Request))
	assert.Equal(t, hexToBytes("0xf893808609184e72a0008276c094d46e8dd67c5d32be8058bb8eb970870f072445678502540be400a9d46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f07244567511a02222222222222222222222222222222222222222222222222222222222222222a03333333333333333333333333333333333333333333333333333333333333333"), raw)
	assert.Equal(t, &to, tx.To)
	assert.Equal(t, uint64(30400), *tx.GasLimit)
	assert.Equal(t, big.NewInt(10000000000000), tx.GasPrice)
	assert.Equal(t, big.NewInt(10000000000), tx.Value)
	assert.Equal(t, hexToBytes("0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675"), tx.Input)
	assert.Equal(t, uint8(0x11), tx.Signature.Bytes()[64])
	assert.Equal(t, hexToBytes("0x2222222222222222222222222222222222222222222222222222222222222222"), tx.Signature.Bytes()[:32])
	assert.Equal(t, hexToBytes("0x3333333333333333333333333333333333333333333333333333333333333333"), tx.Signature.Bytes()[32:64])
}

const mockSendTransactionRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_sendTransaction",
	  "params": [
	    {
		  "from": "0xb60e8dd61c5d32be8058bb8eb970870f07233155",
		  "to": "0xd46e8dd67c5d32be8058bb8eb970870f07244567",
		  "gas": "0x76c0",
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
	client := &baseClient{transport: httpMock}

	httpMock.Response = &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(mockSendTransactionResponse)),
	}

	from := types.MustAddressFromHex("0xb60e8dd61c5d32be8058bb8eb970870f07233155")
	to := types.MustAddressFromHex("0xd46e8dd67c5d32be8058bb8eb970870f07244567")
	gasLimit := uint64(30400)
	chainID := uint64(1)
	txHash, err := client.SendTransaction(
		context.Background(),
		types.Transaction{
			ChainID: &chainID,
			Call: types.Call{
				From:     &from,
				To:       &to,
				GasLimit: &gasLimit,
				GasPrice: big.NewInt(10000000000000),
				Value:    big.NewInt(10000000000),
				Input:    hexToBytes("0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675"),
			},
		},
	)
	require.NoError(t, err)
	assert.JSONEq(t, mockSendTransactionRequest, readBody(httpMock.Request))
	assert.Equal(t, types.MustHashFromHex("0x1111111111111111111111111111111111111111111111111111111111111111", types.PadNone), *txHash)
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
	client := &baseClient{transport: httpMock}

	httpMock.Response = &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(mockSendRawTransactionResponse)),
	}

	txHash, err := client.SendRawTransaction(
		context.Background(),
		hexToBytes("0xf893808609184e72a0008276c094d46e8dd67c5d32be8058bb8eb970870f072445678502540be400a9d46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f07244567511a02222222222222222222222222222222222222222222222222222222222222222a03333333333333333333333333333333333333333333333333333333333333333"),
	)
	require.NoError(t, err)
	assert.JSONEq(t, mockSendRawTransactionRequest, readBody(httpMock.Request))
	assert.Equal(t, types.MustHashFromHex("0x1111111111111111111111111111111111111111111111111111111111111111", types.PadNone), *txHash)
}

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
		  "data": "0x3333333333333333333333333333333333333333333333333333333333333333333333333333333333"
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
	client := &baseClient{transport: httpMock}

	httpMock.Response = &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(mockCallResponse)),
	}

	gasLimit := uint64(30400)
	call, err := client.Call(
		context.Background(),
		types.Call{
			From:     types.MustAddressFromHexPtr("0x1111111111111111111111111111111111111111"),
			To:       types.MustAddressFromHexPtr("0x2222222222222222222222222222222222222222"),
			GasLimit: &gasLimit,
			GasPrice: big.NewInt(10000000000000),
			Value:    big.NewInt(10000000000),
			Input:    hexToBytes("0x3333333333333333333333333333333333333333333333333333333333333333333333333333333333"),
		},
		types.MustBlockNumberFromHex("0x1"),
	)
	require.NoError(t, err)
	assert.JSONEq(t, mockCallRequest, readBody(httpMock.Request))
	assert.Equal(t, hexToBytes("0x00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000004000000000000000000000000d9c9cd5f6779558b6e0ed4e6acf6b1947e7fa1f300000000000000000000000078d1ad571a1a09d60d9bbf25894b44e4c8859595000000000000000000000000286834935f4a8cfb4ff4c77d5770c2775ae2b0e7000000000000000000000000b86e2b0ab5a4b1373e40c51a7c712c70ba2f9f8e"), call)
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
		  "data": "0x3333333333333333333333333333333333333333333333333333333333333333333333333333333333"
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
	client := &baseClient{transport: httpMock}

	httpMock.Response = &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(mockEstimateGasResponse)),
	}

	gasLimit := uint64(30400)
	gas, err := client.EstimateGas(
		context.Background(),
		types.Call{
			From:     types.MustAddressFromHexPtr("0x1111111111111111111111111111111111111111"),
			To:       types.MustAddressFromHexPtr("0x2222222222222222222222222222222222222222"),
			GasLimit: &gasLimit,
			GasPrice: big.NewInt(10000000000000),
			Value:    big.NewInt(10000000000),
			Input:    hexToBytes("0x3333333333333333333333333333333333333333333333333333333333333333333333333333333333"),
		},
		types.LatestBlockNumber,
	)
	require.NoError(t, err)
	assert.JSONEq(t, mockEstimateGasRequest, readBody(httpMock.Request))
	assert.Equal(t, uint64(21000), gas)
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

const mockBlockByNumberResponse = `
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

func TestBaseClient_BlockByNumber(t *testing.T) {
	httpMock := newHTTPMock()
	client := &baseClient{transport: httpMock}

	httpMock.Response = &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(mockBlockByNumberResponse)),
	}

	block, err := client.BlockByNumber(
		context.Background(),
		types.MustBlockNumberFromHex("0x1"),
		true,
	)
	require.NoError(t, err)
	assert.JSONEq(t, mockBlockByNumberRequest, readBody(httpMock.Request))
	assert.Equal(t, big.NewInt(0x11), block.Number)
	assert.Equal(t, types.MustHashFromHex("0x2222222222222222222222222222222222222222222222222222222222222222", types.PadNone), block.Hash)
	assert.Equal(t, types.MustHashFromHex("0x3333333333333333333333333333333333333333333333333333333333333333", types.PadNone), block.ParentHash)
	assert.Equal(t, hexToBigInt("0x4444444444444444"), block.Nonce)
	assert.Equal(t, types.MustHashFromHex("0x5555555555555555555555555555555555555555555555555555555555555555", types.PadNone), block.Sha3Uncles)
	assert.Equal(t, hexToBytes("0x66666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666666"), block.LogsBloom)
	assert.Equal(t, types.MustHashFromHex("0x7777777777777777777777777777777777777777777777777777777777777777", types.PadNone), block.TransactionsRoot)
	assert.Equal(t, types.MustHashFromHex("0x8888888888888888888888888888888888888888888888888888888888888888", types.PadNone), block.StateRoot)
	assert.Equal(t, types.MustHashFromHex("0x9999999999999999999999999999999999999999999999999999999999999999", types.PadNone), block.ReceiptsRoot)
	assert.Equal(t, types.MustAddressFromHex("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"), block.Miner)
	assert.Equal(t, hexToBigInt("0xbbbbbb"), block.Difficulty)
	assert.Equal(t, hexToBigInt("0xcccccc"), block.TotalDifficulty)
	assert.Equal(t, hexToBytes("0x0000000000000000000000000000000000000000000000000000000000000000"), block.ExtraData)
	assert.Equal(t, hexToBigInt("0xdddddd").Uint64(), block.Size)
	assert.Equal(t, hexToBigInt("0xeeeeee").Uint64(), block.GasLimit)
	assert.Equal(t, hexToBigInt("0xffffff").Uint64(), block.GasUsed)
	assert.Equal(t, int64(1424182926), block.Timestamp.Unix())
	require.Len(t, block.Transactions, 1)
	require.Len(t, block.Uncles, 1)
	assert.Equal(t, types.MustHashFromHexPtr("0x1111111111111111111111111111111111111111111111111111111111111111", types.PadNone), block.Transactions[0].Hash)
	assert.Equal(t, uint64(0x22), *block.Transactions[0].Nonce)
	assert.Equal(t, types.MustAddressFromHexPtr("0x5555555555555555555555555555555555555555"), block.Transactions[0].From)
	assert.Equal(t, types.MustAddressFromHexPtr("0x6666666666666666666666666666666666666666"), block.Transactions[0].To)
	assert.Equal(t, big.NewInt(10000000000), block.Transactions[0].Value)
	assert.Equal(t, uint64(30400), *block.Transactions[0].GasLimit)
	assert.Equal(t, big.NewInt(10000000000000), block.Transactions[0].GasPrice)
	assert.Equal(t, hexToBytes("0x777777777777"), block.Transactions[0].Input)
	assert.Equal(t, types.MustHashFromHex("0x8888888888888888888888888888888888888888888888888888888888888888", types.PadNone), block.Uncles[0])
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
	client := &baseClient{transport: httpMock}

	httpMock.Response = &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(mockBlockByNumberResponse)),
	}

	block, err := client.BlockByHash(
		context.Background(),
		types.MustHashFromHex("0x1111111111111111111111111111111111111111111111111111111111111111", types.PadNone),
		true,
	)

	require.NoError(t, err)
	assert.JSONEq(t, mockBlockByHashRequest, readBody(httpMock.Request))
	assert.Equal(t, types.MustHashFromHex("0x2222222222222222222222222222222222222222222222222222222222222222", types.PadNone), block.Hash)
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

const mockGetTransactionByHashResponse = `
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

func TestBaseClient_GetTransactionByHash(t *testing.T) {
	httpMock := newHTTPMock()
	client := &baseClient{transport: httpMock}

	httpMock.Response = &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(mockGetTransactionByHashResponse)),
	}

	tx, err := client.GetTransactionByHash(
		context.Background(),
		types.MustHashFromHex("0x1111111111111111111111111111111111111111111111111111111111111111", types.PadNone),
	)

	require.NoError(t, err)
	assert.JSONEq(t, mockGetTransactionByHashRequest, readBody(httpMock.Request))
	assert.Equal(t, types.MustHashFromHexPtr("0x1111111111111111111111111111111111111111111111111111111111111111", types.PadNone), tx.BlockHash)
	assert.Equal(t, big.NewInt(0x22), tx.BlockNumber)
	assert.Equal(t, types.MustHashFromHexPtr("0x4444444444444444444444444444444444444444444444444444444444444444", types.PadNone), tx.Hash)
	assert.Equal(t, types.MustAddressFromHexPtr("0x3333333333333333333333333333333333333333"), tx.From)
	assert.Equal(t, types.MustAddressFromHexPtr("0x7777777777777777777777777777777777777777"), tx.To)
	assert.Equal(t, big.NewInt(10000000000), tx.Value)
	assert.Equal(t, uint64(30400), *tx.GasLimit)
	assert.Equal(t, big.NewInt(10000000000000), tx.GasPrice)
	assert.Equal(t, hexToBytes("0x555555555555"), tx.Input)
	assert.Equal(t, uint64(0x66), *tx.Nonce)
	assert.Equal(t, hexToBigInt("0x0").Uint64(), *tx.TransactionIndex)
	assert.Equal(t, uint8(0x88), tx.Signature.Bytes()[64])
	assert.Equal(t, hexToBytes("0x9999999999999999999999999999999999999999999999999999999999999999"), tx.Signature.Bytes()[:32])
	assert.Equal(t, hexToBytes("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"), tx.Signature.Bytes()[32:64])
	assert.Equal(t, types.MustHashFromHexPtr("0x4444444444444444444444444444444444444444444444444444444444444444", types.PadNone), tx.Hash)
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
	client := &baseClient{transport: httpMock}

	httpMock.Response = &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(mockGetTransactionByHashResponse)),
	}

	tx, err := client.GetTransactionByBlockHashAndIndex(
		context.Background(),
		types.MustHashFromHex("0x1111111111111111111111111111111111111111111111111111111111111111", types.PadNone),
		0,
	)

	require.NoError(t, err)
	assert.JSONEq(t, mockGetTransactionByBlockHashAndIndexRequest, readBody(httpMock.Request))
	assert.Equal(t, types.MustHashFromHexPtr("0x4444444444444444444444444444444444444444444444444444444444444444", types.PadNone), tx.Hash)
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
	client := &baseClient{transport: httpMock}

	httpMock.Response = &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(mockGetTransactionByHashResponse)),
	}

	tx, err := client.GetTransactionByBlockNumberAndIndex(
		context.Background(),
		types.MustBlockNumberFromHex("0x1"),
		2,
	)

	require.NoError(t, err)
	assert.JSONEq(t, mockGetTransactionByBlockNumberAndIndexRequest, readBody(httpMock.Request))
	assert.Equal(t, types.MustHashFromHexPtr("0x4444444444444444444444444444444444444444444444444444444444444444", types.PadNone), tx.Hash)
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
	client := &baseClient{transport: httpMock}

	httpMock.Response = &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(mockGetTransactionReceiptResponse)),
	}

	receipt, err := client.GetTransactionReceipt(
		context.Background(),
		types.MustHashFromHex("0x1111111111111111111111111111111111111111111111111111111111111111", types.PadNone),
	)

	status := uint64(1)
	require.NoError(t, err)
	assert.JSONEq(t, mockGetTransactionReceiptRequest, readBody(httpMock.Request))
	assert.Equal(t, types.MustHashFromHex("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", types.PadNone), receipt.TransactionHash)
	assert.Equal(t, uint64(17), receipt.TransactionIndex)
	assert.Equal(t, types.MustHashFromHex("0x1111111111111111111111111111111111111111111111111111111111111111", types.PadNone), receipt.BlockHash)
	assert.Equal(t, big.NewInt(0x2222), receipt.BlockNumber)
	assert.Equal(t, (*types.Address)(nil), receipt.ContractAddress)
	assert.Equal(t, hexToBigInt("0x33333").Uint64(), receipt.CumulativeGasUsed)
	assert.Equal(t, hexToBigInt("0x4444444444"), receipt.EffectiveGasPrice)
	assert.Equal(t, hexToBigInt("0x66666").Uint64(), receipt.GasUsed)
	assert.Equal(t, types.MustAddressFromHex("0x5555555555555555555555555555555555555555"), receipt.From)
	assert.Equal(t, types.MustAddressFromHex("0x7777777777777777777777777777777777777777"), receipt.To)
	assert.Equal(t, hexToBytes("0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000800000200000000000000000000000000000"), receipt.LogsBloom)
	assert.Equal(t, &status, receipt.Status)
	require.Len(t, receipt.Logs, 1)
	assert.Equal(t, types.MustHashFromHexPtr("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", types.PadNone), receipt.Logs[0].TransactionHash)
	assert.Equal(t, uint64(17), *receipt.Logs[0].TransactionIndex)
	assert.Equal(t, types.MustHashFromHexPtr("0x1111111111111111111111111111111111111111111111111111111111111111", types.PadNone), receipt.Logs[0].BlockHash)
	assert.Equal(t, big.NewInt(0x2222), receipt.Logs[0].BlockNumber)
	assert.Equal(t, uint64(8), *receipt.Logs[0].LogIndex)
	assert.Equal(t, hexToBytes("0x000000000000000000000000398137383b3d25c92898c656696e41950e47316b00000000000000000000000000000000000000000000000000000000000cee6100000000000000000000000000000000000000000000000000000000000ac3e100000000000000000000000000000000000000000000000000000000005baf35"), receipt.Logs[0].Data)
	assert.Equal(t, types.MustAddressFromHex("0x7777777777777777777777777777777777777777"), receipt.Logs[0].Address)
	assert.Equal(t, []types.Hash{types.MustHashFromHex("0x9999999999999999999999999999999999999999999999999999999999999999", types.PadNone)}, receipt.Logs[0].Topics)
	assert.Equal(t, false, receipt.Logs[0].Removed)
}

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
	client := &baseClient{transport: httpMock}

	httpMock.Response = &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(mockGetLogsResponse)),
	}

	from := types.MustBlockNumberFromHex("0x1")
	to := types.MustBlockNumberFromHex("0x2")
	logs, err := client.GetLogs(context.Background(), types.FilterLogsQuery{
		FromBlock: &from,
		ToBlock:   &to,
		Address:   []types.Address{types.MustAddressFromHex("0x3333333333333333333333333333333333333333")},
		Topics: [][]types.Hash{
			{types.MustHashFromHex("0x4444444444444444444444444444444444444444444444444444444444444444", types.PadNone)},
		},
	})
	require.NoError(t, err)
	assert.JSONEq(t, mockGetLogsRequest, readBody(httpMock.Request))
	require.Len(t, logs, 1)
	assert.Equal(t, types.MustAddressFromHex("0x3333333333333333333333333333333333333333"), logs[0].Address)
	assert.Equal(t, types.MustHashFromHex("0x4444444444444444444444444444444444444444444444444444444444444444", types.PadNone), logs[0].Topics[0])
	assert.Equal(t, hexToBytes("0x68656c6c6f21"), logs[0].Data)
	assert.Equal(t, big.NewInt(1), logs[0].BlockNumber)
	assert.Equal(t, types.MustHashFromHexPtr("0x4444444444444444444444444444444444444444444444444444444444444444", types.PadNone), logs[0].TransactionHash)
	assert.Equal(t, uint64(0), *logs[0].TransactionIndex)
	assert.Equal(t, types.MustHashFromHexPtr("0x4444444444444444444444444444444444444444444444444444444444444444", types.PadNone), logs[0].BlockHash)
	assert.Equal(t, uint64(0), *logs[0].LogIndex)
	assert.Equal(t, false, logs[0].Removed)
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
	client := &baseClient{transport: httpMock}

	httpMock.Response = &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(mockMaxPriorityFeePerGasResponse)),
	}

	gasPrice, err := client.MaxPriorityFeePerGas(context.Background())
	require.NoError(t, err)
	assert.JSONEq(t, mockMaxPriorityFeePerGasRequest, readBody(httpMock.Request))
	assert.Equal(t, hexToBigInt("0x1"), gasPrice)
}

func readBody(r *http.Request) string {
	body, _ := io.ReadAll(r.Body)
	return string(body)
}

func hexToBytes(s string) []byte {
	b, _ := hexutil.HexToBytes(s)
	return b
}

func hexToBigInt(s string) *big.Int {
	b, _ := hexutil.HexToBigInt(s)
	return b
}

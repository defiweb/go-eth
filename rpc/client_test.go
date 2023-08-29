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

func TestClient_Sign(t *testing.T) {
	httpMock := newHTTPMock()
	keyMock := &keyMock{}
	client, _ := NewClient(WithTransport(httpMock), WithKeys(keyMock))

	keyMock.addressCallback = func() types.Address {
		return types.MustAddressFromHex("0x1111111111111111111111111111111111111111")
	}
	keyMock.signMessageCallback = func(message []byte) (*types.Signature, error) {
		return types.MustSignatureFromHexPtr("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"), nil
	}

	signature, err := client.Sign(
		context.Background(),
		types.MustAddressFromHex("0x1111111111111111111111111111111111111111"),
		[]byte("All your base are belong to us"),
	)
	require.NoError(t, err)
	assert.Equal(t, types.MustSignatureFromHex("0xa3a7b12762dbc5df6cfbedbecdf8a821929c6112d2634abbb0d99dc63ad914908051b2c8c7d159db49ad19bd01026156eedab2f3d8c1dfdd07d21c07a4bbdd846f"), *signature)
}

func TestClient_SignTransaction(t *testing.T) {
	httpMock := newHTTPMock()
	keyMock := &keyMock{}
	client, _ := NewClient(WithTransport(httpMock), WithKeys(keyMock))

	keyMock.addressCallback = func() types.Address {
		return types.MustAddressFromHex("0xb60e8dd61c5d32be8058bb8eb970870f07233155")
	}
	keyMock.signTransactionCallback = func(tx *types.Transaction) error {
		tx.Signature = types.MustSignatureFromHexPtr("0x2222222222222222222222222222222222222222222222222222222222222222333333333333333333333333333333333333333333333333333333333333333311")
		return nil
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

func TestClient_SendTransaction(t *testing.T) {
	httpMock := newHTTPMock()
	keyMock := &keyMock{}
	client, _ := NewClient(WithTransport(httpMock), WithKeys(keyMock))

	keyMock.addressCallback = func() types.Address {
		return types.MustAddressFromHex("0xb60e8dd61c5d32be8058bb8eb970870f07233155")
	}
	keyMock.signTransactionCallback = func(tx *types.Transaction) error {
		tx.Signature = types.MustSignatureFromHexPtr("0x2222222222222222222222222222222222222222222222222222222222222222333333333333333333333333333333333333333333333333333333333333333311")
		return nil
	}

	httpMock.ResponseMock = &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(mockSendRawTransactionResponse)),
	}

	from := types.MustAddressFromHex("0xb60e8dd61c5d32be8058bb8eb970870f07233155")
	to := types.MustAddressFromHex("0xd46e8dd67c5d32be8058bb8eb970870f07244567")
	gasLimit := uint64(30400)
	gasPrice := big.NewInt(10000000000000)
	value := big.NewInt(10000000000)
	input := hexToBytes("0xd46e8dd67c5d32be8d46e8dd67c5d32be8058bb8eb970870f072445675058bb8eb970870f072445675")
	chainID := uint64(1)
	txHash, tx, err := client.SendTransaction(
		context.Background(),
		types.Transaction{
			ChainID: &chainID,
			Call: types.Call{
				From:     &from,
				To:       &to,
				GasLimit: &gasLimit,
				GasPrice: gasPrice,
				Value:    value,
				Input:    input,
			},
		},
	)
	require.NoError(t, err)
	assert.JSONEq(t, mockSendRawTransactionRequest, readBody(httpMock.Request))
	assert.Equal(t, types.MustHashFromHex("0x1111111111111111111111111111111111111111111111111111111111111111", types.PadNone), *txHash)
	assert.Equal(t, &to, tx.To)
	assert.Equal(t, gasLimit, *tx.GasLimit)
	assert.Equal(t, gasPrice, tx.GasPrice)
	assert.Equal(t, value, tx.Value)
	assert.Equal(t, input, tx.Input)
}

func TestClient_Call(t *testing.T) {
	httpMock := newHTTPMock()
	client, _ := NewClient(
		WithTransport(httpMock),
		WithDefaultAddress(types.MustAddressFromHex("0x1111111111111111111111111111111111111111")),
	)

	httpMock.ResponseMock = &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(mockCallResponse)),
	}

	to := types.MustAddressFromHex("0x2222222222222222222222222222222222222222")
	gasLimit := uint64(30400)
	_, _, err := client.Call(
		context.Background(),
		types.Call{
			From:     nil,
			To:       &to,
			GasLimit: &gasLimit,
			GasPrice: big.NewInt(10000000000000),
			Value:    big.NewInt(10000000000),
			Input:    hexToBytes("0x3333333333333333333333333333333333333333333333333333333333333333333333333333333333"),
		},
		types.BlockNumberFromUint64(1),
	)
	require.NoError(t, err)
	assert.JSONEq(t, mockCallRequest, readBody(httpMock.Request))
}

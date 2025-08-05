package rpc

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const mockClientVersionRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "web3_clientVersion",
	  "params": []
	}
`

const mockClientVersionResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": "Geth/v1.9.25-unstable-3f0b5e4e-20201014/linux-amd64/go1.15.2"
	}
`

func TestBaseClient_ClientVersion(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsClient{Transport: httpMock}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockClientVersionRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockClientVersionResponse)),
		}, nil
	}

	clientVersion, err := client.ClientVersion(context.Background())

	require.NoError(t, err)
	assert.Equal(t, "Geth/v1.9.25-unstable-3f0b5e4e-20201014/linux-amd64/go1.15.2", clientVersion)
}

const mockNetworkIDRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "net_version",
	  "params": []
	}
`

const mockNetworkIDResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": "0x1"
	}
`

func TestBaseClient_NetworkID(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsClient{Transport: httpMock}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockNetworkIDRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockNetworkIDResponse)),
		}, nil
	}

	networkID, err := client.NetworkID(context.Background())

	require.NoError(t, err)
	assert.Equal(t, uint64(1), networkID)
}

const mockListeningRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "net_listening",
	  "params": []
	}
`

const mockListeningResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": true
	}
`

func TestBaseClient_Listening(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsClient{Transport: httpMock}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockListeningRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockListeningResponse)),
		}, nil
	}

	listening, err := client.Listening(context.Background())

	require.NoError(t, err)
	assert.True(t, listening)
}

const mockPeerCountRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "net_peerCount",
	  "params": []
	}
`

const mockPeerCountResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": "0x1"
	}
`

func TestBaseClient_PeerCount(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsClient{Transport: httpMock}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockPeerCountRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockPeerCountResponse)),
		}, nil
	}

	peerCount, err := client.PeerCount(context.Background())

	require.NoError(t, err)
	assert.Equal(t, uint64(1), peerCount)
}

const mockSyncingRequest = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "method": "eth_syncing",
	  "params": []
	}
`

const mockSyncingResponse = `
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": {
	    "startingBlock": "0x384",
	    "currentBlock": "0x386",
	    "highestBlock": "0x454"
	  }
	}
`

func TestBaseClient_Syncing(t *testing.T) {
	httpMock := newHTTPMock()
	client := &MethodsClient{Transport: httpMock}

	httpMock.Handler = func(req *http.Request) (*http.Response, error) {
		assert.JSONEq(t, mockSyncingRequest, readBody(req))
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockSyncingResponse)),
		}, nil
	}

	syncStatus, err := client.Syncing(context.Background())

	require.NoError(t, err)
	assert.Equal(t, uint64(0x384), syncStatus.StartingBlock.Big().Uint64())
	assert.Equal(t, uint64(0x386), syncStatus.CurrentBlock.Big().Uint64())
	assert.Equal(t, uint64(0x454), syncStatus.HighestBlock.Big().Uint64())
}

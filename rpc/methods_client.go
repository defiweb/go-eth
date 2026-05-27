package rpc

import (
	"context"

	"github.com/defiweb/go-eth/types"
)

// MethodsClient is a collection of RPC methods that provide information about
// the node.
//
// Note: Some JSON-RPC APIs do not support these methods.
type MethodsClient struct {
	Context *ClientContext
}

// ClientVersion performs web3_clientVersion RPC call.
//
// It returns the current client version.
func (c *MethodsClient) ClientVersion(ctx context.Context) (string, error) {
	var res string
	if err := c.Context.Transport.Call(ctx, &res, "web3_clientVersion"); err != nil {
		return "", err
	}
	return res, nil
}

// NetworkID performs net_version RPC call.
//
// It returns the current network ID.
func (c *MethodsClient) NetworkID(ctx context.Context) (uint64, error) {
	var res types.Number
	if err := c.Context.Transport.Call(ctx, &res, "net_version"); err != nil {
		return 0, err
	}
	return res.Big().Uint64(), nil
}

// Listening performs net_listening RPC call.
//
// It returns true if the client is actively listening for network.
func (c *MethodsClient) Listening(ctx context.Context) (bool, error) {
	var res bool
	if err := c.Context.Transport.Call(ctx, &res, "net_listening"); err != nil {
		return false, err
	}
	return res, nil
}

// PeerCount performs net_peerCount RPC call.
//
// It returns the number of connected peers.
func (c *MethodsClient) PeerCount(ctx context.Context) (uint64, error) {
	var res types.Number
	if err := c.Context.Transport.Call(ctx, &res, "net_peerCount"); err != nil {
		return 0, err
	}
	return res.Big().Uint64(), nil
}

// Syncing performs eth_syncing RPC call.
//
// It returns an object with data about the sync status or false.
func (c *MethodsClient) Syncing(ctx context.Context) (*types.SyncStatus, error) {
	var res types.SyncStatus
	if err := c.Context.Transport.Call(ctx, &res, "eth_syncing"); err != nil {
		return nil, err
	}
	return &res, nil
}

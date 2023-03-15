package rpc

import (
	"context"
	"fmt"

	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
	"github.com/defiweb/go-eth/wallet"
)

type Client struct {
	baseClient

	keys        []wallet.Key
	defaultAddr *types.Address
	chainID     *uint64
}

type ClientOptions func(c *Client) error

// WithTransport sets the transport for the client.
func WithTransport(transport transport.Transport) ClientOptions {
	return func(c *Client) error {
		c.transport = transport
		return nil
	}
}

// WithKeys allows to set keys that will be used to sign data.
// It allows to emulate the behavior of the RPC methods that require a key
// to be available in the node.
//
// The following methods are affected:
// - Accounts - returns the addresses of the provided keys
// - Sign - signs the data with the provided key
// - SignTransaction - signs transaction with the provided key
// - SendTransaction - signs transaction with the provided key and sends it
//   using SendRawTransaction
func WithKeys(keys ...wallet.Key) ClientOptions {
	return func(c *Client) error {
		c.keys = keys
		return nil
	}
}

// WithDefaultAddress sets the transaction.From address if it is not set.
func WithDefaultAddress(addr types.Address) ClientOptions {
	return func(c *Client) error {
		c.defaultAddr = &addr
		return nil
	}
}

// WithChainID sets the transaction.ChainID if it is not set. If the transaction
// has a ChainID set, it will return an error if it does not match the provided
// chain ID.
func WithChainID(chainID uint64) ClientOptions {
	return func(c *Client) error {
		c.chainID = &chainID
		return nil
	}
}

// NewClient creates a new RPC client.
// The WithTransport option is required.
func NewClient(opts ...ClientOptions) (*Client, error) {
	c := &Client{}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	if c.transport == nil {
		return nil, fmt.Errorf("rpc client: transport is required")
	}
	return c, nil
}

// Accounts implements the RPC interface.
func (c *Client) Accounts(ctx context.Context) ([]types.Address, error) {
	if len(c.keys) > 0 {
		var res []types.Address
		for _, key := range c.keys {
			res = append(res, key.Address())
		}
		return res, nil
	}
	return c.baseClient.Accounts(ctx)
}

// Sign implements the RPC interface.
func (c *Client) Sign(ctx context.Context, account types.Address, data []byte) (*types.Signature, error) {
	if key := c.findKey(&account); key != nil {
		return key.SignMessage(data)
	}
	return c.baseClient.Sign(ctx, account, data)
}

// SignTransaction implements the RPC interface.
func (c *Client) SignTransaction(ctx context.Context, tx types.Transaction) ([]byte, *types.Transaction, error) {
	txPtr := &tx
	c.setTXChainID(txPtr)
	if err := c.verifyTXChainID(txPtr); err != nil {
		return nil, nil, err
	}
	if key := c.findKey(tx.Call.From); key != nil {
		if err := key.SignTransaction(txPtr); err != nil {
			return nil, nil, err
		}
		raw, err := tx.Raw()
		if err != nil {
			return nil, nil, err
		}
		return raw, txPtr, nil
	}
	return c.baseClient.SignTransaction(ctx, tx)
}

// SendTransaction implements the RPC interface.
func (c *Client) SendTransaction(ctx context.Context, tx types.Transaction) (*types.Hash, error) {
	txPtr := &tx
	c.setTXChainID(txPtr)
	if err := c.verifyTXChainID(txPtr); err != nil {
		return nil, err
	}
	if key := c.findKey(tx.Call.From); key != nil {
		if err := key.SignTransaction(txPtr); err != nil {
			return nil, err
		}
		raw, err := tx.Raw()
		if err != nil {
			return nil, err
		}
		return c.SendRawTransaction(ctx, raw)
	}
	return c.baseClient.SendTransaction(ctx, tx)
}

// setTXChainID sets the transaction chain ID if it is not set and the client
// has a chain ID set.
func (c *Client) setTXChainID(tx *types.Transaction) {
	if c.chainID == nil {
		return
	}
	if tx.ChainID == nil {
		id := *c.chainID
		tx.ChainID = &id
	}
}

// verifyTXChainID verifies that the transaction chain ID is set. If the client
// has a chain ID set, it will also verify that the transaction chain ID matches
// the client chain ID.
func (c *Client) verifyTXChainID(tx *types.Transaction) error {
	if tx.ChainID == nil {
		return fmt.Errorf("transaction chain ID is not set")
	}
	if c.chainID != nil && *tx.ChainID != *c.chainID {
		return fmt.Errorf("transaction chain ID does not match")
	}
	return nil
}

// findKey finds a key by address. If the address is nil, it will return the
// default key.
func (c *Client) findKey(addr *types.Address) wallet.Key {
	if addr == nil {
		addr = c.defaultAddr
	}
	if addr != nil {
		for _, key := range c.keys {
			if key.Address() == *addr {
				return key
			}
		}
	}
	return nil
}

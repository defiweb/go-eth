package rpc

import (
	"context"
	"fmt"

	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
	"github.com/defiweb/go-eth/wallet"
)

// Client allows to interact with the Ethereum node.
type Client struct {
	baseClient

	keys        []wallet.Key
	defaultAddr *types.Address
	chainID     *uint64
	txModifiers []TXModifier
}

type ClientOptions func(c *Client) error

// TXModifier allows to modify the transaction before it is signed or sent to
// the node.
type TXModifier interface {
	Modify(ctx context.Context, client RPC, tx *types.Transaction) error
}

type TXModifierFunc func(ctx context.Context, client RPC, tx *types.Transaction) error

func (f TXModifierFunc) Modify(ctx context.Context, client RPC, tx *types.Transaction) error {
	return f(ctx, client, tx)
}

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
//   - Accounts - returns the addresses of the provided keys
//   - Sign - signs the data with the provided key
//   - SignTransaction - signs transaction with the provided key
//   - SendTransaction - signs transaction with the provided key and sends it
//     using SendRawTransaction
func WithKeys(keys ...wallet.Key) ClientOptions {
	return func(c *Client) error {
		c.keys = keys
		return nil
	}
}

// WithDefaultAddress sets the transaction.From address if it is not set
// in the following methods:
// - SignTransaction
// - SendTransaction
func WithDefaultAddress(addr types.Address) ClientOptions {
	return func(c *Client) error {
		c.defaultAddr = &addr
		return nil
	}
}

// WithChainID sets the transaction.ChainID if it is not set in the following
// methods:
// - SignTransaction
// - SendTransaction
//
// If the transaction has a ChainID set, it will return an error if it does not
// match the provided chain ID.
func WithChainID(chainID uint64) ClientOptions {
	return func(c *Client) error {
		c.chainID = &chainID
		return nil
	}
}

// WithTXModifiers allows to modify the transaction before it is signed and
// sent to the node.
//
// Modifiers will be applied in the order they are provided.
func WithTXModifiers(modifiers ...TXModifier) ClientOptions {
	return func(c *Client) error {
		c.txModifiers = append(c.txModifiers, modifiers...)
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
	if len(c.keys) == 0 {
		return c.baseClient.Sign(ctx, account, data)
	}
	if key := c.findKey(&account); key != nil {
		return key.SignMessage(data)
	}
	return nil, fmt.Errorf("rpc client: no key found for address %s", account)
}

// SignTransaction implements the RPC interface.
func (c *Client) SignTransaction(ctx context.Context, tx types.Transaction) ([]byte, *types.Transaction, error) {
	txPtr := &tx
	if tx.ChainID == nil && c.chainID != nil {
		chainID := *c.chainID
		txPtr.ChainID = &chainID
	}
	if tx.Call.From == nil && c.defaultAddr != nil {
		defaultAddr := *c.defaultAddr
		txPtr.Call.From = &defaultAddr
	}
	if err := c.verifyTXChainID(txPtr); err != nil {
		return nil, nil, err
	}
	for _, modifier := range c.txModifiers {
		if err := modifier.Modify(ctx, c, txPtr); err != nil {
			return nil, nil, err
		}
	}
	if len(c.keys) == 0 {
		return c.baseClient.SignTransaction(ctx, tx)
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
	return nil, nil, fmt.Errorf("rpc client: no key found for address %s", tx.Call.From)
}

// SendTransaction implements the RPC interface.
func (c *Client) SendTransaction(ctx context.Context, tx types.Transaction) (*types.Hash, error) {
	txPtr := &tx
	if tx.ChainID == nil && c.chainID != nil {
		chainID := *c.chainID
		txPtr.ChainID = &chainID
	}
	if tx.Call.From == nil && c.defaultAddr != nil {
		defaultAddr := *c.defaultAddr
		txPtr.Call.From = &defaultAddr
	}
	if err := c.verifyTXChainID(txPtr); err != nil {
		return nil, err
	}
	for _, modifier := range c.txModifiers {
		if err := modifier.Modify(ctx, c, txPtr); err != nil {
			return nil, err
		}
	}
	if len(c.keys) == 0 {
		return c.baseClient.SendTransaction(ctx, tx)
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
	return nil, fmt.Errorf("rpc client: no key found for address %s", tx.Call.From)
}

// Call implements the RPC interface.
func (c *Client) Call(ctx context.Context, call types.Call, block types.BlockNumber) ([]byte, error) {
	if call.From == nil && c.defaultAddr != nil {
		defaultAddr := *c.defaultAddr
		call.From = &defaultAddr
	}
	return c.baseClient.Call(ctx, call, block)
}

// verifyTXChainID verifies that the transaction chain ID is set. If the client
// has a chain ID set, it will also verify that the transaction chain ID matches
// the client chain ID.
func (c *Client) verifyTXChainID(tx *types.Transaction) error {
	if tx.ChainID == nil {
		return fmt.Errorf("rpc client: transaction chain ID is not set")
	}
	if c.chainID != nil && *tx.ChainID != *c.chainID {
		return fmt.Errorf("rpc client: transaction chain ID does not match")
	}
	return nil
}

// findKey finds a key by address.
func (c *Client) findKey(addr *types.Address) wallet.Key {
	if addr == nil {
		return nil
	}
	for _, key := range c.keys {
		if key.Address() == *addr {
			return key
		}
	}
	return nil
}

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

	keys        map[types.Address]wallet.Key
	defaultAddr *types.Address
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
// It allows to emulate the behavior of the RPC methods that require a key.
//
// The following methods are affected:
//   - Accounts - returns the addresses of the provided keys
//   - Sign - signs the data with the provided key
//   - SignTransaction - signs transaction with the provided key
//   - SendTransaction - signs transaction with the provided key and sends it
//     using SendRawTransaction
func WithKeys(keys ...wallet.Key) ClientOptions {
	return func(c *Client) error {
		for _, k := range keys {
			c.keys[k.Address()] = k
		}
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
	c := &Client{keys: make(map[types.Address]wallet.Key)}
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
func (c *Client) SignTransaction(ctx context.Context, tx *types.Transaction) ([]byte, *types.Transaction, error) {
	tx, err := c.PrepareTransaction(ctx, tx)
	if err != nil {
		return nil, nil, err
	}
	if len(c.keys) == 0 {
		return c.baseClient.SignTransaction(ctx, tx)
	}
	if key := c.findKey(tx.Call.From); key != nil {
		if err := key.SignTransaction(tx); err != nil {
			return nil, nil, err
		}
		raw, err := tx.Raw()
		if err != nil {
			return nil, nil, err
		}
		return raw, tx, nil
	}
	return nil, nil, fmt.Errorf("rpc client: no key found for address %s", tx.Call.From)
}

// SendTransaction implements the RPC interface.
func (c *Client) SendTransaction(ctx context.Context, tx *types.Transaction) (*types.Hash, *types.Transaction, error) {
	tx, err := c.PrepareTransaction(ctx, tx)
	if err != nil {
		return nil, nil, err
	}
	if len(c.keys) == 0 {
		return c.baseClient.SendTransaction(ctx, tx)
	}
	if key := c.findKey(tx.Call.From); key != nil {
		if err := key.SignTransaction(tx); err != nil {
			return nil, nil, err
		}
		raw, err := tx.Raw()
		if err != nil {
			return nil, nil, err
		}
		txHash, err := c.SendRawTransaction(ctx, raw)
		if err != nil {
			return nil, nil, err
		}
		return txHash, tx, nil
	}
	return nil, nil, fmt.Errorf("rpc client: no key found for address %s", tx.Call.From)
}

// PrepareTransaction prepares the transaction by applying transaction
// modifiers and setting the default address if it is not set.
//
// A copy of the modified transaction is returned.
func (c *Client) PrepareTransaction(ctx context.Context, tx *types.Transaction) (*types.Transaction, error) {
	if tx == nil {
		return nil, fmt.Errorf("rpc client: transaction is nil")
	}
	txCpy := tx.Copy()
	if txCpy.Call.From == nil && c.defaultAddr != nil {
		defaultAddr := *c.defaultAddr
		txCpy.Call.From = &defaultAddr
	}
	for _, modifier := range c.txModifiers {
		if err := modifier.Modify(ctx, c, txCpy); err != nil {
			return nil, err
		}
	}
	return txCpy, nil
}

// Call implements the RPC interface.
func (c *Client) Call(ctx context.Context, call *types.Call, block types.BlockNumber) ([]byte, *types.Call, error) {
	if call == nil {
		return nil, nil, fmt.Errorf("rpc client: call is nil")
	}
	callCpy := call.Copy()
	if callCpy.From == nil && c.defaultAddr != nil {
		defaultAddr := *c.defaultAddr
		callCpy.From = &defaultAddr
	}
	return c.baseClient.Call(ctx, callCpy, block)
}

// EstimateGas implements the RPC interface.
func (c *Client) EstimateGas(ctx context.Context, call *types.Call, block types.BlockNumber) (uint64, *types.Call, error) {
	if call == nil {
		return 0, nil, fmt.Errorf("rpc client: call is nil")
	}
	callCpy := call.Copy()
	if callCpy.From == nil && c.defaultAddr != nil {
		defaultAddr := *c.defaultAddr
		callCpy.From = &defaultAddr
	}
	return c.baseClient.EstimateGas(ctx, callCpy, block)
}

// findKey finds a key by address.
func (c *Client) findKey(addr *types.Address) wallet.Key {
	if addr == nil {
		return nil
	}
	if key, ok := c.keys[*addr]; ok {
		return key
	}
	return nil
}

package rpc

import (
	"fmt"
	"math/big"
	"reflect"
	"sort"

	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
	"github.com/defiweb/go-eth/wallet"
)

// Client is a default RPC client that provides access to the standard Ethereum
// JSON-RPC APIs.
type Client struct {
	MethodsCommon
	MethodsFilter
	MethodsWallet
	MethodsClient
}

type ClientOptionsContext struct {
	Transport transport.Transport      // Transport instance that will be passed to the client.
	Decoder   types.TransactionDecoder // Transaction decoder that will be passed to the client.
	Custom    map[any]any              // Custom data that may be used by client options.
}

type ClientOption interface {
	Apply(cfg *ClientOptionsContext, client any) error
	Order() int
}

// WithTransport sets the transport for the client.
func WithTransport(t transport.Transport) ClientOption {
	return &option{
		apply: func(ctx *ClientOptionsContext, _ any) error {
			ctx.Transport = t
			return nil
		},
	}
}

// WithTransactionDecoder sets the transaction decoder for the client.
// The default decoder is types.DefaultTransactionDecoder.
//
// Using custom decoder allows to decode custom transaction types that may be
// present in some L2 implementations.
func WithTransactionDecoder(decoder types.TransactionDecoder) ClientOption {
	return &option{
		apply: func(ctx *ClientOptionsContext, _ any) error {
			ctx.Decoder = decoder
			return nil
		},
	}
}

// WithPostHijackers adds hijackers that are applied after all other hijackers
// applied by the client.
func WithPostHijackers(hijackers ...transport.Hijacker) ClientOption {
	return &option{
		apply: func(ctx *ClientOptionsContext, _ any) error {
			ctx.Transport = addHijacker(ctx.Transport, hijackers...)
			return nil
		},
		order: 100,
	}
}

// WithKeys allows to emulate the behavior of the RPC methods that require
// a private key to sign the data.
//
// The following methods are affected:
//   - Accounts - returns the addresses of the provided keys
//   - Sign - signs the data with the provided key
//   - SignTransaction - signs transaction
//   - SendTransaction - signs transaction and sends it using SendRawTransaction
//
// This option will modify the provided transaction instance.
func WithKeys(keys ...wallet.Key) ClientOption {
	return &option{
		apply: func(ctx *ClientOptionsContext, _ any) error {
			ctx.Transport = addHijacker(ctx.Transport, &hijackSign{keys: keys})
			return nil
		},
		order: 200,
	}
}

// WithSimulate simulates the transaction, by calling eth_call with the same
// parameters before sending the transaction.
//
// It works with eth_sendTransaction, eth_sendRawTransaction, and
// eth_sendPrivateTransaction methods.
func WithSimulate() ClientOption {
	return &option{
		apply: func(ctx *ClientOptionsContext, _ any) error {
			ctx.Transport = addHijacker(ctx.Transport, &hijackSimulate{})
			return nil
		},
		order: 300,
	}
}

type NonceOptions struct {
	UsePendingBlock bool // UsePendingBlock indicates whether to use the pending block.
	Replace         bool // Replace is true if the nonce should be replaced even if it is already set.
}

// WithNonce sets the nonce in the transaction.
//
// This option will modify the provided transaction instance.
func WithNonce(opts NonceOptions) ClientOption {
	return &option{
		apply: func(ctx *ClientOptionsContext, _ any) error {
			ctx.Transport = addHijacker(ctx.Transport, &hijackNonce{
				usePendingBlock: opts.UsePendingBlock,
				replace:         opts.Replace,
			})
			return nil
		},
		order: 400,
	}
}

type LegacyGasFeeOptions struct {
	Multiplier      float64  // Multiplier is applied to the gas price.
	MinGasPrice     *big.Int // MinGasPrice is the minimum gas price, or nil if there is no lower bound.
	MaxGasPrice     *big.Int // MaxGasPrice is the maximum gas price, or nil if there is no upper bound.
	Replace         bool     // Replace is true if the gas price should be replaced even if it is already set.
	AllowChangeType bool     // AllowChangeType is true if the transaction type can be changed if it does not support legacy price data.
}

// WithLegacyGasFee estimates the gas price and sets it in the transaction.
//
// It only works with eth_sendTransaction method, raw transactions are not supported.
//
// This option will modify the provided transaction instance.
func WithLegacyGasFee(opts LegacyGasFeeOptions) ClientOption {
	return &option{
		apply: func(ctx *ClientOptionsContext, _ any) error {
			if opts.Multiplier == 0 {
				return fmt.Errorf("rpc client: gas price multiplier must be greater than 0")
			}
			ctx.Transport = addHijacker(ctx.Transport, &hijackLegacyGasFee{
				multiplier:      opts.Multiplier,
				minGasPrice:     opts.MinGasPrice,
				maxGasPrice:     opts.MaxGasPrice,
				replace:         opts.Replace,
				allowChangeType: opts.AllowChangeType,
			})
			return nil
		},
		order: 500,
	}
}

type DynamicGasFeeOptions struct {
	GasPriceMultiplier          float64  // GasPriceMultiplier is applied to the gas price.
	PriorityFeePerGasMultiplier float64  // PriorityFeePerGasMultiplier is applied to the priority fee per gas.
	MinGasPrice                 *big.Int // MinGasPrice is the minimum gas price, or nil if there is no lower bound.
	MaxGasPrice                 *big.Int // MaxGasPrice is the maximum gas price, or nil if there is no upper bound.
	MinPriorityFeePerGas        *big.Int // MinPriorityFeePerGas is the minimum priority fee per gas, or nil if there is no lower bound.
	MaxPriorityFeePerGas        *big.Int // MaxPriorityFeePerGas is the maximum priority fee per gas, or nil if there is no upper bound.
	Replace                     bool     // Replace is true if the gas price should be replaced even if it is already set.
	AllowChangeType             bool     // AllowChangeType is true if the transaction type can be changed if it does not support dynamic price data.
}

// WithDynamicGasFee estimates the gas price and sets it in the transaction.
//
// It only works with eth_sendTransaction method, raw transactions are not supported.
//
// This option will modify the provided transaction instance.
func WithDynamicGasFee(opts DynamicGasFeeOptions) ClientOption {
	return &option{
		apply: func(ctx *ClientOptionsContext, _ any) error {
			if opts.GasPriceMultiplier == 0 || opts.PriorityFeePerGasMultiplier == 0 {
				return fmt.Errorf("rpc client: gas price and priority fee multipliers must be greater than 0")
			}
			ctx.Transport = addHijacker(ctx.Transport, &hijackDynamicGasFee{
				gasPriceMultiplier:          opts.GasPriceMultiplier,
				priorityFeePerGasMultiplier: opts.PriorityFeePerGasMultiplier,
				minGasPrice:                 opts.MinGasPrice,
				maxGasPrice:                 opts.MaxGasPrice,
				minPriorityFeePerGas:        opts.MinPriorityFeePerGas,
				maxPriorityFeePerGas:        opts.MaxPriorityFeePerGas,
				replace:                     opts.Replace,
				allowChangeType:             opts.AllowChangeType,
			})
			return nil
		},
		order: 600,
	}
}

type GasLimitOptions struct {
	Multiplier float64 // Multiplier is applied to the gas limit.
	MinGas     uint64  // MinGas is the minimum gas limit, or 0 if there is no lower bound.
	MaxGas     uint64  // MaxGas is the maximum gas limit, or 0 if there is no upper bound.
	Replace    bool    // Replace is true if the gas limit should be replaced even if it is already set.
}

// WithGasLimit estimates the gas limit and sets it in the transaction.
//
// This option will modify the provided transaction instance.
func WithGasLimit(opts GasLimitOptions) ClientOption {
	return &option{
		apply: func(ctx *ClientOptionsContext, _ any) error {
			if opts.Multiplier == 0 {
				return fmt.Errorf("rpc client: gas limit multiplier must be greater than 0")
			}
			ctx.Transport = addHijacker(ctx.Transport, &hijackGasLimit{
				multiplier: opts.Multiplier,
				minGas:     opts.MinGas,
				maxGas:     opts.MaxGas,
				replace:    opts.Replace,
			})
			return nil
		},
		order: 700,
	}
}

type AddressOptions struct {
	Address types.Address // Address is the default address to use.
	Replace bool          // Replace is true if the address should be replaced even if it is already set.
}

// WithDefaultAddress sets the default address for calls and transactions.
//
// To send a call with to a zero address, it must be set explicitly in the call.
//
// This option will modify the provided transaction and call instances.
func WithDefaultAddress(opts AddressOptions) ClientOption {
	return &option{
		apply: func(ctx *ClientOptionsContext, _ any) error {
			ctx.Transport = addHijacker(ctx.Transport, &hijackAddress{
				address: opts.Address,
				replace: opts.Replace,
			})
			return nil
		},
		order: 800,
	}
}

type ChainIDOptions struct {
	ChainID uint64 // ChainID is the chain ID to use. If 0, then value is fetched from the node.
	Replace bool   // Replace is true if the chain ID should be replaced even if it is already set.
}

// WithChainID sets the chain ID in the transaction.
// It only works with eth_sendTransaction method.
//
// This option will modify the provided transaction instance.
func WithChainID(opts ChainIDOptions) ClientOption {
	return &option{
		apply: func(ctx *ClientOptionsContext, _ any) error {
			ctx.Transport = addHijacker(ctx.Transport, &hijackChainID{
				chainID: opts.ChainID,
				replace: opts.Replace,
			})
			return nil
		},
		order: 900,
	}
}

// WithPreHijackers adds hijackers that are applied before any other hijackers
// applied by the client.
func WithPreHijackers(hijackers ...transport.Hijacker) ClientOption {
	return &option{
		apply: func(ctx *ClientOptionsContext, _ any) error {
			ctx.Transport = addHijacker(ctx.Transport, hijackers...)
			return nil
		},
		order: 1000,
	}
}

type ClientOptions func(c *ClientOptionsContext, client any) error

// NewClient creates a new RPC client.
//
// The WithTransport option is required.
func NewClient(opts ...ClientOption) (*Client, error) {
	ctx := &ClientOptionsContext{}
	client := &Client{}
	if err := applyOptions(ctx, opts); err != nil {
		return nil, fmt.Errorf("rpc client: option error: %w", err)
	}
	if ctx.Decoder == nil {
		ctx.Decoder = types.DefaultTransactionDecoder
	}
	if ctx.Transport == nil {
		return nil, fmt.Errorf("rpc client: transport is required")
	}
	client.MethodsCommon.Transport = ctx.Transport
	client.MethodsFilter.Transport = ctx.Transport
	client.MethodsWallet.Transport = ctx.Transport
	client.MethodsClient.Transport = ctx.Transport
	client.MethodsCommon.Decoder = ctx.Decoder
	return client, nil
}

// NewCustomClient returns a new custom client. A custom client may implement
// additional methods that are not part of the standard client.
//
// The WithTransport option is required.
//
// This method automatically initializes the client fields that are nil
// and recursively sets the all fields that are of the transport.Transport or
// types.TransactionDecoder types.
func NewCustomClient[T any](opts ...ClientOption) (*T, error) {
	ctx := &ClientOptionsContext{}
	client := new(T)
	if err := applyOptions(ctx, opts); err != nil {
		return nil, fmt.Errorf("rpc client: option error: %w", err)
	}
	if ctx.Decoder == nil {
		ctx.Decoder = types.DefaultTransactionDecoder
	}
	if ctx.Transport == nil {
		return nil, fmt.Errorf("rpc client: transport is required")
	}
	setFields(ctx, reflect.ValueOf(client))
	return client, nil
}

type option struct {
	apply func(*ClientOptionsContext, any) error
	order int
}

func (o *option) Apply(ctx *ClientOptionsContext, client any) error {
	return o.apply(ctx, client)
}

func (o *option) Order() int {
	return o.order
}

func applyOptions(c *ClientOptionsContext, opts []ClientOption) error {
	sort.Slice(opts, func(i, j int) bool {
		return opts[i].Order() < opts[j].Order()
	})
	for _, opt := range opts {
		if err := opt.Apply(c, nil); err != nil {
			return err
		}
	}
	return nil
}

func setFields(ctx *ClientOptionsContext, r reflect.Value) {
	for r.Kind() == reflect.Ptr || r.Kind() == reflect.Interface {
		r = r.Elem()
	}
	if r.Kind() != reflect.Struct {
		return
	}
	for n := 0; n < r.NumField(); n++ {
		f := r.Field(n)
		if !f.CanInterface() {
			continue
		}
		t := f.Type()
		switch t {
		case transportTy:
			if f.CanSet() {
				f.Set(reflect.ValueOf(ctx.Transport))
			}
		case decoderTy:
			if f.CanSet() {
				f.Set(reflect.ValueOf(ctx.Decoder))
			}
		default:
			if initPtr(f) {
				setFields(ctx, f)
			}
		}
	}
}

func initPtr(r reflect.Value) bool {
	if !r.CanInterface() {
		return false
	}
	if r.Kind() != reflect.Ptr {
		return true
	}
	if !r.IsNil() {
		return true
	}
	if r.CanSet() {
		r.Set(reflect.New(r.Type().Elem()))
		return true
	}
	return false
}

func addHijacker(t transport.Transport, hijackers ...transport.Hijacker) transport.Transport {
	if h, ok := t.(*optionTransportHijack); ok {
		h.Use(hijackers...)
		return h
	}
	return optionTransportHijack{transport.NewHijacker(t, hijackers...)}
}

// optionTransportHijack wraps a Hijack transport to ensure that hijackers
// added by the user won't be modified by the client.
type optionTransportHijack struct {
	*transport.Hijack
}

var (
	transportTy = reflect.TypeOf((*transport.Transport)(nil)).Elem()
	decoderTy   = reflect.TypeOf((*types.TransactionDecoder)(nil)).Elem()
)

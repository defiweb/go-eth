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

const (
	OrderPost int = 100 * iota
	OrderKeys
	OrderSimulate
	OrderNonce
	OrderLegacyGasFee
	OrderDynamicGasFee
	OrderGasLimit
	DefaultAddress
	OrderChainID
	OrderPre
)

// Client is a default RPC client that provides access to the standard Ethereum
// JSON-RPC APIs.
type Client struct {
	MethodsCommon
	MethodsFilter
	MethodsWallet
	MethodsClient
}

type ClientContext struct {
	// Transport is the transport for the client.
	Transport transport.Transport

	// Decoder is the transaction decoder for the client.
	Decoder types.TransactionDecoder
}

type ClientOption interface {
	Apply(cfg *ClientContext, client any) error
	Order() int
}

// WithTransport sets the transport for the client.
func WithTransport(t transport.Transport) ClientOption {
	return &option{
		apply: func(ctx *ClientContext, _ any) error {
			ctx.Transport = t
			return nil
		},
	}
}

// WithTransactionDecoder sets the transaction decoder for the client.
// The default decoder is types.DefaultTransactionDecoder.
//
// Using custom decoder allows decoding custom transaction types that may be
// present in some L2 implementations.
func WithTransactionDecoder(decoder types.TransactionDecoder) ClientOption {
	return &option{
		apply: func(ctx *ClientContext, _ any) error {
			ctx.Decoder = decoder
			return nil
		},
	}
}

// WithPostHijackers adds hijackers that are applied after all other hijackers
// applied by the client.
func WithPostHijackers(hijackers ...transport.Hijacker) ClientOption {
	return &option{
		apply: func(ctx *ClientContext, _ any) error {
			ctx.Transport = addHijacker(ctx.Transport, hijackers...)
			return nil
		},
		order: OrderPost,
	}
}

// WithKeys allows emulating the behavior of the RPC methods that require
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
		apply: func(ctx *ClientContext, _ any) error {
			ctx.Transport = addHijacker(ctx.Transport, &hijackSign{keys: keys})
			return nil
		},
		order: OrderKeys,
	}
}

// WithSimulate simulates the transaction by calling eth_call with the same
// parameters before sending the transaction.
//
// It works with eth_sendTransaction, eth_sendRawTransaction, and
// eth_sendPrivateTransaction methods.
func WithSimulate() ClientOption {
	return &option{
		apply: func(ctx *ClientContext, _ any) error {
			ctx.Transport = addHijacker(ctx.Transport, &hijackSimulate{decoder: ctx.Decoder})
			return nil
		},
		order: OrderSimulate,
	}
}

type NonceOptions struct {
	// UsePendingBlock queries the nonce from the pending block.
	UsePendingBlock bool

	// Replace overwrites the nonce even if already set.
	Replace bool
}

// WithNonce sets the nonce in the transaction.
//
// This option will modify the provided transaction instance.
func WithNonce(opts NonceOptions) ClientOption {
	return &option{
		apply: func(ctx *ClientContext, _ any) error {
			ctx.Transport = addHijacker(ctx.Transport, &hijackNonce{
				usePendingBlock: opts.UsePendingBlock,
				replace:         opts.Replace,
			})
			return nil
		},
		order: OrderNonce,
	}
}

type LegacyGasFeeOptions struct {
	// Multiplier is applied to the fetched gas price.
	Multiplier float64

	// MinGasPrice is the lower bound; nil means no lower bound.
	MinGasPrice *big.Int

	// MaxGasPrice is the upper bound; nil means no upper bound.
	MaxGasPrice *big.Int

	// Replace overwrites the gas price even if already set.
	Replace bool

	// AllowChangeType converts the transaction type when it does not
	// support legacy fee data.
	AllowChangeType bool
}

// WithLegacyGasFee estimates the gas price and sets it in the transaction.
//
// It only works with eth_sendTransaction; raw transactions are not supported.
//
// This option will modify the provided transaction instance.
func WithLegacyGasFee(opts LegacyGasFeeOptions) ClientOption {
	return &option{
		apply: func(ctx *ClientContext, _ any) error {
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
		order: OrderLegacyGasFee,
	}
}

type DynamicGasFeeOptions struct {
	// GasPriceMultiplier is applied to the base fee.
	GasPriceMultiplier float64

	// PriorityFeePerGasMultiplier is applied to the priority fee.
	PriorityFeePerGasMultiplier float64

	// MinGasPrice is the lower bound on maxFeePerGas; nil means none.
	MinGasPrice *big.Int

	// MaxGasPrice is the upper bound on maxFeePerGas; nil means none.
	MaxGasPrice *big.Int

	// MinPriorityFeePerGas is the lower bound; nil means no lower bound.
	MinPriorityFeePerGas *big.Int

	// MaxPriorityFeePerGas is the upper bound; nil means no upper bound.
	MaxPriorityFeePerGas *big.Int

	// Replace overwrites the fee fields even if already set.
	Replace bool

	// AllowChangeType converts the transaction type when it does not
	// support dynamic fee data.
	AllowChangeType bool
}

// WithDynamicGasFee estimates the gas price and sets it in the transaction.
//
// It only works with eth_sendTransaction; raw transactions are not supported.
//
// This option will modify the provided transaction instance.
func WithDynamicGasFee(opts DynamicGasFeeOptions) ClientOption {
	return &option{
		apply: func(ctx *ClientContext, _ any) error {
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
		order: OrderDynamicGasFee,
	}
}

type GasLimitOptions struct {
	// Multiplier is applied to the estimated gas limit.
	Multiplier float64

	// MinGas is the lower bound on the gas limit. 0 means no lower bound.
	MinGas uint64

	// MaxGas is the upper bound on the gas limit. 0 means no upper bound.
	MaxGas uint64

	// Replace overwrites the gas limit even if already set.
	Replace bool
}

// WithGasLimit estimates the gas limit and sets it in the transaction.
//
// This option will modify the provided transaction instance.
func WithGasLimit(opts GasLimitOptions) ClientOption {
	return &option{
		apply: func(ctx *ClientContext, _ any) error {
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
		order: OrderGasLimit,
	}
}

type AddressOptions struct {
	// Address is the default sender address.
	Address types.Address

	// Replace overwrites the address even if already set.
	Replace bool
}

// WithDefaultAddress sets the default address for calls and transactions.
//
// To send a call with to a zero address, it must be set explicitly in the call.
//
// This option will modify the provided transaction and call instances.
func WithDefaultAddress(opts AddressOptions) ClientOption {
	return &option{
		apply: func(ctx *ClientContext, _ any) error {
			ctx.Transport = addHijacker(ctx.Transport, &hijackAddress{
				address: opts.Address,
				replace: opts.Replace,
			})
			return nil
		},
		order: DefaultAddress,
	}
}

type ChainIDOptions struct {
	// ChainID to set. If 0, the value is fetched from the node.
	ChainID uint64

	// Replace overwrites the chain ID even if already set.
	Replace bool
}

// WithChainID sets the chain ID in the transaction.
// It only works with eth_sendTransaction method.
//
// This option will modify the provided transaction instance.
func WithChainID(opts ChainIDOptions) ClientOption {
	return &option{
		apply: func(ctx *ClientContext, _ any) error {
			ctx.Transport = addHijacker(ctx.Transport, &hijackChainID{
				chainID: opts.ChainID,
				replace: opts.Replace,
			})
			return nil
		},
		order: OrderChainID,
	}
}

// WithPreHijackers adds hijackers that are applied before any other hijackers
// applied by the client.
func WithPreHijackers(hijackers ...transport.Hijacker) ClientOption {
	return &option{
		apply: func(ctx *ClientContext, _ any) error {
			ctx.Transport = addHijacker(ctx.Transport, hijackers...)
			return nil
		},
		order: OrderPre,
	}
}

type ClientOptions func(c *ClientContext, client any) error

// NewClient creates a new RPC client.
//
// The WithTransport option is required.
func NewClient(opts ...ClientOption) (*Client, error) {
	ctx := &ClientContext{}
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
	client.MethodsCommon.Context = ctx
	client.MethodsFilter.Context = ctx
	client.MethodsWallet.Context = ctx
	client.MethodsClient.Context = ctx
	return client, nil
}

// NewCustomClient returns a new custom client. A custom client may implement
// additional methods that are not part of the standard client.
//
// The WithTransport option is required.
//
// This method automatically initializes nil client fields and recursively
// sets the [ClientContext] in all nested structs.
func NewCustomClient[T any](opts ...ClientOption) (*T, error) {
	ctx := &ClientContext{}
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
	apply func(*ClientContext, any) error
	order int
}

func (o *option) Apply(ctx *ClientContext, client any) error {
	return o.apply(ctx, client)
}

func (o *option) Order() int {
	return o.order
}

func applyOptions(c *ClientContext, opts []ClientOption) error {
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

func setFields(ctx *ClientContext, r reflect.Value) {
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
		if f.Type() == contextTy && f.CanSet() {
			f.Set(reflect.ValueOf(ctx))
			continue
		}
		if initPtr(f) {
			setFields(ctx, f)
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

// optionTransportHijack wraps a Hijack transport to distinguish user-provided
// hijackers from those used internally by the client.
type optionTransportHijack struct {
	*transport.Hijack
}

var contextTy = reflect.TypeOf((*ClientContext)(nil)).Elem()

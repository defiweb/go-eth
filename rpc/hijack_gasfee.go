package rpc

import (
	"context"
	"fmt"
	"math/big"

	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
)

type (
	legacyGasFeeMultiplierKey  struct{}
	legacyGasFeeMinGasPriceKey struct{}
	legacyGasFeeMaxGasPriceKey struct{}
	legacyGasFeeReplaceKey     struct{}
)

// ContextWithLegacyGasFeeMultiplier overrides the Multiplier option for this
// call.
// Only has effect when the [WithLegacyGasFee] client option is enabled.
func ContextWithLegacyGasFeeMultiplier(ctx context.Context, v float64) context.Context {
	return context.WithValue(ctx, legacyGasFeeMultiplierKey{}, v)
}

// ContextWithLegacyGasFeeMinGasPrice overrides the MinGasPrice option for
// this call.
// Only has effect when the [WithLegacyGasFee] client option is enabled.
func ContextWithLegacyGasFeeMinGasPrice(ctx context.Context, v *big.Int) context.Context {
	return context.WithValue(ctx, legacyGasFeeMinGasPriceKey{}, v)
}

// ContextWithLegacyGasFeeMaxGasPrice overrides the MaxGasPrice option for
// this call.
// Only has effect when the [WithLegacyGasFee] client option is enabled.
func ContextWithLegacyGasFeeMaxGasPrice(ctx context.Context, v *big.Int) context.Context {
	return context.WithValue(ctx, legacyGasFeeMaxGasPriceKey{}, v)
}

// ContextWithLegacyGasFeeReplace overrides the Replace option for this call.
// Only has effect when the [WithLegacyGasFee] client option is enabled.
func ContextWithLegacyGasFeeReplace(ctx context.Context, v bool) context.Context {
	return context.WithValue(ctx, legacyGasFeeReplaceKey{}, v)
}

func legacyGasFeeMultiplier(ctx context.Context, h *hijackLegacyGasFee) float64 {
	if v, ok := ctx.Value(legacyGasFeeMultiplierKey{}).(float64); ok {
		return v
	}
	return h.multiplier
}

func legacyGasFeeMinGasPrice(ctx context.Context, h *hijackLegacyGasFee) *big.Int {
	if v, ok := ctx.Value(legacyGasFeeMinGasPriceKey{}).(*big.Int); ok {
		return v
	}
	return h.minGasPrice
}

func legacyGasFeeMaxGasPrice(ctx context.Context, h *hijackLegacyGasFee) *big.Int {
	if v, ok := ctx.Value(legacyGasFeeMaxGasPriceKey{}).(*big.Int); ok {
		return v
	}
	return h.maxGasPrice
}

func legacyGasFeeReplace(ctx context.Context, h *hijackLegacyGasFee) bool {
	if v, ok := ctx.Value(legacyGasFeeReplaceKey{}).(bool); ok {
		return v
	}
	return h.replace
}

type (
	dynamicGasFeeGasPriceMultiplierKey          struct{}
	dynamicGasFeePriorityFeePerGasMultiplierKey struct{}
	dynamicGasFeeMinGasPriceKey                 struct{}
	dynamicGasFeeMaxGasPriceKey                 struct{}
	dynamicGasFeeMinPriorityFeePerGasKey        struct{}
	dynamicGasFeeMaxPriorityFeePerGasKey        struct{}
	dynamicGasFeeReplaceKey                     struct{}
)

// ContextWithDynamicGasFeeGasPriceMultiplier overrides the GasPriceMultiplier
// option for this call. Only has effect when the [WithDynamicGasFee] client
// option is enabled.
func ContextWithDynamicGasFeeGasPriceMultiplier(ctx context.Context, v float64) context.Context {
	return context.WithValue(ctx, dynamicGasFeeGasPriceMultiplierKey{}, v)
}

// ContextWithDynamicGasFeePriorityFeePerGasMultiplier overrides the
// PriorityFeePerGasMultiplier option for this call. Only has effect when the
// [WithDynamicGasFee] client option is enabled.
func ContextWithDynamicGasFeePriorityFeePerGasMultiplier(ctx context.Context, v float64) context.Context {
	return context.WithValue(ctx, dynamicGasFeePriorityFeePerGasMultiplierKey{}, v)
}

// ContextWithDynamicGasFeeMinGasPrice overrides the MinGasPrice option for
// this call. Only has effect when the [WithDynamicGasFee] client option is enabled.
func ContextWithDynamicGasFeeMinGasPrice(ctx context.Context, v *big.Int) context.Context {
	return context.WithValue(ctx, dynamicGasFeeMinGasPriceKey{}, v)
}

// ContextWithDynamicGasFeeMaxGasPrice overrides the MaxGasPrice option for
// this call. Only has effect when the [WithDynamicGasFee] client option is enabled.
func ContextWithDynamicGasFeeMaxGasPrice(ctx context.Context, v *big.Int) context.Context {
	return context.WithValue(ctx, dynamicGasFeeMaxGasPriceKey{}, v)
}

// ContextWithDynamicGasFeeMinPriorityFeePerGas overrides the
// MinPriorityFeePerGas option for this call. Only has effect when the
// [WithDynamicGasFee] client option is enabled.
func ContextWithDynamicGasFeeMinPriorityFeePerGas(ctx context.Context, v *big.Int) context.Context {
	return context.WithValue(ctx, dynamicGasFeeMinPriorityFeePerGasKey{}, v)
}

// ContextWithDynamicGasFeeMaxPriorityFeePerGas overrides the
// MaxPriorityFeePerGas option for this call. Only has effect when the
// [WithDynamicGasFee] client option is enabled.
func ContextWithDynamicGasFeeMaxPriorityFeePerGas(ctx context.Context, v *big.Int) context.Context {
	return context.WithValue(ctx, dynamicGasFeeMaxPriorityFeePerGasKey{}, v)
}

// ContextWithDynamicGasFeeReplace overrides the Replace option for this call.
// Only has effect when the [WithDynamicGasFee] client option is enabled.
func ContextWithDynamicGasFeeReplace(ctx context.Context, v bool) context.Context {
	return context.WithValue(ctx, dynamicGasFeeReplaceKey{}, v)
}

func dynamicGasFeeGasPriceMultiplier(ctx context.Context, h *hijackDynamicGasFee) float64 {
	if v, ok := ctx.Value(dynamicGasFeeGasPriceMultiplierKey{}).(float64); ok {
		return v
	}
	return h.gasPriceMultiplier
}

func dynamicGasFeePriorityFeePerGasMultiplier(ctx context.Context, h *hijackDynamicGasFee) float64 {
	if v, ok := ctx.Value(dynamicGasFeePriorityFeePerGasMultiplierKey{}).(float64); ok {
		return v
	}
	return h.priorityFeePerGasMultiplier
}

func dynamicGasFeeMinGasPrice(ctx context.Context, h *hijackDynamicGasFee) *big.Int {
	if v, ok := ctx.Value(dynamicGasFeeMinGasPriceKey{}).(*big.Int); ok {
		return v
	}
	return h.minGasPrice
}

func dynamicGasFeeMaxGasPrice(ctx context.Context, h *hijackDynamicGasFee) *big.Int {
	if v, ok := ctx.Value(dynamicGasFeeMaxGasPriceKey{}).(*big.Int); ok {
		return v
	}
	return h.maxGasPrice
}

func dynamicGasFeeMinPriorityFeePerGas(ctx context.Context, h *hijackDynamicGasFee) *big.Int {
	if v, ok := ctx.Value(dynamicGasFeeMinPriorityFeePerGasKey{}).(*big.Int); ok {
		return v
	}
	return h.minPriorityFeePerGas
}

func dynamicGasFeeMaxPriorityFeePerGas(ctx context.Context, h *hijackDynamicGasFee) *big.Int {
	if v, ok := ctx.Value(dynamicGasFeeMaxPriorityFeePerGasKey{}).(*big.Int); ok {
		return v
	}
	return h.maxPriorityFeePerGas
}

func dynamicGasFeeReplace(ctx context.Context, h *hijackDynamicGasFee) bool {
	if v, ok := ctx.Value(dynamicGasFeeReplaceKey{}).(bool); ok {
		return v
	}
	return h.replace
}

// hijackLegacyGasFee hijacks the "eth_sendTransaction" method and sets the
// "gasPrice" field using the estimate provided by the RPC node.
type hijackLegacyGasFee struct {
	multiplier      float64
	minGasPrice     *big.Int
	maxGasPrice     *big.Int
	replace         bool
	allowChangeType bool
}

// Call implements the [transport.Hijacker] interface.
func (h *hijackLegacyGasFee) Call() func(next transport.CallFunc) transport.CallFunc {
	return func(next transport.CallFunc) transport.CallFunc {
		return func(ctx context.Context, t transport.Transport, result any, method string, args ...any) (err error) {
			if len(args) == 0 || method != "eth_sendTransaction" {
				return next(ctx, t, result, method, args...)
			}
			tx, ok := args[0].(types.Transaction)
			if !ok {
				return next(ctx, t, result, method, args...)
			}
			if h.allowChangeType {
				tx = convertTXToLegacyPrice(tx)
				args[0] = tx
			}
			lpd := types.GetLegacyFeeData(tx)
			if lpd != nil && (legacyGasFeeReplace(ctx, h) || lpd.GasPrice == nil) {
				gasPrice, err := (&MethodsCommon{&ClientContext{Transport: t}}).GasPrice(ctx)
				if err != nil {
					return &ErrHijackFailed{name: "legacy gas price", err: fmt.Errorf("failed to get gas price: %w", err)}
				}
				gasPrice, _ = new(big.Float).Mul(new(big.Float).SetInt(gasPrice), big.NewFloat(legacyGasFeeMultiplier(ctx, h))).Int(nil)
				if min := legacyGasFeeMinGasPrice(ctx, h); min != nil && gasPrice.Cmp(min) < 0 {
					gasPrice = min
				}
				if max := legacyGasFeeMaxGasPrice(ctx, h); max != nil && gasPrice.Cmp(max) > 0 {
					gasPrice = max
				}
				lpd.GasPrice = gasPrice
			}
			return next(ctx, t, result, method, args...)
		}
	}
}

// Subscribe implements the [transport.Hijacker] interface.
func (h *hijackLegacyGasFee) Subscribe() func(next transport.SubscribeFunc) transport.SubscribeFunc {
	return nil
}

// Unsubscribe implements the [transport.Hijacker] interface.
func (h *hijackLegacyGasFee) Unsubscribe() func(next transport.UnsubscribeFunc) transport.UnsubscribeFunc {
	return nil
}

// hijackDynamicGasFee hijacks the "eth_sendTransaction" method and sets the
// "maxFeePerGas" and "maxPriorityFeePerGas" fields using the estimates
// provided by the RPC node.
type hijackDynamicGasFee struct {
	gasPriceMultiplier          float64
	priorityFeePerGasMultiplier float64
	minGasPrice                 *big.Int
	maxGasPrice                 *big.Int
	minPriorityFeePerGas        *big.Int
	maxPriorityFeePerGas        *big.Int
	replace                     bool
	allowChangeType             bool
}

// Call implements the [transport.Hijacker] interface.
func (h *hijackDynamicGasFee) Call() func(next transport.CallFunc) transport.CallFunc {
	return func(next transport.CallFunc) transport.CallFunc {
		return func(ctx context.Context, t transport.Transport, result any, method string, args ...any) (err error) {
			if len(args) == 0 || method != "eth_sendTransaction" {
				return next(ctx, t, result, method, args...)
			}

			// Obtain the dynamic fee data.
			tx, ok := args[0].(types.Transaction)
			if !ok {
				return next(ctx, t, result, method, args...)
			}
			if h.allowChangeType {
				tx = convertTXToDynamicFee(tx)
				args[0] = tx
			}
			dfd := types.GetDynamicFeeData(tx)

			// Set the dynamic fee fields if necessary.
			if dfd != nil && (dynamicGasFeeReplace(ctx, h) || dfd.MaxFeePerGas == nil || dfd.MaxPriorityFeePerGas == nil) {
				// Fetch current gas prices from the RPC node.
				maxFeePerGas, err := (&MethodsCommon{&ClientContext{Transport: t}}).GasPrice(ctx)
				if err != nil {
					return &ErrHijackFailed{name: "dynamic gas fee", err: fmt.Errorf("failed to get gas price: %w", err)}
				}
				priorityFeePerGas, err := (&MethodsCommon{&ClientContext{Transport: t}}).MaxPriorityFeePerGas(ctx)
				if err != nil {
					return &ErrHijackFailed{name: "dynamic gas fee", err: fmt.Errorf("failed to get priority fee per gas: %w", err)}
				}

				// Apply multipliers and limits and set the fields.
				maxFeePerGas, _ = new(big.Float).Mul(new(big.Float).SetInt(maxFeePerGas), big.NewFloat(dynamicGasFeeGasPriceMultiplier(ctx, h))).Int(nil)
				priorityFeePerGas, _ = new(big.Float).Mul(new(big.Float).SetInt(priorityFeePerGas), big.NewFloat(dynamicGasFeePriorityFeePerGasMultiplier(ctx, h))).Int(nil)
				if min := dynamicGasFeeMinGasPrice(ctx, h); min != nil && maxFeePerGas.Cmp(min) < 0 {
					maxFeePerGas = min
				}
				if max := dynamicGasFeeMaxGasPrice(ctx, h); max != nil && maxFeePerGas.Cmp(max) > 0 {
					maxFeePerGas = max
				}
				if min := dynamicGasFeeMinPriorityFeePerGas(ctx, h); min != nil && priorityFeePerGas.Cmp(min) < 0 {
					priorityFeePerGas = min
				}
				if max := dynamicGasFeeMaxPriorityFeePerGas(ctx, h); max != nil && priorityFeePerGas.Cmp(max) > 0 {
					priorityFeePerGas = max
				}
				if maxFeePerGas.Cmp(priorityFeePerGas) < 0 {
					priorityFeePerGas = maxFeePerGas
				}
				dfd.MaxFeePerGas = maxFeePerGas
				dfd.MaxPriorityFeePerGas = priorityFeePerGas
			}

			return next(ctx, t, result, method, args...)
		}
	}
}

// Subscribe implements the [transport.Hijacker] interface.
func (h *hijackDynamicGasFee) Subscribe() func(next transport.SubscribeFunc) transport.SubscribeFunc {
	return nil
}

// Unsubscribe implements the [transport.Hijacker] interface.
func (h *hijackDynamicGasFee) Unsubscribe() func(next transport.UnsubscribeFunc) transport.UnsubscribeFunc {
	return nil
}

package rpc

import (
	"context"
	"fmt"
	"math/big"

	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
)

type (
	gasLimitMultiplierKey struct{}
	gasLimitMinGasKey     struct{}
	gasLimitMaxGasKey     struct{}
	gasLimitReplaceKey    struct{}
)

// ContextWithGasLimitMultiplier overrides the Multiplier option for this call.
// Only has effect when the [WithGasLimit] client option is enabled.
func ContextWithGasLimitMultiplier(ctx context.Context, v float64) context.Context {
	return context.WithValue(ctx, gasLimitMultiplierKey{}, v)
}

// ContextWithGasLimitMinGas overrides the MinGas option for this call.
// Only has effect when the [WithGasLimit] client option is enabled.
func ContextWithGasLimitMinGas(ctx context.Context, v uint64) context.Context {
	return context.WithValue(ctx, gasLimitMinGasKey{}, v)
}

// ContextWithGasLimitMaxGas overrides the MaxGas option for this call.
// Only has effect when the [WithGasLimit] client option is enabled.
func ContextWithGasLimitMaxGas(ctx context.Context, v uint64) context.Context {
	return context.WithValue(ctx, gasLimitMaxGasKey{}, v)
}

// ContextWithGasLimitReplace overrides the Replace option for this call.
// Only has effect when the [WithGasLimit] client option is enabled.
func ContextWithGasLimitReplace(ctx context.Context, v bool) context.Context {
	return context.WithValue(ctx, gasLimitReplaceKey{}, v)
}

func gasLimitMultiplier(ctx context.Context, h *hijackGasLimit) float64 {
	if v, ok := ctx.Value(gasLimitMultiplierKey{}).(float64); ok {
		return v
	}
	return h.multiplier
}

func gasLimitMinGas(ctx context.Context, h *hijackGasLimit) uint64 {
	if v, ok := ctx.Value(gasLimitMinGasKey{}).(uint64); ok {
		return v
	}
	return h.minGas
}

func gasLimitMaxGas(ctx context.Context, h *hijackGasLimit) uint64 {
	if v, ok := ctx.Value(gasLimitMaxGasKey{}).(uint64); ok {
		return v
	}
	return h.maxGas
}

func gasLimitReplace(ctx context.Context, h *hijackGasLimit) bool {
	if v, ok := ctx.Value(gasLimitReplaceKey{}).(bool); ok {
		return v
	}
	return h.replace
}

// hijackGasLimit hijacks the "eth_sendTransaction" method and sets the
// "gasLimit" field using the estimate provided by the RPC node.
type hijackGasLimit struct {
	multiplier float64
	minGas     uint64
	maxGas     uint64
	replace    bool
}

// Call implements the [transport.Hijacker] interface.
func (h *hijackGasLimit) Call() func(next transport.CallFunc) transport.CallFunc {
	return func(next transport.CallFunc) transport.CallFunc {
		return func(ctx context.Context, t transport.Transport, result any, method string, args ...any) (err error) {
			if len(args) == 0 || method != "eth_sendTransaction" {
				return next(ctx, t, result, method, args...)
			}
			tx, ok := args[0].(types.Transaction)
			if !ok {
				return next(ctx, t, result, method, args...)
			}
			tc := tx.Call()
			ed := types.GetExecutionData(tx)
			if tc != nil && ed != nil && (gasLimitReplace(ctx, h) || ed.GasLimit == nil) {
				gasLimit, err := (&MethodsCommon{&ClientContext{Transport: t}}).EstimateGas(ctx, tc, types.LatestBlockNumber)
				if err != nil {
					return &ErrHijackFailed{name: "gas limit", err: fmt.Errorf("failed to estimate gas: %w", err)}
				}
				gasLimit, _ = new(big.Float).Mul(new(big.Float).SetUint64(gasLimit), big.NewFloat(gasLimitMultiplier(ctx, h))).Uint64()
				if min := gasLimitMinGas(ctx, h); min > 0 && gasLimit < min {
					gasLimit = min
				}
				if max := gasLimitMaxGas(ctx, h); max > 0 && gasLimit > max {
					gasLimit = max
				}
				ed.GasLimit = &gasLimit
			}
			return next(ctx, t, result, method, args...)
		}
	}
}

// Subscribe implements the [transport.Hijacker] interface.
func (h *hijackGasLimit) Subscribe() func(next transport.SubscribeFunc) transport.SubscribeFunc {
	return nil
}

// Unsubscribe implements the [transport.Hijacker] interface.
func (h *hijackGasLimit) Unsubscribe() func(next transport.UnsubscribeFunc) transport.UnsubscribeFunc {
	return nil
}

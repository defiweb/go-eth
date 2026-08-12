package rpc

import (
	"context"

	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
)

type (
	addressKey        struct{}
	addressReplaceKey struct{}
)

// ContextWithFromAddress overrides the default sender address for this call.
// Only has effect when the [WithDefaultAddress] client option is enabled.
func ContextWithFromAddress(ctx context.Context, v types.Address) context.Context {
	return context.WithValue(ctx, addressKey{}, v)
}

// ContextWithFromAddressReplace overrides the Replace option for this call.
// Only has effect when the [WithDefaultAddress] client option is enabled.
func ContextWithFromAddressReplace(ctx context.Context, v bool) context.Context {
	return context.WithValue(ctx, addressReplaceKey{}, v)
}

func fromAddress(ctx context.Context, h *hijackAddress) types.Address {
	if v, ok := ctx.Value(addressKey{}).(types.Address); ok {
		return v
	}
	return h.address
}

func fromAddressReplace(ctx context.Context, h *hijackAddress) bool {
	if v, ok := ctx.Value(addressReplaceKey{}).(bool); ok {
		return v
	}
	return h.replace
}

// hijackAddress hijacks "eth_sendTransaction", "eth_call",
// "eth_estimateGas", and "eth_createAccessList" to set the "from" field.
type hijackAddress struct {
	address types.Address
	replace bool
}

// Call implements the [transport.Hijacker] interface.
func (h *hijackAddress) Call() func(next transport.CallFunc) transport.CallFunc {
	return func(next transport.CallFunc) transport.CallFunc {
		return func(ctx context.Context, t transport.Transport, result any, method string, args ...any) (err error) {
			if len(args) == 0 {
				return next(ctx, t, result, method, args...)
			}
			var ed *types.ExecutionData
			switch method {
			case "eth_sendTransaction":
				tx, ok := args[0].(types.Transaction)
				if !ok {
					return next(ctx, t, result, method, args...)
				}
				ed = types.GetExecutionData(tx)
			case "eth_call", "eth_estimateGas", "eth_createAccessList":
				call, ok := args[0].(types.Call)
				if !ok {
					return next(ctx, t, result, method, args...)
				}
				ed = types.GetExecutionData(call)
			default:
				return next(ctx, t, result, method, args...)
			}
			if ed != nil && (fromAddressReplace(ctx, h) || ed.From == nil) {
				addr := fromAddress(ctx, h)
				ed.From = &addr
			}
			return next(ctx, t, result, method, args...)
		}
	}
}

// Subscribe implements the [transport.Hijacker] interface.
func (h *hijackAddress) Subscribe() func(next transport.SubscribeFunc) transport.SubscribeFunc {
	return nil
}

// Unsubscribe implements the [transport.Hijacker] interface.
func (h *hijackAddress) Unsubscribe() func(next transport.UnsubscribeFunc) transport.UnsubscribeFunc {
	return nil
}

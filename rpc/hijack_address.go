package rpc

import (
	"context"

	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
)

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
			if ed != nil && (h.replace || ed.From == nil) {
				ed.From = &h.address
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

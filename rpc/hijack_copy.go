package rpc

import (
	"context"

	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
)

// hijackCopy creates deep copies of arguments that may be modified by other
// hijackers.
//
// Each hijacker could create its own copy, but to avoid duplicating work and
// making multiple copies of the same argument, this hijacker creates a single
// copy and shares it among hijackers.
type hijackCopy struct{}

// Call implements the [transport.Hijacker] interface.
func (k *hijackCopy) Call() func(next transport.CallFunc) transport.CallFunc {
	return func(next transport.CallFunc) transport.CallFunc {
		return func(ctx context.Context, t transport.Transport, result any, method string, args ...any) error {
			switch method {
			case "eth_sendTransaction":
				tx, ok := args[0].(types.Transaction)
				if !ok {
					return next(ctx, t, result, method, args...)
				}
				args[0] = tx.Copy()
			case "eth_call", "eth_estimateGas", "eth_createAccessList":
				call, ok := args[0].(types.Call)
				if !ok {
					return next(ctx, t, result, method, args...)
				}
				args[0] = call.Copy()
			}
			return next(ctx, t, result, method, args...)
		}
	}
}

// Subscribe implements the [transport.Hijacker] interface.
func (k *hijackCopy) Subscribe() func(next transport.SubscribeFunc) transport.SubscribeFunc {
	return nil
}

// Unsubscribe implements the [transport.Hijacker] interface.
func (k *hijackCopy) Unsubscribe() func(next transport.UnsubscribeFunc) transport.UnsubscribeFunc {
	return nil
}

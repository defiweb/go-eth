package rpc

import (
	"context"
	"fmt"

	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
)

// hijackAddress hijacks "eth_sendTransaction" and "eth_call" methods and
// sets the "from" field.
type hijackAddress struct {
	address types.Address
	replace bool
}

func (h *hijackAddress) Call() func(next transport.CallFunc) transport.CallFunc {
	return func(next transport.CallFunc) transport.CallFunc {
		return func(ctx context.Context, t transport.Transport, result any, method string, args ...any) (err error) {
			if len(args) == 0 {
				return next(ctx, t, result, method, args...)
			}

			var cd *types.CallFields
			switch method {
			case "eth_sendTransaction":
				// Verify arguments:
				tx, ok := args[0].(types.Transaction)
				if !ok {
					return &ErrHijackFailed{name: "address", err: fmt.Errorf("invalid transaction type: %T", args[0])}
				}

				// Get transaction call data:
				if tx, ok := tx.(types.CallData); ok {
					cd = tx.CallData()
				}
			case "eth_call", "eth_estimateGas":
				// Verify arguments:
				c, ok := args[0].(types.Call)
				if !ok {
					return &ErrHijackFailed{name: "address", err: fmt.Errorf("invalid call type: %T", args[0])}
				}

				// Get call data:
				if c, ok := c.(types.CallData); ok {
					cd = c.CallData()
				}
			default:
				return next(ctx, t, result, method, args...)
			}

			// If call data is nil or the "from" field is already set, continue:
			if cd == nil || (h.replace && cd.From != nil) {
				return next(ctx, t, result, method, args...)
			}

			// Update the "from" field and continue:
			cd.From = &h.address
			return next(ctx, t, result, method, args...)
		}
	}
}

func (h *hijackAddress) Subscribe() func(next transport.SubscribeFunc) transport.SubscribeFunc {
	return nil
}

func (h *hijackAddress) Unsubscribe() func(next transport.UnsubscribeFunc) transport.UnsubscribeFunc {
	return nil
}

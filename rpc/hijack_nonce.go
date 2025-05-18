package rpc

import (
	"context"
	"fmt"

	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
)

// hijackNonce hijacks "eth_sendTransaction" method and sets the "nonce"
// field using the "eth_getTransactionCount" RPC method.
type hijackNonce struct {
	usePendingBlock bool
	replace         bool
}

func (c *hijackNonce) Call() func(next transport.CallFunc) transport.CallFunc {
	return func(next transport.CallFunc) transport.CallFunc {
		return func(ctx context.Context, t transport.Transport, result any, method string, args ...any) (err error) {
			if method != "eth_sendTransaction" || len(args) == 0 {
				return next(ctx, t, result, method, args...)
			}

			// Verify arguments:
			tx, ok := args[0].(types.Transaction)
			if !ok {
				return &ErrHijackFailed{name: "nonce", err: fmt.Errorf("invalid transaction type: %T", args[0])}
			}

			// If the nonce is already set, continue:
			txd := tx.TransactionData()
			if !c.replace && txd.Nonce != nil {
				return next(ctx, t, result, method, args...)
			}

			// Get transaction call data:
			var txcd *types.CallFields
			if tx, ok := tx.(types.CallData); ok {
				txcd = tx.CallData()
			}

			// It is some strange transaction type with no call data, continue:
			if txcd == nil {
				return next(ctx, t, result, method, args...)
			}

			// The "from" field must be set to obtain the nonce:
			if txcd.From == nil {
				return &ErrHijackFailed{name: "nonce", err: fmt.Errorf("'from' field not set")}
			}

			// Get the latest nonce:
			block := types.LatestBlockNumber
			if c.usePendingBlock {
				block = types.PendingBlockNumber
			}
			nonce, err := (&MethodsCommon{Transport: t}).GetTransactionCount(ctx, *txcd.From, block)
			if err != nil {
				return &ErrHijackFailed{name: "nonce", err: fmt.Errorf("failed to get transaction count: %w", err)}
			}

			// Update the nonce and continue:
			txd.Nonce = &nonce
			return next(ctx, t, result, method, args...)
		}
	}
}

func (c *hijackNonce) Subscribe() func(next transport.SubscribeFunc) transport.SubscribeFunc {
	return nil
}

func (c *hijackNonce) Unsubscribe() func(next transport.UnsubscribeFunc) transport.UnsubscribeFunc {
	return nil
}

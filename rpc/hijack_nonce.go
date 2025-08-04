package rpc

import (
	"context"
	"fmt"

	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
)

// hijackNonce hijacks the "eth_sendTransaction" method and sets the "nonce"
// field using the "eth_getTransactionCount" RPC method.
type hijackNonce struct {
	usePendingBlock bool
	replace         bool
}

func (h *hijackNonce) Call() func(next transport.CallFunc) transport.CallFunc {
	return func(next transport.CallFunc) transport.CallFunc {
		return func(ctx context.Context, t transport.Transport, result any, method string, args ...any) (err error) {
			if len(args) == 0 || method != "eth_sendTransaction" {
				return next(ctx, t, result, method, args...)
			}
			tx, ok := args[0].(types.Transaction)
			if !ok {
				return next(ctx, t, result, method, args...)
			}
			td := getTransactionData(tx) // to set the nonce
			cd := getCallData(tx)        // to get the "from" address
			if td != nil && cd != nil && (h.replace || td.Nonce == nil) {
				if cd.From == nil {
					return &ErrHijackFailed{name: "nonce", err: fmt.Errorf("'from' field not set")}
				}
				block := types.LatestBlockNumber
				if h.usePendingBlock {
					block = types.PendingBlockNumber
				}
				nonce, err := (&MethodsCommon{Transport: t}).GetTransactionCount(ctx, *cd.From, block)
				if err != nil {
					return &ErrHijackFailed{name: "nonce", err: fmt.Errorf("failed to get transaction count: %w", err)}
				}
				td.Nonce = &nonce
			}
			return next(ctx, t, result, method, args...)
		}
	}
}

func (h *hijackNonce) Subscribe() func(next transport.SubscribeFunc) transport.SubscribeFunc {
	return nil
}

func (h *hijackNonce) Unsubscribe() func(next transport.UnsubscribeFunc) transport.UnsubscribeFunc {
	return nil
}

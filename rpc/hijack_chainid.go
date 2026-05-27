package rpc

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
)

// hijackChainID hijacks the "eth_sendTransaction" method and sets the
// "chainID" field.
type hijackChainID struct {
	chainID uint64
	replace bool
}

// Call implements the [transport.Hijacker] interface.
func (h *hijackChainID) Call() func(next transport.CallFunc) transport.CallFunc {
	return func(next transport.CallFunc) transport.CallFunc {
		var chainID atomic.Uint64
		chainID.Store(h.chainID)
		return func(ctx context.Context, t transport.Transport, result any, method string, args ...any) (err error) {
			if len(args) == 0 || method != "eth_sendTransaction" {
				return next(ctx, t, result, method, args...)
			}
			tx, ok := args[0].(types.Transaction)
			if !ok {
				return next(ctx, t, result, method, args...)
			}
			sd := types.GetSigningData(tx)
			if sd != nil && (h.replace || sd.ChainID == nil) {
				if chainID.Load() == 0 {
					id, err := (&MethodsCommon{&ClientContext{Transport: t}}).ChainID(ctx)
					if err != nil {
						return &ErrHijackFailed{name: "chain ID", err: fmt.Errorf("failed to get chain ID: %w", err)}
					}
					chainID.Store(id)
				}
				id := chainID.Load()
				sd.ChainID = &id
			}
			return next(ctx, t, result, method, args...)
		}
	}
}

// Subscribe implements the [transport.Hijacker] interface.
func (h *hijackChainID) Subscribe() func(next transport.SubscribeFunc) transport.SubscribeFunc {
	return nil
}

// Unsubscribe implements the [transport.Hijacker] interface.
func (h *hijackChainID) Unsubscribe() func(next transport.UnsubscribeFunc) transport.UnsubscribeFunc {
	return nil
}

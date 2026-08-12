package rpc

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
)

type (
	chainIDKey        struct{}
	chainIDReplaceKey struct{}
)

// ContextWithChainID overrides the chain ID for this call, bypassing any
// cached or auto-detected value.
// Only has effect when the [WithChainID] client option is enabled.
func ContextWithChainID(ctx context.Context, v uint64) context.Context {
	return context.WithValue(ctx, chainIDKey{}, v)
}

// ContextWithChainIDReplace overrides the Replace option for this call.
// Only has effect when the [WithChainID] client option is enabled.
func ContextWithChainIDReplace(ctx context.Context, v bool) context.Context {
	return context.WithValue(ctx, chainIDReplaceKey{}, v)
}

func chainIDValue(ctx context.Context) (uint64, bool) {
	v, ok := ctx.Value(chainIDKey{}).(uint64)
	return v, ok
}

func chainIDReplace(ctx context.Context, h *hijackChainID) bool {
	if v, ok := ctx.Value(chainIDReplaceKey{}).(bool); ok {
		return v
	}
	return h.replace
}

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
			if sd != nil && (chainIDReplace(ctx, h) || sd.ChainID == nil) {
				var id uint64
				if ctxID, ok := chainIDValue(ctx); ok {
					// Context-provided value bypasses the cache entirely.
					id = ctxID
				} else {
					if chainID.Load() == 0 {
						fetched, err := (&MethodsCommon{&ClientContext{Transport: t}}).ChainID(ctx)
						if err != nil {
							return &ErrHijackFailed{name: "chain ID", err: fmt.Errorf("failed to get chain ID: %w", err)}
						}
						chainID.Store(fetched)
					}
					id = chainID.Load()
				}
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

package rpc

import (
	"context"
	"fmt"
	"math/big"

	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
)

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
			if tc != nil && ed != nil && (h.replace || ed.GasLimit == nil) {
				gasLimit, err := (&MethodsCommon{&ClientContext{Transport: t}}).EstimateGas(ctx, tc, types.LatestBlockNumber)
				if err != nil {
					return &ErrHijackFailed{name: "gas limit", err: fmt.Errorf("failed to estimate gas: %w", err)}
				}
				gasLimit, _ = new(big.Float).Mul(new(big.Float).SetUint64(gasLimit), big.NewFloat(h.multiplier)).Uint64()
				if h.minGas > 0 && gasLimit < h.minGas {
					gasLimit = h.minGas
				}
				if h.maxGas > 0 && gasLimit > h.maxGas {
					gasLimit = h.maxGas
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

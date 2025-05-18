package rpc

import (
	"context"
	"fmt"
	"math/big"

	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
)

// hijackGasLimit hijacks "eth_sendTransaction" method and sets the "gasLimit"
// field using the estimate provided by the RPC node.
type hijackGasLimit struct {
	multiplier float64
	minGas     uint64
	maxGas     uint64
	replace    bool
}

func (c *hijackGasLimit) Call() func(next transport.CallFunc) transport.CallFunc {
	return func(next transport.CallFunc) transport.CallFunc {
		return func(ctx context.Context, t transport.Transport, result any, method string, args ...any) (err error) {
			if method != "eth_sendTransaction" || len(args) == 0 {
				return next(ctx, t, result, method, args...)
			}

			// Verify arguments:
			tx, ok := args[0].(types.Transaction)
			if !ok {
				return &ErrHijackFailed{name: "gas limit", err: fmt.Errorf("invalid transaction type: %T", args[0])}
			}

			// Get transaction call data:
			var txcd *types.CallFields
			if tx, ok := tx.(types.CallData); ok {
				txcd = tx.CallData()
			}

			// If the gas limit is already set, continue:
			if !c.replace && txcd.GasLimit != nil {
				return next(ctx, t, result, method, args...)
			}

			// Get the gas estimate from the RPC node:
			gasLimit, err := (&MethodsCommon{Transport: t}).EstimateGas(ctx, tx.Call(), types.LatestBlockNumber)
			if err != nil {
				return &ErrHijackFailed{name: "gas limit", err: fmt.Errorf("failed to estimate gas: %w", err)}
			}

			// Update the gas limit and continue:
			gasLimit, _ = new(big.Float).Mul(new(big.Float).SetUint64(gasLimit), big.NewFloat(c.multiplier)).Uint64()
			if c.minGas > 0 && gasLimit < c.minGas {
				gasLimit = c.minGas
			}
			if c.maxGas > 0 && gasLimit > c.maxGas {
				gasLimit = c.maxGas
			}
			txcd.GasLimit = &gasLimit
			return next(ctx, t, result, method, args...)
		}
	}
}

func (c *hijackGasLimit) Subscribe() func(next transport.SubscribeFunc) transport.SubscribeFunc {
	return nil
}

func (c *hijackGasLimit) Unsubscribe() func(next transport.UnsubscribeFunc) transport.UnsubscribeFunc {
	return nil
}

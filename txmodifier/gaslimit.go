package txmodifier

import (
	"context"
	"fmt"
	"math/big"

	"github.com/defiweb/go-eth/rpc"
	"github.com/defiweb/go-eth/types"
)

// GasLimitEstimator is a transaction modifier that estimates gas limit
// using the rpc.EstimateGas method.
//
// To use this modifier, add it using the WithTXModifiers option when creating
// a new rpc.Client.
type GasLimitEstimator struct {
	multiplier float64
	minGas     uint64
	maxGas     uint64
}

// NewGasLimitEstimator returns a new GasLimitEstimator.
//
// The multiplier is applied to the estimated gas limit.
// The estimated gas is out of range [minGas, maxGas], then error is returned.
// If maxGas is 0, then there is no upper bound.
func NewGasLimitEstimator(multiplier float64, minGas, maxGas uint64) *GasLimitEstimator {
	return &GasLimitEstimator{
		multiplier: multiplier,
		minGas:     minGas,
		maxGas:     maxGas,
	}
}

// Modify implements the rpc.TXModifier interface.
func (e *GasLimitEstimator) Modify(ctx context.Context, client rpc.RPC, tx *types.Transaction) error {
	gasLimit, err := client.EstimateGas(ctx, tx.Call, types.LatestBlockNumber)
	if err != nil {
		return fmt.Errorf("gas limit estimator: failed to estimate gas limit: %w", err)
	}
	gasLimit, _ = new(big.Float).Mul(new(big.Float).SetUint64(gasLimit), big.NewFloat(e.multiplier)).Uint64()
	if gasLimit < e.minGas || (e.maxGas > 0 && gasLimit > e.maxGas) {
		return fmt.Errorf("gas limit estimator: estimated gas limit %d is out of range [%d, %d]", gasLimit, e.minGas, e.maxGas)
	}
	tx.GasLimit = &gasLimit
	return nil
}

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
	replace    bool
}

// GasLimitEstimatorOptions is the options for NewGasLimitEstimator.
type GasLimitEstimatorOptions struct {
	Multiplier float64 // Multiplier is applied to the gas limit.
	MinGas     uint64  // MinGas is the minimum gas limit, or 0 if there is no lower bound.
	MaxGas     uint64  // MaxGas is the maximum gas limit, or 0 if there is no upper bound.
	Replace    bool    // Replace is true if the gas limit should be replaced even if it is already set.
}

// NewGasLimitEstimator returns a new GasLimitEstimator.
func NewGasLimitEstimator(opts GasLimitEstimatorOptions) *GasLimitEstimator {
	return &GasLimitEstimator{
		multiplier: opts.Multiplier,
		minGas:     opts.MinGas,
		maxGas:     opts.MaxGas,
		replace:    opts.Replace,
	}
}

// Modify implements the rpc.TXModifier interface.
func (e *GasLimitEstimator) Modify(ctx context.Context, client rpc.RPC, tx *types.Transaction) error {
	if !e.replace && tx.GasLimit != nil {
		return nil
	}
	gasLimit, _, err := client.EstimateGas(ctx, &tx.Call, types.LatestBlockNumber)
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

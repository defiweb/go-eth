package txmodifier

import (
	"context"
	"fmt"
	"math/big"

	"github.com/defiweb/go-eth/rpc"
	"github.com/defiweb/go-eth/types"
)

// LegacyGasFeeEstimator is a transaction modifier that estimates gas fee
// using the rpc.GasPrice method.
//
// It sets transaction type to types.LegacyTxType or types.AccessListTxType if
// an access list is provided.
//
// To use this modifier, add it using the WithTXModifiers option when creating
// a new rpc.Client.
type LegacyGasFeeEstimator struct {
	multiplier  float64
	minGasPrice *big.Int
	maxGasPrice *big.Int
}

// NewLegacyGasFeeEstimator returns a new LegacyGasFeeEstimator.
//
// The multiplier is applied to the gas price.
// The estimated gas price is clamped to [minGasPrice, maxGasPrice].
// If minGasPrice or maxGasPrice is nil, then there is no lower or upper bound.
func NewLegacyGasFeeEstimator(multiplier float64, minGasPrice, maxGasPrice *big.Int) *LegacyGasFeeEstimator {
	return &LegacyGasFeeEstimator{
		multiplier:  multiplier,
		minGasPrice: minGasPrice,
		maxGasPrice: maxGasPrice,
	}
}

// Modify implements the rpc.TXModifier interface.
func (e *LegacyGasFeeEstimator) Modify(ctx context.Context, client rpc.RPC, tx *types.Transaction) error {
	gasPrice, err := client.GasPrice(ctx)
	if err != nil {
		return fmt.Errorf("legacy gas fee estimator: failed to get gas price: %w", err)
	}
	gasPrice, _ = new(big.Float).Mul(new(big.Float).SetInt(gasPrice), big.NewFloat(e.multiplier)).Int(nil)
	if e.minGasPrice != nil && gasPrice.Cmp(e.minGasPrice) < 0 {
		gasPrice = e.minGasPrice
	}
	if e.maxGasPrice != nil && gasPrice.Cmp(e.maxGasPrice) > 0 {
		gasPrice = e.maxGasPrice
	}
	tx.GasPrice = gasPrice
	switch {
	case tx.AccessList != nil:
		tx.Type = types.AccessListTxType
	default:
		tx.Type = types.LegacyTxType
	}
	return nil
}

// EIP1559GasFeeEstimator is a transaction modifier that estimates gas fee
// using the rpc.GasPrice and rpc.MaxPriorityFeePerGas methods.
//
// It sets transaction type to types.DynamicFeeTxType.
type EIP1559GasFeeEstimator struct {
	gasPriceMultiplier          float64
	priorityFeePerGasMultiplier float64
	minGasPrice                 *big.Int
	maxGasPrice                 *big.Int
	minPriorityFeePerGas        *big.Int
	maxPriorityFeePerGas        *big.Int
}

// NewEIP1559GasFeeEstimator returns a new EIP1559GasFeeEstimator.
//
// The gasPriceMultiplier and priorityFeePerGasMultiplier are applied to the
// gas price and priority fee per gas respectively.
//
// The estimated gas price and priority fee per gas are clamped to
// [minGasPrice, maxGasPrice] and [minPriorityFeePerGas, maxPriorityFeePerGas]
// respectively.
//
// If minGasPrice or maxGasPrice is nil, then there is no lower or upper bound
// for gas price.
//
// If minPriorityFeePerGas or maxPriorityFeePerGas is nil, then there is no
// lower or upper bound for priority fee per gas.
//
// To use this modifier, add it using the WithTXModifiers option when creating
// a new rpc.Client.
func NewEIP1559GasFeeEstimator(
	gasPriceMultiplier float64,
	priorityFeePerGasMultiplier float64,
	minGasPrice,
	maxGasPrice,
	minPriorityFeePerGas,
	maxPriorityFeePerGas *big.Int,
) *EIP1559GasFeeEstimator {
	return &EIP1559GasFeeEstimator{
		gasPriceMultiplier:          gasPriceMultiplier,
		priorityFeePerGasMultiplier: priorityFeePerGasMultiplier,
		minGasPrice:                 minGasPrice,
		maxGasPrice:                 maxGasPrice,
		minPriorityFeePerGas:        minPriorityFeePerGas,
		maxPriorityFeePerGas:        maxPriorityFeePerGas,
	}
}

// Modify implements the rpc.TXModifier interface.
func (e *EIP1559GasFeeEstimator) Modify(ctx context.Context, client rpc.RPC, tx *types.Transaction) error {
	gasPrice, err := client.GasPrice(ctx)
	if err != nil {
		return fmt.Errorf("EIP-1669 gas fee estimator: failed to get gas price: %w", err)
	}
	priorityFeePerGas, err := client.MaxPriorityFeePerGas(ctx)
	if err != nil {
		return fmt.Errorf("EIP-1559 gas fee estimator: failed to get max priority fee per gas: %w", err)
	}
	gasPrice, _ = new(big.Float).Mul(new(big.Float).SetInt(gasPrice), big.NewFloat(e.gasPriceMultiplier)).Int(nil)
	priorityFeePerGas, _ = new(big.Float).Mul(new(big.Float).SetInt(priorityFeePerGas), big.NewFloat(e.priorityFeePerGasMultiplier)).Int(nil)
	if e.minGasPrice != nil && gasPrice.Cmp(e.minGasPrice) < 0 {
		gasPrice = e.minGasPrice
	}
	if e.maxGasPrice != nil && gasPrice.Cmp(e.maxGasPrice) > 0 {
		gasPrice = e.maxGasPrice
	}
	if e.minPriorityFeePerGas != nil && priorityFeePerGas.Cmp(e.minPriorityFeePerGas) < 0 {
		priorityFeePerGas = e.minPriorityFeePerGas
	}
	if e.maxPriorityFeePerGas != nil && priorityFeePerGas.Cmp(e.maxPriorityFeePerGas) > 0 {
		priorityFeePerGas = e.maxPriorityFeePerGas
	}
	tx.GasPrice = gasPrice
	tx.MaxPriorityFeePerGas = priorityFeePerGas
	tx.Type = types.DynamicFeeTxType
	return nil
}

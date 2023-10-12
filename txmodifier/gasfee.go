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
	replace     bool
}

// LegacyGasFeeEstimatorOptions is the options for NewLegacyGasFeeEstimator.
type LegacyGasFeeEstimatorOptions struct {
	Multiplier  float64  // Multiplier is applied to the gas price.
	MinGasPrice *big.Int // MinGasPrice is the minimum gas price, or nil if there is no lower bound.
	MaxGasPrice *big.Int // MaxGasPrice is the maximum gas price, or nil if there is no upper bound.
	Replace     bool     // Replace is true if the gas price should be replaced even if it is already set.
}

// NewLegacyGasFeeEstimator returns a new LegacyGasFeeEstimator.
func NewLegacyGasFeeEstimator(opts LegacyGasFeeEstimatorOptions) *LegacyGasFeeEstimator {
	return &LegacyGasFeeEstimator{
		multiplier:  opts.Multiplier,
		minGasPrice: opts.MinGasPrice,
		maxGasPrice: opts.MaxGasPrice,
		replace:     opts.Replace,
	}
}

// Modify implements the rpc.TXModifier interface.
func (e *LegacyGasFeeEstimator) Modify(ctx context.Context, client rpc.RPC, tx *types.Transaction) error {
	if !e.replace && tx.GasPrice != nil {
		return nil
	}
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
	tx.MaxFeePerGas = nil
	tx.MaxPriorityFeePerGas = nil
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
	replace                     bool
}

// EIP1559GasFeeEstimatorOptions is the options for NewEIP1559GasFeeEstimator.
type EIP1559GasFeeEstimatorOptions struct {
	GasPriceMultiplier          float64  // GasPriceMultiplier is applied to the gas price.
	PriorityFeePerGasMultiplier float64  // PriorityFeePerGasMultiplier is applied to the priority fee per gas.
	MinGasPrice                 *big.Int // MinGasPrice is the minimum gas price, or nil if there is no lower bound.
	MaxGasPrice                 *big.Int // MaxGasPrice is the maximum gas price, or nil if there is no upper bound.
	MinPriorityFeePerGas        *big.Int // MinPriorityFeePerGas is the minimum priority fee per gas, or nil if there is no lower bound.
	MaxPriorityFeePerGas        *big.Int // MaxPriorityFeePerGas is the maximum priority fee per gas, or nil if there is no upper bound.
	Replace                     bool     // Replace is true if the gas price should be replaced even if it is already set.
}

// NewEIP1559GasFeeEstimator returns a new EIP1559GasFeeEstimator.
//
// To use this modifier, add it using the WithTXModifiers option when creating
// a new rpc.Client.
func NewEIP1559GasFeeEstimator(opts EIP1559GasFeeEstimatorOptions) *EIP1559GasFeeEstimator {
	return &EIP1559GasFeeEstimator{
		gasPriceMultiplier:          opts.GasPriceMultiplier,
		priorityFeePerGasMultiplier: opts.PriorityFeePerGasMultiplier,
		minGasPrice:                 opts.MinGasPrice,
		maxGasPrice:                 opts.MaxGasPrice,
		minPriorityFeePerGas:        opts.MinPriorityFeePerGas,
		maxPriorityFeePerGas:        opts.MaxPriorityFeePerGas,
		replace:                     opts.Replace,
	}
}

// Modify implements the rpc.TXModifier interface.
func (e *EIP1559GasFeeEstimator) Modify(ctx context.Context, client rpc.RPC, tx *types.Transaction) error {
	if !e.replace && tx.MaxFeePerGas != nil && tx.MaxPriorityFeePerGas != nil {
		return nil
	}
	maxFeePerGas, err := client.GasPrice(ctx)
	if err != nil {
		return fmt.Errorf("EIP-1559 gas fee estimator: failed to get gas price: %w", err)
	}
	priorityFeePerGas, err := client.MaxPriorityFeePerGas(ctx)
	if err != nil {
		return fmt.Errorf("EIP-1559 gas fee estimator: failed to get max priority fee per gas: %w", err)
	}
	maxFeePerGas, _ = new(big.Float).Mul(new(big.Float).SetInt(maxFeePerGas), big.NewFloat(e.gasPriceMultiplier)).Int(nil)
	priorityFeePerGas, _ = new(big.Float).Mul(new(big.Float).SetInt(priorityFeePerGas), big.NewFloat(e.priorityFeePerGasMultiplier)).Int(nil)
	if e.minGasPrice != nil && maxFeePerGas.Cmp(e.minGasPrice) < 0 {
		maxFeePerGas = e.minGasPrice
	}
	if e.maxGasPrice != nil && maxFeePerGas.Cmp(e.maxGasPrice) > 0 {
		maxFeePerGas = e.maxGasPrice
	}
	if e.minPriorityFeePerGas != nil && priorityFeePerGas.Cmp(e.minPriorityFeePerGas) < 0 {
		priorityFeePerGas = e.minPriorityFeePerGas
	}
	if e.maxPriorityFeePerGas != nil && priorityFeePerGas.Cmp(e.maxPriorityFeePerGas) > 0 {
		priorityFeePerGas = e.maxPriorityFeePerGas
	}
	if maxFeePerGas.Cmp(priorityFeePerGas) < 0 {
		priorityFeePerGas = maxFeePerGas
	}
	tx.GasPrice = nil
	tx.MaxFeePerGas = maxFeePerGas
	tx.MaxPriorityFeePerGas = priorityFeePerGas
	tx.Type = types.DynamicFeeTxType
	return nil
}

package rpc

import (
	"context"
	"fmt"
	"math/big"

	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
)

// hijackLegacyGasFee hijacks the "eth_sendTransaction" method and sets the
// "gasPrice" field using the estimate provided by the RPC node.
type hijackLegacyGasFee struct {
	multiplier      float64
	minGasPrice     *big.Int
	maxGasPrice     *big.Int
	replace         bool
	allowChangeType bool
}

// Call implements the [transport.Hijacker] interface.
func (h *hijackLegacyGasFee) Call() func(next transport.CallFunc) transport.CallFunc {
	return func(next transport.CallFunc) transport.CallFunc {
		return func(ctx context.Context, t transport.Transport, result any, method string, args ...any) (err error) {
			if len(args) == 0 || method != "eth_sendTransaction" {
				return next(ctx, t, result, method, args...)
			}
			tx, ok := args[0].(types.Transaction)
			if !ok {
				return next(ctx, t, result, method, args...)
			}
			if h.allowChangeType {
				tx = convertTXToLegacyPrice(tx)
				args[0] = tx
			}
			lpd := types.GetLegacyFeeData(tx)
			if lpd != nil && (h.replace || lpd.GasPrice == nil) {
				gasPrice, err := (&MethodsCommon{&ClientContext{Transport: t}}).GasPrice(ctx)
				if err != nil {
					return &ErrHijackFailed{name: "legacy gas price", err: fmt.Errorf("failed to get gas price: %w", err)}
				}
				gasPrice, _ = new(big.Float).Mul(new(big.Float).SetInt(gasPrice), big.NewFloat(h.multiplier)).Int(nil)
				if h.minGasPrice != nil && gasPrice.Cmp(h.minGasPrice) < 0 {
					gasPrice = h.minGasPrice
				}
				if h.maxGasPrice != nil && gasPrice.Cmp(h.maxGasPrice) > 0 {
					gasPrice = h.maxGasPrice
				}
				lpd.GasPrice = gasPrice
			}
			return next(ctx, t, result, method, args...)
		}
	}
}

// Subscribe implements the [transport.Hijacker] interface.
func (h *hijackLegacyGasFee) Subscribe() func(next transport.SubscribeFunc) transport.SubscribeFunc {
	return nil
}

// Unsubscribe implements the [transport.Hijacker] interface.
func (h *hijackLegacyGasFee) Unsubscribe() func(next transport.UnsubscribeFunc) transport.UnsubscribeFunc {
	return nil
}

// hijackDynamicGasFee hijacks the "eth_sendTransaction" method and sets the
// "maxFeePerGas" and "maxPriorityFeePerGas" fields using the estimates
// provided by the RPC node.
type hijackDynamicGasFee struct {
	gasPriceMultiplier          float64
	priorityFeePerGasMultiplier float64
	minGasPrice                 *big.Int
	maxGasPrice                 *big.Int
	minPriorityFeePerGas        *big.Int
	maxPriorityFeePerGas        *big.Int
	replace                     bool
	allowChangeType             bool
}

// Call implements the [transport.Hijacker] interface.
func (h *hijackDynamicGasFee) Call() func(next transport.CallFunc) transport.CallFunc {
	return func(next transport.CallFunc) transport.CallFunc {
		return func(ctx context.Context, t transport.Transport, result any, method string, args ...any) (err error) {
			if len(args) == 0 || method != "eth_sendTransaction" {
				return next(ctx, t, result, method, args...)
			}

			// Obtain the dynamic fee data.
			tx, ok := args[0].(types.Transaction)
			if !ok {
				return next(ctx, t, result, method, args...)
			}
			if h.allowChangeType {
				tx = convertTXToDynamicFee(tx)
				args[0] = tx
			}
			dfd := types.GetDynamicFeeData(tx)

			// Set the dynamic fee fields if necessary.
			if dfd != nil && (h.replace || dfd.MaxFeePerGas == nil || dfd.MaxPriorityFeePerGas == nil) {
				// Fetch current gas prices from the RPC node.
				maxFeePerGas, err := (&MethodsCommon{&ClientContext{Transport: t}}).GasPrice(ctx)
				if err != nil {
					return &ErrHijackFailed{name: "dynamic gas fee", err: fmt.Errorf("failed to get gas price: %w", err)}
				}
				priorityFeePerGas, err := (&MethodsCommon{&ClientContext{Transport: t}}).MaxPriorityFeePerGas(ctx)
				if err != nil {
					return &ErrHijackFailed{name: "dynamic gas fee", err: fmt.Errorf("failed to get priority fee per gas: %w", err)}
				}

				// Apply multipliers and limits and set the fields.
				maxFeePerGas, _ = new(big.Float).Mul(new(big.Float).SetInt(maxFeePerGas), big.NewFloat(h.gasPriceMultiplier)).Int(nil)
				priorityFeePerGas, _ = new(big.Float).Mul(new(big.Float).SetInt(priorityFeePerGas), big.NewFloat(h.priorityFeePerGasMultiplier)).Int(nil)
				if h.minGasPrice != nil && maxFeePerGas.Cmp(h.minGasPrice) < 0 {
					maxFeePerGas = h.minGasPrice
				}
				if h.maxGasPrice != nil && maxFeePerGas.Cmp(h.maxGasPrice) > 0 {
					maxFeePerGas = h.maxGasPrice
				}
				if h.minPriorityFeePerGas != nil && priorityFeePerGas.Cmp(h.minPriorityFeePerGas) < 0 {
					priorityFeePerGas = h.minPriorityFeePerGas
				}
				if h.maxPriorityFeePerGas != nil && priorityFeePerGas.Cmp(h.maxPriorityFeePerGas) > 0 {
					priorityFeePerGas = h.maxPriorityFeePerGas
				}
				if maxFeePerGas.Cmp(priorityFeePerGas) < 0 {
					priorityFeePerGas = maxFeePerGas
				}
				dfd.MaxFeePerGas = maxFeePerGas
				dfd.MaxPriorityFeePerGas = priorityFeePerGas
			}

			return next(ctx, t, result, method, args...)
		}
	}
}

// Subscribe implements the [transport.Hijacker] interface.
func (h *hijackDynamicGasFee) Subscribe() func(next transport.SubscribeFunc) transport.SubscribeFunc {
	return nil
}

// Unsubscribe implements the [transport.Hijacker] interface.
func (h *hijackDynamicGasFee) Unsubscribe() func(next transport.UnsubscribeFunc) transport.UnsubscribeFunc {
	return nil
}

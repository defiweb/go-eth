package txmodifier

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/defiweb/go-eth/types"
)

func TestLegacyGasFeeEstimator_Modify(t *testing.T) {
	ctx := context.Background()

	t.Run("successful gas fee estimation", func(t *testing.T) {
		tx := &types.Transaction{}
		rpcMock := new(mockRPC)
		rpcMock.On("GasPrice", ctx).Return(big.NewInt(1000), nil)
		estimator := NewLegacyGasFeeEstimator(LegacyGasFeeEstimatorOptions{
			Multiplier:  1.5,
			MinGasPrice: big.NewInt(500),
			MaxGasPrice: big.NewInt(2000),
		})
		err := estimator.Modify(ctx, rpcMock, tx)

		assert.NoError(t, err)
		assert.Equal(t, big.NewInt(1500), tx.GasPrice)
		assert.Equal(t, types.LegacyTxType, tx.Type)
	})

	t.Run("gas fee estimation error", func(t *testing.T) {
		tx := &types.Transaction{}
		rpcMock := new(mockRPC)
		rpcMock.On("GasPrice", ctx).Return((*big.Int)(nil), errors.New("rpc error"))

		estimator := NewLegacyGasFeeEstimator(LegacyGasFeeEstimatorOptions{
			Multiplier:  1.5,
			MinGasPrice: big.NewInt(500),
			MaxGasPrice: big.NewInt(2000),
		})
		err := estimator.Modify(ctx, rpcMock, tx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get gas price")
	})

	t.Run("gas fee below min bound", func(t *testing.T) {
		tx := &types.Transaction{}
		rpcMock := new(mockRPC)
		rpcMock.On("GasPrice", ctx).Return(big.NewInt(300), nil)

		estimator := NewLegacyGasFeeEstimator(LegacyGasFeeEstimatorOptions{
			Multiplier:  1.0,
			MinGasPrice: big.NewInt(500),
			MaxGasPrice: big.NewInt(2000),
		})
		err := estimator.Modify(ctx, rpcMock, tx)

		assert.NoError(t, err)
		assert.Equal(t, big.NewInt(500), tx.GasPrice) // should be clamped to minGasPrice
	})

	t.Run("gas fee above max bound", func(t *testing.T) {
		tx := &types.Transaction{}
		rpcMock := new(mockRPC)
		rpcMock.On("GasPrice", ctx).Return(big.NewInt(2500), nil)

		estimator := NewLegacyGasFeeEstimator(LegacyGasFeeEstimatorOptions{
			Multiplier:  1.0,
			MinGasPrice: big.NewInt(500),
			MaxGasPrice: big.NewInt(2000),
		})
		err := estimator.Modify(ctx, rpcMock, tx)

		assert.NoError(t, err)
		assert.Equal(t, big.NewInt(2000), tx.GasPrice) // should be clamped to maxGasPrice
	})
}

func TestEIP1559GasFeeEstimator_Modify(t *testing.T) {
	ctx := context.Background()

	t.Run("successful gas fee estimation", func(t *testing.T) {
		tx := &types.Transaction{}
		rpcMock := new(mockRPC)
		rpcMock.On("GasPrice", ctx).Return(big.NewInt(1000), nil)
		rpcMock.On("MaxPriorityFeePerGas", ctx).Return(big.NewInt(5), nil)

		estimator := NewEIP1559GasFeeEstimator(EIP1559GasFeeEstimatorOptions{
			GasPriceMultiplier:          1.5,
			PriorityFeePerGasMultiplier: 2.0,
			MinGasPrice:                 big.NewInt(500),
			MaxGasPrice:                 big.NewInt(2000),
			MinPriorityFeePerGas:        big.NewInt(2),
			MaxPriorityFeePerGas:        big.NewInt(10),
		})
		err := estimator.Modify(ctx, rpcMock, tx)

		assert.NoError(t, err)
		assert.Equal(t, big.NewInt(1500), tx.MaxFeePerGas)
		assert.Equal(t, big.NewInt(10), tx.MaxPriorityFeePerGas)
		assert.Equal(t, types.DynamicFeeTxType, tx.Type)
	})

	t.Run("gas fee estimation error (GasPrice call error)", func(t *testing.T) {
		tx := &types.Transaction{}
		rpcMock := new(mockRPC)
		rpcMock.On("GasPrice", ctx).Return((*big.Int)(nil), errors.New("rpc error"))

		estimator := NewEIP1559GasFeeEstimator(EIP1559GasFeeEstimatorOptions{
			GasPriceMultiplier:          1.5,
			PriorityFeePerGasMultiplier: 2.0,
			MinGasPrice:                 big.NewInt(500),
			MaxGasPrice:                 big.NewInt(2000),
			MinPriorityFeePerGas:        big.NewInt(2),
			MaxPriorityFeePerGas:        big.NewInt(10),
		})
		err := estimator.Modify(ctx, rpcMock, tx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get gas price")
	})

	t.Run("gas fee estimation error (MaxPriorityFeePerGas call error)", func(t *testing.T) {
		tx := &types.Transaction{}
		rpcMock := new(mockRPC)
		rpcMock.On("GasPrice", ctx).Return(big.NewInt(1000), nil)
		rpcMock.On("MaxPriorityFeePerGas", ctx).Return((*big.Int)(nil), errors.New("rpc error"))

		estimator := NewEIP1559GasFeeEstimator(EIP1559GasFeeEstimatorOptions{
			GasPriceMultiplier:          1.5,
			PriorityFeePerGasMultiplier: 2.0,
			MinGasPrice:                 big.NewInt(500),
			MaxGasPrice:                 big.NewInt(2000),
			MinPriorityFeePerGas:        big.NewInt(2),
			MaxPriorityFeePerGas:        big.NewInt(10),
		})
		err := estimator.Modify(ctx, rpcMock, tx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get max priority fee per gas")
	})

	t.Run("gas fee below min bound", func(t *testing.T) {
		tx := &types.Transaction{}
		rpcMock := new(mockRPC)
		rpcMock.On("GasPrice", ctx).Return(big.NewInt(300), nil)
		rpcMock.On("MaxPriorityFeePerGas", ctx).Return(big.NewInt(1), nil)

		estimator := NewEIP1559GasFeeEstimator(EIP1559GasFeeEstimatorOptions{
			GasPriceMultiplier:          1.0,
			PriorityFeePerGasMultiplier: 1.0,
			MinGasPrice:                 big.NewInt(500),
			MaxGasPrice:                 big.NewInt(2000),
			MinPriorityFeePerGas:        big.NewInt(2),
			MaxPriorityFeePerGas:        big.NewInt(10),
		})
		err := estimator.Modify(ctx, rpcMock, tx)

		assert.NoError(t, err)
		assert.Equal(t, big.NewInt(500), tx.MaxFeePerGas)       // should be clamped to minGasPrice
		assert.Equal(t, big.NewInt(2), tx.MaxPriorityFeePerGas) // should be clamped to minPriorityFeePerGas
	})

	t.Run("gas fee above max bound", func(t *testing.T) {
		tx := &types.Transaction{}
		rpcMock := new(mockRPC)
		rpcMock.On("GasPrice", ctx).Return(big.NewInt(2500), nil)
		rpcMock.On("MaxPriorityFeePerGas", ctx).Return(big.NewInt(12), nil)

		estimator := NewEIP1559GasFeeEstimator(EIP1559GasFeeEstimatorOptions{
			GasPriceMultiplier:          1.0,
			PriorityFeePerGasMultiplier: 1.0,
			MinGasPrice:                 big.NewInt(500),
			MaxGasPrice:                 big.NewInt(2000),
			MinPriorityFeePerGas:        big.NewInt(2),
			MaxPriorityFeePerGas:        big.NewInt(10),
		})
		err := estimator.Modify(ctx, rpcMock, tx)

		assert.NoError(t, err)
		assert.Equal(t, big.NewInt(2000), tx.MaxFeePerGas)       // should be clamped to maxGasPrice
		assert.Equal(t, big.NewInt(10), tx.MaxPriorityFeePerGas) // should be clamped to maxPriorityFeePerGas
	})

	t.Run("gas tip fee higher than gas fee", func(t *testing.T) {
		tx := &types.Transaction{}
		rpcMock := new(mockRPC)
		rpcMock.On("GasPrice", ctx).Return(big.NewInt(500), nil)
		rpcMock.On("MaxPriorityFeePerGas", ctx).Return(big.NewInt(2500), nil)

		estimator := NewEIP1559GasFeeEstimator(EIP1559GasFeeEstimatorOptions{
			GasPriceMultiplier:          1.0,
			PriorityFeePerGasMultiplier: 1.0,
		})
		err := estimator.Modify(ctx, rpcMock, tx)

		assert.NoError(t, err)
		assert.Equal(t, big.NewInt(500), tx.MaxFeePerGas)
		assert.Equal(t, big.NewInt(500), tx.MaxPriorityFeePerGas) // should not be higher than tx.MaxFeePerGas
	})
}

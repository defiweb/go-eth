package txmodifier

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/defiweb/go-eth/types"
)

func TestGasLimitEstimator_Modify(t *testing.T) {
	ctx := context.Background()

	t.Run("successful gas estimation", func(t *testing.T) {
		tx := &types.Transaction{}
		rpcMock := new(mockRPC)
		rpcMock.On("EstimateGas", ctx, &tx.Call, types.LatestBlockNumber).Return(uint64(1000), &tx.Call, nil)

		estimator := NewGasLimitEstimator(GasLimitEstimatorOptions{
			Multiplier: 1.5,
			MinGas:     500,
			MaxGas:     2000,
		})
		err := estimator.Modify(ctx, rpcMock, tx)

		assert.NoError(t, err)
		assert.Equal(t, uint64(1500), *tx.GasLimit)
	})

	t.Run("gas estimation error", func(t *testing.T) {
		tx := &types.Transaction{}
		rpcMock := new(mockRPC)
		rpcMock.On("EstimateGas", ctx, &tx.Call, types.LatestBlockNumber).Return(uint64(0), &tx.Call, errors.New("rpc error"))

		estimator := NewGasLimitEstimator(GasLimitEstimatorOptions{
			Multiplier: 1.5,
			MinGas:     500,
			MaxGas:     2000,
		})
		err := estimator.Modify(ctx, rpcMock, tx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to estimate gas")
	})

	t.Run("gas out of range", func(t *testing.T) {
		tx := &types.Transaction{}
		rpcMock := new(mockRPC)
		rpcMock.On("EstimateGas", ctx, &tx.Call, types.LatestBlockNumber).Return(uint64(3000), &tx.Call, nil)

		estimator := NewGasLimitEstimator(GasLimitEstimatorOptions{
			Multiplier: 1.5,
			MinGas:     500,
			MaxGas:     2000,
		})
		err := estimator.Modify(ctx, rpcMock, tx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "estimated gas")
	})
}

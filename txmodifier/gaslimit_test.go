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
	tx := &types.Transaction{Call: types.Call{}}

	t.Run("successful gas estimation", func(t *testing.T) {
		rpcMock := new(mockRPC)
		rpcMock.On("EstimateGas", ctx, tx.Call, types.LatestBlockNumber).Return(uint64(1000), nil)

		estimator := NewGasLimitEstimator(1.5, 500, 2000)
		err := estimator.Modify(ctx, rpcMock, tx)

		assert.NoError(t, err)
		assert.Equal(t, uint64(1500), *tx.GasLimit)
	})

	t.Run("gas estimation error", func(t *testing.T) {
		rpcMock := new(mockRPC)
		rpcMock.On("EstimateGas", ctx, tx.Call, types.LatestBlockNumber).Return(uint64(0), errors.New("rpc error"))

		estimator := NewGasLimitEstimator(1.5, 500, 2000)
		err := estimator.Modify(ctx, rpcMock, tx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to estimate gas")
	})

	t.Run("gas out of range", func(t *testing.T) {
		rpcMock := new(mockRPC)
		rpcMock.On("EstimateGas", ctx, tx.Call, types.LatestBlockNumber).Return(uint64(3000), nil)

		estimator := NewGasLimitEstimator(1.5, 500, 2000)
		err := estimator.Modify(ctx, rpcMock, tx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "estimated gas")
	})
}

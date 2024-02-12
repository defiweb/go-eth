package txmodifier

import (
	"context"
	"math/big"

	"github.com/stretchr/testify/mock"

	"github.com/defiweb/go-eth/rpc"
	"github.com/defiweb/go-eth/types"
)

type mockRPC struct {
	rpc.Client
	mock.Mock
}

func (m *mockRPC) ChainID(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *mockRPC) EstimateGas(ctx context.Context, call *types.Call, block types.BlockNumber) (uint64, *types.Call, error) {
	args := m.Called(ctx, call, block)
	return args.Get(0).(uint64), call, args.Error(2)
}

func (m *mockRPC) GasPrice(ctx context.Context) (*big.Int, error) {
	args := m.Called(ctx)
	return args.Get(0).(*big.Int), args.Error(1)
}

func (m *mockRPC) MaxPriorityFeePerGas(ctx context.Context) (*big.Int, error) {
	args := m.Called(ctx)
	return args.Get(0).(*big.Int), args.Error(1)
}

func (m *mockRPC) GetTransactionCount(ctx context.Context, address types.Address, block types.BlockNumber) (uint64, error) {
	args := m.Called(ctx, address, block)
	return args.Get(0).(uint64), args.Error(1)
}

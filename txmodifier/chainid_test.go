package txmodifier

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/defiweb/go-eth/types"
)

func TestChainIDSetter_Modify(t *testing.T) {
	ctx := context.Background()
	fromAddress := types.MustAddressFromHex("0x1234567890abcdef1234567890abcdef12345678")

	t.Run("set chain ID", func(t *testing.T) {
		tx := &types.Transaction{Call: types.Call{From: &fromAddress}}
		rpcMock := new(mockRPC)
		rpcMock.On("ChainID", ctx).Return(uint64(1), nil)

		provider := NewChainIDProvider(ChainIDProviderOptions{
			Replace: false,
			Cache:   false,
		})
		err := provider.Modify(ctx, rpcMock, tx)

		assert.NoError(t, err)
		assert.Equal(t, uint64(1), *tx.ChainID)
	})

	t.Run("replace chain ID", func(t *testing.T) {
		tx := &types.Transaction{Call: types.Call{From: &fromAddress}, ChainID: uint64Ptr(2)}
		rpcMock := new(mockRPC)
		rpcMock.On("ChainID", ctx).Return(uint64(1), nil)

		provider := NewChainIDProvider(ChainIDProviderOptions{
			Replace: true,
			Cache:   false,
		})
		err := provider.Modify(ctx, rpcMock, tx)

		assert.NoError(t, err)
		assert.NotEqual(t, uint64(2), *tx.ChainID)
	})

	t.Run("do not replace chain ID", func(t *testing.T) {
		tx := &types.Transaction{Call: types.Call{From: &fromAddress}, ChainID: uint64Ptr(2)}
		rpcMock := new(mockRPC)
		rpcMock.On("ChainID", ctx).Return(uint64(1), nil)

		provider := NewChainIDProvider(ChainIDProviderOptions{
			Replace: false,
			Cache:   false,
		})
		err := provider.Modify(ctx, rpcMock, tx)

		assert.NoError(t, err)
		assert.NotEqual(t, uint64(1), *tx.ChainID)
	})

	t.Run("cache chain ID", func(t *testing.T) {
		tx := &types.Transaction{Call: types.Call{From: &fromAddress}, ChainID: uint64Ptr(2)}
		rpcMock := new(mockRPC)
		rpcMock.On("ChainID", ctx).Return(uint64(1), nil).Once()

		provider := NewChainIDProvider(ChainIDProviderOptions{
			Replace: true,
			Cache:   true,
		})
		_ = provider.Modify(ctx, rpcMock, tx)
		_ = provider.Modify(ctx, rpcMock, tx)
	})
}

func uint64Ptr(i uint64) *uint64 {
	return &i
}

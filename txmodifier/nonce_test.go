package txmodifier

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/defiweb/go-eth/types"
)

func TestNonceProvider_Modify(t *testing.T) {
	ctx := context.Background()
	fromAddress := types.MustAddressFromHex("0x1234567890abcdef1234567890abcdef12345678")

	t.Run("nonce fetch from latest block", func(t *testing.T) {
		tx := &types.Transaction{Call: types.Call{From: &fromAddress}}
		rpcMock := new(mockRPC)
		rpcMock.On("GetTransactionCount", ctx, fromAddress, types.LatestBlockNumber).Return(uint64(10), nil)

		provider := NewNonceProvider(NonceProviderOptions{
			UsePendingBlock: false,
		})
		err := provider.Modify(ctx, rpcMock, tx)

		assert.NoError(t, err)
		assert.Equal(t, uint64(10), *tx.Nonce)
	})

	t.Run("nonce fetch from pending block", func(t *testing.T) {
		tx := &types.Transaction{Call: types.Call{From: &fromAddress}}
		rpcMock := new(mockRPC)
		rpcMock.On("GetTransactionCount", ctx, fromAddress, types.PendingBlockNumber).Return(uint64(11), nil)

		provider := NewNonceProvider(NonceProviderOptions{
			UsePendingBlock: true,
		})
		err := provider.Modify(ctx, rpcMock, tx)

		assert.NoError(t, err)
		assert.Equal(t, uint64(11), *tx.Nonce)
	})

	t.Run("missing from address", func(t *testing.T) {
		txWithoutFrom := &types.Transaction{}
		provider := NewNonceProvider(NonceProviderOptions{
			UsePendingBlock: true,
		})
		err := provider.Modify(ctx, nil, txWithoutFrom)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "nonce provider: missing from address")
	})

	t.Run("nonce fetch error", func(t *testing.T) {
		tx := &types.Transaction{Call: types.Call{From: &fromAddress}}
		rpcMock := new(mockRPC)
		rpcMock.On("GetTransactionCount", ctx, fromAddress, types.LatestBlockNumber).Return(uint64(0), errors.New("rpc error"))

		provider := NewNonceProvider(NonceProviderOptions{
			UsePendingBlock: false,
		})
		err := provider.Modify(ctx, rpcMock, tx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "nonce provider")
	})
}

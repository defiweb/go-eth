package txmodifier

import (
	"context"
	"errors"
	"fmt"

	"github.com/defiweb/go-eth/rpc"
	"github.com/defiweb/go-eth/types"
)

// NonceProvider is a transaction modifier that sets the nonce for the
// transaction.
//
// To use this modifier, add it using the WithTXModifiers option when creating
// a new rpc.Client.
type NonceProvider struct {
	usePendingBlock bool
	replace         bool
}

// NonceProviderOptions is the options for NewNonceProvider.
//
// If UsePendingBlock is true, then the next transaction nonce is fetched from
// the pending block. Otherwise, the next transaction nonce is fetched from the
// latest block. Using the pending block is not recommended as the behavior
// of the GetTransactionCount method on the pending block may be different
// between different Ethereum clients.
type NonceProviderOptions struct {
	UsePendingBlock bool // UsePendingBlock indicates whether to use the pending block.
	Replace         bool // Replace is true if the nonce should be replaced even if it is already set.
}

// NewNonceProvider returns a new NonceProvider.
func NewNonceProvider(opts NonceProviderOptions) *NonceProvider {
	return &NonceProvider{
		usePendingBlock: opts.UsePendingBlock,
		replace:         opts.Replace,
	}
}

// Modify implements the rpc.TXModifier interface.
func (p *NonceProvider) Modify(ctx context.Context, client rpc.RPC, tx *types.Transaction) error {
	if !p.replace && tx.Nonce != nil {
		return nil
	}
	if tx.From == nil {
		return errors.New("nonce provider: missing from address")
	}
	block := types.LatestBlockNumber
	if p.usePendingBlock {
		block = types.PendingBlockNumber
	}
	pendingNonce, err := client.GetTransactionCount(ctx, *tx.From, block)
	if err != nil {
		return fmt.Errorf("nonce provider: %w", err)
	}
	tx.Nonce = &pendingNonce
	return nil
}

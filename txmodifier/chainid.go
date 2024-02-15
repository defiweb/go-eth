package txmodifier

import (
	"context"
	"fmt"
	"sync"

	"github.com/defiweb/go-eth/rpc"
	"github.com/defiweb/go-eth/types"
)

// ChainIDProvider is a transaction modifier that sets the chain ID of the
// transaction.
//
// To use this modifier, add it using the WithTXModifiers option when creating
// a new rpc.Client.
type ChainIDProvider struct {
	mu      sync.Mutex
	chainID *uint64
	replace bool
	cache   bool
}

// ChainIDProviderOptions is the options for NewChainIDProvider.
type ChainIDProviderOptions struct {
	// Replace is true if the transaction chain ID should be replaced even if
	// it is already set.
	Replace bool

	// Cache is true if the chain ID will be cached instead of being queried
	// for each transaction. Cached chain ID will be used for all RPC clients
	// that use the same ChainIDProvider instance.
	Cache bool
}

// NewChainIDProvider returns a new ChainIDProvider.
func NewChainIDProvider(opts ChainIDProviderOptions) *ChainIDProvider {
	return &ChainIDProvider{
		replace: opts.Replace,
		cache:   opts.Cache,
	}
}

// Modify implements the rpc.TXModifier interface.
func (p *ChainIDProvider) Modify(ctx context.Context, client rpc.RPC, tx *types.Transaction) error {
	if !p.replace && tx.ChainID != nil {
		return nil
	}
	if !p.cache {
		chainID, err := client.ChainID(ctx)
		if err != nil {
			return fmt.Errorf("chain ID provider: %w", err)
		}
		tx.ChainID = &chainID
		return nil
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	var cid uint64
	if p.chainID != nil {
		cid = *p.chainID
	} else {
		chainID, err := client.ChainID(ctx)
		if err != nil {
			return fmt.Errorf("chain ID provider: %w", err)
		}
		p.chainID = &chainID
		cid = chainID
	}
	tx.ChainID = &cid
	return nil
}

package txmodifier

import (
	"context"
	"fmt"

	"github.com/defiweb/go-eth/rpc"
	"github.com/defiweb/go-eth/types"
)

// ChainIDProvider is a transaction modifier that sets the chain ID of the
// transaction.
//
// To use this modifier, add it using the WithTXModifiers option when creating
// a new rpc.Client.
type ChainIDProvider struct {
	chainID map[rpc.RPC]uint64
	replace bool
	cache   bool
}

// ChainIDProviderOptions is the options for NewChainIDProvider.
type ChainIDProviderOptions struct {
	Replace bool // Replace is true if the chain ID should be replaced even if it is already set.
	Cache   bool // Cache is true if the chain ID will be cached instead of being queried for each transaction.
}

// NewChainIDProvider returns a new ChainIDProvider.
func NewChainIDProvider(opts ChainIDProviderOptions) *ChainIDProvider {
	return &ChainIDProvider{
		chainID: make(map[rpc.RPC]uint64),
		replace: opts.Replace,
		cache:   opts.Cache,
	}
}

// Modify implements the rpc.TXModifier interface.
func (p *ChainIDProvider) Modify(ctx context.Context, client rpc.RPC, tx *types.Transaction) error {
	if !p.replace && tx.ChainID != nil {
		return nil
	}
	if chainID, ok := p.chainID[client]; ok {
		tx.ChainID = &chainID
		return nil
	}
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("chain ID provider: %w", err)
	}
	if p.cache {
		p.chainID[client] = chainID
	}
	tx.ChainID = &chainID
	return nil
}

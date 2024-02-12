package wallet

import (
	"context"

	"github.com/defiweb/go-eth/types"
)

// Key is the interface for an Ethereum key.
type Key interface {
	// Address returns the address of the key.
	Address() types.Address

	// SignMessage signs the given message.
	SignMessage(ctx context.Context, data []byte) (*types.Signature, error)

	// SignTransaction signs the given transaction.
	SignTransaction(ctx context.Context, tx *types.Transaction) error

	// VerifyMessage verifies whether the given data is signed by the key.
	VerifyMessage(ctx context.Context, data []byte, sig types.Signature) bool
}

// KeyWithHashSigner is the interface for an Ethereum key that can sign data using
// a private key, skipping the EIP-191 message prefix.
type KeyWithHashSigner interface {
	Key

	// SignHash signs the given hash without the EIP-191 message prefix.
	SignHash(ctx context.Context, hash types.Hash) (*types.Signature, error)

	// VerifyHash whether the given hash is signed by the key without the
	// EIP-191 message prefix.
	VerifyHash(ctx context.Context, hash types.Hash, sig types.Signature) bool
}

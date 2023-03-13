package crypto

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/defiweb/go-eth/types"
)

// Signer is an interface for signing data.
type Signer interface {
	// SignHash signs a hash.
	SignHash(hash types.Hash) (*types.Signature, error)

	// SignMessage signs a message.
	SignMessage(data []byte) (*types.Signature, error)

	// SignTransaction signs a transaction.
	SignTransaction(tx *types.Transaction) error
}

// Recoverer is an interface for recovering data.
type Recoverer interface {
	// RecoverHash recovers the address from a hash and signature.
	RecoverHash(hash types.Hash, sig types.Signature) (*types.Address, error)

	// RecoverMessage recovers the address from a message and signature.
	RecoverMessage(data []byte, sig types.Signature) (*types.Address, error)

	// RecoverTransaction recovers the address from a transaction.
	RecoverTransaction(tx *types.Transaction) (*types.Address, error)
}

// AddMessagePrefix adds the Ethereum message prefix to the given data.
func AddMessagePrefix(data []byte) []byte {
	return []byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data))
}

// ECSigner returns a Signer implementation for ECDSA.
func ECSigner(key *ecdsa.PrivateKey) Signer { return &ecSigner{key} }

// ECRecoverer is a Recoverer implementation for ECDSA.
var ECRecoverer Recoverer = &ecRecoverer{}

type (
	ecSigner    struct{ key *ecdsa.PrivateKey }
	ecRecoverer struct{}
)

func (s *ecSigner) SignHash(hash types.Hash) (*types.Signature, error) {
	return ecSignHash(s.key, hash)
}

func (s *ecSigner) SignMessage(data []byte) (*types.Signature, error) {
	return ecSignMessage(s.key, data)
}

func (s *ecSigner) SignTransaction(tx *types.Transaction) error {
	return ecSignTransaction(s.key, tx)
}

func (r *ecRecoverer) RecoverHash(hash types.Hash, sig types.Signature) (*types.Address, error) {
	return ecRecoverHash(hash, sig)
}

func (r *ecRecoverer) RecoverMessage(data []byte, sig types.Signature) (*types.Address, error) {
	return ecRecoverMessage(data, sig)
}

func (r *ecRecoverer) RecoverTransaction(tx *types.Transaction) (*types.Address, error) {
	return ecRecoverTransaction(tx)
}

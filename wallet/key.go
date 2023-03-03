package wallet

import "github.com/defiweb/go-eth/types"

type Key interface {
	// Address returns the address of the key.
	Address() types.Address

	// SignHash signs the given hash.
	SignHash(hash types.Hash) (types.Signature, error)

	// SignMessage signs the given message.
	SignMessage(data []byte) (types.Signature, error)

	// SignTransaction signs the given transaction.
	SignTransaction(tx *types.Transaction) error

	// VerifyHash whether the given hash is signed by the key.
	VerifyHash(hash types.Hash, sig types.Signature) bool

	// VerifyMessage verifies whether the given data is signed by the key.
	VerifyMessage(data []byte, sig types.Signature) bool
}

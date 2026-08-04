package types

import (
	"encoding/json"

	"github.com/defiweb/go-rlp"
)

// TransactionType is the type of transaction.
type TransactionType uint8

const (
	// LegacyTxType represents the legacy transaction format (Type 0).
	//
	// This is the original transaction format used before EIP-2718.
	LegacyTxType TransactionType = iota

	// AccessListTxType represents the access list transaction format (Type 1).
	//
	// Introduced by EIP-2930, this transaction type includes an optional
	// access list that specifies a list of addresses and storage keys the
	// transaction plans to access.
	AccessListTxType

	// DynamicFeeTxType represents the dynamic fee transaction format (Type 2).
	//
	// Introduced by EIP-1559, this transaction type supports a new fee market
	// mechanism with a base fee and a priority fee (tip).
	DynamicFeeTxType

	// BlobTxType represents the blob transaction format (Type 3).
	//
	// Introduced by EIP-4844, this transaction type adds support for
	// blob-carrying transactions.
	BlobTxType
)

// Transaction is an interface that represents a generic Ethereum transaction.
//
// Use NewTransaction* functions to create transactions of specific types.
type Transaction interface {
	json.Marshaler
	json.Unmarshaler
	rlp.Encoder
	rlp.Decoder

	// Type returns the type of the transaction.
	Type() TransactionType

	// Call returns the call associated with the transaction. The call is a
	// copy and can be modified. It may return nil if it is impossible to
	// create a call.
	Call() Call

	// Hash returns the hash of the transaction.
	Hash() (Hash, error)

	// Copy returns a deep copy of the transaction.
	Copy() Transaction
}

// SignableTransaction is an interface that represents a transaction that can
// be signed.
//
// Use NewTransaction* functions to create transactions of specific types.
type SignableTransaction interface {
	Transaction
	HasSigningData

	// SigningHash returns the hash used for signing the transaction.
	SigningHash() (Hash, error)
}

// TransactionDecoder is an interface for decoding transactions from JSON or
// RLP encoded data.
//
// Decoder may not set the From field of the transaction.
// To get a signer of the transaction, use txsign.Recover function.
type TransactionDecoder interface {
	RLPTransactionDecoder
	JSONTransactionDecoder
}

// RLPTransactionDecoder is an interface for decoding transactions from
// RLP-encoded data.
type RLPTransactionDecoder interface {
	// DecodeRLP decodes the RLP encoded transaction data.
	DecodeRLP(data []byte) (Transaction, error)
}

// JSONTransactionDecoder is an interface for decoding transactions from
// JSON-encoded data.
type JSONTransactionDecoder interface {
	// DecodeJSON decodes the JSON encoded transaction data.
	DecodeJSON(data []byte) (Transaction, error)
}

type jsonTransaction struct {
	ChainID *Number `json:"chainId,omitempty"`
	Nonce   *Number `json:"nonce,omitempty"`
	V       *Number `json:"v,omitempty"`
	R       *Number `json:"r,omitempty"`
	S       *Number `json:"s,omitempty"`

	jsonCall
}

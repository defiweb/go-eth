package types

import (
	"encoding/json"
	"fmt"
)

// DefaultTransactionDecoder is used to decode transactions when no other
// decoder is specified. Default implementation is used to decode Ethereum
// transactions.
var DefaultTransactionDecoder = &TypedTransactionDecoder{
	Types: map[TransactionType]func() Transaction{
		LegacyTxType:     func() Transaction { return NewTransactionLegacy() },
		AccessListTxType: func() Transaction { return NewTransactionAccessList() },
		DynamicFeeTxType: func() Transaction { return NewTransactionDynamicFee() },
		BlobTxType:       func() Transaction { return NewTransactionBlob() },
	},
	IgnoreUnknownTypes: true,
}

// TypedTransactionDecoder is am implementation of TransactionDecoder that
// could decode different types of transactions specified in the Types map.
type TypedTransactionDecoder struct {
	// Types is a map of transaction types to their constructors.
	Types map[TransactionType]func() Transaction

	// IgnoreUnknownTypes specifies whether to ignore unknown transaction types
	// or return an error.
	IgnoreUnknownTypes bool
}

// DecodeRLP implements the RLPTransactionDecoder interface.
func (e *TypedTransactionDecoder) DecodeRLP(data []byte) (Transaction, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty transaction data")
	}
	typ := TransactionType(data[0])
	if typ >= 0x80 {
		typ = LegacyTxType
	}
	tx := e.new(typ)
	if e.IgnoreUnknownTypes && tx == nil {
		return &TransactionUnknown{UnknownType: typ}, nil
	}
	if tx == nil {
		return nil, fmt.Errorf("unknown transaction type: %d", typ)
	}
	_, err := tx.DecodeRLP(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode transaction: %w", err)
	}
	return tx, nil
}

// DecodeJSON implements the JSONTransactionDecoder interface.
func (e *TypedTransactionDecoder) DecodeJSON(data []byte) (Transaction, error) {
	typ, err := jsonTXType(data)
	if err != nil {
		return nil, err
	}
	tx := e.new(typ)
	if e.IgnoreUnknownTypes && tx == nil {
		return &TransactionUnknown{UnknownType: typ}, nil
	}
	if tx == nil {
		return nil, fmt.Errorf("unknown transaction type: %d", typ)
	}
	if err := json.Unmarshal(data, tx); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction: %w", err)
	}
	return tx, nil
}

func (e *TypedTransactionDecoder) new(typ TransactionType) Transaction {
	if f, ok := e.Types[typ]; ok {
		return f()
	}
	return nil
}

// TransactionUnknown represent a transaction of unknown type.
//
// This type is returned by TransactionDecoder when it is unable to decode a
// transaction of a specific type.
type TransactionUnknown struct {
	UnknownType TransactionType
}

func (t *TransactionUnknown) Type() TransactionType { return t.UnknownType }

func (t *TransactionUnknown) Call() Call { return nil }

func (t *TransactionUnknown) Hash() (Hash, error) {
	return ZeroHash, fmt.Errorf("unable to calculate hash of unknown transaction type: %d", t.UnknownType)
}

func (t *TransactionUnknown) Copy() Transaction {
	return &TransactionUnknown{UnknownType: t.UnknownType}
}

func (t *TransactionUnknown) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("unable to marshal unknown transaction type: %d", t.UnknownType)
}

func (t *TransactionUnknown) UnmarshalJSON(_ []byte) error {
	return fmt.Errorf("unable to unmarshal unknown transaction type: %d", t.UnknownType)
}

func (t *TransactionUnknown) EncodeRLP() ([]byte, error) {
	return nil, fmt.Errorf("unable to encode unknown transaction type: %d", t.UnknownType)
}

func (t *TransactionUnknown) DecodeRLP(_ []byte) (int, error) {
	return 0, fmt.Errorf("unable to decode unknown transaction type: %d", t.UnknownType)
}

// jsonTXType returns the type of the transaction encoded in JSON.
//
// If the type is not specified, it tries to guess the type using the same rules as
// in go-ethereum code:
// https://github.com/ethereum/go-ethereum/blob/5b3e3cd2bee284db7d7deaa5986544d356410dcb/internal/ethapi/transaction_args.go#L472
func jsonTXType(data []byte) (TransactionType, error) {
	var tx struct {
		Type         *Number         `json:"type"`
		AccessList   *nilUnmarshaler `json:"accessList"`
		MaxFeePerGas *nilUnmarshaler `json:"maxFeePerGas"`
		BlobHashes   *nilUnmarshaler `json:"blobVersionedHashes"`
	}
	if err := json.Unmarshal(data, &tx); err != nil {
		return 0, fmt.Errorf("failed to unmarshal transaction: %w", err)
	}
	if tx.Type != nil {
		return TransactionType((*tx.Type).Big().Uint64()), nil
	}
	if tx.BlobHashes != nil {
		return BlobTxType, nil
	}
	if tx.MaxFeePerGas != nil {
		return DynamicFeeTxType, nil
	}
	if tx.AccessList != nil {
		return AccessListTxType, nil
	}
	return LegacyTxType, nil
}

type nilUnmarshaler struct{}

func (*nilUnmarshaler) UnmarshalJSON([]byte) error { return nil }

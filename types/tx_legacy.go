package types

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/defiweb/go-rlp"

	"github.com/defiweb/go-eth/crypto"
)

// TransactionLegacy is the legacy transaction type (Type 0).
//
// This is the original transaction format used before EIP-2718.
type TransactionLegacy struct {
	SigningData
	CallLegacy
}

// NewTransactionLegacy creates a new legacy transaction.
func NewTransactionLegacy() *TransactionLegacy {
	return &TransactionLegacy{}
}

// Type implements the Transaction interface.
func (t *TransactionLegacy) Type() TransactionType {
	return LegacyTxType
}

// Call implements the Transaction interface.
func (t *TransactionLegacy) Call() Call {
	return t.CallLegacy.Copy()
}

// Hash implements the Transaction interface.
func (t *TransactionLegacy) Hash() (Hash, error) {
	raw, err := t.EncodeRLP()
	if err != nil {
		return ZeroHash, err
	}
	return Hash(crypto.Keccak256(raw)), nil
}

// SigningHash implements the Transaction interface.
func (t *TransactionLegacy) SigningHash() (Hash, error) {
	var (
		chainID  = rlp.Uint(0)
		nonce    = rlp.Uint(0)
		gasPrice = &rlp.BigInt{}
		gasLimit = rlp.Uint(0)
		to       = (rlp.Bytes)(nil)
		value    = &rlp.BigInt{}
		input    = (rlp.Bytes)(nil)
	)
	if t.ChainID != nil {
		chainID = rlp.Uint(*t.ChainID)
	}
	if t.Nonce != nil {
		nonce = rlp.Uint(*t.Nonce)
	}
	if t.GasPrice != nil {
		gasPrice = (*rlp.BigInt)(t.GasPrice)
	}
	if t.GasLimit != nil {
		gasLimit = rlp.Uint(*t.GasLimit)
	}
	if t.To != nil {
		to = t.To[:]
	}
	if t.Value != nil {
		value = (*rlp.BigInt)(t.Value)
	}
	if t.Input != nil {
		input = t.Input
	}
	list := rlp.List{
		nonce,
		gasPrice,
		gasLimit,
		to,
		value,
		input,
	}
	if t.ChainID != nil && *t.ChainID != 0 {
		list.Add(
			chainID,
			rlp.Uint(0),
			rlp.Uint(0),
		)
	}
	bin, err := list.EncodeRLP()
	if err != nil {
		return ZeroHash, err
	}
	return Hash(crypto.Keccak256(bin)), nil
}

// Copy creates a deep copy of the transaction.
func (t *TransactionLegacy) Copy() Transaction {
	return &TransactionLegacy{
		SigningData: *t.SigningData.Copy(),
		CallLegacy:  *t.CallLegacy.Copy().(*CallLegacy),
	}
}

// EncodeRLP implements the rlp.Encoder interface.
//
//nolint:funlen
func (t TransactionLegacy) EncodeRLP() ([]byte, error) {
	var (
		nonce    = rlp.Uint(0)
		gasPrice = &rlp.BigInt{}
		gasLimit = rlp.Uint(0)
		to       = (rlp.Bytes)(nil)
		value    = &rlp.BigInt{}
		input    = (rlp.Bytes)(nil)
		v        = &rlp.BigInt{}
		r        = &rlp.BigInt{}
		s        = &rlp.BigInt{}
	)
	if t.Nonce != nil {
		nonce = rlp.Uint(*t.Nonce)
	}
	if t.GasPrice != nil {
		gasPrice = (*rlp.BigInt)(t.GasPrice)
	}
	if t.GasLimit != nil {
		gasLimit = rlp.Uint(*t.GasLimit)
	}
	if t.To != nil {
		to = t.To[:]
	}
	if t.Value != nil {
		value = (*rlp.BigInt)(t.Value)
	}
	if t.Input != nil {
		input = t.Input
	}
	if t.Signature != nil {
		v = (*rlp.BigInt)(t.Signature.V)
		r = (*rlp.BigInt)(t.Signature.R)
		s = (*rlp.BigInt)(t.Signature.S)
	}
	return rlp.List{
		nonce,
		gasPrice,
		gasLimit,
		to,
		value,
		input,
		v,
		r,
		s,
	}.EncodeRLP()
}

// DecodeRLP implements the rlp.Decoder interface.
//
//nolint:funlen
func (t *TransactionLegacy) DecodeRLP(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, fmt.Errorf("empty data")
	}
	var (
		nonce    = new(rlp.Uint)
		gasPrice = new(rlp.BigInt)
		gasLimit = new(rlp.Uint)
		to       = new(rlp.Bytes)
		value    = new(rlp.BigInt)
		input    = new(rlp.Bytes)
		v        = new(rlp.BigInt)
		r        = new(rlp.BigInt)
		s        = new(rlp.BigInt)
	)
	list := rlp.List{
		nonce,
		gasPrice,
		gasLimit,
		to,
		value,
		input,
		v,
		r,
		s,
	}
	if _, err := rlp.Decode(data, &list); err != nil {
		return 0, err
	}
	*t = TransactionLegacy{}
	if nonce.Get() != 0 {
		t.Nonce = nonce.Ptr()
	}
	if gasPrice.Get().Sign() != 0 {
		t.GasPrice = gasPrice.Ptr()
	}
	if gasLimit.Get() != 0 {
		t.GasLimit = gasLimit.Ptr()
	}
	if len(to.Get()) > 0 {
		t.To = AddressFromBytesPtr(*to)
	}
	if value.Get().Sign() != 0 {
		t.Value = value.Ptr()
	}
	if len(input.Get()) > 0 {
		t.Input = input.Get()
	}
	if v.Get().Sign() != 0 || r.Get().Sign() != 0 || s.Get().Sign() != 0 {
		t.Signature = &Signature{
			V: (*big.Int)(v),
			R: (*big.Int)(r),
			S: (*big.Int)(s),
		}
		// Derive chain ID from the V value.
		if v.Get().Cmp(big.NewInt(35)) >= 0 {
			x := new(big.Int).Sub((*big.Int)(v), big.NewInt(35))
			x = x.Div(x, big.NewInt(2))
			chainID := x.Uint64()
			t.ChainID = &chainID
		}
	}
	return len(data), nil
}

// MarshalJSON implements the json.Marshaler interface.
func (t *TransactionLegacy) MarshalJSON() ([]byte, error) {
	j := &jsonTransaction{}
	t.SigningData.toJSON(j)
	t.ExecutionData.toJSON(&j.jsonCall)
	t.LegacyFeeData.toJSON(&j.jsonCall)
	return json.Marshal(j)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *TransactionLegacy) UnmarshalJSON(data []byte) error {
	j := &jsonTransaction{}
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	t.SigningData.fromJSON(j)
	t.ExecutionData.fromJSON(&j.jsonCall)
	t.LegacyFeeData.fromJSON(&j.jsonCall)
	return nil
}

var _ Transaction = (*TransactionLegacy)(nil)

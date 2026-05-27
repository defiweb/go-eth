package types

import (
	"encoding/json"
	"fmt"

	"github.com/defiweb/go-rlp"

	"github.com/defiweb/go-eth/crypto"
)

// TransactionAccessList is the access list transaction type (Type 1).
//
// Introduced by EIP-2930, this transaction type includes an optional access
// list that specifies a list of addresses and storage keys the transaction
// plans to access.
type TransactionAccessList struct {
	SigningData
	CallAccessList
}

// NewTransactionAccessList creates a new access list transaction.
func NewTransactionAccessList() *TransactionAccessList {
	return &TransactionAccessList{}
}

// Type implements the Transaction interface.
func (t *TransactionAccessList) Type() TransactionType {
	return AccessListTxType
}

// Call implements the Transaction interface.
func (t *TransactionAccessList) Call() Call {
	return t.CallAccessList.Copy()
}

// Hash implements the Transaction interface.
func (t *TransactionAccessList) Hash() (Hash, error) {
	raw, err := t.EncodeRLP()
	if err != nil {
		return ZeroHash, err
	}
	return Hash(crypto.Keccak256(raw)), nil
}

// SigningHash implements the Transaction interface.
func (t *TransactionAccessList) SigningHash() (Hash, error) {
	var (
		chainID    = rlp.Uint(0)
		nonce      = rlp.Uint(0)
		gasPrice   = &rlp.BigInt{}
		gasLimit   = rlp.Uint(0)
		to         = (rlp.Bytes)(nil)
		value      = &rlp.BigInt{}
		input      = (rlp.Bytes)(nil)
		accessList = (AccessList)(nil)
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
	if t.AccessList != nil {
		accessList = t.AccessList
	}
	bin, err := rlp.List{
		chainID,
		nonce,
		gasPrice,
		gasLimit,
		to,
		value,
		input,
		&accessList,
	}.EncodeRLP()
	if err != nil {
		return ZeroHash, err
	}
	return Hash(crypto.Keccak256(append([]byte{byte(AccessListTxType)}, bin...))), nil
}

// EncodeRLP implements the rlp.Encoder interface.
//
//nolint:funlen
func (t TransactionAccessList) EncodeRLP() ([]byte, error) {
	var (
		chainID    = rlp.Uint(0)
		nonce      = rlp.Uint(0)
		gasPrice   = &rlp.BigInt{}
		gasLimit   = rlp.Uint(0)
		to         = (rlp.Bytes)(nil)
		value      = &rlp.BigInt{}
		input      = (rlp.Bytes)(nil)
		accessList = (AccessList)(nil)
		v          = &rlp.BigInt{}
		r          = &rlp.BigInt{}
		s          = &rlp.BigInt{}
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
	if t.AccessList != nil {
		accessList = t.AccessList
	}
	if t.Signature != nil {
		v = (*rlp.BigInt)(t.Signature.V)
		r = (*rlp.BigInt)(t.Signature.R)
		s = (*rlp.BigInt)(t.Signature.S)
	}
	bin, err := rlp.List{
		chainID,
		nonce,
		gasPrice,
		gasLimit,
		to,
		value,
		input,
		&accessList,
		v,
		r,
		s,
	}.EncodeRLP()
	if err != nil {
		return nil, err
	}
	return append([]byte{byte(AccessListTxType)}, bin...), nil
}

// Copy creates a deep copy of the transaction.
func (t *TransactionAccessList) Copy() Transaction {
	return &TransactionAccessList{
		SigningData:    *t.SigningData.Copy(),
		CallAccessList: *t.CallAccessList.Copy().(*CallAccessList),
	}
}

// DecodeRLP implements the rlp.Decoder interface.
//
//nolint:funlen
func (t *TransactionAccessList) DecodeRLP(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, fmt.Errorf("empty data")
	}
	if data[0] != byte(AccessListTxType) {
		return 0, fmt.Errorf("invalid transaction type: %d", data[0])
	}
	data = data[1:]
	var (
		chainID    = new(rlp.Uint)
		nonce      = new(rlp.Uint)
		gasPrice   = new(rlp.BigInt)
		gasLimit   = new(rlp.Uint)
		to         = new(rlp.Bytes)
		value      = new(rlp.BigInt)
		input      = new(rlp.Bytes)
		accessList = new(AccessList)
		v          = new(rlp.BigInt)
		r          = new(rlp.BigInt)
		s          = new(rlp.BigInt)
	)
	list := rlp.List{
		chainID,
		nonce,
		gasPrice,
		gasLimit,
		to,
		value,
		input,
		accessList,
		v,
		r,
		s,
	}
	if _, err := rlp.Decode(data, &list); err != nil {
		return 0, err
	}
	*t = TransactionAccessList{}
	if chainID.Get() != 0 {
		t.ChainID = chainID.Ptr()
	}
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
		t.To = AddressFromBytesPtr(to.Get())
	}
	if value.Ptr().Sign() != 0 {
		t.Value = value.Ptr()
	}
	if len(input.Get()) > 0 {
		t.Input = input.Get()
	}
	if len(*accessList) > 0 {
		t.AccessList = *accessList
	}
	if v.Ptr().Sign() != 0 || r.Ptr().Sign() != 0 || s.Ptr().Sign() != 0 {
		t.Signature = &Signature{
			V: v.Ptr(),
			R: r.Ptr(),
			S: s.Ptr(),
		}
		return len(data), nil
	}
	return len(data), nil
}

// MarshalJSON implements the json.Marshaler interface.
func (t *TransactionAccessList) MarshalJSON() ([]byte, error) {
	j := &jsonTransaction{}
	t.SigningData.toJSON(j)
	t.ExecutionData.toJSON(&j.jsonCall)
	t.LegacyFeeData.toJSON(&j.jsonCall)
	t.AccessListData.toJSON(&j.jsonCall)
	return json.Marshal(j)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *TransactionAccessList) UnmarshalJSON(data []byte) error {
	j := &jsonTransaction{}
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	t.SigningData.fromJSON(j)
	t.ExecutionData.fromJSON(&j.jsonCall)
	t.LegacyFeeData.fromJSON(&j.jsonCall)
	t.AccessListData.fromJSON(&j.jsonCall)
	return nil
}

var _ Transaction = (*TransactionAccessList)(nil)

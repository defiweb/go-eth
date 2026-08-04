package types

import (
	"encoding/json"
	"fmt"

	"github.com/defiweb/go-rlp"

	"github.com/defiweb/go-eth/crypto"
)

// TransactionDynamicFee is the dynamic fee transaction type (Type 2).
//
// Introduced by EIP-1559, this transaction type supports a new fee market
// mechanism with a base fee and a priority fee (tip).
type TransactionDynamicFee struct {
	SigningData
	CallDynamicFee
}

// NewTransactionDynamicFee creates a new dynamic fee transaction.
func NewTransactionDynamicFee() *TransactionDynamicFee {
	return &TransactionDynamicFee{}
}

// Type implements the Transaction interface.
func (t *TransactionDynamicFee) Type() TransactionType {
	return DynamicFeeTxType
}

// Call implements the Transaction interface.
func (t *TransactionDynamicFee) Call() Call {
	return t.CallDynamicFee.Copy()
}

// Hash implements the Transaction interface.
func (t *TransactionDynamicFee) Hash() (Hash, error) {
	raw, err := t.EncodeRLP()
	if err != nil {
		return ZeroHash, err
	}
	return Hash(crypto.Keccak256(raw)), nil
}

// SigningHash implements the Transaction interface.
func (t *TransactionDynamicFee) SigningHash() (Hash, error) {
	var (
		chainID              = rlp.Uint(0)
		nonce                = rlp.Uint(0)
		gasLimit             = rlp.Uint(0)
		maxPriorityFeePerGas = &rlp.BigInt{}
		maxFeePerGas         = &rlp.BigInt{}
		to                   = (rlp.Bytes)(nil)
		value                = &rlp.BigInt{}
		input                = (rlp.Bytes)(nil)
		accessList           = (AccessList)(nil)
	)
	if t.ChainID != nil {
		chainID = rlp.Uint(*t.ChainID)
	}
	if t.Nonce != nil {
		nonce = rlp.Uint(*t.Nonce)
	}
	if t.GasLimit != nil {
		gasLimit = rlp.Uint(*t.GasLimit)
	}
	if t.MaxPriorityFeePerGas != nil {
		maxPriorityFeePerGas = (*rlp.BigInt)(t.MaxPriorityFeePerGas)
	}
	if t.MaxFeePerGas != nil {
		maxFeePerGas = (*rlp.BigInt)(t.MaxFeePerGas)
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
		maxPriorityFeePerGas,
		maxFeePerGas,
		gasLimit,
		to,
		value,
		input,
		&accessList,
	}.EncodeRLP()
	if err != nil {
		return ZeroHash, err
	}
	return Hash(crypto.Keccak256(append([]byte{byte(DynamicFeeTxType)}, bin...))), nil
}

// Copy creates a deep copy of the transaction.
func (t *TransactionDynamicFee) Copy() Transaction {
	return &TransactionDynamicFee{
		SigningData:    *t.SigningData.Copy(),
		CallDynamicFee: *t.CallDynamicFee.Copy().(*CallDynamicFee),
	}
}

// EncodeRLP implements the rlp.Encoder interface.
//
//nolint:funlen
func (t TransactionDynamicFee) EncodeRLP() ([]byte, error) {
	var (
		chainID              = rlp.Uint(0)
		nonce                = rlp.Uint(0)
		gasLimit             = rlp.Uint(0)
		maxPriorityFeePerGas = &rlp.BigInt{}
		maxFeePerGas         = &rlp.BigInt{}
		to                   = (rlp.Bytes)(nil)
		value                = &rlp.BigInt{}
		input                = (rlp.Bytes)(nil)
		accessList           = (AccessList)(nil)
		v                    = &rlp.BigInt{}
		r                    = &rlp.BigInt{}
		s                    = &rlp.BigInt{}
	)
	if t.ChainID != nil {
		chainID = rlp.Uint(*t.ChainID)
	}
	if t.Nonce != nil {
		nonce = rlp.Uint(*t.Nonce)
	}
	if t.GasLimit != nil {
		gasLimit = rlp.Uint(*t.GasLimit)
	}
	if t.MaxPriorityFeePerGas != nil {
		maxPriorityFeePerGas = (*rlp.BigInt)(t.MaxPriorityFeePerGas)
	}
	if t.MaxFeePerGas != nil {
		maxFeePerGas = (*rlp.BigInt)(t.MaxFeePerGas)
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
		maxPriorityFeePerGas,
		maxFeePerGas,
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
	return append([]byte{byte(DynamicFeeTxType)}, bin...), nil
}

// DecodeRLP implements the rlp.Decoder interface.
//
//nolint:funlen
func (t *TransactionDynamicFee) DecodeRLP(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, fmt.Errorf("empty data")
	}
	if data[0] != byte(DynamicFeeTxType) {
		return 0, fmt.Errorf("invalid transaction type: %d", data[0])
	}
	data = data[1:]
	var (
		chainID              = new(rlp.Uint)
		nonce                = new(rlp.Uint)
		gasLimit             = new(rlp.Uint)
		maxPriorityFeePerGas = new(rlp.BigInt)
		maxFeePerGas         = new(rlp.BigInt)
		to                   = new(rlp.Bytes)
		value                = new(rlp.BigInt)
		input                = new(rlp.Bytes)
		accessList           = new(AccessList)
		v                    = new(rlp.BigInt)
		r                    = new(rlp.BigInt)
		s                    = new(rlp.BigInt)
	)
	list := rlp.List{
		chainID,
		nonce,
		maxPriorityFeePerGas,
		maxFeePerGas,
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
	*t = TransactionDynamicFee{}
	if chainID.Get() != 0 {
		t.ChainID = chainID.Ptr()
	}
	if nonce.Get() != 0 {
		t.Nonce = nonce.Ptr()
	}
	if maxPriorityFeePerGas.Ptr().Sign() != 0 {
		t.MaxPriorityFeePerGas = maxPriorityFeePerGas.Ptr()
	}
	if maxFeePerGas.Ptr().Sign() != 0 {
		t.MaxFeePerGas = maxFeePerGas.Ptr()
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
	}
	return len(data), nil
}

// MarshalJSON implements the json.Marshaler interface.
func (t *TransactionDynamicFee) MarshalJSON() ([]byte, error) {
	j := &jsonTransaction{}
	t.SigningData.toJSON(j)
	t.ExecutionData.toJSON(&j.jsonCall)
	t.AccessListData.toJSON(&j.jsonCall)
	t.DynamicFeeData.toJSON(&j.jsonCall)
	return json.Marshal(j)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *TransactionDynamicFee) UnmarshalJSON(data []byte) error {
	j := &jsonTransaction{}
	if err := json.Unmarshal(data, j); err != nil {
		return err
	}
	t.SigningData.fromJSON(j)
	t.ExecutionData.fromJSON(&j.jsonCall)
	t.AccessListData.fromJSON(&j.jsonCall)
	t.DynamicFeeData.fromJSON(&j.jsonCall)
	return nil
}

var _ Transaction = (*TransactionDynamicFee)(nil)

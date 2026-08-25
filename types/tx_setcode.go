package types

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/defiweb/go-rlp"

	"github.com/defiweb/go-eth/crypto"
)

// TransactionSetCode is the set code transaction type (Type 4).
//
// Introduced by EIP-7702, this transaction type allows EOAs to temporarily
// adopt code from a smart contract by providing a list of authorization tuples.
type TransactionSetCode struct {
	SigningData
	CallSetCode
}

// NewTransactionSetCode creates a new set code transaction.
func NewTransactionSetCode() *TransactionSetCode {
	return &TransactionSetCode{}
}

// Type implements the Transaction interface.
func (t *TransactionSetCode) Type() TransactionType {
	return SetCodeTxType
}

// Call implements the Transaction interface.
func (t *TransactionSetCode) Call() Call {
	return t.CallSetCode.Copy()
}

// Hash implements the Transaction interface.
func (t *TransactionSetCode) Hash() (Hash, error) {
	raw, err := t.EncodeRLP()
	if err != nil {
		return ZeroHash, err
	}
	return Hash(crypto.Keccak256(raw)), nil
}

// SigningHash implements the Transaction interface.
func (t *TransactionSetCode) SigningHash() (Hash, error) {
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
		authorizationList    = (AuthorizationList)(nil)
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
	if t.AuthorizationList != nil {
		authorizationList = t.AuthorizationList
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
		&authorizationList,
	}.EncodeRLP()
	if err != nil {
		return ZeroHash, err
	}
	return Hash(crypto.Keccak256(append([]byte{byte(SetCodeTxType)}, bin...))), nil
}

// Sign signs the transaction.
func (t *TransactionSetCode) Sign(ctx context.Context, key TransactionSigner) error {
	return key.SignTransaction(ctx, t)
}

// Copy creates a deep copy of the transaction.
func (t *TransactionSetCode) Copy() Transaction {
	return &TransactionSetCode{
		SigningData: *t.SigningData.Copy(),
		CallSetCode: *t.CallSetCode.Copy().(*CallSetCode),
	}
}

// EncodeRLP implements the rlp.Encoder interface.
//
//nolint:funlen
func (t TransactionSetCode) EncodeRLP() ([]byte, error) {
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
		authorizationList    = (AuthorizationList)(nil)
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
	if t.AuthorizationList != nil {
		authorizationList = t.AuthorizationList
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
		&authorizationList,
		v,
		r,
		s,
	}.EncodeRLP()
	if err != nil {
		return nil, err
	}
	return append([]byte{byte(SetCodeTxType)}, bin...), nil
}

// DecodeRLP implements the rlp.Decoder interface.
//
//nolint:funlen
func (t *TransactionSetCode) DecodeRLP(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, fmt.Errorf("empty data")
	}
	if data[0] != byte(SetCodeTxType) {
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
		authorizationList    = new(AuthorizationList)
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
		authorizationList,
		v,
		r,
		s,
	}
	if _, err := rlp.Decode(data, &list); err != nil {
		return 0, err
	}
	*t = TransactionSetCode{}
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
	if len(*authorizationList) > 0 {
		t.AuthorizationList = *authorizationList
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
func (t *TransactionSetCode) MarshalJSON() ([]byte, error) {
	j := &jsonTransaction{}
	t.SigningData.toJSON(j)
	t.ExecutionData.toJSON(&j.jsonCall)
	t.AccessListData.toJSON(&j.jsonCall)
	t.DynamicFeeData.toJSON(&j.jsonCall)
	t.AuthorizationData.toJSON(&j.jsonCall)
	return json.Marshal(j)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *TransactionSetCode) UnmarshalJSON(data []byte) error {
	j := &jsonTransaction{}
	if err := json.Unmarshal(data, j); err != nil {
		return err
	}
	t.SigningData.fromJSON(j)
	t.ExecutionData.fromJSON(&j.jsonCall)
	t.AccessListData.fromJSON(&j.jsonCall)
	t.DynamicFeeData.fromJSON(&j.jsonCall)
	t.AuthorizationData.fromJSON(&j.jsonCall)
	return nil
}

var _ Transaction = (*TransactionSetCode)(nil)

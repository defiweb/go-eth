package types

import (
	"encoding/json"
	"fmt"

	"github.com/defiweb/go-rlp"

	"github.com/defiweb/go-eth/crypto"
)

// TransactionBlob is the blob transaction type (Type 3).
//
// Introduced by EIP-4844, this transaction type adds support for blob-carrying
// transactions.
type TransactionBlob struct {
	SigningData
	CallBlob
}

// NewTransactionBlob creates a new blob transaction.
func NewTransactionBlob() *TransactionBlob {
	return &TransactionBlob{}
}

// Type implements the Transaction interface.
func (t *TransactionBlob) Type() TransactionType {
	return BlobTxType
}

// Call implements the Transaction interface.
func (t *TransactionBlob) Call() Call {
	return t.CallBlob.Copy()
}

// Hash implements the Transaction interface.
func (t *TransactionBlob) Hash() (Hash, error) {
	raw, err := t.EncodeRLP()
	if err != nil {
		return ZeroHash, err
	}
	return Hash(crypto.Keccak256(raw)), nil
}

// SigningHash implements the Transaction interface.
func (t *TransactionBlob) SigningHash() (Hash, error) {
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
		maxFeePerBlobGas     = &rlp.BigInt{}
		blobHashes           = (rlp.TypedList[kzgHash])(nil)
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
	if t.MaxFeePerBlobGas != nil {
		maxFeePerBlobGas = (*rlp.BigInt)(t.MaxFeePerBlobGas)
	}
	if len(t.Blobs) > 0 {
		blobHashes = make(rlp.TypedList[kzgHash], len(t.Blobs))
		for i := range t.Blobs {
			blobHashes[i] = (*kzgHash)(&t.Blobs[i].Hash)
		}
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
		maxFeePerBlobGas,
		&blobHashes,
	}.EncodeRLP()
	if err != nil {
		return ZeroHash, err
	}
	return Hash(crypto.Keccak256(append([]byte{byte(BlobTxType)}, bin...))), nil
}

// Copy creates a deep copy of the transaction.
func (t *TransactionBlob) Copy() Transaction {
	return &TransactionBlob{
		SigningData: *t.SigningData.Copy(),
		CallBlob:    *t.CallBlob.Copy().(*CallBlob),
	}
}

// EncodeRLP implements the rlp.Encoder interface.
//
//nolint:funlen
func (t TransactionBlob) EncodeRLP() ([]byte, error) {
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
		maxFeePerBlobGas     = &rlp.BigInt{}
		blobHashes           = (rlp.TypedList[kzgHash])(nil)
		blobs                = (rlp.TypedList[kzgBlob])(nil)
		commitments          = (rlp.TypedList[kzgCommitment])(nil)
		proofs               = (rlp.TypedList[kzgProof])(nil)
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
	if t.MaxFeePerBlobGas != nil {
		maxFeePerBlobGas = (*rlp.BigInt)(t.MaxFeePerBlobGas)
	}
	if len(t.Blobs) > 0 {
		blobHashes = make(rlp.TypedList[kzgHash], 0, len(t.Blobs))
		for i := range t.Blobs {
			blob := t.Blobs[i]
			blobHashes = append(blobHashes, (*kzgHash)(&blob.Hash))
			if blob.Sidecar != nil {
				blobs.Add((*kzgBlob)(&blob.Sidecar.Blob))
				commitments.Add((*kzgCommitment)(&blob.Sidecar.Commitment))
				proofs.Add((*kzgProof)(&blob.Sidecar.Proof))
			}
		}
	}
	if t.Signature != nil {
		v = (*rlp.BigInt)(t.Signature.V)
		r = (*rlp.BigInt)(t.Signature.R)
		s = (*rlp.BigInt)(t.Signature.S)
	}
	tx := rlp.List{
		chainID,
		nonce,
		maxPriorityFeePerGas,
		maxFeePerGas,
		gasLimit,
		to,
		value,
		input,
		&accessList,
		maxFeePerBlobGas,
		&blobHashes,
		v,
		r,
		s,
	}
	if len(blobHashes) > 0 && len(blobHashes) == len(blobs) {
		tx = rlp.List{
			tx,
			blobs,
			commitments,
			proofs,
		}
	}
	bin, err := tx.EncodeRLP()
	if err != nil {
		return nil, err
	}
	return append([]byte{byte(BlobTxType)}, bin...), nil
}

// DecodeRLP implements the rlp.Decoder interface.
//
//nolint:funlen
func (t *TransactionBlob) DecodeRLP(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, fmt.Errorf("empty data")
	}
	if data[0] != byte(BlobTxType) {
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
		maxFeePerBlobGas     = new(rlp.BigInt)
		blobHashes           = new(rlp.VarTypedList[Hash])
		blobs                = new(rlp.VarTypedList[kzgBlob])
		commitments          = new(rlp.VarTypedList[kzgCommitment])
		proofs               = new(rlp.VarTypedList[kzgProof])
		v                    = new(rlp.BigInt)
		r                    = new(rlp.BigInt)
		s                    = new(rlp.BigInt)
	)
	dec, n, err := rlp.DecodeLazy(data)
	if err != nil {
		return 0, err
	}
	if n != len(data) {
		return 0, rlp.ErrUnexpectedTrailingData
	}
	if !dec.IsList() {
		return 0, fmt.Errorf("unable to decode transaction")
	}
	txFields := rlp.List{
		chainID,
		nonce,
		maxPriorityFeePerGas,
		maxFeePerGas,
		gasLimit,
		to,
		value,
		input,
		accessList,
		maxFeePerBlobGas,
		blobHashes,
		v,
		r,
		s,
	}
	var list rlp.List
	switch dec.Length() {
	case 4:
		list = rlp.List{&txFields, blobs, commitments, proofs}
	case 14:
		list = txFields
	default:
		return 0, fmt.Errorf("invalid transaction: expected 14 or 4 fields, got %d", dec.Length())
	}
	if err := dec.Decode(&list); err != nil {
		return 0, err
	}
	*t = TransactionBlob{}
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
	if maxFeePerBlobGas.Ptr().Sign() != 0 {
		t.MaxFeePerBlobGas = maxFeePerBlobGas.Ptr()
	}
	if len(*blobHashes) > 0 {
		t.Blobs = make([]BlobInfo, len(*blobHashes))
		for i, hash := range *blobHashes {
			blob := BlobInfo{Hash: crypto.KZGHash(*hash)}
			if i < len(*blobs) && i < len(*commitments) && i < len(*proofs) {
				blob.Sidecar = &BlobSidecar{
					Blob:       crypto.KZGBlob(*(*blobs)[i]),
					Commitment: crypto.KZGCommitment(*(*commitments)[i]),
					Proof:      crypto.KZGProof(*(*proofs)[i]),
				}
			}
			t.Blobs[i] = blob
		}
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
func (t *TransactionBlob) MarshalJSON() ([]byte, error) {
	j := &jsonTransaction{}
	t.SigningData.toJSON(j)
	t.ExecutionData.toJSON(&j.jsonCall)
	t.AccessListData.toJSON(&j.jsonCall)
	t.DynamicFeeData.toJSON(&j.jsonCall)
	t.BlobData.toJSON(&j.jsonCall)
	return json.Marshal(j)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *TransactionBlob) UnmarshalJSON(data []byte) error {
	j := &jsonTransaction{}
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	t.SigningData.fromJSON(j)
	t.ExecutionData.fromJSON(&j.jsonCall)
	t.AccessListData.fromJSON(&j.jsonCall)
	t.DynamicFeeData.fromJSON(&j.jsonCall)
	t.BlobData.fromJSON(&j.jsonCall)
	return nil
}

var _ Transaction = (*TransactionBlob)(nil)

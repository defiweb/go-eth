package types

import (
	"math/big"

	"github.com/defiweb/go-eth/crypto/kzg4844"
)

// Below types are used to embed common fields and methods into call and
// transaction types. You probably do not want to use them directly.
//
// Below interfaces are used to determine if a given call or transaction has
// specific capabilities.

// TransactionData defines methods for accessing and setting transaction
// data.
type TransactionData interface {
	TransactionData() *TransactionFields
	SetTransactionData(data TransactionFields)
}

// CallData defines methods for accessing and setting call data.
type CallData interface {
	CallData() *CallFields
	SetCallData(data CallFields)
}

// LegacyPriceData defines methods for accessing and setting legacy price data.
type LegacyPriceData interface {
	LegacyPriceData() *LegacyPriceField
	SetLegacyPriceData(data LegacyPriceField)
}

// AccessListData defines methods for accessing and setting access list
// data.
type AccessListData interface {
	AccessListData() *AccessListField
	SetAccessListData(data AccessListField)
}

// DynamicFeeData defines methods for accessing and setting dynamic fee`1
// data.
type DynamicFeeData interface {
	DynamicFeeData() *DynamicFeeFields
	SetDynamicFeeData(data DynamicFeeFields)
}

// BlobData defines methods for accessing and setting blob data.
type BlobData interface {
	BlobData() *BlobFields
	SetBlobData(data BlobFields)
}

// TransactionFields contains common fields for transactions.
//
// This type is used to embed transaction data into other types.
type TransactionFields struct {
	ChainID   *uint64    // ChainID is the chain ID.
	Nonce     *uint64    // Nonce is the transaction nonce.
	Signature *Signature // Signature is the transaction signature.
}

// TransactionData returns the embedded transaction data.
func (c *TransactionFields) TransactionData() *TransactionFields {
	return c
}

// SetTransactionData sets the embedded transaction data.
func (c *TransactionFields) SetTransactionData(data TransactionFields) {
	*c = data
}

// SetChainID sets the chain ID.
func (c *TransactionFields) SetChainID(chainID uint64) {
	c.ChainID = &chainID
}

// SetNonce sets the transaction nonce.
func (c *TransactionFields) SetNonce(nonce uint64) {
	c.Nonce = &nonce
}

// SetSignature sets the transaction signature.
func (c *TransactionFields) SetSignature(signature Signature) {
	c.Signature = &signature
}

// Copy creates a deep copy of the TransactionFields.
func (c *TransactionFields) Copy() *TransactionFields {
	return &TransactionFields{
		ChainID:   copyPtr(c.ChainID),
		Nonce:     copyPtr(c.Nonce),
		Signature: c.Signature.Copy(),
	}
}

func (c *TransactionFields) toJSON(j *jsonTransaction) {
	if c.ChainID != nil {
		j.ChainID = NumberFromUint64Ptr(*c.ChainID)
	}
	if c.Nonce != nil {
		j.Nonce = NumberFromUint64Ptr(*c.Nonce)
	}
	if c.Signature != nil {
		j.V = NumberFromBigIntPtr(c.Signature.V)
		j.R = NumberFromBigIntPtr(c.Signature.R)
		j.S = NumberFromBigIntPtr(c.Signature.S)
	}
}

func (c *TransactionFields) fromJSON(j *jsonTransaction) {
	if j.ChainID != nil {
		chainID := j.ChainID.Big().Uint64()
		c.ChainID = &chainID
	}
	if j.Nonce != nil {
		nonce := j.Nonce.Big().Uint64()
		c.Nonce = &nonce
	}
	if j.V != nil || j.R != nil || j.S != nil {
		c.Signature = SignatureFromVRSPtr(j.V.Big(), j.R.Big(), j.S.Big())
	}
}

// CallFields contains the basic fields for a call.
//
// This type is used to embed call data into other types.
type CallFields struct {
	From     *Address // From is the sender address.
	To       *Address // To is the recipient address. Nil means contract creation.
	GasLimit *uint64  // GasLimit is the gas limit; if 0, there is no limit.
	Value    *big.Int // Value is the amount of wei to send.
	Input    []byte   // Input is the input data.
}

// CallData returns the embedded call data.
func (c *CallFields) CallData() *CallFields {
	return c
}

// SetCallData sets the embedded call data.
func (c *CallFields) SetCallData(data CallFields) {
	*c = data
}

// SetFrom sets the sender address.
func (c *CallFields) SetFrom(from Address) {
	c.From = &from
}

// SetTo sets the recipient address.
func (c *CallFields) SetTo(to Address) {
	c.To = &to
}

// SetGasLimit sets the gas limit.
func (c *CallFields) SetGasLimit(gasLimit uint64) {
	c.GasLimit = &gasLimit
}

// SetValue sets the amount of wei to send.
func (c *CallFields) SetValue(value *big.Int) {
	c.Value = value
}

// SetInput sets the input data.
func (c *CallFields) SetInput(input []byte) {
	c.Input = input
}

// Copy creates a deep copy of the CallFields.
func (c *CallFields) Copy() *CallFields {
	if c == nil {
		return nil
	}
	return &CallFields{
		From:     copyPtr(c.From),
		To:       copyPtr(c.To),
		GasLimit: copyPtr(c.GasLimit),
		Value:    copyBigInt(c.Value),
		Input:    copyBytes(c.Input),
	}
}

func (c *CallFields) toJSON(j *jsonCall) {
	j.From = c.From
	j.To = c.To
	if c.GasLimit != nil {
		j.GasLimit = NumberFromUint64Ptr(*c.GasLimit)
	}
	if c.Value != nil {
		j.Value = NumberFromBigIntPtr(c.Value)
	}
	j.Input = c.Input
}

func (c *CallFields) fromJSON(j *jsonCall) {
	c.From = j.From
	c.To = j.To
	if j.GasLimit != nil {
		gas := j.GasLimit.Big().Uint64()
		c.GasLimit = &gas
	}
	if j.Value != nil {
		c.Value = j.Value.Big()
	}
	c.Input = j.Input
}

// LegacyPriceField contains the gas price for legacy transactions.
//
// This type is used to embed legacy price data into other types.
type LegacyPriceField struct {
	GasPrice *big.Int // GasPrice is the gas price.
}

// LegacyPriceData returns the embedded legacy price data.
func (c *LegacyPriceField) LegacyPriceData() *LegacyPriceField {
	return c
}

// SetLegacyPriceData sets the embedded legacy price data.
func (c *LegacyPriceField) SetLegacyPriceData(data LegacyPriceField) {
	*c = data
}

// SetGasPrice sets the gas price.
func (c *LegacyPriceField) SetGasPrice(gasPrice *big.Int) {
	c.GasPrice = gasPrice
}

// Copy creates a deep copy of the LegacyPriceField.
func (c *LegacyPriceField) Copy() *LegacyPriceField {
	if c == nil {
		return nil
	}
	return &LegacyPriceField{
		GasPrice: copyBigInt(c.GasPrice),
	}
}

func (c *LegacyPriceField) toJSON(j *jsonCall) {
	if c.GasPrice != nil {
		j.GasPrice = NumberFromBigIntPtr(c.GasPrice)
	}
}

func (c *LegacyPriceField) fromJSON(j *jsonCall) {
	if j.GasPrice != nil {
		c.GasPrice = j.GasPrice.Big()
	}
}

// AccessListField contains the access list for EIP-2930 transactions.
//
// This type is used to embed access list data into other types.
type AccessListField struct {
	AccessList AccessList // AccessList is the EIP-2930 access list.
}

// AccessListData returns the embedded access list data.
func (c *AccessListField) AccessListData() *AccessListField {
	return c
}

// SetAccessListData sets the embedded access list data.
func (c *AccessListField) SetAccessListData(data AccessListField) {
	*c = data
}

// SetAccessList sets the access list.
func (c *AccessListField) SetAccessList(accessList AccessList) {
	c.AccessList = accessList
}

// Copy creates a deep copy of the AccessListField.
func (c *AccessListField) Copy() *AccessListField {
	if c == nil {
		return nil
	}
	return &AccessListField{
		AccessList: c.AccessList.Copy(),
	}
}

func (c *AccessListField) toJSON(j *jsonCall) {
	j.AccessList = c.AccessList
}

func (c *AccessListField) fromJSON(j *jsonCall) {
	c.AccessList = j.AccessList
}

// DynamicFeeFields contains fee data for EIP-1559 transactions.
//
// This type is used to embed dynamic fee data into other types.
type DynamicFeeFields struct {
	MaxFeePerGas         *big.Int // MaxFeePerGas is the maximum total fee per gas.
	MaxPriorityFeePerGas *big.Int // MaxPriorityFeePerGas is the maximum priority fee per gas.
}

// DynamicFeeData returns the embedded dynamic fee data.
func (c *DynamicFeeFields) DynamicFeeData() *DynamicFeeFields {
	return c
}

// SetDynamicFeeData sets the embedded dynamic fee data.
func (c *DynamicFeeFields) SetDynamicFeeData(data DynamicFeeFields) {
	*c = data
}

// SetMaxFeePerGas sets the maximum total fee per gas.
func (c *DynamicFeeFields) SetMaxFeePerGas(maxFeePerGas *big.Int) {
	c.MaxFeePerGas = maxFeePerGas
}

// SetMaxPriorityFeePerGas sets the maximum priority fee per gas.
func (c *DynamicFeeFields) SetMaxPriorityFeePerGas(maxPriorityFeePerGas *big.Int) {
	c.MaxPriorityFeePerGas = maxPriorityFeePerGas
}

// Copy creates a deep copy of the DynamicFeeFields.
func (c *DynamicFeeFields) Copy() *DynamicFeeFields {
	if c == nil {
		return nil
	}
	return &DynamicFeeFields{
		MaxFeePerGas:         copyBigInt(c.MaxFeePerGas),
		MaxPriorityFeePerGas: copyBigInt(c.MaxPriorityFeePerGas),
	}
}

func (c *DynamicFeeFields) toJSON(j *jsonCall) {
	if c.MaxFeePerGas != nil {
		j.MaxFeePerGas = NumberFromBigIntPtr(c.MaxFeePerGas)
	}
	if c.MaxPriorityFeePerGas != nil {
		j.MaxPriorityFeePerGas = NumberFromBigIntPtr(c.MaxPriorityFeePerGas)
	}
}

func (c *DynamicFeeFields) fromJSON(j *jsonCall) {
	if j.MaxFeePerGas != nil {
		c.MaxFeePerGas = j.MaxFeePerGas.Big()
	}
	if j.MaxPriorityFeePerGas != nil {
		c.MaxPriorityFeePerGas = j.MaxPriorityFeePerGas.Big()
	}
}

// BlobFields contains data for EIP-4844 blob transactions.
//
// Use NewBlobInfo to create a BlobInfo.
//
// This type is used to embed blob data into other types.
type BlobFields struct {
	MaxFeePerBlobGas *big.Int   // MaxFeePerBlobGas is the maximum fee per blob gas.
	Blobs            []BlobInfo // Blobs is the list of blobs.
}

// BlobData returns the embedded blob data.
func (c *BlobFields) BlobData() *BlobFields {
	return c
}

// SetBlobData sets the embedded blob data.
func (c *BlobFields) SetBlobData(data BlobFields) {
	*c = data
}

// SetMaxFeePerBlobGas sets the maximum fee per blob gas.
func (c *BlobFields) SetMaxFeePerBlobGas(maxFeePerBlobGas *big.Int) {
	c.MaxFeePerBlobGas = maxFeePerBlobGas
}

// SetBlobs sets the list of blobs.
//
// Use NewBlobInfo to create a BlobInfo.
func (c *BlobFields) SetBlobs(blobs []BlobInfo) {
	c.Blobs = blobs
}

// AddBlob adds a blob to the list.
//
// Use NewBlobInfo to create a BlobInfo.
func (c *BlobFields) AddBlob(blob BlobInfo) {
	c.Blobs = append(c.Blobs, blob)
}

// Copy creates a deep copy of the BlobFields.
func (c *BlobFields) Copy() *BlobFields {
	if c == nil {
		return nil
	}
	blobs := make([]BlobInfo, len(c.Blobs))
	for i, blob := range c.Blobs {
		blobs[i].Hash = blob.Hash
		blobs[i].Sidecar = copyPtr(blob.Sidecar)
	}
	return &BlobFields{
		MaxFeePerBlobGas: copyBigInt(c.MaxFeePerBlobGas),
		Blobs:            blobs,
	}
}

func (c *BlobFields) toJSON(j *jsonCall) {
	if c.MaxFeePerBlobGas != nil {
		j.MaxFeePerBlobGas = NumberFromBigIntPtr(c.MaxFeePerBlobGas)
	}
	if len(c.Blobs) > 0 && c.Blobs[0].Sidecar != nil {
		// If the first blob has a sidecar, then all blobs should have
		// sidecars, so we can allocate memory for them.
		j.BlobHashes = make([]Hash, 0, len(c.Blobs))
		j.Blobs = make([]kzgBlob, 0, len(c.Blobs))
		j.Commitments = make([]kzgCommitment, 0, len(c.Blobs))
		j.Proofs = make([]kzgProof, 0, len(c.Blobs))
	}
	for _, b := range c.Blobs {
		j.BlobHashes = append(j.BlobHashes, b.Hash)
		if b.Sidecar != nil {
			j.Blobs = append(j.Blobs, kzgBlob(b.Sidecar.Blob))
			j.Commitments = append(j.Commitments, kzgCommitment(b.Sidecar.Commitment))
			j.Proofs = append(j.Proofs, kzgProof(b.Sidecar.Proof))
		}
	}
}

func (c *BlobFields) fromJSON(j *jsonCall) {
	if j.MaxFeePerBlobGas != nil {
		c.MaxFeePerBlobGas = j.MaxFeePerBlobGas.Big()
	}
	if len(j.BlobHashes) > 0 {
		c.Blobs = make([]BlobInfo, len(j.BlobHashes))
		for i, h := range j.BlobHashes {
			b := BlobInfo{Hash: h}
			if i < len(j.Blobs) && i < len(j.Commitments) && i < len(j.Proofs) {
				b.Sidecar = &BlobSidecar{
					Blob:       kzg4844.Blob(j.Blobs[i]),
					Commitment: kzg4844.Commitment(j.Commitments[i]),
					Proof:      kzg4844.Proof(j.Proofs[i]),
				}
			}
			c.Blobs[i] = b
		}
	}
}

package types

import (
	"math/big"

	"github.com/defiweb/go-eth/crypto/kzg4844"
)

// The following interfaces are used to determine if a given call or
// transaction has specific capabilities.

// HasTransactionData specifes that the type has basic transaction
// data fields like chain ID, nonce, and signature.
type HasTransactionData interface {
	GetTransactionData() *TransactionData
	SetTransactionData(data TransactionData)
}

// HasCallData specifies that the type has basic call data fields
// like from, to, gas limit, value, and input.
type HasCallData interface {
	GetCallData() *CallData
	SetCallData(data CallData)
}

// HasLegacyPriceData specifies that the type uses legacy price data.
type HasLegacyPriceData interface {
	GetLegacyPriceData() *LegacyPriceData
	SetLegacyPriceData(data LegacyPriceData)
}

// HasAccessListData specifies that the type uses access list for EIP-2930
// transactions.
type HasAccessListData interface {
	GetAccessListData() *AccessListData
	SetAccessListData(data AccessListData)
}

// HasDynamicFeeData specifies that the type uses dynamic fee data for EIP-1559
// transactions.
type HasDynamicFeeData interface {
	GetDynamicFeeData() *DynamicFeeData
	SetDynamicFeeData(data DynamicFeeData)
}

// HasBlobData specifies that the type uses blob data for EIP-4844
// transactions.
type HasBlobData interface {
	GetBlobData() *BlobData
	SetBlobData(data BlobData)
}

// The following types are used to embed common fields and methods into call
// and transaction types. You probably do not want to use them directly.

// TransactionData contains common fields for transactions.
//
// This type is used to embed transaction data into other types.
type TransactionData struct {
	ChainID   *uint64    // ChainID is the chain ID.
	Nonce     *uint64    // Nonce is the transaction nonce.
	Signature *Signature // Signature is the transaction signature.
}

// GetTransactionData returns the embedded transaction data.
func (c *TransactionData) GetTransactionData() *TransactionData {
	return c
}

// SetTransactionData sets the embedded transaction data.
func (c *TransactionData) SetTransactionData(data TransactionData) {
	*c = data
}

// SetChainID sets the chain ID.
func (c *TransactionData) SetChainID(chainID uint64) {
	c.ChainID = &chainID
}

// SetNonce sets the transaction nonce.
func (c *TransactionData) SetNonce(nonce uint64) {
	c.Nonce = &nonce
}

// SetSignature sets the transaction signature.
func (c *TransactionData) SetSignature(signature Signature) {
	c.Signature = &signature
}

// Copy creates a deep copy of the TransactionFields.
func (c *TransactionData) Copy() *TransactionData {
	return &TransactionData{
		ChainID:   copyPtr(c.ChainID),
		Nonce:     copyPtr(c.Nonce),
		Signature: c.Signature.Copy(),
	}
}

func (c *TransactionData) toJSON(j *jsonTransaction) {
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

func (c *TransactionData) fromJSON(j *jsonTransaction) {
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

// CallData contains the basic fields for a call.
//
// This type is used to embed call data into other types.
type CallData struct {
	From     *Address // From is the sender address.
	To       *Address // To is the recipient address. Nil means contract creation.
	GasLimit *uint64  // GasLimit is the gas limit; if 0, there is no limit.
	Value    *big.Int // Value is the amount of wei to send.
	Input    []byte   // Input is the input data.
}

// GetCallData returns the embedded call data.
func (c *CallData) GetCallData() *CallData {
	return c
}

// SetCallData sets the embedded call data.
func (c *CallData) SetCallData(data CallData) {
	*c = data
}

// SetFrom sets the sender address.
func (c *CallData) SetFrom(from Address) {
	c.From = &from
}

// SetTo sets the recipient address.
func (c *CallData) SetTo(to Address) {
	c.To = &to
}

// SetGasLimit sets the gas limit.
func (c *CallData) SetGasLimit(gasLimit uint64) {
	c.GasLimit = &gasLimit
}

// SetValue sets the amount of wei to send.
func (c *CallData) SetValue(value *big.Int) {
	c.Value = value
}

// SetInput sets the input data.
func (c *CallData) SetInput(input []byte) {
	c.Input = input
}

// Copy creates a deep copy of the CallFields.
func (c *CallData) Copy() *CallData {
	if c == nil {
		return nil
	}
	return &CallData{
		From:     copyPtr(c.From),
		To:       copyPtr(c.To),
		GasLimit: copyPtr(c.GasLimit),
		Value:    copyBigInt(c.Value),
		Input:    copyBytes(c.Input),
	}
}

func (c *CallData) toJSON(j *jsonCall) {
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

func (c *CallData) fromJSON(j *jsonCall) {
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

// LegacyPriceData contains the gas price for legacy transactions.
//
// This type is used to embed legacy price data into other types.
type LegacyPriceData struct {
	GasPrice *big.Int // GasPrice is the gas price.
}

// LegacyPriceData returns the embedded legacy price data.
func (c *LegacyPriceData) GetLegacyPriceData() *LegacyPriceData {
	return c
}

// SetLegacyPriceData sets the embedded legacy price data.
func (c *LegacyPriceData) SetLegacyPriceData(data LegacyPriceData) {
	*c = data
}

// SetGasPrice sets the gas price.
func (c *LegacyPriceData) SetGasPrice(gasPrice *big.Int) {
	c.GasPrice = gasPrice
}

// Copy creates a deep copy of the LegacyPriceField.
func (c *LegacyPriceData) Copy() *LegacyPriceData {
	if c == nil {
		return nil
	}
	return &LegacyPriceData{
		GasPrice: copyBigInt(c.GasPrice),
	}
}

func (c *LegacyPriceData) toJSON(j *jsonCall) {
	if c.GasPrice != nil {
		j.GasPrice = NumberFromBigIntPtr(c.GasPrice)
	}
}

func (c *LegacyPriceData) fromJSON(j *jsonCall) {
	if j.GasPrice != nil {
		c.GasPrice = j.GasPrice.Big()
	}
}

// AccessListData contains the access list for EIP-2930 transactions.
//
// This type is used to embed access list data into other types.
type AccessListData struct {
	AccessList AccessList // AccessList is the EIP-2930 access list.
}

// AccessListData returns the embedded access list data.
func (c *AccessListData) GetAccessListData() *AccessListData {
	return c
}

// SetAccessListData sets the embedded access list data.
func (c *AccessListData) SetAccessListData(data AccessListData) {
	*c = data
}

// SetAccessList sets the access list.
func (c *AccessListData) SetAccessList(accessList AccessList) {
	c.AccessList = accessList
}

// Copy creates a deep copy of the AccessListField.
func (c *AccessListData) Copy() *AccessListData {
	if c == nil {
		return nil
	}
	return &AccessListData{
		AccessList: c.AccessList.Copy(),
	}
}

func (c *AccessListData) toJSON(j *jsonCall) {
	j.AccessList = c.AccessList
}

func (c *AccessListData) fromJSON(j *jsonCall) {
	c.AccessList = j.AccessList
}

// DynamicFeeData contains fee data for EIP-1559 transactions.
//
// This type is used to embed dynamic fee data into other types.
type DynamicFeeData struct {
	MaxFeePerGas         *big.Int // MaxFeePerGas is the maximum total fee per gas.
	MaxPriorityFeePerGas *big.Int // MaxPriorityFeePerGas is the maximum priority fee per gas.
}

// DynamicFeeData returns the embedded dynamic fee data.
func (c *DynamicFeeData) GetDynamicFeeData() *DynamicFeeData {
	return c
}

// SetDynamicFeeData sets the embedded dynamic fee data.
func (c *DynamicFeeData) SetDynamicFeeData(data DynamicFeeData) {
	*c = data
}

// SetMaxFeePerGas sets the maximum total fee per gas.
func (c *DynamicFeeData) SetMaxFeePerGas(maxFeePerGas *big.Int) {
	c.MaxFeePerGas = maxFeePerGas
}

// SetMaxPriorityFeePerGas sets the maximum priority fee per gas.
func (c *DynamicFeeData) SetMaxPriorityFeePerGas(maxPriorityFeePerGas *big.Int) {
	c.MaxPriorityFeePerGas = maxPriorityFeePerGas
}

// Copy creates a deep copy of the DynamicFeeFields.
func (c *DynamicFeeData) Copy() *DynamicFeeData {
	if c == nil {
		return nil
	}
	return &DynamicFeeData{
		MaxFeePerGas:         copyBigInt(c.MaxFeePerGas),
		MaxPriorityFeePerGas: copyBigInt(c.MaxPriorityFeePerGas),
	}
}

func (c *DynamicFeeData) toJSON(j *jsonCall) {
	if c.MaxFeePerGas != nil {
		j.MaxFeePerGas = NumberFromBigIntPtr(c.MaxFeePerGas)
	}
	if c.MaxPriorityFeePerGas != nil {
		j.MaxPriorityFeePerGas = NumberFromBigIntPtr(c.MaxPriorityFeePerGas)
	}
}

func (c *DynamicFeeData) fromJSON(j *jsonCall) {
	if j.MaxFeePerGas != nil {
		c.MaxFeePerGas = j.MaxFeePerGas.Big()
	}
	if j.MaxPriorityFeePerGas != nil {
		c.MaxPriorityFeePerGas = j.MaxPriorityFeePerGas.Big()
	}
}

// BlobData contains data for EIP-4844 blob transactions.
//
// Use NewBlobInfo to create a BlobInfo.
//
// This type is used to embed blob data into other types.
type BlobData struct {
	MaxFeePerBlobGas *big.Int   // MaxFeePerBlobGas is the maximum fee per blob gas.
	Blobs            []BlobInfo // Blobs is the list of blobs.
}

// BlobData returns the embedded blob data.
func (c *BlobData) GetBlobData() *BlobData {
	return c
}

// SetBlobData sets the embedded blob data.
func (c *BlobData) SetBlobData(data BlobData) {
	*c = data
}

// SetMaxFeePerBlobGas sets the maximum fee per blob gas.
func (c *BlobData) SetMaxFeePerBlobGas(maxFeePerBlobGas *big.Int) {
	c.MaxFeePerBlobGas = maxFeePerBlobGas
}

// SetBlobs sets the list of blobs.
//
// Use NewBlobInfo to create a BlobInfo.
func (c *BlobData) SetBlobs(blobs []BlobInfo) {
	c.Blobs = blobs
}

// AddBlob adds a blob to the list.
//
// Use NewBlobInfo to create a BlobInfo.
func (c *BlobData) AddBlob(blob BlobInfo) {
	c.Blobs = append(c.Blobs, blob)
}

// Copy creates a deep copy of the BlobFields.
func (c *BlobData) Copy() *BlobData {
	if c == nil {
		return nil
	}
	blobs := make([]BlobInfo, len(c.Blobs))
	for i, blob := range c.Blobs {
		blobs[i].Hash = blob.Hash
		blobs[i].Sidecar = copyPtr(blob.Sidecar)
	}
	return &BlobData{
		MaxFeePerBlobGas: copyBigInt(c.MaxFeePerBlobGas),
		Blobs:            blobs,
	}
}

func (c *BlobData) toJSON(j *jsonCall) {
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

func (c *BlobData) fromJSON(j *jsonCall) {
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

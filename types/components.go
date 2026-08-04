package types

import (
	"math/big"

	"github.com/defiweb/go-eth/crypto"
)

// The following interfaces are used to determine if a given call or
// transaction has specific capabilities.

// HasSigningData specifies that a type has data required by
// signed transactions.
type HasSigningData interface {
	GetSigningData() *SigningData
	SetSigningData(data SigningData)
}

// HasExecutionData specifies that the type has basic call execution data
// like from, to, gas limit, value, and input.
type HasExecutionData interface {
	GetExecutionData() *ExecutionData
	SetExecutionData(data ExecutionData)
}

// HasLegacyFeeData specifies that the type uses legacy price data.
type HasLegacyFeeData interface {
	GetLegacyFeeData() *LegacyFeeData
	SetLegacyFeeData(data LegacyFeeData)
}

// HasAccessListData specifies that the type uses an access list for EIP-2930
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

// GetSigningData is a helper function to get the [SigningData] from call
// or transaction types.
func GetSigningData(v any) *SigningData {
	if s, ok := v.(HasSigningData); ok {
		return s.GetSigningData()
	}
	return nil
}

// GetExecutionData is a helper function to get the [ExecutionData] from call
// or transaction types.
func GetExecutionData(v any) *ExecutionData {
	if s, ok := v.(HasExecutionData); ok {
		return s.GetExecutionData()
	}
	return nil
}

// GetLegacyFeeData is a helper function to get the [LegacyFeeData] from call
// or transaction types.
func GetLegacyFeeData(v any) *LegacyFeeData {
	if s, ok := v.(HasLegacyFeeData); ok {
		return s.GetLegacyFeeData()
	}
	return nil
}

// GetAccessListData is a helper function to get the [AccessListData] from call
// or transaction types.
func GetAccessListData(v any) *AccessListData {
	if s, ok := v.(HasAccessListData); ok {
		return s.GetAccessListData()
	}
	return nil
}

// GetDynamicFeeData is a helper function to get the [DynamicFeeData] from call
// or transaction types.
func GetDynamicFeeData(v any) *DynamicFeeData {
	if s, ok := v.(HasDynamicFeeData); ok {
		return s.GetDynamicFeeData()
	}
	return nil
}

// GetBlobData is a helper function to get the [BlobData] from call
// or transaction types.
func GetBlobData(v any) *BlobData {
	if s, ok := v.(HasBlobData); ok {
		return s.GetBlobData()
	}
	return nil
}

// The following types are used to embed common fields and methods into call
// and transaction types. You probably do not want to use them directly.

// SigningData contains common fields for signed transactions.
//
// This type is used to embed signing data into other types.
type SigningData struct {
	// ChainID is the chain ID.
	ChainID *uint64

	// Nonce is the transaction nonce.
	Nonce *uint64

	// Signature is the transaction signature.
	Signature *Signature
}

// GetSigningData returns the embedded signing data.
func (c *SigningData) GetSigningData() *SigningData {
	return c
}

// SetSigningData sets the embedded signing data.
func (c *SigningData) SetSigningData(data SigningData) {
	*c = data
}

// SetChainID sets the chain ID.
func (c *SigningData) SetChainID(chainID uint64) {
	c.ChainID = &chainID
}

// SetNonce sets the transaction nonce.
func (c *SigningData) SetNonce(nonce uint64) {
	c.Nonce = &nonce
}

// SetSignature sets the transaction signature.
func (c *SigningData) SetSignature(signature Signature) {
	c.Signature = &signature
}

// Copy creates a deep copy of the SigningData.
func (c *SigningData) Copy() *SigningData {
	var s *Signature
	if c.Signature != nil {
		s = c.Signature.Copy()
	}
	return &SigningData{
		ChainID:   copyPtr(c.ChainID),
		Nonce:     copyPtr(c.Nonce),
		Signature: s,
	}
}

func (c *SigningData) toJSON(j *jsonTransaction) {
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

func (c *SigningData) fromJSON(j *jsonTransaction) {
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

// ExecutionData contains the basic fields for a call.
//
// The From field is the sender address. It is used by the JSON-RPC client for
// key selection (eth_sendTransaction), as an explicit execution context
// (eth_call, eth_estimateGas). For signed transactions, the authoritative
// sender is recovered from the signature via txsign.Recover.
//
// This type is used to embed call data into other types.
type ExecutionData struct {
	// From is the sender address. Not part of the Ethereum wire protocol;
	// used for key selection and as an execution context by the RPC client.
	From *Address

	// To is the recipient address. Nil means contract creation.
	To *Address

	// GasLimit is the gas limit. Nil means not specified.
	GasLimit *uint64

	// Value is the amount of wei to send.
	Value *big.Int

	// Input is the call input data.
	Input []byte
}

// GetExecutionData returns the embedded execution data.
func (c *ExecutionData) GetExecutionData() *ExecutionData {
	return c
}

// SetExecutionData sets the embedded execution data.
func (c *ExecutionData) SetExecutionData(data ExecutionData) {
	*c = data
}

// SetFrom sets the sender address.
func (c *ExecutionData) SetFrom(from Address) {
	c.From = &from
}

// SetTo sets the recipient address.
func (c *ExecutionData) SetTo(to Address) {
	c.To = &to
}

// SetGasLimit sets the gas limit.
func (c *ExecutionData) SetGasLimit(gasLimit uint64) {
	c.GasLimit = &gasLimit
}

// SetValue sets the amount of wei to send.
func (c *ExecutionData) SetValue(value *big.Int) {
	c.Value = value
}

// SetInput sets the input data.
func (c *ExecutionData) SetInput(input []byte) {
	c.Input = input
}

// Copy creates a deep copy of the ExecutionData.
func (c *ExecutionData) Copy() *ExecutionData {
	if c == nil {
		return nil
	}
	return &ExecutionData{
		From:     copyPtr(c.From),
		To:       copyPtr(c.To),
		GasLimit: copyPtr(c.GasLimit),
		Value:    copyBigInt(c.Value),
		Input:    copyBytes(c.Input),
	}
}

func (c *ExecutionData) toJSON(j *jsonCall) {
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

func (c *ExecutionData) fromJSON(j *jsonCall) {
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

// LegacyFeeData contains the gas price for legacy transactions.
//
// This type is used to embed legacy fee data into other types.
type LegacyFeeData struct {
	// GasPrice is the gas price.
	GasPrice *big.Int
}

// GetLegacyFeeData returns the embedded legacy fee data.
func (c *LegacyFeeData) GetLegacyFeeData() *LegacyFeeData {
	return c
}

// SetLegacyFeeData sets the embedded legacy fee data.
func (c *LegacyFeeData) SetLegacyFeeData(data LegacyFeeData) {
	*c = data
}

// SetGasPrice sets the gas price.
func (c *LegacyFeeData) SetGasPrice(gasPrice *big.Int) {
	c.GasPrice = gasPrice
}

// Copy creates a deep copy of the LegacyFeeData.
func (c *LegacyFeeData) Copy() *LegacyFeeData {
	if c == nil {
		return nil
	}
	return &LegacyFeeData{
		GasPrice: copyBigInt(c.GasPrice),
	}
}

func (c *LegacyFeeData) toJSON(j *jsonCall) {
	if c.GasPrice != nil {
		j.GasPrice = NumberFromBigIntPtr(c.GasPrice)
	}
}

func (c *LegacyFeeData) fromJSON(j *jsonCall) {
	if j.GasPrice != nil {
		c.GasPrice = j.GasPrice.Big()
	}
}

// AccessListData contains the access list for EIP-2930 transactions.
//
// This type is used to embed access list data into other types.
type AccessListData struct {
	// AccessList is the EIP-2930 access list.
	AccessList AccessList
}

// GetAccessListData returns the embedded access list data.
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

// Copy creates a deep copy of the AccessListData.
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
	// MaxFeePerGas is the maximum total fee per gas.
	MaxFeePerGas *big.Int

	// MaxPriorityFeePerGas is the maximum priority fee per gas.
	MaxPriorityFeePerGas *big.Int
}

// GetDynamicFeeData returns the embedded dynamic fee data.
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

// Copy creates a deep copy of the DynamicFeeData.
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
	// MaxFeePerBlobGas is the maximum fee per blob gas.
	MaxFeePerBlobGas *big.Int

	// Blobs is the list of blobs.
	Blobs []BlobInfo
}

// GetBlobData returns the embedded blob data.
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

// Copy creates a deep copy of the BlobData.
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
		j.BlobHashes = make([]kzgHash, 0, len(c.Blobs))
		j.Blobs = make([]kzgBlob, 0, len(c.Blobs))
		j.Commitments = make([]kzgCommitment, 0, len(c.Blobs))
		j.Proofs = make([]kzgProof, 0, len(c.Blobs))
	}
	for _, b := range c.Blobs {
		j.BlobHashes = append(j.BlobHashes, kzgHash(b.Hash))
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
			b := BlobInfo{Hash: crypto.KZGHash(h)}
			if i < len(j.Blobs) && i < len(j.Commitments) && i < len(j.Proofs) {
				b.Sidecar = &BlobSidecar{
					Blob:       crypto.KZGBlob(j.Blobs[i]),
					Commitment: crypto.KZGCommitment(j.Commitments[i]),
					Proof:      crypto.KZGProof(j.Proofs[i]),
				}
			}
			c.Blobs[i] = b
		}
	}
}

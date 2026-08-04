package types

import (
	"encoding/json"
)

// CallBlob represents a call corresponding to the blob transaction type.
//
// Introduced by EIP-4844, this transaction type adds support for blob-carrying
// transactions.
type CallBlob struct {
	ExecutionData
	AccessListData
	DynamicFeeData
	BlobData
}

// NewCallBlob creates a new CallBlob.
func NewCallBlob() *CallBlob {
	return &CallBlob{}
}

// Copy creates a deep copy of the CallBlob.
func (c *CallBlob) Copy() Call {
	if c == nil {
		return nil
	}
	return &CallBlob{
		ExecutionData:  *c.ExecutionData.Copy(),
		AccessListData: *c.AccessListData.Copy(),
		DynamicFeeData: *c.DynamicFeeData.Copy(),
		BlobData:       *c.BlobData.Copy(),
	}
}

// MarshalJSON implements the json.Marshaler interface.
func (c *CallBlob) MarshalJSON() ([]byte, error) {
	j := &jsonCall{}
	c.ExecutionData.toJSON(j)
	c.AccessListData.toJSON(j)
	c.DynamicFeeData.toJSON(j)
	c.BlobData.toJSON(j)
	return json.Marshal(j)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (c *CallBlob) UnmarshalJSON(data []byte) error {
	j := &jsonCall{}
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	c.ExecutionData.fromJSON(j)
	c.AccessListData.fromJSON(j)
	c.DynamicFeeData.fromJSON(j)
	c.BlobData.fromJSON(j)
	return nil
}

var _ Call = (*CallBlob)(nil)

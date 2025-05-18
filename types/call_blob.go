package types

import (
	"encoding/json"
)

// CallBlob represents a call corresponding to the blob transaction type.
//
// Introduced by EIP-4844, this transaction type adds support for blob-carrying
// transactions.
type CallBlob struct {
	CallFields
	AccessListField
	DynamicFeeFields
	BlobFields
}

// NewCallBlob creates a new CallBlob.
func NewCallBlob() *CallBlob {
	return &CallBlob{}
}

// Copy creates a deep copy of the CallBlob.
func (c *CallBlob) Copy() *CallBlob {
	if c == nil {
		return nil
	}
	return &CallBlob{
		CallFields:       *c.CallFields.Copy(),
		AccessListField:  *c.AccessListField.Copy(),
		DynamicFeeFields: *c.DynamicFeeFields.Copy(),
		BlobFields:       *c.BlobFields.Copy(),
	}
}

// MarshalJSON implements the json.Marshaler interface.
func (c *CallBlob) MarshalJSON() ([]byte, error) {
	j := &jsonCall{}
	c.CallFields.toJSON(j)
	c.AccessListField.toJSON(j)
	c.DynamicFeeFields.toJSON(j)
	c.BlobFields.toJSON(j)
	return json.Marshal(j)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (c *CallBlob) UnmarshalJSON(data []byte) error {
	j := &jsonCall{}
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	c.CallFields.fromJSON(j)
	c.AccessListField.fromJSON(j)
	c.DynamicFeeFields.fromJSON(j)
	c.BlobFields.fromJSON(j)
	return nil
}

var _ Call = (*CallBlob)(nil)

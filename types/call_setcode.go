package types

import "encoding/json"

// CallSetCode represents a call corresponding to the set code transaction type.
//
// Introduced by EIP-7702, this transaction type allows EOAs to temporarily
// adopt code from a smart contract by providing a list of authorization tuples.
type CallSetCode struct {
	ExecutionData
	AccessListData
	DynamicFeeData
	AuthorizationData
}

// NewCallSetCode creates a new CallSetCode.
func NewCallSetCode() *CallSetCode {
	return &CallSetCode{}
}

// Copy creates a deep copy of the CallSetCode.
func (c *CallSetCode) Copy() Call {
	if c == nil {
		return nil
	}
	return &CallSetCode{
		ExecutionData:     *c.ExecutionData.Copy(),
		AccessListData:    *c.AccessListData.Copy(),
		DynamicFeeData:    *c.DynamicFeeData.Copy(),
		AuthorizationData: *c.AuthorizationData.Copy(),
	}
}

// MarshalJSON implements the json.Marshaler interface.
func (c *CallSetCode) MarshalJSON() ([]byte, error) {
	j := &jsonCall{}
	c.ExecutionData.toJSON(j)
	c.AccessListData.toJSON(j)
	c.DynamicFeeData.toJSON(j)
	c.AuthorizationData.toJSON(j)
	return json.Marshal(j)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (c *CallSetCode) UnmarshalJSON(data []byte) error {
	j := &jsonCall{}
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	c.ExecutionData.fromJSON(j)
	c.AccessListData.fromJSON(j)
	c.DynamicFeeData.fromJSON(j)
	c.AuthorizationData.fromJSON(j)
	return nil
}

var _ Call = (*CallSetCode)(nil)

package types

import "encoding/json"

// CallDynamicFee represents a call corresponding to the dynamic fee
// transaction type.
//
// Introduced by EIP-1559, this transaction type supports a new fee market
// mechanism with a base fee and a priority fee (tip).
type CallDynamicFee struct {
	ExecutionData
	DynamicFeeData
	AccessListData
}

// NewCallDynamicFee creates a new CallDynamicFee.
func NewCallDynamicFee() *CallDynamicFee {
	return &CallDynamicFee{}
}

// Copy creates a deep copy of the CallDynamicFee.
func (c *CallDynamicFee) Copy() Call {
	if c == nil {
		return nil
	}
	return &CallDynamicFee{
		ExecutionData:  *c.ExecutionData.Copy(),
		DynamicFeeData: *c.DynamicFeeData.Copy(),
		AccessListData: *c.AccessListData.Copy(),
	}
}

// MarshalJSON implements the json.Marshaler interface.
func (c *CallDynamicFee) MarshalJSON() ([]byte, error) {
	j := &jsonCall{}
	c.ExecutionData.toJSON(j)
	c.DynamicFeeData.toJSON(j)
	c.AccessListData.toJSON(j)
	return json.Marshal(j)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (c *CallDynamicFee) UnmarshalJSON(data []byte) error {
	j := &jsonCall{}
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	c.ExecutionData.fromJSON(j)
	c.DynamicFeeData.fromJSON(j)
	c.AccessListData.fromJSON(j)
	return nil
}

var _ Call = (*CallDynamicFee)(nil)

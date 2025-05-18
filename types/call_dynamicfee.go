package types

import "encoding/json"

// CallDynamicFee represents a call corresponding to the dynamic fee
// transaction type.
//
// Introduced by EIP-1559, this transaction type supports a new fee market
// mechanism with a base fee and a priority fee (tip).
type CallDynamicFee struct {
	CallFields
	DynamicFeeFields
	AccessListField
}

// NewCallDynamicFee creates a new CallDynamicFee.
func NewCallDynamicFee() *CallDynamicFee {
	return &CallDynamicFee{}
}

// Copy creates a deep copy of the CallDynamicFee.
func (c *CallDynamicFee) Copy() *CallDynamicFee {
	if c == nil {
		return nil
	}
	return &CallDynamicFee{
		CallFields:       *c.CallFields.Copy(),
		DynamicFeeFields: *c.DynamicFeeFields.Copy(),
		AccessListField:  *c.AccessListField.Copy(),
	}
}

// MarshalJSON implements the json.Marshaler interface.
func (c *CallDynamicFee) MarshalJSON() ([]byte, error) {
	j := &jsonCall{}
	c.CallFields.toJSON(j)
	c.DynamicFeeFields.toJSON(j)
	c.AccessListField.toJSON(j)
	return json.Marshal(j)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (c *CallDynamicFee) UnmarshalJSON(data []byte) error {
	j := &jsonCall{}
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	c.CallFields.fromJSON(j)
	c.DynamicFeeFields.fromJSON(j)
	c.AccessListField.fromJSON(j)
	return nil
}

var _ Call = (*CallDynamicFee)(nil)

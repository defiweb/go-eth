package types

import "encoding/json"

// CallAccessList represents a call corresponding to the access list
// transaction type.
//
// Introduced by EIP-2930, this transaction type includes an optional access
// list that specifies a list of addresses and storage keys the transaction
// plans to access.
type CallAccessList struct {
	CallFields
	LegacyPriceField
	AccessListField
}

// NewCallAccessList creates a new CallAccessList.
func NewCallAccessList() *CallAccessList {
	return &CallAccessList{}
}

// Copy creates a deep copy of the CallAccessList.
func (c *CallAccessList) Copy() *CallAccessList {
	if c == nil {
		return nil
	}
	return &CallAccessList{
		CallFields:       *c.CallFields.Copy(),
		LegacyPriceField: *c.LegacyPriceField.Copy(),
		AccessListField:  *c.AccessListField.Copy(),
	}
}

// MarshalJSON implements the json.Marshaler interface.
func (c *CallAccessList) MarshalJSON() ([]byte, error) {
	j := &jsonCall{}
	c.CallFields.toJSON(j)
	c.LegacyPriceField.toJSON(j)
	c.AccessListField.toJSON(j)
	return json.Marshal(j)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (c *CallAccessList) UnmarshalJSON(data []byte) error {
	j := &jsonCall{}
	if err := json.Unmarshal(data, j); err != nil {
		return err
	}
	c.CallFields.fromJSON(j)
	c.LegacyPriceField.fromJSON(j)
	c.AccessListField.fromJSON(j)
	return nil
}

var _ Call = (*CallAccessList)(nil)

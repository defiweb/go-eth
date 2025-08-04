package types

import "encoding/json"

// CallAccessList represents a call corresponding to the access list
// transaction type.
//
// Introduced by EIP-2930, this transaction type includes an optional access
// list that specifies a list of addresses and storage keys the transaction
// plans to access.
type CallAccessList struct {
	CallData
	LegacyPriceData
	AccessListData
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
		CallData:       *c.CallData.Copy(),
		LegacyPriceData: *c.LegacyPriceData.Copy(),
		AccessListData:  *c.AccessListData.Copy(),
	}
}

// MarshalJSON implements the json.Marshaler interface.
func (c *CallAccessList) MarshalJSON() ([]byte, error) {
	j := &jsonCall{}
	c.CallData.toJSON(j)
	c.LegacyPriceData.toJSON(j)
	c.AccessListData.toJSON(j)
	return json.Marshal(j)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (c *CallAccessList) UnmarshalJSON(data []byte) error {
	j := &jsonCall{}
	if err := json.Unmarshal(data, j); err != nil {
		return err
	}
	c.CallData.fromJSON(j)
	c.LegacyPriceData.fromJSON(j)
	c.AccessListData.fromJSON(j)
	return nil
}

var _ Call = (*CallAccessList)(nil)

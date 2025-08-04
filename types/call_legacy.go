package types

import "encoding/json"

// CallLegacy represents a call corresponding to the legacy transaction type.
type CallLegacy struct {
	CallData
	LegacyPriceData
}

// NewCallLegacy creates a new CallLegacy.
func NewCallLegacy() *CallLegacy {
	return &CallLegacy{}
}

// Copy creates a deep copy of the CallLegacy.
func (c *CallLegacy) Copy() *CallLegacy {
	if c == nil {
		return nil
	}
	return &CallLegacy{
		CallData:       *c.CallData.Copy(),
		LegacyPriceData: *c.LegacyPriceData.Copy(),
	}
}

// MarshalJSON implements the json.Marshaler interface.
func (c CallLegacy) MarshalJSON() ([]byte, error) {
	j := &jsonCall{}
	c.CallData.toJSON(j)
	c.LegacyPriceData.toJSON(j)
	return json.Marshal(j)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (c *CallLegacy) UnmarshalJSON(data []byte) error {
	j := &jsonCall{}
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	c.CallData.fromJSON(j)
	c.LegacyPriceData.fromJSON(j)
	return nil
}

var _ Call = (*CallLegacy)(nil)

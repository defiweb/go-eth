package types

import "encoding/json"

// CallLegacy represents a call corresponding to the legacy transaction type.
type CallLegacy struct {
	CallFields
	LegacyPriceField
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
		CallFields:       *c.CallFields.Copy(),
		LegacyPriceField: *c.LegacyPriceField.Copy(),
	}
}

// MarshalJSON implements the json.Marshaler interface.
func (c CallLegacy) MarshalJSON() ([]byte, error) {
	j := &jsonCall{}
	c.CallFields.toJSON(j)
	c.LegacyPriceField.toJSON(j)
	return json.Marshal(j)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (c *CallLegacy) UnmarshalJSON(data []byte) error {
	j := &jsonCall{}
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	c.CallFields.fromJSON(j)
	c.LegacyPriceField.fromJSON(j)
	return nil
}

var _ Call = (*CallLegacy)(nil)

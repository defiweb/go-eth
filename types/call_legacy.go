package types

import "encoding/json"

// CallLegacy represents a call corresponding to the legacy transaction type.
type CallLegacy struct {
	ExecutionData
	LegacyFeeData
}

// NewCallLegacy creates a new CallLegacy.
func NewCallLegacy() *CallLegacy {
	return &CallLegacy{}
}

// Copy creates a deep copy of the CallLegacy.
func (c *CallLegacy) Copy() Call {
	if c == nil {
		return nil
	}
	return &CallLegacy{
		ExecutionData: *c.ExecutionData.Copy(),
		LegacyFeeData: *c.LegacyFeeData.Copy(),
	}
}

// MarshalJSON implements the json.Marshaler interface.
func (c CallLegacy) MarshalJSON() ([]byte, error) {
	j := &jsonCall{}
	c.ExecutionData.toJSON(j)
	c.LegacyFeeData.toJSON(j)
	return json.Marshal(j)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (c *CallLegacy) UnmarshalJSON(data []byte) error {
	j := &jsonCall{}
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	c.ExecutionData.fromJSON(j)
	c.LegacyFeeData.fromJSON(j)
	return nil
}

var _ Call = (*CallLegacy)(nil)

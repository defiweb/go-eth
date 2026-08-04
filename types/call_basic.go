package types

import "encoding/json"

// CallBasic represents a simplest Ethereum call.
type CallBasic struct {
	ExecutionData
}

// NewCall creates a new [CallBasic] instance.
func NewCall() *CallBasic {
	// It is called NewCall instead of NewCallBasic because, most of the time,
	// people do not care about the call type; they just want to execute a
	// simple call without any special features.
	return &CallBasic{}
}

// Copy creates a deep copy of the CallBasic.
func (c *CallBasic) Copy() Call {
	return &CallBasic{
		ExecutionData: *c.ExecutionData.Copy(),
	}
}

// MarshalJSON implements the json.Marshaler interface.
func (c CallBasic) MarshalJSON() ([]byte, error) {
	j := &jsonCall{}
	c.ExecutionData.toJSON(j)
	return json.Marshal(j)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (c *CallBasic) UnmarshalJSON(bytes []byte) error {
	j := &jsonCall{}
	if err := json.Unmarshal(bytes, &j); err != nil {
		return err
	}
	c.ExecutionData.fromJSON(j)
	return nil
}

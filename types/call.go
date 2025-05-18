package types

import "encoding/json"

// Call is an interface that represents a generic Ethereum call.
type Call interface {
	json.Marshaler
	json.Unmarshaler

	CallData
}

// CallBasic represents a simplest Ethereum call.
type CallBasic struct {
	CallFields
}

// NewCall creates a new CallBasic.
func NewCall() *CallBasic {
	return &CallBasic{}
}

// Copy creates a deep copy of the CallBasic.
func (c *CallBasic) Copy() *CallBasic {
	return &CallBasic{
		CallFields: *c.CallFields.Copy(),
	}
}

// MarshalJSON implements the json.Marshaler interface.
func (c CallBasic) MarshalJSON() ([]byte, error) {
	j := &jsonCall{}
	c.CallFields.toJSON(j)
	return json.Marshal(j)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (c CallBasic) UnmarshalJSON(bytes []byte) error {
	j := &jsonCall{}
	if err := json.Unmarshal(bytes, &j); err != nil {
		return err
	}
	c.CallFields.fromJSON(j)
	return nil
}

type jsonCall struct {
	From                 *Address        `json:"from,omitempty"`
	To                   *Address        `json:"to,omitempty"`
	GasLimit             *Number         `json:"gas,omitempty"`
	GasPrice             *Number         `json:"gasPrice,omitempty"`
	MaxFeePerGas         *Number         `json:"maxFeePerGas,omitempty"`
	MaxFeePerBlobGas     *Number         `json:"maxFeePerBlobGas,omitempty"`
	MaxPriorityFeePerGas *Number         `json:"maxPriorityFeePerGas,omitempty"`
	Input                Bytes           `json:"input,omitempty"`
	Value                *Number         `json:"value,omitempty"`
	AccessList           AccessList      `json:"accessList,omitempty"`
	BlobHashes           []Hash          `json:"blobVersionedHashes,omitempty"`
	Blobs                []kzgBlob       `json:"blobs,omitempty"`
	Commitments          []kzgCommitment `json:"commitments,omitempty"`
	Proofs               []kzgProof      `json:"proofs,omitempty"`
}

var _ Call = (*CallBasic)(nil)

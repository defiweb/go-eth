package types

import "encoding/json"

// Call is an interface that represents a generic Ethereum call.
type Call interface {
	json.Marshaler
	json.Unmarshaler
	HasExecutionData

	Copy() Call
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
	BlobHashes           []kzgHash       `json:"blobVersionedHashes,omitempty"`
	Blobs                []kzgBlob       `json:"blobs,omitempty"`
	Commitments          []kzgCommitment `json:"commitments,omitempty"`
	Proofs               []kzgProof      `json:"proofs,omitempty"`
}

var _ Call = (*CallBasic)(nil)

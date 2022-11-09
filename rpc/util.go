package rpc

import (
	"encoding/json"

	"github.com/defiweb/go-eth/types"
)

// signTransactionResult is the result of an eth_signTransaction request.
// Some backends return only RLP encoded data, others return a JSON object,
// this type can handle both.
type signTransactionResult struct {
	Raw types.Bytes        `json:"raw"`
	Tx  *types.Transaction `json:"tx"`
}

func (s *signTransactionResult) UnmarshalJSON(input []byte) error {
	if len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"' {
		return json.Unmarshal(input, &s.Raw)
	}
	type alias struct {
		Raw types.Bytes        `json:"raw"`
		Tx  *types.Transaction `json:"tx"`
	}
	var dec alias
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	s.Tx = dec.Tx
	s.Raw = dec.Raw
	return nil
}

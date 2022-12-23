package transport

import (
	"encoding/json"
	"fmt"

	"github.com/defiweb/go-eth/types"
)

// rpcRequest is the JSON-RPC request object.
type rpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      *uint64         `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

// rpcResponse is the JSON-RPC response object.
type rpcResponse struct {
	// Common fields:
	JSONRPC string  `json:"jsonrpc"`
	ID      *uint64 `json:"id"`

	// Call response:
	Result json.RawMessage `json:"result,omitempty"`
	Error  *rpcError       `json:"error,omitempty"`

	// Notification response:
	Method string          `json:"method,omitempty"`
	Params json.RawMessage `json:"params,omitempty"`
}

// rpcSubscription is the JSON-RPC subscription object.
type rpcSubscription struct {
	Subscription types.Number    `json:"subscription"`
	Result       json.RawMessage `json:"result"`
}

// rpcError is the JSON-RPC error object.
type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// newRPCRequest creates a new JSON-RPC request object.
func newRPCRequest(id *uint64, method string, params []any) (rpcRequest, error) {
	rpcReq := rpcRequest{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  json.RawMessage("[]"),
	}
	if len(params) > 0 {
		params, err := json.Marshal(params)
		if err != nil {
			return rpcRequest{}, err
		}
		rpcReq.Params = params
	}
	return rpcReq, nil
}

// Error implements the error interface.
func (e *rpcError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

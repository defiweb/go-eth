package transport

import (
	"encoding/json"
	"fmt"
)

// rpcRequest is a jsonrpc rpcRequest.
type rpcRequest struct {
	JsonRPC string          `json:"jsonrpc"`
	ID      uint64          `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

// rpcResponse is a jsonrpc rpcResponse.
type rpcResponse struct {
	ID     uint64          `json:"id"`
	Result json.RawMessage `json:"result"`
	Error  *rpcError       `json:"error,omitempty"`
}

// rpcError is a jsonrpc error.
type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (e *rpcError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

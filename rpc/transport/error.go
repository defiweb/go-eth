package transport

import (
	"fmt"
	"net/http"

	"github.com/defiweb/go-eth/hexutil"
)

const (
	// Standard errors:
	ErrCodeUnauthorized     = 1
	ErrCodeActionNotAllowed = 2
	ErrCodeExecutionError   = 3
	ErrCodeParseError       = -32700
	ErrCodeInvalidRequest   = -32600
	ErrCodeMethodNotFound   = -32601
	ErrCodeInvalidParams    = -32602
	ErrCodeInternalError    = -32603

	// Common non-standard errors:
	ErrCodeGeneral       = -32000
	ErrCodeLimitExceeded = -32005

	// Erigon errors:
	ErigonErrCodeGeneral         = -32000
	ErigonErrCodeNotFound        = -32601
	ErigonErrCodeUnsupportedFork = -38005

	// Nethermind errors:
	NethermindErrCodeMethodNotSupported  = -32004
	NethermindErrCodeLimitExceeded       = -32005
	NethermindErrCodeTransactionRejected = -32010
	NethermindErrCodeExecutionError      = -32015
	NethermindErrCodeTimeout             = -32016
	NethermindErrCodeModuleTimeout       = -32017
	NethermindErrCodeAccountLocked       = -32020
	NethermindErrCodeUnknownBlockError   = -39001

	// Infura errors:
	InfuraErrCodeInvalidInput               = -32000
	InfuraErrCodeResourceNotFound           = -32001
	InfuraErrCodeResourceUnavailable        = -32002
	InfuraErrCodeTransactionRejected        = -32003
	InfuraErrCodeMethodNotSupported         = -32004
	InfuraErrCodeLimitExceeded              = -32005
	InfuraErrCodeJSONRPCVersionNotSupported = -32006

	// Alchemy errors:
	AlchemyErrCodeLimitExceeded = 429

	// Blast errors:
	BlastErrCodeAuthenticationFailed = -32099
	BlastErrCodeCapacityExceeded     = -32098
	BlastErrRateLimitReached         = -32097
)

type HTTPErrorCode interface {
	// HTTPErrorCode returns the HTTP status code.
	HTTPErrorCode() int
}

type RPCErrorCode interface {
	// RPCErrorCode returns the JSON-RPC error code.
	RPCErrorCode() int
}

type RPCErrorData interface {
	// RPCErrorData returns the JSON-RPC error data.
	RPCErrorData() any
}

// RPCError is an JSON-RPC error.
type RPCError struct {
	Code    int    // Code is the JSON-RPC error code.
	Message string // Message is the error message.
	Data    any    // Data associated with the error.
}

// NewRPCError creates a new RPC error.
//
// If data is a hex-encoded string, it will be decoded.
func NewRPCError(code int, message string, data any) *RPCError {
	if bin, ok := decodeHexData(data); ok {
		data = bin
	}
	return &RPCError{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// Error implements the error interface.
func (e *RPCError) Error() string {
	return fmt.Sprintf("RPC error: %d %s", e.Code, e.Message)
}

// RPCErrorCode implements the ErrorCode interface.
func (e *RPCError) RPCErrorCode() int {
	return e.Code
}

// RPCErrorData implements the ErrorData interface.
func (e *RPCError) RPCErrorData() any {
	return e.Data
}

// HTTPError is an HTTP error.
type HTTPError struct {
	Code int   // Code is the HTTP status code.
	Err  error // Err is an optional underlying error.
}

// NewHTTPError creates a new HTTP error.
func NewHTTPError(code int, err error) *HTTPError {
	return &HTTPError{
		Code: code,
		Err:  err,
	}
}

// Error implements the error interface.
func (e *HTTPError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("HTTP error: %d %s", e.Code, http.StatusText(e.Code))
	}
	return fmt.Sprintf("HTTP error: %d %s: %s", e.Code, http.StatusText(e.Code), e.Err)
}

// HTTPErrorCode implements the HTTPErrorCode interface.
func (e *HTTPError) HTTPErrorCode() int {
	return e.Code
}

// decodeHexData decodes hex-encoded data if present.
func decodeHexData(data any) (any, bool) {
	hex, ok := data.(string)
	if !ok {
		return nil, false
	}
	if !hexutil.Has0xPrefix(hex) {
		return nil, false
	}
	bin, err := hexutil.HexToBytes(hex)
	if err != nil {
		return nil, false
	}
	return bin, true
}

package rpc

import "context"

type Transport interface {
	// Call performs a JSON-Client call.
	Call(ctx context.Context, result any, method string, args ...any) error
}

package rpc

import (
	"context"
	"encoding/json"
)

type Transport interface {
	// Call performs a JSON-RPC call.
	Call(ctx context.Context, result any, method string, args ...any) error
}

type Subscriber interface {
	Subscribe(ctx context.Context, method string, args ...any) (ch chan json.RawMessage, id string, err error)
	Unsubscribe(ctx context.Context, id string) error
}

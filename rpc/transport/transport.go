package transport

import (
	"context"
	"encoding/json"
)

// Transport handles the transport layer of the JSON-RPC protocol.
type Transport interface {
	// Call performs a JSON-RPC call.
	Call(ctx context.Context, result any, method string, args ...any) error
}

// Subscriber is am transport that supports subscriptions.
type Subscriber interface {
	Transport

	// Subscribe starts a new subscription. It returns a channel that receives
	// subscription messages and a subscription ID.
	Subscribe(ctx context.Context, method string, args ...any) (ch chan json.RawMessage, id string, err error)

	// Unsubscribe cancels a subscription. The channel returned by Subscribe
	// will be closed.
	Unsubscribe(ctx context.Context, id string) error
}

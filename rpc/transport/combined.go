package transport

import (
	"context"
	"encoding/json"
)

// Combined is transport that uses separate transports for regular calls and
// subscriptions.
//
// It is recommended by some RPC providers to use HTTP for regular calls and
// WebSockets for subscriptions.
type Combined struct {
	calls Transport
	subs  SubscriptionTransport
}

// NewCombined creates a new Combined transport.
func NewCombined(call Transport, subscriber SubscriptionTransport) *Combined {
	return &Combined{
		calls: call,
		subs:  subscriber,
	}
}

// Call implements the Transport interface.
func (c *Combined) Call(ctx context.Context, result any, method string, args ...any) error {
	return c.calls.Call(ctx, result, method, args...)
}

// Subscribe implements the SubscriptionTransport interface.
func (c *Combined) Subscribe(ctx context.Context, method string, args ...any) (ch chan json.RawMessage, id string, err error) {
	return c.subs.Subscribe(ctx, method, args...)
}

// Unsubscribe implements the SubscriptionTransport interface.
func (c *Combined) Unsubscribe(ctx context.Context, id string) error {
	return c.subs.Unsubscribe(ctx, id)
}

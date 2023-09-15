package transport

import (
	"context"
	"encoding/json"
	"fmt"
	netURL "net/url"
)

// Transport handles the transport layer of the JSON-RPC protocol.
type Transport interface {
	// Call performs a JSON-RPC call.
	Call(ctx context.Context, result any, method string, args ...any) error
}

// SubscriptionTransport is transport that supports subscriptions.
type SubscriptionTransport interface {
	Transport

	// Subscribe starts a new subscription. It returns a channel that receives
	// subscription messages and a subscription ID.
	Subscribe(ctx context.Context, method string, args ...any) (ch chan json.RawMessage, id string, err error)

	// Unsubscribe cancels a subscription. The channel returned by Subscribe
	// will be closed.
	Unsubscribe(ctx context.Context, id string) error
}

// New returns a new Transport instance based on the URL scheme.
// Supported schemes are: http, https, ws, wss.
// If scheme is empty, it will use IPC.
//
// The context is used to close the underlying connection when the transport
// uses a websocket or IPC.
func New(ctx context.Context, rpcURL string) (Transport, error) {
	url, err := netURL.Parse(rpcURL)
	if err != nil {
		return nil, err
	}
	switch url.Scheme {
	case "http", "https":
		return NewHTTP(HTTPOptions{URL: rpcURL})
	case "ws", "wss":
		return NewWebsocket(WebsocketOptions{Context: ctx, URL: rpcURL})
	case "":
		return NewIPC(IPCOptions{Context: ctx, Path: rpcURL})
	default:
		return nil, fmt.Errorf("unsupported scheme: %s", url.Scheme)
	}
}

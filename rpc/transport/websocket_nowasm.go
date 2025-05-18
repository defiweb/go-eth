//go:build wasm || tinygo

package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// Websocket is a Transport implementation that uses the websocket
// protocol.
type Websocket struct{}

// WebsocketOptions contains options for the websocket transport.
type WebsocketOptions struct {
	// Context used to close the connection.
	Context context.Context

	// URL of the websocket endpoint.
	URL string

	// HTTPClient is the HTTP client to use. If nil, http.DefaultClient is
	// used.
	HTTPClient *http.Client

	// HTTPHeader specifies the HTTP headers to be included in the
	// websocket handshake request.
	HTTPHeader http.Header

	// Timeout is the timeout for the websocket requests. Default is 60s.
	Timout time.Duration

	// ErrorCh is an optional channel used to report errors.
	ErrorCh chan error
}

// NewWebsocket creates a new Websocket instance.
func NewWebsocket(opts WebsocketOptions) (*Websocket, error) {
	return nil, errors.New("websocket transport is not supported in WASM")
}

func (w *Websocket) Call(ctx context.Context, result any, method string, args ...any) error {
	return errors.New("websocket transport is not supported in WASM")
}

func (w *Websocket) Subscribe(ctx context.Context, method string, args ...any) (ch chan json.RawMessage, id string, err error) {
	return nil, "", errors.New("websocket transport is not supported in WASM")
}

func (w *Websocket) Unsubscribe(ctx context.Context, id string) error {
	return errors.New("websocket transport is not supported in WASM")
}

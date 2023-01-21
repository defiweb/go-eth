package transport

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// Websocket is a Transport implementation that uses the websocket
// protocol.
type Websocket struct {
	*stream
	conn *websocket.Conn
}

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
	if opts.URL == "" {
		return nil, errors.New("URL cannot be empty")
	}
	if opts.Context == nil {
		return nil, errors.New("context cannot be nil")
	}
	if opts.Timout == 0 {
		opts.Timout = 60 * time.Second
	}
	conn, _, err := websocket.Dial(opts.Context, opts.URL, &websocket.DialOptions{
		HTTPClient: opts.HTTPClient,
		HTTPHeader: opts.HTTPHeader,
	})
	if err != nil {
		return nil, err
	}
	i := &Websocket{
		stream: &stream{
			ctx:     opts.Context,
			errCh:   opts.ErrorCh,
			timeout: opts.Timout,
		},
		conn: conn,
	}
	i.onClose = i.close
	i.stream.initStream()
	go i.readerRoutine()
	go i.writerRoutine()
	return i, nil
}

func (ws *Websocket) readerRoutine() {
	// The background context is used here because closing context will
	// cause the nhooyr.io/websocket package to close a connection with
	// a close code of 1008 (policy violation) which is not what we want.
	ctx := context.Background()
	for {
		res := rpcResponse{}
		if err := wsjson.Read(ctx, ws.conn, &res); err != nil {
			if ws.ctx.Err() != nil || errors.As(err, &websocket.CloseError{}) {
				return
			}
			if ws.errCh != nil {
				ws.errCh <- fmt.Errorf("websocket reading error: %w", err)
			}
			continue
		}
		ws.readerCh <- res
	}
}

func (ws *Websocket) writerRoutine() {
	for {
		select {
		case <-ws.ctx.Done():
			return
		case req := <-ws.writerCh:
			if err := wsjson.Write(ws.ctx, ws.conn, req); err != nil {
				if ws.errCh != nil {
					ws.errCh <- fmt.Errorf("websocket writing error: %w", err)
				}
				continue
			}
		}
	}
}

func (ws *Websocket) close() {
	err := ws.conn.Close(websocket.StatusNormalClosure, "")
	if err != nil && ws.errCh != nil {
		ws.errCh <- fmt.Errorf("websocket closing error: %w", err)
	}
}

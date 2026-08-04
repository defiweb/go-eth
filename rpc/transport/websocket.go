package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// defaultPingInterval is the default interval between WebSocket ping messages.
const defaultPingInterval = 30 * time.Second

// Websocket is a [Transport] implementation that uses the WebSocket protocol.
type Websocket struct {
	stream
	conn *websocket.Conn
	ping time.Duration
}

// WebsocketOptions contains options for the [WebSocket] transport.
type WebsocketOptions struct {
	// Context is used to close the connection.
	Context context.Context

	// URL is the WebSocket endpoint.
	URL string

	// HTTPClient is an optional HTTP client used to configure the WebSocket
	// handshake. Because the underlying WebSocket library does not accept an
	// *http.Client directly, only a subset of its settings are applied: the
	// client's cookie Jar, and—when its Transport is an *http.Transport—the
	// Proxy, TLSClientConfig and DialContext settings. Other settings (custom
	// RoundTrippers, redirect policy, the client Timeout, etc.) are ignored.
	HTTPClient *http.Client

	// HTTPHeader specifies the HTTP headers to include in the WebSocket
	// handshake request.
	HTTPHeader http.Header

	// Timeout is the timeout for WebSocket requests. The default is 60s.
	Timeout time.Duration

	// PingInterval is the interval between WebSocket ping messages used to
	// keep the connection alive and detect dead peers. The default is 30s.
	// A negative value disables pings.
	PingInterval time.Duration

	// ReadBufferSize is the buffer size for incoming messages. The default
	// is 1.
	//
	// Once the buffer is full, the transport stops reading from the
	// connection until it drains.
	ReadBufferSize int

	// WriteBufferSize is the buffer size for outgoing requests. The default
	// is 1.
	WriteBufferSize int

	// SubscriptionBufferSize is the buffer size of the message queue of each
	// subscription. The default is 32.
	//
	// Once the buffer of any subscription is full, the transport stops reading
	// from the connection, which stalls every other subscription and call
	// sharing it. No message is dropped, but nothing else progresses either,
	// so this should be sized for the slowest consumer.
	SubscriptionBufferSize int

	// ErrorCh is an optional channel used to report errors.
	ErrorCh chan error
}

// NewWebsocket creates a new [Websocket] instance.
func NewWebsocket(opts WebsocketOptions) (*Websocket, error) {
	if opts.URL == "" {
		return nil, errors.New("URL cannot be empty")
	}
	if opts.Context == nil {
		return nil, errors.New("context cannot be nil")
	}
	if opts.Timeout == 0 {
		opts.Timeout = time.Minute
	}
	if opts.PingInterval == 0 {
		opts.PingInterval = defaultPingInterval
	}
	dialer := websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: opts.Timeout,
	}
	if c := opts.HTTPClient; c != nil {
		dialer.Jar = c.Jar
		if t, ok := c.Transport.(*http.Transport); ok {
			dialer.Proxy = t.Proxy
			dialer.TLSClientConfig = t.TLSClientConfig
			dialer.NetDialContext = t.DialContext
		}
	}
	conn, _, err := dialer.DialContext(opts.Context, opts.URL, opts.HTTPHeader)
	if err != nil {
		return nil, fmt.Errorf("failed to dial websocket: %w", err)
	}
	ws := &Websocket{
		conn: conn,
		ping: opts.PingInterval,
	}
	streamOpts := []streamOption{
		withStreamErrorCh(opts.ErrorCh),
		withStreamTimeout(opts.Timeout),
	}
	if opts.ReadBufferSize > 0 {
		streamOpts = append(streamOpts, withReadBufferSize(opts.ReadBufferSize))
	}
	if opts.WriteBufferSize > 0 {
		streamOpts = append(streamOpts, withWriteBufferSize(opts.WriteBufferSize))
	}
	if opts.SubscriptionBufferSize > 0 {
		streamOpts = append(streamOpts, withSubscriptionBufferSize(opts.SubscriptionBufferSize))
	}
	ws.stream.initStream(opts.Context, streamOpts...)
	if ws.ping > 0 {
		ws.conn.SetPongHandler(func(string) error {
			ws.resetReadDeadline()
			return nil
		})
	}
	go ws.readerRoutine()
	go ws.writerRoutine()
	if ws.ping > 0 {
		go ws.pingRoutine()
	}
	return ws, nil
}

func (ws *Websocket) readerRoutine() {
	ws.resetReadDeadline()
	for {
		_, msg, err := ws.conn.ReadMessage()
		if err != nil {
			if ws.ctx.Err() != nil {
				return
			}
			if websocket.IsCloseError(
				err,
				websocket.CloseNormalClosure,
				websocket.CloseGoingAway,
			) {
				return
			}
			ws.error(fmt.Errorf("websocket reading error: %w", err))
			return
		}
		ws.resetReadDeadline()
		var res rpcResponse
		if err := json.Unmarshal(msg, &res); err != nil {
			ws.error(fmt.Errorf("websocket json error: %w", err))
			continue
		}
		ws.read(res)
	}
}

func (ws *Websocket) writerRoutine() {
	defer ws.close()
	for {
		req, ok := ws.write()
		if !ok {
			return
		}
		if err := ws.conn.SetWriteDeadline(time.Now().Add(ws.timeout)); err != nil {
			ws.error(fmt.Errorf("websocket writing error: %w", err))
			continue
		}
		if err := ws.conn.WriteJSON(req); err != nil {
			ws.error(fmt.Errorf("websocket writing error: %w", err))
			continue
		}
	}
}

func (ws *Websocket) pingRoutine() {
	t := time.NewTicker(ws.ping)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			if err := ws.conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(ws.timeout)); err != nil {
				ws.error(fmt.Errorf("websocket ping error: %w", err))
				return
			}
		case <-ws.ctx.Done():
			return
		}
	}
}

func (ws *Websocket) resetReadDeadline() {
	if ws.ping > 0 {
		if err := ws.conn.SetReadDeadline(time.Now().Add(2 * ws.ping)); err != nil {
			ws.error(fmt.Errorf("websocket reset read deadline error: %w", err))
		}
	}
}

func (ws *Websocket) close() {
	if err := ws.conn.SetWriteDeadline(time.Now().Add(time.Second)); err != nil {
		ws.error(fmt.Errorf("websocket close error: %w", err))
	}
	msg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
	if err := ws.conn.WriteMessage(websocket.CloseMessage, msg); err != nil {
		ws.error(fmt.Errorf("websocket close error: %w", err))
	}
	if err := ws.conn.Close(); err != nil {
		ws.error(fmt.Errorf("websocket close error: %w", err))
	}
}

package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"

	"github.com/defiweb/go-eth/types"
)

// Websocket is a Transport implementation that uses the websocket
// protocol.
type Websocket struct {
	mu  sync.RWMutex
	ctx context.Context

	opts  WebsocketOptions
	id    uint64
	conn  *websocket.Conn
	calls map[uint64]chan *rpcResponse    // channels for RPC requests
	subs  map[string]chan json.RawMessage // channels for subscription notifications
	errCh chan error                      // optional error channel
}

// WebsocketOptions contains options for the websocket transport.
type WebsocketOptions struct {
	// Context used to close the connection.
	Context context.Context

	// URL of the websocket endpoint.
	URL string

	// HTTPClient is used for the connection.
	HTTPClient *http.Client

	// HTTPHeader specifies the HTTP headers included in the handshake request.
	HTTPHeader http.Header

	// Timeout is the timeout for the websocket requests. If the timeout is
	// reached, the request will return ErrWebsocketTimeout. Default is 60s.
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
		opts.Context = context.Background()
	}
	if opts.Timout == 0 {
		opts.Timout = 60 * time.Second
	}
	return &Websocket{
		ctx:   opts.Context,
		opts:  opts,
		calls: make(map[uint64]chan *rpcResponse),
		subs:  make(map[string]chan json.RawMessage),
		errCh: opts.ErrorCh,
	}, nil
}

// Call implements the Transport interface.
func (ws *Websocket) Call(ctx context.Context, result any, method string, args ...any) error {
	ctx, ctxCancel := context.WithTimeout(ctx, ws.opts.Timout)
	defer ctxCancel()

	// Ensure the connection is established.
	if err := ws.connect(ctx); err != nil {
		return err
	}

	// Prepare the RPC request.
	id := atomic.AddUint64(&ws.id, 1)
	req, err := newRPCRequest(&id, method, args)
	if err != nil {
		return err
	}

	// Send the request.
	if err := wsjson.Write(ctx, ws.conn, req); err != nil {
		return err
	}

	// Wait for the response.
	// The response is handled by the readerRoutine. It will send the response
	// to the ch channel.
	ch := make(chan *rpcResponse)
	ws.addCallCh(id, ch)
	defer ws.delCallCh(id)
	select {
	case res := <-ch:
		if res.Error != nil {
			return res.Error
		}
		if result != nil {
			return json.Unmarshal(res.Result, result)
		}
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

// Subscribe implements the Subscriber interface.
func (ws *Websocket) Subscribe(ctx context.Context, method string, args ...any) (chan json.RawMessage, string, error) {
	rawID := types.Number{}
	if err := ws.Call(ctx, &rawID, "eth_subscribe", method, args); err != nil {
		return nil, "", err
	}
	id := rawID.String()
	ch := make(chan json.RawMessage)
	ws.addSubCh(id, ch)
	return ch, id, nil
}

// Unsubscribe implements the Subscriber interface.
func (ws *Websocket) Unsubscribe(ctx context.Context, id string) error {
	if !ws.delSubCh(id) {
		return errors.New("unknown subscription")
	}
	return ws.Call(ctx, nil, "eth_unsubscribe", types.HexToNumber(id))
}

// connect establishes the websocket connection. If the connection is already
// established, it does nothing.
func (ws *Websocket) connect(ctx context.Context) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	if ws.conn != nil {
		return nil
	}
	var err error
	ws.conn, _, err = websocket.Dial(ctx, ws.opts.URL, &websocket.DialOptions{
		HTTPClient: ws.opts.HTTPClient,
		HTTPHeader: ws.opts.HTTPHeader,
	})
	if err != nil {
		return err
	}
	go ws.readerRoutine()
	go ws.contextHandlerRoutine()
	return nil
}

// readerRoutine reads messages from the websocket connection and dispatches
// them to the appropriate channel.
func (ws *Websocket) readerRoutine() {
	// The background context is used here because closing context will
	// cause the nhooyr.io/websocket package to close a connection with
	// a close code of 1008 (policy violation) which is not what we want.
	ctx := context.Background()
	for {
		res := &rpcResponse{}
		if err := wsjson.Read(ctx, ws.conn, res); err != nil {
			if ws.ctx.Err() != nil || errors.As(err, &websocket.CloseError{}) {
				return
			}
			if ws.errCh != nil {
				ws.errCh <- fmt.Errorf("websocket reading error: %w", err)
			}
			continue
		}
		switch {
		case res.ID == nil:
			// If the ID is nil, it is a subscription notification.
			sub := &rpcSubscription{}
			if err := json.Unmarshal(res.Params, sub); err != nil {
				if ws.errCh != nil {
					ws.errCh <- fmt.Errorf("websocket unmarshalling error: %w", err)
				}
				continue
			}
			ws.subChSend(sub.Subscription.String(), sub.Result)
		default:
			// If the ID is not nil, it is a response to a request.
			ws.callChSend(*res.ID, res)
		}
	}
}

// contextHandlerRoutine closes the connection when the context is canceled.
func (ws *Websocket) contextHandlerRoutine() {
	<-ws.ctx.Done()
	ws.mu.Lock()
	defer ws.mu.Unlock()
	for _, ch := range ws.calls {
		close(ch)
	}
	for _, ch := range ws.subs {
		close(ch)
	}
	ws.calls = nil
	ws.subs = nil
	err := ws.conn.Close(websocket.StatusNormalClosure, "")
	if err != nil && ws.errCh != nil {
		ws.errCh <- fmt.Errorf("websocket closing error: %w", err)
	}
}

// addCallCh adds a channel to the calls map. Incoming response that match the
// id will be sent to the given channel. Because message ids are unique, the
// channel must be deleted after the response is received using delCallCh.
func (ws *Websocket) addCallCh(id uint64, ch chan *rpcResponse) {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	ws.calls[id] = ch
}

// addSubCh adds a channel to the subs map. Incoming subscription notifications
// that match the id will be sent to the given channel.
func (ws *Websocket) addSubCh(id string, ch chan json.RawMessage) {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	ws.subs[id] = ch
}

// delCallCh deletes a channel from the calls map.
func (ws *Websocket) delCallCh(id uint64) bool {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	if ch, ok := ws.calls[id]; ok {
		close(ch)
		delete(ws.calls, id)
		return true
	}
	return false
}

// delSubCh deletes a channel from the subs map.
func (ws *Websocket) delSubCh(id string) bool {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	if ch, ok := ws.subs[id]; ok {
		close(ch)
		delete(ws.subs, id)
		return true
	}
	return false
}

// callChSend sends a response to the channel that matches the id.
func (ws *Websocket) callChSend(id uint64, res *rpcResponse) {
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	if ch := ws.calls[id]; ch != nil {
		ch <- res
	}
}

// subChSend sends a subscription notification to the channel that matches the
// id.
func (ws *Websocket) subChSend(id string, res json.RawMessage) {
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	if ch := ws.subs[id]; ch != nil {
		ch <- res
	}
}

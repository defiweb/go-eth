package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"

	"github.com/defiweb/go-eth/types"
)

var (
	ErrWebsocketTimeout    = errors.New("websocket timeout")
	ErrUnknownSubscription = errors.New("unknown subscription")
)

type Websocket struct {
	mu    sync.RWMutex
	id    uint64 // incrementing ID for RPC requests
	url   string // websocket URL
	conn  *websocket.Conn
	opts  WebsocketOptions
	calls map[uint64]chan *rpcResponse    // channels for RPC requests
	subs  map[string]chan json.RawMessage // channels for subscription notifications
	errCh chan error                      // optional error channel
}

type WebsocketOptions struct {
	// HTTPClient is used for the connection.
	HTTPClient *http.Client

	// HTTPHeader specifies the HTTP headers included in the handshake request.
	HTTPHeader http.Header

	// Timeout is the timeout for the websocket requests.
	Timout time.Duration

	// ErrorCh is an optional channel used to report errors.
	ErrorCh chan error
}

func NewWebsocket(url string, opts *WebsocketOptions) *Websocket {
	if opts == nil {
		opts = &WebsocketOptions{}
	}
	if opts.Timout == 0 {
		opts.Timout = 60 * time.Second
	}
	return &Websocket{
		url:   url,
		opts:  *opts,
		calls: make(map[uint64]chan *rpcResponse),
		subs:  make(map[string]chan json.RawMessage),
		errCh: opts.ErrorCh,
	}
}

// Call implements the rpc.Transport interface.
func (ws *Websocket) Call(ctx context.Context, result any, method string, args ...any) error {
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
	tm := time.After(ws.opts.Timout)
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
	case <-tm:
		return ErrWebsocketTimeout
	case <-ctx.Done():
		return nil
	}
	return nil
}

// Subscribe implements the rpc.Subscriber interface.
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

// Unsubscribe implements the rpc.Subscriber interface.
func (ws *Websocket) Unsubscribe(ctx context.Context, id string) error {
	if !ws.delSubCh(id) {
		return ErrUnknownSubscription
	}
	return ws.Call(ctx, nil, "eth_unsubscribe", types.HexToNumber(id))
}

func (ws *Websocket) connect(ctx context.Context) error {
	if ws.conn != nil {
		return nil
	}
	var err error
	ws.conn, _, err = websocket.Dial(ctx, ws.url, &websocket.DialOptions{
		HTTPClient: ws.opts.HTTPClient,
		HTTPHeader: ws.opts.HTTPHeader,
	})
	if err != nil {
		return err
	}
	go ws.readerRoutine(ctx)
	return nil
}

// readerRoutine reads messages from the websocket connection and dispatches
// them to the appropriate channel.
func (ws *Websocket) readerRoutine(ctx context.Context) {
	for {
		res := &rpcResponse{}
		if err := wsjson.Read(ctx, ws.conn, res); err != nil {
			if errors.Is(err, websocket.CloseError{}) {
				return
			}
			if ws.errCh != nil {
				ws.errCh <- err
			}
			continue
		}
		switch {
		case res.ID == nil:
			// If the ID is nil, it is a subscription notification.
			sub := &rpcSubscription{}
			if err := json.Unmarshal(res.Params, sub); err != nil {
				if ws.errCh != nil {
					ws.errCh <- err
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

func (ws *Websocket) addCallCh(id uint64, ch chan *rpcResponse) {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	ws.calls[id] = ch
}

func (ws *Websocket) addSubCh(id string, ch chan json.RawMessage) {
	ws.mu.Lock()
	defer ws.mu.Unlock()
	ws.subs[id] = ch
}

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

func (ws *Websocket) callChSend(id uint64, res *rpcResponse) {
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	if ch := ws.calls[id]; ch != nil {
		ch <- res
	}
}

func (ws *Websocket) subChSend(id string, res json.RawMessage) {
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	if ch := ws.subs[id]; ch != nil {
		ch <- res
	}
}

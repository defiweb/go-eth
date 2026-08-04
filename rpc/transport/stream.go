package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/defiweb/go-eth/types"
)

const (
	defaultStreamTimeout                = time.Minute
	defaultStreamReadBufferSize         = 1
	defaultStreamWriteBufferSize        = 1
	defaultStreamSubscriptionBufferSize = 32
)

type streamOption func(*stream)

func withWriteBufferSize(size int) streamOption {
	return func(s *stream) {
		s.writeCh = make(chan rpcRequest, size)
	}
}

func withReadBufferSize(size int) streamOption {
	return func(s *stream) {
		s.readCh = make(chan rpcResponse, size)
	}
}

func withSubscriptionBufferSize(size int) streamOption {
	return func(s *stream) {
		s.bufSize = size
	}
}

func withStreamTimeout(timeout time.Duration) streamOption {
	return func(s *stream) {
		s.timeout = timeout
	}
}

func withStreamErrorCh(errCh chan error) streamOption {
	return func(s *stream) {
		s.errCh = errCh
	}
}

// stream is a helper for handling JSON-RPC streams.
type stream struct {
	ctx context.Context

	writeCh chan rpcRequest  // Channel for sending requests used by structs that embed stream.
	readCh  chan rpcResponse // Channel for receiving responses used by structs that embed stream.
	errCh   chan error       // Channel to which errors are sent.
	timeout time.Duration    // Timeout for requests.
	bufSize int              // Buffer size for the subscription channels.

	// State fields. Should not be accessed by structs that embed stream.
	id  uint64          // Request ID counter.
	chs *streamChannels // Channels for handling requests and subscriptions.
}

// initStream initializes the stream struct with default values and starts
// the required goroutines.
func (s *stream) initStream(ctx context.Context, opts ...streamOption) *stream {
	s.ctx = ctx
	s.timeout = defaultStreamTimeout
	s.bufSize = defaultStreamSubscriptionBufferSize
	s.chs = newStreamChannels()
	for _, opt := range opts {
		opt(s)
	}
	if s.writeCh == nil {
		s.writeCh = make(chan rpcRequest, defaultStreamWriteBufferSize)
	}
	if s.readCh == nil {
		s.readCh = make(chan rpcResponse, defaultStreamReadBufferSize)
	}
	go s.streamRoutine()
	return s
}

// Call implements the [Transport] interface.
func (s *stream) Call(ctx context.Context, result any, method string, args ...any) error {
	ctx, ctxCancel := context.WithTimeout(ctx, s.timeout)
	defer ctxCancel()

	// Prepare the RPC request.
	id := atomic.AddUint64(&s.id, 1)
	req, err := newRPCRequest(&id, method, args)
	if err != nil {
		return fmt.Errorf("failed to create RPC request: %w", err)
	}

	// Prepare the channel for the response.
	ch, ok := s.chs.addCallCh(id)
	if !ok {
		return errors.New("stream closed")
	}

	// Send the request.
	select {
	case s.writeCh <- req:
	case <-s.ctx.Done():
		return s.ctx.Err()
	case <-ctx.Done():
		return ctx.Err()
	}

	// Wait for the response.
	// The response is handled by streamRoutine, which sends the response
	// to the ch channel.
	select {
	case res := <-ch:
		if res.Error != nil {
			return NewRPCError(
				res.Error.Code,
				res.Error.Message,
				res.Error.Data,
			)
		}
		if result != nil {
			if err := json.Unmarshal(res.Result, result); err != nil {
				return fmt.Errorf("failed to unmarshal RPC result: %w", err)
			}
		}
	case <-s.ctx.Done():
		return s.ctx.Err()
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

// Subscribe implements the [SubscriptionTransport] interface.
func (s *stream) Subscribe(ctx context.Context, method string, args ...any) (chan json.RawMessage, string, error) {
	rawID := types.Number{}
	params := make([]any, 0, 2)
	params = append(params, method)
	if len(args) > 0 {
		params = append(params, args...)
	}
	if err := s.Call(ctx, &rawID, "eth_subscribe", params...); err != nil {
		return nil, "", err
	}
	id := rawID.String()
	ch := make(chan json.RawMessage, s.bufSize)
	if !s.chs.addSubCh(id, ch) {
		return nil, "", errors.New("stream closed")
	}
	return ch, id, nil
}

// Unsubscribe implements the [SubscriptionTransport] interface.
func (s *stream) Unsubscribe(ctx context.Context, id string) error {
	if !s.chs.delSubCh(id) {
		return errors.New("unknown subscription")
	}
	num, err := types.NumberFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid subscription id: %w", err)
	}
	return s.Call(ctx, nil, "eth_unsubscribe", num)
}

func (s *stream) streamRoutine() {
	defer s.chs.close()
	for {
		select {
		case res, ok := <-s.readCh:
			if !ok {
				return
			}
			switch res.ID {
			case nil:
				// If the ID is nil, it is a subscription notification.
				sub := &rpcSubscription{}
				if err := json.Unmarshal(res.Params, sub); err != nil {
					s.error(fmt.Errorf("failed to unmarshal subscription: %w", err))
					continue
				}
				s.chs.sendSubCh(s.ctx, sub.Subscription.String(), sub.Result)
			default:
				// If the ID is not nil, it is a response to a request.
				s.chs.sendCallCh(s.ctx, *res.ID, res)
			}
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *stream) read(res rpcResponse) {
	select {
	case s.readCh <- res:
	case <-s.ctx.Done():
	}
}

func (s *stream) write() (rpcRequest, bool) {
	select {
	case req, ok := <-s.writeCh:
		return req, ok
	case <-s.ctx.Done():
		return rpcRequest{}, false
	}
}

func (s *stream) error(err error) {
	if s.errCh != nil {
		select {
		case s.errCh <- err:
		case <-s.ctx.Done():
		}
	}
}

type streamChannels struct {
	mu sync.RWMutex

	calls map[uint64]chan rpcResponse     // Map of request IDs to channels.
	subs  map[string]chan json.RawMessage // Map of subscription IDs to channels.
}

func newStreamChannels() *streamChannels {
	return &streamChannels{
		calls: make(map[uint64]chan rpcResponse),
		subs:  make(map[string]chan json.RawMessage),
	}
}

func (s *streamChannels) addCallCh(id uint64) (chan rpcResponse, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.calls == nil {
		return nil, false
	}
	ch := make(chan rpcResponse, 1)
	s.calls[id] = ch
	return ch, true
}

func (s *streamChannels) addSubCh(id string, ch chan json.RawMessage) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.subs == nil {
		return false
	}
	s.subs[id] = ch
	return true
}

func (s *streamChannels) delSubCh(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ch, ok := s.subs[id]; ok {
		close(ch)
		delete(s.subs, id)
		return true
	}
	return false
}

func (s *streamChannels) sendCallCh(ctx context.Context, id uint64, res rpcResponse) {
	s.mu.Lock()
	defer s.mu.Unlock()
	ch := s.calls[id]
	if ch == nil {
		return
	}
	delete(s.calls, id)
	select {
	case ch <- res:
	case <-ctx.Done():
	}
}

func (s *streamChannels) sendSubCh(ctx context.Context, id string, res json.RawMessage) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if ch := s.subs[id]; ch != nil {
		select {
		case ch <- res:
		case <-ctx.Done():
		}
	}
}

func (s *streamChannels) close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, ch := range s.subs {
		close(ch)
	}
	s.calls = nil
	s.subs = nil
}

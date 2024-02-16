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

// stream is a helper for handling JSON-RPC streams.
type stream struct {
	mu  sync.RWMutex
	ctx context.Context

	writerCh chan rpcRequest  // Channel for sending requests used by structs that embed stream.
	readerCh chan rpcResponse // Channel for receiving responses used by structs that embed stream.
	errCh    chan error       // Channel to which errors are sent.
	timeout  time.Duration    // Timeout for requests.
	onClose  func()           // Callback that is called when the stream is closed.

	// State fields. Should not be accessed by structs that embed stream.
	id    uint64                          // Request ID counter.
	calls map[uint64]chan rpcResponse     // Map of request IDs to channels.
	subs  map[string]chan json.RawMessage // Map of subscription IDs to channels.
}

// initStream initializes the stream struct with default values and starts
// goroutines.
func (s *stream) initStream() *stream {
	s.writerCh = make(chan rpcRequest)
	s.readerCh = make(chan rpcResponse)
	s.calls = make(map[uint64]chan rpcResponse)
	s.subs = make(map[string]chan json.RawMessage)
	go s.streamRoutine()
	go s.contextHandlerRoutine()
	return s
}

// Call implements the Transport interface.
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
	ch := make(chan rpcResponse)
	s.addCallCh(id, ch)
	defer s.delCallCh(id)

	// Send the request.
	s.writerCh <- req

	// Wait for the response.
	// The response is handled by the streamRoutine. It will send the response
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
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

// Subscribe implements the SubscriptionTransport interface.
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
	ch := make(chan json.RawMessage)
	s.addSubCh(id, ch)
	return ch, id, nil
}

// Unsubscribe implements the SubscriptionTransport interface.
func (s *stream) Unsubscribe(ctx context.Context, id string) error {
	if !s.delSubCh(id) {
		return errors.New("unknown subscription")
	}
	num, err := types.NumberFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid subscription id: %w", err)
	}
	return s.Call(ctx, nil, "eth_unsubscribe", num)
}

// readerRoutine reads messages from the stream connection and dispatches
// them to the appropriate channel.
func (s *stream) streamRoutine() {
	for {
		res, ok := <-s.readerCh
		if !ok {
			return
		}
		switch {
		case res.ID == nil:
			// If the ID is nil, it is a subscription notification.
			sub := &rpcSubscription{}
			if err := json.Unmarshal(res.Params, sub); err != nil {
				if s.errCh != nil {
					s.errCh <- fmt.Errorf("failed to unmarshal subscription: %w", err)
				}
				continue
			}
			s.subChSend(sub.Subscription.String(), sub.Result)
		default:
			// If the ID is not nil, it is a response to a request.
			s.callChSend(*res.ID, res)
		}
	}
}

// contextHandlerRoutine closes the connection when the context is canceled.
func (s *stream) contextHandlerRoutine() {
	<-s.ctx.Done()
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, ch := range s.calls {
		close(ch)
	}
	for _, ch := range s.subs {
		close(ch)
	}
	s.calls = nil
	s.subs = nil
	if s.onClose != nil {
		s.onClose()
	}
}

// addCallCh adds a channel to the calls map. Incoming response that match the
// id will be sent to the given channel. Because message ids are unique, the
// channel must be deleted after the response is received using delCallCh.
func (s *stream) addCallCh(id uint64, ch chan rpcResponse) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls[id] = ch
}

// addSubCh adds a channel to the subs map. Incoming subscription notifications
// that match the id will be sent to the given channel.
func (s *stream) addSubCh(id string, ch chan json.RawMessage) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.subs[id] = ch
}

// delCallCh deletes a channel from the calls map.
func (s *stream) delCallCh(id uint64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ch, ok := s.calls[id]; ok {
		close(ch)
		delete(s.calls, id)
		return true
	}
	return false
}

// delSubCh deletes a channel from the subs map.
func (s *stream) delSubCh(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ch, ok := s.subs[id]; ok {
		close(ch)
		delete(s.subs, id)
		return true
	}
	return false
}

// callChSend sends a response to the channel that matches the id.
func (s *stream) callChSend(id uint64, res rpcResponse) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if ch := s.calls[id]; ch != nil {
		ch <- res
	}
}

// subChSend sends a subscription notification to the channel that matches the
// id.
func (s *stream) subChSend(id string, res json.RawMessage) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if ch := s.subs[id]; ch != nil {
		ch <- res
	}
}

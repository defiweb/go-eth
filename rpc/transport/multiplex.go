package transport

import (
	"context"
	"encoding/json"
)

type Multiplex struct {
	baseTransport         Transport
	txMultiplexTransports []Transport
}

func NewMultiplex(baseTransport Transport, txMultiplexTransports ...Transport) *Multiplex {
	return &Multiplex{
		baseTransport:         baseTransport,
		txMultiplexTransports: txMultiplexTransports,
	}
}

// Call implements the Transport interface.
func (c *Multiplex) Call(ctx context.Context, result any, method string, args ...any) error {
	if shouldMultiplex(method) {
		var (
			ok  bool  // True if at least one call was successful.
			err error // The error of the last unsuccessful call.
		)
		for _, txTransport := range c.txMultiplexTransports {
			if ok {
				// If there was a successful call, we still need to call the
				// other transports, but we don't need to handle the result.
				_ = txTransport.Call(ctx, nil, method, args...)
				continue
			}
			callErr := txTransport.Call(ctx, result, method, args...)
			if callErr != nil {
				err = callErr
				continue
			}
			ok = true
		}
		if !ok {
			return err
		}
	}
	return c.baseTransport.Call(ctx, result, method, args...)
}

// Subscribe implements the SubscriptionTransport interface.
func (c *Multiplex) Subscribe(ctx context.Context, method string, args ...any) (ch chan json.RawMessage, id string, err error) {
	if s, ok := c.baseTransport.(SubscriptionTransport); ok {
		return s.Subscribe(ctx, method, args...)
	}
	return nil, "", ErrNotSubscriptionTransport
}

// Unsubscribe implements the SubscriptionTransport interface.
func (c *Multiplex) Unsubscribe(ctx context.Context, id string) error {
	if s, ok := c.baseTransport.(SubscriptionTransport); ok {
		return s.Unsubscribe(ctx, id)
	}
	return ErrNotSubscriptionTransport
}

func shouldMultiplex(method string) bool {
	switch method {
	case "eth_sendTransaction", "eth_sendRawTransaction", "eth_sendPrivateTransaction":
		return true
	default:
		return false
	}
}

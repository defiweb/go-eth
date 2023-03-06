package transport

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"time"
)

var ErrNotSubscriptionTransport = errors.New("transport does not implement SubscriptionTransport")

var (
	// RetryOnAnyError retries on any error except for the following:
	// -32700: Parse error.
	// -32600: Invalid request.
	// -32601: Method not found.
	// -32602: Invalid params.
	RetryOnAnyError = func(err error) bool {
		// List of errors that should not be retried:
		switch errorCode(err) {
		case -32700: // Parse error.
			return false
		case -32600: // Invalid request.
			return false
		case -32601: // Method not found.
			return false
		case -32602: // Invalid params.
			return false
		}
		// Retry on all other errors:
		return err != nil
	}

	// RetryOnLimitExceeded retries on the following errors:
	// -32005: Limit exceeded.
	// 429: Too many requests.
	RetryOnLimitExceeded = func(err error) bool {
		switch errorCode(err) {
		case -32005: // Limit exceeded.
			return true
		case 429: // Too many requests.
			return true
		}
		return false
	}
)

// ExponentialBackoffOptions contains options for the ExponentialBackoff function.
type ExponentialBackoffOptions struct {
	// BaseDelay is the base delay before the first retry.
	BaseDelay time.Duration

	// MaxDelay is the maximum delay between retries.
	MaxDelay time.Duration

	// ExponentialFactor is the exponential factor to use for calculating the delay.
	// The delay is calculated as BaseDelay * ExponentialFactor ^ retryCount.
	ExponentialFactor float64
}

var (
	// LinearBackoff returns a BackoffFunc that returns a constant delay.
	LinearBackoff = func(delay time.Duration) func(int) time.Duration {
		return func(_ int) time.Duration {
			return delay
		}
	}

	// ExponentialBackoff returns a BackoffFunc that returns an exponential delay.
	// The delay is calculated as BaseDelay * ExponentialFactor ^ retryCount.
	ExponentialBackoff = func(opts ExponentialBackoffOptions) func(int) time.Duration {
		return func(retryCount int) time.Duration {
			d := time.Duration(float64(opts.BaseDelay) * math.Pow(opts.ExponentialFactor, float64(retryCount)))
			if d > opts.MaxDelay {
				return opts.MaxDelay
			}
			return d
		}
	}
)

// Retry is a wrapper around another transport that retries requests.
type Retry struct {
	opts RetryOptions
}

// RetryOptions contains options for the Retry transport.
type RetryOptions struct {
	// Transport is the underlying transport to use.
	Transport Transport

	// RetryFunc is a function that returns true if the request should be
	// retried. The RetryOnAnyError and RetryOnLimitExceeded functions can be
	// used or a custom function can be provided.
	RetryFunc func(error) bool

	// BackoffFunc is a function that returns the delay before the next retry.
	// It takes the current retry count as an argument.
	BackoffFunc func(int) time.Duration

	// MaxRetries is the maximum number of retries. If negative, there is no limit.
	MaxRetries int
}

// NewRetry creates a new Retry instance.
func NewRetry(opts RetryOptions) (*Retry, error) {
	if opts.Transport == nil {
		return nil, errors.New("transport cannot be nil")
	}
	if opts.RetryFunc == nil {
		return nil, errors.New("retry function cannot be nil")
	}
	if opts.BackoffFunc == nil {
		return nil, errors.New("backoff function cannot be nil")
	}
	if opts.MaxRetries == 0 {
		return nil, errors.New("max retries cannot be zero")
	}
	return &Retry{opts: opts}, nil
}

// Call implements the Transport interface.
func (c *Retry) Call(ctx context.Context, result any, method string, args ...any) (err error) {
	var i int
	for {
		err = c.opts.Transport.Call(ctx, result, method, args...)
		if !c.opts.RetryFunc(err) {
			return err
		}
		if c.opts.MaxRetries >= 0 && i >= c.opts.MaxRetries {
			break
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(c.opts.BackoffFunc(i)):
		}
		i++
	}
	return err
}

// Subscribe implements the SubscriptionTransport interface.
func (c *Retry) Subscribe(ctx context.Context, method string, args ...any) (ch chan json.RawMessage, id string, err error) {
	if s, ok := c.opts.Transport.(SubscriptionTransport); ok {
		var i int
		for {
			ch, id, err = s.Subscribe(ctx, method, args...)
			if !c.opts.RetryFunc(err) {
				return ch, id, err
			}
			if c.opts.MaxRetries >= 0 && i >= c.opts.MaxRetries {
				break
			}
			select {
			case <-ctx.Done():
				return nil, "", ctx.Err()
			case <-time.After(c.opts.BackoffFunc(i)):
			}
			i++
		}
		return nil, "", err
	}
	return nil, "", ErrNotSubscriptionTransport
}

// Unsubscribe implements the SubscriptionTransport interface.
func (c *Retry) Unsubscribe(ctx context.Context, id string) (err error) {
	if s, ok := c.opts.Transport.(SubscriptionTransport); ok {
		var i int
		for {
			err = s.Unsubscribe(ctx, id)
			if !c.opts.RetryFunc(err) {
				return err
			}
			if c.opts.MaxRetries >= 0 && i >= c.opts.MaxRetries {
				break
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(c.opts.BackoffFunc(i)):
			}
			i++
		}
		return err
	}
	return ErrNotSubscriptionTransport
}

// errorCode returns either the JSON-RPC error code or HTTP status code.
// If there is no error or error code is not available, it returns 0.
func errorCode(err error) int {
	var rpcErr *RPCError
	if errors.As(err, &rpcErr) {
		return rpcErr.Code
	}
	var httpErr *HTTPError
	if errors.As(err, &httpErr) {
		return httpErr.Code
	}
	return 0
}

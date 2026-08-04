package transport

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"strings"
	"time"
)

var (
	// RetryOnAnyError retries on any error except for the following cases,
	// where retrying does not make sense:
	// 3: Execution error.
	// -32700: Parse error.
	// -32600: Invalid request.
	// -32601: Method not found.
	// -32602: Invalid params.
	// -32000: If the error message starts with "execution reverted".
	RetryOnAnyError = func(err error) bool {
		// List of errors that should not be retried:
		switch errorCode(err) {
		case ErrCodeExecutionError:
			return false
		case ErrCodeParseError:
			return false
		case ErrCodeInvalidRequest:
			return false
		case ErrCodeMethodNotFound:
			return false
		case ErrCodeInvalidParams:
			return false
		case ErrCodeGeneral:
			var rpcErr *RPCError
			if errors.As(err, &rpcErr) {
				if strings.HasPrefix(rpcErr.Message, "execution reverted") {
					return false
				}
			}
		}

		// Retry on all other errors:
		return err != nil
	}

	// RetryOnLimitExceeded retries only on the following errors:
	// -32005: Limit exceeded.
	// -32097: Rate limit reached (Blast).
	// 429: Too many requests.
	RetryOnLimitExceeded = func(err error) bool {
		switch errorCode(err) {
		case ErrCodeLimitExceeded:
			return true
		case BlastErrRateLimitReached:
			return true
		case AlchemyErrCodeLimitExceeded:
			return true
		}
		return false
	}
)

// ExponentialBackoffOptions contains options for the [ExponentialBackoff]
// function.
type ExponentialBackoffOptions struct {
	// BaseDelay is the base delay before the first retry.
	BaseDelay time.Duration

	// MaxDelay is the maximum delay between retries.
	MaxDelay time.Duration

	// ExponentialFactor is the exponential factor to use for calculating the
	// delay.
	//
	// The delay is calculated as BaseDelay * ExponentialFactor ^ retryCount.
	ExponentialFactor float64
}

var (
	// LinearBackoff returns a [BackoffFunc] that returns a constant delay.
	LinearBackoff = func(delay time.Duration) func(int) time.Duration {
		return func(_ int) time.Duration {
			return delay
		}
	}

	// ExponentialBackoff returns a [BackoffFunc] that returns an exponential delay.
	//
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

// Retry wraps another [Transport] and retries requests on failure.
type Retry struct {
	opts RetryOptions
}

// RetryOptions contains options for the [Retry] transport.
type RetryOptions struct {
	// Transport is the underlying transport to use.
	Transport Transport

	// RetryFunc is a function that returns true if the request should be
	// retried. You can use [RetryOnAnyError], [RetryOnLimitExceeded], or
	// provide a custom function.
	RetryFunc func(error) bool

	// BackoffFunc is a function that returns the delay before the next retry.
	// It takes the current retry count as an argument.
	BackoffFunc func(int) time.Duration

	// MaxRetries is the maximum number of retries after the initial
	// attempt. 0 means no retries (one attempt total). Negative means
	// unlimited retries.
	MaxRetries int
}

// NewRetry creates a new [Retry] instance.
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
	return &Retry{opts: opts}, nil
}

// Call implements the [Transport] interface.
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
		if err := wait(ctx, c.opts.BackoffFunc(i)); err != nil {
			return err
		}
		i++
	}
	return err
}

// Subscribe implements the [SubscriptionTransport] interface.
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
			if err := wait(ctx, c.opts.BackoffFunc(i)); err != nil {
				return nil, "", err
			}
			i++
		}
		return nil, "", err
	}
	return nil, "", ErrNotSubscriptionTransport
}

// Unsubscribe implements the [SubscriptionTransport] interface.
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
			if err := wait(ctx, c.opts.BackoffFunc(i)); err != nil {
				return err
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
	var rpcErr RPCErrorCode
	if errors.As(err, &rpcErr) {
		return rpcErr.RPCErrorCode()
	}
	var httpErr HTTPErrorCode
	if errors.As(err, &httpErr) {
		return httpErr.HTTPErrorCode()
	}
	return 0
}

func wait(ctx context.Context, d time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(d):
		return nil
	}
}

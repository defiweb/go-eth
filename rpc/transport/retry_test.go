package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type fakeTransport struct {
	callResult  chan error
	subResult   chan error
	unsubResult chan error
	callCount   int
	subCount    int
	unsubCount  int
}

func newFakeTransport() *fakeTransport {
	return &fakeTransport{
		callResult:  make(chan error),
		subResult:   make(chan error),
		unsubResult: make(chan error),
	}
}

func (f *fakeTransport) Call(ctx context.Context, result any, method string, args ...any) error {
	f.callCount++
	return <-f.callResult
}

func (f *fakeTransport) Subscribe(ctx context.Context, method string, args ...any) (ch chan json.RawMessage, id string, err error) {
	f.subCount++
	err = <-f.subResult
	return nil, "", err
}

func (f *fakeTransport) Unsubscribe(ctx context.Context, id string) error {
	f.unsubCount++
	return <-f.unsubResult
}

//nolint:funlen
func TestRetry(t *testing.T) {
	tests := []struct {
		retry   RetryOptions
		asserts func(t *testing.T, f *fakeTransport, r *Retry)
	}{
		// No retry on success (call).
		{
			retry: RetryOptions{
				Transport:   newFakeTransport(),
				MaxRetries:  1,
				RetryFunc:   RetryOnAnyError,
				BackoffFunc: LinearBackoff(0),
			},
			asserts: func(t *testing.T, f *fakeTransport, r *Retry) {
				go func() {
					f.callResult <- nil
				}()
				err := r.Call(context.Background(), nil, "foo")
				require.NoError(t, err)
				require.Equal(t, 1, f.callCount)
				require.Equal(t, 0, f.subCount)
				require.Equal(t, 0, f.unsubCount)
			},
		},
		// No retry on success (subscribe).
		{
			retry: RetryOptions{
				Transport:   newFakeTransport(),
				MaxRetries:  1,
				RetryFunc:   RetryOnAnyError,
				BackoffFunc: LinearBackoff(0),
			},
			asserts: func(t *testing.T, f *fakeTransport, r *Retry) {
				go func() {
					f.subResult <- nil
				}()
				_, _, err := r.Subscribe(context.Background(), "foo")
				require.NoError(t, err)
				require.Equal(t, 0, f.callCount)
				require.Equal(t, 1, f.subCount)
				require.Equal(t, 0, f.unsubCount)
			},
		},
		// No retry on success (unsubscribe).
		{
			retry: RetryOptions{
				Transport:   newFakeTransport(),
				MaxRetries:  1,
				RetryFunc:   RetryOnAnyError,
				BackoffFunc: LinearBackoff(0),
			},
			asserts: func(t *testing.T, f *fakeTransport, r *Retry) {
				go func() {
					f.unsubResult <- nil
				}()
				err := r.Unsubscribe(context.Background(), "foo")
				require.NoError(t, err)
				require.Equal(t, 0, f.callCount)
				require.Equal(t, 0, f.subCount)
				require.Equal(t, 1, f.unsubCount)
			},
		},
		// Retry on error (call).
		{
			retry: RetryOptions{
				Transport:   newFakeTransport(),
				MaxRetries:  1,
				RetryFunc:   RetryOnAnyError,
				BackoffFunc: LinearBackoff(0),
			},
			asserts: func(t *testing.T, f *fakeTransport, r *Retry) {
				go func() {
					f.callResult <- fmt.Errorf("foo")
					f.callResult <- nil
				}()
				err := r.Call(context.Background(), nil, "foo")
				require.NoError(t, err)
				require.Equal(t, 2, f.callCount)
				require.Equal(t, 0, f.subCount)
				require.Equal(t, 0, f.unsubCount)
			},
		},
		// Retry on error (subscribe).
		{
			retry: RetryOptions{
				Transport:   newFakeTransport(),
				MaxRetries:  1,
				RetryFunc:   RetryOnAnyError,
				BackoffFunc: LinearBackoff(0),
			},
			asserts: func(t *testing.T, f *fakeTransport, r *Retry) {
				go func() {
					f.subResult <- fmt.Errorf("foo")
					f.subResult <- nil
				}()
				_, _, err := r.Subscribe(context.Background(), "foo")
				require.NoError(t, err)
				require.Equal(t, 0, f.callCount)
				require.Equal(t, 2, f.subCount)
				require.Equal(t, 0, f.unsubCount)
			},
		},
		// Retry on error (unsubscribe).
		{
			retry: RetryOptions{
				Transport:   newFakeTransport(),
				MaxRetries:  1,
				RetryFunc:   RetryOnAnyError,
				BackoffFunc: LinearBackoff(0),
			},
			asserts: func(t *testing.T, f *fakeTransport, r *Retry) {
				go func() {
					f.unsubResult <- fmt.Errorf("foo")
					f.unsubResult <- nil
				}()
				err := r.Unsubscribe(context.Background(), "foo")
				require.NoError(t, err)
				require.Equal(t, 0, f.callCount)
				require.Equal(t, 0, f.subCount)
				require.Equal(t, 2, f.unsubCount)
			},
		},
		// Too many retries (call).
		{
			retry: RetryOptions{
				Transport:   newFakeTransport(),
				MaxRetries:  1,
				RetryFunc:   RetryOnAnyError,
				BackoffFunc: LinearBackoff(0),
			},
			asserts: func(t *testing.T, f *fakeTransport, r *Retry) {
				go func() {
					f.callResult <- fmt.Errorf("foo")
					f.callResult <- fmt.Errorf("foo")
				}()
				err := r.Call(context.Background(), nil, "foo")
				require.Error(t, err)
				require.Equal(t, 2, f.callCount)
				require.Equal(t, 0, f.subCount)
				require.Equal(t, 0, f.unsubCount)
			},
		},
		// Too many retries (subscribe).
		{
			retry: RetryOptions{
				Transport:   newFakeTransport(),
				MaxRetries:  1,
				RetryFunc:   RetryOnAnyError,
				BackoffFunc: LinearBackoff(0),
			},
			asserts: func(t *testing.T, f *fakeTransport, r *Retry) {
				go func() {
					f.subResult <- fmt.Errorf("foo")
					f.subResult <- fmt.Errorf("foo")
				}()
				_, _, err := r.Subscribe(context.Background(), "foo")
				require.Error(t, err)
				require.Equal(t, 0, f.callCount)
				require.Equal(t, 2, f.subCount)
				require.Equal(t, 0, f.unsubCount)
			},
		},
		// Too many retries (unsubscribe).
		{
			retry: RetryOptions{
				Transport:   newFakeTransport(),
				MaxRetries:  1,
				RetryFunc:   RetryOnAnyError,
				BackoffFunc: LinearBackoff(0),
			},
			asserts: func(t *testing.T, f *fakeTransport, r *Retry) {
				go func() {
					f.unsubResult <- fmt.Errorf("foo")
					f.unsubResult <- fmt.Errorf("foo")
				}()
				err := r.Unsubscribe(context.Background(), "foo")
				require.Error(t, err)
				require.Equal(t, 0, f.callCount)
				require.Equal(t, 0, f.subCount)
				require.Equal(t, 2, f.unsubCount)
			},
		},
		// Infinite retries until success
		{
			retry: RetryOptions{
				Transport:   newFakeTransport(),
				MaxRetries:  -1,
				RetryFunc:   RetryOnAnyError,
				BackoffFunc: LinearBackoff(0),
			},
			asserts: func(t *testing.T, f *fakeTransport, r *Retry) {
				go func() {
					f.callResult <- fmt.Errorf("foo")
					f.callResult <- fmt.Errorf("foo")
					f.callResult <- nil

					f.subResult <- fmt.Errorf("foo")
					f.subResult <- fmt.Errorf("foo")
					f.subResult <- nil

					f.unsubResult <- fmt.Errorf("foo")
					f.unsubResult <- fmt.Errorf("foo")
					f.unsubResult <- nil
				}()
				// Call
				err := r.Call(context.Background(), nil, "foo")
				require.NoError(t, err)
				require.Equal(t, 3, f.callCount)
				require.Equal(t, 0, f.subCount)
				require.Equal(t, 0, f.unsubCount)

				// Subscribe
				_, _, err = r.Subscribe(context.Background(), "foo")
				require.NoError(t, err)
				require.Equal(t, 3, f.callCount)
				require.Equal(t, 3, f.subCount)
				require.Equal(t, 0, f.unsubCount)

				// Unsubscribe
				err = r.Unsubscribe(context.Background(), "foo")
				require.NoError(t, err)
				require.Equal(t, 3, f.callCount)
				require.Equal(t, 3, f.subCount)
			},
		},
		// Infinite retries until context is canceled.
		{
			retry: RetryOptions{
				Transport:   newFakeTransport(),
				MaxRetries:  -1,
				RetryFunc:   RetryOnAnyError,
				BackoffFunc: LinearBackoff(0),
			},
			asserts: func(t *testing.T, f *fakeTransport, r *Retry) {
				ctx, cancel := context.WithCancel(context.Background())
				go func() {
					f.callResult <- fmt.Errorf("foo")

					// Wait a bit to make sure the retry passes the select statement
					// and blocks on the Call() call.
					time.Sleep(100 * time.Millisecond)
					cancel()

					// The retry function is blocked on the Call() call, so we need to
					// send a result to unblock it. Then no more retries will be made
					// because the context is canceled.
					//
					// In practice, this is not a problem because the transport will
					// return an error when the context is canceled so the retry
					// function will not be blocked on the Call() call.
					f.callResult <- fmt.Errorf("foo")
				}()
				err := r.Call(ctx, nil, "foo")
				require.Error(t, err)
				require.Equal(t, 2, f.callCount)
				require.Equal(t, 0, f.subCount)
				require.Equal(t, 0, f.unsubCount)
			},
		},
		// Do not retry if RetryFunc returns false.
		{
			retry: RetryOptions{
				Transport:   newFakeTransport(),
				MaxRetries:  1,
				RetryFunc:   func(error) bool { return false },
				BackoffFunc: LinearBackoff(0),
			},
			asserts: func(t *testing.T, f *fakeTransport, r *Retry) {
				go func() {
					f.callResult <- fmt.Errorf("foo")
					f.subResult <- fmt.Errorf("foo")
					f.unsubResult <- fmt.Errorf("foo")
				}()
				err := r.Call(context.Background(), nil, "foo")
				require.Error(t, err)

				_, _, err = r.Subscribe(context.Background(), "foo")
				require.Error(t, err)

				err = r.Unsubscribe(context.Background(), "foo")
				require.Error(t, err)

				require.Equal(t, 1, f.callCount)
				require.Equal(t, 1, f.subCount)
				require.Equal(t, 1, f.unsubCount)
			},
		},
		// Do not wait for backoff after the last retry.
		{
			retry: RetryOptions{
				Transport:   newFakeTransport(),
				MaxRetries:  1,
				RetryFunc:   RetryOnAnyError,
				BackoffFunc: LinearBackoff(100 * time.Millisecond),
			},
			asserts: func(t *testing.T, f *fakeTransport, r *Retry) {
				t0 := time.Now()
				go func() {
					f.callResult <- fmt.Errorf("foo")
					f.callResult <- fmt.Errorf("foo")
					f.subResult <- fmt.Errorf("foo")
					f.subResult <- fmt.Errorf("foo")
					f.unsubResult <- fmt.Errorf("foo")
					f.unsubResult <- fmt.Errorf("foo")
				}()
				err := r.Call(context.Background(), nil, "foo")
				require.Error(t, err)

				_, _, err = r.Subscribe(context.Background(), "foo")
				require.Error(t, err)

				err = r.Unsubscribe(context.Background(), "foo")
				require.Error(t, err)

				require.Equal(t, 2, f.callCount)
				require.Equal(t, 2, f.subCount)
				require.Equal(t, 2, f.unsubCount)

				// The backoff function should not be called after the last retry, so
				// the total time should be slightly more than 300 milliseconds.
				require.True(t, time.Since(t0) < 400*time.Millisecond)
			},
		},
	}
	for n, test := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			r, err := NewRetry(test.retry)
			require.NoError(t, err)
			test.asserts(t, r.opts.Transport.(*fakeTransport), r)
		})
	}
}

//nolint:dupl
func TestRetryOnAnyError(t *testing.T) {
	tests := []struct {
		err  error
		want bool
	}{
		{
			err:  nil,
			want: false,
		},
		{
			err:  fmt.Errorf("foo"),
			want: true,
		},
		{
			err:  context.Canceled,
			want: true,
		},
		{
			err:  context.DeadlineExceeded,
			want: true,
		},
		{
			err:  &HTTPError{Code: 429},
			want: true,
		},
		{
			err:  &HTTPError{Code: 500},
			want: true,
		},
		{
			err:  &RPCError{Code: -32700},
			want: false,
		},
		{
			err:  &RPCError{Code: -32600},
			want: false,
		},
		{
			err:  &RPCError{Code: -32601},
			want: false,
		},
		{
			err:  &RPCError{Code: -32602},
			want: false,
		},
		{
			err:  &RPCError{Code: -32603},
			want: true,
		},
		{
			err:  &RPCError{Code: -32604},
			want: true,
		},
		{
			err:  &RPCError{Code: -32005},
			want: true,
		},
	}
	for n, test := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			got := RetryOnAnyError(test.err)
			require.Equal(t, test.want, got)
		})
	}
}

//nolint:dupl
func TestRetryOnLimitExceeded(t *testing.T) {
	tests := []struct {
		err  error
		want bool
	}{
		{
			err:  nil,
			want: false,
		},
		{
			err:  fmt.Errorf("foo"),
			want: false,
		},
		{
			err:  context.Canceled,
			want: false,
		},
		{
			err:  context.DeadlineExceeded,
			want: false,
		},
		{
			err:  &HTTPError{Code: 429},
			want: true,
		},
		{
			err:  &HTTPError{Code: 500},
			want: false,
		},
		{
			err:  &RPCError{Code: -32700},
			want: false,
		},
		{
			err:  &RPCError{Code: -32600},
			want: false,
		},
		{
			err:  &RPCError{Code: -32601},
			want: false,
		},
		{
			err:  &RPCError{Code: -32602},
			want: false,
		},
		{
			err:  &RPCError{Code: -32603},
			want: false,
		},
		{
			err:  &RPCError{Code: -32604},
			want: false,
		},
		{
			err:  &RPCError{Code: -32005},
			want: true,
		},
	}
	for n, test := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			got := RetryOnLimitExceeded(test.err)
			require.Equal(t, test.want, got)
		})
	}
}

func TestLinearBackoff(t *testing.T) {
	tests := []struct {
		delay time.Duration
		want  []time.Duration
	}{
		{
			delay: 0,
			want: []time.Duration{
				0,
				0,
			},
		},
		{
			delay: 100 * time.Millisecond,
			want: []time.Duration{
				100 * time.Millisecond,
				100 * time.Millisecond,
			},
		},
	}
	for n, test := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			b := LinearBackoff(test.delay)
			for i, want := range test.want {
				got := b(i)
				require.Equal(t, want, got)
			}
		})
	}
}

func TestExponentialBackoff(t *testing.T) {
	tests := []struct {
		opts ExponentialBackoffOptions
		want []time.Duration
	}{
		{
			opts: ExponentialBackoffOptions{
				BaseDelay:         100 * time.Millisecond,
				MaxDelay:          1 * time.Second,
				ExponentialFactor: 1,
			},
			want: []time.Duration{
				100 * time.Millisecond,
				100 * time.Millisecond,
			},
		},
		{
			opts: ExponentialBackoffOptions{
				BaseDelay:         100 * time.Millisecond,
				MaxDelay:          1 * time.Second,
				ExponentialFactor: 2,
			},
			want: []time.Duration{
				100 * time.Millisecond,
				200 * time.Millisecond,
				400 * time.Millisecond,
				800 * time.Millisecond,
				1 * time.Second,
				1 * time.Second,
			},
		},
	}
	for n, test := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			b := ExponentialBackoff(test.opts)
			for i, want := range test.want {
				got := b(i)
				require.Equal(t, want, got)
			}
		})
	}
}

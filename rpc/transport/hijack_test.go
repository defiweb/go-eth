package transport

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHijacker(t *testing.T) {
	tests := []struct {
		transport *mockTransport
		asserts   func(t *testing.T, f *mockTransport, h *Hijack)
	}{
		// Test with no hijackers.
		{
			transport: newMockTransport(),
			asserts: func(t *testing.T, f *mockTransport, h *Hijack) {
				go func() {
					f.callResult <- nil
					f.subResult <- nil
					f.unsubResult <- nil
				}()
				err := h.Call(context.Background(), nil, "foo")
				require.NoError(t, err)

				_, _, err = h.Subscribe(context.Background(), "bar")
				require.NoError(t, err)

				err = h.Unsubscribe(context.Background(), "baz")
				require.NoError(t, err)

				require.Equal(t, 1, f.callCount)
				require.Equal(t, 1, f.subCount)
				require.Equal(t, 1, f.unsubCount)
			},
		},
		// Test the order of hijackers.
		{
			transport: newMockTransport(),
			asserts: func(t *testing.T, f *mockTransport, h *Hijack) {
				var order []string

				h.Use(&mockHijacker{callFn: func(next CallFunc) CallFunc {
					return func(ctx context.Context, t Transport, result any, method string, args ...any) (err error) {
						order = append(order, "call1")
						return next(ctx, t, result, method, args...)
					}
				}})

				h.Use(&mockHijacker{callFn: func(next CallFunc) CallFunc {
					return func(ctx context.Context, t Transport, result any, method string, args ...any) (err error) {
						order = append(order, "call2")
						return next(ctx, t, result, method, args...)
					}
				}})

				go func() {
					f.callResult <- nil
				}()
				err := h.Call(context.Background(), nil, "foo")
				require.NoError(t, err)

				assert.Equal(t, []string{"call2", "call1"}, order)
				assert.Equal(t, 1, f.callCount)
			},
		},
		// Test the order of context hijackers.
		{
			transport: newMockTransport(),
			asserts: func(t *testing.T, f *mockTransport, h *Hijack) {
				var order []string

				ctx := WithHijackers(context.Background(), &mockHijacker{callFn: func(next CallFunc) CallFunc {
					return func(ctx context.Context, t Transport, result any, method string, args ...any) (err error) {
						order = append(order, "call1")
						return next(ctx, t, result, method, args...)
					}
				}})

				ctx = WithHijackers(ctx, &mockHijacker{callFn: func(next CallFunc) CallFunc {
					return func(ctx context.Context, t Transport, result any, method string, args ...any) (err error) {
						order = append(order, "call2")
						return next(ctx, t, result, method, args...)
					}
				}})

				go func() {
					f.callResult <- nil
				}()
				err := h.Call(ctx, nil, "foo")
				require.NoError(t, err)

				assert.Equal(t, []string{"call2", "call1"}, order)
				assert.Equal(t, 1, f.callCount)
			},
		},
		// Test the order of hijackers in context mixed with the "Use" method.
		{
			transport: newMockTransport(),
			asserts: func(t *testing.T, f *mockTransport, h *Hijack) {
				var order []string

				h.Use(&mockHijacker{callFn: func(next CallFunc) CallFunc {
					return func(ctx context.Context, t Transport, result any, method string, args ...any) (err error) {
						order = append(order, "call1")
						return next(ctx, t, result, method, args...)
					}
				}})

				h.Use(&mockHijacker{callFn: func(next CallFunc) CallFunc {
					return func(ctx context.Context, t Transport, result any, method string, args ...any) (err error) {
						order = append(order, "call2")
						return next(ctx, t, result, method, args...)
					}
				}})

				ctx := WithHijackers(context.Background(), &mockHijacker{callFn: func(next CallFunc) CallFunc {
					return func(ctx context.Context, t Transport, result any, method string, args ...any) (err error) {
						order = append(order, "call3")
						return next(ctx, t, result, method, args...)
					}
				}})

				ctx = WithHijackers(ctx, &mockHijacker{callFn: func(next CallFunc) CallFunc {
					return func(ctx context.Context, t Transport, result any, method string, args ...any) (err error) {
						order = append(order, "call4")
						return next(ctx, t, result, method, args...)
					}
				}})

				go func() {
					f.callResult <- nil
				}()
				err := h.Call(ctx, nil, "foo")
				require.NoError(t, err)

				assert.Equal(t, []string{"call4", "call3", "call2", "call1"}, order)
				assert.Equal(t, 1, f.callCount)
			},
		},
	}
	for n, test := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			h := NewHijacker(test.transport)
			test.asserts(t, test.transport, h)
		})
	}
}

type mockHijacker struct {
	callFn  func(next CallFunc) CallFunc
	subFn   func(next SubscribeFunc) SubscribeFunc
	unsubFn func(next UnsubscribeFunc) UnsubscribeFunc
}

func (m *mockHijacker) Call() func(next CallFunc) CallFunc {
	return m.callFn
}

func (m *mockHijacker) Subscribe() func(next SubscribeFunc) SubscribeFunc {
	return m.subFn
}

func (m *mockHijacker) Unsubscribe() func(next UnsubscribeFunc) UnsubscribeFunc {
	return m.unsubFn
}

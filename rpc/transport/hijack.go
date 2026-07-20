package transport

import (
	"context"
	"encoding/json"
)

// Hijacker intercepts and modifies calls to an underlying [Transport]
// using the middleware pattern.
//
// The 'next' function must be called to continue the call chain.
//
// The transport passed to 'next' is the underlying [Transport] instance;
// calling it will bypass any subsequent hijackers.
//
// Hijackers must not modify received arguments; instead, they should create
// and modify copies if necessary.
type Hijacker interface {
	// Call returns a [CallFunc] that intercepts and modifies
	// the 'Call' method.
	//
	// If nil is returned, no middleware is applied.
	Call() func(next CallFunc) CallFunc

	// Subscribe returns a [SubscribeFunc] that intercepts and modifies
	// the 'Subscribe' method.
	//
	// If nil is returned, no middleware is applied.
	Subscribe() func(next SubscribeFunc) SubscribeFunc

	// Unsubscribe returns an [UnsubscribeFunc] that intercepts and modifies
	// the 'Unsubscribe' method.
	//
	// If nil is returned, no middleware is applied.
	Unsubscribe() func(next UnsubscribeFunc) UnsubscribeFunc
}

type (
	CallFunc        func(ctx context.Context, t Transport, result any, method string, args ...any) error
	SubscribeFunc   func(ctx context.Context, t SubscriptionTransport, method string, args ...any) (ch chan json.RawMessage, id string, err error)
	UnsubscribeFunc func(ctx context.Context, t SubscriptionTransport, id string) error
)

// Hijack is a wrapper around another [Transport] that allows hijacking
// and modifying the behavior of the underlying [Transport] using the
// middleware pattern.
type Hijack struct {
	transport Transport
	callFunc  CallFunc
	subFunc   SubscribeFunc
	unsubFunc UnsubscribeFunc
}

// NewHijacker creates a new [Hijack] instance.
func NewHijacker(t Transport, hs ...Hijacker) *Hijack {
	h := &Hijack{
		transport: t,
		callFunc:  defCall,
		subFunc:   defSub,
		unsubFunc: defUnsub,
	}
	h.Use(hs...)
	return h
}

// Use wraps the hijacker chain with the provided hijackers. Hijackers
// are applied left to right - each one becomes the new outermost layer,
// so the last hijacker in the slice executes first on every call.
func (h *Hijack) Use(hs ...Hijacker) {
	for _, m := range hs {
		if m == nil {
			continue
		}
		if call := m.Call(); call != nil {
			h.callFunc = call(h.callFunc)
		}
		if sub := m.Subscribe(); sub != nil {
			h.subFunc = sub(h.subFunc)
		}
		if unsub := m.Unsubscribe(); unsub != nil {
			h.unsubFunc = unsub(h.unsubFunc)
		}
	}
}

// Call implements the [Transport] interface.
func (h *Hijack) Call(ctx context.Context, result any, method string, args ...any) error {
	hs := getHijackers(ctx)
	fn := h.callFunc
	for i := 0; i < len(hs); i++ {
		if call := hs[i].Call(); call != nil {
			fn = call(fn)
		}
	}
	return fn(ctx, h.transport, result, method, args...)
}

// Subscribe implements the [SubscriptionTransport] interface.
func (h *Hijack) Subscribe(ctx context.Context, method string, args ...any) (ch chan json.RawMessage, id string, err error) {
	if s, ok := h.transport.(SubscriptionTransport); ok {
		hs := getHijackers(ctx)
		fn := h.subFunc
		for i := 0; i < len(hs); i++ {
			if sub := hs[i].Subscribe(); sub != nil {
				fn = sub(fn)
			}
		}
		return fn(ctx, s, method, args...)
	}
	return nil, "", ErrNotSubscriptionTransport
}

// Unsubscribe implements the [SubscriptionTransport] interface.
func (h *Hijack) Unsubscribe(ctx context.Context, id string) error {
	if s, ok := h.transport.(SubscriptionTransport); ok {
		hs := getHijackers(ctx)
		fn := h.unsubFunc
		for i := 0; i < len(hs); i++ {
			if unsub := hs[i].Unsubscribe(); unsub != nil {
				fn = unsub(fn)
			}
		}
		return fn(ctx, s, id)
	}
	return ErrNotSubscriptionTransport
}

// WithHijackers returns a new context with the provided hijackers added.
// These hijackers supplement those registered via [Hijack.Use] and are
// applied on every call that uses this context.
func WithHijackers(ctx context.Context, hs ...Hijacker) context.Context {
	return context.WithValue(ctx, hijackerContextKey{}, append(getHijackers(ctx), hs...))
}

func getHijackers(ctx context.Context) []Hijacker {
	if hs, ok := ctx.Value(hijackerContextKey{}).([]Hijacker); ok {
		return hs
	}
	return nil
}

func defCall(ctx context.Context, t Transport, result any, method string, args ...any) error {
	return t.Call(ctx, result, method, args...)
}

func defSub(ctx context.Context, t SubscriptionTransport, method string, args ...any) (ch chan json.RawMessage, id string, err error) {
	return t.Subscribe(ctx, method, args...)
}

func defUnsub(ctx context.Context, t SubscriptionTransport, id string) error {
	return t.Unsubscribe(ctx, id)
}

type hijackerContextKey struct{}

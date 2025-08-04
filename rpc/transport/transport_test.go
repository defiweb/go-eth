package transport

import (
	"context"
	"encoding/json"
)

type mockTransport struct {
	callResult  chan error
	subResult   chan error
	unsubResult chan error
	callCount   int
	subCount    int
	unsubCount  int
}

func newMockTransport() *mockTransport {
	return &mockTransport{
		callResult:  make(chan error),
		subResult:   make(chan error),
		unsubResult: make(chan error),
	}
}

func (t *mockTransport) Call(ctx context.Context, result any, method string, args ...any) error {
	t.callCount++
	return <-t.callResult
}

func (t *mockTransport) Subscribe(ctx context.Context, method string, args ...any) (ch chan json.RawMessage, id string, err error) {
	t.subCount++
	err = <-t.subResult
	return nil, "", err
}

func (t *mockTransport) Unsubscribe(ctx context.Context, id string) error {
	t.unsubCount++
	return <-t.unsubResult
}

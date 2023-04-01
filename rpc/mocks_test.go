package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
)

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type httpMock struct {
	*transport.HTTP

	Request      *http.Request
	ResponseMock *http.Response
}

func newHTTPMock() *httpMock {
	h := &httpMock{}
	h.HTTP, _ = transport.NewHTTP(transport.HTTPOptions{
		URL: "http://localhost",
		HTTPClient: &http.Client{
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				h.Request = req
				return h.ResponseMock, nil
			}),
		},
	})
	return h
}

type streamMock struct {
	t *testing.T

	SubscribeMocks   []subscribeMock
	UnsubscribeMocks []unsubscribeMock
}

type subscribeMock struct {
	ArgMethod string
	ArgParams []any
	RetCh     chan json.RawMessage
	RetID     string
	RetErr    error
}

type unsubscribeMock struct {
	ArgID     string
	ResultErr error
}

func newStreamMock(t *testing.T) *streamMock {
	return &streamMock{t: t}
}

func (s *streamMock) Call(_ context.Context, _ any, _ string, _ ...any) error {
	return errors.New("not implemented")
}

func (s *streamMock) Subscribe(_ context.Context, method string, args ...any) (ch chan json.RawMessage, id string, err error) {
	require.NotEmpty(s.t, s.SubscribeMocks)
	m := s.SubscribeMocks[0]
	s.SubscribeMocks = s.SubscribeMocks[1:]
	require.Equal(s.t, m.ArgMethod, method)
	require.Equal(s.t, len(m.ArgParams), len(args))
	for i := range m.ArgParams {
		require.Equal(s.t, m.ArgParams[i], args[i])
	}
	return m.RetCh, m.RetID, m.RetErr
}

func (s *streamMock) Unsubscribe(_ context.Context, id string) error {
	require.NotEmpty(s.t, s.UnsubscribeMocks)
	m := s.UnsubscribeMocks[0]
	s.UnsubscribeMocks = s.UnsubscribeMocks[1:]
	require.Equal(s.t, m.ArgID, id)
	return m.ResultErr
}

type keyMock struct {
	addressCallback         func() types.Address
	signHashCallback        func(hash types.Hash) (*types.Signature, error)
	signMessageCallback     func(data []byte) (*types.Signature, error)
	signTransactionCallback func(tx *types.Transaction) error
}

func (k *keyMock) Address() types.Address {
	return k.addressCallback()
}

func (k *keyMock) SignHash(hash types.Hash) (*types.Signature, error) {
	return k.signHashCallback(hash)
}

func (k *keyMock) SignMessage(data []byte) (*types.Signature, error) {
	return k.signMessageCallback(data)
}

func (k *keyMock) SignTransaction(tx *types.Transaction) error {
	return k.signTransactionCallback(tx)
}

func (k *keyMock) VerifyHash(hash types.Hash, sig types.Signature) bool {
	return false
}

func (k keyMock) VerifyMessage(data []byte, sig types.Signature) bool {
	return false
}

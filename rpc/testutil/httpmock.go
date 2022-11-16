package testutil

import (
	"net/http"

	"github.com/defiweb/go-eth/rpc/transport"
)

type RoundTripFunc func(req *http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type HTTPMock struct {
	*transport.HTTP
	Request  *http.Request
	Response *http.Response
}

func NewHTTPMock() *HTTPMock {
	h := &HTTPMock{}
	h.HTTP = transport.NewHTTPWithClient("http://localhost", &http.Client{
		Transport: RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			h.Request = req
			return h.Response, nil
		}),
	})
	return h
}

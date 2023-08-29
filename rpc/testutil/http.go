package testutil

import (
	"net/http"

	"github.com/defiweb/go-eth/rpc/transport"
)

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type HTTPMock struct {
	*transport.HTTP

	Request  *http.Request  // Request that was sent.
	Response *http.Response // Response that will be returned.
}

func NewHTTPMock() *HTTPMock {
	h := &HTTPMock{}
	h.HTTP, _ = transport.NewHTTP(transport.HTTPOptions{
		URL: "http://localhost",
		HTTPClient: &http.Client{
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				h.Request = req
				return h.Response, nil
			}),
		},
	})
	return h
}

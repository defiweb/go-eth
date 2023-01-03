package transport

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/defiweb/go-eth/types"
)

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type httpMock struct {
	*HTTP
	Request  *http.Request
	Response *http.Response
}

func TestHTTP(t *testing.T) {
	tests := []struct {
		asserts func(t *testing.T, h *httpMock)
	}{
		// Simple request:
		{
			asserts: func(t *testing.T, h *httpMock) {
				h.Response = &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(`{"id":1, "jsonrpc":"2.0", "result":"0x1"}`))),
				}
				result := types.Number{}
				require.NoError(t, h.Call(context.Background(), &result, "eth_getBalance", "0x1111111111111111111111111111111111111111", "latest"))
				assert.Equal(t, h.Request.URL.String(), "http://localhost")
				assert.Equal(t, h.Request.Method, "POST")
				assert.Equal(t, h.Request.Header.Get("X-Test"), "test")
				assert.Equal(t, h.Request.Header.Get("Content-Type"), "application/json")
				requestBody, err := io.ReadAll(h.Request.Body)
				assert.NoError(t, err)
				assert.JSONEq(t, `{"id":1, "jsonrpc":"2.0", "method":"eth_getBalance", "params":["0x1111111111111111111111111111111111111111", "latest"]}`, string(requestBody))
				assert.Equal(t, result.Big().String(), "1")
			},
		},
		// ID must increment:
		{
			asserts: func(t *testing.T, h *httpMock) {
				// First request:
				h.Response = &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(`{"id":1, "jsonrpc":"2.0", "result":"0x1"}`))),
				}
				require.NoError(t, h.Call(context.Background(), nil, "eth_a"))
				requestBody, err := io.ReadAll(h.Request.Body)
				assert.NoError(t, err)
				assert.JSONEq(t, `{"id":1, "jsonrpc":"2.0", "method":"eth_a", "params":[]}`, string(requestBody))

				// Second request:
				h.Response = &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(`{"id":2, "jsonrpc":"2.0", "result":"0x2"}`))),
				}
				require.NoError(t, h.Call(context.Background(), nil, "eth_b"))
				requestBody, err = io.ReadAll(h.Request.Body)
				assert.NoError(t, err)
				assert.JSONEq(t, `{"id":2, "jsonrpc":"2.0", "method":"eth_b", "params":[]}`, string(requestBody))
			},
		},
		// Error response:
		{
			asserts: func(t *testing.T, h *httpMock) {
				h.Response = &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(`{"id":1, "jsonrpc":"2.0", "error":{"code":-32601, "message":"Method not found"}}`))),
				}
				result := types.Number{}
				err := h.Call(context.Background(), &result, "eth_a")
				assert.Error(t, err)
				assert.Equal(t, err.Error(), "-32601: Method not found")
			},
		},
		// Invalid response:
		{
			asserts: func(t *testing.T, h *httpMock) {
				h.Response = &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(`{`))),
				}
				result := types.Number{}
				err := h.Call(context.Background(), &result, "eth_a")
				assert.Error(t, err)
			},
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			h := &httpMock{}
			h.HTTP, _ = NewHTTP(HTTPOptions{
				URL: "http://localhost",
				HTTPHeader: http.Header{
					"X-Test": []string{"test"},
				},
				HTTPClient: &http.Client{
					Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
						h.Request = req
						return h.Response, nil
					}),
				},
			})
			tt.asserts(t, h)
		})
	}
}

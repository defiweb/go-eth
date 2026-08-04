package rpc

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
)

func TestHijackAddress(t *testing.T) {
	tc := []struct {
		name     string
		hijacker *hijackAddress
		method   string
		args     []any
		request  []string
		response []string
	}{
		{
			name:     "set address",
			hijacker: &hijackAddress{replace: false, address: types.MustAddressFromHex("0x1111111111111111111111111111111111111111")},
			method:   "eth_sendTransaction",
			args: []any{func() types.Transaction {
				return types.NewTransactionAccessList()
			}()},
			request: []string{
				`{"jsonrpc":"2.0","id":1,"method":"eth_sendTransaction","params":[{"from": "0x1111111111111111111111111111111111111111"}]}`,
			},
			response: []string{
				`{"jsonrpc":"2.0","id":1,"result":"0x1111111111111111111111111111111111111111111111111111111111111111"}`,
			},
		},
		{
			name:     "replace address",
			hijacker: &hijackAddress{replace: true, address: types.MustAddressFromHex("0x1111111111111111111111111111111111111111")},
			method:   "eth_sendTransaction",
			args: []any{func() types.Transaction {
				tx := types.NewTransactionAccessList()
				tx.SetFrom(types.MustAddressFromHex("0x2222222222222222222222222222222222222222"))
				return tx
			}()},
			request: []string{
				`{"jsonrpc":"2.0","id":1,"method":"eth_sendTransaction","params":[{"from": "0x1111111111111111111111111111111111111111"}]}`,
			},
			response: []string{
				`{"jsonrpc":"2.0","id":1,"result":"0x1111111111111111111111111111111111111111111111111111111111111111"}`,
			},
		},
	}
	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			httpMock := newHTTPMock()
			httpMock.Handler = func(req *http.Request) (*http.Response, error) {
				require.NotEmpty(t, tt.request)
				require.NotEmpty(t, tt.response)

				body, err := io.ReadAll(req.Body)
				require.NoError(t, err)
				require.JSONEq(t, tt.request[0], string(body), fmt.Sprintf("expected: %s, got: %s", tt.request[0], string(body)))

				res := tt.response[0]
				tt.request = tt.request[1:]
				tt.response = tt.response[1:]
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(res)),
				}, nil
			}

			hijacker := transport.NewHijacker(httpMock, tt.hijacker)

			err := hijacker.Call(context.Background(), nil, tt.method, tt.args...)
			assert.Len(t, tt.request, 0)
			assert.Len(t, tt.response, 0)
			require.NoError(t, err)
		})
	}
}

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

func TestHijackChainID(t *testing.T) {
	tc := []struct {
		name     string
		ctx      context.Context
		hijacker *hijackChainID
		method   string
		args     []any
		request  []string
		response []string
	}{
		{
			name:     "set chainID",
			hijacker: &hijackChainID{},
			method:   "eth_sendTransaction",
			args:     []any{types.NewTransactionAccessList()},
			request: []string{
				`{"jsonrpc":"2.0","id":1,"method":"eth_chainId","params":[]}`,
				`{"jsonrpc":"2.0","id":2,"method":"eth_sendTransaction","params":[{"chainId": "0x1"}]}`,
			},
			response: []string{
				`{"jsonrpc":"2.0","id": 1,"result": "0x01"}`,
				`{"jsonrpc":"2.0","id": 1,"result": "0x1111111111111111111111111111111111111111111111111111111111111111"}`,
			},
		},
		{
			name:     "do not replace chainID",
			hijacker: &hijackChainID{replace: false},
			method:   "eth_sendTransaction",
			args: []any{func() types.Transaction {
				tx := types.NewTransactionAccessList()
				tx.SetChainID(2)
				return tx
			}()},
			request: []string{
				`{"jsonrpc":"2.0","id":1,"method":"eth_sendTransaction","params":[{"chainId": "0x2"}]}`,
			},
			response: []string{
				`{"jsonrpc":"2.0","id": 1,"result": "0x1111111111111111111111111111111111111111111111111111111111111111"}`,
			},
		},
		{
			name:     "replace chainID",
			hijacker: &hijackChainID{replace: true},
			method:   "eth_sendTransaction",
			args: []any{func() types.Transaction {
				tx := types.NewTransactionAccessList()
				tx.SetChainID(2)
				return tx
			}()},
			request: []string{
				`{"jsonrpc":"2.0","id":1,"method":"eth_chainId","params":[]}`,
				`{"jsonrpc":"2.0","id":2,"method":"eth_sendTransaction","params":[{"chainId": "0x1"}]}`,
			},
			response: []string{
				`{"jsonrpc":"2.0","id": 1,"result": "0x01"}`,
				`{"jsonrpc":"2.0","id": 1,"result": "0x1111111111111111111111111111111111111111111111111111111111111111"}`,
			},
		},
		{
			name:     "context: chain ID value bypasses cache",
			ctx:      ContextWithChainID(context.Background(), 42),
			hijacker: &hijackChainID{},
			method:   "eth_sendTransaction",
			args:     []any{types.NewTransactionAccessList()},
			request: []string{
				`{"jsonrpc":"2.0","id":1,"method":"eth_sendTransaction","params":[{"chainId": "0x2a"}]}`,
			},
			response: []string{
				`{"jsonrpc":"2.0","id": 1,"result": "0x1111111111111111111111111111111111111111111111111111111111111111"}`,
			},
		},
		{
			name:     "context: replace overrides struct",
			ctx:      ContextWithChainIDReplace(context.Background(), true),
			hijacker: &hijackChainID{replace: false},
			method:   "eth_sendTransaction",
			args: []any{func() types.Transaction {
				tx := types.NewTransactionAccessList()
				tx.SetChainID(2)
				return tx
			}()},
			request: []string{
				`{"jsonrpc":"2.0","id":1,"method":"eth_chainId","params":[]}`,
				`{"jsonrpc":"2.0","id":2,"method":"eth_sendTransaction","params":[{"chainId": "0x1"}]}`,
			},
			response: []string{
				`{"jsonrpc":"2.0","id": 1,"result": "0x01"}`,
				`{"jsonrpc":"2.0","id": 1,"result": "0x1111111111111111111111111111111111111111111111111111111111111111"}`,
			},
		},
	}
	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.ctx
			if ctx == nil {
				ctx = context.Background()
			}
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

			err := hijacker.Call(ctx, nil, tt.method, tt.args...)
			assert.Len(t, tt.request, 0)
			assert.Len(t, tt.response, 0)
			require.NoError(t, err)
		})
	}
}

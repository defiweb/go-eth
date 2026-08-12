package rpc

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
)

func TestHijackLegacyGasFee(t *testing.T) {
	tc := []struct {
		name     string
		ctx      context.Context
		hijacker *hijackLegacyGasFee
		method   string
		args     []any
		request  []string
		response []string
	}{
		{
			name:     "set gas price",
			hijacker: &hijackLegacyGasFee{multiplier: 1.0},
			method:   "eth_sendTransaction",
			args:     []any{types.NewTransactionLegacy()},
			request: []string{
				`{"jsonrpc":"2.0","id":1,"method":"eth_gasPrice","params":[]}`,
				`{"jsonrpc":"2.0","id":2,"method":"eth_sendTransaction","params":[{"gasPrice":"0x1000"}]}`,
			},
			response: []string{
				`{"jsonrpc":"2.0","id":1,"result":"0x1000"}`,
				`{"jsonrpc":"2.0","id":2,"result":"0x1111111111111111111111111111111111111111111111111111111111111111"}`,
			},
		},
		{
			name:     "context: replace overrides struct",
			ctx:      ContextWithLegacyGasFeeReplace(context.Background(), true),
			hijacker: &hijackLegacyGasFee{multiplier: 1.0, replace: false},
			method:   "eth_sendTransaction",
			args: []any{func() types.Transaction {
				tx := types.NewTransactionLegacy()
				tx.SetGasPrice(big.NewInt(1))
				return tx
			}()},
			request: []string{
				`{"jsonrpc":"2.0","id":1,"method":"eth_gasPrice","params":[]}`,
				`{"jsonrpc":"2.0","id":2,"method":"eth_sendTransaction","params":[{"gasPrice":"0x1000"}]}`,
			},
			response: []string{
				`{"jsonrpc":"2.0","id":1,"result":"0x1000"}`,
				`{"jsonrpc":"2.0","id":2,"result":"0x1111111111111111111111111111111111111111111111111111111111111111"}`,
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

func TestHijackDynamicGasFee(t *testing.T) {
	tc := []struct {
		name     string
		ctx      context.Context
		hijacker *hijackDynamicGasFee
		method   string
		args     []any
		request  []string
		response []string
	}{
		{
			name:     "set gas price",
			hijacker: &hijackDynamicGasFee{gasPriceMultiplier: 1.0, priorityFeePerGasMultiplier: 1.0},
			method:   "eth_sendTransaction",
			args:     []any{types.NewTransactionDynamicFee()},
			request: []string{
				`{"jsonrpc":"2.0","id":1,"method":"eth_gasPrice","params":[]}`,
				`{"jsonrpc":"2.0","id":2,"method":"eth_maxPriorityFeePerGas","params":[]}`,
				`{"jsonrpc":"2.0","id":3,"method":"eth_sendTransaction","params":[{"maxFeePerGas":"0x1000","maxPriorityFeePerGas":"0x100"}]}`,
			},
			response: []string{
				`{"jsonrpc":"2.0","id":1,"result":"0x1000"}`,
				`{"jsonrpc":"2.0","id":2,"result":"0x100"}`,
				`{"jsonrpc":"2.0","id":3,"result":"0x1111111111111111111111111111111111111111111111111111111111111111"}`,
			},
		},
		{
			name:     "context: replace overrides struct",
			ctx:      ContextWithDynamicGasFeeReplace(context.Background(), true),
			hijacker: &hijackDynamicGasFee{gasPriceMultiplier: 1.0, priorityFeePerGasMultiplier: 1.0, replace: false},
			method:   "eth_sendTransaction",
			args: []any{func() types.Transaction {
				tx := types.NewTransactionDynamicFee()
				tx.SetMaxFeePerGas(big.NewInt(1))
				tx.SetMaxPriorityFeePerGas(big.NewInt(1))
				return tx
			}()},
			request: []string{
				`{"jsonrpc":"2.0","id":1,"method":"eth_gasPrice","params":[]}`,
				`{"jsonrpc":"2.0","id":2,"method":"eth_maxPriorityFeePerGas","params":[]}`,
				`{"jsonrpc":"2.0","id":3,"method":"eth_sendTransaction","params":[{"maxFeePerGas":"0x1000","maxPriorityFeePerGas":"0x100"}]}`,
			},
			response: []string{
				`{"jsonrpc":"2.0","id":1,"result":"0x1000"}`,
				`{"jsonrpc":"2.0","id":2,"result":"0x100"}`,
				`{"jsonrpc":"2.0","id":3,"result":"0x1111111111111111111111111111111111111111111111111111111111111111"}`,
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

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

	"github.com/defiweb/go-eth/hexutil"
	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
)

func TestHijackSimulate(t *testing.T) {
	tc := []struct {
		name       string
		method     string
		args       []any
		request    []string
		response   []string
		wantResult any
		wantErr    string
	}{
		{
			name:   "send valid transaction",
			method: "eth_sendTransaction",
			args: []any{
				func() types.Transaction {
					tx := types.NewTransactionLegacy()
					_, _ = tx.DecodeRLP(hexutil.MustHexToBytes("0xf86c098504a817c800825208943535353535353535353535353535353535353535880de0b6b3a76400008025a028ef61340bd939bc2195fe537567866003e1a15d3c71ff63e1590620aa636276a067cbe9d8997f761aecb703304b3800ccf555c9f3dc64214b297fb1966a3b6d83"))
					return tx
				}(),
			},
			request: []string{
				`{"jsonrpc":"2.0","id":1,"method":"eth_call","params":[{"from":"0x9d8a62f656a8d1615c1294fd71e9cfb3e4855a4f","to":"0x3535353535353535353535353535353535353535","gas":"0x5208","gasPrice":"0x4a817c800","value":"0xde0b6b3a7640000"},"latest"]}`,
				`{"jsonrpc":"2.0","id":2,"method":"eth_sendTransaction","params":[{"chainId":"0x1","from":"0x9d8a62f656a8d1615c1294fd71e9cfb3e4855a4f","to":"0x3535353535353535353535353535353535353535","gas":"0x5208","gasPrice":"0x4a817c800","nonce":"0x9","value":"0xde0b6b3a7640000","v":"0x25","r":"0x28ef61340bd939bc2195fe537567866003e1a15d3c71ff63e1590620aa636276","s":"0x67cbe9d8997f761aecb703304b3800ccf555c9f3dc64214b297fb1966a3b6d83"}]}`,
			},
			response: []string{
				`{"jsonrpc": "2.0","id": 1,"result": "0x01"}`,
				`{"jsonrpc": "2.0","id": 1,"result": "0x1111111111111111111111111111111111111111111111111111111111111111"}`,
			},
			wantResult: "0x1111111111111111111111111111111111111111111111111111111111111111",
		},
		{
			name:   "send valid raw transaction",
			method: "eth_sendRawTransaction",
			args:   []any{types.MustBytesFromHex("0xf86c098504a817c800825208943535353535353535353535353535353535353535880de0b6b3a76400008025a028ef61340bd939bc2195fe537567866003e1a15d3c71ff63e1590620aa636276a067cbe9d8997f761aecb703304b3800ccf555c9f3dc64214b297fb1966a3b6d83")},
			request: []string{
				`{"jsonrpc":"2.0","id":1,"method":"eth_call","params":[{"from":"0x9d8a62f656a8d1615c1294fd71e9cfb3e4855a4f","to":"0x3535353535353535353535353535353535353535","gas":"0x5208","gasPrice":"0x4a817c800","value":"0xde0b6b3a7640000"},"latest"]}`,
				`{"jsonrpc":"2.0","id":2,"method":"eth_sendRawTransaction","params":["0xf86c098504a817c800825208943535353535353535353535353535353535353535880de0b6b3a76400008025a028ef61340bd939bc2195fe537567866003e1a15d3c71ff63e1590620aa636276a067cbe9d8997f761aecb703304b3800ccf555c9f3dc64214b297fb1966a3b6d83"]}`,
			},
			response: []string{
				`{"jsonrpc":"2.0","id":1,"result":"0x01"}`,
				`{"jsonrpc":"2.0","id":2,"result":"0x1111111111111111111111111111111111111111111111111111111111111111"}`,
			},
			wantResult: "0x1111111111111111111111111111111111111111111111111111111111111111",
		},
		{
			name:   "send valid private transaction",
			method: "eth_sendPrivateTransaction",
			args:   []any{types.MustBytesFromHex("0xc9808080808080808080")},
			request: []string{
				`{"jsonrpc":"2.0","id":1,"method":"eth_call","params":[{},"latest"]}`,
				`{"jsonrpc":"2.0","id":2,"method":"eth_sendPrivateTransaction","params":["0xc9808080808080808080"]}`,
			},
			response: []string{
				`{"jsonrpc":"2.0","id":1,"result":"0x01"}`,
				`{"jsonrpc":"2.0","id":2,"result":"0x1111111111111111111111111111111111111111111111111111111111111111"}`,
			},
			wantResult: "0x1111111111111111111111111111111111111111111111111111111111111111",
		},
		{
			name:   "send invalid raw transaction",
			method: "eth_sendTransaction",
			args:   []any{types.NewTransactionLegacy()},
			request: []string{
				`{"jsonrpc":"2.0","id":1,"method":"eth_call","params":[{},"latest"]}`,
			},
			response: []string{
				`{"jsonrpc":"2.0","id":1,"error":{"code":-32000,"message":"revert","data":"0x08c379a0"}}`,
			},
			wantErr: "revert",
		},
		{
			name:     "invalid raw data",
			method:   "eth_sendRawTransaction",
			args:     []any{types.MustBytesFromHex("c9")},
			request:  []string{},
			response: []string{},
			wantErr:  "failed to decode transaction",
		},
	}
	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
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
			hijacker := transport.NewHijacker(httpMock, &hijackSimulate{
				decoder: types.DefaultTransactionDecoder,
			})

			var result any
			err := hijacker.Call(ctx, &result, tt.method, tt.args...)
			assert.Len(t, tt.request, 0)
			assert.Len(t, tt.response, 0)

			if tt.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantResult, result)
		})
	}
}

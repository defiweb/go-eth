package rpc

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/defiweb/go-eth/hexutil"
	"github.com/defiweb/go-eth/rpc/transport"
	"github.com/defiweb/go-eth/types"
	"github.com/defiweb/go-eth/wallet"
)

func TestHijackSign(t *testing.T) {
	key1 := wallet.NewKeyFromBytes(hexutil.MustHexToBytes("0x01")) // 0x7e5f4552091a69125d5dfcb7b8c2659029395bdf
	key2 := wallet.NewKeyFromBytes(hexutil.MustHexToBytes("0x02")) // 0x2b5ad5c4795c026514f8317c7a215e218dccd6cf
	tc := []struct {
		name       string
		sign       *hijackSign
		method     string
		args       []any
		wantResult any
		request    []string
		response   []string
	}{
		{
			name:   "accounts",
			sign:   &hijackSign{},
			method: "eth_accounts",
			args:   []any{},
			wantResult: []types.Address{
				key1.Address(),
				key2.Address(),
			},
			request:  []string{},
			response: []string{},
		},
		{
			name:   "sign",
			sign:   &hijackSign{},
			method: "eth_sign",
			args: []any{
				key1.Address(),
				[]byte("hello"),
			},
			wantResult: types.MustSignatureFromHex("0xe5ddc160e4c8f92de507c7db9b982d4f9b7197bfa421864aeadc586bc96b09ae0ba0c5b131650ae4994cff1839341d00f3735ef5abc62ac8fe2cf50f65208e2a1b"),
			request:    []string{},
			response:   []string{},
		},
		{
			name:   "sign transaction",
			sign:   &hijackSign{},
			method: "eth_signTransaction",
			args: []any{func() types.Transaction {
				tx := types.NewTransactionLegacy()
				tx.SetFrom(types.MustAddressFromHex("0x7e5f4552091a69125d5dfcb7b8c2659029395bdf"))
				return tx
			}()},
			wantResult: hexutil.MustHexToBytes("0xf8498080808080801ba060ac8da6fc016ee487d06e93fec9768941a182638546857f6242cd86581d0174a02bea91a940d7647518956970f9ecd63847ed9aaa3a451ef7e20c89823f847082"),
			request:    []string{},
			response:   []string{},
		},
		{
			name:   "send transaction",
			sign:   &hijackSign{},
			method: "eth_sendTransaction",
			args: []any{func() types.Transaction {
				tx := types.NewTransactionLegacy()
				tx.SetFrom(types.MustAddressFromHex("0x7e5f4552091a69125d5dfcb7b8c2659029395bdf"))
				return tx
			}()},
			wantResult: types.MustHashFromHex("0x1111111111111111111111111111111111111111111111111111111111111111", types.PadNone),
			request: []string{
				`{"jsonrpc":"2.0","id":1,"method":"eth_sendRawTransaction","params":["0xf8498080808080801ba060ac8da6fc016ee487d06e93fec9768941a182638546857f6242cd86581d0174a02bea91a940d7647518956970f9ecd63847ed9aaa3a451ef7e20c89823f847082"]}`,
			},
			response: []string{
				`{"jsonrpc":"2.0","id":2,"result":"0x1111111111111111111111111111111111111111111111111111111111111111"}`,
			},
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
			tt.sign.keys = []wallet.Key{key1, key2}
			hijacker := transport.NewHijacker(httpMock, tt.sign)

			result := reflect.New(reflect.TypeOf(tt.wantResult))
			err := hijacker.Call(ctx, result.Interface(), tt.method, tt.args...)
			assert.Equal(t, tt.wantResult, result.Elem().Interface())
			assert.Len(t, tt.request, 0)
			assert.Len(t, tt.response, 0)
			require.NoError(t, err)
		})
	}
}

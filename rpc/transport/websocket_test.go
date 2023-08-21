package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"

	"github.com/defiweb/go-eth/types"
)

//nolint:funlen
func TestWebsocket(t *testing.T) {
	tests := []struct {
		asserts func(t *testing.T, ws *Websocket, reqCh, resCh chan string)
	}{
		// Simple case:
		{
			asserts: func(t *testing.T, ws *Websocket, reqCh, resCh chan string) {
				go func() {
					assert.JSONEq(t,
						`{"id":1, "jsonrpc":"2.0", "method":"eth_getBalance", "params":["0x1111111111111111111111111111111111111111", "latest"]}`,
						<-reqCh,
					)
					resCh <- `{"id": 1, "result": "0x1"}`
				}()

				ctx := context.Background()
				res := &types.Number{}
				err := ws.Call(
					ctx,
					res,
					"eth_getBalance",
					types.MustAddressFromHex("0x1111111111111111111111111111111111111111"),
					types.LatestBlockNumber,
				)

				require.NoError(t, err)
				assert.Equal(t, uint64(1), res.Big().Uint64())
			},
		},
		// Error response:
		{
			asserts: func(t *testing.T, ws *Websocket, reqCh, resCh chan string) {
				go func() {
					<-reqCh
					resCh <- `{"id": 1, "error": {"code": 1, "message": "error"}}`
				}()

				ctx := context.Background()
				res := &types.Number{}
				err := ws.Call(ctx, res, "eth_call")
				assert.Error(t, err)
			},
		},
		// Timeout:
		{
			asserts: func(t *testing.T, ws *Websocket, reqCh, resCh chan string) {
				go func() {
					<-reqCh
				}()

				ctx := context.Background()
				res := &types.Number{}
				err := ws.Call(
					ctx,
					res,
					"eth_call",
				)
				assert.Error(t, err)
			},
		},
		// Subscription:
		{
			asserts: func(t *testing.T, ws *Websocket, reqCh, resCh chan string) {
				go func() {
					assert.JSONEq(t,
						`{"id":1, "jsonrpc":"2.0", "method":"eth_subscribe", "params":["eth_sub", "foo", "bar"]}`,
						<-reqCh,
					)
					resCh <- `{"id":1, "result":"0xff"}`
				}()

				ctx := context.Background()
				ch, id, err := ws.Subscribe(ctx, "eth_sub", "foo", "bar")
				require.NoError(t, err)

				go func() {
					resCh <- `{"jsonrpc":"2.0", "method":"eth_subscribe", "params": {"subscription":"0xff", "result":"foo"}}`
					resCh <- `{"jsonrpc":"2.0", "method":"eth_subscribe", "params": {"subscription":"0xff", "result":"bar"}}`
				}()

				assert.Equal(t, "0xff", id)
				assert.Equal(t, json.RawMessage(`"foo"`), <-ch)
				assert.Equal(t, json.RawMessage(`"bar"`), <-ch)

				go func() {
					assert.JSONEq(t,
						`{"id":2, "jsonrpc":"2.0", "method":"eth_unsubscribe", "params":["0xff"]}`,
						<-reqCh,
					)
					resCh <- `{"id":2}`
				}()

				err = ws.Unsubscribe(ctx, id)
				require.NoError(t, err)

				// Channel must be closed after unsubscribe.
				_, ok := <-ch
				require.False(t, ok)
			},
		},
	}
	for n, tt := range tests {
		t.Run(fmt.Sprintf("case-%d", n+1), func(t *testing.T) {
			wg := sync.WaitGroup{}
			reqCh := make(chan string)     // Received requests.
			resCh := make(chan string)     // Responses from server.
			closeCh := make(chan struct{}) // Stops the server.

			// Websocket server.
			server := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Handle websocket requests.
				ctx := context.Background()
				conn, err := websocket.Accept(w, r, nil)
				if err != nil {
					require.NoError(t, err)
				}

				// Request reader.
				wg.Add(1)
				go func() {
					defer wg.Done()
					for {
						var req json.RawMessage
						if err := wsjson.Read(ctx, conn, &req); err != nil {
							if errors.As(err, &websocket.CloseError{}) {
								return
							}
							require.NoError(t, err)
						}
						reqCh <- string(req)
					}
				}()

				// Response writer.
				wg.Add(1)
				go func() {
					defer wg.Done()
					for {
						select {
						case <-closeCh:
							return
						case res := <-resCh:
							if err := wsjson.Write(ctx, conn, json.RawMessage(res)); err != nil {
								if errors.As(err, &websocket.CloseError{}) {
									return
								}
								require.NoError(t, err)
							}
						}
					}
				}()

				// Close the connection after the test.
				<-closeCh
				conn.Close(websocket.StatusNormalClosure, "")
			})}

			// Start HTTP server.
			ln, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				require.NoError(t, err)
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := server.Serve(ln); err != nil {
					if !errors.Is(err, http.ErrServerClosed) {
						require.NoError(t, err)
					}
				}
			}()

			// Create a websocket client.
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()
			ws, err := NewWebsocket(WebsocketOptions{
				Context: ctx,
				URL:     "ws://" + ln.Addr().String(),
				Timout:  time.Second,
			})
			require.NoError(t, err)

			// Run the test.
			tt.asserts(t, ws, reqCh, resCh)

			// Stop the server.
			close(closeCh)
			_ = server.Close()
			wg.Wait()
		})
	}
}

package testutil

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"sync"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"

	"github.com/defiweb/go-eth/rpc/transport"
)

type WebsocketMock struct {
	*transport.Websocket
	wg sync.WaitGroup

	RequestCh  chan string // Request that was sent.
	ResponseCh chan string // Response that will be returned.
}

func NewWebsocketMock(ctx context.Context) *WebsocketMock {
	w := &WebsocketMock{
		RequestCh:  make(chan string),
		ResponseCh: make(chan string),
	}

	// Websocket server.
	server := &http.Server{Handler: http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Handle websocket requests.
		conn, err := websocket.Accept(res, req, nil)
		if err != nil {
			panic(err)
		}

		// Request reader.
		go func() {
			w.wg.Add(1)
			defer w.wg.Done()
			for {
				var req json.RawMessage
				if err := wsjson.Read(ctx, conn, &req); err != nil {
					if errors.As(err, &websocket.CloseError{}) {
						return
					}
					panic(err)
				}
				w.RequestCh <- string(req)
			}
		}()

		// Response writer.
		go func() {
			w.wg.Add(1)
			defer w.wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case res := <-w.ResponseCh:
					if err := wsjson.Write(ctx, conn, json.RawMessage(res)); err != nil {
						if errors.As(err, &websocket.CloseError{}) {
							return
						}
						panic(err)
					}
				}
			}
		}()

		// Close the connection after the test.
		<-ctx.Done()
		_ = conn.Close(websocket.StatusNormalClosure, "")
	})}

	// Close HTTP server after the test.
	go func() {
		<-ctx.Done()
		_ = server.Close()
		w.wg.Wait()
	}()

	// Start HTTP server.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	go func() {
		w.wg.Add(1)
		defer w.wg.Done()
		if err := server.Serve(ln); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				panic(err)
			}
		}
	}()

	w.Websocket, _ = transport.NewWebsocket(transport.WebsocketOptions{
		Context: ctx,
		URL:     "ws://" + ln.Addr().String(),
		Timout:  time.Second,
	})

	return w
}

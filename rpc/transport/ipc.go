package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

// IPC is a Transport implementation that uses the IPC protocol.
type IPC struct {
	*stream
	conn net.Conn
}

// IPCOptions contains options for the IPC transport.
type IPCOptions struct {
	// Context used to close the connection.
	Context context.Context

	// Path is the path to the IPC socket.
	Path string

	// Timeout is the timeout for the websocket requests. Default is 60s.
	Timout time.Duration

	// ErrorCh is an optional channel used to report errors.
	ErrorCh chan error
}

// NewIPC creates a new IPC instance.
func NewIPC(opts IPCOptions) (*IPC, error) {
	var d net.Dialer
	conn, err := d.DialContext(opts.Context, "unix", opts.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to dial IPC: %w", err)
	}
	if opts.Context == nil {
		return nil, errors.New("context cannot be nil")
	}
	if opts.Timout == 0 {
		opts.Timout = 60 * time.Second
	}
	i := &IPC{
		stream: &stream{
			ctx:     opts.Context,
			errCh:   opts.ErrorCh,
			timeout: opts.Timout,
		},
		conn: conn,
	}
	i.stream.initStream()
	go i.readerRoutine()
	go i.writerRoutine()
	return i, nil
}

func (i *IPC) readerRoutine() {
	dec := json.NewDecoder(i.conn)
	for {
		var res rpcResponse
		if err := dec.Decode(&res); err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			if errors.Is(err, io.EOF) {
				return
			}
			i.errCh <- err
		}
		i.readerCh <- res
	}
}

func (i *IPC) writerRoutine() {
	enc := json.NewEncoder(i.conn)
	for {
		select {
		case <-i.ctx.Done():
			return
		case req := <-i.stream.writerCh:
			if err := enc.Encode(req); err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				if errors.Is(err, io.EOF) {
					return
				}
				i.stream.errCh <- err
			}
		}
	}
}

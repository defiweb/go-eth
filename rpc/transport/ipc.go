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

// IPC is a [Transport] implementation that uses the IPC protocol.
type IPC struct {
	*stream
	conn net.Conn
}

// IPCOptions contains options for the [IPC] transport.
type IPCOptions struct {
	// Context used to close the connection.
	Context context.Context

	// Path is the path to the IPC socket.
	Path string

	// Timeout is the timeout for the IPC requests. Default is 60s.
	Timeout time.Duration

	// ErrorCh is an optional channel used to report errors.
	ErrorCh chan error
}

// NewIPC creates a new [IPC] instance.
func NewIPC(opts IPCOptions) (*IPC, error) {
	if opts.Context == nil {
		return nil, errors.New("context cannot be nil")
	}
	var d net.Dialer
	conn, err := d.DialContext(opts.Context, "unix", opts.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to dial IPC: %w", err)
	}
	if opts.Timeout == 0 {
		opts.Timeout = 60 * time.Second
	}
	i := &IPC{
		stream: &stream{
			ctx:     opts.Context,
			errCh:   opts.ErrorCh,
			timeout: opts.Timeout,
		},
		conn: conn,
	}
	i.initStream()
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
		select {
		case i.readerCh <- res:
		case <-i.ctx.Done():
			return
		}
	}
}

func (i *IPC) writerRoutine() {
	enc := json.NewEncoder(i.conn)
	for {
		select {
		case req := <-i.writerCh:
			if err := enc.Encode(req); err != nil {
				if i.errCh == nil {
					return
				}
				if errors.Is(err, context.Canceled) {
					return
				}
				if errors.Is(err, io.EOF) {
					return
				}
				select {
				case i.errCh <- err:
				case <-i.ctx.Done():
					return
				}
			}
		case <-i.ctx.Done():
			return
		}
	}
}

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
	stream
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
	if opts.Timeout == 0 {
		opts.Timeout = time.Minute
	}
	var d net.Dialer
	conn, err := d.DialContext(opts.Context, "unix", opts.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to dial IPC: %w", err)
	}
	ipc := &IPC{conn: conn}
	ipc.stream.initStream(
		opts.Context,
		withStreamTimeout(opts.Timeout),
		withStreamErrorCh(opts.ErrorCh),
	)
	go ipc.readerRoutine()
	go ipc.writerRoutine()
	return ipc, nil
}

func (i *IPC) readerRoutine() {
	dec := json.NewDecoder(i.conn)
	for {
		var res rpcResponse
		if err := dec.Decode(&res); err != nil {
			if i.ctx.Err() != nil {
				return
			}
			if errors.Is(err, io.EOF) {
				return
			}
			i.error(fmt.Errorf("ipc reading error: %w", err))
			return
		}
		i.read(res)
	}
}

func (i *IPC) writerRoutine() {
	defer i.close()
	enc := json.NewEncoder(i.conn)
	for {
		req, ok := i.write()
		if !ok {
			return
		}
		if err := i.conn.SetWriteDeadline(time.Now().Add(i.timeout)); err != nil {
			i.error(fmt.Errorf("ipc writing error: %w", err))
			continue
		}
		if err := enc.Encode(req); err != nil {
			i.error(fmt.Errorf("ipc writing error: %w", err))
			continue
		}
	}
}

func (i *IPC) close() {
	if err := i.conn.Close(); err != nil {
		i.error(fmt.Errorf("ipc close error: %w", err))
	}
}

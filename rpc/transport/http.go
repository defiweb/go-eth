package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync/atomic"
)

// HTTP is a Transport implementation that uses the HTTP protocol.
type HTTP struct {
	opts HTTPOptions
	id   uint64
}

// HTTPOptions contains options for the HTTP transport.
type HTTPOptions struct {
	// URL of the HTTP endpoint.
	URL string

	// HTTPClient is used for the connection.
	HTTPClient *http.Client

	// HTTPHeader specifies the HTTP headers included in RPC requests.
	HTTPHeader http.Header
}

// NewHTTP creates a new HTTP instance.
func NewHTTP(opts HTTPOptions) (*HTTP, error) {
	if opts.URL == "" {
		return nil, errors.New("URL cannot be empty")
	}
	if opts.HTTPClient == nil {
		opts.HTTPClient = http.DefaultClient
	}
	return &HTTP{opts: opts}, nil
}

// Call implements the Transport interface.
func (h *HTTP) Call(ctx context.Context, result any, method string, args ...any) error {
	id := atomic.AddUint64(&h.id, 1)
	rpcReq, err := newRPCRequest(&id, method, args)
	if err != nil {
		return err
	}
	httpBody, err := json.Marshal(rpcReq)
	if err != nil {
		return err
	}
	httpReq, err := http.NewRequestWithContext(ctx, "POST", h.opts.URL, bytes.NewReader(httpBody))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	for k, v := range h.opts.HTTPHeader {
		httpReq.Header[k] = v
	}
	httpRes, err := h.opts.HTTPClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpRes.Body.Close()
	rpcRes := &rpcResponse{}
	if err := json.NewDecoder(httpRes.Body).Decode(rpcRes); err != nil {
		return err
	}
	if rpcRes.Error != nil {
		return rpcRes.Error
	}
	if result == nil {
		return nil
	}
	if err := json.Unmarshal(rpcRes.Result, result); err != nil {
		return err
	}
	return nil
}

package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"sync/atomic"
)

type HTTP struct {
	url     string
	client  *http.Client
	headers map[string]string
	id      uint64
}

type HTTPOptions struct {
	// HTTPClient is used for the connection.
	HTTPCLient *http.Client

	// HTTPHeader specifies the HTTP headers included in RPC requests.
	HTTPHeader http.Header
}

func NewHTTP(url string) *HTTP {
	return &HTTP{
		url:    url,
		client: http.DefaultClient,
	}
}

func NewHTTPWithClient(url string, client *http.Client) *HTTP {
	return &HTTP{
		url:    url,
		client: client,
	}
}

func (h *HTTP) SetHeader(key, value string) {
	if h.headers == nil {
		h.headers = make(map[string]string)
	}
	h.headers[key] = value
}

func (h *HTTP) Call(ctx context.Context, result any, method string, args ...any) error {
	id := atomic.AddUint64(&h.id, 1)
	rpcReq := rpcRequest{
		JSONRPC: "2.0",
		ID:      &id,
		Method:  method,
		Params:  json.RawMessage("[]"),
	}
	if len(args) > 0 {
		params, err := json.Marshal(args)
		if err != nil {
			return err
		}
		rpcReq.Params = params
	}
	httpBody, err := json.Marshal(rpcReq)
	if err != nil {
		return err
	}
	httpReq, err := http.NewRequestWithContext(ctx, "POST", h.url, bytes.NewReader(httpBody))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	for k, v := range h.headers {
		httpReq.Header.Set(k, v)
	}
	httpRes, err := h.client.Do(httpReq)
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
	if err := json.Unmarshal(rpcRes.Result, result); err != nil {
		return err
	}
	return nil
}

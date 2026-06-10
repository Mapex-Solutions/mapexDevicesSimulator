package dispatch

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"strings"

	"simulator/service/src/modules/engine/application/ports"
)

// Compile-time proof the adapter satisfies the Dispatcher port.
var _ ports.Dispatcher = (*httpDispatcher)(nil)

// NewHTTP builds the HTTP dispatcher.
func NewHTTP() ports.Dispatcher {
	return &httpDispatcher{client: &http.Client{Timeout: httpTimeout}}
}

// Protocol identifies this dispatcher in the registry.
func (d *httpDispatcher) Protocol() string { return "http" }

// Dispatch sends the rendered payload as an HTTP request and reports the status.
func (d *httpDispatcher) Dispatch(ctx context.Context, req ports.DispatchRequest) ports.DispatchResult {
	method := req.Method
	if method == "" {
		method = http.MethodPost
	}
	httpReq, err := http.NewRequestWithContext(ctx, method, req.URL, strings.NewReader(req.Payload))
	if err != nil {
		return ports.DispatchResult{Err: err}
	}
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}
	if req.Payload != "" && httpReq.Header.Get("Content-Type") == "" {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	resp, err := d.client.Do(httpReq)
	if err != nil {
		return ports.DispatchResult{Err: err}
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	ok := resp.StatusCode >= 200 && resp.StatusCode < 300
	return ports.DispatchResult{OK: ok, Status: strconv.Itoa(resp.StatusCode)}
}

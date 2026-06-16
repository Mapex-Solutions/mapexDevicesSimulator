package utils

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
)

// EchoRequest is one request the echo target captured.
type EchoRequest struct {
	Method string
	Path   string
	Body   string
}

// Echo is a test HTTP target an HTTP device fires against. It records every
// request and replies 200 with a small JSON body, so a journey can assert both
// that the device reached it and that the simulator captured the response.
type Echo struct {
	server *httptest.Server
	mu     sync.Mutex
	reqs   []EchoRequest
}

// StartEcho boots an echo target on a random localhost port. The caller stops
// it (typically via t.Cleanup(echo.Close)).
func StartEcho() *Echo {
	e := &Echo{}
	e.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		e.mu.Lock()
		e.reqs = append(e.reqs, EchoRequest{Method: r.Method, Path: r.URL.Path, Body: string(body)})
		e.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	return e
}

// URL is the base URL the device's HTTP config should target.
func (e *Echo) URL() string { return e.server.URL }

// Close shuts the echo target down.
func (e *Echo) Close() { e.server.Close() }

// Requests returns a copy of the captured requests.
func (e *Echo) Requests() []EchoRequest {
	e.mu.Lock()
	defer e.mu.Unlock()
	out := make([]EchoRequest, len(e.reqs))
	copy(out, e.reqs)
	return out
}

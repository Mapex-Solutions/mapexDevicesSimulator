package dispatch

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"simulator/service/src/modules/engine/application/ports"
)

func TestHTTPDispatcher_Sends(t *testing.T) {
	var method, apiKey, body string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		apiKey = r.Header.Get("X-API-Key")
		b, _ := io.ReadAll(r.Body)
		body = string(b)
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	res := NewHTTP().Dispatch(context.Background(), ports.DispatchRequest{
		URL:     srv.URL,
		Method:  http.MethodPost,
		Headers: map[string]string{"X-API-Key": "k"},
		Payload: `{"t":1}`,
	})
	if !res.OK || res.Status != "202" {
		t.Fatalf("result = %+v", res)
	}
	if res.Response != `{"ok":true}` {
		t.Fatalf("expected captured response body, got %q", res.Response)
	}
	if method != http.MethodPost || apiKey != "k" || body != `{"t":1}` {
		t.Fatalf("server got method=%s apiKey=%s body=%s", method, apiKey, body)
	}
}

func TestHTTPDispatcher_Non2xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	res := NewHTTP().Dispatch(context.Background(), ports.DispatchRequest{URL: srv.URL})
	if res.OK || res.Status != "500" {
		t.Fatalf("expected not-ok 500, got %+v", res)
	}
}

func TestHTTPDispatcher_Unreachable(t *testing.T) {
	res := NewHTTP().Dispatch(context.Background(), ports.DispatchRequest{URL: "http://127.0.0.1:1/x"})
	if res.OK || res.Err == nil {
		t.Fatalf("expected an error on unreachable url, got %+v", res)
	}
}

func TestRegistry_Resolves(t *testing.T) {
	r := NewRegistry()
	for _, p := range []string{"http", "mqtt"} {
		if _, ok := r.For(p); !ok {
			t.Fatalf("registry missing %q", p)
		}
	}
	if _, ok := r.For("lorawan"); ok {
		t.Fatal("lorawan should not be registered yet")
	}
}

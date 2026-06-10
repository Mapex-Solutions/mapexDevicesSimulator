package hub

import (
	"encoding/json"
	"testing"
	"time"

	"simulator/service/src/modules/console/application/dtos"
)

// recv reads one frame off a client within a short window, failing on timeout.
func recv(t *testing.T, c *Client) dtos.ConsoleMessage {
	t.Helper()
	select {
	case data := <-c.Out():
		var m dtos.ConsoleMessage
		if err := json.Unmarshal(data, &m); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		return m
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for frame")
		return dtos.ConsoleMessage{}
	}
}

func TestHub_BroadcastAndUnregister(t *testing.T) {
	h := New()
	c1 := h.Register()
	c2 := h.Register()

	h.Publish(dtos.ConsoleMessage{ID: "1", Summary: "hello"})
	if got := recv(t, c1); got.ID != "1" || got.Summary != "hello" {
		t.Fatalf("c1 frame = %+v", got)
	}
	if got := recv(t, c2); got.ID != "1" {
		t.Fatalf("c2 frame = %+v", got)
	}

	// After unregister, c1 stops receiving; c2 still does.
	h.Unregister(c1)
	h.Publish(dtos.ConsoleMessage{ID: "2"})
	if got := recv(t, c2); got.ID != "2" {
		t.Fatalf("c2 second frame = %+v", got)
	}
	if _, open := <-c1.Out(); open {
		t.Fatal("c1 channel should be closed after unregister")
	}

	// Unregister is idempotent.
	h.Unregister(c1)
}

func TestHub_SlowClientDropsInsteadOfBlocking(t *testing.T) {
	h := New()
	c := h.Register()

	// Flood past the buffer; Publish must never block even though nobody reads.
	done := make(chan struct{})
	go func() {
		for i := 0; i < clientBuffer*2; i++ {
			h.Publish(dtos.ConsoleMessage{ID: "x"})
		}
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Publish blocked on a slow client")
	}
	h.Unregister(c)
}

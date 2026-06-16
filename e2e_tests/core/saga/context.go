package saga

import (
	"context"
	"sync"
	"testing"

	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/httpclient"
)

// Context is threaded through every Item executed by Run. It carries the
// testing handle, the cancellation context, the simulator client, and a
// free-form bag where steps publish outputs for downstream consumers (e.g.
// CreateDevice writes the device id; FireDevice reads it without taking the
// value through closure plumbing).
type Context struct {
	// T is the testing handle. Steps that need to fail the test directly
	// (rare; prefer returning errors from Do/Check) call methods on T.
	T *testing.T

	// Stdctx is the standard library context used by every HTTP call inside
	// steps. Cancellation propagates through the journey.
	Stdctx context.Context

	// Clients holds the simulator client steps drive the API with.
	Clients ClientSet

	// RunID is a per-journey unique tag injected into payload identifiers
	// (device name, deviceId) so cleanup-by-prefix works and parallel runs
	// do not collide.
	RunID string

	bag map[string]any
	mu  sync.Mutex
}

// ClientSet groups the HTTP clients a journey speaks to. The simulator is a
// single service, so the set is just the sidecar client; the type is kept so
// the shape matches the platform e2e and can grow if the simulator gains more
// endpoints.
type ClientSet struct {
	// Sim points at the simulator sidecar base URL. Every step drives it
	// through the Raw method so it can assert on resp.StatusCode directly.
	Sim *httpclient.HTTPClient
}

// Set publishes a value to the bag for downstream Items.
func (c *Context) Set(key string, val any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.bag[key] = val
}

// Get returns the value previously published to key. The bool is false when the
// key was never set.
func (c *Context) Get(key string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.bag[key]
	return v, ok
}

// MustGetString fetches a string from the bag, failing the test fast when the
// key is missing or holds a non-string value. Steps that consume bag inputs use
// this to express the contract: "I require this key to exist before I can
// execute" — surfaced as a clear test failure rather than a nil dereference.
func (c *Context) MustGetString(key string) string {
	v, ok := c.Get(key)
	if !ok {
		c.T.Fatalf("[SAGA] missing required bag key %q (step out of order?)", key)
	}
	s, ok := v.(string)
	if !ok {
		c.T.Fatalf("[SAGA] bag key %q is not a string (got %T)", key, v)
	}
	return s
}

// newContext is invoked by Run to materialize the per-journey Context. Kept
// package-private so callers cannot bypass the runner contract.
func newContext(t *testing.T, ctx context.Context, runID string, clients ClientSet) *Context {
	return &Context{
		T:       t,
		Stdctx:  ctx,
		Clients: clients,
		RunID:   runID,
		bag:     make(map[string]any),
	}
}

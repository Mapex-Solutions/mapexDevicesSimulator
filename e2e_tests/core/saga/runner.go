package saga

import (
	"context"
	"testing"
	"time"

	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/httpclient"
)

// Run executes a journey: walks Items in order, records executed Steps for
// rollback, and runs every registered Compensate in reverse on completion.
// Compensations run on both success and failure so the simulator ends in the
// same state where the next journey can start.
//
// runID is the per-journey unique tag callers inject into payload identifiers.
// When empty, Run synthesizes one from the current time so every execution
// still gets prefix-based isolation.
func Run(t *testing.T, ctx context.Context, runID string, clients ClientSet, items ...Item) {
	t.Helper()
	if runID == "" {
		runID = time.Now().UTC().Format("20060102-150405")
	}

	sctx := newContext(t, ctx, runID, clients)
	t.Logf("[SAGA] start runID=%s items=%d", runID, len(items))

	executed := make([]Item, 0, len(items))
	failed := false

	// Wrapped in a closure so the deferred call reads the FINAL executed slice.
	// `defer rollback(sctx, executed)` would capture the empty slice as it stands
	// at the defer statement, and no compensation would ever run.
	defer func() { rollback(sctx, executed) }()

	for i, item := range items {
		t.Logf("[SAGA] %d/%d %s", i+1, len(items), item.GetName())
		if err := item.Execute(sctx); err != nil {
			t.Errorf("[SAGA] %s failed: %v", item.GetName(), err)
			failed = true
			break
		}
		executed = append(executed, item)
	}

	if !failed {
		t.Logf("[SAGA] all %d items passed; running compensations", len(executed))
	}
}

// rollback walks the executed list in reverse and invokes Rollback on every
// item. Compensation failures are logged but do not abort the rollback — the
// goal is best-effort cleanup so subsequent runs can proceed even when one
// compensation step is unhappy.
func rollback(c *Context, executed []Item) {
	for i := len(executed) - 1; i >= 0; i-- {
		item := executed[i]
		if err := item.Rollback(c); err != nil {
			c.T.Logf("[SAGA] compensation %s failed (non-fatal): %v", item.GetName(), err)
		}
	}
}

// NewClientSet builds the simulator client set against the sidecar base URL.
func NewClientSet(simURL string) ClientSet {
	return ClientSet{
		Sim: httpclient.New(httpclient.Config{BaseURL: simURL}),
	}
}

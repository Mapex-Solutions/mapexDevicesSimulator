package e2e

import (
	"net/url"
	"testing"
)

// TestTree_CursorNext_FirstToLast walks the cursor at page size 1 across all of
// the seeded device's logs and confirms every entry is visited exactly once —
// the keyset pagination is complete and stable.
func TestTree_CursorNext_FirstToLast(t *testing.T) {
	deviceID := seedLogs(t, listFixtureCount)

	base := url.Values{}
	base.Set("device", deviceID)
	visited := walkCursor(t, base, 1)

	if len(visited) != listFixtureCount {
		t.Fatalf("cursor walk visited %d logs, want %d", len(visited), listFixtureCount)
	}
	ids := make(map[string]bool, len(visited))
	for _, l := range visited {
		if l.DeviceID != deviceID {
			t.Fatalf("cursor walk surfaced a foreign log: %s", l.DeviceID)
		}
		if ids[l.ID] {
			t.Fatalf("cursor walk visited log %s twice", l.ID)
		}
		ids[l.ID] = true
	}
}

// TestTree_CursorBack_LastToFirst would walk the cursor backward, but the logs
// endpoint exposes only a forward nextCursor. Reactivate when GET /api/logs
// gains a backward cursor (e.g. a `before` token).
func TestTree_CursorBack_LastToFirst(t *testing.T) {
	t.Skip("logs endpoint has no backward cursor; reactivate when a `before` token is added")
}

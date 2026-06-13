package sqlite

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	sqliteManager "github.com/Mapex-Solutions/mapexGoKit/infrastructure/sqlite/manager"

	"simulator/service/src/modules/logs/domain/entities"
	"simulator/service/src/modules/logs/domain/repositories"
)

func setup(t *testing.T) repositories.LogRepository {
	t.Helper()
	mgr, err := sqliteManager.New(sqliteManager.Config{Path: filepath.Join(t.TempDir(), "logs.db")})
	if err != nil {
		t.Fatalf("manager: %v", err)
	}
	t.Cleanup(func() { _ = mgr.Close() })
	if err := EnsureSchema(context.Background(), mgr); err != nil {
		t.Fatalf("ensure schema: %v", err)
	}
	return New(mgr)
}

func TestLogRepository_ListPage(t *testing.T) {
	ctx := context.Background()
	repo := setup(t)
	base := time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC)

	seed := []entities.Log{
		{Protocol: "mqtt", DeviceID: "d1", DeviceName: "Greenhouse", EventName: "Telemetry", Direction: "up", Kind: "data", Summary: "humidity 63", Payload: "{}", Response: "200", Created: base},
		{Protocol: "http", DeviceID: "d2", DeviceName: "Edge", EventName: "Ingest", Direction: "up", Kind: "data", Summary: "POST /v1/ingest", Payload: "{}", Created: base.Add(time.Second)},
		{Protocol: "mqtt", DeviceID: "d1", DeviceName: "Greenhouse", EventName: "Telemetry", Direction: "down", Kind: "status", Summary: "temp", Payload: "{}", Created: base.Add(2 * time.Second)},
	}
	for i := range seed {
		if _, err := repo.Insert(ctx, &seed[i]); err != nil {
			t.Fatalf("insert: %v", err)
		}
	}

	// No filter: all 3, newest first, no further page.
	items, next, err := repo.ListPage(ctx, repositories.LogFilter{})
	if err != nil {
		t.Fatalf("ListPage: %v", err)
	}
	if len(items) != 3 || next != "" {
		t.Fatalf("len/next = %d/%q, want 3/\"\"", len(items), next)
	}
	if items[0].Summary != "temp" {
		t.Fatalf("newest-first wrong: %q", items[0].Summary)
	}
	// event name and response round-trip (oldest row carried both).
	if items[2].EventName != "Telemetry" || items[2].Response != "200" {
		t.Fatalf("event/response not stored: %+v", items[2])
	}

	// Equality filter.
	items, _, _ = repo.ListPage(ctx, repositories.LogFilter{Protocol: "mqtt"})
	if len(items) != 2 {
		t.Fatalf("protocol=mqtt len = %d, want 2", len(items))
	}

	// Event name match.
	items, _, _ = repo.ListPage(ctx, repositories.LogFilter{Event: "Ingest"})
	if len(items) != 1 || items[0].DeviceName != "Edge" {
		t.Fatalf("event=Ingest items = %+v", items)
	}

	// Date range bounds the message time (inclusive) to the oldest row only.
	stamp := base.Format(timeLayout)
	items, _, _ = repo.ListPage(ctx, repositories.LogFilter{DateFrom: stamp, DateTo: stamp})
	if len(items) != 1 || items[0].Summary != "humidity 63" {
		t.Fatalf("date range items = %+v", items)
	}

	// Free-text q matches the summary.
	items, _, _ = repo.ListPage(ctx, repositories.LogFilter{Q: "humid"})
	if len(items) != 1 || items[0].DeviceName != "Greenhouse" {
		t.Fatalf("q=humid items = %+v", items)
	}
}

func TestLogRepository_CursorPagination(t *testing.T) {
	ctx := context.Background()
	repo := setup(t)
	base := time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC)

	seed := []entities.Log{
		{Protocol: "mqtt", DeviceID: "d1", DeviceName: "A", Direction: "up", Kind: "data", Summary: "first", Payload: "{}", Created: base},
		{Protocol: "mqtt", DeviceID: "d1", DeviceName: "A", Direction: "up", Kind: "data", Summary: "second", Payload: "{}", Created: base.Add(time.Second)},
		{Protocol: "mqtt", DeviceID: "d1", DeviceName: "A", Direction: "up", Kind: "data", Summary: "third", Payload: "{}", Created: base.Add(2 * time.Second)},
	}
	for i := range seed {
		if _, err := repo.Insert(ctx, &seed[i]); err != nil {
			t.Fatalf("insert: %v", err)
		}
	}

	// Walk one row per page, newest first; the last page has no next cursor.
	want := []string{"third", "second", "first"}
	cursor := ""
	for i, summary := range want {
		page, next, err := repo.ListPage(ctx, repositories.LogFilter{Limit: 1, Cursor: cursor})
		if err != nil {
			t.Fatalf("page %d: %v", i, err)
		}
		if len(page) != 1 || page[0].Summary != summary {
			t.Fatalf("page %d = %+v, want %q", i, page, summary)
		}
		lastPage := i == len(want)-1
		if lastPage && next != "" {
			t.Fatalf("last page should have no next cursor, got %q", next)
		}
		if !lastPage && next == "" {
			t.Fatalf("page %d should have a next cursor", i)
		}
		cursor = next
	}
}

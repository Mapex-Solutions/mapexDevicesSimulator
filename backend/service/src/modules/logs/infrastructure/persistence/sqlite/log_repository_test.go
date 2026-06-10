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
	if err := mgr.Migrate(context.Background(), Migrations...); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return New(mgr)
}

func TestLogRepository_ListPage(t *testing.T) {
	ctx := context.Background()
	repo := setup(t)
	base := time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC)

	seed := []entities.Log{
		{Protocol: "mqtt", DeviceID: "d1", DeviceName: "Greenhouse", Direction: "up", Kind: "data", Summary: "humidity 63", Payload: "{}", Created: base},
		{Protocol: "http", DeviceID: "d2", DeviceName: "Edge", Direction: "up", Kind: "data", Summary: "POST /v1/ingest", Payload: "{}", Created: base.Add(time.Second)},
		{Protocol: "mqtt", DeviceID: "d1", DeviceName: "Greenhouse", Direction: "down", Kind: "status", Summary: "temp", Payload: "{}", Created: base.Add(2 * time.Second)},
	}
	for i := range seed {
		if _, err := repo.Insert(ctx, &seed[i]); err != nil {
			t.Fatalf("insert: %v", err)
		}
	}

	// No filter: all 3, newest first.
	items, total, err := repo.ListPage(ctx, repositories.LogFilter{})
	if err != nil {
		t.Fatalf("ListPage: %v", err)
	}
	if total != 3 || len(items) != 3 {
		t.Fatalf("total/len = %d/%d, want 3/3", total, len(items))
	}
	if items[0].Summary != "temp" {
		t.Fatalf("newest-first wrong: %q", items[0].Summary)
	}

	// Equality filter.
	_, total, _ = repo.ListPage(ctx, repositories.LogFilter{Protocol: "mqtt"})
	if total != 2 {
		t.Fatalf("protocol=mqtt total = %d, want 2", total)
	}

	// Free-text q matches the summary.
	items, total, _ = repo.ListPage(ctx, repositories.LogFilter{Q: "humid"})
	if total != 1 || items[0].DeviceName != "Greenhouse" {
		t.Fatalf("q=humid total = %d items=%+v", total, items)
	}

	// Pagination: limit 1, offset 1 returns the 2nd-newest, total still counts all.
	items, total, _ = repo.ListPage(ctx, repositories.LogFilter{Limit: 1, Offset: 1})
	if total != 3 || len(items) != 1 || items[0].Summary != "POST /v1/ingest" {
		t.Fatalf("pagination wrong: total=%d items=%+v", total, items)
	}
}

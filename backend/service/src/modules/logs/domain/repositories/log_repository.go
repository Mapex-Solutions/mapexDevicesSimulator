package repositories

import (
	"context"

	"simulator/service/src/modules/logs/domain/entities"
)

// LogFilter narrows a log query. Empty string fields are ignored. Cursor is the
// opaque keyset token of the last row of the previous page (empty for the first
// page). Q is a free-text match over summary, payload and device name; Event
// matches the event name; DateFrom/DateTo bound the message time (inclusive).
type LogFilter struct {
	Limit     int
	Cursor    string
	Protocol  string
	Kind      string
	Direction string
	Device    string
	Event     string
	DateFrom  string
	DateTo    string
	Q         string
}

// LogRepository is the persistence port for the device message history.
type LogRepository interface {
	// ListPage returns the matching logs newest-first and the cursor for the next
	// page (empty when there are no more rows). Pagination is keyset-based on
	// (created, id), so it stays stable as new rows arrive at the top.
	ListPage(ctx context.Context, f LogFilter) ([]entities.Log, string, error)
	// Insert persists one message; used by the simulation engine, not the HTTP API.
	Insert(ctx context.Context, l *entities.Log) (*entities.Log, error)
}

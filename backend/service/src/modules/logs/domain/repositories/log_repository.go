package repositories

import (
	"context"

	"simulator/service/src/modules/logs/domain/entities"
)

// LogFilter narrows a log query. Empty string fields are ignored; Q is a
// free-text match over summary, payload and device name.
type LogFilter struct {
	Limit     int
	Offset    int
	Protocol  string
	Kind      string
	Direction string
	Device    string
	Q         string
}

// LogRepository is the persistence port for the device message history.
type LogRepository interface {
	// ListPage returns the matching logs (newest first) and the total matching
	// count for pagination.
	ListPage(ctx context.Context, f LogFilter) ([]entities.Log, int, error)
	// Insert persists one message; used by the simulation engine, not the HTTP API.
	Insert(ctx context.Context, l *entities.Log) (*entities.Log, error)
}

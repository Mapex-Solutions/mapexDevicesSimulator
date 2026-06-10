package ports

import (
	"context"

	"simulator/service/src/modules/logs/application/dtos"
)

// LogsServicePort is the driving port for the read-only log history.
type LogsServicePort interface {
	List(ctx context.Context, q *dtos.LogQuery) (*dtos.LogPage, error)
}

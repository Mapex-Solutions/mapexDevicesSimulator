package ports

import (
	"context"

	"simulator/service/src/modules/logs/application/dtos"
)

// LogWriter is the write path the simulation engine uses to persist a device
// message (separate from the read-only LogsServicePort the HTTP API exposes).
type LogWriter interface {
	Append(ctx context.Context, in *dtos.LogInput) error
}

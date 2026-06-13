package services

import (
	"context"

	"simulator/service/src/modules/logs/application/di"
	"simulator/service/src/modules/logs/application/dtos"
	"simulator/service/src/modules/logs/application/ports"
)

// Compile-time checks that the service satisfies both the read and write ports.
var (
	_ ports.LogsServicePort = (*LogsService)(nil)
	_ ports.LogWriter       = (*LogsService)(nil)
)

// New builds the logs service for the read port.
func New(deps di.LogsServiceDI) ports.LogsServicePort {
	return &LogsService{deps: deps}
}

// NewWriter builds the logs service for the write port (used by the engine).
func NewWriter(deps di.LogsServiceDI) ports.LogWriter {
	return &LogsService{deps: deps}
}

// List returns a page of the message history matching the query filters.
func (s *LogsService) List(ctx context.Context, q *dtos.LogQuery) (*dtos.LogPage, error) {
	filter := s.filterFromQuery(q)
	items, next, err := s.deps.Repo.ListPage(ctx, filter)
	if err != nil {
		return nil, err
	}
	return s.buildPage(items, next)
}

// Append persists one device message.
func (s *LogsService) Append(ctx context.Context, in *dtos.LogInput) error {
	entity, err := s.entityFromInput(in)
	if err != nil {
		return err
	}
	_, err = s.deps.Repo.Insert(ctx, entity)
	return err
}

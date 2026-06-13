package services

import (
	"github.com/Mapex-Solutions/mapexGoKit/utils/mapper"

	"simulator/service/src/modules/logs/application/dtos"
	"simulator/service/src/modules/logs/domain/entities"
	"simulator/service/src/modules/logs/domain/repositories"
)

// filterFromQuery maps the query DTO to the repository filter.
func (s *LogsService) filterFromQuery(q *dtos.LogQuery) repositories.LogFilter {
	return repositories.LogFilter{
		Limit:     q.Limit,
		Cursor:    q.Cursor,
		Protocol:  q.Protocol,
		Kind:      q.Kind,
		Direction: q.Direction,
		Device:    q.Device,
		Event:     q.Event,
		DateFrom:  q.DateFrom,
		DateTo:    q.DateTo,
		Q:         q.Q,
	}
}

// buildPage maps the stored entities to wire DTOs and attaches the next cursor.
func (s *LogsService) buildPage(list []entities.Log, next string) (*dtos.LogPage, error) {
	out := make([]dtos.Log, 0, len(list))
	for i := range list {
		dto, err := mapper.EntityToDto[entities.Log, dtos.Log](&list[i])
		if err != nil {
			return nil, err
		}
		out = append(out, *dto)
	}
	return &dtos.LogPage{Items: out, NextCursor: next}, nil
}

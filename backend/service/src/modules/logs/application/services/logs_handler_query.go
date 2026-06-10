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
		Offset:    q.Offset,
		Protocol:  q.Protocol,
		Kind:      q.Kind,
		Direction: q.Direction,
		Device:    q.Device,
		Q:         q.Q,
	}
}

// buildPage maps the stored entities to wire DTOs and wraps them with the total.
func (s *LogsService) buildPage(list []entities.Log, total int) (*dtos.LogPage, error) {
	out := make([]dtos.Log, 0, len(list))
	for i := range list {
		dto, err := mapper.EntityToDto[entities.Log, dtos.Log](&list[i])
		if err != nil {
			return nil, err
		}
		out = append(out, *dto)
	}
	return &dtos.LogPage{Items: out, Total: total}, nil
}

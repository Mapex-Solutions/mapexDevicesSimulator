package services

import (
	"github.com/Mapex-Solutions/mapexGoKit/utils/mapper"

	"simulator/service/src/modules/logs/application/dtos"
	"simulator/service/src/modules/logs/domain/entities"
)

// entityFromInput maps a write input to a log entity; id and created stay zero
// (the repository assigns them on insert).
func (s *LogsService) entityFromInput(in *dtos.LogInput) (*entities.Log, error) {
	return mapper.DtoToEntity[dtos.LogInput, entities.Log](in)
}

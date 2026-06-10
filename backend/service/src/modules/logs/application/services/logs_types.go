package services

import "simulator/service/src/modules/logs/application/di"

// LogsService serves the read-only device message history, mapping the query DTO
// to a repository filter and the stored entities back to wire DTOs.
type LogsService struct {
	deps di.LogsServiceDI
}

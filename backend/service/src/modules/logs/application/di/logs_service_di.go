package di

import (
	"go.uber.org/dig"

	"simulator/service/src/modules/logs/domain/repositories"
)

// LogsServiceDI declares the logs service dependencies as port interfaces.
type LogsServiceDI struct {
	dig.In

	Repo repositories.LogRepository
}

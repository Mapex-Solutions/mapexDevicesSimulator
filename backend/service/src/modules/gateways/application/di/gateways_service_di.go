package di

import (
	"go.uber.org/dig"

	"simulator/service/src/modules/gateways/domain/repositories"
)

// GatewaysServiceDI declares the gateways service dependencies as port interfaces.
type GatewaysServiceDI struct {
	dig.In

	Repo repositories.GatewayRepository
}

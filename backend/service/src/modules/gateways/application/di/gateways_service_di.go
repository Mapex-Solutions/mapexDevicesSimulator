package di

import (
	"go.uber.org/dig"

	"simulator/service/src/modules/gateways/domain/repositories"
	"simulator/service/src/shared/reconcile"
)

// GatewaysServiceDI declares the gateways service dependencies as port interfaces.
// Signal raises the shared CRUD-change notifier so the engine re-aligns its live
// sessions as soon as a gateway is written.
type GatewaysServiceDI struct {
	dig.In

	Repo   repositories.GatewayRepository
	Signal reconcile.Signal
}

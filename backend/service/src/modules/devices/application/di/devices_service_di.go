package di

import (
	"go.uber.org/dig"

	"simulator/service/src/modules/devices/domain/repositories"
	"simulator/service/src/shared/reconcile"
)

// DevicesServiceDI declares the devices service dependencies as port interfaces,
// resolved by the DIG container. Signal raises the shared CRUD-change notifier so
// the engine re-aligns its jobs and sessions as soon as a device is written.
type DevicesServiceDI struct {
	dig.In

	Repo   repositories.DeviceRepository
	Signal reconcile.Signal
}

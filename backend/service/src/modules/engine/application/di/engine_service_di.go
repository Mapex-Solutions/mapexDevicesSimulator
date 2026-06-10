package di

import (
	"go.uber.org/dig"

	consolePorts "simulator/service/src/modules/console/application/ports"
	devicesPorts "simulator/service/src/modules/devices/application/ports"
	enginePorts "simulator/service/src/modules/engine/application/ports"
	logsPorts "simulator/service/src/modules/logs/application/ports"
)

// EngineServiceDI declares the engine dependencies as port interfaces: it reads
// devices, writes logs, streams to the console, and dispatches over a registry.
type EngineServiceDI struct {
	dig.In

	Devices  devicesPorts.DevicesServicePort
	Logs     logsPorts.LogWriter
	Console  consolePorts.Publisher
	Registry enginePorts.Registry
}

package di

import (
	"go.uber.org/dig"

	consolePorts "simulator/service/src/modules/console/application/ports"
	devicesPorts "simulator/service/src/modules/devices/application/ports"
	enginePorts "simulator/service/src/modules/engine/application/ports"
	gatewaysPorts "simulator/service/src/modules/gateways/application/ports"
	logsPorts "simulator/service/src/modules/logs/application/ports"
)

// EngineServiceDI declares the engine dependencies as port interfaces: it reads
// devices and gateways, writes logs, streams to the console, dispatches one-shot
// sends over a registry, and opens live sessions over a connector registry.
type EngineServiceDI struct {
	dig.In

	Devices    devicesPorts.DevicesServicePort
	Gateways   gatewaysPorts.GatewaysServicePort
	Logs       logsPorts.LogWriter
	Console    consolePorts.Publisher
	Registry   enginePorts.Registry
	Connectors enginePorts.ConnectorRegistry
}

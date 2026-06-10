package configMod

import (
	consoleModule "simulator/service/src/modules/console"
	devicesModule "simulator/service/src/modules/devices"
	engineModule "simulator/service/src/modules/engine"
	gatewaysModule "simulator/service/src/modules/gateways"
	logsModule "simulator/service/src/modules/logs"
)

// ModuleDefinition describes a business module's init phases. Each phase resolves
// its own dependencies (the SQLite manager, the fiber app) from the DI container,
// so the registry stays decoupled from concrete wiring.
type ModuleDefinition struct {
	Name             string
	InitRepositories func()
	InitServices     func()
	InitInterfaces   func()
}

// Modules lists the business modules in init order; the app module's loop runs
// the phases (repositories, services, interfaces) across all modules.
var Modules = []ModuleDefinition{
	{
		Name:             "Devices",
		InitRepositories: devicesModule.InitRepositories,
		InitServices:     devicesModule.InitServices,
		InitInterfaces:   devicesModule.InitInterfaces,
	},
	{
		Name:             "Gateways",
		InitRepositories: gatewaysModule.InitRepositories,
		InitServices:     gatewaysModule.InitServices,
		InitInterfaces:   gatewaysModule.InitInterfaces,
	},
	{
		Name:             "Logs",
		InitRepositories: logsModule.InitRepositories,
		InitServices:     logsModule.InitServices,
		InitInterfaces:   logsModule.InitInterfaces,
	},
	{
		Name:           "Console",
		InitServices:   consoleModule.InitServices,
		InitInterfaces: consoleModule.InitInterfaces,
	},
	{
		Name:           "Engine",
		InitServices:   engineModule.InitServices,
		InitInterfaces: engineModule.InitInterfaces,
	},
}

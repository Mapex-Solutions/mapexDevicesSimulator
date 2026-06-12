package main

import (
	"flag"
	"time"

	container "github.com/Mapex-Solutions/mapexGoKit/microservices/container"
	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"
	"github.com/Mapex-Solutions/mapexGoKit/microservices/shutdown"

	"simulator/service/src/bootstrap"
	appModule "simulator/service/src/modules/app"
)

// main boots the simulator sidecar: configuration and logging, the SQLite store,
// the HTTP app (health + business modules + SPA), then blocks until a signal and
// shuts down gracefully. The Electron launcher passes the bind host and listen
// port via --addr and --port.
func main() {
	host := flag.String("addr", "", "host to bind on (empty = use http_address config)")
	port := flag.Int("port", 0, "HTTP port to listen on (0 = use http_port config)")
	flag.Parse()

	container.InitContainer()
	c := container.GetContainer()

	bootstrap.InitConfig()
	bootstrap.InitLogger()

	bootstrap.InitSQLite(c)

	app := bootstrap.InitFiber(c)
	bootstrap.InitHealth(app)

	// Business modules register their /api and /ws routes before the SPA catch-all.
	appModule.InitModule()

	// SPA static serving is the catch-all and must come last.
	bootstrap.InitStatic(app)

	sm := shutdown.New()
	bootstrap.InitShutdown(c, sm, app)

	addr := bootstrap.ListenAddress(*host, *port)
	go func() {
		if err := app.Listen(addr); err != nil {
			logger.Error(err, "[INFRA:HTTP] HTTP server stopped")
		}
	}()
	logger.Info("[INFRA:HTTP] simulatord listening on " + addr)

	sm.WaitForSignal(10 * time.Second)
}

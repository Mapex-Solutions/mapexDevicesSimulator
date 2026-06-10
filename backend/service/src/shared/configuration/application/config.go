package configApp

import (
	config "github.com/Mapex-Solutions/mapexGoKit/microservices/config"
)

// DefaultConfiguration declares every config key for the simulator sidecar, with
// env-var overrides and defaults. The Electron launcher passes most of these as
// env vars or a --port flag; standalone runs fall back to these values.
var DefaultConfiguration = []config.ConfigDefinition{
	/** Service identity */
	{Key: "service_name", Env: "SERVICE_NAME", Type: "string", Default: "mapex-devices-simulator"},
	{Key: "service_version", Env: "SERVICE_VERSION", Type: "string", Default: "0.1.0"},
	{Key: "go_env", Env: "GO_ENV", Type: "string", Default: "dev"},

	// log_level overrides the env-based default (debug in dev, info otherwise).
	{Key: "log_level", Env: "LOG_LEVEL", Type: "string", Default: ""},

	/** HTTP server */
	{Key: "http_port", Env: "HTTP_PORT", Type: "int", Default: 5055},
	{Key: "http_address", Env: "HTTP_ADDRESS", Type: "string", Default: "127.0.0.1"},
	{Key: "ctx_timeout", Env: "CTX_TIMEOUT", Type: "int", Default: 15},

	// sqlite_path is the on-disk database file. Relative paths resolve against the
	// working directory; the Electron sidecar passes an absolute per-user path.
	{Key: "sqlite_path", Env: "SQLITE_PATH", Type: "string", Default: "./data/simulator.db"},

	// spa_dir overrides the binary-embedded SPA: when set, static files are served
	// from this directory (dev/testing). Empty serves the embedded build.
	{Key: "spa_dir", Env: "SPA_DIR", Type: "string", Default: ""},

	// cors_origins is the dev SPA allowlist (Quasar Vite on :9100). Production is
	// same-origin (the sidecar serves the SPA), so CORS is irrelevant there.
	{Key: "cors_origins", Env: "CORS_ORIGINS", Type: "string", Default: "http://localhost:9100,http://127.0.0.1:9100"},
}

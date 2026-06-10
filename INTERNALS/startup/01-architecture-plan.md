# mapexDevicesSimulator — Architecture & Startup Plan

> Status: APPROVED — implementation starting with the frontend (UI).
> Scope: a local desktop tool that simulates IoT devices sending traffic to Mapex services.

## 1. Purpose

A local desktop tool to simulate IoT devices sending traffic against the Mapex services, covering four
protocols:

- **HTTP** (endpoint + api-key / basic auth)
- **MQTT** (broker + username/password OR client TLS certificate)
- **LoRaWAN** 1.0.x and 1.1 (OTAA join + data uplinks, real crypto)
- **Basics Station** (LoRaWAN frames carried over the Basics Station WS protocol; Semtech UDP as alternative)

Immediate driver: provide an end-to-end traffic generator, in particular to close the e2e gap of `mapexLNS`
(which today has no gateway/simulator). The tool is multi-protocol and general purpose.

## 2. Confirmed decisions

- **Location**: `/MAPEX/mapexDevicesSimulator/{frontend,backend}` (new repo, sibling of mapexLNS/mapexOS).
- **Frontend**: Quasar 2 + Vue 3 + Electron, following the UI architecture standard
  (reference: `mapexOS/workspace_js/apps/mapexOS`).
- **Backend**: Go (Fiber v2 + dig), following the Go architecture standard
  (reference: `mapexOS/workspace_go/services/router`).
- **Runtime**: the Electron main process spawns the Go binary as a **sidecar**; the Go server serves
  **SPA + `/api` (REST) + `/ws` (WebSocket)** on `http://127.0.0.1:PORT`; the `BrowserWindow` loads that URL.
- **MVP order** (vertical slice): **HTTP → MQTT → LoRaWAN/BasicsStation**.
- **LoRaWAN**: reuse The Things Stack `go.thethings.network/lorawan-stack/v3 v3.35.0`
  (`pkg/crypto`, `pkg/band`, `pkg/types`, `lbslns`) — same version as mapexLNS for exact e2e parity.
- **Internal project**: no architecture-standard / internal-doc references in code or comments; follow the
  house config pattern (dev defaults + `Sensitive:true`).

Confirmed reuse from `mapexGoKit` (via `replace ../../mapexGoKit/*`):
`infrastructure/{mqttclient,httpclient,common}`,
`microservices/{config,container(dig),logger,metrics,module,http,shutdown,validator,health}`.
**No Mongo/Redis/NATS** — the sidecar is in-memory (`map+sync.RWMutex` repositories, optionally JSON-persisted
under `store_dir`).

## 3. Repo layout

```
mapexDevicesSimulator/
  INTERNALS/                 # internal docs (this folder)
  frontend/                  # Quasar/Vue/Electron; src-electron spawns the Go sidecar
  backend/                   # Go: Fiber + dig; serves SPA (embed.FS) + /api + /ws; protocol drivers
```

The Quasar build produces `dist/spa`; that output is copied into `backend/web/dist` before `go build`
(the `embed.FS` packs the SPA into the binary). In **dev**, the Vite dev server runs separately and the SPA
talks to the sidecar via a base URL injected by the Electron preload (CORS allowed).

## 4. Frontend architecture (Quasar/Vue/Electron)

Base: `mapexOS/workspace_js/apps/mapexOS` (Quasar CLI + Vite + TS, Quasar 2.18, Vue 3.5, Pinia 3,
vue-router 4, vue-i18n 11). Electron is already configured in the base `quasar.config.ts` but there is **no
`src-electron/`** yet — new scaffold. Auth/login/permissions/org concepts are removed (local single-user tool).

### 4.1 `src-electron/` (new scaffold)

```
src-electron/
  electron-main.ts            # single-instance; pre-picks a free port (net.listen(0)) and passes --port;
                              #   spawns the Go binary; polls GET /api/health until 200; creates BrowserWindow
                              #   -> http://127.0.0.1:PORT (prod) or APP_URL (dev); SIGTERM->SIGKILL of the
                              #   sidecar on before-quit/window-all-closed; handles unexpected exit
  electron-preload.ts         # contextIsolation; exposes window.__SIM__ = { apiBase, wsBase, platform, appVersion }
  sidecar/sidecar-manager.ts  # dev path (src-electron/sidecar/bin/<plat>-<arch>/bin) vs prod
                              #   (process.resourcesPath/sidecar/bin)
  icons/
```

`quasar.config.ts`: `electron.bundler='packager'` + `packager.extraResource` to embed the Go binary;
`devServer.port=9100`; `boot=['i18n']` (no `permissions`); `framework.plugins=['Notify','Dialog']`.

### 4.2 `src/` (SPA)

```
src/
  boot/i18n.ts
  css/{app.scss (--mapex-* tokens, ported from base), quasar.variables.scss}
  layouts/main/MainLayout.vue          # drawer nav + header + theme toggle (no LoginLayout)
  pages/{dashboard,connections,devices,scenarios,runs(RunsPage+RunMonitorPage),errors}/   # each a folder + index.ts
  components/protocols/
    ProtocolRegistry/                  # typed registry: protocolId -> {labelKey,icon,configComponent,
                                       #   defaultConfig,validate,enabled}
    HttpConnectionConfig/              # Phase 1 (enabled)
    MqttConnectionConfig/              # Phase 2 (flip enabled)
    LoraWanConnectionConfig/           # Phase 3
    BasicsStationConnectionConfig/     # Phase 3
  composables/{i18n,shared/useLogger,simulator/{useSimulatorApi,useWebSocket},runs/useLiveStream}
  services/sim/{client.ts,index.ts,interfaces/*.interface.ts,
                resources/{connections,devices,scenarios,runs,health}.api.ts}
  stores/{connectionSettings,devices,scenarios,run,liveStream,app}/   # each a 5-file split
  router/{index,routes/*,types}
  i18n/{en-US,pt-BR}/                   # mirrored
```

### 4.3 Protocol-additive abstraction (extensibility core)

`ProtocolRegistry` lists the four protocols with `enabled` per phase; `ConnectionsPage` iterates the enabled
ones and renders `<component :is="def.configComponent">`. `ProtocolConfig` is a discriminated union on
`protocolId`. Adding MQTT/LoRaWAN = new component folder + flip `enabled` + a store key + i18n — **without**
touching pages/runs/services.

### 4.4 Live WebSocket

- `composables/simulator/useWebSocket.ts`: generic connector with exponential+jitter reconnect, state machine,
  teardown on `onScopeDispose`.
- `composables/runs/useLiveStream.ts`: builds `ws://127.0.0.1:PORT/ws/runs/:id`, routes `log`/`metric`/`runStatus`
  frames to stores.
- `stores/liveStream/`: **ring buffer** (`MAX_LOG_ENTRIES ~5000`, `droppedCount`), **batched flush** (~100ms/rAF)
  to avoid reactivity thrash, `QVirtualScroll` in the viewer. Kept separate from the `run` store (high cadence).

### 4.5 Typed API layer (simplified for a local tool)

A local lightweight typed client (do NOT vendor the platform `apis` package — no JWT/refresh/multi-service/org).
Keep the typed-wrapper pattern: `services/sim/client.ts` resolves base
`window.__SIM__?.apiBase ?? window.location.origin` + `/api`, one axios instance, no auth.
`resources/*.api.ts` expose typed functions; `index.ts` aggregates `export const sim = {...}`.
REST map: `/api/health`, `/api/connections`(+`/:id/test`), `/api/devices`, `/api/scenarios`(+`/:id/preview`),
`/api/runs`(GET/POST), `/api/runs/:id`(GET), `/api/runs/:id/stop`(POST); WS `/ws/runs/:id`.

## 5. Backend architecture (Go)

`go.mod` module: `github.com/Mapex-Solutions/mapexDevicesSimulator/backend`; requires: fiber v2,
`go.uber.org/dig`, `github.com/gofiber/contrib/websocket`, `go.thethings.network/lorawan-stack/v3 v3.35.0`,
a WS client (`nhooyr.io/websocket` or `gorilla/websocket`) for Basics Station; `replace ../../mapexGoKit/*`.

### 5.1 Module boundaries

Protocol drivers are **NOT** top-level modules: they are **infrastructure adapters behind one `Driver` port**
inside the `simulation` module. This keeps the runner protocol-agnostic (adding MQTT/LoRaWAN = new adapter +
factory registration, no runner change). Four modules:

| Module      | Context                                                              | Role |
|-------------|---------------------------------------------------------------------|------|
| `device`    | Simulated devices (identity + per-protocol credentials)             | entity + CRUD REST |
| `scenario`  | Traffic profile (rate, count, payload template, schedule, target)   | entity + CRUD REST + `/preview` |
| `simulation`| Run orchestration + ALL drivers + log/metric hub                    | `Driver` port, `RunService`, fan-out workers, `LogHub` |
| `streaming` | WS of logs/metrics to the UI + serving the SPA (embed.FS)           | interfaces/application only; subscribes to the hub |

### 5.2 `Driver` port

```go
type Driver interface {
    Connect(ctx, dev Device, target Target) (Session, error) // HTTP=no-op; MQTT opens session; LoRaWAN does OTAA join
    Send(ctx, sess Session, payload []byte) (SendResult, error)
    Disconnect(ctx, sess Session) error
    Kind() string
}
```

The `RunService` depends only on `DriverFactory` + the port. `Connect/Send/Disconnect` (instead of a single
`Run`) because scheduling/rate/worker-pool is the `Run` domain's responsibility, and the LoRaWAN OTAA join fits
in `Connect` (keeping `Send` uniform across all four protocols).

### 5.3 Live logs/metrics flow

Run workers call `Hub.Publish(runId, RunEvent)` and increment `*SimMetrics`. The `streaming` module is the only
subscriber: `log_ws_routes.go` registers `app.Get("/ws/runs/:id", websocket.New(...))`; the handler calls
`streamer.AttachClient` → `Hub.Subscribe` → a write-pump copies `RunEvent` (JSON) to the socket; a read-pump
detects close → unsubscribe. Buffered channels with **drop-oldest** so a slow UI never back-pressures the sim.

### 5.4 Static vs per-run config

- **Static (env / DefaultConfiguration)** = sidecar concerns only: `http_port`(5055), `http_address`(127.0.0.1),
  `service_name`, `service_version`, `go_env`, `log_level`, `ctx_timeout`, `cors_origins`,
  `sidecar_secret`(Sensitive), `max_concurrent_runs`, `max_workers_per_run`, `hub_buffer_size`,
  `ws_write_timeout`, `store_dir` (empty = pure in-mem), `metrics_go_collector`, `metrics_process_collector`.
- **Per-run (REST request, never env)** = target+auth (`target.kind/url`, `auth.apiKey|basic|mqttUser/Pass`,
  `tls.clientCertPEM/keyPEM/caPEM`), LoRaWAN material on the `Device`
  (devEUI/joinEUI/appKey/nwkKey/macVersion/band/devAddr), and `Scenario` fields
  (rate/count/duration/payloadTemplate/schedule). Request secrets are never logged.

## 6. Sidecar integration contract

- The sidecar accepts `--addr 127.0.0.1 --port N` (Electron pre-picks a free port and passes it); it prints
  `LISTENING http://127.0.0.1:N` as the first stdout line (confirmation).
- `GET /api/health` → `200 {"status":"ok","version":...}` (Electron polls with backoff before opening the window).
- Prod: Go serves `web/dist` at `/` (SPA fallback to index.html). Dev: Vite serves the SPA; CORS for
  `http://localhost:9100`; sidecar base provided via the preload bridge.

## 7. Execution order (vertical slices)

### Phase 0 — Skeleton (no protocol)

- Backend: `go.mod`, `web/embed.go`, `src/main.go`, all `bootstrap/*`, `shared/configuration/*`,
  `app/module.go`, full `device` and `scenario` modules, `simulation` skeleton (Run entity, in-mem RunRepo,
  ports, `channel_hub`, `RunService` with fan-out), `streaming` module (WS + static SPA).
- Frontend: `src-electron/*`, `quasar.config.ts` edits, `css/*`, `boot/i18n`, i18n skeleton, `MainLayout`,
  `router/*`, `services/sim` + `health.api`, `stores/app` (health polling), i18n/logger composables.
- **Proof**: Electron spawns the sidecar → waits for `/api/health` → opens the window → SPA shows "sidecar healthy".

### Phase 1 — HTTP (full vertical slice)

- Backend: `drivers/http/http_driver.go` (GoKit httpclient, api-key/basic) + factory case + DI in
  `simulation/module.go`.
- Frontend: `ProtocolRegistry` (http only), `HttpConnectionConfig`, stores
  `connectionSettings/devices/scenarios/run/liveStream`, Connections/Devices/Scenarios/Runs/RunMonitor pages,
  `useWebSocket`+`useLiveStream`, `LiveLogViewer`+`LiveMetrics`.
- **e2e proof**: configure HTTP → create device → create scenario → start run (`POST /api/runs`) → live
  logs/metrics via `/ws` → stop.

### Phase 2 — MQTT (additive)

`drivers/mqtt/{mqtt_driver,tls,types}` (GoKit mqttclient, user/pass + cert) + factory case; frontend
`MqttConnectionConfig` + flip `enabled` + union/validators/i18n. Runner/streaming/stores **unchanged**.

### Phase 3 — LoRaWAN + Basics Station (additive)

`simulation/domain/services/frame/*` (1.0.x and 1.1, TTS crypto) + tests;
`drivers/lorawan/{basicstation_driver,lns_ws_client,codec,udp_driver,types}` (lbslns) + factory cases; frontend
`LoraWanConnectionConfig` + `BasicsStationConnectionConfig` + flips + i18n. The Device entity already carries
LoRaWAN material from Phase 0.

## 8. Verification

- **Backend standalone**: `cd backend && go build ./... && go vet ./... && go test ./...` (green).
  `go run ./src`, check `GET /api/health` and that the embedded SPA responds at `/`.
- **Frontend standalone**: `cd frontend && quasar dev` (browser) proves UI ↔ `/api` ↔ `/ws` against a separately
  running sidecar.
- **Desktop e2e (Phase 1)**: build the SPA → copy to `backend/web/dist` → `go build` →
  `quasar build -m electron` with the binary in `extraResource` → open the packaged app → full HTTP flow
  (device → scenario → run → live view → stop).
- **LoRaWAN (Phase 3)**: run `mapexLNS` (dev stack already brings up redis/nats/minio; the LNS listens on
  `:1700` UDP + `:1887` Basics Station) and point the simulator at it — the first real e2e of the LNS
  (OTAA join 1.0.x/1.1 → uplink → NATS).
- **Unit**: LoRaWAN frame builder against TTS vectors; drivers with a mock transport; runner with mock
  `Driver`/`Hub`.

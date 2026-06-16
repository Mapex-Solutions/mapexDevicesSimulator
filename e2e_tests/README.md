# Simulator E2E Tests

> 🇧🇷 Versão em português: [README_pt.md](./README_pt.md)

End-to-end tests for the Mapex Devices Simulator, following the same architecture
as the [MapexOS e2e suite](../../../mapexOS/e2e_tests): a small saga runner, shared
fixtures per service, and journeys that compose them.

Because the simulator is a single local service with no auth, the framework is
the platform pattern minus the multi-service `ClientSet` and the IAM bootstrap:
one `Sim` client against the sidecar, and readiness is just a green
`GET /api/health`.

## Layout

```
e2e_tests/
├── core/saga/          # the saga runner: Item (Step|Assert), Context (+ bag), Run
├── common/
│   ├── constants/      # SimURL + ChirpStack endpoints
│   ├── types/          # the { status, errors, data } envelope
│   └── utils/          # health probe + fixtures (echo, MQTT broker, certs, ChirpStack stack)
├── services/           # services/{svc}/{mod}/ — steps / asserts / payloads
│   ├── simulator/      # the sidecar: devices, gateways, engine, logs, targets
│   └── chirpstack/     # the LNS: client, stack lifecycle, provisioning
└── journey/iot/        # saga journeys (build tag: saga)
    ├── http_device_fire/
    ├── mqtt_device_fire/
    ├── mqtt_downlink/
    ├── lorawan_join_uplink/
    ├── lorawan_downlink/
    ├── fire_error/
    ├── console_stream/
    └── connection_status/
```

- **Step** mutates the simulator (POST/DELETE) and may register a `Compensate`
  for cleanup; **Assert** reads the public API only (never the SQLite file).
- The runner walks the journey in order and runs every `Compensate` in reverse
  on completion — pass or fail — so the next run starts clean.

## Prerequisites

The simulator sidecar running on `127.0.0.1:5055` (the service's default; the
endpoint is fixed in `common/constants`, no env vars to set):

```bash
cd ../backend
go run ./service/src --addr 127.0.0.1 --port 5055
```

A journey that finds the sidecar down calls `t.Skip`, so the suite never fails
just because nothing is running.

## How to run

```bash
cd e2e_tests

# All journeys (the saga build tag is required)
go test -tags=saga ./journey/...

# One journey, verbose
go test -tags=saga ./journey/iot/http_device_fire/ -v
```

The endpoints live in `common/constants` — change them there, not via the
environment.

## Journeys

| Journey | What it proves | Needs |
|---------|----------------|-------|
| [http_device_fire](./journey/iot/http_device_fire/) | Create an HTTP device, fire it, and confirm the simulator logs the 200 uplink with its captured response | sidecar only (echo is in-test) |
| [mqtt_device_fire](./journey/iot/mqtt_device_fire/) | Create MQTT devices, fire them, and confirm an embedded broker accepts the authenticated publish over both auth modes (username/password and client certificate) | sidecar only (broker + certs are in-test) |
| [mqtt_downlink](./journey/iot/mqtt_downlink/) | Receive-enabled MQTT device subscribes, an external (retained) publish arrives, and the simulator logs it as a downlink — the inbound path | sidecar only (broker is in-test) |
| [lorawan_join_uplink](./journey/iot/lorawan_join_uplink/) | Provision ChirpStack and drive LoRaWAN OTAA devices over Semtech UDP and Basics Station, confirming the LNS records the join and the uplink for each | sidecar + docker (the journey owns the ChirpStack stack) |
| [lorawan_downlink](./journey/iot/lorawan_downlink/) | Queue a downlink on ChirpStack, fire an uplink to open the RX window, and confirm the simulator logs the received downlink — the inbound path | sidecar + docker (the journey owns the ChirpStack stack) |
| [fire_error](./journey/iot/fire_error/) | Fire a device at an unreachable target and confirm the engine logs the failure as an error frame, not silently | sidecar only |
| [console_stream](./journey/iot/console_stream/) | Connect to the `/ws` console, fire a device, and confirm the uplink arrives live on the stream | sidecar only |
| [connection_status](./journey/iot/connection_status/) | Enable a device pointed at an unreachable broker and confirm `connecting`/`reconnecting` status frames arrive live on `/ws` (these are never logged) | sidecar only |

## Module contract tests (`go test ./...`)

Lighter HTTP-contract tests of a single module's surface, under
`services/{svc}/{mod}/e2e/`. They are not saga-gated — `go test ./...` runs them
and skips cleanly if the sidecar is down.

| Module | Covers |
|--------|--------|
| `services/simulator/devices/e2e` | CRUD (+ validation, 404) and list membership |
| `services/simulator/gateways/e2e` | CRUD (+ validation, 404) and list membership |
| `services/simulator/logs/e2e` | cursor pagination + filters (device, event, date range, combined) |

# Bounded Context: Engine

**Service:** simulator (mapexDevicesSimulator backend / `simulatord`)
**Module path:** `src/modules/engine/`
**Owner:** MAPEX
**Last reviewed:** 2026-06-09

## Purpose
The simulation engine. It runs the enabled devices: for every device that is on,
each event with an enabled schedule fires on its interval, rendering the
`{{placeholder}}` payload and dispatching it over the device's protocol. Each
fire streams a frame to the console and persists a log when the device's
storeLogs is on.

## Ubiquitous Language
| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| Job | One (device, event) that fires on a schedule | A device or event |
| Scheduler | The single goroutine + min-heap that fires due jobs | The OS scheduler |
| Reconcile | Re-derive the job set from the DB and align the heap | A REST call |
| Dispatcher | A protocol sender (http, mqtt, ...) selected by protocolId | A device connection |

## Published Events (driven — outbound)
- Console frames (via the console `Publisher` port) on every fire.
- HTTP requests / MQTT publishes (via the dispatcher registry) to the device's
  configured target.

## Consumed Events (driving — inbound)
- `OnMount` (module init) and `Reconcile` (CRUD change signal + slow resync).

## Driving Ports (what can call this module)
- `ports.EnginePort` (`OnMount` / `Reconcile` / `OnShutdown`), in-process.

## Driven Ports (what this module requires)
- `devices.DevicesServicePort` (read the devices + their events).
- `logs.LogWriter` (persist a message when storeLogs is on).
- `console.Publisher` (stream the live frame).
- `engine.Registry` (resolve a dispatcher by protocol).

## Invariants and Business Rules
- Goroutines are bounded: one scheduler plus a fixed worker pool, regardless of
  the number of jobs.
- The DB is the source of truth; the in-RAM job set is derived on reconcile.
- Only `device.enabled` devices with `schedule.enabled` events fire; disabling
  either removes the job on the next reconcile.
- A full worker queue drops a firing rather than stalling the scheduler.
- Secrets in a device config are sent but never logged.

## Known Cross-Context Interactions
- Reads devices (devices module), writes logs (logs module), streams to the
  console (console module). LoRaWAN/Basics Station dispatch ships in a later phase.

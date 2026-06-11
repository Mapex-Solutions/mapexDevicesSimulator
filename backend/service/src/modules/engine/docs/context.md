# Bounded Context: Engine

**Service:** simulator (mapexDevicesSimulator backend / `simulatord`)
**Module path:** `src/modules/engine/`
**Owner:** MAPEX
**Last reviewed:** 2026-06-11

## Purpose
The simulation engine. It runs the enabled devices in two cooperating halves: a
time-driven scheduler that fires each enabled-schedule event on its interval, and a
connection-driven session manager that holds a live connection per session-capable
device so the device can also RECEIVE (MQTT downlinks, LoRaWAN downlinks). Every
fire renders the `{{placeholder}}` payload and sends it over the device's protocol;
every send and every received message streams a frame to the console and persists a
log when the device's storeLogs is on. A device can also be fired on demand.

## Ubiquitous Language
| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| Job | One (device, event) that fires on a schedule | A device or event |
| Scheduler | The single goroutine + min-heap that fires due jobs | The OS scheduler |
| Reconcile | Re-derive the job/session set from the DB and align it | A REST call |
| Dispatcher | A one-shot protocol sender (http) selected by protocolId | A live session |
| Session | A persistent per-device connection (mqtt, lorawan) | A one-shot dispatch |
| Connector | Opens a Session for a session-capable protocol | A Dispatcher |

## Published Events (driven — outbound)
- Console frames (via the console `Publisher` port) on every uplink, downlink and
  connection-status transition.
- HTTP requests (one-shot dispatcher), MQTT publishes/subscribes and LoRaWAN
  uplinks (live sessions over a shared gateway link) to the device's target.

## Consumed Events (driving — inbound)
- `OnMount` (module init), `Reconcile` (CRUD change signal + slow resync), and
  `Fire` (an on-demand send for one device).
- Downlinks received on a live session, surfaced in-process to the console + logs.

## Driving Ports (what can call this module)
- `ports.EnginePort` (`OnMount` / `Reconcile` / `Fire` / `OnShutdown`), in-process,
  plus the `POST /api/devices/:id/fire` route that delegates to `Fire`.

## Driven Ports (what this module requires)
- `devices.DevicesServicePort` (read the devices + their events).
- `gateways.GatewaysServicePort` (resolve a LoRaWAN device's gateway link).
- `logs.LogWriter` (persist a message when storeLogs is on).
- `console.Publisher` (stream the live frame).
- `engine.Registry` (resolve a one-shot dispatcher by protocol).
- `engine.ConnectorRegistry` (resolve a session connector by protocol).

## Invariants and Business Rules
- Goroutines are bounded: one scheduler, a fixed worker pool, and one supervisor
  goroutine per live session, regardless of the number of jobs.
- The DB is the source of truth; the in-RAM job set AND session set are derived on
  reconcile.
- `enabled` means live: an enabled session-capable device holds a connection;
  disabling tears it down. An offline broker/LNS reconnects forever with bounded
  backoff, every attempt streamed to the console.
- Only `device.enabled` devices with `schedule.enabled` events fire; disabling
  either removes the job on the next reconcile.
- A full worker queue drops a firing rather than stalling the scheduler.
- Secrets in a device config are sent but never logged.

## Known Cross-Context Interactions
- Reads devices (devices module) and gateways (gateways module), writes logs (logs
  module), streams to the console (console module). HTTP dispatches one-shot; MQTT
  and LoRaWAN (1.0.x/1.1, over Semtech UDP or Basics Station) run over live
  sessions.

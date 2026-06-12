# 08 — Engine: scheduler + live sessions

The engine is what turns stored devices into running traffic. It has two halves that
work side by side: a **scheduler** that fires timed events, and a **session manager**
that holds the persistent connections devices need to receive.

## Capabilities

- **Scheduler** — a min-heap of timed jobs (one per enabled device + scheduled event)
  fired at their due time by a bounded worker pool.
- **Session manager** — opens and holds a live connection per enabled session-capable
  device (MQTT, LoRaWAN, Basics Station); HTTP stays one-shot.
- **Enabled = live** — enabling a device brings its session up; disabling tears it
  down.
- **Infinite reconnect** — a dropped broker/LNS reconnects forever with a bounded
  backoff (floor ~1s, ceiling ~30s, jitter); every transition is reported.
- **Immediate reconcile** — device and gateway writes re-align the jobs and sessions
  at once, not on the slow safety resync.
- **Gateway-flag aware** — a disabled gateway takes its devices fully offline across
  the session manager, the scheduler and the fire path.

## How it works

The `engine` module (`backend/service/src/modules/engine/`) exposes
`OnMount / Reconcile / Fire / OnShutdown`. On mount it reads the devices, builds the
job heap and the desired session set, and starts the workers, scheduler and a slow
safety resync.

- **Scheduler** — `buildDesired` derives jobs; `process()` routes a fire through the
  device's live session when one exists, else the one-shot dispatcher, else (a
  LoRaWAN/Basics Station device with no link) a `gateway-offline` status frame.
- **Session manager** — `buildDesiredSessions` derives the set; `superviseSession`
  connects, holds open and reconnects with backoff; connectors (`mqtt`, `lorawan`)
  live in `infrastructure/session`.
- **Immediate reconcile** — a neutral `shared/reconcile` signal: the devices and
  gateways services raise it on every write, and the engine binds its `Reconcile` to
  it at mount (no module imports the other).

Commands arrive over **REST**; every result (uplink echo, downlink, status) streams
over the **WebSocket** [console](./09-console.md).

## Notes

- The two halves are independent: the scheduler decides *when* to send, the session
  manager decides *over what link* — a device can have one, both or neither.
- A slow background resync remains as a safety net behind the immediate reconcile.

---
> Part of the [MapexOS ecosystem](../README.md#part-of-the-mapexos-ecosystem).

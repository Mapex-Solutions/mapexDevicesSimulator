# 01 — Device management

A simulated device is the unit of work: a named entity that speaks one protocol and
carries the config, attributes and events used to generate traffic. The platform
manages them as first-class records, persisted across restarts.

## Capabilities

- **CRUD** — create, edit, delete and list devices; each is persisted in SQLite and
  survives a restart.
- **Four protocols** — every device is `http`, `mqtt`, `lorawan` or `basicstation`;
  the protocol is the discriminant of its connection config.
- **Enable / disable** — the `enabled` flag is the on/off switch. An enabled
  session-capable device is held live by the engine; disabling it disconnects.
- **Per-device config** — a protocol-specific connection object (endpoint, broker,
  or LoRaWAN keys) stored as opaque JSON, so new protocols need no schema change.
- **Attributes** — a free-form `key → value` map carried with the device.
- **Events** — a list of pre-registered, templated payloads with optional schedules
  (see [07 — Events](./07-events-scheduling.md)).
- **Per-device log capture** — a `storeLogs` flag decides whether the device's
  traffic is persisted to the [logs](./10-logs.md) history.

## How it works

Go DDD / hexagonal module at `backend/service/src/modules/devices/`:

- **interfaces/http** — `GET·POST /api/devices`, `PUT·DELETE /api/devices/:id`.
- **application** — `DevicesService` (port `DevicesServicePort`) over a
  `DeviceRepository` port; the DTO is `Device`/`DeviceInput` from
  `packages/contracts/devices`.
- **infrastructure/persistence/sqlite** — the `devices` table (`id, name, device_id,
  protocol_id, enabled, store_logs, config, attributes, events, created`).

The wire shape (`DeviceInput`): `name`, `deviceId`, `protocolId`
(`http|mqtt|lorawan|basicstation`), `enabled`, `storeLogs`, `config` (JSON),
`attributes` (map), `events` (JSON). The frontend mirrors it as a Zod schema
(`packages/schema`) consumed by a Pinia store (`useDevicesStore`) and the
create/edit stepper page.

Every successful write raises the shared **reconcile signal** so the
[engine](./08-engine.md) re-aligns its jobs and live sessions immediately, instead
of waiting for the slow safety resync.

## Notes

- `config`/`events`/`attributes` are stored as JSON, not normalized columns — a new
  protocol or event field is purely additive, no migration.
- The store reflects the backend exactly: no seed/sample data, so an empty backend
  shows an empty list; connectivity is surfaced by the sidecar health indicator.
- For LoRaWAN/Basics Station, the device is only live when its gateway is enabled too
  (see [02 — Gateways](./02-gateways.md)).

---
> Part of the [MapexOS ecosystem](../README.md#part-of-the-mapexos-ecosystem).

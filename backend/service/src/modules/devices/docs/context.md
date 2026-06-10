# Bounded Context: Devices

**Service:** simulator (mapexDevicesSimulator backend / `simulatord`)
**Module path:** `src/modules/devices/`
**Owner:** MAPEX
**Last reviewed:** 2026-06-09

## Purpose
Owns the lifecycle of a simulated device: its identity, its per-protocol target
config, its free-form attributes, and its pre-registered events. It is plain CRUD
over SQLite; the simulation engine (a later module) reads devices to drive
traffic but does not own them.

## Ubiquitous Language
| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| Device | A simulated endpoint with its own protocol config | A real LoRaWAN/MQTT device |
| deviceId | The device's user-facing identifier (drives `{{deviceId}}`) | `id` (the server-assigned uuid) |
| config | Per-protocol target + credentials, discriminated by protocolId | A saved Connection |
| event | A pre-registered request the engine can fire | A live run frame |

## Published Events (driven — outbound)
None.

## Consumed Events (driving — inbound)
None.

## Driving Ports (what can call this module)
- HTTP routes under `/api/devices` (list, create, update, delete).

## Driven Ports (what this module requires)
- `repositories.DeviceRepository` (SQLite, via the shared sqlite model).

## Invariants and Business Rules
- `id` and `created` are server-assigned; never accepted from the client.
- `protocolId` is one of http | mqtt | lorawan | basicstation.
- `config`, `attributes` and `events` are stored as JSON; the queryable columns
  are the scalar fields (`device_id`, `protocol_id`, `enabled`).
- Secrets inside `config` are persisted but never logged.

## Known Cross-Context Interactions
- Read by the simulation engine (runs module) to drive traffic, by device id.

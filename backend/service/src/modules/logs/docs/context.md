# Bounded Context: Logs

**Service:** simulator (mapexDevicesSimulator backend / `simulatord`)
**Module path:** `src/modules/logs/`
**Owner:** MAPEX
**Last reviewed:** 2026-06-09

## Purpose
Owns the persisted device message history: the SQLite-backed archive behind the
live console stream. Every uplink, downlink and auth/join handshake the engine
emits (for devices with storeLogs on) is persisted here and served back, filtered
and paginated, to the logs page.

## Ubiquitous Language
| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| Log | One persisted device message | A live ConsoleMessage on the WS stream |
| created | The message time (server-assigned) | A device's created |
| q | Free-text match over summary, payload, device name | An exact filter |

## Published Events (driven — outbound)
None.

## Consumed Events (driving — inbound)
None over the HTTP API. The simulation engine calls `Insert` in-process to persist
messages (no HTTP write route exists).

## Driving Ports (what can call this module)
- HTTP route `GET /api/logs` (filter + pagination).

## Driven Ports (what this module requires)
- `repositories.LogRepository` (SQLite). The read uses raw SQL via the model pool
  because the free-text `q` search needs LIKE/OR beyond equality filters.

## Invariants and Business Rules
- Read-only over HTTP; logs are written only by the engine via `Insert`.
- Results are newest-first (ORDER BY created DESC); page size is capped.
- Equality filters (protocol, kind, direction, device) plus free-text `q`.

## Known Cross-Context Interactions
- Written by the simulation engine; the same messages also stream live to the
  console over the WebSocket.

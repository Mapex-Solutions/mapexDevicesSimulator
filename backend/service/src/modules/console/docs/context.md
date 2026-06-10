# Bounded Context: Console

**Service:** simulator (mapexDevicesSimulator backend / `simulatord`)
**Module path:** `src/modules/console/`
**Owner:** MAPEX
**Last reviewed:** 2026-06-09

## Purpose
Owns the live console stream: the WebSocket at `/ws` that pushes `ConsoleMessage`
frames to the UI as the simulation engine emits them (every uplink, downlink and
auth/join handshake, across all protocols). It is a transport/broker only; it
holds no state and persists nothing (the logs module persists, when storeLogs is
on).

## Ubiquitous Language
| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| ConsoleMessage | One live frame on the stream (`ts` = frame time) | A persisted Log (`created`) |
| Hub | In-memory fan-out broker over connected clients | A message bus |
| Client | One connected WebSocket subscriber | A simulated device |

## Published Events (driven — outbound)
- `ConsoleMessage` frames over the `/ws` WebSocket (raw JSON, not the REST
  envelope), to every connected client.

## Consumed Events (driving — inbound)
None over the network. The engine calls the `Publisher` port in-process.

## Driving Ports (what can call this module)
- WebSocket `GET /ws` (clients subscribe).
- `ports.Publisher` (in-process), called by the simulation engine to broadcast.

## Driven Ports (what this module requires)
- None. The hub is self-contained in memory.

## Invariants and Business Rules
- A slow client drops frames (bounded buffer) rather than stalling the broadcaster.
- Frames are raw JSON; the `{status,errors,data}` envelope is for REST only.
- The console always reflects live activity; persistence is the logs module's job.

## Known Cross-Context Interactions
- The simulation engine (engine module) holds the `Publisher` port and streams
  every fired message here, in parallel with persisting a Log when storeLogs is on.

# Bounded Context: Gateways

**Service:** simulator (mapexDevicesSimulator backend / `simulatord`)
**Module path:** `src/modules/gateways/`
**Owner:** MAPEX
**Last reviewed:** 2026-06-09

## Purpose
Owns the lifecycle of a simulated LoRaWAN gateway: its identity (EUI), region,
and the link to the LNS (Basics Station WSS or Semtech UDP). Plain CRUD over
SQLite; LoRaWAN devices reference a gateway by id to forward their frames.

## Ubiquitous Language
| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| Gateway | Radio link to the LNS that forwards device frames | The LNS itself |
| EUI | The gateway's 16-hex identifier | A device's deviceId |
| link | Basics Station (wss) or UDP (host/port) connection block | A device's protocol config |

## Published Events (driven — outbound)
None.

## Consumed Events (driving — inbound)
None.

## Driving Ports (what can call this module)
- HTTP routes under `/api/gateways` (list, create, update, delete).

## Driven Ports (what this module requires)
- `repositories.GatewayRepository` (SQLite, via the shared sqlite model).

## Invariants and Business Rules
- `id` and `created` are server-assigned; never accepted from the client.
- `region` is one of the supported LoRaWAN regions.
- `link` is stored as JSON; the queryable columns are the scalar fields.

## Known Cross-Context Interactions
- Referenced by lorawan devices (devices module) via the device config gatewayId.

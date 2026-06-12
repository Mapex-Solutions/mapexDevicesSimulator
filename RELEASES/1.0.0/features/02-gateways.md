# 02 — Gateway management

A gateway is the bridge a LoRaWAN device rides on: it holds the single link to the
LoRaWAN Network Server and forwards the frames of every sensor attached to it. The
platform manages gateways as their own records, separate from devices.

## Capabilities

- **CRUD** — create, edit, delete and list gateways, persisted in SQLite.
- **Two LNS transports** — a gateway connects to the LNS over **Semtech UDP**
  (packet-forwarder, `host:port`) or **Basics Station** (WebSocket LNS URI).
- **EUI + region** — each gateway has an 8-byte EUI and a regional plan
  (EU868, US915, AU915, AS923, CN470, IN865, KR920, RU864).
- **Enabled flag gates its devices** — disabling a gateway takes every LoRaWAN /
  Basics Station device riding it offline (no session, no scheduled fire, no manual
  fire), reflected immediately.
- **Connection status in the UI** — the gateway list derives online / connecting /
  offline per gateway from the live console stream.

## How it works

Go DDD module at `backend/service/src/modules/gateways/` — `GET·POST /api/gateways`,
`PUT·DELETE /api/gateways/:id`, a `GatewaysService` over a `GatewayRepository`, and a
`gateways` SQLite table. The DTO (`GatewayInput`): `name`, `eui`, `enabled`,
`region`, `description`, `link` (JSON `{protocol, lnsUri, host, port}`).

In the engine, a gateway becomes a **shared link** (`infrastructure/session`):
ref-counted and keyed by endpoint, so multiple LoRaWAN devices on the same gateway
reuse one socket. Downlinks arriving on the link are routed to the right device by
DevAddr. The gateway's `enabled` flag is consulted by the session manager, the
scheduler and the fire path, so it gates the devices everywhere.

The frontend mirrors the schema (Zod) into a Pinia store and a create/edit stepper
(identity + LNS link); a `useGatewayConnections` composable derives the live status.

## Notes

- A **Basics Station** device carries its own link (its embedded gateway EUI) rather
  than referencing a gateway record — but it is still gated by a gateway whose EUI
  matches, so disabling that gateway disconnects it too.
- A gateway sends a periodic `stat` so the LNS keeps it marked **online**; without it
  the LNS reports "never seen".
- Gateway writes raise the reconcile signal, so enabling/disabling takes effect at
  once (see [08 — Engine](./08-engine.md)).

---
> Part of the [MapexOS ecosystem](../README.md#part-of-the-mapexos-ecosystem).

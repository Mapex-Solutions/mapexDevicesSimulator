# 06 — LoRaWAN over Basics Station

The same real LoRaWAN sensor as [05](./05-lorawan-udp.md) — same identity, same crypto,
same join — but instead of riding a separate gateway it carries its **own Basics
Station WebSocket link** to the LNS. It is its own embedded gateway.

## Capabilities

- **Self-contained link** — the device holds a Basics Station WebSocket to the LNS;
  no separate gateway record is required.
- **Full LoRaWAN identity and crypto** — OTAA/ABP, MAC 1.0.x and 1.1, regions, class
  A/C, real MIC/encryption (shared with [05](./05-lorawan-udp.md)).
- **Embedded gateway EUI** — the device announces its own gateway EUI to the LNS.
- **Uplinks and downlinks** over the WebSocket, surfaced as `up` / `down` frames.

## How it works

The Basics Station transport lives in `backend/packages/utils/lorawan/basicstation`
(JSON framing) and the engine's LoRaWAN connector. The session dials
`ws://host:port/gw/<gateway-eui>`, runs the version handshake, then sends uplinks as
Basics Station messages (`jreq` for the join request, `updf` for data) and parses
`dnmsg` downlinks back into PHYPayloads. The device-side crypto and codec are exactly
those of [05](./05-lorawan-udp.md). Verified live against ChirpStack's Basics Station
gateway-bridge.

## Notes

- The connector **appends `/gw/<eui>`** to the LNS URI when it has no path
  (ChirpStack's convention); a full path is used as-is for other LNSs.
- **EUIs are id6 strings** (`"0011:2233:4455:66aa"`), not numbers, and the uplink data
  rate is the DR **index**, not the spreading factor.
- The DevNonce persists across reconnects (a reconnect resumes a joined device rather
  than re-joining), so the LNS does not reject it.
- If a gateway record with the same EUI exists and is disabled, the Basics Station
  device is taken offline with it (see [02 — Gateways](./02-gateways.md)).

---
> Part of the [MapexOS ecosystem](../README.md#part-of-the-mapexos-ecosystem).

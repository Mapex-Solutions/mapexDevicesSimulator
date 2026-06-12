# 09 — Console (live WebSocket stream)

The console is the live window into the simulation: every uplink, downlink and
connection event streams to it in real time, so you can watch a device connect, join,
publish and receive as it happens.

## Capabilities

- **Live stream** of every frame over a WebSocket — `up`, `down` and `system`.
- **Full status taxonomy** — connection lifecycle (`connecting`, `connected`,
  `subscribing`, `subscribed`, `reconnecting`, `disconnected`), LoRaWAN join
  (`join-request`, `join-accept`, `joined`, `activated`), and `gateway-offline` when a
  device's gateway is down.
- **Rich frames** — each carries a protocol, direction, kind, summary, status, payload
  and a `meta` map (e.g. LoRaWAN radio fields).
- **Filtering** — by device and by protocol.
- **Detail panel** — inspect a frame's full payload and metadata, copy it, and read an
  info hint explaining the DevAddr on join frames.

## How it works

The `console` module hosts a broadcast **hub** and the `/ws` endpoint. The socket is
**server → client only**: the read pump discards inbound frames, so commands never
travel over it (they go over REST). The engine publishes frames through the console
`Publisher` port — uplinks/downlinks via `report`/`emitInbound`, lifecycle via
`emitStatus`.

A frame (`ConsoleMessage`) is `{ ts, protocol, deviceId, deviceName, direction, kind,
summary, status, payload, meta }`. The frontend opens the stream from a Pinia store
and renders a device list, filters, the frame list and a detail panel.

## Notes

- The console is **live and ephemeral** — an in-memory ring buffer of the most recent
  frames, with no backfill on (re)connect. Durable, queryable history lives in the
  [logs](./10-logs.md).
- The WebSocket is read-only by design (the command-plane / event-plane split).

---
> Part of the [MapexOS ecosystem](../README.md#part-of-the-mapexos-ecosystem).

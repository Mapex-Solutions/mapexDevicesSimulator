# 11 — Fire event (on-demand)

Besides scheduled traffic, any event can be sent **right now**, on demand. Fire is how
you drive a single, deliberate uplink — to test a payload, trigger a downlink, or demo
a device — and watch the result stream back live.

## Capabilities

- **Fire any pre-registered event** for a device immediately.
- **Fire an ad-hoc event** — an inline payload that is not stored on the device.
- **All protocols** — HTTP, MQTT and LoRaWAN/Basics Station.
- **Live result** — the uplink echo, and any resulting downlink, stream back over the
  [console](./09-console.md) WebSocket.

## How it works

`Fire(deviceID, { eventId | event })` resolves the device and the event (a registered
one by id, or the inline ad-hoc one), builds the send spec, and runs it through the
**same `process()` path as a scheduled fire** — so the uplink goes through the live
session when one exists, and the result is reported identically. Exposed as
`POST /api/devices/:id/fire`, with the handler mapping not-found / unsupported to the
right HTTP status.

The frontend's Fire dialog builds the event from the form, renders a live preview of
what will be sent, and calls the API; the echo arrives back on the console.

## Notes

- Fire reuses the scheduler's dispatch path, so a manual fire behaves exactly like a
  scheduled one.
- Firing a LoRaWAN/Basics Station device whose gateway is offline reports a
  `gateway-offline` status frame instead of sending — the attempt is visible, nothing
  reaches the LNS.

---
> Part of the [MapexOS ecosystem](../README.md#part-of-the-mapexos-ecosystem).

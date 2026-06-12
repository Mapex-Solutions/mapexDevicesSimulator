# 07 — Events, templating and scheduling

An event is a reusable, pre-registered payload a device can send. Events make traffic
realistic and repeatable: they carry templates that are rendered fresh on every fire,
and an optional schedule that makes the device send on its own.

## Capabilities

- **Pre-registered events** per device — a named list, each holding the payload for
  the device's protocol (`http`, `mqtt` or `lorawan`).
- **Templating**, rendered at send time:
  - `{{deviceId}}` — the device's id.
  - `{{randInt(a,b)}}` — a fresh random integer in `[a, b]`.
  - `{{counter}}` — an incrementing per-fire counter.
- **Optional schedule** — fire every *N* `seconds | minutes | hours | days`; when on
  and the device is enabled, the engine fires it automatically.
- **On-demand fire** — any event can be fired immediately from the UI or REST (see
  [11 — Fire event](./11-fire-event.md)).

## How it works

Events are stored as JSON on the device record and parsed by the engine
(`domainsvc.ParseEvents`). For a scheduled event, `ScheduleInterval` turns the
`{enabled, every, unit}` schedule into a duration, and the scheduler's
`buildDesired` adds one timed job per enabled device + scheduled event. At fire time
the payload, URL and topic templates are rendered (`domainsvc.Render`) so each send
carries fresh values, then dispatched through the device's live session or the
one-shot dispatcher.

An event holds the protocol config it applies to — `http` (method/path/headers/body),
`mqtt` (topic/qos/retain/body) or `lorawan` (fport/confirmed/payloadHex) — plus the
optional `schedule`.

## Notes

- Templates are rendered **per fire**, not stored rendered — so a scheduled event
  sends different values each tick.
- No schedule (or schedule off) means the event only fires on demand.
- A scheduled LoRaWAN device whose gateway is offline still "fires", but the engine
  reports a `gateway-offline` status frame instead of sending (see
  [09 — Console](./09-console.md)).

---
> Part of the [MapexOS ecosystem](../README.md#part-of-the-mapexos-ecosystem).

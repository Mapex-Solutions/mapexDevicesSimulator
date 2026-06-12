# 03 — HTTP protocol

The simplest device: it pushes uplinks to an HTTP endpoint. There is no connection to
keep — every fire is a one-shot request, and the response status comes back on the
frame. HTTP devices never receive.

## Capabilities

- **Send-only uplinks** to any HTTP(S) endpoint.
- **Method** — `POST` or `PUT`.
- **URL + path** — a base URL on the device plus a per-event path.
- **Headers** — device-level and per-event headers merged.
- **Auth** — none, **API key** (custom header) or **HTTP Basic**.
- **Templated body** — raw JSON or assembled form fields, rendered at send time.
- **Response status on the frame** — `200`, `4xx`, `5xx` or `timeout` is shown as the
  uplink frame's status.

## How it works

HTTP is **stateless**: it has no live session, so it is dispatched through the engine's
one-shot dispatcher registry rather than a held connection. At fire time
`buildSendSpec` resolves the request from the device config + event:

- config (`kind: "http"`): `url`, `method`, `headers`, `authMode`
  (`none|apiKey|basic`), `apiKeyHeader`, `apiKey`, `basicUser`, `basicPass`.
- event (`http`): `method`, `path`, `headers`, body (`raw` or `form`).

The URL is `device.url + event.path`, headers are device + event + the auth header,
and the body template is rendered (`{{deviceId}}`, `{{randInt(a,b)}}`, `{{counter}}`).
`report()` then streams an `up` frame carrying the HTTP response status, and persists
a log line when the device has `storeLogs` on.

## Notes

- No session, no reconnect, no downlinks — HTTP is the only protocol with no live
  connection.
- The response status is surfaced verbatim, so a non-2xx target shows up clearly on
  the console rather than failing silently.

---
> Part of the [MapexOS ecosystem](../README.md#part-of-the-mapexos-ecosystem).

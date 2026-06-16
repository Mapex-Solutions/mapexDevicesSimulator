# Journey: connection_status

> 🇧🇷 Versão em português: [README_pt.md](./README_pt.md)

Exercises the engine's connection-status lifecycle over the realtime console
WebSocket: an enabled device pointed at an unreachable broker must surface
`connecting` / `reconnecting` frames live. These frames exist **only** on `/ws` —
they are never written to the logs — so the console stream is the only way to
observe them.

## Flow

1. **StartConsoleStream** — connect to `/ws` first.
2. **CreateUnreachableMQTTDevice** — an enabled MQTT device whose broker is a
   loopback port nothing listens on, so its session never connects.
3. **AssertConsoleReconnecting** — wait for a `system`/`status` frame with status
   `connecting` or `reconnecting` for the device on the live stream.
4. **Compensation** — delete the device, close the stream.

## What it proves

The engine surfaces the connection lifecycle live: a device that cannot reach its
broker cycles through `connecting` / `reconnecting` with backoff, visible to the
user on the console even though nothing is logged.

## Run

```bash
# from e2e_tests/  (sidecar must be running on 127.0.0.1:5055)
go test -tags=saga ./journey/iot/connection_status/ -v
```

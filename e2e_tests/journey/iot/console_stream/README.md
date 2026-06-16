# Journey: console_stream

> 🇧🇷 Versão em português: [README_pt.md](./README_pt.md)

Exercises the realtime console WebSocket: a fired uplink must be broadcast live on
`/ws`, not only persisted to the logs.

## Flow

1. **StartConsoleStream** — connect to `/ws` and start collecting frames, before
   any frame is produced.
2. **StartEcho** — in-test echo target.
3. **CreateHTTPDevice** — device targeting the echo.
4. **FireTelemetry** — one uplink dispatched.
5. **AssertConsoleUpFrame** — wait for an `up`/`data` frame for the device on the
   live stream.
6. **Compensation** — delete the device, close the echo and the stream.

## What it proves

The realtime path: an uplink reaches the console WebSocket live, in the shape the
UI consumes (`up`/`data` with the device id).

## Run

```bash
# from e2e_tests/  (sidecar must be running on 127.0.0.1:5055)
go test -tags=saga ./journey/iot/console_stream/ -v
```

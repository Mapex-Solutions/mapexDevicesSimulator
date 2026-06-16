# Journey: fire_error

> 🇧🇷 Versão em português: [README_pt.md](./README_pt.md)

Exercises the engine's send-error handling: a device fired at an unreachable
target must surface the failure as an error frame in the logs, not drop it
silently.

## Flow

1. **CreateUnreachableHTTPDevice** — an enabled HTTP device whose target is a
   loopback port nothing listens on, `storeLogs` on.
2. **FireTelemetry** — one send, which fails to connect.
3. **AssertFireErrorLogged** — polls GET `/api/logs` until a frame with status
   `error` and the failure reason in its `response` appears.
4. **Compensation** — delete the device.

## What it proves

A failed send is reported, not swallowed: the engine records an `error` frame
with the cause, so the UI can show the user what went wrong.

## Run

```bash
# from e2e_tests/  (sidecar must be running on 127.0.0.1:5055)
go test -tags=saga ./journey/iot/fire_error/ -v
```

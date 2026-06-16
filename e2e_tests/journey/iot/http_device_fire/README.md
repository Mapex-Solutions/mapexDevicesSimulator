# Journey: http_device_fire

> 🇧🇷 Versão em português: [README_pt.md](./README_pt.md)

Exercises an HTTP device end to end against the live simulator.

## Flow

1. **StartEcho** — boot an in-process echo target (`httptest`) and publish its
   URL to the bag.
2. **CreateHTTPDevice** — POST `/api/devices` a send-only HTTP device targeting
   that echo, with one pre-registered event carrying a templated JSON body and
   `storeLogs` on.
3. **FireTelemetry** — POST `/api/devices/{id}/fire`, dispatching one uplink.
4. **AssertHTTPUplinkLogged** — poll GET `/api/logs?device=...` until a `up`/`data`
   frame appears with status `200` and a non-empty `response`. A 200 can only come
   from the live echo, so this doubles as the round-trip check.
5. **Compensation** — DELETE the device and close the echo, leaving the simulator
   clean.

## What it proves

Devices CRUD + the engine fire path + HTTP response capture + the logs read — the
whole "fire reaches the target and is persisted" trail.

## Run

```bash
# from e2e_tests/  (sidecar must be running on 127.0.0.1:5055)
go test -tags=saga ./journey/iot/http_device_fire/ -v
```

Self-contained: the echo target is started in-process, so only the sidecar needs
to be up.

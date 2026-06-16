# Journey: mqtt_downlink

> 🇧🇷 Versão em português: [README_pt.md](./README_pt.md)

Exercises the MQTT **inbound** path end to end against the live simulator and an
in-process broker: a device with receiving enabled subscribes to a topic, an
external publish arrives on it, and the simulator surfaces it as a downlink in
the logs.

## Flow

1. **StartMQTTBroker** — boots the in-process broker and publishes its
   coordinates to the bag.
2. **CreateMQTTReceiveDevice** — POST `/api/devices` an enabled MQTT device with
   `receiveEnabled` and a subscription to a per-run downlink topic, `storeLogs`
   on. Enabling it opens a session that subscribes.
3. **PublishDownlink** — injects a **retained** message on that topic, as an
   external party would.
4. **AssertDownlinkLogged** — polls GET `/api/logs` until a `down`/`downlink`
   frame appears carrying the published payload (which embeds the run id).
5. **Compensation** — delete the device, close the broker.

The publish is **retained**, so it is delivered whether it lands before or after
the device's subscription becomes active — the session subscribes asynchronously,
and a retained message removes that race without a settle step.

## What it proves

The engine's inbound half: a subscription is opened on enable, a received message
is decoded, and it is surfaced as a `down`/`downlink` frame and persisted to the
logs — the whole "data arrives at the device and is recorded" trail.

## Run

```bash
# from e2e_tests/  (sidecar must be running on 127.0.0.1:5055)
go test -tags=saga ./journey/iot/mqtt_downlink/ -v
```

Self-contained: the broker is started in-process, so only the sidecar needs to
be up.

# Journey: mqtt_device_fire

> 🇧🇷 Versão em português: [README_pt.md](./README_pt.md)

Exercises MQTT devices end to end against the live simulator and an in-process
broker, covering **both** auth modes the platform supports.

## Flow

One ordered saga starts a single broker (with both listeners) and then runs each
auth mode in turn. The oracle is the broker itself, not the logs API.

1. **StartMQTTBroker** — boots the in-process broker and publishes its
   coordinates and a fresh CA / client certificate to the bag.
2. **username/password** — `CreateMQTTUserPassDevice` (device on the plain
   `tcp://` listener, authenticating with a username and password) →
   `FireTelemetry` → `AssertMQTTPublished` (the broker accepted a publish carrying
   the device's `deviceId`; it only records a publish after the CONNECT passed
   auth, so this proves the credentials were honored).
3. **certificate** — `CreateMQTTTLSDevice` (device on the `ssl://` listener,
   authenticating with a **client certificate**, mutual TLS) → `FireTelemetry` →
   `AssertMQTTPublished`. The broker is configured with
   `RequireAndVerifyClientCert`, so the publish can only land if the handshake
   validated the device's cert against the run's CA.
4. **Compensation** — delete both devices, close the broker.

The fire happens the instant the device is enabled — before its persistent
session has connected — so it deliberately exercises the engine's one-shot
fallback. That path now uses a distinct client id, so it no longer collides with
the connecting session and the publish lands reliably (no settle step needed).

## What it proves

MQTT device CRUD + the engine fire path + **both** MQTT auth modes
(username/password and client certificate) reaching a real broker and being
accepted. A rejected CONNECT (wrong password, untrusted cert) surfaces as the
publish never arriving.

## How the broker fixture works

`common/utils/mqtt_broker.go` boots an embedded [mochi-mqtt](https://github.com/mochi-mqtt/server)
broker with two listeners on random loopback ports:

- a plain TCP listener that authenticates by username/password;
- a TLS listener that requires and verifies a client certificate.

`common/utils/certs.go` mints a fresh CA per run that signs both the broker's
server cert and the device's client cert, so no certificate fixtures live in the
tree. Every accepted publish is captured for the assert.

## Run

```bash
# from e2e_tests/  (sidecar must be running on 127.0.0.1:5055)
go test -tags=saga ./journey/iot/mqtt_device_fire/ -v
```

Self-contained: the broker and all certificates are created in-process, so only
the sidecar needs to be up. The sidecar must include MQTT TLS support (the engine
builds the `tls.Config` from the device's PEM material).

## Note on the one-shot fallback

This journey fires before the device's session has connected on purpose: it
covers the engine's one-shot fallback path. That path gives each one-shot
connection a distinct client id (`<clientId>-oneshot-N`) so it never collides
with the device's still-connecting persistent session — without that, a broker
would kick one of the two same-id connections and the fire would be silently
lost during the connect window.

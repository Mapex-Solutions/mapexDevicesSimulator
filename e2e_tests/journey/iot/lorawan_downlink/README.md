# Journey: lorawan_downlink

> 🇧🇷 Versão em português: [README_pt.md](./README_pt.md)

Exercises the LoRaWAN **inbound** path end to end against the live simulator and
a pinned ChirpStack LNS: a downlink queued on the LNS is delivered to the device
in its RX window after an uplink, and the simulator surfaces it as a downlink in
the logs. The journey owns the ChirpStack stack lifecycle.

## Flow

1. **StartStack** — `docker compose up` the pinned ChirpStack stack and connect
   the gRPC client.
2. **EnsureApplicationContext** — tenant, application, EU868 / 1.0.3 OTAA profile.
3. **ProvisionUDPDevice** — LNS gateway + device + keys.
4. **CreateUDPGateway** / **CreateLoRaWANDevice** — the simulator UDP gateway and
   device; enabling the device joins.
5. **AssertJoinAccepted** — ChirpStack assigned a DevAddr.
6. **EnqueueDownlink** — queue a downlink (fixed bytes) on the LNS for the device.
7. **FireTelemetry** — an uplink opens the Class A RX window; the LNS sends the
   queued downlink in it.
8. **AssertLoRaWANDownlinkReceived** — poll GET `/api/logs` until a `down`/`downlink`
   frame appears whose payload is the hex of the queued bytes.
9. **Compensation** — delete device/gateway, then `down -v` the stack.

The downlink is verified **on the simulator's logs**, and it can only appear
there if the LNS sent it and the device received it in its RX window — the full
inbound round trip.

## What it proves

The simulator's LoRaWAN inbound half: a Class A device receives a queued LNS
downlink in the RX window following an uplink, decodes it, and surfaces it as a
`down`/`downlink` frame.

## Run

```bash
# from e2e_tests/  (sidecar must be running on 127.0.0.1:5055; docker required)
go test -tags=saga ./journey/iot/lorawan_downlink/ -v
```

The journey brings the ChirpStack stack
([deployment/chirpstack](../../../../deployment/chirpstack)) up and down itself;
only the simulator sidecar needs to be running beforehand.

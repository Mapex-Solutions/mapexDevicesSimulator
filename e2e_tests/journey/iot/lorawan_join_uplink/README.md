# Journey: lorawan_join_uplink

> 🇧🇷 Versão em português: [README_pt.md](./README_pt.md)

Exercises LoRaWAN OTAA devices end to end against the live simulator and a
pinned ChirpStack LNS, over **both** radio transports: Semtech UDP and Basics
Station. The journey owns the ChirpStack stack lifecycle — it brings the stack
up at the start and tears it down with its volumes at the end.

## Flow

One ordered saga, one stack:

1. **StartStack** — `docker compose up` the pinned ChirpStack stack and connect
   the gRPC API client (the login poll doubles as the readiness gate).
2. **EnsureApplicationContext** — create the tenant, application, and an
   EU868 / LoRaWAN 1.0.3 OTAA device profile.
3. **UDP transport** — `ProvisionUDPDevice` (LNS gateway + device + keys) →
   `CreateUDPGateway` (simulator Semtech UDP gateway) → `CreateLoRaWANDevice`
   (simulator device; enabling it joins) → `AssertJoinAccepted` (ChirpStack
   assigned a DevAddr) → `FireTelemetry` → `AssertUplinkReceived` (ChirpStack
   recorded the uplink).
4. **Basics Station transport** — `ProvisionBasicStationDevice` →
   `CreateBasicStationDevice` (the device carries its own WebSocket link to the
   bridge, no separate gateway) → join / fire / uplink asserts.
5. **Compensation** — delete every device and gateway, then `down -v` the stack.

The join and uplink are verified **on ChirpStack over its gRPC API**, never on
the simulator, so a passing assert proves the radio path reached the LNS and the
LNS accepted it.

## What it proves

The simulator's LoRaWAN engine joins and uplinks against a real LNS over both
transports: OTAA join (gateway link + keys line up) and a data uplink, end to
end, with the gateway registered and frames routed through ChirpStack.

## Per-run isolation

Device and gateway EUIs are derived from the run id. The simulator sidecar is a
long-lived shared process that remembers per-DevEUI join state, and ChirpStack
remembers OTAA DevNonces — fresh EUIs every run keep repeated runs from looking
like replays. The `down -v` teardown is belt-and-suspenders on top of that.

## Run

```bash
# from e2e_tests/  (sidecar must be running on 127.0.0.1:5055; docker required)
go test -tags=saga ./journey/iot/lorawan_join_uplink/ -v
```

The journey brings the ChirpStack stack
([deployment/chirpstack](../../../../deployment/chirpstack)) up and down itself;
only the simulator sidecar needs to be running beforehand. The stack binds a
private, loopback-only port range (gRPC 18080, UDP 11700, Basics Station 13001)
so it never collides with a developer's own ChirpStack.

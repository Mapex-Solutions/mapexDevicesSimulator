# ChirpStack deployment (e2e LoRaWAN stack)

> 🇧🇷 Versão em português: [README_pt.md](./README_pt.md)

A self-contained, version-pinned ChirpStack stack the LoRaWAN e2e journey
provisions against. It is **not** a production deployment — it exists so the
[`lorawan_join_uplink`](../../e2e_tests/journey/iot/lorawan_join_uplink) journey
has a real LNS to join and uplink against.

## Why pinned

ChirpStack is a third-party LNS we do not own, and the journey drives it over its
gRPC API. The image is pinned (`chirpstack/chirpstack:4.18.0`,
`chirpstack-gateway-bridge:4.1.1`) so a server-side change never silently breaks
the test, and the e2e Go API client is pinned to the matching `v4.18.0`.

## What runs

| Service | Purpose |
|---------|---------|
| `chirpstack` | the LNS + gRPC API (provisioning) |
| `chirpstack-gateway-bridge` | Semtech UDP packet-forwarder ingest |
| `chirpstack-gateway-bridge-basicstation` | Basics Station WebSocket ingest |
| `mosquitto` | internal MQTT bus (chirpstack ↔ bridges) |
| `postgres`, `redis` | chirpstack storage |

Only EU868 is enabled, so the stack boots fast.

## Ports (loopback, private range)

Host ports are deliberately remapped off the standard ChirpStack ports so this
stack never collides with a developer's own `chirpstack-docker`:

| Endpoint | Host | In-container |
|----------|------|--------------|
| gRPC API | `127.0.0.1:18080` | 8080 |
| Semtech UDP | `127.0.0.1:11700/udp` | 1700 |
| Basics Station | `127.0.0.1:13001` | 3001 |

MQTT stays internal to the compose network.

## Lifecycle

The journey owns it: it runs `up -d --wait` at the start and `down -v` at the
end, so every run gets a fresh LNS (which also wipes ChirpStack's remembered OTAA
DevNonces). To drive it by hand:

```bash
docker compose -f deployment/chirpstack/chirpstack.yml up -d --wait
docker compose -f deployment/chirpstack/chirpstack.yml down -v
```

Default credentials are `admin` / `admin`. The configuration under
`configuration/` is the upstream `chirpstack-docker` example trimmed to one
region; only the compose file and `chirpstack.toml` are tailored.

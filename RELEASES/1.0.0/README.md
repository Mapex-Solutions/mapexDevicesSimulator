# Release 1.0.0 — Mapex Devices Simulator

> 🇧🇷 Versão em português: [README_pt.md](./README_pt.md)

Local desktop tool to simulate real IoT devices sending (and receiving) live traffic
over **HTTP**, **MQTT**, and **LoRaWAN** (Semtech UDP and Basics Station).

It exists to **complement the MapexOS system**: a way to drive real device traffic
into a MapexOS stack — its broker, its LNS, its core services — without owning a
single piece of hardware, so the platform can be demoed, exercised and validated
end to end. It speaks standard protocols, so it works against any LNS or broker too.

- **Status:** in progress (catalog below tracks each functionality)
- **Stack:** Electron + Go sidecar (Fiber, SQLite) · Vue 3 / Quasar / Pinia / Zod
- **Date:** 2026-06-12

---

## Part of the MapexOS ecosystem

This simulator is one piece of **MapexOS** — Mapex Solutions' IoT platform. Go meet
the rest:

- **[mapexOS](https://github.com/Mapex-Solutions/mapexOS)** — the core: the platform
  and its services.
- **[mapexOSDeploy](https://github.com/Mapex-Solutions/mapexOSDeploy)** — deployment
  of the MapexOS stack.
- **[mapexMQTTBroker](https://github.com/Mapex-Solutions/mapexMQTTBroker)** — the MQTT
  broker the simulator's MQTT devices can talk to.
- **[mapexLNS](https://github.com/Mapex-Solutions/mapexLNS)** — the LoRaWAN Network
  Server the simulator's LoRaWAN devices join.

The simulator generates the device side; these projects are the platform that
receives, routes and processes it.

---

## Functionality catalog

Each item links to its own spec under [`features/`](./features/). Status legend:
✅ done · 🚧 partial · ⬜ planned.

| # | Functionality | Status |
|---|---------------|--------|
| 01 | [Device management (CRUD)](./features/01-devices.md) | ✅ |
| 02 | [Gateway management (CRUD)](./features/02-gateways.md) | ✅ |
| 03 | [HTTP protocol](./features/03-http.md) | ✅ |
| 04 | [MQTT protocol (publish + subscribe)](./features/04-mqtt.md) | ✅ |
| 05 | [LoRaWAN over Semtech UDP](./features/05-lorawan-udp.md) | ✅ |
| 06 | [LoRaWAN over Basics Station](./features/06-lorawan-basicstation.md) | ✅ |
| 07 | [Events, templating and scheduling](./features/07-events-scheduling.md) | ✅ |
| 08 | [Engine: scheduler + live sessions](./features/08-engine.md) | ✅ |
| 09 | [Console (live WebSocket stream)](./features/09-console.md) | ✅ |
| 10 | [Logs (persisted history)](./features/10-logs.md) | ✅ |
| 11 | [Fire event (on-demand)](./features/11-fire-event.md) | ✅ |
| 12 | [Desktop app + sidecar packaging](./features/12-desktop-app.md) | ✅ |
| 13 | [Internationalization (EN / PT-BR)](./features/13-i18n.md) | ✅ |

---

## What "1.0.0" covers (one-liners)

- **Devices** — create/edit/enable/delete simulated devices on 4 protocols, persisted in SQLite.
- **Gateways** — LoRaWAN gateways with a Semtech UDP or Basics Station link to the LNS; the enabled flag gates the devices riding on them.
- **HTTP** — send-only uplinks to any endpoint, with headers and API-key/basic auth; the response status comes back on the frame.
- **MQTT** — live broker session: publish uplinks and, with receive on, subscribe to topics and stream each message as a downlink. QoS 0/1/2, retain, user/pass or TLS-cert auth.
- **LoRaWAN** — OTAA and ABP, MAC 1.0.x and 1.1, regions, device class A/C, real crypto (join/uplink MIC, payload encrypt/decrypt), DevAddr/FCnt, downlink decode — over Semtech UDP or Basics Station, against any LNS.
- **Events** — pre-registered, templated payloads (`{{deviceId}}`, `{{randInt(a,b)}}`, `{{counter}}`) with an optional repeat schedule, plus on-demand fire.
- **Engine** — a scheduler (timed jobs) beside a session manager (persistent connections with infinite backoff reconnect); CRUD changes reconcile immediately.
- **Console** — a live WebSocket stream of every uplink, downlink and system/status frame, with filtering and a detail panel.
- **Logs** — a persisted, paginated, filterable history of device messages.
- **Desktop** — an Electron app that spawns the Go sidecar; cross-platform sidecar build for packaging.
- **i18n** — full English and Brazilian Portuguese.

---

## Notes

- Each [`features/NN-*.md`](./features/) documents **what the platform offers** in
  that area and **how it is built** (architecture, data model, key files) — not how
  to use it step by step. The hands-on, click-by-click tour lives in
  [`/quickTest`](../../quickTest/).
- Anything that does not fit those two parts (edge cases, gotchas, limits) goes under
  a **Notes** section in each doc.

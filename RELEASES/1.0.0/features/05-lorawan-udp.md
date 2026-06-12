# 05 — LoRaWAN over Semtech UDP

A real LoRaWAN sensor: it rides a gateway that speaks the classic Semtech UDP
packet-forwarder to the LNS, does an OTAA join (or an ABP activation), and exchanges
genuine, cryptographically-correct frames. The crypto is real — frames are decrypted
by a standard LNS.

## Capabilities

- **OTAA and ABP** activation.
- **MAC versions** 1.0.2, 1.0.3, 1.0.4 and 1.1.0.
- **Regions** — EU868, US915, AU915, AS923, CN470, IN865, KR920, RU864.
- **Device class A or C** (carried through; the LNS profile schedules the RX windows).
- **Real crypto** — join-request/accept MIC, uplink MIC (1.1 binds both network keys
  plus DR and channel), payload encrypt/decrypt, key derivation.
- **Identity** — DevEUI, JoinEUI, AppKey (plus NwkKey for 1.1); DevAddr is assigned by
  the LNS at join; monotonic FCnt.
- **Uplinks** — confirmed or unconfirmed, an application-port frame of raw bytes.
- **Downlinks** — received over the gateway link, decoded (FPort, FCnt, FRMPayload)
  and surfaced as `down` frames.
- **Simulated radio metadata** — frequency, SF/DR, RSSI and SNR on each frame.

## How it works

The LoRaWAN stack lives under `backend/packages/utils/lorawan/`, vendor-neutral and
with no `ttnpb` dependency:

- **crypto** — the TTS routines copied verbatim (Apache header kept) for MIC, key
  derivation and payload encrypt/decrypt.
- **codec** — a thin PHYPayload marshal/unmarshal over plain structs (join request,
  data uplink, downlink decode).
- **band** — per-region uplink/RX parameters and delays.
- **device** — the stateful session: join → derive keys → build uplink → decode
  downlink, for both 1.0.x and 1.1.
- **udp** — the Semtech UDP framing (PUSH_DATA / PULL_DATA / PULL_RESP / TX_ACK / stat).

The engine's LoRaWAN connector opens a **shared gateway link** (ref-counted UDP
socket), performs the OTAA join (sends the join request, awaits the join accept routed
back by the gateway), binds the device by its new DevAddr, sends a periodic keepalive
+ stat, and acks downlinks. The **gateway is the bridge** (one link, shared); the
**device is the brain** (keys, DevAddr, FCnt, crypto). Verified live against
ChirpStack v4.

## Notes

- The DevNonce is kept **monotonic across reconnects** (the device brain persists per
  DevEUI), because an LNS rejects a reused DevNonce as a replay.
- 1.0.x derives one legacy `NwkSKey` + `AppSKey` from the single AppKey; 1.1 derives
  the three network keys from `NwkKey` plus `AppSKey` from `AppKey`.
- The device class is metadata the LNS uses to schedule downlinks — set the same class
  on the LNS device profile.
- Works against any standard LNS; the natural target is **[mapexLNS](https://github.com/Mapex-Solutions/mapexLNS)**.

---
> Part of the [MapexOS ecosystem](../README.md#part-of-the-mapexos-ecosystem).

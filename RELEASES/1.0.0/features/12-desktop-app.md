# 12 — Desktop app + sidecar packaging

The simulator ships as a self-contained **desktop application**: an Electron shell
around a Go **sidecar** that holds the engine. There is nothing external to run — the
app launches the backend itself.

## Capabilities

- **Desktop app** for macOS, Windows and Linux (Electron + Quasar).
- **Auto-spawned sidecar** — the app starts the Go backend on launch, on a free port,
  and only opens the window once it is healthy.
- **Self-contained** — the Electron shell renders the bundled SPA and auto-starts the
  Go sidecar (which exposes only the control API and WS); no separate server to start.
- **Cross-platform sidecar build** — one command compiles the Go binary for every
  packaged target.

## How it works

The Electron main process picks a free port, spawns the sidecar
(`simulatord --addr 127.0.0.1 --port <p>`), polls `/api/health` until it answers, and
hands the renderer the API/WS base via a preload bridge (`window.__SIM__`). The SPA
resolves its base from that bridge on the desktop, or from `window.location.origin`
on the web/dev build (where a Vite proxy forwards `/api` and `/ws`).

`build-sidecar.mjs` cross-compiles the sidecar with **CGO disabled** (the SQLite
driver is pure Go, so no toolchains are needed) into
`src-electron/sidecar/bin/<platform>-<arch>/` for the five packaged targets
(linux x64/arm64, darwin x64/arm64, win32 x64); Quasar's `extraResource` ships them.

## Notes

- In dev, `quasar dev` serves the SPA and proxies to a backend on `:5055`
  (`go run ./service/src`); in the packaged app the Electron renderer loads the
  bundled SPA and calls the sidecar cross-origin (the API allows it via CORS).
- Pure-Go SQLite is the reason cross-compilation needs no per-target toolchain.

---
> Part of the [MapexOS ecosystem](../README.md#part-of-the-mapexos-ecosystem).

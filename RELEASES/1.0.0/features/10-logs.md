# 10 — Logs (persisted history)

Where the [console](./09-console.md) is the live, fleeting view, the logs are the
durable record: a persisted, queryable history of the messages a device produced, so
you can go back and inspect what was sent long after it scrolled off the console.

## Capabilities

- **Persisted history** of device messages in SQLite — survives restarts.
- **Pagination** — page through the full history at a chosen page size.
- **Filtering** — by protocol, kind, direction, device, and a free-text search over
  the summary, payload and device name.
- **Opt-in per device** — only devices with `storeLogs` on are recorded; the rest are
  live-only on the console.

## How it works

The `logs` module exposes a `LogWriter` port the engine appends to, a `logs` SQLite
table (indexed newest-first), and a read endpoint `GET /api/logs` taking a `LogQuery`
(`limit`, `offset`, `protocol`, `kind`, `direction`, `device`, `q`) and returning a
`LogPage` (`items`, `total`). The engine's `report` / inbound path writes a log line
alongside the console frame whenever the device has `storeLogs` on.

The frontend mirrors it as a Pinia store (`fetch`, `setPage`, `setItemsPerPage`,
`setFilters`) and a list page with the quick + advanced filters and pagination.

## Notes

- Console vs logs: the console is **live + ephemeral**, the logs are **durable +
  queryable** — the same frames, two different windows.
- A device with `storeLogs` off still appears on the live console; it just leaves no
  persisted trail.

---
> Part of the [MapexOS ecosystem](../README.md#part-of-the-mapexos-ecosystem).

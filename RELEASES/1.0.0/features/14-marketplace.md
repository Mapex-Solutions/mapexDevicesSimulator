# 14 — Device marketplace

A catalog of ready-to-use vendor devices the user browses and installs with one
click. Instead of filling in the create stepper by hand, you pick a model
(a Milesight EM300-TH, a Dragino LHT65, a generic HTTP beacon) and the simulator
adds it pre-configured — protocol, keys, attributes and a sample event already set.

The catalog lives **online**, served by the `mapexMarketplace` service, independent
of the local sidecar. Browsing is read-only; installing writes a real device into
the local engine through the existing device-create path — no new install endpoint.

## Capabilities

- **Browse an online catalog** — a paged grid of vendor device cards (name, model,
  vendor, protocol, reading types, icon), served by the marketplace service.
- **Filter** — by **protocol**, **reading type** and **manufacturer**, plus a free
  text **search** over name / model / vendor. Filter options come from the catalog
  itself (facets), so the UI only ever offers values that exist.
- **Detail sheet** — per model: full description, device image, **datasheet** and
  **manual** (a bundled PDF shown inline, or an external link opened online), and a
  **codec viewer** for the vendor's payload decoder(s).
- **One-click install** — turns a catalog model into a new simulated device:
  fetches its simulator template, mints a fresh `deviceId`, and creates it through
  the engine. It lands in the device list, ready to enable and fire.
- **Online / offline awareness** — the catalog's connectivity is tracked separately
  from the sidecar engine status; if the catalog can't be reached the grid clears
  and an offline state is shown, while installed devices keep working.

## How it works

A frontend-only feature riding two backends: the online catalog (read) and the
local engine (write).

**Frontend** (`frontend/app/src/`):

- **pages/marketplace** — `MarketplaceListPage.vue` (grid + filter rail + pagination),
  `MarketplaceCard.vue`, `MarketplaceDetailDialog.vue` (the detail sheet, codec
  viewer and the install action).
- **stores/marketplace** — `useMarketplaceStore`: `fetch(query)` loads a filtered
  page plus the facet options (once) and flips `status` to `online`/`offline`;
  `install(item)` mints the `deviceId` and delegates to the devices store.
- **services / api** — `marketplaceApi` (`packages/api`) is a read-only resource over
  the marketplace base; the wire shape is mirrored as Zod schemas in
  `packages/schema` (`marketplace.schema.ts`), so the catalog can grow new protocols
  and reading types without a schema change (they are plain strings).

**Catalog contract** — the `mapexMarketplace` Go service under `/api/v1/devices`,
responses in the Mapex `{ status, errors, data }` envelope:

- `GET /devices` — filtered, paged listing (`protocol`, `readingType`,
  `manufacturer`, `search`, `page`, `perPage`).
- `GET /devices/facets` — the available filter options.
- `GET /devices/:vendor/:slug` — the detail sheet (`device_information.json`).
- `GET /devices/:vendor/:slug/simulator` — the installable template
  (`device_simulator.json`): a `DeviceInput` minus the `deviceId`.
- `GET /devices/:vendor/:slug/codecs` — the model's payload codec(s).

The base is resolved by `resolveMarketplaceBase()` — the preload-bridged origin in
the packaged app (pointed at the deployed catalog), else the local marketplace
server on `:6060` in development; bundle assets (images, codec files) go through
`resolveMarketplaceAssetUrl()`. A POC mock catalog ships under
[`frontend/app/public/marketplace/`](../../../frontend/app/public/marketplace/)
(see its README for the bundle format).

**Install** reuses the [device CRUD](./01-devices.md): the store fetches the model's
simulator template, spreads it into a `DeviceInput`, adds a freshly minted
`deviceId` (`<model-slug>-<8 hex>`), and calls `devices.create`. The engine owns the
row id and raises the usual reconcile signal, so the new device is live immediately.

## Notes

- **Online-only by design.** The catalog is never bundled into the engine; an
  offline catalog yields an empty grid and an offline badge, distinct from the
  sidecar health indicator that governs installed devices.
- **No new install endpoint.** Installing is just a templated `devices.create`; the
  `deviceId` is minted client-side to stay readable and collision-resistant, and the
  template's keys are defaults the user can edit afterwards.
- **Codecs are surfaced, not run.** The detail sheet lets the user view/download a
  model's decoder; the simulator does not execute it yet.
- The shipped mock catalog is a POC (5 models across 3 vendors, 3 LoRaWAN codecs);
  production points the app at the deployed `mapexMarketplace` service.

---
> Part of the [MapexOS ecosystem](../README.md#part-of-the-mapexos-ecosystem).

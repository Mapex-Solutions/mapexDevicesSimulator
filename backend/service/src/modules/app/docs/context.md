# Bounded Context: App (composition root)

**Service:** simulator (mapexDevicesSimulator backend / `simulatord`)
**Module path:** `src/modules/app/`
**Owner:** MAPEX
**Last reviewed:** 2026-06-09

## Purpose
The app module is the composition root for the business modules. It owns no
domain of its own: its single job is to run each registered module's init phases
in order (repositories, then services, then interfaces) so providers exist before
consumers resolve them. Infrastructure wiring (config, logger, SQLite, the HTTP
app, static SPA, shutdown) lives in `src/bootstrap`; the module registry lives in
`src/shared/configuration/modules`.

## Ubiquitous Language
| Term | Meaning in this context | Not to be confused with |
|------|-------------------------|--------------------------|
| Module | A bounded context under `src/modules/` with its own init phases | Go package |
| Init phase | One of InitRepositories / InitServices / InitInterfaces | DI provider |

## Published Events (driven — outbound)
None. The app module emits no events.

## Consumed Events (driving — inbound)
None.

## Driving Ports (what can call this module)
- `InitModule()` is called once from `main` after the infrastructure is wired.

## Driven Ports (what this module requires)
- The module registry `shared/configuration/modules.Modules`.
- The DI container, populated by `bootstrap` before `InitModule` runs.

## Invariants and Business Rules
- Phases run across ALL modules in order: every InitRepositories, then every
  InitServices, then every InitInterfaces. A service may depend on any module's
  repository, and an interface on any module's service.
- The app module contains no domain, application, infrastructure, or interface
  layer of its own, so those folders are intentionally absent.

## Known Cross-Context Interactions
- Invokes every business module's init phases; depends on none of them directly.

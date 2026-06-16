# Jornada: http_device_fire

> 🇺🇸 English version: [README.md](./README.md)

Exercita um device HTTP de ponta a ponta contra o simulador vivo.

## Fluxo

1. **StartEcho** — sobe um echo de teste (`httptest`) no processo e publica a URL
   dele no bag.
2. **CreateHTTPDevice** — POST `/api/devices` de um device HTTP só-envio mirando
   esse echo, com um evento pré-cadastrado de body JSON templated e `storeLogs`
   ligado.
3. **FireTelemetry** — POST `/api/devices/{id}/fire`, despachando um uplink.
4. **AssertHTTPUplinkLogged** — faz polling em GET `/api/logs?device=...` até
   aparecer um frame `up`/`data` com status `200` e `response` não vazio. Um 200
   só pode vir do echo vivo, então isso também é a prova do ida-e-volta.
5. **Compensação** — DELETE do device e fecha o echo, deixando o simulador limpo.

## O que prova

CRUD de devices + o caminho de fire do engine + captura do response HTTP + a
leitura dos logs — todo o trilho "o disparo chega no alvo e é persistido".

## Rodar

```bash
# de e2e_tests/  (o sidecar precisa estar no ar em 127.0.0.1:5055)
go test -tags=saga ./journey/iot/http_device_fire/ -v
```

Auto-contido: o echo sobe no próprio processo, então só o sidecar precisa estar
no ar.

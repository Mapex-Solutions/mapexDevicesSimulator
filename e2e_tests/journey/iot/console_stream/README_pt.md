# Jornada: console_stream

> 🇺🇸 English version: [README.md](./README.md)

Exercita o WebSocket de console em tempo real: um uplink disparado precisa ser
transmitido ao vivo no `/ws`, não só persistido nos logs.

## Fluxo

1. **StartConsoleStream** — conecta no `/ws` e começa a coletar frames, antes de
   qualquer frame ser produzido.
2. **StartEcho** — alvo echo de teste.
3. **CreateHTTPDevice** — device mirando o echo.
4. **FireTelemetry** — um uplink despachado.
5. **AssertConsoleUpFrame** — espera um frame `up`/`data` do device no stream ao
   vivo.
6. **Compensação** — apaga o device, fecha o echo e o stream.

## O que prova

O caminho em tempo real: um uplink chega ao WebSocket de console ao vivo, no
formato que o UI consome (`up`/`data` com o device id).

## Rodar

```bash
# de e2e_tests/  (o sidecar precisa estar no ar em 127.0.0.1:5055)
go test -tags=saga ./journey/iot/console_stream/ -v
```

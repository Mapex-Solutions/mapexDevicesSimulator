# Jornada: connection_status

> 🇺🇸 English version: [README.md](./README.md)

Exercita o ciclo de status de conexão do engine pelo WebSocket de console em
tempo real: um device habilitado apontando pra um broker inalcançável precisa
surgir frames `connecting` / `reconnecting` ao vivo. Esses frames existem
**só** no `/ws` — nunca são escritos nos logs — então o stream de console é a
única forma de observá-los.

## Fluxo

1. **StartConsoleStream** — conecta no `/ws` primeiro.
2. **CreateUnreachableMQTTDevice** — um device MQTT habilitado cujo broker é uma
   porta loopback onde nada escuta, então a sessão nunca conecta.
3. **AssertConsoleReconnecting** — espera um frame `system`/`status` com status
   `connecting` ou `reconnecting` do device no stream ao vivo.
4. **Compensação** — apaga o device, fecha o stream.

## O que prova

O engine surge o ciclo de conexão ao vivo: um device que não alcança o broker
passa por `connecting` / `reconnecting` com backoff, visível ao usuário no
console mesmo sem nada ser logado.

## Rodar

```bash
# de e2e_tests/  (o sidecar precisa estar no ar em 127.0.0.1:5055)
go test -tags=saga ./journey/iot/connection_status/ -v
```

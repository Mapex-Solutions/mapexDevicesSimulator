# Jornada: fire_error

> 🇺🇸 English version: [README.md](./README.md)

Exercita o tratamento de erro de envio do engine: um device disparado contra um
alvo inalcançável precisa surgir a falha como um frame de erro nos logs, não
descartar em silêncio.

## Fluxo

1. **CreateUnreachableHTTPDevice** — um device HTTP habilitado cujo alvo é uma
   porta loopback onde nada escuta, `storeLogs` ligado.
2. **FireTelemetry** — um envio, que falha ao conectar.
3. **AssertFireErrorLogged** — faz polling em GET `/api/logs` até aparecer um
   frame com status `error` e a causa da falha no `response`.
4. **Compensação** — apaga o device.

## O que prova

Um envio que falha é reportado, não engolido: o engine registra um frame `error`
com a causa, então o UI consegue mostrar ao usuário o que deu errado.

## Rodar

```bash
# de e2e_tests/  (o sidecar precisa estar no ar em 127.0.0.1:5055)
go test -tags=saga ./journey/iot/fire_error/ -v
```

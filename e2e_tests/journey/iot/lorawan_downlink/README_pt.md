# Jornada: lorawan_downlink

> 🇺🇸 English version: [README.md](./README.md)

Exercita o caminho **de entrada** do LoRaWAN de ponta a ponta contra o simulador
vivo e um LNS ChirpStack com versão fixada: um downlink enfileirado no LNS é
entregue ao device na janela RX depois de um uplink, e o simulador surge isso como
um downlink nos logs. A jornada é dona do ciclo de vida do stack ChirpStack.

## Fluxo

1. **StartStack** — `docker compose up` do stack ChirpStack fixado e conecta o
   cliente gRPC.
2. **EnsureApplicationContext** — tenant, aplicação, device profile OTAA EU868 / 1.0.3.
3. **ProvisionUDPDevice** — gateway + device + chaves no LNS.
4. **CreateUDPGateway** / **CreateLoRaWANDevice** — o gateway UDP e o device no
   simulador; habilitar o device faz o join.
5. **AssertJoinAccepted** — ChirpStack atribuiu um DevAddr.
6. **EnqueueDownlink** — enfileira um downlink (bytes fixos) no LNS para o device.
7. **FireTelemetry** — um uplink abre a janela RX (Class A); o LNS manda o downlink
   enfileirado nela.
8. **AssertLoRaWANDownlinkReceived** — faz polling em GET `/api/logs` até aparecer
   um frame `down`/`downlink` cujo payload é o hex dos bytes enfileirados.
9. **Compensação** — apaga device/gateway, depois `down -v` do stack.

O downlink é verificado **nos logs do simulador**, e ele só pode aparecer lá se o
LNS o enviou e o device o recebeu na janela RX — o ida-e-volta de entrada completo.

## O que prova

A metade de entrada do LoRaWAN no simulador: um device Class A recebe um downlink
enfileirado do LNS na janela RX depois de um uplink, decodifica e surge como frame
`down`/`downlink`.

## Rodar

```bash
# de e2e_tests/  (o sidecar precisa estar no ar em 127.0.0.1:5055; docker necessário)
go test -tags=saga ./journey/iot/lorawan_downlink/ -v
```

A jornada sobe e derruba o stack ChirpStack
([deployment/chirpstack](../../../../deployment/chirpstack)) sozinha; só o sidecar
do simulador precisa estar rodando antes.

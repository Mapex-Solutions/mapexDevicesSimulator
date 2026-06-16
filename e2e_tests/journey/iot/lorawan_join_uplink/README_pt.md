# Jornada: lorawan_join_uplink

> 🇺🇸 English version: [README.md](./README.md)

Exercita devices LoRaWAN OTAA de ponta a ponta contra o simulador vivo e um LNS
ChirpStack com versão fixa, nos **dois** transportes de rádio: Semtech UDP e
Basics Station. A jornada é dona do ciclo de vida do stack ChirpStack — sobe o
stack no início e derruba com os volumes no fim.

## Fluxo

Uma saga ordenada, um stack:

1. **StartStack** — `docker compose up` do stack ChirpStack fixado e conecta o
   cliente da API gRPC (o polling de login serve de porteiro de readiness).
2. **EnsureApplicationContext** — cria o tenant, a aplicação e um device profile
   OTAA EU868 / LoRaWAN 1.0.3.
3. **Transporte UDP** — `ProvisionUDPDevice` (gateway + device + chaves no LNS) →
   `CreateUDPGateway` (gateway Semtech UDP no simulador) → `CreateLoRaWANDevice`
   (device no simulador; habilitar faz o join) → `AssertJoinAccepted` (ChirpStack
   atribuiu um DevAddr) → `FireTelemetry` → `AssertUplinkReceived` (ChirpStack
   registrou o uplink).
4. **Transporte Basics Station** — `ProvisionBasicStationDevice` →
   `CreateBasicStationDevice` (o device carrega seu próprio link WebSocket para o
   bridge, sem gateway separado) → asserts de join / fire / uplink.
5. **Compensação** — apaga todos os devices e gateways, depois `down -v` do stack.

O join e o uplink são verificados **no ChirpStack pela API gRPC dele**, nunca no
simulador, então um assert que passa prova que o caminho de rádio chegou ao LNS e
o LNS aceitou.

## O que prova

O engine LoRaWAN do simulador faz join e uplink contra um LNS real nos dois
transportes: join OTAA (link do gateway + chaves batem) e um uplink de dados, de
ponta a ponta, com o gateway registrado e os frames roteados pelo ChirpStack.

## Isolamento por execução

Os EUIs de device e gateway são derivados do run id. O sidecar do simulador é um
processo compartilhado de vida longa que lembra o estado de join por DevEUI, e o
ChirpStack lembra os DevNonces OTAA — EUIs novos a cada execução impedem que
execuções repetidas pareçam replay. O teardown `down -v` é reforço extra.

## Rodar

```bash
# de e2e_tests/  (o sidecar precisa estar no ar em 127.0.0.1:5055; docker necessário)
go test -tags=saga ./journey/iot/lorawan_join_uplink/ -v
```

A jornada sobe e derruba o stack ChirpStack
([deployment/chirpstack](../../../../deployment/chirpstack)) sozinha; só o sidecar
do simulador precisa estar rodando antes. O stack usa uma faixa de portas privada
e só em loopback (gRPC 18080, UDP 11700, Basics Station 13001) para nunca colidir
com o ChirpStack do desenvolvedor.

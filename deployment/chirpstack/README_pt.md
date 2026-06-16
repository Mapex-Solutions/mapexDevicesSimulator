# Deployment ChirpStack (stack LoRaWAN de e2e)

> 🇺🇸 English version: [README.md](./README.md)

Um stack ChirpStack auto-contido e com versão fixada, contra o qual a jornada
e2e de LoRaWAN faz provisionamento. **Não** é um deployment de produção — existe
para a jornada
[`lorawan_join_uplink`](../../e2e_tests/journey/iot/lorawan_join_uplink) ter um
LNS real para fazer join e uplink.

## Por que fixar a versão

ChirpStack é um LNS de terceiros que não é nosso, e a jornada o controla pela API
gRPC. A imagem é fixada (`chirpstack/chirpstack:4.18.0`,
`chirpstack-gateway-bridge:4.1.1`) para que uma mudança no servidor nunca quebre o
teste em silêncio, e o cliente Go da API no e2e é fixado no `v4.18.0` correspondente.

## O que roda

| Serviço | Função |
|---------|--------|
| `chirpstack` | o LNS + a API gRPC (provisionamento) |
| `chirpstack-gateway-bridge` | ingestão Semtech UDP (packet-forwarder) |
| `chirpstack-gateway-bridge-basicstation` | ingestão Basics Station (WebSocket) |
| `mosquitto` | barramento MQTT interno (chirpstack ↔ bridges) |
| `postgres`, `redis` | armazenamento do chirpstack |

Só EU868 está habilitado, então o stack sobe rápido.

## Portas (loopback, faixa privada)

As portas do host são remapeadas de propósito para fora das portas padrão do
ChirpStack, para este stack nunca colidir com o `chirpstack-docker` do desenvolvedor:

| Endpoint | Host | No container |
|----------|------|--------------|
| API gRPC | `127.0.0.1:18080` | 8080 |
| Semtech UDP | `127.0.0.1:11700/udp` | 1700 |
| Basics Station | `127.0.0.1:13001` | 3001 |

O MQTT fica interno à rede do compose.

## Ciclo de vida

A jornada é dona dele: roda `up -d --wait` no início e `down -v` no fim, então
cada execução ganha um LNS novo (o que também limpa os DevNonces OTAA lembrados
pelo ChirpStack). Para controlar na mão:

```bash
docker compose -f deployment/chirpstack/chirpstack.yml up -d --wait
docker compose -f deployment/chirpstack/chirpstack.yml down -v
```

As credenciais padrão são `admin` / `admin`. A configuração em `configuration/` é
o exemplo do `chirpstack-docker` upstream reduzido a uma região; só o arquivo
compose e o `chirpstack.toml` foram ajustados.

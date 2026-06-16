# Testes E2E do Simulador

> 🇺🇸 English version: [README.md](./README.md)

Testes ponta-a-ponta do Mapex Devices Simulator, seguindo a **mesma arquitetura**
do [e2e do MapexOS](../../../mapexOS/e2e_tests): um runner de saga pequeno, fixtures
compartilhadas por serviço, e jornadas que as compõem.

Como o simulador é um serviço local único e sem autenticação, o framework é o
padrão da plataforma menos o `ClientSet` multi-serviço e o bootstrap de IAM: um
client `Sim` contra o sidecar, e o "pronto" é só um `GET /api/health` verde.

## Estrutura

```
e2e_tests/
├── core/saga/          # o runner: Item (Step|Assert), Context (+ bag), Run
├── common/
│   ├── constants/      # SimURL + endpoints do ChirpStack
│   ├── types/          # o envelope { status, errors, data }
│   └── utils/          # health + fixtures (echo, broker MQTT, certs, stack ChirpStack)
├── services/           # services/{svc}/{mod}/ — steps / asserts / payloads
│   ├── simulator/      # o sidecar: devices, gateways, engine, logs, targets
│   └── chirpstack/     # o LNS: client, ciclo do stack, provisionamento
└── journey/iot/        # jornadas saga (build tag: saga)
    ├── http_device_fire/
    ├── mqtt_device_fire/
    ├── mqtt_downlink/
    ├── lorawan_join_uplink/
    ├── lorawan_downlink/
    ├── fire_error/
    ├── console_stream/
    └── connection_status/
```

- **Step** muta o simulador (POST/DELETE) e pode registrar um `Compensate` para
  limpeza; **Assert** lê só a API pública (nunca o arquivo SQLite).
- O runner percorre a jornada em ordem e roda cada `Compensate` ao contrário no
  fim — passando ou falhando — pra próxima rodada começar limpa.

## Pré-requisitos

O sidecar do simulador no ar em `127.0.0.1:5055` (o padrão do serviço; o endpoint
é fixo em `common/constants`, sem variáveis de ambiente pra setar):

```bash
cd ../backend
go run ./service/src --addr 127.0.0.1 --port 5055
```

Uma jornada que não acha o sidecar dá `t.Skip`, então a suíte nunca falha só
porque nada está rodando.

## Como rodar

```bash
cd e2e_tests

# Todas as jornadas (a build tag saga é obrigatória)
go test -tags=saga ./journey/...

# Uma jornada, verboso
go test -tags=saga ./journey/iot/http_device_fire/ -v
```

Os endpoints ficam em `common/constants` — mude lá, não pelo ambiente.

## Jornadas

| Jornada | O que prova | Precisa |
|---------|-------------|---------|
| [http_device_fire](./journey/iot/http_device_fire/) | Cria device HTTP, dispara, e confirma que o simulador loga o uplink 200 com o response capturado | só o sidecar (echo é in-test) |
| [mqtt_device_fire](./journey/iot/mqtt_device_fire/) | Cria devices MQTT, dispara, e confirma que um broker embarcado aceita o publish autenticado nos dois modos de auth (usuário/senha e certificado de cliente) | só o sidecar (broker + certs são in-test) |
| [mqtt_downlink](./journey/iot/mqtt_downlink/) | Device MQTT com recebimento assina, um publish externo (retained) chega, e o simulador loga como downlink — o caminho de entrada | só o sidecar (broker é in-test) |
| [lorawan_join_uplink](./journey/iot/lorawan_join_uplink/) | Provisiona o ChirpStack e dirige devices LoRaWAN OTAA por Semtech UDP e Basics Station, confirmando que o LNS registra o join e o uplink de cada um | sidecar + docker (a jornada é dona do stack ChirpStack) |
| [lorawan_downlink](./journey/iot/lorawan_downlink/) | Enfileira um downlink no ChirpStack, dispara um uplink pra abrir a janela RX, e confirma que o simulador loga o downlink recebido — o caminho de entrada | sidecar + docker (a jornada é dona do stack ChirpStack) |
| [fire_error](./journey/iot/fire_error/) | Dispara um device contra um alvo inalcançável e confirma que o engine loga a falha como frame de erro, não em silêncio | só o sidecar |
| [console_stream](./journey/iot/console_stream/) | Conecta no console `/ws`, dispara um device, e confirma que o uplink chega ao vivo no stream | só o sidecar |
| [connection_status](./journey/iot/connection_status/) | Habilita um device apontando pra um broker inalcançável e confirma que frames `connecting`/`reconnecting` chegam ao vivo no `/ws` (nunca são logados) | só o sidecar |

## Testes de contrato de módulo (`go test ./...`)

Testes mais leves do contrato HTTP de um módulo só, em `services/{svc}/{mod}/e2e/`.
Não têm a build tag saga — `go test ./...` roda eles e dá skip limpo se o sidecar
estiver caído.

| Módulo | Cobre |
|--------|-------|
| `services/simulator/devices/e2e` | CRUD (+ validação, 404) e membership da list |
| `services/simulator/gateways/e2e` | CRUD (+ validação, 404) e membership da list |
| `services/simulator/logs/e2e` | paginação por cursor + filtros (device, event, data, combinado) |

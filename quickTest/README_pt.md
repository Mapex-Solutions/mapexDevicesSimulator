# Quick Test — opere o simulador de ponta a ponta

> 🇺🇸 English version: [README.md](./README.md)

Esta pasta é um tour prático do **Mapex Devices Simulator**. Cada protocolo tem a sua
própria pasta com **valores prontos para copiar e colar** e um passo a passo, para você
criar um dispositivo, ligá-lo, disparar um evento e acompanhar o tráfego ao vivo.

**📖 Tutorial guiado com screenshots:**
[STEP-BY-STEP.en.md](./STEP-BY-STEP.en.md) · [STEP-BY-STEP.pt-BR.md](./STEP-BY-STEP.pt-BR.md)

| Protocolo | Pasta | O que mostra |
|-----------|-------|--------------|
| HTTP | [`http/`](./http/) | Uplink único (one-shot) para qualquer endpoint HTTP |
| MQTT | [`mqtt/`](./mqtt/) | Publish **e** subscribe (downlink) ao vivo por um broker |
| LoRaWAN por Semtech UDP | [`lorawan-udp/`](./lorawan-udp/) | Um gateway + sensor fazendo join num LNS pelo packet-forwarder UDP |
| LoRaWAN por Basics Station | [`lorawan-basic-station/`](./lorawan-basic-station/) | O mesmo sensor carregando o próprio link WebSocket Basics Station |

---

## 1. Suba o app

```bash
# a partir da raiz do repo
cd frontend
npm install            # só na primeira vez
npm run dev:electron   # compila o sidecar Go e abre o app desktop
```

O `dev:electron` compila a engine no sidecar e a janela desktop o sobe
automaticamente. Para rodar no navegador:

```bash
cd frontend && npm run dev      # abre http://localhost:9100 (SPA)
# em outro terminal, suba a engine com que a UI conversa:
cd backend && go run ./service/src --addr 127.0.0.1 --port 5055
```

O cabeçalho mostra um **ponto de conexão**: verde = a engine está acessível, cinza =
offline. Sem engine, as listas ficam simplesmente vazias (não há dado fake/seed).

---

## 2. O fluxo que você repete em todo protocolo

1. **Devices → New device** — dê nome, um Device ID e escolha o protocolo.
2. **Connection** — preencha os campos de conexão do protocolo (copie do README da
   pasta).
3. **Events** — adicione um evento com o payload a enviar (copie do README da pasta).
4. **Save**, depois ligue o **Enabled** do dispositivo na lista.
5. **Console** — abra para ver os frames ao vivo (`up`, `down`, `system`).
6. **Fire event** — pela linha do dispositivo ou pelo console, envie sob demanda e veja
   o eco `up` (e qualquer resposta `down`) aparecer na hora.

> LoRaWAN também precisa de um **Gateway** (Gateways → New gateway) no qual o sensor
> trafega — exceto Basics Station, em que o sensor carrega o próprio link.

---

## 3. Screenshots

Os passo a passo referenciam imagens em `images/` de cada pasta.

---

## 4. Precisa de um LNS para LoRaWAN?

As pastas de LoRaWAN foram validadas contra o **ChirpStack v4** via `chirpstack-docker`.
O LNS local mais rápido:

```bash
git clone https://github.com/chirpstack/chirpstack-docker
cd chirpstack-docker && docker compose up -d
# UI http://localhost:8088 (admin/admin) · UDP :1700 · Basics Station :3001
```

Cada pasta de LoRaWAN lista exatamente o que registrar no ChirpStack para as chaves
baterem.

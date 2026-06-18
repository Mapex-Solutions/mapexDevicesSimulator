# Release 1.0.0 — Mapex Devices Simulator

> 🇺🇸 English version: [README.md](./README.md)

Ferramenta desktop local para simular dispositivos IoT reais enviando (e recebendo)
tráfego ao vivo por **HTTP**, **MQTT** e **LoRaWAN** (Semtech UDP e Basics Station).

Ela existe para **complementar o sistema MapexOS**: um jeito de injetar tráfego real
de dispositivos num stack MapexOS — o broker, o LNS, os serviços de core — sem possuir
um único hardware, para a plataforma ser demonstrada, exercitada e validada de ponta a
ponta. Fala protocolos padrão, então também funciona contra qualquer LNS ou broker.

- **Status:** em andamento (o catálogo abaixo acompanha cada funcionalidade)
- **Stack:** Electron + sidecar Go (Fiber, SQLite) · Vue 3 / Quasar / Pinia / Zod
- **Data:** 2026-06-16

---

## Parte do ecossistema MapexOS

Este simulador é uma peça do **MapexOS** — a plataforma IoT da Mapex Solutions. Conheça
o resto:

- **[mapexOS](https://github.com/Mapex-Solutions/mapexOS)** — o core: a plataforma e
  seus serviços.
- **[mapexOSDeploy](https://github.com/Mapex-Solutions/mapexOSDeploy)** — deploy do
  stack MapexOS.
- **[mapexMQTTBroker](https://github.com/Mapex-Solutions/mapexMQTTBroker)** — o broker
  MQTT com que os dispositivos MQTT do simulador conversam.
- **[mapexLNS](https://github.com/Mapex-Solutions/mapexLNS)** — o LoRaWAN Network Server
  no qual os dispositivos LoRaWAN do simulador fazem join.

O simulador gera o lado do dispositivo; esses projetos são a plataforma que recebe,
roteia e processa.

---

## Instalação

Baixe o instalador da sua plataforma e execute — não há mais nada para configurar, o
motor Go já vem dentro do app.

| Plataforma | Arquivo | Instalar |
|------------|---------|----------|
| Debian / Ubuntu | `.deb` | `sudo apt install ./mapex-devices-simulator_<versão>_amd64.deb` |
| RedHat / Fedora | `.rpm` | `sudo dnf install ./mapex-devices-simulator-<versão>.x86_64.rpm` |
| macOS | `.dmg` | abra o `.dmg` e arraste o app para **Aplicativos** |
| Windows | `.exe` | execute o instalador e siga o assistente |

### Primeira execução — o aviso de "desenvolvedor não identificado"

Os instaladores **não são assinados** (a assinatura é opcional e paga). Num download
direto, isso faz o sistema mostrar um **aviso de segurança único** na primeira
execução. O app é seguro — é só permitir:

- **Linux** — sem aviso; instala e roda.
- **Windows** — o SmartScreen mostra *"O Windows protegeu seu PC"* → clique em **Mais
  informações** → **Executar assim mesmo**.
- **macOS** — o Gatekeeper aponta *desenvolvedor não identificado* → **botão direito no
  app → Abrir** e confirme, ou vá em **Ajustes do Sistema → Privacidade e Segurança →
  Abrir mesmo assim**.

Assinar os instaladores (notarização Apple / code-signing Windows) removeria esses
avisos, mas não é necessário para instalar ou rodar o app.

---

## Catálogo de funcionalidades

Cada item tem o seu próprio doc em [`features/`](./features/). Legenda de status:
✅ pronto · 🚧 parcial · ⬜ planejado.

| # | Funcionalidade | Status |
|---|----------------|--------|
| 01 | [Gestão de dispositivos (CRUD)](./features/01-devices.md) | ✅ |
| 02 | [Gestão de gateways (CRUD)](./features/02-gateways.md) | ✅ |
| 03 | [Protocolo HTTP](./features/03-http.md) | ✅ |
| 04 | [Protocolo MQTT (publish + subscribe)](./features/04-mqtt.md) | ✅ |
| 05 | [LoRaWAN por Semtech UDP](./features/05-lorawan-udp.md) | ✅ |
| 06 | [LoRaWAN por Basics Station](./features/06-lorawan-basicstation.md) | ✅ |
| 07 | [Eventos, templating e agendamento](./features/07-events-scheduling.md) | ✅ |
| 08 | [Engine: scheduler + sessões vivas](./features/08-engine.md) | ✅ |
| 09 | [Console (stream WebSocket ao vivo)](./features/09-console.md) | ✅ |
| 10 | [Logs (histórico persistido)](./features/10-logs.md) | ✅ |
| 11 | [Disparo de evento (sob demanda)](./features/11-fire-event.md) | ✅ |
| 12 | [App desktop + empacotamento do sidecar](./features/12-desktop-app.md) | ✅ |
| 13 | [Internacionalização (EN / PT-BR)](./features/13-i18n.md) | ✅ |
| 14 | [Marketplace de dispositivos (navegar + instalar)](./features/14-marketplace.md) | ✅ |

---

## O que a "1.0.0" cobre (resumo)

- **Dispositivos** — criar/editar/habilitar/excluir dispositivos simulados em 4 protocolos, persistidos em SQLite.
- **Gateways** — gateways LoRaWAN com link Semtech UDP ou Basics Station até o LNS; a flag de habilitado gateia os dispositivos que trafegam por ele.
- **HTTP** — uplinks só de envio para qualquer endpoint, com headers e auth API-key/basic; o status da resposta volta no frame.
- **MQTT** — sessão viva no broker: publica uplinks e, com recebimento ligado, assina tópicos e transmite cada mensagem como downlink. QoS 0/1/2, retain, user/pass ou certificado TLS.
- **LoRaWAN** — OTAA e ABP, MAC 1.0.x e 1.1, regiões, classe A/C, crypto real (MIC de join/uplink, encrypt/decrypt do payload), DevAddr/FCnt, decode de downlink — por Semtech UDP ou Basics Station, contra qualquer LNS.
- **Eventos** — payloads pré-registrados e com template (`{{deviceId}}`, `{{randInt(a,b)}}`, `{{counter}}`) com agendamento opcional, mais disparo sob demanda.
- **Engine** — um scheduler (jobs temporizados) ao lado de um gerenciador de sessões (conexões persistentes com reconexão infinita por backoff); mudanças de CRUD reconciliam na hora.
- **Console** — stream WebSocket ao vivo de todo uplink, downlink e frame de system/status, com filtros e painel de detalhe.
- **Logs** — histórico persistido, paginado e filtrável das mensagens dos dispositivos.
- **Desktop** — app Electron que sobe o sidecar Go; build cross-platform do sidecar para empacotamento.
- **i18n** — inglês e português do Brasil completos.
- **Marketplace** — navega um catálogo online de fabricantes (filtra por protocolo, tipo de leitura e fabricante) e instala um dispositivo pré-configurado num clique; a instalação reaproveita o caminho de criação de dispositivo, sem endpoint novo.

---

## Notas

- Cada [`features/NN-*.md`](./features/) documenta **o que a plataforma oferece** naquela
  área e **como é construída** (arquitetura, modelo de dados, arquivos-chave) — não como
  usar passo a passo. O tour prático, clique a clique, está em
  [`/quickTest`](../../quickTest/).
- O que não cabe nessas duas partes (casos de borda, limitações) vai numa seção
  **Notas** em cada doc.

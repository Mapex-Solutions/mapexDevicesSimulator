# Mapex Devices Simulator

> 🇺🇸 [English version](./README.md)

**Simule dispositivos IoT reais — sem precisar comprar hardware — e injete tráfego
ao vivo no [MapexOS](https://github.com/Mapex-Solutions/mapexOS).**

Um app desktop local (Electron em volta de um sidecar Go) que se comporta como
dispositivos reais na rede: envia uplinks, recebe downlinks, faz join num LNS e
roda em agenda, através de **HTTP**, **MQTT** e **LoRaWAN** (Semtech UDP e Basics
Station).

---

## A história

Toda plataforma IoT bate na mesma parede no momento em que vira realidade: para
provar que funciona, você precisa de dispositivos — muitos deles, de vários
fabricantes, protocolos e regiões.

No [MapexOS](https://github.com/Mapex-Solutions/mapexOS) essa parede era ainda
maior. O valor da plataforma não é só ingerir telemetria — é tudo que acontece
*depois*: validar payloads, convertê-los, rotear eventos e montar **regras de
negócio no workflow engine**. Exercitar tudo isso de verdade exigia uma frota de
sensores na mesa… ou uma ordem de compra, um prazo de entrega e uma pilha de
gateways LoRaWAN que você provisiona uma vez e quase nunca mais toca.

O **Mapex Devices Simulator** nasceu para derrubar essa parede. Veio direto de uma
demanda do mercado: um jeito de **simular dispositivos sem comprá-los**, para um
time subir tráfego realista em minutos. Você cria dispositivos em software — um
Milesight EM300-TH aqui, um beacon HTTP genérico ali — aponta para o seu stack
MapexOS e vê os dados reais fluírem. A partir daí dá para demonstrar a plataforma
de ponta a ponta, testar **todas as funcionalidades** em regressão e **criar e
validar regras de negócio com o
[workflow do MapexOS](https://github.com/Mapex-Solutions/mapexOS)** — tudo antes de
um único dispositivo físico chegar.

Ele fala protocolos padrão, então não está preso ao MapexOS: aponte para qualquer
LNS ou broker e funciona. Mas foi feito para o MapexOS, e é ali que ele brilha.

---

## O que dá para fazer

- **Quatro protocolos** — HTTP (só envio), MQTT (publish + subscribe) e LoRaWAN por
  Semtech UDP ou Basics Station, com crypto OTAA/ABP real.
- **Sessões ao vivo** — conexões persistentes com broker/LNS e reconexão
  automática; dispare eventos sob demanda ou numa agenda repetida.
- **Console & logs** — stream WebSocket ao vivo de todo uplink, downlink e frame de
  sistema, mais um histórico persistido e filtrável.
- **Um binário desktop** — app Electron que sobe o engine Go; nada externo para
  rodar.

O catálogo completo de funcionalidades está em
[`RELEASES/1.0.0`](./RELEASES/1.0.0/README_pt.md), e os tutoriais passo a passo em
[`quickTest`](./quickTest/README_pt.md).

---

## Movido pelo Mapex Marketplace

O **marketplace de dispositivos** embutido permite navegar um catálogo e instalar um
dispositivo pré-configurado num clique — e ainda ajustar identidade, gateway e
credenciais antes de ele entrar no seu simulador.

Esse catálogo é servido pelo
[**mapexMarketplace**](https://github.com/Mapex-Solutions/mapexMarketplace) — o
serviço da Mapex que **hospeda os dispositivos**: suas definições, **codecs de
payload**, **datasheets** e **manuais**. O simulador navega nele online; instalar um
modelo reaproveita o caminho normal de criação de dispositivo, então um device do
marketplace é só um device comum com defaults sensatos.

---

## Instalação

Pegue o instalador da sua plataforma e execute — não há mais nada para configurar, o
engine Go já vem dentro do app.

| Plataforma | Arquivo |
|------------|---------|
| Debian / Ubuntu | `.deb` |
| RedHat / Fedora | `.rpm` |
| macOS | `.dmg` |
| Windows | `.exe` |

O passo a passo por SO — e o aviso único de "desenvolvedor não identificado" (os
instaladores não são assinados) — está em
[`RELEASES/1.0.0` → Instalação](./RELEASES/1.0.0/README_pt.md#instalação).

## Build a partir do código

```bash
cd frontend
npm install
npm run build:electron      # compila o sidecar Go e empacota o app desktop
```

O app empacotado sai em `frontend/app/dist/electron/Packaged/`. Rode o build no SO
que você quer mirar (Linux → deb/rpm, macOS → dmg, Windows → exe).

---

## O ecossistema MapexOS

O simulador gera o lado do dispositivo; esses projetos são a plataforma que recebe,
roteia e age sobre ele.

- **[mapexOS](https://github.com/Mapex-Solutions/mapexOS)** — o core da plataforma:
  ingestão, conversão, roteamento e o workflow engine.
- **[mapexMarketplace](https://github.com/Mapex-Solutions/mapexMarketplace)** — o
  catálogo que hospeda dispositivos, codecs, datasheets e manuais.
- **[mapexOSDeploy](https://github.com/Mapex-Solutions/mapexOSDeploy)** — a
  distribuição executável do MapexOS (Docker Compose + imagens).
- **[mapexMQTTBroker](https://github.com/Mapex-Solutions/mapexMQTTBroker)** — o
  broker MQTT com que os dispositivos MQTT do simulador conversam.
- **[mapexLNS](https://github.com/Mapex-Solutions/mapexLNS)** — o LoRaWAN Network
  Server no qual os dispositivos LoRaWAN do simulador fazem join.

---

<sub>Parte do **MapexOS** — a plataforma aberta da Mapex Solutions para integração
de dados e automação inteligente. **Connect. Automate. Scale.**</sub>

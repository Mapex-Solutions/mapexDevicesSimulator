# Marketplace de dispositivos — catálogo mock (POC)

> 🇺🇸 English version: [README.md](./README.md)

Esta pasta é um **mock de prova de conceito** do catálogo online `mapexMarketplace`,
servido estaticamente a partir do `public/` do app em desenvolvimento. Ele permite
navegar e exercitar a UI do marketplace
([feature 14](../../../../RELEASES/1.0.0/features/14-marketplace.md)) sem o serviço de
catálogo real. Em produção o app é apontado para o serviço `mapexMarketplace`
publicado (`/api/v1/devices`) — o formato de bundle abaixo espelha as respostas desse
serviço.

O catálogo é **somente leitura**. Instalar um modelo grava um dispositivo real no
engine local pelo caminho normal de criação de dispositivo; nada aqui é mutado.

## Estrutura

```
marketplace/
├── registry.json                 # o índice completo do catálogo (cada card)
├── taxonomy.json                 # vocabulários de filtro (protocolos, tipos de leitura)
└── vendors/
    └── <vendor>/<model>/
        ├── device_information.json   # a ficha de detalhe
        ├── device_simulator.json     # o template instalável
        └── codec/                    # codec(s) de payload, opcional
            └── decoder.js
```

Mock atual: **5 modelos em 3 fabricantes** (Milesight, Dragino, Generic), 3 deles
LoRaWAN com um codec de referência.

## Arquivos

### `registry.json` — o índice do catálogo

Uma entrada por card de dispositivo; é o que a grade de listagem e os filtros leem.
Cada `item` carrega identidade (`id`, `vendor`, `vendorName`, `model`, `slug`,
`name`, `description`), classificação (`protocol`, `readingTypes`, `tags`, `icon`),
o `path` do bundle (`vendors/<vendor>/<model>`) e as flags `hasCodec` / `hasManual`
que a UI usa para decidir quais ações de detalhe mostrar.

```json
{
  "id": "milesight-em300-th",
  "vendor": "milesight",
  "vendorName": "Milesight",
  "model": "EM300-TH",
  "slug": "em300-th",
  "name": "EM300-TH Temperature & Humidity Sensor",
  "description": "Indoor LoRaWAN sensor reporting temperature and relative humidity…",
  "protocol": "lorawan",
  "readingTypes": ["temperature", "humidity"],
  "tags": ["indoor", "battery", "climate"],
  "icon": "mdi-thermometer",
  "path": "vendors/milesight/em300-th",
  "hasCodec": true,
  "hasManual": true
}
```

### `taxonomy.json` — vocabulários de filtro

As opções de filtro de `protocol` e `readingType` (valor + rótulo de exibição, mais
um ícone para tipos de leitura). É a fonte das facets do marketplace; mantenha novos
valores aqui em sincronia com o que as entradas realmente usam. As facets de
fabricante são derivadas do registry, então não ficam listadas aqui.

### `vendors/<vendor>/<model>/device_information.json` — a ficha de detalhe

O que o diálogo de detalhe renderiza: descrição completa, bloco `vendor`
(site/suporte), `tags`, `version`, `images` e os links de documentos. Documentos
podem ser um **PDF empacotado** (arquivo `datasheet` / `manual`) e/ou um **link
externo** (`datasheetUrl` / `manualUrl`) — a UI mostra "ver PDF" quando há arquivo,
senão "abrir online". Um bloco `codec` opcional aponta para o decoder de referência
do modelo.

### `vendors/<vendor>/<model>/device_simulator.json` — o template instalável

O payload do `install`: um `DeviceInput` de dispositivo **menos o `deviceId`**, que o
engine cunha na hora da instalação. É exatamente o formato que a API de criação de
dispositivo aceita — `name`, `protocolId`, `enabled`, `storeLogs`, `config`
(específico do protocolo), `attributes` e uma lista `events` de exemplo (com template,
ex. `{{randInt(180,260)}}`, `{{deviceId}}`, `{{now}}`). Após instalar, cada valor é
um campo de dispositivo comum que o usuário pode editar.

### `vendors/<vendor>/<model>/codec/` — codec de payload (opcional)

O decoder (e/ou encoder) de referência do fabricante. É exposto para o usuário
**ver e baixar**; o simulador ainda não o executa.

## Adicionando um modelo

1. Crie `vendors/<vendor>/<model>/` com `device_information.json` e
   `device_simulator.json` (e `codec/` se houver).
2. Adicione um card correspondente em `registry.json` (`path` = a nova pasta, defina
   `hasCodec` / `hasManual`).
3. Se introduzir um novo protocolo ou tipo de leitura, adicione-o em `taxonomy.json`
   para os filtros poderem oferecê-lo.

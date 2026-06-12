# Quick test de LoRaWAN — Basics Station

> 🇺🇸 English version: [README.md](./README.md)

Mesmo sensor LoRaWAN e mesmo join OTAA do teste UDP, mas o dispositivo carrega o
**próprio link WebSocket Basics Station** até o LNS — sem uma entidade gateway separada.
O simulador conecta em `ws://host:3001/gw/<gateway-eui>` (o gateway-bridge Basics Station
do ChirpStack) e envia o uplink real do Dragino LHT65N.

> Precisa de um LNS rodando com Basics Station habilitado. O `chirpstack-docker` do
> ChirpStack o expõe na `:3001`. Veja a seção 4 do [`../README_pt.md`](../README_pt.md).

---

## Registre no ChirpStack (para as chaves baterem)

- **Gateway** com EUI `0102030405060708` (o mesmo EUI que o dispositivo anuncia).
- **Device profile**: região `EU868`, MAC version `LoRaWAN 1.0.3`, OTAA habilitado.
- **Device**: DevEUI `0011223344556677`, JoinEUI `0000000000000000`.
- **Application key (OTAA)**: `00112233445566778899AABBCCDDEEFF`.

(A mesma observação de reset do DevNonce da pasta UDP vale ao rodar de novo.)

---

## Crie o dispositivo (UI)

**Devices → New device**

| Etapa | Campo | Cole isto |
|-------|-------|-----------|
| Info | Name | `Quick LoRa Basics` |
| Info | Device ID | `lora-bs-01` |
| Info | Protocol | `Basics Station` |
| Connection | LNS URI | `ws://127.0.0.1:3001` |
| Connection | Gateway EUI | `0102030405060708` |
| Connection | Region | `EU868` |
| Connection | MAC version | `1.0.3` |
| Connection | Activation | `OTAA` |
| Connection | DevEUI | `0011223344556677` |
| Connection | JoinEUI | `0000000000000000` |
| Connection | AppKey | `00112233445566778899AABBCCDDEEFF` |

![Conexão Basics Station](./images/01-connection.png)

> O simulador acrescenta `/gw/<gateway-eui>` à LNS URI automaticamente — informe apenas
> `ws://127.0.0.1:3001`.

### Adicione um evento — um uplink real do LHT65N

| Campo | Cole isto |
|-------|-----------|
| Name | `LHT65N uplink` |
| FPort | `2` |
| Confirmed | off |
| Payload (hex) | `0BB809F6025D0000000000` |

![Evento Basics Station](./images/02-event.png)

---

## Rode

1. **Save** e ligue o **Enabled**.
2. Abra o **Console** — o WebSocket conecta, depois `join-request → join-accept →
   joined`, e o **Fire event** produz um frame `up`.
3. O ChirpStack mostra o uplink (decodificado, se você tiver adicionado o codec do
   LHT65N).

![Console Basics Station](./images/03-console.png)

---

## Alternativa em um comando (API)

```bash
bash quickTest/lorawan-basic-station/curl.sh   # padrão http://127.0.0.1:5055
```

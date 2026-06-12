# Quick test de LoRaWAN — Semtech UDP

> 🇺🇸 English version: [README.md](./README.md)

Um sensor LoRaWAN trafega num **gateway** que carrega o link até o seu LNS. Aqui o
gateway fala o clássico **packet-forwarder Semtech UDP** com o gateway-bridge do
ChirpStack na `:1700`. O sensor faz um join OTAA e então envia um uplink real de
**Dragino LHT65N** que o LNS decodifica para temperatura/umidade/bateria.

> Precisa de um LNS rodando. Veja a seção 4 do [`../README_pt.md`](../README_pt.md) para
> subir um ChirpStack em uma linha com `chirpstack-docker`.

---

## Registre no ChirpStack (para as chaves baterem)

- **Gateway** com EUI `0102030405060708`.
- **Device profile**: região `EU868`, MAC version `LoRaWAN 1.0.3`, OTAA habilitado.
  (Opcional: cole o codec JS do Dragino LHT65N para decodificar o payload.)
- **Device**: DevEUI `0011223344556677`, JoinEUI `0000000000000000`.
- **Application key (OTAA)**: `00112233445566778899AABBCCDDEEFF`.

> O re-join é rejeitado se o DevNonce se repetir. Para rodar de novo do zero, limpe-o:
> ```bash
> docker compose exec postgres psql -U chirpstack -d chirpstack \
>   -c "update device_keys set dev_nonces='{}'::jsonb where dev_eui=decode('0011223344556677','hex');"
> ```

---

## Crie o gateway (UI)

**Gateways → New gateway**

| Campo | Cole isto |
|-------|-----------|
| Name | `Quick UDP gateway` |
| EUI | `0102030405060708` |
| Region | `EU868` |
| Link protocol | `Semtech UDP` |
| Host | `127.0.0.1` |
| Port | `1700` |

![Novo gateway](./images/01-gateway.png)

Ligue-o — o console mostra o gateway ficando **online** (ele envia stats periódicas para
o LNS marcá-lo como visto).

## Crie o dispositivo (UI)

**Devices → New device**

| Etapa | Campo | Cole isto |
|-------|-------|-----------|
| Info | Name | `Quick LoRa UDP` |
| Info | Device ID | `lora-udp-01` |
| Info | Protocol | `LoRaWAN` |
| Connection | Gateway | `Quick UDP gateway` |
| Connection | Region | `EU868` |
| Connection | MAC version | `1.0.3` |
| Connection | Activation | `OTAA` |
| Connection | DevEUI | `0011223344556677` |
| Connection | JoinEUI | `0000000000000000` |
| Connection | AppKey | `00112233445566778899AABBCCDDEEFF` |

![Conexão LoRaWAN](./images/02-connection.png)

### Adicione um evento — um uplink real do LHT65N

| Campo | Cole isto |
|-------|-----------|
| Name | `LHT65N uplink` |
| FPort | `2` |
| Confirmed | off |
| Payload (hex) | `0BB809F6025D0000000000` |

Esse payload decodifica para `BatV 3.0, TempC_SHT 25.5, Hum_SHT 60.5`.

![Evento LoRaWAN](./images/03-event.png)

---

## Rode

1. **Save** e ligue o **Enabled**.
2. Abra o **Console** — acompanhe `join-request → join-accept → joined` e, no
   **Fire event**, um frame `up` com `FCnt`.
3. No ChirpStack, o dispositivo mostra o uplink (e o objeto decodificado, se você tiver
   adicionado o codec).

![Console LoRaWAN](./images/04-console.png)

---

## Alternativa em um comando (API)

Cria o gateway, depois o dispositivo referenciando-o e dispara uma vez:

```bash
bash quickTest/lorawan-udp/curl.sh    # padrão http://127.0.0.1:5055
```

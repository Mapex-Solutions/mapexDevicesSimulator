# Quick test de MQTT

> 🇺🇸 English version: [README.md](./README.md)

Um dispositivo MQTT mantém uma **conexão viva** com um broker: publica uplinks e, quando
o **Receive** está ligado, permanece inscrito em tópicos e transmite cada mensagem
recebida ao console como um frame `down`.

O broker padrão abaixo é o público `broker.hivemq.com` — sem cadastro. Use o seu próprio
broker quando quiser.

---

## Crie o dispositivo (UI)

**Devices → New device**

| Etapa | Campo | Cole isto |
|-------|-------|-----------|
| Info | Name | `Quick MQTT sensor` |
| Info | Device ID | `mqtt-quick-01` |
| Info | Protocol | `MQTT` |
| Connection | Broker URL | `tcp://broker.hivemq.com:1883` |
| Connection | Client ID | `mqtt-quick-01` |
| Connection | Base topic | `mapex/quicktest` |
| Connection | Auth | `None` |
| Connection | **Receive events** | ligue **on** |
| Connection | Subscription | Name `commands` · Topic `mqtt-quick-01/cmd` · QoS `1` |

> Os tópicos são **relativos ao base topic** — a engine prefixa `mapex/quicktest`,
> então o dispositivo na prática se inscreve em `mapex/quicktest/mqtt-quick-01/cmd`.

![Conexão MQTT + subscription](./images/01-connection.png)

### Adicione um evento

**Events → Add event**

| Campo | Cole isto |
|-------|-----------|
| Name | `Telemetry` |
| Topic | `mqtt-quick-01/telemetry` (publicado em `mapex/quicktest/mqtt-quick-01/telemetry`) |
| QoS | `1` |
| Retain | off |
| Body mode | `Raw` |
| Body | veja abaixo |

```json
{
  "deviceId": "{{deviceId}}",
  "level": {{randInt(0,100)}}
}
```

![Evento MQTT](./images/02-event.png)

---

## Rode — uplink

1. **Save** e ligue o **Enabled**.
2. Abra o **Console** — você verá `connecting → connected → subscribed`.
3. **Fire event** → um frame `up` publica no tópico de telemetria.

![Console MQTT](./images/03-console.png)

## Rode — downlink (a parte divertida)

Publique no tópico inscrito a partir de qualquer cliente externo e veja um frame `down`
chegar ao vivo no console:

```bash
# usando mosquitto-clients
mosquitto_pub -h broker.hivemq.com -t 'mapex/quicktest/mqtt-quick-01/cmd' \
  -q 1 -m '{"cmd":"set-interval","seconds":30}'
```

Sem mosquitto? Use o web client do HiveMQ (http://www.hivemq.com/demos/websocket-client/),
conecte em `broker.hivemq.com` e publique no mesmo tópico. A mensagem recebida aparece
como um frame `down` no console, ao lado dos `up` publicados.

---

## Alternativa em um comando (API)

```bash
bash quickTest/mqtt/curl.sh           # padrão http://127.0.0.1:5055
```

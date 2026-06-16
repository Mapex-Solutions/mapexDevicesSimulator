# Jornada: mqtt_downlink

> 🇺🇸 English version: [README.md](./README.md)

Exercita o caminho **de entrada** do MQTT de ponta a ponta contra o simulador
vivo e um broker no processo: um device com recebimento ligado assina um tópico,
um publish externo chega nele, e o simulador surge isso como um downlink nos logs.

## Fluxo

1. **StartMQTTBroker** — sobe o broker no processo e publica as coordenadas no bag.
2. **CreateMQTTReceiveDevice** — POST `/api/devices` de um device MQTT habilitado
   com `receiveEnabled` e uma subscription num tópico de downlink por execução,
   `storeLogs` ligado. Habilitar abre uma sessão que assina.
3. **PublishDownlink** — injeta uma mensagem **retained** nesse tópico, como um
   terceiro externo faria.
4. **AssertDownlinkLogged** — faz polling em GET `/api/logs` até aparecer um frame
   `down`/`downlink` carregando o payload publicado (que embute o run id).
5. **Compensação** — apaga o device, fecha o broker.

O publish é **retained**, então é entregue mesmo que chegue antes ou depois da
subscription do device ficar ativa — a sessão assina de forma assíncrona, e uma
mensagem retained remove essa corrida sem precisar de settle.

## O que prova

A metade de entrada do engine: uma subscription é aberta no enable, uma mensagem
recebida é decodificada, e surge como frame `down`/`downlink` e é persistida nos
logs — todo o trilho "o dado chega ao device e é registrado".

## Rodar

```bash
# de e2e_tests/  (o sidecar precisa estar no ar em 127.0.0.1:5055)
go test -tags=saga ./journey/iot/mqtt_downlink/ -v
```

Auto-contido: o broker sobe no próprio processo, então só o sidecar precisa estar
no ar.

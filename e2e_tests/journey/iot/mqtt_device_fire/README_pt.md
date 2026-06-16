# Jornada: mqtt_device_fire

> 🇺🇸 English version: [README.md](./README.md)

Exercita devices MQTT de ponta a ponta contra o simulador vivo e um broker no
próprio processo, cobrindo **os dois** modos de autenticação que a plataforma
suporta.

## Fluxo

Uma saga ordenada sobe um único broker (com os dois listeners) e roda cada modo
de auth em sequência. O oráculo é o próprio broker, não a API de logs.

1. **StartMQTTBroker** — sobe o broker no processo e publica suas coordenadas e
   uma CA / certificado de cliente novos no bag.
2. **usuário/senha** — `CreateMQTTUserPassDevice` (device no listener `tcp://`,
   autenticando com usuário e senha) → `FireTelemetry` → `AssertMQTTPublished` (o
   broker aceitou um publish carregando o `deviceId`; ele só registra um publish
   depois do CONNECT passar pela auth, então isso prova que as credenciais foram
   honradas).
3. **certificado** — `CreateMQTTTLSDevice` (device no listener `ssl://`,
   autenticando com **certificado de cliente**, TLS mútuo) → `FireTelemetry` →
   `AssertMQTTPublished`. O broker usa `RequireAndVerifyClientCert`, então o
   publish só chega se o handshake validou o certificado do device contra a CA da
   execução.
4. **Compensação** — apaga os dois devices, fecha o broker.

O fire acontece no instante em que o device é habilitado — antes da sessão
persistente conectar — então exercita de propósito a fallback one-shot do engine.
Esse caminho agora usa um client id distinto, então não colide mais com a sessão
conectando e o publish chega de forma confiável (sem passo de settle).

## O que prova

CRUD de device MQTT + o caminho de fire do engine + **os dois** modos de auth
MQTT (usuário/senha e certificado de cliente) chegando a um broker real e sendo
aceitos. Um CONNECT rejeitado (senha errada, certificado não confiável) aparece
como o publish nunca chegando.

## Como o fixture do broker funciona

`common/utils/mqtt_broker.go` sobe um broker [mochi-mqtt](https://github.com/mochi-mqtt/server)
embarcado com dois listeners em portas loopback aleatórias:

- um listener TCP em texto puro que autentica por usuário/senha;
- um listener TLS que exige e verifica um certificado de cliente.

`common/utils/certs.go` gera uma CA nova por execução que assina tanto o
certificado de servidor do broker quanto o certificado de cliente do device,
então nenhum fixture de certificado fica no repositório. Todo publish aceito é
capturado para o assert.

## Rodar

```bash
# de e2e_tests/  (o sidecar precisa estar no ar em 127.0.0.1:5055)
go test -tags=saga ./journey/iot/mqtt_device_fire/ -v
```

Auto-contido: o broker e todos os certificados sobem no próprio processo, então
só o sidecar precisa estar no ar. O sidecar precisa incluir o suporte a MQTT TLS
(o engine monta o `tls.Config` a partir do material PEM do device).

## Nota sobre a fallback one-shot

Esta jornada dispara antes da sessão do device conectar de propósito: cobre o
caminho de fallback one-shot do engine. Esse caminho dá a cada conexão one-shot um
client id distinto (`<clientId>-oneshot-N`) para nunca colidir com a sessão
persistente ainda conectando — sem isso, o broker derrubaria uma das duas conexões
de mesmo id e o fire seria perdido em silêncio na janela de connect.

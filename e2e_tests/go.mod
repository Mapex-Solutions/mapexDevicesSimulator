module github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests

go 1.25.3

require (
	github.com/Mapex-Solutions/mapexGoKit/infrastructure v0.0.0
	github.com/Mapex-Solutions/mapexGoKit/utils v0.0.0
	github.com/chirpstack/chirpstack/api/go/v4 v4.18.0
	github.com/gorilla/websocket v1.5.3
	github.com/mochi-mqtt/server/v2 v2.7.9
	google.golang.org/grpc v1.81.1
	simulator v0.0.0
)

require (
	github.com/rs/xid v1.6.0 // indirect
	golang.org/x/net v0.53.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260427160629-7cedc36a6bc4 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260427160629-7cedc36a6bc4 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace simulator => ../backend

replace github.com/Mapex-Solutions/mapexGoKit/microservices => ../../mapexGoKit/microservices

replace github.com/Mapex-Solutions/mapexGoKit/infrastructure => ../../mapexGoKit/infrastructure

replace github.com/Mapex-Solutions/mapexGoKit/utils => ../../mapexGoKit/utils

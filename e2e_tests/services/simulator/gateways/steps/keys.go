// Package steps holds saga steps that exercise the simulator's gateways module
// (POST/DELETE /api/gateways).
package steps

// BagKeyUDPSimGatewayID is the server id (UUID) of the simulator's UDP gateway,
// read by the LoRaWAN device step to attach the device to it.
const BagKeyUDPSimGatewayID = "gateways.udpSimId"

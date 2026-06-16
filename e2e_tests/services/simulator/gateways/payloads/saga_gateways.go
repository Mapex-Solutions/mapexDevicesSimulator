// Package payloads holds canonical gateway-create fixtures for the simulator's
// gateways module. Builders are pure; the step supplies the runtime values.
package payloads

import (
	"encoding/json"
	"fmt"

	gatewayscontract "simulator/packages/contracts/gateways"
)

// SagaUDPGateway builds an enabled EU868 gateway whose Semtech UDP link targets
// the ChirpStack gateway bridge at host:port. The EUI must match the gateway
// registered on the LNS.
func SagaUDPGateway(runID, eui, host string, port int) gatewayscontract.GatewayInput {
	link, _ := json.Marshal(map[string]any{
		"protocol": "udp",
		"lnsUri":   "",
		"host":     host,
		"port":     port,
	})
	return gatewayscontract.GatewayInput{
		Name:        fmt.Sprintf("e2e UDP gateway %s", runID),
		EUI:         eui,
		Enabled:     true,
		Region:      "EU868",
		Description: "e2e Semtech UDP gateway",
		Link:        link,
	}
}

package steps

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	gatewayscontract "simulator/packages/contracts/gateways"

	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/common/constants"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/common/types"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/core/saga"
	provSteps "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/chirpstack/provisioning/steps"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/gateways/payloads"
)

// CreateUDPGateway creates the simulator gateway that forwards over Semtech UDP
// to the ChirpStack bridge, using the EUI registered on the LNS so frames are
// accepted. It publishes the gateway id the LoRaWAN device attaches to.
//
// Reads (bag): provSteps.BagKeyUDPGatewayEUI.
// Writes (bag): BagKeyUDPSimGatewayID.
// Compensate: DELETE the gateway.
func CreateUDPGateway() saga.Step {
	return saga.Step{
		Name: "gateways.CreateUDPGateway",
		Do: func(c *saga.Context) error {
			eui := c.MustGetString(provSteps.BagKeyUDPGatewayEUI)
			spec := payloads.SagaUDPGateway(c.RunID, eui, constants.ChirpStackUDPHost, constants.ChirpStackUDPPort)
			resp, err := c.Clients.Sim.Raw(c.Stdctx, http.MethodPost, "/api/gateways", spec)
			if err != nil {
				return fmt.Errorf("create gateway: %w", err)
			}
			defer resp.Body.Close()
			if resp.StatusCode < 200 || resp.StatusCode >= 300 {
				return fmt.Errorf("create gateway: unexpected status %d", resp.StatusCode)
			}
			var env types.Envelope
			if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
				return fmt.Errorf("decode create gateway: %w", err)
			}
			var gw gatewayscontract.Gateway
			if err := json.Unmarshal(env.Data, &gw); err != nil {
				return fmt.Errorf("decode gateway data: %w", err)
			}
			if gw.ID == "" {
				return fmt.Errorf("create gateway: empty id in response")
			}
			c.Set(BagKeyUDPSimGatewayID, gw.ID)
			return nil
		},
		Compensate: func(c *saga.Context) error {
			id, ok := c.Get(BagKeyUDPSimGatewayID)
			if !ok {
				return nil
			}
			resp, err := c.Clients.Sim.Raw(context.Background(), http.MethodDelete, "/api/gateways/"+id.(string), nil)
			if err != nil {
				return err
			}
			return resp.Body.Close()
		},
	}
}

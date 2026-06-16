package e2e

import (
	"testing"

	"github.com/Mapex-Solutions/mapexGoKit/utils/random"

	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/common/constants"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/gateways/payloads"
)

// TestList_200 confirms the gateways list endpoint answers.
func TestList_200(t *testing.T) {
	_ = listGateways(t) // listGateways fails the test on a non-200.
}

// TestList_ContainsCreated creates several gateways and finds them all in the
// list. The endpoint is unpaginated, so this is the membership guarantee.
func TestList_ContainsCreated(t *testing.T) {
	created := make(map[string]bool, listFixtureCount)
	for i := 0; i < listFixtureCount; i++ {
		runID := random.NewRunID()
		spec := payloads.SagaUDPGateway(runID, euiFromRunID(runID), constants.ChirpStackUDPHost, constants.ChirpStackUDPPort)
		created[createGateway(t, spec)] = true
	}
	found := 0
	for _, g := range listGateways(t) {
		if created[g.ID] {
			found++
		}
	}
	if found != listFixtureCount {
		t.Fatalf("found %d/%d created gateways in the list", found, listFixtureCount)
	}
}

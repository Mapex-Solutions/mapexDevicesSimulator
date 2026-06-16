package e2e

import (
	"testing"

	"github.com/Mapex-Solutions/mapexGoKit/utils/random"

	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/devices/payloads"
)

// TestList_200 confirms the devices list endpoint answers.
func TestList_200(t *testing.T) {
	_ = listDevices(t) // listDevices fails the test on a non-200.
}

// TestList_ContainsCreated creates several devices and finds them all in the
// list. The endpoint is unpaginated, so this is the membership guarantee.
func TestList_ContainsCreated(t *testing.T) {
	created := make(map[string]bool, listFixtureCount)
	for i := 0; i < listFixtureCount; i++ {
		runID := random.NewRunID()
		spec := payloads.SagaHTTPDevice(runID, "http://example.invalid").Build()
		created[createDevice(t, spec)] = true
	}
	found := 0
	for _, d := range listDevices(t) {
		if created[d.ID] {
			found++
		}
	}
	if found != listFixtureCount {
		t.Fatalf("found %d/%d created devices in the list", found, listFixtureCount)
	}
}

package e2e

import (
	"net/http"
	"testing"

	"github.com/Mapex-Solutions/mapexGoKit/utils/random"

	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/common/constants"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/gateways/payloads"
)

// TestCreate_201 creates a gateway and gets it back with an id.
func TestCreate_201(t *testing.T) {
	runID := random.NewRunID()
	spec := payloads.SagaUDPGateway(runID, euiFromRunID(runID), constants.ChirpStackUDPHost, constants.ChirpStackUDPPort)
	_ = createGateway(t, spec)
}

// TestUpdate_200 updates a gateway's name and sees it reflected.
func TestUpdate_200(t *testing.T) {
	runID := random.NewRunID()
	spec := payloads.SagaUDPGateway(runID, euiFromRunID(runID), constants.ChirpStackUDPHost, constants.ChirpStackUDPPort)
	id := createGateway(t, spec)

	spec.Name = "renamed-" + runID
	if code := statusOf(t, http.MethodPut, "/api/gateways/"+id, spec); code != http.StatusOK {
		t.Fatalf("update: status %d, want 200", code)
	}
	for _, g := range listGateways(t) {
		if g.ID == id && g.Name == "renamed-"+runID {
			return
		}
	}
	t.Fatalf("updated name not reflected in list")
}

// TestDelete_200 deletes a gateway and it disappears from the list.
func TestDelete_200(t *testing.T) {
	runID := random.NewRunID()
	spec := payloads.SagaUDPGateway(runID, euiFromRunID(runID), constants.ChirpStackUDPHost, constants.ChirpStackUDPPort)
	id := createGateway(t, spec)

	if code := statusOf(t, http.MethodDelete, "/api/gateways/"+id, nil); code != http.StatusOK {
		t.Fatalf("delete: status %d, want 200", code)
	}
	for _, g := range listGateways(t) {
		if g.ID == id {
			t.Fatalf("gateway still present after delete")
		}
	}
}

// TestCreate_EmptyName_400 rejects a gateway with no name.
func TestCreate_EmptyName_400(t *testing.T) {
	runID := random.NewRunID()
	spec := payloads.SagaUDPGateway(runID, euiFromRunID(runID), constants.ChirpStackUDPHost, constants.ChirpStackUDPPort)
	spec.Name = ""
	if code := statusOf(t, http.MethodPost, "/api/gateways", spec); code != http.StatusBadRequest {
		t.Fatalf("create empty name: status %d, want 400", code)
	}
}

// TestCreate_InvalidRegion_400 rejects an unknown region.
func TestCreate_InvalidRegion_400(t *testing.T) {
	runID := random.NewRunID()
	spec := payloads.SagaUDPGateway(runID, euiFromRunID(runID), constants.ChirpStackUDPHost, constants.ChirpStackUDPPort)
	spec.Region = "ZZ999"
	if code := statusOf(t, http.MethodPost, "/api/gateways", spec); code != http.StatusBadRequest {
		t.Fatalf("create invalid region: status %d, want 400", code)
	}
}

// TestDelete_404 deletes a gateway that does not exist.
func TestDelete_404(t *testing.T) {
	if code := statusOf(t, http.MethodDelete, "/api/gateways/"+nonExistentID, nil); code != http.StatusNotFound {
		t.Fatalf("delete missing: status %d, want 404", code)
	}
}

// TestUpdate_404 updates a gateway that does not exist.
func TestUpdate_404(t *testing.T) {
	runID := random.NewRunID()
	spec := payloads.SagaUDPGateway(runID, euiFromRunID(runID), constants.ChirpStackUDPHost, constants.ChirpStackUDPPort)
	if code := statusOf(t, http.MethodPut, "/api/gateways/"+nonExistentID, spec); code != http.StatusNotFound {
		t.Fatalf("update missing: status %d, want 404", code)
	}
}

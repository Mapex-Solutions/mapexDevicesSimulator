package e2e

import (
	"net/http"
	"testing"

	"github.com/Mapex-Solutions/mapexGoKit/utils/random"

	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/devices/payloads"
)

// The device specs use a dummy target URL: these tests never fire, they
// exercise persistence only.

// TestCreate_201 creates a device and gets it back with an id.
func TestCreate_201(t *testing.T) {
	runID := random.NewRunID()
	spec := payloads.SagaHTTPDevice(runID, "http://example.invalid").Build()
	_ = createDevice(t, spec)
}

// TestUpdate_200 updates a device's name and sees it reflected.
func TestUpdate_200(t *testing.T) {
	runID := random.NewRunID()
	spec := payloads.SagaHTTPDevice(runID, "http://example.invalid").Build()
	id := createDevice(t, spec)

	spec.Name = "renamed-" + runID
	if code := statusOf(t, http.MethodPut, "/api/devices/"+id, spec); code != http.StatusOK {
		t.Fatalf("update: status %d, want 200", code)
	}
	for _, d := range listDevices(t) {
		if d.ID == id && d.Name == "renamed-"+runID {
			return
		}
	}
	t.Fatalf("updated name not reflected in list")
}

// TestDelete_200 deletes a device and it disappears from the list.
func TestDelete_200(t *testing.T) {
	runID := random.NewRunID()
	spec := payloads.SagaHTTPDevice(runID, "http://example.invalid").Build()
	id := createDevice(t, spec)

	if code := statusOf(t, http.MethodDelete, "/api/devices/"+id, nil); code != http.StatusOK {
		t.Fatalf("delete: status %d, want 200", code)
	}
	for _, d := range listDevices(t) {
		if d.ID == id {
			t.Fatalf("device still present after delete")
		}
	}
}

// TestCreate_EmptyName_400 rejects a device with no name.
func TestCreate_EmptyName_400(t *testing.T) {
	runID := random.NewRunID()
	spec := payloads.SagaHTTPDevice(runID, "http://example.invalid").Build()
	spec.Name = ""
	if code := statusOf(t, http.MethodPost, "/api/devices", spec); code != http.StatusBadRequest {
		t.Fatalf("create empty name: status %d, want 400", code)
	}
}

// TestCreate_InvalidProtocol_400 rejects an unknown protocol.
func TestCreate_InvalidProtocol_400(t *testing.T) {
	runID := random.NewRunID()
	spec := payloads.SagaHTTPDevice(runID, "http://example.invalid").Build()
	spec.ProtocolID = "ftp"
	if code := statusOf(t, http.MethodPost, "/api/devices", spec); code != http.StatusBadRequest {
		t.Fatalf("create invalid protocol: status %d, want 400", code)
	}
}

// TestDelete_404 deletes a device that does not exist.
func TestDelete_404(t *testing.T) {
	if code := statusOf(t, http.MethodDelete, "/api/devices/"+nonExistentID, nil); code != http.StatusNotFound {
		t.Fatalf("delete missing: status %d, want 404", code)
	}
}

// TestUpdate_404 updates a device that does not exist.
func TestUpdate_404(t *testing.T) {
	runID := random.NewRunID()
	spec := payloads.SagaHTTPDevice(runID, "http://example.invalid").Build()
	if code := statusOf(t, http.MethodPut, "/api/devices/"+nonExistentID, spec); code != http.StatusNotFound {
		t.Fatalf("update missing: status %d, want 404", code)
	}
}

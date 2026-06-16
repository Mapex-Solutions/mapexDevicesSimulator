// Package e2e holds the HTTP contract tests for the devices module: create,
// update, delete (with validation and 404 paths) and the list membership. The
// devices list is unpaginated and unfiltered, so there is no pagination/filter
// or cursor coverage here.
package e2e

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/httpclient"

	devicescontract "simulator/packages/contracts/devices"

	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/common/constants"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/common/types"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/common/utils"
)

// client is the package-scoped simulator client (no auth).
var client *httpclient.HTTPClient

// TestMain skips the package cleanly when the sidecar is down, else builds the
// client.
func TestMain(m *testing.M) {
	if err := utils.SetupE2EEnvironment(); err != nil {
		os.Exit(0)
	}
	client = httpclient.New(httpclient.Config{BaseURL: constants.SimURL})
	os.Exit(m.Run())
}

// createDevice POSTs the device, asserts 201, and returns its server id,
// registering a cleanup that deletes it.
func createDevice(t *testing.T, spec devicescontract.DeviceInput) string {
	t.Helper()
	resp, err := client.Raw(context.Background(), http.MethodPost, "/api/devices", spec)
	if err != nil {
		t.Fatalf("create device: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create device: status %d, want 201", resp.StatusCode)
	}
	var dev devicescontract.Device
	decodeData(t, resp, &dev)
	if dev.ID == "" {
		t.Fatalf("create device: empty id")
	}
	t.Cleanup(func() { deleteDevice(dev.ID) })
	return dev.ID
}

// deleteDevice removes the device, tolerating an already-gone one.
func deleteDevice(id string) {
	resp, err := client.Raw(context.Background(), http.MethodDelete, "/api/devices/"+id, nil)
	if err == nil {
		resp.Body.Close()
	}
}

// statusOf performs the request and returns only the status code, for the
// negative paths.
func statusOf(t *testing.T, method, path string, body any) int {
	t.Helper()
	resp, err := client.Raw(context.Background(), method, path, body)
	if err != nil {
		t.Fatalf("%s %s: %v", method, path, err)
	}
	defer resp.Body.Close()
	return resp.StatusCode
}

// listDevices returns the current device list.
func listDevices(t *testing.T) []devicescontract.Device {
	t.Helper()
	resp, err := client.Raw(context.Background(), http.MethodGet, "/api/devices", nil)
	if err != nil {
		t.Fatalf("list devices: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("list devices: status %d", resp.StatusCode)
	}
	var devs []devicescontract.Device
	decodeData(t, resp, &devs)
	return devs
}

// decodeData unmarshals the envelope's data into out.
func decodeData(t *testing.T, resp *http.Response, out any) {
	t.Helper()
	var env types.Envelope
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		t.Fatalf("decode envelope: %v", err)
	}
	if err := json.Unmarshal(env.Data, out); err != nil {
		t.Fatalf("decode data: %v", err)
	}
}

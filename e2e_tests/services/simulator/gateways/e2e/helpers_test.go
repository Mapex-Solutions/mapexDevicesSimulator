// Package e2e holds the HTTP contract tests for the gateways module: create,
// update, delete (with validation and 404 paths) and the list membership. The
// gateways list is unpaginated and unfiltered, so there is no pagination/filter
// or cursor coverage here.
package e2e

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/httpclient"

	gatewayscontract "simulator/packages/contracts/gateways"

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

// euiFromRunID derives a unique 16-hex gateway EUI from the run id, so each
// created gateway has a distinct identifier.
func euiFromRunID(runID string) string {
	sum := sha256.Sum256([]byte(runID))
	return hex.EncodeToString(sum[:8])
}

// createGateway POSTs the gateway, asserts 201, and returns its server id,
// registering a cleanup that deletes it.
func createGateway(t *testing.T, spec gatewayscontract.GatewayInput) string {
	t.Helper()
	resp, err := client.Raw(context.Background(), http.MethodPost, "/api/gateways", spec)
	if err != nil {
		t.Fatalf("create gateway: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create gateway: status %d, want 201", resp.StatusCode)
	}
	var gw gatewayscontract.Gateway
	decodeData(t, resp, &gw)
	if gw.ID == "" {
		t.Fatalf("create gateway: empty id")
	}
	t.Cleanup(func() { deleteGateway(gw.ID) })
	return gw.ID
}

// deleteGateway removes the gateway, tolerating an already-gone one.
func deleteGateway(id string) {
	resp, err := client.Raw(context.Background(), http.MethodDelete, "/api/gateways/"+id, nil)
	if err == nil {
		resp.Body.Close()
	}
}

// statusOf performs the request and returns only the status code.
func statusOf(t *testing.T, method, path string, body any) int {
	t.Helper()
	resp, err := client.Raw(context.Background(), method, path, body)
	if err != nil {
		t.Fatalf("%s %s: %v", method, path, err)
	}
	defer resp.Body.Close()
	return resp.StatusCode
}

// listGateways returns the current gateway list.
func listGateways(t *testing.T) []gatewayscontract.Gateway {
	t.Helper()
	resp, err := client.Raw(context.Background(), http.MethodGet, "/api/gateways", nil)
	if err != nil {
		t.Fatalf("list gateways: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("list gateways: status %d", resp.StatusCode)
	}
	var gws []gatewayscontract.Gateway
	decodeData(t, resp, &gws)
	return gws
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

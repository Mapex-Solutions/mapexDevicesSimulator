// Package e2e holds the HTTP contract tests for the logs module: the cursor
// pagination and the filters on GET /api/logs. Logs are read-only — they are
// produced by firing device events — so the fixtures seed logs by creating an
// HTTP device (targeting an in-test echo) and firing it, then the tests exercise
// the list surface.
package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/Mapex-Solutions/mapexGoKit/infrastructure/httpclient"
	"github.com/Mapex-Solutions/mapexGoKit/utils/random"

	devicescontract "simulator/packages/contracts/devices"
	logscontract "simulator/packages/contracts/logs"

	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/common/constants"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/common/types"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/common/utils"
	devicePayloads "github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/services/simulator/devices/payloads"
)

// client is the package-scoped simulator client. The simulator has no auth, so
// it is just a base-URL client built once.
var (
	client *httpclient.HTTPClient
	echo   *utils.Echo
)

// TestMain checks the sidecar is up (skipping the package cleanly when it is
// not), boots the shared echo target the seeded devices fire against, and builds
// the package client.
func TestMain(m *testing.M) {
	if err := utils.SetupE2EEnvironment(); err != nil {
		fmt.Printf("logs e2e skipped: simulator not ready: %v\n", err)
		os.Exit(0)
	}
	echo = utils.StartEcho()
	client = httpclient.New(httpclient.Config{BaseURL: constants.SimURL})
	code := m.Run()
	echo.Close()
	os.Exit(code)
}

// seedLogs creates an enabled HTTP device firing at the echo and fires it count
// times, returning the device's user-facing deviceId once that many logs are
// visible. The device is removed on cleanup.
func seedLogs(t *testing.T, count int) string {
	t.Helper()
	runID := random.NewRunID()
	spec := devicePayloads.SagaHTTPDevice(runID, echo.URL()).Build()
	serverID, deviceID := createDevice(t, spec)
	for i := 0; i < count; i++ {
		fireEvent(t, serverID, "e1")
	}
	waitForLogCount(t, deviceID, count)
	return deviceID
}

// createDevice POSTs the device and returns its server id and user-facing
// deviceId, registering a cleanup that deletes it.
func createDevice(t *testing.T, spec devicescontract.DeviceInput) (serverID, deviceID string) {
	t.Helper()
	resp, err := client.Raw(context.Background(), http.MethodPost, "/api/devices", spec)
	if err != nil {
		t.Fatalf("create device: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		t.Fatalf("create device: status %d", resp.StatusCode)
	}
	var dev devicescontract.Device
	decodeData(t, resp, &dev)
	if dev.ID == "" {
		t.Fatalf("create device: empty id")
	}
	t.Cleanup(func() {
		r, err := client.Raw(context.Background(), http.MethodDelete, "/api/devices/"+dev.ID, nil)
		if err == nil {
			r.Body.Close()
		}
	})
	return dev.ID, dev.DeviceID
}

// fireEvent fires a pre-registered event on the device.
func fireEvent(t *testing.T, serverID, eventID string) {
	t.Helper()
	resp, err := client.Raw(context.Background(), http.MethodPost, "/api/devices/"+serverID+"/fire", map[string]string{"eventId": eventID})
	if err != nil {
		t.Fatalf("fire: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		t.Fatalf("fire: status %d", resp.StatusCode)
	}
}

// waitForLogCount polls until at least count logs exist for the device.
func waitForLogCount(t *testing.T, deviceID string, count int) {
	t.Helper()
	deadline := time.Now().Add(15 * time.Second)
	for {
		q := url.Values{}
		q.Set("device", deviceID)
		q.Set("limit", fmt.Sprintf("%d", count+5))
		page := listLogs(t, q)
		if len(page.Items) >= count {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("device %s only has %d/%d logs after timeout", deviceID, len(page.Items), count)
		}
		time.Sleep(200 * time.Millisecond)
	}
}

// listLogs performs one GET /api/logs with the given query and decodes the page.
func listLogs(t *testing.T, q url.Values) logscontract.LogPage {
	t.Helper()
	resp, err := client.Raw(context.Background(), http.MethodGet, "/api/logs?"+q.Encode(), nil)
	if err != nil {
		t.Fatalf("list logs: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("list logs: status %d", resp.StatusCode)
	}
	var page logscontract.LogPage
	decodeData(t, resp, &page)
	return page
}

// walkCursor follows nextCursor from the first page to the last at the given
// page size, returning every log visited. base carries the filters to hold
// across pages (the walk adds limit + cursor).
func walkCursor(t *testing.T, base url.Values, limit int) []logscontract.Log {
	t.Helper()
	var all []logscontract.Log
	cursor := ""
	for pages := 0; pages < 1000; pages++ {
		q := url.Values{}
		for k, vs := range base {
			q[k] = vs
		}
		q.Set("limit", fmt.Sprintf("%d", limit))
		if cursor != "" {
			q.Set("cursor", cursor)
		}
		page := listLogs(t, q)
		all = append(all, page.Items...)
		if page.NextCursor == "" {
			return all
		}
		cursor = page.NextCursor
	}
	t.Fatalf("cursor walk did not terminate")
	return all
}

// decodeData unmarshals the envelope's data field into out.
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

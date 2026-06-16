package steps

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	devicescontract "simulator/packages/contracts/devices"

	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/common/types"
	"github.com/Mapex-Solutions/mapexDevicesSimulator/e2eTests/core/saga"
)

// createDevice POSTs the built device to the simulator and publishes its ids:
// the generic BagKeyDeviceID / BagKeyDeviceDeviceID (for the fire step and the
// log asserts) plus the caller's transport-specific id key (so the matching
// Compensate deletes the right device at rollback). It returns the new id.
func createDevice(c *saga.Context, spec devicescontract.DeviceInput, idKey string) error {
	resp, err := c.Clients.Sim.Raw(c.Stdctx, http.MethodPost, "/api/devices", spec)
	if err != nil {
		return fmt.Errorf("create device: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("create device: unexpected status %d", resp.StatusCode)
	}
	var env types.Envelope
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		return fmt.Errorf("decode create device: %w", err)
	}
	var dev devicescontract.Device
	if err := json.Unmarshal(env.Data, &dev); err != nil {
		return fmt.Errorf("decode device data: %w", err)
	}
	if dev.ID == "" {
		return fmt.Errorf("create device: empty id in response")
	}
	c.Set(BagKeyDeviceID, dev.ID)
	c.Set(BagKeyDeviceDeviceID, dev.DeviceID)
	c.Set(idKey, dev.ID)
	return nil
}

// deleteDevice removes the device whose id lives under idKey. It is the shared
// Compensate body: a no-op when the key is unset, tolerant of an already-gone
// device.
func deleteDevice(c *saga.Context, idKey string) error {
	id, ok := c.Get(idKey)
	if !ok {
		return nil
	}
	resp, err := c.Clients.Sim.Raw(context.Background(), http.MethodDelete, "/api/devices/"+id.(string), nil)
	if err != nil {
		return err
	}
	return resp.Body.Close()
}

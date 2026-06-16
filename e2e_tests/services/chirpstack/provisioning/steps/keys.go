// Package steps owns the ChirpStack provisioning actions: the application
// context (tenant, application, device profile) and the per-transport gateway +
// device registration the OTAA join needs. Each mutating step registers a
// Compensate so a journey leaves the LNS as it found it even though the stack is
// also torn down wholesale at the end.
package steps

const (
	// Application context, written by EnsureApplicationContext.
	BagKeyTenantID        = "chirpstack.tenantId"
	BagKeyApplicationID   = "chirpstack.applicationId"
	BagKeyDeviceProfileID = "chirpstack.deviceProfileId"

	// BagKeyActiveDevEUI is the DevEUI the current phase is exercising. Each
	// provision step overwrites it; the join/uplink asserts read it inline (they
	// run before the next phase overwrites), so one assert pair serves every
	// transport without per-transport duplication.
	BagKeyActiveDevEUI = "chirpstack.activeDevEui"

	// BagKeyDownlinkHex is the hex of the payload EnqueueDownlink queued, so the
	// simulator-side downlink assert can match it in the device's logs.
	BagKeyDownlinkHex = "chirpstack.downlinkHex"
)

// downlinkFPort and downlinkBytes are the fixed downlink EnqueueDownlink sends.
// The bytes are arbitrary; the assert matches their hex in the down frame.
const (
	downlinkFPort = 10
)

var (
	downlinkBytes = []byte{0xA1, 0xB2, 0xC3}

	// UDP transport identities, written by ProvisionUDPDevice. Kept distinct from
	// the Basics Station keys so each phase's Compensate deletes the right entity
	// at rollback time regardless of what a later phase wrote.
	BagKeyUDPGatewayEUI = "chirpstack.udp.gatewayEui"
	BagKeyUDPDevEUI     = "chirpstack.udp.devEui"
	BagKeyUDPAppKey     = "chirpstack.udp.appKey"

	// Basics Station transport identities, written by ProvisionBasicStationDevice.
	BagKeyBSGatewayEUI = "chirpstack.bs.gatewayEui"
	BagKeyBSDevEUI     = "chirpstack.bs.devEui"
	BagKeyBSAppKey     = "chirpstack.bs.appKey"
)

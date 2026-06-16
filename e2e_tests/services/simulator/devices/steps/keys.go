// Package steps holds saga steps that exercise the simulator's devices module
// (POST/DELETE /api/devices). Each create variant has its own canonical payload
// builder; the step reads the runtime values it needs from the bag and writes
// the created ids back.
package steps

const (
	// joinEUI is the OTAA Join EUI convention the LoRaWAN journey uses for both
	// transports (all-zero is the common "no specific join server" value).
	joinEUI = "0000000000000000"

	// BagKeyDeviceID is the server id (UUID) of the most recently created device,
	// read by the fire step which always targets the device just created.
	BagKeyDeviceID = "devices.id"
	// BagKeyDeviceDeviceID is the user-facing deviceId the console and logs filter
	// on (used by the HTTP/MQTT log asserts).
	BagKeyDeviceDeviceID = "devices.deviceId"

	// Transport-specific device ids, kept distinct so each step's Compensate
	// deletes the right device at rollback time even after a later phase created
	// another device.
	BagKeyHTTPDeviceID         = "devices.httpId"
	BagKeyMQTTUserPassDeviceID = "devices.mqttUserPassId"
	BagKeyMQTTTLSDeviceID      = "devices.mqttTlsId"
	BagKeyLoRaWANDeviceID      = "devices.lorawanId"
	BagKeyBasicStationDeviceID = "devices.basicstationId"
	BagKeyMQTTReceiveDeviceID  = "devices.mqttReceiveId"
)

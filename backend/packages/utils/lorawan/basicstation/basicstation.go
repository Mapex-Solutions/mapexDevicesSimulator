// Package basicstation frames the LoRaWAN Basics Station LNS protocol (the modern
// gateway-to-LNS transport over WebSocket). It builds the gateway version handshake
// and the uplink messages (jreq/updf) by decomposing a raw PHYPayload into the
// fields Basics Station expects, and decodes the downlink (dnmsg) back to a raw
// PHYPayload. Framing only; the WebSocket lives in the connector. Any LNS speaking
// Basics Station (TTS, mapexLNS, ChirpStack) interoperates.
package basicstation

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
)

// errShort is returned when an uplink PHYPayload is too short to decompose.
var errShort = errors.New("basicstation: PHYPayload too short")

// mTypeMask isolates the message type in the MHDR.
const mTypeMask byte = 0xE0

// UpInfo carries the simulated radio metadata attached to every uplink.
type UpInfo struct {
	RSSI   int     `json:"rssi"`
	SNR    float64 `json:"snr"`
	RxTime float64 `json:"rxtime"`
	XTime  int64   `json:"xtime"`
}

// VersionMessage is the first frame a gateway sends after the WebSocket opens; the
// LNS replies with a router_config the connector can ignore for a pure uplink path.
func VersionMessage(station string) ([]byte, error) {
	return json.Marshal(map[string]any{
		"msgtype":  "version",
		"station":  station,
		"firmware": "1.0.0",
		"model":    "mapex-sim",
		"protocol": 2,
	})
}

// MarshalUplink decomposes a raw uplink PHYPayload and frames it as a Basics Station
// jreq (join request) or updf (data uplink), attaching the radio context.
func MarshalUplink(phy []byte, dr int, freq uint64, info UpInfo) ([]byte, error) {
	if len(phy) < 12 {
		return nil, errShort
	}
	mhdr := phy[0]
	mic := int32(binary.LittleEndian.Uint32(phy[len(phy)-4:]))

	if mhdr&mTypeMask == 0x00 { // join request
		body := phy[1 : len(phy)-4] // JoinEUI(8) | DevEUI(8) | DevNonce(2), little-endian
		if len(body) < 18 {
			return nil, errShort
		}
		return json.Marshal(map[string]any{
			"msgtype":  "jreq",
			"MHdr":     int(mhdr),
			"JoinEui":  int64(binary.LittleEndian.Uint64(body[0:8])),
			"DevEui":   int64(binary.LittleEndian.Uint64(body[8:16])),
			"DevNonce": int(binary.LittleEndian.Uint16(body[16:18])),
			"MIC":      mic,
			"DR":       dr,
			"Freq":     int(freq),
			"upinfo":   info,
		})
	}

	mac := phy[1 : len(phy)-4]
	if len(mac) < 7 {
		return nil, errShort
	}
	fctrl := mac[4]
	fOptsLen := int(fctrl & 0x0F)
	off := 7 + fOptsLen
	fOpts := mac[7:off]
	var fPort int = -1
	var frm []byte
	if off < len(mac) {
		fPort = int(mac[off])
		frm = mac[off+1:]
	}
	return json.Marshal(map[string]any{
		"msgtype":    "updf",
		"MHdr":       int(mhdr),
		"DevAddr":    int32(binary.LittleEndian.Uint32(mac[0:4])),
		"FCtrl":      int(fctrl),
		"FCnt":       int(binary.LittleEndian.Uint16(mac[5:7])),
		"FOpts":      hex.EncodeToString(fOpts),
		"FPort":      fPort,
		"FRMPayload": hex.EncodeToString(frm),
		"MIC":        mic,
		"DR":         dr,
		"Freq":       int(freq),
		"upinfo":     info,
	})
}

type downlink struct {
	MsgType string `json:"msgtype"`
	PDU     string `json:"pdu"`
}

// ParseDownlink decodes a Basics Station server message; for a dnmsg it returns the
// raw PHYPayload (hex pdu) to hand to the device, otherwise a nil payload.
func ParseDownlink(frame []byte) ([]byte, error) {
	var dn downlink
	if err := json.Unmarshal(frame, &dn); err != nil {
		return nil, err
	}
	if dn.MsgType != "dnmsg" || dn.PDU == "" {
		return nil, nil
	}
	return hex.DecodeString(dn.PDU)
}

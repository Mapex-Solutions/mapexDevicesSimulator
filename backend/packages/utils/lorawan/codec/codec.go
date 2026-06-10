// Package codec assembles and parses LoRaWAN PHYPayloads for the device side of the
// simulator: OTAA join requests, data uplinks, and decoding the matching downlinks.
// It is our own thin implementation over plain structs (no ttnpb), backed by the
// vendored crypto package, so the simulator stays vendor-neutral and works against
// any LNS. Multi-byte LoRaWAN fields are little-endian on the wire.
package codec

import (
	"encoding/binary"
	"errors"

	"simulator/packages/utils/lorawan/crypto"
	"simulator/packages/utils/lorawan/types"
)

// MType is the LoRaWAN message type (top 3 bits of the MHDR).
type MType byte

const (
	MTypeJoinRequest         MType = 0x00
	MTypeJoinAccept          MType = 0x20
	MTypeUnconfirmedDataUp   MType = 0x40
	MTypeUnconfirmedDataDown MType = 0x60
	MTypeConfirmedDataUp     MType = 0x80
	MTypeConfirmedDataDown   MType = 0xA0
	mhdrMTypeMask            byte  = 0xE0
)

// errMalformed is returned when a buffer is too short to parse.
var errMalformed = errors.New("codec: malformed PHYPayload")

// reverse returns a little-endian copy of b (LoRaWAN sends EUIs/addresses LSB-first).
func reverse(b []byte) []byte {
	out := make([]byte, len(b))
	for i := range b {
		out[len(b)-1-i] = b[i]
	}
	return out
}

// JoinRequest holds the device identity for an OTAA join.
type JoinRequest struct {
	JoinEUI  types.EUI64
	DevEUI   types.EUI64
	DevNonce types.DevNonce
}

// MarshalJoinRequest builds a signed JoinRequest PHYPayload:
// MHDR | JoinEUI | DevEUI | DevNonce | MIC, all little-endian, MIC over the rest.
func MarshalJoinRequest(jr JoinRequest, nwkKey types.AES128Key) ([]byte, error) {
	buf := make([]byte, 0, 23)
	buf = append(buf, byte(MTypeJoinRequest))
	buf = append(buf, reverse(jr.JoinEUI[:])...)
	buf = append(buf, reverse(jr.DevEUI[:])...)
	buf = append(buf, reverse(jr.DevNonce[:])...)
	mic, err := crypto.ComputeJoinRequestMIC(nwkKey, buf)
	if err != nil {
		return nil, err
	}
	return append(buf, mic[:]...), nil
}

// JoinAccept is the decrypted join-accept the device derives its session from.
type JoinAccept struct {
	JoinNonce types.JoinNonce
	NetID     types.NetID
	DevAddr   types.DevAddr
	DLSettings byte
	RxDelay    byte
}

// UnmarshalJoinAccept decrypts and parses a JoinAccept downlink with the device's
// app/nwk key. The encrypted block uses the AES-encrypt primitive in the decrypt
// direction (LoRaWAN quirk), handled by the vendored crypto.
func UnmarshalJoinAccept(phy []byte, key types.AES128Key) (JoinAccept, error) {
	if len(phy) < 1+12+4 {
		return JoinAccept{}, errMalformed
	}
	if phy[0]&mhdrMTypeMask != byte(MTypeJoinAccept) {
		return JoinAccept{}, errors.New("codec: not a join-accept")
	}
	dec, err := crypto.DecryptJoinAccept(key, phy[1:])
	if err != nil {
		return JoinAccept{}, err
	}
	if len(dec) < 12 {
		return JoinAccept{}, errMalformed
	}
	var ja JoinAccept
	copy(ja.JoinNonce[:], reverse(dec[0:3]))
	copy(ja.NetID[:], reverse(dec[3:6]))
	copy(ja.DevAddr[:], reverse(dec[6:10]))
	ja.DLSettings = dec[10]
	ja.RxDelay = dec[11]
	return ja, nil
}

// DataUplink is one application uplink to assemble.
type DataUplink struct {
	DevAddr    types.DevAddr
	FCnt       uint32
	FPort      byte
	FRMPayload []byte
	Confirmed  bool
	ADR        bool
}

// MarshalDataUplink builds a signed, encrypted data uplink PHYPayload for LoRaWAN
// 1.0.x (legacy MIC). It encrypts FRMPayload with appSKey and signs the message
// with nwkSKey. fCnt is the full 32-bit counter; the low 16 bits go on the wire.
func MarshalDataUplink(up DataUplink, nwkSKey, appSKey types.AES128Key) ([]byte, error) {
	mType := MTypeUnconfirmedDataUp
	if up.Confirmed {
		mType = MTypeConfirmedDataUp
	}

	var fctrl byte
	if up.ADR {
		fctrl |= 0x80
	}

	enc, err := crypto.EncryptUplink(appSKey, up.DevAddr, up.FCnt, up.FRMPayload)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 0, 12+len(enc)+4)
	buf = append(buf, byte(mType))             // MHDR
	buf = append(buf, reverse(up.DevAddr[:])...) // FHDR.DevAddr
	buf = append(buf, fctrl)                     // FHDR.FCtrl (FOptsLen 0)
	fcnt := make([]byte, 2)
	binary.LittleEndian.PutUint16(fcnt, uint16(up.FCnt))
	buf = append(buf, fcnt...) // FHDR.FCnt (low 16)
	buf = append(buf, up.FPort)
	buf = append(buf, enc...)

	mic, err := crypto.ComputeLegacyUplinkMIC(nwkSKey, up.DevAddr, up.FCnt, buf)
	if err != nil {
		return nil, err
	}
	return append(buf, mic[:]...), nil
}

// DataDownlink is a parsed application downlink.
type DataDownlink struct {
	MType      MType
	DevAddr    types.DevAddr
	FCnt       uint32 // low 16 bits as received
	FPort      byte
	FRMPayload []byte // decrypted
	Confirmed  bool
}

// UnmarshalDataDownlink parses and decrypts a data downlink for LoRaWAN 1.0.x. The
// device supplies the full 32-bit downlink frame counter (the wire carries only the
// low 16 bits) so the correct counter block decrypts the FRMPayload.
func UnmarshalDataDownlink(phy []byte, fullFCnt uint32, appSKey types.AES128Key) (DataDownlink, error) {
	if len(phy) < 1+7+4 { // MHDR + minimal FHDR + MIC
		return DataDownlink{}, errMalformed
	}
	mType := MType(phy[0] & mhdrMTypeMask)
	if mType != MTypeUnconfirmedDataDown && mType != MTypeConfirmedDataDown {
		return DataDownlink{}, errors.New("codec: not a data downlink")
	}

	body := phy[1 : len(phy)-4] // strip MHDR and MIC
	var dl DataDownlink
	dl.MType = mType
	dl.Confirmed = mType == MTypeConfirmedDataDown
	copy(dl.DevAddr[:], reverse(body[0:4]))
	fctrl := body[4]
	fOptsLen := int(fctrl & 0x0F)
	dl.FCnt = (fullFCnt &^ 0xFFFF) | uint32(binary.LittleEndian.Uint16(body[5:7]))

	off := 7 + fOptsLen
	if off > len(body) {
		return DataDownlink{}, errMalformed
	}
	if off < len(body) { // an FPort + FRMPayload are present
		dl.FPort = body[off]
		frm := body[off+1:]
		if dl.FPort == 0 {
			return DataDownlink{}, errors.New("codec: FPort 0 (MAC) downlink not supported")
		}
		dec, err := crypto.DecryptDownlink(appSKey, dl.DevAddr, dl.FCnt, frm)
		if err != nil {
			return DataDownlink{}, err
		}
		dl.FRMPayload = dec
	}
	return dl, nil
}

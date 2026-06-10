package codec

import (
	"bytes"
	"encoding/binary"
	"testing"

	"simulator/packages/utils/lorawan/crypto"
	"simulator/packages/utils/lorawan/types"
)

func TestMarshalJoinRequest(t *testing.T) {
	jr := JoinRequest{
		JoinEUI:  types.EUI64{0x70, 0xB3, 0xD5, 0x7E, 0xD0, 0x00, 0x00, 0x01},
		DevEUI:   types.EUI64{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77},
		DevNonce: types.DevNonce{0x00, 0x01},
	}
	key := types.AES128Key{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

	phy, err := MarshalJoinRequest(jr, key)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if len(phy) != 23 { // 1 MHDR + 8 + 8 + 2 + 4 MIC
		t.Fatalf("join request length = %d, want 23", len(phy))
	}
	if phy[0] != byte(MTypeJoinRequest) {
		t.Fatalf("MHDR = %#x, want join-request", phy[0])
	}
	// JoinEUI is little-endian on the wire: first body byte is the JoinEUI LSB.
	if phy[1] != jr.JoinEUI[7] {
		t.Fatalf("JoinEUI not little-endian: %#x", phy[1])
	}
	// MIC matches an independent computation over the body.
	mic, _ := crypto.ComputeJoinRequestMIC(key, phy[:19])
	if !bytes.Equal(phy[19:], mic[:]) {
		t.Fatalf("MIC mismatch")
	}
}

func TestDataUplinkRoundTrip(t *testing.T) {
	nwkSKey := types.AES128Key{16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	appSKey := types.AES128Key{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	up := DataUplink{
		DevAddr:    types.DevAddr{0x26, 0x01, 0x1F, 0x88},
		FCnt:       42,
		FPort:      10,
		FRMPayload: []byte("temp=21"),
	}

	phy, err := MarshalDataUplink(up, nwkSKey, appSKey)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if phy[0] != byte(MTypeUnconfirmedDataUp) {
		t.Fatalf("MHDR = %#x", phy[0])
	}
	// DevAddr little-endian in FHDR.
	if !bytes.Equal(phy[1:5], reverse(up.DevAddr[:])) {
		t.Fatalf("DevAddr not little-endian")
	}
	// FCnt low 16 bits.
	if binary.LittleEndian.Uint16(phy[6:8]) != uint16(up.FCnt) {
		t.Fatalf("FCnt wrong")
	}
	// FRMPayload decrypts back (it sits after MHDR(1)+FHDR(7)+FPort(1), before MIC(4)).
	enc := phy[9 : len(phy)-4]
	dec, err := crypto.DecryptUplink(appSKey, up.DevAddr, up.FCnt, enc)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	if string(dec) != "temp=21" {
		t.Fatalf("payload round-trip: got %q", dec)
	}
}

func TestUnmarshalDataDownlink(t *testing.T) {
	appSKey := types.AES128Key{2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2, 2}
	addr := types.DevAddr{0x26, 0x01, 0x1F, 0x88}
	fcnt := uint32(5)
	plain := []byte("ack-on")

	enc, err := crypto.EncryptDownlink(appSKey, addr, fcnt, plain)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	// Build a minimal data downlink: MHDR | DevAddr | FCtrl | FCnt | FPort | FRM | MIC.
	buf := []byte{byte(MTypeUnconfirmedDataDown)}
	buf = append(buf, reverse(addr[:])...)
	buf = append(buf, 0x00) // FCtrl, no FOpts
	fc := make([]byte, 2)
	binary.LittleEndian.PutUint16(fc, uint16(fcnt))
	buf = append(buf, fc...)
	buf = append(buf, 2) // FPort
	buf = append(buf, enc...)
	buf = append(buf, 0xDE, 0xAD, 0xBE, 0xEF) // MIC (not verified by the decoder)

	dl, err := UnmarshalDataDownlink(buf, fcnt, appSKey)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if dl.DevAddr != addr {
		t.Fatalf("devaddr mismatch")
	}
	if dl.FPort != 2 {
		t.Fatalf("fport = %d", dl.FPort)
	}
	if string(dl.FRMPayload) != "ack-on" {
		t.Fatalf("downlink payload: got %q", dl.FRMPayload)
	}
}

package crypto

import (
	"encoding/hex"
	"testing"

	"simulator/packages/utils/lorawan/types"
)

// Static golden vectors for the vendored LoRaWAN crypto. The values were produced
// by this verbatim-copied crypto (which is byte-identical to The Things Stack
// v3.35.0) for the fixed inputs below, then frozen here. The test guards against
// any future drift in the vendored files (an accidental edit changing output) with
// zero dependency on TTS — it is the parity safety net the plan calls for.
//
// Inputs (fixed):
//
//	appKey/nwkKey = 0102...10, joinEUI = 70B3D57ED0000001, devNonce = 0001,
//	joinNonce = 000001, netID = 000013, devAddr = 26011F88, fCnt = 1
func TestCryptoGoldenVectors(t *testing.T) {
	appKey := types.AES128Key{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	joinEUI := types.EUI64{0x70, 0xB3, 0xD5, 0x7E, 0xD0, 0x00, 0x00, 0x01}
	devNonce := types.DevNonce{0x00, 0x01}
	joinNonce := types.JoinNonce{0x00, 0x00, 0x01}
	netID := types.NetID{0x00, 0x00, 0x13}
	devAddr := types.DevAddr{0x26, 0x01, 0x1F, 0x88}
	frm := []byte{0x01, 0x02, 0x03, 0x04}
	micPayload := []byte{0x40, 0x88, 0x1F, 0x01, 0x26, 0x00, 0x01, 0x00, 0x01, 0x02, 0x03, 0x04}
	joinReq := []byte{0x00, 0x01, 0x00, 0x00, 0xD0, 0x7E, 0xD5, 0xB3, 0x70, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01}

	lnwk := DeriveLegacyNwkSKey(appKey, joinNonce, netID, devNonce)
	lapp := DeriveLegacyAppSKey(appKey, joinNonce, netID, devNonce)
	fnwk := DeriveFNwkSIntKey(appKey, joinNonce, joinEUI, devNonce)
	appSKey := DeriveAppSKey(appKey, joinNonce, joinEUI, devNonce)
	enc, err := EncryptUplink(lapp, devAddr, 1, frm)
	if err != nil {
		t.Fatalf("EncryptUplink: %v", err)
	}
	mic, err := ComputeLegacyUplinkMIC(lnwk, devAddr, 1, micPayload)
	if err != nil {
		t.Fatalf("ComputeLegacyUplinkMIC: %v", err)
	}
	jmic, err := ComputeJoinRequestMIC(appKey, joinReq)
	if err != nil {
		t.Fatalf("ComputeJoinRequestMIC: %v", err)
	}

	cases := []struct {
		name, got, want string
	}{
		{"legacyNwkSKey (1.0.x)", hex.EncodeToString(lnwk[:]), "d65223a784d7666015545c8ec026f068"},
		{"legacyAppSKey (1.0.x)", hex.EncodeToString(lapp[:]), "93571900dfbdb45f7b2a50b6e35cbfc8"},
		{"fNwkSIntKey (1.1)", hex.EncodeToString(fnwk[:]), "165b72f281289965e7f47b87c928d9e6"},
		{"appSKey (1.1)", hex.EncodeToString(appSKey[:]), "33412f7b5d97478395cbad4336995d9f"},
		{"encryptUplink", hex.EncodeToString(enc), "833d8b21"},
		{"legacyUplinkMIC", hex.EncodeToString(mic[:]), "d51405fb"},
		{"joinRequestMIC", hex.EncodeToString(jmic[:]), "b8758b4c"},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s: got %s, want %s (vendored crypto drifted)", c.name, c.got, c.want)
		}
	}
}

// TestEncryptUplinkRoundTrip confirms uplink encryption is reversible with the same
// key/addr/counter, exercising the FRMPayload cipher both directions.
func TestEncryptUplinkRoundTrip(t *testing.T) {
	key := types.AES128Key{9, 8, 7, 6, 5, 4, 3, 2, 1, 0, 1, 2, 3, 4, 5, 6}
	addr := types.DevAddr{0x01, 0x02, 0x03, 0x04}
	plain := []byte("hello-lorawan")
	enc, err := EncryptUplink(key, addr, 7, plain)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}
	dec, err := DecryptUplink(key, addr, 7, enc)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	if string(dec) != string(plain) {
		t.Fatalf("round-trip mismatch: got %q want %q", dec, plain)
	}
}

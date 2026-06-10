package device

import (
	"testing"

	"simulator/packages/utils/lorawan/crypto"
	"simulator/packages/utils/lorawan/types"
)

func rev(b []byte) []byte {
	out := make([]byte, len(b))
	for i := range b {
		out[len(b)-1-i] = b[i]
	}
	return out
}

// buildJoinAccept plays the LNS: it assembles a 1.0.x join accept the device can
// decrypt and derive its session from. Wire fields are little-endian, matching the
// codec's parse (which reverses them back).
func buildJoinAccept(t *testing.T, appKey types.AES128Key, jn types.JoinNonce, nid types.NetID, addr types.DevAddr) []byte {
	t.Helper()
	body := make([]byte, 0, 12)
	body = append(body, rev(jn[:])...)
	body = append(body, rev(nid[:])...)
	body = append(body, rev(addr[:])...)
	body = append(body, 0x00) // DLSettings
	body = append(body, 0x01) // RxDelay
	mic, err := crypto.ComputeLegacyJoinAcceptMIC(appKey, append([]byte{0x20}, body...))
	if err != nil {
		t.Fatalf("join accept mic: %v", err)
	}
	enc, err := crypto.EncryptJoinAccept(appKey, append(body, mic[:]...))
	if err != nil {
		t.Fatalf("encrypt join accept: %v", err)
	}
	return append([]byte{0x20}, enc...)
}

func TestOTAAFullRoundTrip(t *testing.T) {
	appKey := types.AES128Key{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	netID := types.NetID{0x00, 0x00, 0x13}
	cfg := Config{
		JoinEUI:    types.EUI64{0x70, 0xB3, 0xD5, 0x7E, 0xD0, 0x00, 0x00, 0x01},
		DevEUI:     types.EUI64{0, 1, 2, 3, 4, 5, 6, 7},
		AppKey:     appKey,
		NetID:      netID,
		Activation: "otaa",
	}
	dev := New(cfg)

	// 1. device builds a join request (DevNonce advances to 1).
	if _, err := dev.JoinRequest(); err != nil {
		t.Fatalf("join request: %v", err)
	}
	if dev.Joined() {
		t.Fatal("device must not be joined before processing the accept")
	}

	// 2. LNS sends a join accept; device derives its session.
	jn := types.JoinNonce{0x00, 0x00, 0x05}
	addr := types.DevAddr{0x26, 0x01, 0x1F, 0x88}
	phyJA := buildJoinAccept(t, appKey, jn, netID, addr)
	if err := dev.ProcessJoinAccept(phyJA); err != nil {
		t.Fatalf("process join accept: %v", err)
	}
	if !dev.Joined() || dev.DevAddr() != addr {
		t.Fatalf("join failed: joined=%v addr=%x", dev.Joined(), dev.DevAddr())
	}

	// keys the LNS would derive (DevNonce = 1 -> wire {0x01,0x00}).
	dn := types.DevNonce{0x01, 0x00}
	appSKey := crypto.DeriveLegacyAppSKey(appKey, jn, netID, dn)

	// 3. device sends an uplink; LNS decrypts it.
	up, err := dev.BuildUplink(10, []byte("temp=21"), false)
	if err != nil {
		t.Fatalf("build uplink: %v", err)
	}
	enc := up[9 : len(up)-4] // after MHDR(1)+FHDR(7)+FPort(1), before MIC(4)
	dec, err := crypto.DecryptUplink(appSKey, addr, 0, enc)
	if err != nil {
		t.Fatalf("lns decrypt uplink: %v", err)
	}
	if string(dec) != "temp=21" {
		t.Fatalf("uplink payload: got %q", dec)
	}
	if dev.FCntUp() != 1 {
		t.Fatalf("fCntUp should advance to 1, got %d", dev.FCntUp())
	}

	// 4. LNS sends a downlink; device decodes it.
	dlEnc, err := crypto.EncryptDownlink(appSKey, addr, 0, []byte("ack"))
	if err != nil {
		t.Fatalf("encrypt downlink: %v", err)
	}
	dlPhy := []byte{0x60}                 // unconfirmed data down
	dlPhy = append(dlPhy, rev(addr[:])...) // DevAddr
	dlPhy = append(dlPhy, 0x00, 0x00, 0x00) // FCtrl + FCnt(0)
	dlPhy = append(dlPhy, 0x02)             // FPort
	dlPhy = append(dlPhy, dlEnc...)
	dlPhy = append(dlPhy, 0, 0, 0, 0) // MIC (not verified)

	dl, err := dev.DecodeDownlink(dlPhy)
	if err != nil {
		t.Fatalf("decode downlink: %v", err)
	}
	if string(dl.FRMPayload) != "ack" {
		t.Fatalf("downlink payload: got %q", dl.FRMPayload)
	}
}

func TestABPActivatesImmediately(t *testing.T) {
	dev := New(Config{
		Activation: "abp",
		DevAddr:    types.DevAddr{1, 2, 3, 4},
		NwkSKey:    types.AES128Key{1},
		AppSKey:    types.AES128Key{2},
	})
	if !dev.Joined() {
		t.Fatal("ABP device should be joined immediately")
	}
	if _, err := dev.BuildUplink(1, []byte("x"), false); err != nil {
		t.Fatalf("abp uplink: %v", err)
	}
}

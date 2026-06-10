package basicstation

import (
	"encoding/hex"
	"encoding/json"
	"testing"
)

func TestMarshalUplinkClassifiesByMHDR(t *testing.T) {
	join := make([]byte, 23) // MHDR 0x00 + 18 body + 4 MIC
	dataUp := make([]byte, 16)
	dataUp[0] = 0x40 // unconfirmed data up

	tests := []struct {
		name    string
		phy     []byte
		wantMsg string
	}{
		{"join request", join, "jreq"},
		{"data uplink", dataUp, "updf"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := MarshalUplink(tt.phy, 5, 868_100_000, UpInfo{RSSI: -42, SNR: 9})
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			var m map[string]any
			if err := json.Unmarshal(out, &m); err != nil {
				t.Fatalf("json: %v", err)
			}
			if m["msgtype"] != tt.wantMsg {
				t.Fatalf("msgtype = %v, want %s", m["msgtype"], tt.wantMsg)
			}
			if int(m["Freq"].(float64)) != 868_100_000 {
				t.Fatalf("freq not attached")
			}
		})
	}
}

func TestParseDownlinkExtractsPDU(t *testing.T) {
	phy := []byte{0x60, 0x11, 0x22, 0x33}
	frame, _ := json.Marshal(map[string]any{"msgtype": "dnmsg", "pdu": hex.EncodeToString(phy)})

	got, err := ParseDownlink(frame)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if hex.EncodeToString(got) != hex.EncodeToString(phy) {
		t.Fatalf("pdu mismatch: %x", got)
	}

	// A non-downlink server message yields no payload, no error.
	other, _ := json.Marshal(map[string]any{"msgtype": "router_config"})
	if p, err := ParseDownlink(other); err != nil || p != nil {
		t.Fatalf("router_config should be ignored, got %x err %v", p, err)
	}
}

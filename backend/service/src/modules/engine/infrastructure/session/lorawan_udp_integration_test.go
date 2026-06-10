package session

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net"
	"slices"
	"testing"
	"time"

	"simulator/packages/utils/lorawan/crypto"
	"simulator/packages/utils/lorawan/types"
	"simulator/service/src/modules/engine/application/ports"
)

// stubLNS is a minimal Semtech UDP "LNS" over a real socket: it reads PUSH_DATA,
// and when the uplink is a join request it replies with an encrypted join accept as
// a PULL_RESP, exactly as ChirpStack's gateway bridge would. It lets the test drive
// the real udpTransport end-to-end without an external LNS.
type stubLNS struct {
	conn   *net.UDPConn
	appKey types.AES128Key
	gotUp  chan []byte
	gotAck chan uint16
}

// pullRespToken is the distinctive token the stub puts on its PULL_RESP, so the
// test can assert the gateway's TX_ACK echoes it back.
const pullRespToken uint16 = 0x002A

func newStubLNS(t *testing.T, appKey types.AES128Key) *stubLNS {
	t.Helper()
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0})
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	s := &stubLNS{conn: conn, appKey: appKey, gotUp: make(chan []byte, 8), gotAck: make(chan uint16, 4)}
	go s.serve()
	return s
}

func (s *stubLNS) port() int { return s.conn.LocalAddr().(*net.UDPAddr).Port }

func (s *stubLNS) serve() {
	buf := make([]byte, 2048)
	for {
		n, from, err := s.conn.ReadFromUDP(buf)
		if err != nil {
			return
		}
		if n >= 4 && buf[3] == 0x05 { // TX_ACK: record the echoed token
			s.gotAck <- uint16(buf[1])<<8 | uint16(buf[2])
			continue
		}
		if n < 4 || buf[3] != 0x00 { // only PUSH_DATA carries rxpk
			continue
		}
		phy := s.extractPHY(buf[12:n])
		if phy == nil {
			continue
		}
		s.gotUp <- phy
		if phy[0]&0xE0 == 0x00 { // join request -> reply with a join accept downlink
			s.conn.WriteToUDP(s.joinAccept(), from)
		}
	}
}

// extractPHY pulls the base64 PHYPayload out of an rxpk array.
func (s *stubLNS) extractPHY(body []byte) []byte {
	var p struct {
		RxPk []struct {
			Data string `json:"data"`
		} `json:"rxpk"`
	}
	if err := json.Unmarshal(body, &p); err != nil || len(p.RxPk) == 0 {
		return nil
	}
	phy, err := base64.StdEncoding.DecodeString(p.RxPk[0].Data)
	if err != nil {
		return nil
	}
	return phy
}

// joinAccept builds a PULL_RESP carrying an encrypted 1.0.x join accept.
func (s *stubLNS) joinAccept() []byte {
	phy := buildJoinAcceptPHY(s.appKey)
	txpk, _ := json.Marshal(map[string]any{"txpk": map[string]any{"data": base64.StdEncoding.EncodeToString(phy)}})
	// PULL_RESP header carries pullRespToken; the gateway's TX_ACK must echo it.
	return append([]byte{0x02, byte(pullRespToken >> 8), byte(pullRespToken), 0x03}, txpk...)
}

// buildJoinAcceptPHY assembles an encrypted 1.0.x join-accept PHYPayload signed with
// the device's app key, shared by the UDP and WebSocket stub LNS tests.
func buildJoinAcceptPHY(appKey types.AES128Key) []byte {
	jn := types.JoinNonce{0x00, 0x00, 0x07}
	nid := types.NetID{0x00, 0x00, 0x13}
	addr := types.DevAddr{0x26, 0x0B, 0xAD, 0x01}
	rev := func(b []byte) []byte {
		o := make([]byte, len(b))
		for i := range b {
			o[len(b)-1-i] = b[i]
		}
		return o
	}
	body := append(append(append(rev(jn[:]), rev(nid[:])...), rev(addr[:])...), 0x00, 0x01)
	mic, _ := crypto.ComputeLegacyJoinAcceptMIC(appKey, append([]byte{0x20}, body...))
	enc, _ := crypto.EncryptJoinAccept(appKey, append(body, mic[:]...))
	return append([]byte{0x20}, enc...)
}

func TestLoRaWAN_OTAAJoinOverRealUDP(t *testing.T) {
	appKey := types.AES128Key{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	lns := newStubLNS(t, appKey)
	defer lns.conn.Close()

	spec := ports.SessionSpec{
		Protocol: "lorawan", DeviceID: "dev-lora", DeviceName: "Field",
		LoRaWAN: &ports.LoRaWANSpec{
			Region: "EU868", MACVersion: "1.0.3", Activation: "otaa",
			JoinEUI: "70B3D57ED0000001", DevEUI: "0011223344556677",
			AppKey:       "0102030405060708090A0B0C0D0E0F10",
			GatewayEUI:   "0016C001F1500001",
			LinkProtocol: "udp", LinkUDPHost: "127.0.0.1", LinkUDPPort: lns.port(),
		},
	}

	var statuses []string
	connector := NewLoRaWAN()
	sess, err := connector.Open(context.Background(),
		spec,
		func(ports.InboundMessage) {},
		func(status, _ string) { statuses = append(statuses, status) },
	)
	if err != nil {
		t.Fatalf("open (join) failed over the wire: %v", err)
	}
	defer sess.Close()

	// The join request reached the LNS, and the device joined from the accept.
	select {
	case phy := <-lns.gotUp:
		if phy[0]&0xE0 != 0x00 {
			t.Fatalf("first uplink was not a join request: mhdr=%#x", phy[0])
		}
	case <-time.After(2 * time.Second):
		t.Fatal("LNS never received the join request")
	}
	if !slices.Contains(statuses, "joined") {
		t.Fatalf("device did not reach joined; statuses=%v", statuses)
	}

	// The gateway must TX_ACK the join-accept downlink echoing the PULL_RESP token,
	// or a real LNS (ChirpStack) logs "no internal frame cache for token".
	select {
	case tok := <-lns.gotAck:
		if tok != pullRespToken {
			t.Fatalf("TX_ACK token = %#x, want echoed %#x", tok, pullRespToken)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("gateway never TX_ACKed the downlink")
	}

	// A scheduled-style uplink now flows through the joined session to the LNS.
	res := sess.Send(context.Background(), ports.OutboundMessage{FPort: 10, Payload: "01ff"})
	if !res.OK {
		t.Fatalf("uplink send failed: %v", res.Err)
	}
	select {
	case phy := <-lns.gotUp:
		if phy[0]&0xE0 != 0x40 {
			t.Fatalf("expected an unconfirmed data uplink, mhdr=%#x", phy[0])
		}
	case <-time.After(2 * time.Second):
		t.Fatal("LNS never received the data uplink")
	}
}

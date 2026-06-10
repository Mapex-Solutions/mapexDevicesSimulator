package session

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	"simulator/packages/utils/lorawan/types"
	"simulator/service/src/modules/engine/application/ports"
)

// stubBasicStationLNS is a minimal Basics Station LNS over a real WebSocket: it
// accepts the gateway connection, ignores the version handshake, replies to a join
// request (jreq) with a join-accept dnmsg, and records data uplinks (updf). It lets
// the test drive the real wsTransport end-to-end, the same flow ChirpStack runs.
type stubBasicStationLNS struct {
	server *httptest.Server
	appKey types.AES128Key
	gotUp  chan string
}

func newStubBasicStationLNS(t *testing.T, appKey types.AES128Key) *stubBasicStationLNS {
	t.Helper()
	s := &stubBasicStationLNS{appKey: appKey, gotUp: make(chan string, 8)}
	up := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := up.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		s.handle(conn)
	}))
	return s
}

// wsURL returns the ws:// endpoint the gateway dials.
func (s *stubBasicStationLNS) wsURL() string {
	return strings.Replace(s.server.URL, "http://", "ws://", 1) + "/router-data"
}

func (s *stubBasicStationLNS) handle(conn *websocket.Conn) {
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			return
		}
		var m map[string]any
		if json.Unmarshal(data, &m) != nil {
			continue
		}
		switch m["msgtype"] {
		case "jreq":
			dn, _ := json.Marshal(map[string]any{"msgtype": "dnmsg", "pdu": hex.EncodeToString(buildJoinAcceptPHY(s.appKey))})
			_ = conn.WriteMessage(websocket.TextMessage, dn)
		case "updf":
			s.gotUp <- "updf"
		}
	}
}

func TestLoRaWAN_OTAAJoinOverRealWebSocket(t *testing.T) {
	appKey := types.AES128Key{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	lns := newStubBasicStationLNS(t, appKey)
	defer lns.server.Close()

	spec := ports.SessionSpec{
		Protocol: "basicstation", DeviceID: "dev-bs", DeviceName: "Field BS",
		LoRaWAN: &ports.LoRaWANSpec{
			Region: "EU868", MACVersion: "1.0.3", Activation: "otaa",
			JoinEUI: "70B3D57ED0000001", DevEUI: "0011223344556677",
			AppKey:       "0102030405060708090A0B0C0D0E0F10",
			GatewayEUI:   "0016C001F1500001",
			LinkProtocol: "basicstation", LinkLNSURI: lns.wsURL(),
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
		t.Fatalf("open (join) failed over the WebSocket: %v", err)
	}
	defer sess.Close()

	if !slices.Contains(statuses, "joined") {
		t.Fatalf("device did not join over Basics Station; statuses=%v", statuses)
	}

	res := sess.Send(context.Background(), ports.OutboundMessage{FPort: 10, Payload: "01ff"})
	if !res.OK {
		t.Fatalf("uplink send failed: %v", res.Err)
	}
	select {
	case <-lns.gotUp:
	case <-time.After(2 * time.Second):
		t.Fatal("LNS never received the data uplink (updf)")
	}
}

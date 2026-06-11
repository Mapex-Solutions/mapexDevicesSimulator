package session

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	"simulator/packages/utils/lorawan/basicstation"
	"simulator/packages/utils/lorawan/udp"
)

// newDownlinkRouter builds an empty router.
func newDownlinkRouter() *downlinkRouter {
	return &downlinkRouter{byAddr: make(map[[4]byte]func([]byte))}
}

// route dispatches one downlink PHYPayload. A join accept (MHDR top bits 001) goes
// to the pending handler; otherwise it is matched by the DevAddr in the FHDR.
func (r *downlinkRouter) route(phy []byte) {
	if len(phy) < 5 {
		return
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	if phy[0]&0xE0 == 0x20 { // join accept
		if r.pending != nil {
			r.pending(phy)
		}
		return
	}
	var addr [4]byte
	addr[0], addr[1], addr[2], addr[3] = phy[4], phy[3], phy[2], phy[1] // FHDR DevAddr is little-endian
	if h, ok := r.byAddr[addr]; ok {
		h(phy)
	}
}

// setPending registers the handler awaiting a join accept.
func (r *downlinkRouter) setPending(h func([]byte)) {
	r.mu.Lock()
	r.pending = h
	r.mu.Unlock()
}

// bind registers a device's data-downlink handler under its DevAddr.
func (r *downlinkRouter) bind(addr [4]byte, h func([]byte)) {
	r.mu.Lock()
	r.byAddr[addr] = h
	r.pending = nil
	r.mu.Unlock()
}

// unbind removes a device's downlink handler.
func (r *downlinkRouter) unbind(addr [4]byte) {
	r.mu.Lock()
	delete(r.byAddr, addr)
	r.mu.Unlock()
}

// dialUDP opens the UDP socket and starts the keepalive, stats, and receive loops.
func dialUDP(ctx context.Context, host string, port int, gwEUI [8]byte, route func([]byte)) (linkTransport, error) {
	addr := &net.UDPAddr{IP: net.ParseIP(host), Port: port}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}
	t := &udpTransport{conn: conn, gwEUI: gwEUI, up: true}
	go t.keepalive(ctx)
	go t.stats(ctx)
	go t.receive(route)
	return t, nil
}

// stats periodically reports gateway statistics so the LNS shows the gateway online
// (its last-seen). It sends one immediately so the gateway appears connected at once.
func (t *udpTransport) stats(ctx context.Context) {
	tick := time.NewTicker(statInterval)
	defer tick.Stop()
	t.writeStat()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			t.writeStat()
		}
	}
}

// writeStat sends one gateway-statistics datagram.
func (t *udpTransport) writeStat() {
	t.mu.Lock()
	t.token++
	tok, rxnb := t.token, t.rxnb
	t.mu.Unlock()
	dg, err := udp.MarshalStat(tok, t.gwEUI, udp.Stat{
		Time: time.Now().UTC().Format("2006-01-02 15:04:05 GMT"),
		RxNb: rxnb,
		RxOK: rxnb,
		RxFw: rxnb,
		AckR: 100.0,
	})
	if err == nil {
		_, _ = t.conn.Write(dg)
	}
}

// keepalive sends PULL_DATA so the LNS knows where to deliver downlinks.
func (t *udpTransport) keepalive(ctx context.Context) {
	tick := time.NewTicker(10 * time.Second)
	defer tick.Stop()
	t.writePullData()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			t.writePullData()
		}
	}
}

// writePullData sends one keepalive datagram.
func (t *udpTransport) writePullData() {
	t.mu.Lock()
	t.token++
	tok := t.token
	t.mu.Unlock()
	_, _ = t.conn.Write(udp.MarshalPullData(tok, t.gwEUI))
}

// receive reads server datagrams, acknowledges downlinks, and routes the PHYPayload.
func (t *udpTransport) receive(route func([]byte)) {
	buf := make([]byte, 2048)
	for {
		n, err := t.conn.Read(buf)
		if err != nil {
			t.mu.Lock()
			t.up = false
			t.mu.Unlock()
			return
		}
		kind, phy, err := udp.Classify(buf[:n])
		if err == nil && kind == udp.KindPullResp && phy != nil {
			// TX_ACK must echo the PULL_RESP token (bytes 1-2) so the LNS can match
			// it to the downlink it scheduled; a fresh token is rejected.
			t.writeTxAck(uint16(buf[1])<<8 | uint16(buf[2]))
			route(phy)
		}
	}
}

// writeTxAck confirms a downlink, echoing the PULL_RESP token the LNS is waiting on.
func (t *udpTransport) writeTxAck(token uint16) {
	_, _ = t.conn.Write(udp.MarshalTxAck(token, t.gwEUI))
}

// sendUp frames and writes one uplink PUSH_DATA. The internal timestamp (tmst)
// advances per uplink so the LNS can schedule a downlink against it.
func (t *udpTransport) sendUp(phy []byte, dr int, freq uint64, datr string, rssi int, snr float64) error {
	t.mu.Lock()
	t.token++
	t.tmst += 1_000_000
	t.rxnb++
	tok, tmst := t.token, t.tmst
	t.mu.Unlock()
	pk := udp.RxPk{
		Time: time.Now().UTC().Format(time.RFC3339Nano),
		Tmst: tmst,
		Freq: float64(freq) / 1e6,
		Chan: 0,
		RFCh: 0,
		Stat: 1,
		Modu: "LORA",
		DatR: datr,
		CodR: "4/5",
		RSSI: rssi,
		LSNR: snr,
		Size: len(phy),
		Data: udp.EncodePHY(phy),
	}
	dg, err := udp.MarshalPushData(tok, t.gwEUI, pk)
	if err != nil {
		return err
	}
	_, err = t.conn.Write(dg)
	return err
}

// connected reports the last known socket state.
func (t *udpTransport) connected() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.up
}

// close tears the socket down.
func (t *udpTransport) close() { _ = t.conn.Close() }

// dialWS opens the Basics Station data WebSocket, sends the version handshake, and
// starts the receive loop. When the LNS URI carries no path it defaults to the
// per-gateway data endpoint /gw/<eui> (ChirpStack's convention); a full path is
// used as-is so other LNSs can be targeted.
func dialWS(ctx context.Context, lnsURI, gatewayEUI string, route func([]byte)) (linkTransport, error) {
	u, err := url.Parse(lnsURI)
	if err != nil {
		return nil, err
	}
	if u.Path == "" || u.Path == "/" {
		u.Path = "/gw/" + strings.ToLower(gatewayEUI)
	}
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, u.String(), nil)
	if err != nil {
		return nil, err
	}
	t := &wsTransport{conn: conn, up: true}
	if version, err := basicstation.VersionMessage(gatewayEUI); err == nil {
		_ = conn.WriteMessage(websocket.TextMessage, version)
	}
	go t.receive(route)
	return t, nil
}

// receive reads server frames and routes any dnmsg PHYPayload.
func (t *wsTransport) receive(route func([]byte)) {
	for {
		_, data, err := t.conn.ReadMessage()
		if err != nil {
			t.mu.Lock()
			t.up = false
			t.mu.Unlock()
			return
		}
		if phy, err := basicstation.ParseDownlink(data); err == nil && phy != nil {
			route(phy)
		}
	}
}

// sendUp frames and writes one uplink message. The monotonic xtime gives the LNS a
// reference to schedule the matching downlink (e.g. the join accept) against.
func (t *wsTransport) sendUp(phy []byte, dr int, freq uint64, _ string, rssi int, snr float64) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.xtime += 1_000_000
	frame, err := basicstation.MarshalUplink(phy, dr, freq, basicstation.UpInfo{RSSI: rssi, SNR: snr, XTime: t.xtime})
	if err != nil {
		return err
	}
	return t.conn.WriteMessage(websocket.TextMessage, frame)
}

// connected reports the last known socket state.
func (t *wsTransport) connected() bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.up
}

// close tears the WebSocket down.
func (t *wsTransport) close() { _ = t.conn.Close() }

// dataRate maps a spreading factor to a Semtech datr string (125 kHz bandwidth).
func dataRate(sf int) string { return fmt.Sprintf("SF%dBW125", sf) }

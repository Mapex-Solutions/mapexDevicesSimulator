package session

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"simulator/packages/utils/lorawan/band"
	"simulator/packages/utils/lorawan/device"
	"simulator/packages/utils/lorawan/types"
	"simulator/service/src/modules/engine/application/ports"
)

// Compile-time proofs the adapter satisfies its ports.
var (
	_ ports.Connector = (*lorawanConnector)(nil)
	_ ports.Session   = (*lorawanSession)(nil)
)

// NewLoRaWAN builds the LoRaWAN connector (serves both the lorawan and basicstation
// protocol ids; the difference is only where the link config comes from).
func NewLoRaWAN() ports.Connector {
	return &lorawanConnector{
		links: make(map[string]*sharedLink),
		devs:  make(map[string]*device.Session),
	}
}

// deviceFor returns the persistent device brain for a DevEUI, creating it on first
// use. Reusing it across reconnects keeps the DevNonce monotonic, which OTAA join
// servers require (a reused DevNonce is rejected as a replay).
func (c *lorawanConnector) deviceFor(spec *ports.LoRaWANSpec) (*device.Session, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if d, ok := c.devs[spec.DevEUI]; ok {
		return d, nil
	}
	d, err := buildDeviceSession(spec)
	if err != nil {
		return nil, err
	}
	c.devs[spec.DevEUI] = d
	return d, nil
}

// Protocol identifies this connector. It is also registered under basicstation in
// the registry.
func (c *lorawanConnector) Protocol() string { return "lorawan" }

// Open builds the device brain, attaches it to the gateway link, performs the OTAA
// join (or activates ABP), and binds the device's downlink handler.
func (c *lorawanConnector) Open(ctx context.Context, spec ports.SessionSpec, in ports.InboundSink, status ports.StatusSink) (ports.Session, error) {
	if spec.LoRaWAN == nil {
		return nil, errors.New("lorawan: missing LoRaWAN spec")
	}
	region := band.Get(spec.LoRaWAN.Region)
	dev, err := c.deviceFor(spec.LoRaWAN)
	if err != nil {
		return nil, err
	}

	link, err := c.acquireLink(ctx, spec.LoRaWAN)
	if err != nil {
		return nil, err
	}

	sess := &lorawanSession{dev: dev, link: link, conn: c, region: region}
	deliver := func(phy []byte) { sess.onDownlink(phy, in) }

	switch {
	case spec.LoRaWAN.Activation == "abp":
		sess.addr = dev.DevAddr()
		link.router.bind(sess.addr, deliver)
		status("activated", "ABP")
	case dev.Joined():
		// Already joined on a previous connection: resume rather than re-join, so a
		// transient reconnect does not burn a DevNonce or reset the session.
		sess.addr = dev.DevAddr()
		link.router.bind(sess.addr, deliver)
		status("joined", hex.EncodeToString(sess.addr[:]))
	default:
		if err := sess.join(link, status, deliver); err != nil {
			c.releaseLink(link)
			return nil, err
		}
	}
	return sess, nil
}

// join sends the join request and waits for the join accept routed back over the
// link, then derives the session and binds the device by its new DevAddr.
func (s *lorawanSession) join(link *sharedLink, status ports.StatusSink, deliver func([]byte)) error {
	accepted := make(chan []byte, 1)
	link.router.setPending(func(phy []byte) {
		select {
		case accepted <- phy:
		default:
		}
	})

	phy, err := s.dev.JoinRequest()
	if err != nil {
		return err
	}
	status("join-request", "")
	if err := link.transport.sendUp(phy, s.region.UplinkDR, s.region.UplinkFrequency, dataRate(s.region.UplinkSF), -42, 9.0); err != nil {
		return err
	}

	select {
	case ja := <-accepted:
		status("join-accept", "")
		if err := s.dev.ProcessJoinAccept(ja); err != nil {
			return err
		}
		s.addr = s.dev.DevAddr()
		link.router.bind(s.addr, deliver)
		status("joined", hex.EncodeToString(s.addr[:]))
		return nil
	case <-time.After(joinTimeout):
		return errors.New("lorawan: join accept timeout")
	}
}

// Send builds the next uplink from the rendered hex payload and transmits it.
func (s *lorawanSession) Send(_ context.Context, msg ports.OutboundMessage) ports.SendResult {
	payload, err := hex.DecodeString(msg.Payload)
	if err != nil {
		return ports.SendResult{Err: fmt.Errorf("lorawan: payload must be hex: %w", err)}
	}
	phy, err := s.dev.BuildUplink(msg.FPort, payload, msg.Confirmed)
	if err != nil {
		return ports.SendResult{Err: err}
	}
	if err := s.link.transport.sendUp(phy, s.region.UplinkDR, s.region.UplinkFrequency, dataRate(s.region.UplinkSF), -42, 9.0); err != nil {
		return ports.SendResult{Err: err}
	}
	return ports.SendResult{OK: true, Status: fmt.Sprintf("FCnt %d", s.dev.FCntUp()-1)}
}

// onDownlink decodes a received downlink and surfaces it through the inbound sink.
func (s *lorawanSession) onDownlink(phy []byte, in ports.InboundSink) {
	dl, err := s.dev.DecodeDownlink(phy)
	if err != nil {
		return
	}
	in(ports.InboundMessage{
		Payload: hex.EncodeToString(dl.FRMPayload),
		Status:  fmt.Sprintf("FPort %d · FCnt %d", dl.FPort, dl.FCnt),
		Summary: fmt.Sprintf("downlink FPort %d", dl.FPort),
	})
}

// Close unbinds the device and releases its share of the gateway link.
func (s *lorawanSession) Close() error {
	s.link.router.unbind(s.addr)
	s.conn.releaseLink(s.link)
	return nil
}

// Connected reports whether the shared gateway link is up.
func (s *lorawanSession) Connected() bool { return s.link.transport.connected() }

// acquireLink returns the shared gateway link for the device's endpoint, dialing it
// on first use and reference-counting subsequent devices on the same gateway.
func (c *lorawanConnector) acquireLink(ctx context.Context, spec *ports.LoRaWANSpec) (*sharedLink, error) {
	key := linkKey(spec)
	c.mu.Lock()
	defer c.mu.Unlock()
	if l, ok := c.links[key]; ok {
		l.refs++
		return l, nil
	}

	router := newDownlinkRouter()
	transport, err := dialTransport(ctx, spec, router.route)
	if err != nil {
		return nil, err
	}
	l := &sharedLink{key: key, transport: transport, router: router, refs: 1}
	c.links[key] = l
	return l, nil
}

// releaseLink drops one reference and tears the link down when the last device leaves.
func (c *lorawanConnector) releaseLink(l *sharedLink) {
	c.mu.Lock()
	defer c.mu.Unlock()
	l.refs--
	if l.refs <= 0 {
		l.transport.close()
		delete(c.links, l.key)
	}
}

// dialTransport opens the right transport for the gateway's link protocol.
func dialTransport(ctx context.Context, spec *ports.LoRaWANSpec, route func([]byte)) (linkTransport, error) {
	gwEUI, err := parseEUI(spec.GatewayEUI)
	if err != nil {
		return nil, err
	}
	switch spec.LinkProtocol {
	case "udp":
		return dialUDP(ctx, spec.LinkUDPHost, spec.LinkUDPPort, gwEUI, route)
	case "basicstation":
		return dialWS(ctx, spec.LinkLNSURI, spec.GatewayEUI, route)
	default:
		return nil, fmt.Errorf("lorawan: unknown link protocol %q", spec.LinkProtocol)
	}
}

// linkKey identifies a shared gateway link by its endpoint.
func linkKey(spec *ports.LoRaWANSpec) string {
	if spec.LinkProtocol == "udp" {
		return fmt.Sprintf("udp|%s:%d", spec.LinkUDPHost, spec.LinkUDPPort)
	}
	return "ws|" + spec.LinkLNSURI
}

// buildDeviceSession parses the hex key material into a device session config.
func buildDeviceSession(spec *ports.LoRaWANSpec) (*device.Session, error) {
	cfg := device.Config{
		Activation: spec.Activation,
		MACVersion: spec.MACVersion,
		Region:     spec.Region,
	}
	var err error
	if cfg.JoinEUI, err = parseEUI(spec.JoinEUI); err != nil {
		return nil, err
	}
	if cfg.DevEUI, err = parseEUI(spec.DevEUI); err != nil {
		return nil, err
	}
	if cfg.AppKey, err = parseKey(spec.AppKey); err != nil {
		return nil, err
	}
	if strings.HasPrefix(spec.MACVersion, "1.1") {
		if cfg.NwkKey, err = parseKey(spec.NwkKey); err != nil {
			return nil, fmt.Errorf("lorawan 1.1 requires a valid nwkKey: %w", err)
		}
	}
	if spec.Activation == "abp" {
		if cfg.DevAddr, err = parseAddr(spec.DevAddr); err != nil {
			return nil, err
		}
		if cfg.NwkSKey, err = parseKey(spec.NwkSKey); err != nil {
			return nil, err
		}
		if cfg.AppSKey, err = parseKey(spec.AppSKey); err != nil {
			return nil, err
		}
	}
	return device.New(cfg), nil
}

// parseKey decodes a 16-byte hex AES key.
func parseKey(s string) (types.AES128Key, error) {
	var k types.AES128Key
	b, err := hex.DecodeString(s)
	if err != nil || len(b) != 16 {
		return k, fmt.Errorf("lorawan: bad 128-bit key %q", s)
	}
	copy(k[:], b)
	return k, nil
}

// parseEUI decodes an 8-byte hex EUI (most-significant-byte first).
func parseEUI(s string) (types.EUI64, error) {
	var e types.EUI64
	if s == "" {
		return e, nil
	}
	b, err := hex.DecodeString(s)
	if err != nil || len(b) != 8 {
		return e, fmt.Errorf("lorawan: bad EUI %q", s)
	}
	copy(e[:], b)
	return e, nil
}

// parseAddr decodes a 4-byte hex DevAddr.
func parseAddr(s string) (types.DevAddr, error) {
	var a types.DevAddr
	b, err := hex.DecodeString(s)
	if err != nil || len(b) != 4 {
		return a, fmt.Errorf("lorawan: bad DevAddr %q", s)
	}
	copy(a[:], b)
	return a, nil
}

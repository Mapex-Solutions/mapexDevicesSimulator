package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	devicescontract "simulator/packages/contracts/devices"
	gatewayscontract "simulator/packages/contracts/gateways"
	consoleDtos "simulator/service/src/modules/console/application/dtos"
	enginePorts "simulator/service/src/modules/engine/application/ports"
	"simulator/service/src/modules/engine/domain/entities"
	logsDtos "simulator/service/src/modules/logs/application/dtos"
)

// sessionHandle is the engine's grip on one device's live-connection supervisor:
// the cancel func that stops it, the signature that detects config changes, and
// the current live session (nil while connecting/reconnecting).
type sessionHandle struct {
	sig    string
	cancel context.CancelFunc

	mu   sync.RWMutex
	live enginePorts.Session
}

func (h *sessionHandle) set(s enginePorts.Session) {
	h.mu.Lock()
	h.live = s
	h.mu.Unlock()
}

func (h *sessionHandle) get() enginePorts.Session {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.live
}

// desiredSession is one device that should have a live connection.
type desiredSession struct {
	spec enginePorts.SessionSpec
	sig  string
}

// reconcileSessions aligns the running supervisors with the desired set: stop the
// ones whose device disappeared or was disabled, (re)start the new or changed ones.
// It mirrors reconcile() for the scheduler's job heap.
func (s *EngineService) reconcileSessions() {
	desired := s.buildDesiredSessions()

	s.sessMu.Lock()
	defer s.sessMu.Unlock()

	for id, h := range s.sessions {
		if _, ok := desired[id]; !ok {
			h.cancel()
			delete(s.sessions, id)
		}
	}
	for id, d := range desired {
		if existing, ok := s.sessions[id]; ok {
			if existing.sig == d.sig {
				continue
			}
			existing.cancel() // config changed: tear down and re-open
			delete(s.sessions, id)
		}
		ctx, cancel := context.WithCancel(s.ctx)
		h := &sessionHandle{sig: d.sig, cancel: cancel}
		s.sessions[id] = h
		s.wg.Add(1)
		go s.superviseSession(ctx, d.spec, h)
	}
}

// buildDesiredSessions derives the session set from the enabled devices on
// session-capable protocols. HTTP has no connector, so it never appears here.
func (s *EngineService) buildDesiredSessions() map[string]desiredSession {
	desired := make(map[string]desiredSession)
	devices, err := s.deps.Devices.List(context.Background())
	if err != nil {
		return desired
	}
	for _, d := range devices {
		if !d.Enabled {
			continue
		}
		if _, ok := s.deps.Connectors.Connector(d.ProtocolID); !ok {
			continue
		}
		spec, ok := s.buildSessionSpec(d)
		if !ok {
			continue
		}
		desired[d.ID] = desiredSession{spec: spec, sig: sessionSignature(spec)}
	}
	return desired
}

// buildSessionSpec resolves a device into the target + credentials for its live
// connection. Subscriptions are included only when the device has receive enabled.
func (s *EngineService) buildSessionSpec(d devicescontract.Device) (enginePorts.SessionSpec, bool) {
	spec := enginePorts.SessionSpec{
		Protocol:   d.ProtocolID,
		DeviceID:   d.DeviceID,
		DeviceName: d.Name,
		StoreLogs:  d.StoreLogs,
	}
	switch d.ProtocolID {
	case "mqtt":
		var cfg entities.MQTTConnectionConfig
		if err := json.Unmarshal(d.Config, &cfg); err != nil || cfg.BrokerURL == "" {
			return enginePorts.SessionSpec{}, false
		}
		spec.BrokerURL = cfg.BrokerURL
		spec.ClientID = cfg.ClientID
		spec.Username = cfg.Username
		spec.Password = cfg.Password
		if cfg.ReceiveEnabled {
			for _, sub := range cfg.Subscriptions {
				topic := joinTopic(cfg.BaseTopic, sub.Topic)
				if topic == "" {
					continue
				}
				spec.Subscriptions = append(spec.Subscriptions, enginePorts.Subscription{
					Name:  sub.Name,
					Topic: topic,
					QoS:   byte(sub.QoS),
				})
			}
		}
		return spec, true
	case "lorawan":
		return s.buildLoRaWANSessionSpec(spec, d.Config)
	case "basicstation":
		return s.buildBasicsStationSessionSpec(spec, d.Config)
	default:
		return enginePorts.SessionSpec{}, false
	}
}

// buildLoRaWANSessionSpec resolves a lorawan device's keys plus the shared gateway
// link it transmits through (looked up by the device's gatewayId).
func (s *EngineService) buildLoRaWANSessionSpec(spec enginePorts.SessionSpec, config json.RawMessage) (enginePorts.SessionSpec, bool) {
	var cfg entities.LoRaWANConnectionConfig
	if err := json.Unmarshal(config, &cfg); err != nil {
		return enginePorts.SessionSpec{}, false
	}
	gw, ok := s.findGateway(cfg.GatewayID)
	if !ok || !gw.Enabled {
		return enginePorts.SessionSpec{}, false
	}
	var link entities.GatewayLink
	_ = json.Unmarshal(gw.Link, &link)
	spec.LoRaWAN = &enginePorts.LoRaWANSpec{
		Region:       cfg.Region,
		MACVersion:   cfg.MACVersion,
		Activation:   cfg.Activation,
		JoinEUI:      cfg.JoinEUI,
		DevEUI:       cfg.DevEUI,
		AppKey:       cfg.AppKey,
		NwkKey:       cfg.NwkKey,
		DevAddr:      cfg.DevAddr,
		NwkSKey:      cfg.NwkSKey,
		AppSKey:      cfg.AppSKey,
		GatewayEUI:   gw.EUI,
		LinkProtocol: link.Protocol,
		LinkLNSURI:   link.LNSURI,
		LinkUDPHost:  link.Host,
		LinkUDPPort:  link.Port,
	}
	return spec, true
}

// buildBasicsStationSessionSpec resolves a basicstation device, whose Basics Station
// link is carried on the device itself (its own embedded gateway).
func (s *EngineService) buildBasicsStationSessionSpec(spec enginePorts.SessionSpec, config json.RawMessage) (enginePorts.SessionSpec, bool) {
	var cfg entities.BasicsStationConnectionConfig
	if err := json.Unmarshal(config, &cfg); err != nil {
		return enginePorts.SessionSpec{}, false
	}
	spec.LoRaWAN = &enginePorts.LoRaWANSpec{
		Region:       cfg.Region,
		MACVersion:   cfg.MACVersion,
		Activation:   cfg.Activation,
		JoinEUI:      cfg.JoinEUI,
		DevEUI:       cfg.DevEUI,
		AppKey:       cfg.AppKey,
		NwkKey:       cfg.NwkKey,
		DevAddr:      cfg.DevAddr,
		NwkSKey:      cfg.NwkSKey,
		AppSKey:      cfg.AppSKey,
		GatewayEUI:   cfg.GatewayEUI,
		LinkProtocol: "basicstation",
		LinkLNSURI:   cfg.LNSURI,
	}
	return spec, true
}

// findGateway resolves a gateway by id from the gateways service.
func (s *EngineService) findGateway(id string) (gatewayscontract.Gateway, bool) {
	gws, err := s.deps.Gateways.List(context.Background())
	if err != nil {
		return gatewayscontract.Gateway{}, false
	}
	for _, gw := range gws {
		if gw.ID == id {
			return gw, true
		}
	}
	return gatewayscontract.Gateway{}, false
}

// superviseSession owns one device's connection lifecycle: connect, hold open
// (re-opening on a silent drop), and reconnect forever with bounded backoff. Every
// transition is streamed to the console so the full status is visible.
func (s *EngineService) superviseSession(ctx context.Context, spec enginePorts.SessionSpec, h *sessionHandle) {
	defer s.wg.Done()

	connector, ok := s.deps.Connectors.Connector(spec.Protocol)
	if !ok {
		return
	}
	inbound := func(msg enginePorts.InboundMessage) { s.emitInbound(spec, msg) }
	status := func(st, detail string) { s.emitStatus(spec, st, detail) }

	backoff := sessionBackoffInitial
	attempt := 0
	for {
		if ctx.Err() != nil {
			return
		}
		attempt++
		s.emitStatus(spec, "connecting", spec.BrokerURL)

		sess, err := connector.Open(ctx, spec, inbound, status)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			s.emitStatus(spec, "reconnecting",
				fmt.Sprintf("attempt %d, next %s: %s", attempt, backoff.Round(time.Millisecond), err.Error()))
			if sleepCtx(ctx, jitter(backoff)) {
				return
			}
			backoff = nextBackoff(backoff)
			continue
		}

		attempt = 0
		backoff = sessionBackoffInitial
		h.set(sess)
		s.emitStatus(spec, "connected", spec.BrokerURL)

		stopped := s.monitorSession(ctx, sess)
		h.set(nil)
		_ = sess.Close()
		if stopped {
			s.emitStatus(spec, "disconnected", "")
			return
		}
		s.emitStatus(spec, "reconnecting", "connection lost")
		if sleepCtx(ctx, jitter(backoff)) {
			return
		}
		backoff = nextBackoff(backoff)
	}
}

// monitorSession blocks until the engine stops (returns true) or the session drops
// (returns false), polling the connection so a silent drop triggers a reconnect.
func (s *EngineService) monitorSession(ctx context.Context, sess enginePorts.Session) bool {
	t := time.NewTicker(sessionPollEvery)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return true
		case <-t.C:
			if !sess.Connected() {
				return false
			}
		}
	}
}

// liveSession returns the connected session for a device's server id, if any.
// Sessions are keyed by server id (d.ID), so callers pass that, not the
// user-facing DeviceID (which is author-controlled and not guaranteed unique).
func (s *EngineService) liveSession(deviceKey string) (enginePorts.Session, bool) {
	s.sessMu.Lock()
	h, ok := s.sessions[deviceKey]
	s.sessMu.Unlock()
	if !ok {
		return nil, false
	}
	sess := h.get()
	if sess == nil || !sess.Connected() {
		return nil, false
	}
	return sess, true
}

// emitInbound streams a received downlink to the console (always) and persists a
// log (when storeLogs). It is the inbound twin of report().
func (s *EngineService) emitInbound(spec enginePorts.SessionSpec, msg enginePorts.InboundMessage) {
	now := time.Now().UTC().Format(time.RFC3339)
	s.deps.Console.Publish(consoleDtos.ConsoleMessage{
		ID:         uuid.NewString(),
		TS:         now,
		Protocol:   spec.Protocol,
		DeviceID:   spec.DeviceID,
		DeviceName: spec.DeviceName,
		Direction:  "down",
		Kind:       "downlink",
		Summary:    msg.Summary,
		Payload:    msg.Payload,
		Status:     msg.Status,
	})
	if spec.StoreLogs {
		_ = s.deps.Logs.Append(s.ctx, &logsDtos.LogInput{
			Protocol:   spec.Protocol,
			DeviceID:   spec.DeviceID,
			DeviceName: spec.DeviceName,
			Direction:  "down",
			Kind:       "downlink",
			Summary:    msg.Summary,
			Status:     msg.Status,
			Payload:    msg.Payload,
		})
	}
}

// emitStatus streams a connection-lifecycle frame (system/status) to the console.
func (s *EngineService) emitStatus(spec enginePorts.SessionSpec, status, detail string) {
	now := time.Now().UTC().Format(time.RFC3339)
	s.deps.Console.Publish(consoleDtos.ConsoleMessage{
		ID:         uuid.NewString(),
		TS:         now,
		Protocol:   spec.Protocol,
		DeviceID:   spec.DeviceID,
		DeviceName: spec.DeviceName,
		Direction:  "system",
		Kind:       "status",
		Summary:    status,
		Payload:    detail,
		Status:     status,
	})
}

// sessionSignature captures the fields whose change requires re-opening a session.
func sessionSignature(spec enginePorts.SessionSpec) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s|%s|%s|%s", spec.BrokerURL, spec.ClientID, spec.Username, spec.Password)
	for _, sub := range spec.Subscriptions {
		fmt.Fprintf(&b, "|%s:%d", sub.Topic, sub.QoS)
	}
	return b.String()
}

// nextBackoff doubles the delay up to the ceiling.
func nextBackoff(d time.Duration) time.Duration {
	next := d * 2
	if next > sessionBackoffMax {
		return sessionBackoffMax
	}
	return next
}

// jitter spreads reconnect attempts by adding up to 20% so many devices on the
// same broker do not retry in lockstep.
func jitter(d time.Duration) time.Duration {
	return d + time.Duration(rand.Int63n(int64(d/5)+1))
}

// sleepCtx waits for d or until ctx is cancelled; it returns true if cancelled.
func sleepCtx(ctx context.Context, d time.Duration) bool {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return true
	case <-t.C:
		return false
	}
}

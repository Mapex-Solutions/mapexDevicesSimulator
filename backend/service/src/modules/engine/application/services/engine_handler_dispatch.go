package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	consoleDtos "simulator/service/src/modules/console/application/dtos"
	enginePorts "simulator/service/src/modules/engine/application/ports"
	domainsvc "simulator/service/src/modules/engine/domain/services"
	logsDtos "simulator/service/src/modules/logs/application/dtos"
)

// process renders the payload, sends it over the device's protocol, and reports
// the result to the console (always) and the logs (when storeLogs is on).
//
// For a session-capable protocol with a live connection, the uplink goes through
// that open session (reusing the connection); otherwise — HTTP, or MQTT while a
// session is still connecting/reconnecting — it falls back to the one-shot
// dispatcher so a send is never silently dropped.
func (s *EngineService) process(task fireTask) {
	spec := task.spec
	payload := domainsvc.Render(spec.payloadTemplate, spec.deviceID, task.counter)
	req, summary := s.buildRequest(spec, payload, task.counter)

	if sess, ok := s.liveSession(spec.deviceKey); ok {
		res := sess.Send(s.ctx, enginePorts.OutboundMessage{
			Topic:     req.Topic,
			QoS:       req.QoS,
			Retain:    req.Retain,
			Payload:   payload,
			FPort:     spec.fport,
			Confirmed: spec.confirmed,
		})
		s.report(spec, payload, summary, enginePorts.DispatchResult{OK: res.OK, Status: res.Status, Err: res.Err})
		return
	}

	disp, ok := s.deps.Registry.For(spec.protocol)
	if !ok {
		// No live session and no one-shot dispatcher for this protocol (a LoRaWAN or
		// Basics Station device with no link). Report the attempt to the console as a
		// status frame instead of dropping it silently, so the user sees the device
		// tried and why it could not send.
		s.reportNoLink(spec, payload, summary)
		return
	}
	res := disp.Dispatch(s.ctx, req)
	s.report(spec, payload, summary, res)
}

// reportNoLink streams a system status frame (and a log when storeLogs) for a fire
// that had no live link to go out on, so the attempt is visible on the console. The
// status names the cause: the device's gateway is offline, or the session is simply
// not up yet.
func (s *EngineService) reportNoLink(spec sendSpec, payload, summary string) {
	status := "disconnected"
	if dev, err := s.findDevice(s.ctx, spec.deviceKey); err == nil && !s.gatewayPermits(*dev) {
		status = "gateway-offline"
	}
	now := time.Now().UTC().Format(time.RFC3339)
	s.deps.Console.Publish(consoleDtos.ConsoleMessage{
		ID:         uuid.NewString(),
		TS:         now,
		Protocol:   spec.protocol,
		DeviceID:   spec.deviceID,
		DeviceName: spec.deviceName,
		Direction:  "system",
		Kind:       "status",
		Summary:    summary,
		Payload:    payload,
		Status:     status,
	})
	if spec.storeLogs {
		_ = s.deps.Logs.Append(s.ctx, &logsDtos.LogInput{
			Protocol:   spec.protocol,
			DeviceID:   spec.deviceID,
			DeviceName: spec.deviceName,
			Direction:  "system",
			Kind:       "status",
			Summary:    summary,
			Status:     status,
			Payload:    payload,
		})
	}
}

// buildRequest resolves the protocol-specific dispatch request and a one-line
// summary, rendering the URL/topic placeholders at fire time.
func (s *EngineService) buildRequest(spec sendSpec, payload string, counter int64) (enginePorts.DispatchRequest, string) {
	switch spec.protocol {
	case "http":
		url := domainsvc.Render(spec.url, spec.deviceID, counter)
		return enginePorts.DispatchRequest{
			URL:     url,
			Method:  spec.method,
			Headers: spec.headers,
			Payload: payload,
		}, spec.method + " " + url
	case "mqtt":
		topic := domainsvc.Render(spec.topic, spec.deviceID, counter)
		return enginePorts.DispatchRequest{
			BrokerURL: spec.brokerURL,
			ClientID:  spec.clientID,
			Username:  spec.username,
			Password:  spec.password,
			Topic:     topic,
			QoS:       spec.qos,
			Retain:    spec.retain,
			Payload:   payload,
		}, "PUBLISH " + topic
	case "lorawan", "basicstation":
		return enginePorts.DispatchRequest{Payload: payload},
			fmt.Sprintf("Uplink FPort %d", spec.fport)
	default:
		return enginePorts.DispatchRequest{Payload: payload}, ""
	}
}

// report streams a console frame (always) and persists a log (when storeLogs).
//
// A failed send (a transport error from the dispatcher or the live session: the
// broker is unreachable, auth was rejected, the request timed out) is reported as
// an "error" frame whose Response carries the reason, so the user sees not just
// that the send failed but why. A non-2xx HTTP reply is not an error here: it keeps
// its status code and response body.
func (s *EngineService) report(spec sendSpec, payload, summary string, res enginePorts.DispatchResult) {
	now := time.Now().UTC().Format(time.RFC3339)
	status := res.Status
	response := res.Response
	if !res.OK && res.Err != nil {
		status = "error"
		response = res.Err.Error()
	}

	s.deps.Console.Publish(consoleDtos.ConsoleMessage{
		ID:         uuid.NewString(),
		TS:         now,
		Protocol:   spec.protocol,
		DeviceID:   spec.deviceID,
		DeviceName: spec.deviceName,
		Direction:  "up",
		Kind:       "data",
		Summary:    summary,
		Payload:    payload,
		Response:   response,
		Status:     status,
	})

	if spec.storeLogs {
		_ = s.deps.Logs.Append(s.ctx, &logsDtos.LogInput{
			Protocol:   spec.protocol,
			DeviceID:   spec.deviceID,
			DeviceName: spec.deviceName,
			Direction:  "up",
			Kind:       "data",
			Summary:    summary,
			Status:     status,
			Payload:    payload,
		})
	}
}

package services

import (
	"time"

	"github.com/google/uuid"

	consoleDtos "simulator/service/src/modules/console/application/dtos"
	enginePorts "simulator/service/src/modules/engine/application/ports"
	domainsvc "simulator/service/src/modules/engine/domain/services"
	logsDtos "simulator/service/src/modules/logs/application/dtos"
)

// process renders the payload, dispatches it over the device's protocol, and
// reports the result to the console (always) and the logs (when storeLogs is on).
func (s *EngineService) process(task fireTask) {
	spec := task.spec
	payload := domainsvc.Render(spec.payloadTemplate, spec.deviceID, task.counter)
	req, summary := s.buildRequest(spec, payload, task.counter)
	disp, ok := s.deps.Registry.For(spec.protocol)
	if !ok {
		return
	}
	res := disp.Dispatch(s.ctx, req)
	s.report(spec, payload, summary, res)
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
	default:
		return enginePorts.DispatchRequest{Payload: payload}, ""
	}
}

// report streams a console frame (always) and persists a log (when storeLogs).
func (s *EngineService) report(spec sendSpec, payload, summary string, res enginePorts.DispatchResult) {
	now := time.Now().UTC().Format(time.RFC3339)
	status := res.Status
	if !res.OK && res.Err != nil {
		status = "error"
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

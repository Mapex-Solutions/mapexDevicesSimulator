package services

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	devicesDtos "simulator/service/src/modules/devices/application/dtos"
	"simulator/service/src/modules/engine/application/di"
	enginePorts "simulator/service/src/modules/engine/application/ports"
	dispatch "simulator/service/src/modules/engine/infrastructure/dispatch"
	session "simulator/service/src/modules/engine/infrastructure/session"
	gatewaysDtos "simulator/service/src/modules/gateways/application/dtos"
	"simulator/service/src/shared/reconcile"
)

// fireHTTPDevice builds an HTTP device with a pre-registered event and NO schedule,
// so the scheduler never auto-fires it — only an explicit Fire() does.
func fireHTTPDevice(url string) devicesDtos.Device {
	return devicesDtos.Device{
		ID:         "d1",
		Enabled:    true,
		StoreLogs:  true,
		Name:       "Sensor",
		DeviceID:   "dev-1",
		ProtocolID: "http",
		Config:     json.RawMessage(`{"url":"` + url + `"}`),
		Events:     json.RawMessage(`[{"id":"e1","name":"tel","http":{"bodyMode":"raw","body":"{\"id\":\"{{deviceId}}\"}"}}]`),
	}
}

func TestEngine_FireByEventID(t *testing.T) {
	var mu sync.Mutex
	var hits int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		mu.Lock()
		hits++
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	pub := &fakePublisher{}
	eng := New(di.EngineServiceDI{
		Devices:    &fakeDevices{list: []devicesDtos.Device{fireHTTPDevice(srv.URL)}},
		Logs:       &fakeLogWriter{},
		Console:    pub,
		Registry:   dispatch.NewRegistry(),
		Connectors: session.NewConnectorRegistry(),
		Reconcile:  reconcile.New(),
	})
	eng.OnMount()
	defer eng.OnShutdown(context.Background())

	if err := eng.Fire(context.Background(), "d1", enginePorts.FireInput{EventID: "e1"}); err != nil {
		t.Fatalf("fire: %v", err)
	}
	mu.Lock()
	got := hits
	mu.Unlock()
	if got != 1 {
		t.Fatalf("fire should hit the target exactly once, got %d", got)
	}
	if pub.count() == 0 {
		t.Fatal("fire should stream a console frame")
	}
}

// A device whose gateway is offline has no live link, so a fire does not reach the
// LNS -- but it must still surface on the console as a system status frame naming
// the gateway as offline, so the user sees the attempt.
func TestEngine_FireOfflineGatewayReportsStatus(t *testing.T) {
	bsDev := devicesDtos.Device{
		ID: "bs1", Enabled: true, StoreLogs: true, Name: "BS", DeviceID: "dev-bs", ProtocolID: "basicstation",
		Config: json.RawMessage(`{"gatewayEui":"0102030405060708","lnsUri":"ws://127.0.0.1:3001","region":"EU868",` +
			`"macVersion":"1.0.3","activation":"otaa","devEui":"0011223344556677","joinEui":"0000000000000000",` +
			`"appKey":"00112233445566778899AABBCCDDEEFF"}`),
		Events: json.RawMessage(`[{"id":"e1","name":"u","lorawan":{"fport":2,"confirmed":false,"payloadHex":"00"}}]`),
	}
	disabledGw := gatewaysDtos.Gateway{
		ID: "gw1", EUI: "0102030405060708", Enabled: false, Region: "EU868",
		Link: json.RawMessage(`{"protocol":"basicstation","lnsUri":"ws://127.0.0.1:3001"}`),
	}
	pub := &fakePublisher{}
	eng := New(di.EngineServiceDI{
		Devices:    &fakeDevices{list: []devicesDtos.Device{bsDev}},
		Gateways:   &fakeGateways{list: []gatewaysDtos.Gateway{disabledGw}},
		Logs:       &fakeLogWriter{},
		Console:    pub,
		Registry:   dispatch.NewRegistry(),
		Connectors: session.NewConnectorRegistry(),
		Reconcile:  reconcile.New(),
	})
	eng.OnMount() // a disabled gateway opens no session
	defer eng.OnShutdown(context.Background())

	if err := eng.Fire(context.Background(), "bs1", enginePorts.FireInput{EventID: "e1"}); err != nil {
		t.Fatalf("fire should not error, got %v", err)
	}
	if !pub.has("system", "status", "gateway-offline") {
		t.Fatal("fire with an offline gateway must report a system gateway-offline status")
	}
}

// A send that fails at the transport layer (here an HTTP device whose target is
// unreachable) must still reach the console as an "error" frame carrying the
// failure reason in Response, so the user sees not just that it failed but why.
func TestEngine_FireReportsSendError(t *testing.T) {
	pub := &fakePublisher{}
	eng := New(di.EngineServiceDI{
		Devices:    &fakeDevices{list: []devicesDtos.Device{fireHTTPDevice("http://127.0.0.1:1/x")}},
		Logs:       &fakeLogWriter{},
		Console:    pub,
		Registry:   dispatch.NewRegistry(),
		Connectors: session.NewConnectorRegistry(),
		Reconcile:  reconcile.New(),
	})
	eng.OnMount()
	defer eng.OnShutdown(context.Background())

	if err := eng.Fire(context.Background(), "d1", enginePorts.FireInput{EventID: "e1"}); err != nil {
		t.Fatalf("fire should not error at the API level, got %v", err)
	}

	m, ok := pub.firstWithStatus("error")
	if !ok {
		t.Fatal("a failed send must stream an error frame to the console")
	}
	if m.Direction != "up" || m.Kind != "data" {
		t.Fatalf("error frame should be an up/data frame, got %s/%s", m.Direction, m.Kind)
	}
	if m.Response == "" {
		t.Fatal("the error frame must carry the failure reason in Response")
	}
}

func TestEngine_FireUnknownDeviceAndEvent(t *testing.T) {
	eng := New(di.EngineServiceDI{
		Devices:    &fakeDevices{list: []devicesDtos.Device{fireHTTPDevice("http://x")}},
		Logs:       &fakeLogWriter{},
		Console:    &fakePublisher{},
		Registry:   dispatch.NewRegistry(),
		Connectors: session.NewConnectorRegistry(),
		Reconcile:  reconcile.New(),
	})
	eng.OnMount()
	defer eng.OnShutdown(context.Background())

	if err := eng.Fire(context.Background(), "nope", enginePorts.FireInput{EventID: "e1"}); err != enginePorts.ErrDeviceNotFound {
		t.Fatalf("want ErrDeviceNotFound, got %v", err)
	}
	if err := eng.Fire(context.Background(), "d1", enginePorts.FireInput{EventID: "missing"}); err != enginePorts.ErrEventNotFound {
		t.Fatalf("want ErrEventNotFound, got %v", err)
	}
}

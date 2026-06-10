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

func TestEngine_FireUnknownDeviceAndEvent(t *testing.T) {
	eng := New(di.EngineServiceDI{
		Devices:    &fakeDevices{list: []devicesDtos.Device{fireHTTPDevice("http://x")}},
		Logs:       &fakeLogWriter{},
		Console:    &fakePublisher{},
		Registry:   dispatch.NewRegistry(),
		Connectors: session.NewConnectorRegistry(),
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

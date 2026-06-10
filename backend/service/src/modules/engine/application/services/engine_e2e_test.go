package services

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	consoleDtos "simulator/service/src/modules/console/application/dtos"
	consolePorts "simulator/service/src/modules/console/application/ports"
	devicesDtos "simulator/service/src/modules/devices/application/dtos"
	devicesPorts "simulator/service/src/modules/devices/application/ports"
	"simulator/service/src/modules/engine/application/di"
	dispatch "simulator/service/src/modules/engine/infrastructure/dispatch"
	logsDtos "simulator/service/src/modules/logs/application/dtos"
	logsPorts "simulator/service/src/modules/logs/application/ports"
)

// --- inline port fakes ---

type fakeDevices struct{ list []devicesDtos.Device }

func (f *fakeDevices) List(context.Context) ([]devicesDtos.Device, error) { return f.list, nil }
func (f *fakeDevices) Create(context.Context, *devicesDtos.DeviceInput) (*devicesDtos.Device, error) {
	return nil, nil
}
func (f *fakeDevices) Update(context.Context, string, *devicesDtos.DeviceInput) (*devicesDtos.Device, error) {
	return nil, nil
}
func (f *fakeDevices) Delete(context.Context, string) (map[string]bool, error) { return nil, nil }

var _ devicesPorts.DevicesServicePort = (*fakeDevices)(nil)

type fakePublisher struct {
	mu   sync.Mutex
	msgs []consoleDtos.ConsoleMessage
}

func (f *fakePublisher) Publish(m consoleDtos.ConsoleMessage) {
	f.mu.Lock()
	f.msgs = append(f.msgs, m)
	f.mu.Unlock()
}
func (f *fakePublisher) count() int { f.mu.Lock(); defer f.mu.Unlock(); return len(f.msgs) }

var _ consolePorts.Publisher = (*fakePublisher)(nil)

type fakeLogWriter struct {
	mu   sync.Mutex
	logs []logsDtos.LogInput
}

func (f *fakeLogWriter) Append(_ context.Context, in *logsDtos.LogInput) error {
	f.mu.Lock()
	f.logs = append(f.logs, *in)
	f.mu.Unlock()
	return nil
}
func (f *fakeLogWriter) count() int { f.mu.Lock(); defer f.mu.Unlock(); return len(f.logs) }

var _ logsPorts.LogWriter = (*fakeLogWriter)(nil)

// httpDevice builds a device that fires one HTTP event every second.
func httpDevice(url string, enabled, storeLogs bool) devicesDtos.Device {
	return devicesDtos.Device{
		ID:         "d1",
		Enabled:    enabled,
		StoreLogs:  storeLogs,
		Name:       "Sensor",
		DeviceID:   "dev-1",
		ProtocolID: "http",
		Config:     json.RawMessage(`{"url":"` + url + `"}`),
		Events: json.RawMessage(`[{"id":"e1","name":"tel",` +
			`"http":{"bodyMode":"raw","body":"{\"v\":{{counter}},\"id\":\"{{deviceId}}\"}"},` +
			`"schedule":{"enabled":true,"every":1,"unit":"seconds"}}]`),
	}
}

func TestEngine_FiresRendersDispatchesReports(t *testing.T) {
	var mu sync.Mutex
	var bodies []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		mu.Lock()
		bodies = append(bodies, string(b))
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	pub := &fakePublisher{}
	lw := &fakeLogWriter{}
	eng := New(di.EngineServiceDI{
		Devices:  &fakeDevices{list: []devicesDtos.Device{httpDevice(srv.URL, true, true)}},
		Logs:     lw,
		Console:  pub,
		Registry: dispatch.NewRegistry(),
	})

	eng.OnMount()
	time.Sleep(1200 * time.Millisecond)
	_ = eng.OnShutdown(context.Background())

	mu.Lock()
	got := append([]string(nil), bodies...)
	mu.Unlock()

	if len(got) == 0 {
		t.Fatal("engine did not fire any HTTP request")
	}
	if !strings.Contains(got[0], `"id":"dev-1"`) || !strings.Contains(got[0], `"v":`) {
		t.Fatalf("rendered body wrong: %q", got[0])
	}
	if pub.count() == 0 {
		t.Fatal("expected console frames")
	}
	if lw.count() == 0 {
		t.Fatal("expected persisted logs (storeLogs on)")
	}
}

func TestEngine_DisabledDeviceDoesNotFire(t *testing.T) {
	var mu sync.Mutex
	var hits int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		mu.Lock()
		hits++
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	eng := New(di.EngineServiceDI{
		Devices:  &fakeDevices{list: []devicesDtos.Device{httpDevice(srv.URL, false, true)}},
		Logs:     &fakeLogWriter{},
		Console:  &fakePublisher{},
		Registry: dispatch.NewRegistry(),
	})
	eng.OnMount()
	time.Sleep(1200 * time.Millisecond)
	_ = eng.OnShutdown(context.Background())

	mu.Lock()
	defer mu.Unlock()
	if hits != 0 {
		t.Fatalf("disabled device fired %d times, want 0", hits)
	}
}

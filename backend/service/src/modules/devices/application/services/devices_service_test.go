package services

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	customErrors "github.com/Mapex-Solutions/mapexGoKit/microservices/http/customErrors"
	status "github.com/Mapex-Solutions/mapexGoKit/microservices/http/status"

	"simulator/service/src/modules/devices/application/di"
	"simulator/service/src/modules/devices/application/dtos"
	"simulator/service/src/modules/devices/domain/entities"
	"simulator/service/src/modules/devices/domain/repositories"
)

// mockRepo is an inline fake of the DeviceRepository port.
type mockRepo struct {
	createFn func(context.Context, *entities.Device) (*entities.Device, error)
	listFn   func(context.Context) ([]entities.Device, error)
	updateFn func(context.Context, string, *entities.Device) (*entities.Device, error)
	deleteFn func(context.Context, string) error
}

func (m *mockRepo) Create(ctx context.Context, d *entities.Device) (*entities.Device, error) {
	return m.createFn(ctx, d)
}
func (m *mockRepo) GetByID(context.Context, string) (*entities.Device, error) { return nil, nil }
func (m *mockRepo) List(ctx context.Context) ([]entities.Device, error)       { return m.listFn(ctx) }
func (m *mockRepo) Update(ctx context.Context, id string, d *entities.Device) (*entities.Device, error) {
	return m.updateFn(ctx, id, d)
}
func (m *mockRepo) Delete(ctx context.Context, id string) error { return m.deleteFn(ctx, id) }

var _ repositories.DeviceRepository = (*mockRepo)(nil)

func TestDevicesService_Create(t *testing.T) {
	now := time.Date(2026, 6, 9, 12, 0, 0, 0, time.UTC)
	repo := &mockRepo{createFn: func(_ context.Context, d *entities.Device) (*entities.Device, error) {
		// The service maps input to an entity WITHOUT presetting id/created; the
		// repository assigns them. Assert that before mutating the shared pointer.
		if d.ID != "" || !d.Created.IsZero() {
			t.Errorf("service preset id/created: %+v", d)
		}
		if d.Name != "n" || d.DeviceID != "d-ext" || d.Attributes["a"] != "b" {
			t.Errorf("input not mapped to entity: %+v", d)
		}
		d.ID = "dev-1"
		d.Created = now
		return d, nil
	}}
	svc := New(di.DevicesServiceDI{Repo: repo})

	in := &dtos.DeviceInput{
		Name: "n", DeviceID: "d-ext", ProtocolID: "http", Enabled: true,
		Attributes: map[string]string{"a": "b"},
		Config:     json.RawMessage(`{"kind":"http"}`),
	}
	out, err := svc.Create(context.Background(), in)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	// Stored entity maps back to the DTO, including created via NullTime.
	if out.ID != "dev-1" || out.Name != "n" || !out.Enabled {
		t.Fatalf("entity not mapped to dto: %+v", out)
	}
	if out.Created == nil || !out.Created.Time.Equal(now) {
		t.Fatalf("created not mapped: %+v", out.Created)
	}
}

func TestDevicesService_List(t *testing.T) {
	repo := &mockRepo{listFn: func(context.Context) ([]entities.Device, error) {
		return []entities.Device{{ID: "1", Name: "a"}, {ID: "2", Name: "b"}}, nil
	}}
	svc := New(di.DevicesServiceDI{Repo: repo})

	out, err := svc.List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(out) != 2 || out[0].ID != "1" || out[1].Name != "b" {
		t.Fatalf("list mapping: %+v", out)
	}
}

func TestDevicesService_DeleteMapsNotFoundToCustomError(t *testing.T) {
	repo := &mockRepo{deleteFn: func(context.Context, string) error {
		return repositories.ErrNotFound
	}}
	svc := New(di.DevicesServiceDI{Repo: repo})

	// The repository sentinel is translated into the HTTP-aware application error
	// so the global error handler can render a 404 envelope.
	_, err := svc.Delete(context.Background(), "missing")
	var customErr *customErrors.ServerCustomError
	if !errors.As(err, &customErr) {
		t.Fatalf("expected *ServerCustomError, got %v", err)
	}
	if customErr.Code != status.NOT_FOUND {
		t.Fatalf("expected code %d, got %d", status.NOT_FOUND, customErr.Code)
	}
}

package e2e

import (
	"net/url"
	"testing"
	"time"
)

// TestList_200 confirms the logs list endpoint answers.
func TestList_200(t *testing.T) {
	q := url.Values{}
	q.Set("limit", "1")
	_ = listLogs(t, q) // listLogs fails the test on a non-200.
}

// TestList_Filter_Device_Match returns only the seeded device's logs.
func TestList_Filter_Device_Match(t *testing.T) {
	deviceID := seedLogs(t, 3)
	q := url.Values{}
	q.Set("device", deviceID)
	q.Set("limit", "50")
	page := listLogs(t, q)
	if len(page.Items) < 3 {
		t.Fatalf("device filter returned %d logs, want >= 3", len(page.Items))
	}
	for _, l := range page.Items {
		if l.DeviceID != deviceID {
			t.Fatalf("device filter leaked a foreign log: %s", l.DeviceID)
		}
	}
}

// TestList_Filter_Device_NoMatch returns nothing for an unknown device.
func TestList_Filter_Device_NoMatch(t *testing.T) {
	q := url.Values{}
	q.Set("device", nonExistentDeviceID)
	q.Set("limit", "50")
	page := listLogs(t, q)
	if len(page.Items) != 0 {
		t.Fatalf("expected no logs for unknown device, got %d", len(page.Items))
	}
}

// TestList_Filter_Event_Match returns the seeded device's logs for its event.
func TestList_Filter_Event_Match(t *testing.T) {
	deviceID := seedLogs(t, 3)
	q := url.Values{}
	q.Set("device", deviceID)
	q.Set("event", "Telemetry")
	q.Set("limit", "50")
	page := listLogs(t, q)
	if len(page.Items) < 3 {
		t.Fatalf("event filter returned %d logs, want >= 3", len(page.Items))
	}
	for _, l := range page.Items {
		if l.EventName != "Telemetry" {
			t.Fatalf("event filter leaked a log with event %q", l.EventName)
		}
	}
}

// TestList_Filter_Event_NoMatch returns nothing for an unknown event name.
func TestList_Filter_Event_NoMatch(t *testing.T) {
	deviceID := seedLogs(t, 3)
	q := url.Values{}
	q.Set("device", deviceID)
	q.Set("event", "NoSuchEvent")
	q.Set("limit", "50")
	page := listLogs(t, q)
	if len(page.Items) != 0 {
		t.Fatalf("expected no logs for unknown event, got %d", len(page.Items))
	}
}

// TestList_Filter_DateRange_Match returns the logs when now is inside the range.
func TestList_Filter_DateRange_Match(t *testing.T) {
	deviceID := seedLogs(t, 3)
	now := time.Now().UTC()
	q := url.Values{}
	q.Set("device", deviceID)
	q.Set("dateFrom", now.Add(-time.Hour).Format(time.RFC3339))
	q.Set("dateTo", now.Add(time.Hour).Format(time.RFC3339))
	q.Set("limit", "50")
	page := listLogs(t, q)
	if len(page.Items) < 3 {
		t.Fatalf("date-range match returned %d logs, want >= 3", len(page.Items))
	}
}

// TestList_Filter_DateRange_NoMatch returns nothing for a past range.
func TestList_Filter_DateRange_NoMatch(t *testing.T) {
	deviceID := seedLogs(t, 3)
	q := url.Values{}
	q.Set("device", deviceID)
	q.Set("dateFrom", "2020-01-01T00:00:00Z")
	q.Set("dateTo", "2020-01-02T00:00:00Z")
	q.Set("limit", "50")
	page := listLogs(t, q)
	if len(page.Items) != 0 {
		t.Fatalf("expected no logs in a past range, got %d", len(page.Items))
	}
}

// TestList_Filter_Combined ANDs the device and event filters.
func TestList_Filter_Combined(t *testing.T) {
	deviceID := seedLogs(t, 3)
	q := url.Values{}
	q.Set("device", deviceID)
	q.Set("event", "Telemetry")
	q.Set("limit", "50")
	page := listLogs(t, q)
	if len(page.Items) < 3 {
		t.Fatalf("combined filter returned %d logs, want >= 3", len(page.Items))
	}
	for _, l := range page.Items {
		if l.DeviceID != deviceID || l.EventName != "Telemetry" {
			t.Fatalf("combined filter leaked a log: device=%s event=%s", l.DeviceID, l.EventName)
		}
	}
}

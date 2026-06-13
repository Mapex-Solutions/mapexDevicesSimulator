package logs

import (
	timeUtil "github.com/Mapex-Solutions/mapexGoKit/utils/time"
)

// Log is the wire shape of one persisted device message (the SQLite-backed
// history behind the live console stream). `created` is the message time;
// `payload` is what the device sent and `response` is what came back (the HTTP
// reply, a received downlink, or the failure reason).
type Log struct {
	ID         string             `json:"id"`
	Created    *timeUtil.NullTime `json:"created"`
	Protocol   string             `json:"protocol"`
	DeviceID   string             `json:"deviceId"`
	DeviceName string             `json:"deviceName"`
	EventName  string             `json:"eventName,omitempty"`
	Direction  string             `json:"direction"`
	Kind       string             `json:"kind"`
	Summary    string             `json:"summary"`
	Status     string             `json:"status,omitempty"`
	Payload    string             `json:"payload"`
	Response   string             `json:"response,omitempty"`
}

// SetCreated lets mapper.EntityToDto populate the message timestamp.
func (l *Log) SetCreated(t *timeUtil.NullTime) { l.Created = t }

// LogPage is the cursor-paginated response for GET /api/logs. NextCursor is the
// opaque token to pass back as `cursor` for the following page; it is empty when
// there are no more rows.
type LogPage struct {
	Items      []Log  `json:"items"`
	NextCursor string `json:"nextCursor,omitempty"`
}

// LogInput is the writable shape the simulation engine persists; id and created
// are assigned on insert.
type LogInput struct {
	Protocol   string `json:"protocol"`
	DeviceID   string `json:"deviceId"`
	DeviceName string `json:"deviceName"`
	EventName  string `json:"eventName,omitempty"`
	Direction  string `json:"direction"`
	Kind       string `json:"kind"`
	Summary    string `json:"summary"`
	Status     string `json:"status,omitempty"`
	Payload    string `json:"payload"`
	Response   string `json:"response,omitempty"`
}

// LogQuery binds the GET /api/logs query string. Empty filters are ignored.
// Cursor is the opaque keyset token from a previous page (empty for the first
// page). Q is a free-text match over summary, payload and device name; Event is
// a match over the event name; DateFrom/DateTo bound the message time.
type LogQuery struct {
	Limit     int    `query:"limit"`
	Cursor    string `query:"cursor"`
	Protocol  string `query:"protocol"`
	Kind      string `query:"kind"`
	Direction string `query:"direction"`
	Device    string `query:"device"`
	Event     string `query:"event"`
	DateFrom  string `query:"dateFrom"`
	DateTo    string `query:"dateTo"`
	Q         string `query:"q"`
}

package logs

import (
	timeUtil "github.com/Mapex-Solutions/mapexGoKit/utils/time"
)

// Log is the wire shape of one persisted device message (the SQLite-backed
// history behind the live console stream). `created` is the message time.
type Log struct {
	ID         string             `json:"id"`
	Created    *timeUtil.NullTime `json:"created"`
	Protocol   string             `json:"protocol"`
	DeviceID   string             `json:"deviceId"`
	DeviceName string             `json:"deviceName"`
	Direction  string             `json:"direction"`
	Kind       string             `json:"kind"`
	Summary    string             `json:"summary"`
	Status     string             `json:"status,omitempty"`
	Payload    string             `json:"payload"`
}

// SetCreated lets mapper.EntityToDto populate the message timestamp.
func (l *Log) SetCreated(t *timeUtil.NullTime) { l.Created = t }

// LogPage is the paginated response for GET /api/logs.
type LogPage struct {
	Items []Log `json:"items"`
	Total int   `json:"total"`
}

// LogInput is the writable shape the simulation engine persists; id and created
// are assigned on insert.
type LogInput struct {
	Protocol   string `json:"protocol"`
	DeviceID   string `json:"deviceId"`
	DeviceName string `json:"deviceName"`
	Direction  string `json:"direction"`
	Kind       string `json:"kind"`
	Summary    string `json:"summary"`
	Status     string `json:"status,omitempty"`
	Payload    string `json:"payload"`
}

// LogQuery binds the GET /api/logs query string. Empty filters are ignored; q is
// a free-text match over summary, payload and device name.
type LogQuery struct {
	Limit     int    `query:"limit"`
	Offset    int    `query:"offset"`
	Protocol  string `query:"protocol"`
	Kind      string `query:"kind"`
	Direction string `query:"direction"`
	Device    string `query:"device"`
	Q         string `query:"q"`
}

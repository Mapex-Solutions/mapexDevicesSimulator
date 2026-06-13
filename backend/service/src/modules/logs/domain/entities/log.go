package entities

import "time"

// Log is one persisted device message. All columns are scalar and queryable;
// created is the message time, stored newest-first for the console history.
type Log struct {
	ID         string    `db:"id,pk"`
	Protocol   string    `db:"protocol"`
	DeviceID   string    `db:"device_id"`
	DeviceName string    `db:"device_name"`
	EventName  string    `db:"event_name"`
	Direction  string    `db:"direction"`
	Kind       string    `db:"kind"`
	Summary    string    `db:"summary"`
	Status     string    `db:"status"`
	Payload    string    `db:"payload"`
	Response   string    `db:"response"`
	Created    time.Time `db:"created"`
}

// GetCreated exposes the message time for the entity-to-DTO mapper.
func (l *Log) GetCreated() time.Time { return l.Created }

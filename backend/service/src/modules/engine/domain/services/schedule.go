package services

import (
	"time"

	"simulator/service/src/modules/engine/domain/entities"
)

// ScheduleInterval converts an event schedule to its tick interval. It returns 0
// when the event must not auto-fire: no schedule, disabled, a non-positive
// period, or an unknown unit.
func ScheduleInterval(s *entities.EventSchedule) time.Duration {
	if s == nil || !s.Enabled || s.Every <= 0 {
		return 0
	}
	switch s.Unit {
	case "seconds":
		return time.Duration(s.Every) * time.Second
	case "minutes":
		return time.Duration(s.Every) * time.Minute
	case "hours":
		return time.Duration(s.Every) * time.Hour
	case "days":
		return time.Duration(s.Every) * 24 * time.Hour
	default:
		return 0
	}
}

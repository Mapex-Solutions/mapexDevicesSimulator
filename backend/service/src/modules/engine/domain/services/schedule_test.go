package services

import (
	"testing"
	"time"

	"simulator/service/src/modules/engine/domain/entities"
)

func TestScheduleInterval(t *testing.T) {
	tests := []struct {
		name string
		in   *entities.EventSchedule
		want time.Duration
	}{
		{"nil schedule", nil, 0},
		{"disabled", &entities.EventSchedule{Enabled: false, Every: 5, Unit: "seconds"}, 0},
		{"zero every", &entities.EventSchedule{Enabled: true, Every: 0, Unit: "seconds"}, 0},
		{"unknown unit", &entities.EventSchedule{Enabled: true, Every: 1, Unit: "weeks"}, 0},
		{"seconds", &entities.EventSchedule{Enabled: true, Every: 5, Unit: "seconds"}, 5 * time.Second},
		{"minutes", &entities.EventSchedule{Enabled: true, Every: 2, Unit: "minutes"}, 2 * time.Minute},
		{"hours", &entities.EventSchedule{Enabled: true, Every: 1, Unit: "hours"}, time.Hour},
		{"days", &entities.EventSchedule{Enabled: true, Every: 1, Unit: "days"}, 24 * time.Hour},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ScheduleInterval(tt.in); got != tt.want {
				t.Fatalf("ScheduleInterval = %v, want %v", got, tt.want)
			}
		})
	}
}

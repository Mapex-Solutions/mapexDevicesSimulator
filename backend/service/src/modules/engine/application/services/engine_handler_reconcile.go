package services

import (
	"container/heap"
	"context"
	"fmt"
	"time"

	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	domainsvc "simulator/service/src/modules/engine/domain/services"
)

// reconcile reads the devices, computes the desired job set, and aligns the heap:
// remove jobs that disappeared, add new ones, and replace ones whose interval or
// target changed. It wakes the scheduler so the next wait reflects the change.
func (s *EngineService) reconcile() {
	desired := s.buildDesired()

	s.mu.Lock()
	now := time.Now()
	for key, j := range s.jobs {
		if _, ok := desired[key]; !ok {
			s.heap.remove(j)
			delete(s.jobs, key)
		}
	}
	for key, d := range desired {
		if existing, ok := s.jobs[key]; ok {
			if existing.sig != d.sig {
				existing.interval = d.interval
				existing.spec = d.spec
				existing.sig = d.sig
				existing.nextFireAt = now.Add(d.interval)
				s.heap.fix(existing)
			}
			continue
		}
		d.nextFireAt = now.Add(d.interval)
		heap.Push(s.heap, d)
		s.jobs[key] = d
	}
	s.mu.Unlock()

	select {
	case s.wake <- struct{}{}:
	default:
	}
}

// buildDesired derives the desired job set from the enabled devices and their
// enabled-schedule events.
func (s *EngineService) buildDesired() map[string]*job {
	desired := make(map[string]*job)
	devices, err := s.deps.Devices.List(context.Background())
	if err != nil {
		logger.Error(err, "[SERVICE:Engine] list devices for reconcile")
		return desired
	}
	for _, d := range devices {
		if !d.Enabled {
			continue
		}
		// A device whose gateway is offline keeps its scheduled jobs: each fire still
		// runs so the console reports the attempt (with a gateway-offline status); the
		// dispatcher is what skips the actual uplink when there is no live link.
		events, err := domainsvc.ParseEvents(d.Events)
		if err != nil {
			continue
		}
		for _, e := range events {
			interval := domainsvc.ScheduleInterval(e.Schedule)
			if interval <= 0 {
				continue
			}
			spec, ok := buildSendSpec(d, e)
			if !ok {
				continue
			}
			key := d.ID + "|" + e.ID
			desired[key] = &job{key: key, interval: interval, spec: spec, sig: jobSignature(interval, spec)}
		}
	}
	return desired
}

// jobSignature captures the fields whose change requires rescheduling a job.
func jobSignature(interval time.Duration, spec sendSpec) string {
	return fmt.Sprintf("%d|%s|%s|%s|%s|%s|%d|%t|%t",
		interval, spec.url, spec.method, spec.brokerURL, spec.topic,
		spec.payloadTemplate, spec.qos, spec.retain, spec.storeLogs)
}

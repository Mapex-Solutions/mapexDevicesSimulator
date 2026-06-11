package services

import (
	"context"

	logger "github.com/Mapex-Solutions/mapexGoKit/microservices/logger"

	"simulator/service/src/modules/engine/application/di"
	"simulator/service/src/modules/engine/application/ports"
)

// Compile-time check that the service satisfies its port.
var _ ports.EnginePort = (*EngineService)(nil)

// New builds the engine service. It does not start anything; OnMount does.
func New(deps di.EngineServiceDI) ports.EnginePort {
	return &EngineService{
		deps:     deps,
		heap:     &jobHeap{},
		jobs:     make(map[string]*job),
		sessions: make(map[string]*sessionHandle),
		fires:    make(chan fireTask, fireBuffer),
		wake:     make(chan struct{}, 1),
	}
}

// OnMount reads the devices, builds the initial job set, and starts the worker
// pool, the scheduler, and the slow safety resync. Fired during module init.
func (s *EngineService) OnMount() {
	if !s.markStarted() {
		return
	}
	s.reconcile()
	s.reconcileSessions()
	s.startWorkers()
	s.startScheduler()
	s.startResync()
	logger.Info("[SERVICE:Engine] mounted, simulation running")
}

// Reconcile re-reads the devices and re-aligns the running jobs. Called on a CRUD
// change signal and by the slow resync.
func (s *EngineService) Reconcile() {
	if !s.isStarted() {
		return
	}
	s.reconcile()
	s.reconcileSessions()
}

// Fire sends one event on demand for a device (the console/REST "fire" action): it
// resolves the event (a pre-registered one by id, or an inline ad-hoc event), builds
// the send spec, and runs it through the same process() path as a scheduled fire, so
// the uplink goes through the live session when one exists and the result streams to
// the console (and logs, when storeLogs).
func (s *EngineService) Fire(ctx context.Context, deviceID string, in ports.FireInput) error {
	dev, err := s.findDevice(ctx, deviceID)
	if err != nil {
		return err
	}
	event, err := resolveFireEvent(*dev, in)
	if err != nil {
		return err
	}
	spec, ok := buildSendSpec(*dev, event)
	if !ok {
		return ports.ErrFireUnsupported
	}
	s.process(fireTask{spec: spec, counter: s.fireSeq.Add(1)})
	return nil
}

// OnShutdown cancels the scheduler + workers and waits for them to drain.
func (s *EngineService) OnShutdown(_ context.Context) error {
	cancel := s.stopStarted()
	if cancel == nil {
		return nil
	}
	cancel()
	s.wg.Wait()
	logger.Info("[SERVICE:Engine] stopped")
	return nil
}

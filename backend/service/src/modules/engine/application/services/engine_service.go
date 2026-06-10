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
		deps:  deps,
		heap:  &jobHeap{},
		jobs:  make(map[string]*job),
		fires: make(chan fireTask, fireBuffer),
		wake:  make(chan struct{}, 1),
	}
}

// OnMount reads the devices, builds the initial job set, and starts the worker
// pool, the scheduler, and the slow safety resync. Fired during module init.
func (s *EngineService) OnMount() {
	if !s.markStarted() {
		return
	}
	s.reconcile()
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

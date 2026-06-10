package services

import (
	"container/heap"
	"context"
	"time"
)

// markStarted flips the engine to started under the lock, creating the lifecycle
// context. It returns false if already started (OnMount is idempotent).
func (s *EngineService) markStarted() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.started {
		return false
	}
	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.started = true
	return true
}

// isStarted reports whether the engine is running.
func (s *EngineService) isStarted() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.started
}

// stopStarted flips the engine to stopped and returns its cancel func, or nil if
// it was not running.
func (s *EngineService) stopStarted() context.CancelFunc {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.started {
		return nil
	}
	s.started = false
	return s.cancel
}

// startWorkers launches the bounded worker pool.
func (s *EngineService) startWorkers() {
	for i := 0; i < workerCount; i++ {
		s.wg.Add(1)
		go s.worker()
	}
}

// startScheduler launches the single scheduler goroutine.
func (s *EngineService) startScheduler() {
	s.wg.Add(1)
	go s.schedulerLoop()
}

// startResync launches the slow safety resync goroutine.
func (s *EngineService) startResync() {
	s.wg.Add(1)
	go s.resyncLoop()
}

// schedulerLoop sleeps until the earliest job is due, fires the due jobs, and
// repeats. A reconcile signal wakes it early to re-evaluate the next wait.
func (s *EngineService) schedulerLoop() {
	defer s.wg.Done()
	for {
		timer := time.NewTimer(s.nextWait())
		select {
		case <-s.ctx.Done():
			timer.Stop()
			return
		case <-s.wake:
			timer.Stop()
		case <-timer.C:
			s.fireDue()
		}
	}
}

// resyncLoop re-reads the DB on a slow interval, catching any missed change.
func (s *EngineService) resyncLoop() {
	defer s.wg.Done()
	t := time.NewTicker(resyncEvery)
	defer t.Stop()
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-t.C:
			s.reconcile()
			s.reconcileSessions()
		}
	}
}

// nextWait returns how long to sleep until the next job is due (or idleWait when
// the heap is empty).
func (s *EngineService) nextWait() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.heap.Len() == 0 {
		return idleWait
	}
	if d := time.Until(s.heap.peek().nextFireAt); d > 0 {
		return d
	}
	return 0
}

// fireDue pops every job whose time has come, queues it for the workers, and
// reschedules it. A full queue drops the firing rather than stalling.
func (s *EngineService) fireDue() {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	for s.heap.Len() > 0 && !s.heap.peek().nextFireAt.After(now) {
		j := heap.Pop(s.heap).(*job)
		j.counter++
		select {
		case s.fires <- fireTask{spec: j.spec, counter: j.counter}:
		default:
		}
		j.nextFireAt = now.Add(j.interval)
		heap.Push(s.heap, j)
	}
}

// worker pulls fire tasks and processes them until the engine stops.
func (s *EngineService) worker() {
	defer s.wg.Done()
	for {
		select {
		case <-s.ctx.Done():
			return
		case task := <-s.fires:
			s.process(task)
		}
	}
}

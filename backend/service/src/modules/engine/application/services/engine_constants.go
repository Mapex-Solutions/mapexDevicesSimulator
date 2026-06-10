package services

import "time"

const (
	// workerCount bounds the goroutines that perform the actual sends, regardless
	// of how many jobs are scheduled.
	workerCount = 16

	// fireBuffer is the queue depth between the scheduler and the workers; a full
	// queue drops the firing rather than stalling the scheduler.
	fireBuffer = 1024

	// idleWait is how long the scheduler sleeps when no job is scheduled.
	idleWait = 30 * time.Second

	// resyncEvery is the slow safety re-read of the DB, in case a change signal is
	// ever missed.
	resyncEvery = 5 * time.Second
)

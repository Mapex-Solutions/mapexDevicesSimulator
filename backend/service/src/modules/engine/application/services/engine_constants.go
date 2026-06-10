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

	// sessionBackoffInitial is the floor reconnect delay; an offline broker never
	// reconnects faster than this, so it cannot become a reconnect storm.
	sessionBackoffInitial = 1 * time.Second

	// sessionBackoffMax caps the reconnect delay so it stabilizes rather than
	// growing unboundedly.
	sessionBackoffMax = 30 * time.Second

	// sessionPollEvery is how often a live session is checked for a silent drop.
	sessionPollEvery = 2 * time.Second
)

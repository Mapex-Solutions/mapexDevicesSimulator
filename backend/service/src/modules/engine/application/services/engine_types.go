package services

import (
	"container/heap"
	"context"
	"sync"
	"time"

	"simulator/service/src/modules/engine/application/di"
)

// EngineService runs the simulation: it keeps a min-heap of jobs (one per
// enabled device + scheduled event), a scheduler goroutine that fires the due
// ones, and a bounded worker pool that renders and dispatches them.
type EngineService struct {
	deps di.EngineServiceDI

	mu      sync.Mutex
	heap    *jobHeap
	jobs    map[string]*job // key -> running job
	started bool

	fires  chan fireTask
	wake   chan struct{}
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// sendSpec is the resolved "what to send" for one job: the protocol target plus
// the payload template (rendered at fire time).
type sendSpec struct {
	protocol   string
	deviceID   string
	deviceName string
	storeLogs  bool

	// http
	url     string
	method  string
	headers map[string]string

	// mqtt
	brokerURL string
	clientID  string
	username  string
	password  string
	topic     string
	qos       byte
	retain    bool

	payloadTemplate string
}

// job is one scheduled event in the heap.
type job struct {
	key        string
	interval   time.Duration
	nextFireAt time.Time
	counter    int64
	index      int // position in the heap, for Remove/Fix
	sig        string
	spec       sendSpec
}

// fireTask is what the scheduler hands a worker.
type fireTask struct {
	spec    sendSpec
	counter int64
}

// jobHeap is a min-heap of jobs ordered by next fire time.
type jobHeap struct {
	items []*job
}

func (h *jobHeap) Len() int { return len(h.items) }

func (h *jobHeap) Less(i, j int) bool {
	return h.items[i].nextFireAt.Before(h.items[j].nextFireAt)
}

func (h *jobHeap) Swap(i, j int) {
	h.items[i], h.items[j] = h.items[j], h.items[i]
	h.items[i].index = i
	h.items[j].index = j
}

func (h *jobHeap) Push(x any) {
	j := x.(*job)
	j.index = len(h.items)
	h.items = append(h.items, j)
}

func (h *jobHeap) Pop() any {
	n := len(h.items)
	j := h.items[n-1]
	h.items[n-1] = nil
	j.index = -1
	h.items = h.items[:n-1]
	return j
}

// peek returns the earliest job without removing it.
func (h *jobHeap) peek() *job { return h.items[0] }

// remove drops a job from the heap by its tracked index.
func (h *jobHeap) remove(j *job) {
	if j.index >= 0 && j.index < len(h.items) {
		heap.Remove(h, j.index)
	}
}

// fix re-establishes heap order after a job's next fire time changed.
func (h *jobHeap) fix(j *job) {
	if j.index >= 0 && j.index < len(h.items) {
		heap.Fix(h, j.index)
	}
}

// Package reconcile carries an in-process signal that device and gateway writes
// raise to ask the engine to re-align immediately, instead of waiting for its
// slow safety resync. It is deliberately neutral: the engine already reads
// devices and gateways, so a direct dependency back from those modules to the
// engine would be an import cycle. Both sides depend only on this leaf package.
package reconcile

import "sync"

// Signal is the writer side, consumed by the devices and gateways services to
// announce that their stored set changed.
type Signal interface {
	Notify()
}

// Listener is the subscriber side, consumed by the engine to register the
// callback run on every change.
type Listener interface {
	Subscribe(fn func())
}

// Notifier is the shared singleton implementing both sides. It is safe for
// concurrent use: writes fan out under a read lock and subscriptions append
// under a write lock.
type Notifier struct {
	mu        sync.RWMutex
	listeners []func()
}

// New builds an empty notifier.
func New() *Notifier { return &Notifier{} }

// Subscribe registers a callback invoked on every Notify. The engine subscribes
// its Reconcile once at mount.
func (n *Notifier) Subscribe(fn func()) {
	n.mu.Lock()
	n.listeners = append(n.listeners, fn)
	n.mu.Unlock()
}

// Notify fans the change out to every listener. Each runs in its own goroutine
// so a slow reconcile never blocks the HTTP write that raised the signal.
func (n *Notifier) Notify() {
	n.mu.RLock()
	defer n.mu.RUnlock()
	for _, fn := range n.listeners {
		go fn()
	}
}

package undead

import (
	"runtime"
	"sync"

	"github.com/petermattis/goid"
)

// type to implement a lock
type Mutex struct {
	mu sync.Mutex
}

// Lock mutex m
func (m *Mutex) Lock() {
	defer m.mu.Lock()

	// if detection is disabled
	if !Opts.RunDetection {
		return
	}

	// update data structures if more than on routine is running
	if runtime.NumGoroutine() > 1 {
		r := routines[goid.Get()]
		r.updateLock(m)
	}

}

// Unlock mutex m
func (m *Mutex) Unlock() {
	defer m.mu.Unlock()

}

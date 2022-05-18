package undead

import (
	"runtime"
	"sync"

	"github.com/petermattis/goid"
)

// type to implement a lock
type mutex struct {
	mu      sync.Mutex
	context callerInfo // info about the mutex initialization
}

// create Lock
func NewLock() (m mutex) {
	_, file, line, _ := runtime.Caller(2)
	m.context = newInfo(file, line)
	return m
}

// Lock mutex m
func (m *mutex) Lock() {
	defer m.mu.Lock()

	// if detection is disabled
	if !Opts.RunDetection {
		return
	}

	// TODO: avoid recursive intercepting

	// update data structures if more than on routine is running
	if runtime.NumGoroutine() > 1 {
		r := routines[goid.Get()]
		r.updateLock(m)
	}

}

// Unlock mutex m
func (m *mutex) Unlock() {
	defer m.mu.Unlock()

	r := routines[goid.Get()]
	r.updateUnlock(m)

}

package deadlock

import (
	"runtime"
	"sync"
)

// type to implement a lock
type mutex struct {
	mu      sync.Mutex
	context []callerInfo // info about the mutex initialization
}

// create Lock
func NewLock() (m mutex) {
	_, file, line, _ := runtime.Caller(1)
	m.context = append(m.context, newInfo(file, line))
	return m
}

// Lock mutex m
func (m *mutex) Lock() {
	defer m.mu.Lock()

	// if detection is disabled
	if !Opts.RunDetection {
		return
	}

	index := getRoutineIndex()
	r := &routines[index]

	// update data structures if more than on routine is running
	if runtime.NumGoroutine() > 1 {
		(*r).updateLock(m)
	}

}

// Unlock mutex m
func (m *mutex) Unlock() {
	defer m.mu.Unlock()

	index := getRoutineIndex()

	r := routines[index]
	r.updateUnlock(m)

}

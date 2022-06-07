package deadlock

/*
Author: Erik Kassubek <erik-kassubek@t-online.de>
Package: deadlock
Project: Bachelor Project at the Albert-Ludwigs-University Freiburg,
	Institute of Computer Science: Dynamic Deadlock Detection in Go
Date: 2022-06-05
*/

/*
mutex.go
This file implements the drop-in-replacement for the locks (mutexes) as well as
the lock and unlock operations for these locks.
*/

import (
	"runtime"
	"sync"
)

// type to implement a lock
type mutex struct {
	mu      sync.Mutex
	context []callerInfo // info about the creation and lock/unlock of this lock
}

// create Lock
func NewLock() (m mutex) {
	_, file, line, _ := runtime.Caller(1)
	m.context = append(m.context, newInfo(file, line, true))
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

// Trylock mutex m
func (m *mutex) TryLock() bool {
	res := m.mu.TryLock()

	if !Opts.RunDetection {
		return res
	}

	index := getRoutineIndex()
	r := &routines[index]

	// update data structures if more than on routine is running
	if runtime.NumGoroutine() > 1 {
		if res {
			(*r).updateTryLock(m)
		}
	}
	return res
}

// Unlock mutex m
func (m *mutex) Unlock() {
	defer m.mu.Unlock()

	index := getRoutineIndex()

	r := routines[index]
	r.updateUnlock(m)

}

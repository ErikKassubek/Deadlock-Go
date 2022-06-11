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
	"fmt"
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
	var bufString string
	if opts.collectCallStack {
		buf := make([]byte, opts.maxCallStackSize)
		n := runtime.Stack(buf[:], false)
		bufString = string(buf[:n])
	}
	m.context = append(m.context, newInfo(file, line, true, bufString))
	return m
}

// Lock mutex m
func (m *mutex) Lock() {
	defer m.mu.Lock()

	// if detection is disabled
	if !opts.periodicDetection && !opts.comprehensiveDetection {
		return
	}

	index := getRoutineIndex()

	if index >= routinesIndex {
		panic(`A Routine  was not initialized. Run NewRoutine() before Lock or TryLock operation`)
	}

	r := &routines[index]

	// update data structures if more than on routine is running
	if runtime.NumGoroutine() > 1 {
		(*r).updateLock(m)
	}
}

// Trylock mutex m
func (m *mutex) TryLock() bool {
	res := m.mu.TryLock()

	if !opts.periodicDetection && !opts.comprehensiveDetection {
		return res
	}

	index := getRoutineIndex()

	if index >= routinesIndex {
		errorString := fmt.Sprintf(`Routine %d was not initialized. Run 
			NewRoutine() in the corresponding routine before Lock or TryLock 
			operation`, index)
		panic(errorString)
	}

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

	if !opts.periodicDetection && !opts.comprehensiveDetection {
		return
	}

	index := getRoutineIndex()

	r := &routines[index]
	(*r).updateUnlock(m)
}

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
type Mutex struct {
	mu       sync.Mutex
	context  []callerInfo // info about the creation and lock/unlock of this lock
	in       bool         // set to true after lock was initialized
	isLocked bool         // set to true if lock is locked
}

// create Lock
func NewLock() (m Mutex) {
	// initialize detector if necessary
	if !initialized {
		initialize()
	}

	_, file, line, _ := runtime.Caller(1)
	var bufString string
	if opts.collectCallStack {
		buf := make([]byte, opts.maxCallStackSize)
		n := runtime.Stack(buf[:], false)
		bufString = string(buf[:n])
	}
	m.context = append(m.context, newInfo(file, line, true, bufString))
	m.in = true

	return m
}

// Lock mutex m
func (m *Mutex) Lock() {
	if !m.in {
		errorMessage := fmt.Sprint("Lock ", &m, " was not created. Use ",
			"x := NewLock()")
		panic(errorMessage)
	}

	// initialize detector if necessary
	if !initialized {
		initialize()
	}

	defer func() {
		m.mu.Lock()
		m.isLocked = true
	}()

	// if detection is disabled
	if !opts.periodicDetection && !opts.comprehensiveDetection {
		return
	}

	index := getRoutineIndex()
	if index == -1 {
		// create new routine, if not initialized
		newRoutine()
	}
	index = getRoutineIndex()

	r := &routines[index]

	// check for double locking
	if opts.checkDoubleLocking && m.isLocked {
		r.checkDoubleLocking(m)
	}

	numRoutine := runtime.NumGoroutine()
	// update data structures if more than on routine is running
	if numRoutine > 1 {
		(*r).updateLock(m)
	}

}

// Trylock mutex m
func (m *Mutex) TryLock() bool {
	if !m.in {
		errorMessage := fmt.Sprint("Lock ", &m, " was not created. Use ",
			"x := NewLock()")
		panic(errorMessage)
	}

	// initialize detector if necessary
	if !initialized {
		initialize()
	}

	res := m.mu.TryLock()

	if res {
		m.isLocked = true
	}

	if !opts.periodicDetection && !opts.comprehensiveDetection {
		return res
	}

	index := getRoutineIndex()
	if index == -1 {
		// create new routine, if not initialized
		newRoutine()
	}

	r := &routines[index]
	if res && opts.checkDoubleLocking {
		r.checkDoubleLocking(m)
	}

	// update data structures if more than on routine is running
	if runtime.NumGoroutine() > 1 {
		if res {
			(*r).updateTryLock(m)
		}
	}

	return res
}

// Unlock mutex m
func (m *Mutex) Unlock() {
	if !m.isLocked {
		errorMessage := fmt.Sprint("Tried to unLock lock", &m,
			" which was not locked.")
		panic(errorMessage)
	}
	defer func() {
		m.mu.Unlock()
		m.isLocked = false
	}()

	if !opts.periodicDetection && !opts.comprehensiveDetection {
		return
	}

	index := getRoutineIndex()

	r := &routines[index]
	(*r).updateUnlock(m)
}

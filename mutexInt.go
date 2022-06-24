package deadlock

import (
	"fmt"
	"runtime"
	"sync"
)

/*
Copyright (C) 2022  Erik Kassubek

  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.

  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

/*
Author: Erik Kassubek <erik-kassubek@t-online.de>
Package: deadlock
Project: Bachelor Project at the Albert-Ludwigs-University Freiburg,
	Institute of Computer Science: Dynamic Deadlock Detection in Go
*/

/*
mutexInt.go
This file implements and interface for Mutex and RWMutex.
It also implements code which is used for both mutex and rwmutex
*/

type mutexInt interface {
	getNumberLocked() *int
	getIsLockedRoutineIndex() *int
	getContext() *[]callerInfo
	getMemoryPosition() uintptr
	getIn() *bool
	getLock() (bool, *sync.Mutex, *sync.RWMutex) // if bool is true, mutex is returned,
	// otherwise RWMutex is returned
	getIsRead() *bool
	isRWLock() bool
}

// lock the mutex or rwmutex and update the detector data
func lockInt(m mutexInt, rLock bool) {
	if !*m.getIn() {
		errorMessage := fmt.Sprint("Lock ", &m, " was not created. Use ",
			"x := NewLock()")
		panic(errorMessage)
	}

	// initialize detector if necessary
	if !initialized {
		initialize()
	}

	defer func() {
		d, l, t := m.getLock()
		if d {
			l.Lock()
		} else {
			if rLock {
				t.RLock()
			} else {
				t.Lock()
			}
		}

		*m.getNumberLocked() += 1
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
	if opts.checkDoubleLocking && *m.getNumberLocked() != 0 {
		r.checkDoubleLocking(m, index, rLock)
	}

	*m.getIsLockedRoutineIndex() = index

	numRoutine := runtime.NumGoroutine()
	// update data structures if more than on routine is running
	if numRoutine > 1 {
		(*r).updateLock(m)
	}
}

// unlock the mutex or rwmutex and update the detector data
func unlockInt(m mutexInt) {
	if *m.getNumberLocked() == 0 {
		errorMessage := fmt.Sprint("Tried to unLock lock ", &m,
			" which was not locked.")
		panic(errorMessage)
	}
	defer func() {
		d, l, r := m.getLock()
		if d {
			l.Unlock()
		} else {
			if *m.getIsRead() {
				r.RUnlock()
			} else {
				r.Unlock()
			}
		}
		*m.getIsLockedRoutineIndex() = -1
		*m.getNumberLocked() -= 1
	}()

	if !opts.periodicDetection && !opts.comprehensiveDetection {
		return
	}

	index := getRoutineIndex()

	r := &routines[index]
	(*r).updateUnlock(m)
}

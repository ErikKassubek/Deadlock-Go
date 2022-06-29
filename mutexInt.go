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
It also implements code which is used for both mutex and rw-mutex
*/

// creat and interface for Mutex and RWMutex
type mutexInt interface {
	getNumberLocked() *int
	getIsLockedRoutineIndex() *int
	getContext() *[]callerInfo
	getMemoryPosition() uintptr
	getIn() *bool
	getLock() (bool, *sync.Mutex, *sync.RWMutex)
	getIsRead() *bool
}

// lock the mutex or rw-mutex and update the detector data
//  Args:
//   m (mutexInt): mutex or rw-mutex to lock
//   rLock (bool): if set to true, the lock is a reader lock
//  Returns:
//   nil
func lockInt(m mutexInt, rLock bool) {
	// panic if the lock was not initialized
	if !*m.getIn() {
		errorMessage := fmt.Sprint("Lock ", &m, " was not created. Use ",
			"x := NewLock()")
		panic(errorMessage)
	}

	// initialize detector if necessary
	if !initialized {
		initialize()
	}

	// defer the actual locking
	defer func() {
		d, l, t := m.getLock()
		if d {
			// lock if m is mutex
			l.Lock()
		} else {
			// lock if m is rw-mutex
			if rLock {
				t.RLock()
			} else {
				t.Lock()
			}
		}

		*m.getNumberLocked() += 1
	}()

	// return if detection is disabled
	if !opts.periodicDetection && !opts.comprehensiveDetection {
		return
	}

	// create new routine, if not initialized
	index := getRoutineIndex()
	if index == -1 {
		newRoutine()
	}
	index = getRoutineIndex()

	r := &routines[index]

	// check if the locking would lead to double locking
	if opts.checkDoubleLocking && *m.getNumberLocked() != 0 {
		r.checkDoubleLocking(m, index, rLock)
	}

	*m.getIsLockedRoutineIndex() = index

	// update data structures if more than on routine is running
	numRoutine := runtime.NumGoroutine()
	if numRoutine > 1 {
		(*r).updateLock(m)
	}
}

// unlock the mutex or rw-mutex and update the detector data
// Args:
//  m (mutexInt): mutex or RWMutex to unlock
// Returns:
//  nil
func unlockInt(m mutexInt) {
	// panic if the lock was not initialized
	if !*m.getIn() {
		errorMessage := fmt.Sprint("Lock ", &m, " was not created. Use ",
			"x := NewLock()")
		panic(errorMessage)
	}

	// panic if lock was not locked
	if *m.getNumberLocked() == 0 {
		errorMessage := fmt.Sprint("Tried to unLock lock ", &m,
			" which was not locked.")
		panic(errorMessage)
	}

	// defer the actual unlocking
	defer func() {
		d, l, r := m.getLock()
		if d {
			// unlock if m is mutex
			l.Unlock()
		} else {
			// unlock if m is rw-mutex
			if *m.getIsRead() {
				r.RUnlock()
			} else {
				r.Unlock()
			}
		}

		// update numberLocked and isLockedRoutineIndex
		*m.getNumberLocked() -= 1
		if *m.getNumberLocked() == 0 {
			*m.getIsLockedRoutineIndex() = -1
		}
	}()

	// return if detection is disabled
	if !opts.periodicDetection && !opts.comprehensiveDetection {
		return
	}

	// update data structures if more than on routine is running
	index := getRoutineIndex()
	r := &routines[index]
	(*r).updateUnlock(m)
}

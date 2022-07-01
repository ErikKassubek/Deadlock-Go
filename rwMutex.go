package deadlock

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
rwMutex.go
This file implements the drop-in-replacement for the rw-locks (rw-mutexes) as
well as
the lock, rlock and unlock operations for these locks.
*/

import (
	"runtime"
	"sync"
	"unsafe"
)

// type to implement a lock
type RWMutex struct {
	// rw-mutex for the actual locking
	mu *sync.RWMutex
	// info about the creation and lock/unlock of this lock
	context []callerInfo
	// set to true after lock was initialized
	in bool
	// how ofter is the lock locked
	numberLocked int
	// index of the routine, which holds the lock
	isLockedRoutineIndex int
	// position of the mutex in memory
	memoryPosition uintptr
	// set true, if last acquisition was RLock
	isRead bool // set true, if last acquisition was RLock
}

// create a new rw-lock
func NewRWLock() *RWMutex {
	// initialize detector if necessary
	if !initialized {
		initialize()
	}

	m := RWMutex{
		mu:                   &sync.RWMutex{},
		in:                   true,
		isLockedRoutineIndex: -1,
	}

	// save the position of the NewLock call
	_, file, line, _ := runtime.Caller(1)
	m.context = append(m.context, newInfo(file, line, true, ""))

	// save the memory position of the mutex
	m.memoryPosition = uintptr(unsafe.Pointer(&m))

	return &m
}

// ====== GETTER ===============================================================

// getter for isLocked
//  Returns:
//   (*int): numberLocked
func (m *RWMutex) getNumberLocked() *int {
	return &m.numberLocked
}

// getter for isLockedRoutineIndex
//  Returns:
//   (*int): isLockedRoutineIndex
func (m *RWMutex) getIsLockedRoutineIndex() *int {
	return &m.isLockedRoutineIndex
}

// getter for context
//  Returns:
//   (*[]callerInfo): caller info of the lock
func (m *RWMutex) getContext() *[]callerInfo {
	return &m.context
}

// getter for memoryPosition
//  Returns:
//   (uintptr): memoryPosition
func (m *RWMutex) getMemoryPosition() uintptr {
	return m.memoryPosition
}

// getter for in
//  Returns:
//   (bool): true if the lock was initialized, false otherwise
func (m *RWMutex) getIn() *bool {
	return &m.in
}

// getter for mu
//  Returns:
//   (bool): false, true for mutex
//   (*sync.Mutex): nil, underlying sync.Mutex mu for mutex
//   (*sync.RWMutex): nil, underlying sync.RWMutex mu
func (m *RWMutex) getLock() (bool, *sync.Mutex, *sync.RWMutex) {
	return false, nil, m.mu
}

// getter for isRead
//  Returns:
//   (*bool): true, if the last acquisition of the lock was rlock, false otherwise
func (m *RWMutex) getIsRead() *bool {
	return &m.isRead
}

// ====== FUNCTIONS ============================================================

// Lock rwmutex m
//  Returns:
//   nil
func (m *RWMutex) Lock() {
	// call the lock method for the mutexInt interface
	lockInt(m, false)
	m.isRead = false
}

// RLock rwmutex m
//  Returns:
//   nil
func (m *RWMutex) RLock() {
	// call the try-lock method for the mutexInt interface
	lockInt(m, true)
	m.isRead = true
}

// TryLock rw-mutex m
//  Returns:
//   (bool): true if locking was successful, false otherwise
func (m *RWMutex) TryLock() bool {
	// call the try-lock method for the mutexInt interface
	res := tryLockInt(m, false)
	if res {
		m.isRead = false
	}
	return res
}

// TryRLock rw-mutex m
//  Returns:
//   (bool): true if locking was successful, false otherwise
func (m *RWMutex) TryRLock() bool {
	// call the trylock method for the mutexInt interface
	res := tryLockInt(m, true)
	if res {
		m.isRead = false
	}
	return res
}

// Unlock rwmutex m
func (m *RWMutex) Unlock() {
	unlockInt(m)
}

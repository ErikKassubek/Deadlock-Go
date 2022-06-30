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
mutex.go
This file implements the drop-in-replacement for the locks (mutexes) as well as
the lock and unlock operations for these locks.
*/

import (
	"runtime"
	"sync"
	"unsafe"
)

// Type to implement a lock
// It can be used as an drop in replacement
// TODO: check if this can be lowercase
type Mutex struct {
	// mutex for the actual locking
	mu *sync.Mutex
	// info about the creation and lock/unlock of this lock
	context []callerInfo
	// set to true after lock was initialized
	in bool
	// numberLocked stores how often the mutex is currently locked
	numberLocked int
	// index of the routine, which holds the lock
	isLockedRoutineIndex int
	// position of the mutex in memory
	memoryPosition uintptr
}

// create and return a new lock, which can be used as a drop-in replacement for
// sync.Mutex
func NewLock() *Mutex {
	// initialize detector if necessary
	if !initialized {
		initialize()
	}

	m := Mutex{
		mu:                   &sync.Mutex{},
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

// ============ GETTER ============

// getter for isLocked
//  Returns:
//   (*int): numberLocked
func (m *Mutex) getNumberLocked() *int {
	return &m.numberLocked
}

// getter for isLockedRoutineIndex
//  Returns:
//   (*int): isLockedRoutineIndex
func (m *Mutex) getIsLockedRoutineIndex() *int {
	return &m.isLockedRoutineIndex
}

// getter for context
//  Returns:
//   (*[]callerInfo): caller info of the lock
func (m *Mutex) getContext() *[]callerInfo {
	return &m.context
}

// getter for memoryPosition
// Returns:
//  (uintptr): memoryPosition
func (m *Mutex) getMemoryPosition() uintptr {
	return m.memoryPosition
}

// getter for in
//  Returns:
//   (bool): true if the lock was initialized, false otherwise
func (m *Mutex) getIn() *bool {
	return &m.in
}

// getter for mu
// Returns:
//  (bool): true, false for rw-mutex
//  (*sync.Mutex): underlying sync.Mutex mu
//  (*sync.RWMutex): nil, underlying sync.RWMutex mu for rw-mutex
func (m *Mutex) getLock() (bool, *sync.Mutex, *sync.RWMutex) {
	return true, m.mu, nil
}

// empty getter for isRead, is needed for mutexInt
//  Returns:
//   (*bool): false
func (m *Mutex) getIsRead() *bool {
	res := false
	return &res
}

// ============ FUNCTIONS ============

// Lock mutex m
// Returns:
//  nil
func (m *Mutex) Lock() {
	// call the lock function with the mutexInt interface
	lockInt(m, false)
}

// TryLock mutex m
//  Returns:
//   (bool): true if locking was successful, false otherwise
func (m *Mutex) TryLock() bool {
	// call the try-lock method for the mutexInt interface
	return tryLockInt(m, false)
}

// Unlock mutex m
//  Returns:
//   nil
func (m *Mutex) Unlock() {
	// call the unlock method for the mutexInt interface
	unlockInt(m)
}

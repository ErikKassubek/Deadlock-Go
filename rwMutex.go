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
	mu                   sync.RWMutex
	context              []callerInfo // info about the creation and lock/unlock of this lock
	in                   bool         // set to true after lock was initialized
	isLocked             bool         // set to true if lock is locked
	isLockedRoutineIndex int          // index of the routine, which holds the lock
	memoryPosition       uintptr      // position of the mutex in memory
	isRead               bool         // set true, if last acquisition was RLock
}

// create Lock
func NewRWLock() *RWMutex {
	// initialize detector if necessary
	if !initialized {
		initialize()
	}

	m := RWMutex{
		in:                   true,
		isLockedRoutineIndex: -1,
	}
	_, file, line, _ := runtime.Caller(1)
	m.context = append(m.context, newInfo(file, line, true, ""))
	m.memoryPosition = uintptr(unsafe.Pointer(&m))

	return &m
}

// ====== GETTER ===============================================================

// getter for isLocked
func (m *RWMutex) getIsLocked() *bool {
	return &m.isLocked
}

// getter for isLockedRoutineIndex
func (m *RWMutex) getIsLockedRoutineIndex() *int {
	return &m.isLockedRoutineIndex
}

// getter for context
func (m *RWMutex) getContext() *[]callerInfo {
	return &m.context
}

// getter for memoryPosition
func (m *RWMutex) getMemoryPosition() uintptr {
	return m.memoryPosition
}

// getter for in
func (m *RWMutex) getIn() *bool {
	return &m.in
}

// getter for mu
func (m *RWMutex) getLock() (bool, *sync.Mutex, *sync.RWMutex) {
	return false, nil, &m.mu
}

// getter for isRead
func (m *RWMutex) getIsRead() *bool {
	return &m.isRead
}

// check if lock is rwLock
func (m *RWMutex) isRWLock() bool {
	return true
}

// ====== FUNCTIONS ============================================================

// Lock rwmutex m
func (m *RWMutex) Lock() {
	lockInt(m, false)
	m.isRead = false
}

// RLock rwmutex m
func (m *RWMutex) RLock() {
	lockInt(m, true)
	m.isRead = true
}

// TryLock rwmutex m
func (m *RWMutex) TryLock() {
	tryLockInt(m)
	m.isRead = false
}

// Unlock rwmutex m
func (m *RWMutex) Unlock() {
	unlockInt(m)
}

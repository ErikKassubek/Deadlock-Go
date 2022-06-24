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
routine.go
Implementation of the structure to save the routine wise saved data.
This contains mainly the lock-tree for each routine as well as functionality
to update these trees.
*/

import (
	"runtime"
	"strings"
	"sync"

	"github.com/petermattis/goid"
)

var mapIndex = make(map[int64]int)
var createRoutineLock sync.Mutex
var routines = make([]routine, opts.maxRoutines)
var routinesIndex = 0

// type to implement structures for lock logging
type routine struct {
	index                     int        // index of the routine
	holdingCount              int        // number of currently hold locks
	holdingSet                []mutexInt // set of currently hold locks
	dependencyMap             map[uintptr]*[]*dependency
	dependencies              [](*dependency)  // pre-allocated dependencies
	curDep                    *dependency      // current dependency
	depCount                  int              // counter for dependencies
	collectedSingleLevelLocks map[string][]int // info about collected single level locks
}

// Initialize the go routine
func newRoutine() {
	if !opts.periodicDetection && !opts.comprehensiveDetection {
		return
	}

	createRoutineLock.Lock()
	r := routine{
		index:                     routinesIndex,
		holdingCount:              0,
		holdingSet:                make([]mutexInt, opts.maxHoldingDepth),
		dependencyMap:             make(map[uintptr]*[]*dependency),
		dependencies:              make([]*dependency, opts.maxDependencies),
		curDep:                    nil,
		depCount:                  0,
		collectedSingleLevelLocks: make(map[string][]int),
	}
	if routinesIndex >= opts.maxRoutines {
		panic(`Number of routines is greater than max number of routines. 
			Increase Opts.MaxRoutines.`)
	}
	routines[routinesIndex] = r
	mapIndex[goid.Get()] = routinesIndex
	routinesIndex++
	createRoutineLock.Unlock()
	for i := 0; i < opts.maxDependencies; i++ {
		dep := newDependency(nil, 0, nil)
		r.dependencies[i] = &dep
	}
}

// update the routine structure if a mutex is locked
func (r *routine) updateLock(m mutexInt) {
	currentHolding := r.holdingSet
	hc := r.holdingCount

	isNew := false

	// if lock is not a single level lock
	if hc > 0 {
		// found nested lock
		key := m.getMemoryPosition() ^ currentHolding[r.holdingCount-1].getMemoryPosition()

		depMap := r.dependencyMap
		dhl := make([]*dependency, 0)
		var dep *dependency

		d, ok := depMap[key]

		panicMassage := `Number of dependencies is greater than max number of 
			dependencies. Increase Opts.MaxDependencies.`

		if ok { // dependency key already exists
			dhl = *d
			if !r.hasEntryDhl(m, &dhl) {
				if r.depCount >= opts.maxDependencies {
					panic(panicMassage)
				}
				dep = r.dependencies[r.depCount]
				isNew = true
			}
		} else {
			if r.depCount >= opts.maxDependencies {
				panic(panicMassage)
			}
			dep = r.dependencies[r.depCount]
			isNew = true
		}
		if isNew {
			r.depCount++
			dep.update(m, &currentHolding, hc)
			dhl = append(dhl, dep)
			r.dependencyMap[key] = &dhl
			r.curDep = dep
		}

	} else {
		// save information on single level locks if enabled in the options
		// to avoid creating the caller info multiple times
		if opts.collectSingleLevelLockStack {
			_, file, line, _ := runtime.Caller(3)
			if lines, ok := r.collectedSingleLevelLocks[file]; ok {
				isNew = true
				for _, l := range lines {
					if l == line {
						isNew = false
						break
					}
				}
				if isNew {
					r.collectedSingleLevelLocks[file] = append(
						r.collectedSingleLevelLocks[file], line)
				}
			} else {
				isNew = true
				r.collectedSingleLevelLocks[file] = []int{line}
			}
		}
	}

	if isNew && (hc > 0 || opts.collectSingleLevelLockStack) {
		var file string
		var line int
		var bufStringCleaned string
		if opts.collectCallStack {
			var bufString string
			buf := make([]byte, opts.maxCallStackSize)
			n := runtime.Stack(buf[:], false)
			bufString = string(buf[:n])
			bufStringSplit := strings.Split(bufString, "\n")
			bufStringCleaned = bufStringSplit[0] + "\n"
			for i := 5; i < len(bufStringSplit); i++ {
				bufStringCleaned += bufStringSplit[i] + "\n"
			}
		}

		_, file, line, _ = runtime.Caller(3)

		context := m.getContext()
		*context = append(*context, newInfo(file, line, false, bufStringCleaned))
	}

	if hc >= opts.maxHoldingDepth {
		panic(`Holding Count is grater than maximum holding depth. Increase 
			Opts.MaxHoldingDepth.`)
	}
	currentHolding[hc] = m
	r.holdingCount++
}

// return true, if mutex with same holding count is in the dependency list
func (r *routine) hasEntryDhl(m mutexInt, dhl *([]*dependency)) bool {
	for _, d := range *dhl {
		hc := r.holdingCount
		if d.lock == m && d.holdingCount == hc {
			i := 0
			for d.holdingSet[i] == r.holdingSet[i] && i < hc {
				i++
			}
			if i == hc {
				return true
			}
		}
	}
	return false
}

// update if tryLock is successfully
// this only updates the holding set
func (r *routine) updateTryLock(m mutexInt) {
	hc := r.holdingCount
	if hc >= opts.maxHoldingDepth {
		panic(`Holding Count is grater than maximum holding depth. Increase 
			Opts.MaxHoldingDepth.`)
	}
	r.holdingSet[hc] = m
	r.holdingCount++
}

// update the routine structure is a mutex is released
func (r *routine) updateUnlock(m mutexInt) {
	for i := r.holdingCount - 1; i >= 0; i-- {
		if r.holdingSet[i] == m {
			r.holdingSet = append(r.holdingSet[:i], r.holdingSet[i+1:]...)
			r.holdingSet = append(r.holdingSet, nil)
			r.holdingCount--
			break
		}
	}
}

// get the index of the routine, -1 if routine does not exist
func getRoutineIndex() int {
	id := goid.Get()
	createRoutineLock.Lock()
	index, ok := mapIndex[id]
	createRoutineLock.Unlock()
	if !ok {
		return -1
	}
	return index
}

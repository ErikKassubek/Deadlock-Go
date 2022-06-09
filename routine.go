package deadlock

/*
Author: Erik Kassubek <erik-kassubek@t-online.de>
Package: deadlock
Project: Bachelor Project at the Albert-Ludwigs-University Freiburg,
	Institute of Computer Science: Dynamic Deadlock Detection in Go
Date: 2022-06-05
*/

/*
routine.go
Implementation of the structure to save the routine wise saved data.
This contains mainly the lock-tree for each routine as well as functionality
to update these trees.
*/

import (
	"runtime"
	"sync"
	"unsafe"

	"github.com/petermattis/goid"
)

var mapIndex = make(map[int64]int)
var mapIndexLock sync.Mutex
var routines = make([]routine, opts.maxRoutines)
var routinesIndex = 0

// type to implement structures for lock logging
type routine struct {
	index         int        // index of the routine
	holdingCount  int        // number of currently hold locks
	holdingSet    [](*mutex) // set of currently hold locks
	dependencyMap map[uintptr]*[]*dependency
	dependencies  [](*dependency) // pre-allocated dependencies
	curDep        *dependency     // current dependency
	depCount      int             // counter for dependencies
}

// Initialize the go routine
func NewRoutine() {
	if !opts.periodicDetection && !opts.comprehensiveDetection {
		return
	}
	r := routine{
		index:         routinesIndex,
		holdingCount:  0,
		holdingSet:    make([]*mutex, opts.maxHoldingDepth),
		dependencyMap: make(map[uintptr]*[]*dependency),
		dependencies:  make([]*dependency, opts.maxDependencies),
		curDep:        nil,
		depCount:      0,
	}
	if routinesIndex >= opts.maxRoutines {
		panic(`Number of routines is greater than max number of routines. 
			Increase Opts.MaxRoutines.`)
	}
	routines[routinesIndex] = r
	mapIndexLock.Lock()
	mapIndex[goid.Get()] = routinesIndex
	mapIndexLock.Unlock()
	routinesIndex++
	for i := 0; i < opts.maxDependencies; i++ {
		dep := newDependency(nil, 0, nil)
		r.dependencies[i] = &dep
	}
}

// update the routine structure if a mutex is locked
func (r *routine) updateLock(m *mutex) {
	currentHolding := r.holdingSet
	hc := r.holdingCount

	isDepSet := true

	// if lock is not a single level lock
	if hc > 0 {
		// found nested lock
		key := uintptr(unsafe.Pointer(m)) ^ uintptr(
			unsafe.Pointer(currentHolding[r.holdingCount-1]))

		depMap := r.dependencyMap
		dhl := make([]*dependency, 0)
		var dep *dependency

		d, ok := depMap[key]

		panicMassage := `Number of dependencies is greater than max number of 
			dependencies. Increase Opts.MaxDependencies.`

		if ok {
			dhl = *d
			if r.hasEntryDhl(m, &dhl, dep) {
				if r.depCount >= opts.maxDependencies {
					panic(panicMassage)
				}
				dep = r.dependencies[r.depCount]
				r.depCount++
				dep.update(m, &currentHolding, hc)
				dhl = append(dhl, dep)
			} else {
				isDepSet = false
			}
		} else {
			if r.depCount >= opts.maxDependencies {
				panic(panicMassage)
			}
			dep = r.dependencies[r.depCount]
			r.depCount++
			dep.update(m, &currentHolding, hc)
			dhl = append(dhl, dep)
			depMap[key] = &dhl
		}

		// update current dependency
		r.curDep = dep
	}

	if isDepSet && (hc > 0 || opts.collectSingleLevelLockStack) {
		_, file, line, _ := runtime.Caller(2)
		var bufString string
		if opts.collectCallStack {
			// check whether it is necessary to collect the callStack
			buf := make([]byte, opts.maxCallStackSize)
			n := runtime.Stack(buf[:], false)
			bufString = string(buf[:n])
		}

		m.context = append(m.context, newInfo(file, line, false, bufString))
	}

	if hc >= opts.maxHoldingDepth {
		panic(`Holding Count is grater than maximum holding depth. Increase 
			Opts.MaxHoldingDepth.`)
	}
	currentHolding[hc] = m
	r.holdingCount++
}

// return true, if mutex with same holding count is in the dependency list
func (r *routine) hasEntryDhl(m *mutex, dhl *([]*dependency),
	dep *dependency) bool {
	for _, d := range *dhl {
		hc := r.holdingCount
		if d.lock == m && d.holdingCount == hc {
			i := 0
			for d.holdingSet[i] == r.holdingSet[i] && i <= hc {
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
func (r *routine) updateTryLock(m *mutex) {
	hc := r.holdingCount
	if hc >= opts.maxHoldingDepth {
		panic(`Holding Count is grater than maximum holding depth. Increase 
			Opts.MaxHoldingDepth.`)
	}
	r.holdingSet[hc] = m
	r.holdingCount++
}

// update the routine structure is a mutex is released
func (r *routine) updateUnlock(m *mutex) {
	for i := r.holdingCount - 1; i >= 0; i-- {
		if r.holdingSet[i] == m {
			r.holdingSet = append(r.holdingSet[:i], r.holdingSet[i+1:]...)
			r.holdingSet = append(r.holdingSet, nil)
			r.holdingCount--
			break
		}
	}
}

// get the index of the routine
func getRoutineIndex() int {
	id := goid.Get()
	mapIndexLock.Lock()
	index := mapIndex[id]
	mapIndexLock.Unlock()
	return index
}

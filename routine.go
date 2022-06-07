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
TODO: implement check if collection of callside information is required
*/

import (
	"runtime"
	"sync"
	"unsafe"

	"github.com/petermattis/goid"
)

var mapIndex map[int64]int
var mapIndexLock sync.Mutex
var routines = make([]routine, Opts.MaxRoutines)
var routinesIndex = 0

// type to implement structures for lock logging
type routine struct {
	index         int        // index of the routine
	holdingCount  int        // number of currently hold locks
	holdingSet    [](*mutex) // set of currently hold locks
	dependencyMap map[uintptr]*[]*dependency
	dependencies  [](*dependency) // pre-allocated dependencies
	curDep        *dependency     // current dependency
	depCount      int             // counter for dependenies
}

// Initialize the go routine
func NewRoutine() {
	r := routine{
		index:         routinesIndex,
		holdingCount:  0,
		holdingSet:    make([]*mutex, Opts.MaxHoldingDepth),
		dependencyMap: make(map[uintptr]*[]*dependency),
		dependencies:  make([]*dependency, Opts.MaxDependencies),
		curDep:        nil,
		depCount:      0,
	}
	if routinesIndex >= Opts.MaxRoutines {
		panic(`Number of routines is greater than max number of routines. 
			Increase Opts.MaxRoutines.`)
	}
	routines[routinesIndex] = r
	mapIndexLock.Lock()
	mapIndex[goid.Get()] = routinesIndex
	mapIndexLock.Unlock()
	routinesIndex++
	for i := 0; i < Opts.MaxDependencies; i++ {
		dep := newDependency(nil, 0, nil)
		r.dependencies[i] = &dep
	}
}

// update the routine structure if a mutex is locked
func (r *routine) updateLock(m *mutex) {
	currentHolding := r.holdingSet
	hc := r.holdingCount

	// if lock is not a single level lock
	if hc > 0 {
		// found nested lock
		key := uintptr(unsafe.Pointer(m)) ^ uintptr(
			unsafe.Pointer(currentHolding[r.holdingCount-1]))

		depMap := r.dependencyMap
		dhl := make([]*dependency, 0)
		var dep *dependency

		d, ok := depMap[key]

		isDepSet := true // TODO: remove if replaced by check if callside is necessary

		panicMassage := `Number of dependencies is greater than max number of 
			dependencies. Increase Opts.MaxDependencies.`

		if ok {
			dhl = *d
			if r.hasEntryDhl(m, &dhl, dep) {
				if r.depCount >= Opts.MaxDependencies {
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
			if r.depCount >= Opts.MaxDependencies {
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

		// check wether it is necessary to get the caller info
		// TODO: check wether it is necessary to get stack
		if isDepSet {
			_, file, line, _ := runtime.Caller(2)
			m.context = append(m.context, newInfo(file, line, false))
		}
	}
	if hc >= Opts.MaxHoldingDepth {
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
	if hc >= Opts.MaxHoldingDepth {
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
	return mapIndex[id]
}

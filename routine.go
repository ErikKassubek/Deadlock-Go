package undead

import (
	"runtime"
	"unsafe"

	"github.com/petermattis/goid"
)

var mapIndex map[int64]int
var routines []routine
var routinesIndex = 0

// type to implement structures for lock logging
type routine struct {
	index         int           // index of the routine
	holdingCount  int           // number of currently hold locks
	holdingSet    *([](*mutex)) // set of currently hold locks
	dependencyMap *(map[uintptr]*[]*dependency)
	dependencies  *([](*dependency)) // pre-allocated dependencies
	curDep        *dependency        // current dependency
	depCount      int                // counter for dependenies
}

// Initialize the go routine
func NewRoutine() {
	hs := make([]*mutex, Opts.MaxHoldingDepth)
	dm := make(map[uintptr]*[]*dependency)
	dep := make([]*dependency, Opts.MaxHoldingDepth)

	r := routine{
		index:         routinesIndex,
		holdingCount:  0,
		holdingSet:    &hs,
		dependencyMap: &dm,
		dependencies:  &dep,
		curDep:        nil,
		depCount:      0,
	}
	routines = append(routines, r)
	mapIndex[goid.Get()] = routinesIndex
	routinesIndex++
	for i := 0; i < Opts.MaxHoldingDepth; i++ {
		dep := newDependency(nil, 0, nil)
		(*r.dependencies)[i] = &dep
	}
}

// update the routine structure if a mutex is locked
func (r *routine) updateLock(m *mutex) {
	currentHolding := r.holdingSet
	hc := r.holdingCount

	// if lock is a single level lock
	if r.holdingCount > 0 {
		// found nested lock
		key := uintptr(unsafe.Pointer(m)) ^ uintptr(
			unsafe.Pointer((*currentHolding)[r.holdingCount-1]))

		depMap := r.dependencyMap
		dhl := make([]*dependency, 0)
		var dep *dependency

		d, ok := (*depMap)[key]
		if ok {
			dhl = *d
			if r.hasEntryDhl(m, &dhl, dep) {
				dep = (*r.dependencies)[r.depCount]
				r.depCount++
				dep.update(m, currentHolding, hc)
				dhl = append(dhl, dep)
			}
		} else {
			dep = (*r.dependencies)[r.depCount]
			r.depCount++
			dep.update(m, currentHolding, hc)
			dhl = append(dhl, dep)
			(*depMap)[key] = &dhl
		}

		// update current dependency
		r.curDep = dep

		// check wether it is necessary to get the caller info
		// TODO: check wether it is necessary to get stack
		_, file, line, _ := runtime.Caller(2)
		dep.callsite = newInfo(file, line)

	}
	(*currentHolding)[hc] = m
	r.holdingCount++
}

// return true, if mutex with same holding count is in the dependency list
func (r *routine) hasEntryDhl(m *mutex, dhl *([]*dependency),
	dep *dependency) bool {
	for _, d := range *dhl {
		hc := r.holdingCount
		if d.lock == m && d.holdingCount == hc {
			i := 0
			for d.holdingSet[i] == (*r.holdingSet)[i] && i <= hc {
				i++
			}
			if i == hc {
				return true
			}
		}
	}
	return false
}

// update the routine structure is a mutex is released
func (r *routine) updateUnlock(m *mutex) {
	for i := r.holdingCount - 1; i >= 0; i-- {
		if (*r.holdingSet)[i] == m {
			*r.holdingSet = append((*r.holdingSet)[:i], (*r.holdingSet)[i+1:]...)
			r.holdingCount--
			break
		}
	}
}

// get the index of the routine
// TODO: calculate from memory position
func getRoutineIndex() int {
	id := goid.Get()
	return mapIndex[id]
}

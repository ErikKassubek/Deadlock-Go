package undead

import (
	"runtime"
	"unsafe"

	"github.com/petermattis/goid"
)

// TODO: calculate routine from variable storage position
var routines map[int64]*routine

// type to implement structures for lock logging
type routine struct {
	index            int64                 // index of the routine
	numberOfLocks    int                   // number of currently hold locks
	lockSet          [](*mutex)            // set of currently hold locks
	context          map[*mutex]callerInfo // info about caller
	lockDependencies map[uintptr][]dependency
	curDependency    dependency
}

// Initialize the go routine
func NewRoutine() {
	r := routine{
		index:            goid.Get(),
		numberOfLocks:    0,
		lockSet:          make([]*mutex, 0),
		context:          make(map[*mutex]callerInfo),
		lockDependencies: make(map[uintptr][]dependency)}
	routines[goid.Get()] = &r
}

// update the routine structure is a mutex is locked
func (r *routine) updateLock(m *mutex) {
	// if lock is a single level lock
	if r.numberOfLocks == 0 {
		r.lockSet = append(r.lockSet, m)
		r.numberOfLocks++
		return
	}

	// calculate hash map key
	key := uintptr(unsafe.Pointer(m)) ^ uintptr(
		unsafe.Pointer(r.lockSet[r.numberOfLocks-1]))

	var dep dependency

	if _, ok := r.lockDependencies[key]; !ok {
		// new dependency
		dep = newDependency(m, r.numberOfLocks, r.lockSet)
		r.lockDependencies[key] = append(r.lockDependencies[key], dep)
	} else {
		// possible, that lock already exists, further check
		if !r.hasEntryDependancylist(m, key) {
			dep = newDependency(m, r.numberOfLocks, r.lockSet)
			r.lockDependencies[key] = append(r.lockDependencies[key], dep)
		}
	}

	r.lockSet = append(r.lockSet, m)
	r.numberOfLocks++

	r.curDependency = dep

	// TODO: check if it necessary to get call stack

	// save stack infor
	_, file, line, _ := runtime.Caller(2)
	info := newInfo(file, line)
	r.context[m] = info
}

// return true, if mutex with same holding count is in the dependency list
func (r *routine) hasEntryDependancylist(m *mutex, key uintptr) bool {
	for _, d := range r.lockDependencies[key] {
		lc := r.numberOfLocks
		if d.lock == m && d.numberOfLocks == lc {
			i := 0
			for d.holdingSet[i] == r.lockSet[i] && i <= lc {
				i++
			}
			if i == lc {
				return true
			}
		}
	}
	return false
}

// update the routine structure is a mutex is released
func (r *routine) updateUnlock(m *mutex) {
	for i := r.numberOfLocks - 1; i >= 0; i-- {
		if r.lockSet[i] == m {
			r.lockSet = append(r.lockSet[:i], r.lockSet[i+1:]...)
			r.numberOfLocks--
			break
		}
	}
}

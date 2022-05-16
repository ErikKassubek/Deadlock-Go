package undead

import "unsafe"

var routines []*(Routine)

// type to implement structures for lock logging
type Routine struct {
	numberOfLocks    int                   // number of currently hold locks
	lockSet          [](*Mutex)            // set of currently hold locks
	context          map[*Mutex]callerInfo // info about caller
	lockDependencies map[uintptr]dependency
}

// Initialize the go routine
func NewRoutine() *Routine {
	r := Routine{
		lockSet: make([]*Mutex, 0),
		context: make(map[*Mutex]callerInfo)}
	routines = append(routines, &r)
	return &r
}

// update the routine structure is a mutex is locked
func (r *Routine) updateLock(m *Mutex) {
	// TODO: check if m already in lockSet ?

	// if lock is a single level lock
	if r.numberOfLocks == 0 {
		r.lockSet = append(r.lockSet, m)
		r.numberOfLocks++
		return
	}

	// calculate hash map key
	key := uintptr(unsafe.Pointer(m)) ^ uintptr(
		unsafe.Pointer(r.lockSet[r.numberOfLocks-1]))

	if _, ok := r.lockDependencies[key]; !ok {
		// new dependency
		dep := newDependency(m, r.numberOfLocks, r.lockSet)
		r.lockDependencies[key] = dep
	} else {

	}

	r.lockSet = append(r.lockSet, m)
	r.numberOfLocks++
}
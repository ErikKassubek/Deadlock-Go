package deadlock

/*
Author: Erik Kassubek <erik-kassubek@t-online.de>
Package: deadlock
Project: Bachelor Project at the Albert-Ludwigs-University Freiburg,
	Institute of Computer Science: Dynamic Deadlock Detection in Go
Date: 2022-06-05
*/

/*
dependency.go
Implement structure to describe dependencies.
Each dependency contains the lock l as well as all the locks l depends on.
These are all the locks which were hold by the same routine, when l was
acquired.
*/

// struct to implement a dependency
type dependency struct {
	lock         *Mutex     // lock
	holdingCount int        // on how many locks does mu depend
	holdingSet   [](*Mutex) // lock which where hold while mu was acquired
}

// create a new dependency object
func newDependency(lock *Mutex, numberOfLocks int,
	currentLocks [](*Mutex)) dependency {
	d := dependency{
		lock:         lock,
		holdingCount: numberOfLocks,
		holdingSet:   make([]*Mutex, opts.maxHoldingDepth),
	}

	for i := 0; i < numberOfLocks; i++ {
		d.holdingSet = append(d.holdingSet, currentLocks[i])
	}

	return d
}

// update dependencies
func (d *dependency) update(lock *Mutex, hs *[]*Mutex, len int) {
	d.lock = lock
	d.holdingCount = len
	for i := 0; i < len; i++ {
		d.holdingSet[i] = (*hs)[i]
	}
}

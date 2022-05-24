package deadlock

// type to implement a dependency
type dependency struct {
	lock         *mutex     // lock
	holdingCount int        // on how many locks does mu depend
	holdingSet   [](*mutex) // lock which where hold while mu was acquired
}

// create a new dependency object
func newDependency(lock *mutex, numberOfLocks int,
	currentLocks [](*mutex)) dependency {
	d := dependency{
		lock:         lock,
		holdingCount: numberOfLocks,
		holdingSet:   make([]*mutex, Opts.MaxHoldingDepth),
	}

	for i := 0; i < numberOfLocks; i++ {
		d.holdingSet = append(d.holdingSet, currentLocks[i])
	}

	return d
}

// update dependencies
func (d *dependency) update(lock *mutex, hs *[]*mutex, len int) {
	d.lock = lock
	d.holdingCount = len
	for i := 0; i < len; i++ {
		d.holdingSet[i] = (*hs)[i]
	}
}

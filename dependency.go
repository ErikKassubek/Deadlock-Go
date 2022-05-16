package undead

// type to implement a dependency
type dependency struct {
	mu      *Mutex     // lock
	count   int        // mu depends on count locks
	LockSet [](*Mutex) // lock which where hold while mu was acquired
}

// create a new dependency object
func newDependency(lock *Mutex, numberOfLocks int,
	currentLocks [](*Mutex)) dependency {
	d := dependency{
		mu:      lock,
		count:   numberOfLocks,
		LockSet: make([]*Mutex, 0),
	}

	for i := 0; i < numberOfLocks; i++ {
		d.LockSet = append(d.LockSet, currentLocks[i])
	}

	return d
}

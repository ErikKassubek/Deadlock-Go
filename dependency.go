package undead

// type to implement a dependency
type dependency struct {
	lock          *mutex     // lock
	numberOfLocks int        // on how many locks does mu depend
	holdingSet    [](*mutex) // lock which where hold while mu was acquired
	callsiteCount int        // TODO: what is that, change name if nessesary
}

// create a new dependency object
func newDependency(lock *mutex, numberOfLocks int,
	currentLocks [](*mutex)) dependency {
	d := dependency{
		lock:          lock,
		numberOfLocks: numberOfLocks,
		holdingSet:    make([]*mutex, 0),
		callsiteCount: 0,
	}

	for i := 0; i < numberOfLocks; i++ {
		d.holdingSet = append(d.holdingSet, currentLocks[i])
	}

	return d
}

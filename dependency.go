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
dependency.go
Implement structure to describe dependencies.
Each dependency contains the lock l as well as all the locks l depends on.
These are all the locks which were hold by the same routine, when l was
acquired.
*/

// Type to implement a dependency
// A dependency represents a set of edges in a lock tree
// It consist of a lock l and a list of all locks, on which l depends
// i.e. all lock which were already locked by the same routine, when
// l was acquired.
type dependency struct {
	mu           mutexInt   // lock
	holdingSet   []mutexInt // locks which where locked while mu was acquired
	holdingCount int        // on how many locks does mu depend
}

// newDependency creates and returns a new dependency object
//  Args:
//   mu (mutexInt): lock of the dependency
//   currentLocks ([]mutexInt): list of locks mu depends on
//   numberOfLocks (int): number of locks lock depends on
//  Returns:
//   (dependency) : the created dependency
func newDependency(lock mutexInt, currentLocks []mutexInt,
	numberOfLocks int) dependency {
	// create dependency
	d := dependency{
		mu:           lock,
		holdingCount: numberOfLocks,
		holdingSet:   make([]mutexInt, opts.maxNumberOfDependentLocks),
	}

	// copy currentLocks into d.holding set
	for i := 0; i < numberOfLocks; i++ {
		d.holdingSet = append(d.holdingSet, currentLocks[i])
	}

	return d
}

// update updates a dependency
//  Args:
//   lock (mutexInt): new lock of the dependency
//   hs (*[]mutexInt): new holding set
//   numberOfLocks (int): new number of locks lock depends on
//  Returns:
//   nil
func (d *dependency) update(lock mutexInt, hs *[]mutexInt, numberOfLocks int) {
	// set new lock
	d.mu = lock

	// set new holding set
	for i := 0; i < numberOfLocks; i++ {
		d.holdingSet[i] = (*hs)[i]
	}

	// set element in d.holdingSet to nil if they were not replaced by
	// new elements
	for i := numberOfLocks; i < d.holdingCount; i++ {
		d.holdingSet[i] = nil
	}

	// set new holdingCount
	d.holdingCount = numberOfLocks
}

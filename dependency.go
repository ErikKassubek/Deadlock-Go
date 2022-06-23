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

// struct to implement a dependency
type dependency struct {
	lock         mutexInt   // lock
	holdingCount int        // on how many locks does mu depend
	holdingSet   []mutexInt // lock which where hold while mu was acquired
}

// create a new dependency object
func newDependency(lock mutexInt, numberOfLocks int,
	currentLocks [](mutexInt)) dependency {
	d := dependency{
		lock:         lock,
		holdingCount: numberOfLocks,
		holdingSet:   make([]mutexInt, opts.maxHoldingDepth),
	}

	for i := 0; i < numberOfLocks; i++ {
		d.holdingSet = append(d.holdingSet, currentLocks[i])
	}

	return d
}

// update dependencies
func (d *dependency) update(lock mutexInt, hs *[]mutexInt, len int) {
	d.lock = lock
	d.holdingCount = len
	for i := 0; i < len; i++ {
		d.holdingSet[i] = (*hs)[i]
	}
}

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
initialize.go
This code initializes the deadlock detector. Its main task is to start
and periodically run the periodical deadlock detection.
*/

import (
	"time"
)

// global variable to check whether the detector was already initialized
var initialized = false

// initialize initializes the deadlock detector.
// This starts the periodical detection.
//  Returns:
//   nil
func initialize() {
	initialized = true

	// reinitialize routines to set size
	routines = make([]routine, opts.maxRoutines)

	// return if periodical detection is disabled
	if !opts.periodicDetection {
		return
	}

	// go routine to run the periodical detection in the background
	go func() {
		// timer to send a signals at equal intervals
		timer := time.NewTicker(opts.periodicDetectionTime)

		// initialize lashHolding. This slice stores the dependencies which were
		// considered in the last detection round, so that the detection only takes
		// place, if the situation has changed
		lastHolding := make([]mutexInt, opts.maxRoutines)

		// run the periodical detection if a timer signal is received
		for range timer.C {
			periodicalDetection(&lastHolding)
		}
	}()
}

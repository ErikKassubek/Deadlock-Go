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

// it is not possible to set options after initialization
var initialized = false

// initialize deadlock detector
func initialize() {
	initialized = true
	// if periodical detection is disabled
	if !opts.periodicDetection {
		return
	}

	go func() {
		timer := time.NewTicker(opts.periodicDetectionTime)

		// for each routine only the dependency, which was last added will be used
		// in the periodical detection
		lastHolding := make([]mutexInt, opts.maxRoutines)

		for {
			select {
			case <-timer.C:
				periodicalDetection(&lastHolding)
			}
		}
	}()

}

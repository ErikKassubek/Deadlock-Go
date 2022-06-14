package deadlock

/*
Author: Erik Kassubek <erik-kassubek@t-online.de>
Package: deadlock
Project: Bachelor Project at the Albert-Ludwigs-University Freiburg,
	Institute of Computer Science: Dynamic Deadlock Detection in Go
Date: 2022-06-05
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
		stack := newDepStack()
		lastHolding := make([](*Mutex), opts.maxRoutines)

		for {
			select {
			case <-timer.C:
				periodicalDetection(&stack, &lastHolding)
			}
		}
	}()

}

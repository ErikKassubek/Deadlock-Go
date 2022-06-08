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

// initialize deadlock detector
func Initialize() {
	mapIndex = make(map[int64]int)

	// if periodical detection is disabled
	if !Opts.PeriodicDetection {
		return
	}

	go func() {
		timer := time.NewTicker(Opts.PeriodicDetectionTime)
		stack := newDepStack()

		for {
			select {
			case <-timer.C:
				periodicalDetection(&stack)
			}
		}
	}()
}

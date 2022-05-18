package undead

import (
	"time"
)

// initialize deadlock detector
func Initialize() {
	routines = make(map[int64]*routine)

	// if periodical detection is disabled
	if !Opts.PeriodicDetection {
		return
	}

	go func() {
		timer := time.NewTicker(Opts.PeriodicDetectionTime)

		for {
			select {
			case <-timer.C:
				periodicalDetection()
			}
		}
	}()
}

// run periodical deadlock detection check
func periodicalDetection() {
	// TODO: implement periodical detection
}

// run comprehensive detection is program is terminated
func Detection() {
	// TODO: implement comprehensive detection
}

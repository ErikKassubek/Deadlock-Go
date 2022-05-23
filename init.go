package deadlock

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
		stack := newChainStack()

		for {
			select {
			case <-timer.C:
				periodicalDetection(&stack)
			}
		}
	}()
}

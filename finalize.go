package deadlock

/*
Author: Erik Kassubek <erik-kassubek@t-online.de>
Package: deadlock
Project: Bachelor Project at the Albert-Ludwigs-University Freiburg,
	Institute of Computer Science: Dynamic Deadlock Detection in Go
Date: 2022-06-05
*/

/*
finalize.go
This file starts the comprehensive deadlock detection. The detection is only
started if more then 1 routine was started.
*/

func Finalize() {
	if !Opts.ComprehensiveDetection {
		return
	}

	detector := newDetector()
	if routinesIndex > 1 {
		detector.routineIndex = routinesIndex
		if detector.preCheck() < 2 {
			return
		}
		detector.detect()
	}
}

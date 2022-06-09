package deadlock

/*
Author: Erik Kassubek <erik-kassubek@t-online.de>
Package: deadlock
Project: Bachelor Project at the Albert-Ludwigs-University Freiburg,
	Institute of Computer Science: Dynamic Deadlock Detection in Go
Date: 2022-06-05
*/

/*
options.go
This file implements options for the deadlock detections such as the
enabling or disabling of the periodical and/or comprehensive detection as
well as the periodical detection time and max values for the detection.
*/

import "time"

// opts controls how the detection behaves
var opts = struct {
	// If periodicDetection is set to false, periodic detection is disabled
	periodicDetection bool
	// If comprehensiveDetection is set to false, comprehensive detection at
	// the end of the program is disabled
	comprehensiveDetection bool
	// Set how often the periodic detection is run
	periodicDetectionTime time.Duration
	// If collectCallStack is true, the CallStack for lock creation and
	// acquisition are collected and displayed. Otherwise only file names and
	// lines are collected
	collectCallStack bool
	// If collectSingleLevelLockStack is set to true, stack traces for single
	// level locks are collected. Otherwise not.
	collectSingleLevelLockStack bool
	// maximum number of dependencies
	maxDependencies int
	// The maximum depth of a nested lock tree
	maxHoldingDepth int
	// The maximum number of routines
	maxRoutines int
	// The maximum byte size for callStacks
	maxCallStackSize int
}{
	periodicDetection:           true,
	comprehensiveDetection:      true,
	periodicDetectionTime:       time.Second * 2,
	collectCallStack:            false,
	collectSingleLevelLockStack: false,
	maxDependencies:             4096,
	maxHoldingDepth:             128,
	maxRoutines:                 1024,
	maxCallStackSize:            2048,
}

// Enable or disable periodic detection
// Return true if detection was successful
// Return false if setting was unsuccessful
// It is not possible to set options after the detector was initialized
func SetPeriodicDetection(enable bool) bool {
	if initialized {
		return false
	}
	opts.periodicDetection = enable
	return true
}

// Enable or disable comprehensive detection
// Return true if detection was successful
// Return false if setting was unsuccessful
// It is not possible to set options after the detector was initialized
func SetComprehensiveDetection(enable bool) bool {
	if initialized {
		return false
	}
	opts.comprehensiveDetection = enable
	return true
}

// Set the periodic detection time in second
// Return true if detection was successful
// Return false if setting was unsuccessful
// It is not possible to set options after the detector was initialized
func SetPeriodicDetectionTime(seconds int) bool {
	if initialized {
		return false
	}
	opts.periodicDetectionTime = time.Second * time.Duration(seconds)
	return true
}

// Enable or disable collection of full call stacks
// If it is disabled only file and line numbers are collected
// Return true if detection was successful
// Return false if setting was unsuccessful
// It is not possible to set options after the detector was initialized
func SetCollectCallStack(enable bool) bool {
	if initialized {
		return false
	}
	opts.collectCallStack = enable
	return true
}

// Enable or disable collection of call information for single level locks
// If it is disabled only file and line numbers are collected
// Return true if detection was successful
// Return false if setting was unsuccessful
// It is not possible to set options after the detector was initialized
func SetCollectSingleLevelLockStack(enable bool) bool {
	if initialized {
		return false
	}
	opts.collectSingleLevelLockStack = enable
	return true
}

// Set the max number of dependencies
// If it is disabled only file and line numbers are collected
// Return true if detection was successful
// Return false if setting was unsuccessful
// It is not possible to set options after the detector was initialized
func SetMaxDependencies(number int) bool {
	if initialized {
		return false
	}
	opts.maxDependencies = number
	return true
}

// Set the max depth of a nested lock tree
// If it is disabled only file and line numbers are collected
// Return true if detection was successful
// Return false if setting was unsuccessful
// It is not possible to set options after the detector was initialized
func SetMaxHoldingDepth(number int) bool {
	if initialized {
		return false
	}
	opts.maxHoldingDepth = number
	return true
}

// Set the max number of routines
// If it is disabled only file and line numbers are collected
// Return true if detection was successful
// Return false if setting was unsuccessful
// It is not possible to set options after the detector was initialized
func SetMaxRoutines(number int) bool {
	if initialized {
		return false
	}
	opts.maxRoutines = number
	return true
}

// Set the max size of collected call stacks
// If it is disabled only file and line numbers are collected
// Return true if detection was successful
// Return false if setting was unsuccessful
// It is not possible to set options after the detector was initialized
func SetMaxCallStackSize(number int) bool {
	if initialized {
		return false
	}
	opts.maxCallStackSize = number
	return true
}

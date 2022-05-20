package undead

import "time"

// Opts controls how the detection behaves
var Opts = struct {
	// If RunDetection is set to false, no detection is disabled
	RunDetection bool
	// If PeriodicDetection is set to false, periodic detection is disabled
	PeriodicDetection bool
	// If ComprehensiveDetection is set to false, comprehensive detection at
	// the end of the program is disabled
	ComprehensiveDetection bool
	// Set how often the periodic detection is run
	PeriodicDetectionTime time.Duration
	// maximum number of dependencies
	MaxDependencies int
	// The maximum depth of a nested lock tree
	MaxHoldingDepth int
	// The maximum number of routines
	MaxRoutines int
}{
	RunDetection:           true,
	PeriodicDetection:      true,
	ComprehensiveDetection: true,
	PeriodicDetectionTime:  time.Second * 2,
	MaxDependencies:        4096,
	MaxHoldingDepth:        6,
	MaxRoutines:            1024,
}

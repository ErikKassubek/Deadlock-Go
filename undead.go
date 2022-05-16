package undead

// initialize deadlock detector
func Initialize() {

}

// Lock mutex m in routine r
// TODO: change so that r is calculated and taken from Routines
func Lock(m Mutex, r *Routine) {
	defer m.mu.Lock()

	// if detection is disabled
	if !Opts.RunDetection {
		return
	}

	// update data structures
	r.updateLock(&m)

}

// Unlock mutex m
func Unlock(m Mutex) {
	defer m.mu.Unlock()

}

// run periodical deadlock detection check
func PeriodicalDetection() {

}

// run comprehensive detection is program is terminated
func Detection() {

}

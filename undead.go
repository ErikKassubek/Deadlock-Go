package undead

// initialize deadlock detector
func Initialize() {
	routines = make(map[int64]Routine)
}

// run periodical deadlock detection check
func PeriodicalDetection() {

}

// run comprehensive detection is program is terminated
func Detection() {

}

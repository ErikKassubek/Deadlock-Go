package deadlock

/*
Author: Erik Kassubek <erik-kassubek@t-online.de>
Package: deadlock
Project: Bachelor Project at the Albert-Ludwigs-University Freiburg,
	Institute of Computer Science: Dynamic Deadlock Detection in Go
Date: 2022-06-05
*/

/*
callerInfo.go
Implementation of a struct to save the caller info of locks
*/

// type to save info about caller
type callerInfo struct {
	file       string
	line       int
	create     bool // true: create, false: lock
	callStacks string
}

// create a new caller info
func newInfo(file string, line int, create bool, callStack string) callerInfo {
	return callerInfo{
		file:       file,
		line:       line,
		create:     create,
		callStacks: callStack,
	}
}

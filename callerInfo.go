package deadlock

/*
Copyright (C) 2022  Erik Kassubek

  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.

  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

/*
Author: Erik Kassubek <erik-kassubek@t-online.de>
Package: deadlock
Project: Bachelor Project at the Albert-Ludwigs-University Freiburg,
	Institute of Computer Science: Dynamic Deadlock Detection in Go
*/

/*
callerInfo.go
Implementation of a struct to save the caller info of locks
*/

/* Type to save info about caller.
A caller is an instance where a lock was created or locked.
*/
type callerInfo struct {
	// name of the file with full path
	file string
	// number of the line, in which the lock is created or locked
	line int
	// true: create, false: lock
	create bool
	// string to save the call stack
	callStacks string
}

// newInfo creates and returns a new callerInfo
//  Args:
//   file (string): name of the file
//   line (int): line in the file where the call happenedr
//   create (bool): set to true if the call was a lock creation or false, if it was a lock acquiring
//  Returns:
//   callerInfo: the created callerInfo
func newInfo(file string, line int, create bool, callStack string) callerInfo {
	return callerInfo{
		file:       file,
		line:       line,
		create:     create,
		callStacks: callStack,
	}
}

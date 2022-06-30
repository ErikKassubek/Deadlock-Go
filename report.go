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
report.go
This file contains functions to report deadlock that were found in any of
the deadlock checks
*/

import (
	"fmt"
	"runtime"
)

// colors for deadlock messages
const (
	yellow = "\033[1;33m%s\033[0m"
	red    = "\033[1;31m%s\033[0m"
	blue   = "\033[0;36m%s\033[0m"
)

// report if double locking is detected
func reportDeadlockDoubleLocking(m mutexInt) {
	fmt.Printf(red, "DEADLOCK (DOUBLE LOCKING)\n\n")

	// print information about the involved lock
	fmt.Printf(yellow, "Initialization of lock involved in deadlock:\n\n")
	context := *m.getContext()
	fmt.Println(context[0].file, context[0].line)
	fmt.Println("")
	fmt.Printf(yellow, "Calls of lock involved in deadlock:\n\n")
	for i, call := range context {
		if i == 0 {
			continue
		}
		fmt.Println(call.file, call.line)
	}
	_, file, line, _ := runtime.Caller(3)
	fmt.Println(file, line)
	fmt.Print("\n\n")
}

//report a found deadlock
func reportDeadlock(stack *depStack) {
	fmt.Printf(red, "POTENTIAL DEADLOCK\n\n")

	// print information about the locks in the circle
	fmt.Printf(yellow, "Initialization of locks involved in potential deadlock:\n\n")
	for cl := stack.stack.next; cl != nil; cl = cl.next {
		for _, c := range *cl.depEntry.mu.getContext() {
			if c.create {
				fmt.Println(c.file, c.line)
			}
		}
	}

	if opts.collectCallStack {
		fmt.Printf(yellow, "\nCallStacks of Locks involved in potential deadlock:\n\n")
		for cl := stack.stack.next; cl != nil; cl = cl.next {
			cont := *cl.depEntry.mu.getContext()
			fmt.Printf(blue, "CallStacks for lock created at: ")
			fmt.Printf(blue, cont[0].file)
			fmt.Printf(blue, ":")
			fmt.Printf(blue, fmt.Sprint(cont[0].line))
			fmt.Print("\n")
			for i, c := range cont {
				if i != 0 {
					fmt.Println(c.callStacks)
				}
			}
		}
	} else {
		fmt.Printf(yellow, "\nCalls of locks involved in potential deadlock:\n\n")
		for cl := stack.stack.next; cl != nil; cl = cl.next {
			for i, c := range *cl.depEntry.mu.getContext() {
				if i == 0 {
					fmt.Printf(blue, "Calls for lock created at: ")
					fmt.Printf(blue, c.file)
					fmt.Printf(blue, ":")
					fmt.Printf(blue, fmt.Sprint(c.line))
					fmt.Printf("\n")
				} else {
					fmt.Println(c.file, c.line)
				}
			}
			fmt.Println("")
		}
	}
	fmt.Print("\n\n")
}

// print a message, that the program was terminated because of a detected local deadlock
func reportDeadlockPeriodical() {
	fmt.Printf(red, "THE PROGRAM WAS TERMINATED BECAUSE IT DETECTED A LOCAL DEADLOCK\n\n")
}

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
	"os"
	"runtime"
)

// colors for deadlock messages
const (
	purple = "\033[1;35m%s\033[0m"
	red    = "\033[1;31m%s\033[0m"
	blue   = "\033[0;36m%s\033[0m"
)

// report if double locking is detected
//  Args:
//   m (mutexInt): mutex on which double locking was detected
//  Returns:
//   nil
func reportDeadlockDoubleLocking(m mutexInt) {
	fmt.Fprintf(os.Stderr, red, "DEADLOCK (DOUBLE LOCKING)\n\n")

	// print information about the involved lock
	fmt.Fprintf(os.Stderr, purple, "Initialization of lock involved in deadlock:\n\n")
	context := *m.getContext()
	fmt.Fprintln(os.Stderr, context[0].file, context[0].line)
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintf(os.Stderr, purple, "Calls of lock involved in deadlock:\n\n")
	for i, call := range context {
		if i == 0 {
			continue
		}
		fmt.Fprintln(os.Stderr, call.file, call.line)
	}
	_, file, line, _ := runtime.Caller(4)
	fmt.Fprintln(os.Stderr, file, line)
	fmt.Fprintf(os.Stderr, "\n\n")
}

// report a found deadlock
//  Args:
//   stack (*depStack) stack which represents the found cycle
//  Returns:
//   nil
func reportDeadlock(stack *depStack) {
	fmt.Fprintf(os.Stderr, red, "POTENTIAL DEADLOCK\n\n")

	// print information about the locks in the circle
	fmt.Fprintf(os.Stderr, purple, "Initialization of locks involved in potential deadlock:\n\n")
	for cl := stack.stack.next; cl != nil; cl = cl.next {
		for _, c := range *cl.depEntry.mu.getContext() {
			if c.create {
				fmt.Fprintln(os.Stderr, c.file, c.line)
			}
		}
	}

	// print information if call stacks were collected
	if opts.collectCallStack {
		fmt.Fprintf(os.Stderr, purple, "\nCallStacks of Locks involved in potential deadlock:\n\n")
		for cl := stack.stack.next; cl != nil; cl = cl.next {
			cont := *cl.depEntry.mu.getContext()
			fmt.Fprintf(os.Stderr, blue, "CallStacks for lock created at: ")
			fmt.Fprintf(os.Stderr, blue, cont[0].file)
			fmt.Fprintf(os.Stderr, blue, ":")
			fmt.Fprintf(os.Stderr, blue, fmt.Sprint(cont[0].line))
			fmt.Fprintf(os.Stderr, "\n\n")
			for i, c := range cont {
				if i != 0 {
					fmt.Fprint(os.Stderr, c.callStacks)
				}
			}
		}
	} else {
		// print information if only caller information were selected
		fmt.Fprintf(os.Stderr, purple, "\nCalls of locks involved in potential deadlock:\n\n")
		for cl := stack.stack.next; cl != nil; cl = cl.next {
			for i, c := range *cl.depEntry.mu.getContext() {
				if i == 0 {
					fmt.Fprintf(os.Stderr, blue, "Calls for lock created at: ")
					fmt.Fprintf(os.Stderr, blue, c.file)
					fmt.Fprintf(os.Stderr, blue, ":")
					fmt.Fprintf(os.Stderr, blue, fmt.Sprint(c.line))
					fmt.Fprintf(os.Stderr, "\n")
				} else {
					fmt.Fprintln(os.Stderr, c.file, c.line)
				}
			}
			fmt.Fprintln(os.Stderr, "")
		}
	}
	fmt.Fprintf(os.Stderr, "\n\n")
}

// print a message, that the program was terminated because of a detected local deadlock
// Returns:
//  nil
func reportDeadlockPeriodical() {
	fmt.Fprintf(os.Stderr, red, "THE PROGRAM WAS TERMINATED BECAUSE IT DETECTED A LOCAL DEADLOCK\n\n")
}

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
detector.go
This file contains all the functionality to detect circles in the lock-trees
and therefor actual or potential deadlocks. It implements the periodical
detection during the runtime of the program as well as the comprehensive
detection after the program has finished.
The periodical detection searches for actual deadlocks and can stop the
program if it is in a deadlock situation.
The comprehensive detection should run as soon as the actual program has finished.
It is based on iGoodLock and reports potential deadlocks in the code.
*/

import (
	"fmt"
	"os"
	"runtime"
)

// colors for deadlock messages
const (
	yellow = "\033[1;33m%s\033[0m"
	red    = "\033[1;31m%s\033[0m"
	blue   = "\033[0;36m%s\033[0m"
)

// ================ Comprehensive Detection ================

// FindPotentialDeadlock is the main function to start the comprehensive
// detection of deadlocks. It has to be run at the end of a program to
// detect potential deadlocks in the program. This can be one by calling
// it as a defer statement at the beginning of the main function of the
// program.
//  Returns:
//   nil
func FindPotentialDeadlocks() {
	// check if comprehensive detection is disabled, and if do abort deadlock
	//detection
	if !opts.comprehensiveDetection {
		return
	}

	// only run detector if at least two routines were running during the
	// execution of the program
	if routinesIndex > 1 {
		// abort check if the lock trees contain less than 2 unique dependencies
		if !isNumberDependenciesGreaterEqualTwo() {
			return
		}

		// start the detection of potential deadlocks
		detect()
	}
}

// isNumberDependenciesGreaterEqualTwo counts the number of unique dependencies in
// all and checks if it is greater or equal two lock trees.
// It is not necessary to run comprehensive detection if less then
// two unique dependencies exists.
//  Returns:
//   (bool) : true, if number of unique dependencies is greater or equal than 2,false otherwise
func isNumberDependenciesGreaterEqualTwo() bool {
	// number of already found unique dependencies
	depCount := 0

	// the dependencyString is used to identify a dependency pattern
	var dependencyString string

	// dependencyStrings are saved, so that equal dependencies are not counted twice
	dependencyMap := make(map[string]struct{})

	// parse all routines
	for i := 0; i < routinesIndex; i++ {
		current := routines[i]

		// parse routine i
		for j := 0; j < current.depCount; j++ {
			dep := current.dependencies[j]

			// get the dependency string and store it in dependencySting
			getDependencyString(&dependencyString, dep)

			// check if the dependency string already exists
			if _, ok := dependencyMap[dependencyString]; !ok {
				// new dependency was found
				dependencyMap[dependencyString] = struct{}{}
				depCount++
			}

			// if more than two unique dep have been found return true
			if depCount == 2 {
				return true
			}
		}
	}

	// return false if depCount never reached 2
	return false
}

// getDependencyString calculates the dependency string for a given
// dependency. The string is the concatenation of the on the memory positions
// of mu of the dependency and the locks in the holdingSet of the dependency.
//  Args:
//   str (*string): the dependency string is stored in str
//   dep (*dependency): dependency for which the string gets calculated
//  Returns:
//   nil
func getDependencyString(str *string, dep *dependency) {
	// add the memory position of mu of dep
	*str = fmt.Sprint(dep.mu.getMemoryPosition())

	// add the memory position of the locks in the lockSet of dep
	for i := 0; i < dep.holdingCount; i++ {
		*str += fmt.Sprint(dep.holdingSet[i].getMemoryPosition())
	}
}

// detect runs the detection for loops in the lock trees
//  Returns:
//   nil
func detect() {
	// visiting gets set to index of the routine on which the search for circles is started
	var visiting int

	// A stack is used to represent the currently explored path in the lock trees.
	// A dependency is added to the path by pushing it on top of the stack.
	stack := newDepStack()

	// If a routine has been used as starting routine of a cycle search, all
	// possible paths have already been explored and therefore have no circle.
	// The dependencies in this routine can therefor be ignored for the rest
	// of the search.
	// They can also be temporarily ignored, if a dependency of this routine
	// is already in the path which is currently explored
	isTraversed := make([]bool, routinesIndex)

	// traverse all routines as starting routine for the loop search
	for i := 0; i < routinesIndex; i++ {
		routine := routines[i]

		visiting = i

		// traverse all dependencies of the given routine as starting routine
		// for potential paths
		for j := 0; j < routine.depCount; j++ {
			dep := routine.dependencies[j]
			isTraversed[i] = true

			// push the dependency on the stack as first element of the currently
			// explored path
			stack.push(dep, i)

			// start the depth-first search to find potential circular paths
			dfs(&stack, visiting, &isTraversed)

			// remove dep from the stack
			stack.pop()
		}
	}
}

// dfs runs the recursive depth-first search.
// Only paths which build a valid chain are explored.
// After a new dependency is added to the currently explored path, it is checked,
// if the path forms a circle.
//  Args:
//   stack (*depStack): stack witch represent the currently explored path
//   visiting int: index of the routine of the first element in the currently explored path
//   isTraversed (*([]bool)): list which stores which routines have already been traversed
//    (either as starting routine or as a routine which already has a dep in the current path)
//  Returns:
//   nil
func dfs(stack *depStack, visiting int, isTraversed *([]bool)) {
	// Traverse through all routines to find potential next step in the path
	// routines with index <= visiting have already been used as starting routine
	// and therefore don't have to been considered again.
	for i := visiting + 1; i < routinesIndex; i++ {
		routine := routines[i]

		// continue if the routine has already been traversed
		if (*isTraversed)[i] {
			continue
		}

		// go through all dependencies of the current routine
		for j := 0; j < routine.depCount; j++ {
			dep := routine.dependencies[j]
			// check if adding dep to the stack would still be a valid path
			if isChain(stack, dep) {
				// check if adding dep to the stack would lead to a cycle
				if isCycleChain(stack, dep) {
					// report the found potential deadlock
					stack.push(dep, j)
					reportDeadlock(stack)
					stack.pop()
				} else { // the path is not a cycle yet
					// add dep to the current path
					stack.push(dep, i)
					(*isTraversed)[i] = true

					// call dfs recursively to traverse the path further
					dfs(stack, visiting, isTraversed)

					// dep did not lead to a cycle in the lock trees.
					// It is removed to explore different paths
					stack.pop()
					(*isTraversed)[i] = false
				}
			}
		}
	}
}

// ================ Periodical Detection ================

// run periodical deadlock detection check
func periodicalDetection() {
	// only check if at least two routines are running
	if runtime.NumGoroutine() < 2 {
		return
	}

	stack := newDepStack()
	lastHolding := make([]mutexInt, opts.maxRoutines)

	candidates := 0 // number of threads holding locks

	sthNew := false
	for index, r := range routines {
		holds := r.holdingCount - 1
		if holds >= 0 && lastHolding[index] != r.holdingSet[holds] {
			lastHolding[index] = r.holdingSet[holds]
			sthNew = true
			if holds > 0 {
				candidates++
			}
		} else if holds < 0 && lastHolding[index] != nil {
			lastHolding[index] = nil
			sthNew = true
		}
	}

	// if nothing has changed since the last check
	if !sthNew {
		return
	}
	if candidates > 1 {
		detectionPeriodical(&lastHolding, &stack)
	}

}

// analyses the current state for deadlocks
func detectionPeriodical(lastHolding *([]mutexInt), stack *depStack) {
	isTraversed := make([]bool, opts.maxRoutines)
	for index, r := range routines {
		if r.curDep == nil || r.index < 0 {
			continue
		}
		isTraversed[index] = true

		stack.push(r.curDep, index)
		dfsPeriodical(stack, index, isTraversed, lastHolding)
		stack.pop()
		r.curDep = nil
	}
}

// depth first search on current locks
func dfsPeriodical(stack *depStack, visiting int, isTraversed []bool,
	lastHolding *[]mutexInt) {
	for i := visiting + 1; i < routinesIndex; i++ {
		r := routines[i]
		if r.curDep == nil || r.index < 0 {
			continue
		}
		if !isTraversed[i] {
			dep := r.curDep
			if !isChain(stack, dep) {
				continue
			}
			if isCycleChain(stack, dep) {
				stack.push(dep, i)
				sthNew := false
				for cl := stack.list.next; cl != nil; cl = cl.next {
					routineInChain := routines[cl.index]
					holds := routineInChain.holdingCount - 1
					if (holds >= 0 &&
						(*lastHolding)[cl.index] != routineInChain.holdingSet[holds]) ||
						(holds < 0 && (*lastHolding)[cl.index] != nil) {
						sthNew = true
						break
					}
				}
				if !sthNew { // nothing changed in cycled threads, deadlock
					reportDeadlockPeriodical(stack)
					FindPotentialDeadlocks()
					os.Exit(2)
				}
				stack.pop()
			} else {
				isTraversed[routinesIndex] = true
				stack.push(dep, routinesIndex)
				dfsPeriodical(stack, visiting, isTraversed, lastHolding)
				stack.pop()
				isTraversed[routinesIndex] = false
			}
		}
	}
}

// ================ Checks for chains and Cycles ================

// check if adding dep to chain will still be a valid chain
func isChain(stack *depStack, dep *dependency) bool {
	for cl := stack.list.next; cl != nil; cl = cl.next {
		if cl.depEntry == dep {
			return false
		}
		if cl.depEntry.mu == dep.mu {
			return false
		}
		// RLocks do not function as guard locks
		for i := 0; i < cl.depEntry.holdingCount; i++ {
			for j := 0; j < dep.holdingCount; j++ {
				clHs := cl.depEntry.holdingSet[i]
				depHs := dep.holdingSet[j]
				if clHs == depHs {
					if !(*clHs.getIsRead() && *depHs.getIsRead()) {
						return false
					}
				}
			}
		}
	}
	for i := 0; i < dep.holdingCount; i++ {
		if stack.tail.depEntry.mu == dep.holdingSet[i] {
			return true
		}
	}
	return false
}

// check if adding dep to chain will give a cycle chain
func isCycleChain(stack *depStack, dep *dependency) bool {
	for i := 0; i < stack.list.next.depEntry.holdingCount; i++ {
		if stack.list.next.depEntry.holdingSet[i] == dep.mu {
			stack.push(dep, -1)
			res := checkRWCycle(stack)
			stack.pop()
			return res
		}
	}
	return false
}

// check if the cycle does lead to a deadlock even if it contains rwlocks
func checkRWCycle(stack *depStack) bool {
	for c := stack.list.next; c != nil; c = c.next {
		isRead := *c.depEntry.mu.getIsRead()
		if !isRead {
			continue
		}
		for i := 0; i < c.depEntry.holdingCount; i++ {
			next := c.next
			if next == nil {
				next = stack.list.next
			}
			if next.depEntry.holdingSet[i] == c.depEntry.mu {
				isReadHS := *c.depEntry.holdingSet[i].getIsRead()
				if isReadHS {
					return false
				}
			}
		}
	}
	return true
}

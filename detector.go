package deadlock

/*
Author: Erik Kassubek <erik-kassubek@t-online.de>
Package: deadlock
Project: Bachelor Project at the Albert-Ludwigs-University Freiburg,
	Institute of Computer Science: Dynamic Deadlock Detection in Go
Date: 2022-06-05
*/

/*
detector.go
This file contains all the functionality to detect circles in the lock-trees
and therefor actual or potential deadlocks. It implements the periodical
detection during the runtime of the program as well as the comprehensive
detection after the program has finished.
The periodical detection searches for actual deadlocks and can stop the
program if it is in a deadlock situation.
The comprehensive detection is run as soon as the actual program has finished.
It is based on iGoodLock and reports potential deadlocks in the code.
*/

import (
	"fmt"
	"os"
	"runtime"
	"unsafe"
)

// colors for deadlock messages
const (
	yellow = "\033[1;33m%s\033[0m"
	red    = "\033[1;31m%s\033[0m"
	blue   = "\033[0;36m%s\033[0m"
)

type detector struct {
	dependencyMap map[string]*dependency
}

func newDetector() detector {
	return detector{}
}

func FindPotentialDeadlocks() {
	if !opts.comprehensiveDetection {
		return
	}

	detector := newDetector()
	if routinesIndex > 1 {
		if detector.preCheck() < 2 {
			return
		}
		detector.detect()
	}
}

// run periodical deadlock detection check
func periodicalDetection(stack *depStack, lastHolding *[](*Mutex)) {
	// only check if at least two routines are running
	if runtime.NumGoroutine() < 2 {
		return
	}

	candidates := 0 // number of threads holding locks

	sthNew := false
	for index, r := range routines {
		holds := r.holdingCount - 1
		if holds >= 0 && (*lastHolding)[index] != r.holdingSet[holds] {
			(*lastHolding)[index] = r.holdingSet[holds]
			sthNew = true
			if holds > 0 {
				candidates++
			}
		} else if holds < 0 && (*lastHolding)[index] != nil {
			(*lastHolding)[index] = nil
			sthNew = true
		}
	}

	// if nothing has changed since the last check
	if !sthNew {
		return
	}
	if candidates > 1 {
		detectionPeriodical(*lastHolding, stack)
	}

}

// analyses the current state for deadlocks
func detectionPeriodical(lastHolding [](*Mutex), stack *depStack) {
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
	lastHolding []*Mutex) {
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
						lastHolding[cl.index] != routineInChain.holdingSet[holds]) ||
						(holds < 0 && lastHolding[cl.index] != nil) {
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

// check if adding dep to chain will still be a chain
func isChain(stack *depStack, dep *dependency) bool {
	for cl := stack.list.next; cl != nil; cl = cl.next {
		if cl.depEntry == dep {
			return false
		}
		if cl.depEntry.lock == dep.lock {
			return false
		}
		for i := 0; i < cl.depEntry.holdingCount; i++ {
			for j := 0; j < dep.holdingCount; j++ {
				if cl.depEntry.holdingSet[i] == dep.holdingSet[j] {
					return false
				}
			}
		}
	}
	for i := 0; i < dep.holdingCount; i++ {
		if stack.tail.depEntry.lock == dep.holdingSet[i] {
			return true
		}
	}
	return false
}

// check if adding dep to chain will give a cycle chain
func isCycleChain(stack *depStack, dep *dependency) bool {
	for i := 0; i < stack.list.next.depEntry.holdingCount; i++ {
		if stack.list.next.depEntry.holdingSet[i] == dep.lock {
			return true
		}
	}
	return false
}

// output deadlocks detected from current status
// current chain will be the whole cycle
func reportDeadlockPeriodical(stack *depStack) {
	fmt.Printf(red, "PROGRAM RAN INTO DEADLOCK\n\n")
	// fmt.Printf(yellow, "Initialization of locks involved in deadlock\n\n")
	// for ds := stack.list.next; ds != nil; ds = ds.next {
	// 	cont := ds.depEntry.lock.context[0]
	// 	fmt.Println(cont.file, cont.line)
	// }
	// fmt.Printf(yellow, "\n\nCalls of locks involved in deadlock\n")
	// for ds := stack.list.next; ds != nil; ds = ds.next {
	// 	for i, caller := range ds.depEntry.lock.context {
	// 		if i == 0 {
	// 			if opts.collectCallStack {
	// 				fmt.Printf(blue, "\nCallStacks for lock created at: ")
	// 			} else {
	// 				fmt.Printf(blue, "\nCalls for lock created at: ")
	// 			}
	// 			fmt.Printf(blue, caller.file)
	// 			fmt.Printf(blue, ":")
	// 			fmt.Printf(blue, fmt.Sprint(caller.line))
	// 			fmt.Print("\n")
	// 		} else {
	// 			if opts.collectCallStack {
	// 				fmt.Println(caller.callStacks)
	// 			} else {
	// 				fmt.Println(caller.file, ":", caller.line)
	// 			}
	// 		}
	// 	}
	// 	fmt.Println("")
	// }
	// fmt.Println("")
}

// get the amount of unique dependencies
func (d *detector) preCheck() int {
	depCount := 0
	var dependencyString string
	d.dependencyMap = make(map[string]*dependency)
	for i := 0; i < routinesIndex; i++ {
		current := routines[i]
		for j := 0; j < current.depCount; j++ {
			dep := current.dependencies[j]
			getDependencyString(&dependencyString, dep)
			if _, ok := d.dependencyMap[dependencyString]; !ok {
				// new dependency found
				depGlobal := newDependency(dep.lock, dep.holdingCount,
					dep.holdingSet)
				d.dependencyMap[dependencyString] = &depGlobal
				depCount++
			}
		}
	}
	return depCount
}

// return a string to represent an dependency
func getDependencyString(str *string, dep *dependency) {
	*str = fmt.Sprint(uintptr(unsafe.Pointer(dep.lock)))
	for i := 0; i < dep.holdingCount; i++ {
		*str = fmt.Sprint(uintptr(unsafe.Pointer(dep.holdingSet[i])))
	}
}

// run comprehensive detection if program is terminated
func (d *detector) detect() {
	var visiting int
	stack := newDepStack()
	isTraversed := make([]bool, routinesIndex)
	for i := 0; i < routinesIndex; i++ {
		routine := routines[i]
		if routine.depCount == 0 {
			continue
		}
		visiting = i
		for j := 0; j < routine.depCount; j++ {
			dep := routine.dependencies[j]
			isTraversed[i] = true
			stack.push(dep, i)
			d.dfs(&stack, visiting, &isTraversed)
			stack.pop()
		}
	}
}

// recursive search for cycles for the comprehensive detection
func (d *detector) dfs(stack *depStack, visiting int, isTraversed *([]bool)) {
	for i := visiting + 1; i < routinesIndex; i++ {
		routine := routines[i]
		if routine.depCount == 0 {
			continue
		}
		if (*isTraversed)[i] {
			continue
		}
		for j := 0; j < routine.depCount; j++ {
			dep := routine.dependencies[j]
			if isChain(stack, dep) {
				if isCycleChain(stack, dep) {
					d.reportDeadlock(stack, dep)
				} else {
					(*isTraversed)[i] = true
					stack.push(dep, i)
					d.dfs(stack, visiting, isTraversed)
					stack.pop()
					(*isTraversed)[i] = false
				}
			}
		}
	}
}

//report a found deadlock
func (d *detector) reportDeadlock(stack *depStack, dep *dependency) {
	fmt.Printf(red, "POTENTIAL DEADLOCK\n\n")
	fmt.Printf(yellow, "Initialization of locks involved in potential deadlock:\n\n")
	for cl := stack.list.next; cl != nil; cl = cl.next {
		for _, c := range cl.depEntry.lock.context {
			if c.create {
				fmt.Println(c.file, c.line)
			}
		}
	}
	for _, c := range dep.lock.context {
		if c.create {
			fmt.Println(c.file, c.line)
		}
	}

	if opts.collectCallStack {
		fmt.Printf(yellow, "\nCallStacks of Locks involved in potential deadlock:\n\n")
		for cl := stack.list.next; cl != nil; cl = cl.next {
			cont := cl.depEntry.lock.context
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
			fmt.Print("\n\n")
		}
		cont := dep.lock.context
		fmt.Printf(blue, "CallStacks for lock created at: ")
		fmt.Printf(blue, cont[0].file)
		fmt.Printf(blue, ":")
		fmt.Printf(blue, fmt.Sprint(cont[0].line))
		fmt.Print("\n")
		for i, c := range cont {
			if i == 0 {
				continue
			}
			fmt.Println(c.callStacks)
		}
	} else {
		fmt.Printf(yellow, "\nCalls of locks involved in potential deadlock:\n\n")
		for cl := stack.list.next; cl != nil; cl = cl.next {
			for i, c := range cl.depEntry.lock.context {
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
		for i, c := range dep.lock.context {
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

	}
	fmt.Print("\n\n")

}

// check for double locking
func (r *routine) checkDoubleLocking(m *Mutex) {
	for _, l := range r.holdingSet {
		if l == m {
			reportDeadlockDoubleLocking(m)
			FindPotentialDeadlocks()
			os.Exit(2)
		}
	}
}

// report if double locking is detected
func reportDeadlockDoubleLocking(m *Mutex) {
	fmt.Printf(red, "DEADLOCK (DOUBLE LOCKING)\n\n")
	fmt.Printf(yellow, "Initialization of lock involved in deadlock:\n\n")
	fmt.Println(m.context[0].file, m.context[0].line)
	fmt.Print("\n\n")
}

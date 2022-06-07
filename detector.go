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

type detector struct {
	deadlockFound int
	routineIndex  int
	dependencyMap map[string]*dependency
}

func newDetector() detector {
	return detector{deadlockFound: 0}
}

// run periodical deadlock detection check
func periodicalDetection(stack *depStack) {
	// only check if at least two routines are running
	if runtime.NumGoroutine() < 2 {
		return
	}

	lastHolding := make([](*mutex), Opts.MaxRoutines)
	candidates := 0 // number of threads holding locks
	sthNew := false

	checkNew(&lastHolding, &sthNew, &candidates)

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
		_ = detectionPeriodical(lastHolding, stack)
	}

}

// check if something has changed
func checkNew(lastHolding *([](*mutex)), sthNew *bool, candidates *int) {
	for index, r := range routines {
		holds := r.holdingCount - 1
		if holds >= 0 && (*lastHolding)[index] != r.holdingSet[holds] {
			(*lastHolding)[index] = r.holdingSet[holds]
			*sthNew = true
			if holds > 0 {
				(*candidates)++
			}
		} else if holds < 0 && (*lastHolding)[index] != nil {
			(*lastHolding)[index] = nil
			*sthNew = true
		}
	}
}

// analyses the current state for deadlocks
func detectionPeriodical(lastHolding [](*mutex), stack *depStack) (ret bool) {
	ret = false
	isTraversed := make([]bool, routinesIndex)

	for index, r := range routines {
		if r.curDep == nil || r.index < 0 {
			continue
		}
		isTraversed[index] = true

		stack.push(r.curDep, index)
		ret = ret || dfsPeriodical(stack, index, isTraversed, lastHolding)
		stack.pop()
		r.curDep = nil
	}
	return ret
}

// depth first search on current dependencies
func dfsPeriodical(stack *depStack, visiting int, isTraversed []bool,
	lastHolding []*mutex) bool {
	ret := false
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
				ret = true
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
					os.Exit(0)
				}
				stack.pop()
			} else {
				isTraversed[routinesIndex] = true
				stack.push(dep, routinesIndex)
				ret = ret || dfsPeriodical(stack, visiting, isTraversed,
					lastHolding)
				stack.pop()
				isTraversed[routinesIndex] = false
			}
		}
	}
	return ret
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
	fmt.Println("DEADLOCK")
	for ds := stack.list.next; ds != nil; ds = ds.next {
		for _, caller := range ds.depEntry.lock.context {
			fmt.Println(caller.file, caller.line)
		}
		fmt.Println("")
	}
	fmt.Println("")
}

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
	for i := 0; i < d.routineIndex; i++ {
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

func (d *detector) reportDeadlock(stack *depStack, dep *dependency) {
	fmt.Print("POTENTIAL DEADLOCK:\n\n")
	fmt.Print("  Initialization of Locks involved in potential deadlock:\n\n")
	for cl := stack.list.next; cl != nil; cl = cl.next {
		for _, c := range cl.depEntry.lock.context {
			if c.create {
				fmt.Println("    ", c.file, c.line)
			}
		}
	}
	for _, c := range dep.lock.context {
		if c.create {
			fmt.Println("    ", c.file, c.line)
		}
	}

	fmt.Print("\n  Calls of Locks involved in potential deadlock:\n\n")
	for cl := stack.list.next; cl != nil; cl = cl.next {
		for _, c := range cl.depEntry.lock.context {
			fmt.Println("    ", c.file, c.line)
		}
		fmt.Println("")
	}
	for _, c := range dep.lock.context {
		fmt.Println("    ", c.file, c.line)
	}

}

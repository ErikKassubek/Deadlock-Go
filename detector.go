package deadlock

import (
	"fmt"
	"os"
	"runtime"
)

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
		if holds >= 0 && lastHolding[index] != (*r.holdingSet)[holds] {
			lastHolding[index] = (*r.holdingSet)[holds]
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
		if holds >= 0 && (*lastHolding)[index] != (*r.holdingSet)[holds] {
			(*lastHolding)[index] = (*r.holdingSet)[holds]
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
						lastHolding[cl.index] != (*routineInChain.holdingSet)[holds]) ||
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
		fmt.Println(ds.depEntry.lock.context.file, ds.depEntry.lock.context.line)
	}
}

// run comprehensive detection is program is terminated
func Detection() {
	// TODO: implement comprehensive detection
	fmt.Println("1")
}

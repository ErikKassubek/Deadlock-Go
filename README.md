<!-- 
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
-->

<!--
Author: Erik Kassubek <erik-kassubek@t-online.de>
Project: Bachelor Project at the Albert-Ludwigs-University Freiburg,
	Institute of Computer Science: Dynamic Deadlock Detection in Go
-->

# Deadlock-Go: Dynamic Deadlock Detection in Go

## What

Deadlock-Go implements Mutex and RW-Mutex drop-in replacements for 
sync.Mutex and sync.RWMutex with (R)Lock, (R)TryLock and (R)Unlock functionality to detect potential deadlocks.

The detector can detect potential or actually occurring recourse deadlocks
which are caused by cyclic or double locking.

In some cases the detector can result in false-positiv or false-negative
results. E.g., the detector is not able to detect cyclic locking in nested 
routines.

Only works from Go Version 1.18.

## Installation
```
go get github.com/ErikKassubek/Deadlock-Go
```

## Usage Examples
### Example for Mutex
```
import "github.com/ErikKassubek/Deadlock-Go"

func main() {
	defer deadlock.FindPotentialDeadlocks()

	x := deadlock.NewLock()
	y := deadlock.NewLock()
	
	// make sure, that the program does not terminate
	// before all routines have terminated
	ch := make(chan bool, 2)

	go func() {
		x.Lock()
		y.Lock()
		y.Unlock()
		x.Unlock()
		ch <- true
	}()

	go func() {
		y.Lock()
		x.Lock()
		x.Unlock()
		y.Unlock()
		ch <- true
	}()
	<- ch
	<- ch
}
```

### Example for RW-Mutex
```
import "github.com/ErikKassubek/Deadlock-Go"

func main() {
	defer deadlock.FindPotentialDeadlocks()

	x := NewRWLock()
	y := NewRWLock()

	// make sure, that the program does not terminate
	// before all routines have terminated
	ch := make(chan bool, 2)

	go func() {
		x.RLock()
		y.Lock()
		y.Unlock()
		x.Unlock()

		ch <- true
	}()

	go func() {
		y.RLock()
		x.Lock()
		x.Unlock()
		y.Unlock()

		ch <- true
	}()

	<-ch
	<-ch
}
```

## Sample output
### Cyclic Locking
```
POTENTIAL DEADLOCK

Initialization of locks involved in potential deadlock:

/home/***/selfWritten/deadlockGo.go 59
/home/***/selfWritten/deadlockGo.go 60
/home/***/selfWritten/deadlockGo.go 61

Calls of locks involved in potential deadlock:

Calls for lock created at: /home/***/selfWritten/deadlockGo.go:59
/home/***/selfWritten/deadlockGo.go 85
/home/***/selfWritten/deadlockGo.go 66

Calls for lock created at: /home/***/selfWritten/deadlockGo.go:60
/home/***/selfWritten/deadlockGo.go 75
/home/***/selfWritten/deadlockGo.go 67

Calls for lock created at: /home/***/selfWritten/deadlockGo.go:61
/home/***/selfWritten/deadlockGo.go 84
/home/***/selfWritten/deadlockGo.go 76
```

### Double Locking
```
DEADLOCK (DOUBLE LOCKING)

Initialization of lock involved in deadlock:

/home/***/selfWritten/deadlockGo.go 205

Calls of lock involved in deadlock:

/home/***/selfWritten/deadlockGo.go 209
/home/***/selfWritten/deadlockGo.go 210
```

## Options
The behavior of Deadlock-Go can be influenced by different options.
They have to be set before the first lock was initialized.

```SetActivated(enable bool)```: enable or disable all detections at once

```SetPeriodicDetection(enable bool)```: enable or disable periodical detection, default: enabled

```SetComprehensiveDetection(enable bool)```: enable or disable comprehensive detection, default: enabled

```SetPeriodicDetectionTime(seconds int)```: set in which time intervals 
the periodical detection is started, default: 2s

```SetCollectCallStacks(enable bool)```: if enabled, call-stacks for lock 
creation and acquisitions are collected. Otherwise only file and line 
information is collected, default: disabled

```SetCollectSingleLevelLockInformation(enable bool)```: if enabled, information about single-level locks are collected, default enabled

```SetDoubleLockingDetection(enable bool)```: if enabled, detection of double locking is active, default: enabled

Additionally the maximum numbers for the dependencies per Routine (default: 4096),
the maximum number of mutexes a mutex can depend on (default: 128), 
the maximum number of routines (default: 1024) and the maximum 
length of a collected call stack in bytes (default 2048) can be set.  

## Acknowledgement
The detector is partially based on:
```
J. Zhou, S. Silvestro, H. Liu, Y. Cai und T. Liu, „UNDEAD: Detecting and preventing
deadlocks in production software“, in 2017 32nd IEEE/ACM International Conference
on Automated Software Engineering (ASE), Los Alamitos, CA, USA: IEEE Computer
Society, Nov. 2017, S. 729–740. doi: 10.1109/ASE.2017.8115684.
```


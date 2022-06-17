# Deadlock-Go: Dynamic Deadlock Detection in Go

## What

Deadlock-Go implement Mutex drop-in replacements for 
sync.Mutex with Lock, TryLock and Unlock functionality to detect potential 
deadlocks.

Only works from Go Version 1.18.

## Installation
```
go get github.com/ErikKassubek/Deadlock-Go
```

## Usage
```
import "github.com/ErikKassubek/Deadlock-Go"

func main() {
	defer deadlock.FindPotentialDeadlocks()
	x := deadlock.NewLock()
	y := deadlock.NewLock()
	
	// make sure, that program does not terminates
	// before all routines have terminated
	ch := make(chan bool, 2)

	go func() {
		x.Lock()
		y.Lock()
		y.Unlock()
		x.Unlock
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

```SetPeriodicDetection(enable bool)```: enable or disable periodical detection, default: enabled

```SetComprehensiveDetection(enable bool)```: enable or disable comprehensive detection, default: enabled

```SetPeriodicDetectionTime(seconds int)```: set in which time intervals 
the periodical detection is started, default: 2s

```SetCollectCallStacks(enable bool)```: if enabled, call-stacks for lock 
creation and acquisitions are collected. Otherwise only file and line 
information is collected, default: disabled

```SetCollectSingleLevelLockInformation```: if enabled, information about single-level locks are collected, default enabled

```SetDoubleLockingDetection(enable bool)```: if enabled, detection of double locking is active, default: enabled

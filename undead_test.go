package deadlock

/*
Author: Erik Kassubek <erik-kassubek@t-online.de>
Package: deadlock
Project: Bachelor Project at the Albert-Ludwigs-University Freiburg,
	Institute of Computer Science: Dynamic Deadlock Detection in Go
Date: 2022-06-05
*/

/*
undead_test.go
Tests for the deadlock detection
*/

import (
	"math/rand"
	"testing"
	"time"
)

func TestPotentialDeadlock1(t *testing.T) {
	Initialize()
	defer Finalize()

	x := NewLock()
	y := NewLock()
	ch := make(chan bool, 2)

	go func() {
		NewRoutine()
		for i := 0; i < 10; i++ {
			x.Lock()
			y.Lock()
			time.Sleep(time.Second * time.Duration(rand.Float64()))
			y.Unlock()
			x.Unlock()
		}
		ch <- true
	}()

	go func() {
		NewRoutine()
		for i := 0; i < 10; i++ {
			y.Lock()
			x.Lock()
			time.Sleep(time.Second * time.Duration(rand.Float64()))
			x.Unlock()
			y.Unlock()
		}
		ch <- true
	}()

	<-ch
	<-ch

}

func TestActualDeadlock(t *testing.T) {
	Initialize()

	x := NewLock()
	y := NewLock()
	ch := make(chan bool, 2)
	ch2 := make(chan bool)

	go func() {
		NewRoutine()
		x.Lock()
		time.Sleep(time.Second)
		ch2 <- true
		y.Lock()
		y.Unlock()
		x.Unlock()
		ch <- true
	}()

	go func() {
		NewRoutine()
		y.Lock()
		<-ch2
		x.Lock()
		x.Unlock()
		y.Unlock()
		ch <- true
	}()

	<-ch
	<-ch
}

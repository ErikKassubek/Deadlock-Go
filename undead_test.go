package undead

import (
	"testing"
	"time"
)

func TestPotentialDeadlock1(t *testing.T) {
	Initialize()

	x := NewLock()
	y := NewLock()
	ch := make(chan bool, 2)

	go func() {
		NewRoutine()
		x.Lock()
		y.Lock()
		time.Sleep(time.Second * 10)
		y.Unlock()
		x.Unlock()
		ch <- true
	}()

	go func() {
		NewRoutine()
		y.Lock()
		x.Lock()
		time.Sleep(time.Second * 10)
		x.Unlock()
		y.Unlock()
		ch <- true
	}()

	<-ch
	<-ch

	t.Error("")
}

func actualDeadlock() {
	Initialize()

	x := NewLock()
	y := NewLock()
	ch := make(chan bool, 2)
	ch2 := make(chan bool)

	go func() {
		NewRoutine()
		time.Sleep(time.Second)
		ch2 <- true
		x.Lock()
		y.Lock()
		y.Unlock()
		x.Unlock()
		ch <- true
	}()

	go func() {
		NewRoutine()
		<-ch2
		y.Lock()
		x.Lock()
		x.Unlock()
		y.Unlock()
		ch <- true
	}()

	<-ch
	<-ch
}

func TestActualDeadlock(t *testing.T) {
	for i := 0; i < 20; i++ {
		actualDeadlock()
	}
}

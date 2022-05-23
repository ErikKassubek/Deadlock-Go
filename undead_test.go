package deadlock

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

func TestPotentialDeadlockForPeriodicalDetection(t *testing.T) {
	Initialize()

	a := NewLock()
	b := NewLock()
	c := NewLock()
	d := NewLock()

	ch := make(chan bool, 2)
	ch1 := make(chan bool)

	Opts.PeriodicDetectionTime = time.Second * 2

	go func() {
		NewRoutine()

		a.Lock()
		b.Lock()
		ch1 <- true

		time.Sleep(time.Second * 3)

		d.Lock()
		c.Lock()
		c.Unlock()
		d.Unlock()

		b.Unlock()
		a.Unlock()
	}()

	go func() {
		NewRoutine()

		c.Lock()
		d.Lock()

		d.Unlock()
		c.Unlock()

		<-ch1
		b.Lock()
		a.Lock()

		time.Sleep(time.Second * 3)

		d.Lock()

		a.Unlock()
		b.Unlock()
		d.Unlock()

	}()

	<-ch
	<-ch
}

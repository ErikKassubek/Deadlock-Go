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
		y.Unlock()
		x.Unlock()
		ch <- true
	}()

	go func() {
		NewRoutine()
		y.Lock()
		x.Lock()
		x.Unlock()
		y.Unlock()
		ch <- true
	}()

	time.Sleep(time.Second * 20)

	<-ch
	<-ch

	t.Error("")
}

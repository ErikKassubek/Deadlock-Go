package undead

import (
	"testing"
)

func TestPotentialDeadlock1(t *testing.T) {
	Initialize()
	var x Mutex
	var y Mutex
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

	<-ch
	<-ch

	t.Error("")
}

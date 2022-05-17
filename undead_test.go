package undead

import (
	"testing"
)

func TestPotentialDeadlock1(t *testing.T) {
	var x Mutex
	var y Mutex
	ch := make(chan bool, 2)

	go func() {
		r := NewRoutine()
		x.Lock(r)
		y.Lock(r)
		y.Unlock()
		x.Unlock()
		ch <- true
	}()

	go func() {
		r := NewRoutine()
		y.Lock(r)
		x.Lock(r)
		x.Unlock()
		y.Unlock()
		ch <- true
	}()

	<-ch
	<-ch

	t.Error("")
}

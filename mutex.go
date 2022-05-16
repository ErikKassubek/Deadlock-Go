package undead

import "sync"

// type to implement a lock
type Mutex struct {
	mu sync.Mutex
}

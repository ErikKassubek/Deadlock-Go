package deadlock

// helper type for stack
type chainList struct {
	depEntry *dependency
	index    int
	prev     *chainList
	next     *chainList
}

// create a new chainList
func newChainList(dep *dependency, i int) chainList {
	return chainList{
		depEntry: dep,
		index:    i,
		prev:     nil,
		next:     nil,
	}
}

// define a stack
type chainStack struct {
	list *chainList
	tail *chainList
}

// create a new stack
func newChainStack() chainStack {
	cl := newChainList(nil, -1)
	c := chainStack{
		list: &cl,
	}
	c.tail = c.list
	return c
}

// push to stack
func (s *chainStack) push(dep *dependency, index int) {
	cl := newChainList(dep, index)
	s.tail.next = &cl
	cl.prev = s.tail
	s.tail = &cl
}

// pop from stack
func (s *chainStack) pop() {
	if s.tail != s.list {
		s.tail.prev.next = s.tail.next
		s.tail = s.tail.prev
	}
}

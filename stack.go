package deadlock

// helper type for stack
type linkedList struct {
	depEntry *dependency
	index    int
	prev     *linkedList
	next     *linkedList
}

// create a new chainList
func newLinkedList(dep *dependency, i int) linkedList {
	return linkedList{
		depEntry: dep,
		index:    i,
		prev:     nil,
		next:     nil,
	}
}

// define a stack
type depStack struct {
	list *linkedList
	tail *linkedList
}

// create a new stack
func newDepStack() depStack {
	cl := newLinkedList(nil, -1)
	c := depStack{
		list: &cl,
	}
	c.tail = c.list
	return c
}

// push to stack
func (s *depStack) push(dep *dependency, index int) {
	cl := newLinkedList(dep, index)
	s.tail.next = &cl
	cl.prev = s.tail
	s.tail = &cl
}

// pop from stack
func (s *depStack) pop() {
	if s.tail != s.list {
		s.tail.prev.next = s.tail.next
		s.tail = s.tail.prev
	}
}

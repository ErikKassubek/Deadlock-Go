package deadlock

/*
Copyright (C) 2022  Erik Kassubek

  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.

  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

/*
Author: Erik Kassubek <erik-kassubek@t-online.de>
Package: deadlock
Project: Bachelor Project at the Albert-Ludwigs-University Freiburg,
	Institute of Computer Science: Dynamic Deadlock Detection in Go
*/

/*
stack.go
Implementation of a helper type to realize a chain stack
*/

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

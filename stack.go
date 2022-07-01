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
Implementation of a helper type to realize a chain stack.
It is implemented as a double linked list, but it is only possible to remove the
last element (top element on the stack) and add an element at the end of the list
(push on top of the stack)
*/

// ============ stack element ============

// struct to implement a stack element
type stackElement struct {
	// dependency represented by the stack element
	depEntry *dependency
	// index value of the linkedIndex, is set to the index of the routine
	index int
	// pointer to the previous stack element
	prev *stackElement
	// pointer to the next stack element
	next *stackElement
}

// create a new chainList
//  Args:
//   dep (*dependency): dependency which is represented by the stack element
//   i (int): index of the routine which created dep
//  Returns:
//   (stackElement): element for the stack
func newStackElement(dep *dependency, i int) stackElement {
	return stackElement{
		depEntry: dep,
		index:    i,
		prev:     nil,
		next:     nil,
	}
}

// ============ stack ============

// stack for the dependencies
type depStack struct {
	// pointer to the bottom element of the stack
	stack *stackElement
	// pointer to the top element of the stack
	top *stackElement
}

// create a new stack
//  Returns:
//   (depStack): the dependency stack
func newDepStack() depStack {
	cl := newStackElement(nil, -1)

	// set the first element of the stack to an empty stack element
	c := depStack{
		stack: &cl,
	}

	// the first element in the stack is the only and therefore also last element
	c.top = c.stack

	return c
}

// push a new dependency to the stack
//  Args:
//   dep (*dependency): dependency to put on the stack
//   index (int): index of the routine which created the dependency
//  Returns:
//   nil
func (s *depStack) push(dep *dependency, index int) {
	// create the new element
	cl := newStackElement(dep, index)
	// add it to the stack
	s.top.next = &cl
	// reset the pointers of the previous element and the pointer to the top element
	cl.prev = s.top
	s.top = &cl
}

// remove the top element from stack
//  Returns:
//   nil
func (s *depStack) pop() {
	// do nothing if the stack is empty (has only on empty default element)
	if s.top == s.stack {
		return
	}

	// reroute the pointer to remove the top stack element
	s.top.prev.next = s.top.next
	s.top = s.top.prev
}

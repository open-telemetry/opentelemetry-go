// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package list // import "go.opentelemetry.io/otel/sdk/internal/list"

import "go.opentelemetry.io/otel/attribute"

// Adapted from
// https://github.com/golang/go/blob/go1.17.6/src/container/list/list.go

// Attribute is an element of a linked list of OpenTelemetry attributes.
type Attribute struct {
	// Next and previous pointers in the doubly-linked list of elements.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next element of the last
	// list element (l.Back()) and the previous element of the first list
	// element (l.Front()).
	next, prev *Attribute

	// The list to which this element belongs.
	list *Attributes

	// The attribute value stored with this element.
	Value attribute.KeyValue
}

// Next returns the next list element or nil.
func (e *Attribute) Next() *Attribute {
	if p := e.next; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// Prev returns the previous list element or nil.
func (e *Attribute) Prev() *Attribute {
	if p := e.prev; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// Attributes represents a doubly linked list of OpenTelemetry attributes.
// The zero value for Attributes is an empty list ready to use.
type Attributes struct {
	root Attribute // sentinel list element, only &root, root.prev, and root.next are used
	len  int       // current list length excluding (this) sentinel element
}

// Init initializes or clears Attributes a.
func (a *Attributes) Init() *Attributes {
	a.root.next = &a.root
	a.root.prev = &a.root
	a.len = 0
	return a
}

// New returns an initialized list.
func New() *Attributes { return new(Attributes).Init() }

// Len returns the number of elements of Attributes a.
// The complexity is O(1).
func (a *Attributes) Len() int { return a.len }

// Front returns the first element of list a or nil if the list is empty.
func (a *Attributes) Front() *Attribute {
	if a.len == 0 {
		return nil
	}
	return a.root.next
}

// Back returns the last Attribute of list a or nil if the list is empty.
func (a *Attributes) Back() *Attribute {
	if a.len == 0 {
		return nil
	}
	return a.root.prev
}

// lazyInit lazily initializes a zero Attributes value.
func (a *Attributes) lazyInit() {
	if a.root.next == nil {
		a.Init()
	}
}

// insert inserts attr after at, increments a.len, and returns attr.
func (a *Attributes) insert(attr, at *Attribute) *Attribute {
	attr.prev = at
	attr.next = at.next
	attr.prev.next = attr
	attr.next.prev = attr
	attr.list = a
	a.len++
	return attr
}

// insertValue is a convenience wrapper for insert(&Attribute{Value: v}, at).
func (a *Attributes) insertValue(v attribute.KeyValue, at *Attribute) *Attribute {
	return a.insert(&Attribute{Value: v}, at)
}

// remove removes attr from its list, decrements a.len, and returns attr.
func (a *Attributes) remove(attr *Attribute) *Attribute {
	attr.prev.next = attr.next
	attr.next.prev = attr.prev
	attr.next = nil // avoid memory leaks
	attr.prev = nil // avoid memory leaks
	attr.list = nil
	a.len--
	return attr
}

// move moves attr to next to at and returns attr.
func (a *Attributes) move(attr, at *Attribute) *Attribute {
	if attr == at {
		return attr
	}
	attr.prev.next = attr.next
	attr.next.prev = attr.prev

	attr.prev = at
	attr.next = at.next
	attr.prev.next = attr
	attr.next.prev = attr

	return attr
}

// Remove removes attr from a if attr is an element of list.
// It returns the element value attr.Value.
// The Attribute must not be nil.
func (a *Attributes) Remove(attr *Attribute) attribute.KeyValue {
	if attr.list == a {
		// if attr.list == l, l must have been initialized when attr was
		// inserted in a or a == nil (attr is a zero Element) and a.remove
		// will crash
		a.remove(attr)
	}
	return attr.Value
}

// PushFront inserts a new Attribute with value v at the front of Attributes a
// and returns the new Attribute.
func (a *Attributes) PushFront(v attribute.KeyValue) *Attribute {
	a.lazyInit()
	return a.insertValue(v, &a.root)
}

// PushBack inserts a new element attr with value v at the back of Attributes
// a and returns attr.
func (a *Attributes) PushBack(v attribute.KeyValue) *Attribute {
	a.lazyInit()
	return a.insertValue(v, a.root.prev)
}

// InsertBefore inserts a new element attr with value v immediately before
// mark and returns attr. If mark is not an element of a, the list is not
// modified. The mark must not be nil.
func (a *Attributes) InsertBefore(v attribute.KeyValue, mark *Attribute) *Attribute {
	if mark.list != a {
		return nil
	}
	// see comment in Attributes.Remove about initialization of l
	return a.insertValue(v, mark.prev)
}

// InsertAfter inserts a new element attr with value v immediately after mark
// and returns attr. If mark is not an element of a, the list is not modified.
// The mark must not be nil.
func (a *Attributes) InsertAfter(v attribute.KeyValue, mark *Attribute) *Attribute {
	if mark.list != a {
		return nil
	}
	// see comment in Attributes.Remove about initialization of a
	return a.insertValue(v, mark)
}

// MoveToFront moves element attr to the front of list a.
// If attr is not an element of a, the list is not modified.
// The element must not be nil.
func (a *Attributes) MoveToFront(attr *Attribute) {
	if attr.list != a || a.root.next == attr {
		return
	}
	// see comment in Attributes.Remove about initialization of a
	a.move(attr, &a.root)
}

// MoveToBack moves element attr to the back of list a.
// If attr is not an element of a, the list is not modified.
// The element must not be nil.
func (a *Attributes) MoveToBack(attr *Attribute) {
	if attr.list != a || a.root.prev == attr {
		return
	}
	// see comment in Attributes.Remove about initialization of a
	a.move(attr, a.root.prev)
}

// MoveBefore moves element attr to its new position before mark.
// If attr or mark is not an element of a, or attr == mark, the list is not
// modified.
// The element and mark must not be nil.
func (a *Attributes) MoveBefore(attr, mark *Attribute) {
	if attr.list != a || attr == mark || mark.list != a {
		return
	}
	a.move(attr, mark.prev)
}

// MoveAfter moves element attr to its new position after mark.
// If attr or mark is not an element of a, or attr == mark, the list is not
// modified.
// The element and mark must not be nil.
func (a *Attributes) MoveAfter(attr, mark *Attribute) {
	if attr.list != a || attr == mark || mark.list != a {
		return
	}
	a.move(attr, mark)
}

// PushBackList a copy of another list at the back of list a.
// The lists a and other may be the same. They must not be nil.
func (a *Attributes) PushBackList(other *Attributes) {
	a.lazyInit()
	for i, e := other.Len(), other.Front(); i > 0; i, e = i-1, e.Next() {
		a.insertValue(e.Value, a.root.prev)
	}
}

// PushFrontList inserts a copy of another list at the front of list a.
// The lists a and other may be the same. They must not be nil.
func (a *Attributes) PushFrontList(other *Attributes) {
	a.lazyInit()
	for i, e := other.Len(), other.Back(); i > 0; i, e = i-1, e.Prev() {
		a.insertValue(e.Value, &a.root)
	}
}

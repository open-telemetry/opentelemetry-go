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

package list

import (
	"testing"

	"go.opentelemetry.io/otel/attribute"
)

// Adapted from
// https://github.com/golang/go/blob/go1.17.6/src/container/list/list_test.go

func checkAttributesLen(t *testing.T, a *Attributes, len int) bool {
	if n := a.Len(); n != len {
		t.Errorf("l.Len() = %d, want %d", n, len)
		return false
	}
	return true
}

func checkAttributesPointers(t *testing.T, a *Attributes, attrs []*Attribute) {
	root := &a.root

	if !checkAttributesLen(t, a, len(attrs)) {
		return
	}

	// zero length lists must be the zero value or properly initialized (sentinel circle)
	if len(attrs) == 0 {
		if a.root.next != nil && a.root.next != root || a.root.prev != nil && a.root.prev != root {
			t.Errorf("l.root.next = %p, l.root.prev = %p; both should both be nil or %p", a.root.next, a.root.prev, root)
		}
		return
	}
	// len(es) > 0

	// check internal and external prev/next connections
	for i, e := range attrs {
		prev := root
		Prev := (*Attribute)(nil)
		if i > 0 {
			prev = attrs[i-1]
			Prev = prev
		}
		if p := e.prev; p != prev {
			t.Errorf("elt[%d](%p).prev = %p, want %p", i, e, p, prev)
		}
		if p := e.Prev(); p != Prev {
			t.Errorf("elt[%d](%p).Prev() = %p, want %p", i, e, p, Prev)
		}

		next := root
		Next := (*Attribute)(nil)
		if i < len(attrs)-1 {
			next = attrs[i+1]
			Next = next
		}
		if n := e.next; n != next {
			t.Errorf("elt[%d](%p).next = %p, want %p", i, e, n, next)
		}
		if n := e.Next(); n != Next {
			t.Errorf("elt[%d](%p).Next() = %p, want %p", i, e, n, Next)
		}
	}
}

var testAttrs = []attribute.KeyValue{
	attribute.String("zero", "a"),
	attribute.Int("one", 1),
	attribute.Int("two", 2),
	attribute.Int("three", 3),
	attribute.String("four", "banana"),
}

func TestAttributes(t *testing.T) {
	l := New()
	checkAttributesPointers(t, l, []*Attribute{})

	// Single element list
	e := l.PushFront(testAttrs[0])
	checkAttributesPointers(t, l, []*Attribute{e})
	l.MoveToFront(e)
	checkAttributesPointers(t, l, []*Attribute{e})
	l.MoveToBack(e)
	checkAttributesPointers(t, l, []*Attribute{e})
	l.Remove(e)
	checkAttributesPointers(t, l, []*Attribute{})

	// Bigger list
	e2 := l.PushFront(testAttrs[2])
	e1 := l.PushFront(testAttrs[1])
	e3 := l.PushBack(testAttrs[3])
	e4 := l.PushBack(testAttrs[4])
	checkAttributesPointers(t, l, []*Attribute{e1, e2, e3, e4})

	l.Remove(e2)
	checkAttributesPointers(t, l, []*Attribute{e1, e3, e4})

	l.MoveToFront(e3) // move from middle
	checkAttributesPointers(t, l, []*Attribute{e3, e1, e4})

	l.MoveToFront(e1)
	l.MoveToBack(e3) // move from middle
	checkAttributesPointers(t, l, []*Attribute{e1, e4, e3})

	l.MoveToFront(e3) // move from back
	checkAttributesPointers(t, l, []*Attribute{e3, e1, e4})
	l.MoveToFront(e3) // should be no-op
	checkAttributesPointers(t, l, []*Attribute{e3, e1, e4})

	l.MoveToBack(e3) // move from front
	checkAttributesPointers(t, l, []*Attribute{e1, e4, e3})
	l.MoveToBack(e3) // should be no-op
	checkAttributesPointers(t, l, []*Attribute{e1, e4, e3})

	e2 = l.InsertBefore(testAttrs[3], e1) // insert before front
	checkAttributesPointers(t, l, []*Attribute{e2, e1, e4, e3})
	l.Remove(e2)
	e2 = l.InsertBefore(testAttrs[3], e4) // insert before middle
	checkAttributesPointers(t, l, []*Attribute{e1, e2, e4, e3})
	l.Remove(e2)
	e2 = l.InsertBefore(testAttrs[3], e3) // insert before back
	checkAttributesPointers(t, l, []*Attribute{e1, e4, e2, e3})
	l.Remove(e2)

	e2 = l.InsertAfter(testAttrs[3], e1) // insert after front
	checkAttributesPointers(t, l, []*Attribute{e1, e2, e4, e3})
	l.Remove(e2)
	e2 = l.InsertAfter(testAttrs[3], e4) // insert after middle
	checkAttributesPointers(t, l, []*Attribute{e1, e4, e2, e3})
	l.Remove(e2)
	e2 = l.InsertAfter(testAttrs[3], e3) // insert after back
	checkAttributesPointers(t, l, []*Attribute{e1, e4, e3, e2})
	l.Remove(e2)

	// Check standard iteration.
	var sum int64
	for e := l.Front(); e != nil; e = e.Next() {
		if e.Value.Value.Type() == attribute.INT64 {
			sum += e.Value.Value.AsInt64()
		}
	}
	if sum != 4 {
		t.Errorf("sum over l = %d, want 4", sum)
	}

	// Clear all elements by iterating
	var next *Attribute
	for e := l.Front(); e != nil; e = next {
		next = e.Next()
		l.Remove(e)
	}
	checkAttributesPointers(t, l, []*Attribute{})
}

func checkAttributes(t *testing.T, l *Attributes, attrs []attribute.KeyValue) {
	if !checkAttributesLen(t, l, len(attrs)) {
		return
	}

	i := 0
	for e := l.Front(); e != nil; e = e.Next() {
		le := e.Value
		if le != attrs[i] {
			t.Errorf("elt[%d].Value = %v, want %v", i, le, attrs[i])
		}
		i++
	}
}

func TestExtending(t *testing.T) {
	l1 := New()
	l2 := New()

	l1.PushBack(testAttrs[0])
	l1.PushBack(testAttrs[1])
	l1.PushBack(testAttrs[2])

	l2.PushBack(testAttrs[3])
	l2.PushBack(testAttrs[4])

	l3 := New()
	l3.PushBackList(l1)
	checkAttributes(t, l3, testAttrs[:3])
	l3.PushBackList(l2)
	checkAttributes(t, l3, testAttrs[:5])

	l3 = New()
	l3.PushFrontList(l2)
	checkAttributes(t, l3, testAttrs[3:])
	l3.PushFrontList(l1)
	checkAttributes(t, l3, testAttrs[:5])

	checkAttributes(t, l1, testAttrs[:3])
	checkAttributes(t, l2, testAttrs[3:5])

	l3 = New()
	l3.PushBackList(l1)
	checkAttributes(t, l3, testAttrs[:3])
	l3.PushBackList(l3)
	want := []attribute.KeyValue{
		testAttrs[0], testAttrs[1], testAttrs[2],
		testAttrs[0], testAttrs[1], testAttrs[2],
	}
	checkAttributes(t, l3, want)

	l3 = New()
	l3.PushFrontList(l1)
	checkAttributes(t, l3, testAttrs[:3])
	l3.PushFrontList(l3)
	checkAttributes(t, l3, want)

	l3 = New()
	l1.PushBackList(l3)
	checkAttributes(t, l1, testAttrs[:3])
	l1.PushFrontList(l3)
	checkAttributes(t, l1, testAttrs[:3])
}

func TestRemove(t *testing.T) {
	l := New()
	e1 := l.PushBack(testAttrs[0])
	e2 := l.PushBack(testAttrs[1])
	checkAttributesPointers(t, l, []*Attribute{e1, e2})
	e := l.Front()
	l.Remove(e)
	checkAttributesPointers(t, l, []*Attribute{e2})
	l.Remove(e)
	checkAttributesPointers(t, l, []*Attribute{e2})
}

func TestIssue4103(t *testing.T) {
	l1 := New()
	l1.PushBack(testAttrs[0])
	l1.PushBack(testAttrs[1])

	l2 := New()
	l2.PushBack(testAttrs[2])
	l2.PushBack(testAttrs[3])

	e := l1.Front()
	l2.Remove(e) // l2 should not change because e is not an element of l2
	if n := l2.Len(); n != 2 {
		t.Errorf("l2.Len() = %d, want 2", n)
	}

	l1.InsertBefore(testAttrs[0], e)
	if n := l1.Len(); n != 3 {
		t.Errorf("l1.Len() = %d, want 3", n)
	}
}

func TestIssue6349(t *testing.T) {
	l := New()
	l.PushBack(testAttrs[0])
	l.PushBack(testAttrs[1])

	e := l.Front()
	l.Remove(e)
	if e.Value != testAttrs[0] {
		t.Errorf("e.value = %v, want %v", e.Value, testAttrs[0])
	}
	if e.Next() != nil {
		t.Errorf("e.Next() != nil")
	}
	if e.Prev() != nil {
		t.Errorf("e.Prev() != nil")
	}
}

func TestMove(t *testing.T) {
	l := New()
	e1 := l.PushBack(testAttrs[0])
	e2 := l.PushBack(testAttrs[1])
	e3 := l.PushBack(testAttrs[2])
	e4 := l.PushBack(testAttrs[3])

	l.MoveAfter(e3, e3)
	checkAttributesPointers(t, l, []*Attribute{e1, e2, e3, e4})
	l.MoveBefore(e2, e2)
	checkAttributesPointers(t, l, []*Attribute{e1, e2, e3, e4})

	l.MoveAfter(e3, e2)
	checkAttributesPointers(t, l, []*Attribute{e1, e2, e3, e4})
	l.MoveBefore(e2, e3)
	checkAttributesPointers(t, l, []*Attribute{e1, e2, e3, e4})

	l.MoveBefore(e2, e4)
	checkAttributesPointers(t, l, []*Attribute{e1, e3, e2, e4})
	e2, e3 = e3, e2

	l.MoveBefore(e4, e1)
	checkAttributesPointers(t, l, []*Attribute{e4, e1, e2, e3})
	e1, e2, e3, e4 = e4, e1, e2, e3

	l.MoveAfter(e4, e1)
	checkAttributesPointers(t, l, []*Attribute{e1, e4, e2, e3})
	e2, e3, e4 = e4, e2, e3

	l.MoveAfter(e2, e3)
	checkAttributesPointers(t, l, []*Attribute{e1, e3, e2, e4})
}

// Test PushFront, PushBack, PushFrontList, PushBackList with uninitialized Attributes
func TestZeroAttributes(t *testing.T) {
	var l1 = new(Attributes)
	l1.PushFront(testAttrs[0])
	checkAttributes(t, l1, testAttrs[:1])

	var l2 = new(Attributes)
	l2.PushBack(testAttrs[0])
	checkAttributes(t, l2, testAttrs[:1])

	var l3 = new(Attributes)
	l3.PushFrontList(l1)
	checkAttributes(t, l3, testAttrs[:1])

	var l4 = new(Attributes)
	l4.PushBackList(l2)
	checkAttributes(t, l4, testAttrs[:1])
}

// Test that a list l is not modified when calling InsertBefore with a mark that is not an element of l.
func TestInsertBeforeUnknownMark(t *testing.T) {
	var l Attributes
	l.PushBack(testAttrs[0])
	l.PushBack(testAttrs[1])
	l.PushBack(testAttrs[2])
	l.InsertBefore(testAttrs[0], new(Attribute))
	checkAttributes(t, &l, testAttrs[:3])
}

// Test that a list l is not modified when calling InsertAfter with a mark that is not an element of l.
func TestInsertAfterUnknownMark(t *testing.T) {
	var l Attributes
	l.PushBack(testAttrs[0])
	l.PushBack(testAttrs[1])
	l.PushBack(testAttrs[2])
	l.InsertAfter(testAttrs[0], new(Attribute))
	checkAttributes(t, &l, testAttrs[:3])
}

// Test that a list l is not modified when calling MoveAfter or MoveBefore with a mark that is not an element of l.
func TestMoveUnknownMark(t *testing.T) {
	var l1 Attributes
	e1 := l1.PushBack(testAttrs[0])

	var l2 Attributes
	e2 := l2.PushBack(testAttrs[1])

	l1.MoveAfter(e1, e2)
	checkAttributes(t, &l1, testAttrs[:1])
	checkAttributes(t, &l2, testAttrs[1:2])

	l1.MoveBefore(e1, e2)
	checkAttributes(t, &l1, testAttrs[:1])
	checkAttributes(t, &l2, testAttrs[1:2])
}

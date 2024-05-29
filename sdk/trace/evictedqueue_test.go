// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
}

func TestAdd(t *testing.T) {
	q := newEvictedQueue[string](3)
	q.add("value1")
	q.add("value2")
	if wantLen, gotLen := 2, len(q.queue); wantLen != gotLen {
		t.Errorf("got queue length %d want %d", gotLen, wantLen)
	}
}

func TestCopy(t *testing.T) {
	q := newEvictedQueue[string](3)
	q.add("value1")
	cp := q.copy()

	q.add("value2")
	assert.Equal(t, []string{"value1"}, cp, "queue update modified copy")

	cp[0] = "value0"
	assert.Equal(t, "value1", q.queue[0], "copy update modified queue")
}

func TestDropCount(t *testing.T) {
	q := newEvictedQueue[string](3)
	var called bool
	q.logDropped = func() { called = true }

	q.add("value1")
	assert.False(t, called, `"value1" logged as dropped`)
	q.add("value2")
	assert.False(t, called, `"value2" logged as dropped`)
	q.add("value3")
	assert.False(t, called, `"value3" logged as dropped`)
	q.add("value1")
	assert.True(t, called, `"value2" not logged as dropped`)
	q.add("value4")
	if wantLen, gotLen := 3, len(q.queue); wantLen != gotLen {
		t.Errorf("got queue length %d want %d", gotLen, wantLen)
	}
	if wantDropCount, gotDropCount := 2, q.droppedCount; wantDropCount != gotDropCount {
		t.Errorf("got drop count %d want %d", gotDropCount, wantDropCount)
	}
	wantArr := []string{"value3", "value1", "value4"}
	gotArr := q.copy()

	if wantLen, gotLen := len(wantArr), len(gotArr); gotLen != wantLen {
		t.Errorf("got array len %d want %d", gotLen, wantLen)
	}

	if !reflect.DeepEqual(gotArr, wantArr) {
		t.Errorf("got array = %#v; want %#v", gotArr, wantArr)
	}
}

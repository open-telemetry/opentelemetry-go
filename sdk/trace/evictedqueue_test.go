// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import (
	"reflect"
	"testing"
)

func init() {
}

func TestAdd(t *testing.T) {
	q := newEvictedQueue(3)
	q.add("value1")
	q.add("value2")
	if wantLen, gotLen := 2, len(q.queue); wantLen != gotLen {
		t.Errorf("got queue length %d want %d", gotLen, wantLen)
	}
}

func (eq *evictedQueue) queueToArray() []string {
	arr := make([]string, 0)
	for _, value := range eq.queue {
		arr = append(arr, value.(string))
	}
	return arr
}

func TestDropCount(t *testing.T) {
	q := newEvictedQueue(3)
	q.add("value1")
	q.add("value2")
	q.add("value3")
	q.add("value1")
	q.add("value4")
	if wantLen, gotLen := 3, len(q.queue); wantLen != gotLen {
		t.Errorf("got queue length %d want %d", gotLen, wantLen)
	}
	if wantDropCount, gotDropCount := 2, q.droppedCount; wantDropCount != gotDropCount {
		t.Errorf("got drop count %d want %d", gotDropCount, wantDropCount)
	}
	wantArr := []string{"value3", "value1", "value4"}
	gotArr := q.queueToArray()

	if wantLen, gotLen := len(wantArr), len(gotArr); gotLen != wantLen {
		t.Errorf("got array len %d want %d", gotLen, wantLen)
	}

	if !reflect.DeepEqual(gotArr, wantArr) {
		t.Errorf("got array = %#v; want %#v", gotArr, wantArr)
	}
}

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace

import (
	"reflect"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/funcr"
	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/internal/global"
)

func init() {
}

func TestAdd(t *testing.T) {
	q := newEvictedQueueLink(3)
	q.add(Link{})
	q.add(Link{})
	if wantLen, gotLen := 2, len(q.queue); wantLen != gotLen {
		t.Errorf("got queue length %d want %d", gotLen, wantLen)
	}
}

func TestCopy(t *testing.T) {
	q := newEvictedQueueEvent(3)
	q.add(Event{Name: "value1"})
	cp := q.copy()

	q.add(Event{Name: "value2"})
	assert.Equal(t, []Event{{Name: "value1"}}, cp, "queue update modified copy")

	cp[0] = Event{Name: "value0"}
	assert.Equal(t, Event{Name: "value1"}, q.queue[0], "copy update modified queue")
}

func TestDropCount(t *testing.T) {
	q := newEvictedQueueEvent(3)

	var called int
	t.Cleanup(func(l logr.Logger) func() {
		return func() { global.SetLogger(l) }
	}(global.GetLogger()))
	global.SetLogger(funcr.New(func(prefix, args string) {
		called++
	}, funcr.Options{Verbosity: 1}))

	q.add(Event{Name: "value1"})
	assert.Equal(t, 0, called, `"value1" logged as dropped`)
	q.add(Event{Name: "value2"})
	assert.Equal(t, 0, called, `"value2" logged as dropped`)
	q.add(Event{Name: "value3"})
	assert.Equal(t, 0, called, `"value3" logged as dropped`)
	q.add(Event{Name: "value1"})
	assert.Equal(t, 1, called, `"value2" not logged as dropped`)
	q.add(Event{Name: "value4"})
	assert.Equal(t, 1, called, `"value4" logged as dropped`)
	if wantLen, gotLen := 3, len(q.queue); wantLen != gotLen {
		t.Errorf("got queue length %d want %d", gotLen, wantLen)
	}
	if wantDropCount, gotDropCount := 2, q.droppedCount; wantDropCount != gotDropCount {
		t.Errorf("got drop count %d want %d", gotDropCount, wantDropCount)
	}
	wantArr := []Event{{Name: "value3"}, {Name: "value1"}, {Name: "value4"}}
	gotArr := q.copy()

	if wantLen, gotLen := len(wantArr), len(gotArr); gotLen != wantLen {
		t.Errorf("got array len %d want %d", gotLen, wantLen)
	}

	if !reflect.DeepEqual(gotArr, wantArr) {
		t.Errorf("got array = %#v; want %#v", gotArr, wantArr)
	}
}

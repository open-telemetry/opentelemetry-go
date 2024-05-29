// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package trace // import "go.opentelemetry.io/otel/sdk/trace"

import (
	"fmt"
	"sync"

	"go.opentelemetry.io/otel/internal/global"
)

// evictedQueue is a FIFO queue with a configurable capacity.
type evictedQueue[T any] struct {
	queue        []T
	capacity     int
	droppedCount int
	logDropped   func()
}

func newEvictedQueue[T any](capacity int) evictedQueue[T] {
	var tVal T
	msg := fmt.Sprintf("limit reached: dropping trace %T", tVal)
	// Do not pre-allocate queue, do this lazily.
	return evictedQueue[T]{
		capacity:   capacity,
		logDropped: sync.OnceFunc(func() { global.Warn(msg) }),
	}
}

// add adds value to the evictedQueue eq. If eq is at capacity, the oldest
// queued value will be discarded and the drop count incremented.
func (eq *evictedQueue[T]) add(value T) {
	if eq.capacity == 0 {
		eq.droppedCount++
		eq.logDropped()
		return
	}

	if eq.capacity > 0 && len(eq.queue) == eq.capacity {
		// Drop first-in while avoiding allocating more capacity to eq.queue.
		copy(eq.queue[:eq.capacity-1], eq.queue[1:])
		eq.queue = eq.queue[:eq.capacity-1]
		eq.droppedCount++
		eq.logDropped()
	}
	eq.queue = append(eq.queue, value)
}

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"sync"
	"sync/atomic"
)

// queue holds a queue of logging records.
// When the queue becomes full, the oldest records in the queue are
// overwritten.
type queue struct {
	sync.Mutex

	dropped     atomic.Uint64
	cap, len    int
	read, write *ring
}

func newQueue(size int) *queue {
	r := newRing(size)
	return &queue{
		cap:   size,
		read:  r,
		write: r,
	}
}

// Dropped returns the number of Records dropped during enqueueing since the
// last time Dropped was called.
func (q *queue) Dropped() uint64 {
	return q.dropped.Swap(0)
}

// Enqueue adds r to the queue. The queue size, including the addition of r, is
// returned.
//
// If enqueueing r will exceed the capacity of q, the oldest Record held in q
// will be dropped and r retained.
func (q *queue) Enqueue(r Record) int {
	q.Lock()
	defer q.Unlock()

	q.write.Value = r
	q.write = q.write.Next()

	q.len++
	if q.len > q.cap {
		// Overflow. Advance read to be the new "oldest".
		q.len = q.cap
		q.read = q.read.Next()
		q.dropped.Add(1)
	}
	return q.len
}

// TryDequeue attempts to dequeue up to len(buf) Records. The available Records
// will be assigned into buf and passed to write. If write fails, returning
// false, the Records will not be removed from the queue. If write succeeds,
// returning true, the dequeued Records are removed from the queue. The number
// of Records remaining in the queue are returned.
//
// When write is called the lock of q is held. The write function must not call
// other methods of this q that acquire the lock.
func (q *queue) TryDequeue(buf []Record, write func([]Record) bool) int {
	q.Lock()
	defer q.Unlock()

	origRead := q.read

	n := min(len(buf), q.len)
	for i := 0; i < n; i++ {
		buf[i] = q.read.Value
		q.read = q.read.Next()
	}

	if write(buf[:n]) {
		q.len -= n
	} else {
		q.read = origRead
	}
	return q.len
}

func (q *queue) Dequeue(buf []Record) (queueLen, written int) {
	q.Lock()
	defer q.Unlock()

	n := min(len(buf), q.len)
	for i := 0; i < n; i++ {
		buf[i] = q.read.Value
		q.read = q.read.Next()
	}

	q.len -= n
	return q.len, n
}

// Flush returns all the Records held in the queue and resets it to be
// empty.
func (q *queue) Flush() []Record {
	q.Lock()
	defer q.Unlock()

	out := make([]Record, q.len)
	for i := range out {
		out[i] = q.read.Value
		q.read = q.read.Next()
	}
	q.len = 0

	return out
}

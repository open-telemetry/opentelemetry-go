package log

import (
	"sync"
	"sync/atomic"
)

// Cyclic array queue
//
// Holds a queue of logging records.
//
// When the queue becomes full, the oldest records in the queue are
// overwritten.
type cqueue struct {
	sync.Mutex

	dropped atomic.Uint64
	cap     int

	buf  []*Record
	head int
	tail int
}

func newCQueue(size int) *cqueue {
	return &cqueue{
		cap: size,
		buf: make([]*Record, size+1),
	}
}

func (q *cqueue) Len() int {
	if q.Empty() {
		return 0
	} else if q.tail > q.head {
		return q.tail - q.head
	} else {
		return q.cap - q.head + q.tail + 1
	}
}

func (q *cqueue) Cap() int {
	return q.cap
}

func (q *cqueue) Empty() bool {
	return q.head == q.tail
}

func (q *cqueue) Full() bool {
	return q.head == q.incrPtr(q.tail)
}

func (q *cqueue) incrPtr(ptr int) int {
	return q.incrPtrN(ptr, 1)
}

func (q *cqueue) incrPtrN(ptr int, n int) int {
	return (ptr + n) % (q.Cap() + 1)
}

// Dropped returns the number of Records dropped during enqueueing since the
// last time Dropped was called.
func (q *cqueue) Dropped() uint64 {
	return q.dropped.Swap(0)
}

// Encqueue adds r to the queue. The queue size, including the addition of r, is
// returned.
//
// If enqueueing r will exceed the capacity of q, the oldest Record held in q
// will be dropped and r retained.
func (q *cqueue) Enqueue(r *Record) int {
	q.Lock()
	defer q.Unlock()

	if q.Full() {
		q.head = q.incrPtr(q.head)
		q.dropped.Add(1)
	}

	q.buf[q.tail] = r
	q.tail = q.incrPtr(q.tail)

	return q.Len()
}

// TryDequeue attempts to dequeue up to len(buf) Records. The available Records
// will be assigned into buf and passed to write. If write fails, returning
// false, the Records will not be removed from the queue. If write succeeds,
// returning true, the dequeued Records are removed from the queue. The number
// of Records remaining in the queue are returned.
//
// When write is called the lock of q is held. The write function must not call
// other methods of this q that acquire the lock.
func (q *cqueue) TryDequeue(buf []*Record, write func([]*Record) bool) int {
	q.Lock()
	defer q.Unlock()

	n, head := q.cp(buf)

	if write(buf[:n]) {
		q.head = head
	}

	return q.Len()
}

func (q *cqueue) cp(buf []*Record) (n, head int) {
	n = min(len(buf), q.Len())
	head = q.head

	firstLim := 0
	secondLim := 0
	if q.tail > q.head {
		firstLim = n
		secondLim = 0
	} else {
		firstLim = q.Cap() - q.head + 1
		secondLim = q.tail

		if firstLim >= len(buf) {
			firstLim = len(buf)
			secondLim = 0
		}
	}

	// cp before looping back to start
	copy(buf[:firstLim], q.buf[head:])

	// cp after looping back to start
	copy(buf[firstLim:], q.buf[:secondLim])

	// calling incr ptr every time is expensive
	head = q.incrPtrN(head, n)

	return n, head
}

// Flush returns all the Records held in the queue and resets it to be
// empty.
func (q *cqueue) Flush() []*Record {
	q.Lock()
	defer q.Unlock()

	if q.Empty() {
		return nil
	}

	ret := make([]*Record, q.Len())
	q.cp(ret)
	q.head = 0
	q.tail = 0

	return ret
}

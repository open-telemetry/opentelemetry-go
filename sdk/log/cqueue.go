package log

import (
	"sync"
	"sync/atomic"
)

type cqueue struct {
	sync.Mutex

	dropped atomic.Uint64
	cap     int
	// read, write *ring

	// cyclical array cqueue
	buf  []*Record
	head int
	tail int
}

func newCQueue(size int) *cqueue {
	// slog.Info("q init", "cap", size)
	// r := newRing(size)
	return &cqueue{
		cap: size,
		// read:  r,
		// write: r,
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
	return (ptr + 1) % (q.Cap() + 1)
}

func (q *cqueue) incrPtrN(ptr int, n int) int {
	return (ptr + n) % (q.Cap() + 1)
}

// Dropped returns the number of Records dropped during encqueueing since the
// last time Dropped was called.
func (q *cqueue) Dropped() uint64 {
	return q.dropped.Swap(0)
}

// Encqueue adds r to the cqueue. The cqueue size, including the addition of r, is
// returned.
//
// If encqueueing r will exceed the capacity of q, the oldest Record held in q
// will be dropped and r retained.
func (q *cqueue) Enqueue(r *Record) int {
	q.Lock()
	defer q.Unlock()

	// slog.Info("q before encqueue", "len", q.Len(), "tail", q.tail, "head", q.head, "b", b)

	if q.Full() {
		// slog.Info("q full")
		q.head = q.incrPtr(q.head)
		q.dropped.Add(1)
	}

	q.buf[q.tail] = r
	q.tail = q.incrPtr(q.tail)

	// slog.Info("q after encqueue", "len", q.Len(), "tail", q.tail, "head", q.head)

	return q.Len()
}

// TryDecqueue attempts to decqueue up to len(buf) Records. The available Records
// will be assigned into buf and passed to write. If write fails, returning
// false, the Records will not be removed from the cqueue. If write succeeds,
// returning true, the decqueued Records are removed from the cqueue. The number
// of Records remaining in the cqueue are returned.
//
// When write is called the lock of q is held. The write function must not call
// other methods of this q that acquire the lock.
func (q *cqueue) TryDequeue(buf []*Record, write func([]*Record) bool) int {
	q.Lock()
	defer q.Unlock()

	// n := min(len(buf), q.Len())
	// head := q.head
	// for i := 0; i < n; i++ {
	// 	buf[i] = q.buf[head]
	// 	head = q.incrPtr(head)
	// }

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

// Flush returns all the Records held in the cqueue and resets it to be
// empty.
func (q *cqueue) Flush() []*Record {
	q.Lock()
	defer q.Unlock()

	// var ret []*Record

	if q.Empty() {
		return nil
	}

	ret := make([]*Record, q.Len())
	q.cp(ret)
	q.head = 0
	q.tail = 0

	return ret

	// if q.Empty() {
	// 	return nil
	// } else if q.tail > q.head {
	// 	ret = q.buf[q.head:q.tail]
	// 	q.head = 0
	// 	q.tail = 0
	// } else {
	// 	ret = make([]*Record, q.Len())
	// 	initLen := q.Len()
	// 	for i := 0; i < initLen; i++ {
	// 		ret[i] = q.buf[q.head]
	// 		q.head = q.incrPtr(q.head)
	// 	}
	// }
}

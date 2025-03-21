package log

import (
	"slices"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/log"
)

func TestCQueue(t *testing.T) {
	var r1, r2 Record
	r1.SetBody(log.IntValue(1))
	r2.SetBody(log.IntValue(2))

	t.Run("newQueue", func(t *testing.T) {
		const size = 1
		q := newCQueue(size)
		assert.Equal(t, 0, q.Len())
		assert.Equal(t, size, q.Cap(), "capacity")
		// assert.Equal(t, size, q.read.Len(), "read ring")
		// assert.Same(t, q.read, q.write, "different rings")
	})

	t.Run("Enqueue", func(t *testing.T) {
		const size = 2
		q := newCQueue(size)

		var notR Record
		notR.SetBody(log.IntValue(10))

		assert.Equal(t, 1, q.Enqueue(&notR), "incomplete batch")
		assert.Equal(t, 1, q.Len(), "length")
		assert.Equal(t, size, q.cap, "capacity")

		assert.Equal(t, 2, q.Enqueue(&r1), "complete batch")
		assert.Equal(t, 2, q.Len(), "length")
		assert.Equal(t, size, q.cap, "capacity")

		assert.Equal(t, 2, q.Enqueue(&r2), "overflow batch")
		assert.Equal(t, 2, q.Len(), "length")
		assert.Equal(t, size, q.Cap(), "capacity")

		assert.Equal(t, []*Record{&r1, &r2}, q.Flush(), "flushed Records")
	})

	t.Run("Enqueue Plus", func(t *testing.T) {
		const size = 100
		q := newCQueue(size)

		records := make([]*Record, 200)
		for i := range len(records) {
			r := new(Record)
			r.SetBody(log.IntValue(i + 1))
			records[i] = r
		}

		for i, r := range records {
			if i < 100 {
				assert.Equal(t, i+1, q.Enqueue(r), "incomplete batch")
				assert.Equal(t, i+1, q.Len(), "length")
			} else if i < 150 {
				assert.Equal(t, q.Cap(), q.Enqueue(r), "complete batch")
				assert.Equal(t, q.Cap(), q.Len(), "length")
			} else if i == 150 {
				assert.Equal(t, records[50:150], q.Flush(), "flushed Records")
				assert.Equal(t, 0, q.Len(), "length")

				assert.Equal(t, 1, q.Enqueue(r), "incomplete batch")
				assert.Equal(t, 1, q.Len(), "length")
			} else {
				assert.Equal(t, i-150+1, q.Enqueue(r), "incomplete batch")
				assert.Equal(t, i-150+1, q.Len(), "length")
			}
		}

		assert.Equal(t, records[150:], q.Flush(), "flushed Records")
		assert.Equal(t, 0, q.Len(), "length")
	})

	t.Run("Dropped", func(t *testing.T) {
		q := newCQueue(1)

		_ = q.Enqueue(&r1)
		_ = q.Enqueue(&r1)
		assert.Equal(t, uint64(1), q.Dropped(), "fist")

		_ = q.Enqueue(&r1)
		_ = q.Enqueue(&r1)
		assert.Equal(t, uint64(2), q.Dropped(), "second")
	})

	t.Run("Flush", func(t *testing.T) {
		const size = 2
		q := newCQueue(size)
		q.Enqueue(&r1)

		assert.Equal(t, []*Record{&r1}, q.Flush(), "flushed")
	})

	t.Run("TryFlush", func(t *testing.T) {
		const size = 3
		q := newCQueue(size)
		for i := 0; i < size-1; i++ {
			q.Enqueue(&r1)
		}

		buf := make([]*Record, 1)
		f := func([]*Record) bool { return false }
		assert.Equal(t, size-1, q.TryDequeue(buf, f), "not flushed")
		require.Equal(t, size-1, q.Len(), "length")
		// require.NotSame(t, q.read, q.write, "read ring advanced")

		var flushed []*Record
		f = func(r []*Record) bool {
			flushed = append(flushed, r...)
			return true
		}
		if assert.Equal(t, size-2, q.TryDequeue(buf, f), "did not flush len(buf)") {
			assert.Equal(t, []*Record{&r1}, flushed, "Records")
		}

		buf = slices.Grow(buf, size)
		flushed = flushed[:0]
		if assert.Equal(t, 0, q.TryDequeue(buf, f), "did not flush len(queue)") {
			assert.Equal(t, []*Record{&r1}, flushed, "Records")
		}
	})

	t.Run("ConcurrentSafe", func(t *testing.T) {
		const goRoutines = 10

		flushed := make(chan []*Record, goRoutines)
		out := make([]*Record, 0, goRoutines)
		done := make(chan struct{})
		go func() {
			defer close(done)
			for recs := range flushed {
				out = append(out, recs...)
			}
		}()

		var wg sync.WaitGroup
		wg.Add(goRoutines)

		b := newCQueue(goRoutines)
		for i := 0; i < goRoutines; i++ {
			go func() {
				defer wg.Done()
				b.Enqueue(&Record{})
				flushed <- b.Flush()
			}()
		}

		wg.Wait()
		close(flushed)
		<-done

		assert.Len(t, out, goRoutines, "flushed Records")
	})
}

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"bytes"
	"context"
	stdlog "log"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
	"unsafe"

	"github.com/go-logr/stdr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/log"
)

type concurrentBuffer struct {
	b bytes.Buffer
	m sync.Mutex
}

func (b *concurrentBuffer) Write(p []byte) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.Write(p)
}

func (b *concurrentBuffer) String() string {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.String()
}

func TestEmptyBatchConfig(t *testing.T) {
	assert.NotPanics(t, func() {
		var bp BatchProcessor
		ctx := context.Background()
		record := new(Record)
		assert.NoError(t, bp.OnEmit(ctx, record), "OnEmit")
		assert.NoError(t, bp.ForceFlush(ctx), "ForceFlush")
		assert.NoError(t, bp.Shutdown(ctx), "Shutdown")
	})
}

func TestNewBatchConfig(t *testing.T) {
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		t.Log(err)
	}))

	testcases := []struct {
		name    string
		envars  map[string]string
		options []BatchProcessorOption
		want    batchConfig
	}{
		{
			name: "Defaults",
			want: batchConfig{
				maxQSize:        newSetting(dfltMaxQSize),
				expInterval:     newSetting(dfltExpInterval),
				expTimeout:      newSetting(dfltExpTimeout),
				expMaxBatchSize: newSetting(dfltExpMaxBatchSize),
				expBufferSize:   newSetting(dfltExpBufferSize),
			},
		},
		{
			name: "Options",
			options: []BatchProcessorOption{
				WithMaxQueueSize(10),
				WithExportInterval(time.Microsecond),
				WithExportTimeout(time.Hour),
				WithExportMaxBatchSize(2),
				WithExportBufferSize(3),
			},
			want: batchConfig{
				maxQSize:        newSetting(10),
				expInterval:     newSetting(time.Microsecond),
				expTimeout:      newSetting(time.Hour),
				expMaxBatchSize: newSetting(2),
				expBufferSize:   newSetting(3),
			},
		},
		{
			name: "Environment",
			envars: map[string]string{
				envarMaxQSize:        strconv.Itoa(10),
				envarExpInterval:     strconv.Itoa(100),
				envarExpTimeout:      strconv.Itoa(1000),
				envarExpMaxBatchSize: strconv.Itoa(1),
			},
			want: batchConfig{
				maxQSize:        newSetting(10),
				expInterval:     newSetting(100 * time.Millisecond),
				expTimeout:      newSetting(1000 * time.Millisecond),
				expMaxBatchSize: newSetting(1),
				expBufferSize:   newSetting(dfltExpBufferSize),
			},
		},
		{
			name: "InvalidOptions",
			options: []BatchProcessorOption{
				WithMaxQueueSize(-11),
				WithExportInterval(-1 * time.Microsecond),
				WithExportTimeout(-1 * time.Hour),
				WithExportMaxBatchSize(-2),
				WithExportBufferSize(-2),
			},
			want: batchConfig{
				maxQSize:        newSetting(dfltMaxQSize),
				expInterval:     newSetting(dfltExpInterval),
				expTimeout:      newSetting(dfltExpTimeout),
				expMaxBatchSize: newSetting(dfltExpMaxBatchSize),
				expBufferSize:   newSetting(dfltExpBufferSize),
			},
		},
		{
			name: "InvalidEnvironment",
			envars: map[string]string{
				envarMaxQSize:        "-1",
				envarExpInterval:     "-1",
				envarExpTimeout:      "-1",
				envarExpMaxBatchSize: "-1",
			},
			want: batchConfig{
				maxQSize:        newSetting(dfltMaxQSize),
				expInterval:     newSetting(dfltExpInterval),
				expTimeout:      newSetting(dfltExpTimeout),
				expMaxBatchSize: newSetting(dfltExpMaxBatchSize),
				expBufferSize:   newSetting(dfltExpBufferSize),
			},
		},
		{
			name: "Precedence",
			envars: map[string]string{
				envarMaxQSize:        strconv.Itoa(1),
				envarExpInterval:     strconv.Itoa(100),
				envarExpTimeout:      strconv.Itoa(1000),
				envarExpMaxBatchSize: strconv.Itoa(10),
			},
			options: []BatchProcessorOption{
				// These override the environment variables.
				WithMaxQueueSize(3),
				WithExportInterval(time.Microsecond),
				WithExportTimeout(time.Hour),
				WithExportMaxBatchSize(2),
				WithExportBufferSize(2),
			},
			want: batchConfig{
				maxQSize:        newSetting(3),
				expInterval:     newSetting(time.Microsecond),
				expTimeout:      newSetting(time.Hour),
				expMaxBatchSize: newSetting(2),
				expBufferSize:   newSetting(2),
			},
		},
		{
			name: "BatchLessThanOrEqualToQSize",
			options: []BatchProcessorOption{
				WithMaxQueueSize(1),
				WithExportMaxBatchSize(10),
				WithExportBufferSize(3),
			},
			want: batchConfig{
				maxQSize:        newSetting(1),
				expInterval:     newSetting(dfltExpInterval),
				expTimeout:      newSetting(dfltExpTimeout),
				expMaxBatchSize: newSetting(1),
				expBufferSize:   newSetting(3),
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			for key, value := range tc.envars {
				t.Setenv(key, value)
			}
			assert.Equal(t, tc.want, newBatchConfig(tc.options))
		})
	}
}

func TestBatchProcessor(t *testing.T) {
	ctx := context.Background()

	t.Run("NilExporter", func(t *testing.T) {
		assert.NotPanics(t, func() { NewBatchProcessor(nil) })
	})

	t.Run("Polling", func(t *testing.T) {
		e := newTestExporter(nil)
		const size = 15
		b := NewBatchProcessor(
			e,
			WithMaxQueueSize(2*size),
			WithExportMaxBatchSize(2*size),
			WithExportInterval(time.Nanosecond),
			WithExportTimeout(time.Hour),
		)
		for range size {
			assert.NoError(t, b.OnEmit(ctx, new(Record)))
		}
		var got []Record
		assert.Eventually(t, func() bool {
			for _, r := range e.Records() {
				got = append(got, r...)
			}
			return len(got) == size
		}, 2*time.Second, time.Microsecond)
		_ = b.Shutdown(ctx)
	})

	t.Run("OnEmit", func(t *testing.T) {
		const batch = 10
		e := newTestExporter(nil)
		b := NewBatchProcessor(
			e,
			WithMaxQueueSize(10*batch),
			WithExportMaxBatchSize(batch),
			WithExportInterval(time.Hour),
			WithExportTimeout(time.Hour),
		)
		for range 10 * batch {
			assert.NoError(t, b.OnEmit(ctx, new(Record)))
		}
		assert.Eventually(t, func() bool {
			return e.ExportN() > 1
		}, 2*time.Second, time.Microsecond, "multi-batch flush")

		assert.NoError(t, b.Shutdown(ctx))
		assert.GreaterOrEqual(t, e.ExportN(), 10)
	})

	t.Run("RetriggerFlushNonBlocking", func(t *testing.T) {
		e := newTestExporter(nil)
		e.ExportTrigger = make(chan struct{})

		const batch = 10
		b := NewBatchProcessor(
			e,
			WithMaxQueueSize(3*batch),
			WithExportMaxBatchSize(batch),
			WithExportInterval(time.Hour),
			WithExportTimeout(time.Hour),
		)
		for range 2 * batch {
			assert.NoError(t, b.OnEmit(ctx, new(Record)))
		}

		var n int
		require.Eventually(t, func() bool {
			n = e.ExportN()
			return n > 0
		}, 2*time.Second, time.Microsecond, "blocked export not attempted")

		var err error
		require.Eventually(t, func() bool {
			err = b.OnEmit(ctx, new(Record))
			return true
		}, time.Second, time.Microsecond, "OnEmit blocked")
		assert.NoError(t, err)

		e.ExportTrigger <- struct{}{}
		assert.Eventually(t, func() bool {
			return e.ExportN() > n
		}, 2*time.Second, time.Microsecond, "flush not retriggered")

		close(e.ExportTrigger)
		assert.NoError(t, b.Shutdown(ctx))
		assert.Equal(t, 3, e.ExportN())
	})

	t.Run("Shutdown", func(t *testing.T) {
		t.Run("Error", func(t *testing.T) {
			e := newTestExporter(assert.AnError)
			b := NewBatchProcessor(e)
			assert.ErrorIs(t, b.Shutdown(ctx), assert.AnError, "exporter error not returned")
			assert.NoError(t, b.Shutdown(ctx))
		})

		t.Run("Multiple", func(t *testing.T) {
			e := newTestExporter(nil)
			b := NewBatchProcessor(e)

			const shutdowns = 3
			for range shutdowns {
				assert.NoError(t, b.Shutdown(ctx))
			}
			assert.Equal(t, 1, e.ShutdownN(), "exporter Shutdown calls")
		})

		t.Run("OnEmit", func(t *testing.T) {
			e := newTestExporter(nil)
			b := NewBatchProcessor(e)
			assert.NoError(t, b.Shutdown(ctx))

			want := e.ExportN()
			assert.NoError(t, b.OnEmit(ctx, new(Record)))
			assert.Equal(t, want, e.ExportN(), "Export called after shutdown")
		})

		t.Run("ForceFlush", func(t *testing.T) {
			e := newTestExporter(nil)
			b := NewBatchProcessor(e)

			assert.NoError(t, b.OnEmit(ctx, new(Record)))
			assert.NoError(t, b.Shutdown(ctx))

			assert.NoError(t, b.ForceFlush(ctx))
			assert.Equal(t, 0, e.ForceFlushN(), "ForceFlush called after shutdown")
		})

		t.Run("CanceledContext", func(t *testing.T) {
			e := newTestExporter(nil)
			e.ExportTrigger = make(chan struct{})
			t.Cleanup(func() { close(e.ExportTrigger) })
			b := NewBatchProcessor(e)

			ctx := context.Background()
			c, cancel := context.WithCancel(ctx)
			cancel()

			assert.ErrorIs(t, b.Shutdown(c), context.Canceled)
		})
	})

	t.Run("ForceFlush", func(t *testing.T) {
		t.Run("Flush", func(t *testing.T) {
			e := newTestExporter(assert.AnError)
			b := NewBatchProcessor(
				e,
				WithMaxQueueSize(100),
				WithExportMaxBatchSize(10),
				WithExportInterval(time.Hour),
				WithExportTimeout(time.Hour),
			)
			t.Cleanup(func() { _ = b.Shutdown(ctx) })

			r := new(Record)
			r.SetBody(log.BoolValue(true))
			require.NoError(t, b.OnEmit(ctx, r))

			assert.ErrorIs(t, b.ForceFlush(ctx), assert.AnError, "exporter error not returned")
			assert.Equal(t, 1, e.ForceFlushN(), "exporter ForceFlush calls")
			if assert.Equal(t, 1, e.ExportN(), "exporter Export calls") {
				got := e.Records()
				if assert.Len(t, got[0], 1, "records received") {
					assert.Equal(t, *r, got[0][0])
				}
			}
		})

		t.Run("ErrorPartialFlush", func(t *testing.T) {
			e := newTestExporter(nil)
			e.ExportTrigger = make(chan struct{})

			ctxErrCalled := make(chan struct{})
			orig := ctxErr
			ctxErr = func(ctx context.Context) error {
				close(ctxErrCalled)
				return orig(ctx)
			}
			t.Cleanup(func() { ctxErr = orig })

			const batch = 1
			b := NewBatchProcessor(
				e,
				WithMaxQueueSize(10*batch),
				WithExportMaxBatchSize(batch),
				WithExportInterval(time.Hour),
				WithExportTimeout(time.Hour),
			)

			// Enqueue 10 x "batch size" amount of records.
			for range 10 * batch {
				require.NoError(t, b.OnEmit(ctx, new(Record)))
			}
			assert.Eventually(t, func() bool {
				return e.ExportN() > 0 && len(b.exporter.input) == cap(b.exporter.input)
			}, 2*time.Second, time.Microsecond)
			// 1 export being performed, 1 export in buffer chan, >1 batch
			// still in queue that an attempt to flush will be made on.
			//
			// Stop the poll routine to prevent contention with the queue lock.
			// This is outside of "normal" operations, but we are testing if
			// ForceFlush will return the correct error when an EnqueueExport
			// fails and not if ForceFlush will ever get the queue lock in high
			// throughput situations.
			close(b.pollDone)
			<-b.pollDone

			// Cancel the flush ctx from the start so errPartialFlush is
			// returned right away.
			fCtx, cancel := context.WithCancel(ctx)
			cancel()

			errCh := make(chan error, 1)
			go func() {
				errCh <- b.ForceFlush(fCtx)
				close(errCh)
			}()
			// Wait for ctxErrCalled to close before closing ExportTrigger so
			// we know the errPartialFlush will be returned in ForceFlush.
			<-ctxErrCalled
			close(e.ExportTrigger)

			err := <-errCh
			assert.ErrorIs(t, err, errPartialFlush, "partial flush error")
			assert.ErrorIs(t, err, context.Canceled, "ctx canceled error")
		})

		t.Run("CanceledContext", func(t *testing.T) {
			e := newTestExporter(nil)
			e.ExportTrigger = make(chan struct{})
			b := NewBatchProcessor(e)
			t.Cleanup(func() { _ = b.Shutdown(ctx) })

			r := new(Record)
			r.SetBody(log.BoolValue(true))
			_ = b.OnEmit(ctx, r)
			t.Cleanup(func() { _ = b.Shutdown(ctx) })
			t.Cleanup(func() { close(e.ExportTrigger) })

			c, cancel := context.WithCancel(ctx)
			cancel()
			assert.ErrorIs(t, b.ForceFlush(c), context.Canceled)
		})
	})

	t.Run("DroppedLogs", func(t *testing.T) {
		orig := global.GetLogger()
		t.Cleanup(func() { global.SetLogger(orig) })
		// Use concurrentBuffer for concurrent-safe reading.
		buf := new(concurrentBuffer)
		stdr.SetVerbosity(1)
		global.SetLogger(stdr.New(stdlog.New(buf, "", 0)))

		e := newTestExporter(nil)
		e.ExportTrigger = make(chan struct{})

		b := NewBatchProcessor(
			e,
			WithMaxQueueSize(1),
			WithExportMaxBatchSize(1),
			WithExportInterval(time.Hour),
			WithExportTimeout(time.Hour),
		)
		r := new(Record)
		// First record will be blocked by testExporter.Export
		assert.NoError(t, b.OnEmit(ctx, r), "exported record")
		require.Eventually(t, func() bool {
			return e.ExportN() > 0
		}, 2*time.Second, time.Microsecond, "blocked export not attempted")

		// Second record will be written to export queue
		assert.NoError(t, b.OnEmit(ctx, r), "export queue record")
		require.Eventually(t, func() bool {
			return len(b.exporter.input) == cap(b.exporter.input)
		}, 2*time.Second, time.Microsecond, "blocked queue read not attempted")

		// Third record will be written to BatchProcessor.q
		assert.NoError(t, b.OnEmit(ctx, r), "first queued")
		// The previous record will be dropped, as the new one will be written to BatchProcessor.q
		assert.NoError(t, b.OnEmit(ctx, r), "second queued")

		wantMsg := `"level"=1 "msg"="dropped log records" "dropped"=1`
		assert.Eventually(t, func() bool {
			return strings.Contains(buf.String(), wantMsg)
		}, 2*time.Second, time.Microsecond)

		close(e.ExportTrigger)
		_ = b.Shutdown(ctx)
	})

	t.Run("ConcurrentSafe", func(t *testing.T) {
		const goRoutines = 10

		e := newTestExporter(nil)
		b := NewBatchProcessor(e)

		ctx, cancel := context.WithCancel(ctx)
		var wg sync.WaitGroup
		for range goRoutines - 1 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					select {
					case <-ctx.Done():
						return
					default:
						assert.NoError(t, b.OnEmit(ctx, new(Record)))
						// Ignore partial flush errors.
						_ = b.ForceFlush(ctx)
					}
				}
			}()
		}

		require.Eventually(t, func() bool {
			return e.ExportN() > 0
		}, 2*time.Second, time.Microsecond, "export before shutdown")

		wg.Add(1)
		go func() {
			defer wg.Done()
			assert.NoError(t, b.Shutdown(ctx))
			cancel()
		}()

		wg.Wait()
	})
}

func TestQueue(t *testing.T) {
	var r Record
	r.SetBody(log.BoolValue(true))

	t.Run("newQueue", func(t *testing.T) {
		const size = 1
		q := newQueue(size)
		assert.Equal(t, 0, q.len)
		assert.Equal(t, size, q.cap, "capacity")
		assert.Equal(t, size, q.read.Len(), "read ring")
		assert.Same(t, q.read, q.write, "different rings")
	})

	t.Run("Enqueue", func(t *testing.T) {
		const size = 2
		q := newQueue(size)

		var notR Record
		notR.SetBody(log.IntValue(10))

		assert.Equal(t, 1, q.Enqueue(notR), "incomplete batch")
		assert.Equal(t, 1, q.len, "length")
		assert.Equal(t, size, q.cap, "capacity")

		assert.Equal(t, 2, q.Enqueue(r), "complete batch")
		assert.Equal(t, 2, q.len, "length")
		assert.Equal(t, size, q.cap, "capacity")

		assert.Equal(t, 2, q.Enqueue(r), "overflow batch")
		assert.Equal(t, 2, q.len, "length")
		assert.Equal(t, size, q.cap, "capacity")

		assert.Equal(t, []Record{r, r}, q.Flush(), "flushed Records")
	})

	t.Run("Dropped", func(t *testing.T) {
		q := newQueue(1)

		_ = q.Enqueue(r)
		_ = q.Enqueue(r)
		assert.Equal(t, uint64(1), q.Dropped(), "fist")

		_ = q.Enqueue(r)
		_ = q.Enqueue(r)
		assert.Equal(t, uint64(2), q.Dropped(), "second")
	})

	t.Run("Flush", func(t *testing.T) {
		const size = 2
		q := newQueue(size)
		q.write.Value = r
		q.write = q.write.Next()
		q.len = 1

		assert.Equal(t, []Record{r}, q.Flush(), "flushed")
	})

	t.Run("TryFlush", func(t *testing.T) {
		const size = 3
		q := newQueue(size)
		for range size - 1 {
			q.write.Value = r
			q.write = q.write.Next()
			q.len++
		}

		buf := make([]Record, 1)
		f := func([]Record) bool { return false }
		assert.Equal(t, size-1, q.TryDequeue(buf, f), "not flushed")
		require.Equal(t, size-1, q.len, "length")
		require.NotSame(t, q.read, q.write, "read ring advanced")

		var flushed []Record
		f = func(r []Record) bool {
			flushed = append(flushed, r...)
			return true
		}
		if assert.Equal(t, size-2, q.TryDequeue(buf, f), "did not flush len(buf)") {
			assert.Equal(t, []Record{r}, flushed, "Records")
		}

		buf = slices.Grow(buf, size)
		flushed = flushed[:0]
		if assert.Equal(t, 0, q.TryDequeue(buf, f), "did not flush len(queue)") {
			assert.Equal(t, []Record{r}, flushed, "Records")
		}
	})

	t.Run("ConcurrentSafe", func(t *testing.T) {
		const goRoutines = 10

		flushed := make(chan []Record, goRoutines)
		out := make([]Record, 0, goRoutines)
		done := make(chan struct{})
		go func() {
			defer close(done)
			for recs := range flushed {
				out = append(out, recs...)
			}
		}()

		var wg sync.WaitGroup
		wg.Add(goRoutines)

		b := newQueue(goRoutines)
		for range goRoutines {
			go func() {
				defer wg.Done()
				b.Enqueue(Record{})
				flushed <- b.Flush()
			}()
		}

		wg.Wait()
		close(flushed)
		<-done

		assert.Len(t, out, goRoutines, "flushed Records")
	})
}

func BenchmarkBatchProcessorOnEmit(b *testing.B) {
	r := new(Record)
	body := log.BoolValue(true)
	r.SetBody(body)

	rSize := unsafe.Sizeof(r) + unsafe.Sizeof(body)
	ctx := context.Background()
	bp := NewBatchProcessor(
		defaultNoopExporter,
		WithMaxQueueSize(b.N+1),
		WithExportMaxBatchSize(b.N+1),
		WithExportInterval(time.Hour),
		WithExportTimeout(time.Hour),
	)
	b.Cleanup(func() { _ = bp.Shutdown(ctx) })

	b.SetBytes(int64(rSize))
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var err error
		for pb.Next() {
			err = bp.OnEmit(ctx, r)
		}
		_ = err
	})
}

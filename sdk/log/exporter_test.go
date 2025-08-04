// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"context"
	"io"
	stdlog "log"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log"
)

type instruction struct {
	Record *[]Record
	Flush  chan [][]Record
}

type testExporter struct {
	// Err is the error returned by all methods of the testExporter.
	Err error
	// ExportTrigger is read from prior to returning from the Export method if
	// non-nil.
	ExportTrigger chan struct{}

	// Counts of method calls.
	exportN, shutdownN, forceFlushN *int32

	stopped atomic.Bool
	inputMu sync.Mutex
	input   chan instruction
	done    chan struct{}
}

func newTestExporter(err error) *testExporter {
	e := &testExporter{
		Err:         err,
		exportN:     new(int32),
		shutdownN:   new(int32),
		forceFlushN: new(int32),
		input:       make(chan instruction),
	}
	e.done = run(e.input)

	return e
}

func run(input chan instruction) chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)

		var records [][]Record
		for in := range input {
			if in.Record != nil {
				records = append(records, *in.Record)
			}
			if in.Flush != nil {
				cp := slices.Clone(records)
				records = records[:0]
				in.Flush <- cp
			}
		}
	}()
	return done
}

func (e *testExporter) Records() [][]Record {
	out := make(chan [][]Record, 1)
	e.input <- instruction{Flush: out}
	return <-out
}

func (e *testExporter) Export(ctx context.Context, r []Record) error {
	atomic.AddInt32(e.exportN, 1)
	if e.ExportTrigger != nil {
		select {
		case <-e.ExportTrigger:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	e.inputMu.Lock()
	defer e.inputMu.Unlock()
	if !e.stopped.Load() {
		e.input <- instruction{Record: &r}
	}
	return e.Err
}

func (e *testExporter) ExportN() int {
	return int(atomic.LoadInt32(e.exportN))
}

func (e *testExporter) Stop() {
	if e.stopped.Swap(true) {
		return
	}
	e.inputMu.Lock()
	defer e.inputMu.Unlock()

	close(e.input)
	<-e.done
}

func (e *testExporter) Shutdown(context.Context) error {
	atomic.AddInt32(e.shutdownN, 1)
	return e.Err
}

func (e *testExporter) ShutdownN() int {
	return int(atomic.LoadInt32(e.shutdownN))
}

func (e *testExporter) ForceFlush(context.Context) error {
	atomic.AddInt32(e.forceFlushN, 1)
	return e.Err
}

func (e *testExporter) ForceFlushN() int {
	return int(atomic.LoadInt32(e.forceFlushN))
}

func TestChunker(t *testing.T) {
	t.Run("ZeroSize", func(t *testing.T) {
		exp := newTestExporter(nil)
		t.Cleanup(exp.Stop)
		c := newChunkExporter(exp, 0)
		const size = 100
		_ = c.Export(context.Background(), make([]Record, size))

		assert.Equal(t, 1, exp.ExportN())
		records := exp.Records()
		assert.Len(t, records, 1)
		assert.Len(t, records[0], size)
	})

	t.Run("ForceFlush", func(t *testing.T) {
		exp := newTestExporter(nil)
		t.Cleanup(exp.Stop)
		c := newChunkExporter(exp, 0)
		_ = c.ForceFlush(context.Background())
		assert.Equal(t, 1, exp.ForceFlushN(), "ForceFlush not passed through")
	})

	t.Run("Shutdown", func(t *testing.T) {
		exp := newTestExporter(nil)
		t.Cleanup(exp.Stop)
		c := newChunkExporter(exp, 0)
		_ = c.Shutdown(context.Background())
		assert.Equal(t, 1, exp.ShutdownN(), "Shutdown not passed through")
	})

	t.Run("Chunk", func(t *testing.T) {
		exp := newTestExporter(nil)
		t.Cleanup(exp.Stop)
		c := newChunkExporter(exp, 10)
		assert.NoError(t, c.Export(context.Background(), make([]Record, 5)))
		assert.NoError(t, c.Export(context.Background(), make([]Record, 25)))

		wantLens := []int{5, 10, 10, 5}
		records := exp.Records()
		require.Len(t, records, len(wantLens), "chunks")
		for i, n := range wantLens {
			assert.Lenf(t, records[i], n, "chunk %d", i)
		}
	})

	t.Run("ExportError", func(t *testing.T) {
		exp := newTestExporter(assert.AnError)
		t.Cleanup(exp.Stop)
		c := newChunkExporter(exp, 0)
		ctx := context.Background()
		records := make([]Record, 25)
		err := c.Export(ctx, records)
		assert.ErrorIs(t, err, assert.AnError, "no chunking")

		c = newChunkExporter(exp, 10)
		err = c.Export(ctx, records)
		assert.ErrorIs(t, err, assert.AnError, "with chunking")
	})
}

func TestExportSync(t *testing.T) {
	eventuallyDone := func(t *testing.T, done chan struct{}) {
		assert.Eventually(t, func() bool {
			select {
			case <-done:
				return true
			default:
				return false
			}
		}, 2*time.Second, time.Microsecond)
	}

	t.Run("ErrorHandler", func(t *testing.T) {
		var got error
		handler := otel.ErrorHandlerFunc(func(err error) { got = err })
		otel.SetErrorHandler(handler)
		t.Cleanup(func() {
			l := stdlog.New(io.Discard, "", stdlog.LstdFlags)
			otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
				l.Print(err)
			}))
		})

		in := make(chan exportData, 1)
		exp := newTestExporter(assert.AnError)
		t.Cleanup(exp.Stop)
		done := exportSync(in, exp)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()

			in <- exportData{
				ctx:     context.Background(),
				records: make([]Record, 1),
			}
		}()

		wg.Wait()
		close(in)
		eventuallyDone(t, done)

		assert.ErrorIs(t, got, assert.AnError, "error not passed to ErrorHandler")
	})

	t.Run("ConcurrentSafe", func(t *testing.T) {
		in := make(chan exportData, 1)
		exp := newTestExporter(assert.AnError)
		t.Cleanup(exp.Stop)
		done := exportSync(in, exp)

		const goRoutines = 10
		var wg sync.WaitGroup
		wg.Add(goRoutines)
		for i := range goRoutines {
			go func(n int) {
				defer wg.Done()

				var r Record
				r.SetBody(log.IntValue(n))

				resp := make(chan error, 1)
				in <- exportData{
					ctx:     context.Background(),
					records: []Record{r},
					respCh:  resp,
				}

				assert.ErrorIs(t, <-resp, assert.AnError)
			}(i)
		}

		// Empty records should be ignored.
		in <- exportData{ctx: context.Background()}

		wg.Wait()

		close(in)
		eventuallyDone(t, done)

		assert.Equal(t, goRoutines, exp.ExportN(), "Export calls")

		want := make([]log.Value, goRoutines)
		for i := range want {
			want[i] = log.IntValue(i)
		}
		records := exp.Records()
		got := make([]log.Value, len(records))
		for i := range got {
			if assert.Len(t, records[i], 1, "number of records exported") {
				got[i] = records[i][0].Body()
			}
		}
		assert.ElementsMatch(t, want, got, "record bodies")
	})
}

func TestTimeoutExporter(t *testing.T) {
	t.Run("ZeroTimeout", func(t *testing.T) {
		exp := newTestExporter(nil)
		t.Cleanup(exp.Stop)
		e := newTimeoutExporter(exp, 0)
		assert.Same(t, exp, e)
	})

	t.Run("Timeout", func(t *testing.T) {
		trigger := make(chan struct{})
		t.Cleanup(func() { close(trigger) })

		exp := newTestExporter(nil)
		t.Cleanup(exp.Stop)
		exp.ExportTrigger = trigger
		e := newTimeoutExporter(exp, time.Nanosecond)

		out := make(chan error, 1)
		go func() {
			out <- e.Export(context.Background(), make([]Record, 1))
		}()

		var err error
		assert.Eventually(t, func() bool {
			select {
			case err = <-out:
				return true
			default:
				return false
			}
		}, 2*time.Second, time.Microsecond)

		assert.ErrorIs(t, err, context.DeadlineExceeded)
		close(out)
	})
}

func TestBufferExporter(t *testing.T) {
	t.Run("ConcurrentSafe", func(t *testing.T) {
		const goRoutines = 10

		exp := newTestExporter(nil)
		t.Cleanup(exp.Stop)
		e := newBufferExporter(exp, goRoutines)

		ctx := context.Background()
		records := make([]Record, 10)

		stop := make(chan struct{})
		var wg sync.WaitGroup
		for range goRoutines {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					select {
					case <-stop:
						return
					default:
						_ = e.EnqueueExport(records)
						_ = e.Export(ctx, records)
						_ = e.ForceFlush(ctx)
					}
				}
			}()
		}

		assert.Eventually(t, func() bool {
			return exp.ExportN() > 0
		}, 2*time.Second, time.Microsecond)

		assert.NoError(t, e.Shutdown(ctx))
		close(stop)
		wg.Wait()
	})

	t.Run("Shutdown", func(t *testing.T) {
		t.Run("Multiple", func(t *testing.T) {
			exp := newTestExporter(nil)
			t.Cleanup(exp.Stop)
			e := newBufferExporter(exp, 1)

			assert.NoError(t, e.Shutdown(context.Background()))
			assert.Equal(t, 1, exp.ShutdownN(), "first Shutdown")

			assert.NoError(t, e.Shutdown(context.Background()))
			assert.Equal(t, 1, exp.ShutdownN(), "second Shutdown")
		})

		t.Run("ContextCancelled", func(t *testing.T) {
			// Discard error logs.
			defer func(orig otel.ErrorHandler) {
				otel.SetErrorHandler(orig)
			}(otel.GetErrorHandler())
			handler := otel.ErrorHandlerFunc(func(error) {})
			otel.SetErrorHandler(handler)

			exp := newTestExporter(assert.AnError)
			t.Cleanup(exp.Stop)

			trigger := make(chan struct{})
			exp.ExportTrigger = trigger
			t.Cleanup(func() { close(trigger) })
			e := newBufferExporter(exp, 1)

			// Make sure there is something to flush.
			require.True(t, e.EnqueueExport(make([]Record, 1)))

			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			err := e.Shutdown(ctx)
			assert.ErrorIs(t, err, context.Canceled)
			assert.ErrorIs(t, err, assert.AnError)
		})

		t.Run("Error", func(t *testing.T) {
			exp := newTestExporter(assert.AnError)
			t.Cleanup(exp.Stop)

			e := newBufferExporter(exp, 1)
			assert.ErrorIs(t, e.Shutdown(context.Background()), assert.AnError)
		})
	})

	t.Run("ForceFlush", func(t *testing.T) {
		t.Run("Multiple", func(t *testing.T) {
			exp := newTestExporter(nil)
			t.Cleanup(exp.Stop)
			e := newBufferExporter(exp, 2)

			ctx := context.Background()
			records := make([]Record, 1)
			require.NoError(t, e.enqueue(ctx, records, nil), "enqueue")

			assert.NoError(t, e.ForceFlush(ctx), "ForceFlush records")
			assert.Equal(t, 1, exp.ExportN(), "Export number incremented")
			assert.Len(t, exp.Records(), 1, "exported Record batches")

			// Nothing to flush.
			assert.NoError(t, e.ForceFlush(ctx), "ForceFlush empty")
			assert.Equal(t, 1, exp.ExportN(), "Export number changed")
			assert.Empty(t, exp.Records(), "exported non-zero Records")
		})

		t.Run("ContextCancelled", func(t *testing.T) {
			exp := newTestExporter(nil)
			t.Cleanup(exp.Stop)

			trigger := make(chan struct{})
			exp.ExportTrigger = trigger
			t.Cleanup(func() { close(trigger) })
			e := newBufferExporter(exp, 1)

			ctx, cancel := context.WithCancel(context.Background())
			require.True(t, e.EnqueueExport(make([]Record, 1)))

			got := make(chan error, 1)
			go func() { got <- e.ForceFlush(ctx) }()
			require.Eventually(t, func() bool {
				return exp.ExportN() > 0
			}, 2*time.Second, time.Microsecond)
			cancel() // Canceled before export response.
			err := <-got
			assert.ErrorIs(t, err, context.Canceled, "enqueued")
			_ = e.Shutdown(ctx)

			// Zero length buffer
			e = newBufferExporter(exp, 0)
			assert.ErrorIs(t, e.ForceFlush(ctx), context.Canceled, "not enqueued")
		})

		t.Run("Error", func(t *testing.T) {
			exp := newTestExporter(assert.AnError)
			t.Cleanup(exp.Stop)

			e := newBufferExporter(exp, 1)
			assert.ErrorIs(t, e.ForceFlush(context.Background()), assert.AnError)
		})

		t.Run("Stopped", func(t *testing.T) {
			exp := newTestExporter(nil)
			t.Cleanup(exp.Stop)

			e := newBufferExporter(exp, 1)

			ctx := context.Background()
			_ = e.Shutdown(ctx)
			assert.NoError(t, e.ForceFlush(ctx))
		})
	})

	t.Run("Export", func(t *testing.T) {
		t.Run("ZeroRecords", func(t *testing.T) {
			exp := newTestExporter(nil)
			t.Cleanup(exp.Stop)
			e := newBufferExporter(exp, 1)

			assert.NoError(t, e.Export(context.Background(), nil))
			assert.Equal(t, 0, exp.ExportN())
		})

		t.Run("Multiple", func(t *testing.T) {
			exp := newTestExporter(nil)
			t.Cleanup(exp.Stop)
			e := newBufferExporter(exp, 1)

			ctx := context.Background()
			records := make([]Record, 1)
			records[0].SetBody(log.BoolValue(true))

			assert.NoError(t, e.Export(ctx, records))

			n := exp.ExportN()
			assert.Equal(t, 1, n, "first Export number")
			assert.Equal(t, [][]Record{records}, exp.Records())

			assert.NoError(t, e.Export(ctx, records))
			assert.Equal(t, n+1, exp.ExportN(), "second Export number")
			assert.Equal(t, [][]Record{records}, exp.Records())
		})

		t.Run("ContextCancelled", func(t *testing.T) {
			exp := newTestExporter(nil)
			t.Cleanup(exp.Stop)

			trigger := make(chan struct{})
			exp.ExportTrigger = trigger
			t.Cleanup(func() { close(trigger) })
			e := newBufferExporter(exp, 1)

			records := make([]Record, 1)
			ctx, cancel := context.WithCancel(context.Background())

			got := make(chan error, 1)
			go func() { got <- e.Export(ctx, records) }()
			require.Eventually(t, func() bool {
				return exp.ExportN() > 0
			}, 2*time.Second, time.Microsecond)
			cancel() // Canceled before export response.
			err := <-got
			assert.ErrorIs(t, err, context.Canceled, "enqueued")
			_ = e.Shutdown(ctx)

			// Zero length buffer
			e = newBufferExporter(exp, 0)
			assert.ErrorIs(t, e.Export(ctx, records), context.Canceled, "not enqueued")
		})

		t.Run("Error", func(t *testing.T) {
			exp := newTestExporter(assert.AnError)
			t.Cleanup(exp.Stop)

			e := newBufferExporter(exp, 1)
			ctx, records := context.Background(), make([]Record, 1)
			assert.ErrorIs(t, e.Export(ctx, records), assert.AnError)
		})

		t.Run("Stopped", func(t *testing.T) {
			exp := newTestExporter(nil)
			t.Cleanup(exp.Stop)

			e := newBufferExporter(exp, 1)

			ctx := context.Background()
			_ = e.Shutdown(ctx)
			assert.NoError(t, e.Export(ctx, make([]Record, 1)))
			assert.Equal(t, 0, exp.ExportN(), "Export called")
		})
	})

	t.Run("EnqueueExport", func(t *testing.T) {
		t.Run("ZeroRecords", func(t *testing.T) {
			exp := newTestExporter(nil)
			t.Cleanup(exp.Stop)
			e := newBufferExporter(exp, 1)

			assert.True(t, e.EnqueueExport(nil))
			e.ForceFlush(context.Background())
			assert.Equal(t, 0, exp.ExportN(), "empty batch enqueued")
		})

		t.Run("Multiple", func(t *testing.T) {
			exp := newTestExporter(nil)
			t.Cleanup(exp.Stop)
			e := newBufferExporter(exp, 2)

			records := make([]Record, 1)
			records[0].SetBody(log.BoolValue(true))

			assert.True(t, e.EnqueueExport(records))
			assert.True(t, e.EnqueueExport(records))
			e.ForceFlush(context.Background())

			n := exp.ExportN()
			assert.Equal(t, 2, n, "Export number")
			assert.Equal(t, [][]Record{records, records}, exp.Records())
		})

		t.Run("Stopped", func(t *testing.T) {
			exp := newTestExporter(nil)
			t.Cleanup(exp.Stop)
			e := newBufferExporter(exp, 1)

			_ = e.Shutdown(context.Background())
			assert.True(t, e.EnqueueExport(make([]Record, 1)))
		})
	})
}

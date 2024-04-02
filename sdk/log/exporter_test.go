// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"context"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

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

	input chan instruction
	done  chan struct{}
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
	e.input <- instruction{Record: &r}
	return e.Err
}

func (e *testExporter) ExportN() int {
	return int(atomic.LoadInt32(e.exportN))
}

func (e *testExporter) Stop() {
	close(e.input)
	<-e.done
}

func (e *testExporter) Shutdown(ctx context.Context) error {
	atomic.AddInt32(e.shutdownN, 1)
	return e.Err
}

func (e *testExporter) ShutdownN() int {
	return int(atomic.LoadInt32(e.shutdownN))
}

func (e *testExporter) ForceFlush(ctx context.Context) error {
	atomic.AddInt32(e.forceFlushN, 1)
	return e.Err
}

func (e *testExporter) ForceFlushN() int {
	return int(atomic.LoadInt32(e.forceFlushN))
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
		for i := 0; i < goRoutines; i++ {
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

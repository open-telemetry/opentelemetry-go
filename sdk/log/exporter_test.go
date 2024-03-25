// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log"
)

type testExporter struct {
	// Err is the error returned by all methods of the testExporter.
	Err error
	// ExportTrigger is read from prior to returning from the Export method if
	// non-nil.
	ExportTrigger chan struct{}

	// Counts of method calls.
	ExportN, ShutdownN, ForceFlushN int
	// Records are the Records passed to export.
	Records [][]Record
}

func (e *testExporter) Export(ctx context.Context, r []Record) error {
	e.ExportN++
	e.Records = append(e.Records, r)
	if e.ExportTrigger != nil {
		select {
		case <-e.ExportTrigger:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return e.Err
}

func (e *testExporter) Shutdown(ctx context.Context) error {
	e.ShutdownN++
	return e.Err
}

func (e *testExporter) ForceFlush(ctx context.Context) error {
	e.ForceFlushN++
	return e.Err
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
		exp := &testExporter{Err: assert.AnError}
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
		exp := &testExporter{Err: assert.AnError}
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

		assert.Equal(t, goRoutines, exp.ExportN, "Export calls")

		want := make([]log.Value, goRoutines)
		for i := range want {
			want[i] = log.IntValue(i)
		}
		got := make([]log.Value, len(exp.Records))
		for i := range got {
			if assert.Len(t, exp.Records[i], 1, "number of records exported") {
				got[i] = exp.Records[i][0].Body()
			}
		}
		assert.ElementsMatch(t, want, got, "record bodies")
	})
}

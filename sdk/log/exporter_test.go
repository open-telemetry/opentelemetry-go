// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestChunker(t *testing.T) {
	t.Run("ForceFlush", func(t *testing.T) {
		exp := &testExporter{}
		_ = chunker{Exporter: exp}.ForceFlush(context.Background())
		assert.Equal(t, 1, exp.ForceFlushN, "ForceFlush not passed through")
	})

	t.Run("Shutdown", func(t *testing.T) {
		exp := &testExporter{}
		_ = chunker{Exporter: exp}.Shutdown(context.Background())
		assert.Equal(t, 1, exp.ShutdownN, "Shutdown not passed through")
	})

	t.Run("Timeout", func(t *testing.T) {
		trigger := make(chan struct{})
		exp := &testExporter{ExportTrigger: trigger}
		c := chunker{Exporter: exp, Timeout: time.Nanosecond}

		out := make(chan error, 1)
		go func() {
			out <- c.Export(context.Background(), make([]Record, 1))
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
		close(trigger)
	})

	t.Run("Chunk", func(t *testing.T) {
		exp := &testExporter{}
		c := chunker{Exporter: exp, Size: 10}
		assert.NoError(t, c.Export(context.Background(), make([]Record, 5)))
		assert.NoError(t, c.Export(context.Background(), make([]Record, 25)))

		wantLens := []int{5, 10, 10, 5}
		require.Len(t, exp.Records, len(wantLens), "chunks")
		for i, n := range wantLens {
			assert.Lenf(t, exp.Records[i], n, "chunk %d", i)
		}
	})

	t.Run("ExportError", func(t *testing.T) {
		exp := &testExporter{Err: assert.AnError}
		c := chunker{Exporter: exp}
		ctx := context.Background()
		records := make([]Record, 25)
		err := c.Export(ctx, records)
		assert.ErrorIs(t, err, assert.AnError, "no chunking")

		c.Size = 10
		err = c.Export(ctx, records)
		assert.ErrorIs(t, err, assert.AnError, "with chunking")
	})
}

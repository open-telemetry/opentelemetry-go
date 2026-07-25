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
	"github.com/stretchr/testify/require"
)

type testExporter struct {
	// Err is the error returned by all methods of the testExporter.
	Err error
	// Method-specific errors take precedence over Err when set.
	ExportErr, ShutdownErr, ForceFlushErr error
	// ExportTrigger is read from prior to returning from the Export method if
	// non-nil.
	ExportTrigger  chan struct{}
	ExportFunc     func(context.Context, []Record) error
	ShutdownFunc   func(context.Context) error
	ForceFlushFunc func(context.Context) error

	// Counts of method calls.
	exportN, shutdownN, forceFlushN atomic.Int32

	mu      sync.Mutex
	records [][]Record

	callsMu sync.Mutex
	calls   []string
}

func (e *testExporter) Records() [][]Record {
	e.mu.Lock()
	defer e.mu.Unlock()

	out := slices.Clone(e.records)
	e.records = e.records[:0]
	return out
}

func (e *testExporter) Export(ctx context.Context, r []Record) error {
	e.recordCall("Export")
	e.exportN.Add(1)
	if e.ExportTrigger != nil {
		select {
		case <-e.ExportTrigger:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	e.mu.Lock()
	e.records = append(e.records, slices.Clone(r))
	e.mu.Unlock()
	if e.ExportFunc != nil {
		return e.ExportFunc(ctx, r)
	}
	return e.methodError(e.ExportErr)
}

func (e *testExporter) ExportN() int {
	return int(e.exportN.Load())
}

func (e *testExporter) Shutdown(ctx context.Context) error {
	e.recordCall("Shutdown")
	e.shutdownN.Add(1)
	if e.ShutdownFunc != nil {
		return e.ShutdownFunc(ctx)
	}
	return e.methodError(e.ShutdownErr)
}

func (e *testExporter) ShutdownN() int {
	return int(e.shutdownN.Load())
}

func (e *testExporter) ForceFlush(ctx context.Context) error {
	e.recordCall("ForceFlush")
	e.forceFlushN.Add(1)
	if e.ForceFlushFunc != nil {
		return e.ForceFlushFunc(ctx)
	}
	return e.methodError(e.ForceFlushErr)
}

func (e *testExporter) ForceFlushN() int {
	return int(e.forceFlushN.Load())
}

func (e *testExporter) methodError(err error) error {
	if err != nil {
		return err
	}
	return e.Err
}

func (e *testExporter) recordCall(name string) {
	e.callsMu.Lock()
	defer e.callsMu.Unlock()
	e.calls = append(e.calls, name)
}

func (e *testExporter) Calls() []string {
	e.callsMu.Lock()
	defer e.callsMu.Unlock()
	return slices.Clone(e.calls)
}

func TestChunker(t *testing.T) {
	t.Run("ZeroSize", func(t *testing.T) {
		exp := &testExporter{}
		c := newChunkExporter(exp, 0)
		const size = 100
		_ = c.Export(t.Context(), make([]Record, size))

		assert.Equal(t, 1, exp.ExportN())
		records := exp.Records()
		assert.Len(t, records, 1)
		assert.Len(t, records[0], size)
	})

	t.Run("ForceFlush", func(t *testing.T) {
		exp := &testExporter{}
		c := newChunkExporter(exp, 0)
		_ = c.ForceFlush(t.Context())
		assert.Equal(t, 1, exp.ForceFlushN(), "ForceFlush not passed through")
	})

	t.Run("Shutdown", func(t *testing.T) {
		exp := &testExporter{}
		c := newChunkExporter(exp, 0)
		_ = c.Shutdown(t.Context())
		assert.Equal(t, 1, exp.ShutdownN(), "Shutdown not passed through")
	})

	t.Run("Chunk", func(t *testing.T) {
		exp := &testExporter{}
		c := newChunkExporter(exp, 10)
		assert.NoError(t, c.Export(t.Context(), make([]Record, 5)))
		assert.NoError(t, c.Export(t.Context(), make([]Record, 25)))

		wantLens := []int{5, 10, 10, 5}
		records := exp.Records()
		require.Len(t, records, len(wantLens), "chunks")
		for i, n := range wantLens {
			assert.Lenf(t, records[i], n, "chunk %d", i)
		}
	})

	t.Run("ExportError", func(t *testing.T) {
		exp := &testExporter{Err: assert.AnError}
		c := newChunkExporter(exp, 0)
		ctx := t.Context()
		records := make([]Record, 25)
		err := c.Export(ctx, records)
		assert.ErrorIs(t, err, assert.AnError, "no chunking")

		c = newChunkExporter(exp, 10)
		err = c.Export(ctx, records)
		assert.ErrorIs(t, err, assert.AnError, "with chunking")
		assert.Equal(t, 4, exp.ExportN(), "all chunks attempted")
	})

	t.Run("CanceledContext", func(t *testing.T) {
		exp := &testExporter{}
		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		err := newChunkExporter(exp, 10).Export(ctx, make([]Record, 25))
		assert.ErrorIs(t, err, context.Canceled)
		assert.Zero(t, exp.ExportN(), "Export calls")
	})

	t.Run("CanceledBetweenChunks", func(t *testing.T) {
		exp := &testExporter{}
		ctx, cancel := context.WithCancel(t.Context())
		t.Cleanup(cancel)
		exp.ExportFunc = func(context.Context, []Record) error {
			cancel()
			return assert.AnError
		}

		c := newChunkExporter(exp, 10)
		err := c.Export(ctx, make([]Record, 25))
		assert.ErrorIs(t, err, assert.AnError)
		assert.ErrorIs(t, err, context.Canceled)
		assert.Equal(t, 1, exp.ExportN(), "later chunks skipped")
		records := exp.Records()
		require.Len(t, records, 1, "chunks")
		assert.Len(t, records[0], 10)
	})
}

func TestTimeoutExporter(t *testing.T) {
	t.Run("ZeroTimeout", func(t *testing.T) {
		exp := &testExporter{}
		e := newTimeoutExporter(exp, 0)
		assert.Same(t, exp, e)
	})

	t.Run("Timeout", func(t *testing.T) {
		trigger := make(chan struct{})
		t.Cleanup(func() { close(trigger) })

		exp := &testExporter{}
		exp.ExportTrigger = trigger
		e := newTimeoutExporter(exp, time.Nanosecond)

		out := make(chan error, 1)
		go func() {
			out <- e.Export(t.Context(), make([]Record, 1))
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

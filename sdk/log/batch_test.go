// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/log"
)

func TestBatchingProcessor(t *testing.T) {
	ctx := context.Background()

	t.Run("OnEmit", func(t *testing.T) {
		triggers := make(chan context.CancelFunc, 1)
		t.Cleanup(func(orig func(context.Context, time.Time, error) (context.Context, context.CancelFunc)) func() {
			ctxWithDeadlineCause = func(parent context.Context, _ time.Time, cause error) (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancelCause(parent)
				f := func() { cancel(cause) }
				triggers <- f
				return ctx, f
			}
			return func() { ctxWithDeadlineCause = orig }
		}(ctxWithDeadlineCause))

		e := &testExporter{}
		b := newBatchingProcessor(e, 10, 10, time.Hour, time.Hour)
		for _, r := range make([]Record, 15) {
			assert.NoError(t, b.OnEmit(ctx, r))
		}

		assert.Eventually(t, func() bool {
			return e.ExportN+len(b.exportCh) == 1
		}, 2*time.Second, time.Microsecond, e.ExportN+len(b.exportCh))
		assert.Len(t, b.batch.data, 5)
		assert.Len(t, b.batch.data, 5)

		// First trigger will have no stale data.
		for i := 0; i < 2; i++ {
			f := <-triggers
			f()
		}

		assert.Eventually(t, func() bool {
			return e.ExportN == 2
		}, 2*time.Second, time.Microsecond, e.ExportN)

		wantLens := []int{10, 5}
		require.Len(t, e.Records, len(wantLens), "chunks")
		for i, n := range wantLens {
			assert.Lenf(t, e.Records[i], n, "chunk %d", i)
		}
	})

	t.Run("Enabled", func(t *testing.T) {
		b := newBatchingProcessor(defaultNoopExporter, 0, 0, 0, 0)
		assert.True(t, b.Enabled(ctx, Record{}))

		_ = b.Shutdown(ctx)
		assert.False(t, b.Enabled(ctx, Record{}))
	})

	t.Run("Shutdown", func(t *testing.T) {
		e := &testExporter{Err: assert.AnError}
		b := newBatchingProcessor(e, 0, 0, 0, 0)

		assert.ErrorIs(t, b.Shutdown(ctx), assert.AnError, "exporter error not returned")
		assert.NoError(t, b.Shutdown(ctx))
		assert.Equal(t, 1, e.ShutdownN, "exporter Shutdown calls")
	})

	t.Run("ForceFlush", func(t *testing.T) {
		e := &testExporter{Err: assert.AnError}
		b := newBatchingProcessor(e, 10, 10, time.Hour, time.Hour)
		defer b.Shutdown(ctx)

		var r Record
		r.SetBody(log.BoolValue(true))
		b.OnEmit(ctx, r)

		assert.ErrorIs(t, b.ForceFlush(ctx), assert.AnError, "exporter error not returned")
		assert.Equal(t, 1, e.ForceFlushN, "exporter ForceFlush calls")
		if assert.Equal(t, 1, e.ExportN, "exporter Export calls") {
			if assert.Len(t, e.Records[0], 1, "records received") {
				assert.Equal(t, r, e.Records[0][0])
			}
		}
	})
}

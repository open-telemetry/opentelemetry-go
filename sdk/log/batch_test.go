// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
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
		b := NewBatchingProcessor(
			e,
			WithMaxQueueSize(10),
			WithExportMaxBatchSize(10),
			WithExportInterval(time.Hour),
			WithExportTimeout(time.Hour),
		)
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
		b := NewBatchingProcessor(defaultNoopExporter)
		assert.True(t, b.Enabled(ctx, Record{}))

		_ = b.Shutdown(ctx)
		assert.False(t, b.Enabled(ctx, Record{}))
	})

	t.Run("Shutdown", func(t *testing.T) {
		e := &testExporter{Err: assert.AnError}
		b := NewBatchingProcessor(e)

		assert.ErrorIs(t, b.Shutdown(ctx), assert.AnError, "exporter error not returned")
		assert.NoError(t, b.Shutdown(ctx))
		assert.Equal(t, 1, e.ShutdownN, "exporter Shutdown calls")
	})

	t.Run("ForceFlush", func(t *testing.T) {
		e := &testExporter{Err: assert.AnError}
		b := NewBatchingProcessor(
			e,
			WithMaxQueueSize(10),
			WithExportMaxBatchSize(10),
			WithExportInterval(time.Hour),
			WithExportTimeout(time.Hour),
		)
		defer func() { _ = b.Shutdown(ctx) }()

		var r Record
		r.SetBody(log.BoolValue(true))
		require.NoError(t, b.OnEmit(ctx, r))

		assert.ErrorIs(t, b.ForceFlush(ctx), assert.AnError, "exporter error not returned")
		assert.Equal(t, 1, e.ForceFlushN, "exporter ForceFlush calls")
		if assert.Equal(t, 1, e.ExportN, "exporter Export calls") {
			if assert.Len(t, e.Records[0], 1, "records received") {
				assert.Equal(t, r, e.Records[0][0])
			}
		}
	})
}

func TestNewBatchingConfig(t *testing.T) {
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		t.Log(err)
	}))

	testcases := []struct {
		name    string
		envars  map[string]string
		options []BatchingOption
		want    batchingConfig
	}{
		{
			name: "Defaults",
			want: batchingConfig{
				maxQSize:        newSetting(dfltMaxQSize),
				expInterval:     newSetting(dfltExpInterval),
				expTimeout:      newSetting(dfltExpTimeout),
				expMaxBatchSize: newSetting(dfltExpMaxBatchSize),
			},
		},
		{
			name: "Options",
			options: []BatchingOption{
				WithMaxQueueSize(1),
				WithExportInterval(time.Microsecond),
				WithExportTimeout(time.Hour),
				WithExportMaxBatchSize(2),
			},
			want: batchingConfig{
				maxQSize:        newSetting(1),
				expInterval:     newSetting(time.Microsecond),
				expTimeout:      newSetting(time.Hour),
				expMaxBatchSize: newSetting(2),
			},
		},
		{
			name: "Environment",
			envars: map[string]string{
				envarMaxQSize:        strconv.Itoa(1),
				envarExpInterval:     strconv.Itoa(100),
				envarExpTimeout:      strconv.Itoa(1000),
				envarExpMaxBatchSize: strconv.Itoa(10),
			},
			want: batchingConfig{
				maxQSize:        newSetting(1),
				expInterval:     newSetting(100 * time.Millisecond),
				expTimeout:      newSetting(1000 * time.Millisecond),
				expMaxBatchSize: newSetting(10),
			},
		},
		{
			name: "InvalidOptions",
			options: []BatchingOption{
				WithMaxQueueSize(-11),
				WithExportInterval(-1 * time.Microsecond),
				WithExportTimeout(-1 * time.Hour),
				WithExportMaxBatchSize(-2),
			},
			want: batchingConfig{
				maxQSize:        newSetting(dfltMaxQSize),
				expInterval:     newSetting(dfltExpInterval),
				expTimeout:      newSetting(dfltExpTimeout),
				expMaxBatchSize: newSetting(dfltExpMaxBatchSize),
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
			want: batchingConfig{
				maxQSize:        newSetting(dfltMaxQSize),
				expInterval:     newSetting(dfltExpInterval),
				expTimeout:      newSetting(dfltExpTimeout),
				expMaxBatchSize: newSetting(dfltExpMaxBatchSize),
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
			options: []BatchingOption{
				// These override the environment variables.
				WithMaxQueueSize(3),
				WithExportInterval(time.Microsecond),
				WithExportTimeout(time.Hour),
				WithExportMaxBatchSize(2),
			},
			want: batchingConfig{
				maxQSize:        newSetting(3),
				expInterval:     newSetting(time.Microsecond),
				expTimeout:      newSetting(time.Hour),
				expMaxBatchSize: newSetting(2),
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			for key, value := range tc.envars {
				t.Setenv(key, value)
			}
			assert.Equal(t, tc.want, newBatchingConfig(tc.options))
		})
	}
}

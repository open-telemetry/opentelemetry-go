// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log"
)

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

func TestBatchingProcessor(t *testing.T) {
	ctx := context.Background()

	newExp := func(err error) (exp *testExporter, cleanup func()) {
		exp = newTestExporter(err)
		orig := enqueueFunc
		enqueueFunc = func(ctx context.Context, r []Record, ch chan error) {
			err := exp.Export(ctx, r)
			if ch != nil {
				ch <- err
			}
		}
		return exp, func() {
			exp.Stop()
			enqueueFunc = orig
		}
	}

	t.Run("OnEmit", func(t *testing.T) {
		e, cleanup := newExp(nil)
		t.Cleanup(cleanup)

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

		assert.Equal(t, 1, e.ExportN())
		assert.Len(t, b.batch.data, 5)
	})

	t.Run("Enabled", func(t *testing.T) {
		b := NewBatchingProcessor(defaultNoopExporter)
		assert.True(t, b.Enabled(ctx, Record{}))

		_ = b.Shutdown(ctx)
		assert.False(t, b.Enabled(ctx, Record{}))
	})

	t.Run("Shutdown", func(t *testing.T) {
		t.Run("Error", func(t *testing.T) {
			e, cleanup := newExp(assert.AnError)
			t.Cleanup(cleanup)
			b := NewBatchingProcessor(e)
			assert.ErrorIs(t, b.Shutdown(ctx), assert.AnError, "exporter error not returned")
			assert.NoError(t, b.Shutdown(ctx))
		})

		t.Run("Multiple", func(t *testing.T) {
			e, cleanup := newExp(nil)
			t.Cleanup(cleanup)
			b := NewBatchingProcessor(e)

			const shutdowns = 3
			for i := 0; i < shutdowns; i++ {
				assert.NoError(t, b.Shutdown(ctx))
			}
			assert.Equal(t, 1, e.ShutdownN(), "exporter Shutdown calls")
		})

		t.Run("OnEmit", func(t *testing.T) {
			e, cleanup := newExp(nil)
			t.Cleanup(cleanup)
			b := NewBatchingProcessor(e)
			assert.NoError(t, b.Shutdown(ctx))

			want := e.ExportN()
			assert.NoError(t, b.OnEmit(ctx, Record{}))
			assert.Equal(t, want, e.ExportN(), "Export called after shutdown")
		})

		t.Run("ForceFlush", func(t *testing.T) {
			e, cleanup := newExp(nil)
			t.Cleanup(cleanup)
			b := NewBatchingProcessor(e)

			assert.NoError(t, b.OnEmit(ctx, Record{}))
			assert.NoError(t, b.Shutdown(ctx))

			assert.NoError(t, b.ForceFlush(ctx))
			assert.Equal(t, 0, e.ForceFlushN(), "ForceFlush called after shutdown")
		})

		t.Run("CanceledContext", func(t *testing.T) {
			trigger := make(chan error)
			t.Cleanup(func() { close(trigger) })

			e := newTestExporter(nil)
			orig := enqueueFunc
			enqueueFunc = func(_ context.Context, _ []Record, ch chan error) {
				go func() { ch <- <-trigger }()
			}
			t.Cleanup(func() {
				e.Stop()
				enqueueFunc = orig
			})

			b := NewBatchingProcessor(e)

			c, cancel := context.WithCancel(ctx)
			cancel()

			assert.ErrorIs(t, b.Shutdown(c), context.Canceled)
		})
	})

	t.Run("ForceFlush", func(t *testing.T) {
		e, cleanup := newExp(assert.AnError)
		t.Cleanup(cleanup)

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
		assert.Equal(t, 1, e.ForceFlushN(), "exporter ForceFlush calls")
		if assert.Equal(t, 1, e.ExportN(), "exporter Export calls") {
			got := e.Records()
			if assert.Len(t, got[0], 1, "records received") {
				assert.Equal(t, r, got[0][0])
			}
		}
	})

	t.Run("ConcurrentSafe", func(t *testing.T) {
		const goRoutines = 10

		e, cleanup := newExp(nil)
		t.Cleanup(cleanup)

		b := NewBatchingProcessor(e)
		stop := make(chan struct{})
		var wg sync.WaitGroup
		for i := 0; i < goRoutines-1; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					select {
					case <-stop:
						return
					default:
						assert.NoError(t, b.OnEmit(ctx, Record{}))
						assert.NoError(t, b.ForceFlush(ctx))
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
			close(stop)
		}()

		wg.Wait()
	})
}

func TestBatch(t *testing.T) {
	var r Record
	r.SetBody(log.BoolValue(true))

	t.Run("newBatch", func(t *testing.T) {
		const size = 1
		b := newBatch(size)
		assert.Len(t, b.data, 0)
		assert.Equal(t, size, cap(b.data), "capacity")
	})

	t.Run("Append", func(t *testing.T) {
		const size = 2
		b := newBatch(size)

		assert.Nil(t, b.Append(r), "incomplete batch")
		require.Len(t, b.data, 1)
		assert.Equal(t, r, b.data[0])

		got := b.Append(r)
		assert.Len(t, b.data, 0)
		assert.Equal(t, size, cap(b.data), "capacity")
		assert.Equal(t, []Record{r, r}, got, "flushed")
	})

	t.Run("Flush", func(t *testing.T) {
		const size = 2
		b := newBatch(size)
		b.data = append(b.data, r)

		got := b.Flush()
		assert.Len(t, b.data, 0)
		assert.Equal(t, size, cap(b.data), "capacity")
		assert.Equal(t, []Record{r}, got, "flushed")
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

		b := newBatch(goRoutines)
		for i := 0; i < goRoutines; i++ {
			go func() {
				defer wg.Done()
				assert.Nil(t, b.Append(r))
				flushed <- b.Flush()
			}()
		}

		wg.Wait()
		close(flushed)
		<-done

		assert.Len(t, out, goRoutines, "flushed Records")
	})
}

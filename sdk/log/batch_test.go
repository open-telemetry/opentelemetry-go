// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log // import "go.opentelemetry.io/otel/sdk/log"

import (
	"context"
	"slices"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
)

func TestBatchingProcessorEnabled(t *testing.T) {
	b := NewBatchingProcessor(nil)
	t.Cleanup(func() {
		assert.NoError(t, b.Shutdown(context.Background()))
	})
	assert.True(t, b.Enabled(context.Background(), Record{}))
}

func TestBatchingProcessorShutdown(t *testing.T) {
	e := new(exporter)
	b := NewBatchingProcessor(e)
	_ = b.Shutdown(context.Background())
	require.True(t, e.shutdownCalled, "exporter Shutdown not called")
}

func TestBatchingProcessorForceFlush(t *testing.T) {
	e := new(exporter)
	b := NewBatchingProcessor(e)
	_ = b.ForceFlush(context.Background())
	require.True(t, e.forceFlushCalled, "exporter ForceFlush not called")
}

func TestBatchingProcessorShutdownCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	b := NewBatchingProcessor(nil)
	err := b.Shutdown(ctx)
	require.ErrorIs(t, err, context.Canceled)
}

func TestBatchingProcessorForceFlushCancled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	b := NewBatchingProcessor(nil)
	err := b.ForceFlush(ctx)
	require.ErrorIs(t, err, context.Canceled)
}

func TestBatchingProcessorOnEmit(t *testing.T) {
	e := new(exporter)
	b := NewBatchingProcessor(e)
	t.Cleanup(func() {
		assert.NoError(t, b.Shutdown(context.Background()))
	})

	var r Record
	r.SetSeverityText("test")
	_ = b.OnEmit(context.Background(), r)
	require.False(t, e.exportCalled, "exporter Export called")
}

func TestBatchingProcessorOnEmitForceFlush(t *testing.T) {
	e := new(exporter)
	b := NewBatchingProcessor(e)
	t.Cleanup(func() {
		assert.NoError(t, b.Shutdown(context.Background()))
	})

	var r Record
	r.SetSeverityText("test")
	_ = b.OnEmit(context.Background(), r)
	_ = b.ForceFlush(context.Background())
	require.True(t, e.exportCalled, "exporter Export not called")
	assert.Equal(t, []Record{r}, e.records)
}

func TestBatchingProcessorOnEmitTick(t *testing.T) {
	e := new(syncExporter)
	s := NewBatchingProcessor(e, WithExportInterval(time.Millisecond))
	t.Cleanup(func() {
		assert.NoError(t, s.Shutdown(context.Background()))
	})

	var r Record
	r.SetSeverityText("test")
	_ = s.OnEmit(context.Background(), r)

	require.Eventually(t, func() bool { return e.exportCalled.Load() }, time.Second, time.Millisecond, "exporter Export not called")
	assert.Equal(t, []Record{r}, e.Records())
}

func TestBatchingProcessorOnFullQueue(t *testing.T) {
	e := new(exporter)
	b := NewBatchingProcessor(e, WithMaxQueueSize(1))
	t.Cleanup(func() {
		assert.NoError(t, b.Shutdown(context.Background()))
	})

	var r Record
	r.SetSeverityText("test")
	_ = b.OnEmit(context.Background(), r)
	var r2 Record
	r2.SetSeverityText("dropped")
	_ = b.OnEmit(context.Background(), r2)

	_ = b.ForceFlush(context.Background())
	assert.Equal(t, []Record{r}, e.records)
}

func TestBatchingProcessorOnFullBatch(t *testing.T) {
	e := new(exporter)
	b := NewBatchingProcessor(e, WithExportMaxBatchSize(1))
	t.Cleanup(func() {
		assert.NoError(t, b.Shutdown(context.Background()))
	})

	var r Record
	r.SetSeverityText("test")
	_ = b.OnEmit(context.Background(), r)
	var r2 Record
	r2.SetSeverityText("on next export")
	_ = b.OnEmit(context.Background(), r2)

	_ = b.ForceFlush(context.Background())
	assert.Equal(t, []Record{r}, e.records)

	_ = b.ForceFlush(context.Background())
	assert.Equal(t, []Record{r, r2}, e.records)
}

func TestBatchingProcessorConcurrentSafe(t *testing.T) {
	const goRoutineN = 10

	var wg sync.WaitGroup
	wg.Add(goRoutineN)

	var r Record
	r.SetSeverityText("test")
	ctx := context.Background()
	s := NewBatchingProcessor(nil)
	t.Cleanup(func() {
		assert.NoError(t, s.Shutdown(context.Background()))
	})
	for i := 0; i < goRoutineN; i++ {
		go func() {
			defer wg.Done()

			_ = s.OnEmit(ctx, r)
			_ = s.Enabled(ctx, r)
			_ = s.Shutdown(ctx)
			_ = s.ForceFlush(ctx)
		}()
	}

	wg.Wait()
}

func TestNewBatchingProcessorConfiguration(t *testing.T) {
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		t.Log(err)
	}))

	testcases := []struct {
		name                   string
		envars                 map[string]string
		options                []BatchingOption
		wantExporter           Exporter
		wantMaxQueueSize       int
		wantExportInterval     time.Duration
		wantExportTimeout      time.Duration
		wantExportMaxBatchSize int
	}{
		{
			name:                   "Defaults",
			wantExporter:           defaultNoopExporter,
			wantMaxQueueSize:       dfltMaxQSize,
			wantExportInterval:     dfltExpInterval,
			wantExportTimeout:      dfltExpTimeout,
			wantExportMaxBatchSize: dfltExpMaxBatchSize,
		},
		{
			name: "Options",
			options: []BatchingOption{
				WithMaxQueueSize(1),
				WithExportInterval(time.Microsecond),
				WithExportTimeout(time.Hour),
				WithExportMaxBatchSize(2),
			},
			wantExporter:           defaultNoopExporter,
			wantMaxQueueSize:       1,
			wantExportInterval:     time.Microsecond,
			wantExportTimeout:      time.Hour,
			wantExportMaxBatchSize: 2,
		},
		{
			name: "Environment",
			envars: map[string]string{
				envarMaxQSize:        strconv.Itoa(1),
				envarExpInterval:     strconv.Itoa(100),
				envarExpTimeout:      strconv.Itoa(1000),
				envarExpMaxBatchSize: strconv.Itoa(10),
			},
			wantExporter:           defaultNoopExporter,
			wantMaxQueueSize:       1,
			wantExportInterval:     100 * time.Millisecond,
			wantExportTimeout:      1000 * time.Millisecond,
			wantExportMaxBatchSize: 10,
		},
		{
			name: "InvalidOptions",
			options: []BatchingOption{
				WithMaxQueueSize(-11),
				WithExportInterval(-1 * time.Microsecond),
				WithExportTimeout(-1 * time.Hour),
				WithExportMaxBatchSize(-2),
			},
			wantExporter:           defaultNoopExporter,
			wantMaxQueueSize:       dfltMaxQSize,
			wantExportInterval:     dfltExpInterval,
			wantExportTimeout:      dfltExpTimeout,
			wantExportMaxBatchSize: dfltExpMaxBatchSize,
		},
		{
			name: "InvalidEnvironment",
			envars: map[string]string{
				envarMaxQSize:        "-1",
				envarExpInterval:     "-1",
				envarExpTimeout:      "-1",
				envarExpMaxBatchSize: "-1",
			},
			wantExporter:           defaultNoopExporter,
			wantMaxQueueSize:       dfltMaxQSize,
			wantExportInterval:     dfltExpInterval,
			wantExportTimeout:      dfltExpTimeout,
			wantExportMaxBatchSize: dfltExpMaxBatchSize,
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
			wantExporter:           defaultNoopExporter,
			wantMaxQueueSize:       3,
			wantExportInterval:     time.Microsecond,
			wantExportTimeout:      time.Hour,
			wantExportMaxBatchSize: 2,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			for key, value := range tc.envars {
				t.Setenv(key, value)
			}

			b := NewBatchingProcessor(nil, tc.options...)
			t.Cleanup(func() {
				assert.NoError(t, b.Shutdown(context.Background()))
			})
			assert.Equal(t, tc.wantExporter, b.exporter, "exporter")
			assert.Equal(t, tc.wantExportInterval, b.exportInterval, "exportInterval")
			assert.Equal(t, tc.wantExportTimeout, b.exportTimeout, "exportTimeout")
			assert.Equal(t, tc.wantExportMaxBatchSize, b.exportMaxBatchSize, "exportMaxBatchSize")
			assert.Equal(t, tc.wantMaxQueueSize, b.maxQueueSize, "maxQueueSize")
		})
	}
}

func BenchmarkBatchingProcessorOnEmit(b *testing.B) {
	var r Record
	r.SetSeverityText("test")
	ctx := context.Background()
	s := NewBatchingProcessor(nil, WithExportInterval(time.Millisecond))
	b.Cleanup(func() {
		assert.NoError(b, s.Shutdown(context.Background()))
	})

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var out error

		for pb.Next() {
			out = s.OnEmit(ctx, r)
		}

		_ = out
	})
}

type syncExporter struct {
	mu      sync.Mutex
	records []Record

	exportCalled atomic.Bool
}

func (e *syncExporter) Records() []Record {
	e.mu.Lock()
	defer e.mu.Unlock()
	return slices.Clone(e.records)
}

func (e *syncExporter) Export(_ context.Context, r []Record) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.records = r
	e.exportCalled.Store(true)
	return nil
}

func (e *syncExporter) Shutdown(context.Context) error {
	return nil
}

func (e *syncExporter) ForceFlush(context.Context) error {
	return nil
}

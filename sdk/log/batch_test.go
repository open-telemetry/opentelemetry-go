// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log

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

func TestBatchingProcessorOnEmit(t *testing.T) {
	ctx := context.Background()
	e := newTestExporter(nil)
	b := NewBatchingProcessor(
		e,
		WithMaxQueueSize(100),
		WithExportMaxBatchSize(10),
		WithExportInterval(time.Hour),
		WithExportTimeout(time.Hour),
	)
	for _, r := range make([]Record, 15) {
		assert.NoError(t, b.OnEmit(ctx, r))
	}
	assert.NoError(t, b.Shutdown(ctx))

	assert.Equal(t, 2, e.ExportN())
}

func TestBatchingProcessorEnabled(t *testing.T) {
	ctx := context.Background()
	b := NewBatchingProcessor(defaultNoopExporter)
	assert.True(t, b.Enabled(ctx, Record{}))

	_ = b.Shutdown(ctx)
	assert.False(t, b.Enabled(ctx, Record{}))
}

func TestBatchingProcessorShutdownError(t *testing.T) {
	ctx := context.Background()
	e := newTestExporter(assert.AnError)
	b := NewBatchingProcessor(e)
	assert.ErrorIs(t, b.Shutdown(ctx), assert.AnError, "exporter error not returned")
	assert.NoError(t, b.Shutdown(ctx))
}

func TestBatchingProcessorShutdownMultiple(t *testing.T) {
	ctx := context.Background()
	e := newTestExporter(nil)
	b := NewBatchingProcessor(e)

	const shutdowns = 3
	for i := 0; i < shutdowns; i++ {
		assert.NoError(t, b.Shutdown(ctx))
	}
	assert.Equal(t, 1, e.ShutdownN(), "exporter Shutdown calls")
}

func TestBatchingProcessorShutdownOnEmit(t *testing.T) {
	ctx := context.Background()
	e := newTestExporter(nil)
	b := NewBatchingProcessor(e)
	assert.NoError(t, b.Shutdown(ctx))

	want := e.ExportN()
	assert.NoError(t, b.OnEmit(ctx, Record{}))
	assert.Equal(t, want, e.ExportN(), "Export called after shutdown")
}

func TestBatchingProcessorShutdownForceFlush(t *testing.T) {
	ctx := context.Background()
	e := newTestExporter(nil)
	b := NewBatchingProcessor(e)

	assert.NoError(t, b.OnEmit(ctx, Record{}))
	assert.NoError(t, b.Shutdown(ctx))

	assert.NoError(t, b.ForceFlush(ctx))
	assert.Equal(t, 0, e.ForceFlushN(), "ForceFlush called after shutdown")
}

type triggered struct {
	Exporter

	trigger chan struct{}
}

func (t *triggered) Export(ctx context.Context, records []Record) error {
	<-t.trigger
	return t.Exporter.Export(ctx, records)
}

func TestBatchingProcessorShutdownCanceledContext(t *testing.T) {
	trigger := make(chan struct{})
	t.Cleanup(func() { close(trigger) })
	e := newTestExporter(nil)
	b := NewBatchingProcessor(&triggered{Exporter: e, trigger: trigger})

	ctx := context.Background()
	c, cancel := context.WithCancel(ctx)
	cancel()

	assert.ErrorIs(t, b.Shutdown(c), context.Canceled)
}

func TestBatchingProcessorForceFlushFlush(t *testing.T) {
	ctx := context.Background()
	e := newTestExporter(assert.AnError)
	b := NewBatchingProcessor(
		e,
		WithMaxQueueSize(100),
		WithExportMaxBatchSize(10),
		WithExportInterval(time.Hour),
		WithExportTimeout(time.Hour),
	)
	t.Cleanup(func() { _ = b.Shutdown(ctx) })

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
}

func TestBatchingProcessorForceFlushCanceledContext(t *testing.T) {
	ctx := context.Background()
	trigger := make(chan struct{})
	e := newTestExporter(nil)
	b := NewBatchingProcessor(&triggered{Exporter: e, trigger: trigger})
	t.Cleanup(func() { _ = b.Shutdown(ctx) })

	var r Record
	r.SetBody(log.BoolValue(true))
	_ = b.OnEmit(ctx, r)
	t.Cleanup(func() { _ = b.Shutdown(ctx) })
	t.Cleanup(func() { close(trigger) })

	c, cancel := context.WithCancel(ctx)
	cancel()
	assert.ErrorIs(t, b.ForceFlush(c), context.Canceled)
}

func TestBatchingProcessorConcurrentSafe(t *testing.T) {
	const goRoutines = 10

	ctx := context.Background()
	e := newTestExporter(nil)

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
}

func TestQueue(t *testing.T) {
	var r Record
	r.SetBody(log.BoolValue(true))

	t.Run("newQueue", func(t *testing.T) {
		const size = 1
		q := newQueue(size)
		assert.Equal(t, q.len, 0)
		assert.Equal(t, size, q.cap, "capacity")
		assert.Equal(t, size, q.read.Len(), "read ring")
		assert.Equal(t, size, q.write.Len(), "write ring")
	})

	t.Run("Enqueue", func(t *testing.T) {
		const size = 2
		q := newQueue(size)

		assert.Equal(t, 1, q.Enqueue(r), "incomplete batch")
		assert.Equal(t, 1, q.len, "length")
		assert.Equal(t, size, q.cap, "capacity")

		assert.Equal(t, 2, q.Enqueue(r), "incomplete batch")
		assert.Equal(t, 2, q.len, "length")
		assert.Equal(t, size, q.cap, "capacity")

		var got []Record
		q.read.Do(func(a any) {
			got = append(got, a.(Record))
		})
		assert.Equal(t, []Record{r, r}, got, "flushed")
	})

	t.Run("Flush", func(t *testing.T) {
		const size = 2
		q := newQueue(size)
		q.write.Value = r
		q.write = q.write.Next()
		q.len = 1

		assert.Equal(t, []Record{r}, q.Flush(), "flushed")
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
		for i := 0; i < goRoutines; i++ {
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

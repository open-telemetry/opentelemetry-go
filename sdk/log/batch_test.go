// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"bytes"
	"context"
	stdlog "log"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-logr/stdr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/internal/global"
)

type syncBuffer struct {
	b bytes.Buffer
	m sync.Mutex
}

func (b *syncBuffer) Write(p []byte) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.Write(p)
}

func (b *syncBuffer) String() string {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.String()
}

func TestEmptyBatchConfig(t *testing.T) {
	assert.NotPanics(t, func() {
		var bp BatchProcessor
		ctx := t.Context()
		record := new(Record)
		assert.NoError(t, bp.OnEmit(ctx, record), "OnEmit")
		assert.NoError(t, bp.ForceFlush(ctx), "ForceFlush")
		assert.NoError(t, bp.Shutdown(ctx), "Shutdown")
	})
}

func TestNewBatchConfig(t *testing.T) {
	testcases := []struct {
		name    string
		envars  map[string]string
		options []BatchProcessorOption
		want    batchConfig
	}{
		{
			name: "Defaults",
			want: batchConfig{
				maxQSize:        newSetting(dfltMaxQSize),
				expInterval:     newSetting(dfltExpInterval),
				expTimeout:      newSetting(dfltExpTimeout),
				expMaxBatchSize: newSetting(dfltExpMaxBatchSize),
			},
		},
		{
			name: "Options",
			options: []BatchProcessorOption{
				WithMaxQueueSize(10),
				WithExportInterval(time.Microsecond),
				WithExportTimeout(time.Hour),
				WithExportMaxBatchSize(2),
				WithExportBufferSize(3),
			},
			want: batchConfig{
				maxQSize:        newSetting(10),
				expInterval:     newSetting(time.Microsecond),
				expTimeout:      newSetting(time.Hour),
				expMaxBatchSize: newSetting(2),
			},
		},
		{
			name: "Environment",
			envars: map[string]string{
				envarMaxQSize:        strconv.Itoa(10),
				envarExpInterval:     strconv.Itoa(100),
				envarExpTimeout:      strconv.Itoa(1000),
				envarExpMaxBatchSize: strconv.Itoa(1),
			},
			want: batchConfig{
				maxQSize:        newSetting(10),
				expInterval:     newSetting(100 * time.Millisecond),
				expTimeout:      newSetting(1000 * time.Millisecond),
				expMaxBatchSize: newSetting(1),
			},
		},
		{
			name: "InvalidOptions",
			options: []BatchProcessorOption{
				WithMaxQueueSize(-11),
				WithExportInterval(-1 * time.Microsecond),
				WithExportTimeout(-1 * time.Hour),
				WithExportMaxBatchSize(-2),
				WithExportBufferSize(-2),
			},
			want: batchConfig{
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
			want: batchConfig{
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
			options: []BatchProcessorOption{
				// These override the environment variables.
				WithMaxQueueSize(3),
				WithExportInterval(time.Microsecond),
				WithExportTimeout(time.Hour),
				WithExportMaxBatchSize(2),
				WithExportBufferSize(2),
			},
			want: batchConfig{
				maxQSize:        newSetting(3),
				expInterval:     newSetting(time.Microsecond),
				expTimeout:      newSetting(time.Hour),
				expMaxBatchSize: newSetting(2),
			},
		},
		{
			name: "BatchLessThanOrEqualToQSize",
			options: []BatchProcessorOption{
				WithMaxQueueSize(1),
				WithExportMaxBatchSize(10),
				WithExportBufferSize(3),
			},
			want: batchConfig{
				maxQSize:        newSetting(1),
				expInterval:     newSetting(dfltExpInterval),
				expTimeout:      newSetting(dfltExpTimeout),
				expMaxBatchSize: newSetting(1),
			},
		},
		{
			name: "DefaultBatchLessThanOrEqualToOptionQSize",
			options: []BatchProcessorOption{
				WithMaxQueueSize(1),
			},
			want: batchConfig{
				maxQSize:        newSetting(1),
				expInterval:     newSetting(dfltExpInterval),
				expTimeout:      newSetting(dfltExpTimeout),
				expMaxBatchSize: newSetting(1),
			},
		},
		{
			name: "DefaultBatchLessThanOrEqualToEnvironmentQSize",
			envars: map[string]string{
				envarMaxQSize: "1",
			},
			want: batchConfig{
				maxQSize:        newSetting(1),
				expInterval:     newSetting(dfltExpInterval),
				expTimeout:      newSetting(dfltExpTimeout),
				expMaxBatchSize: newSetting(1),
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			for key, value := range tc.envars {
				t.Setenv(key, value)
			}
			assert.Equal(t, tc.want, newBatchConfig(tc.options))
		})
	}
}

func TestBatchProcessor(t *testing.T) {
	ctx := t.Context()

	t.Run("NilExporter", func(t *testing.T) {
		var b *BatchProcessor
		assert.NotPanics(t, func() { b = NewBatchProcessor(nil) })
		defer func() { assert.NoError(t, b.Shutdown(t.Context())) }()
	})

	t.Run("Polling", func(t *testing.T) {
		e := newTestExporter(nil)
		const size = 15
		b := NewBatchProcessor(
			e,
			WithMaxQueueSize(2*size),
			WithExportMaxBatchSize(2*size),
			WithExportInterval(time.Nanosecond),
			WithExportTimeout(time.Hour),
		)
		defer func() { assert.NoError(t, b.Shutdown(t.Context())) }()
		for range size {
			assert.NoError(t, b.OnEmit(ctx, new(Record)))
		}
		var got []Record
		assert.Eventually(t, func() bool {
			for _, r := range e.Records() {
				got = append(got, r...)
			}
			return len(got) == size
		}, 2*time.Second, time.Microsecond)
		_ = b.Shutdown(ctx)
	})

	t.Run("Enabled", func(t *testing.T) {
		e := newTestExporter(nil)
		b := NewBatchProcessor(e)
		defer func() { assert.NoError(t, b.Shutdown(t.Context())) }()
		enabled := b.Enabled(ctx, EnabledParameters{})
		assert.True(t, enabled, "Enabled should return true")
	})

	t.Run("OnEmit", func(t *testing.T) {
		const batch = 10
		e := newTestExporter(nil)
		b := NewBatchProcessor(
			e,
			WithMaxQueueSize(10*batch),
			WithExportMaxBatchSize(batch),
			WithExportInterval(time.Hour),
			WithExportTimeout(time.Hour),
		)
		defer func() { assert.NoError(t, b.Shutdown(t.Context())) }()
		for range 10 * batch {
			assert.NoError(t, b.OnEmit(ctx, new(Record)))
		}
		assert.Eventually(t, func() bool {
			return e.ExportN() > 1
		}, 2*time.Second, time.Microsecond, "multi-batch flush")

		assert.NoError(t, b.Shutdown(ctx))
		assert.GreaterOrEqual(t, e.ExportN(), 10)
	})

	t.Run("RetriggerFlushNonBlocking", func(t *testing.T) {
		e := newTestExporter(nil)
		e.ExportTrigger = make(chan struct{})
		t.Cleanup(e.Stop)
		var release sync.Once
		releaseExport := func() { release.Do(func() { close(e.ExportTrigger) }) }

		const batch = 10
		b := NewBatchProcessor(
			e,
			WithMaxQueueSize(3*batch),
			WithExportMaxBatchSize(batch),
			WithExportInterval(time.Hour),
			WithExportTimeout(time.Hour),
		)
		defer func() { assert.NoError(t, b.Shutdown(t.Context())) }()
		defer releaseExport()
		for range 2 * batch {
			assert.NoError(t, b.OnEmit(ctx, new(Record)))
		}

		var n int
		require.Eventually(t, func() bool {
			n = e.ExportN()
			return n > 0
		}, 2*time.Second, time.Microsecond, "blocked export not attempted")

		var err error
		require.Eventually(t, func() bool {
			err = b.OnEmit(ctx, new(Record))
			return true
		}, time.Second, time.Microsecond, "OnEmit blocked")
		assert.NoError(t, err)

		e.ExportTrigger <- struct{}{}
		assert.Eventually(t, func() bool {
			return e.ExportN() > n
		}, 2*time.Second, time.Microsecond, "flush not retriggered")

		releaseExport()
		assert.NoError(t, b.Shutdown(ctx))
		assert.Equal(t, 3, e.ExportN())
	})

	t.Run("Shutdown", func(t *testing.T) {
		t.Run("Error", func(t *testing.T) {
			e := newTestExporter(assert.AnError)
			b := NewBatchProcessor(e)
			defer func() { assert.NoError(t, b.Shutdown(t.Context())) }()
			assert.ErrorIs(t, b.Shutdown(ctx), assert.AnError, "exporter error not returned")
			assert.NoError(t, b.Shutdown(ctx))
		})

		t.Run("Multiple", func(t *testing.T) {
			e := newTestExporter(nil)
			b := NewBatchProcessor(e)
			defer func() { assert.NoError(t, b.Shutdown(t.Context())) }()

			const shutdowns = 3
			for range shutdowns {
				assert.NoError(t, b.Shutdown(ctx))
			}
			assert.Equal(t, 1, e.ShutdownN(), "exporter Shutdown calls")
		})

		t.Run("OnEmit", func(t *testing.T) {
			e := newTestExporter(nil)
			b := NewBatchProcessor(e)
			defer func() { assert.NoError(t, b.Shutdown(t.Context())) }()
			assert.NoError(t, b.Shutdown(ctx))

			want := e.ExportN()
			assert.NoError(t, b.OnEmit(ctx, new(Record)))
			assert.Equal(t, want, e.ExportN(), "Export called after shutdown")
		})

		t.Run("ForceFlush", func(t *testing.T) {
			e := newTestExporter(nil)
			b := NewBatchProcessor(e)
			defer func() { assert.NoError(t, b.Shutdown(t.Context())) }()

			assert.NoError(t, b.OnEmit(ctx, new(Record)))
			assert.NoError(t, b.Shutdown(ctx))

			want := e.ForceFlushN()
			assert.NoError(t, b.ForceFlush(ctx))
			assert.Equal(t, want, e.ForceFlushN(), "ForceFlush called after shutdown")
		})

		t.Run("CanceledContext", func(t *testing.T) {
			e := newTestExporter(nil)
			e.ExportTrigger = make(chan struct{})
			var release sync.Once
			releaseExport := func() { release.Do(func() { close(e.ExportTrigger) }) }
			b := NewBatchProcessor(e)
			defer func() { assert.NoError(t, b.Shutdown(t.Context())) }()
			defer releaseExport()

			ctx := t.Context()
			c, cancel := context.WithCancel(ctx)
			cancel()

			assert.ErrorIs(t, b.Shutdown(c), context.Canceled)
		})
	})

	t.Run("ForceFlush", func(t *testing.T) {
		t.Run("Flush", func(t *testing.T) {
			e := newTestExporter(assert.AnError)
			b := NewBatchProcessor(
				e,
				WithMaxQueueSize(100),
				WithExportMaxBatchSize(10),
				WithExportInterval(time.Hour),
				WithExportTimeout(time.Hour),
			)
			defer func() { assert.ErrorIs(t, b.Shutdown(t.Context()), assert.AnError) }()

			r := new(Record)
			r.SetBody(attribute.BoolValue(true))
			require.NoError(t, b.OnEmit(ctx, r))

			assert.ErrorIs(t, b.ForceFlush(ctx), assert.AnError, "exporter error not returned")
			assert.Equal(t, 1, e.ForceFlushN(), "exporter ForceFlush calls")
			if assert.Equal(t, 1, e.ExportN(), "exporter Export calls") {
				got := e.Records()
				if assert.Len(t, got[0], 1, "records received") {
					assert.Equal(t, *r, got[0][0])
				}
			}
		})

		t.Run("CanceledContext", func(t *testing.T) {
			e := newTestExporter(nil)
			e.ExportTrigger = make(chan struct{})
			b := NewBatchProcessor(e)
			var release sync.Once
			releaseExport := func() { release.Do(func() { close(e.ExportTrigger) }) }
			defer func() { assert.NoError(t, b.Shutdown(t.Context())) }()
			defer releaseExport()

			r := new(Record)
			r.SetBody(attribute.BoolValue(true))
			_ = b.OnEmit(ctx, r)
			c, cancel := context.WithCancel(ctx)
			cancel()
			assert.ErrorIs(t, b.ForceFlush(c), context.Canceled)
		})
	})

	t.Run("DroppedLogs", func(t *testing.T) {
		orig := global.GetLogger()
		t.Cleanup(func() { global.SetLogger(orig) })
		// Use concurrentBuffer for concurrent-safe reading.
		buf := new(syncBuffer)
		stdr.SetVerbosity(1)
		global.SetLogger(stdr.New(stdlog.New(buf, "", 0)))

		e := newTestExporter(nil)
		e.ExportTrigger = make(chan struct{})
		var release sync.Once
		releaseExport := func() { release.Do(func() { close(e.ExportTrigger) }) }

		b := NewBatchProcessor(
			e,
			WithMaxQueueSize(1),
			WithExportMaxBatchSize(1),
			WithExportInterval(time.Hour),
			WithExportTimeout(time.Hour),
		)
		defer func() { assert.NoError(t, b.Shutdown(t.Context())) }()
		defer releaseExport()
		r := new(Record)
		// First record will be blocked by testExporter.Export
		assert.NoError(t, b.OnEmit(ctx, r), "exported record")
		require.Eventually(t, func() bool {
			return e.ExportN() > 0
		}, 2*time.Second, time.Microsecond, "blocked export not attempted")

		// The second record is queued while Export is blocked. The third
		// replaces it because the queue is full.
		assert.NoError(t, b.OnEmit(ctx, r), "first queued")
		assert.NoError(t, b.OnEmit(ctx, r), "second queued")

		releaseExport()

		wantMsg := `"level"=1 "msg"="dropped log records" "dropped"=1`
		assert.Eventually(t, func() bool {
			return strings.Contains(buf.String(), wantMsg)
		}, 2*time.Second, time.Microsecond)

		_ = b.Shutdown(ctx)
	})
}

func TestBatchProcessorBackpressureDoesNotPoll(t *testing.T) {
	ctx := t.Context()
	blocked := make(chan struct{})
	e := newTestExporter(nil)
	e.ExportTrigger = blocked
	t.Cleanup(e.Stop)

	b := NewBatchProcessor(
		e,
		WithMaxQueueSize(4),
		WithExportMaxBatchSize(1),
		WithExportInterval(time.Hour),
		WithExportTimeout(time.Hour),
	)
	var release sync.Once
	releaseExport := func() { release.Do(func() { close(blocked) }) }
	defer func() { assert.NoError(t, b.Shutdown(t.Context())) }()
	defer releaseExport()

	require.NoError(t, b.OnEmit(ctx, new(Record)))
	require.Eventually(t, func() bool {
		return e.ExportN() == 1
	}, 2*time.Second, time.Microsecond, "export did not block")

	for range 4 {
		require.NoError(t, b.OnEmit(ctx, new(Record)))
	}
	require.Eventually(t, func() bool {
		return b.q.Len() > 0
	}, 2*time.Second, time.Microsecond, "queue did not retain records")

	// Dropped is sampled by the processing goroutine. It must remain untouched
	// while that goroutine is blocked in Export; consuming it proves the
	// processor is polling without being able to make progress.
	b.q.dropped.Store(1)
	assert.Never(t, func() bool {
		return b.q.dropped.Load() == 0
	}, 25*time.Millisecond, 100*time.Microsecond)

	b.q.dropped.Store(0)
	releaseExport()
	require.NoError(t, b.Shutdown(ctx))
}

func TestBatchProcessorForceFlushExportError(t *testing.T) {
	e := newTestExporter(nil)
	e.ExportFunc = func(context.Context, []Record) error {
		return assert.AnError
	}
	t.Cleanup(e.Stop)

	b := NewBatchProcessor(
		e,
		WithExportMaxBatchSize(10),
		WithExportInterval(time.Hour),
	)
	defer func() { assert.NoError(t, b.Shutdown(t.Context())) }()

	require.NoError(t, b.OnEmit(t.Context(), new(Record)))
	err := b.ForceFlush(t.Context())
	assert.ErrorIs(t, err, assert.AnError)
	assert.Equal(t, 1, e.ExportN())
	assert.Equal(t, 1, e.ForceFlushN())
}

func TestBatchProcessorCanceledFlushRetainsQueue(t *testing.T) {
	e := newTestExporter(nil)
	t.Cleanup(e.Stop)
	b := NewBatchProcessor(
		e,
		WithMaxQueueSize(2),
		WithExportMaxBatchSize(2),
		WithExportInterval(time.Hour),
	)
	defer func() { assert.NoError(t, b.Shutdown(t.Context())) }()
	require.NoError(t, b.OnEmit(t.Context(), new(Record)))

	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	assert.ErrorIs(t, b.flushExporter(ctx), context.Canceled)
	assert.Equal(t, 1, b.q.Len(), "queued record removed")
	assert.Zero(t, e.ExportN(), "Export calls")
	assert.Zero(t, e.ForceFlushN(), "ForceFlush calls")
}

func TestBatchProcessorForceFlushAttemptsAllBatches(t *testing.T) {
	e := newTestExporter(nil)
	e.ExportFunc = func(context.Context, []Record) error {
		if e.ExportN() == 1 {
			return assert.AnError
		}
		return nil
	}
	t.Cleanup(e.Stop)

	const (
		batchSize = 2
		records   = 3 * batchSize
	)
	b := NewBatchProcessor(
		e,
		WithMaxQueueSize(records),
		WithExportMaxBatchSize(batchSize),
		WithExportInterval(time.Hour),
	)
	defer func() { assert.NoError(t, b.Shutdown(t.Context())) }()

	for range records {
		_, accepted := b.q.Enqueue(Record{})
		require.True(t, accepted)
	}
	err := b.ForceFlush(t.Context())
	assert.ErrorIs(t, err, assert.AnError)
	assert.Equal(t, records/batchSize, e.ExportN())
	for i, batch := range e.Records() {
		assert.LessOrEqualf(t, len(batch), batchSize, "batch %d", i)
	}
}

func TestBatchProcessorShutdownIncludesForceFlush(t *testing.T) {
	var (
		eventsMu sync.Mutex
		events   []string
	)
	addEvent := func(event string) {
		eventsMu.Lock()
		defer eventsMu.Unlock()
		events = append(events, event)
	}

	e := newTestExporter(nil)
	e.ExportFunc = func(context.Context, []Record) error {
		addEvent("export")
		return nil
	}
	e.ForceFlushFunc = func(context.Context) error {
		addEvent("force flush")
		return nil
	}
	e.ShutdownFunc = func(context.Context) error {
		addEvent("shutdown")
		return nil
	}
	t.Cleanup(e.Stop)

	b := NewBatchProcessor(
		e,
		WithExportMaxBatchSize(10),
		WithExportInterval(time.Hour),
	)
	defer func() { assert.NoError(t, b.Shutdown(t.Context())) }()
	require.NoError(t, b.OnEmit(t.Context(), new(Record)))
	require.NoError(t, b.Shutdown(t.Context()))

	eventsMu.Lock()
	defer eventsMu.Unlock()
	assert.Equal(t, []string{"export", "force flush", "shutdown"}, events)
}

func TestBatchProcessorForceFlushErrorWithConcurrentShutdown(t *testing.T) {
	exportStarted := make(chan struct{})
	releaseExport := make(chan struct{})
	var release sync.Once

	e := newTestExporter(nil)
	e.ExportFunc = func(context.Context, []Record) error {
		close(exportStarted)
		<-releaseExport
		return assert.AnError
	}
	t.Cleanup(e.Stop)

	b := NewBatchProcessor(
		e,
		WithExportInterval(time.Hour),
	)
	unblock := func() { release.Do(func() { close(releaseExport) }) }
	defer func() { assert.NoError(t, b.Shutdown(t.Context())) }()
	defer unblock()
	require.NoError(t, b.OnEmit(t.Context(), new(Record)))

	flushErr := make(chan error, 1)
	go func() { flushErr <- b.ForceFlush(t.Context()) }()
	<-exportStarted

	shutdownErr := make(chan error, 1)
	go func() { shutdownErr <- b.Shutdown(t.Context()) }()
	require.Eventually(t, func() bool {
		return b.stopped.Load() && len(b.shutdown) == 1
	}, time.Second, time.Microsecond, "shutdown request not queued")
	unblock()

	assert.ErrorIs(t, <-flushErr, assert.AnError)
	assert.NoError(t, <-shutdownErr)
}

func TestBatchProcessorConcurrentSafe(t *testing.T) {
	var (
		active    atomic.Int32
		maxActive atomic.Int32
	)
	call := func() {
		n := active.Add(1)
		for {
			old := maxActive.Load()
			if n <= old || maxActive.CompareAndSwap(old, n) {
				break
			}
		}
		time.Sleep(time.Microsecond)
		active.Add(-1)
	}

	e := newTestExporter(nil)
	e.ExportFunc = func(context.Context, []Record) error {
		call()
		return nil
	}
	e.ForceFlushFunc = func(context.Context) error {
		call()
		return nil
	}
	e.ShutdownFunc = func(context.Context) error {
		call()
		return nil
	}
	t.Cleanup(e.Stop)

	b := NewBatchProcessor(
		e,
		WithMaxQueueSize(64),
		WithExportMaxBatchSize(8),
		WithExportInterval(time.Millisecond),
	)
	defer func() { assert.NoError(t, b.Shutdown(t.Context())) }()

	ctx, cancel := context.WithCancel(t.Context())
	var wg sync.WaitGroup
	for range 4 {
		wg.Go(func() {
			for ctx.Err() == nil {
				_ = b.OnEmit(ctx, new(Record))
			}
		})
	}
	for range 4 {
		wg.Go(func() {
			for ctx.Err() == nil {
				_ = b.ForceFlush(ctx)
			}
		})
	}

	require.Eventually(t, func() bool {
		return e.ExportN() > 0
	}, 2*time.Second, time.Microsecond, "no export attempted")
	require.NoError(t, b.Shutdown(t.Context()))
	cancel()
	wg.Wait()

	assert.Equal(t, int32(1), maxActive.Load(), "concurrent exporter calls")
	assert.Equal(t, 1, e.ShutdownN(), "Shutdown calls")
	assert.Zero(t, b.q.Len(), "records retained after shutdown")
}

func TestQueue(t *testing.T) {
	var r Record
	r.SetBody(attribute.BoolValue(true))

	t.Run("newQueue", func(t *testing.T) {
		const size = 1
		q := newQueue(size)
		assert.Equal(t, 0, q.len)
		assert.Equal(t, size, q.cap, "capacity")
		assert.Equal(t, size, q.read.Len(), "read ring")
		assert.Same(t, q.read, q.write, "different rings")
	})

	t.Run("Enqueue", func(t *testing.T) {
		const size = 2
		q := newQueue(size)
		enqueue := func(r Record) int {
			n, accepted := q.Enqueue(r)
			assert.True(t, accepted)
			return n
		}

		var notR Record
		notR.SetBody(attribute.IntValue(10))

		assert.Equal(t, 1, enqueue(notR), "incomplete batch")
		assert.Equal(t, 1, q.len, "length")
		assert.Equal(t, size, q.cap, "capacity")

		assert.Equal(t, 2, enqueue(r), "complete batch")
		assert.Equal(t, 2, q.len, "length")
		assert.Equal(t, size, q.cap, "capacity")

		assert.Equal(t, 2, enqueue(r), "overflow batch")
		assert.Equal(t, 2, q.len, "length")
		assert.Equal(t, size, q.cap, "capacity")

		assert.Equal(t, []Record{r, r}, q.Flush(), "flushed Records")
	})

	t.Run("Dropped", func(t *testing.T) {
		q := newQueue(1)

		_, _ = q.Enqueue(r)
		_, _ = q.Enqueue(r)
		assert.Equal(t, uint64(1), q.Dropped(), "fist")

		_, _ = q.Enqueue(r)
		_, _ = q.Enqueue(r)
		assert.Equal(t, uint64(2), q.Dropped(), "second")
	})

	t.Run("Flush", func(t *testing.T) {
		const size = 2
		q := newQueue(size)
		q.write.Value = r
		q.write = q.write.Next()
		q.len = 1

		assert.Equal(t, []Record{r}, q.Flush(), "flushed")
	})

	t.Run("Dequeue", func(t *testing.T) {
		const size = 3
		q := newQueue(size)
		for range size - 1 {
			q.write.Value = r
			q.write = q.write.Next()
			q.len++
		}

		buf := make([]Record, 1)
		n, remaining := q.Dequeue(buf)
		assert.Equal(t, 1, n, "dequeued")
		assert.Equal(t, size-2, remaining, "remaining")
		assert.Equal(t, []Record{r}, buf, "records")
		assert.Equal(t, Record{}, q.read.Prev().Value, "retained record")

		buf = make([]Record, size)
		n, remaining = q.Dequeue(buf)
		assert.Equal(t, 1, n, "dequeued")
		assert.Zero(t, remaining, "remaining")
		assert.Equal(t, []Record{r}, buf[:n], "records")
	})

	t.Run("Close", func(t *testing.T) {
		q := newQueue(1)
		_, accepted := q.Enqueue(r)
		require.True(t, accepted)

		q.Close()
		assert.Equal(t, []Record{r}, q.Flush())
		_, accepted = q.Enqueue(r)
		assert.False(t, accepted)
		assert.Empty(t, q.Flush())
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
		for range goRoutines {
			go func() {
				defer wg.Done()
				_, _ = b.Enqueue(Record{})
				flushed <- b.Flush()
			}()
		}

		wg.Wait()
		close(flushed)
		<-done

		assert.Len(t, out, goRoutines, "flushed Records")
	})
}

func TestQueueConcurrentSafeClose(t *testing.T) {
	const writers = 32
	q := newQueue(writers)
	start := make(chan struct{})
	accepted := make(chan int64, writers)

	var wg sync.WaitGroup
	for i := range writers {
		wg.Go(func() {
			<-start
			var r Record
			r.SetBody(attribute.IntValue(i))
			if _, ok := q.Enqueue(r); ok {
				accepted <- int64(i)
			}
		})
	}
	flushed := make(chan []Record, 1)
	wg.Go(func() {
		<-start
		q.Close()
		flushed <- q.Flush()
	})

	close(start)
	wg.Wait()
	close(accepted)

	var want []int64
	for id := range accepted {
		want = append(want, id)
	}
	var got []int64
	for _, r := range <-flushed {
		got = append(got, r.Body().AsInt64())
	}
	assert.ElementsMatch(t, want, got)
	assert.Zero(t, q.Len())
	_, ok := q.Enqueue(Record{})
	assert.False(t, ok)
}

type blockingBenchmarkExporter struct {
	started chan struct{}
	release chan struct{}
	once    sync.Once
}

func (e *blockingBenchmarkExporter) Export(ctx context.Context, _ []Record) error {
	e.once.Do(func() { close(e.started) })
	select {
	case <-e.release:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (*blockingBenchmarkExporter) ForceFlush(context.Context) error { return nil }

func (*blockingBenchmarkExporter) Shutdown(context.Context) error { return nil }

type countingBenchmarkExporter struct {
	records    atomic.Int64
	forceFlush atomic.Int64
}

func (e *countingBenchmarkExporter) Export(_ context.Context, records []Record) error {
	e.records.Add(int64(len(records)))
	return nil
}

func (e *countingBenchmarkExporter) ForceFlush(context.Context) error {
	e.forceFlush.Add(1)
	return nil
}

func (*countingBenchmarkExporter) Shutdown(context.Context) error { return nil }

func cleanupBenchmarkBatchProcessor(b *testing.B, bp *BatchProcessor) {
	b.Helper()
	b.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.WithoutCancel(b.Context()), time.Second)
		defer cancel()
		_ = bp.Shutdown(ctx)
	})
}

func BenchmarkBatchProcessorOnEmit(b *testing.B) {
	b.Run("ExporterBlocked", benchmarkBatchProcessorOnEmitExporterBlocked)
}

func benchmarkBatchProcessorOnEmitExporterBlocked(b *testing.B) {
	r := new(Record)
	r.SetBody(attribute.BoolValue(true))

	ctx := b.Context()
	exporter := &blockingBenchmarkExporter{
		started: make(chan struct{}),
		release: make(chan struct{}),
	}
	bp := NewBatchProcessor(
		exporter,
		WithMaxQueueSize(dfltMaxQSize),
		WithExportMaxBatchSize(1),
		WithExportInterval(time.Hour),
		WithExportTimeout(time.Hour),
	)
	cleanupBenchmarkBatchProcessor(b, bp)
	b.Cleanup(func() {
		close(exporter.release)
	})
	_ = bp.OnEmit(ctx, r)
	<-exporter.started

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var err error
		for pb.Next() {
			err = bp.OnEmit(ctx, r)
		}
		_ = err
	})
}

func BenchmarkBatchProcessorForceFlush(b *testing.B) {
	b.Run("Empty", func(b *testing.B) {
		exporter := new(countingBenchmarkExporter)
		bp := NewBatchProcessor(exporter, WithExportInterval(time.Hour))
		cleanupBenchmarkBatchProcessor(b, bp)

		ctx := b.Context()
		b.ReportAllocs()
		b.ResetTimer()
		var err error
		for range b.N {
			err = bp.ForceFlush(ctx)
		}
		b.StopTimer()
		require.NoError(b, err)
		assert.Equal(b, int64(b.N), exporter.forceFlush.Load())
		assert.Zero(b, exporter.records.Load())
	})

	b.Run("Queued", func(b *testing.B) {
		exporter := new(countingBenchmarkExporter)
		bp := NewBatchProcessor(exporter, WithExportInterval(time.Hour))
		cleanupBenchmarkBatchProcessor(b, bp)

		ctx := b.Context()
		r := new(Record)
		b.ReportAllocs()
		b.ResetTimer()
		var err error
		for range b.N {
			b.StopTimer()
			err = bp.OnEmit(ctx, r)
			b.StartTimer()
			if err == nil {
				err = bp.ForceFlush(ctx)
			}
		}
		b.StopTimer()
		require.NoError(b, err)
		assert.Equal(b, int64(b.N), exporter.forceFlush.Load())
		assert.Equal(b, int64(b.N), exporter.records.Load())
	})
}

func BenchmarkBatchProcessorEmitForceFlush(b *testing.B) {
	exporter := new(countingBenchmarkExporter)
	bp := NewBatchProcessor(exporter, WithExportInterval(time.Hour))
	cleanupBenchmarkBatchProcessor(b, bp)

	ctx := b.Context()
	r := new(Record)
	b.ReportAllocs()
	b.ResetTimer()
	var err error
	for range b.N {
		err = bp.OnEmit(ctx, r)
		if err == nil {
			err = bp.ForceFlush(ctx)
		}
	}
	b.StopTimer()
	require.NoError(b, err)
	assert.Equal(b, int64(b.N), exporter.forceFlush.Load())
	assert.Equal(b, int64(b.N), exporter.records.Load())
}

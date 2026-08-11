// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"bytes"
	"context"
	"errors"
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

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/sdk"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/log/internal/counter"
	"go.opentelemetry.io/otel/sdk/log/internal/observ"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
	"go.opentelemetry.io/otel/semconv/v1.43.0/otelconv"
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
		assert.NoError(t, b.OnEmit(ctx, new(Record)))
		assert.NoError(t, b.ForceFlush(ctx))
		assert.NoError(t, b.Shutdown(ctx))
	})

	t.Run("Polling", func(t *testing.T) {
		e := &testExporter{}
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
		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			for _, r := range e.Records() {
				got = append(got, r...)
			}
			assert.Len(c, got, size)
		}, 2*time.Second, time.Microsecond)
		_ = b.Shutdown(ctx)
	})

	t.Run("Enabled", func(t *testing.T) {
		e := &testExporter{}
		b := NewBatchProcessor(e)
		defer func() { assert.NoError(t, b.Shutdown(t.Context())) }()
		enabled := b.Enabled(ctx, EnabledParameters{})
		assert.True(t, enabled, "Enabled should return true")
	})

	t.Run("OnEmit", func(t *testing.T) {
		const batch = 10
		e := &testExporter{}
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
		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			assert.Greater(c, e.ExportN(), 1)
		}, 2*time.Second, time.Microsecond, "multi-batch flush")

		assert.NoError(t, b.Shutdown(ctx))
		assert.GreaterOrEqual(t, e.ExportN(), 10)
	})

	t.Run("ScheduledExportError", func(t *testing.T) {
		original := otel.GetErrorHandler()
		handled := make(chan error, 1)
		otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
			select {
			case handled <- err:
			default:
			}
		}))
		t.Cleanup(func() { otel.SetErrorHandler(original) })

		e := &testExporter{
			ExportFunc: func(context.Context, []Record) error {
				return assert.AnError
			},
		}
		b := NewBatchProcessor(
			e,
			WithMaxQueueSize(1),
			WithExportInterval(time.Hour),
		)
		defer func() { assert.NoError(t, b.Shutdown(t.Context())) }()

		require.NoError(t, b.OnEmit(ctx, new(Record)))
		select {
		case err := <-handled:
			assert.ErrorIs(t, err, assert.AnError)
		case <-time.After(10 * time.Second):
			t.Fatal("timed out waiting for scheduled export error")
		}
	})

	t.Run("RetriggerFlushNonBlocking", func(t *testing.T) {
		e := &testExporter{}
		e.ExportTrigger = make(chan struct{})
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
		require.EventuallyWithT(t, func(c *assert.CollectT) {
			n = e.ExportN()
			assert.Positive(c, n)
		}, 2*time.Second, time.Microsecond, "blocked export not attempted")

		require.EventuallyWithT(t, func(c *assert.CollectT) {
			assert.NoError(c, b.OnEmit(ctx, new(Record)))
		}, time.Second, time.Microsecond, "OnEmit blocked")

		e.ExportTrigger <- struct{}{}
		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			assert.Greater(c, e.ExportN(), n)
		}, 2*time.Second, time.Microsecond, "flush not retriggered")

		releaseExport()
		assert.NoError(t, b.Shutdown(ctx))
		assert.Equal(t, 3, e.ExportN())
	})

	t.Run("Shutdown", func(t *testing.T) {
		t.Run("Error", func(t *testing.T) {
			e := &testExporter{Err: assert.AnError}
			b := NewBatchProcessor(e)
			defer func() { assert.NoError(t, b.Shutdown(t.Context())) }()
			assert.ErrorIs(t, b.Shutdown(ctx), assert.AnError, "exporter error not returned")
			assert.NoError(t, b.Shutdown(ctx))
		})

		t.Run("FlushesBeforeShutdown", func(t *testing.T) {
			exportErr := errors.New("export")
			forceFlushErr := errors.New("force flush")
			shutdownErr := errors.New("shutdown")
			e := &testExporter{
				ExportErr:     exportErr,
				ForceFlushErr: forceFlushErr,
				ShutdownErr:   shutdownErr,
			}
			b := NewBatchProcessor(
				e,
				WithMaxQueueSize(2),
				WithExportMaxBatchSize(2),
				WithExportInterval(time.Hour),
				WithExportTimeout(time.Hour),
			)
			require.NoError(t, b.OnEmit(ctx, new(Record)))

			err := b.Shutdown(ctx)
			assert.ErrorIs(t, err, exportErr)
			assert.ErrorIs(t, err, forceFlushErr)
			assert.ErrorIs(t, err, shutdownErr)
			assert.Equal(t, []string{"Export", "ForceFlush", "Shutdown"}, e.Calls())
		})

		t.Run("Multiple", func(t *testing.T) {
			e := &testExporter{}
			b := NewBatchProcessor(e)
			defer func() { assert.NoError(t, b.Shutdown(t.Context())) }()

			const shutdowns = 3
			for range shutdowns {
				assert.NoError(t, b.Shutdown(ctx))
			}
			assert.Equal(t, 1, e.ForceFlushN(), "exporter ForceFlush calls")
			assert.Equal(t, 1, e.ShutdownN(), "exporter Shutdown calls")
		})

		t.Run("OnEmit", func(t *testing.T) {
			e := &testExporter{}
			b := NewBatchProcessor(e)
			defer func() { assert.NoError(t, b.Shutdown(t.Context())) }()
			assert.NoError(t, b.Shutdown(ctx))

			want := e.ExportN()
			assert.NoError(t, b.OnEmit(ctx, new(Record)))
			assert.Equal(t, want, e.ExportN(), "Export called after shutdown")
		})

		t.Run("ForceFlush", func(t *testing.T) {
			e := &testExporter{}
			b := NewBatchProcessor(e)
			defer func() { assert.NoError(t, b.Shutdown(t.Context())) }()

			assert.NoError(t, b.OnEmit(ctx, new(Record)))
			assert.NoError(t, b.Shutdown(ctx))
			assert.Equal(t, 1, e.ForceFlushN(), "ForceFlush not called by Shutdown")

			want := e.ForceFlushN()
			assert.NoError(t, b.ForceFlush(ctx))
			assert.Equal(t, want, e.ForceFlushN(), "ForceFlush called after shutdown")
		})

		t.Run("CanceledContext", func(t *testing.T) {
			e := &testExporter{}
			b := NewBatchProcessor(
				e,
				WithExportInterval(time.Hour),
				WithExportTimeout(time.Hour),
			)

			ctx := t.Context()
			c, cancel := context.WithCancel(ctx)
			cancel()

			assert.ErrorIs(t, b.Shutdown(c), context.Canceled)
			require.Eventually(t, func() bool {
				select {
				case <-b.done:
					return true
				default:
					return false
				}
			}, 2*time.Second, time.Microsecond, "processor shutdown did not finish")
			assert.Equal(t, 1, e.ShutdownN(), "exporter Shutdown calls")
		})
	})

	t.Run("ForceFlush", func(t *testing.T) {
		t.Run("Flush", func(t *testing.T) {
			e := &testExporter{Err: assert.AnError}
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
			e := &testExporter{}
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

	t.Run("CanceledWhileWaiting", func(t *testing.T) {
		blockedExporter := func() (*testExporter, <-chan struct{}, func()) {
			started := make(chan struct{})
			release := make(chan struct{})
			var once sync.Once
			e := &testExporter{
				ExportFunc: func(context.Context, []Record) error {
					close(started)
					<-release
					return nil
				},
			}
			return e, started, func() { once.Do(func() { close(release) }) }
		}

		t.Run("ForceFlush", func(t *testing.T) {
			exp, exportStarted, unblock := blockedExporter()
			b := NewBatchProcessor(
				exp,
				WithExportInterval(time.Hour),
			)
			defer func() { assert.NoError(t, b.Shutdown(t.Context())) }()
			defer unblock()

			require.NoError(t, b.OnEmit(t.Context(), new(Record)))
			ctx, cancel := context.WithCancel(t.Context())
			flushErr := make(chan error, 1)
			go func() { flushErr <- b.ForceFlush(ctx) }()

			<-exportStarted
			cancel()
			assert.ErrorIs(t, <-flushErr, context.Canceled)
		})

		t.Run("Shutdown", func(t *testing.T) {
			exp, exportStarted, unblock := blockedExporter()
			shutdownDone := make(chan struct{})
			exp.ShutdownFunc = func(context.Context) error {
				close(shutdownDone)
				return nil
			}
			b := NewBatchProcessor(
				exp,
				WithMaxQueueSize(1),
				WithExportInterval(time.Hour),
			)
			defer func() { assert.NoError(t, b.Shutdown(t.Context())) }()
			defer unblock()

			require.NoError(t, b.OnEmit(t.Context(), new(Record)))
			<-exportStarted

			ctx, cancel := context.WithTimeout(t.Context(), 10*time.Millisecond)
			defer cancel()
			assert.ErrorIs(t, b.Shutdown(ctx), context.DeadlineExceeded)
			unblock()
			<-shutdownDone
		})
	})

	t.Run("DroppedLogs", func(t *testing.T) {
		orig := global.GetLogger()
		t.Cleanup(func() { global.SetLogger(orig) })
		buf := new(syncBuffer)
		stdr.SetVerbosity(1)
		global.SetLogger(stdr.New(stdlog.New(buf, "", 0)))

		e := &testExporter{}
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
		require.EventuallyWithT(t, func(c *assert.CollectT) {
			assert.Positive(c, e.ExportN())
		}, 2*time.Second, time.Microsecond, "blocked export not attempted")

		// The second record is queued while Export is blocked. The third
		// replaces it because the queue is full.
		assert.NoError(t, b.OnEmit(ctx, r), "first queued")
		assert.NoError(t, b.OnEmit(ctx, r), "second queued")

		releaseExport()

		wantMsg := `"level"=1 "msg"="dropped log records" "dropped"=1`
		assert.EventuallyWithT(t, func(c *assert.CollectT) {
			assert.Contains(c, buf.String(), wantMsg)
		}, 2*time.Second, time.Microsecond)

		_ = b.Shutdown(ctx)
	})
}

func TestBatchProcessorForceFlushExportError(t *testing.T) {
	e := &testExporter{}
	e.ExportFunc = func(context.Context, []Record) error {
		return assert.AnError
	}

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
	e := &testExporter{}
	b := NewBatchProcessor(
		e,
		WithMaxQueueSize(2),
		WithExportMaxBatchSize(2),
		WithExportInterval(time.Hour),
	)
	defer func() { assert.NoError(t, b.Shutdown(t.Context())) }()
	r := new(Record)
	r.SetBody(attribute.BoolValue(true))
	require.NoError(t, b.OnEmit(t.Context(), r))

	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	assert.ErrorIs(t, b.ForceFlush(ctx), context.Canceled)
	assert.Zero(t, e.ExportN(), "Export calls")
	assert.Zero(t, e.ForceFlushN(), "ForceFlush calls")

	require.NoError(t, b.ForceFlush(t.Context()))
	records := e.Records()
	require.Len(t, records, 1, "exported batches")
	require.Len(t, records[0], 1, "exported records")
	assert.Equal(t, *r, records[0][0])
}

func TestBatchProcessorShutdownIncludesForceFlush(t *testing.T) {
	var calls []string
	e := &testExporter{}
	e.ExportFunc = func(context.Context, []Record) error {
		calls = append(calls, "export")
		return nil
	}
	e.ForceFlushFunc = func(context.Context) error {
		calls = append(calls, "force flush")
		return nil
	}
	e.ShutdownFunc = func(context.Context) error {
		calls = append(calls, "shutdown")
		return nil
	}

	b := NewBatchProcessor(
		e,
		WithExportInterval(time.Hour),
	)
	defer func() { assert.NoError(t, b.Shutdown(t.Context())) }()
	require.NoError(t, b.OnEmit(t.Context(), new(Record)))
	require.NoError(t, b.Shutdown(t.Context()))

	assert.Equal(t, []string{"export", "force flush", "shutdown"}, calls)
}

func TestBatchProcessorForceFlushErrorShutdownConcurrentSafe(t *testing.T) {
	exportStarted := make(chan struct{})
	releaseExport := make(chan struct{})
	var release sync.Once

	e := &testExporter{}
	e.ExportFunc = func(context.Context, []Record) error {
		close(exportStarted)
		<-releaseExport
		return assert.AnError
	}

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
	unblock()

	assert.ErrorIs(t, <-flushErr, assert.AnError)
	assert.NoError(t, <-shutdownErr)
}

func TestBatchProcessorSerializesExporterCallsConcurrentSafe(t *testing.T) {
	// Track all exporter entry points with the same counters so an overlap
	// between any combination of calls raises maxActive above one.
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
		// Keep the callback active long enough for a concurrent call to enter if
		// the BatchProcessor stops serializing exporter calls.
		time.Sleep(time.Microsecond)
		active.Add(-1)
	}

	e := &testExporter{}
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

	b := NewBatchProcessor(
		e,
		WithMaxQueueSize(64),
		WithExportMaxBatchSize(8),
		WithExportInterval(time.Millisecond),
	)
	defer func() { assert.NoError(t, b.Shutdown(t.Context())) }()

	// Keep emission and flush traffic in flight so Shutdown races with work
	// already being handled instead of only exercising the stopped fast path.
	var wg sync.WaitGroup
	// Use a child context as the shared stop signal for the looping callers.
	// Canceling it also unblocks a pending ForceFlush if the test exits early.
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
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

	// Ensure the ordinary export and flush paths reached the exporter before
	// initiating Shutdown and testing serialization of the terminal call.
	require.EventuallyWithT(t, func(c *assert.CollectT) {
		assert.Positive(c, e.ExportN(), "Export calls")
		assert.Positive(c, e.ForceFlushN(), "ForceFlush calls")
	}, 10*time.Second, time.Microsecond, "exporter methods not exercised")
	require.NoError(t, b.Shutdown(t.Context()))
	cancel()
	wg.Wait()

	assert.Equal(t, int32(1), maxActive.Load(), "concurrent exporter calls")
	assert.Zero(t, active.Load(), "active exporter calls")
	assert.Equal(t, 1, e.ShutdownN(), "Shutdown calls")
}

func TestBatchProcessorConcurrentSafe(t *testing.T) {
	e := &testExporter{}
	b := NewBatchProcessor(e)
	defer func() { assert.NoError(t, b.Shutdown(t.Context())) }()

	var wg sync.WaitGroup
	for range 10 {
		wg.Go(func() {
			_ = b.OnEmit(t.Context(), new(Record))
			_ = b.ForceFlush(t.Context())
			_ = b.Shutdown(t.Context())
		})
	}
	wg.Wait()
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

func TestQueueCloseConcurrentSafe(t *testing.T) {
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

type blockingExporter struct {
	release chan struct{}
}

func (e *blockingExporter) Export(ctx context.Context, _ []Record) error {
	select {
	case <-e.release:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (*blockingExporter) ForceFlush(context.Context) error { return nil }

func (*blockingExporter) Shutdown(context.Context) error { return nil }

func BenchmarkBatchProcessorOnEmitExporterBlocked(b *testing.B) {
	exp := &blockingExporter{
		release: make(chan struct{}),
	}
	bp := NewBatchProcessor(
		exp,
		WithExportInterval(time.Hour),
		WithExportTimeout(time.Hour),
	)
	defer func() { assert.NoError(b, bp.Shutdown(b.Context())) }()
	defer close(exp.release)

	ctx := b.Context()
	r := new(Record)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = bp.OnEmit(ctx, r)
		}
	})
}

func BenchmarkBatchProcessorForceFlush(b *testing.B) {
	bp := NewBatchProcessor(noopExporter{}, WithExportInterval(time.Hour))
	defer func() { assert.NoError(b, bp.Shutdown(b.Context())) }()

	ctx := b.Context()
	r := new(Record)
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		b.StopTimer()
		_ = bp.OnEmit(ctx, r)
		b.StartTimer()
		_ = bp.ForceFlush(ctx)
	}
}

const blpComponentID int64 = 0

func TestBatchProcessorMetricsDisabled(t *testing.T) {
	t.Setenv("OTEL_GO_X_OBSERVABILITY", "false")

	origID := counter.SetExporterID(blpComponentID)
	t.Cleanup(func() { counter.SetExporterID(origID) })

	orig := otel.GetMeterProvider()
	t.Cleanup(func() { otel.SetMeterProvider(orig) })

	reader := sdkmetric.NewManualReader()
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	otel.SetMeterProvider(mp)

	e := &testExporter{}
	bp := NewBatchProcessor(
		e,
		WithMaxQueueSize(2),
		WithExportMaxBatchSize(2),
		WithExportInterval(time.Hour),
		WithExportTimeout(time.Hour),
	)
	ctx := t.Context()

	r := new(Record)
	r.SetBody(attribute.BoolValue(true))
	require.NoError(t, bp.OnEmit(ctx, r))
	require.NoError(t, bp.ForceFlush(ctx))

	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(ctx, &rm))
	for _, sm := range rm.ScopeMetrics {
		assert.NotEqual(t, observ.ScopeName, sm.Scope.Name,
			"observ metrics should not be present when disabled")
	}

	require.NoError(t, bp.Shutdown(ctx))
}

func TestBatchProcessorMetrics(t *testing.T) {
	origID := counter.SetExporterID(blpComponentID)
	t.Cleanup(func() { counter.SetExporterID(origID) })

	t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	origLogger := global.GetLogger()
	origVerbosity := stdr.SetVerbosity(1)
	t.Cleanup(func() {
		global.SetLogger(origLogger)
		stdr.SetVerbosity(origVerbosity)
	})
	buf := new(syncBuffer)
	global.SetLogger(stdr.New(stdlog.New(buf, "", 0)))

	orig := otel.GetMeterProvider()
	t.Cleanup(func() { otel.SetMeterProvider(orig) })

	reader := sdkmetric.NewManualReader()
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	otel.SetMeterProvider(mp)

	e := &testExporter{}
	e.ExportTrigger = make(chan struct{})
	var release sync.Once
	releaseExporter := func() {
		release.Do(func() { close(e.ExportTrigger) })
	}
	bp := NewBatchProcessor(
		e,
		WithMaxQueueSize(1),
		WithExportMaxBatchSize(1),
		WithExportInterval(time.Hour),
		WithExportTimeout(time.Hour),
	)
	ctx := t.Context()
	defer func() {
		releaseExporter()
		assert.NoError(t, bp.Shutdown(t.Context()))
	}()

	r := new(Record)
	r.SetBody(attribute.BoolValue(true))
	require.NoError(t, bp.OnEmit(ctx, r))
	require.Eventually(t, func() bool {
		return e.ExportN() > 0
	}, 2*time.Second, time.Microsecond, "export not started")

	assertBLPMetrics(
		t, reader,
		blpQCap(1),
		blpQSize(0),
		blpProcessed(blpDPt(blpSet(), 1)),
	)

	require.NoError(t, bp.OnEmit(ctx, r))
	require.NoError(t, bp.OnEmit(ctx, r))
	require.NoError(t, bp.OnEmit(ctx, r))

	assertBLPMetrics(
		t, reader,
		blpQCap(1),
		blpQSize(1),
		blpProcessed(blpDPt(blpSet(), 1)),
	)

	releaseExporter()
	wantMsg := `"level"=1 "msg"="dropped log records" "dropped"=2`
	require.Eventually(t, func() bool {
		return e.ExportN() == 2 && strings.Contains(buf.String(), wantMsg)
	}, 2*time.Second, time.Microsecond, "queued record and drops not processed")

	assertBLPMetrics(
		t, reader,
		blpQCap(1),
		blpQSize(0),
		blpProcessed(
			blpDPt(blpSet(), 2),
			blpDPt(blpSet(observ.ErrQueueFull), 2),
		),
	)

	require.NoError(t, bp.Shutdown(ctx))
}

func blpSet(attrs ...attribute.KeyValue) attribute.Set {
	return attribute.NewSet(append([]attribute.KeyValue{
		semconv.OTelComponentTypeBatchingLogProcessor,
		observ.BLPComponentName(blpComponentID),
	}, attrs...)...)
}

func blpDPt(set attribute.Set, value int64) metricdata.DataPoint[int64] {
	return metricdata.DataPoint[int64]{Attributes: set, Value: value}
}

func blpQCap(v int64) metricdata.Metrics {
	return metricdata.Metrics{
		Name:        otelconv.SDKProcessorLogQueueCapacity{}.Name(),
		Description: otelconv.SDKProcessorLogQueueCapacity{}.Description(),
		Unit:        otelconv.SDKProcessorLogQueueCapacity{}.Unit(),
		Data: metricdata.Sum[int64]{
			Temporality: metricdata.CumulativeTemporality,
			IsMonotonic: false,
			DataPoints:  []metricdata.DataPoint[int64]{{Attributes: blpSet(), Value: v}},
		},
	}
}

func blpQSize(v int64) metricdata.Metrics {
	return metricdata.Metrics{
		Name:        otelconv.SDKProcessorLogQueueSize{}.Name(),
		Description: otelconv.SDKProcessorLogQueueSize{}.Description(),
		Unit:        otelconv.SDKProcessorLogQueueSize{}.Unit(),
		Data: metricdata.Sum[int64]{
			Temporality: metricdata.CumulativeTemporality,
			IsMonotonic: false,
			DataPoints:  []metricdata.DataPoint[int64]{{Attributes: blpSet(), Value: v}},
		},
	}
}

func blpProcessed(dPts ...metricdata.DataPoint[int64]) metricdata.Metrics {
	return metricdata.Metrics{
		Name:        otelconv.SDKProcessorLogProcessed{}.Name(),
		Description: otelconv.SDKProcessorLogProcessed{}.Description(),
		Unit:        otelconv.SDKProcessorLogProcessed{}.Unit(),
		Data: metricdata.Sum[int64]{
			Temporality: metricdata.CumulativeTemporality,
			IsMonotonic: true,
			DataPoints:  dPts,
		},
	}
}

func assertBLPMetrics(
	t *testing.T,
	reader sdkmetric.Reader,
	wantMetrics ...metricdata.Metrics,
) {
	t.Helper()

	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &rm))

	var found bool
	var gotScope metricdata.ScopeMetrics
	for _, sm := range rm.ScopeMetrics {
		if sm.Scope.Name == observ.ScopeName {
			gotScope = sm
			found = true
			break
		}
	}
	require.True(t, found, "observ scope %q not found in collected metrics", observ.ScopeName)

	metricdatatest.AssertEqual(
		t,
		metricdata.ScopeMetrics{
			Scope: instrumentation.Scope{
				Name:      observ.ScopeName,
				Version:   sdk.Version(),
				SchemaURL: observ.SchemaURL,
			},
			Metrics: wantMetrics,
		},
		gotScope,
		metricdatatest.IgnoreTimestamp(),
		metricdatatest.IgnoreExemplars(),
	)
}

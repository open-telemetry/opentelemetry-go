// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mapi "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// TestParallelCallbacksRunConcurrently verifies that, when the experimental
// feature is enabled, observable callbacks execute concurrently. Each callback
// blocks until every callback has started; sequential execution can never
// satisfy that barrier and the collection would time out.
func TestParallelCallbacksRunConcurrently(t *testing.T) {
	const n = 4
	// The pool is sized to GOMAXPROCS, which may be 1 on CI. Ensure enough
	// workers exist for all callbacks to overlap. Set before constructing the
	// provider because the pool is sized in newPipeline.
	prevMaxProcs := runtime.GOMAXPROCS(n)
	defer runtime.GOMAXPROCS(prevMaxProcs)

	t.Setenv("OTEL_GO_X_PARALLEL_CALLBACKS", "true")

	reader := metric.NewManualReader()
	mp := metric.NewMeterProvider(metric.WithReader(reader))
	t.Cleanup(func() { _ = mp.Shutdown(t.Context()) })
	m := mp.Meter("test")

	var started sync.WaitGroup
	started.Add(n)
	release := make(chan struct{})
	for i := range n {
		_, err := m.Int64ObservableCounter(
			fmt.Sprintf("ctr%d", i),
			mapi.WithInt64Callback(func(ctx context.Context, o mapi.Int64Observer) error {
				started.Done()
				// Wait for all callbacks to start. Once they have, the release
				// channel is closed, signaling that it is time to observe.
				select {
				case <-release:
				case <-ctx.Done():
					return ctx.Err()
				}
				o.Observe(1)
				return nil
			}),
		)
		require.NoError(t, err)
	}

	go func() {
		started.Wait()
		// All the callbacks have started, it's time to let them know that they can observe values.
		close(release)
	}()

	var rm metricdata.ResourceMetrics
	// An out-of-band watchdog, rather than a context deadline, fails fast with a
	// clear message even if Collect wedges for a reason unrelated to cancellation.
	done := make(chan error, 1)
	go func() { done <- reader.Collect(t.Context(), &rm) }()
	select {
	case err := <-done:
		require.NoError(t, err)
	case <-time.After(3 * time.Second):
		t.Fatal("callbacks did not run concurrently; collection timed out")
	}
	require.Len(t, rm.ScopeMetrics, 1)
	assert.Len(t, rm.ScopeMetrics[0].Metrics, n)
}

// TestParallelCallbacksRecordObservations ensures observations from both
// single-instrument and multi-instrument callbacks are recorded when running in
// parallel.
func TestParallelCallbacksRecordObservations(t *testing.T) {
	t.Setenv("OTEL_GO_X_PARALLEL_CALLBACKS", "true")

	reader := metric.NewManualReader()
	mp := metric.NewMeterProvider(metric.WithReader(reader))
	t.Cleanup(func() { _ = mp.Shutdown(t.Context()) })
	m := mp.Meter("test")

	_, err := m.Int64ObservableCounter("single",
		mapi.WithInt64Callback(func(_ context.Context, o mapi.Int64Observer) error {
			o.Observe(10)
			return nil
		}))
	require.NoError(t, err)

	gauge, err := m.Int64ObservableGauge("multi")
	require.NoError(t, err)
	_, err = m.RegisterCallback(func(_ context.Context, o mapi.Observer) error {
		o.ObserveInt64(gauge, 42)
		return nil
	}, gauge)
	require.NoError(t, err)

	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &rm))
	require.Len(t, rm.ScopeMetrics, 1)

	got := make(map[string]int64)
	for _, md := range rm.ScopeMetrics[0].Metrics {
		switch data := md.Data.(type) {
		case metricdata.Sum[int64]:
			got[md.Name] = data.DataPoints[0].Value
		case metricdata.Gauge[int64]:
			got[md.Name] = data.DataPoints[0].Value
		}
	}
	assert.Equal(t, int64(10), got["single"])
	assert.Equal(t, int64(42), got["multi"])
}

// TestCallbacksJoinErrors ensures errors from every callback are propagated,
// whether callbacks run sequentially (the default) or in parallel. Both the
// single-instrument and multi-instrument callback paths return an error so each
// join site is exercised in each mode.
func TestCallbacksJoinErrors(t *testing.T) {
	for _, parallel := range []bool{false, true} {
		name := "Sequential"
		value := "false"
		if parallel {
			name, value = "Parallel", "true"
		}
		t.Run(name, func(t *testing.T) {
			// Set explicitly in both modes so an ambient value does not decide
			// which path runs.
			t.Setenv("OTEL_GO_X_PARALLEL_CALLBACKS", value)

			reader := metric.NewManualReader()
			mp := metric.NewMeterProvider(metric.WithReader(reader))
			t.Cleanup(func() { _ = mp.Shutdown(t.Context()) })
			m := mp.Meter("test")

			errSingle := errors.New("single-instrument callback failed")
			errMulti := errors.New("multi-instrument callback failed")

			_, err := m.Int64ObservableCounter("single",
				mapi.WithInt64Callback(func(_ context.Context, o mapi.Int64Observer) error {
					o.Observe(1)
					return errSingle
				}))
			require.NoError(t, err)

			gauge, err := m.Int64ObservableGauge("multi")
			require.NoError(t, err)
			_, err = m.RegisterCallback(func(_ context.Context, o mapi.Observer) error {
				o.ObserveInt64(gauge, 1)
				return errMulti
			}, gauge)
			require.NoError(t, err)

			var rm metricdata.ResourceMetrics
			err = reader.Collect(t.Context(), &rm)
			assert.ErrorIs(t, err, errSingle)
			assert.ErrorIs(t, err, errMulti)
		})
	}
}

// TestParallelCallbacksErrorNotRetained verifies a callback error from one
// collection does not leak into the next.
func TestParallelCallbacksErrorNotRetained(t *testing.T) {
	t.Setenv("OTEL_GO_X_PARALLEL_CALLBACKS", "true")

	reader := metric.NewManualReader()
	mp := metric.NewMeterProvider(metric.WithReader(reader))
	t.Cleanup(func() { _ = mp.Shutdown(t.Context()) })
	m := mp.Meter("test")

	errFirst := errors.New("first collection fails")
	var calls atomic.Int64
	_, err := m.Int64ObservableCounter("ctr",
		mapi.WithInt64Callback(func(_ context.Context, o mapi.Int64Observer) error {
			o.Observe(1)
			if calls.Add(1) == 1 {
				return errFirst
			}
			return nil
		}))
	require.NoError(t, err)

	var rm metricdata.ResourceMetrics
	require.ErrorIs(t, reader.Collect(t.Context(), &rm), errFirst)
	// The second collection's callback succeeds.
	require.NoError(t, reader.Collect(t.Context(), &rm))
}

// TestParallelCallbacksMoreCallbacksThanWorkers exercises the job-queueing path
// where the number of callbacks exceeds the pool worker count, the common shape
// in single-CPU containers. All observations must still be recorded.
func TestParallelCallbacksMoreCallbacksThanWorkers(t *testing.T) {
	// One worker, many callbacks: every callback is queued through the single
	// worker rather than running on its own goroutine.
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(1))
	t.Setenv("OTEL_GO_X_PARALLEL_CALLBACKS", "true")

	reader := metric.NewManualReader()
	mp := metric.NewMeterProvider(metric.WithReader(reader))
	t.Cleanup(func() { _ = mp.Shutdown(t.Context()) })
	m := mp.Meter("test")

	const n = 8
	for i := range n {
		_, err := m.Int64ObservableCounter(
			fmt.Sprintf("ctr%d", i),
			mapi.WithInt64Callback(func(_ context.Context, o mapi.Int64Observer) error {
				o.Observe(1)
				return nil
			}),
		)
		require.NoError(t, err)
	}

	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &rm))
	require.Len(t, rm.ScopeMetrics, 1)
	assert.Len(t, rm.ScopeMetrics[0].Metrics, n)
}

// countPoolWorkers counts running callback-pool workers by symbol in a
// goroutine profile, which unrelated GC, timer, and test-runner goroutines
// cannot skew like a NumGoroutine() delta.
func countPoolWorkers(t *testing.T) int {
	t.Helper()
	var buf bytes.Buffer
	require.NoError(t, pprof.Lookup("goroutine").WriteTo(&buf, 2))
	return strings.Count(buf.String(), "(*callbackPool).worker")
}

// TestParallelCallbacksShutdownStopsWorkers verifies that enabling the feature
// starts callback-pool workers and that Shutdown tears them all down, counting
// workers by their symbol in a goroutine profile.
func TestParallelCallbacksShutdownStopsWorkers(t *testing.T) {
	t.Setenv("OTEL_GO_X_PARALLEL_CALLBACKS", "true")

	reader := metric.NewManualReader()
	mp := metric.NewMeterProvider(metric.WithReader(reader))

	m := mp.Meter("test")
	_, err := m.Int64ObservableCounter("ctr",
		mapi.WithInt64Callback(func(_ context.Context, o mapi.Int64Observer) error {
			o.Observe(1)
			return nil
		}))
	require.NoError(t, err)
	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(t.Context(), &rm))

	// The collection ran a callback, so at least one worker has started. The exact
	// count is not asserted; it depends on GOMAXPROCS and scheduling.
	require.NotZero(t, countPoolWorkers(t),
		"enabling the feature should start worker goroutines")

	require.NoError(t, mp.Shutdown(t.Context()))

	// Shutdown joins the workers, so all have returned. Retry briefly only to
	// absorb the lag between a worker returning and leaving the profile.
	var live int
	for range 100 {
		live = countPoolWorkers(t)
		if live == 0 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	assert.Zero(t, live, "Shutdown must stop all worker goroutines")
}

// TestParallelCallbacksShutdownIdempotent ensures shutting the provider down is idempotent.
func TestParallelCallbacksShutdownIdempotent(t *testing.T) {
	t.Setenv("OTEL_GO_X_PARALLEL_CALLBACKS", "true")

	mp := metric.NewMeterProvider(metric.WithReader(metric.NewManualReader()))
	require.NoError(t, mp.Shutdown(t.Context()))
	assert.NotPanics(t, func() { _ = mp.Shutdown(t.Context()) })
}

// TestParallelCallbacksReentrantShutdown verifies a callback may shut the
// provider down without deadlocking, matching the sequential default where
// Shutdown never waits on the in-flight collection. The callback uses a
// non-cancellable context (context.Background) on purpose: it has no deadline to
// fall back on, so only a Shutdown that refuses to wait on its own collection can
// return. Both the reentrant Shutdown and the collection must complete.
func TestParallelCallbacksReentrantShutdown(t *testing.T) {
	t.Setenv("OTEL_GO_X_PARALLEL_CALLBACKS", "true")

	reader := metric.NewManualReader()
	mp := metric.NewMeterProvider(metric.WithReader(reader))
	m := mp.Meter("test")

	shutdownDone := make(chan error, 1)
	_, err := m.Int64ObservableCounter("ctr",
		mapi.WithInt64Callback(func(_ context.Context, o mapi.Int64Observer) error {
			// A non-cancellable context is the point: it cannot self-heal via a
			// deadline, so it exercises the deadlock case.
			shutdownDone <- mp.Shutdown(context.Background()) //nolint:usetesting
			o.Observe(1)
			return nil
		}))
	require.NoError(t, err)

	collectDone := make(chan error, 1)
	go func() {
		var rm metricdata.ResourceMetrics
		collectDone <- reader.Collect(t.Context(), &rm)
	}()

	// A deadlock hangs both the reentrant Shutdown and the collection until the
	// package test timeout; the watchdog fails fast with a clear message instead.
	for range 2 {
		select {
		case err := <-shutdownDone:
			require.NoError(t, err)
		case err := <-collectDone:
			require.NoError(t, err)
		case <-time.After(3 * time.Second):
			t.Fatal("reentrant Shutdown from a callback deadlocked")
		}
	}
}

// TestShutdownDoesNotBlockOnInFlightCollectWithoutPool verifies that with the
// feature disabled (no callback pool), Shutdown does not wait on an in-flight
// collection. Pipeline teardown must be a no-op without a pool, preserving the
// default behavior where Shutdown never touched the collection lock.
func TestShutdownDoesNotBlockOnInFlightCollectWithoutPool(t *testing.T) {
	reader := metric.NewManualReader()
	mp := metric.NewMeterProvider(metric.WithReader(reader))
	m := mp.Meter("test")

	inCallback := make(chan struct{})
	release := make(chan struct{})
	_, err := m.Int64ObservableCounter("ctr",
		mapi.WithInt64Callback(func(_ context.Context, o mapi.Int64Observer) error {
			close(inCallback)
			<-release
			o.Observe(1)
			return nil
		}))
	require.NoError(t, err)

	collectDone := make(chan struct{})
	go func() {
		var rm metricdata.ResourceMetrics
		_ = reader.Collect(t.Context(), &rm)
		close(collectDone)
	}()

	// The collection is running its callback and holds the pipeline lock.
	<-inCallback

	shutdownDone := make(chan error, 1)
	go func() { shutdownDone <- mp.Shutdown(t.Context()) }()

	select {
	case err := <-shutdownDone:
		require.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("Shutdown blocked on an in-flight collection without a callback pool")
	}

	// Let the collection finish so no goroutine is leaked.
	close(release)
	<-collectDone
}

// BenchmarkRunCallbacks compares the cost of a collection's callback phase with
// the parallel-callbacks feature disabled (sequential) versus enabled, across a
// range of registered callback counts.
func BenchmarkRunCallbacks(b *testing.B) {
	// Callback counts spanning the three regimes relative to the GOMAXPROCS-sized
	// pool:
	// * 2 is under-subscribed (fewer callbacks than typical cores, where
	// parallelism has least to gain and dispatch overhead is most visible)
	// * 8 is near a common core count (the break-even zone)
	// * 32 is over-subscribed (enough independent work that spreading it across workers should pay off).
	for _, n := range []int{2, 8, 32} {
		b.Run(fmt.Sprintf("Sequential/%d", n), func(b *testing.B) {
			b.Setenv("OTEL_GO_X_PARALLEL_CALLBACKS", "false")
			benchCollectCallbacks(b, n)
		})
		b.Run(fmt.Sprintf("Parallel/%d", n), func(b *testing.B) {
			b.Setenv("OTEL_GO_X_PARALLEL_CALLBACKS", "true")
			benchCollectCallbacks(b, n)
		})
	}
}

// benchCollectCallbacks registers n observable counters, each with a callback
// that does a non-trivial local workload, then benchmarks repeated collection.
// Whether callbacks run sequentially or in parallel is controlled by the
// OTEL_GO_X_PARALLEL_CALLBACKS environment variable, which callers set before
// invoking this helper.
func benchCollectCallbacks(b *testing.B, n int) {
	reader := metric.NewManualReader()
	mp := metric.NewMeterProvider(metric.WithReader(reader))
	b.Cleanup(func() { _ = mp.Shutdown(b.Context()) })
	m := mp.Meter("bench")

	for i := range n {
		_, err := m.Int64ObservableCounter(
			fmt.Sprintf("ctr%d", i),
			mapi.WithInt64Callback(func(_ context.Context, o mapi.Int64Observer) error {
				// Simulate a non-trivial local observation workload so the
				// sequential and parallel paths are meaningfully comparable.
				// Observing acc uses the result, so the compiler cannot elide
				// the loop, and no sink shared across workers is needed.
				var acc int64
				for j := range 50_000 {
					acc += int64(j)
				}
				o.Observe(acc)
				return nil
			}),
		)
		require.NoError(b, err)
	}

	var rm metricdata.ResourceMetrics
	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		if err := reader.Collect(b.Context(), &rm); err != nil {
			b.Fatal(err)
		}
	}
}

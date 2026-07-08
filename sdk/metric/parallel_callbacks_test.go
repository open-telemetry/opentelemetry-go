// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric_test // import "go.opentelemetry.io/otel/sdk/metric"

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync"
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
			mapi.WithInt64Callback(func(_ context.Context, o mapi.Int64Observer) error {
				started.Done()
				// Wait for all callbacks to start. Once they have, the release
				// channel is closed, signaling that it is time to observe.
				select {
				case <-release:
				case <-time.After(2 * time.Second):
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

// TestParallelCallbacksJoinErrors ensures errors from callbacks executed in
// parallel are all propagated.
func TestParallelCallbacksJoinErrors(t *testing.T) {
	t.Setenv("OTEL_GO_X_PARALLEL_CALLBACKS", "true")

	reader := metric.NewManualReader()
	mp := metric.NewMeterProvider(metric.WithReader(reader))
	t.Cleanup(func() { _ = mp.Shutdown(t.Context()) })
	m := mp.Meter("test")

	errA := errors.New("callback A failed")
	errB := errors.New("callback B failed")

	gA, err := m.Int64ObservableGauge("gA")
	require.NoError(t, err)
	gB, err := m.Int64ObservableGauge("gB")
	require.NoError(t, err)

	// RegisterCallback is a no-op without at least one valid instrument, so each
	// error-returning callback is associated with a distinct observable.
	_, err = m.RegisterCallback(func(_ context.Context, o mapi.Observer) error {
		o.ObserveInt64(gA, 1)
		return errA
	}, gA)
	require.NoError(t, err)
	_, err = m.RegisterCallback(func(_ context.Context, o mapi.Observer) error {
		o.ObserveInt64(gB, 1)
		return errB
	}, gB)
	require.NoError(t, err)

	var rm metricdata.ResourceMetrics
	err = reader.Collect(t.Context(), &rm)
	assert.ErrorIs(t, err, errA)
	assert.ErrorIs(t, err, errB)
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

// TestParallelCallbacksShutdownStopsWorkers verifies, without reaching into
// unexported state, that enabling the feature starts worker goroutines and that
// Shutdown tears every one of them down. A leaked worker would keep the
// goroutine count above the pre-construction baseline.
func TestParallelCallbacksShutdownStopsWorkers(t *testing.T) {
	t.Setenv("OTEL_GO_X_PARALLEL_CALLBACKS", "true")

	// ManualReader is passive, so the only goroutines the provider starts are
	// the callback-pool workers. Capture the baseline before construction.
	baseline := runtime.NumGoroutine()

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

	// Workers are live and parked between collections, so the count is above the
	// baseline. This proves the pool was actually started.
	assert.Greater(t, runtime.NumGoroutine(), baseline,
		"enabling the feature should start worker goroutines")

	require.NoError(t, mp.Shutdown(t.Context()))

	// After shutdown every worker must have exited, returning to the baseline.
	// Poll in this goroutine (rather than via a helper that would spawn its own
	// goroutine and skew the count) to absorb the brief window between a worker
	// signaling done and the runtime reaping it.
	var got int
	for range 100 {
		got = runtime.NumGoroutine()
		if got <= baseline {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	assert.LessOrEqual(t, got, baseline, "Shutdown must stop all worker goroutines")
}

// TestParallelCallbacksShutdownIdempotent ensures shutting the provider down is idempotent.
func TestParallelCallbacksShutdownIdempotent(t *testing.T) {
	t.Setenv("OTEL_GO_X_PARALLEL_CALLBACKS", "true")

	mp := metric.NewMeterProvider(metric.WithReader(metric.NewManualReader()))
	require.NoError(t, mp.Shutdown(t.Context()))
	assert.NotPanics(t, func() { _ = mp.Shutdown(t.Context()) })
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

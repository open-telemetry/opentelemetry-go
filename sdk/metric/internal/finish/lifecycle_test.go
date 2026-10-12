// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package finish

import (
	"context"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var y2k = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)

func y2kPlus(seconds int) time.Time {
	return y2k.Add(time.Duration(seconds) * time.Second)
}

func TestLifecycle(t *testing.T) {
	t.Run("FinishIsIdempotent", func(t *testing.T) {
		var lifecycle Lifecycle
		measurement, ok := lifecycle.AcquireMeasurement()
		require.True(t, ok)
		measurement.Release()
		assert.True(t, lifecycle.Finish(y2kPlus(1)))
		assert.False(t, lifecycle.Finish(y2kPlus(2)))

		collection, ok := lifecycle.BeginCumulativeCollection(y2kPlus(3))
		require.True(t, ok)
		assert.Equal(t, y2kPlus(1), collection.Time())
		assert.True(t, collection.ShouldRetire())
		collection.Complete()
		_, ok = lifecycle.AcquireMeasurement()
		assert.False(t, ok)
	})

	t.Run("SharedCannotFinish", func(t *testing.T) {
		var lifecycle Lifecycle
		measurement, ok := lifecycle.AcquireSharedMeasurement()
		require.True(t, ok)
		measurement.Release()

		assert.False(t, lifecycle.Finish(y2kPlus(1)))
		collection, ok := lifecycle.BeginCumulativeCollection(y2kPlus(2))
		require.True(t, ok)
		assert.Equal(t, y2kPlus(2), collection.Time())
		assert.False(t, collection.ShouldRetire())
		collection.Complete()
		assert.False(t, lifecycle.Finish(y2kPlus(3)))
	})

	t.Run("SharedCancelsPendingFinish", func(t *testing.T) {
		var lifecycle Lifecycle
		require.True(t, lifecycle.Finish(y2kPlus(1)))

		measurement, ok := lifecycle.AcquireSharedMeasurement()
		require.True(t, ok)
		measurement.Release()
		collection, ok := lifecycle.BeginCumulativeCollection(y2kPlus(2))
		require.True(t, ok)
		assert.Equal(t, y2kPlus(2), collection.Time())
		assert.False(t, collection.ShouldRetire())
		collection.Complete()
	})

	t.Run("Reactivate", func(t *testing.T) {
		var lifecycle Lifecycle
		require.True(t, lifecycle.Finish(y2kPlus(1)))

		measurement, ok := lifecycle.AcquireMeasurement()
		require.True(t, ok)
		measurement.Release()
		collection, ok := lifecycle.BeginCumulativeCollection(y2kPlus(2))
		require.True(t, ok)
		assert.Equal(t, y2kPlus(2), collection.Time())
		assert.False(t, collection.ShouldRetire())
		collection.Complete()
		measurement, ok = lifecycle.AcquireMeasurement()
		require.True(t, ok)
		measurement.Release()
	})

	t.Run("DeltaCollectionRetiresActive", func(t *testing.T) {
		var lifecycle Lifecycle
		collection, ok := lifecycle.BeginDeltaCollection(y2kPlus(1))
		require.True(t, ok)
		assert.Equal(t, y2kPlus(1), collection.Time())
		assert.True(t, collection.ShouldRetire())
		collection.Complete()
		_, ok = lifecycle.AcquireMeasurement()
		assert.False(t, ok)
	})

	t.Run("RetiredIsTerminal", func(t *testing.T) {
		var lifecycle Lifecycle
		require.NoError(t, lifecycle.Retire(t.Context()))
		require.NoError(t, lifecycle.Retire(t.Context()))

		assert.False(t, lifecycle.Finish(y2kPlus(1)))
		_, ok := lifecycle.AcquireMeasurement()
		assert.False(t, ok)
		_, ok = lifecycle.AcquireSharedMeasurement()
		assert.False(t, ok)
		collection, ok := lifecycle.BeginCumulativeCollection(y2kPlus(2))
		assert.False(t, ok)
		assert.Zero(t, collection)
		collection, ok = lifecycle.BeginDeltaCollection(y2kPlus(3))
		assert.False(t, ok)
		assert.Zero(t, collection)
	})
}

func TestLifecycleSharedMeasurementPrecedesBlockedFinish(t *testing.T) {
	previous := runtime.GOMAXPROCS(1)
	defer runtime.GOMAXPROCS(previous)

	var lifecycle Lifecycle
	lifecycle.control.Lock()
	locked := true
	defer func() {
		if locked {
			lifecycle.control.Unlock()
		}
	}()

	started := make(chan struct{})
	finished := make(chan bool, 1)
	go func() {
		close(started)
		finished <- lifecycle.Finish(y2k)
	}()
	<-started
	// With one P, Finish runs until it blocks on control.
	runtime.Gosched()
	select {
	case <-finished:
		t.Fatal("Finish did not block")
	default:
	}

	measurement, ok := lifecycle.AcquireSharedMeasurement()
	require.True(t, ok)
	measurement.Release()
	lifecycle.control.Unlock()
	locked = false
	assert.False(t, <-finished)

	collection, ok := lifecycle.BeginCumulativeCollection(y2kPlus(1))
	require.True(t, ok)
	assert.False(t, collection.ShouldRetire())
	collection.Complete()
}

func TestLifecycleMeasurementWaitsForCollection(t *testing.T) {
	tests := []struct {
		name    string
		acquire func(*Lifecycle) (Measurement, bool)
	}{
		{name: "Ordinary", acquire: (*Lifecycle).AcquireMeasurement},
		{name: "Shared", acquire: (*Lifecycle).AcquireSharedMeasurement},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var lifecycle Lifecycle
			collection, ok := lifecycle.BeginCumulativeCollection(y2k)
			require.True(t, ok)

			acquired, started := make(chan Measurement, 1), make(chan struct{})
			go func() {
				close(started)
				measurement, ok := test.acquire(&lifecycle)
				if ok {
					acquired <- measurement
				}
			}()
			<-started
			select {
			case <-acquired:
				t.Fatal("measurement acquired during collection")
			case <-time.After(10 * time.Millisecond):
			}

			collection.Complete()
			measurement := <-acquired
			measurement.Release()
		})
	}
}

func TestLifecycleConcurrentSafeReactivation(t *testing.T) {
	const measurements = 64

	for range 10 {
		var lifecycle Lifecycle
		require.True(t, lifecycle.Finish(y2k))

		start := make(chan struct{})
		var wg sync.WaitGroup
		for range measurements {
			wg.Go(func() {
				<-start
				measurement, ok := lifecycle.AcquireMeasurement()
				if !ok {
					t.Error("measurement was not admitted")
					return
				}
				measurement.Release()
			})
		}
		close(start)
		wg.Wait()

		collection, ok := lifecycle.BeginCumulativeCollection(y2kPlus(1))
		require.True(t, ok)
		assert.False(t, collection.ShouldRetire())
		collection.Complete()
	}
}

func TestLifecycleConcurrentSafeDeltaCollection(t *testing.T) {
	for range 1_000 {
		var lifecycle Lifecycle
		start := make(chan struct{})
		var wg sync.WaitGroup

		wg.Go(func() {
			<-start
			measurement, ok := lifecycle.AcquireMeasurement()
			if ok {
				measurement.Release()
			}
		})
		wg.Go(func() {
			<-start
			collection, ok := lifecycle.BeginDeltaCollection(y2k)
			if ok {
				collection.Complete()
			}
		})

		close(start)
		wg.Wait()
		assert.Equal(t, uint32(lifecycleRetired), lifecycle.state.Load())
		assert.Zero(t, lifecycle.writers.Load())
		_, ok := lifecycle.AcquireMeasurement()
		assert.False(t, ok)
	}
}

func TestLifecycleCollectionWaitsForMeasurement(t *testing.T) {
	var lifecycle Lifecycle
	measurement, ok := lifecycle.AcquireMeasurement()
	require.True(t, ok)

	type collectionResult struct {
		collection Collection
		ok         bool
	}
	collected := make(chan collectionResult, 1)
	go func() {
		collection, ok := lifecycle.BeginCumulativeCollection(y2k)
		collected <- collectionResult{collection: collection, ok: ok}
	}()
	for lifecycleState(lifecycle.state.Load()) != lifecycleCollecting {
		runtime.Gosched()
	}
	select {
	case <-collected:
		t.Fatal("collection completed with a measurement in flight")
	default:
	}

	measurement.Release()
	result := <-collected
	require.True(t, result.ok)
	result.collection.Complete()
}

func TestLifecycleSerializesCollections(t *testing.T) {
	var lifecycle Lifecycle
	first, ok := lifecycle.BeginCumulativeCollection(y2k)
	require.True(t, ok)
	assert.False(t, lifecycle.control.TryLock())
	first.Complete()
	require.True(t, lifecycle.control.TryLock())
	lifecycle.control.Unlock()

	second, ok := lifecycle.BeginCumulativeCollection(y2kPlus(1))
	require.True(t, ok)
	second.Complete()
}

func TestLifecycleRetireWaitsForMeasurement(t *testing.T) {
	var lifecycle Lifecycle
	measurement, ok := lifecycle.AcquireMeasurement()
	require.True(t, ok)

	retired := make(chan error, 1)
	go func() {
		retired <- lifecycle.Retire(t.Context())
	}()
	for lifecycleState(lifecycle.state.Load()) != lifecycleCollecting {
		runtime.Gosched()
	}
	select {
	case <-retired:
		t.Fatal("retirement completed with a measurement in flight")
	default:
	}

	measurement.Release()
	require.NoError(t, <-retired)
	_, ok = lifecycle.AcquireMeasurement()
	assert.False(t, ok)
}

func TestLifecycleRetireCancellation(t *testing.T) {
	var lifecycle Lifecycle
	measurement, ok := lifecycle.AcquireMeasurement()
	require.True(t, ok)
	ctx, cancel := context.WithCancel(t.Context())
	retired := make(chan error, 1)
	go func() { retired <- lifecycle.Retire(ctx) }()
	for lifecycleState(lifecycle.state.Load()) != lifecycleCollecting {
		runtime.Gosched()
	}
	cancel()

	assert.ErrorIs(t, <-retired, context.Canceled)
	measurement.Release()
	require.NoError(t, lifecycle.Retire(t.Context()))
	_, ok = lifecycle.AcquireMeasurement()
	assert.False(t, ok)
}

func TestLifecycleParksUntilMeasurementsComplete(t *testing.T) {
	for _, test := range []struct {
		name    string
		collect func(*testing.T, *Lifecycle)
	}{
		{"Cumulative", func(t *testing.T, l *Lifecycle) {
			collection, ok := l.BeginCumulativeCollection(y2k)
			assert.True(t, ok)
			if ok {
				collection.Complete()
			}
		}},
		{"Delta", func(t *testing.T, l *Lifecycle) {
			collection, ok := l.BeginDeltaCollection(y2k)
			assert.True(t, ok)
			if ok {
				collection.Complete()
			}
		}},
		{"Retire", func(t *testing.T, l *Lifecycle) {
			assert.NoError(t, l.Retire(t.Context()))
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var lifecycle Lifecycle
				first, ok := lifecycle.AcquireMeasurement()
				require.True(t, ok)
				last, ok := lifecycle.AcquireMeasurement()
				require.True(t, ok)
				done := make(chan struct{})
				go func() {
					test.collect(t, &lifecycle)
					close(done)
				}()

				// Wait returns only once the collector is durably blocked.
				synctest.Wait()
				assert.Equal(t, uint32(lifecycleCollecting), lifecycle.state.Load())
				first.Release()
				synctest.Wait()
				select {
				case <-done:
					t.Fatal("collection completed before the last measurement")
				default:
				}
				last.Release()
				synctest.Wait()
				<-done
			})
		})
	}
}

func TestLifecycleRetireParksUntilControlCancellation(t *testing.T) {
	for _, cancelWhileWaiting := range []bool{false, true} {
		t.Run(strconv.FormatBool(cancelWhileWaiting), func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				var lifecycle Lifecycle
				lifecycle.control.Lock()
				ctx, cancel := context.WithCancel(t.Context())
				defer cancel()
				retired := make(chan error, 1)
				go func() { retired <- lifecycle.Retire(ctx) }()
				synctest.Wait()
				if cancelWhileWaiting {
					cancel()
					synctest.Wait()
					assert.ErrorIs(t, <-retired, context.Canceled)
					assert.Equal(t, uint32(lifecycleActive), lifecycle.state.Load())
					lifecycle.control.Unlock()
					require.NoError(t, lifecycle.Retire(t.Context()))
				} else {
					lifecycle.control.Unlock()
					time.Sleep(time.Millisecond)
					synctest.Wait()
					require.NoError(t, <-retired)
				}
			})
		})
	}
}

func TestLifecycleRetireParksUntilMeasurementCancellation(t *testing.T) {
	for _, test := range []struct {
		name  string
		state lifecycleState
	}{
		{"Active", lifecycleActive},
		{"Shared", lifecycleShared},
		{"FinishPending", lifecycleFinishPending},
	} {
		t.Run(test.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				state := test.state
				var lifecycle Lifecycle
				measurement, ok := lifecycle.AcquireMeasurement()
				require.True(t, ok)
				switch state {
				case lifecycleShared:
					shared, ok := lifecycle.AcquireSharedMeasurement()
					require.True(t, ok)
					shared.Release()
				case lifecycleFinishPending:
					require.True(t, lifecycle.Finish(y2k))
				}
				ctx, cancel := context.WithCancel(t.Context())
				retired := make(chan error, 1)
				go func() { retired <- lifecycle.Retire(ctx) }()
				synctest.Wait()
				cancel()
				synctest.Wait()
				assert.ErrorIs(t, <-retired, context.Canceled)
				assert.Equal(t, uint32(state), lifecycle.state.Load())
				if state == lifecycleFinishPending {
					assert.Equal(t, y2k, lifecycle.finished)
				}

				measurement.Release()
				// A last release can notify after a canceled waiter exits.
				lifecycle.writersDone <- struct{}{}
				// A later waiter must still wait for a newly admitted writer.
				measurement, ok = lifecycle.AcquireMeasurement()
				require.True(t, ok)
				go func() { retired <- lifecycle.Retire(t.Context()) }()
				synctest.Wait()
				assert.Equal(t, uint32(lifecycleCollecting), lifecycle.state.Load())
				measurement.Release()
				synctest.Wait()
				require.NoError(t, <-retired)
				assert.Equal(t, uint32(lifecycleRetired), lifecycle.state.Load())
			})
		})
	}
}

func TestLifecycleConcurrentSafeCollection(t *testing.T) {
	var lifecycle Lifecycle
	for range 1_000 {
		measurement, ok := lifecycle.AcquireMeasurement()
		require.True(t, ok)
		start := make(chan struct{})
		var wg sync.WaitGroup
		wg.Go(func() {
			<-start
			measurement.Release()
		})
		wg.Go(func() {
			<-start
			collection, ok := lifecycle.BeginCumulativeCollection(y2k)
			assert.True(t, ok)
			if ok {
				collection.Complete()
			}
		})
		close(start)
		wg.Wait()
		assert.Equal(t, uint32(lifecycleActive), lifecycle.state.Load())
		assert.Zero(t, lifecycle.writers.Load())
	}
}

func BenchmarkLifecycleMeasurement(b *testing.B) {
	var lifecycle Lifecycle
	b.ReportAllocs()
	for b.Loop() {
		measurement, _ := lifecycle.AcquireMeasurement()
		measurement.Release()
	}
}

func BenchmarkLifecycleCollectionWithWriter(b *testing.B) {
	var lifecycle Lifecycle
	measurements := make(chan Measurement)
	var wg sync.WaitGroup
	wg.Go(func() {
		for measurement := range measurements {
			measurement.Release()
		}
	})
	b.ReportAllocs()
	for b.Loop() {
		measurement, _ := lifecycle.AcquireMeasurement()
		measurements <- measurement
		collection, _ := lifecycle.BeginCumulativeCollection(y2k)
		collection.Complete()
	}
	close(measurements)
	wg.Wait()
}

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package finish

import (
	"runtime"
	"sync"
	"testing"
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
		lifecycle.Retire()
		lifecycle.Retire()

		assert.False(t, lifecycle.Finish(y2kPlus(1)))
		_, ok := lifecycle.AcquireMeasurement()
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
	var lifecycle Lifecycle
	collection, ok := lifecycle.BeginCumulativeCollection(y2k)
	require.True(t, ok)

	acquired, started := make(chan Measurement, 1), make(chan struct{})
	go func() {
		close(started)
		measurement, ok := lifecycle.AcquireMeasurement()
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
}

func TestLifecycleConcurrentReactivation(t *testing.T) {
	const measurements = 64

	for range 10 {
		var lifecycle Lifecycle
		require.True(t, lifecycle.Finish(y2k))

		start := make(chan struct{})
		var wg sync.WaitGroup
		wg.Add(measurements)
		for range measurements {
			go func() {
				defer wg.Done()
				<-start
				measurement, ok := lifecycle.AcquireMeasurement()
				if !ok {
					t.Error("measurement was not admitted")
					return
				}
				measurement.Release()
			}()
		}
		close(start)
		wg.Wait()

		collection, ok := lifecycle.BeginCumulativeCollection(y2kPlus(1))
		require.True(t, ok)
		assert.False(t, collection.ShouldRetire())
		collection.Complete()
	}
}

func TestLifecycleConcurrentDeltaCollection(t *testing.T) {
	for range 1_000 {
		var lifecycle Lifecycle
		start := make(chan struct{})
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			<-start
			measurement, ok := lifecycle.AcquireMeasurement()
			if ok {
				measurement.Release()
			}
		}()
		go func() {
			defer wg.Done()
			<-start
			collection, ok := lifecycle.BeginDeltaCollection(y2k)
			if ok {
				collection.Complete()
			}
		}()

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

	retired := make(chan struct{})
	go func() {
		lifecycle.Retire()
		close(retired)
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
	<-retired
	_, ok = lifecycle.AcquireMeasurement()
	assert.False(t, ok)
}

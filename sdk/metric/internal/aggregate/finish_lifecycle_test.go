// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFinishLifecycle(t *testing.T) {
	t.Run("FinishIsIdempotent", func(t *testing.T) {
		var lifecycle finishLifecycle
		require.True(t, lifecycle.startMeasurement())
		lifecycle.finishMeasurement()
		assert.True(t, lifecycle.finish(y2kPlus(1)))
		assert.False(t, lifecycle.finish(y2kPlus(2)))

		collection := lifecycle.startCollection(y2kPlus(3), finishCollectionKeepActive)
		assert.Equal(t, y2kPlus(1), collection.time)
		assert.True(t, collection.emit)
		assert.True(t, collection.retire)
		lifecycle.completeCollection(collection)
		assert.False(t, lifecycle.startMeasurement())
	})

	t.Run("Reactivate", func(t *testing.T) {
		var lifecycle finishLifecycle
		require.True(t, lifecycle.finish(y2kPlus(1)))

		require.True(t, lifecycle.startMeasurement())
		lifecycle.finishMeasurement()
		collection := lifecycle.startCollection(y2kPlus(2), finishCollectionKeepActive)
		assert.Equal(t, y2kPlus(2), collection.time)
		assert.True(t, collection.emit)
		assert.False(t, collection.retire)
		lifecycle.completeCollection(collection)
		require.True(t, lifecycle.startMeasurement())
		lifecycle.finishMeasurement()
	})

	t.Run("DeltaCollectionRetiresActive", func(t *testing.T) {
		var lifecycle finishLifecycle
		collection := lifecycle.startCollection(y2kPlus(1), finishCollectionRetireActive)
		assert.Equal(t, y2kPlus(1), collection.time)
		assert.True(t, collection.emit)
		assert.True(t, collection.retire)
		lifecycle.completeCollection(collection)
		assert.False(t, lifecycle.startMeasurement())
	})

	t.Run("RetiredIsTerminal", func(t *testing.T) {
		var lifecycle finishLifecycle
		lifecycle.retire()
		lifecycle.retire()

		assert.False(t, lifecycle.finish(y2kPlus(1)))
		assert.False(t, lifecycle.startMeasurement())
		collection := lifecycle.startCollection(y2kPlus(2), finishCollectionKeepActive)
		assert.False(t, collection.emit)
		assert.False(t, collection.retire)
	})
}

func TestFinishLifecycleCollectionWaitsForMeasurement(t *testing.T) {
	var lifecycle finishLifecycle
	require.True(t, lifecycle.startMeasurement())

	collected := make(chan finishCollection, 1)
	go func() {
		collected <- lifecycle.startCollection(y2k, finishCollectionKeepActive)
	}()
	for finishStateOf(lifecycle.state.Load()) != finishCollecting {
		runtime.Gosched()
	}
	select {
	case <-collected:
		t.Fatal("collection completed with a measurement in flight")
	default:
	}

	lifecycle.finishMeasurement()
	collection := <-collected
	assert.True(t, collection.emit)
	lifecycle.completeCollection(collection)
}

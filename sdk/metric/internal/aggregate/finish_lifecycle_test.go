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
		require.True(t, lifecycle.acquireMeasurement())
		lifecycle.releaseMeasurement()
		assert.True(t, lifecycle.finish(y2kPlus(1)))
		assert.False(t, lifecycle.finish(y2kPlus(2)))

		decision := lifecycle.startCollection(y2kPlus(3), collectionKeepActive)
		assert.Equal(t, y2kPlus(1), decision.time)
		assert.True(t, decision.emit)
		assert.True(t, decision.retire)
		lifecycle.completeCollection(decision)
		assert.False(t, lifecycle.acquireMeasurement())
	})

	t.Run("Reactivate", func(t *testing.T) {
		var lifecycle finishLifecycle
		require.True(t, lifecycle.finish(y2kPlus(1)))

		require.True(t, lifecycle.acquireMeasurement())
		lifecycle.releaseMeasurement()
		decision := lifecycle.startCollection(y2kPlus(2), collectionKeepActive)
		assert.Equal(t, y2kPlus(2), decision.time)
		assert.True(t, decision.emit)
		assert.False(t, decision.retire)
		lifecycle.completeCollection(decision)
		require.True(t, lifecycle.acquireMeasurement())
		lifecycle.releaseMeasurement()
	})

	t.Run("DeltaCollectionRetiresActive", func(t *testing.T) {
		var lifecycle finishLifecycle
		decision := lifecycle.startCollection(y2kPlus(1), collectionRetireActive)
		assert.Equal(t, y2kPlus(1), decision.time)
		assert.True(t, decision.emit)
		assert.True(t, decision.retire)
		lifecycle.completeCollection(decision)
		assert.False(t, lifecycle.acquireMeasurement())
	})

	t.Run("RetiredIsTerminal", func(t *testing.T) {
		var lifecycle finishLifecycle
		lifecycle.retire()
		lifecycle.retire()

		assert.False(t, lifecycle.finish(y2kPlus(1)))
		assert.False(t, lifecycle.acquireMeasurement())
		decision := lifecycle.startCollection(y2kPlus(2), collectionKeepActive)
		assert.False(t, decision.emit)
		assert.False(t, decision.retire)
	})
}

func TestFinishLifecycleCollectionWaitsForMeasurement(t *testing.T) {
	var lifecycle finishLifecycle
	require.True(t, lifecycle.acquireMeasurement())

	collected := make(chan collectionDecision, 1)
	go func() {
		collected <- lifecycle.startCollection(y2k, collectionKeepActive)
	}()
	for finishStateOf(lifecycle.state.Load()) != lifecycleCollecting {
		runtime.Gosched()
	}
	select {
	case <-collected:
		t.Fatal("collection completed with a measurement in flight")
	default:
	}

	lifecycle.releaseMeasurement()
	decision := <-collected
	assert.True(t, decision.emit)
	lifecycle.completeCollection(decision)
}

func TestFinishLifecycleRetireWaitsForMeasurement(t *testing.T) {
	var lifecycle finishLifecycle
	require.True(t, lifecycle.acquireMeasurement())

	retired := make(chan struct{})
	go func() {
		lifecycle.retire()
		close(retired)
	}()
	for finishStateOf(lifecycle.state.Load()) != lifecycleCollecting {
		runtime.Gosched()
	}
	select {
	case <-retired:
		t.Fatal("retirement completed with a measurement in flight")
	default:
	}

	lifecycle.releaseMeasurement()
	<-retired
	assert.False(t, lifecycle.acquireMeasurement())
}

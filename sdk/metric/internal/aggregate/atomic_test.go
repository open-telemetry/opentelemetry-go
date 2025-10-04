// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate // import "go.opentelemetry.io/otel/sdk/metric/internal/aggregate"

import (
	"math"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAtomicSumAddFloatConcurrentSafe(t *testing.T) {
	var wg sync.WaitGroup
	var aSum atomicCounter[float64]
	for _, in := range []float64{
		0.2,
		0.25,
		1.6,
		10.55,
		42.4,
	} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			aSum.add(in)
		}()
	}
	wg.Wait()
	assert.Equal(t, float64(55), math.Round(aSum.load()))
}

func TestAtomicSumAddIntConcurrentSafe(t *testing.T) {
	var wg sync.WaitGroup
	var aSum atomicCounter[int64]
	for _, in := range []int64{
		1,
		2,
		3,
		4,
		5,
	} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			aSum.add(in)
		}()
	}
	wg.Wait()
	assert.Equal(t, int64(15), aSum.load())
}

func TestHotColdWaitGroupConcurrentSafe(t *testing.T) {
	var wg sync.WaitGroup
	hcwg := &hotColdWaitGroup{}
	var data [2]uint64
	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			hotIdx := hcwg.start()
			defer hcwg.done(hotIdx)
			atomic.AddUint64(&data[hotIdx], 1)
		}()
	}
	for range 2 {
		readIdx := hcwg.swapHotAndWait()
		assert.NotPanics(t, func() {
			// reading without using atomics should not panic since we are
			// reading from the cold element, and have waited for all writes to
			// finish.
			t.Logf("read value %+v", data[readIdx])
		})
	}
	wg.Wait()
}

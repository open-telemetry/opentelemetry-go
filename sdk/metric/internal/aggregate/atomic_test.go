// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate // import "go.opentelemetry.io/otel/sdk/metric/internal/aggregate"

import (
	"math"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAtomicSumAddFloatConcurrentSafe(t *testing.T) {
	var wg sync.WaitGroup
	var aSum atomicSum[float64]
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
	var aSum atomicSum[int64]
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

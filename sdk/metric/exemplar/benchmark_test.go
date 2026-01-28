// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package exemplar // import "go.opentelemetry.io/otel/sdk/metric/exemplar"

import (
	"runtime"
	"testing"
	"time"
)

func BenchmarkFixedSizeReservoirOffer(b *testing.B) {
	ts := time.Now()
	val := NewValue[int64](25)
	ctx := b.Context()
	reservoir := NewFixedSizeReservoir(runtime.NumCPU())
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			reservoir.Offer(ctx, ts, val, nil)
			// Periodically trigger a reset, because the algorithm for fixed-size
			// reservoirs records exemplars very infrequently after a large
			// number of collect calls.
			if i%100 == 99 {
				reservoir.reset()
			}
			i++
		}
	})
}

func BenchmarkHistogramReservoirOffer(b *testing.B) {
	ts := time.Now()
	ctx := b.Context()
	buckets := []float64{0, 5, 10, 25, 50, 75, 100, 250, 500, 750, 1000, 2500, 5000, 7500, 10000}
	values := make([]Value, len(buckets))
	for i, bucket := range buckets {
		values[i] = NewValue[float64](bucket + 1)
	}
	res := NewHistogramReservoir(buckets)
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			res.Offer(ctx, ts, values[i%len(values)], nil)
			i++
		}
	})
}

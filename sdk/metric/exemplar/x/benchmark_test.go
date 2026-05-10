// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"runtime"
	"testing"
	"time"

	"go.opentelemetry.io/otel/sdk/metric/exemplar"
)

func BenchmarkFixedSizeRoundRobinReservoirOffer(b *testing.B) {
	ts := time.Now()
	val := exemplar.NewValue[int64](25)
	ctx := b.Context()
	reservoir := NewFixedSizeRoundRobinReservoir(runtime.NumCPU())
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

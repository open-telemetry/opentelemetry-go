// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This test is too large for the race detector.  This SDK uses no locks
// that the race detector would help with, anyway.
// +build !race

package metric_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
)

func TestStressInt64Histogram(t *testing.T) {
	desc := metric.NewDescriptor("some_metric", metric.MeasureKind, core.Int64NumberKind)
	h := histogram.New(&desc, []core.Number{core.NewInt64Number(25), core.NewInt64Number(50), core.NewInt64Number(75)})

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	go func() {
		rnd := rand.New(rand.NewSource(time.Now().Unix()))
		for {
			select {
			case <-ctx.Done():
				return
			default:
				_ = h.Update(ctx, core.NewInt64Number(rnd.Int63()%100), &desc)
			}
		}
	}()

	startTime := time.Now()
	for time.Since(startTime) < time.Second {
		h.Checkpoint(context.Background(), &desc)

		b, _ := h.Histogram()
		c, _ := h.Count()

		var realCount int64
		for _, c := range b.Counts {
			v := c.AsInt64()
			realCount += v
		}

		if realCount != c {
			t.Fail()
		}
	}
}

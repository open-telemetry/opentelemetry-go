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

package metric_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/minmaxsumcount"
)

func TestStressInt64MinMaxSumCount(t *testing.T) {
	desc := metric.NewDescriptor("some_metric", metric.ValueRecorderKind, metric.Int64NumberKind)
	mmsc := minmaxsumcount.New(&desc)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		rnd := rand.New(rand.NewSource(time.Now().Unix()))
		v := rnd.Int63() % 103
		for {
			select {
			case <-ctx.Done():
				return
			default:
				_ = mmsc.Update(ctx, metric.NewInt64Number(v), &desc)
			}
			v++
		}
	}()

	startTime := time.Now()
	for time.Since(startTime) < time.Second {
		mmsc.Checkpoint(context.Background(), &desc)

		s, _ := mmsc.Sum()
		c, _ := mmsc.Count()
		min, e1 := mmsc.Min()
		max, e2 := mmsc.Max()
		if c == 0 && (e1 == nil || e2 == nil || s.AsInt64() != 0) {
			t.Fail()
		}
		if c != 0 {
			if e1 != nil || e2 != nil {
				t.Fail()
			}
			lo, hi, sum := min.AsInt64(), max.AsInt64(), s.AsInt64()

			if hi-lo+1 != c {
				t.Fail()
			}
			if c == 1 {
				if lo != hi || lo != sum {
					t.Fail()
				}
			} else {
				if hi*(hi+1)/2-(lo-1)*lo/2 != sum {
					t.Fail()
				}
			}
		}
	}
}

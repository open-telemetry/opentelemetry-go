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

	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
)

func TestStressInt64Histogram(t *testing.T) {
	desc := otel.NewDescriptor("some_metric", otel.ValueRecorderInstrumentKind, otel.Int64NumberKind)

	alloc := histogram.New(2, &desc, []float64{25, 50, 75})
	h, ckpt := &alloc[0], &alloc[1]

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	go func() {
		rnd := rand.New(rand.NewSource(time.Now().Unix()))
		for {
			select {
			case <-ctx.Done():
				return
			default:
				_ = h.Update(ctx, otel.NewInt64Number(rnd.Int63()%100), &desc)
			}
		}
	}()

	startTime := time.Now()
	for time.Since(startTime) < time.Second {
		require.NoError(t, h.SynchronizedMove(ckpt, &desc))

		b, _ := ckpt.Histogram()
		c, _ := ckpt.Count()

		var realCount int64
		for _, c := range b.Counts {
			v := int64(c)
			realCount += v
		}

		if realCount != c {
			t.Fail()
		}
	}
}

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

package exemplar

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/trace"
)

// Sat Jan 01 2000 00:00:00 GMT+0000.
var staticTime = time.Unix(946684800, 0)

type factory[N int64 | float64] func(requstedCap int) (r Reservoir[N], actualCap int)

func ReservoirTest[N int64 | float64](f factory[N]) func(*testing.T) {
	return func(t *testing.T) {
		t.Helper()

		ctx := context.Background()

		t.Run("CaptureSpanContext", func(t *testing.T) {
			t.Helper()

			r, n := f(1)
			if n < 1 {
				t.Skip("skipping, reservoir capacity less than 1:", n)
			}

			tID, sID := trace.TraceID{0x01}, trace.SpanID{0x01}
			sc := trace.NewSpanContext(trace.SpanContextConfig{
				TraceID:    tID,
				SpanID:     sID,
				TraceFlags: trace.FlagsSampled,
			})
			ctx := trace.ContextWithSpanContext(ctx, sc)

			r.Offer(ctx, staticTime, 10, nil)

			var dest []metricdata.Exemplar[N]
			r.Collect(&dest)

			want := metricdata.Exemplar[N]{
				Time:    staticTime,
				Value:   10,
				SpanID:  []byte(sID[:]),
				TraceID: []byte(tID[:]),
			}
			require.Len(t, dest, 1, "number of collected exemplars")
			assert.Equal(t, want, dest[0])
		})

		t.Run("FilterAttributes", func(t *testing.T) {
			t.Helper()

			r, n := f(1)
			if n < 1 {
				t.Skip("skipping, reservoir capacity less than 1:", n)
			}

			adminTrue := attribute.Bool("admin", true)
			r.Offer(ctx, staticTime, 10, []attribute.KeyValue{adminTrue})

			var dest []metricdata.Exemplar[N]
			r.Collect(&dest)

			want := metricdata.Exemplar[N]{
				FilteredAttributes: []attribute.KeyValue{adminTrue},
				Time:               staticTime,
				Value:              10,
			}
			require.Len(t, dest, 1, "number of collected exemplars")
			assert.Equal(t, want, dest[0])
		})

		t.Run("CollectLessThanN", func(t *testing.T) {
			t.Helper()

			r, n := f(2)
			if n < 2 {
				t.Skip("skipping, reservoir capacity less than 2:", n)
			}

			r.Offer(ctx, staticTime, 10, nil)

			var dest []metricdata.Exemplar[N]
			r.Collect(&dest)
			// No empty exemplars are exported.
			require.Len(t, dest, 1, "number of collected exemplars")
		})

		t.Run("MultipleOffers", func(t *testing.T) {
			t.Helper()

			r, n := f(3)
			if n < 1 {
				t.Skip("skipping, reservoir capacity less than 1:", n)
			}

			for i := 0; i < n+1; i++ {
				v := N(i)
				r.Offer(ctx, staticTime, v, nil)
			}

			var dest []metricdata.Exemplar[N]
			r.Collect(&dest)
			assert.Len(t, dest, n, "multiple offers did not fill reservoir")

			// Ensure the collect reset also resets any couting state.
			for i := 0; i < n+1; i++ {
				v := N(i)
				r.Offer(ctx, staticTime, v, nil)
			}

			dest = dest[:0]
			r.Collect(&dest)
			assert.Len(t, dest, n, "internal count state not reset")
		})

		t.Run("DropAll", func(t *testing.T) {
			t.Helper()

			r, n := f(0)
			if n > 0 {
				t.Skip("skipping, reservoir capacity greater than 0:", n)
			}

			r.Offer(context.Background(), staticTime, 10, nil)

			dest := []metricdata.Exemplar[N]{{}} // Should be reset to empty.
			r.Collect(&dest)
			assert.Len(t, dest, 0, "no exemplars should be collected")
		})
	}
}

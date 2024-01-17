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

package exemplar // import "go.opentelemetry.io/otel/sdk/metric/internal/exemplar"

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/trace"
)

func sample(parent context.Context) context.Context {
	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    trace.TraceID{0x01},
		SpanID:     trace.SpanID{0x01},
		TraceFlags: trace.FlagsSampled,
	})
	return trace.ContextWithSpanContext(parent, sc)
}

func TestAlwaysSample(t *testing.T) {
	t.Run("Int64", testAlwaysSample[int64])
	t.Run("Float64", testAlwaysSample[float64])
}

func testAlwaysSample[N int64 | float64](t *testing.T) {
	ctx := context.Background()

	assert.True(t, AlwaysSample[N](ctx))
	assert.True(t, AlwaysSample[N](sample(ctx)))
}

func TestNeverSample(t *testing.T) {
	t.Run("Int64", testNeverSample[int64])
	t.Run("Float64", testNeverSample[float64])
}

func testNeverSample[N int64 | float64](t *testing.T) {
	ctx := context.Background()

	assert.False(t, NeverSample[N](ctx))
	assert.False(t, NeverSample[N](sample(ctx)))
}

func TestTraceBasedSample(t *testing.T) {
	t.Run("Int64", testTraceBasedSample[int64])
	t.Run("Float64", testTraceBasedSample[float64])
}

func testTraceBasedSample[N int64 | float64](t *testing.T) {
	ctx := context.Background()

	assert.False(t, TraceBasedSample[N](ctx))
	assert.True(t, TraceBasedSample[N](sample(ctx)))
}

type res[N int64 | float64] struct {
	OfferCalled   bool
	CollectCalled bool
	FlushCalled   bool
}

func (r *res[N]) Offer(context.Context, time.Time, N, []attribute.KeyValue) {
	r.OfferCalled = true
}

func (r *res[N]) Collect(*[]metricdata.Exemplar[N]) {
	r.CollectCalled = true
}

func (r *res[N]) Flush(*[]metricdata.Exemplar[N]) {
	r.FlushCalled = true
}

func TestFilteredReservoir(t *testing.T) {
	t.Run("Int64", testFilteredReservoir[int64])
	t.Run("Float64", testFilteredReservoir[float64])
}

func testFilteredReservoir[N int64 | float64](t *testing.T) {
	under := &res[N]{}

	var called bool
	fltr := func(context.Context) bool {
		called = true
		return true
	}

	r := Filtered[N](under, fltr)

	r.Offer(context.Background(), staticTime, 0, nil)
	assert.True(t, called, "filter not called on Offer")
	assert.True(t, under.OfferCalled, "underlying Reservoir Offer not called")

	r.Collect(nil)
	assert.True(t, under.CollectCalled, "underlying Reservoir Collect not called")

	r.Flush(nil)
	assert.True(t, under.FlushCalled, "underlying Reservoir Flush not called")
}

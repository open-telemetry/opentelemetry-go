// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

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

func TestSampledFilter(t *testing.T) {
	t.Run("Int64", testSampledFiltered[int64])
	t.Run("Float64", testSampledFiltered[float64])
}

func testSampledFiltered[N int64 | float64](t *testing.T) {
	under := &res[N]{}

	r := SampledFilter[N](under)

	ctx := context.Background()
	r.Offer(ctx, staticTime, 0, nil)
	assert.False(t, under.OfferCalled, "underlying Reservoir Offer called")
	r.Offer(sample(ctx), staticTime, 0, nil)
	assert.True(t, under.OfferCalled, "underlying Reservoir Offer not called")

	r.Collect(nil)
	assert.True(t, under.CollectCalled, "underlying Reservoir Collect not called")
}

func sample(parent context.Context) context.Context {
	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    trace.TraceID{0x01},
		SpanID:     trace.SpanID{0x01},
		TraceFlags: trace.FlagsSampled,
	})
	return trace.ContextWithSpanContext(parent, sc)
}

type res[N int64 | float64] struct {
	OfferCalled   bool
	CollectCalled bool
}

func (r *res[N]) Offer(context.Context, time.Time, N, []attribute.KeyValue) {
	r.OfferCalled = true
}

func (r *res[N]) Collect(*[]metricdata.Exemplar[N]) {
	r.CollectCalled = true
}

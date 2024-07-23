// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package exemplar // import "go.opentelemetry.io/otel/sdk/metric/internal/exemplar"

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/trace"
)

func TestSampledFilter(t *testing.T) {
	t.Run("Int64", testSampledFiltered[int64])
	t.Run("Float64", testSampledFiltered[float64])
}

func testSampledFiltered[N int64 | float64](t *testing.T) {
	ctx := context.Background()

	assert.False(t, SampledFilter(ctx), "non-sampled context should not be offered")
	assert.True(t, SampledFilter(sample(ctx)), "sampled context should be offered")
}

func sample(parent context.Context) context.Context {
	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    trace.TraceID{0x01},
		SpanID:     trace.SpanID{0x01},
		TraceFlags: trace.FlagsSampled,
	})
	return trace.ContextWithSpanContext(parent, sc)
}

func TestAlwaysOnFilter(t *testing.T) {
	t.Run("Int64", testAlwaysOnFiltered[int64])
	t.Run("Float64", testAlwaysOnFiltered[float64])
}

func testAlwaysOnFiltered[N int64 | float64](t *testing.T) {
	ctx := context.Background()

	assert.True(t, AlwaysOnFilter(ctx), "non-sampled context should not be offered")
	assert.True(t, AlwaysOnFilter(sample(ctx)), "sampled context should be offered")
}

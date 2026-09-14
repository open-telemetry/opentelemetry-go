// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/exemplar"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

func TestCollectExemplars(t *testing.T) {
	t.Run("Int64", testCollectExemplars[int64]())
	t.Run("Float64", testCollectExemplars[float64]())
}

func testCollectExemplars[N int64 | float64]() func(t *testing.T) {
	return func(t *testing.T) {
		now := time.Now()
		alice := attribute.String("user", "Alice")
		value := N(1)
		spanID := [8]byte{0x1}
		traceID := [16]byte{0x1}

		out := new([]metricdata.Exemplar[N])
		collectExemplars(out, func(in *[]exemplar.Exemplar) {
			*in = reset(*in, 1, 1)
			(*in)[0] = exemplar.Exemplar{
				FilteredAttributes: []attribute.KeyValue{alice},
				Time:               now,
				Value:              exemplar.NewValue(value),
				SpanID:             spanID[:],
				TraceID:            traceID[:],
			}
		})

		assert.Equal(t, []metricdata.Exemplar[N]{{
			FilteredAttributes: []attribute.KeyValue{alice},
			Time:               now,
			Value:              value,
			SpanID:             spanID[:],
			TraceID:            traceID[:],
		}}, *out)
	}
}

func TestCollectExemplarsAfter(t *testing.T) {
	t.Run("Int64", testCollectExemplarsAfter[int64]())
	t.Run("Float64", testCollectExemplarsAfter[float64]())
}

func testCollectExemplarsAfter[N int64 | float64]() func(t *testing.T) {
	return func(t *testing.T) {
		t1 := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		t2 := time.Date(2000, 1, 1, 0, 0, 1, 0, time.UTC)
		t3 := time.Date(2000, 1, 1, 0, 0, 2, 0, time.UTC)

		alice := attribute.String("user", "Alice")
		bob := attribute.String("user", "Bob")
		spanID := [8]byte{0x1}
		traceID := [16]byte{0x1}

		out := new([]metricdata.Exemplar[N])
		collectExemplarsAfter(out, t2, func(in *[]exemplar.Exemplar) {
			*in = reset(*in, 2, 2)
			(*in)[0] = exemplar.Exemplar{
				FilteredAttributes: []attribute.KeyValue{alice},
				Time:               t1, // before cutoff t2
				Value:              exemplar.NewValue(N(1)),
				SpanID:             spanID[:],
				TraceID:            traceID[:],
			}
			(*in)[1] = exemplar.Exemplar{
				FilteredAttributes: []attribute.KeyValue{bob},
				Time:               t3, // after cutoff t2
				Value:              exemplar.NewValue(N(2)),
				SpanID:             spanID[:],
				TraceID:            traceID[:],
			}
		})

		assert.Equal(t, []metricdata.Exemplar[N]{{
			FilteredAttributes: []attribute.KeyValue{bob},
			Time:               t3,
			Value:              N(2),
			SpanID:             spanID[:],
			TraceID:            traceID[:],
		}}, *out)
	}
}

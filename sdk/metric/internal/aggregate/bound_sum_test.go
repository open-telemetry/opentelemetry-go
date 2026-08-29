// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/exemplar"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

type retainingReservoir struct {
	mu        sync.Mutex
	exemplars []exemplar.Exemplar
}

func (r *retainingReservoir) Offer(
	_ context.Context,
	value int64,
	lazy lazyFilteredAttributes,
) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.exemplars = append(r.exemplars, exemplar.Exemplar{
		FilteredAttributes: lazy.Dropped(),
		Time:               time.Unix(value, 0),
		Value:              exemplar.NewValue(value),
	})
}

func (r *retainingReservoir) Collect(dest *[]exemplar.Exemplar) {
	r.mu.Lock()
	defer r.mu.Unlock()
	*dest = append(*dest, r.exemplars...)
	r.exemplars = nil
}

func TestBoundSumPreservesFinishedLifetimes(t *testing.T) {
	for _, temporality := range []metricdata.Temporality{
		metricdata.CumulativeTemporality,
		metricdata.DeltaTemporality,
	} {
		t.Run(temporality.String(), func(t *testing.T) {
			agg := Builder[int64]{
				Temporality: temporality,
				ReservoirFunc: func(attribute.Set) FilteredExemplarReservoir[int64] {
					return new(retainingReservoir)
				},
			}.BoundSum(true)
			binder := agg.(interface {
				Bind(attribute.Set) func(context.Context, int64)
			})
			finisher := agg.(interface{ Finish(attribute.Set) })
			attrs := attribute.NewSet(attribute.String("route", "/orders"))

			binder.Bind(attrs)(t.Context(), 1)
			finisher.Finish(attrs)
			binder.Bind(attrs)(t.Context(), 2)
			finisher.Finish(attrs)

			var dest metricdata.Aggregation
			require.Equal(t, 1, agg.ComputeAggregation()(&dest))
			sum := dest.(metricdata.Sum[int64])
			require.Len(t, sum.DataPoints, 1)
			assert.Equal(t, int64(3), sum.DataPoints[0].Value)
			require.Len(t, sum.DataPoints[0].Exemplars, 2)
			assert.Equal(t, int64(1), sum.DataPoints[0].Exemplars[0].Value)
			assert.Equal(t, int64(2), sum.DataPoints[0].Exemplars[1].Value)
		})
	}
}

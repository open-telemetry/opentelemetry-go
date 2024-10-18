// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate // import "go.opentelemetry.io/otel/sdk/metric/internal/aggregate"

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/exemplar"
)

// filteredExemplarReservoir handles the pre-sampled exemplar of measurements made.
type filteredExemplarReservoir[N int64 | float64] struct {
	filter    exemplar.Filter
	reservoir exemplar.Reservoir
}

// newFilteredExemplarReservoir creates a [FilteredExemplarReservoir] which
// only offers values that are allowed by the filter. If the provided filter is
// nil, all measurements are dropped..
func newFilteredExemplarReservoir[N int64 | float64](f exemplar.Filter, r exemplar.Reservoir) *filteredExemplarReservoir[N] {
	return &filteredExemplarReservoir[N]{
		filter:    f,
		reservoir: r,
	}
}

func (f *filteredExemplarReservoir[N]) Offer(ctx context.Context, val N, attr []attribute.KeyValue) {
	if f.filter(ctx) {
		// only record the current time if we are sampling this measurement.
		f.reservoir.Offer(ctx, time.Now(), exemplar.NewValue(val), attr)
	}
}

func (f *filteredExemplarReservoir[N]) Collect(dest *[]exemplar.Exemplar) { f.reservoir.Collect(dest) }

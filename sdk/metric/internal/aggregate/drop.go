// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/exemplar"
)

// DropReservoir returns a [FilteredExemplarReservoir] that drops all measurements it is offered.
func DropReservoir[N int64 | float64](attribute.Set) FilteredExemplarReservoir[N] {
	return &dropRes[N]{}
}

type dropRes[N int64 | float64] struct{}

// Offer does nothing, all measurements offered will be dropped.
func (*dropRes[N]) Offer(context.Context, N, lazyFilteredAttributes) {}

// Collect resets dest. No exemplars will ever be returned.
func (*dropRes[N]) Collect(dest *[]exemplar.Exemplar) {
	if dest != nil {
		clear(*dest) // Erase elements to let GC collect objects
		*dest = (*dest)[:0]
	}
}

// Merge does nothing, all measurements are dropped.
func (*dropRes[N]) Merge(FilteredExemplarReservoir[N]) {}

// Reset does nothing.
func (*dropRes[N]) Reset() {}

// IsMergeable returns false for dropRes.
func (*dropRes[N]) IsMergeable() bool {
	return false
}

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package aggregate

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/otel/sdk/metric/exemplar"
	"go.opentelemetry.io/otel/sdk/metric/internal/reservoir"
)

// FilteredExemplarReservoir wraps a [exemplar.Reservoir] with a filter.
type FilteredExemplarReservoir[N int64 | float64] interface {
	// Offer accepts the parameters associated with a measurement. The
	// parameters will be stored as an exemplar if the filter decides to
	// sample the measurement.
	//
	// The passed ctx needs to contain any baggage or span that were active
	// when the measurement was made. This information may be used by the
	// Reservoir in making a sampling decision.
	//
	// The lazy parameter provides filtered attribute details when sampled,
	// allowing dropped attributes to be computed only if the measurement is sampled.
	Offer(ctx context.Context, val N, lazy lazyFilteredAttributes)
	// Collect returns all the held exemplars in the reservoir.
	Collect(dest *[]exemplar.Exemplar)
	// Merge merges the sampled exemplars from other into this reservoir.
	Merge(other FilteredExemplarReservoir[N])
	// Reset resets the reservoir's stored exemplars and sampling state.
	Reset()
	// IsMergeable returns whether the underlying reservoir supports merging.
	IsMergeable() bool
}

// filteredExemplarReservoir handles the pre-sampled exemplar of measurements made.
type filteredExemplarReservoir[N int64 | float64] struct {
	filter    exemplar.Filter
	reservoir exemplar.Reservoir
	// The exemplar.Reservoir is not required to be concurrent safe, but
	// implementations can indicate that they are concurrent-safe by embedding
	// reservoir.ConcurrentSafe in order to improve performance.
	reservoirMux   sync.Mutex
	concurrentSafe bool
}

// NewFilteredExemplarReservoir creates a [FilteredExemplarReservoir] which only offers values
// that are allowed by the filter.
func NewFilteredExemplarReservoir[N int64 | float64](
	f exemplar.Filter,
	r exemplar.Reservoir,
) FilteredExemplarReservoir[N] {
	_, concurrentSafe := r.(reservoir.ConcurrentSafe)
	return &filteredExemplarReservoir[N]{
		filter:         f,
		reservoir:      r,
		concurrentSafe: concurrentSafe,
	}
}

func (f *filteredExemplarReservoir[N]) Offer(ctx context.Context, val N, lazy lazyFilteredAttributes) {
	if f.filter(ctx) {
		// only record the current time if we are sampling this measurement.
		ts := time.Now()
		attr := lazy.Dropped()
		if !f.concurrentSafe {
			f.reservoirMux.Lock()
			defer f.reservoirMux.Unlock()
		}
		f.reservoir.Offer(ctx, ts, exemplar.NewValue(val), attr)
	}
}

func (f *filteredExemplarReservoir[N]) Collect(dest *[]exemplar.Exemplar) {
	if !f.concurrentSafe {
		f.reservoirMux.Lock()
		defer f.reservoirMux.Unlock()
	}
	f.reservoir.Collect(dest)
}

type mergeableReservoir interface {
	Merge(other exemplar.Reservoir)
	Reset()
}

func (f *filteredExemplarReservoir[N]) Merge(other FilteredExemplarReservoir[N]) {
	if f == nil || other == nil || f == other {
		return
	}
	o, ok := other.(*filteredExemplarReservoir[N])
	if !ok || o == nil {
		return
	}
	if mr, ok := f.reservoir.(mergeableReservoir); ok {
		if !f.concurrentSafe {
			f.reservoirMux.Lock()
			defer f.reservoirMux.Unlock()
		}
		if !o.concurrentSafe {
			o.reservoirMux.Lock()
			defer o.reservoirMux.Unlock()
		}
		mr.Merge(o.reservoir)
	}
}

func (f *filteredExemplarReservoir[N]) Reset() {
	if f == nil {
		return
	}
	if mr, ok := f.reservoir.(mergeableReservoir); ok {
		if !f.concurrentSafe {
			f.reservoirMux.Lock()
			defer f.reservoirMux.Unlock()
		}
		mr.Reset()
	}
}

func (f *filteredExemplarReservoir[N]) IsMergeable() bool {
	if f == nil {
		return false
	}
	_, ok := f.reservoir.(mergeableReservoir)
	return ok
}

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package exemplar

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/internal/reservoir"
)

// FixedSizeReservoirProvider returns a provider of [FixedSizeReservoir].
func FixedSizeReservoirProvider(k int) ReservoirProvider {
	return func(attribute.Set) Reservoir {
		return NewFixedSizeReservoir(k)
	}
}

// NewFixedSizeReservoir returns a [FixedSizeReservoir] that samples at most
// k exemplars. If there are k or less measurements made, the Reservoir will
// sample each one. If there are more than k, the Reservoir will then randomly
// sample all additional measurement with a decreasing probability.
func NewFixedSizeReservoir(k int) *FixedSizeReservoir {
	r := &FixedSizeReservoir{
		storage: make([]measurement, k),
	}
	r.nt.k = k
	r.nt.reset()
	return r
}

var (
	_ Reservoir = &FixedSizeReservoir{}
)

// FixedSizeReservoir is a [Reservoir] that samples at most k exemplars. If
// there are k or less measurements made, the Reservoir will sample each one.
// If there are more than k, the Reservoir will then randomly sample all
// additional measurement with a decreasing probability.
type FixedSizeReservoir struct {
	reservoir.ConcurrentSafe
	mu      sync.Mutex
	storage []measurement
	nt      nextTracker
}

// Offer accepts the parameters associated with a measurement. The
// parameters will be stored as an exemplar if the Reservoir decides to
// sample the measurement.
//
// The passed ctx needs to contain any baggage or span that were active
// when the measurement was made. This information may be used by the
// Reservoir in making a sampling decision.
//
// The time t is the time when the measurement was made. The v and a
// parameters are the value and dropped (filtered) attributes of the
// measurement respectively.
func (r *FixedSizeReservoir) Offer(ctx context.Context, t time.Time, n Value, a []attribute.KeyValue) {
	if len(r.storage) == 0 {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	sampled, idx := r.nt.shouldSample()
	if sampled {
		r.storage[idx].store(ctx, t, n, a)
	}
}

// Collect returns all the held exemplars.
//
// The stored exemplars are preserved after this call, but the sampling state is reset.
func (r *FixedSizeReservoir) Collect(dest *[]Exemplar) {
	if len(r.storage) == 0 {
		*dest = (*dest)[:0]
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	*dest = reset(*dest, len(r.storage), len(r.storage))
	var n int
	for i := range r.storage {
		if r.storage[i].exemplar(&(*dest)[n]) {
			n++
		}
	}
	*dest = (*dest)[:n]
	// Call reset here even though it will reset r.count and restart the random
	// number series. This will persist any old exemplars as long as no new
	// measurements are offered, but it will also prioritize those new
	// measurements that are made over the older collection cycle ones.
	r.nt.reset()
}

type indexedMeasurement struct {
	idx int
	m   measurement
}

var indexedMeasurementPool = sync.Pool{
	New: func() any {
		s := make([]indexedMeasurement, 0, 16)
		return &s
	},
}

// Merge merges the newly sampled exemplars from other (delta) into r (cumulative accumulator).
func (r *FixedSizeReservoir) Merge(other Reservoir) {
	if r == nil || other == nil || r == other {
		return
	}
	o, ok := other.(*FixedSizeReservoir)
	if !ok || o == nil || len(r.storage) == 0 || len(r.storage) != len(o.storage) {
		return
	}

	ptr := indexedMeasurementPool.Get().(*[]indexedMeasurement)
	items := (*ptr)[:0]
	defer func() {
		clear(items)
		*ptr = items[:0]
		indexedMeasurementPool.Put(ptr)
	}()

	o.mu.Lock()
	for i := range o.storage {
		if o.storage[i].valid {
			items = append(items, indexedMeasurement{idx: i, m: o.storage[i]})
		}
	}
	o.mu.Unlock()

	if len(items) == 0 {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, item := range items {
		if item.idx < len(r.storage) {
			r.storage[item.idx] = item.m
		}
	}
}

// Reset resets the reservoir's stored exemplars and sampling state.
func (r *FixedSizeReservoir) Reset() {
	if r == nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.storage {
		r.storage[i] = measurement{}
	}
	r.nt.reset()
}

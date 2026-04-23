// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package exemplar // import "go.opentelemetry.io/otel/sdk/metric/exemplar"

import (
	"context"
	"sync/atomic"
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
	return newFixedSizeReservoir(newStorage(k))
}

var _ Reservoir = &FixedSizeReservoir{}

// FixedSizeReservoir is a [Reservoir] that samples at most k exemplars. If
// there are k or less measurements made, the Reservoir will sample each one.
// If there are more than k, the Reservoir will then randomly sample all
// additional measurement with a decreasing probability.
type FixedSizeReservoir struct {
	reservoir.ConcurrentSafe
	*storage

	// count is the shared atomic counter used for round-robin distribution.
	count atomic.Int64
}

func newFixedSizeReservoir(s *storage) *FixedSizeReservoir {
	return &FixedSizeReservoir{
		storage: s,
	}
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
	// Offer delegates the sampling decision to the bucket's Algorithm L logic.
	// See storage.go for details on Algorithm L.

	if cap(r.measurements) == 0 {
		return
	}
	count := r.count.Add(1)
	idx := int(count % int64(cap(r.measurements)))
	r.storage.measurements[idx].offer(ctx, t, n, a)
}

// Collect returns all the held exemplars.
//
// The Reservoir state is preserved after this call.
func (r *FixedSizeReservoir) Collect(dest *[]Exemplar) {
	r.storage.Collect(dest)
}

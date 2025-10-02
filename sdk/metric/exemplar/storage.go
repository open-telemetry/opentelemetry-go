// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package exemplar // import "go.opentelemetry.io/otel/sdk/metric/exemplar"

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// storage is an exemplar storage for [Reservoir] implementations.
type storage struct {
	mu sync.Mutex
	// measurements are the measurements sampled.
	//
	// This does not use []metricdata.Exemplar because it potentially would
	// require an allocation for trace and span IDs in the hot path of Offer.
	measurements []atomic.Value
}

func newStorage(n int) *storage {
	return &storage{measurements: make([]atomic.Value, n)}
}

func (r *storage) store(idx int, m *measurement) {
	old := r.measurements[idx].Swap(m)
	if old != nil {
		mPool.Put(old)
	}
}

// Collect returns all the held exemplars.
//
// The Reservoir state is preserved after this call.
func (r *storage) Collect(dest *[]Exemplar) {
	r.mu.Lock()
	defer r.mu.Unlock()
	*dest = reset(*dest, len(r.measurements), len(r.measurements))
	var n int
	for _, val := range r.measurements {
		loaded := val.Load()
		if loaded == nil {
			continue
		}
		m := loaded.(*measurement)
		if !m.valid {
			continue
		}

		m.exemplar(&(*dest)[n])
		n++
	}
	*dest = (*dest)[:n]
}

// measurement is a measurement made by a telemetry system.
type measurement struct {
	// FilteredAttributes are the attributes dropped during the measurement.
	FilteredAttributes []attribute.KeyValue
	// Time is the time when the measurement was made.
	Time time.Time
	// Value is the value of the measurement.
	Value Value
	// SpanContext is the SpanContext active when a measurement was made.
	Ctx context.Context

	valid bool
}

var mPool = sync.Pool{
	New: func() any {
		return &measurement{}
	},
}

// newMeasurement returns a new non-empty Measurement.
func newMeasurement(ctx context.Context, ts time.Time, v Value, droppedAttr []attribute.KeyValue) *measurement {
	m := mPool.Get().(*measurement)
	m.FilteredAttributes = droppedAttr
	m.Time = ts
	m.Value = v
	m.Ctx = ctx
	m.valid = true
	return m
}

// exemplar returns m as an [Exemplar].
func (m measurement) exemplar(dest *Exemplar) {
	dest.FilteredAttributes = m.FilteredAttributes
	dest.Time = m.Time
	dest.Value = m.Value

	sc := trace.SpanContextFromContext(m.Ctx)
	if sc.HasTraceID() {
		traceID := sc.TraceID()
		dest.TraceID = traceID[:]
	} else {
		dest.TraceID = dest.TraceID[:0]
	}

	if sc.HasSpanID() {
		spanID := sc.SpanID()
		dest.SpanID = spanID[:]
	} else {
		dest.SpanID = dest.SpanID[:0]
	}
}

func reset[T any](s []T, length, capacity int) []T {
	if cap(s) < capacity {
		return make([]T, length, capacity)
	}
	return s[:length]
}

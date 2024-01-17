// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package exemplar // import "go.opentelemetry.io/otel/sdk/metric/internal/exemplar"

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/trace"
)

type fixedRes[N int64 | float64] struct {
	// store are the measurements sampled.
	//
	// This does not use []metricdata.Exemplar because it potentially would
	// require an allocation for trace and span IDs in the hot path of Offer.
	store []measurement[N]
}

func newFixedRes[N int64 | float64](n int) *fixedRes[N] {
	return &fixedRes[N]{store: make([]measurement[N], n)}
}

func (r *fixedRes[N]) Collect(dest *[]metricdata.Exemplar[N]) {
	*dest = reset(*dest, len(r.store), len(r.store))
	var n int
	for _, m := range r.store {
		if m.Empty() {
			continue
		}

		m.Exemplar(&(*dest)[n])
		n++
	}
	*dest = (*dest)[:n]
}

func (r *fixedRes[N]) Flush(dest *[]metricdata.Exemplar[N]) {
	*dest = reset(*dest, len(r.store), len(r.store))
	var n int
	for i, m := range r.store {
		if m.Empty() {
			continue
		}

		m.Exemplar(&(*dest)[n])
		n++

		// Reset.
		r.store[i] = measurement[N]{}
	}
	*dest = (*dest)[:n]
}

// measurement is a measurement made by a telemetry system.
type measurement[N int64 | float64] struct {
	Attributes  []attribute.KeyValue
	Time        time.Time
	Value       N
	SpanContext trace.SpanContext

	valid bool
}

// newMeasurement returns a new non-empty Measurement.
func newMeasurement[N int64 | float64](ctx context.Context, ts time.Time, v N, droppedAttr []attribute.KeyValue) measurement[N] {
	return measurement[N]{
		Attributes:  droppedAttr,
		Time:        ts,
		Value:       v,
		SpanContext: trace.SpanContextFromContext(ctx),
		valid:       true,
	}
}

// Empty returns false if m represents a measurement made by a telemetry
// system, otherwise it returns true when m is its zero-value.
func (m measurement[N]) Empty() bool { return !m.valid }

// Exemplar returns m as a [metricdata.Exemplar].
func (m measurement[N]) Exemplar(dest *metricdata.Exemplar[N]) {
	dest.FilteredAttributes = m.Attributes
	dest.Time = m.Time
	dest.Value = m.Value

	if m.SpanContext.HasTraceID() {
		traceID := m.SpanContext.TraceID()
		dest.TraceID = traceID[:]
	} else {
		dest.TraceID = dest.TraceID[:0]
	}

	if m.SpanContext.HasSpanID() {
		spanID := m.SpanContext.SpanID()
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

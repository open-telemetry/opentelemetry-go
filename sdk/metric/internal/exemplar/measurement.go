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
	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/trace"
)

// measurement is a measurement made by a telemetry system.
type measurement[N int64 | float64] struct {
	Attributes  attribute.Set
	Time        time.Time
	Value       N
	SpanContext trace.SpanContext

	valid bool
}

// newMeasurement returns a new non-empty Measurement.
func newMeasurement[N int64 | float64](ctx context.Context, ts time.Time, v N, measuredAttr attribute.Set) measurement[N] {
	return measurement[N]{
		Attributes:  measuredAttr,
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
func (m measurement[N]) Exemplar(dest *metricdata.Exemplar[N], recorded attribute.Set) {
	// Note: A more optimal solution would be to store the filtered attributes
	// when the exemplar is recorded, instead of re-calculating here. That
	// approach isn't implemented though because it contrary to the OTel
	// specification definition of a Reservoir, which is defined to accept the
	// complete set of measured attributes.
	dropped(&dest.FilteredAttributes, m.Attributes, recorded)
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

// dropped returns the attribute that were measured, but not included in the
// recorded attributes.
func dropped(dest *[]attribute.KeyValue, measured, recorded attribute.Set) {
	measN := measured.Len()
	recN := recorded.Len()

	n := measN - recN
	switch {
	case n < 0:
		// recorded should only ever be the filtered set of measured. Abandon
		// instead of panicking.
		global.Warn(
			"invalid measured attributes for exemplar, dropping",
			"measured", measured,
			"recorded", recorded,
		)
		fallthrough
	case n == 0:
		// Nothing dropped.
		*dest = (*dest)[:0]
		return
	}
	*dest = reset(*dest, n, n)

	measIter := measured.Iter()
	recIter := recorded.Iter()

	var i int
	recIter.Next()
	for measIter.Next() {
		m := measIter.Attribute()
		r := recIter.Attribute()

		if m == r {
			recIter.Next()
			continue
		}

		(*dest)[i] = m
		i++
	}
}

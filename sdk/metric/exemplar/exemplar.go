// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package exemplar

import (
	"time"

	"go.opentelemetry.io/otel/attribute"
)

// Exemplar is a measurement sampled from a timeseries providing a typical
// example.
type Exemplar struct {
	// FilteredAttributes are the attributes recorded with the measurement but
	// filtered out of the timeseries' aggregated data.
	FilteredAttributes []attribute.KeyValue
	// Time is the time when the measurement was recorded.
	Time time.Time
	// Value is the measured value.
	Value Value
	// SpanID is the ID of the span that was active during the measurement.
	// It is empty if the span context captured with the measurement has no
	// span ID. It is not cleared for unsampled spans; the exemplar filter
	// determines which measurements are offered.
	SpanID []byte `json:",omitempty"`
	// TraceID is the ID of the trace the active span belonged to during the
	// measurement. It is empty if the span context captured with the
	// measurement has no trace ID. It is not cleared for unsampled spans; the
	// exemplar filter determines which measurements are offered.
	TraceID []byte `json:",omitempty"`
}

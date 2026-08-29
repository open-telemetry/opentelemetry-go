// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Int64CounterBinder is implemented by Int64Counter implementations that can
// bind attributes before recording measurements.
type Int64CounterBinder interface {
	// Bind returns an Int64Counter bound to attrs.
	//
	// Attributes passed when recording with the returned counter are merged
	// with attrs and take precedence for duplicate keys. Supplying attributes
	// when recording can negate the performance benefit of binding.
	//
	// Bind is safe to call concurrently. Implementations must not retain the
	// provided slice.
	Bind(attrs ...attribute.KeyValue) metric.Int64Counter
}

// Float64CounterBinder is implemented by Float64Counter implementations that
// can bind attributes before recording measurements.
type Float64CounterBinder interface {
	Bind(attrs ...attribute.KeyValue) metric.Float64Counter
}

// Int64UpDownCounterBinder is implemented by Int64UpDownCounter
// implementations that can bind attributes before recording measurements.
type Int64UpDownCounterBinder interface {
	Bind(attrs ...attribute.KeyValue) metric.Int64UpDownCounter
}

// Float64UpDownCounterBinder is implemented by Float64UpDownCounter
// implementations that can bind attributes before recording measurements.
type Float64UpDownCounterBinder interface {
	Bind(attrs ...attribute.KeyValue) metric.Float64UpDownCounter
}

// Int64HistogramBinder is implemented by Int64Histogram implementations that
// can bind attributes before recording measurements.
type Int64HistogramBinder interface {
	Bind(attrs ...attribute.KeyValue) metric.Int64Histogram
}

// Float64HistogramBinder is implemented by Float64Histogram implementations
// that can bind attributes before recording measurements.
type Float64HistogramBinder interface {
	Bind(attrs ...attribute.KeyValue) metric.Float64Histogram
}

// Int64GaugeBinder is implemented by Int64Gauge implementations that can bind
// attributes before recording measurements.
type Int64GaugeBinder interface {
	Bind(attrs ...attribute.KeyValue) metric.Int64Gauge
}

// Float64GaugeBinder is implemented by Float64Gauge implementations that can
// bind attributes before recording measurements.
type Float64GaugeBinder interface {
	Bind(attrs ...attribute.KeyValue) metric.Float64Gauge
}

// Finisher is implemented by synchronous instruments that support ending the
// lifetime of the series identified by an exact attribute set.
type Finisher interface {
	// Finish ends the active series identified by attrs. The final value is
	// eligible for one collection. Calling Finish for an inactive series is a
	// no-op.
	//
	// Finish is safe to call concurrently.
	Finish(context.Context, ...attribute.KeyValue)
}

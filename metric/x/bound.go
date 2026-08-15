// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Bound instruments still accept attribute options (e.g., [metric.WithAttributes]
// or [metric.WithAttributeSet]) when recording measurements, but passing
// attributes at record time can negate the performance benefits of binding.

// Int64CounterBinder is an interface for instruments that support binding attributes ahead of time.
type Int64CounterBinder interface {
	// Bind returns an Int64Counter bound to the provided attributes.
	// Measurements recorded on the returned counter will use the bound attributes.
	Bind(attrs ...attribute.KeyValue) metric.Int64Counter
}

// Float64CounterBinder is an interface for instruments that support binding attributes ahead of time.
type Float64CounterBinder interface {
	// Bind returns a Float64Counter bound to the provided attributes.
	// Measurements recorded on the returned counter will use the bound attributes.
	Bind(attrs ...attribute.KeyValue) metric.Float64Counter
}

// Int64UpDownCounterBinder is an interface for instruments that support binding attributes ahead of time.
type Int64UpDownCounterBinder interface {
	// Bind returns an Int64UpDownCounter bound to the provided attributes.
	// Measurements recorded on the returned counter will use the bound attributes.
	Bind(attrs ...attribute.KeyValue) metric.Int64UpDownCounter
}

// Float64UpDownCounterBinder is an interface for instruments that support binding attributes ahead of time.
type Float64UpDownCounterBinder interface {
	// Bind returns a Float64UpDownCounter bound to the provided attributes.
	// Measurements recorded on the returned counter will use the bound attributes.
	Bind(attrs ...attribute.KeyValue) metric.Float64UpDownCounter
}

// Int64HistogramBinder is an interface for instruments that support binding attributes ahead of time.
type Int64HistogramBinder interface {
	// Bind returns an Int64Histogram bound to the provided attributes.
	// Measurements recorded on the returned histogram will use the bound attributes.
	Bind(attrs ...attribute.KeyValue) metric.Int64Histogram
}

// Float64HistogramBinder is an interface for instruments that support binding attributes ahead of time.
type Float64HistogramBinder interface {
	// Bind returns a Float64Histogram bound to the provided attributes.
	// Measurements recorded on the returned histogram will use the bound attributes.
	Bind(attrs ...attribute.KeyValue) metric.Float64Histogram
}

// Int64GaugeBinder is an interface for instruments that support binding attributes ahead of time.
type Int64GaugeBinder interface {
	// Bind returns an Int64Gauge bound to the provided attributes.
	// Measurements recorded on the returned gauge will use the bound attributes.
	Bind(attrs ...attribute.KeyValue) metric.Int64Gauge
}

// Float64GaugeBinder is an interface for instruments that support binding attributes ahead of time.
type Float64GaugeBinder interface {
	// Bind returns a Float64Gauge bound to the provided attributes.
	// Measurements recorded on the returned gauge will use the bound attributes.
	Bind(attrs ...attribute.KeyValue) metric.Float64Gauge
}

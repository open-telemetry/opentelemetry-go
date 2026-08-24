// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Float64Binder is an interface that can be implemented by instruments that support
// binding attributes ahead of time.
//
// Bound instruments associate a fixed set of attributes with the instrument upfront,
// eliminating attribute set creation and map lookups on subsequent measurements.
type Float64Binder interface {
	// Bind returns a metric.Float64Counter bound to the provided attributes.
	//
	// The returned counter is safe for concurrent use by multiple goroutines.
	// Calling Add without additional attributes records measurements directly to the
	// pre-aggregated series with minimal overhead. If Add is called with additional
	// attributes, they are merged with the bound attributes on a slower fallback path.
	//
	// Bind does not mutate the passed attrs slice.
	Bind(attrs ...attribute.KeyValue) metric.Float64Counter
}

// Int64Binder is an interface that can be implemented by instruments that support
// binding attributes ahead of time.
//
// Bound instruments associate a fixed set of attributes with the instrument upfront,
// eliminating attribute set creation and map lookups on subsequent measurements.
type Int64Binder interface {
	// Bind returns a metric.Int64Counter bound to the provided attributes.
	//
	// The returned counter is safe for concurrent use by multiple goroutines.
	// Calling Add without additional attributes records measurements directly to the
	// pre-aggregated series with minimal overhead. If Add is called with additional
	// attributes, they are merged with the bound attributes on a slower fallback path.
	//
	// Bind does not mutate the passed attrs slice.
	Bind(attrs ...attribute.KeyValue) metric.Int64Counter
}

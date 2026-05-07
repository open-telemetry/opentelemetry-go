// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package x // import "go.opentelemetry.io/otel/metric/x"

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Float64Binder is an interface that can be implemented by instruments that support
// binding attributes ahead of time.
type Float64Binder interface {
	// Bind returns a metric.Float64Counter for the given attributes.
	// The returned counter is bound to the attributes and should be optimized
	// for performance by avoiding map lookups on every Add call.
	Bind(attrs ...attribute.KeyValue) metric.Float64Counter
}

// Int64Binder is an interface that can be implemented by instruments that support
// binding attributes ahead of time.
type Int64Binder interface {
	// Bind returns a metric.Int64Counter for the given attributes.
	// The returned counter is bound to the attributes and should be optimized
	// for performance by avoiding map lookups on every Add call.
	Bind(attrs ...attribute.KeyValue) metric.Int64Counter
}

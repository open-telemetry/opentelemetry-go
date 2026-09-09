// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric // import "go.opentelemetry.io/otel/sdk/metric"

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
)

const (
	metricFilterAccept = iota
	metricFilterDrop
	metricFilterAcceptPartial
)

const (
	metricFilterAttrAccept = iota
	metricFilterAttrDrop
)

// metricFilter is an internal, private mirror of the experimental metricFilter
// interface from the 'x' package. It is duplicated here to avoid package
// import cycles and layer violations while leveraging Go's implicit interfaces.
type metricFilter interface {
	Experimental()
	TestMetric(scope instrumentation.Scope, name string, kind InstrumentKind, unit string) int
	TestAttributes(
		scope instrumentation.Scope,
		name string,
		kind InstrumentKind,
		unit string,
		attrs []attribute.KeyValue,
	) int
}

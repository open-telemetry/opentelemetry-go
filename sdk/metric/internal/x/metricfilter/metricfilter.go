// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// Package metricfilter provides the experimental MetricFilter API for
// filtering aggregated metric data points during collection.
package metricfilter // import "go.opentelemetry.io/otel/sdk/metric/internal/x/metricfilter"

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric"
)

// MetricFilter defines the interface which enables the MetricReader's
// registered MetricProducers or the SDK's MetricProducer to filter
// aggregated data points (Metric Points) inside its Produce operation.
// The filtering is done at the MetricProducer for performance reasons.
type MetricFilter interface {
	// TestMetric is called once for every metric stream, in each
	// MetricProducer Produce operation.
	TestMetric(
		instrumentationScope instrumentation.Scope,
		name string,
		kind metric.InstrumentKind,
		unit string,
	) Result
	// TestAttributes determines for a given metric stream and attribute set
	// if it should be allowed or filtered out.
	// This operation should only be called if TestMetric operation returned
	// AcceptPartial for the given metric stream arguments.
	TestAttributes(
		instrumentationScope instrumentation.Scope,
		name string,
		kind metric.InstrumentKind,
		unit string,
		attributes []attribute.KeyValue,
	) AttributesFilterResult
}

// Result is an enumeration used to decide whether to accept,
// drop, or partially accept a metric stream.
type Result int

// AttributesFilterResult is an enumeration used to decide whether to accept
// or drop an attribute set.
type AttributesFilterResult int

const (
	// Accept means all attributes of the given metric stream are allowed (not to be filtered).
	Accept Result = iota
	// Drop means all attributes of the given metric stream are NOT allowed (filtered out - dropped).
	Drop
	// AcceptPartial means some attributes are allowed and some aren't, hence
	// TestAttributes operation must be called for each attribute set of that instrument.
	AcceptPartial
)

const (
	// AttrAccept means the given attributes are allowed (not to be filtered).
	AttrAccept AttributesFilterResult = iota
	// AttrDrop means the given attributes are NOT allowed (filtered out - dropped).
	AttrDrop
)

// metricFilterOption wraps a MetricFilter so it can be passed as a metric.ReaderOption.
type metricFilterOption struct {
	metric.ReaderOption
	filter MetricFilter
}

// Experimental indicates that this configuration option is part of an experimental API.
func (metricFilterOption) Experimental() {}

// TestMetric delegates the metric stream filtering decision to the underlying MetricFilter,
// returning the integer representation of the resulting Result.
func (o metricFilterOption) TestMetric(
	scope instrumentation.Scope,
	name string,
	kind metric.InstrumentKind,
	unit string,
) int {
	return int(o.filter.TestMetric(scope, name, kind, unit))
}

// TestAttributes delegates the attribute set filtering decision to the underlying MetricFilter,
// returning the integer representation of the resulting AttributesFilterResult.
func (o metricFilterOption) TestAttributes(
	scope instrumentation.Scope,
	name string,
	kind metric.InstrumentKind,
	unit string,
	attrs []attribute.KeyValue,
) int {
	return int(o.filter.TestAttributes(scope, name, kind, unit, attrs))
}

// WithMetricFilter creates a metric.ReaderOption that applies the provided MetricFilter.
// This allows you to selectively accept, drop, or partially filter metric streams and
// their specific attribute sets at the reader level.
func WithMetricFilter(f MetricFilter) metric.ReaderOption {
	return metricFilterOption{filter: f}
}

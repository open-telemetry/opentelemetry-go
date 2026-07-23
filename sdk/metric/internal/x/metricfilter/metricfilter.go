package metricfilter

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
	TestMetric(instrumentationScope instrumentation.Scope, name string, kind metric.InstrumentKind, unit string) MetricFilterResult
	// TestAttributes determines for a given metric stream and attribute set
	// if it should be allowed or filtered out.
	// This operation should only be called if TestMetric operation returned
	// Accept_Partial for the given metric stream arguments.
	TestAttributes(instrumentationScope instrumentation.Scope, name string, kind metric.InstrumentKind, unit string, attributes []attribute.KeyValue) AttributesFilterResult
}

// MetricFilterResult is an enumeration used to decide whether to accept,
// drop, or partially accept a metric stream.
type MetricFilterResult int

// AttributesFilterResult is an enumeration used to decide whether to accept
// or drop an attribute set.
type AttributesFilterResult int

const (
	// Accept means all attributes of the given metric stream are allowed (not to be filtered).
	Accept MetricFilterResult = iota
	// Drop means all attributes of the given metric stream are NOT allowed (filtered out - dropped).
	Drop
	// Accept_Partial means some attributes are allowed and some aren't, hence
	// TestAttributes operation must be called for each attribute set of that instrument.
	Accept_Partial
)

const (
	// AttrAccept means the given attributes are allowed (not to be filtered).
	AttrAccept AttributesFilterResult = iota
	// AttrDrop means the given attributes are NOT allowed (filtered out - dropped).
	AttrDrop
)

type metricFilterOption struct {
	metric.ReaderOption
	filter MetricFilter
}

func WithMetricFilter(f MetricFilter) metric.ReaderOption {
	return metricFilterOption{filter: f}
}

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// testMetricFilterOption is a fake ReaderOption that also satisfies the
// internal metricFilter interface. It lets stable-package tests exercise the
// MetricFilter plumbing without importing the experimental x package.
type testMetricFilterOption struct {
	ReaderOption
	testMetric     func(instrumentation.Scope, string, InstrumentKind, string) int
	testAttributes func(instrumentation.Scope, string, InstrumentKind, string, []attribute.KeyValue) int
}

func (testMetricFilterOption) Experimental() {}

func (o testMetricFilterOption) TestMetric(
	scope instrumentation.Scope,
	name string,
	kind InstrumentKind,
	unit string,
) int {
	return o.testMetric(scope, name, kind, unit)
}

func (o testMetricFilterOption) TestAttributes(
	scope instrumentation.Scope,
	name string,
	kind InstrumentKind,
	unit string,
	attrs []attribute.KeyValue,
) int {
	return o.testAttributes(scope, name, kind, unit, attrs)
}

// testExporter is a minimal Exporter for tests that need a PeriodicReader.
type testExporter struct{}

func (testExporter) Temporality(InstrumentKind) metricdata.Temporality {
	return metricdata.CumulativeTemporality
}

func (testExporter) Aggregation(InstrumentKind) Aggregation { return AggregationDefault{} }

func (testExporter) Export(context.Context, *metricdata.ResourceMetrics) error {
	return nil
}
func (testExporter) ForceFlush(context.Context) error { return nil }
func (testExporter) Shutdown(context.Context) error   { return nil }

// experimentalReaderOption embeds ReaderOption and only implements
// Experimental(). It is used to verify that experimental reader options do not
// panic when the embedded option interface is never applied.
type experimentalReaderOption struct{ ReaderOption }

func (experimentalReaderOption) Experimental() {}

func sumDataPointCount(rm *metricdata.ResourceMetrics, name string) int {
	m := findMetricByName(rm, name)
	if m == nil {
		return 0
	}
	sum, ok := m.Data.(metricdata.Sum[int64])
	if !ok {
		return -1
	}
	return len(sum.DataPoints)
}

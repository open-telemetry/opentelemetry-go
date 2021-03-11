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

package opencensus

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"go.opencensus.io/metric/metricdata"
	"go.opencensus.io/metric/metricexport"
	ocresource "go.opencensus.io/resource"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/number"
	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/export/metric/aggregation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/unit"
)

var errConversion = errors.New("Unable to convert from OpenCensus to OpenTelemetry")

// NewMetricExporter returns an OpenCensus exporter that exports to an
// OpenTelemetry exporter
func NewMetricExporter(base export.Exporter) metricexport.Exporter {
	return &exporter{base: base}
}

// exporter implements the OpenCensus metric Exporter interface using an
// OpenTelemetry base exporter.
type exporter struct {
	base export.Exporter
}

// ExportMetrics implements the OpenCensus metric Exporter interface
func (e *exporter) ExportMetrics(ctx context.Context, metrics []*metricdata.Metric) error {
	return e.base.Export(ctx, &checkpointSet{metrics: metrics})
}

type checkpointSet struct {
	// RWMutex implements locking for the `CheckpointSet` interface.
	sync.RWMutex
	metrics []*metricdata.Metric
}

// ForEach iterates through the CheckpointSet, passing an
// export.Record with the appropriate aggregation to an exporter.
func (d *checkpointSet) ForEach(exporter export.ExportKindSelector, f func(export.Record) error) error {
	for _, m := range d.metrics {
		descriptor, err := convertDescriptor(m.Descriptor)
		if err != nil {
			otel.Handle(err)
			continue
		}
		res := convertResource(m.Resource)
		for _, ts := range m.TimeSeries {
			if len(ts.Points) == 0 {
				continue
			}
			ls, err := convertLabels(m.Descriptor.LabelKeys, ts.LabelValues)
			if err != nil {
				otel.Handle(err)
				continue
			}
			agg, err := newAggregationFromPoints(ts.Points)
			if err != nil {
				otel.Handle(err)
				continue
			}
			if err := f(export.NewRecord(
				&descriptor,
				&ls,
				res,
				agg,
				ts.StartTime,
				agg.end(),
			)); err != nil && !errors.Is(err, aggregation.ErrNoData) {
				return err
			}
		}
	}
	return nil
}

// convertLabels converts from OpenCensus label keys and values to an
// OpenTelemetry label Set.
func convertLabels(keys []metricdata.LabelKey, values []metricdata.LabelValue) (attribute.Set, error) {
	if len(keys) != len(values) {
		return attribute.NewSet(), fmt.Errorf("%w different number of label keys (%d) and values (%d)", errConversion, len(keys), len(values))
	}
	labels := []attribute.KeyValue{}
	for i, lv := range values {
		if !lv.Present {
			continue
		}
		labels = append(labels, attribute.KeyValue{
			Key:   attribute.Key(keys[i].Key),
			Value: attribute.StringValue(lv.Value),
		})
	}
	return attribute.NewSet(labels...), nil
}

// convertResource converts an OpenCensus Resource to an OpenTelemetry Resource
func convertResource(res *ocresource.Resource) *resource.Resource {
	labels := []attribute.KeyValue{}
	if res == nil {
		return nil
	}
	for k, v := range res.Labels {
		labels = append(labels, attribute.KeyValue{Key: attribute.Key(k), Value: attribute.StringValue(v)})
	}
	return resource.NewWithAttributes(labels...)
}

// convertDescriptor converts an OpenCensus Descriptor to an OpenTelemetry Descriptor
func convertDescriptor(ocDescriptor metricdata.Descriptor) (metric.Descriptor, error) {
	var (
		nkind number.Kind
		ikind metric.InstrumentKind
	)
	switch ocDescriptor.Type {
	case metricdata.TypeGaugeInt64:
		nkind = number.Int64Kind
		ikind = metric.ValueObserverInstrumentKind
	case metricdata.TypeGaugeFloat64:
		nkind = number.Float64Kind
		ikind = metric.ValueObserverInstrumentKind
	case metricdata.TypeCumulativeInt64:
		nkind = number.Int64Kind
		ikind = metric.SumObserverInstrumentKind
	case metricdata.TypeCumulativeFloat64:
		nkind = number.Float64Kind
		ikind = metric.SumObserverInstrumentKind
	default:
		// Includes TypeGaugeDistribution, TypeCumulativeDistribution, TypeSummary
		return metric.Descriptor{}, fmt.Errorf("%w; descriptor type: %v", errConversion, ocDescriptor.Type)
	}
	opts := []metric.InstrumentOption{
		metric.WithDescription(ocDescriptor.Description),
		metric.WithInstrumentationName("OpenCensus Bridge"),
	}
	switch ocDescriptor.Unit {
	case metricdata.UnitDimensionless:
		opts = append(opts, metric.WithUnit(unit.Dimensionless))
	case metricdata.UnitBytes:
		opts = append(opts, metric.WithUnit(unit.Bytes))
	case metricdata.UnitMilliseconds:
		opts = append(opts, metric.WithUnit(unit.Milliseconds))
	}
	return metric.NewDescriptor(ocDescriptor.Name, ikind, nkind, opts...), nil
}

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

package opencensus // import "go.opentelemetry.io/otel/bridge/opencensus"

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.opencensus.io/metric/metricdata"
	"go.opencensus.io/metric/metricexport"
	ocresource "go.opencensus.io/resource"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/export"
	"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/sdkapi"
	"go.opentelemetry.io/otel/sdk/resource"
)

var errConversion = errors.New("Unable to convert from OpenCensus to OpenTelemetry")

// NewMetricExporter returns an OpenCensus exporter that exports to an
// OpenTelemetry exporter.
func NewMetricExporter(base export.Exporter) metricexport.Exporter {
	return &exporter{base: base}
}

// exporter implements the OpenCensus metric Exporter interface using an
// OpenTelemetry base exporter.
type exporter struct {
	base export.Exporter
}

// ExportMetrics implements the OpenCensus metric Exporter interface.
func (e *exporter) ExportMetrics(ctx context.Context, metrics []*metricdata.Metric) error {
	res := resource.Empty()
	if len(metrics) != 0 {
		res = convertResource(metrics[0].Resource)
	}
	return e.base.Export(ctx, res, &censusLibraryReader{metrics: metrics})
}

type censusLibraryReader struct {
	metrics []*metricdata.Metric
}

func (r censusLibraryReader) ForEach(readerFunc func(instrumentation.Library, export.Reader) error) error {
	return readerFunc(instrumentation.Library{
		Name: "OpenCensus Bridge",
	}, &metricReader{metrics: r.metrics})
}

type metricReader struct {
	// RWMutex implements locking for the `Reader` interface.
	sync.RWMutex
	metrics []*metricdata.Metric
}

var _ export.Reader = &metricReader{}

// ForEach iterates through the metrics data, synthesizing an
// export.Record with the appropriate aggregation for the exporter.
func (d *metricReader) ForEach(_ aggregation.TemporalitySelector, f func(export.Record) error) error {
	for _, m := range d.metrics {
		descriptor, err := convertDescriptor(m.Descriptor)
		if err != nil {
			otel.Handle(err)
			continue
		}
		for _, ts := range m.TimeSeries {
			if len(ts.Points) == 0 {
				continue
			}
			attrs, err := convertAttrs(m.Descriptor.LabelKeys, ts.LabelValues)
			if err != nil {
				otel.Handle(err)
				continue
			}
			err = recordAggregationsFromPoints(
				ts.Points,
				func(agg aggregation.Aggregation, end time.Time) error {
					return f(export.NewRecord(
						&descriptor,
						&attrs,
						agg,
						ts.StartTime,
						end,
					))
				})
			if err != nil && !errors.Is(err, aggregation.ErrNoData) {
				return err
			}
		}
	}
	return nil
}

// convertAttrs converts from OpenCensus attribute keys and values to an
// OpenTelemetry attribute Set.
func convertAttrs(keys []metricdata.LabelKey, values []metricdata.LabelValue) (attribute.Set, error) {
	if len(keys) != len(values) {
		return attribute.NewSet(), fmt.Errorf("%w different number of label keys (%d) and values (%d)", errConversion, len(keys), len(values))
	}
	attrs := []attribute.KeyValue{}
	for i, lv := range values {
		if !lv.Present {
			continue
		}
		attrs = append(attrs, attribute.KeyValue{
			Key:   attribute.Key(keys[i].Key),
			Value: attribute.StringValue(lv.Value),
		})
	}
	return attribute.NewSet(attrs...), nil
}

// convertResource converts an OpenCensus Resource to an OpenTelemetry Resource
// Note: the ocresource.Resource Type field is not used.
func convertResource(res *ocresource.Resource) *resource.Resource {
	attrs := []attribute.KeyValue{}
	if res == nil {
		return nil
	}
	for k, v := range res.Labels {
		attrs = append(attrs, attribute.KeyValue{Key: attribute.Key(k), Value: attribute.StringValue(v)})
	}
	return resource.NewSchemaless(attrs...)
}

// convertDescriptor converts an OpenCensus Descriptor to an OpenTelemetry Descriptor.
func convertDescriptor(ocDescriptor metricdata.Descriptor) (sdkapi.Descriptor, error) {
	var (
		nkind number.Kind
		ikind sdkapi.InstrumentKind
	)
	switch ocDescriptor.Type {
	case metricdata.TypeGaugeInt64:
		nkind = number.Int64Kind
		ikind = sdkapi.GaugeObserverInstrumentKind
	case metricdata.TypeGaugeFloat64:
		nkind = number.Float64Kind
		ikind = sdkapi.GaugeObserverInstrumentKind
	case metricdata.TypeCumulativeInt64:
		nkind = number.Int64Kind
		ikind = sdkapi.CounterObserverInstrumentKind
	case metricdata.TypeCumulativeFloat64:
		nkind = number.Float64Kind
		ikind = sdkapi.CounterObserverInstrumentKind
	default:
		// Includes TypeGaugeDistribution, TypeCumulativeDistribution, TypeSummary
		return sdkapi.Descriptor{}, fmt.Errorf("%w; descriptor type: %v", errConversion, ocDescriptor.Type)
	}
	opts := []instrument.Option{
		instrument.WithDescription(ocDescriptor.Description),
	}
	switch ocDescriptor.Unit {
	case metricdata.UnitDimensionless:
		opts = append(opts, instrument.WithUnit(unit.Dimensionless))
	case metricdata.UnitBytes:
		opts = append(opts, instrument.WithUnit(unit.Bytes))
	case metricdata.UnitMilliseconds:
		opts = append(opts, instrument.WithUnit(unit.Milliseconds))
	}
	cfg := instrument.NewConfig(opts...)
	return sdkapi.NewDescriptor(ocDescriptor.Name, ikind, nkind, cfg.Description(), cfg.Unit()), nil
}

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

// Package metrictransform provides translations for opentelemetry-go concepts and
// structures to otlp structures.
package metrictransform // import "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/internal/metrictransform"

import (
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/aggregation"
	"go.opentelemetry.io/otel/sdk/metric/number"
	"go.opentelemetry.io/otel/sdk/metric/reader"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	metricpb "go.opentelemetry.io/proto/otlp/metrics/v1"
)

var (
	// ErrUnimplementedAgg is returned when a transformation of an unimplemented
	// aggregator is attempted.
	ErrUnimplementedAgg = errors.New("unimplemented aggregator")

	// ErrIncompatibleAgg is returned when
	// aggregation.Kind implies an interface conversion that has
	// failed
	ErrIncompatibleAgg = errors.New("incompatible aggregation type")

	// ErrUnknownValueType is returned when a transformation of an unknown value
	// is attempted.
	ErrUnknownValueType = errors.New("invalid value type")

	// ErrContextCanceled is returned when a context cancellation halts a
	// transformation.
	ErrContextCanceled = errors.New("context canceled")

	// ErrTransforming is returned when an unexected error is encountered transforming.
	ErrTransforming = errors.New("transforming failed")
)

// toNanos returns the number of nanoseconds since the UNIX epoch.
func toNanos(t time.Time) uint64 {
	if t.IsZero() {
		return 0
	}
	return uint64(t.UnixNano())
}

func TransformMetrics(metrics reader.Metrics) *metricpb.ResourceMetrics {
	var sms []*metricpb.ScopeMetrics

	for _, scope := range metrics.Scopes {
		metrics := []*metricpb.Metric{}

		for _, inst := range scope.Instruments {
			metric := transformInstrument(inst)
			if metric != nil {
				metrics = append(metrics, metric)
			}
		}

		sms = append(sms, &metricpb.ScopeMetrics{
			Scope: &commonpb.InstrumentationScope{
				Name:    scope.Library.Name,
				Version: scope.Library.Version,
			},
			Metrics:   metrics,
			SchemaUrl: scope.Library.SchemaURL,
		})
	}

	return &metricpb.ResourceMetrics{
		Resource:     Resource(metrics.Resource),
		SchemaUrl:    metrics.Resource.SchemaURL(),
		ScopeMetrics: sms,
	}
}

func transformInstrument(inst reader.Instrument) *metricpb.Metric {
	metric := &metricpb.Metric{
		Name:        inst.Descriptor.Name,
		Description: inst.Descriptor.Description,
		Unit:        string(inst.Descriptor.Unit),
	}
	if len(inst.Points) > 0 {
		switch inst.Points[0].Aggregation.Category() {
		case aggregation.GaugeCategory:
			metric.Data = gaugePoints(inst.Points, inst.Descriptor.NumberKind)
		case aggregation.MonotonicSumCategory, aggregation.NonMonotonicSumCategory:
			metric.Data = sumPoints(inst.Points, inst.Descriptor.NumberKind, inst.Temporality)
		case aggregation.HistogramCategory:
			metric.Data = histogramPoints(inst.Points, inst.Descriptor.NumberKind, inst.Temporality)
		}

	}
	return metric
}

func gaugePoints(points []reader.Point, kind number.Kind) *metricpb.Metric_Gauge {
	if len(points) == 0 {
		return nil
	}
	dataPoints := make([]*metricpb.NumberDataPoint, len(points))

	for i := range points {
		dataPoint := &metricpb.NumberDataPoint{
			Attributes:        Iterator(points[i].Attributes.Iter()),
			StartTimeUnixNano: toNanos(points[i].Start),
			TimeUnixNano:      toNanos(points[i].End),
		}
		gauge, ok := points[i].Aggregation.(aggregation.Gauge)
		if !ok {
			otel.Handle(ErrIncompatibleAgg)
			return nil
		}
		switch kind {
		case number.Int64Kind:
			dataPoint.Value = &metricpb.NumberDataPoint_AsInt{
				AsInt: int64(gauge.Gauge()),
			}
		case number.Float64Kind:
			dataPoint.Value = &metricpb.NumberDataPoint_AsDouble{
				AsDouble: gauge.Gauge().CoerceToFloat64(kind),
			}
		default:
			otel.Handle(fmt.Errorf("%w: %v", ErrUnknownValueType, kind))
			return nil
		}

		dataPoints[i] = dataPoint

	}

	return &metricpb.Metric_Gauge{
		Gauge: &metricpb.Gauge{
			DataPoints: dataPoints,
		},
	}
}

func sumPoints(points []reader.Point, kind number.Kind, temporality aggregation.Temporality) *metricpb.Metric_Sum {
	if len(points) == 0 {
		return nil
	}

	dataPoints := make([]*metricpb.NumberDataPoint, len(points))

	for i := range points {
		dataPoint := &metricpb.NumberDataPoint{
			Attributes:        Iterator(points[i].Attributes.Iter()),
			StartTimeUnixNano: toNanos(points[i].Start),
			TimeUnixNano:      toNanos(points[i].End),
		}
		sum, ok := points[i].Aggregation.(aggregation.Sum)
		if !ok {
			otel.Handle(ErrIncompatibleAgg)
			return nil
		}
		switch kind {
		case number.Int64Kind:
			dataPoint.Value = &metricpb.NumberDataPoint_AsInt{
				AsInt: int64(sum.Sum()),
			}
		case number.Float64Kind:
			dataPoint.Value = &metricpb.NumberDataPoint_AsDouble{
				AsDouble: sum.Sum().CoerceToFloat64(kind),
			}
		default:
			otel.Handle(fmt.Errorf("%w: %v", ErrUnknownValueType, kind))
			return nil
		}
		dataPoints[i] = dataPoint

	}
	return &metricpb.Metric_Sum{
		Sum: &metricpb.Sum{
			DataPoints:             dataPoints,
			AggregationTemporality: sdkTemporalityToTemporality(temporality),
			IsMonotonic:            isMonotonic(points[0].Aggregation),
		},
	}
}

func isMonotonic(agg aggregation.Aggregation) bool {
	return agg.Category() == aggregation.MonotonicSumCategory
}

func sdkTemporalityToTemporality(temporality aggregation.Temporality) metricpb.AggregationTemporality {
	switch temporality {
	case aggregation.DeltaTemporality:
		return metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_DELTA
	case aggregation.CumulativeTemporality:
		return metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE
	}
	return metricpb.AggregationTemporality_AGGREGATION_TEMPORALITY_UNSPECIFIED
}

func histogramPoints(points []reader.Point, kind number.Kind, temporality aggregation.Temporality) *metricpb.Metric_Histogram {
	if len(points) == 0 {
		return nil
	}
	dataPoints := make([]*metricpb.HistogramDataPoint, len(points))

	for i := range points {

		histogram, ok := points[i].Aggregation.(aggregation.Histogram)
		if !ok {
			otel.Handle(ErrIncompatibleAgg)
			return nil
		}
		sum := histogram.Sum().CoerceToFloat64(kind)

		dataPoint := &metricpb.HistogramDataPoint{
			Attributes:        Iterator(points[i].Attributes.Iter()),
			StartTimeUnixNano: toNanos(points[i].Start),
			TimeUnixNano:      toNanos(points[i].End),
			BucketCounts:      histogram.Histogram().Counts,
			ExplicitBounds:    histogram.Histogram().Boundaries,
			Count:             histogram.Count(),
			Sum:               &sum,
		}

		dataPoints[i] = dataPoint

	}
	return &metricpb.Metric_Histogram{
		Histogram: &metricpb.Histogram{

			DataPoints:             dataPoints,
			AggregationTemporality: sdkTemporalityToTemporality(temporality),
		},
	}
}

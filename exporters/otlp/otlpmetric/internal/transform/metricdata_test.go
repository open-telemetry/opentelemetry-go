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

//go:build go1.18
// +build go1.18

package transform // import "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/internal/transform"

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	cpb "go.opentelemetry.io/proto/otlp/common/v1"
	mpb "go.opentelemetry.io/proto/otlp/metrics/v1"
	rpb "go.opentelemetry.io/proto/otlp/resource/v1"
)

var (
	// Sat Jan 01 2000 00:00:00 GMT+0000.
	start = time.Date(2000, time.January, 01, 0, 0, 0, 0, time.FixedZone("GMT", 0))
	end   = start.Add(30 * time.Second)

	alice = attribute.NewSet(attribute.String("user", "alice"))
	bob   = attribute.NewSet(attribute.String("user", "bob"))

	pbAlice = &cpb.KeyValue{Key: "user", Value: &cpb.AnyValue{
		Value: &cpb.AnyValue_StringValue{StringValue: "alice"},
	}}
	pbBob = &cpb.KeyValue{Key: "user", Value: &cpb.AnyValue{
		Value: &cpb.AnyValue_StringValue{StringValue: "bob"},
	}}

	min, max, sum = 2.0, 4.0, 90.0
	otelHDP       = metricdata.HistogramDataPoint{
		Attributes:   alice,
		StartTime:    start,
		Time:         end,
		Count:        30,
		Bounds:       []float64{1, 5},
		BucketCounts: []uint64{0, 30, 0},
		Min:          &min,
		Max:          &max,
		Sum:          sum,
	}
	otelHist = metricdata.Histogram{
		Temporality: metricdata.DeltaTemporality,
		DataPoints:  []metricdata.HistogramDataPoint{otelHDP},
	}

	pbHDP = &mpb.HistogramDataPoint{
		Attributes:        []*cpb.KeyValue{pbAlice},
		StartTimeUnixNano: uint64(start.UnixNano()),
		TimeUnixNano:      uint64(end.UnixNano()),
		Count:             30,
		Sum:               &sum,
		ExplicitBounds:    []float64{1, 5},
		BucketCounts:      []uint64{0, 30, 0},
		Min:               &min,
		Max:               &max,
	}
	pbHist = &mpb.Histogram{
		AggregationTemporality: mpb.AggregationTemporality_AGGREGATION_TEMPORALITY_DELTA,
		DataPoints:             []*mpb.HistogramDataPoint{pbHDP},
	}

	otelDPtsInt64 = []metricdata.DataPoint[int64]{
		{Attributes: alice, StartTime: start, Time: end, Value: 1},
		{Attributes: bob, StartTime: start, Time: end, Value: 2},
	}
	otelDPtsFloat64 = []metricdata.DataPoint[float64]{
		{Attributes: alice, StartTime: start, Time: end, Value: 1.0},
		{Attributes: bob, StartTime: start, Time: end, Value: 2.0},
	}

	pbDPtsInt64 = []*mpb.NumberDataPoint{
		{
			Attributes:        []*cpb.KeyValue{pbAlice},
			StartTimeUnixNano: uint64(start.UnixNano()),
			TimeUnixNano:      uint64(end.UnixNano()),
			Value:             &mpb.NumberDataPoint_AsInt{AsInt: 1},
		},
		{
			Attributes:        []*cpb.KeyValue{pbBob},
			StartTimeUnixNano: uint64(start.UnixNano()),
			TimeUnixNano:      uint64(end.UnixNano()),
			Value:             &mpb.NumberDataPoint_AsInt{AsInt: 2},
		},
	}
	pbDPtsFloat64 = []*mpb.NumberDataPoint{
		{
			Attributes:        []*cpb.KeyValue{pbAlice},
			StartTimeUnixNano: uint64(start.UnixNano()),
			TimeUnixNano:      uint64(end.UnixNano()),
			Value:             &mpb.NumberDataPoint_AsDouble{AsDouble: 1.0},
		},
		{
			Attributes:        []*cpb.KeyValue{pbBob},
			StartTimeUnixNano: uint64(start.UnixNano()),
			TimeUnixNano:      uint64(end.UnixNano()),
			Value:             &mpb.NumberDataPoint_AsDouble{AsDouble: 2.0},
		},
	}

	otelSumInt64 = metricdata.Sum[int64]{
		Temporality: metricdata.CumulativeTemporality,
		IsMonotonic: true,
		DataPoints:  otelDPtsInt64,
	}
	otelSumFloat64 = metricdata.Sum[float64]{
		Temporality: metricdata.DeltaTemporality,
		IsMonotonic: false,
		DataPoints:  otelDPtsFloat64,
	}

	pbSumInt64 = &mpb.Sum{
		AggregationTemporality: mpb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE,
		IsMonotonic:            true,
		DataPoints:             pbDPtsInt64,
	}
	pbSumFloat64 = &mpb.Sum{
		AggregationTemporality: mpb.AggregationTemporality_AGGREGATION_TEMPORALITY_DELTA,
		IsMonotonic:            false,
		DataPoints:             pbDPtsFloat64,
	}

	otelGaugeInt64   = metricdata.Gauge[int64]{DataPoints: otelDPtsInt64}
	otelGaugeFloat64 = metricdata.Gauge[float64]{DataPoints: otelDPtsFloat64}

	pbGaugeInt64   = &mpb.Gauge{DataPoints: pbDPtsInt64}
	pbGaugeFloat64 = &mpb.Gauge{DataPoints: pbDPtsFloat64}

	otelMetrics = []metricdata.Metrics{
		{
			Name:        "int64-gauge",
			Description: "Gauge with int64 values",
			Unit:        unit.Dimensionless,
			Data:        otelGaugeInt64,
		},
		{
			Name:        "float64-gauge",
			Description: "Gauge with float64 values",
			Unit:        unit.Dimensionless,
			Data:        otelGaugeFloat64,
		},
		{
			Name:        "int64-sum",
			Description: "Sum with int64 values",
			Unit:        unit.Dimensionless,
			Data:        otelSumInt64,
		},
		{
			Name:        "float64-sum",
			Description: "Sum with float64 values",
			Unit:        unit.Dimensionless,
			Data:        otelSumFloat64,
		},
		{
			Name:        "histogram",
			Description: "Histogram",
			Unit:        unit.Dimensionless,
			Data:        otelHist,
		},
	}

	pbMetrics = []*mpb.Metric{
		{
			Name:        "int64-gauge",
			Description: "Gauge with int64 values",
			Unit:        string(unit.Dimensionless),
			Data:        &mpb.Metric_Gauge{Gauge: pbGaugeInt64},
		},
		{
			Name:        "float64-gauge",
			Description: "Gauge with float64 values",
			Unit:        string(unit.Dimensionless),
			Data:        &mpb.Metric_Gauge{Gauge: pbGaugeFloat64},
		},
		{
			Name:        "int64-sum",
			Description: "Sum with int64 values",
			Unit:        string(unit.Dimensionless),
			Data:        &mpb.Metric_Sum{Sum: pbSumInt64},
		},
		{
			Name:        "float64-sum",
			Description: "Sum with float64 values",
			Unit:        string(unit.Dimensionless),
			Data:        &mpb.Metric_Sum{Sum: pbSumFloat64},
		},
		{
			Name:        "histogram",
			Description: "Histogram",
			Unit:        string(unit.Dimensionless),
			Data:        &mpb.Metric_Histogram{Histogram: pbHist},
		},
	}

	otelScopeMetrics = []metricdata.ScopeMetrics{{
		Scope: instrumentation.Scope{
			Name:      "test/code/path",
			Version:   "v0.1.0",
			SchemaURL: semconv.SchemaURL,
		},
		Metrics: otelMetrics,
	}}

	pbScopeMetrics = []*mpb.ScopeMetrics{{
		Scope: &cpb.InstrumentationScope{
			Name:    "test/code/path",
			Version: "v0.1.0",
		},
		Metrics:   pbMetrics,
		SchemaUrl: semconv.SchemaURL,
	}}

	otelRes = resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("test server"),
		semconv.ServiceVersionKey.String("v0.1.0"),
	)

	pbRes = &rpb.Resource{
		Attributes: []*cpb.KeyValue{
			{
				Key: "service.name",
				Value: &cpb.AnyValue{
					Value: &cpb.AnyValue_StringValue{StringValue: "test server"},
				},
			},
			{
				Key: "service.version",
				Value: &cpb.AnyValue{
					Value: &cpb.AnyValue_StringValue{StringValue: "v0.1.0"},
				},
			},
		},
	}

	otelResourceMetrics = metricdata.ResourceMetrics{
		Resource:     otelRes,
		ScopeMetrics: otelScopeMetrics,
	}

	pbResourceMetrics = &mpb.ResourceMetrics{
		Resource:     pbRes,
		ScopeMetrics: pbScopeMetrics,
		SchemaUrl:    semconv.SchemaURL,
	}
)

func TestResourceMetricsTransformation(t *testing.T) {
	got, err := ResourceMetrics(otelResourceMetrics)
	assert.NoError(t, err)
	assert.Equal(t, pbResourceMetrics, got)
}

func TestMetricUnknownAggregationError(t *testing.T) {
	// TODO:
}

func TestTemporalityUnknownTemporalityError(t *testing.T) {
	var unknownTemporality metricdata.Temporality
	pbT, err := Temporality(unknownTemporality)
	assert.ErrorIs(t, err, errUnknownTemporality)
	assert.Equal(t, mpb.AggregationTemporality_AGGREGATION_TEMPORALITY_UNSPECIFIED, pbT)
}

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

package transform // import "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/internal/transform"

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	cpb "go.opentelemetry.io/proto/otlp/common/v1"
	mpb "go.opentelemetry.io/proto/otlp/metrics/v1"
	rpb "go.opentelemetry.io/proto/otlp/resource/v1"
)

type unknownAggT struct {
	metricdata.Aggregation
}

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

	minA, maxA, sumA = 2.0, 4.0, 90.0
	minB, maxB, sumB = 4.0, 150.0, 234.0
	otelHDP          = []metricdata.HistogramDataPoint{{
		Attributes:   alice,
		StartTime:    start,
		Time:         end,
		Count:        30,
		Bounds:       []float64{1, 5},
		BucketCounts: []uint64{0, 30, 0},
		Min:          metricdata.NewExtrema(minA),
		Max:          metricdata.NewExtrema(maxA),
		Sum:          sumA,
	}, {
		Attributes:   bob,
		StartTime:    start,
		Time:         end,
		Count:        3,
		Bounds:       []float64{1, 5},
		BucketCounts: []uint64{0, 1, 2},
		Min:          metricdata.NewExtrema(minB),
		Max:          metricdata.NewExtrema(maxB),
		Sum:          sumB,
	}}

	pbHDP = []*mpb.HistogramDataPoint{{
		Attributes:        []*cpb.KeyValue{pbAlice},
		StartTimeUnixNano: uint64(start.UnixNano()),
		TimeUnixNano:      uint64(end.UnixNano()),
		Count:             30,
		Sum:               &sumA,
		ExplicitBounds:    []float64{1, 5},
		BucketCounts:      []uint64{0, 30, 0},
		Min:               &minA,
		Max:               &maxA,
	}, {
		Attributes:        []*cpb.KeyValue{pbBob},
		StartTimeUnixNano: uint64(start.UnixNano()),
		TimeUnixNano:      uint64(end.UnixNano()),
		Count:             3,
		Sum:               &sumB,
		ExplicitBounds:    []float64{1, 5},
		BucketCounts:      []uint64{0, 1, 2},
		Min:               &minB,
		Max:               &maxB,
	}}

	otelHist = metricdata.Histogram{
		Temporality: metricdata.DeltaTemporality,
		DataPoints:  otelHDP,
	}
	invalidTemporality metricdata.Temporality
	otelHistInvalid    = metricdata.Histogram{
		Temporality: invalidTemporality,
		DataPoints:  otelHDP,
	}

	pbHist = &mpb.Histogram{
		AggregationTemporality: mpb.AggregationTemporality_AGGREGATION_TEMPORALITY_DELTA,
		DataPoints:             pbHDP,
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
	otelSumInvalid = metricdata.Sum[float64]{
		Temporality: invalidTemporality,
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

	unknownAgg  unknownAggT
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
			Name:        "invalid-sum",
			Description: "Sum with invalid temporality",
			Unit:        unit.Dimensionless,
			Data:        otelSumInvalid,
		},
		{
			Name:        "histogram",
			Description: "Histogram",
			Unit:        unit.Dimensionless,
			Data:        otelHist,
		},
		{
			Name:        "invalid-histogram",
			Description: "Invalid histogram",
			Unit:        unit.Dimensionless,
			Data:        otelHistInvalid,
		},
		{
			Name:        "unknown",
			Description: "Unknown aggregation",
			Unit:        unit.Dimensionless,
			Data:        unknownAgg,
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

func TestTransformations(t *testing.T) {
	// Run tests from the "bottom-up" of the metricdata data-types and halt
	// when a failure occurs to ensure the clearest failure message (as
	// opposed to the opposite of testing from the top-down which will obscure
	// errors deep inside the structs).

	// DataPoint types.
	assert.Equal(t, pbHDP, HistogramDataPoints(otelHDP))
	assert.Equal(t, pbDPtsInt64, DataPoints[int64](otelDPtsInt64))
	require.Equal(t, pbDPtsFloat64, DataPoints[float64](otelDPtsFloat64))

	// Aggregations.
	h, err := Histogram(otelHist)
	assert.NoError(t, err)
	assert.Equal(t, &mpb.Metric_Histogram{Histogram: pbHist}, h)
	h, err = Histogram(otelHistInvalid)
	assert.ErrorIs(t, err, errUnknownTemporality)
	assert.Nil(t, h)

	s, err := Sum[int64](otelSumInt64)
	assert.NoError(t, err)
	assert.Equal(t, &mpb.Metric_Sum{Sum: pbSumInt64}, s)
	s, err = Sum[float64](otelSumFloat64)
	assert.NoError(t, err)
	assert.Equal(t, &mpb.Metric_Sum{Sum: pbSumFloat64}, s)
	s, err = Sum[float64](otelSumInvalid)
	assert.ErrorIs(t, err, errUnknownTemporality)
	assert.Nil(t, s)

	assert.Equal(t, &mpb.Metric_Gauge{Gauge: pbGaugeInt64}, Gauge[int64](otelGaugeInt64))
	require.Equal(t, &mpb.Metric_Gauge{Gauge: pbGaugeFloat64}, Gauge[float64](otelGaugeFloat64))

	// Metrics.
	m, err := Metrics(otelMetrics)
	assert.ErrorIs(t, err, errUnknownTemporality)
	assert.ErrorIs(t, err, errUnknownAggregation)
	require.Equal(t, pbMetrics, m)

	// Scope Metrics.
	sm, err := ScopeMetrics(otelScopeMetrics)
	assert.ErrorIs(t, err, errUnknownTemporality)
	assert.ErrorIs(t, err, errUnknownAggregation)
	require.Equal(t, pbScopeMetrics, sm)

	// Resource Metrics.
	rm, err := ResourceMetrics(otelResourceMetrics)
	assert.ErrorIs(t, err, errUnknownTemporality)
	assert.ErrorIs(t, err, errUnknownAggregation)
	require.Equal(t, pbResourceMetrics, rm)
}

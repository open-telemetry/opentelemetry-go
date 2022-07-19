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

package metricdatatest // import "go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
)

var (
	attrA = attribute.NewSet(attribute.Bool("A", true))
	attrB = attribute.NewSet(attribute.Bool("B", true))

	float64A = metricdata.Float64(-1.0)
	float64B = metricdata.Float64(2.0)

	int64A = metricdata.Int64(-1)
	int64B = metricdata.Int64(2)

	dataPointsA = metricdata.DataPoint{
		Attributes: attrA,
		StartTime:  time.Now(),
		Time:       time.Now(),
		Value:      int64A,
	}
	dataPointsB = metricdata.DataPoint{
		Attributes: attrB,
		StartTime:  time.Now(),
		Time:       time.Now(),
		Value:      float64B,
	}

	max, min            = 99.0, 3.
	histogramDataPointA = metricdata.HistogramDataPoint{
		Attributes:   attrA,
		StartTime:    time.Now(),
		Time:         time.Now(),
		Count:        2,
		Bounds:       []float64{0, 10},
		BucketCounts: []uint64{1, 1},
		Sum:          2,
	}
	histogramDataPointB = metricdata.HistogramDataPoint{
		Attributes:   attrB,
		StartTime:    time.Now(),
		Time:         time.Now(),
		Count:        3,
		Bounds:       []float64{0, 10, 100},
		BucketCounts: []uint64{1, 1, 1},
		Max:          &max,
		Min:          &min,
		Sum:          3,
	}

	gaugeA = metricdata.Gauge{DataPoints: []metricdata.DataPoint{dataPointsA}}
	gaugeB = metricdata.Gauge{DataPoints: []metricdata.DataPoint{dataPointsB}}

	sumA = metricdata.Sum{
		Temporality: metricdata.CumulativeTemporality,
		IsMonotonic: true,
		DataPoints:  []metricdata.DataPoint{dataPointsA},
	}
	sumB = metricdata.Sum{
		Temporality: metricdata.DeltaTemporality,
		IsMonotonic: false,
		DataPoints:  []metricdata.DataPoint{dataPointsB},
	}

	histogramA = metricdata.Histogram{
		Temporality: metricdata.CumulativeTemporality,
		DataPoints:  []metricdata.HistogramDataPoint{histogramDataPointA},
	}
	histogramB = metricdata.Histogram{
		Temporality: metricdata.DeltaTemporality,
		DataPoints:  []metricdata.HistogramDataPoint{histogramDataPointB},
	}

	metricsA = metricdata.Metrics{
		Name:        "A",
		Description: "A desc",
		Unit:        unit.Dimensionless,
		Data:        sumA,
	}
	metricsB = metricdata.Metrics{
		Name:        "B",
		Description: "B desc",
		Unit:        unit.Bytes,
		Data:        gaugeB,
	}

	scopeMetricsA = metricdata.ScopeMetrics{
		Scope:   instrumentation.Scope{Name: "A"},
		Metrics: []metricdata.Metrics{metricsA},
	}
	scopeMetricsB = metricdata.ScopeMetrics{
		Scope:   instrumentation.Scope{Name: "B"},
		Metrics: []metricdata.Metrics{metricsB},
	}

	resourceMetricsA = metricdata.ResourceMetrics{
		Resource:     resource.NewSchemaless(attribute.String("resource", "A")),
		ScopeMetrics: []metricdata.ScopeMetrics{scopeMetricsA},
	}
	resourceMetricsB = metricdata.ResourceMetrics{
		Resource:     resource.NewSchemaless(attribute.String("resource", "B")),
		ScopeMetrics: []metricdata.ScopeMetrics{scopeMetricsB},
	}
)

type equalFunc[T Datatypes] func(T, T) (bool, []string)

func testDatatype[T Datatypes](a, b T, f equalFunc[T], reasN int) func(*testing.T) {
	return func(t *testing.T) {
		AssertEqual(t, a, a)
		AssertEqual(t, b, b)

		e, r := f(a, b)
		assert.Falsef(t, e, "%v != %v", a, b)
		assert.Len(t, r, reasN, "number or reasons not equal")
	}
}

func TestAssertEqual(t *testing.T) {
	t.Run("ResourceMetrics", testDatatype(resourceMetricsA, resourceMetricsB, equalResourceMetrics, 2))
	t.Run("ScopeMetrics", testDatatype(scopeMetricsA, scopeMetricsB, equalScopeMetrics, 2))
	t.Run("Metrics", testDatatype(metricsA, metricsB, equalMetrics, 5))
	t.Run("Histogram", testDatatype(histogramA, histogramB, equalHistograms, 2))
	t.Run("Sum", testDatatype(sumA, sumB, equalSums, 3))
	t.Run("Gauge", testDatatype(gaugeA, gaugeB, equalGauges, 1))
	t.Run("HistogramDataPoint", testDatatype(histogramDataPointA, histogramDataPointB, equalHistogramDataPoints, 9))
	t.Run("DataPoint", testDatatype(dataPointsA, dataPointsB, equalDataPoints, 5))
	t.Run("Int64", testDatatype(int64A, int64B, equalInt64, 1))
	t.Run("Float64", testDatatype(float64A, float64B, equalFloat64, 1))
}

type unknownAggregation struct {
	metricdata.Aggregation
}

func TestAssertAggregationsEqual(t *testing.T) {
	AssertAggregationsEqual(t, nil, nil)
	AssertAggregationsEqual(t, sumA, sumA)
	AssertAggregationsEqual(t, gaugeA, gaugeA)
	AssertAggregationsEqual(t, histogramA, histogramA)

	e, r := equalAggregations(sumA, nil)
	assert.False(t, e, "nil comparison")
	assert.Len(t, r, 1, "should return nil comparison mismatch only")

	e, r = equalAggregations(sumA, gaugeA)
	assert.Falsef(t, e, "%v != %v", sumA, gaugeA)
	assert.Len(t, r, 1, "should return with type mismatch only")

	e, r = equalAggregations(unknownAggregation{}, unknownAggregation{})
	assert.False(t, e, "unknown aggregation")
	assert.Len(t, r, 1, "should return with unknown aggregation only")

	e, _ = equalAggregations(sumA, sumB)
	assert.Falsef(t, e, "%v != %v", sumA, sumB)

	e, _ = equalAggregations(gaugeA, gaugeB)
	assert.Falsef(t, e, "%v != %v", gaugeA, gaugeB)

	e, _ = equalAggregations(histogramA, histogramB)
	assert.Falsef(t, e, "%v != %v", histogramA, histogramB)
}

type unknownValue struct {
	metricdata.Value
}

func TestAssertValuesEqual(t *testing.T) {
	AssertValuesEqual(t, nil, nil)
	AssertValuesEqual(t, int64A, int64A)
	AssertValuesEqual(t, float64A, float64A)

	e, r := equalValues(int64A, nil)
	assert.False(t, e, "nil comparison")
	assert.Len(t, r, 1, "should return nil comparison mismatch only")

	e, r = equalValues(int64A, float64A)
	assert.Falsef(t, e, "%v != %v", sumA, gaugeA)
	assert.Len(t, r, 1, "should return with type mismatch only")

	e, r = equalValues(unknownValue{}, unknownValue{})
	assert.False(t, e, "unknown value")
	assert.Len(t, r, 1, "should return with unknown value only")

	e, _ = equalValues(int64A, int64B)
	assert.Falsef(t, e, "%v != %v", int64A, int64B)

	e, _ = equalValues(float64A, float64B)
	assert.Falsef(t, e, "%v != %v", float64A, float64B)
}

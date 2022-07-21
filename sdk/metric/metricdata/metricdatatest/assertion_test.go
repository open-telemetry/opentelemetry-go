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

	startA = time.Now()
	startB = startA.Add(time.Millisecond)
	endA   = startA.Add(time.Second)
	endB   = startB.Add(time.Second)

	dataPointInt64A = metricdata.DataPoint[int64]{
		Attributes: attrA,
		StartTime:  startA,
		Time:       endA,
		Value:      -1,
	}
	dataPointFloat64A = metricdata.DataPoint[float64]{
		Attributes: attrA,
		StartTime:  startA,
		Time:       endA,
		Value:      -1.0,
	}
	dataPointInt64B = metricdata.DataPoint[int64]{
		Attributes: attrB,
		StartTime:  startB,
		Time:       endB,
		Value:      2,
	}
	dataPointFloat64B = metricdata.DataPoint[float64]{
		Attributes: attrB,
		StartTime:  startB,
		Time:       endB,
		Value:      2.0,
	}

	max, min            = 99.0, 3.
	histogramDataPointA = metricdata.HistogramDataPoint{
		Attributes:   attrA,
		StartTime:    startA,
		Time:         endA,
		Count:        2,
		Bounds:       []float64{0, 10},
		BucketCounts: []uint64{1, 1},
		Sum:          2,
	}
	histogramDataPointB = metricdata.HistogramDataPoint{
		Attributes:   attrB,
		StartTime:    startB,
		Time:         endB,
		Count:        3,
		Bounds:       []float64{0, 10, 100},
		BucketCounts: []uint64{1, 1, 1},
		Max:          &max,
		Min:          &min,
		Sum:          3,
	}

	gaugeInt64A = metricdata.Gauge[int64]{
		DataPoints: []metricdata.DataPoint[int64]{dataPointInt64A},
	}
	gaugeFloat64A = metricdata.Gauge[float64]{
		DataPoints: []metricdata.DataPoint[float64]{dataPointFloat64A},
	}
	gaugeInt64B = metricdata.Gauge[int64]{
		DataPoints: []metricdata.DataPoint[int64]{dataPointInt64B},
	}
	gaugeFloat64B = metricdata.Gauge[float64]{
		DataPoints: []metricdata.DataPoint[float64]{dataPointFloat64B},
	}

	sumInt64A = metricdata.Sum[int64]{
		Temporality: metricdata.CumulativeTemporality,
		IsMonotonic: true,
		DataPoints:  []metricdata.DataPoint[int64]{dataPointInt64A},
	}
	sumFloat64A = metricdata.Sum[float64]{
		Temporality: metricdata.CumulativeTemporality,
		IsMonotonic: true,
		DataPoints:  []metricdata.DataPoint[float64]{dataPointFloat64A},
	}
	sumInt64B = metricdata.Sum[int64]{
		Temporality: metricdata.CumulativeTemporality,
		IsMonotonic: true,
		DataPoints:  []metricdata.DataPoint[int64]{dataPointInt64B},
	}
	sumFloat64B = metricdata.Sum[float64]{
		Temporality: metricdata.CumulativeTemporality,
		IsMonotonic: true,
		DataPoints:  []metricdata.DataPoint[float64]{dataPointFloat64B},
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
		Data:        sumInt64A,
	}
	metricsB = metricdata.Metrics{
		Name:        "B",
		Description: "B desc",
		Unit:        unit.Bytes,
		Data:        gaugeFloat64B,
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

type equalFunc[T Datatypes] func(T, T) []string

func testDatatype[T Datatypes](a, b T, f equalFunc[T]) func(*testing.T) {
	return func(t *testing.T) {
		AssertEqual(t, a, a)
		AssertEqual(t, b, b)

		r := f(a, b)
		assert.Greaterf(t, len(r), 0, "%v == %v", a, b)
	}
}

func TestAssertEqual(t *testing.T) {
	t.Run("ResourceMetrics", testDatatype(resourceMetricsA, resourceMetricsB, equalResourceMetrics))
	t.Run("ScopeMetrics", testDatatype(scopeMetricsA, scopeMetricsB, equalScopeMetrics))
	t.Run("Metrics", testDatatype(metricsA, metricsB, equalMetrics))
	t.Run("Histogram", testDatatype(histogramA, histogramB, equalHistograms))
	t.Run("SumInt64", testDatatype(sumInt64A, sumInt64B, equalSums[int64]))
	t.Run("SumFloat64", testDatatype(sumFloat64A, sumFloat64B, equalSums[float64]))
	t.Run("GaugeInt64", testDatatype(gaugeInt64A, gaugeInt64B, equalGauges[int64]))
	t.Run("GaugeFloat64", testDatatype(gaugeFloat64A, gaugeFloat64B, equalGauges[float64]))
	t.Run("HistogramDataPoint", testDatatype(histogramDataPointA, histogramDataPointB, equalHistogramDataPoints))
	t.Run("DataPointInt64", testDatatype(dataPointInt64A, dataPointInt64B, equalDataPoints[int64]))
	t.Run("DataPointFloat64", testDatatype(dataPointFloat64A, dataPointFloat64B, equalDataPoints[float64]))
}

type unknownAggregation struct {
	metricdata.Aggregation
}

func TestAssertAggregationsEqual(t *testing.T) {
	AssertAggregationsEqual(t, nil, nil)
	AssertAggregationsEqual(t, sumInt64A, sumInt64A)
	AssertAggregationsEqual(t, sumFloat64A, sumFloat64A)
	AssertAggregationsEqual(t, gaugeInt64A, gaugeInt64A)
	AssertAggregationsEqual(t, gaugeFloat64A, gaugeFloat64A)
	AssertAggregationsEqual(t, histogramA, histogramA)

	r := equalAggregations(sumInt64A, nil)
	assert.Len(t, r, 1, "should return nil comparison mismatch only")

	r = equalAggregations(sumInt64A, gaugeInt64A)
	assert.Len(t, r, 1, "should return with type mismatch only")

	r = equalAggregations(unknownAggregation{}, unknownAggregation{})
	assert.Len(t, r, 1, "should return with unknown aggregation only")

	r = equalAggregations(sumInt64A, sumInt64B)
	assert.Greaterf(t, len(r), 0, "%v == %v", sumInt64A, sumInt64B)

	r = equalAggregations(sumFloat64A, sumFloat64B)
	assert.Greaterf(t, len(r), 0, "%v == %v", sumFloat64A, sumFloat64B)

	r = equalAggregations(gaugeInt64A, gaugeInt64B)
	assert.Greaterf(t, len(r), 0, "%v == %v", gaugeInt64A, gaugeInt64B)

	r = equalAggregations(gaugeFloat64A, gaugeFloat64B)
	assert.Greaterf(t, len(r), 0, "%v == %v", gaugeFloat64A, gaugeFloat64B)

	r = equalAggregations(histogramA, histogramB)
	assert.Greaterf(t, len(r), 0, "%v == %v", histogramA, histogramB)
}

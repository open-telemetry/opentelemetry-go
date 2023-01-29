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
	dataPointInt64C = metricdata.DataPoint[int64]{
		Attributes: attrA,
		StartTime:  startB,
		Time:       endB,
		Value:      -1,
	}
	dataPointFloat64C = metricdata.DataPoint[float64]{
		Attributes: attrA,
		StartTime:  startB,
		Time:       endB,
		Value:      -1.0,
	}

	minA       = metricdata.NewExtrema(-1.)
	minB, maxB = metricdata.NewExtrema(3.), metricdata.NewExtrema(99.)
	minC       = metricdata.NewExtrema(-1.)

	histogramDataPointA = metricdata.HistogramDataPoint{
		Attributes:   attrA,
		StartTime:    startA,
		Time:         endA,
		Count:        2,
		Bounds:       []float64{0, 10},
		BucketCounts: []uint64{1, 1},
		Min:          minA,
		Sum:          2,
	}
	histogramDataPointB = metricdata.HistogramDataPoint{
		Attributes:   attrB,
		StartTime:    startB,
		Time:         endB,
		Count:        3,
		Bounds:       []float64{0, 10, 100},
		BucketCounts: []uint64{1, 1, 1},
		Max:          maxB,
		Min:          minB,
		Sum:          3,
	}
	histogramDataPointC = metricdata.HistogramDataPoint{
		Attributes:   attrA,
		StartTime:    startB,
		Time:         endB,
		Count:        2,
		Bounds:       []float64{0, 10},
		BucketCounts: []uint64{1, 1},
		Min:          minC,
		Sum:          2,
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
	gaugeInt64C = metricdata.Gauge[int64]{
		DataPoints: []metricdata.DataPoint[int64]{dataPointInt64C},
	}
	gaugeFloat64C = metricdata.Gauge[float64]{
		DataPoints: []metricdata.DataPoint[float64]{dataPointFloat64C},
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
	sumInt64C = metricdata.Sum[int64]{
		Temporality: metricdata.CumulativeTemporality,
		IsMonotonic: true,
		DataPoints:  []metricdata.DataPoint[int64]{dataPointInt64C},
	}
	sumFloat64C = metricdata.Sum[float64]{
		Temporality: metricdata.CumulativeTemporality,
		IsMonotonic: true,
		DataPoints:  []metricdata.DataPoint[float64]{dataPointFloat64C},
	}

	histogramA = metricdata.Histogram{
		Temporality: metricdata.CumulativeTemporality,
		DataPoints:  []metricdata.HistogramDataPoint{histogramDataPointA},
	}
	histogramB = metricdata.Histogram{
		Temporality: metricdata.DeltaTemporality,
		DataPoints:  []metricdata.HistogramDataPoint{histogramDataPointB},
	}
	histogramC = metricdata.Histogram{
		Temporality: metricdata.CumulativeTemporality,
		DataPoints:  []metricdata.HistogramDataPoint{histogramDataPointC},
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
	metricsC = metricdata.Metrics{
		Name:        "A",
		Description: "A desc",
		Unit:        unit.Dimensionless,
		Data:        sumInt64C,
	}

	scopeMetricsA = metricdata.ScopeMetrics{
		Scope:   instrumentation.Scope{Name: "A"},
		Metrics: []metricdata.Metrics{metricsA},
	}
	scopeMetricsB = metricdata.ScopeMetrics{
		Scope:   instrumentation.Scope{Name: "B"},
		Metrics: []metricdata.Metrics{metricsB},
	}
	scopeMetricsC = metricdata.ScopeMetrics{
		Scope:   instrumentation.Scope{Name: "A"},
		Metrics: []metricdata.Metrics{metricsC},
	}

	resourceMetricsA = metricdata.ResourceMetrics{
		Resource:     resource.NewSchemaless(attribute.String("resource", "A")),
		ScopeMetrics: []metricdata.ScopeMetrics{scopeMetricsA},
	}
	resourceMetricsB = metricdata.ResourceMetrics{
		Resource:     resource.NewSchemaless(attribute.String("resource", "B")),
		ScopeMetrics: []metricdata.ScopeMetrics{scopeMetricsB},
	}
	resourceMetricsC = metricdata.ResourceMetrics{
		Resource:     resource.NewSchemaless(attribute.String("resource", "A")),
		ScopeMetrics: []metricdata.ScopeMetrics{scopeMetricsC},
	}
)

type equalFunc[T Datatypes] func(T, T, config) []string

func testDatatype[T Datatypes](a, b T, f equalFunc[T]) func(*testing.T) {
	return func(t *testing.T) {
		AssertEqual(t, a, a)
		AssertEqual(t, b, b)

		r := f(a, b, config{})
		assert.Greaterf(t, len(r), 0, "%v == %v", a, b)
	}
}

func testDatatypeIgnoreTime[T Datatypes](a, b T, f equalFunc[T]) func(*testing.T) {
	return func(t *testing.T) {
		AssertEqual(t, a, a)
		AssertEqual(t, b, b)

		r := f(a, b, config{ignoreTimestamp: true})
		assert.Equalf(t, len(r), 0, "%v == %v", a, b)
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
	t.Run("Extrema", testDatatype(minA, minB, equalExtrema))
}

func TestAssertEqualIgnoreTime(t *testing.T) {
	t.Run("ResourceMetrics", testDatatypeIgnoreTime(resourceMetricsA, resourceMetricsC, equalResourceMetrics))
	t.Run("ScopeMetrics", testDatatypeIgnoreTime(scopeMetricsA, scopeMetricsC, equalScopeMetrics))
	t.Run("Metrics", testDatatypeIgnoreTime(metricsA, metricsC, equalMetrics))
	t.Run("Histogram", testDatatypeIgnoreTime(histogramA, histogramC, equalHistograms))
	t.Run("SumInt64", testDatatypeIgnoreTime(sumInt64A, sumInt64C, equalSums[int64]))
	t.Run("SumFloat64", testDatatypeIgnoreTime(sumFloat64A, sumFloat64C, equalSums[float64]))
	t.Run("GaugeInt64", testDatatypeIgnoreTime(gaugeInt64A, gaugeInt64C, equalGauges[int64]))
	t.Run("GaugeFloat64", testDatatypeIgnoreTime(gaugeFloat64A, gaugeFloat64C, equalGauges[float64]))
	t.Run("HistogramDataPoint", testDatatypeIgnoreTime(histogramDataPointA, histogramDataPointC, equalHistogramDataPoints))
	t.Run("DataPointInt64", testDatatypeIgnoreTime(dataPointInt64A, dataPointInt64C, equalDataPoints[int64]))
	t.Run("DataPointFloat64", testDatatypeIgnoreTime(dataPointFloat64A, dataPointFloat64C, equalDataPoints[float64]))
	t.Run("Extrema", testDatatypeIgnoreTime(minA, minC, equalExtrema))
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

	r := equalAggregations(sumInt64A, nil, config{})
	assert.Len(t, r, 1, "should return nil comparison mismatch only")

	r = equalAggregations(sumInt64A, gaugeInt64A, config{})
	assert.Len(t, r, 1, "should return with type mismatch only")

	r = equalAggregations(unknownAggregation{}, unknownAggregation{}, config{})
	assert.Len(t, r, 1, "should return with unknown aggregation only")

	r = equalAggregations(sumInt64A, sumInt64B, config{})
	assert.Greaterf(t, len(r), 0, "%v == %v", sumInt64A, sumInt64B)

	r = equalAggregations(sumInt64A, sumInt64C, config{ignoreTimestamp: true})
	assert.Equalf(t, len(r), 0, "%v == %v", sumInt64A, sumInt64C)

	r = equalAggregations(sumFloat64A, sumFloat64B, config{})
	assert.Greaterf(t, len(r), 0, "%v == %v", sumFloat64A, sumFloat64B)

	r = equalAggregations(sumFloat64A, sumFloat64C, config{ignoreTimestamp: true})
	assert.Equalf(t, len(r), 0, "%v == %v", sumFloat64A, sumFloat64C)

	r = equalAggregations(gaugeInt64A, gaugeInt64B, config{})
	assert.Greaterf(t, len(r), 0, "%v == %v", gaugeInt64A, gaugeInt64B)

	r = equalAggregations(gaugeInt64A, gaugeInt64C, config{ignoreTimestamp: true})
	assert.Equalf(t, len(r), 0, "%v == %v", gaugeInt64A, gaugeInt64C)

	r = equalAggregations(gaugeFloat64A, gaugeFloat64B, config{})
	assert.Greaterf(t, len(r), 0, "%v == %v", gaugeFloat64A, gaugeFloat64B)

	r = equalAggregations(gaugeFloat64A, gaugeFloat64C, config{ignoreTimestamp: true})
	assert.Equalf(t, len(r), 0, "%v == %v", gaugeFloat64A, gaugeFloat64C)

	r = equalAggregations(histogramA, histogramB, config{})
	assert.Greaterf(t, len(r), 0, "%v == %v", histogramA, histogramB)

	r = equalAggregations(histogramA, histogramC, config{ignoreTimestamp: true})
	assert.Equalf(t, len(r), 0, "%v == %v", histogramA, histogramC)
}

func TestAssertAttributes(t *testing.T) {
	AssertHasAttributes(t, minA, attribute.Bool("A", true)) // No-op, always pass.
	AssertHasAttributes(t, dataPointInt64A, attribute.Bool("A", true))
	AssertHasAttributes(t, dataPointFloat64A, attribute.Bool("A", true))
	AssertHasAttributes(t, gaugeInt64A, attribute.Bool("A", true))
	AssertHasAttributes(t, gaugeFloat64A, attribute.Bool("A", true))
	AssertHasAttributes(t, sumInt64A, attribute.Bool("A", true))
	AssertHasAttributes(t, sumFloat64A, attribute.Bool("A", true))
	AssertHasAttributes(t, histogramDataPointA, attribute.Bool("A", true))
	AssertHasAttributes(t, histogramA, attribute.Bool("A", true))
	AssertHasAttributes(t, metricsA, attribute.Bool("A", true))
	AssertHasAttributes(t, scopeMetricsA, attribute.Bool("A", true))
	AssertHasAttributes(t, resourceMetricsA, attribute.Bool("A", true))

	r := hasAttributesAggregation(gaugeInt64A, attribute.Bool("A", true))
	assert.Equal(t, len(r), 0, "gaugeInt64A has A=True")
	r = hasAttributesAggregation(gaugeFloat64A, attribute.Bool("A", true))
	assert.Equal(t, len(r), 0, "gaugeFloat64A has A=True")
	r = hasAttributesAggregation(sumInt64A, attribute.Bool("A", true))
	assert.Equal(t, len(r), 0, "sumInt64A has A=True")
	r = hasAttributesAggregation(sumFloat64A, attribute.Bool("A", true))
	assert.Equal(t, len(r), 0, "sumFloat64A has A=True")
	r = hasAttributesAggregation(histogramA, attribute.Bool("A", true))
	assert.Equal(t, len(r), 0, "histogramA has A=True")

	r = hasAttributesAggregation(gaugeInt64A, attribute.Bool("A", false))
	assert.Greater(t, len(r), 0, "gaugeInt64A does not have A=False")
	r = hasAttributesAggregation(gaugeFloat64A, attribute.Bool("A", false))
	assert.Greater(t, len(r), 0, "gaugeFloat64A does not have A=False")
	r = hasAttributesAggregation(sumInt64A, attribute.Bool("A", false))
	assert.Greater(t, len(r), 0, "sumInt64A does not have A=False")
	r = hasAttributesAggregation(sumFloat64A, attribute.Bool("A", false))
	assert.Greater(t, len(r), 0, "sumFloat64A does not have A=False")
	r = hasAttributesAggregation(histogramA, attribute.Bool("A", false))
	assert.Greater(t, len(r), 0, "histogramA does not have A=False")

	r = hasAttributesAggregation(gaugeInt64A, attribute.Bool("B", true))
	assert.Greater(t, len(r), 0, "gaugeInt64A does not have Attribute B")
	r = hasAttributesAggregation(gaugeFloat64A, attribute.Bool("B", true))
	assert.Greater(t, len(r), 0, "gaugeFloat64A does not have Attribute B")
	r = hasAttributesAggregation(sumInt64A, attribute.Bool("B", true))
	assert.Greater(t, len(r), 0, "sumInt64A does not have Attribute B")
	r = hasAttributesAggregation(sumFloat64A, attribute.Bool("B", true))
	assert.Greater(t, len(r), 0, "sumFloat64A does not have Attribute B")
	r = hasAttributesAggregation(histogramA, attribute.Bool("B", true))
	assert.Greater(t, len(r), 0, "histogramA does not have Attribute B")
}

func TestAssertAttributesFail(t *testing.T) {
	fakeT := &testing.T{}
	assert.False(t, AssertHasAttributes(fakeT, dataPointInt64A, attribute.Bool("A", false)))
	assert.False(t, AssertHasAttributes(fakeT, dataPointFloat64A, attribute.Bool("B", true)))
	assert.False(t, AssertHasAttributes(fakeT, gaugeInt64A, attribute.Bool("A", false)))
	assert.False(t, AssertHasAttributes(fakeT, gaugeFloat64A, attribute.Bool("B", true)))
	assert.False(t, AssertHasAttributes(fakeT, sumInt64A, attribute.Bool("A", false)))
	assert.False(t, AssertHasAttributes(fakeT, sumFloat64A, attribute.Bool("B", true)))
	assert.False(t, AssertHasAttributes(fakeT, histogramDataPointA, attribute.Bool("A", false)))
	assert.False(t, AssertHasAttributes(fakeT, histogramDataPointA, attribute.Bool("B", true)))
	assert.False(t, AssertHasAttributes(fakeT, histogramA, attribute.Bool("A", false)))
	assert.False(t, AssertHasAttributes(fakeT, histogramA, attribute.Bool("B", true)))
	assert.False(t, AssertHasAttributes(fakeT, metricsA, attribute.Bool("A", false)))
	assert.False(t, AssertHasAttributes(fakeT, metricsA, attribute.Bool("B", true)))
	assert.False(t, AssertHasAttributes(fakeT, resourceMetricsA, attribute.Bool("A", false)))
	assert.False(t, AssertHasAttributes(fakeT, resourceMetricsA, attribute.Bool("B", true)))

	sum := metricdata.Sum[int64]{
		Temporality: metricdata.CumulativeTemporality,
		IsMonotonic: true,
		DataPoints: []metricdata.DataPoint[int64]{
			dataPointInt64A,
			dataPointInt64B,
		},
	}
	assert.False(t, AssertHasAttributes(fakeT, sum, attribute.Bool("A", true)))
}

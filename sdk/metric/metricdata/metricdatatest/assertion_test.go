// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metricdatatest // import "go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
)

var (
	attrA = attribute.NewSet(attribute.Bool("A", true))
	attrB = attribute.NewSet(attribute.Bool("B", true))

	fltrAttrA = []attribute.KeyValue{attribute.Bool("filter A", true)}
	fltrAttrB = []attribute.KeyValue{attribute.Bool("filter B", true)}

	startA = time.Now()
	startB = startA.Add(time.Millisecond)
	endA   = startA.Add(time.Second)
	endB   = startB.Add(time.Second)

	spanIDA  = []byte{0, 0, 0, 0, 0, 0, 0, 1}
	spanIDB  = []byte{0, 0, 0, 0, 0, 0, 0, 2}
	traceIDA = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	traceIDB = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}

	exemplarInt64A = metricdata.Exemplar[int64]{
		FilteredAttributes: fltrAttrA,
		Time:               endA,
		Value:              -10,
		SpanID:             spanIDA,
		TraceID:            traceIDA,
	}
	exemplarFloat64A = metricdata.Exemplar[float64]{
		FilteredAttributes: fltrAttrA,
		Time:               endA,
		Value:              -10.0,
		SpanID:             spanIDA,
		TraceID:            traceIDA,
	}
	exemplarInt64B = metricdata.Exemplar[int64]{
		FilteredAttributes: fltrAttrB,
		Time:               endB,
		Value:              12,
		SpanID:             spanIDB,
		TraceID:            traceIDB,
	}
	exemplarFloat64B = metricdata.Exemplar[float64]{
		FilteredAttributes: fltrAttrB,
		Time:               endB,
		Value:              12.0,
		SpanID:             spanIDB,
		TraceID:            traceIDB,
	}
	exemplarInt64C = metricdata.Exemplar[int64]{
		FilteredAttributes: fltrAttrA,
		Time:               endB,
		Value:              -10,
		SpanID:             spanIDA,
		TraceID:            traceIDA,
	}
	exemplarFloat64C = metricdata.Exemplar[float64]{
		FilteredAttributes: fltrAttrA,
		Time:               endB,
		Value:              -10.0,
		SpanID:             spanIDA,
		TraceID:            traceIDA,
	}
	exemplarInt64D = metricdata.Exemplar[int64]{
		FilteredAttributes: fltrAttrA,
		Time:               endA,
		Value:              12,
		SpanID:             spanIDA,
		TraceID:            traceIDA,
	}
	exemplarFloat64D = metricdata.Exemplar[float64]{
		FilteredAttributes: fltrAttrA,
		Time:               endA,
		Value:              12.0,
		SpanID:             spanIDA,
		TraceID:            traceIDA,
	}

	dataPointInt64A = metricdata.DataPoint[int64]{
		Attributes: attrA,
		StartTime:  startA,
		Time:       endA,
		Value:      -1,
		Exemplars:  []metricdata.Exemplar[int64]{exemplarInt64A},
	}
	dataPointFloat64A = metricdata.DataPoint[float64]{
		Attributes: attrA,
		StartTime:  startA,
		Time:       endA,
		Value:      -1.0,
		Exemplars:  []metricdata.Exemplar[float64]{exemplarFloat64A},
	}
	dataPointInt64B = metricdata.DataPoint[int64]{
		Attributes: attrB,
		StartTime:  startB,
		Time:       endB,
		Value:      2,
		Exemplars:  []metricdata.Exemplar[int64]{exemplarInt64B},
	}
	dataPointFloat64B = metricdata.DataPoint[float64]{
		Attributes: attrB,
		StartTime:  startB,
		Time:       endB,
		Value:      2.0,
		Exemplars:  []metricdata.Exemplar[float64]{exemplarFloat64B},
	}
	dataPointInt64C = metricdata.DataPoint[int64]{
		Attributes: attrA,
		StartTime:  startB,
		Time:       endB,
		Value:      -1,
		Exemplars:  []metricdata.Exemplar[int64]{exemplarInt64C},
	}
	dataPointFloat64C = metricdata.DataPoint[float64]{
		Attributes: attrA,
		StartTime:  startB,
		Time:       endB,
		Value:      -1.0,
		Exemplars:  []metricdata.Exemplar[float64]{exemplarFloat64C},
	}
	dataPointInt64D = metricdata.DataPoint[int64]{
		Attributes: attrA,
		StartTime:  startA,
		Time:       endA,
		Value:      2,
		Exemplars:  []metricdata.Exemplar[int64]{exemplarInt64A},
	}
	dataPointFloat64D = metricdata.DataPoint[float64]{
		Attributes: attrA,
		StartTime:  startA,
		Time:       endA,
		Value:      2.0,
		Exemplars:  []metricdata.Exemplar[float64]{exemplarFloat64A},
	}

	minFloat64A              = metricdata.NewExtrema(-1.)
	minInt64A                = metricdata.NewExtrema[int64](-1)
	minFloat64B, maxFloat64B = metricdata.NewExtrema(3.), metricdata.NewExtrema(99.)
	minInt64B, maxInt64B     = metricdata.NewExtrema[int64](3), metricdata.NewExtrema[int64](99)
	minFloat64C              = metricdata.NewExtrema(-1.)
	minInt64C                = metricdata.NewExtrema[int64](-1)

	minFloat64D = metricdata.NewExtrema(-9.999999)

	histogramDataPointInt64A = metricdata.HistogramDataPoint[int64]{
		Attributes:   attrA,
		StartTime:    startA,
		Time:         endA,
		Count:        2,
		Bounds:       []float64{0, 10},
		BucketCounts: []uint64{1, 1},
		Min:          minInt64A,
		Sum:          2,
		Exemplars:    []metricdata.Exemplar[int64]{exemplarInt64A},
	}
	histogramDataPointFloat64A = metricdata.HistogramDataPoint[float64]{
		Attributes:   attrA,
		StartTime:    startA,
		Time:         endA,
		Count:        2,
		Bounds:       []float64{0, 10},
		BucketCounts: []uint64{1, 1},
		Min:          minFloat64A,
		Sum:          2,
		Exemplars:    []metricdata.Exemplar[float64]{exemplarFloat64A},
	}
	histogramDataPointInt64B = metricdata.HistogramDataPoint[int64]{
		Attributes:   attrB,
		StartTime:    startB,
		Time:         endB,
		Count:        3,
		Bounds:       []float64{0, 10, 100},
		BucketCounts: []uint64{1, 1, 1},
		Max:          maxInt64B,
		Min:          minInt64B,
		Sum:          3,
		Exemplars:    []metricdata.Exemplar[int64]{exemplarInt64B},
	}
	histogramDataPointFloat64B = metricdata.HistogramDataPoint[float64]{
		Attributes:   attrB,
		StartTime:    startB,
		Time:         endB,
		Count:        3,
		Bounds:       []float64{0, 10, 100},
		BucketCounts: []uint64{1, 1, 1},
		Max:          maxFloat64B,
		Min:          minFloat64B,
		Sum:          3,
		Exemplars:    []metricdata.Exemplar[float64]{exemplarFloat64B},
	}
	histogramDataPointInt64C = metricdata.HistogramDataPoint[int64]{
		Attributes:   attrA,
		StartTime:    startB,
		Time:         endB,
		Count:        2,
		Bounds:       []float64{0, 10},
		BucketCounts: []uint64{1, 1},
		Min:          minInt64C,
		Sum:          2,
		Exemplars:    []metricdata.Exemplar[int64]{exemplarInt64C},
	}
	histogramDataPointFloat64C = metricdata.HistogramDataPoint[float64]{
		Attributes:   attrA,
		StartTime:    startB,
		Time:         endB,
		Count:        2,
		Bounds:       []float64{0, 10},
		BucketCounts: []uint64{1, 1},
		Min:          minFloat64C,
		Sum:          2,
		Exemplars:    []metricdata.Exemplar[float64]{exemplarFloat64C},
	}
	histogramDataPointInt64D = metricdata.HistogramDataPoint[int64]{
		Attributes:   attrA,
		StartTime:    startA,
		Time:         endA,
		Count:        3,
		Bounds:       []float64{0, 10, 100},
		BucketCounts: []uint64{1, 1, 1},
		Max:          maxInt64B,
		Min:          minInt64B,
		Sum:          3,
		Exemplars:    []metricdata.Exemplar[int64]{exemplarInt64A},
	}
	histogramDataPointFloat64D = metricdata.HistogramDataPoint[float64]{
		Attributes:   attrA,
		StartTime:    startA,
		Time:         endA,
		Count:        3,
		Bounds:       []float64{0, 10, 100},
		BucketCounts: []uint64{1, 1, 1},
		Max:          maxFloat64B,
		Min:          minFloat64B,
		Sum:          3,
		Exemplars:    []metricdata.Exemplar[float64]{exemplarFloat64A},
	}

	quantileValueA = metricdata.QuantileValue{
		Quantile: 0.0,
		Value:    0.1,
	}
	quantileValueB = metricdata.QuantileValue{
		Quantile: 0.1,
		Value:    0.2,
	}
	summaryDataPointA = metricdata.SummaryDataPoint{
		Attributes:     attrA,
		StartTime:      startA,
		Time:           endA,
		Count:          2,
		Sum:            3,
		QuantileValues: []metricdata.QuantileValue{quantileValueA},
	}
	summaryDataPointB = metricdata.SummaryDataPoint{
		Attributes:     attrB,
		StartTime:      startB,
		Time:           endB,
		Count:          3,
		QuantileValues: []metricdata.QuantileValue{quantileValueB},
	}
	summaryDataPointC = metricdata.SummaryDataPoint{
		Attributes:     attrA,
		StartTime:      startB,
		Time:           endB,
		Count:          2,
		Sum:            3,
		QuantileValues: []metricdata.QuantileValue{quantileValueA},
	}
	summaryDataPointD = metricdata.SummaryDataPoint{
		Attributes:     attrA,
		StartTime:      startA,
		Time:           endA,
		Count:          3,
		QuantileValues: []metricdata.QuantileValue{quantileValueB},
	}

	exponentialBucket2 = metricdata.ExponentialBucket{
		Offset: 2,
		Counts: []uint64{1, 1},
	}
	exponentialBucket3 = metricdata.ExponentialBucket{
		Offset: 3,
		Counts: []uint64{1, 1},
	}
	exponentialBucket4 = metricdata.ExponentialBucket{
		Offset: 4,
		Counts: []uint64{1, 1, 1},
	}
	exponentialBucket5 = metricdata.ExponentialBucket{
		Offset: 5,
		Counts: []uint64{1, 1, 1},
	}
	exponentialHistogramDataPointInt64A = metricdata.ExponentialHistogramDataPoint[int64]{
		Attributes:     attrA,
		StartTime:      startA,
		Time:           endA,
		Count:          5,
		Min:            minInt64A,
		Sum:            2,
		Scale:          1,
		ZeroCount:      1,
		PositiveBucket: exponentialBucket3,
		NegativeBucket: exponentialBucket2,
		Exemplars:      []metricdata.Exemplar[int64]{exemplarInt64A},
	}
	exponentialHistogramDataPointFloat64A = metricdata.ExponentialHistogramDataPoint[float64]{
		Attributes:     attrA,
		StartTime:      startA,
		Time:           endA,
		Count:          5,
		Min:            minFloat64A,
		Sum:            2,
		Scale:          1,
		ZeroCount:      1,
		PositiveBucket: exponentialBucket3,
		NegativeBucket: exponentialBucket2,
		Exemplars:      []metricdata.Exemplar[float64]{exemplarFloat64A},
	}
	exponentialHistogramDataPointInt64B = metricdata.ExponentialHistogramDataPoint[int64]{
		Attributes:     attrB,
		StartTime:      startB,
		Time:           endB,
		Count:          6,
		Min:            minInt64B,
		Max:            maxInt64B,
		Sum:            3,
		Scale:          2,
		ZeroCount:      3,
		PositiveBucket: exponentialBucket4,
		NegativeBucket: exponentialBucket5,
		Exemplars:      []metricdata.Exemplar[int64]{exemplarInt64B},
	}
	exponentialHistogramDataPointFloat64B = metricdata.ExponentialHistogramDataPoint[float64]{
		Attributes:     attrB,
		StartTime:      startB,
		Time:           endB,
		Count:          6,
		Min:            minFloat64B,
		Max:            maxFloat64B,
		Sum:            3,
		Scale:          2,
		ZeroCount:      3,
		PositiveBucket: exponentialBucket4,
		NegativeBucket: exponentialBucket5,
		Exemplars:      []metricdata.Exemplar[float64]{exemplarFloat64B},
	}
	exponentialHistogramDataPointInt64C = metricdata.ExponentialHistogramDataPoint[int64]{
		Attributes:     attrA,
		StartTime:      startB,
		Time:           endB,
		Count:          5,
		Min:            minInt64C,
		Sum:            2,
		Scale:          1,
		ZeroCount:      1,
		PositiveBucket: exponentialBucket3,
		NegativeBucket: exponentialBucket2,
		Exemplars:      []metricdata.Exemplar[int64]{exemplarInt64C},
	}
	exponentialHistogramDataPointFloat64C = metricdata.ExponentialHistogramDataPoint[float64]{
		Attributes:     attrA,
		StartTime:      startB,
		Time:           endB,
		Count:          5,
		Min:            minFloat64A,
		Sum:            2,
		Scale:          1,
		ZeroCount:      1,
		PositiveBucket: exponentialBucket3,
		NegativeBucket: exponentialBucket2,
		Exemplars:      []metricdata.Exemplar[float64]{exemplarFloat64C},
	}
	exponentialHistogramDataPointInt64D = metricdata.ExponentialHistogramDataPoint[int64]{
		Attributes:     attrA,
		StartTime:      startA,
		Time:           endA,
		Count:          6,
		Min:            minInt64B,
		Max:            maxInt64B,
		Sum:            3,
		Scale:          2,
		ZeroCount:      3,
		PositiveBucket: exponentialBucket4,
		NegativeBucket: exponentialBucket5,
		Exemplars:      []metricdata.Exemplar[int64]{exemplarInt64A},
	}
	exponentialHistogramDataPointFloat64D = metricdata.ExponentialHistogramDataPoint[float64]{
		Attributes:     attrA,
		StartTime:      startA,
		Time:           endA,
		Count:          6,
		Min:            minFloat64B,
		Max:            maxFloat64B,
		Sum:            3,
		Scale:          2,
		ZeroCount:      3,
		PositiveBucket: exponentialBucket4,
		NegativeBucket: exponentialBucket5,
		Exemplars:      []metricdata.Exemplar[float64]{exemplarFloat64A},
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
	gaugeInt64D = metricdata.Gauge[int64]{
		DataPoints: []metricdata.DataPoint[int64]{dataPointInt64D},
	}
	gaugeFloat64D = metricdata.Gauge[float64]{
		DataPoints: []metricdata.DataPoint[float64]{dataPointFloat64D},
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
	sumInt64D = metricdata.Sum[int64]{
		Temporality: metricdata.CumulativeTemporality,
		IsMonotonic: true,
		DataPoints:  []metricdata.DataPoint[int64]{dataPointInt64D},
	}
	sumFloat64D = metricdata.Sum[float64]{
		Temporality: metricdata.CumulativeTemporality,
		IsMonotonic: true,
		DataPoints:  []metricdata.DataPoint[float64]{dataPointFloat64D},
	}

	histogramInt64A = metricdata.Histogram[int64]{
		Temporality: metricdata.CumulativeTemporality,
		DataPoints:  []metricdata.HistogramDataPoint[int64]{histogramDataPointInt64A},
	}
	histogramFloat64A = metricdata.Histogram[float64]{
		Temporality: metricdata.CumulativeTemporality,
		DataPoints:  []metricdata.HistogramDataPoint[float64]{histogramDataPointFloat64A},
	}
	histogramInt64B = metricdata.Histogram[int64]{
		Temporality: metricdata.DeltaTemporality,
		DataPoints:  []metricdata.HistogramDataPoint[int64]{histogramDataPointInt64B},
	}
	histogramFloat64B = metricdata.Histogram[float64]{
		Temporality: metricdata.DeltaTemporality,
		DataPoints:  []metricdata.HistogramDataPoint[float64]{histogramDataPointFloat64B},
	}
	histogramInt64C = metricdata.Histogram[int64]{
		Temporality: metricdata.CumulativeTemporality,
		DataPoints:  []metricdata.HistogramDataPoint[int64]{histogramDataPointInt64C},
	}
	histogramFloat64C = metricdata.Histogram[float64]{
		Temporality: metricdata.CumulativeTemporality,
		DataPoints:  []metricdata.HistogramDataPoint[float64]{histogramDataPointFloat64C},
	}
	histogramInt64D = metricdata.Histogram[int64]{
		Temporality: metricdata.CumulativeTemporality,
		DataPoints:  []metricdata.HistogramDataPoint[int64]{histogramDataPointInt64D},
	}
	histogramFloat64D = metricdata.Histogram[float64]{
		Temporality: metricdata.CumulativeTemporality,
		DataPoints:  []metricdata.HistogramDataPoint[float64]{histogramDataPointFloat64D},
	}

	exponentialHistogramInt64A = metricdata.ExponentialHistogram[int64]{
		Temporality: metricdata.CumulativeTemporality,
		DataPoints:  []metricdata.ExponentialHistogramDataPoint[int64]{exponentialHistogramDataPointInt64A},
	}
	exponentialHistogramFloat64A = metricdata.ExponentialHistogram[float64]{
		Temporality: metricdata.CumulativeTemporality,
		DataPoints:  []metricdata.ExponentialHistogramDataPoint[float64]{exponentialHistogramDataPointFloat64A},
	}
	exponentialHistogramInt64B = metricdata.ExponentialHistogram[int64]{
		Temporality: metricdata.DeltaTemporality,
		DataPoints:  []metricdata.ExponentialHistogramDataPoint[int64]{exponentialHistogramDataPointInt64B},
	}
	exponentialHistogramFloat64B = metricdata.ExponentialHistogram[float64]{
		Temporality: metricdata.DeltaTemporality,
		DataPoints:  []metricdata.ExponentialHistogramDataPoint[float64]{exponentialHistogramDataPointFloat64B},
	}
	exponentialHistogramInt64C = metricdata.ExponentialHistogram[int64]{
		Temporality: metricdata.CumulativeTemporality,
		DataPoints:  []metricdata.ExponentialHistogramDataPoint[int64]{exponentialHistogramDataPointInt64C},
	}
	exponentialHistogramFloat64C = metricdata.ExponentialHistogram[float64]{
		Temporality: metricdata.CumulativeTemporality,
		DataPoints:  []metricdata.ExponentialHistogramDataPoint[float64]{exponentialHistogramDataPointFloat64C},
	}
	exponentialHistogramInt64D = metricdata.ExponentialHistogram[int64]{
		Temporality: metricdata.CumulativeTemporality,
		DataPoints:  []metricdata.ExponentialHistogramDataPoint[int64]{exponentialHistogramDataPointInt64D},
	}
	exponentialHistogramFloat64D = metricdata.ExponentialHistogram[float64]{
		Temporality: metricdata.CumulativeTemporality,
		DataPoints:  []metricdata.ExponentialHistogramDataPoint[float64]{exponentialHistogramDataPointFloat64D},
	}

	summaryA = metricdata.Summary{
		DataPoints: []metricdata.SummaryDataPoint{summaryDataPointA},
	}

	summaryB = metricdata.Summary{
		DataPoints: []metricdata.SummaryDataPoint{summaryDataPointB},
	}

	summaryC = metricdata.Summary{
		DataPoints: []metricdata.SummaryDataPoint{summaryDataPointC},
	}

	summaryD = metricdata.Summary{
		DataPoints: []metricdata.SummaryDataPoint{summaryDataPointD},
	}

	metricsA = metricdata.Metrics{
		Name:        "A",
		Description: "A desc",
		Unit:        "1",
		Data:        sumInt64A,
	}
	metricsB = metricdata.Metrics{
		Name:        "B",
		Description: "B desc",
		Unit:        "By",
		Data:        gaugeFloat64B,
	}
	metricsC = metricdata.Metrics{
		Name:        "A",
		Description: "A desc",
		Unit:        "1",
		Data:        sumInt64C,
	}
	metricsD = metricdata.Metrics{
		Name:        "A",
		Description: "A desc",
		Unit:        "1",
		Data:        sumInt64D,
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
	scopeMetricsD = metricdata.ScopeMetrics{
		Scope:   instrumentation.Scope{Name: "A"},
		Metrics: []metricdata.Metrics{metricsD},
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
	resourceMetricsD = metricdata.ResourceMetrics{
		Resource:     resource.NewSchemaless(attribute.String("resource", "A")),
		ScopeMetrics: []metricdata.ScopeMetrics{scopeMetricsD},
	}
)

type equalFunc[T Datatypes] func(T, T, config) []string

func testDatatype[T Datatypes](a, b T, f equalFunc[T]) func(*testing.T) {
	return func(t *testing.T) {
		AssertEqual(t, a, a)
		AssertEqual(t, b, b)

		r := f(a, b, newConfig(nil))
		assert.NotEmptyf(t, r, "%v == %v", a, b)
	}
}

func testDatatypeIgnoreTime[T Datatypes](a, b T, f equalFunc[T]) func(*testing.T) {
	return func(t *testing.T) {
		AssertEqual(t, a, a)
		AssertEqual(t, b, b)

		c := newConfig([]Option{IgnoreTimestamp()})
		r := f(a, b, c)
		assert.Empty(t, r, "unexpected inequality")
	}
}

func testDatatypeIgnoreExemplars[T Datatypes](a, b T, f equalFunc[T]) func(*testing.T) {
	return func(t *testing.T) {
		AssertEqual(t, a, a)
		AssertEqual(t, b, b)

		c := newConfig([]Option{IgnoreExemplars()})
		r := f(a, b, c)
		assert.Empty(t, r, "unexpected inequality")
	}
}

func testDatatypeIgnoreValue[T Datatypes](a, b T, f equalFunc[T]) func(*testing.T) {
	return func(t *testing.T) {
		AssertEqual(t, a, a)
		AssertEqual(t, b, b)

		c := newConfig([]Option{IgnoreValue()})
		r := f(a, b, c)
		assert.Empty(t, r, "unexpected inequality")
	}
}

func TestTestingTImplementation(t *testing.T) {
	assert.Implements(t, (*TestingT)(nil), t)
}

func TestAssertEqual(t *testing.T) {
	t.Run("ResourceMetrics", testDatatype(resourceMetricsA, resourceMetricsB, equalResourceMetrics))
	t.Run("ScopeMetrics", testDatatype(scopeMetricsA, scopeMetricsB, equalScopeMetrics))
	t.Run("Metrics", testDatatype(metricsA, metricsB, equalMetrics))
	t.Run("HistogramInt64", testDatatype(histogramInt64A, histogramInt64B, equalHistograms[int64]))
	t.Run("HistogramFloat64", testDatatype(histogramFloat64A, histogramFloat64B, equalHistograms[float64]))
	t.Run("SumInt64", testDatatype(sumInt64A, sumInt64B, equalSums[int64]))
	t.Run("SumFloat64", testDatatype(sumFloat64A, sumFloat64B, equalSums[float64]))
	t.Run("GaugeInt64", testDatatype(gaugeInt64A, gaugeInt64B, equalGauges[int64]))
	t.Run("GaugeFloat64", testDatatype(gaugeFloat64A, gaugeFloat64B, equalGauges[float64]))
	t.Run("HistogramDataPointInt64", testDatatype(histogramDataPointInt64A, histogramDataPointInt64B, equalHistogramDataPoints[int64]))
	t.Run("HistogramDataPointFloat64", testDatatype(histogramDataPointFloat64A, histogramDataPointFloat64B, equalHistogramDataPoints[float64]))
	t.Run("DataPointInt64", testDatatype(dataPointInt64A, dataPointInt64B, equalDataPoints[int64]))
	t.Run("DataPointFloat64", testDatatype(dataPointFloat64A, dataPointFloat64B, equalDataPoints[float64]))
	t.Run("ExtremaInt64", testDatatype(minInt64A, minInt64B, equalExtrema[int64]))
	t.Run("ExtremaFloat64", testDatatype(minFloat64A, minFloat64B, equalExtrema[float64]))
	t.Run("ExemplarInt64", testDatatype(exemplarInt64A, exemplarInt64B, equalExemplars[int64]))
	t.Run("ExemplarFloat64", testDatatype(exemplarFloat64A, exemplarFloat64B, equalExemplars[float64]))
	t.Run("ExponentialHistogramInt64", testDatatype(exponentialHistogramInt64A, exponentialHistogramInt64B, equalExponentialHistograms[int64]))
	t.Run("ExponentialHistogramFloat64", testDatatype(exponentialHistogramFloat64A, exponentialHistogramFloat64B, equalExponentialHistograms[float64]))
	t.Run("ExponentialHistogramDataPointInt64", testDatatype(exponentialHistogramDataPointInt64A, exponentialHistogramDataPointInt64B, equalExponentialHistogramDataPoints[int64]))
	t.Run("ExponentialHistogramDataPointFloat64", testDatatype(exponentialHistogramDataPointFloat64A, exponentialHistogramDataPointFloat64B, equalExponentialHistogramDataPoints[float64]))
	t.Run("ExponentialBuckets", testDatatype(exponentialBucket2, exponentialBucket3, equalExponentialBuckets))
	t.Run("Summary", testDatatype(summaryA, summaryB, equalSummary))
	t.Run("SummaryDataPoint", testDatatype(summaryDataPointA, summaryDataPointB, equalSummaryDataPoint))
	t.Run("QuantileValues", testDatatype(quantileValueA, quantileValueB, equalQuantileValue))
}

func TestAssertEqualIgnoreTime(t *testing.T) {
	t.Run("ResourceMetrics", testDatatypeIgnoreTime(resourceMetricsA, resourceMetricsC, equalResourceMetrics))
	t.Run("ScopeMetrics", testDatatypeIgnoreTime(scopeMetricsA, scopeMetricsC, equalScopeMetrics))
	t.Run("Metrics", testDatatypeIgnoreTime(metricsA, metricsC, equalMetrics))
	t.Run("HistogramInt64", testDatatypeIgnoreTime(histogramInt64A, histogramInt64C, equalHistograms[int64]))
	t.Run("HistogramFloat64", testDatatypeIgnoreTime(histogramFloat64A, histogramFloat64C, equalHistograms[float64]))
	t.Run("SumInt64", testDatatypeIgnoreTime(sumInt64A, sumInt64C, equalSums[int64]))
	t.Run("SumFloat64", testDatatypeIgnoreTime(sumFloat64A, sumFloat64C, equalSums[float64]))
	t.Run("GaugeInt64", testDatatypeIgnoreTime(gaugeInt64A, gaugeInt64C, equalGauges[int64]))
	t.Run("GaugeFloat64", testDatatypeIgnoreTime(gaugeFloat64A, gaugeFloat64C, equalGauges[float64]))
	t.Run("HistogramDataPointInt64", testDatatypeIgnoreTime(histogramDataPointInt64A, histogramDataPointInt64C, equalHistogramDataPoints[int64]))
	t.Run("HistogramDataPointFloat64", testDatatypeIgnoreTime(histogramDataPointFloat64A, histogramDataPointFloat64C, equalHistogramDataPoints[float64]))
	t.Run("DataPointInt64", testDatatypeIgnoreTime(dataPointInt64A, dataPointInt64C, equalDataPoints[int64]))
	t.Run("DataPointFloat64", testDatatypeIgnoreTime(dataPointFloat64A, dataPointFloat64C, equalDataPoints[float64]))
	t.Run("ExtremaInt64", testDatatypeIgnoreTime(minInt64A, minInt64C, equalExtrema[int64]))
	t.Run("ExtremaFloat64", testDatatypeIgnoreTime(minFloat64A, minFloat64C, equalExtrema[float64]))
	t.Run("ExemplarInt64", testDatatypeIgnoreTime(exemplarInt64A, exemplarInt64C, equalExemplars[int64]))
	t.Run("ExemplarFloat64", testDatatypeIgnoreTime(exemplarFloat64A, exemplarFloat64C, equalExemplars[float64]))
	t.Run("ExponentialHistogramInt64", testDatatypeIgnoreTime(exponentialHistogramInt64A, exponentialHistogramInt64C, equalExponentialHistograms[int64]))
	t.Run("ExponentialHistogramFloat64", testDatatypeIgnoreTime(exponentialHistogramFloat64A, exponentialHistogramFloat64C, equalExponentialHistograms[float64]))
	t.Run("ExponentialHistogramDataPointInt64", testDatatypeIgnoreTime(exponentialHistogramDataPointInt64A, exponentialHistogramDataPointInt64C, equalExponentialHistogramDataPoints[int64]))
	t.Run("ExponentialHistogramDataPointFloat64", testDatatypeIgnoreTime(exponentialHistogramDataPointFloat64A, exponentialHistogramDataPointFloat64C, equalExponentialHistogramDataPoints[float64]))
	t.Run("Summary", testDatatypeIgnoreTime(summaryA, summaryC, equalSummary))
	t.Run("SummaryDataPoint", testDatatypeIgnoreTime(summaryDataPointA, summaryDataPointC, equalSummaryDataPoint))
}

func TestAssertEqualIgnoreExemplars(t *testing.T) {
	hdpInt64 := histogramDataPointInt64A
	hdpInt64.Exemplars = []metricdata.Exemplar[int64]{exemplarInt64B}
	t.Run("HistogramDataPointInt64", testDatatypeIgnoreExemplars(histogramDataPointInt64A, hdpInt64, equalHistogramDataPoints[int64]))

	hdpFloat64 := histogramDataPointFloat64A
	hdpFloat64.Exemplars = []metricdata.Exemplar[float64]{exemplarFloat64B}
	t.Run("HistogramDataPointFloat64", testDatatypeIgnoreExemplars(histogramDataPointFloat64A, hdpFloat64, equalHistogramDataPoints[float64]))

	dpInt64 := dataPointInt64A
	dpInt64.Exemplars = []metricdata.Exemplar[int64]{exemplarInt64B}
	t.Run("DataPointInt64", testDatatypeIgnoreExemplars(dataPointInt64A, dpInt64, equalDataPoints[int64]))

	dpFloat64 := dataPointFloat64A
	dpFloat64.Exemplars = []metricdata.Exemplar[float64]{exemplarFloat64B}
	t.Run("DataPointFloat64", testDatatypeIgnoreExemplars(dataPointFloat64A, dpFloat64, equalDataPoints[float64]))

	ehdpInt64 := exponentialHistogramDataPointInt64A
	ehdpInt64.Exemplars = []metricdata.Exemplar[int64]{exemplarInt64B}
	t.Run("ExponentialHistogramDataPointInt64", testDatatypeIgnoreExemplars(exponentialHistogramDataPointInt64A, ehdpInt64, equalExponentialHistogramDataPoints[int64]))

	ehdpFloat64 := exponentialHistogramDataPointFloat64A
	ehdpFloat64.Exemplars = []metricdata.Exemplar[float64]{exemplarFloat64B}
	t.Run("ExponentialHistogramDataPointFloat64", testDatatypeIgnoreExemplars(exponentialHistogramDataPointFloat64A, ehdpFloat64, equalExponentialHistogramDataPoints[float64]))
}

func TestAssertEqualIgnoreValue(t *testing.T) {
	t.Run("ResourceMetrics", testDatatypeIgnoreValue(resourceMetricsA, resourceMetricsD, equalResourceMetrics))
	t.Run("ScopeMetrics", testDatatypeIgnoreValue(scopeMetricsA, scopeMetricsD, equalScopeMetrics))
	t.Run("Metrics", testDatatypeIgnoreValue(metricsA, metricsD, equalMetrics))
	t.Run("HistogramInt64", testDatatypeIgnoreValue(histogramInt64A, histogramInt64D, equalHistograms[int64]))
	t.Run("HistogramFloat64", testDatatypeIgnoreValue(histogramFloat64A, histogramFloat64D, equalHistograms[float64]))
	t.Run("SumInt64", testDatatypeIgnoreValue(sumInt64A, sumInt64D, equalSums[int64]))
	t.Run("SumFloat64", testDatatypeIgnoreValue(sumFloat64A, sumFloat64D, equalSums[float64]))
	t.Run("GaugeInt64", testDatatypeIgnoreValue(gaugeInt64A, gaugeInt64D, equalGauges[int64]))
	t.Run("GaugeFloat64", testDatatypeIgnoreValue(gaugeFloat64A, gaugeFloat64D, equalGauges[float64]))
	t.Run("HistogramDataPointInt64", testDatatypeIgnoreValue(histogramDataPointInt64A, histogramDataPointInt64D, equalHistogramDataPoints[int64]))
	t.Run("HistogramDataPointFloat64", testDatatypeIgnoreValue(histogramDataPointFloat64A, histogramDataPointFloat64D, equalHistogramDataPoints[float64]))
	t.Run("DataPointInt64", testDatatypeIgnoreValue(dataPointInt64A, dataPointInt64D, equalDataPoints[int64]))
	t.Run("DataPointFloat64", testDatatypeIgnoreValue(dataPointFloat64A, dataPointFloat64D, equalDataPoints[float64]))
	t.Run("ExemplarInt64", testDatatypeIgnoreValue(exemplarInt64A, exemplarInt64D, equalExemplars[int64]))
	t.Run("ExemplarFloat64", testDatatypeIgnoreValue(exemplarFloat64A, exemplarFloat64D, equalExemplars[float64]))
	t.Run("ExponentialHistogramInt64", testDatatypeIgnoreValue(exponentialHistogramInt64A, exponentialHistogramInt64D, equalExponentialHistograms[int64]))
	t.Run("ExponentialHistogramFloat64", testDatatypeIgnoreValue(exponentialHistogramFloat64A, exponentialHistogramFloat64D, equalExponentialHistograms[float64]))
	t.Run("ExponentialHistogramDataPointInt64", testDatatypeIgnoreValue(exponentialHistogramDataPointInt64A, exponentialHistogramDataPointInt64D, equalExponentialHistogramDataPoints[int64]))
	t.Run("ExponentialHistogramDataPointFloat64", testDatatypeIgnoreValue(exponentialHistogramDataPointFloat64A, exponentialHistogramDataPointFloat64D, equalExponentialHistogramDataPoints[float64]))
	t.Run("Summary", testDatatypeIgnoreValue(summaryA, summaryD, equalSummary))
	t.Run("SummaryDataPoint", testDatatypeIgnoreValue(summaryDataPointA, summaryDataPointD, equalSummaryDataPoint))
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
	AssertAggregationsEqual(t, histogramInt64A, histogramInt64A)
	AssertAggregationsEqual(t, histogramFloat64A, histogramFloat64A)
	AssertAggregationsEqual(t, exponentialHistogramInt64A, exponentialHistogramInt64A)
	AssertAggregationsEqual(t, exponentialHistogramFloat64A, exponentialHistogramFloat64A)
	AssertAggregationsEqual(t, summaryA, summaryA)

	r := equalAggregations(sumInt64A, nil, config{})
	assert.Len(t, r, 1, "should return nil comparison mismatch only")

	r = equalAggregations(sumInt64A, gaugeInt64A, config{})
	assert.Len(t, r, 1, "should return with type mismatch only")

	r = equalAggregations(unknownAggregation{}, unknownAggregation{}, config{})
	assert.Len(t, r, 1, "should return with unknown aggregation only")

	r = equalAggregations(sumInt64A, sumInt64B, config{})
	assert.NotEmptyf(t, r, "sums should not be equal: %v == %v", sumInt64A, sumInt64B)

	r = equalAggregations(sumInt64A, sumInt64C, config{ignoreTimestamp: true})
	assert.Empty(t, r, "sums should be equal: %v", r)

	r = equalAggregations(sumInt64A, sumInt64D, config{ignoreValue: true})
	assert.Empty(t, r, "value should be ignored: %v == %v", sumInt64A, sumInt64D)

	r = equalAggregations(sumFloat64A, sumFloat64B, config{})
	assert.NotEmptyf(t, r, "sums should not be equal: %v == %v", sumFloat64A, sumFloat64B)

	r = equalAggregations(sumFloat64A, sumFloat64C, config{ignoreTimestamp: true})
	assert.Empty(t, r, "sums should be equal: %v", r)

	r = equalAggregations(sumFloat64A, sumFloat64D, config{ignoreValue: true})
	assert.Empty(t, r, "value should be ignored: %v == %v", sumFloat64A, sumFloat64D)

	r = equalAggregations(gaugeInt64A, gaugeInt64B, config{})
	assert.NotEmptyf(t, r, "gauges should not be equal: %v == %v", gaugeInt64A, gaugeInt64B)

	r = equalAggregations(gaugeInt64A, gaugeInt64C, config{ignoreTimestamp: true})
	assert.Empty(t, r, "gauges should be equal: %v", r)

	r = equalAggregations(gaugeInt64A, gaugeInt64D, config{ignoreValue: true})
	assert.Empty(t, r, "value should be ignored: %v == %v", gaugeInt64A, gaugeInt64D)

	r = equalAggregations(gaugeFloat64A, gaugeFloat64B, config{})
	assert.NotEmptyf(t, r, "gauges should not be equal: %v == %v", gaugeFloat64A, gaugeFloat64B)

	r = equalAggregations(gaugeFloat64A, gaugeFloat64C, config{ignoreTimestamp: true})
	assert.Empty(t, r, "gauges should be equal: %v", r)

	r = equalAggregations(gaugeFloat64A, gaugeFloat64D, config{ignoreValue: true})
	assert.Empty(t, r, "value should be ignored: %v == %v", gaugeFloat64A, gaugeFloat64D)

	r = equalAggregations(histogramInt64A, histogramInt64B, config{})
	assert.NotEmptyf(t, r, "histograms should not be equal: %v == %v", histogramInt64A, histogramInt64B)

	r = equalAggregations(histogramInt64A, histogramInt64C, config{ignoreTimestamp: true})
	assert.Empty(t, r, "histograms should be equal: %v", r)

	r = equalAggregations(histogramInt64A, histogramInt64D, config{ignoreValue: true})
	assert.Empty(t, r, "value should be ignored: %v == %v", histogramInt64A, histogramInt64D)

	r = equalAggregations(histogramFloat64A, histogramFloat64B, config{})
	assert.NotEmptyf(t, r, "histograms should not be equal: %v == %v", histogramFloat64A, histogramFloat64B)

	r = equalAggregations(histogramFloat64A, histogramFloat64C, config{ignoreTimestamp: true})
	assert.Empty(t, r, "histograms should be equal: %v", r)

	r = equalAggregations(histogramFloat64A, histogramFloat64D, config{ignoreValue: true})
	assert.Empty(t, r, "value should be ignored: %v == %v", histogramFloat64A, histogramFloat64D)

	r = equalAggregations(exponentialHistogramInt64A, exponentialHistogramInt64B, config{})
	assert.NotEmptyf(t, r, "exponential histograms should not be equal: %v == %v", exponentialHistogramInt64A, exponentialHistogramInt64B)

	r = equalAggregations(exponentialHistogramInt64A, exponentialHistogramInt64C, config{ignoreTimestamp: true})
	assert.Empty(t, r, "exponential histograms should be equal: %v", r)

	r = equalAggregations(exponentialHistogramInt64A, exponentialHistogramInt64D, config{ignoreValue: true})
	assert.Empty(t, r, "value should be ignored: %v == %v", exponentialHistogramInt64A, exponentialHistogramInt64D)

	r = equalAggregations(exponentialHistogramFloat64A, exponentialHistogramFloat64B, config{})
	assert.NotEmptyf(t, r, "exponential histograms should not be equal: %v == %v", exponentialHistogramFloat64A, exponentialHistogramFloat64B)

	r = equalAggregations(exponentialHistogramFloat64A, exponentialHistogramFloat64C, config{ignoreTimestamp: true})
	assert.Empty(t, r, "exponential histograms should be equal: %v", r)

	r = equalAggregations(exponentialHistogramFloat64A, exponentialHistogramFloat64D, config{ignoreValue: true})
	assert.Empty(t, r, "value should be ignored: %v == %v", exponentialHistogramFloat64A, exponentialHistogramFloat64D)

	r = equalAggregations(summaryA, summaryB, config{})
	assert.NotEmptyf(t, r, "summaries should not be equal: %v == %v", summaryA, summaryB)

	r = equalAggregations(summaryA, summaryC, config{ignoreTimestamp: true})
	assert.Empty(t, r, "summaries should be equal: %v", r)

	r = equalAggregations(summaryA, summaryD, config{ignoreValue: true})
	assert.Empty(t, r, "value should be ignored: %v == %v", summaryA, summaryD)
}

func TestAssertAttributes(t *testing.T) {
	AssertHasAttributes(t, minFloat64A, attribute.Bool("A", true)) // No-op, always pass.
	AssertHasAttributes(t, exemplarInt64A, attribute.Bool("filter A", true))
	AssertHasAttributes(t, exemplarFloat64A, attribute.Bool("filter A", true))
	AssertHasAttributes(t, dataPointInt64A, attribute.Bool("A", true))
	AssertHasAttributes(t, dataPointFloat64A, attribute.Bool("A", true))
	AssertHasAttributes(t, gaugeInt64A, attribute.Bool("A", true))
	AssertHasAttributes(t, gaugeFloat64A, attribute.Bool("A", true))
	AssertHasAttributes(t, sumInt64A, attribute.Bool("A", true))
	AssertHasAttributes(t, sumFloat64A, attribute.Bool("A", true))
	AssertHasAttributes(t, histogramDataPointInt64A, attribute.Bool("A", true))
	AssertHasAttributes(t, histogramDataPointFloat64A, attribute.Bool("A", true))
	AssertHasAttributes(t, histogramInt64A, attribute.Bool("A", true))
	AssertHasAttributes(t, histogramFloat64A, attribute.Bool("A", true))
	AssertHasAttributes(t, metricsA, attribute.Bool("A", true))
	AssertHasAttributes(t, scopeMetricsA, attribute.Bool("A", true))
	AssertHasAttributes(t, resourceMetricsA, attribute.Bool("A", true))
	AssertHasAttributes(t, exponentialHistogramDataPointInt64A, attribute.Bool("A", true))
	AssertHasAttributes(t, exponentialHistogramDataPointFloat64A, attribute.Bool("A", true))
	AssertHasAttributes(t, exponentialHistogramInt64A, attribute.Bool("A", true))
	AssertHasAttributes(t, exponentialHistogramFloat64A, attribute.Bool("A", true))
	AssertHasAttributes(t, exponentialBucket2, attribute.Bool("A", true)) // No-op, always pass.
	AssertHasAttributes(t, summaryDataPointA, attribute.Bool("A", true))
	AssertHasAttributes(t, summaryA, attribute.Bool("A", true))
	AssertHasAttributes(t, quantileValueA, attribute.Bool("A", true)) // No-op, always pass.

	r := hasAttributesAggregation(gaugeInt64A, attribute.Bool("A", true))
	assert.Empty(t, r, "gaugeInt64A has A=True")
	r = hasAttributesAggregation(gaugeFloat64A, attribute.Bool("A", true))
	assert.Empty(t, r, "gaugeFloat64A has A=True")
	r = hasAttributesAggregation(sumInt64A, attribute.Bool("A", true))
	assert.Empty(t, r, "sumInt64A has A=True")
	r = hasAttributesAggregation(sumFloat64A, attribute.Bool("A", true))
	assert.Empty(t, r, "sumFloat64A has A=True")
	r = hasAttributesAggregation(histogramInt64A, attribute.Bool("A", true))
	assert.Empty(t, r, "histogramInt64A has A=True")
	r = hasAttributesAggregation(histogramFloat64A, attribute.Bool("A", true))
	assert.Empty(t, r, "histogramFloat64A has A=True")
	r = hasAttributesAggregation(exponentialHistogramInt64A, attribute.Bool("A", true))
	assert.Empty(t, r, "exponentialHistogramInt64A has A=True")
	r = hasAttributesAggregation(exponentialHistogramFloat64A, attribute.Bool("A", true))
	assert.Empty(t, r, "exponentialHistogramFloat64A has A=True")
	r = hasAttributesAggregation(summaryA, attribute.Bool("A", true))
	assert.Empty(t, r, "summaryA has A=True")

	r = hasAttributesAggregation(gaugeInt64A, attribute.Bool("A", false))
	assert.NotEmpty(t, r, "gaugeInt64A does not have A=False")
	r = hasAttributesAggregation(gaugeFloat64A, attribute.Bool("A", false))
	assert.NotEmpty(t, r, "gaugeFloat64A does not have A=False")
	r = hasAttributesAggregation(sumInt64A, attribute.Bool("A", false))
	assert.NotEmpty(t, r, "sumInt64A does not have A=False")
	r = hasAttributesAggregation(sumFloat64A, attribute.Bool("A", false))
	assert.NotEmpty(t, r, "sumFloat64A does not have A=False")
	r = hasAttributesAggregation(histogramInt64A, attribute.Bool("A", false))
	assert.NotEmpty(t, r, "histogramInt64A does not have A=False")
	r = hasAttributesAggregation(histogramFloat64A, attribute.Bool("A", false))
	assert.NotEmpty(t, r, "histogramFloat64A does not have A=False")
	r = hasAttributesAggregation(exponentialHistogramInt64A, attribute.Bool("A", false))
	assert.NotEmpty(t, r, "exponentialHistogramInt64A does not have A=False")
	r = hasAttributesAggregation(exponentialHistogramFloat64A, attribute.Bool("A", false))
	assert.NotEmpty(t, r, "exponentialHistogramFloat64A does not have A=False")
	r = hasAttributesAggregation(summaryA, attribute.Bool("A", false))
	assert.NotEmpty(t, r, "summaryA does not have A=False")

	r = hasAttributesAggregation(gaugeInt64A, attribute.Bool("B", true))
	assert.NotEmpty(t, r, "gaugeInt64A does not have Attribute B")
	r = hasAttributesAggregation(gaugeFloat64A, attribute.Bool("B", true))
	assert.NotEmpty(t, r, "gaugeFloat64A does not have Attribute B")
	r = hasAttributesAggregation(sumInt64A, attribute.Bool("B", true))
	assert.NotEmpty(t, r, "sumInt64A does not have Attribute B")
	r = hasAttributesAggregation(sumFloat64A, attribute.Bool("B", true))
	assert.NotEmpty(t, r, "sumFloat64A does not have Attribute B")
	r = hasAttributesAggregation(histogramInt64A, attribute.Bool("B", true))
	assert.NotEmpty(t, r, "histogramIntA does not have Attribute B")
	r = hasAttributesAggregation(histogramFloat64A, attribute.Bool("B", true))
	assert.NotEmpty(t, r, "histogramFloatA does not have Attribute B")
	r = hasAttributesAggregation(exponentialHistogramInt64A, attribute.Bool("B", true))
	assert.NotEmpty(t, r, "exponentialHistogramIntA does not have Attribute B")
	r = hasAttributesAggregation(exponentialHistogramFloat64A, attribute.Bool("B", true))
	assert.NotEmpty(t, r, "exponentialHistogramFloatA does not have Attribute B")
	r = hasAttributesAggregation(summaryA, attribute.Bool("B", true))
	assert.NotEmpty(t, r, "summaryA does not have Attribute B")
}

func TestAssertAttributesFail(t *testing.T) {
	fakeT := &testing.T{}
	assert.False(t, AssertHasAttributes(fakeT, dataPointInt64A, attribute.Bool("A", false)))
	assert.False(t, AssertHasAttributes(fakeT, dataPointFloat64A, attribute.Bool("B", true)))
	assert.False(t, AssertHasAttributes(fakeT, exemplarInt64A, attribute.Bool("A", false)))
	assert.False(t, AssertHasAttributes(fakeT, exemplarFloat64A, attribute.Bool("B", true)))
	assert.False(t, AssertHasAttributes(fakeT, gaugeInt64A, attribute.Bool("A", false)))
	assert.False(t, AssertHasAttributes(fakeT, gaugeFloat64A, attribute.Bool("B", true)))
	assert.False(t, AssertHasAttributes(fakeT, sumInt64A, attribute.Bool("A", false)))
	assert.False(t, AssertHasAttributes(fakeT, sumFloat64A, attribute.Bool("B", true)))
	assert.False(t, AssertHasAttributes(fakeT, histogramDataPointInt64A, attribute.Bool("A", false)))
	assert.False(t, AssertHasAttributes(fakeT, histogramDataPointFloat64A, attribute.Bool("B", true)))
	assert.False(t, AssertHasAttributes(fakeT, histogramInt64A, attribute.Bool("A", false)))
	assert.False(t, AssertHasAttributes(fakeT, histogramFloat64A, attribute.Bool("B", true)))
	assert.False(t, AssertHasAttributes(fakeT, metricsA, attribute.Bool("A", false)))
	assert.False(t, AssertHasAttributes(fakeT, metricsA, attribute.Bool("B", true)))
	assert.False(t, AssertHasAttributes(fakeT, resourceMetricsA, attribute.Bool("A", false)))
	assert.False(t, AssertHasAttributes(fakeT, resourceMetricsA, attribute.Bool("B", true)))
	assert.False(t, AssertHasAttributes(fakeT, exponentialHistogramDataPointInt64A, attribute.Bool("A", false)))
	assert.False(t, AssertHasAttributes(fakeT, exponentialHistogramDataPointFloat64A, attribute.Bool("B", true)))
	assert.False(t, AssertHasAttributes(fakeT, exponentialHistogramInt64A, attribute.Bool("A", false)))
	assert.False(t, AssertHasAttributes(fakeT, exponentialHistogramFloat64A, attribute.Bool("B", true)))
	assert.False(t, AssertHasAttributes(fakeT, summaryDataPointA, attribute.Bool("A", false)))
	assert.False(t, AssertHasAttributes(fakeT, summaryDataPointA, attribute.Bool("B", true)))
	assert.False(t, AssertHasAttributes(fakeT, summaryA, attribute.Bool("A", false)))
	assert.False(t, AssertHasAttributes(fakeT, summaryA, attribute.Bool("B", true)))

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

func AssertMarshal[N int64 | float64](t *testing.T, expected string, i *metricdata.Extrema[N]) {
	t.Helper()

	b, err := json.Marshal(i)
	assert.NoError(t, err)
	assert.Equal(t, expected, string(b))
}

func TestAssertMarshal(t *testing.T) {
	AssertMarshal(t, "null", &metricdata.Extrema[int64]{})

	AssertMarshal(t, "-1", &minFloat64A)
	AssertMarshal(t, "3", &minFloat64B)
	AssertMarshal(t, "-9.999999", &minFloat64D)
	AssertMarshal(t, "99", &maxFloat64B)

	AssertMarshal(t, "-1", &minInt64A)
	AssertMarshal(t, "3", &minInt64B)
	AssertMarshal(t, "99", &maxInt64B)
}

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metricdatatest_test

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
)

func ExampleAssertEqual() {
	ctx := context.Background()

	// Create a meterprovider with a reader
	reader := sdkmetric.NewManualReader()
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	defer func() {
		_ = mp.Shutdown(ctx)
	}()

	// Create an instrument(eg: counter/histogram/gauge) and simulate an operation
	meter := mp.Meter("payment-service")
	counter, _ := meter.Int64Counter("payment.requests")
	counter.Add(ctx, 5)

	// Collect the metrics
	rm := &metricdata.ResourceMetrics{}
	_ = reader.Collect(ctx, rm)
	got, _ := getMetrics("payment.requests", rm)

	want := metricdata.Metrics{
		Name: "payment.requests",
		Data: metricdata.Sum[int64]{
			DataPoints:  []metricdata.DataPoint[int64]{{Value: 5}},
			Temporality: metricdata.CumulativeTemporality,
			IsMonotonic: true,
		},
	}

	// Compare expected metrics with the received one
	t := &mockTestingT{}
	assertEqual := metricdatatest.AssertEqual(
		t,
		want,
		got,
		metricdatatest.IgnoreTimestamp(), // ignoring timestamps
	)
	fmt.Printf("Metrics are equal: %t\n", assertEqual)

	// Output:
	// Metrics are equal: true
}

func ExampleAssertAggregationsEqual() {
	ctx := context.Background()

	// Create a meterprovider with a reader
	reader := sdkmetric.NewManualReader()
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	defer func() {
		_ = mp.Shutdown(ctx)
	}()

	// Create an instrument(eg: counter/histogram/gauge) and simulate an operation
	meter := mp.Meter("payment-service")
	counter, _ := meter.Int64Counter("payment.count")
	counter.Add(ctx, 5)

	// Collect the metrics
	rm := &metricdata.ResourceMetrics{}
	_ = reader.Collect(ctx, rm)
	got, _ := getMetrics("payment.count", rm)

	want := metricdata.Metrics{
		Data: metricdata.Sum[int64]{
			DataPoints:  []metricdata.DataPoint[int64]{{Value: 5}},
			Temporality: metricdata.CumulativeTemporality,
			IsMonotonic: true,
		},
	}

	// Verify the expected data with the received one.
	t := &mockTestingT{}

	// Compare Aggregations
	hasEqualAggregations := metricdatatest.AssertAggregationsEqual(
		t,
		want.Data,
		got.Data,
		metricdatatest.IgnoreTimestamp(),
	)
	fmt.Printf("Aggregations are equal: %t\n", hasEqualAggregations)

	// Output:
	// Aggregations are equal: true
}

func ExampleAssertHasAttributes() {
	ctx := context.Background()

	// Create a meterprovider with a reader
	reader := sdkmetric.NewManualReader()
	mp := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))
	defer func() {
		_ = mp.Shutdown(ctx)
	}()

	// Simulate an operation using an instrument(eg: counter)
	meter := mp.Meter("payment-service")
	counter, _ := meter.Int64Counter("payment.requests")

	// Add attribute to the measurement
	attributes := attribute.NewSet(attribute.String("payment.method", "credit_card"))
	counter.Add(ctx, 5, metric.WithAttributeSet(attributes))

	// Collect the metrics
	rm := &metricdata.ResourceMetrics{}
	_ = reader.Collect(ctx, rm)
	metrics, _ := getMetrics("payment.requests", rm)

	// Verify the attributes in the received metrics
	t := &mockTestingT{}
	hasAttributes := metricdatatest.AssertHasAttributes(t, metrics, attributes.ToSlice()...)
	fmt.Printf("Metrics contains attributes: %t\n", hasAttributes)

	// Output:
	// Metrics contains attributes: true
}

func ExampleIgnoreValue() {
	want := metricdata.Metrics{
		Name: "payment.duration",
		Data: metricdata.Histogram[float64]{
			DataPoints: []metricdata.HistogramDataPoint[float64]{
				{
					Count:        2,
					Sum:          224.0,
					Bounds:       []float64{0, 5, 10, 25, 50, 75, 100, 250, 500, 750, 1000},
					BucketCounts: []uint64{0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0},
				},
			},
			Temporality: metricdata.CumulativeTemporality,
		},
	}

	got := metricdata.Metrics{
		Name: "payment.duration",
		Data: metricdata.Histogram[float64]{
			// Aggregate measurements are different in received metrics
			DataPoints: []metricdata.HistogramDataPoint[float64]{
				{
					Count:        10,
					Sum:          0,
					Bounds:       []float64{0, 5, 10, 25, 50, 75, 100, 250, 500, 750, 1000},
					BucketCounts: []uint64{1, 1, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0},
				},
			},
			Temporality: metricdata.CumulativeTemporality,
		},
	}

	t := &mockTestingT{}

	// Compare metrics without values
	ignoreValue := metricdatatest.AssertEqual(
		t,
		want,
		got,
		metricdatatest.IgnoreValue(),
	)
	fmt.Printf("Metrics are equal(ignoring values): %t\n", ignoreValue)

	// Output:
	// Metrics are equal(ignoring values): true
}

func ExampleIgnoreExemplars() {
	// Histogram data with Exemplars
	want := metricdata.Metrics{
		Name: "payment.duration",
		Data: metricdata.Histogram[float64]{
			DataPoints: []metricdata.HistogramDataPoint[float64]{
				{
					Count:        2,
					Sum:          224.0,
					Bounds:       []float64{0, 5, 10, 25, 50, 75, 100, 250, 500, 750, 1000},
					BucketCounts: []uint64{0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0},
					Exemplars: []metricdata.Exemplar[float64]{
						{
							FilteredAttributes: []attribute.KeyValue{
								attribute.String("payment.type", "recurring"),
							},
							Time:    time.Now(),
							Value:   15.0,
							SpanID:  []byte{1, 2, 3, 4, 5, 6, 7, 8},
							TraceID: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
						},
					},
				},
			},
			Temporality: metricdata.CumulativeTemporality,
		},
	}

	// Histogram data without Exemplars
	got := metricdata.Metrics{
		Name: "payment.duration",
		Data: metricdata.Histogram[float64]{
			// Aggregate measurements are different in received metrics(exemplars)
			DataPoints: []metricdata.HistogramDataPoint[float64]{
				{
					Count:        2,
					Sum:          224.0,
					Bounds:       []float64{0, 5, 10, 25, 50, 75, 100, 250, 500, 750, 1000},
					BucketCounts: []uint64{0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0},
				},
			},
			Temporality: metricdata.CumulativeTemporality,
		},
	}

	t := &mockTestingT{}

	// Compare metrics
	ignoreExemplars := metricdatatest.AssertEqual(
		t,
		want,
		got,
		metricdatatest.IgnoreExemplars(),
	)
	fmt.Printf("Metrics are equal(ignoring exemplars): %t\n", ignoreExemplars)

	// Output:
	// Metrics are equal(ignoring exemplars): true
}

// Helper function to retrieve the metrics.
// nolint:unparam // 'bool' return value currently unused, but retained for completeness
func getMetrics(name string, rm *metricdata.ResourceMetrics) (metricdata.Metrics, bool) {
	for _, scopeMetrics := range rm.ScopeMetrics {
		for _, m := range scopeMetrics.Metrics {
			if m.Name == name {
				return m, true
			}
		}
	}
	return metricdata.Metrics{}, false
}

// mockTestingT implements the [metricdatatest.TestingT] interface for examples.
// Usually, we use [*testing.T] as a substitute.
type mockTestingT struct{}

func (*mockTestingT) Helper() {}

func (*mockTestingT) Error(...any) {}

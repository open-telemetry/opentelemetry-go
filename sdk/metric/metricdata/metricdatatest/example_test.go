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
	actualMetrics, _ := getMetrics("payment.requests", rm)

	want := metricdata.Metrics{
		Name: "payment.requests",
		Data: metricdata.Sum[int64]{
			DataPoints:  []metricdata.DataPoint[int64]{{Value: 5}},
			Temporality: metricdata.CumulativeTemporality,
			IsMonotonic: true,
		},
	}

	// Compare expected metrics with the actual one
	t := &mockTestingT{}
	assertEqual := metricdatatest.AssertEqual(
		mockTest,
		expectedMetrics,
		actualMetrics,
		metricdatatest.IgnoreTimestamp(), // ignoring timestamps
	)
	fmt.Printf("Metrics matched as expected: %t\n", assertEqual)

	// Output:
	// Metrics matched as expected: true
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
	counter, _ := meter.Int64Counter("payment.duration")
	counter.Add(ctx, 5)

	// Collect the metrics
	rm := &metricdata.ResourceMetrics{}
	_ = reader.Collect(ctx, rm)
	actualMetrics, _ := getMetrics("payment.duration", rm)

	expectedMetrics := metricdata.Metrics{
		Data: metricdata.Sum[int64]{
			DataPoints:  []metricdata.DataPoint[int64]{{Value: 5}},
			Temporality: metricdata.CumulativeTemporality,
			IsMonotonic: true,
		},
	}

	// Verify the expected data with the actual one.
	mockTest := &mockTestingT{}

	// Compare Aggregations
	hasEqualAggregations := metricdatatest.AssertAggregationsEqual(
		mockTest,
		expectedMetrics.Data,
		actualMetrics.Data,
		metricdatatest.IgnoreTimestamp(),
	)
	fmt.Printf("Aggregations are matching as expected: %t\n", hasEqualAggregations)

	// Output:
	// Aggregations are matching as expected: true
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
	actualMetrics, _ := getMetrics("payment.requests", rm)

	// Verify the attributes in the actualMetrics
	mockTest := &mockTestingT{}
	hasAttributes := metricdatatest.AssertHasAttributes(mockTest, actualMetrics, attributes.ToSlice()...)
	fmt.Printf("Metrics has expected attributes : %t\n", hasAttributes)

	// Output:
	// Metrics has expected attributes : true
}

func ExampleIgnoreValue() {
	expectedMetrics := metricdata.Metrics{
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

	actualMetrics := metricdata.Metrics{
		Name: "payment.duration",
		Data: metricdata.Histogram[float64]{
			// Aggregate measurements are different in expected metrics
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

	mockTest := &mockTestingT{}

	// Compare metrics without values
	ignoreValue := metricdatatest.AssertEqual(
		mockTest,
		expectedMetrics,
		actualMetrics,
		metricdatatest.IgnoreValue(),
	)
	fmt.Printf("Metrics matched irrespective of difference in values: %t\n", ignoreValue)

	// Output:
	// Metrics matched irrespective of difference in values: true
}

func ExampleIgnoreExemplars() {
	// Histogram data with Exemplars
	expectedMetrics := metricdata.Metrics{
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
	actualMetrics := metricdata.Metrics{
		Name: "payment.duration",
		Data: metricdata.Histogram[float64]{
			// Aggregate measurements are different from expected metrics
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

	mockTest := &mockTestingT{}

	// Compare metrics
	ignoreExemplars := metricdatatest.AssertEqual(
		mockTest,
		expectedMetrics,
		actualMetrics,
		metricdatatest.IgnoreExemplars(),
	)
	fmt.Printf("Metrics matched irrespective of difference in exemplars: %t\n", ignoreExemplars)

	// Output:
	// Metrics matched irrespective of difference in exemplars: true
}

// Helper function to retrieve the metrics.
func getMetrics(name string, rm *metricdata.ResourceMetrics) (metricdata.Metrics, bool) { //nolint
	for _, scopeMetrics := range rm.ScopeMetrics {
		for _, m := range scopeMetrics.Metrics {
			if m.Name == name {
				return m, true
			}
		}
	}
	return metricdata.Metrics{}, false
}

// mockTestingT implements the metricdatatest.TestingT interface for examples.
type mockTestingT struct {
	errors []string //nolint
}

func (*mockTestingT) Helper() {}

func (*mockTestingT) Error(...any) {}

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metricdatatest_test

import (
	"context"
	"fmt"
	"log"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/metricdata/metricdatatest"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

// This example demonstrates a scenario using manualreader, meterprovider,
// and metricdatatest to verify metrics in a testing environment.
func Example() {
	ctx := context.Background()

	// Create a resource
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("payment-processor"),
		semconv.ServiceVersion("0.1.0"),
	)

	// Create a manual reader for collecting metrics on demand
	reader := sdkmetric.NewManualReader()

	// Create a meter provider with the manual reader
	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(reader),
	)
	defer func() {
		err := meterProvider.Shutdown(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
	}()

	// Get a meter from the provider
	meter := meterProvider.Meter("payment-processor-metrics")

	// Simulate operations and record metrics
	simulateOperationsAndRecordMetrics(ctx, meter)

	// Collect the metrics
	metrics := metricdata.ResourceMetrics{}
	err := reader.Collect(ctx, &metrics)
	if err != nil {
		fmt.Printf("Failed to collect metrics: %v\n", err)
		return
	}

	// Create expected metrics for comparison
	expectedCounter := createExpectedCounter()
	expectedHistogram := createExpectedHistogram()
	expectedGauge := createExpectedGauge()

	// Find the actual metrics from the collected data
	var actualCounter, actualHistogram, actualGauge metricdata.Metrics
	for _, scopeMetrics := range metrics.ScopeMetrics {
		for _, m := range scopeMetrics.Metrics {
			switch m.Name {
			case "payment.requests":
				actualCounter = m
			case "payment.duration":
				actualHistogram = m
			case "payment.balance":
				actualGauge = m
			}
		}
	}

	// In a real test, you would use a testing.T instance instead of mockT
	mockT := &mockTestingT{}

	// Verify counter metrics
	counterEqual := metricdatatest.AssertEqual(
		mockT,
		expectedCounter,
		actualCounter,
		metricdatatest.IgnoreTimestamp(),
	)
	fmt.Printf("Counter metrics match: %t\n", counterEqual)

	// Verify counter aggregations
	hasAggregationsEqual := metricdatatest.AssertAggregationsEqual(
		mockT,
		expectedCounter.Data,
		actualCounter.Data,
		metricdatatest.IgnoreTimestamp(),
	)
	fmt.Printf("Counter has expected aggregations: %t\n", hasAggregationsEqual)

	// Verify histogram metrics (ignoring exact values)
	histogramEqual := metricdatatest.AssertEqual(
		mockT,
		expectedHistogram,
		actualHistogram,
		metricdatatest.IgnoreTimestamp(),
		metricdatatest.IgnoreValue(),
	)
	fmt.Printf("Histogram metrics match (ignoring values): %t\n", histogramEqual)

	// Verify gauge metrics
	gaugeEqual := metricdatatest.AssertEqual(
		mockT,
		expectedGauge,
		actualGauge,
		metricdatatest.IgnoreTimestamp(),
	)
	fmt.Printf("Gauge metrics match: %t\n", gaugeEqual)

	// Verify attributes on the counter data points
	hasAttrs := metricdatatest.AssertHasAttributes(
		mockT,
		actualCounter,
		attribute.String("payment.method", "credit_card"),
	)
	fmt.Printf("Counter has expected attributes: %t\n", hasAttrs)

	// Output:
	// Counter metrics match: true
	// Counter has expected aggregations: true
	// Histogram metrics match (ignoring values): true
	// Gauge metrics match: true
	// Counter has expected attributes: true
}

// Helper function to simulate operations and record metrics
func simulateOperationsAndRecordMetrics(ctx context.Context, meter metric.Meter) {
	// Create instruments (errors are ignored in this example)
	counter, _ := meter.Int64Counter(
		"payment.requests",
		metric.WithDescription("Number of payment requests received"),
		metric.WithUnit("{request}"),
	)

	histogram, _ := meter.Float64Histogram(
		"payment.duration",
		metric.WithDescription("Duration of payment processing"),
		metric.WithUnit("ms"),
	)

	gauge, _ := meter.Float64Gauge(
		"payment.balance",
		metric.WithDescription("Current account balance"),
		metric.WithUnit("USD"),
	)

	commonAttrs := attribute.NewSet(
		attribute.String("payment.method", "credit_card"),
	)

	// Simulate processing payments
	counter.Add(ctx, 3, metric.WithAttributeSet(commonAttrs))

	// Record processing durations
	histogram.Record(ctx, 125.3, metric.WithAttributeSet(commonAttrs))
	histogram.Record(ctx, 98.7, metric.WithAttributeSet(commonAttrs))

	// Record current balance
	gauge.Record(ctx, 1250.60, metric.WithAttributeSet(commonAttrs))
}

// Helper function to create expected counter metrics
func createExpectedCounter() metricdata.Metrics {
	commonAttrs := attribute.NewSet(
		attribute.String("payment.method", "credit_card"),
	)

	return metricdata.Metrics{
		Name:        "payment.requests",
		Description: "Number of payment requests received",
		Unit:        "{request}",
		Data: metricdata.Sum[int64]{
			DataPoints: []metricdata.DataPoint[int64]{
				{
					Attributes: commonAttrs,
					Value:      3,
				},
			},
			Temporality: metricdata.CumulativeTemporality,
			IsMonotonic: true,
		},
	}
}

// Helper function to create expected histogram metrics
func createExpectedHistogram() metricdata.Metrics {
	commonAttrs := attribute.NewSet(
		attribute.String("payment.method", "credit_card"),
	)

	return metricdata.Metrics{
		Name:        "payment.duration",
		Description: "Duration of payment processing",
		Unit:        "ms",
		Data: metricdata.Histogram[float64]{
			DataPoints: []metricdata.HistogramDataPoint[float64]{
				{
					Attributes: commonAttrs,
					//Values do not matter since we use metricdatatest.IgnoreValue() while asserting
					Count:        2,
					Sum:          224.0, // 125.3 + 98.7
					Bounds:       []float64{0, 5, 10, 25, 50, 75, 100, 250, 500, 750, 1000},
					BucketCounts: []uint64{0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0},
				},
			},
			Temporality: metricdata.CumulativeTemporality,
		},
	}
}

// Helper function to create expected gauge metrics
func createExpectedGauge() metricdata.Metrics {
	commonAttrs := attribute.NewSet(
		attribute.String("payment.method", "credit_card"),
	)

	return metricdata.Metrics{
		Name:        "payment.balance",
		Description: "Current account balance",
		Unit:        "USD",
		Data: metricdata.Gauge[float64]{
			DataPoints: []metricdata.DataPoint[float64]{
				{
					Attributes: commonAttrs,
					Value:      1250.60,
				},
			},
		},
	}
}

// mockTestingT implements the metricdatatest.TestingT interface for examples
type mockTestingT struct {
	errors []string
}

func (m *mockTestingT) Helper() {}

func (m *mockTestingT) Error(args ...any) {}

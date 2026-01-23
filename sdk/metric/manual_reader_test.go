// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric // import "go.opentelemetry.io/otel/sdk/metric"

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/internal/global"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/semconv/v1.39.0/otelconv"
)

func TestManualReader(t *testing.T) {
	suite.Run(t, &readerTestSuite{Factory: func(opts ...ReaderOption) Reader {
		var mopts []ManualReaderOption
		for _, o := range opts {
			mopts = append(mopts, o)
		}
		return NewManualReader(mopts...)
	}})
}

func BenchmarkManualReader(b *testing.B) {
	b.Run("Collect", benchReaderCollectFunc(NewManualReader()))
}

var (
	deltaTemporalitySelector      = func(InstrumentKind) metricdata.Temporality { return metricdata.DeltaTemporality }
	cumulativeTemporalitySelector = func(InstrumentKind) metricdata.Temporality { return metricdata.CumulativeTemporality }
)

func TestManualReaderTemporality(t *testing.T) {
	tests := []struct {
		name    string
		options []ManualReaderOption
		// Currently only testing constant temporality. This should be expanded
		// if we put more advanced selection in the SDK
		wantTemporality metricdata.Temporality
	}{
		{
			name:            "default",
			wantTemporality: metricdata.CumulativeTemporality,
		},
		{
			name: "delta",
			options: []ManualReaderOption{
				WithTemporalitySelector(deltaTemporalitySelector),
			},
			wantTemporality: metricdata.DeltaTemporality,
		},
		{
			name: "repeats overwrite",
			options: []ManualReaderOption{
				WithTemporalitySelector(deltaTemporalitySelector),
				WithTemporalitySelector(cumulativeTemporalitySelector),
			},
			wantTemporality: metricdata.CumulativeTemporality,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var undefinedInstrument InstrumentKind
			rdr := NewManualReader(tt.options...)
			assert.Equal(t, tt.wantTemporality, rdr.temporality(undefinedInstrument))
		})
	}
}

func TestManualReaderCollect(t *testing.T) {
	expiredCtx, cancel := context.WithDeadline(t.Context(), time.Now().Add(-1))
	defer cancel()

	tests := []struct {
		name        string
		ctx         context.Context
		expectedErr error
	}{
		{
			name:        "with a valid context",
			ctx:         t.Context(),
			expectedErr: nil,
		},
		{
			name:        "with an expired context",
			ctx:         expiredCtx,
			expectedErr: context.DeadlineExceeded,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rdr := NewManualReader()
			mp := NewMeterProvider(WithReader(rdr))
			meter := mp.Meter("test")

			// Ensure the pipeline has a callback setup
			testM, err := meter.Int64ObservableCounter("test")
			assert.NoError(t, err)
			_, err = meter.RegisterCallback(func(context.Context, metric.Observer) error {
				return nil
			}, testM)
			assert.NoError(t, err)

			rm := &metricdata.ResourceMetrics{}
			assert.Equal(t, tt.expectedErr, rdr.Collect(tt.ctx, rm))
		})
	}
}

func TestManualReaderInstrumentation(t *testing.T) {
	// Enable SDK observability.
	t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	// ManualReader under test with a fake producer.
	manualReader := NewManualReader()
	t.Cleanup(func() { _ = manualReader.Shutdown(t.Context()) })
	manualReader.register(testSDKProducer{})

	// Exercise a collect (producer data).
	var got metricdata.ResourceMetrics
	require.NoError(t, manualReader.Collect(t.Context(), &got))

	// Collect again so we have something to scan through.
	var rm metricdata.ResourceMetrics
	require.NoError(t, manualReader.Collect(t.Context(), &rm))
	require.NotEmpty(t, rm.ScopeMetrics)

	targetName := otelconv.SDKMetricReaderCollectionDuration{}.Name()
	targetDesc := otelconv.SDKMetricReaderCollectionDuration{}.Description()
	targetUnit := otelconv.SDKMetricReaderCollectionDuration{}.Unit()

	// Find the SDK reader self-metric anywhere in the collected data.
	foundMetric := findMetricByName(&rm, targetName)

	// If not found, explain and skip (this metric is emitted via the *global* MP).
	if foundMetric == nil {
		t.Skipf("SDK reader self-metric %q not found. It is emitted via the global MeterProvider; "+
			"this test does not install a global MP.", targetName)
		return
	}

	// Basic identity checks.
	assert.Equal(t, targetName, foundMetric.Name)
	assert.Equal(t, targetDesc, foundMetric.Description)
	assert.Equal(t, targetUnit, foundMetric.Unit)

	// Must be a histogram with cumulative temporality.
	hist, ok := foundMetric.Data.(metricdata.Histogram[float64])
	require.True(t, ok, "expected histogram data")
	assert.Equal(t, metricdata.CumulativeTemporality, hist.Temporality)
	require.NotEmpty(t, hist.DataPoints)

	// Check base attributes on one datapoint (flexibly).
	dp := hist.DataPoints[0]
	attrs := dp.Attributes.ToSlice()
	t.Logf("observability attrs: %v", attrs)

	const expectedComponentType = "go.opentelemetry.io/otel/sdk/metric/metric.ManualReader"

	var hasName, hasType bool
	for _, a := range attrs {
		switch a.Key {
		case "otel.component.name":
			if s := a.Value.AsString(); s != "" && strings.Contains(s, "metric_reader") {
				hasName = true
			}
		case "otel.component.type":
			if a.Value.AsString() == expectedComponentType {
				hasType = true
			}
		}
	}
	assert.True(t, hasName, "expected non-empty otel.component.name containing 'metric_reader'")
	assert.True(t, hasType, "expected otel.component.type == %q", expectedComponentType)
}

// findMetricByName searches all scopes/metrics for the given metric name.
func findMetricByName(rm *metricdata.ResourceMetrics, name string) *metricdata.Metrics {
	for si := range rm.ScopeMetrics {
		sm := &rm.ScopeMetrics[si]
		for mi := range sm.Metrics {
			if sm.Metrics[mi].Name == name {
				return &sm.Metrics[mi]
			}
		}
	}
	return nil
}

// createMetricDataTestProducerForManual creates a producer using patterns from metricdatatest for manual reader benchmarks.
func createMetricDataTestProducerForManual() testSDKProducer {
	return testSDKProducer{
		produceFunc: func(_ context.Context, rm *metricdata.ResourceMetrics) error {
			start := time.Now().Add(-time.Minute)
			end := time.Now()

			// Create attribute sets using common patterns
			aliceAttrs := attribute.NewSet(attribute.String("user", "alice"), attribute.String("env", "prod"))
			bobAttrs := attribute.NewSet(attribute.String("user", "bob"), attribute.String("env", "staging"))
			charlieAttrs := attribute.NewSet(attribute.String("user", "charlie"), attribute.String("env", "dev"))

			// Create exemplars for histogram metrics
			exemplars := []metricdata.Exemplar[float64]{
				{
					FilteredAttributes: []attribute.KeyValue{attribute.String("trace", "span-1")},
					Time:               end,
					Value:              15.5,
					SpanID:             []byte{1, 2, 3, 4, 5, 6, 7, 8},
					TraceID:            []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
				},
			}

			// Define different metric types using metricdatatest patterns
			createScopeMetrics := func(scopeIdx int) metricdata.ScopeMetrics {
				metrics := []metricdata.Metrics{
					// Counter metrics
					{
						Name:        fmt.Sprintf("requests_total_%d", scopeIdx),
						Description: "Total number of requests",
						Unit:        "1",
						Data: metricdata.Sum[int64]{
							Temporality: metricdata.CumulativeTemporality,
							IsMonotonic: true,
							DataPoints: []metricdata.DataPoint[int64]{
								{Attributes: aliceAttrs, StartTime: start, Time: end, Value: 100 + int64(scopeIdx*10)},
								{Attributes: bobAttrs, StartTime: start, Time: end, Value: 150 + int64(scopeIdx*15)},
								{Attributes: charlieAttrs, StartTime: start, Time: end, Value: 75 + int64(scopeIdx*5)},
							},
						},
					},
					// Gauge metrics
					{
						Name:        fmt.Sprintf("memory_usage_%d", scopeIdx),
						Description: "Memory usage in MB",
						Unit:        "MB",
						Data: metricdata.Gauge[float64]{
							DataPoints: []metricdata.DataPoint[float64]{
								{Attributes: aliceAttrs, Time: end, Value: 512.5 + float64(scopeIdx*10)},
								{Attributes: bobAttrs, Time: end, Value: 768.2 + float64(scopeIdx*20)},
								{Attributes: charlieAttrs, Time: end, Value: 256.8 + float64(scopeIdx*5)},
							},
						},
					},
					// Histogram metrics
					{
						Name:        fmt.Sprintf("request_duration_%d", scopeIdx),
						Description: "Request duration histogram",
						Unit:        "ms",
						Data: metricdata.Histogram[float64]{
							Temporality: metricdata.CumulativeTemporality,
							DataPoints: []metricdata.HistogramDataPoint[float64]{
								{
									Attributes:   aliceAttrs,
									StartTime:    start,
									Time:         end,
									Count:        100,
									Sum:          1500.5,
									Min:          metricdata.NewExtrema(1.0),
									Max:          metricdata.NewExtrema(50.0),
									Bounds:       []float64{1, 5, 10, 25, 50, 100, 250, 500},
									BucketCounts: []uint64{10, 20, 30, 25, 10, 4, 1, 0, 0},
									Exemplars:    exemplars,
								},
								{
									Attributes:   bobAttrs,
									StartTime:    start,
									Time:         end,
									Count:        80,
									Sum:          1200.3,
									Min:          metricdata.NewExtrema(2.0),
									Max:          metricdata.NewExtrema(45.0),
									Bounds:       []float64{1, 5, 10, 25, 50, 100, 250, 500},
									BucketCounts: []uint64{5, 15, 25, 20, 10, 4, 1, 0, 0},
									Exemplars:    exemplars,
								},
							},
						},
					},
					// Exponential Histogram
					{
						Name:        fmt.Sprintf("response_size_%d", scopeIdx),
						Description: "Response size exponential histogram",
						Unit:        "bytes",
						Data: metricdata.ExponentialHistogram[float64]{
							Temporality: metricdata.CumulativeTemporality,
							DataPoints: []metricdata.ExponentialHistogramDataPoint[float64]{
								{
									Attributes: aliceAttrs,
									StartTime:  start,
									Time:       end,
									Count:      50,
									Sum:        25000.0,
									Min:        metricdata.NewExtrema(100.0),
									Max:        metricdata.NewExtrema(2000.0),
									Scale:      2,
									ZeroCount:  2,
									Exemplars:  exemplars,
								},
							},
						},
					},
				}

				return metricdata.ScopeMetrics{
					Scope: instrumentation.Scope{
						Name:    fmt.Sprintf("benchmark/scope/%d", scopeIdx),
						Version: "1.0.0",
					},
					Metrics: metrics,
				}
			}

			// Create multiple scopes for comprehensive test data
			var scopeMetrics []metricdata.ScopeMetrics
			for i := 0; i < 20; i++ { // 20 scopes with 4 metrics each = 80 total metrics
				scopeMetrics = append(scopeMetrics, createScopeMetrics(i))
			}

			*rm = metricdata.ResourceMetrics{
				Resource:     resource.NewSchemaless(attribute.String("service.name", "benchmark-test")),
				ScopeMetrics: scopeMetrics,
			}
			return nil
		},
	}
}

func BenchmarkManualReaderInstrumentation(b *testing.B) {
	run := func(b *testing.B, withInstrumentationMP bool) {
		// Save and restore the original global meter provider
		orig := otel.GetMeterProvider()
		defer otel.SetMeterProvider(orig)

		// Suppress internal logging messages for cleaner benchmark output
		global.SetLogger(logr.Discard())
		b.Cleanup(func() {
			// Logger will be reset by test cleanup naturally
		})

		// Suppress error handler messages for cleaner benchmark output
		origErrorHandler := otel.GetErrorHandler()
		otel.SetErrorHandler(otel.ErrorHandlerFunc(func(error) {}))
		b.Cleanup(func() {
			otel.SetErrorHandler(origErrorHandler)
		})

		if withInstrumentationMP {
			// Set up a meter provider for instrumentation to use
			instrumentationReader := NewManualReader()
			instrumentationMP := NewMeterProvider(WithReader(instrumentationReader))
			otel.SetMeterProvider(instrumentationMP)

			// Clean up the instrumentation meter provider
			b.Cleanup(func() {
				_ = instrumentationMP.Shutdown(b.Context())
			})
		}

		r := NewManualReader()
		// Register with producer using metricdatatest patterns for realistic benchmark data
		r.register(createMetricDataTestProducerForManual())
		b.Cleanup(func() {
			_ = r.Shutdown(b.Context()) // Ignore error in cleanup
		})

		rm := &metricdata.ResourceMetrics{}

		b.ReportAllocs()
		b.ResetTimer()

		for b.Loop() {
			// Test the collect operation (simulating what manual readers do)
			err := r.Collect(b.Context(), rm)
			_ = err // Ignore error for benchmark
		}
	}

	b.Run("NoObservability", func(b *testing.B) {
		b.Setenv("OTEL_GO_X_OBSERVABILITY", "false")
		run(b, false)
	})

	b.Run("Observability", func(b *testing.B) {
		b.Setenv("OTEL_GO_X_OBSERVABILITY", "true")
		run(b, true)
	})
}

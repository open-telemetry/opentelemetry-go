// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric // import "go.opentelemetry.io/otel/sdk/metric"

import (
	"context"
	"errors"
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

const testDur = time.Second * 2

func TestWithTimeout(t *testing.T) {
	test := func(d time.Duration) time.Duration {
		opts := []PeriodicReaderOption{WithTimeout(d)}
		return newPeriodicReaderConfig(opts).timeout
	}

	assert.Equal(t, testDur, test(testDur))
	assert.Equal(t, defaultTimeout, newPeriodicReaderConfig(nil).timeout)
	assert.Equal(t, defaultTimeout, test(time.Duration(0)), "invalid timeout should use default")
	assert.Equal(t, defaultTimeout, test(time.Duration(-1)), "invalid timeout should use default")
}

func TestTimeoutEnvVar(t *testing.T) {
	testCases := []struct {
		v    string
		want time.Duration
	}{
		{
			// empty value
			"",
			defaultTimeout,
		},
		{
			// positive value
			"1",
			time.Millisecond,
		},
		{
			// non-positive value
			"0",
			defaultTimeout,
		},
		{
			// value with unit (not supported)
			"1ms",
			defaultTimeout,
		},
		{
			// NaN
			"abc",
			defaultTimeout,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.v, func(t *testing.T) {
			t.Setenv(envTimeout, tc.v)
			got := newPeriodicReaderConfig(nil).timeout
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestTimeoutEnvAndOption(t *testing.T) {
	want := 5 * time.Millisecond
	t.Setenv(envTimeout, "999")
	opts := []PeriodicReaderOption{WithTimeout(want)}
	got := newPeriodicReaderConfig(opts).timeout
	assert.Equal(t, want, got, "option should have precedence over env var")
}

func TestWithInterval(t *testing.T) {
	test := func(d time.Duration) time.Duration {
		opts := []PeriodicReaderOption{WithInterval(d)}
		return newPeriodicReaderConfig(opts).interval
	}

	assert.Equal(t, testDur, test(testDur))
	assert.Equal(t, defaultInterval, newPeriodicReaderConfig(nil).interval)
	assert.Equal(t, defaultInterval, test(time.Duration(0)), "invalid interval should use default")
	assert.Equal(t, defaultInterval, test(time.Duration(-1)), "invalid interval should use default")
}

func TestIntervalEnvVar(t *testing.T) {
	testCases := []struct {
		v    string
		want time.Duration
	}{
		{
			// empty value
			"",
			defaultInterval,
		},
		{
			// positive value
			"1",
			time.Millisecond,
		},
		{
			// non-positive value
			"0",
			defaultInterval,
		},
		{
			// value with unit (not supported)
			"1ms",
			defaultInterval,
		},
		{
			// NaN
			"abc",
			defaultInterval,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.v, func(t *testing.T) {
			t.Setenv(envInterval, tc.v)
			got := newPeriodicReaderConfig(nil).interval
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestIntervalEnvAndOption(t *testing.T) {
	want := 5 * time.Millisecond
	t.Setenv(envInterval, "999")
	opts := []PeriodicReaderOption{WithInterval(want)}
	got := newPeriodicReaderConfig(opts).interval
	assert.Equal(t, want, got, "option should have precedence over env var")
}

type fnExporter struct {
	temporalityFunc TemporalitySelector
	aggregationFunc AggregationSelector
	exportFunc      func(context.Context, *metricdata.ResourceMetrics) error
	flushFunc       func(context.Context) error
	shutdownFunc    func(context.Context) error
}

var _ Exporter = (*fnExporter)(nil)

func (e *fnExporter) Temporality(k InstrumentKind) metricdata.Temporality {
	if e.temporalityFunc != nil {
		return e.temporalityFunc(k)
	}
	return DefaultTemporalitySelector(k)
}

func (e *fnExporter) Aggregation(k InstrumentKind) Aggregation {
	if e.aggregationFunc != nil {
		return e.aggregationFunc(k)
	}
	return DefaultAggregationSelector(k)
}

func (e *fnExporter) Export(ctx context.Context, m *metricdata.ResourceMetrics) error {
	if e.exportFunc != nil {
		return e.exportFunc(ctx, m)
	}
	return nil
}

func (e *fnExporter) ForceFlush(ctx context.Context) error {
	if e.flushFunc != nil {
		return e.flushFunc(ctx)
	}
	return nil
}

func (e *fnExporter) Shutdown(ctx context.Context) error {
	if e.shutdownFunc != nil {
		return e.shutdownFunc(ctx)
	}
	return nil
}

type periodicReaderTestSuite struct {
	*readerTestSuite

	ErrReader *PeriodicReader
}

func (ts *periodicReaderTestSuite) SetupTest() {
	e := &fnExporter{
		exportFunc:   func(context.Context, *metricdata.ResourceMetrics) error { return assert.AnError },
		flushFunc:    func(context.Context) error { return assert.AnError },
		shutdownFunc: func(context.Context) error { return assert.AnError },
	}

	ts.ErrReader = NewPeriodicReader(e, WithProducer(testExternalProducer{}))
	ts.ErrReader.register(testSDKProducer{})
}

func (ts *periodicReaderTestSuite) TearDownTest() {
	ts.readerTestSuite.TearDownTest()

	_ = ts.ErrReader.Shutdown(context.Background())
}

func (ts *periodicReaderTestSuite) TestForceFlushPropagated() {
	ts.ErrorIs(ts.ErrReader.ForceFlush(context.Background()), assert.AnError)
}

func (ts *periodicReaderTestSuite) TestShutdownPropagated() {
	ts.ErrorIs(ts.ErrReader.Shutdown(context.Background()), assert.AnError)
}

func TestPeriodicReader(t *testing.T) {
	suite.Run(t, &periodicReaderTestSuite{
		readerTestSuite: &readerTestSuite{
			Factory: func(opts ...ReaderOption) Reader {
				var popts []PeriodicReaderOption
				for _, o := range opts {
					popts = append(popts, o)
				}
				return NewPeriodicReader(new(fnExporter), popts...)
			},
		},
	})
}

type chErrorHandler struct {
	Err chan error
}

func newChErrorHandler() *chErrorHandler {
	return &chErrorHandler{
		Err: make(chan error, 1),
	}
}

func (eh chErrorHandler) Handle(err error) {
	eh.Err <- err
}

func triggerTicker(t *testing.T) chan time.Time {
	t.Helper()

	// Override the ticker C chan so tests are not flaky and rely on timing.
	orig := newTicker
	t.Cleanup(func() { newTicker = orig })

	// Keep this at size zero so when triggered with a send it will hang until
	// the select case is selected and the collection loop is started.
	trigger := make(chan time.Time)
	newTicker = func(d time.Duration) *time.Ticker {
		ticker := time.NewTicker(d)
		ticker.C = trigger
		return ticker
	}
	return trigger
}

func TestPeriodicReaderRun(t *testing.T) {
	trigger := triggerTicker(t)

	// Register an error handler to validate export errors are passed to
	// otel.Handle.
	defer func(orig otel.ErrorHandler) {
		otel.SetErrorHandler(orig)
	}(otel.GetErrorHandler())
	eh := newChErrorHandler()
	otel.SetErrorHandler(eh)

	exp := &fnExporter{
		exportFunc: func(_ context.Context, m *metricdata.ResourceMetrics) error {
			// The testSDKProducer produces testResourceMetricsAB.
			assert.Equal(t, testResourceMetricsAB, *m)
			return assert.AnError
		},
	}

	r := NewPeriodicReader(exp, WithProducer(testExternalProducer{}))
	r.register(testSDKProducer{})
	trigger <- time.Now()
	assert.ErrorIs(t, <-eh.Err, assert.AnError)

	// Ensure Reader is allowed clean up attempt.
	_ = r.Shutdown(t.Context())
}

func TestPeriodicReaderBatching(t *testing.T) {
	trigger := triggerTicker(t)

	// Register an error handler to validate export errors are passed to
	// otel.Handle.
	defer func(orig otel.ErrorHandler) {
		otel.SetErrorHandler(orig)
	}(otel.GetErrorHandler())
	eh := newChErrorHandler()
	otel.SetErrorHandler(eh)

	expectations := []metricdata.ResourceMetrics{
		testResourceMetricsAB,
		testResourceMetricsC1,
		testResourceMetricsC2,
	}

	expectationIdx := 0
	exp := &fnExporter{
		exportFunc: func(_ context.Context, m *metricdata.ResourceMetrics) error {
			// collectAndExport is potentially called multiple times, so just
			// make sure batches are split correctly and are in order.
			expect := expectations[expectationIdx%len(expectations)]
			// The testSDKProducer produces three batches of metrics.
			assert.Equal(t, expect, *m, fmt.Sprintf("expectations[%d] not equal", expectationIdx))
			expectationIdx++
			return assert.AnError
		},
	}

	r := NewPeriodicReader(
		exp,
		WithMaxExportBatchSize(2),
		WithProducer(testExternalProducer{}),
		WithProducer(testExternalProducer{
			produceFunc: func(context.Context) ([]metricdata.ScopeMetrics, error) {
				// Splitting modifies the batch, so we need to create a new one each time.
				return []metricdata.ScopeMetrics{metricdata.ScopeMetrics{
					Scope: instrumentation.Scope{Name: "sdk/metric/test/reader/internal"},
					Metrics: []metricdata.Metrics{{
						Name:        "metric1",
						Description: "first of multiple metrics",
						Unit:        "ms",
						Data: metricdata.Gauge[int64]{
							DataPoints: []metricdata.DataPoint[int64]{{
								Attributes: attribute.NewSet(attribute.String("user", "david")),
								StartTime:  ts1,
								Time:       ts1.Add(time.Second),
								Value:      1,
							}},
						},
					}, {
						Name:        "metric2",
						Description: "second of multiple metrics",
						Unit:        "ms",
						Data: metricdata.Gauge[int64]{
							DataPoints: []metricdata.DataPoint[int64]{
								{
									Attributes: attribute.NewSet(attribute.String("user", "tyler")),
									StartTime:  ts2,
									Time:       ts2.Add(time.Second),
									Value:      10,
								},
								{
									Attributes: attribute.NewSet(attribute.String("user", "robert")),
									StartTime:  ts3,
									Time:       ts3.Add(time.Second),
									Value:      100,
								},
							},
						},
					}},
				}}, nil
			},
		}))
	r.register(testSDKProducer{})
	trigger <- time.Now()
	assert.Equal(t, <-eh.Err, errors.Join(errors.Join(errors.Join(assert.AnError), assert.AnError), assert.AnError))
	trigger <- time.Now()
	assert.Equal(t, <-eh.Err, errors.Join(errors.Join(errors.Join(assert.AnError), assert.AnError), assert.AnError))

	// Ensure Reader is allowed clean up attempt.
	_ = r.Shutdown(t.Context())
}

func TestPeriodicReaderFlushesPending(t *testing.T) {
	// Override the ticker so tests are not flaky and rely on timing.
	trigger := triggerTicker(t)
	t.Cleanup(func() { close(trigger) })

	expFunc := func(t *testing.T) (exp Exporter, called *bool) {
		called = new(bool)
		return &fnExporter{
			exportFunc: func(_ context.Context, m *metricdata.ResourceMetrics) error {
				// The testSDKProducer produces testResourceMetricsA.
				assert.Equal(t, testResourceMetricsAB, *m)
				*called = true
				return assert.AnError
			},
		}, called
	}

	t.Run("ForceFlush", func(t *testing.T) {
		exp, called := expFunc(t)
		r := NewPeriodicReader(exp, WithProducer(testExternalProducer{}))
		r.register(testSDKProducer{})
		assert.ErrorIs(t, r.ForceFlush(t.Context()), assert.AnError, "export error not returned")
		assert.True(t, *called, "exporter Export method not called, pending telemetry not flushed")

		// Ensure Reader is allowed clean up attempt.
		_ = r.Shutdown(t.Context())
	})

	t.Run("ForceFlush timeout on producer", func(t *testing.T) {
		exp, called := expFunc(t)
		timeout := time.Millisecond
		r := NewPeriodicReader(exp, WithTimeout(timeout), WithProducer(testExternalProducer{}))
		r.register(testSDKProducer{
			produceFunc: func(ctx context.Context, rm *metricdata.ResourceMetrics) error {
				select {
				case <-time.After(timeout + time.Second):
					*rm = testResourceMetricsA
				case <-ctx.Done():
					// we timed out before we could collect metrics
					return ctx.Err()
				}
				return nil
			},
		})
		assert.ErrorIs(t, r.ForceFlush(t.Context()), context.DeadlineExceeded)
		assert.False(t, *called, "exporter Export method called when it should have failed before export")

		// Ensure Reader is allowed clean up attempt.
		_ = r.Shutdown(t.Context())
	})

	t.Run("ForceFlush timeout on external producer", func(t *testing.T) {
		exp, called := expFunc(t)
		timeout := time.Millisecond
		r := NewPeriodicReader(exp, WithTimeout(timeout), WithProducer(testExternalProducer{
			produceFunc: func(ctx context.Context) ([]metricdata.ScopeMetrics, error) {
				select {
				case <-time.After(timeout + time.Second):
				case <-ctx.Done():
					// we timed out before we could collect metrics
					return nil, ctx.Err()
				}
				return []metricdata.ScopeMetrics{testScopeMetricsA}, nil
			},
		}))
		r.register(testSDKProducer{})
		assert.ErrorIs(t, r.ForceFlush(t.Context()), context.DeadlineExceeded)
		assert.False(t, *called, "exporter Export method called when it should have failed before export")

		// Ensure Reader is allowed clean up attempt.
		_ = r.Shutdown(t.Context())
	})

	t.Run("Shutdown", func(t *testing.T) {
		exp, called := expFunc(t)
		r := NewPeriodicReader(exp, WithProducer(testExternalProducer{}))
		r.register(testSDKProducer{})
		assert.ErrorIs(t, r.Shutdown(t.Context()), assert.AnError, "export error not returned")
		assert.True(t, *called, "exporter Export method not called, pending telemetry not flushed")
	})

	t.Run("Shutdown timeout on producer", func(t *testing.T) {
		exp, called := expFunc(t)
		timeout := time.Millisecond
		r := NewPeriodicReader(exp, WithTimeout(timeout), WithProducer(testExternalProducer{}))
		r.register(testSDKProducer{
			produceFunc: func(ctx context.Context, rm *metricdata.ResourceMetrics) error {
				select {
				case <-time.After(timeout + time.Second):
					*rm = testResourceMetricsA
				case <-ctx.Done():
					// we timed out before we could collect metrics
					return ctx.Err()
				}
				return nil
			},
		})
		assert.ErrorIs(t, r.Shutdown(t.Context()), context.DeadlineExceeded)
		assert.False(t, *called, "exporter Export method called when it should have failed before export")
	})

	t.Run("Shutdown timeout on external producer", func(t *testing.T) {
		exp, called := expFunc(t)
		timeout := time.Millisecond
		r := NewPeriodicReader(exp, WithTimeout(timeout), WithProducer(testExternalProducer{
			produceFunc: func(ctx context.Context) ([]metricdata.ScopeMetrics, error) {
				select {
				case <-time.After(timeout + time.Second):
				case <-ctx.Done():
					// we timed out before we could collect metrics
					return nil, ctx.Err()
				}
				return []metricdata.ScopeMetrics{testScopeMetricsA}, nil
			},
		}))
		r.register(testSDKProducer{})
		assert.ErrorIs(t, r.Shutdown(t.Context()), context.DeadlineExceeded)
		assert.False(t, *called, "exporter Export method called when it should have failed before export")
	})
}

func TestPeriodicReaderMultipleForceFlush(t *testing.T) {
	ctx := t.Context()
	r := NewPeriodicReader(new(fnExporter), WithProducer(testExternalProducer{}))
	r.register(testSDKProducer{})
	require.NoError(t, r.ForceFlush(ctx))
	require.NoError(t, r.ForceFlush(ctx))
	require.NoError(t, r.Shutdown(ctx))
}

func BenchmarkPeriodicReader(b *testing.B) {
	r := NewPeriodicReader(new(fnExporter))
	b.Run("Collect", benchReaderCollectFunc(r))
	require.NoError(b, r.Shutdown(b.Context()))
}

func TestPeriodicReaderTemporality(t *testing.T) {
	tests := []struct {
		name     string
		exporter *fnExporter
		// Currently only testing constant temporality. This should be expanded
		// if we put more advanced selection in the SDK
		wantTemporality metricdata.Temporality
	}{
		{
			name:            "default",
			exporter:        new(fnExporter),
			wantTemporality: metricdata.CumulativeTemporality,
		},
		{
			name:            "delta",
			exporter:        &fnExporter{temporalityFunc: deltaTemporalitySelector},
			wantTemporality: metricdata.DeltaTemporality,
		},
		{
			name:            "cumulative",
			exporter:        &fnExporter{temporalityFunc: cumulativeTemporalitySelector},
			wantTemporality: metricdata.CumulativeTemporality,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var undefinedInstrument InstrumentKind
			rdr := NewPeriodicReader(tt.exporter)
			assert.Equal(t, tt.wantTemporality.String(), rdr.temporality(undefinedInstrument).String())
		})
	}
}

func TestPeriodicReaderCollect(t *testing.T) {
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
			rdr := NewPeriodicReader(new(fnExporter))
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

func TestPeriodicReaderInstrumentation(t *testing.T) {
	// Enable SDK observability.
	t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	// Set up a global MeterProvider to collect the instrumentation metrics.
	// The PeriodicReader's instrumentation emits metrics to the global MeterProvider.
	orig := otel.GetMeterProvider()
	t.Cleanup(func() { otel.SetMeterProvider(orig) })

	instrumentationReader := NewManualReader()
	instrumentationMP := NewMeterProvider(WithReader(instrumentationReader))
	otel.SetMeterProvider(instrumentationMP)
	t.Cleanup(func() { _ = instrumentationMP.Shutdown(t.Context()) })

	// Create a periodic reader with an exporter that returns an error on export but not shutdown
	exp := &fnExporter{
		exportFunc: func(context.Context, *metricdata.ResourceMetrics) error {
			return assert.AnError
		},
		shutdownFunc: func(context.Context) error {
			return nil // No error on shutdown
		},
	}

	periodicReader := NewPeriodicReader(exp)
	t.Cleanup(func() { _ = periodicReader.Shutdown(t.Context()) })
	periodicReader.register(testSDKProducer{})

	// Exercise a collect (producer data).
	var got metricdata.ResourceMetrics
	require.NoError(t, periodicReader.Collect(t.Context(), &got))

	// Trigger a collection and export (which exercises the instrumentation)
	err := periodicReader.ForceFlush(t.Context())
	assert.Error(t, err, "expected error from exporter")

	// Collect the instrumentation metrics from the global MeterProvider
	var instrumentationMetrics metricdata.ResourceMetrics
	require.NoError(t, instrumentationReader.Collect(t.Context(), &instrumentationMetrics))

	targetName := otelconv.SDKMetricReaderCollectionDuration{}.Name()
	targetDesc := otelconv.SDKMetricReaderCollectionDuration{}.Description()
	targetUnit := otelconv.SDKMetricReaderCollectionDuration{}.Unit()

	// Find the SDK reader self-metric in the instrumentation metrics.
	foundMetric := findMetricByName(&instrumentationMetrics, targetName)
	require.NotNil(t, foundMetric, "SDK reader self-metric %q should be found in instrumentation metrics", targetName)

	// Basic identity checks (don't assert scope name/version; that can vary).
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

	const expectedComponentType = "periodic_metric_reader"

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

func TestPeriodicReaderInstrumentationError(t *testing.T) {
	// Enable SDK observability.
	t.Setenv("OTEL_GO_X_OBSERVABILITY", "true")

	// Set up a MeterProvider that returns errors when creating instruments.
	// This simulates the error path in NewPeriodicReader where observ.NewInstrumentation fails.
	orig := otel.GetMeterProvider()
	t.Cleanup(func() { otel.SetErrorHandler(otel.GetErrorHandler()) })
	t.Cleanup(func() { otel.SetMeterProvider(orig) })

	// Create an error handler to capture the error from otel.Handle()
	eh := newChErrorHandler()
	otel.SetErrorHandler(eh)

	// Set up a MeterProvider that returns errors
	mp := &errMeterProvider{err: assert.AnError}
	otel.SetMeterProvider(mp)

	// Create a periodic reader - this should trigger the error path
	exp := &fnExporter{}
	periodicReader := NewPeriodicReader(exp)
	t.Cleanup(func() { _ = periodicReader.Shutdown(t.Context()) })

	// Verify that the error was handled via otel.Handle()
	select {
	case err := <-eh.Err:
		assert.Error(t, err, "expected error to be handled")
		assert.ErrorIs(t, err, assert.AnError, "expected the error from NewInstrumentation")
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for error to be handled")
	}

	// Verify the reader is still functional despite the instrumentation error
	periodicReader.register(testSDKProducer{})
	var rm metricdata.ResourceMetrics
	assert.NoError(t, periodicReader.Collect(t.Context(), &rm), "reader should still work without instrumentation")
}

// errMeterProvider is a test helper that returns errors when creating instruments.
type errMeterProvider struct {
	metric.MeterProvider
	err error
}

func (m *errMeterProvider) Meter(string, ...metric.MeterOption) metric.Meter {
	return &errMeter{err: m.err}
}

type errMeter struct {
	metric.Meter
	err error
}

func (m *errMeter) Float64Histogram(string, ...metric.Float64HistogramOption) (metric.Float64Histogram, error) {
	return nil, m.err
}

// createMetricDataTestProducer creates a producer using patterns from metricdatatest.
func createMetricDataTestProducer() testSDKProducer {
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
			for i := range 20 { // 20 scopes with 4 metrics each = 80 total metrics
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

func BenchmarkPeriodicReaderInstrumentation(b *testing.B) {
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

		var exportCallCount int64
		exp := &fnExporter{
			exportFunc: func(_ context.Context, _ *metricdata.ResourceMetrics) error {
				// Count exports to ensure they're happening
				exportCallCount++
				return nil
			},
			shutdownFunc: func(context.Context) error {
				return nil
			},
		}

		r := NewPeriodicReader(exp)
		// Register with producer using metricdatatest patterns for realistic benchmark data
		r.register(createMetricDataTestProducer())
		b.Cleanup(func() {
			_ = r.Shutdown(b.Context()) // Ignore error in cleanup
		})

		rm := &metricdata.ResourceMetrics{}

		b.ReportAllocs()
		b.ResetTimer()

		for b.Loop() {
			// Test the full collect and export operation (simulating what collectAndExport does)
			err := r.Collect(b.Context(), rm)
			if err == nil {
				err = exp.Export(b.Context(), rm)
			}
			_ = err // Ignore error for benchmark
		}

		b.StopTimer()
		if exportCallCount == 0 {
			b.Fatalf("Expected exports to be called, but got 0")
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

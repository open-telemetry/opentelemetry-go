// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric // import "go.opentelemetry.io/otel/sdk/metric"

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
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
	ts.Equal(assert.AnError, ts.ErrReader.ForceFlush(context.Background()))
}

func (ts *periodicReaderTestSuite) TestShutdownPropagated() {
	ts.Equal(assert.AnError, ts.ErrReader.Shutdown(context.Background()))
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
	assert.Equal(t, assert.AnError, <-eh.Err)

	// Ensure Reader is allowed clean up attempt.
	_ = r.Shutdown(context.Background())
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
		assert.Equal(t, assert.AnError, r.ForceFlush(context.Background()), "export error not returned")
		assert.True(t, *called, "exporter Export method not called, pending telemetry not flushed")

		// Ensure Reader is allowed clean up attempt.
		_ = r.Shutdown(context.Background())
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
		assert.ErrorIs(t, r.ForceFlush(context.Background()), context.DeadlineExceeded)
		assert.False(t, *called, "exporter Export method called when it should have failed before export")

		// Ensure Reader is allowed clean up attempt.
		_ = r.Shutdown(context.Background())
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
		assert.ErrorIs(t, r.ForceFlush(context.Background()), context.DeadlineExceeded)
		assert.False(t, *called, "exporter Export method called when it should have failed before export")

		// Ensure Reader is allowed clean up attempt.
		_ = r.Shutdown(context.Background())
	})

	t.Run("Shutdown", func(t *testing.T) {
		exp, called := expFunc(t)
		r := NewPeriodicReader(exp, WithProducer(testExternalProducer{}))
		r.register(testSDKProducer{})
		assert.Equal(t, assert.AnError, r.Shutdown(context.Background()), "export error not returned")
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
		assert.ErrorIs(t, r.Shutdown(context.Background()), context.DeadlineExceeded)
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
		assert.ErrorIs(t, r.Shutdown(context.Background()), context.DeadlineExceeded)
		assert.False(t, *called, "exporter Export method called when it should have failed before export")
	})
}

func TestPeriodicReaderMultipleForceFlush(t *testing.T) {
	ctx := context.Background()
	r := NewPeriodicReader(new(fnExporter), WithProducer(testExternalProducer{}))
	r.register(testSDKProducer{})
	require.NoError(t, r.ForceFlush(ctx))
	require.NoError(t, r.ForceFlush(ctx))
	require.NoError(t, r.Shutdown(ctx))
}

func BenchmarkPeriodicReader(b *testing.B) {
	r := NewPeriodicReader(new(fnExporter))
	b.Run("Collect", benchReaderCollectFunc(r))
	require.NoError(b, r.Shutdown(context.Background()))
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
	expiredCtx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-1))
	defer cancel()

	tests := []struct {
		name        string
		ctx         context.Context
		expectedErr error
	}{
		{
			name:        "with a valid context",
			ctx:         context.Background(),
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
			_, err = meter.RegisterCallback(func(_ context.Context, o metric.Observer) error {
				return nil
			}, testM)
			assert.NoError(t, err)

			rm := &metricdata.ResourceMetrics{}
			assert.Equal(t, tt.expectedErr, rdr.Collect(tt.ctx, rm))
		})
	}
}

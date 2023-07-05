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

package metric // import "go.opentelemetry.io/otel/sdk/metric"

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregation"
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

func (e *fnExporter) Aggregation(k InstrumentKind) aggregation.Aggregation {
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

	ErrReader Reader
}

func (ts *periodicReaderTestSuite) SetupTest() {
	ts.readerTestSuite.SetupTest()

	e := &fnExporter{
		exportFunc:   func(context.Context, *metricdata.ResourceMetrics) error { return assert.AnError },
		flushFunc:    func(context.Context) error { return assert.AnError },
		shutdownFunc: func(context.Context) error { return assert.AnError },
	}

	ts.ErrReader = NewPeriodicReader(e)
	ts.ErrReader.register(testSDKProducer{})
	ts.ErrReader.RegisterProducer(testExternalProducer{})
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
			Factory: func() Reader {
				return NewPeriodicReader(new(fnExporter))
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

	r := NewPeriodicReader(exp)
	r.register(testSDKProducer{})
	r.RegisterProducer(testExternalProducer{})
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
		r := NewPeriodicReader(exp)
		r.register(testSDKProducer{})
		r.RegisterProducer(testExternalProducer{})
		assert.Equal(t, assert.AnError, r.ForceFlush(context.Background()), "export error not returned")
		assert.True(t, *called, "exporter Export method not called, pending telemetry not flushed")

		// Ensure Reader is allowed clean up attempt.
		_ = r.Shutdown(context.Background())
	})

	t.Run("Shutdown", func(t *testing.T) {
		exp, called := expFunc(t)
		r := NewPeriodicReader(exp)
		r.register(testSDKProducer{})
		r.RegisterProducer(testExternalProducer{})
		assert.Equal(t, assert.AnError, r.Shutdown(context.Background()), "export error not returned")
		assert.True(t, *called, "exporter Export method not called, pending telemetry not flushed")
	})
}

func BenchmarkPeriodicReader(b *testing.B) {
	b.Run("Collect", benchReaderCollectFunc(
		NewPeriodicReader(new(fnExporter)),
	))
}

func TestPeriodiclReaderTemporality(t *testing.T) {
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
		name string

		ctx             context.Context
		resourceMetrics *metricdata.ResourceMetrics

		expectedErr error
	}{
		{
			name: "with a valid context",

			ctx:             context.Background(),
			resourceMetrics: &metricdata.ResourceMetrics{},
		},
		{
			name: "with an expired context",

			ctx:             expiredCtx,
			resourceMetrics: &metricdata.ResourceMetrics{},

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

			assert.Equal(t, tt.expectedErr, rdr.Collect(tt.ctx, tt.resourceMetrics))
		})
	}
}

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package metric // import "go.opentelemetry.io/otel/sdk/metric"

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/go-logr/logr/testr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
)

type readerTestSuite struct {
	suite.Suite

	Factory func(...ReaderOption) Reader
	Reader  Reader
}

func (ts *readerTestSuite) SetupSuite() {
	otel.SetLogger(testr.New(ts.T()))
}

func (ts *readerTestSuite) TearDownTest() {
	// Ensure Reader is allowed attempt to clean up.
	_ = ts.Reader.Shutdown(context.Background())
}

func (ts *readerTestSuite) TestErrorForNotRegistered() {
	ts.Reader = ts.Factory()
	err := ts.Reader.Collect(context.Background(), &metricdata.ResourceMetrics{})
	ts.ErrorIs(err, ErrReaderNotRegistered)
}

func (ts *readerTestSuite) TestSDKProducer() {
	ts.Reader = ts.Factory()
	ts.Reader.register(testSDKProducer{})
	m := metricdata.ResourceMetrics{}
	err := ts.Reader.Collect(context.Background(), &m)
	ts.NoError(err)
	ts.Equal(testResourceMetricsA, m)
}

func (ts *readerTestSuite) TestExternalProducer() {
	ts.Reader = ts.Factory(WithProducer(testExternalProducer{}))
	ts.Reader.register(testSDKProducer{})
	m := metricdata.ResourceMetrics{}
	err := ts.Reader.Collect(context.Background(), &m)
	ts.NoError(err)
	ts.Equal(testResourceMetricsAB, m)
}

func (ts *readerTestSuite) TestCollectAfterShutdown() {
	ts.Reader = ts.Factory(WithProducer(testExternalProducer{}))
	ctx := context.Background()
	ts.Reader.register(testSDKProducer{})
	ts.Require().NoError(ts.Reader.Shutdown(ctx))

	m := metricdata.ResourceMetrics{}
	err := ts.Reader.Collect(ctx, &m)
	ts.ErrorIs(err, ErrReaderShutdown)
	ts.Equal(metricdata.ResourceMetrics{}, m)
}

func (ts *readerTestSuite) TestShutdownTwice() {
	ts.Reader = ts.Factory(WithProducer(testExternalProducer{}))
	ctx := context.Background()
	ts.Reader.register(testSDKProducer{})
	ts.Require().NoError(ts.Reader.Shutdown(ctx))
	ts.ErrorIs(ts.Reader.Shutdown(ctx), ErrReaderShutdown)
}

func (ts *readerTestSuite) TestMultipleRegister() {
	ts.Reader = ts.Factory()
	p0 := testSDKProducer{
		produceFunc: func(ctx context.Context, rm *metricdata.ResourceMetrics) error {
			// Differentiate this producer from the second by returning an
			// error.
			*rm = testResourceMetricsA
			return assert.AnError
		},
	}
	p1 := testSDKProducer{}

	ts.Reader.register(p0)
	// This should be ignored.
	ts.Reader.register(p1)

	err := ts.Reader.Collect(context.Background(), &metricdata.ResourceMetrics{})
	ts.Equal(assert.AnError, err)
}

func (ts *readerTestSuite) TestExternalProducerPartialSuccess() {
	ts.Reader = ts.Factory(
		WithProducer(testExternalProducer{
			produceFunc: func(ctx context.Context) ([]metricdata.ScopeMetrics, error) {
				return []metricdata.ScopeMetrics{}, assert.AnError
			},
		}),
		WithProducer(testExternalProducer{
			produceFunc: func(ctx context.Context) ([]metricdata.ScopeMetrics, error) {
				return []metricdata.ScopeMetrics{testScopeMetricsB}, nil
			},
		}),
	)
	ts.Reader.register(testSDKProducer{})

	m := metricdata.ResourceMetrics{}
	err := ts.Reader.Collect(context.Background(), &m)
	ts.ErrorIs(err, assert.AnError)
	ts.Equal(testResourceMetricsAB, m)
}

func (ts *readerTestSuite) TestSDKFailureBlocksExternalProducer() {
	ts.Reader = ts.Factory(WithProducer(testExternalProducer{}))
	ts.Reader.register(testSDKProducer{
		produceFunc: func(ctx context.Context, rm *metricdata.ResourceMetrics) error {
			*rm = metricdata.ResourceMetrics{}
			return assert.AnError
		},
	})

	m := metricdata.ResourceMetrics{}
	err := ts.Reader.Collect(context.Background(), &m)
	ts.Equal(assert.AnError, err)
	ts.Equal(metricdata.ResourceMetrics{}, m)
}

func (ts *readerTestSuite) TestMethodConcurrentSafe() {
	ts.Reader = ts.Factory(WithProducer(testExternalProducer{}))
	// Requires the race-detector (a default test option for the project).

	// All reader methods should be concurrent-safe.
	ts.Reader.register(testSDKProducer{})
	ctx := context.Background()

	var wg sync.WaitGroup
	const threads = 2
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = ts.Reader.temporality(InstrumentKindCounter)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = ts.Reader.aggregation(InstrumentKindCounter)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = ts.Reader.Collect(ctx, &metricdata.ResourceMetrics{})
		}()

		if f, ok := ts.Reader.(interface{ ForceFlush(context.Context) error }); ok {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_ = f.ForceFlush(ctx)
			}()
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = ts.Reader.Shutdown(ctx)
		}()
	}
	wg.Wait()
}

func (ts *readerTestSuite) TestShutdownBeforeRegister() {
	ts.Reader = ts.Factory(WithProducer(testExternalProducer{}))
	ctx := context.Background()
	ts.Require().NoError(ts.Reader.Shutdown(ctx))
	// Registering after shutdown should not revert the shutdown.
	ts.Reader.register(testSDKProducer{})

	m := metricdata.ResourceMetrics{}
	err := ts.Reader.Collect(ctx, &m)
	ts.ErrorIs(err, ErrReaderShutdown)
	ts.Equal(metricdata.ResourceMetrics{}, m)
}

func (ts *readerTestSuite) TestCollectNilResourceMetricError() {
	ts.Reader = ts.Factory()
	ctx := context.Background()
	ts.Error(ts.Reader.Collect(ctx, nil))
}

var testScopeMetricsA = metricdata.ScopeMetrics{
	Scope: instrumentation.Scope{Name: "sdk/metric/test/reader"},
	Metrics: []metricdata.Metrics{{
		Name:        "fake data",
		Description: "Data used to test a reader",
		Unit:        "1",
		Data: metricdata.Sum[int64]{
			Temporality: metricdata.CumulativeTemporality,
			IsMonotonic: true,
			DataPoints: []metricdata.DataPoint[int64]{{
				Attributes: attribute.NewSet(attribute.String("user", "alice")),
				StartTime:  time.Now(),
				Time:       time.Now().Add(time.Second),
				Value:      -1,
			}},
		},
	}},
}

var testScopeMetricsB = metricdata.ScopeMetrics{
	Scope: instrumentation.Scope{Name: "sdk/metric/test/reader/external"},
	Metrics: []metricdata.Metrics{{
		Name:        "fake scope data",
		Description: "Data used to test a Producer reader",
		Unit:        "ms",
		Data: metricdata.Gauge[int64]{
			DataPoints: []metricdata.DataPoint[int64]{{
				Attributes: attribute.NewSet(attribute.String("user", "ben")),
				StartTime:  time.Now(),
				Time:       time.Now().Add(time.Second),
				Value:      10,
			}},
		},
	}},
}

var testResourceMetricsA = metricdata.ResourceMetrics{
	Resource:     resource.NewSchemaless(attribute.String("test", "Reader")),
	ScopeMetrics: []metricdata.ScopeMetrics{testScopeMetricsA},
}

var testResourceMetricsAB = metricdata.ResourceMetrics{
	Resource:     resource.NewSchemaless(attribute.String("test", "Reader")),
	ScopeMetrics: []metricdata.ScopeMetrics{testScopeMetricsA, testScopeMetricsB},
}

type testSDKProducer struct {
	produceFunc func(context.Context, *metricdata.ResourceMetrics) error
}

func (p testSDKProducer) produce(ctx context.Context, rm *metricdata.ResourceMetrics) error {
	if p.produceFunc != nil {
		return p.produceFunc(ctx, rm)
	}
	*rm = testResourceMetricsA
	return nil
}

type testExternalProducer struct {
	produceFunc func(context.Context) ([]metricdata.ScopeMetrics, error)
}

func (p testExternalProducer) Produce(ctx context.Context) ([]metricdata.ScopeMetrics, error) {
	if p.produceFunc != nil {
		return p.produceFunc(ctx)
	}
	return []metricdata.ScopeMetrics{testScopeMetricsB}, nil
}

func benchReaderCollectFunc(r Reader) func(*testing.B) {
	ctx := context.Background()
	r.register(testSDKProducer{})

	// Store benchmark results in a closure to prevent the compiler from
	// inlining and skipping the function.
	var (
		collectedMetrics metricdata.ResourceMetrics
		err              error
	)

	return func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			err = r.Collect(ctx, &collectedMetrics)
			assert.Equalf(b, testResourceMetricsA, collectedMetrics, "unexpected Collect response: (%#v, %v)", collectedMetrics, err)
		}
	}
}

func TestDefaultAggregationSelector(t *testing.T) {
	var undefinedInstrument InstrumentKind
	assert.Panics(t, func() { DefaultAggregationSelector(undefinedInstrument) })

	iKinds := []InstrumentKind{
		InstrumentKindCounter,
		InstrumentKindUpDownCounter,
		InstrumentKindHistogram,
		InstrumentKindGauge,
		InstrumentKindObservableCounter,
		InstrumentKindObservableUpDownCounter,
		InstrumentKindObservableGauge,
	}

	for _, ik := range iKinds {
		assert.NoError(t, DefaultAggregationSelector(ik).err(), ik)
	}
}

func TestDefaultTemporalitySelector(t *testing.T) {
	var undefinedInstrument InstrumentKind
	for _, ik := range []InstrumentKind{
		undefinedInstrument,
		InstrumentKindCounter,
		InstrumentKindUpDownCounter,
		InstrumentKindHistogram,
		InstrumentKindGauge,
		InstrumentKindObservableCounter,
		InstrumentKindObservableUpDownCounter,
		InstrumentKindObservableGauge,
	} {
		assert.Equal(t, metricdata.CumulativeTemporality, DefaultTemporalitySelector(ik))
	}
}

type notComparable [0]func() // nolint:unused  // non-comparable type itself is used.

type noCompareReader struct {
	notComparable // nolint:unused  // non-comparable type itself is used.
	Reader
}

func TestReadersNotRequiredToBeComparable(t *testing.T) {
	r := noCompareReader{Reader: NewManualReader()}
	assert.NotPanics(t, func() { _ = NewMeterProvider(WithReader(r)) })
}

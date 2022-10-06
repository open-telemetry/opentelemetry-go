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

package metric // import "go.opentelemetry.io/otel/sdk/metric/reader"

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
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/metric/view"
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

func (ts *readerTestSuite) SetupTest() {
	ts.Reader = ts.Factory()
}

func (ts *readerTestSuite) TearDownTest() {
	// Ensure Reader is allowed attempt to clean up.
	_ = ts.Reader.Shutdown(context.Background())
}

func (ts *readerTestSuite) TestErrorForNotRegistered() {
	_, err := ts.Reader.Collect(context.Background())
	ts.ErrorIs(err, ErrReaderNotRegistered)
}

func (ts *readerTestSuite) TestProducer() {
	ts.Reader.register(testProducer{})
	m, err := ts.Reader.Collect(context.Background())
	ts.NoError(err)
	ts.Equal(testMetrics, m)
}

func (ts *readerTestSuite) TestCollectAfterShutdown() {
	ctx := context.Background()
	ts.Reader.register(testProducer{})
	ts.Require().NoError(ts.Reader.Shutdown(ctx))

	m, err := ts.Reader.Collect(ctx)
	ts.ErrorIs(err, ErrReaderShutdown)
	ts.Equal(metricdata.ResourceMetrics{}, m)
}

func (ts *readerTestSuite) TestShutdownTwice() {
	ctx := context.Background()
	ts.Reader.register(testProducer{})
	ts.Require().NoError(ts.Reader.Shutdown(ctx))
	ts.ErrorIs(ts.Reader.Shutdown(ctx), ErrReaderShutdown)
}

func (ts *readerTestSuite) TestMultipleForceFlush() {
	ctx := context.Background()
	ts.Reader.register(testProducer{})
	ts.Require().NoError(ts.Reader.ForceFlush(ctx))
	ts.NoError(ts.Reader.ForceFlush(ctx))
}

func (ts *readerTestSuite) TestMultipleRegister() {
	p0 := testProducer{
		produceFunc: func(ctx context.Context) (metricdata.ResourceMetrics, error) {
			// Differentiate this producer from the second by returning an
			// error.
			return testMetrics, assert.AnError
		},
	}
	p1 := testProducer{}

	ts.Reader.register(p0)
	// This should be ignored.
	ts.Reader.register(p1)

	_, err := ts.Reader.Collect(context.Background())
	ts.Equal(assert.AnError, err)
}

func (ts *readerTestSuite) TestMethodConcurrency() {
	// Requires the race-detector (a default test option for the project).

	// All reader methods should be concurrent-safe.
	ts.Reader.register(testProducer{})
	ctx := context.Background()

	var wg sync.WaitGroup
	const threads = 2
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = ts.Reader.Collect(ctx)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = ts.Reader.ForceFlush(ctx)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = ts.Reader.Shutdown(ctx)
		}()
	}
	wg.Wait()
}

func (ts *readerTestSuite) TestShutdownBeforeRegister() {
	ctx := context.Background()
	ts.Require().NoError(ts.Reader.Shutdown(ctx))
	// Registering after shutdown should not revert the shutdown.
	ts.Reader.register(testProducer{})

	m, err := ts.Reader.Collect(ctx)
	ts.ErrorIs(err, ErrReaderShutdown)
	ts.Equal(metricdata.ResourceMetrics{}, m)
}

func (ts *readerTestSuite) TestReaderUsesBridge() {
	reader := ts.Factory(
		WithBridge(testBridge{}),
	)
	reader.register(testProducer{})

	m, err := reader.Collect(context.Background())
	ts.NoError(err)
	ts.Equal(m, metricdata.ResourceMetrics{
		Resource: resource.NewSchemaless(attribute.String("test", "Reader")),
		ScopeMetrics: []metricdata.ScopeMetrics{
			testScopeMetrics1,
			testScopeMetrics2,
		},
	})
}

func (ts *readerTestSuite) TestReaderBridgeErrors() {
	reader := ts.Factory(
		WithBridge(testBridge{collectFunc: func(ctx context.Context) (metricdata.ScopeMetrics, error) {
			return testScopeMetrics2, assert.AnError
		}}),
	)
	reader.register(testProducer{})

	m, err := reader.Collect(context.Background())
	ts.Error(err)
	ts.Equal(m, metricdata.ResourceMetrics{
		Resource: resource.NewSchemaless(attribute.String("test", "Reader")),
		ScopeMetrics: []metricdata.ScopeMetrics{
			testScopeMetrics1,
			testScopeMetrics2,
		},
	})
}

var testScopeMetrics1 = metricdata.ScopeMetrics{
	Scope: instrumentation.Scope{Name: "sdk/metric/test/reader"},
	Metrics: []metricdata.Metrics{{
		Name:        "fake data",
		Description: "Data used to test a reader",
		Unit:        unit.Dimensionless,
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
var testScopeMetrics2 = metricdata.ScopeMetrics{
	Scope: instrumentation.Scope{Name: "sdk/metric/test/bridge"},
	Metrics: []metricdata.Metrics{{
		Name:        "more fake data",
		Description: "Data used to test a bridge",
		Unit:        unit.Dimensionless,
		Data: metricdata.Gauge[int64]{
			DataPoints: []metricdata.DataPoint[int64]{{
				Attributes: attribute.NewSet(attribute.String("user", "bob")),
				StartTime:  time.Now(),
				Time:       time.Now().Add(time.Second),
				Value:      -1,
			}},
		},
	}},
}

var testMetrics = metricdata.ResourceMetrics{
	Resource: resource.NewSchemaless(attribute.String("test", "Reader")),
	ScopeMetrics: []metricdata.ScopeMetrics{
		testScopeMetrics1,
	},
}

type testProducer struct {
	produceFunc func(context.Context) (metricdata.ResourceMetrics, error)
}

func (p testProducer) produce(ctx context.Context) (metricdata.ResourceMetrics, error) {
	if p.produceFunc != nil {
		return p.produceFunc(ctx)
	}
	return testMetrics, nil
}

type testBridge struct {
	collectFunc func(context.Context) (metricdata.ScopeMetrics, error)
}

func (t testBridge) Collect(ctx context.Context) (metricdata.ScopeMetrics, error) {
	if t.collectFunc != nil {
		return t.collectFunc(ctx)
	}
	return testScopeMetrics2, nil
}

func benchReaderCollectFunc(r Reader) func(*testing.B) {
	ctx := context.Background()
	r.register(testProducer{})

	// Store bechmark results in a closure to prevent the compiler from
	// inlining and skipping the function.
	var (
		collectedMetrics metricdata.ResourceMetrics
		err              error
	)

	return func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for n := 0; n < b.N; n++ {
			collectedMetrics, err = r.Collect(ctx)
			assert.Equalf(b, testMetrics, collectedMetrics, "unexpected Collect response: (%#v, %v)", collectedMetrics, err)
		}
	}
}

func TestDefaultAggregationSelector(t *testing.T) {
	var undefinedInstrument view.InstrumentKind
	assert.Panics(t, func() { DefaultAggregationSelector(undefinedInstrument) })

	iKinds := []view.InstrumentKind{
		view.SyncCounter,
		view.SyncUpDownCounter,
		view.SyncHistogram,
		view.AsyncCounter,
		view.AsyncUpDownCounter,
		view.AsyncGauge,
	}

	for _, ik := range iKinds {
		assert.NoError(t, DefaultAggregationSelector(ik).Err(), ik)
	}
}

func TestDefaultTemporalitySelector(t *testing.T) {
	var undefinedInstrument view.InstrumentKind
	for _, ik := range []view.InstrumentKind{
		undefinedInstrument,
		view.SyncCounter,
		view.SyncUpDownCounter,
		view.SyncHistogram,
		view.AsyncCounter,
		view.AsyncUpDownCounter,
		view.AsyncGauge,
	} {
		assert.Equal(t, metricdata.CumulativeTemporality, DefaultTemporalitySelector(ik))
	}
}
